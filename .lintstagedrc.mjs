import path from 'path';

// 尝试兼容旧版 ESLint 配置 (如果未来启用 eslint)
process.env.ESLINT_USE_FLAT_CONFIG = 'false';

export default {
  // 处理前端 Vue 代码
  'web/vue/**/*.{js,ts,vue}': (filenames) => {
    if (filenames.length === 0) return [];

    const vueDir = path.join(process.cwd(), 'web/vue');
    
    // 转为相对路径，统一使用正斜杠
    const relativeFiles = filenames.map(f => path.relative(vueDir, f).split(path.sep).join('/'));

    return [
      // 暂时禁用 ESLint，因为项目升级到 ESLint v9 后尚未迁移到 flat config，
      // 且在 lint-staged 环境下加载旧版配置和插件存在路径兼容性问题。
      // 优先保证代码格式化和提交顺畅。
      // `pnpm -C web/vue exec eslint --fix ${relativeFiles.join(' ')}`,
      
      // 仅运行 Prettier
      `pnpm -C web/vue exec prettier --write ${relativeFiles.join(' ')}`
    ];
  },
  
  // 处理根目录及其他通用文件
  '*.{json,md,yml,yaml}': (filenames) => {
      const relativeFiles = filenames.map(f => path.relative(process.cwd(), f).split(path.sep).join('/'));
      return `prettier --write ${relativeFiles.join(' ')}`;
  }
};
