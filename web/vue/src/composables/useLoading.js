import { ref } from 'vue'
import { ElLoading } from 'element-plus'

export function useLoading(initialState = false) {
  const loading = ref(initialState)
  
  const withLoading = async (fn) => {
    loading.value = true
    try {
      return await fn()
    } finally {
      loading.value = false
    }
  }
  
  return { loading, withLoading }
}

// 全屏 loading
export function useFullScreenLoading() {
  let loadingInstance = null
  
  const show = (text = '加载中...') => {
    loadingInstance = ElLoading.service({
      lock: true,
      text,
      background: 'rgba(0, 0, 0, 0.7)'
    })
  }
  
  const hide = () => {
    loadingInstance?.close()
  }
  
  const withLoading = async (fn, text) => {
    show(text)
    try {
      return await fn()
    } finally {
      hide()
    }
  }
  
  return { show, hide, withLoading }
}
