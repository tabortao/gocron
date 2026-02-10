# Gocron 开发环境运行指南

本文档介绍如何在本地开发环境中配置、运行和调试 Gocron 项目。

## 环境要求

在开始之前，请确保您的开发环境满足以下要求：

- **Go**: 1.24.0 或更高版本
- **Node.js**: 推荐 v20+
- **包管理器**: pnpm (前端使用)
- **C 编译器**: 无需安装 (已迁移至纯 Go SQLite 驱动)

## 快速开始 (Windows/Linux/macOS)

我们提供了一个简单的脚本来自动安装依赖并启动服务。

### Windows (PowerShell)

在项目根目录下运行：

```powershell
./run_dev.ps1
```

### 手动运行步骤

如果您更喜欢手动控制每一步，请参考以下说明。

#### 1. 后端 (Go)

后端服务负责 API 接口和任务调度。

1.  **安装 Go 依赖**:

    ```bash
    go mod download
    ```

2.  **设置环境变量 (可选)**:
    默认情况下，系统使用 SQLite 数据库。如果需要自定义，可以设置环境变量：

    ```bash
    # Windows PowerShell
    $env:GOCRON_DB_ENGINE="sqlite"
    $env:GOCRON_DB_DATABASE="gocron.db"

    # Linux/macOS
    export GOCRON_DB_ENGINE=sqlite
    export GOCRON_DB_DATABASE=gocron.db
    ```

3.  **启动服务**:
    ```bash
    # CGO 已不再必须，但默认开启也无妨
    # 如果想强制关闭 CGO: $env:CGO_ENABLED="0"
    go run cmd/gocron/gocron.go web -e dev
    ```
    服务默认运行在 `0.0.0.0:5920`。

#### 2. 前端 (Vue.js)

前端管理界面基于 Vue 3 + Vite。

1.  **进入前端目录**:

    ```bash
    cd web/vue
    ```

2.  **安装依赖 (使用 pnpm)**:

    ```bash
    pnpm install
    ```

3.  **启动开发服务器**:
    ```bash
    pnpm run dev
    ```
    访问终端输出的地址（通常是 `http://localhost:8080`）。

## 数据库配置

### SQLite (默认)

项目默认配置为使用 SQLite。无需额外安装数据库服务，数据文件将存储在运行目录下的 `gocron.db` (或配置的路径)。

- 使用纯 Go 驱动，无需安装 GCC。

### MySQL (兼容模式)

如果您需要连接现有的 MySQL 数据库：

1.  **环境变量方式**:

    ```bash
    $env:GOCRON_DB_ENGINE="mysql"
    $env:GOCRON_DB_HOST="127.0.0.1"
    $env:GOCRON_DB_PORT="3306"
    $env:GOCRON_DB_USER="root"
    $env:GOCRON_DB_PASSWORD="password"
    $env:GOCRON_DB_DATABASE="gocron"
    ```

2.  **配置文件方式**:
    在 `conf/app.ini` 中修改：
    ```ini
    [default]
    db.engine=mysql
    db.host=127.0.0.1
    db.port=3306
    db.user=root
    db.password=password
    db.database=gocron
    ```

## 常见问题

**Q: 启动后端报错 `exec: "gcc": executable file not found in %PATH%`**
A: 请确保您已更新到最新代码。新版已使用纯 Go SQLite 驱动，不再依赖 GCC。如果您仍需使用旧版驱动，请安装 GCC。

**Q: 前端报错 `pnpm: command not found`**
A: 请运行 `npm install -g pnpm` 安装 pnpm。
