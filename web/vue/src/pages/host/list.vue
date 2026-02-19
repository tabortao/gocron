<template>
  <div class="page-container">
    <div class="page-header">
      <div class="page-header__left">
        <h1 class="page-title">{{ t('host.list') }}</h1>
        <p class="page-desc" v-if="hostTotal > 0">
          {{ isZh ? '共' : 'Total' }} <span class="count">{{ hostTotal }}</span>
          {{ isZh ? '个节点' : 'nodes' }}
        </p>
      </div>
      <div class="page-header__right">
        <el-button type="success" v-if="isAdmin" @click="showAgentInstall" class="success-btn">
          <el-icon><Download /></el-icon>
          <span class="btn-text">{{ t('host.autoRegister') }}</span>
        </el-button>
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
              <label>ID</label>
              <el-input v-model.trim="searchParams.id" placeholder="ID" clearable />
            </div>
            <div class="filter-item">
              <label>{{ t('host.name') }}</label>
              <el-input v-model.trim="searchParams.name" :placeholder="t('host.name')" clearable />
            </div>
            <div class="filter-item hidden-mobile">
              <label>{{ t('host.alias') }}</label>
              <el-input
                v-model.trim="searchParams.alias"
                :placeholder="t('host.alias')"
                clearable
              />
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
        <el-table :data="hosts" class="data-table" v-loading="loading">
          <el-table-column prop="id" label="ID" width="70" align="center">
            <template #default="{ row }">
              <span class="id-badge">{{ row.id }}</span>
            </template>
          </el-table-column>
          <el-table-column prop="alias" :label="t('host.alias')" min-width="120">
            <template #default="{ row }">
              <div class="node-name">
                <span class="name-text">{{ row.alias }}</span>
              </div>
            </template>
          </el-table-column>
          <el-table-column prop="name" :label="t('host.name')" min-width="140">
            <template #default="{ row }">
              <code class="host-code">{{ row.name }}</code>
            </template>
          </el-table-column>
          <el-table-column prop="port" :label="t('host.port')" width="100" align="center">
            <template #default="{ row }">
              <el-tag type="info" size="small" effect="plain">{{ row.port }}</el-tag>
            </template>
          </el-table-column>
          <el-table-column :label="t('task.viewLog')" width="120" align="center">
            <template #default="{ row }">
              <el-button type="success" size="small" text @click="toTasks(row)">
                <el-icon><Document /></el-icon>
                {{ t('task.list') }}
              </el-button>
            </template>
          </el-table-column>
          <el-table-column prop="remark" :label="t('host.remark')" min-width="140">
            <template #default="{ row }">
              <span class="remark-text">{{ row.remark || '-' }}</span>
            </template>
          </el-table-column>
          <el-table-column
            :label="t('common.operation')"
            :width="locale === 'zh-CN' ? 200 : 220"
            v-if="isAdmin"
            fixed="right"
          >
            <template #default="{ row }">
              <div class="action-btns">
                <el-button type="primary" size="small" text @click="toEdit(row)">
                  <el-icon><Edit /></el-icon>
                  {{ t('common.edit') }}
                </el-button>
                <el-button type="info" size="small" text @click="ping(row)">
                  <el-icon><Connection /></el-icon>
                  <span class="hidden-mobile">{{ t('system.testSend') }}</span>
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
        <div v-for="host in hosts" :key="host.id" class="node-card">
          <div class="node-card__header">
            <span class="node-card__id">#{{ host.id }}</span>
            <el-tag type="info" size="small" effect="plain">{{ host.port }}</el-tag>
          </div>
          <div class="node-card__body">
            <h3 class="node-card__title">{{ host.alias }}</h3>
            <div class="node-card__meta">
              <span class="meta-item">
                <el-icon><Monitor /></el-icon>
                <code>{{ host.name }}</code>
              </span>
            </div>
            <p class="node-card__remark" v-if="host.remark">{{ host.remark }}</p>
          </div>
          <div class="node-card__footer">
            <el-button type="success" size="small" text @click="toTasks(host)">
              <el-icon><Document /></el-icon>
              {{ t('task.list') }}
            </el-button>
            <el-button type="primary" size="small" text @click="toEdit(host)" v-if="isAdmin">
              <el-icon><Edit /></el-icon>
              {{ t('common.edit') }}
            </el-button>
            <el-button type="info" size="small" text @click="ping(host)" v-if="isAdmin">
              <el-icon><Connection /></el-icon>
              {{ t('system.testSend') }}
            </el-button>
            <el-button type="danger" size="small" text @click="remove(host)" v-if="isAdmin">
              <el-icon><Delete /></el-icon>
            </el-button>
          </div>
        </div>
        <el-empty v-if="!loading && hosts.length === 0" :description="t('message.noData')" />
      </div>

      <div class="pagination-wrapper">
        <el-pagination
          background
          :layout="isMobile ? 'prev, pager, next' : 'total, sizes, prev, pager, next'"
          :total="hostTotal"
          v-model:current-page="searchParams.page"
          v-model:page-size="searchParams.page_size"
          @size-change="changePageSize"
          @current-change="changePage"
          :page-sizes="[10, 20, 50, 100]"
          small
        />
      </div>
    </div>

    <div class="empty-wrapper" v-if="!loading && hosts.length === 0">
      <el-empty :description="t('message.noNodeTip')">
        <el-button type="primary" v-if="isAdmin" @click="toEdit(null)">
          <el-icon><Plus /></el-icon>
          {{ t('common.add') }}
        </el-button>
      </el-empty>
    </div>

    <el-dialog v-model="agentDialogVisible" :title="t('host.agentInstall')" width="750px">
      <div v-if="installCommand">
        <el-alert type="info" :closable="false" style="margin-bottom: 20px" show-icon>
          <template #title>{{ t('host.installTipTitle') }}</template>
          <div style="line-height: 1.6">
            <div>{{ t('host.installTipLine1') }}</div>
            <div>{{ t('host.installTipLine2') }}</div>
            <div style="margin-top: 6px">{{ t('host.installTipLine3') }}</div>
            <div>{{ t('host.installTipLine4') }}</div>
          </div>
        </el-alert>

        <el-tabs v-model="activeTab" type="card">
          <el-tab-pane label="Linux / macOS" name="linux">
            <div style="padding: 15px; background: #f5f7fa; border-radius: 4px">
              <div style="margin-bottom: 10px; color: #606266; font-size: 14px">
                <el-icon style="vertical-align: middle"><Monitor /></el-icon>
                {{ t('host.installModeNormal') }}
              </div>
              <div style="margin-bottom: 10px; color: #909399; font-size: 13px">
                {{ t('host.installModeNormalTip') }}
              </div>
              <el-input
                v-model="installCommand"
                type="textarea"
                :rows="3"
                readonly
                style="font-family: monospace; font-size: 13px"
              />
              <div style="margin-top: 10px; text-align: right">
                <el-button type="primary" @click="copyCommand('linux')" icon="DocumentCopy"
                  >Copy</el-button
                >
              </div>
              <el-divider style="margin: 18px 0" />
              <div style="margin-bottom: 10px; color: #606266; font-size: 14px">
                <el-icon style="vertical-align: middle"><Monitor /></el-icon>
                {{ t('host.installModeAllowRoot') }}
              </div>
              <div style="margin-bottom: 10px; color: #909399; font-size: 13px">
                {{ t('host.installModeAllowRootTip') }}
              </div>
              <el-input
                v-model="installCommandAllowRoot"
                type="textarea"
                :rows="3"
                readonly
                style="font-family: monospace; font-size: 13px"
              />
              <div style="margin-top: 10px; text-align: right">
                <el-button
                  type="primary"
                  @click="copyCommand('linux-allow-root')"
                  icon="DocumentCopy"
                  >Copy</el-button
                >
              </div>
            </div>
          </el-tab-pane>

          <el-tab-pane label="Windows" name="windows">
            <div style="padding: 15px">
              <el-alert type="warning" :closable="false" style="margin-bottom: 15px">
                <template #title>
                  <strong>{{ t('host.windowsManualInstall') }}</strong>
                </template>
                {{ t('host.windowsManualInstallTip') }}
              </el-alert>

              <el-steps direction="vertical" :active="3">
                <el-step
                  :title="t('host.windowsStep1')"
                  :description="t('host.windowsStep1Desc')"
                />
                <el-step
                  :title="t('host.windowsStep2')"
                  :description="t('host.windowsStep2Desc')"
                />
                <el-step
                  :title="t('host.windowsStep3')"
                  :description="t('host.windowsStep3Desc')"
                />
              </el-steps>
            </div>
          </el-tab-pane>
        </el-tabs>

        <el-divider />

        <div style="padding: 10px 0">
          <el-descriptions :column="1" border>
            <el-descriptions-item :label="t('host.tokenExpires')">
              <el-tag type="warning" effect="plain">{{ expiresAt }}</el-tag>
            </el-descriptions-item>
            <el-descriptions-item :label="t('host.tokenUsage')">
              <span style="color: #67c23a">{{ t('host.tokenReusable') }}</span>
            </el-descriptions-item>
          </el-descriptions>
        </div>
      </div>
      <div v-else style="text-align: center; padding: 20px">
        <el-icon class="is-loading" :size="30"><Loading /></el-icon>
        <p>{{ t('common.loading') }}</p>
      </div>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted, watch } from 'vue'
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
  Download,
  Document,
  Connection,
  Monitor,
  Loading
} from '@element-plus/icons-vue'
import hostService from '../../api/host'
import agentService from '../../api/agent'
import { copyText } from '../../utils/clipboard'

