<template>
  <div class="notification-page">
    <notification-tab />

    <div class="config-card">
      <div class="card-header">
        <div class="header-icon" style="background: linear-gradient(135deg, #10b981, #059669)">
          <el-icon><Link /></el-icon>
        </div>
        <div class="header-text">
          <h2>{{ t('system.webhook') }}</h2>
          <p>配置 Webhook 以发送 HTTP 通知</p>
        </div>
      </div>

      <div class="info-banner">
        <el-icon><InfoFilled /></el-icon>
        <span>{{ t('system.webhookTip') }}</span>
      </div>

      <el-form
        ref="formRef"
        :model="form"
        :rules="formRules"
        label-position="top"
        class="config-form"
      >
        <el-form-item :label="t('system.template')" prop="template">
          <el-input
            type="textarea"
            :rows="8"
            :placeholder="webhookPlaceholder"
            v-model.trim="form.template"
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

    <div class="urls-card">
      <div class="card-header">
        <div class="header-icon" style="background: linear-gradient(135deg, #3b82f6, #1d4ed8)">
          <el-icon><Position /></el-icon>
        </div>
        <div class="header-text">
          <h2>{{ t('system.webhookUrls') }}</h2>
          <p>管理 Webhook 接收地址</p>
        </div>
        <el-button type="primary" @click="dialogVisible = true" class="add-btn">
          <el-icon><Plus /></el-icon>
          {{ t('system.addWebhookUrl') }}
        </el-button>
      </div>

      <div class="urls-list" v-if="webhookUrls.length > 0">
        <div class="url-item" v-for="item in webhookUrls" :key="item.id">
          <div class="url-icon">
            <el-icon><Link /></el-icon>
          </div>
          <div class="url-info">
            <span class="url-name">{{ item.name }}</span>
            <span class="url-address">{{ item.url }}</span>
          </div>
          <button class="delete-btn" @click="deleteUrl(item)">
            <el-icon><Close /></el-icon>
          </button>
        </div>
      </div>

      <div class="empty-urls" v-else>
        <el-icon :size="48"><Link /></el-icon>
        <p>{{ t('system.noWebhookUrls') }}</p>
        <el-button type="primary" @click="dialogVisible = true">
          <el-icon><Plus /></el-icon>
          {{ t('system.addWebhookUrl') }}
        </el-button>
      </div>
    </div>

    <el-dialog
      v-model="dialogVisible"
      :title="t('system.addWebhookUrl')"
      :width="isMobile ? '90%' : '400px'"
      class="url-dialog"
    >
      <el-form :model="{ name, url }" label-position="top">
        <el-form-item :label="t('system.webhookName')">
          <el-input v-model.trim="name" :placeholder="t('system.webhookName')">
            <template #prefix>
              <el-icon><CollectionTag /></el-icon>
            </template>
          </el-input>
        </el-form-item>
        <el-form-item label="URL">
          <el-input v-model.trim="url" placeholder="https://example.com/webhook">
            <template #prefix>
              <el-icon><Link /></el-icon>
            </template>
          </el-input>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">{{ t('common.cancel') }}</el-button>
        <el-button type="primary" @click="saveUrl">{{ t('common.confirm') }}</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { ElMessage } from 'element-plus'
import {
  Link,
  InfoFilled,
  Check,
  Position,
  Plus,
  Close,
  CollectionTag
} from '@element-plus/icons-vue'
import notificationTab from './tab.vue'
import notificationService from '../../../api/notification'

const { t } = useI18n()

const formRef = ref()
const saving = ref(false)
const dialogVisible = ref(false)
const isMobile = ref(false)

const form = ref({ template: '' })
const webhookUrls = ref([])
const name = ref('')
const url = ref('')

const checkMobile = () => {
  isMobile.value = window.innerWidth <= 768
}

const webhookPlaceholder = `{"task_id": "{{.TaskId}}", "task_name": "{{.TaskName}}", "status": "{{.Status}}", "result": "{{.Result}}", "remark": "{{.Remark}}"}`

const formRules = computed(() => ({
  template: [{ required: true, message: t('system.pleaseEnterTemplate'), trigger: 'blur' }]
}))

const fetchData = () => {
  notificationService.webhook(data => {
    form.value.template = data.template || ''
    webhookUrls.value = data.webhook_urls || []
  })
}

const submit = async () => {
  if (!formRef.value) return
  await formRef.value.validate(valid => {
    if (!valid) return
    saving.value = true
    notificationService.updateWebHook(form.value, () => {
      ElMessage.success(t('message.updateSuccess'))
      saving.value = false
      fetchData()
    })
  })
}

const saveUrl = () => {
  if (!name.value || !url.value) {
    ElMessage.error(t('system.incompleteParameters'))
    return
  }
  notificationService.createWebhookUrl(
    {
      name: name.value,
      url: url.value
    },
    () => {
      dialogVisible.value = false
      name.value = ''
      url.value = ''
      fetchData()
    }
  )
}

const deleteUrl = item => {
  notificationService.removeWebhookUrl(item.id, () => {
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
.urls-card {
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

.info-banner {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 14px 24px;
  background: linear-gradient(135deg, #ecfdf5 0%, #d1fae5 100%);
  color: #047857;
  font-size: 13px;
}

.info-banner .el-icon {
  font-size: 18px;
}

.config-form {
  padding: 24px;
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

.urls-list {
  padding: 16px 24px 24px;
  display: grid;
  gap: 12px;
}

.url-item {
  display: flex;
  align-items: center;
  gap: 14px;
  padding: 16px;
  background: #f8fafc;
  border-radius: 14px;
  transition: all 0.2s;
}

.url-item:hover {
  background: #f1f5f9;
}

.url-icon {
  width: 40px;
  height: 40px;
  border-radius: 10px;
  background: linear-gradient(135deg, #10b981, #059669);
  color: white;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.url-info {
  flex: 1;
  min-width: 0;
}

.url-name {
  display: block;
  font-weight: 500;
  color: #1e293b;
  margin-bottom: 4px;
}

.url-address {
  display: block;
  font-size: 12px;
  color: #64748b;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  font-family: 'SF Mono', Monaco, monospace;
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

.empty-urls {
  padding: 48px 24px;
  text-align: center;
  color: #94a3b8;
}

.empty-urls .el-icon {
  margin-bottom: 12px;
}

.empty-urls p {
  margin: 0 0 16px 0;
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

  .info-banner {
    padding: 12px 16px;
    font-size: 12px;
  }

  .config-form {
    padding: 16px;
  }

  .urls-list {
    padding: 12px 16px 16px;
  }

  .url-item {
    padding: 12px;
  }
}
</style>
