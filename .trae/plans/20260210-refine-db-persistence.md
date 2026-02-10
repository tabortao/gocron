# Docker Compose 数据库路径优化

## 任务目标

1.  进一步优化 `docker-compose.yml` 中的数据库配置。
2.  将数据库文件路径修改为 `/.gocron/data/gocron.db`，使其在持久化目录中拥有独立的子目录，结构更清晰。
3.  更新相关文档。

## 任务分解

### 1. 修改 `docker-compose.yml`

- [ ] 将 `GOCRON_DB_DATABASE` 从 `/.gocron/gocron.db` 修改为 `/.gocron/data/gocron.db`。
- [ ] 保持 `volumes` 为 `- ./data:/.gocron` 不变。
- [ ] 结果：宿主机的 `./data` 目录下将自动生成 `data` 目录存放数据库文件，与 `conf` 和 `log` 并列。

### 2. 更新文档

- [ ] 检查 `docs/development_guide.md`，更新关于 Docker 运行部分（如果有）。
- [ ] 检查 `README.md`，更新快速开始部分。

## 执行步骤

1.  修改 `docker-compose.yml`。
2.  更新文档。
