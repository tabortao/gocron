#!/bin/bash

# 发布脚本
# 使用方法: ./scripts/release.sh v1.2.3

set -e

VERSION=$1

if [ -z "$VERSION" ]; then
    echo "Usage: $0 <version>"
    echo "Example: $0 v1.2.3"
    exit 1
fi

# 检查是否有未提交的更改
if [ -n "$(git status --porcelain)" ]; then
    echo "Error: There are uncommitted changes. Please commit or stash them first."
    exit 1
fi

# 检查是否在main分支
CURRENT_BRANCH=$(git branch --show-current)
if [ "$CURRENT_BRANCH" != "main" ]; then
    echo "Warning: You are not on the main branch. Current branch: $CURRENT_BRANCH"
    read -p "Do you want to continue? (y/N): " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        exit 1
    fi
fi

# 检查版本格式
if [[ ! $VERSION =~ ^v[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
    echo "Error: Version must be in format vX.Y.Z (e.g., v1.2.3)"
    exit 1
fi

# 检查标签是否已存在
if git rev-parse "$VERSION" >/dev/null 2>&1; then
    echo "Error: Tag $VERSION already exists"
    exit 1
fi

echo "Preparing release $VERSION..."

# 运行测试
echo "Running tests..."
make test

# 构建前端
echo "Building frontend..."
make build-vue

# 生成静态资源
echo "Generating static assets..."
go install github.com/rakyll/statik@latest
go generate ./...

# 本地构建测试
echo "Testing local build..."
make build

echo "Creating git tag..."
git tag -a "$VERSION" -m "Release $VERSION"

echo "Release $VERSION is ready!"
echo ""
echo "To complete the release:"
echo "1. Push the tag: git push origin $VERSION"
echo "2. GitHub Actions will automatically build and create the release"
echo ""
echo "Or to cancel:"
echo "1. Delete the tag: git tag -d $VERSION"