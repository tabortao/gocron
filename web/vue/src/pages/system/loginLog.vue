<template>
  <div class="page-container">
    <div class="page-header">
      <div class="page-header__left">
        <h1 class="page-title">{{ t('system.loginLog') }}</h1>
        <p class="page-desc" v-if="logTotal > 0">
          {{ isZh ? '共' : 'Total' }} <span class="count">{{ logTotal }}</span>
          {{ isZh ? '条记录' : 'records' }}
        </p>
      </div>
      <div class="page-header__right">
        <el-button @click="refresh" class="refresh-btn">
          <el-icon :class="{ 'is-loading': loading }"><Refresh /></el-icon>
          <span class="btn-text">{{ t('common.refresh') }}</span>
        </el-button>
      </div>
    </div>

    <div class="content-card">
      <div class="table-wrapper hidden-mobile">
        <el-table :data="logs" class="data-table" v-loading="loading">
          <el-table-column prop="id" label="ID" width="80" align="center">
            <template #default="{ row }">
              <span class="id-badge">{{ row.id }}</span>
            </template>
          </el-table-column>
          <el-table-column prop="username" :label="t('user.username')" min-width="120">
            <template #default="{ row }">
              <div class="user-cell">
                <div class="user-avatar">{{ row.username?.charAt(0)?.toUpperCase() || 'U' }}</div>
                <span>{{ row.username }}</span>
              </div>
            </template>
          </el-table-column>
          <el-table-column prop="ip" :label="t('system.loginIp')" min-width="140">
            <template #default="{ row }">
              <code class="ip-code">{{ row.ip }}</code>
            </template>
          </el-table-column>
          <el-table-column :label="t('system.loginTime')" min-width="180">
            <template #default="{ row }">
              <div class="time-cell">
                <el-icon><Clock /></el-icon>
                <span>{{ $filters.formatTime(row.created) }}</span>
              </div>
            </template>
          </el-table-column>
        </el-table>
      </div>

      <div class="card-list visible-mobile">
        <div v-for="log in logs" :key="log.id" class="log-card">
          <div class="log-card__header">
            <span class="log-card__id">#{{ log.id }}</span>
            <span class="log-card__time">{{ formatTime(log.created) }}</span>
          </div>
          <div class="log-card__body">
            <div class="log-card__user">
              <div class="user-avatar">{{ log.username?.charAt(0)?.toUpperCase() || 'U' }}</div>
              <span class="user-name">{{ log.username }}</span>
            </div>
            <div class="log-card__ip">
              <el-icon><Position /></el-icon>
              <code>{{ log.ip }}</code>
            </div>
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

    <div class="empty-wrapper" v-if="!loading && logs.length === 0">
      <el-empty :description="isZh ? '暂无登录记录' : 'No login records'">
        <el-button type="primary" @click="refresh">{{ t('common.refresh') }}</el-button>
      </el-empty>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { Refresh, Clock, Position } from '@element-plus/icons-vue'
import systemService from '../../api/system'

const { t, locale } = useI18n()
const isZh = computed(() => locale.value === 'zh-CN')

const isMobile = ref(false)
const loading = ref(false)
const logs = ref([])
const logTotal = ref(0)
const searchParams = ref({
  page_size: 20,
  page: 1
})

const checkMobile = () => {
  isMobile.value = window.innerWidth <= 768
}

const formatTime = time => {
  if (!time) return ''
  const d = new Date(time)
  return `${d.getMonth() + 1}/${d.getDate()} ${d.getHours().toString().padStart(2, '0')}:${d.getMinutes().toString().padStart(2, '0')}`
}

const fetchData = () => {
  loading.value = true
  systemService.loginLogList(searchParams.value, data => {
    logs.value = data.data
    logTotal.value = data.total
    loading.value = false
  })
}

const refresh = () => {
  fetchData()
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

.refresh-btn {
  border: 1px solid #e5e7eb;
  background: #fff;
  color: #374151;
  font-weight: 500;
}

.refresh-btn:hover {
  border-color: #3b82f6;
  color: #3b82f6;
}

.btn-text {
  margin-left: 6px;
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

.user-cell {
  display: flex;
  align-items: center;
  gap: 10px;
}

.user-avatar {
  width: 32px;
  height: 32px;
  border-radius: 8px;
  background: #e0e7ff;
  color: #4f46e5;
  display: flex;
  align-items: center;
  justify-content: center;
  font-weight: 600;
  font-size: 13px;
  flex-shrink: 0;
}

.ip-code {
  background: #f3f4f6;
  padding: 4px 10px;
  border-radius: 6px;
  font-family: ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, monospace;
  font-size: 13px;
  color: #4b5563;
}

.time-cell {
  display: flex;
  align-items: center;
  gap: 6px;
  color: #6b7280;
  font-size: 13px;
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

.log-card__time {
  font-size: 12px;
  color: #9ca3af;
}

.log-card__body {
  padding: 16px;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.log-card__user {
  display: flex;
  align-items: center;
  gap: 10px;
}

.log-card__ip {
  display: flex;
  align-items: center;
  gap: 6px;
  color: #6b7280;
  font-size: 13px;
}

.log-card__ip code {
  background: #f3f4f6;
  padding: 2px 8px;
  border-radius: 4px;
  font-size: 12px;
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
  }

  .refresh-btn {
    width: 100%;
    justify-content: center;
  }

  .page-title {
    font-size: 20px;
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
