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
      <el-form ref="formRef" :model="form" label-width="100px" :rules="formRules">
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
        <el-form-item>
          <el-checkbox v-model="form.rememberMe">{{ t('login.rememberMe') }}</el-checkbox>
        </el-form-item>
        <el-form-item>
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
  background: linear-gradient(135deg, #f5f7fa 0%, #c3cfe2 100%);
  position: relative;
  overflow: hidden;
}

.login-container::before {
  content: '';
  position: absolute;
  top: -50%;
  right: -10%;
  width: 600px;
  height: 600px;
  background: rgba(99, 102, 241, 0.1);
  border-radius: 50%;
  filter: blur(80px);
}

.login-container::after {
  content: '';
  position: absolute;
  bottom: -30%;
  left: -10%;
  width: 500px;
  height: 500px;
  background: rgba(168, 85, 247, 0.08);
  border-radius: 50%;
  filter: blur(80px);
}

.login-box {
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(10px);
  padding: 48px 40px;
  border-radius: 16px;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.08);
  width: 100%;
  max-width: 420px;
  position: relative;
  z-index: 1;
  border: 1px solid rgba(255, 255, 255, 0.8);
}

.language-switcher {
  position: absolute;
  top: 16px;
  left: 16px;
}

.login-title {
  text-align: center;
  margin: 0 0 32px 0;
  font-size: 26px;
  color: #1f2937;
  font-weight: 600;
  letter-spacing: -0.5px;
}

.el-button--large {
  height: 40px;
  line-height: 40px;
  padding: 0 15px;
}

.login-button {
  width: calc(100% + 60px);
  margin-left: -60px;
}
</style>
