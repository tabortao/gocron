# 飞书机器人 Webhook 通知教程（gocron）

本文说明如何把 gocron 的 WebHook 通知对接到飞书「自定义机器人」，以及“收不到消息”时如何排查。

## 1. 前置准备：创建飞书自定义机器人

1. 在飞书群里添加「自定义机器人」。
2. 复制机器人 Webhook 地址（形如 `https://open.feishu.cn/open-apis/bot/v2/hook/xxxx...`）。

注意：该 Webhook 地址等同于密钥，请不要公开传播（截图/日志/工单/群聊）。

## 2. gocron 系统配置：新增 Webhook 地址

进入：系统管理 → 通知设置 → WebHook 通知

1. 点击「新增Webhook地址」
2. 名称：建议写清楚用途，例如「飞书-告警群」
3. URL：粘贴飞书机器人 Webhook 地址
4. 保存

保存后页面下方会出现「Webhook地址列表」，并带有一个 `id`（任务里选择接收目标时实际用的是这个 id）。

## 3. gocron 模板配置：飞书可用的 JSON 模板（推荐）

gocron 的 Webhook 是：

- POST 请求
- Header 固定为 `Content-Type: application/json`
- Body 直接发送“模板渲染后的字符串”，因此模板最终必须是合法 JSON

模板使用 **Go template** 语法，变量名 **区分大小写**，且必须使用 **双大括号**。

模板支持的变量（Go template 语法，直接复制下面这种写法）：

- `{{.TaskId}}`
- `{{.TaskName}}`
- `{{.Status}}`
- `{{.Result}}`
- `{{.Remark}}`

另外，本项目对通知内容做了增强（用于把“节点返回的 JSON”提取成更干净的摘要），你还可以使用：

- `{{.StatusZh}}`：成功/失败
- `{{.IsSuccess}}`：true/false
- `{{.Host}}`：当输出包含 `Host: [...]` 时会提取
- `{{.ResultSummary}}`：优先提取 JSON 中的 `message`（例如 SUCCESS），避免刷屏
- `{{.ResultBody}}`：去掉 Host 行后的主体内容（可能仍是 JSON，会做长度截断）

### 3.0 常见错误（会导致字段为空/发不出来）

下面这些写法都不对（很容易出现你说的“状态/节点/摘要都为空”）：

- 把 `TaskId` 写成了 `Taskld`（字母 l 和大写 I 混淆）
- 写成了 `{.TaskId}` / `(.TaskId)` / `TaskId}}`（不是 Go template 标准写法）
- 少了点号：`{{TaskId}}`（必须是 `{{.TaskId}}`）

你可以用这个“最小可用模板”快速验证变量是否能正确渲染（如果这条都不对，说明模板没按 Go template 写）：

```json
{
  "msg_type": "text",
  "content": {
    "text": "ID={{.TaskId}} | Name={{.TaskName}} | Status={{.Status}}"
  }
}
```

### 3.0.1 换行规则（重要）

飞书文本里要换行，JSON 里写 **`\n`**（一个反斜杠），不要写 **`\\n`**，否则飞书会把它当作普通字符显示成 `\n`。

### 3.1 推荐：飞书交互卡片（更美观）

把系统里的「模板」设置为下面内容（直接复制粘贴）：

```json
{
  "msg_type": "interactive",
  "card": {
    "config": {
      "wide_screen_mode": true
    },
    "header": {
      "template": "{{ if .IsSuccess }}green{{ else }}red{{ end }}",
      "title": {
        "tag": "plain_text",
        "content": "gocron 任务通知"
      }
    },
    "elements": [
      {
        "tag": "div",
        "text": {
          "tag": "lark_md",
          "content": "**任务**：{{.TaskName}}（ID: {{.TaskId}}）\n**状态**：{{.StatusZh}}"
        }
      },
      {
        "tag": "div",
        "text": {
          "tag": "lark_md",
          "content": "{{ if .Host }}**节点**：{{.Host}}{{ end }}"
        }
      },
      {
        "tag": "div",
        "text": {
          "tag": "lark_md",
          "content": "**摘要**：{{.ResultSummary}}"
        }
      },
      {
        "tag": "hr"
      },
      {
        "tag": "div",
        "text": {
          "tag": "lark_md",
          "content": "{{ if .Remark }}**备注**：{{.Remark}}{{ end }}"
        }
      }
    ]
  }
}
```

### 3.2 简洁版：飞书 text 消息模板

把系统里的「模板」设置为下面内容（直接复制粘贴）：

```json
{
  "msg_type": "text",
  "content": {
    "text": "【gocron任务通知】\n任务：{{.TaskName}}（ID: {{.TaskId}}）\n状态：{{.StatusZh}}\n{{ if .Host }}节点：{{.Host}}\n{{ end }}摘要：{{.ResultSummary}}\n{{ if .Remark }}备注：{{.Remark}}{{ end }}"
  }
}
```

说明：

- 这是飞书机器人官方支持的基础消息格式：`msg_type=text`
- 如果模板不是飞书认可的结构，飞书会返回 4xx/错误码，导致你“看起来配置了但收不到消息”

### 3.3 可选：飞书富文本（post）模板

飞书 `post` 结构更复杂，适合做排版；如果你需要可以在上面的 `text` 方案稳定后再改造。

## 4. 任务配置：选择通知条件 + 选择 Webhook 接收目标

进入：任务管理 → 新增/编辑任务

1. 任务通知：选择通知条件
   - 不通知：永远不发
   - 失败通知：仅失败时发
   - 总是通知：每次执行都发
   - 关键字匹配通知：输出包含关键字才发
2. 通知类型：选择「WebHook」
3. 接收目标：勾选你在系统里新增的 Webhook 地址（例如「飞书-告警群」）
4. 保存任务

常见误区：

- 选择了「失败通知」，但你的任务一直成功：自然不会收到
- 选择了「关键字匹配通知」，但输出不包含关键字：自然不会收到

## 5. 收不到消息的排查清单

按下面顺序排查，通常能快速定位：

1. 确认任务通知条件是否会触发（总是通知/失败通知/关键字匹配）
2. 确认任务里已选择正确的 Webhook 接收目标（不是只在系统里保存了 URL）
3. 确认系统里的模板是“飞书认可的 JSON 结构”（推荐先用本文的 `text` 模板）
4. 查看 gocron 服务端日志
   - 关键字通常包含：`#webHook#`
   - 如果飞书返回非 2xx，会记录 `status` 和 `resp`（响应内容会被截断）
5. 用飞书机器人地址做一次手工 POST 验证（仅在本机/安全环境下执行，注意不要把 URL 发到群里）

PowerShell 示例（把 URL 替换成你自己的）：

```powershell
$url = "https://open.feishu.cn/open-apis/bot/v2/hook/你的hook"
$body = @{
  msg_type = "text"
  content  = @{ text = "test from powershell" }
} | ConvertTo-Json -Depth 5

Invoke-RestMethod -Method Post -Uri $url -ContentType "application/json" -Body $body
```

如果手工 POST 能收到，但 gocron 收不到：

- 重点检查任务通知条件/接收目标是否配置正确
- 再检查模板是否渲染成了合法 JSON（例如引号、换行、特殊字符）
- 如果你之前模板里直接用了 `{{.Result}}`，可能会把节点返回的 JSON 整段发到群里；建议改用 `{{.ResultSummary}}`
