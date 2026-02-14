# Bark 通知教程

## 1. Bark 是什么

Bark 是 iOS 上常用的推送通知工具。你可以在手机上安装 Bark App，然后通过一个包含设备 Key 的 URL 接收推送。

gocron 的 Bark 通知发送方式为：

- 使用你配置的 Bark 地址作为基础 URL（通常是 `https://api.day.app/你设备的key`）
- 在发送时拼接 `/{title}/{body}` 并以 GET 请求触发推送

## 2. 获取 Bark Key

1. 在 iPhone 上安装 Bark（App Store）
2. 打开 Bark App，在首页即可看到你的设备 Key

## 3. gocron 中如何配置

进入：系统管理 → 通知设置 → Bark

### 3.1 配置模板

- 标题模板：用于生成推送标题
- 内容模板：用于生成推送正文

模板支持 Go template 变量（与其他通知方式一致），常用变量如下：

- `{{.TaskId}}`：任务 ID
- `{{.TaskName}}`：任务名称
- `{{.Status}}`：状态（Success/Failed）
- `{{.StatusZh}}`：状态中文（成功/失败）
- `{{.Result}}`：完整输出
- `{{.ResultSummary}}`：输出摘要（优先 JSON.message）
- `{{.Host}}`：节点信息（若输出包含 Host 行）
- `{{.Remark}}`：任务备注

### 3.2 新增 Bark 地址

点击“新增 Bark 地址”，填写：

- 名称：自定义即可
- URL：必须为“地址 + 设备 Key”的形式，例如：

```
https://api.day.app/你设备的key
```

说明：

- 不需要手动在 URL 中加 `/{title}/{body}`，gocron 发送时会自动拼接
- 如果你使用的是自建 Bark 服务，请将 URL 换成你的服务地址 + Key（保持同样格式）

## 4. 任务中如何选择 Bark 地址

在任务编辑页开启通知后：

- 通知类型勾选 “Bark”
- Bark 地址可多选，也可以选择“全部 Bark 地址”

## 5. 常见问题

### 5.1 收不到推送

- 检查 Bark 地址是否可访问（网络/代理）
- 确认 Key 是否正确
- 确认 Bark App 已允许通知权限

### 5.2 标题/内容出现乱码或路径异常

gocron 发送时会对 title/body 做 URL 路径编码；若你使用自建服务且对路径编码有特殊处理，请确认服务端兼容标准 URL 编码。
