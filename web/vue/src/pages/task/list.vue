<template>
  <div class="page-container">
    <div class="page-header">
      <div class="page-header__left">
        <h1 class="page-title">{{ t('task.list') }}</h1>
        <p class="page-desc" v-if="taskTotal > 0">
          {{ isZh ? '共' : 'Total' }} <span class="count">{{ taskTotal }}</span>
          {{ isZh ? '个任务' : 'tasks' }}
        </p>
      </div>
      <div class="page-header__right">
        <el-button type="primary" v-if="isAdmin" @click="toEdit(null)" class="primary-btn">
          <el-icon><Plus /></el-icon>
          <span>{{ t('common.add') }}</span>
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
              <el-input v-model.trim="searchParams.id" :placeholder="t('task.id')" clearable />
            </div>
            <div class="filter-item">
              <label>{{ t('task.name') }}</label>
              <el-input v-model.trim="searchParams.name" :placeholder="t('task.name')" clearable />
            </div>
            <div class="filter-item">
              <label>{{ t('task.tag') }}</label>
              <el-input v-model.trim="searchParams.tag" :placeholder="t('task.tag')" clearable />
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
              <label>{{ t('task.taskNode') }}</label>
              <el-select v-model.trim="searchParams.host_id" clearable filterable>
                <el-option :label="t('select')" value="" />
                <el-option
                  v-for="item in hosts"
                  :key="item.id"
                  :label="item.alias"
                  :value="item.id"
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

    <transition name="slide-up">
      <div class="batch-bar" v-if="isAdmin && selectedTasks.length > 0">
        <div class="batch-info">
          <span class="batch-count">{{ selectedTasks.length }}</span>
          <span>{{ t('message.selected') }}</span>
        </div>
        <div class="batch-actions">
          <el-button type="success" size="small" @click="batchEnable">
            <el-icon><VideoPlay /></el-icon>
            <span class="hidden-mobile">{{ t('message.batchEnable') }}</span>
          </el-button>
          <el-button type="warning" size="small" @click="batchDisable">
            <el-icon><VideoPause /></el-icon>
            <span class="hidden-mobile">{{ t('message.batchDisable') }}</span>
          </el-button>
          <el-button type="danger" size="small" @click="batchRemove">
            <el-icon><Delete /></el-icon>
            <span class="hidden-mobile">{{ t('message.batchDelete') }}</span>
          </el-button>
        </div>
      </div>
    </transition>

    <div class="content-card">
      <div class="table-wrapper hidden-mobile">
        <el-table
          :data="tasks"
          @selection-change="handleSelectionChange"
          class="data-table"
          v-loading="loading"
        >
          <el-table-column type="selection" width="50" v-if="isAdmin" />
          <el-table-column type="expand">
            <template #default="{ row }">
              <div class="expand-content">
                <div class="expand-grid">
                  <div class="expand-item">
                    <span class="expand-label">{{ t('message.taskCreatedTime') }}</span>
                    <span class="expand-value">{{ $filters.formatTime(row.created) }}</span>
                  </div>
                  <div class="expand-item">
                    <span class="expand-label">{{ t('message.taskType') }}</span>
                    <span class="expand-value">{{ formatLevel(row.level) }}</span>
                  </div>
                  <div class="expand-item">
                    <span class="expand-label">{{ t('message.singleInstanceRun') }}</span>
                    <span class="expand-value">{{ formatMulti(row.multi) }}</span>
                  </div>
                  <div class="expand-item">
                    <span class="expand-label">{{ t('message.timeoutTime') }}</span>
                    <span class="expand-value">{{ formatTimeout(row.timeout) }}</span>
                  </div>
                </div>
                <div class="expand-full" v-if="row.hosts?.length">
                  <span class="expand-label">{{ t('message.taskNodeLabel') }}</span>
                  <div class="host-tags">
                    <el-tag v-for="h in row.hosts" :key="h.host_id" size="small" type="info">
                      {{ h.alias }}
                    </el-tag>
                  </div>
                </div>
                <div class="expand-full" v-if="row.command">
                  <span class="expand-label">{{ t('message.commandLabel') }}</span>
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
          <el-table-column prop="name" :label="t('task.name')" min-width="140">
            <template #default="{ row }">
              <div class="task-name">
                <span class="name-text">{{ row.name }}</span>
                <el-tag v-if="row.tag" size="small" type="info" effect="plain">{{
                  row.tag
                }}</el-tag>
              </div>
            </template>
          </el-table-column>
          <el-table-column prop="spec" :label="t('task.cronExpression')" width="120">
            <template #default="{ row }">
              <code class="cron-code">{{ row.spec }}</code>
            </template>
          </el-table-column>
          <el-table-column :label="t('task.nextRunTime')" width="160">
            <template #default="{ row }">
              <span class="time-text">{{ $filters.formatTime(row.next_run_time) }}</span>
            </template>
          </el-table-column>
          <el-table-column :label="t('task.protocol')" width="100" align="center">
            <template #default="{ row }">
              <el-tag :type="getProtocolType(row)" size="small" effect="plain">
                {{ formatProtocol(row) }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column :label="t('common.status')" width="80" align="center">
            <template #default="{ row }">
              <el-switch
                v-if="row.level === 1"
                v-model="row.status"
                :active-value="1"
                :inactive-value="0"
                @change="changeStatus(row)"
                :disabled="!isAdmin"
                inline-prompt
                class="status-switch"
              />
            </template>
          </el-table-column>
          <el-table-column :label="t('common.operation')" width="180" v-if="isAdmin" fixed="right">
            <template #default="{ row }">
              <div class="action-btns">
                <el-button type="primary" size="small" text @click="toEdit(row)">
                  <el-icon><Edit /></el-icon>
                  {{ t('common.edit') }}
                </el-button>
                <el-button type="success" size="small" text @click="runTask(row)">
                  <el-icon><VideoPlay /></el-icon>
                </el-button>
                <el-button type="info" size="small" text @click="jumpToLog(row)">
                  <el-icon><Document /></el-icon>
                </el-button>
                <el-button type="danger" size="small" text @click="remove(row)">
                  <el-icon><Delete /></el-icon>
                </el-button>
              </div>
            </template>
          </el-table-column>
        </el-table>
      </div>

      <div class="card-list visible-mobile">
        <div
          v-for="task in tasks"
          :key="task.id"
          class="task-card"
          :class="{ disabled: task.status === 0 }"
        >
          <div class="task-card__header">
            <span class="task-card__id">#{{ task.id }}</span>
            <el-switch
              v-if="task.level === 1 && isAdmin"
              v-model="task.status"
              :active-value="1"
              :inactive-value="0"
              @change="changeStatus(task)"
              size="small"
            />
          </div>
          <div class="task-card__body">
            <h3 class="task-card__title">{{ task.name }}</h3>
            <div class="task-card__meta">
              <span class="meta-item">
                <el-icon><Clock /></el-icon>
                <code>{{ task.spec }}</code>
              </span>
              <el-tag :type="getProtocolType(task)" size="small" effect="plain">
                {{ formatProtocol(task) }}
              </el-tag>
            </div>
            <el-tag v-if="task.tag" size="small" type="info" effect="plain">{{ task.tag }}</el-tag>
          </div>
          <div class="task-card__footer">
            <el-button type="primary" size="small" text @click="toEdit(task)">
              <el-icon><Edit /></el-icon>
              {{ t('common.edit') }}
            </el-button>
            <el-button type="success" size="small" text @click="runTask(task)">
              <el-icon><VideoPlay /></el-icon>
              {{ t('task.manualRun') }}
            </el-button>
            <el-button type="info" size="small" text @click="jumpToLog(task)">
              <el-icon><Document /></el-icon>
              {{ t('task.viewLog') }}
            </el-button>
          </div>
        </div>
        <el-empty v-if="!loading && tasks.length === 0" :description="t('message.noData')" />
      </div>

      <div class="pagination-wrapper">
        <el-pagination
          background
          :layout="isMobile ? 'prev, pager, next' : 'total, sizes, prev, pager, next'"
          :total="taskTotal"
          v-model:current-page="searchParams.page"
          v-model:page-size="searchParams.page_size"
          @size-change="changePageSize"
          @current-change="changePage"
          :page-sizes="[10, 20, 50, 100]"
          small
        />
      </div>
    </div>

    <div class="empty-wrapper" v-if="!loading && tasks.length === 0">
      <el-empty :description="t('message.noTaskTip')">
        <el-button type="primary" v-if="isAdmin" @click="toEdit(null)">
          <el-icon><Plus /></el-icon>
          {{ t('common.add') }}
        </el-button>
      </el-empty>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onActivated, onUnmounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRouter, useRoute } from 'vue-router'
