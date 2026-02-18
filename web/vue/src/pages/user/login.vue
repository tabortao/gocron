<template>
  <div class="login-container">
    <div class="login-box">
      <div class="language-switcher">
        <LanguageSwitcher />
      </div>
      <h2 class="login-title">{{ t('login.title') }}</h2>
      <el-alert
        v-if="errorMessage"
        :title="errorMessage"
        type="error"
        :closable="false"
        style="margin-bottom: 20px"
      />
      <el-form
        ref="formRef"
        :model="form"
        label-width="100px"
        :rules="formRules"
        class="login-form"
      >
        <el-form-item :label="t('login.username')" prop="username">
          <el-input
            v-model.trim="form.username"
            :placeholder="t('login.usernamePlaceholder')"
            autocomplete="username"
            name="username"
            autofocus
            size="large"
          />
        </el-form-item>
        <el-form-item :label="t('login.password')" prop="password">
          <el-input
            v-model.trim="form.password"
            type="password"
            :placeholder="t('login.passwordPlaceholder')"
            autocomplete="current-password"
            name="password"
            show-password
            @keyup.enter="submit"
            size="large"
          />
        </el-form-item>
        <el-form-item :label="t('login.verifyCode')" prop="twoFactorCode" v-if="require2FA">
          <el-input
            v-model.trim="form.twoFactorCode"
            :placeholder="t('login.verifyCodePlaceholder')"
            maxlength="6"
            autocomplete="one-time-code"
            @keyup.enter="submit"
            size="large"
          />
        </el-form-item>
        <el-form-item class="remember-item">
          <el-checkbox v-model="form.rememberMe">{{ t('login.rememberMe') }}</el-checkbox>
        </el-form-item>
        <el-form-item class="submit-item">
          <el-button
            type="primary"
            @click="submit"
            :loading="loading"
            class="login-button"
            size="large"
            >{{ t('login.login') }}</el-button
          >
        </el-form-item>
      </el-form>
    </div>

    <!-- 装饰元素 -->
    <div class="decoration-circle circle-1"></div>
    <div class="decoration-circle circle-2"></div>
    <div class="decoration-circle circle-3"></div>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { useUserStore } from '../../stores/user'
import { useLoading } from '../../composables/useLoading'
import userService from '../../api/user'
import LanguageSwitcher from '../../components/common/LanguageSwitcher.vue'

const { t, locale } = useI18n()

const router = useRouter()
const route = useRoute()
const userStore = useUserStore()
const { loading, withLoading } = useLoading()

const require2FA = ref(false)
const formRef = ref()
const errorMessage = ref('')
const loginPreferenceKey = 'gocron-login-preference'

const form = reactive({
  username: '',
  password: '',
  twoFactorCode: '',
  rememberMe: false
})

onMounted(() => {
  try {
    const raw = localStorage.getItem(loginPreferenceKey)
    if (!raw) return
    const saved = JSON.parse(raw)
    if (saved && typeof saved.username === 'string') {
      form.username = saved.username
    }
    if (saved && typeof saved.rememberMe === 'boolean') {
      form.rememberMe = saved.rememberMe
    }
  } catch (e) {}
})

const syncAutofillValues = () => {
  const usernameInput = document.querySelector('input[name="username"]')
  const passwordInput = document.querySelector('input[name="password"]')
  if (usernameInput && (!form.username || form.username.trim() === '')) {
    form.username = usernameInput.value
  }
  if (passwordInput && (!form.password || form.password.trim() === '')) {
    form.password = passwordInput.value
  }
}

const persistLoginPreference = () => {
  try {
    if (form.rememberMe) {
      localStorage.setItem(
        loginPreferenceKey,
        JSON.stringify({
          username: form.username || '',
          rememberMe: true
        })
      )
    } else {
      localStorage.removeItem(loginPreferenceKey)
    }
  } catch (e) {}
}

const maybeLoadCredential = async () => {
  try {
    if (!navigator.credentials || !window.PasswordCredential) return
    const credential = await navigator.credentials.get({ password: true, mediation: 'optional' })
    if (!credential) return
    if (!form.username) {
      form.username = credential.id || ''
    }
    if (!form.password && credential.password) {
      form.password = credential.password
    }
  } catch (e) {}
}

const maybeStoreCredential = async () => {
  try {
    if (!form.rememberMe) return
    if (!navigator.credentials || !window.PasswordCredential) return
    syncAutofillValues()
    if (!form.username || !form.password) return
    const cred = new window.PasswordCredential({
      id: form.username,
      password: form.password,
      name: form.username
    })
    await navigator.credentials.store(cred)
  } catch (e) {}
}

onMounted(() => {
  maybeLoadCredential()
})

const formRules = computed(() => ({
  username: [{ required: true, message: t('login.usernameRequired'), trigger: 'blur' }],
  password: [{ required: true, message: t('login.passwordRequired'), trigger: 'blur' }],
  twoFactorCode: [{ required: true, message: t('login.verifyCodeRequired'), trigger: 'blur' }]
}))

