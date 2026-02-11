#!/usr/bin/env bash

set -euo pipefail

VERSION=""
SKIP_CHECKS=false

while [[ $# -gt 0 ]]; do
  case "$1" in
    -Version|--version)
      VERSION="$2"
      shift 2
      ;;
    --skip-checks)
      SKIP_CHECKS=true
      shift
      ;;
    *)
      echo "Unknown option: $1"
      exit 1
      ;;
  esac
done

if [[ -z "$VERSION" ]]; then
  echo "Usage: $0 -Version v1.4.10 [--skip-checks]"
  exit 1
fi

if [[ ! "$VERSION" =~ ^v[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
  echo "Invalid version: $VERSION (expected v1.4.10)"
  exit 1
fi

git rev-parse --is-inside-work-tree >/dev/null

if [[ -n "$(git status --porcelain)" ]]; then
  echo "Working tree is dirty. Commit or clean before release."
  exit 1
fi

if [[ "$SKIP_CHECKS" == "false" ]]; then
  pnpm -C web/vue install --frozen-lockfile
  pnpm -C web/vue run build
  go test ./internal/routers/host ./internal/modules/rpc/... ./internal/service
fi

git tag -a "$VERSION" -m "Release $VERSION"
git push origin "$VERSION"

echo ""
echo "Tag pushed: $VERSION"
echo "Check GitHub Actions: Release Packages"

