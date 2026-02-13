# NAS 用户使用教程（gocron）

本文面向 NAS（群晖 / 威联通 / TrueNAS / Linux NAS）用户，介绍如何在 NAS 上安装并接入 gocron 任务节点（gocron-node），以及如何在 Web 控制台创建任务节点与定时任务。

## 0. 基本概念

- **gocron（主控/控制台）**：提供 Web UI（默认 `:5920`），保存任务、触发调度、向任务节点下发执行指令。
- **gocron-node（任务节点/Agent）**：运行在目标机器（NAS）上，提供 gRPC 服务（默认 `:5921`），接收主控下发的命令并在本机执行。

## 1. 前置条件

- NAS 可以访问到 gocron 主控的地址（例如 `http://192.168.3.4:5920`）。
- NAS 上能执行 shell 命令（SSH/终端），且有 sudo 权限（安装脚本会用 sudo 写入 `/opt` 与 systemd 服务）。
- NAS 防火墙/安全组允许主控访问 NAS 的 **5921/tcp**（如果你使用默认端口）。

## 2. 创建任务节点（推荐：自动注册）

1. 登录 gocron Web（例如 `http://192.168.3.4:5920`）。
2. 进入 **任务节点** 页面，点击 **自动注册**。
3. 在弹窗里复制 Linux/macOS 的安装命令，在 NAS 上执行。

示例（两种方式二选一，token 以实际为准）：

```bash
# 方式一：普通权限安装（推荐）
curl -fsSL 'http://192.168.3.4:5920/api/agent/install.sh?token=xxxx' | bash

# 方式二：allow-root 安装（谨慎）
curl -fsSL 'http://192.168.3.4:5920/api/agent/install.sh?token=xxxx' | sudo bash
```

执行成功后：

- NAS 上会安装 `gocron-node` 到 `/opt/gocron-node`
- Linux NAS 会创建并启动 systemd 服务 `gocron-node`
- 主控会自动注册该节点到 **任务节点** 列表

说明：

- 推荐优先使用“普通权限安装”，并通过 `sudo usermod -aG docker <user>` 的方式授权节点用户访问 Docker。
- 仅当你明确需要以 root 运行节点时才使用 allow-root 安装（安装脚本会自动为 `gocron-node` 追加 `-allow-root` 参数）。

## 3. 创建任务节点（手动方式）

如果你不使用自动注册，也可以手动添加：

1. 进入 **任务节点 → 新增节点**
2. 按如下规则填写：
   - **别名**：给你自己看的名字，例如 `NAS-1`、`群晖主机`。
   - **主机名**：填写 NAS 的可达地址（从 gocron 主控所在机器能访问到的 IP/域名），例如 `192.168.3.20`。
   - **端口**：默认 `5921`（除非你启动 gocron-node 时改了）。
3. 保存后，在列表里点 **测试发送**，看到“连接成功”表示主控能连到节点。

## 4. 创建定时任务（示例）

进入 **任务管理 → 新增任务**，常用字段说明：

- **Cron 表达式**：支持标准语法（秒 分 时 天 月 周），也支持快捷语法（界面里有示例按钮）。
- **任务节点**：选择你刚创建/自动注册的 NAS 节点。
- **命令**：在节点机器上执行的命令（Linux NAS 通常是 bash/sh 命令）。
- **超时时间**：秒；建议为耗时任务设置合理超时，避免卡死。
- **通知**：可选（失败通知/总是通知/关键字匹配通知）。

### 示例 A：每天凌晨 02:30 备份某目录到 NAS 备份目录

目标：把 `/volume1/docker` 打包到 `/volume1/backup`，文件名带时间戳。

1. **Cron 表达式**：`0 30 2 * * *`（每天 02:30:00）
2. **任务节点**：选择 NAS 节点
3. **命令**：

```bash
set -e
TS=$(date +"%Y%m%d_%H%M%S")
mkdir -p /volume1/backup
tar -czf "/volume1/backup/docker_${TS}.tar.gz" -C /volume1 docker
echo "backup ok: ${TS}"
```

4. **超时时间**：例如 `3600`（1 小时）

### 示例 B：每 5 分钟检查磁盘剩余空间，低于阈值输出关键字触发通知

1. **Cron 表达式**：`0 */5 * * * *`（每 5 分钟）
2. **命令**：

```bash
set -e
THRESHOLD=15
AVAIL=$(df -P / | awk 'NR==2{print 100-$5}' | tr -d '%')
echo "disk_avail_percent=${AVAIL}"
if [ "$AVAIL" -lt "$THRESHOLD" ]; then
  echo "ALERT: disk is low"
  exit 1
fi
```

3. **通知策略**：
   - 失败通知（或关键字匹配通知，关键字设为 `ALERT`）

## 5. 常见问题

### 5.1 测试发送失败

优先检查：

- NAS 上 `gocron-node` 是否运行、端口是否监听：`ss -lntp | grep 5921`（不同 NAS 命令可能不同）
- gocron 主控到 NAS 的网络是否可达、是否被防火墙拦截（尤其是 5921/tcp）
- 主机名是否填写为“从主控可达”的地址（不要填 NAS 自己的 `127.0.0.1`）

### 5.2 systemd 不可用的 NAS

部分 NAS 没有标准 systemd（或被裁剪）。此时建议：

- 手动启动：`/opt/gocron-node/gocron-node -s 0.0.0.0:5921`
- 使用 NAS 自带的“开机启动/任务计划”功能配置自启动