const { t, locale } = useI18n()
const router = useRouter()
const route = useRoute()
const userStore = useUserStore()

const isZh = computed(() => locale.value === 'zh-CN')
const isAdmin = userStore.isAdmin

const isMobile = ref(false)
const loading = ref(false)
const hosts = ref([])
const hostTotal = ref(0)
const filterExpanded = ref(true)
const agentDialogVisible = ref(false)
const installCommand = ref('')
const expiresAt = ref('')
const activeTab = ref('linux')
const cachedToken = ref(null)
const cachedTokenExpires = ref(null)

const checkMobile = () => {
  isMobile.value = window.innerWidth <= 768
}

const searchParams = ref({
  page_size: 20,
  page: 1,
  id: '',
  name: '',
  alias: ''
})

const installCommandAllowRoot = computed(() => {
  if (!installCommand.value) return ''
  return installCommand.value.replace(/\|\s*bash\s*$/, '| sudo bash')
})

const fetchData = (callback = null) => {
  loading.value = true
  hostService.list(searchParams.value, data => {
    hosts.value = data.data
    hostTotal.value = data.total
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
    name: '',
    alias: ''
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

const remove = item => {
  ElMessageBox.confirm(t('message.confirmDeleteNode'), t('common.tip'), {
    confirmButtonText: t('common.confirm'),
    cancelButtonText: t('common.cancel'),
    type: 'warning',
    center: true
  })
    .then(() => {
      hostService.remove(item.id, () => refresh())
    })
    .catch(() => {})
}

const ping = item => {
  if (!item.id || item.id <= 0) {
    ElMessage.error(t('message.dataNotFound'))
    return
  }
  hostService.ping(item.id, () => {
    ElMessage.success(t('message.connectionSuccess'))
  })
}

const toEdit = item => {
  const path = item ? `/host/edit/${item.id}` : '/host/create'
  router.push(path)
}

const toTasks = item => {
  router.push({
    path: '/task',
    query: {
      host_id: item.id
    }
  })
}

const showAgentInstall = () => {
  agentDialogVisible.value = true

  const now = new Date()
  if (cachedToken.value && cachedTokenExpires.value && now < cachedTokenExpires.value) {
    installCommand.value = cachedToken.value.install_cmd
    expiresAt.value = cachedTokenExpires.value.toLocaleString()
    return
  }

  installCommand.value = ''
  expiresAt.value = ''
  agentService.generateToken(data => {
    installCommand.value = data.install_cmd
    const expiresDate = new Date(data.expires_at)
    expiresAt.value = expiresDate.toLocaleString()

    cachedToken.value = data
    cachedTokenExpires.value = expiresDate
  })
}

const copyCommand = type => {
  const cmd =
    type === 'windows'
      ? installCommand.value
      : type === 'linux-allow-root'
        ? installCommandAllowRoot.value
        : installCommand.value
  copyText(cmd)
    .then(() => {
      ElMessage.success(t('message.copySuccess'))
    })
    .catch(() => {
      ElMessage.error(t('message.copyFailed'))
    })
}

watch(
  () => route.path,
  (to, from) => {
    if (to === '/host' && (from === '/host/create' || from?.startsWith('/host/edit/'))) {
      fetchData()
    }
  }
)

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

.primary-btn,
.success-btn {
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

.node-name {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.name-text {
  font-weight: 500;
  color: #1f2937;
}

.host-code {
  background: #f3f4f6;
  padding: 4px 8px;
  border-radius: 6px;
  font-family: ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, monospace;
  font-size: 12px;
  color: #4b5563;
}

.remark-text {
  font-size: 13px;
  color: #6b7280;
}

.action-btns {
  display: flex;
  gap: 4px;
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

.node-card {
  background: #fff;
  border: 1px solid #e5e7eb;
  border-radius: 10px;
  overflow: hidden;
}

.node-card__header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px 16px;
  background: #f9fafb;
  border-bottom: 1px solid #f3f4f6;
}

.node-card__id {
  font-weight: 600;
  color: #6b7280;
  font-size: 13px;
}

.node-card__body {
  padding: 16px;
}

.node-card__title {
  font-size: 16px;
  font-weight: 600;
  color: #1f2937;
  margin: 0 0 10px 0;
}

.node-card__meta {
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

.node-card__remark {
  font-size: 13px;
  color: #9ca3af;
  margin: 8px 0 0 0;
}

.node-card__footer {
  display: flex;
  border-top: 1px solid #f3f4f6;
  padding: 8px;
  gap: 4px;
  flex-wrap: wrap;
}

.node-card__footer .el-button {
  flex: 1;
  min-width: 60px;
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
    flex-wrap: wrap;
  }

  .primary-btn,
  .success-btn,
  .ghost-btn {
    flex: 1;
    justify-content: center;
    min-width: 100px;
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
}
</style>
