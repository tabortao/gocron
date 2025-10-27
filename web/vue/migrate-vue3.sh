#!/bin/bash

# Vue 3 完整迁移脚本

echo "=========================================="
echo "Vue 2 to Vue 3 自动迁移脚本"
echo "=========================================="

# 备份
echo "创建备份..."
cp -r src src.backup.$(date +%Y%m%d_%H%M%S)

# 查找所有 .vue 和 .js 文件
echo "开始处理文件..."

find ./src -type f \( -name "*.vue" -o -name "*.js" \) | while read file; do
    echo "处理: $file"
    
    # 1. Element UI -> Element Plus
    sed -i '' "s/from 'element-ui'/from 'element-plus'/g" "$file"
    sed -i '' 's/from "element-ui"/from "element-plus"/g' "$file"
    sed -i '' "s/'element-ui\/lib\/theme-chalk\/index.css'/'element-plus\/dist\/index.css'/g" "$file"
    sed -i '' 's/"element-ui\/lib\/theme-chalk\/index.css"/"element-plus\/dist\/index.css"/g' "$file"
    
    # 2. Message 组件
    sed -i '' 's/{Message}/{ElMessage}/g' "$file"
    sed -i '' 's/{ Message }/{ ElMessage }/g' "$file"
    sed -i '' 's/Message\./ElMessage\./g' "$file"
    sed -i '' 's/Message(/ElMessage(/g' "$file"
    
    # 3. MessageBox
    sed -i '' 's/{MessageBox}/{ElMessageBox}/g' "$file"
    sed -i '' 's/MessageBox\./ElMessageBox\./g' "$file"
    
    # 4. Notification
    sed -i '' 's/{Notification}/{ElNotification}/g' "$file"
    sed -i '' 's/Notification\./ElNotification\./g' "$file"
    
    # 5. Loading
    sed -i '' 's/{Loading}/{ElLoading}/g' "$file"
    sed -i '' 's/Loading\./ElLoading\./g' "$file"
    
done

# 处理 .vue 文件的特殊语法
find ./src -name "*.vue" -type f | while read file; do
    # 6. slot 语法: slot="title" -> #title
    sed -i '' 's/slot="title"/#title/g' "$file"
    sed -i '' 's/<template slot="title">/<template #title>/g' "$file"
    
    # 7. v-for 中移除 this
    sed -i '' 's/v-if="this\.\$/v-if="\$/g' "$file"
    sed -i '' 's/v-show="this\.\$/v-show="\$/g' "$file"
    sed -i '' 's/{{this\.\$/{{$/g' "$file"
    
    # 8. el-submenu -> el-sub-menu (Element Plus 命名变更)
    sed -i '' 's/<el-submenu/<el-sub-menu/g' "$file"
    sed -i '' 's/<\/el-submenu>/<\/el-sub-menu>/g' "$file"
    
done

echo ""
echo "=========================================="
echo "自动迁移完成！"
echo "=========================================="
echo ""
echo "⚠️  请注意以下事项："
echo ""
echo "1. 过滤器 (filters) 需要手动迁移"
echo "   之前: {{ value | filterName }}"
echo "   之后: {{ \$filters.filterName(value) }}"
echo ""
echo "2. \$listeners 已被移除，合并到 \$attrs"
echo ""
echo "3. v-model 自定义组件需要检查"
echo ""
echo "4. .sync 修饰符已移除，使用 v-model:propName"
echo ""
echo "5. 某些 Element Plus 组件 API 有变化，需要查阅文档"
echo ""
echo "6. 建议逐个测试组件功能"
echo ""
echo "备份已保存到 src.backup.* 目录"
echo ""
