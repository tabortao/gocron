#!/bin/bash

# Vue 3 组件批量迁移脚本

echo "开始迁移 Vue 组件..."

# 查找所有 .vue 文件
find ./src -name "*.vue" -type f | while read file; do
    echo "处理: $file"
    
    # 1. 替换 Element UI 导入
    sed -i '' 's/from '\''element-ui'\''/from '\''element-plus'\''/g' "$file"
    sed -i '' 's/import {Message}/import {ElMessage as Message}/g' "$file"
    sed -i '' 's/import { Message }/import { ElMessage as Message }/g' "$file"
    
    # 2. 替换 Message 调用
    sed -i '' 's/Message\./ElMessage\./g' "$file"
    sed -i '' 's/this\.\$message/this\.\$message/g' "$file"
    
    # 3. 替换过滤器语法 (简单情况)
    # 注意：复杂的过滤器需要手动处理
    sed -i '' 's/| formatTime/\$filters\.formatTime(/g' "$file"
    sed -i '' 's/}}/)}}/g' "$file"
    
done

echo "批量替换完成！"
echo "请注意："
echo "1. 过滤器的替换可能不完整，需要手动检查"
echo "2. 某些 Element Plus API 变更需要手动调整"
echo "3. 建议逐个测试组件功能"
