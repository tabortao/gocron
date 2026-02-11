# 任务计划：安装检测与 pnpm workspace 迁移（2026-02-11）

## 目标

- 升级/迁移场景下，若检测到已有 SQLite 数据库文件，自动跳过安装引导并指向正确 DB。
- 将前端验证脚本 verify.sh 全面切换为 pnpm。
- 将仓库迁移为 pnpm workspace，统一锁文件与常用脚本，并同步更新相关构建文件（Dockerfile/Makefile）。

## 执行步骤

1. 增强后端安装检测逻辑：当 install.lock 缺失但已存在 DB 文件时，自动设置环境变量并视为已安装。
2. 将 web/vue/verify.sh 从 yarn 改为 pnpm 命令。
3. 添加 pnpm-workspace.yaml 并调整 package.json name，避免 workspace 包名冲突。
4. 删除 web/vue/pnpm-lock.yaml 与 web/vue/yarn.lock，统一由根 pnpm-lock.yaml 管理。
5. 更新 Dockerfile.gocron 与 makefile 的前端构建命令为 pnpm/workspace 方式。

## 验证

- Go：运行 go test ./... 确保后端变更无编译错误。
- Node：在仓库根执行 pnpm install，确认生成/更新根 pnpm-lock.yaml，并能 pnpm -C web/vue build。
- Docker：确认 Dockerfile.gocron 前端阶段不再依赖 web/vue/pnpm-lock.yaml。
