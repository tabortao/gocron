<template>
  <div class="notification-page">
    <notification-tab />

    <div class="config-card">
      <div class="card-header">
        <div class="header-icon" style="background: linear-gradient(135deg, #4a154b, #611f69)">
          <el-icon><ChatDotRound /></el-icon>
        </div>
        <div class="header-text">
          <h2>Slack</h2>
          <p>配置 Slack Webhook 以发送通知</p>
        </div>
      </div>

      <el-form
        ref="formRef"
        :model="form"
        :rules="formRules"
        label-position="top"
        class="config-form"
      >
        <el-form-item :label="t('system.slackUrl')" prop="url">
          <el-input v-model="form.url" placeholder="https://hooks.slack.com/services/...">
            <template #prefix>
              <el-icon><Link /></el-icon>
            </template>
          </el-input>
        </el-form-item>

        <el-form-item :label="t('system.template')" prop="template">
          <el-input
            type="textarea"
            :rows="8"
            :placeholder="slackPlaceholder"
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

    <div class="channels-card">
      <div class="card-header">
        <div class="header-icon" style="background: linear-gradient(135deg, #e11d48, #be123c)">
          <el-icon><Bell /></el-icon>
        </div>
        <div class="header-text">
          <h2>{{ t('system.channels') }}</h2>
          <p>管理 Slack 频道</p>
        </div>
        <el-button type="primary" @click="dialogVisible = true" class="add-btn">
          <el-icon><Plus /></el-icon>
          {{ t('system.addChannel') }}
        </el-button>
      </div>

      <div class="channels-list" v-if="channels.length > 0">
        <div class="channel-item" v-for="item in channels" :key="item.id">
          <div class="channel-icon">
            <el-icon><ChatDotRound /></el-icon>
          </div>
          <div class="channel-info">
            <span class="channel-name"># {{ item.name }}</span>
          </div>
          <button class="delete-btn" @click="deleteChannel(item)">
            <el-icon><Close /></el-icon>
          </button>
        </div>
      </div>

      <div class="empty-channels" v-else>
        <el-icon :size="48"><ChatDotRound /></el-icon>
        <p>{{ t('system.noChannels') }}</p>
        <el-button type="primary" @click="dialogVisible = true">
          <el-icon><Plus /></el-icon>
          {{ t('system.addChannel') }}
        </el-button>
      </div>
    </div>

    <el-dialog
      v-model="dialogVisible"
      :title="t('system.addChannel')"
      :width="isMobile ? '90%' : '400px'"
    >
      <el-form label-position="top">
        <el-form-item :label="t('system.channelName')">
          <el-input v-model.trim="channel" placeholder="general">
            <template #prefix>
              <span>#</span>
            </template>
          </el-input>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">{{ t('common.cancel') }}</el-button>
        <el-button type="primary" @click="saveChannel">{{ t('common.confirm') }}</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { ElMessage } from 'element-plus'
import { ChatDotRound, Link, Check, Bell, Plus, Close } from '@element-plus/icons-vue'
import notificationTab from './tab.vue'
import notificationService from '../../../api/notification'

const { t } = useI18n()

const formRef = ref()
const saving = ref(false)
const dialogVisible = ref(false)
const isMobile = ref(false)

const form = ref({ url: '', template: '' })
const channels = ref([])
const channel = ref('')

const checkMobile = () => {
  isMobile.value = window.innerWidth <= 768
}

const slackPlaceholder = computed(
  () =>
    `${t('system.taskIdVar')}: {{.TaskId}}
${t('system.taskNameVar')}: {{.TaskName}}
${t('system.statusVar')}: {{.Status}}
${t('system.resultVar')}: {{.Result}}
${t('task.remark')}: {{.Remark}}`
)

const formRules = computed(() => ({
  url: [{ type: 'url', required: true, message: t('system.pleaseEnterValidUrl'), trigger: 'blur' }],
  template: [{ required: true, message: t('system.pleaseEnterTemplate'), trigger: 'blur' }]
}))

const fetchData = () => {
  notificationService.slack(data => {
    form.value.url = data.url || ''
    form.value.template = data.template || ''
    channels.value = data.channels || []
  })
}

const submit = async () => {
  if (!formRef.value) return
  await formRef.value.validate(valid => {
    if (!valid) return
    saving.value = true
    notificationService.updateSlack(form.value, () => {
      ElMessage.success(t('message.updateSuccess'))
      saving.value = false
      fetchData()
    })
  })
}

const saveChannel = () => {
  if (!channel.value) {
    ElMessage.error(t('system.pleaseEnterChannelName'))
    return
  }
  notificationService.createSlackChannel(channel.value, () => {
    dialogVisible.value = false
    channel.value = ''
    fetchData()
  })
}

const deleteChannel = item => {
  notificationService.removeSlackChannel(item.id, () => {
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
.channels-card {
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

.channels-list {
  padding: 16px 24px 24px;
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
  gap: 12px;
}

.channel-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 14px 16px;
  background: #f8fafc;
  border-radius: 14px;
  transition: all 0.2s;
}

.channel-item:hover {
  background: #f1f5f9;
}

.channel-icon {
  width: 36px;
  height: 36px;
  border-radius: 10px;
  background: linear-gradient(135deg, #4a154b, #611f69);
  color: white;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.channel-info {
  flex: 1;
  min-width: 0;
}

.channel-name {
  font-weight: 500;
  color: #1e293b;
}

.delete-btn {
  width: 28px;
  height: 28px;
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

.empty-channels {
  padding: 48px 24px;
  text-align: center;
  color: #94a3b8;
}

.empty-channels .el-icon {
  margin-bottom: 12px;
}

.empty-channels p {
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

  .config-form {
    padding: 16px;
  }

  .channels-list {
    padding: 12px 16px 16px;
    grid-template-columns: 1fr;
  }
}
</style>
