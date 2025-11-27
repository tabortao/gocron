// 性能监控工具
export const measurePerformance = (name, fn) => {
  if (import.meta.env.DEV) {
    const start = performance.now()
    const result = fn()
    const end = performance.now()
    console.log(`[Performance] ${name}: ${(end - start).toFixed(2)}ms`)
    return result
  }
  return fn()
}

// 监控路由切换性能
export const measureRouteChange = (to, from) => {
  if (import.meta.env.DEV && performance.mark) {
    performance.mark(`route-${to.path}-start`)
  }
}

export const measureRouteChangeEnd = (to) => {
  if (import.meta.env.DEV && performance.mark && performance.measure) {
    performance.mark(`route-${to.path}-end`)
    try {
      performance.measure(
        `route-${to.path}`,
        `route-${to.path}-start`,
        `route-${to.path}-end`
      )
      const measure = performance.getEntriesByName(`route-${to.path}`)[0]
      console.log(`[Route Performance] ${to.path}: ${measure.duration.toFixed(2)}ms`)
    } catch (e) {
      // ignore
    }
  }
}
