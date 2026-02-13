# allow-root 安装脚本支持

## 目标

- 允许用户在执行 `curl .../api/agent/install.sh?... | bash` 时使用 root 安装，并让 gocron-node 以 root 运行时自动追加 `-allow-root` 参数。
- 保持默认场景下仍可用普通用户安装与运行。

## 方案

- 安装脚本内检测 `id -u`：
  - root：不使用 sudo，设置 `NODE_ARGS="-allow-root"` 并在 systemd/macOS 启动命令中追加该参数。
  - 非 root：要求存在 sudo，用 sudo 执行安装与 systemd 管理操作。

## 变更点

- 后端：`/api/agent/install.sh` 脚本生成逻辑调整（root 时自动启用 `-allow-root`，并兼容无 sudo 场景）。
- 文档：补充 “允许 root 运行” 的说明与注意事项。
- 前端：安装提示文案增加 `-allow-root` 提示。

## 验证

- 非 root 安装：执行安装命令后，systemd 服务以普通用户运行，节点可正常注册。
- root 安装：使用 `sudo bash -c "curl ... | bash"` 安装后，systemd ExecStart 包含 `-allow-root` 且节点可正常运行。
- Docker 控制：节点用户加入 docker 组并重启后，任务可执行 `docker ps` / `docker compose`。
