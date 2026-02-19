<template>
  <div class="page-container">
    <div class="page-header">
      <div class="page-header__left">
        <h1 class="page-title">{{ t('taskLog.list') }}</h1>
        <p class="page-desc" v-if="logTotal > 0">
          {{ isZh ? '共' : 'Total' }} <span class="count">{{ logTotal }}</span>
          {{ isZh ? '条日志' : 'logs' }}
        </p>
      </div>
      <div class="page-header__right">
        <el-button type="danger" v-if="isAdmin" @click="clearLog" class="danger-btn">
          <el-icon><Delete /></el-icon>
          <span class="btn-text">{{ t('message.clearLog') }}</span>
        </el-button>
        <el-button @click="refresh" class="ghost-btn">
          <el-icon :class="{ 'is-loading': loading }"><Refresh /></el-icon>
          <span class="btn-text">{{ t('common.refresh') }}</span>
        </el-button>
      </div>
    </div>

    <div class="filter-card">
      <div class="filter-header" @click="filterExpanded = !filterExpanded">
        <div class="filter-header__left">
          <el-icon><Search /></el-icon>
          <span>{{ t('common.search') }}</span>
        </div>
        <el-icon class="filter-arrow" :class="{ expanded: filterExpanded }">
          <ArrowDown />
        </el-icon>
      </div>
      <el-collapse-transition>
        <div class="filter-body" v-show="filterExpanded">
          <div class="filter-grid">
            <div class="filter-item">
              <label>{{ t('task.id') }}</label>
              <el-input v-model.trim="searchParams.task_id" :placeholder="t('task.id')" clearable />
            </div>
            <div class="filter-item hidden-mobile">
              <label>{{ t('task.protocol') }}</label>
              <el-select v-model.trim="searchParams.protocol" clearable>
                <el-option :label="t('select')" value="" />
                <el-option
                  v-for="item in protocolList"
                  :key="item.value"
                  :label="item.label"
                  :value="item.value"
                />
              </el-select>
            </div>
            <div class="filter-item hidden-mobile">
              <label>{{ t('common.status') }}</label>
              <el-select v-model.trim="searchParams.status" clearable>
                <el-option :label="t('select')" value="" />
                <el-option
                  v-for="item in statusList"
                  :key="item.value"
                  :label="item.label"
                  :value="item.value"
                />
              </el-select>
            </div>
          </div>
          <div class="filter-actions">
            <el-button type="primary" @click="search()">
              <el-icon><Search /></el-icon>
              {{ t('common.search') }}
            </el-button>
            <el-button @click="resetSearch">{{ t('common.reset') }}</el-button>
          </div>
        </div>
      </el-collapse-transition>
    </div>

    <div class="content-card">
      <div class="table-wrapper hidden-mobile">
        <el-table :data="logs" class="data-table" v-loading="loading">
          <el-table-column type="expand">
            <template #default="{ row }">
              <div class="expand-content">
                <div class="expand-grid">
                  <div class="expand-item">
                    <span class="expand-label">{{ t('message.retryCount') }}</span>
                    <span class="expand-value">{{ row.retry_times || 0 }}</span>
                  </div>
                  <div class="expand-item">
                    <span class="expand-label">{{ t('task.cronExpression') }}</span>
                    <code class="expand-code">{{ row.spec }}</code>
                  </div>
                </div>
                <div class="expand-full" v-if="row.command">
                  <span class="expand-label">{{ t('task.command') }}</span>
                  <code class="command-code">{{ row.command }}</code>
                </div>
              </div>
            </template>
          </el-table-column>
          <el-table-column prop="id" label="ID" width="70" align="center">
            <template #default="{ row }">
              <span class="id-badge">{{ row.id }}</span>
            </template>
          </el-table-column>
          <el-table-column prop="task_id" :label="t('task.id')" width="80" align="center">
            <template #default="{ row }">
              <span class="task-id-link" @click="goToTask(row.task_id)">{{ row.task_id }}</span>
            </template>
          </el-table-column>
          <el-table-column prop="name" :label="t('task.name')" min-width="140">
            <template #default="{ row }">
              <div class="task-name">
                <span class="name-text">{{ row.name }}</span>
              </div>
            </template>
          </el-table-column>
          <el-table-column :label="t('task.protocol')" width="100" align="center">
            <template #default="{ row }">
              <el-tag :type="getProtocolType(row)" size="small" effect="plain">
                {{ formatProtocol(row) }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column :label="t('task.taskNode')" width="140">
            <template #default="{ row }">
              <span class="host-text" v-html="row.hostname"></span>
            </template>
          </el-table-column>
          <el-table-column :label="t('taskLog.duration')" width="200">
            <template #default="{ row }">
              <div class="time-info">
                <div class="time-row">
                  <el-icon><Timer /></el-icon>
                  <span
                    >{{ row.total_time > 0 ? row.total_time : 1 }}{{ t('message.seconds') }}</span
                  >
                </div>
                <div class="time-row sub">
                  <span>{{ $filters.formatTime(row.start_time) }}</span>
                </div>
              </div>
            </template>
          </el-table-column>
          <el-table-column :label="t('common.status')" width="100" align="center">
            <template #default="{ row }">
              <el-tag :type="getStatusType(row.status)" size="small" effect="plain">
                {{ getStatusText(row.status) }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column :label="t('taskLog.result')" width="140" align="center">
            <template #default="{ row }">
              <div class="action-btns">
                <el-button
                  :type="getResultBtnType(row.status)"
                  size="small"
                  text
                  @click="showTaskResult(row)"
                >
                  <el-icon><View /></el-icon>
                  {{ t('taskLog.viewOutput') }}
                </el-button>
                <el-button
                  v-if="isAdmin && row.status === 1 && row.protocol === 2"
                  type="danger"
                  size="small"
                  text
                  @click="stopTask(row)"
                >
                  <el-icon><VideoPause /></el-icon>
                </el-button>
              </div>
            </template>
          </el-table-column>
        </el-table>
      </div>

      <div class="card-list visible-mobile">
        <div v-for="log in logs" :key="log.id" class="log-card" :class="getStatusClass(log.status)">
          <div class="log-card__header">
            <span class="log-card__id">#{{ log.id }}</span>
            <el-tag :type="getStatusType(log.status)" size="small" effect="plain">
              {{ getStatusText(log.status) }}
            </el-tag>
          </div>
          <div class="log-card__body">
            <h3 class="log-card__title">{{ log.name }}</h3>
            <div class="log-card__meta">
              <span class="meta-item">
                <el-icon><CollectionTag /></el-icon>
                <span>{{ t('task.id') }}: {{ log.task_id }}</span>
              </span>
              <el-tag :type="getProtocolType(log)" size="small" effect="plain">
                {{ formatProtocol(log) }}
              </el-tag>
            </div>
            <div class="log-card__meta">
              <span class="meta-item">
                <el-icon><Timer /></el-icon>
                <span>{{ log.total_time > 0 ? log.total_time : 1 }}{{ t('message.seconds') }}</span>
              </span>
              <span class="meta-item">
                <el-icon><Clock /></el-icon>
                <span>{{ $filters.formatTime(log.start_time) }}</span>
              </span>
            </div>
            <div class="log-card__host" v-if="log.hostname" v-html="log.hostname"></div>
          </div>
          <div class="log-card__footer">
            <el-button
              :type="getResultBtnType(log.status)"
              size="small"
              text
              @click="showTaskResult(log)"
            >
              <el-icon><View /></el-icon>
              {{ t('taskLog.viewOutput') }}
            </el-button>
            <el-button
              v-if="isAdmin && log.status === 1 && log.protocol === 2"
              type="danger"
              size="small"
              text
              @click="stopTask(log)"
            >
              <el-icon><VideoPause /></el-icon>
              {{ t('message.stopTask') }}
            </el-button>
          </div>
        </div>
        <el-empty v-if="!loading && logs.length === 0" :description="t('message.noData')" />
      </div>

      <div class="pagination-wrapper">
        <el-pagination
          background
          :layout="isMobile ? 'prev, pager, next' : 'total, sizes, prev, pager, next'"
          :total="logTotal"
          v-model:current-page="searchParams.page"
          v-model:page-size="searchParams.page_size"
          @size-change="changePageSize"
          @current-change="changePage"
          :page-sizes="[10, 20, 50, 100]"
          small
        />
      </div>
    </div>

    <el-dialog
      :title="t('message.taskExecutionResult')"
      v-model="dialogVisible"
      width="60%"
      class="result-dialog"
      :fullscreen="isMobile"
    >
      <div class="result-content">
        <div class="result-section" v-if="currentTaskResult.hostname">
          <div class="result-label">{{ t('taskLog.host') }}</div>
          <div class="result-value host-value" v-html="currentTaskResult.hostname"></div>
        </div>
        <div class="result-section">
          <div class="result-label">{{ t('task.command') }}</div>
          <pre class="result-code command">{{ currentTaskResult.command }}</pre>
        </div>
        <div class="result-section">
          <div class="result-label">{{ t('taskLog.output') }}</div>
          <pre class="result-code output" ref="resultPre">{{ currentTaskResult.result }}</pre>
        </div>
      </div>
    </el-dialog>
  </div>
</template>

<script setup>
import {
  ref,
  computed,
  watch,
  onMounted,
  onActivated,
  onDeactivated,
  onUnmounted,
  nextTick
} from 'vue'
import { useI18n } from 'vue-i18n'
import { useRouter, useRoute } from 'vue-router'
import { useUserStore } from '../../stores/user'
import { ElMessageBox, ElMessage } from 'element-plus'
import {
  Delete,
  Refresh,
  Search,
  ArrowDown,
  View,
  VideoPause,
  Timer,
  Clock,
  CollectionTag
} from '@element-plus/icons-vue'
import taskLogService from '../../api/taskLog'

const { t, locale } = useI18n()
const router = useRouter()
const route = useRoute()
const userStore = useUserStore()

const isZh = computed(() => locale.value === 'zh-CN')
const isAdmin = userStore.isAdmin

const isMobile = ref(false)
const loading = ref(false)
const logs = ref([])
const logTotal = ref(0)
const filterExpanded = ref(true)
const dialogVisible = ref(false)
const resultPre = ref(null)

const searchParams = ref({
  page_size: 20,
  page: 1,
  task_id: '',
  protocol: '',
  status: ''
})

const currentTaskResult = ref({
  hostname: '',
  command: '',
  result: ''
})

const currentLogId = ref(0)
const currentLogStatus = ref(0)
const outputRefreshTimer = ref(null)
const outputRefreshInFlight = ref(false)
const autoRefreshTimer = ref(null)
const autoRefreshInFlight = ref(false)

const protocolList = [
  { value: '1', label: 'HTTP' },
  { value: '2', label: 'Shell' }
]

const statusList = computed(() => [
  { value: '2', label: t('taskLog.success') },
  { value: '0', label: t('taskLog.failed') },
  { value: '1', label: t('message.running') },
  { value: '3', label: t('message.cancelled') }
])

const checkMobile = () => {
  isMobile.value = window.innerWidth <= 768
}

const formatProtocol = row => {
  return row.protocol === 1 ? 'HTTP' : 'Shell'
}

const getProtocolType = row => {
  return row.protocol === 2 ? 'warning' : 'primary'
}

const getStatusType = status => {
  const types = {
    0: 'danger',
    1: 'primary',
    2: 'success',
    3: 'info'
  }
  return types[status] || 'info'
}

const getStatusText = status => {
  const texts = {
    0: t('taskLog.failed'),
    1: t('message.running'),
    2: t('taskLog.success'),
    3: t('message.cancelled')
  }
  return texts[status] || '-'
}

const getStatusClass = status => {
  const classes = {
    0: 'status-failed',
    1: 'status-running',
    2: 'status-success',
    3: 'status-cancelled'
  }
  return classes[status] || ''
}

const getResultBtnType = status => {
  const types = {
    0: 'warning',
    1: 'primary',
    2: 'success',
    3: 'info'
  }
  return types[status] || 'primary'
}

const fetchData = (callback = null) => {
  loading.value = true
  taskLogService.list(searchParams.value, data => {
    logs.value = data.data
    logTotal.value = data.total
    loading.value = false
    ensureAutoRefresh()
    if (callback) callback()
  })
}

const search = () => {
  searchParams.value.page = 1
  fetchData()
}

const resetSearch = () => {
  searchParams.value = {
    page_size: 20,
    page: 1,
    task_id: '',
    protocol: '',
    status: ''
  }
  fetchData()
}

const refresh = () => {
  fetchData(() => {
    ElMessage.success(t('message.refreshSuccess'))
  })
}

const changePage = page => {
  searchParams.value.page = page
  fetchData()
}

const changePageSize = size => {
  searchParams.value.page_size = size
  searchParams.value.page = 1
  fetchData()
}

const goToTask = taskId => {
  router.push(`/task/log?task_id=${taskId}`)
}

const clearLog = () => {
  ElMessageBox.confirm(t('message.confirmClearLog'), t('common.tip'), {
    confirmButtonText: t('common.confirm'),
    cancelButtonText: t('common.cancel'),
    type: 'warning',
    center: true
  })
    .then(() => {
      taskLogService.clear(() => {
        searchParams.value.page = 1
        fetchData()
      })
    })
    .catch(() => {})
}

const stopTask = item => {
  taskLogService.stop(item.id, item.task_id, () => {
    fetchData()
  })
}

const showTaskResult = item => {
  dialogVisible.value = true
  currentLogId.value = item.id || 0
  currentLogStatus.value = item.status
  let cleanedCommand = item.command
  if (cleanedCommand) {
    cleanedCommand = cleanedCommand
      .replace(/&quot;/g, '"')
      .replace(/&apos;/g, "'")
      .replace(/&#39;/g, "'")
      .replace(/&lt;/g, '<')
      .replace(/&gt;/g, '>')
      .replace(/&amp;/g, '&')
  }
  currentTaskResult.value = {
    hostname: item.hostname || '',
    command: cleanedCommand,
    result: item.result
  }
  if (item.status === 1) {
    fetchLiveOutput()
    startOutputRefresh()
  } else {
    stopOutputRefresh()
  }
  nextTick(() => {
    scrollOutputToBottom()
  })
}

const fetchLiveOutput = () => {
  if (!currentLogId.value || outputRefreshInFlight.value) return
  outputRefreshInFlight.value = true
  taskLogService.output(currentLogId.value, data => {
    outputRefreshInFlight.value = false
    if (data && data.status !== undefined) {
      currentLogStatus.value = data.status
      if (currentLogStatus.value !== 1) {
        stopOutputRefresh()
      }
    }
    if (data && data.output !== undefined) {
      currentTaskResult.value.result = data.output
      nextTick(() => {
        scrollOutputToBottom()
      })
    }
  })
}

const startOutputRefresh = () => {
  if (outputRefreshTimer.value || currentLogStatus.value !== 1) return
  outputRefreshTimer.value = setInterval(() => {
    if (!dialogVisible.value || currentLogStatus.value !== 1) return
    fetchLiveOutput()
  }, 2000)
}

const stopOutputRefresh = () => {
  if (!outputRefreshTimer.value) return
  clearInterval(outputRefreshTimer.value)
  outputRefreshTimer.value = null
  outputRefreshInFlight.value = false
}

const scrollOutputToBottom = () => {
  const el = resultPre.value
  if (!el) return
  el.scrollTop = el.scrollHeight
}

const ensureAutoRefresh = () => {
  const hasRunning = Array.isArray(logs.value) && logs.value.some(item => item.status === 1)
  if (hasRunning) {
    startAutoRefresh()
  } else {
    stopAutoRefresh()
  }
}

const startAutoRefresh = () => {
  if (autoRefreshTimer.value) return
  autoRefreshTimer.value = setInterval(() => {
    if (autoRefreshInFlight.value) return
    autoRefreshInFlight.value = true
    fetchData(() => {
      autoRefreshInFlight.value = false
    })
  }, 3000)
}

const stopAutoRefresh = () => {
  if (!autoRefreshTimer.value) return
  clearInterval(autoRefreshTimer.value)
  autoRefreshTimer.value = null
  autoRefreshInFlight.value = false
}

const updateTaskIdFromRoute = () => {
  if (route.query.task_id) {
    searchParams.value.task_id = route.query.task_id
    searchParams.value.page = 1
  }
}

watch(
  () => route.query.task_id,
  newTaskId => {
    if (newTaskId !== undefined) {
      searchParams.value.task_id = newTaskId
      searchParams.value.page = 1
      search()
    }
  }
)

watch(dialogVisible, visible => {
  if (!visible) {
    stopOutputRefresh()
  }
})

onMounted(() => {
  checkMobile()
  window.addEventListener('resize', checkMobile)
  updateTaskIdFromRoute()
  fetchData()
})

onActivated(() => {
  updateTaskIdFromRoute()
  fetchData()
})

onDeactivated(() => {
  stopAutoRefresh()
  stopOutputRefresh()
})

onUnmounted(() => {
  window.removeEventListener('resize', checkMobile)
  stopAutoRefresh()
  stopOutputRefresh()
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

.page-desc .count {
  font-weight: 600;
  color: #374151;
}

.page-header__right {
  display: flex;
  gap: 10px;
}

.danger-btn {
  font-weight: 500;
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

.filter-card {
  background: #fff;
  border-radius: 12px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.05);
  margin-bottom: 16px;
  overflow: hidden;
}

.filter-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 14px 20px;
  cursor: pointer;
  border-bottom: 1px solid #f3f4f6;
}

.filter-header__left {
  display: flex;
  align-items: center;
  gap: 8px;
  font-weight: 500;
  color: #374151;
}

.filter-arrow {
  transition: transform 0.3s;
  color: #9ca3af;
}

.filter-arrow.expanded {
  transform: rotate(180deg);
}

.filter-body {
  padding: 20px;
}

.filter-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(180px, 1fr));
  gap: 16px;
  margin-bottom: 16px;
}

.filter-item label {
  display: block;
  font-size: 13px;
  font-weight: 500;
  color: #6b7280;
  margin-bottom: 6px;
}

.filter-item :deep(.el-input),
.filter-item :deep(.el-select) {
  width: 100%;
}

.filter-actions {
  display: flex;
  gap: 10px;
}

.content-card {
  background: #fff;
  border-radius: 12px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.05);
  overflow: hidden;
}

.table-wrapper {
  overflow-x: auto;
}

.data-table {
  --el-table-border-color: #f3f4f6;
  --el-table-header-bg-color: #f9fafb;
}

.data-table :deep(th) {
  font-weight: 600;
  color: #374151;
  font-size: 13px;
}

.id-badge {
  display: inline-block;
  padding: 2px 8px;
  background: #f3f4f6;
  border-radius: 4px;
  font-size: 12px;
  font-weight: 500;
  color: #6b7280;
}

.task-id-link {
  color: #3b82f6;
  cursor: pointer;
  font-weight: 500;
}

.task-id-link:hover {
  text-decoration: underline;
}

.task-name {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.name-text {
  font-weight: 500;
  color: #1f2937;
}

.host-text {
  font-size: 13px;
  color: #6b7280;
}

.time-info {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.time-row {
  display: flex;
  align-items: center;
  gap: 4px;
  font-size: 13px;
  color: #374151;
}

.time-row.sub {
  color: #9ca3af;
  font-size: 12px;
}

.action-btns {
  display: flex;
  gap: 4px;
  justify-content: center;
}

.expand-content {
  padding: 16px 20px;
  background: #f9fafb;
}

.expand-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 12px;
  margin-bottom: 12px;
}

.expand-item {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.expand-label {
  font-size: 12px;
  color: #9ca3af;
}

.expand-value {
  font-size: 13px;
  color: #374151;
  font-weight: 500;
}

.expand-code {
  background: #f3f4f6;
  padding: 4px 8px;
  border-radius: 6px;
  font-family: ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, monospace;
  font-size: 12px;
  color: #4b5563;
}

.expand-full {
  margin-top: 12px;
}

.command-code {
  display: block;
  background: #1f2937;
  color: #a5f3fc;
  padding: 12px 16px;
  border-radius: 8px;
  font-family: ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, monospace;
  font-size: 12px;
  margin-top: 6px;
  overflow-x: auto;
}

.pagination-wrapper {
  padding: 16px 20px;
  display: flex;
  justify-content: flex-end;
  border-top: 1px solid #f3f4f6;
}

.card-list {
  padding: 12px;
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.log-card {
  background: #fff;
  border: 1px solid #e5e7eb;
  border-radius: 10px;
  overflow: hidden;
}

.log-card.status-running {
  border-left: 3px solid #3b82f6;
}

.log-card.status-success {
  border-left: 3px solid #22c55e;
}

.log-card.status-failed {
  border-left: 3px solid #ef4444;
}

.log-card.status-cancelled {
  border-left: 3px solid #9ca3af;
}

.log-card__header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px 16px;
  background: #f9fafb;
  border-bottom: 1px solid #f3f4f6;
}

.log-card__id {
  font-weight: 600;
  color: #6b7280;
  font-size: 13px;
}

.log-card__body {
  padding: 16px;
}

.log-card__title {
  font-size: 16px;
  font-weight: 600;
  color: #1f2937;
  margin: 0 0 10px 0;
}

.log-card__meta {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 8px;
}

.meta-item {
  display: flex;
  align-items: center;
  gap: 4px;
  color: #6b7280;
  font-size: 13px;
}

.log-card__host {
  font-size: 12px;
  color: #9ca3af;
  margin-top: 8px;
}

.log-card__footer {
  display: flex;
  border-top: 1px solid #f3f4f6;
  padding: 8px;
  gap: 4px;
}

.log-card__footer .el-button {
  flex: 1;
}

.result-dialog :deep(.el-dialog__body) {
  padding: 0;
}

.result-content {
  padding: 20px;
}

.result-section {
  margin-bottom: 20px;
}

.result-section:last-child {
  margin-bottom: 0;
}

.result-label {
  font-size: 13px;
  font-weight: 600;
  color: #374151;
  margin-bottom: 8px;
}

.result-value {
  font-size: 14px;
  color: #6b7280;
}

.host-value {
  padding: 10px 12px;
  background: #f9fafb;
  border-radius: 6px;
}

.result-code {
  background: #1f2937;
  color: #e5e7eb;
  padding: 16px;
  border-radius: 8px;
  font-family: ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, monospace;
  font-size: 13px;
  white-space: pre-wrap;
  word-wrap: break-word;
  margin: 0;
  overflow-x: auto;
}

.result-code.command {
  background: #374151;
  color: #a5f3fc;
}

.result-code.output {
  max-height: 50vh;
  overflow: auto;
}

.hidden-mobile {
  display: block;
}

.visible-mobile {
  display: none !important;
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

  .danger-btn {
    flex: 1;
  }

  .ghost-btn {
    flex: 1;
    justify-content: center;
  }

  .page-title {
    font-size: 20px;
  }

  .filter-header {
    padding: 12px 16px;
  }

  .filter-body {
    padding: 16px;
  }

  .filter-grid {
    grid-template-columns: 1fr;
  }

  .filter-actions {
    flex-direction: column;
  }

  .filter-actions .el-button {
    width: 100%;
  }

  .hidden-mobile {
    display: none !important;
  }

  .visible-mobile {
    display: block !important;
  }

  .pagination-wrapper {
    justify-content: center;
  }

  .result-content {
    padding: 16px;
  }

  .result-code {
    font-size: 12px;
    padding: 12px;
  }
}
</style>
