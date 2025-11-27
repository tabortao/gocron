import { ref, watch, onUnmounted } from 'vue'

export function useDebounce(value, delay = 300) {
  const debouncedValue = ref(value.value)
  let timeout = null

  watch(value, (newValue) => {
    if (timeout) clearTimeout(timeout)
    timeout = setTimeout(() => {
      debouncedValue.value = newValue
    }, delay)
  })
  
  // 组件卸载时清理
  onUnmounted(() => {
    if (timeout) clearTimeout(timeout)
  })

  return debouncedValue
}

export function useDebounceFn(fn, delay = 300) {
  let timeout = null
  
  return (...args) => {
    if (timeout) clearTimeout(timeout)
    timeout = setTimeout(() => {
      fn(...args)
    }, delay)
  }
}
