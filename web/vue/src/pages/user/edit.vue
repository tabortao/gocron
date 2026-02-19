<template>
  <div class="page-container">
    <div class="page-header">
      <div class="page-header__left">
        <h1 class="page-title">{{ isEdit ? t('user.editUser') : t('user.addUser') }}</h1>
        <p class="page-desc" v-if="isEdit && form.name">
          {{ t('user.editingUser') }}: <span class="highlight">{{ form.name }}</span>
        </p>
      </div>
      <div class="page-header__right">
        <el-button @click="cancel" class="ghost-btn">
          <el-icon><ArrowLeft /></el-icon>
          <span class="btn-text">{{ t('common.back') }}</span>
        </el-button>
      </div>
    </div>

    <div class="form-card">
      <el-form
        ref="formRef"
        :model="form"
        :rules="formRules"
        label-position="top"
        class="user-form"
      >
        <div class="form-section">
          <div class="form-section__title">{{ t('user.basicInfo') }}</div>
          <div class="form-grid">
            <el-form-item :label="t('user.username')" prop="name">
              <el-input v-model="form.name" :placeholder="t('user.usernamePlaceholder')" />
            </el-form-item>
            <el-form-item :label="t('user.email')" prop="email">
              <el-input v-model="form.email" :placeholder="t('user.emailPlaceholder')" />
            </el-form-item>
          </div>
        </div>

        <div class="form-section" v-if="!isEdit">
          <div class="form-section__title">{{ t('user.passwordInfo') }}</div>
          <div class="form-grid">
            <el-form-item :label="t('user.password')" prop="password">
              <el-input
                v-model="form.password"
                type="password"
                :placeholder="t('user.passwordPlaceholder')"
                show-password
              />
            </el-form-item>
            <el-form-item :label="t('user.confirmPassword')" prop="confirm_password">
              <el-input
                v-model="form.confirm_password"
                type="password"
                :placeholder="t('user.passwordPlaceholder')"
                show-password
              />
            </el-form-item>
          </div>
        </div>

        <div class="form-section">
          <div class="form-section__title">{{ t('user.permissionInfo') }}</div>
          <div class="form-grid">
            <el-form-item :label="t('user.role')" prop="is_admin">
              <el-radio-group v-model="form.is_admin">
                <el-radio :value="0">{{ t('user.normalUser') }}</el-radio>
                <el-radio :value="1">{{ t('user.admin') }}</el-radio>
              </el-radio-group>
            </el-form-item>
            <el-form-item :label="t('common.status')" prop="status">
              <el-radio-group v-model="form.status">
                <el-radio :value="1">{{ t('common.enabled') }}</el-radio>
                <el-radio :value="0">{{ t('common.disabled') }}</el-radio>
              </el-radio-group>
            </el-form-item>
          </div>
        </div>

        <div class="form-actions">
          <el-button @click="cancel">{{ t('common.cancel') }}</el-button>
          <el-button type="primary" @click="submit" :loading="submitting">
            <el-icon><Check /></el-icon>
            {{ t('common.save') }}
          </el-button>
        </div>
      </el-form>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, watch, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRouter, useRoute } from 'vue-router'
import { ElMessage } from 'element-plus'
import { ArrowLeft, Check } from '@element-plus/icons-vue'
import userService from '../../api/user'

const { t } = useI18n()
const router = useRouter()
const route = useRoute()

const formRef = ref(null)
const submitting = ref(false)

const form = ref({
  id: '',
  name: '',
  email: '',
  is_admin: 0,
  password: '',
  confirm_password: '',
  status: 1
})

const isEdit = computed(() => !!route.params.id)

const formRules = computed(() => {
  const rules = {
    name: [{ required: true, message: t('user.usernameRequired'), trigger: 'blur' }],
    email: [{ type: 'email', required: true, message: t('user.emailRequired'), trigger: 'blur' }]
  }

  if (!isEdit.value) {
    rules.password = [{ required: true, message: t('user.passwordRequired'), trigger: 'blur' }]
    rules.confirm_password = [
      { required: true, message: t('user.confirmPasswordRequired'), trigger: 'blur' }
    ]
  }

  return rules
})

const loadUser = () => {
  const id = route.params.id
  if (!id) return

  userService.detail(id, data => {
    if (!data) {
      ElMessage.error(t('message.dataNotFound'))
      router.push('/user')
      return
    }
    form.value = {
      id: data.id,
      name: data.name,
      email: data.email,
      is_admin: data.is_admin,
      password: '',
      confirm_password: '',
      status: data.status
    }
  })
}

const submit = () => {
  formRef.value.validate(valid => {
    if (!valid) return
    save()
  })
}

const save = () => {
  submitting.value = true
  userService.update(form.value, () => {
    submitting.value = false
    ElMessage.success(t('message.saveSuccess'))
    router.push('/user')
  })
}

const cancel = () => {
  router.push('/user')
}

onMounted(() => {
  loadUser()
})
</script>

<style scoped>
.page-container {
  max-width: 1200px;
  margin: 0 auto;
  padding: 24px;
  min-height: 100%;
  background: #f5f7fa;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 20px;
}

.page-title {
  font-size: 22px;
  font-weight: 600;
  color: #1f2937;
  margin: 0 0 6px 0;
}

.page-desc {
  font-size: 14px;
  color: #6b7280;
  margin: 0;
}

.page-desc .highlight {
  font-weight: 600;
  color: #3b82f6;
}

.page-header__right {
  display: flex;
  gap: 10px;
}

.ghost-btn {
  border: 1px solid #e5e7eb;
  background: #fff;
  color: #374151;
  font-weight: 500;
}

.ghost-btn:hover {
  border-color: #3b82f6;
  color: #3b82f6;
}

.btn-text {
  margin-left: 6px;
}

.form-card {
  background: #fff;
  border-radius: 12px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.05);
  padding: 24px;
}

.form-section {
  margin-bottom: 24px;
  padding-bottom: 24px;
  border-bottom: 1px solid #f3f4f6;
}

.form-section:last-of-type {
  margin-bottom: 0;
  padding-bottom: 0;
  border-bottom: none;
}

.form-section__title {
  font-size: 15px;
  font-weight: 600;
  color: #374151;
  margin-bottom: 16px;
  padding-left: 12px;
  border-left: 3px solid #3b82f6;
}

.form-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 20px;
}

.user-form :deep(.el-form-item__label) {
  font-weight: 500;
  color: #374151;
  padding-bottom: 8px;
}

.user-form :deep(.el-radio-group) {
  display: flex;
  gap: 24px;
}

.form-actions {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
  margin-top: 24px;
  padding-top: 24px;
  border-top: 1px solid #f3f4f6;
}

@media screen and (max-width: 768px) {
  .page-container {
    padding: 16px;
  }

  .page-header {
    flex-direction: column;
    gap: 12px;
  }

  .page-header__right {
    width: 100%;
  }

  .ghost-btn {
    width: 100%;
    justify-content: center;
  }

  .page-title {
    font-size: 20px;
  }

  .form-card {
    padding: 16px;
  }

  .form-grid {
    grid-template-columns: 1fr;
  }

  .form-actions {
    flex-direction: column-reverse;
  }

  .form-actions .el-button {
    width: 100%;
  }
}
</style>
