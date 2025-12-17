import { createI18n } from 'vue-i18n'
import zhCN from './zh-CN'
import enUS from './en-US'
import { availableLanguages } from '@/const/index'

const getDefaultLocale = () => {
  const savedLocale = localStorage.getItem('locale')
  return savedLocale || availableLanguages.zhCN.value
}

const i18n = createI18n({
  legacy: false,
  locale: getDefaultLocale(),
  fallbackLocale: availableLanguages.zhCN.value,
  messages: {
    [availableLanguages.zhCN.value]: zhCN,
    [availableLanguages.enUS.value]: enUS
  }
})

export default i18n
