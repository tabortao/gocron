// 环境变量验证和获取
export const env = {
  isDev: import.meta.env.DEV,
  isProd: import.meta.env.PROD,
  apiBaseUrl: import.meta.env.VITE_API_BASE_URL || '/api'
}

// 开发环境日志
export const devLog = (...args) => {
  if (env.isDev) {
    console.log('[Dev]', ...args)
  }
}

export const devWarn = (...args) => {
  if (env.isDev) {
    console.warn('[Dev]', ...args)
  }
}

export const devError = (...args) => {
  if (env.isDev) {
    console.error('[Dev]', ...args)
  }
}
