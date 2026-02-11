# Windows 用户使用教程（gocron）

本文面向 Windows 用户，介绍如何在 Windows 上接入 gocron 任务节点（gocron-node），以及如何在 Web 控制台创建任务节点与定时任务，并给出可直接复制的任务案例。

## 0. 基本概念

- **gocron（主控/控制台）**：提供 Web UI（默认 `:5920`），保存任务、触发调度、向任务节点下发执行指令。
- **gocron-node（任务节点/Agent）**：运行在目标机器（Windows）上，提供 gRPC 服务（默认 `:5921`），接收主控下发的命令并在本机执行。

## 1. 前置条件

- Windows 机器能访问到 gocron 主控地址（例如 `http://192.168.3.4:5920`）。
- Windows 防火墙允许主控访问 **5921/tcp**（默认端口）。
- 如果你在 Windows 本机开发，同时用 Docker 跑 gocron 主控：请注意容器内访问 `127.0.0.1` 指向的是容器自身，而不是宿主机。
- 常见误区：**gocron.exe（主控）只提供 Web/调度（5920），并不会监听 5921**。5921 必须由 **gocron-node** 进程提供；否则“测试发送”会出现 `actively refused`（端口未监听导致被拒绝）。

## 2. 在 Windows 上启动 gocron-node

### 2.1 方式 A：下载 Release 包（推荐）

1. 从 GitHub Release 下载 `gocron-node-windows-amd64.zip`（或 arm64 对应包）。
2. 解压到任意目录，例如 `C:\gocron-node\`
3. 打开 PowerShell，启动：

```powershell
cd C:\gocron-node
.\gocron-node.exe -s 0.0.0.0:5921
```

看到类似 “server listen on 0.0.0.0:5921” 即表示启动成功。

快速自检（建议在 Windows 节点机器上执行）：

```powershell
netstat -ano | findstr 5921
```

能看到 `LISTENING` 才表示节点端口已正确监听。

### 2.2 方式 B：源码运行（本地开发用）

在项目根目录执行（需要 Go 环境）：

```powershell
go run .\cmd\node\node.go -s 0.0.0.0:5921
```

## 3. 创建任务节点（如何填写：别名/主机名/端口）

进入 **任务节点 → 新增节点**：

- **别名**：随意填写，建议写清楚用途，例如 `Windows-DEV`、`Win10-PC`。
- **主机名**：填写“从 gocron 主控所在机器能访问到的 Windows 地址”：
  - 若 gocron 主控就在你这台 Windows 上直接运行（非 Docker）：可以填 `127.0.0.1`。
  - 若 gocron 主控运行在 Docker 容器内：不要填 `127.0.0.1`，请填 Windows 宿主机局域网 IP（例如 `192.168.3.10`），或在某些环境下可以用 `host.docker.internal`（取决于你的部署方式）。
  - 若 gocron 主控在另一台服务器：填这台 Windows 的局域网 IP（例如 `192.168.3.10`）。
- **端口**：默认 `5921`（与 gocron-node 启动参数一致）。

保存后，在节点列表点 **测试发送**：

- 成功：说明主控能连到 gocron-node，并能执行一个简单命令（默认测试命令为 `echo hello`）。
- 失败：优先排查网络连通与防火墙（见“常见问题”）。

## 4. 创建定时任务（示例）

进入 **任务管理 → 新增任务**，常用字段：

- **Cron 表达式**：支持标准语法（秒 分 时 天 月 周）。
- **任务节点**：选择 Windows 节点。
- **命令**：在 Windows 节点上执行的命令（会以 cmd 批处理方式运行）。
- **超时时间**：秒。

### 示例 A：每天 09:00 写入心跳日志

1. **Cron 表达式**：`0 0 9 * * *`
2. **任务节点**：Windows 节点
3. **命令**：

```bat
echo %date% %time% gocron heartbeat >> C:\gocron-node\heartbeat.log
```

4. **超时时间**：`30`

### 示例 B：每 10 分钟清理 7 天前的临时文件

1. **Cron 表达式**：`0 */10 * * * *`
2. **命令**（PowerShell 命令也可以通过 `powershell -NoProfile -Command` 调用）：

```bat
powershell -NoProfile -ExecutionPolicy Bypass -Command "Get-ChildItem -Path $env:TEMP -Recurse -File -ErrorAction SilentlyContinue | Where-Object { $_.LastWriteTime -lt (Get-Date).AddDays(-7) } | Remove-Item -Force -ErrorAction SilentlyContinue; Write-Output 'cleanup done'"
```

3. **超时时间**：`300`

### 示例 C：每分钟检查某端口是否存活（失败则返回非 0 触发失败通知）

1. **Cron 表达式**：`0 * * * * *`
2. **命令**：

```bat
powershell -NoProfile -ExecutionPolicy Bypass -Command "if (-not (Test-NetConnection 127.0.0.1 -Port 5920).TcpTestSucceeded) { Write-Output 'ALERT: port down'; exit 1 } else { Write-Output 'ok' }"
```

## 5. 常见问题

### 5.1 “测试发送”失败但我填了 127.0.0.1:5921

常见原因是 gocron 主控并不在同一个网络命名空间里：

- gocron 主控在 Docker 容器内运行：容器内的 `127.0.0.1` 不是你的 Windows。
- 解决：把节点“主机名”改成 Windows 的局域网 IP（例如 `192.168.3.10`），并确保 Windows 防火墙放行 5921。

### 5.2 Windows 防火墙放行 5921（示例）

可以在 Windows 防火墙高级设置里添加入站规则，或用 PowerShell（管理员）：

```powershell
New-NetFirewallRule -DisplayName "gocron-node 5921" -Direction Inbound -Action Allow -Protocol TCP -LocalPort 5921
```

### 5.3 节点启动了但仍连接不上

排查建议：

- 确认节点监听地址：启动参数建议 `-s 0.0.0.0:5921`，而不是只监听 `127.0.0.1`
- 从主控所在机器测试连通：`Test-NetConnection <windows-ip> -Port 5921`
