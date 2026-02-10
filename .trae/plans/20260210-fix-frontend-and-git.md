# 前端跳转与 CI/CD 配置修复计划

## 任务目标

1.  **前端跳转**：修改前端 Vue 代码，使得点击指定的 SVG 图标跳转到 `https://github.com/tabortao/gocron/`。
2.  **GitHub Actions 文档**：编写文档说明 `DOCKER_HUB_USERNAME`、`DOCKER_HUB_TOKEN`、`DOCKER_HUB_NAMESPACE` 的获取与填写方法，并放入 `docs` 文件夹。
3.  **修复 Git 提交失败**：解决 `lint-staged` 在提交时再次报错的问题。

## 任务分解

### 1. 前端跳转 (Vue)

- [ ] 定位 SVG 所在的文件。根据用户提供的 outer_html，这是一个 GitHub 图标。推测可能在 `Header.vue` 或 `Footer.vue` 中。
- [ ] 搜索代码中的 SVG 路径 `M8 0C3.58 0 0 3.58 0 8c0 3.54 2.29 6.53 5.47 7.59`。
- [ ] 将其包裹在 `<a>` 标签中或添加点击事件，跳转到目标 URL。

### 2. GitHub Actions Secrets 文档

- [ ] 创建 `docs/github_actions_setup.md`。
- [ ] 编写内容：
  - Docker Hub 注册与 Token 生成步骤。
  - GitHub Repository Secrets 配置步骤。
  - 解释三个变量的含义：
    - `DOCKER_HUB_USERNAME`: Docker Hub 用户名。
    - `DOCKER_HUB_TOKEN`: 访问令牌 (Access Token)。
    - `DOCKER_HUB_NAMESPACE`: 镜像仓库的命名空间 (通常也是用户名，或者是组织名)。

### 3. 修复 Git 提交失败 (Lint-staged)

- [ ] 分析错误日志：用户未提供具体的错误信息，只说“又没法提交了”，并引用了之前的错误片段。
- [ ] 回顾之前的修改：之前修改了 `package.json` 中的 `lint-staged` 配置，使用了 `pnpm --dir`。
- [ ] 检查 `package.json`。
- [ ] 问题可能出在 `lint-staged` 对文件列表的传递方式上。当使用 `pnpm --dir web/vue exec ...` 时，`lint-staged` 传递的文件路径是相对于根目录的（如 `web/vue/src/App.vue`）。
- [ ] 如果 `pnpm --dir web/vue` 在子目录执行命令，它可能期望接收相对于该子目录的路径，或者能够处理绝对路径。
- [ ] 另一种可能性：`eslint` 和 `prettier` 在子目录运行时，找不到根目录传递进来的文件（因为路径不匹配）。
- [ ] **解决方案尝试**：
  - 方案 A：不再使用 `--dir`，而是直接在根目录运行 `eslint` 和 `prettier`，并确保它们能正确处理 `web/vue` 下的文件。这通常需要根目录有相应的配置。
  - 方案 B：使用函数式配置（在 `.lintstagedrc.js` 中），手动处理路径，将文件路径转换为绝对路径或相对路径。
  - 方案 C（最简单）：由于项目结构是 Monorepo 风格但没有完全配置 Workspace，尝试将命令改为 `cd web/vue && pnpm exec eslint --fix`，但需要处理 Windows 路径分隔符问题。之前就是因为 Windows `cd` 问题报错。
  - **修正方案**：使用 `git update-index --chmod=+x` 或者调整 `lint-staged` 配置，使其在根目录执行，但指定配置文件路径。
  - 或者，检查是否因为 `husky` 的 `pre-commit` 脚本问题。
  - **重点怀疑**：`lint-staged` 传递给命令的文件列表包含 `web/vue/...`。如果命令是 `pnpm --dir web/vue exec eslint`，`eslint` 接收到的参数是 `web/vue/file.js`。如果 `eslint` 在 `web/vue` 目录下运行，它会去寻找 `web/vue/web/vue/file.js`，这肯定找不到。
  - **正确做法**：应该在根目录运行 eslint，或者使用 `relative` 选项（lint-staged v10+）。
  - 或者，让 `lint-staged` 进入子目录运行：使用 `lerna` 或 `nx` 的方式，或者简单的 `cd`（但要兼容 Windows）。
  - **Windows 兼容的 cd**：在此场景下，最稳妥的是在根目录配置 `eslint` 来检查子目录文件，但这需要修改 eslint 配置。
  - **替代方案**：修改 `package.json`，让命令能够接受根目录相对路径。
  - 如果 `eslint` 在 `web/vue` 运行，它需要相对于 `web/vue` 的路径。
  - 我们可以尝试修改 `lint-staged` 配置，使其匹配 `web/vue` 下的文件时，传递绝对路径（lint-staged 默认传递绝对路径吗？需要确认）。
  - 如果是绝对路径，`eslint` 在任何目录运行都没问题。
  - 如果是相对路径，就有问题。
  - **最终策略**：将 `lint-staged` 命令改回 `cd web/vue && ...` 但确保在 Windows 上能跑。或者使用 `npm-run-all` 等跨平台工具？不，太复杂。
  - **最佳尝试**：使用 `pnpm -C web/vue exec eslint --fix` (等同于 `--dir`)，并确认 lint-staged 传递的是绝对路径。如果传递的是绝对路径，那么 `eslint` 应该能处理。如果报错，可能是 `lint-staged` 版本问题或配置问题。
  - 让我们先读取 `package.json` 确认当前配置。

## 执行步骤

1.  搜索 SVG 并修改。
2.  编写文档。
3.  读取 `package.json`，分析 lint-staged 问题，并尝试修复（可能需要将命令改为在根目录执行，或者使用绝对路径）。
