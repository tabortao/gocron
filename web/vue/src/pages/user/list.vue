<template>
  <div class="page-container">
    <div class="page-header">
      <div class="page-header__left">
        <h1 class="page-title">{{ t('user.list') }}</h1>
        <p class="page-desc" v-if="userTotal > 0">
          {{ isZh ? '共' : 'Total' }} <span class="count">{{ userTotal }}</span>
          {{ isZh ? '个用户' : 'users' }}
        </p>
      </div>
      <div class="page-header__right">
        <el-button type="primary" v-if="isAdmin" @click="toEdit(null)" class="primary-btn">
          <el-icon><Plus /></el-icon>
          <span class="btn-text">{{ t('common.add') }}</span>
        </el-button>
        <el-button @click="refresh" class="ghost-btn">
          <el-icon :class="{ 'is-loading': loading }"><Refresh /></el-icon>
          <span class="btn-text">{{ t('common.refresh') }}</span>
        </el-button>
      </div>
    </div>

    <div class="content-card">
      <div class="table-wrapper hidden-mobile">
        <el-table :data="users" class="data-table" v-loading="loading">
          <el-table-column prop="id" label="ID" width="70" align="center">
            <template #default="{ row }">
              <span class="id-badge">{{ row.id }}</span>
            </template>
          </el-table-column>
          <el-table-column prop="name" :label="t('user.username')" min-width="120">
            <template #default="{ row }">
              <div class="user-info">
                <span class="user-name">{{ row.name }}</span>
                <el-tag v-if="row.is_admin === 1" type="danger" size="small" effect="plain">
                  {{ t('user.admin') }}
                </el-tag>
              </div>
            </template>
          </el-table-column>
          <el-table-column prop="email" :label="t('user.email')" min-width="180">
            <template #default="{ row }">
              <span class="email-text">{{ row.email || '-' }}</span>
            </template>
          </el-table-column>
          <el-table-column :label="t('user.role')" width="120" align="center">
            <template #default="{ row }">
              <el-tag :type="row.is_admin === 1 ? 'danger' : 'info'" size="small" effect="plain">
                {{ formatRole(row) }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column :label="t('common.status')" width="100" align="center">
            <template #default="{ row }">
              <el-switch
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
          <el-table-column :label="t('common.operation')" width="220" v-if="isAdmin" fixed="right">
            <template #default="{ row }">
              <div class="action-btns">
                <el-button type="primary" size="small" text @click="toEdit(row)">
                  <el-icon><Edit /></el-icon>
                  {{ t('common.edit') }}
                </el-button>
                <el-button type="warning" size="small" text @click="editPassword(row)">
                  <el-icon><Key /></el-icon>
                  <span class="hidden-mobile">{{ t('user.changePassword') }}</span>
                </el-button>
                <el-button type="danger" size="small" text @click="remove(row)">
                  <el-icon><Delete /></el-icon>
                  {{ t('common.delete') }}
                </el-button>
              </div>
            </template>
          </el-table-column>
        </el-table>
      </div>

      <div class="card-list visible-mobile">
        <div
          v-for="user in users"
          :key="user.id"
          class="user-card"
          :class="{ disabled: user.status === 0 }"
        >
          <div class="user-card__header">
            <div class="user-card__avatar">
              <el-icon><User /></el-icon>
            </div>
            <div class="user-card__info">
              <div class="user-card__name">
                {{ user.name }}
                <el-tag v-if="user.is_admin === 1" type="danger" size="small" effect="plain">
                  {{ t('user.admin') }}
                </el-tag>
              </div>
              <div class="user-card__email">{{ user.email || '-' }}</div>
            </div>
            <el-switch
              v-model="user.status"
              :active-value="1"
              :inactive-value="0"
              @change="changeStatus(user)"
              v-if="isAdmin"
              size="small"
            />
          </div>
          <div class="user-card__body">
            <div class="user-card__meta">
              <span class="meta-item">
                <el-icon><Stamp /></el-icon>
                <span>{{ t('user.role') }}: {{ formatRole(user) }}</span>
              </span>
            </div>
          </div>
          <div class="user-card__footer" v-if="isAdmin">
            <el-button type="primary" size="small" text @click="toEdit(user)">
              <el-icon><Edit /></el-icon>
              {{ t('common.edit') }}
            </el-button>
            <el-button type="warning" size="small" text @click="editPassword(user)">
              <el-icon><Key /></el-icon>
              {{ t('user.changePassword') }}
            </el-button>
            <el-button type="danger" size="small" text @click="remove(user)">
              <el-icon><Delete /></el-icon>
              {{ t('common.delete') }}
            </el-button>
          </div>
        </div>
        <el-empty v-if="!loading && users.length === 0" :description="t('message.noData')" />
      </div>

      <div class="pagination-wrapper">
        <el-pagination
          background
          :layout="isMobile ? 'prev, pager, next' : 'total, sizes, prev, pager, next'"
          :total="userTotal"
          v-model:current-page="searchParams.page"
          v-model:page-size="searchParams.page_size"
          @size-change="changePageSize"
          @current-change="changePage"
          :page-sizes="[10, 20, 50, 100]"
          small
        />
      </div>
    </div>

    <div class="empty-wrapper" v-if="!loading && users.length === 0">
      <el-empty :description="t('message.noUserTip')">
        <el-button type="primary" v-if="isAdmin" @click="toEdit(null)">
          <el-icon><Plus /></el-icon>
          {{ t('common.add') }}
        </el-button>
      </el-empty>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRouter } from 'vue-router'
