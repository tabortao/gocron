# Gocron Windows 开发环境一键启动脚本

Write-Host "正在检查开发环境..." -ForegroundColor Cyan

# 检查 Go
if (-not (Get-Command go -ErrorAction SilentlyContinue)) {
    Write-Error "未检测到 Go，请先安装 Go 语言环境。"
    exit 1
}

# 检查 GCC (不再必须)
if (Get-Command gcc -ErrorAction SilentlyContinue) {
    Write-Host "检测到 GCC，可以使用 CGO。" -ForegroundColor Gray
} else {
    Write-Host "未检测到 GCC，将使用纯 Go 模式运行。" -ForegroundColor Gray
}

# 检查 pnpm
if (-not (Get-Command pnpm -ErrorAction SilentlyContinue)) {
    Write-Host "未检测到 pnpm，尝试安装..." -ForegroundColor Yellow
    if (Get-Command npm -ErrorAction SilentlyContinue) {
        npm install -g pnpm
    } else {
        Write-Error "未检测到 npm，无法安装 pnpm。"
        exit 1
    }
}

Write-Host "--------------------------------"
Write-Host "1. 安装后端依赖..." -ForegroundColor Green
go mod download
if ($LASTEXITCODE -ne 0) {
    Write-Error "后端依赖安装失败"
    exit 1
}

Write-Host "--------------------------------"
Write-Host "2. 安装前端依赖..." -ForegroundColor Green
Push-Location web/vue
# pnpm install might fail if network issues, but we proceed
cmd /c "pnpm install"
if ($LASTEXITCODE -ne 0) {
    Pop-Location
    Write-Error "前端依赖安装失败"
    exit 1
}
Pop-Location

Write-Host "--------------------------------"
Write-Host "3. 启动服务..." -ForegroundColor Green

# 启动后端 (新窗口)
Write-Host "启动后端服务..."
# 注意: 这里使用 Start-Process 启动新窗口，避免阻塞当前脚本
Start-Process powershell -ArgumentList "-NoExit", "-Command", "$env:CGO_ENABLED='0'; go run cmd/gocron/gocron.go web -e dev"

# 启动前端
Write-Host "启动前端开发服务器..."
Set-Location web/vue
cmd /c "pnpm run dev"
