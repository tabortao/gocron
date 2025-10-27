import { ref } from 'vue'

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
