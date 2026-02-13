# Docker 容器控制与权限排查（gocron + gocron-node）

## 适用场景

- 使用 Docker 部署了 gocron Web（Master，默认 5920）后，希望在任务中执行 `docker ps`、`docker compose up/down` 等命令来控制宿主机容器。
- 在 NAS/Linux 机器上安装并运行 gocron-node（Agent/Worker，默认 5921），让任务在节点机器上执行。

## 执行原理（必读）

- Docker 部署的通常只是 gocron Web（Master）。它负责管理任务、调度、下发执行请求。
- 真正执行 Shell 命令的是 gocron-node（运行在目标机器上）。你在任务日志里看到 `Host: [x.x.x.x:5921]` 就表示在该节点执行。
- 因此：
  - 需要在哪台机器控制 Docker，就必须在那台机器安装并运行 gocron-node。
  - gocron Web 容器内挂载了 `/var/run/docker.sock` 并不等于节点具备 Docker 权限；节点执行时仍以节点进程的用户权限为准。

## 常见报错与含义

- `docker-compose: command not found`
  - 节点环境没有 `docker-compose`（老命令），可能是 Docker Compose V2（新命令）只能用 `docker compose`。
- `permission denied while trying to connect to the Docker daemon socket at unix:///var/run/docker.sock`
  - 节点执行任务的用户无权访问 Docker 套接字 `/var/run/docker.sock`（通常属于 `root:docker`，权限 660）。

## 推荐配置（让普通用户也能执行 docker）

1. 在节点机器确认 Docker 已安装并正常运行：

```bash
docker --version
docker ps
```

2. 将运行 gocron-node 的用户加入 docker 组（示例用户：taozi）：

```bash
sudo usermod -aG docker taozi
```

3. 重启 Docker（不同系统可能命令不同，以下为常见写法）：

```bash
sudo service docker stop
sudo service docker start
```

4. 重启 gocron-node，使其重新加载用户组权限：

```bash
sudo systemctl restart gocron-node
```

5. 验证权限是否生效（以 gocron-node 的运行用户执行）：

```bash
docker ps
docker compose version || true
docker-compose version || true
```

## 任务配置案例：在 NAS 上控制 docker compose 项目

假设你的 compose 工程目录位于：

- `/vol1/1000/docker/memos`

创建一个 Shell 类型任务，脚本示例（兼容 compose v1/v2）：

```bash
#!/bin/bash
set -e
export PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin

cd /vol1/1000/docker/memos

if command -v docker-compose >/dev/null 2>&1; then
  docker-compose down
elif docker compose version >/dev/null 2>&1; then
  docker compose down
else
  echo "Docker Compose not found: need docker-compose or docker compose plugin" >&2
  exit 127
fi
```

### 为什么能控制宿主机容器

- 该脚本在节点机器上执行。
- `docker compose`/`docker-compose` 会通过 `/var/run/docker.sock` 与宿主 Docker daemon 通信。
- 只要节点运行用户具备访问该套接字的权限，就可进行容器管理操作。

## 安全建议

- 让进程拥有 `/var/run/docker.sock` 的访问权限，本质上等价于获得很高的宿主控制权限。
- 建议：
  - 仅让可信节点安装 gocron-node，并限制 Web 端账号权限。
  - 只在必要目录上进行挂载与操作，避免误删关键数据。
  - 不要在文档/截图中公开 Agent token、AuthSecret 等敏感信息，发现泄露需立即作废并重新生成。
