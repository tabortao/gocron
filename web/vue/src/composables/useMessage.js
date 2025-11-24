import { ElMessage, ElMessageBox } from 'element-plus'
import { useI18n } from 'vue-i18n'

export function useMessage() {
  const { t } = useI18n()
  
  const success = (message) => {
    ElMessage.success(message)
  }
  
  const error = (message) => {
    ElMessage.error(message)
  }
  
  const warning = (message) => {
    ElMessage.warning(message)
  }
  
  const info = (message) => {
    ElMessage.info(message)
  }
  
  const confirm = (message, title, options = {}) => {
    return ElMessageBox.confirm(
      message,
      title || t('common.tip'),
      {
        confirmButtonText: t('common.confirm'),
        cancelButtonText: t('common.cancel'),
        type: 'warning',
        center: true,
        ...options
      }
    )
  }
  
  return {
    success,
    error,
    warning,
    info,
    confirm
  }
}
