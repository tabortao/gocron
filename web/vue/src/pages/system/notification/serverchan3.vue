<template>
  <div class="notification-page">
    <notification-tab />

    <div class="config-card">
      <div class="card-header">
        <div class="header-icon" style="background: linear-gradient(135deg, #f59e0b, #d97706)">
          <el-icon><Bell /></el-icon>
        </div>
        <div class="header-text">
          <h2>Server酱³</h2>
          <p>微信推送通知服务</p>
        </div>
      </div>

      <div class="info-banner">
        <el-icon><InfoFilled /></el-icon>
        <span>Server酱³ 是一个微信消息推送服务，支持多种推送渠道</span>
      </div>

      <el-form
        ref="formRef"
        :model="form"
        :rules="formRules"
        label-position="top"
        class="config-form"
      >
        <el-form-item :label="isZh ? '标题模板' : 'Title template'" prop="title_template">
          <el-input
            type="textarea"
            :rows="3"
            :placeholder="titlePlaceholder"
            v-model.trim="form.title_template"
          />
        </el-form-item>

        <el-form-item :label="isZh ? '内容模板' : 'Content template'" prop="body_template">
          <el-input
            type="textarea"
            :rows="8"
            :placeholder="bodyPlaceholder"
            v-model.trim="form.body_template"
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
        <div class="header-icon" style="background: linear-gradient(135deg, #10b981, #059669)">
          <el-icon><Position /></el-icon>
        </div>
        <div class="header-text">
          <h2>{{ isZh ? '推送渠道' : 'Push Channels' }}</h2>
          <p>{{ isZh ? '管理 Server酱³ 推送渠道' : 'Manage ServerChan³ channels' }}</p>
        </div>
        <el-button type="primary" @click="dialogVisible = true" class="add-btn">
          <el-icon><Plus /></el-icon>
          {{ isZh ? '新增渠道' : 'Add Channel' }}
        </el-button>
      </div>

      <div class="channels-list" v-if="channels.length > 0">
        <div class="channel-item" v-for="item in channels" :key="item.id">
          <div class="channel-icon">
            <el-icon><Bell /></el-icon>
          </div>
          <div class="channel-info">
            <span class="channel-name">{{ item.name }}</span>
            <span class="channel-url">{{ item.url }}</span>
          </div>
          <button class="delete-btn" @click="deleteChannel(item)">
            <el-icon><Close /></el-icon>
          </button>
        </div>
      </div>

      <div class="empty-channels" v-else>
        <el-icon :size="48"><Bell /></el-icon>
        <p>{{ isZh ? '暂无推送渠道' : 'No channels' }}</p>
        <el-button type="primary" @click="dialogVisible = true">
          <el-icon><Plus /></el-icon>
          {{ isZh ? '新增渠道' : 'Add Channel' }}
        </el-button>
      </div>
    </div>

    <el-dialog
      v-model="dialogVisible"
      :title="isZh ? '新增推送渠道' : 'Add Channel'"
      :width="isMobile ? '90%' : '400px'"
    >
      <el-form label-position="top">
        <el-form-item :label="isZh ? '渠道名称' : 'Channel Name'">
          <el-input v-model.trim="name" :placeholder="isZh ? '我的微信' : 'My WeChat'">
            <template #prefix>
              <el-icon><CollectionTag /></el-icon>
            </template>
          </el-input>
        </el-form-item>
        <el-form-item label="URL">
          <el-input v-model.trim="url" placeholder="https://sctapi.ftqq.com/your_sendkey.send">
            <template #prefix>
              <el-icon><Link /></el-icon>
            </template>
          </el-input>
          <div class="url-tip">
            {{
              isZh
                ? 'URL 格式：https://sctapi.ftqq.com/{sendkey}.send'
                : 'URL format: https://sctapi.ftqq.com/{sendkey}.send'
            }}
          </div>
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
import {
  Bell,
  InfoFilled,
  Check,
  Position,
  Plus,
  Close,
  Link,
  CollectionTag
} from '@element-plus/icons-vue'
import notificationTab from './tab.vue'
import notificationService from '../../../api/notification'

const { t, locale } = useI18n()

const isZh = computed(() => locale.value === 'zh-CN')
const isMobile = ref(false)
const formRef = ref()
const saving = ref(false)
const dialogVisible = ref(false)

const form = ref({ title_template: '', body_template: '' })
const channels = ref([])
const name = ref('')
const url = ref('')

const checkMobile = () => {
  isMobile.value = window.innerWidth <= 768
}

const titlePlaceholder = '{{.TaskName}} - {{.StatusZh}}'
const bodyPlaceholder = `**任务名称**：{{.TaskName}}
**任务ID**：{{.TaskId}}
**执行状态**：{{.StatusZh}}
{{ if .Host }}**执行节点**：{{.Host}}
{{ end }}**执行结果**：
{{.ResultSummary}}

{{ if .Remark }}**备注**：{{.Remark}}{{ end }}`

const formRules = computed(() => ({
  title_template: [
    { required: true, message: isZh.value ? '请输入标题模板' : 'Required', trigger: 'blur' }
  ],
  body_template: [
    { required: true, message: isZh.value ? '请输入内容模板' : 'Required', trigger: 'blur' }
  ]
}))

const fetchData = () => {
  notificationService.serverchan3(data => {
    form.value.title_template = data.title_template || ''
    form.value.body_template = data.body_template || ''
    channels.value = data.channels || []
  })
}

const submit = async () => {
  if (!formRef.value) return
  await formRef.value.validate(valid => {
    if (!valid) return
    saving.value = true
    notificationService.updateServerchan3(form.value, () => {
      ElMessage.success(t('message.updateSuccess'))
      saving.value = false
      fetchData()
    })
  })
}

const saveChannel = () => {
  if (!name.value || !url.value) {
    ElMessage.error(t('system.incompleteParameters'))
    return
  }
  notificationService.createServerchan3Channel({ name: name.value, url: url.value }, () => {
    dialogVisible.value = false
    name.value = ''
    url.value = ''
    fetchData()
  })
}

const deleteChannel = item => {
  notificationService.removeServerchan3Channel(item.id, () => fetchData())
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

.info-banner {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 14px 24px;
  background: linear-gradient(135deg, #fffbeb 0%, #fef3c7 100%);
  color: #b45309;
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

.channels-list {
  padding: 16px 24px 24px;
  display: grid;
  gap: 12px;
}

.channel-item {
  display: flex;
  align-items: center;
  gap: 14px;
  padding: 16px;
  background: #f8fafc;
  border-radius: 14px;
  transition: all 0.2s;
}

.channel-item:hover {
  background: #f1f5f9;
}

.channel-icon {
  width: 40px;
  height: 40px;
  border-radius: 10px;
  background: linear-gradient(135deg, #f59e0b, #d97706);
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
  display: block;
  font-weight: 500;
  color: #1e293b;
  margin-bottom: 4px;
}

.channel-url {
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

.url-tip {
  margin-top: 8px;
  font-size: 12px;
  color: #64748b;
  line-height: 1.5;
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

  .channels-list {
    padding: 12px 16px 16px;
  }

  .channel-item {
    padding: 12px;
  }
}
</style>
