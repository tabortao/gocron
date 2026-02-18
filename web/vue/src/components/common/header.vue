<template>
  <div class="app-header">
    <div class="header-left">
      <!-- 移动端汉堡菜单按钮 -->
      <el-button v-if="isMobile" class="menu-toggle" :icon="Menu" @click="toggleSidebar" text />
      <span class="page-title">{{ pageTitle }}</span>
    </div>

    <div class="header-right">
      <a
        href="https://github.com/tabortao/gocron/"
        target="_blank"
        class="github-link"
        title="GitHub"
      >
        <svg height="20" width="20" viewBox="0 0 16 16" fill="currentColor">
          <path
            d="M8 0C3.58 0 0 3.58 0 8c0 3.54 2.29 6.53 5.47 7.59.4.07.55-.17.55-.38 0-.19-.01-.82-.01-1.49-2.01.37-2.53-.49-2.69-.94-.09-.23-.48-.94-.82-1.13-.28-.15-.68-.52-.01-.53.63-.01 1.08.58 1.23.82.72 1.21 1.87.87 2.33.66.07-.52.28-.87.51-1.07-1.78-.2-3.64-.89-3.64-3.95 0-.87.31-1.59.82-2.15-.08-.2-.36-1.02.08-2.12 0 0 .67-.21 2.2.82.64-.18 1.32-.27 2-.27.68 0 1.36.09 2 .27 1.53-1.04 2.2-.82 2.2-.82.44 1.1.16 1.92.08 2.12.51.56.82 1.27.82 2.15 0 3.07-1.87 3.75-3.65 3.95.29.25.54.73.54 1.48 0 1.07-.01 1.93-.01 2.2 0 .21.15.46.55.38A8.013 8.013 0 0016 8c0-4.42-3.58-8-8-8z"
          ></path>
        </svg>
      </a>

      <el-dropdown v-if="userStore.isLogin" trigger="click" class="user-dropdown">
        <span class="user-info">
          <el-icon><User /></el-icon>
          <span class="username">{{ userStore.username }}</span>
          <el-icon class="hidden-mobile"><ArrowDown /></el-icon>
        </span>
        <template #dropdown>
          <el-dropdown-menu>
            <el-dropdown-item @click="$router.push('/user/edit-my-password')">
              <el-icon><Lock /></el-icon>
              {{ t('nav.changePassword') }}
            </el-dropdown-item>
            <el-dropdown-item @click="$router.push('/user/two-factor')">
              <el-icon><Key /></el-icon>
              {{ t('nav.twoFactor') }}
            </el-dropdown-item>
            <el-dropdown-item divided @click="logout">
              <el-icon><SwitchButton /></el-icon>
              {{ t('nav.logout') }}
            </el-dropdown-item>
          </el-dropdown-menu>
        </template>
      </el-dropdown>
    </div>
  </div>
</template>

<script setup>
import { computed, inject } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { useUserStore } from '../../stores/user'
import { ArrowDown, User, Lock, Key, SwitchButton, Menu } from '@element-plus/icons-vue'

const emit = defineEmits(['toggle-sidebar'])

const { t } = useI18n()
const route = useRoute()
const router = useRouter()
const userStore = useUserStore()
const isMobile = inject('isMobile', { value: false })

const pageTitle = computed(() => {
  const path = route.path
  if (path.startsWith('/task/log')) return t('task.log')
  if (path.startsWith('/task')) return t('nav.taskManage')
  if (path.startsWith('/statistics')) return t('nav.statistics')
  if (path.startsWith('/host')) return t('nav.taskNode')
  if (path.startsWith('/user')) return t('nav.userManage')
  if (path.startsWith('/system')) return t('nav.systemManage')
  return 'gocron'
})

const toggleSidebar = () => {
  emit('toggle-sidebar')
}

const logout = () => {
  userStore.logout()
  router.push('/user/login').then(() => {
    window.location.reload()
  })
}
</script>

<style scoped>
.app-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  height: 60px;
  padding: 0 20px;
  background-color: #fff;
  border-bottom: 1px solid #e4e7ed;
  box-shadow: 0 1px 4px rgba(0, 0, 0, 0.08);
}

.header-left {
  flex: 1;
  display: flex;
  align-items: center;
  gap: 8px;
  min-width: 0;
}

.menu-toggle {
  padding: 8px;
  font-size: 20px;
  color: #606266;
  flex-shrink: 0;
}

.menu-toggle:hover {
  color: #409eff;
  background-color: #f5f7fa;
}

.page-title {
  font-size: 18px;
  font-weight: 600;
  color: #303133;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.header-right {
  display: flex;
  align-items: center;
  gap: 20px;
  flex-shrink: 0;
}

.github-link {
  display: flex;
  align-items: center;
  justify-content: center;
  color: #606266;
  text-decoration: none;
  transition: all 0.3s;
  padding: 8px;
  border-radius: 4px;
}

.github-link:hover {
  color: #409eff;
  background-color: #f5f7fa;
}

.user-dropdown {
  cursor: pointer;
}

.user-info {
  display: flex;
  align-items: center;
  gap: 8px;
  color: #606266;
  padding: 8px 12px;
  border-radius: 4px;
  transition: all 0.3s;
}

.user-info:hover {
  background-color: #f5f7fa;
  color: #409eff;
}

.username {
  max-width: 120px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

:deep(.el-dropdown-menu__item) {
  display: flex;
  align-items: center;
  gap: 8px;
}

@media screen and (max-width: 768px) {
  .app-header {
    height: 56px;
    padding: 0 12px;
  }

  .page-title {
    font-size: 16px;
  }

  .header-right {
    gap: 12px;
  }

  .github-link {
    padding: 6px;
  }

  .user-info {
    padding: 6px 8px;
  }

  .username {
    max-width: 80px;
    font-size: 14px;
  }
}

@media screen and (max-width: 480px) {
  .app-header {
    padding: 0 8px;
  }

  .page-title {
    font-size: 15px;
  }

  .header-right {
    gap: 8px;
  }

  .username {
    max-width: 60px;
  }
}
</style>
