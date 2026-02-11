# Git 提交与 GitHub Tag Release 指南

本文介绍如何在本项目中提交代码、打 tag，并通过 GitHub Actions 自动发布 Release（上传构建产物）。

## 1. 前置条件

- 已在 GitHub 创建仓库：`https://github.com/tabortao/gocron`
- 本地仓库已设置 `origin` 指向该地址
- 你有权限 push 分支与 tags

验证远程地址：

```bash
git remote -v
```

如需设置：

```bash
git remote add origin https://github.com/tabortao/gocron.git
```

## 2. 提交代码并推送到 GitHub

1. 查看变更：

```bash
git status
```

2. 暂存变更：

```bash
git add -A
```

3. 提交：

```bash
git commit -m "feat: xxx"
```

4. 推送分支（以 main 为例）：

```bash
git push origin main
```

## 3. 通过 tag 触发自动 Release

本项目已包含 GitHub Actions 工作流：当 push 满足 `v*.*.*` 格式的 tag 时，会自动构建并上传 Release 产物。

工作流文件：

- `.github/workflows/release-packages.yml`

### 3.1 创建并推送 tag（推荐：annotated tag）

示例（版本号按语义化版本递增，例如 `v1.4.10`）：

```bash
git tag -a v1.4.10 -m "Release v1.4.10"
git push origin v1.4.10
```

推送 tag 后：

- GitHub Actions 会自动构建（含前端资源）
- 自动创建/更新 GitHub Release
- 自动上传打包文件（zip/tar.gz）

### 3.2 运行自动发布脚本（Windows / Linux）

脚本会做以下事情：

- 检查工作区是否干净（避免把未提交的改动打进 release）
- 可选执行基础检查（go test、前端 build）
- 创建 annotated tag 并 push 到 origin

优先推荐：使用 Python 一键完成（commit/push + tag/push）：

```bash
python ./scripts/release.py --version v1.5.4 --message "chore: release v1.5.4"
```

Windows PowerShell：

```powershell
pwsh .\scripts\release.ps1 -Version v1.4.10
```

如果你的 Windows 没有安装 PowerShell 7（没有 `pwsh` 命令），可以用 Windows PowerShell 5.1：

```powershell
powershell -ExecutionPolicy Bypass -File .\scripts\release.ps1 -Version v1.4.10
```

Linux/macOS（bash）：

```bash
bash ./scripts/release.sh -Version v1.4.10
```

## 4. 常见问题

### 4.1 push 了 tag 但没看到 Release

建议按顺序排查：

- tag 是否符合 `v*.*.*`（例如 `v1.4.10`）
- Actions 页面中 `Release Packages` 工作流是否执行、是否失败
- 是否有 `contents: write` 权限（fork 场景下权限可能受限）

### 4.2 如何重新发布同一个版本

建议做法是发布一个新版本（例如 `v1.4.10` → `v1.4.11`）。如果必须复用同名 tag，需要先删除远程 tag 再重建：

```bash
git tag -d v1.4.10
git push origin :refs/tags/v1.4.10
git tag -a v1.4.10 -m "Release v1.4.10"
git push origin v1.4.10
```
