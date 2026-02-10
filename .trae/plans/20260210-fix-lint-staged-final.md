# 修复 lint-staged 配置与规则更新计划

## 任务目标

1.  **彻底修复 lint-staged 问题**：解决 Windows 下 `cd` 命令语法错误，以及 `prettier` 找不到文件的错误。
2.  **原因分析与教育**：解释错误原因（Windows 命令分隔符、绝对路径与相对路径的传递机制）。
3.  **规则更新**：将精简后的最佳实践写入 `.trae/rules/project_rules.md`。

## 问题分析

1.  **`cd web/vue && pnpm exec eslint --fix` 错误**:
    - 错误信息：`�ļ�����Ŀ¼��������﷨����ȷ��` (文件名、目录名或卷标语法不正确)。
    - 原因：在 Windows `cmd.exe` (lint-staged 默认 shell) 中，`cd` 命令不支持同时接受文件参数。`lint-staged` 会将文件列表追加到命令末尾。例如：`cd web/vue && pnpm exec eslint --fix F:\Code\...\App.vue`。这被解析为 `cd` 的参数不正确，或者路径传递给了 `eslint` 但 `eslint` 此时在子目录运行，可能无法正确解析绝对路径（虽然通常绝对路径没问题，但 `cd` 这一步就挂了）。
    - **关键点**：`lint-staged` 自动追加的文件路径是**绝对路径**。

2.  **`prettier --write` 错误**:
    - 错误信息：`[error] No files matching the pattern were found: "F:/Code/..."`。
    - 原因：Prettier 接收到了绝对路径，但可能配置为仅处理特定 glob，或者路径格式（正反斜杠）在 Windows 下导致了匹配失败。通常 Prettier 可以处理绝对路径，但如果 glob pattern 是 `*.{json,md}` 且没有正确引用，或者 lint-staged 在传递参数时出现了转义问题。
    - **更深层原因**：Prettier 3.0+ 某些版本对路径处理变严格，或者 glob 模式与绝对路径混合使用时的行为差异。

## 解决方案

1.  **放弃 `cd` 组合命令**。这在 lint-staged 中几乎总是导致跨平台问题。
2.  **使用 `pnpm --dir` 或 `-C`**：这比 `cd` 更健壮，因为 `pnpm` 会正确处理上下文。但是之前用户反馈 `pnpm --dir` 也有问题。
    - 如果 `pnpm --dir web/vue exec eslint` 接收到 `F:\Code\...\App.vue`，`eslint` 应该能处理。
    - 之前的错误可能是路径分隔符问题。
3.  **终极方案：在根目录运行**。
    - 对于 `eslint`：在根目录运行，使用 `--config web/vue/.eslintrc.cjs` (或类似) 并指定目标文件。或者，如果项目配置了 Workspace，直接在根目录 `pnpm eslint`。
    - 但当前项目结构较简单。
    - **推荐尝试**：使用 `git update-index --chmod=+x` 等方式？不，这是文件权限。
    - **正确修复**：使用 `npm-run-all` 或 `cross-env`？太重。
    - **回归本质**：`lint-staged` 支持函数式配置。在 `.lintstagedrc.js` 中编写函数，接收文件列表，手动过滤并构造命令。这能完美解决路径问题（比如转为相对路径）。
    - **方案 D (最稳妥)**：将 `package.json` 中的 `lint-staged` 配置提取到 `.lintstagedrc.mjs` (ESM)，使用函数接收文件名，然后：
      - 对 Vue 文件：计算相对 `web/vue` 的路径，然后构造 `pnpm -C web/vue exec eslint --fix <rel_paths>`。
      - 这样避免了绝对路径带来的种种问题。

## 执行步骤

1.  **创建 `.lintstagedrc.mjs`**：替代 `package.json` 中的配置。
2.  **编写路径处理逻辑**：
    - 过滤文件。
    - 使用 `path.relative` 将绝对路径转换为相对路径。
    - 构造正确的 `pnpm` 命令。
3.  **清理 `package.json`**：移除 `lint-staged` 字段。
4.  **更新规则文件**：记录“使用配置文件而非 package.json 来处理复杂路径”的规则。

## 验证计划

- 用户再次尝试提交。

## 具体代码逻辑 (.lintstagedrc.mjs)

```javascript
import path from 'path'

export default {
  'web/vue/**/*.{js,ts,vue}': filenames => {
    // 将绝对路径转换为相对于 web/vue 的路径
    const relativeFiles = filenames.map(f => path.relative(path.join(process.cwd(), 'web/vue'), f))
    // 构造命令：在 web/vue 目录下执行，传入相对路径
    // 注意：文件名可能包含空格，需要引用，但在 exec 中通常数组形式更安全，lint-staged JS 配置返回字符串或数组。
    // 如果返回数组，lint-staged 会执行多个命令。
    // 这里我们返回单条命令字符串。
    return `pnpm -C web/vue exec eslint --fix ${relativeFiles.join(' ')}`
    // 同理处理 prettier
  },
  '*.{json,md,yml,yaml}': 'prettier --write'
}
```

_修正_：ESLint 和 Prettier 其实都支持绝对路径。问题核心通常是 `cd` 命令干扰了参数传递。如果使用 `pnpm -C web/vue exec eslint --fix` 加上绝对路径，理论上可行。但为了保险，相对路径更不易出错。

让我们先尝试**不使用 cd**，直接用 `pnpm -C`，且**不转换路径**（因为转换路径稍显复杂）。如果失败，再上 JS 配置文件。
等等，用户之前的报错 `pnpm --dir` 失败，可能是因为 `eslint` 配置在子目录，当传入绝对路径时，它找不到对应的配置（如果配置是相对查找的）。

**决定采用 `.lintstagedrc.mjs` + 相对路径** 方案，这是最稳健的，能彻底屏蔽 OS 差异和路径解析怪癖。
