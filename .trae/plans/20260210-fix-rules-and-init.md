# 修复规则文件格式与初始化逻辑

## 任务目标

1.  将 `.trae\rules\project_rules.xml` 转换为 Markdown 格式 `.trae\rules\project_rules.md`，方便 Trae 识别。
2.  修复当 `./data/gocron.db` (或其他 SQLite 数据库文件) 存在时，项目启动仍进入初始化界面的问题。

## 任务分解

### 1. 规则文件转换

- [ ] 读取 `.trae/rules/project_rules.xml` 内容。
- [ ] 将其内容转换为 Markdown 格式 (Trae Rules 格式)。
- [ ] 写入 `.trae/rules/project_rules.md`。
- [ ] 删除 `.trae/rules/project_rules.xml`。

### 2. 修复初始化判定逻辑

- [ ] 分析 `internal/modules/app/app.go` 中的 `IsInstalled` 函数。
- [ ] 分析 `cmd/gocron/gocron.go` 的初始化流程。
- [ ] 修改 `IsInstalled` 逻辑：
  - 除了检查 `install.lock` 和环境变量外，增加对 SQLite 数据库文件的检测。
  - 如果 `app.ini` 存在且配置了 SQLite，且对应的 DB 文件存在，应视为已安装。
  - 或者，如果检测到默认位置 (如 `./gocron.db` 或 `./data/gocron.db`) 存在 SQLite 文件，且 `app.ini` 不存在或未配置，是否应尝试自动使用？
  - 更稳妥的方式：如果 `app.ini` 存在，且能成功连接数据库，视为已安装。但 `IsInstalled` 在 DB 初始化之前调用。
  - **策略**：检查 `app.ini` 是否存在。如果存在，读取其中的 DB 配置。如果是 SQLite 且文件存在，返回 `true`。同时，如果 `install.lock` 丢失但 DB 存在，自动补全 `install.lock` 可能是一个好主意，或者直接在内存中视为已安装。

## 执行步骤

1.  转换规则文件。
2.  读取代码 (`app.go`, `setting.go`)。
3.  修改 `app.go` 增强安装检测逻辑。
4.  验证修复。
