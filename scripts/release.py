#!/usr/bin/env python3

import argparse
import re
import subprocess
import sys


SEMVER_TAG_RE = re.compile(r"^v\d+\.\d+\.\d+$")


def run(cmd: list[str], *, check: bool = True) -> subprocess.CompletedProcess:
    return subprocess.run(cmd, check=check, text=True, capture_output=True, encoding="utf-8")


def run_print(cmd: list[str]) -> None:
    p = run(cmd, check=False)
    if p.stdout:
        sys.stdout.write(p.stdout)
    if p.stderr:
        sys.stderr.write(p.stderr)
    if p.returncode != 0:
        raise subprocess.CalledProcessError(p.returncode, cmd, output=p.stdout, stderr=p.stderr)


def get_current_branch() -> str:
    p = run(["git", "rev-parse", "--abbrev-ref", "HEAD"])
    return p.stdout.strip()


def is_dirty() -> bool:
    p = run(["git", "status", "--porcelain"])
    return bool(p.stdout.strip())


def has_staged_changes() -> bool:
    p = run(["git", "diff", "--cached", "--name-only"])
    return bool(p.stdout.strip())


def tag_exists(tag: str) -> bool:
    p = run(["git", "rev-parse", "--verify", f"refs/tags/{tag}"], check=False)
    return p.returncode == 0


def main() -> int:
    parser = argparse.ArgumentParser(description="Auto git commit/push and tag/push for GitHub Release.")
    parser.add_argument("--version", "-v", required=True, help="Release tag like v1.5.4")
    parser.add_argument("--message", "-m", default="", help="Commit message (default: chore: release <tag>)")
    parser.add_argument("--remote", default="origin", help="Git remote name (default: origin)")
    parser.add_argument("--branch", default="", help="Branch to push (default: current branch)")
    parser.add_argument("--no-commit", action="store_true", help="Skip commit even if there are changes")
    args = parser.parse_args()

    tag = args.version.strip()
    if not SEMVER_TAG_RE.match(tag):
        sys.stderr.write(f"Version 格式错误：{tag}（需要类似 v1.5.4）\n")
        return 2

    try:
        run_print(["git", "rev-parse", "--is-inside-work-tree"])
    except subprocess.CalledProcessError:
        sys.stderr.write("当前目录不是 git 仓库。\n")
        return 2

    branch = args.branch.strip() or get_current_branch()
    commit_msg = args.message.strip() or f"chore: release {tag}"

    try:
        if is_dirty():
            if args.no_commit:
                sys.stderr.write("工作区不干净，但指定了 --no-commit；请先自行提交或清理。\n")
                return 2

            run_print(["git", "add", "-A"])
            if has_staged_changes():
                run_print(["git", "commit", "-m", commit_msg])

        run_print(["git", "push", args.remote, branch])

        if tag_exists(tag):
            sys.stderr.write(f"Tag 已存在：{tag}。请换一个版本号，或先手动删除远程 tag。\n")
            return 2

        run_print(["git", "tag", "-a", tag, "-m", f"Release {tag}"])
        run_print(["git", "push", args.remote, tag])
    except subprocess.CalledProcessError as e:
        sys.stderr.write(f"命令执行失败：{' '.join(e.cmd)}\n")
        return e.returncode or 1

    sys.stdout.write(f"\n完成：已推送分支 {branch}，并推送 tag {tag}。\n")
    sys.stdout.write("随后 GitHub Actions 'Release Packages' 会自动构建并发布 Release。\n")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())

