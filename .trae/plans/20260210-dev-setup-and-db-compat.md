# 开发环境搭建与双数据库支持计划

## 任务目标

1. 确保项目同时支持 SQLite 和 MySQL，且默认使用 SQLite，同时兼容旧版 MySQL 配置。
2. 编写开发环境运行教程并放入 `docs` 文件夹。
3. 自动安装依赖并在本地运行项目。

## 任务分解

### 1. 数据库兼容性与默认配置

- [x] 检查 `internal/modules/setting/setting.go`，确认默认值逻辑。
  - 将代码中的默认值 `MustString("mysql")` 修改为 `MustString("sqlite")`，实现了“默认使用 sqlite”。
  - 经确认，如果用户有旧的 MySQL 配置文件，只要文件中明确指定了 `db.engine=mysql`，则不受影响，兼容性良好。
- [x] 检查 `app.ini.sqlite.example` 和其他示例配置。(已存在，无需修改)

### 2. 文档编写

- [x] 创建 `docs` 目录。
- [x] 编写 `docs/development_guide.md`，包含：
  - 环境要求 (Go, Node.js, pnpm, GCC/MinGW for SQLite CGO)。
  - 后端运行步骤。
  - 前端运行步骤 (强调使用 pnpm)。
  - 数据库切换说明。

### 3. 自动安装与运行

- [x] 检查前端 `package.json` 脚本。
- [x] 执行后端依赖安装: `go mod download`。
- [x] 执行前端依赖安装: `cd web/vue && pnpm install`。
  - 检测到 `pnpm` 已安装。
- [x] 启动后端: `go run cmd/gocron/gocron.go web -e dev`。
  - 预先创建了 `web/vue/dist` 目录以满足 `embed` 编译要求。
  - 后端已成功启动在 `0.0.0.0:5920`。
- [x] 启动前端: `cd web/vue && pnpm run dev`。
  - 前端已成功启动在 `http://localhost:8080`。
- [x] (可选) 创建一个方便的启动脚本 `run_dev.ps1`。

## 执行结果

所有任务已完成。项目已在本地运行，文档已生成。
