# Project Rules

## Windows Path Compatibility & Lint-Staged

**Description:**
在配置 npm scripts 或 lint-staged 时，Windows 环境下的路径处理（如反斜杠、绝对路径）和命令组合（如 `cd`）极易出错。

**Solution:**

1.  **避免使用 `cd`**：不要在脚本中使用 `cd dir && command`。
2.  **使用配置文件**：优先使用 JS/MJS 配置文件 (`.lintstagedrc.mjs`) 以便处理路径转换。
3.  **处理 ESLint 兼容性**：如果项目使用了 ESLint v9 但仍保留旧版配置 (`.eslintrc.*`)，需在配置中设置 `process.env.ESLINT_USE_FLAT_CONFIG = 'false'`，或者暂时禁用 ESLint 检查，优先保证 Prettier 格式化。
4.  如需要使用账号进行调试，可以使用测试账号：testuser，测试账号密码：Testuser123

**Example (.lintstagedrc.mjs):**

```javascript
export default {
  'web/vue/**/*.{js,ts,vue}': filenames => {
    const relativeFiles = filenames.map(f => path.relative('web/vue', f))
    return `pnpm -C web/vue exec eslint --fix ${relativeFiles.join(' ')}`
  }
}
```

## Go Embed Directory

**Description:**
使用 `//go:embed` 嵌入目录时，必须确保目录存在且不为空，否则会导致编译失败。

**Solution:**
在目标目录中放置 `.gitkeep` 或占位文件。

**Example:**
`mkdir -p web/vue/dist && touch web/vue/dist/.gitkeep`

## SQLite CGO Dependency

**Description:**
标准 SQLite 驱动 (`go-sqlite3`) 依赖 CGO 和 GCC，增加了跨平台编译的复杂性（尤其是 Windows）。

**Solution:**
优先使用纯 Go 实现的 SQLite 驱动 `github.com/glebarez/sqlite` (基于 `modernc.org/sqlite`)。