import { useUserStore } from '../../stores/user'
import { ElMessageBox, ElMessage } from 'element-plus'
import {
  Plus,
  Refresh,
  Search,
  ArrowDown,
  Edit,
  Delete,
  VideoPlay,
  VideoPause,
  Document,
  Clock
} from '@element-plus/icons-vue'
import taskService from '../../api/task'

const { t, locale } = useI18n()
const router = useRouter()
const route = useRoute()
const userStore = useUserStore()

const isZh = computed(() => locale.value === 'zh-CN')
const isAdmin = userStore.isAdmin

const isMobile = ref(false)
const loading = ref(false)
const tasks = ref([])
const hosts = ref([])
const taskTotal = ref(0)
const selectedTasks = ref([])
const filterExpanded = ref(true)

const checkMobile = () => {
  isMobile.value = window.innerWidth <= 768
}

const searchParams = ref({
  page_size: 20,
  page: 1,
  id: '',
  protocol: '',
  name: '',
  tag: '',
  host_id: '',
  status: ''
})

const protocolList = [
  { value: '1', label: 'HTTP' },
  { value: '2', label: 'Shell' }
]

const statusList = computed(() => [
  { value: '2', label: t('message.activated') },
  { value: '1', label: t('message.stopped') }
])

const formatLevel = value => (value === 1 ? t('task.mainTask') : t('task.childTask'))
const formatTimeout = value => (value > 0 ? value + t('message.seconds') : t('message.noLimit'))
const formatMulti = value => (value > 0 ? t('common.no') : t('common.yes'))

