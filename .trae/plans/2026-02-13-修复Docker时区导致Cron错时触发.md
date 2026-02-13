## 背景

- NAS 容器部署后，任务 cron `0 25 9 * * *` 在北京时间 17:25 触发，表现为按伦敦/UTC 时间触发（9:25）。
- docker-compose 已配置 `TZ=Asia/Shanghai`，但容器内可能缺少时区数据（zoneinfo/tzdata），导致 Go 无法加载 Asia/Shanghai，最终回退到 UTC。

## 目标

- 保证在 Docker 容器中只要设置了 `TZ=Asia/Shanghai`（或 `GOCRON_TIMEZONE`），调度器按对应时区计算触发时间。
- 即使基础镜像缺少 tzdata，也能正确加载常见 IANA 时区（如 Asia/Shanghai）。

## 方案

- [x] 在运行时镜像中安装 tzdata（提供 `/usr/share/zoneinfo`）。
- [x] 在 gocron 启动后读取时区配置并设置 `time.Local`（优先级：`GOCRON_TIMEZONE` > `TZ` > `app.ini`）。
- [x] 在应用初始化阶段引入 `time/tzdata`，确保没有 tzdata 也能加载 IANA 时区。
- [x] 增加最小化测试/自检：单测覆盖时区优先级，并在启动日志输出当前时区与时间。

## 验证

- [ ] Docker Compose 仅设置 `TZ=Asia/Shanghai` 启动服务后，日志输出 `time.Local=Asia/Shanghai`。
- [ ] 创建 cron `0 25 9 * * *`，在北京时间 09:25 执行，而不是 17:25。
