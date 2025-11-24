import axios from 'axios'
import { ElMessage } from 'element-plus'
import router from '../router'
import { useUserStore } from '../stores/user'

const SUCCESS_CODE = 0
const AUTH_ERROR_CODE = 401
const APP_NOT_INSTALL_CODE = 801

const request = axios.create({
  baseURL: '/api',
  timeout: 10000
})

request.interceptors.request.use(
  config => {
    const userStore = useUserStore()
    if (userStore.token) {
      config.headers['Auth-Token'] = userStore.token
    }
    return config
  },
  error => {
    ElMessage.error('请求失败')
    return Promise.reject(error)
  }
)

request.interceptors.response.use(
  response => {
    const { code, message, data } = response.data
    
    if (code === APP_NOT_INSTALL_CODE) {
      router.push('/install')
      return Promise.reject(new Error(message))
    }
    
    if (code === AUTH_ERROR_CODE) {
      const userStore = useUserStore()
      userStore.logout()
      router.push('/user/login')
      return Promise.reject(new Error(message))
    }
    
    if (code !== SUCCESS_CODE) {
      ElMessage.error(message || '请求失败')
      return Promise.reject(new Error(message))
    }
    
    return data
  },
  error => {
    // 网络错误或超时
    if (error.code === 'ECONNABORTED') {
      ElMessage.error('请求超时，请稍后重试')
    } else if (!error.response) {
      ElMessage.error('网络连接失败，请检查网络')
    } else {
      ElMessage.error(error.message || '请求失败')
    }
    return Promise.reject(error)
  }
)

export default request
