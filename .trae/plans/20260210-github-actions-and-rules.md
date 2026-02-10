# GitHub Actions CI/CD 配置与经验总结

## 任务目标

1.  **经验总结**：将最近遇到的问题（lint-staged 在 Windows 下的路径问题、SQLite CGO 依赖问题、embed 目录缺失问题）总结到 `.trae/rules` 目录下的规则文件中，以便后续避免。
2.  **GitHub Actions 配置**：参考用户提供的示例，为本项目创建 GitHub Actions workflow，实现 Docker 镜像的自动构建和推送。

## 任务分解

### 1. 经验总结 (Rules 更新)

- [ ] 创建或更新 `.trae/rules/project_rules.xml` (假设这是规则文件的存储位置，如果没有则创建 `.trae/rules` 目录)。
- [ ] 记录以下规则：
  - **Windows 兼容性**: `lint-staged` 命令避免使用 `cd` 组合命令，应使用工具自带的 `--dir` 或路径参数。
  - **CGO 依赖**: 优先使用纯 Go 实现的库（如 `glebarez/sqlite`）以去除对 GCC 的依赖，提升跨平台兼容性。
  - **Go Embed**: 使用 `go:embed` 时，必须确保目标目录存在且不为空（可放 `.gitkeep` 或占位文件）。

### 2. GitHub Actions Workflow

- [ ] 创建目录 `.github/workflows`。
- [ ] 创建文件 `.github/workflows/docker-build.yml`。
- [ ] 内容适配：
  - **触发条件**: `release/**` 分支推送, `v*.*.*` 标签推送, `workflow_dispatch`。
  - **环境**: `ubuntu-latest`。
  - **步骤**:
    - Checkout。
    - Setup QEMU & Docker Buildx。
    - 提取版本号 (适配 Shell 脚本)。
    - 登录 Docker Hub & GHCR。
    - Docker meta (生成 tags 和 labels)。
    - Build and Push:
      - Context: `.`
      - Dockerfile: `Dockerfile.gocron` (注意本项目 Dockerfile 文件名)。
      - Platforms: `linux/amd64` (根据 env 配置)。
      - Tags/Labels: 引用 meta 输出。
      - Cache: GHA cache。

## 执行步骤

1.  写入规则文件。
2.  创建 Workflow 文件。
3.  验证 Workflow 文件的语法（人工检查）。
