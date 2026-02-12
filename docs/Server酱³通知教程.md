# Server 酱³（SC3）通知教程（gocron）

本文说明如何把 gocron 的任务通知发送到 Server 酱³（SC3），以及常见问题排查。

## 1. 什么是 Server 酱³

Server 酱³提供一个可直接 POST 的 API URL，你可以用 `title` + `desp` 两个字段发送消息（`desp` 支持 Markdown）。

示例（PowerShell）：

```powershell
$url = "https://xxxx.push.ft07.com/send/xxxx.send"
$body = "title=❤️每早计划&desp=😄你好呀，记得创建好今天的计划呀✅"
Invoke-RestMethod -Method Post -Uri $url -ContentType "application/x-www-form-urlencoded" -Body $body
```

注意：API URL 等同于密钥，请不要公开传播（截图/日志/工单/群聊）。

## 2. gocron 系统配置：新增 Server 酱³

进入：系统管理 → 通知设置 → Server 酱³

1. 点击「新增API地址」
2. 名称：例如「Server酱-日程提醒」
3. URL：粘贴你的 SC3 API URL（例如 `https://xxxx.push.ft07.com/send/...send`）
4. 保存

随后在页面里配置：

- 标题模板（title）
- 内容模板（desp）

## 3. 模板语法与可用变量（gocron 特点）

gocron 模板使用 **Go template** 语法，必须使用 `{{.变量}}` 形式，变量名区分大小写。

常用变量：

- `{{.TaskId}}` / `{{.TaskName}}`
- `{{.Status}}`：Success / Failed
- `{{.StatusZh}}`：成功 / 失败
- `{{.IsSuccess}}`：true / false
- `{{.Host}}`：当输出包含 `Host: [...]` 时会提取
- `{{.ResultSummary}}`：输出摘要（优先提取 JSON.message，例如 SUCCESS）
- `{{.ResultBody}}`：去掉 Host 行后的主体输出
- `{{.Remark}}`

## 4. 推荐模板（适配 SC3：title + desp）

### 4.1 标题模板（title）

```txt
{{.TaskName}} - {{.StatusZh}}
```

### 4.2 内容模板（desp）

```txt
**任务**：{{.TaskName}}（ID: {{.TaskId}}）

**状态**：{{.StatusZh}}

{{ if .Host }}**节点**：{{.Host}}

{{ end }}**摘要**：{{.ResultSummary}}

{{ if .Remark }}**备注**：{{.Remark}}{{ end }}
```

说明：

- `desp` 支持 Markdown，所以这里用加粗和空行，阅读更清爽
- 建议优先使用 `{{.ResultSummary}}`，避免把节点返回的整段 JSON 刷屏

## 5. 任务配置：选择 Server 酱³ 通知

进入：任务管理 → 新增/编辑任务

1. 任务通知：选择通知条件（失败通知/总是通知/关键字匹配等）
2. 通知类型：选择「Server 酱³」
3. 接收用户：勾选你在系统里新增的 SC3 API 地址（可多选）
4. 保存任务

## 6. 收不到通知的排查清单

1. 任务通知条件是否会触发（失败通知但任务一直成功；关键字匹配但输出不包含关键字）
2. 任务里是否已选择 Server 酱³ 接收目标（不是只在系统里保存了 URL）
3. SC3 API URL 是否正确可用（建议用 curl/PowerShell 先手工发一条）
4. 查看 gocron 服务端日志（关键字通常包含 `#serverchan3#`）