import { useUserStore } from '../../stores/user'
import { ElMessageBox, ElMessage } from 'element-plus'
import { Plus, Refresh, Edit, Delete, Key, User, Stamp } from '@element-plus/icons-vue'
import userService from '../../api/user'

const { t, locale } = useI18n()
const router = useRouter()
const userStore = useUserStore()

const isZh = computed(() => locale.value === 'zh-CN')
const isAdmin = userStore.isAdmin

const isMobile = ref(false)
const loading = ref(false)
const users = ref([])
const userTotal = ref(0)

const checkMobile = () => {
  isMobile.value = window.innerWidth <= 768
}

const searchParams = ref({
  page_size: 20,
  page: 1
})

const formatRole = row => {
  return row.is_admin === 1 ? t('user.admin') : t('user.normalUser')
}

const fetchData = (callback = null) => {
  loading.value = true
  userService.list(searchParams.value, data => {
    users.value = data.data
    userTotal.value = data.total
    loading.value = false
    if (callback) callback()
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
    userService.enable(item.id)
  } else {
    userService.disable(item.id)
  }
}

const toEdit = item => {
  const path = item ? `/user/edit/${item.id}` : '/user/create'
  router.push(path)
}

const editPassword = item => {
  router.push(`/user/edit-password/${item.id}`)
}

const refresh = () => {
  fetchData(() => {
    ElMessage.success(t('message.refreshSuccess'))
  })
}

const remove = item => {
  ElMessageBox.confirm(t('message.confirmDeleteUser'), t('common.tip'), {
    confirmButtonText: t('common.confirm'),
    cancelButtonText: t('common.cancel'),
    type: 'warning',
    center: true
  })
    .then(() => {
      userService.remove(item.id, () => {
        refresh()
      })
    })
    .catch(() => {})
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

.user-info {
  display: flex;
  align-items: center;
  gap: 8px;
}

.user-name {
  font-weight: 500;
  color: #1f2937;
}

.email-text {
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

.user-card {
  background: #fff;
  border: 1px solid #e5e7eb;
  border-radius: 10px;
  overflow: hidden;
}

.user-card.disabled {
  opacity: 0.6;
}

.user-card__header {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 16px;
  background: #f9fafb;
  border-bottom: 1px solid #f3f4f6;
}

.user-card__avatar {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  background: linear-gradient(135deg, #3b82f6, #8b5cf6);
  display: flex;
  align-items: center;
  justify-content: center;
  color: #fff;
  font-size: 18px;
}

.user-card__info {
  flex: 1;
}

.user-card__name {
  font-size: 15px;
  font-weight: 600;
  color: #1f2937;
  display: flex;
  align-items: center;
  gap: 8px;
}

.user-card__email {
  font-size: 12px;
  color: #9ca3af;
  margin-top: 2px;
}

.user-card__body {
  padding: 12px 16px;
}

.user-card__meta {
  display: flex;
  align-items: center;
  gap: 16px;
}

.meta-item {
  display: flex;
  align-items: center;
  gap: 4px;
  color: #6b7280;
  font-size: 13px;
}

.user-card__footer {
  display: flex;
  border-top: 1px solid #f3f4f6;
  padding: 8px;
  gap: 4px;
}

.user-card__footer .el-button {
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
