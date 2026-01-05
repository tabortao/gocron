#!/bin/bash

echo "🚀 启动Webhook测试服务..."

# 检查Go是否安装
if ! command -v go &> /dev/null; then
    echo "❌ 未找到Go，请先安装Go语言环境"
    exit 1
fi

# 进入webhook-test目录
cd "$(dirname "$0")"

# 启动服务
echo "📡 启动服务在端口8080..."
go run webhook-test-server.go