<template>
  <div class="notification-tabs">
    <div class="tabs-header">
      <button
        v-for="tab in tabs"
        :key="tab.name"
        class="tab-btn"
        :class="{ active: activeName === tab.name }"
        @click="switchTab(tab.name)"
      >
        <span class="tab-icon" :style="{ background: tab.color }">
          <component :is="tab.icon" />
        </span>
        <span class="tab-label">{{ tab.label }}</span>
      </button>
    </div>

    <div class="template-vars-card">
      <div class="vars-header" @click="varsExpanded = !varsExpanded">
        <div class="vars-title">
          <el-icon><InfoFilled /></el-icon>
          <span>{{ t('system.templateVariables') }}</span>
        </div>
        <el-icon class="vars-arrow" :class="{ expanded: varsExpanded }">
          <ArrowDown />
        </el-icon>
      </div>
      <transition name="slide">
        <div class="vars-body" v-show="varsExpanded">
          <div class="vars-grid">
            <div class="var-item" v-for="v in variables" :key="v.name">
              <code class="var-code">{{ v.name }}</code>
              <span class="var-desc">{{ v.desc }}</span>
            </div>
          </div>
        </div>
      </transition>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, watch, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import {
  Message,
  Bell,
  Link,
  ChatDotRound,
  Iphone,
  InfoFilled,
  ArrowDown
} from '@element-plus/icons-vue'

const { t } = useI18n()
const route = useRoute()
const router = useRouter()

const activeName = ref('email')
const varsExpanded = ref(false)

const tabs = [
  {
    name: 'email',
    label: 'Email',
    icon: Message,
    color: 'linear-gradient(135deg, #3b82f6, #1d4ed8)'
  },
  {
    name: 'slack',
    label: 'Slack',
    icon: ChatDotRound,
    color: 'linear-gradient(135deg, #4a154b, #611f69)'
  },
  {
    name: 'webhook',
    label: 'Webhook',
    icon: Link,
    color: 'linear-gradient(135deg, #10b981, #059669)'
  },
  {
    name: 'serverchan3',
    label: 'Server酱³',
    icon: Bell,
    color: 'linear-gradient(135deg, #f59e0b, #d97706)'
  },
  { name: 'bark', label: 'Bark', icon: Iphone, color: 'linear-gradient(135deg, #ec4899, #db2777)' }
]

const variables = computed(() => [
  { name: '{{.TaskId}}', desc: t('system.taskIdVar') },
  { name: '{{.TaskName}}', desc: t('system.taskNameVar') },
  { name: '{{.Status}}', desc: t('system.statusVar') },
  { name: '{{.StatusZh}}', desc: t('system.statusVar') + '（中文）' },
  { name: '{{.IsSuccess}}', desc: t('system.statusVar') + '（true/false）' },
  { name: '{{.Host}}', desc: '节点信息' },
  { name: '{{.Result}}', desc: t('system.resultVar') },
  { name: '{{.ResultSummary}}', desc: '输出摘要' },
  { name: '{{.ResultBody}}', desc: '主体输出' },
  { name: '{{.Remark}}', desc: t('task.remark') }
])

const switchTab = name => {
  activeName.value = name
  router.push(`/system/notification/${name}`)
}

watch(
  () => route.path,
  path => {
    const segments = path.split('/')
    if (segments.length === 4 && segments[2] === 'notification') {
      activeName.value = segments[3]
    }
  },
  { immediate: true }
)
</script>

<style scoped>
.notification-tabs {
  margin-bottom: 24px;
}

.tabs-header {
  display: flex;
  gap: 8px;
  padding: 6px;
  background: #f1f5f9;
  border-radius: 16px;
  overflow-x: auto;
  -webkit-overflow-scrolling: touch;
}

.tab-btn {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 12px 20px;
  border: none;
  background: transparent;
  border-radius: 12px;
  cursor: pointer;
  transition: all 0.3s;
  white-space: nowrap;
  flex: 1;
  justify-content: center;
}

.tab-btn:hover {
  background: rgba(255, 255, 255, 0.6);
}

.tab-btn.active {
  background: white;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.08);
}

.tab-icon {
  width: 32px;
  height: 32px;
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  font-size: 16px;
}

.tab-label {
  font-weight: 500;
  color: #334155;
  font-size: 14px;
}

.template-vars-card {
  background: white;
  border-radius: 16px;
  margin-top: 16px;
  overflow: hidden;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.04);
}

.vars-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px 20px;
  cursor: pointer;
  transition: background 0.2s;
}

.vars-header:hover {
  background: #f8fafc;
}

.vars-title {
  display: flex;
  align-items: center;
  gap: 10px;
  font-weight: 600;
  color: #334155;
}

.vars-title .el-icon {
  color: #3b82f6;
}

.vars-arrow {
  transition: transform 0.3s;
  color: #94a3b8;
}

.vars-arrow.expanded {
  transform: rotate(180deg);
}

.vars-body {
  padding: 0 20px 20px;
}

.vars-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
  gap: 12px;
}

.var-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 10px 14px;
  background: #f8fafc;
  border-radius: 10px;
}

.var-code {
  background: #1e293b;
  color: #a5f3fc;
  padding: 4px 10px;
  border-radius: 6px;
  font-family: 'SF Mono', Monaco, monospace;
  font-size: 12px;
  white-space: nowrap;
}

.var-desc {
  font-size: 13px;
  color: #64748b;
}

.slide-enter-active,
.slide-leave-active {
  transition: all 0.3s ease;
}

.slide-enter-from,
.slide-leave-to {
  opacity: 0;
  transform: translateY(-10px);
}

@media screen and (max-width: 768px) {
  .tabs-header {
    padding: 4px;
    gap: 4px;
  }

  .tab-btn {
    padding: 10px 12px;
    flex-direction: column;
    gap: 6px;
  }

  .tab-icon {
    width: 28px;
    height: 28px;
    font-size: 14px;
  }

  .tab-label {
    font-size: 11px;
  }

  .vars-grid {
    grid-template-columns: 1fr;
  }

  .var-item {
    flex-direction: column;
    align-items: flex-start;
    gap: 6px;
  }
}
</style>
