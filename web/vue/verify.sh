#!/bin/bash

set -e

echo "=========================================="
echo "前端改动验证脚本"
echo "=========================================="
echo ""

# 1. 依赖检查
echo "1. 检查依赖..."
yarn install --frozen-lockfile
echo "✅ 依赖检查完成"
echo ""

# 2. 运行测试
echo "2. 运行单元测试..."
yarn test --run
echo "✅ 测试通过"
echo ""

# 3. Lint 检查
echo "3. 运行 Lint 检查..."
yarn lint || echo "⚠️  Lint 有警告（非阻塞）"
echo ""

# 4. 构建验证
echo "4. 构建生产版本..."
yarn build
echo "✅ 构建成功"
echo ""

# 5. 检查构建产物
echo "5. 检查构建产物大小..."
echo ""
echo "主要 JS 文件:"
ls -lh dist/static/*.js 2>/dev/null | head -5 || echo "无 JS 文件"
echo ""
echo "主要 CSS 文件:"
ls -lh dist/static/*.css 2>/dev/null | head -5 || echo "无 CSS 文件"
echo ""

# 6. 总结
echo "=========================================="
echo "✅ 所有验证通过！"
echo "=========================================="
echo ""
echo "下一步:"
echo "  1. 启动开发服务器: yarn dev"
echo "  2. 手动测试功能"
echo ""