const formatProtocol = row => {
  if (row.protocol === 2) return 'Shell'
  return row.http_method === 1 ? 'GET' : 'POST'
}

const getProtocolType = row => {
  if (row.protocol === 2) return 'warning'
  return row.http_method === 1 ? 'success' : 'primary'
}

const fetchData = (callback = null) => {
  loading.value = true
  taskService.list(searchParams.value, (data, hostList) => {
    tasks.value = data.data
    taskTotal.value = data.total
    hosts.value = hostList
    loading.value = false
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
    id: '',
    protocol: '',
    name: '',
    tag: '',
    host_id: '',
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

const changeStatus = item => {
  if (item.status) {
    taskService.enable(item.id, () => fetchData())
  } else {
    taskService.disable(item.id, () => fetchData())
  }
}

const handleSelectionChange = selection => {
  selectedTasks.value = selection.filter(task => task.level === 1)
}

const toEdit = item => {
  const path = item ? `/task/edit/${item.id}` : '/task/create'
  router.push(path)
}

const runTask = item => {
  ElMessageBox.confirm(
    t('message.confirmRunTask', { name: item.name }),
    t('message.manualRunTask'),
    {
      confirmButtonText: t('message.confirmExecute'),
      cancelButtonText: t('common.cancel'),
      type: 'info'
    }
  )
    .then(() => {
      taskService.run(item.id, () => {
        ElMessage.success(t('message.taskStarted'))
      })
    })
    .catch(() => {})
}

const remove = item => {
  ElMessageBox.confirm(
    t('message.confirmDeleteTask', { name: item.name }),
    t('message.confirmDeleteTitle'),
    {
      confirmButtonText: t('common.confirm'),
      cancelButtonText: t('common.cancel'),
      type: 'warning'
    }
  )
    .then(() => {
      taskService.remove(item.id, () => refresh())
    })
    .catch(() => {})
}

const jumpToLog = item => {
  router.push(`/task/log?task_id=${item.id}`)
}

const batchEnable = () => {
  if (!selectedTasks.value.length) return
  ElMessageBox.confirm(
    t('message.confirmBatchEnable', { count: selectedTasks.value.length }),
    t('message.batchEnable'),
    { confirmButtonText: t('common.confirm'), cancelButtonText: t('common.cancel'), type: 'info' }
  )
    .then(() => {
      taskService.batchEnable(
        selectedTasks.value.map(t => t.id),
        () => {
          ElMessage.success(t('message.batchEnableSuccess'))
          selectedTasks.value = []
          fetchData()
        }
      )
    })
    .catch(() => {})
}

const batchDisable = () => {
  if (!selectedTasks.value.length) return
  ElMessageBox.confirm(
    t('message.confirmBatchDisable', { count: selectedTasks.value.length }),
    t('message.batchDisable'),
    {
      confirmButtonText: t('common.confirm'),
      cancelButtonText: t('common.cancel'),
      type: 'warning'
    }
  )
    .then(() => {
      taskService.batchDisable(
        selectedTasks.value.map(t => t.id),
        () => {
          ElMessage.success(t('message.batchDisableSuccess'))
          selectedTasks.value = []
          fetchData()
        }
      )
    })
    .catch(() => {})
}

const batchRemove = () => {
  if (!selectedTasks.value.length) return
  ElMessageBox.confirm(
    t('message.confirmBatchDelete', { count: selectedTasks.value.length }),
    t('message.batchDelete'),
    {
      confirmButtonText: t('message.confirmDeleteButton'),
      cancelButtonText: t('common.cancel'),
      type: 'error'
    }
  )
    .then(() => {
      taskService.batchRemove(
        selectedTasks.value.map(t => t.id),
        () => {
          ElMessage.success(t('message.batchDeleteSuccess'))
          selectedTasks.value = []
          fetchData()
        }
      )
    })
    .catch(() => {})
}

let isFirstActivate = true
onMounted(() => {
  checkMobile()
  window.addEventListener('resize', checkMobile)
  if (route.query.host_id) {
    searchParams.value.host_id = route.query.host_id
  }
  fetchData()
})

onActivated(() => {
  if (isFirstActivate) {
    isFirstActivate = false
    return
  }
  fetchData()
})

onUnmounted(() => {
  window.removeEventListener('resize', checkMobile)
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

.primary-btn {
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

.batch-bar {
  background: #fef3c7;
  border: 1px solid #fcd34d;
  border-radius: 10px;
  padding: 12px 20px;
  margin-bottom: 16px;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.batch-count {
  background: #f59e0b;
  color: #fff;
  font-weight: 600;
  padding: 2px 10px;
  border-radius: 20px;
  margin-right: 8px;
}

.batch-actions {
  display: flex;
  gap: 8px;
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

.task-name {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.name-text {
  font-weight: 500;
  color: #1f2937;
}

.cron-code {
  background: #f3f4f6;
  padding: 4px 8px;
  border-radius: 6px;
  font-family: ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, monospace;
  font-size: 12px;
  color: #4b5563;
}

.time-text {
  font-size: 13px;
  color: #6b7280;
}

.status-switch {
  --el-switch-on-color: #22c55e;
  --el-switch-off-color: #ef4444;
}

.action-btns {
  display: flex;
  gap: 4px;
}

.expand-content {
  padding: 16px 20px;
  background: #f9fafb;
}

.expand-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
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

.expand-full {
  margin-top: 12px;
}

.host-tags {
  display: flex;
  gap: 6px;
  flex-wrap: wrap;
  margin-top: 6px;
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

.task-card {
  background: #fff;
  border: 1px solid #e5e7eb;
  border-radius: 10px;
  overflow: hidden;
}

.task-card.disabled {
  opacity: 0.6;
}

.task-card__header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px 16px;
  background: #f9fafb;
  border-bottom: 1px solid #f3f4f6;
}

.task-card__id {
  font-weight: 600;
  color: #6b7280;
  font-size: 13px;
}

.task-card__body {
  padding: 16px;
}

.task-card__title {
  font-size: 16px;
  font-weight: 600;
  color: #1f2937;
  margin: 0 0 10px 0;
}

.task-card__meta {
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

.meta-item code {
  background: #f3f4f6;
  padding: 2px 6px;
  border-radius: 4px;
  font-family: ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, monospace;
  font-size: 11px;
}

.task-card__footer {
  display: flex;
  border-top: 1px solid #f3f4f6;
  padding: 8px;
  gap: 4px;
}

.task-card__footer .el-button {
  flex: 1;
}

.empty-wrapper {
  padding: 60px 20px;
}

.hidden-mobile {
  display: block;
}

.visible-mobile {
  display: none !important;
}

.slide-up-enter-active,
.slide-up-leave-active {
  transition: all 0.3s ease;
}

.slide-up-enter-from,
.slide-up-leave-to {
  opacity: 0;
  transform: translateY(-10px);
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

  .primary-btn {
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

  .batch-bar {
    flex-direction: column;
    gap: 12px;
    padding: 12px 16px;
  }

  .batch-actions {
    width: 100%;
    justify-content: center;
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

  .expand-grid {
    grid-template-columns: 1fr 1fr;
  }
}
</style>