const submit = async () => {
  if (!formRef.value) return

  errorMessage.value = ''
  syncAutofillValues()

  await formRef.value.validate(async valid => {
    if (!valid) return

    if (require2FA.value && !form.twoFactorCode) {
      errorMessage.value = t('login.verifyCodeRequired')
      return
    }

    await withLoading(async () => {
      const params = {
        username: form.username,
        password: form.password,
        two_factor_code: form.twoFactorCode || undefined
      }

      userService.login(
        params.username,
        params.password,
        params.two_factor_code,
        form.rememberMe,
        data => {
          if (data.require_2fa) {
            require2FA.value = true
            errorMessage.value = ''
            return
          }

          userStore.setUser({
            token: data.token,
            uid: data.uid,
            username: data.username,
            isAdmin: data.is_admin
          })
          persistLoginPreference()
          maybeStoreCredential()

          router.push(route.query.redirect || '/')
        },
        (code, message) => {
          errorMessage.value = message || '登录失败'
        }
      )
    })
  })
}
</script>

<style scoped>
.login-container {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 100vh;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  position: relative;
  overflow: hidden;
  padding: 20px;
}

.decoration-circle {
  position: absolute;
  border-radius: 50%;
  opacity: 0.1;
  pointer-events: none;
}

.circle-1 {
  width: 600px;
  height: 600px;
  background: rgba(255, 255, 255, 0.3);
  top: -200px;
  right: -200px;
  animation: float 20s ease-in-out infinite;
}

.circle-2 {
  width: 400px;
  height: 400px;
  background: rgba(255, 255, 255, 0.2);
  bottom: -100px;
  left: -100px;
  animation: float 15s ease-in-out infinite reverse;
}

.circle-3 {
  width: 200px;
  height: 200px;
  background: rgba(255, 255, 255, 0.15);
  top: 50%;
  left: 10%;
  animation: float 10s ease-in-out infinite;
}

@keyframes float {
  0%,
  100% {
    transform: translateY(0) rotate(0deg);
  }
  50% {
    transform: translateY(-30px) rotate(180deg);
  }
}

.login-box {
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(20px);
  padding: 48px 40px;
  border-radius: 20px;
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.15);
  width: 100%;
  max-width: 420px;
  position: relative;
  z-index: 1;
  border: 1px solid rgba(255, 255, 255, 0.8);
  animation: slideUp 0.5s ease-out;
}

@keyframes slideUp {
  from {
    opacity: 0;
    transform: translateY(30px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.language-switcher {
  position: absolute;
  top: 16px;
  left: 16px;
}

.login-title {
  text-align: center;
  margin: 0 0 32px 0;
  font-size: 28px;
  color: #1f2937;
  font-weight: 700;
  letter-spacing: -0.5px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

.login-form {
  width: 100%;
}

.login-form :deep(.el-form-item) {
  margin-bottom: 24px;
}

.login-form :deep(.el-form-item__label) {
  font-weight: 500;
  color: #374151;
}

.remember-item {
  margin-bottom: 16px;
}

.remember-item :deep(.el-form-item__content) {
  justify-content: flex-start;
  margin-left: 0 !important;
}

.submit-item {
  margin-bottom: 0;
  margin-top: 32px;
}

.login-button {
  width: 100%;
  height: 48px;
  font-size: 16px;
  font-weight: 600;
  border-radius: 12px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border: none;
  transition: all 0.3s ease;
}

.login-button:hover {
  transform: translateY(-2px);
  box-shadow: 0 8px 25px rgba(102, 126, 234, 0.4);
}

.login-button:active {
  transform: translateY(0);
}

@media screen and (max-width: 768px) {
  .login-container {
    padding: 16px;
    align-items: flex-start;
    padding-top: 10vh;
  }

  .login-box {
    padding: 32px 24px;
    border-radius: 16px;
  }

  .login-title {
    font-size: 24px;
    margin-bottom: 24px;
  }

  .login-form :deep(.el-form-item) {
    margin-bottom: 20px;
  }

  .login-form :deep(.el-form-item__label) {
    padding-bottom: 8px;
  }

  .submit-item {
    margin-top: 24px;
  }

  .login-button {
    height: 44px;
    font-size: 15px;
  }

  .circle-1 {
    width: 400px;
    height: 400px;
    top: -150px;
    right: -150px;
  }

  .circle-2 {
    width: 300px;
    height: 300px;
  }

  .circle-3 {
    display: none;
  }
}

@media screen and (max-width: 480px) {
  .login-container {
    padding: 12px;
    padding-top: 8vh;
  }

  .login-box {
    padding: 24px 20px;
    border-radius: 12px;
  }

  .login-title {
    font-size: 22px;
    margin-bottom: 20px;
  }

  .login-form :deep(.el-form-item__label) {
    font-size: 14px;
  }

  .login-button {
    height: 42px;
    font-size: 14px;
    border-radius: 10px;
  }

  .language-switcher {
    top: 12px;
    left: 12px;
  }
}

@media screen and (max-height: 600px) {
  .login-container {
    align-items: flex-start;
    padding-top: 20px;
  }

  .login-box {
    padding: 24px 20px;
  }

  .login-title {
    margin-bottom: 16px;
  }

  .login-form :deep(.el-form-item) {
    margin-bottom: 16px;
  }
}
</style>
