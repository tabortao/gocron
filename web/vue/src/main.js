import { createApp } from 'vue'
import { createPinia } from 'pinia'
import piniaPluginPersistedstate from 'pinia-plugin-persistedstate'
import { ElMessageBox, ElMessage } from 'element-plus'
import 'element-plus/dist/index.css'
import dayjs from 'dayjs'
import utc from 'dayjs/plugin/utc'
import timezone from 'dayjs/plugin/timezone'
import App from './App.vue'
import router from './router'
import i18n from './locales'

dayjs.extend(utc)
dayjs.extend(timezone)

const app = createApp(App)
const pinia = createPinia()
pinia.use(piniaPluginPersistedstate)

app.use(pinia)
app.use(router)
app.use(i18n)

app.directive('focus', {
  mounted(el) {
    el.focus()
  }
})

app.config.globalProperties.$appConfirm = function (callback) {
  ElMessageBox.confirm(
    i18n.global.t('common.confirmOperation'),
    i18n.global.t('common.tip'),
    {
      confirmButtonText: i18n.global.t('common.confirm'),
      cancelButtonText: i18n.global.t('common.cancel'),
      type: 'warning',
      center: true,
      customClass: 'custom-message-box'
    }
  ).then(() => {
    callback()
  }).catch(() => {})
}

app.config.globalProperties.$message = ElMessage

app.config.globalProperties.$filters = {
  formatTime(time) {
    if (!time) return ''
    return dayjs(time).format('YYYY-MM-DD HH:mm:ss')
  }
}

app.mount('#app')
