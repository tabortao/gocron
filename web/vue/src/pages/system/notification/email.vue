<template>
  <div class="notification-page">
    <notification-tab />

    <!-- 服务器配置卡片 -->
    <div class="config-card">
      <div class="card-header">
        <div class="header-icon" style="background: linear-gradient(135deg, #3b82f6, #1d4ed8)">
          <el-icon><Message /></el-icon>
        </div>
        <div class="header-text">
          <h2>{{ t('system.emailServerConfig') }}</h2>
          <p>配置 SMTP 服务器以发送邮件通知</p>
        </div>
      </div>

      <el-form
        ref="formRef"
        :model="form"
        :rules="formRules"
        label-position="top"
        class="config-form"
      >
        <div class="form-grid">
          <el-form-item :label="t('system.smtpHost')" prop="host" class="form-item-grow">
            <el-input v-model="form.host" placeholder="smtp.example.com">
              <template #prefix>
                <el-icon><Monitor /></el-icon>
              </template>
            </el-input>
          </el-form-item>
          <el-form-item :label="t('host.port')" prop="port">
            <el-input v-model.number="form.port" placeholder="465" type="number">
              <template #prefix>
                <el-icon><Connection /></el-icon>
              </template>
            </el-input>
          </el-form-item>
        </div>

        <div class="form-grid">
          <el-form-item :label="t('user.username')" prop="user">
            <el-input v-model="form.user" :placeholder="t('user.username')">
              <template #prefix>
                <el-icon><User /></el-icon>
              </template>
            </el-input>
          </el-form-item>
          <el-form-item :label="t('user.password')" prop="password">
            <el-input
              v-model="form.password"
              type="password"
              show-password
              :placeholder="t('user.password')"
            >
              <template #prefix>
                <el-icon><Lock /></el-icon>
              </template>
            </el-input>
          </el-form-item>
        </div>

        <el-form-item :label="t('system.template')" prop="template">
          <el-input
            type="textarea"
            :rows="6"
            :placeholder="emailPlaceholder"
            v-model="form.template"
            class="template-input"
          />
        </el-form-item>

        <div class="form-actions">
          <el-button type="primary" @click="submit" :loading="saving">
            <el-icon><Check /></el-icon>
            {{ t('common.save') }}
          </el-button>
        </div>
      </el-form>
    </div>

    <!-- 通知用户卡片 -->
    <div class="users-card">
      <div class="card-header">
        <div class="header-icon" style="background: linear-gradient(135deg, #10b981, #059669)">
          <el-icon><UserFilled /></el-icon>
        </div>
        <div class="header-text">
          <h2>{{ t('system.notificationUsers') }}</h2>
          <p>管理接收邮件通知的用户</p>
        </div>
        <el-button type="primary" @click="dialogVisible = true" class="add-btn">
          <el-icon><Plus /></el-icon>
          {{ t('system.addUser') }}
        </el-button>
      </div>

      <div class="users-list" v-if="receivers.length > 0">
        <div class="user-item" v-for="item in receivers" :key="item.id">
          <div class="user-avatar">{{ item.username?.charAt(0)?.toUpperCase() || 'U' }}</div>
          <div class="user-info">
            <span class="user-name">{{ item.username }}</span>
            <span class="user-email">{{ item.email }}</span>
          </div>
          <button class="delete-btn" @click="deleteUser(item)">
            <el-icon><Close /></el-icon>
          </button>
        </div>
      </div>

      <div class="empty-users" v-else>
        <el-icon :size="48"><User /></el-icon>
        <p>{{ t('system.noNotificationUsers') }}</p>
        <el-button type="primary" @click="dialogVisible = true">
          <el-icon><Plus /></el-icon>
          {{ t('system.addUser') }}
        </el-button>
      </div>
    </div>

    <!-- 添加用户对话框 -->
    <el-dialog
      v-model="dialogVisible"
      :title="t('system.addUser')"
      :width="isMobile ? '90%' : '400px'"
      class="user-dialog"
    >
      <el-form :model="{ username, email }" label-position="top">
        <el-form-item :label="t('user.username')">
          <el-input v-model.trim="username" :placeholder="t('user.username')">
            <template #prefix>
              <el-icon><User /></el-icon>
            </template>
          </el-input>
        </el-form-item>
        <el-form-item :label="t('system.emailAddress')">
          <el-input v-model.trim="email" placeholder="user@example.com">
            <template #prefix>
              <el-icon><Message /></el-icon>
            </template>
          </el-input>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">{{ t('common.cancel') }}</el-button>
        <el-button type="primary" @click="saveUser">{{ t('common.confirm') }}</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { ElMessage } from 'element-plus'
import {
  Message,
  Monitor,
  Connection,
  User,
  Lock,
  Check,
  UserFilled,
  Plus,
  Close
} from '@element-plus/icons-vue'
import notificationTab from './tab.vue'
import notificationService from '../../../api/notification'

const { t } = useI18n()

const formRef = ref()
const saving = ref(false)
const dialogVisible = ref(false)
const isMobile = ref(false)

const form = ref({
  host: '',
  port: 465,
  user: '',
  password: '',
  template: ''
})

const receivers = ref([])
const username = ref('')
const email = ref('')

const checkMobile = () => {
  isMobile.value = window.innerWidth <= 768
}

