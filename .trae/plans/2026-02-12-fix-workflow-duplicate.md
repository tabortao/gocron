# 任务计划 - 修复 GitHub Action 重复触发问题

## 任务背景

在 GitHub 上创建 Release 并同时创建标签（如 `v1.5.9`）时，`release-packages.yml` 工作流被触发了两次。

## 原因分析

查看 [release-packages.yml](file:///f:/Code/Go-WorkSpace/gocron/.github/workflows/release-packages.yml) 的触发配置：

```yaml
on:
  push:
    tags:
      - 'v*.*.*'
  release:
    types:
      - published
  workflow_dispatch:
```

- `push: tags`: 当有匹配 `v*.*.*` 的标签推送到仓库时触发。
- `release: published`: 当发布一个新的 Release 时触发。

在 GitHub Web 界面创建 Release 时，如果选择“创建新标签”，GitHub 会先推送标签，然后发布 Release。这两个操作分别触发了上述两个事件，导致工作流运行两次。

## 解决方案更新 (学习 docker-build.yml)

参考 [docker-build.yml](file:///f:/Code/Go-WorkSpace/gocron/.github/workflows/docker-build.yml) 的触发配置，它只使用了 `push` 触发器：

```yaml
on:
  push:
    branches:
      - 'release/**'
    tags:
      - 'v*.*.*'
```

这样做的好处是：

1. **避免重复**：当在 GitHub Web 界面创建 Release 时，标签推送动作会触发 `push: tags`。因为没有配置 `release: published`，所以不会发生第二次触发。
2. **更通用**：无论是通过 `git push --tags` 还是通过 GitHub Release UI 创建标签，都能触发工作流。

## 计划步骤

1. [x] 分析重复触发原因
2. [x] 创建计划文档
3. [x] (已执行) 初步修复：改为 `release: published` (已完成，但根据反馈需进一步优化)
4. [ ] 优化修复：参考 `docker-build.yml`，改为使用 `push: tags` 触发器，并保持对 `workflow_dispatch` 的支持。
5. [ ] 解释修改原因并告知用户。
