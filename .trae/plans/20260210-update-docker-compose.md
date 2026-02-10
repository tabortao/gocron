# Docker Compose 镜像源与持久化配置更新

## 任务目标

1.  修改 `docker-compose.yml`，将镜像源改为 Docker Hub 的 `tabortao/gocron:latest`。
2.  确保数据库文件持久化（当前配置已使用 volume，需确认是否符合用户期望，或者改为绑定挂载）。
3.  检查并修改相关文件（如 README 或文档），确保文档与新的 docker-compose 配置一致。

## 任务分解

### 1. 修改 `docker-compose.yml`

- [ ] 将 `image: gocron:latest` 改为 `image: tabortao/gocron:latest`。
- [ ] 移除 `build` 块（既然使用 Docker Hub 镜像，通常不再需要本地构建配置，或者将其注释掉以供参考）。
- [ ] 确认 `volumes` 配置。当前是 `gocron-data:/.gocron`，这是一个命名卷。用户要求“数据库文件需要持久化”，命名卷可以满足，但用户可能更倾向于本地目录映射（如 `./data:/.gocron`）以便查看。
  - _决策_：保留命名卷通常更符合 Docker 最佳实践，但为了方便用户管理数据（尤其是 SQLite 文件），改为本地目录挂载 `./data:/app/data` 或类似结构可能更直观。
  - _当前路径_：应用配置目录是 `/.gocron`。
  - _改进_：将 `volumes` 改为 `- ./data:/.gocron`，这样数据会直接保存在当前目录的 `data` 文件夹下，方便用户备份和查看 `gocron.db`。

### 2. 检查相关文件

- [ ] 检查 `README.md` 或 `README_ZH.md` 中关于 Docker 启动的说明，更新镜像名称。
- [ ] 检查 `docs/development_guide.md`，如果有 Docker 相关内容，一并更新。

## 执行步骤

1.  修改 `docker-compose.yml`。
2.  搜索并更新文档中的 Docker 镜像引用。
