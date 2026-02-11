param(
  [Parameter(Mandatory = $true)]
  [string]$Version,

  [switch]$SkipChecks
)

$ErrorActionPreference = "Stop"

function Exec([string]$Cmd) {
  Write-Host $Cmd
  iex $Cmd
}

if ($Version -notmatch '^v\d+\.\d+\.\d+$') {
  throw "Version 格式错误：$Version（需要类似 v1.4.10）"
}

Exec "git rev-parse --is-inside-work-tree | Out-Null"

$status = (git status --porcelain)
if (-not [string]::IsNullOrWhiteSpace($status)) {
  throw "工作区不干净，请先 git add/commit 或清理后再发布。"
}

if (-not $SkipChecks) {
  Exec "pnpm -C web/vue -s install --frozen-lockfile"
  Exec "pnpm -C web/vue -s build"
  Exec "go test ./internal/routers/host ./internal/modules/rpc/... ./internal/service"
}

Exec "git tag -a $Version -m `"Release $Version`""
Exec "git push origin $Version"

Write-Host ""
Write-Host "已推送 tag：$Version"
Write-Host "请到 GitHub Actions 查看 Release Packages 工作流执行状态。"

