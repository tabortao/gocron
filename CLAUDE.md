# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

**gocron** is a lightweight distributed scheduled task management system developed in Go, designed to replace Linux crontab. It features a master-worker architecture with web-based management interface.

### Key Components

- **gocron**: Main web server and scheduler (port 5920)
- **gocron-node**: Worker node that executes tasks (port 5921)
- **Web UI**: Vue.js frontend for task management
- **Database**: Supports MySQL, PostgreSQL, and SQLite
- **RPC**: gRPC communication between master and workers

## Build Commands

### Development Build & Run

```bash
# Build both gocron and gocron-node
make build

# Build and run in development mode
make run

# Build with race detector
make build-race

# Run with race detector
make run-race
```

### Testing

```bash
# Run all tests
make test

# Run tests with race detector
make test-race

# Run tests with coverage report
make test-coverage

# Run specific test (example)
go test ./internal/models -v -run TestUser
```

### Code Quality

```bash
# Format code
make fmt

# Check formatting
make fmt-check

# Run vet
make vet

# Run linter
make lint

# Security scan
make security

# Complete pre-release check
make pre-release
```

### Production Packaging

```bash
# Build packages for current platform
make package

# Build for all platforms
make package-all

# Build Linux packages
make package-linux

# Build Windows packages
make package-windows

# Build macOS packages
make package-darwin
```

### Frontend Development

```bash
# Install Vue dependencies
make install-vue

# Build Vue frontend
make build-vue

# Run Vue dev server
make run-vue
```

### Docker

```bash
# Start with Docker Compose
docker-compose up -d

# Access web interface at http://localhost:5920
```

## Architecture Overview

### Core Architecture

The system follows a **master-worker distributed architecture**:

1. **Master (gocron)**: Web interface, task scheduling, database management
2. **Worker (gocron-node)**: Task execution nodes that register with master
3. **Communication**: gRPC with TLS authentication
4. **Database**: GORM with support for multiple backends
5. **Web Framework**: Gin HTTP router

### Key Package Structure

```
cmd/
├── gocron/          # Main web server and scheduler
└── node/           # Worker node executable

internal/
├── models/         # Database models (User, Task, Host, etc.)
├── modules/        # Core functionality modules
│   ├── app/        # Application initialization
│   ├── logger/     # Async logging system
│   ├── notify/     # Email, Slack, Webhook notifications
│   ├── rpc/        # gRPC client/server communication
│   ├── setting/    # Configuration management
│   └── utils/      # Utility functions
├── routers/        # HTTP route handlers
│   ├── agent/      # Agent management
│   ├── task/       # Task CRUD operations
│   ├── user/       # User management & auth
│   └── ...
└── service/        # Business logic layer
    └── task.go     # Task scheduling and execution

web/
└── vue/            # Vue.js frontend application
```

### Database Models

Key entities in the system:

- **User**: Authentication and authorization
- **Task**: Scheduled tasks with cron expressions
- **Host**: Worker nodes that execute tasks
- **TaskLog**: Execution history and results
- **AgentToken**: Authentication tokens for workers
- **Setting**: System configuration

### Task Execution Flow

1. User creates task via web interface
2. Task stored in database with cron schedule
3. Master scheduler evaluates cron expressions
4. When task is due, master sends execution request via gRPC
5. Worker executes task and returns results
6. Results logged to database and notifications sent

### Configuration

Configuration is managed through:

- `app.ini.sqlite.example`: Example configuration file
- Environment variables
- Database settings table
- Command-line flags for both gocron and gocron-node

## Development Guidelines

### Adding New Features

1. **Database Changes**: Add models to `internal/models/`, run migrations
2. **API Endpoints**: Add routes to appropriate router in `internal/routers/`
3. **Business Logic**: Implement in `internal/service/`
4. **Frontend**: Add Vue components in `web/vue/`
5. **Testing**: Write unit tests alongside implementation

### Testing Strategy

- **Unit Tests**: Test individual functions and modules
- **Integration Tests**: Test database interactions and API endpoints
- **Race Detection**: Use `make test-race` for concurrent code
- **Coverage**: Generate reports with `make test-coverage`

### Code Style

- Follow Go standard formatting (`gofmt`)
- Use `make fmt` before commits
- Run `make vet` and `make lint` for quality checks
- Write tests for new functionality

### Security Considerations

- Worker nodes cannot run as root (unless explicitly allowed)
- gRPC communication uses TLS authentication
- User passwords are hashed
- SQL injection protection via GORM
- Input validation on all API endpoints

### Performance Optimizations

- Async logging system to prevent I/O blocking
- Database connection pooling via GORM
- gRPC connection pooling for worker communication
- Efficient cron scheduling algorithm

## Common Development Tasks

### Adding a New API Endpoint

1. Add route handler in appropriate router file
2. Implement business logic in service layer
3. Add database operations in models
4. Update frontend to consume new endpoint
5. Write tests for the new functionality

### Adding a New Database Model

1. Create model struct in `internal/models/`
2. Add migration logic if needed
3. Update `ensureTables()` in `gocron.go` if auto-creation needed
4. Create corresponding service methods
5. Add API endpoints for CRUD operations

### Debugging Task Execution

1. Check worker node logs: `./gocron-node -log-level debug`
2. Verify gRPC connectivity between master and worker
3. Check task logs in database
4. Verify cron expression parsing
5. Review notification settings if alerts not received

### Adding Notification Methods

1. Add new notifier in `internal/modules/notify/`
2. Implement notification interface
3. Add configuration options
4. Update user interface for new notification type
5. Add tests for notification logic

## Environment Setup

### Prerequisites

- Go 1.24.0+
- Node.js/pnpm for frontend development
- Database (MySQL/PostgreSQL/SQLite)
- Docker (optional, for containerized deployment)

### Development Dependencies

```bash
# Install Go development tools
make dev-deps

# Install frontend dependencies
pnpm install
```

### Configuration Files

- `.air.toml`: Live reload configuration for development
- `commitlint.config.cjs`: Commit message linting
- `app.ini.sqlite.example`: Application configuration template
- `docker-compose.yml`: Local development environment

This architecture enables horizontal scaling of worker nodes while maintaining centralized task management and monitoring through the web interface.