const emailPlaceholder = computed(
  () =>
    `${t('system.taskIdVar')}: {{.TaskId}}
${t('system.taskNameVar')}: {{.TaskName}}
${t('system.statusVar')}: {{.Status}}
${t('system.resultVar')}: {{.Result}}
${t('task.remark')}: {{.Remark}}`
)

const formRules = computed(() => ({
  host: [{ required: true, message: t('system.pleaseEnterEmailServer'), trigger: 'blur' }],
  port: [
    { type: 'number', required: true, message: t('system.pleaseEnterValidPort'), trigger: 'blur' }
  ],
  user: [{ required: true, message: t('system.pleaseEnterUserEmail'), trigger: 'blur' }],
  password: [{ required: true, message: t('user.passwordRequired'), trigger: 'blur' }],
  template: [{ required: true, message: t('system.pleaseEnterTemplate'), trigger: 'blur' }]
}))

const fetchData = () => {
  notificationService.mail(data => {
    form.value.host = data.host || ''
    form.value.port = data.port || 465
    form.value.user = data.user || ''
    form.value.password = data.password || ''
    form.value.template = data.template || ''
    receivers.value = data.mail_users || []
  })
}

const submit = async () => {
  if (!formRef.value) return
  await formRef.value.validate(valid => {
    if (!valid) return
    saving.value = true
    notificationService.updateMail(form.value, () => {
      ElMessage.success(t('message.updateSuccess'))
      saving.value = false
      fetchData()
    })
  })
}

const saveUser = () => {
  if (!username.value || !email.value) {
    ElMessage.error(t('system.incompleteParameters'))
    return
  }
  notificationService.createMailUser(
    {
      username: username.value,
      email: email.value
    },
    () => {
      dialogVisible.value = false
      username.value = ''
      email.value = ''
      fetchData()
    }
  )
}

const deleteUser = item => {
  notificationService.removeMailUser(item.id, () => {
    fetchData()
  })
}

onMounted(() => {
  checkMobile()
  window.addEventListener('resize', checkMobile)
  fetchData()
})

onUnmounted(() => {
  window.removeEventListener('resize', checkMobile)
})
</script>

<style scoped>
.notification-page {
  padding: 20px;
  max-width: 900px;
  margin: 0 auto;
}

.config-card,
.users-card {
  background: white;
  border-radius: 20px;
  overflow: hidden;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.06);
  margin-bottom: 20px;
}

.card-header {
  display: flex;
  align-items: center;
  gap: 16px;
  padding: 24px;
  background: linear-gradient(135deg, #f8fafc 0%, #f1f5f9 100%);
  border-bottom: 1px solid #f1f5f9;
}

.header-icon {
  width: 48px;
  height: 48px;
  border-radius: 14px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  font-size: 22px;
  flex-shrink: 0;
}

.header-text {
  flex: 1;
}

.header-text h2 {
  margin: 0 0 4px 0;
  font-size: 18px;
  font-weight: 600;
  color: #1e293b;
}

.header-text p {
  margin: 0;
  font-size: 13px;
  color: #64748b;
}

.add-btn {
  flex-shrink: 0;
}

.config-form {
  padding: 24px;
}

.form-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 16px;
}

.form-item-grow {
  grid-column: span 1;
}

.template-input :deep(.el-textarea__inner) {
  font-family: 'SF Mono', Monaco, monospace;
  font-size: 13px;
  background: #f8fafc;
  border-radius: 12px;
}

.form-actions {
  display: flex;
  justify-content: flex-end;
  padding-top: 8px;
}

.users-list {
  padding: 16px 24px 24px;
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
  gap: 12px;
}

.user-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 14px 16px;
  background: #f8fafc;
  border-radius: 14px;
  transition: all 0.2s;
}

.user-item:hover {
  background: #f1f5f9;
}

.user-avatar {
  width: 40px;
  height: 40px;
  border-radius: 12px;
  background: linear-gradient(135deg, #3b82f6, #1d4ed8);
  color: white;
  display: flex;
  align-items: center;
  justify-content: center;
  font-weight: 600;
  font-size: 16px;
  flex-shrink: 0;
}

.user-info {
  flex: 1;
  min-width: 0;
}

.user-name {
  display: block;
  font-weight: 500;
  color: #1e293b;
  margin-bottom: 2px;
}

.user-email {
  display: block;
  font-size: 12px;
  color: #64748b;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.delete-btn {
  width: 32px;
  height: 32px;
  border: none;
  background: transparent;
  border-radius: 8px;
  cursor: pointer;
  color: #94a3b8;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.2s;
}

.delete-btn:hover {
  background: #fef2f2;
  color: #ef4444;
}

.empty-users {
  padding: 48px 24px;
  text-align: center;
  color: #94a3b8;
}

.empty-users .el-icon {
  margin-bottom: 12px;
}

.empty-users p {
  margin: 0 0 16px 0;
}

.user-dialog :deep(.el-dialog__body) {
  padding: 20px 24px;
}

@media screen and (max-width: 768px) {
  .notification-page {
    padding: 12px;
  }

  .card-header {
    flex-wrap: wrap;
    padding: 16px;
  }

  .header-icon {
    width: 40px;
    height: 40px;
    font-size: 18px;
  }

  .header-text h2 {
    font-size: 16px;
  }

  .add-btn {
    width: 100%;
    margin-top: 12px;
  }

  .config-form {
    padding: 16px;
  }

  .form-grid {
    grid-template-columns: 1fr;
  }

  .users-list {
    padding: 12px 16px 16px;
    grid-template-columns: 1fr;
  }
}
</style>
