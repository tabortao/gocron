# NAS/Docker 部署偶发异常排障指引

## 常见现象

- Web 能打开，但任务不执行、主机/worker 操作失败、部分页面接口报错。
- 运行一段时间后才出现，重启 master（gocron）后恢复。

这类“部分运行”通常不是单一 Bug，而是下面几类问题叠加导致。

## 一步定位：先看健康检查

升级到包含 `/api/healthz` 的版本后，优先请求：

- `GET /api/healthz`

关注字段：

- `installed`：是否已安装完成
- `db.ok / db.error`：数据库连通性
- `scheduler.running / scheduler.entryCount`：调度器是否启动、已加载任务数
- `scheduler.concurrencyUsed / scheduler.concurrencyCap`：并发队列占用
- `rpc.poolSize`：gRPC 连接池大小

如果 HTTP 状态码为 503，说明 master 处于不健康状态（但接口仍会返回详细 data）。

## 排查优先级

### 1) NAS 资源与宿主机问题（最高优先级）

- 观察容器是否出现 OOM/被系统杀死：
  - 容器重启次数、Docker 事件、NAS 系统日志（OOM killer）
- 观察磁盘：
  - 空间不足、卷只读、I/O 长期打满都会让“写入相关功能”失效（任务日志、配置写入、SQLite 写入）
- 观察时间同步：
  - NAS/NTP 异常会影响定时触发行为（看似“任务不跑”）

### 2) 数据库问题（SQLite 尤其要注意挂载方式）

#### MySQL / PostgreSQL

- 网络抖动、连接空闲过久、NAT/防火墙回收连接，可能出现“部分接口卡住/报错”。
- 建议把 DB 放在稳定网络内，并确保连接池与超时配置合理。

#### SQLite

- 不建议把 SQLite 数据文件放在 NFS/SMB 等网络文件系统上运行。
  - 这类文件系统的锁语义/一致性实现差异大，容易出现间歇性锁住、写失败、超时。
- 更推荐：
  - SQLite 放在本机持久化卷（本地磁盘）；
  - 或直接切换到 MySQL/PostgreSQL。

### 3) gRPC / worker 连接问题（非常符合“部分运行 + 重启恢复”）

典型场景：

- worker 容器重启（IP 变化）或网络短暂不可达；
- master 侧连接池复用旧连接，导致后续 RPC 持续失败；
- 重启 master 后连接池重建，问题消失。

改进建议：

- 使用包含“RPC Unavailable 自动释放并重建连接”的版本：
  - master 在遇到连接不可用时会释放该地址连接，下一次调用重新 Dial，从而自愈。
- 业务侧辅助判断：
  - Web UI 正常但 `/api/host/ping/:id` 失败，优先怀疑 gRPC/worker 通道或 worker 状态。

## 建议的 Docker 运行策略

- 为 gocron 配置 healthcheck（调用 `/api/healthz`），并结合重启策略：
  - 在 master 不健康时自动重启，减少人工干预。
- 为 worker 配置合理的重启策略与资源限制，避免反复重启造成雪崩。

## 收集信息（复现时最有价值）

- master 日志中：
  - `database ping failed`（数据库不可用）
  - RPC 相关报错（Unavailable/timeout）
  - 调度器初始化日志（任务是否成功加载）
- NAS 系统日志：
  - OOM、磁盘只读、I/O 错误
- 容器事件：
  - 是否发生重启、退出码、重启时间点与“部分运行”出现时间点是否一致
