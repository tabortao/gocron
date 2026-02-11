# 2026-02-11 GitHub Actions Release 打包计划

## 目标

在发布（tag/release）时，通过 GitHub Actions 自动构建并上传以下产物到 GitHub Release：

- gocron-<版本>-darwin-amd64.tar.gz
- gocron-<版本>-darwin-arm64.tar.gz
- gocron-<版本>-linux-amd64.tar.gz
- gocron-<版本>-linux-arm64.tar.gz
- gocron-<版本>-windows-amd64.zip
- gocron-node-darwin-amd64.tar.gz
- gocron-node-darwin-arm64.tar.gz
- gocron-node-linux-amd64.tar.gz
- gocron-node-linux-arm64.tar.gz
- gocron-node-windows-amd64.zip

其中 `<版本>` 从 tag 提取，若 tag 为 `v1.4.9` 则 `<版本>` 为 `1.4.9`。

## 约束与实现选择

- 使用现有 `package.sh` 进行跨平台编译与打包，保证产物结构与项目保持一致。
- 前端依赖使用 pnpm，在 Action 内执行 `web/vue` 的构建，确保 `web/vue/dist` 可用于 Go embed。
- darwin 产物在 `macos-latest` 上构建（避免 Linux 交叉编译 macOS 的 cgo/工具链问题）。
- linux + windows 产物在 `ubuntu-latest` 上构建（Go 纯交叉编译即可产出 zip/tar.gz）。
- 最终由单独的 release job 汇总各 job 的 artifact，并上传到 GitHub Release。

## 执行步骤

1. 新增 `.github/workflows/release-packages.yml`：
   - 触发：`push` tag `v*.*.*` + `workflow_dispatch`
   - job1：ubuntu 构建 linux + windows 包
   - job2：macos 构建 darwin 包
   - job3：下载产物并上传到 GitHub Release
2. 在 workflow 中统一提取版本号：
   - `TAG=${{ github.ref_name }}` 或手动输入
   - `VERSION=${TAG#v}`
3. 验证产物命名是否符合目标列表（通过脚本/列目录）。
