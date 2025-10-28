# gocron - Cron Task Management System

English | [简体中文](README.md)

A lightweight cron task management system developed in Go, designed to replace Linux-crontab.

## Features

* Web-based task management interface
* Crontab time expressions with second precision
* Task retry on failure
* Task timeout and forced termination
* Task dependency configuration
* Multi-user and permission control
* Two-Factor Authentication (2FA)
* Task types
    * Shell tasks - Execute shell commands on task nodes
    * HTTP tasks - Access specified URLs
* Task execution log viewing
* Task execution notifications (Email, Slack, Webhook)

## Requirements

* Go 1.23+
* MySQL or PostgreSQL
* Node.js 18+ (for frontend development)

## Quick Start

### Development Environment

```bash
# 1. Clone the project
git clone https://github.com/gocronx-team/gocron.git
cd gocron

# 2. Install dependencies
go mod download

# 3. Configure database
# Edit ~/.gocron/conf/app.ini

# 4. Start backend (with hot reload)
air

# 5. Start frontend (in another terminal)
cd web/vue
npm install
npm run dev
```

Visit http://localhost:8080

### Production Deployment

```bash
# 1. Build
make package

# 2. Start service
./gocron web

# 3. Start task node
./gocron-node
```

Visit http://localhost:5920

## Commands

### gocron

```bash
gocron web              # Start web service
gocron web -p 8080      # Specify port
gocron web -e dev       # Development mode
gocron -v               # Show version
```

### gocron-node

```bash
gocron-node             # Start task node
gocron-node -s :5921    # Specify listening address
gocron-node -enable-tls # Enable TLS
```

## Tech Stack

* Backend: Gin + GORM + gRPC
* Frontend: Vue3 + Element Plus + Vite
* Scheduler: Cron
* Database: MySQL / PostgreSQL

## Development Tools

* `make` - Build project
* `make run` - Build and run
* `air` - Backend hot reload
* `npm run dev` - Frontend hot reload

## License

MIT
