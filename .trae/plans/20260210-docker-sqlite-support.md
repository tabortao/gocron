# Docker 支持 SQLite 改进计划

## 任务目标

分析 `gocron` 项目在 Docker 环境下是否支持 SQLite，如果不支持找出原因，并进行改进以支持 SQLite。

## 现状分析

- [x] 检查 `Dockerfile` 构建配置，确认是否开启 CGO 及安装相关依赖（SQLite 需要 CGO 支持）。
  - 结论：`Dockerfile.gocron` 已开启 `CGO_ENABLED=1` 并安装了 `sqlite-dev`，支持 SQLite。
- [x] 检查 `docker-compose.yml` 配置，查看数据卷挂载和环境变量。
  - 结论：默认只挂载了数据卷，未配置数据库类型，默认使用 MySQL 导致启动失败。
- [x] 检查代码中数据库连接初始化部分，确认 SQLite 驱动是否被正确引入和使用。
  - 结论：代码支持 SQLite，但配置读取逻辑仅支持 ini 文件，不支持环境变量。

## 问题定位

- [x] 确认是否因为 `CGO_ENABLED=0` 导致 SQLite 驱动不可用（Go 的 pure Go SQLite 驱动较少用，通常用 `go-sqlite3` 需要 CGO）。
  - 结论：否，CGO 已开启。
- [x] 确认 Docker 容器中 SQLite 数据库文件的存储路径和持久化问题。
  - 结论：主要问题是缺乏通过环境变量配置数据库的能力，导致无法在 Docker 中方便地切换到 SQLite。

## 改进方案

- [x] 修改 `Dockerfile`：
  - 确保 `CGO_ENABLED=1`。(已存在，无需修改)
- [x] 修改 `internal/modules/setting/setting.go`：
  - 增加从环境变量读取配置的功能 (`GOCRON_DB_ENGINE`, `GOCRON_DB_DATABASE` 等)。
  - 处理配置文件不存在的情况，允许仅通过环境变量启动。
- [x] 修改 `internal/modules/app/app.go`：
  - 允许通过环境变量判定应用已安装，跳过 Web 安装向导。
- [x] 修改 `docker-compose.yml`：
  - 默认配置环境变量使用 SQLite。
  - 挂载路径设置为 `/.gocron/gocron.db`。

## 执行与验证

- [x] 执行文件修改。
- [x] 验证构建脚本（`go vet` 和 `go test`）。
  - `setting_test.go` 新增测试用例通过。
