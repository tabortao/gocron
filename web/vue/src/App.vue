<template>
  <el-container style="height: 100vh" class="app-container">
    <!-- 桌面端侧边栏 -->
    <app-sidebar v-if="userStore.isLogin && !isMobile" class="desktop-sidebar"></app-sidebar>

    <!-- 移动端抽屉式侧边栏 -->
    <el-drawer
      v-if="userStore.isLogin && isMobile"
      v-model="sidebarVisible"
      direction="ltr"
      :with-header="false"
      :size="280"
      :z-index="2000"
      class="mobile-sidebar-drawer"
      @close="handleSidebarClose"
    >
      <app-sidebar @menu-select="handleMenuSelect"></app-sidebar>
    </el-drawer>

    <el-container style="flex-direction: column" class="main-container">
      <el-header v-if="userStore.isLogin" :height="isMobile ? '56px' : '60px'" class="app-header">
        <app-header @toggle-sidebar="toggleSidebar"></app-header>
      </el-header>
      <el-main style="padding: 0; overflow-y: auto" class="app-main">
        <div id="main-container" v-cloak>
          <el-config-provider :locale="activeLang">
            <router-view v-slot="{ Component }">
              <keep-alive>
                <component :is="Component" />
              </keep-alive>
            </router-view>
          </el-config-provider>
        </div>
      </el-main>
    </el-container>
  </el-container>
</template>

<script setup>
import { computed, onMounted, onUnmounted, ref, provide } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from './stores/user'
import installService from './api/install'
import appHeader from './components/common/header.vue'
import appSidebar from './components/common/sidebar.vue'
import { ElConfigProvider } from 'element-plus'
import zhCN from 'element-plus/es/locale/lang/zh-cn'
import en from 'element-plus/es/locale/lang/en'
import { useI18n } from 'vue-i18n'
import { availableLanguages } from './const/index'

const { locale } = useI18n()
const router = useRouter()
const userStore = useUserStore()

const isMobile = ref(false)
const sidebarVisible = ref(false)

const checkMobile = () => {
  isMobile.value = window.innerWidth <= 768
}

const activeLang = computed(() => {
  switch (locale.value) {
    case availableLanguages.enUS.value:
      return en
    case availableLanguages.zhCN.value:
      return zhCN
    default:
      return zhCN
  }
})

const toggleSidebar = () => {
  sidebarVisible.value = !sidebarVisible.value
}

const handleSidebarClose = () => {
  sidebarVisible.value = false
}

const handleMenuSelect = () => {
  if (isMobile.value) {
    sidebarVisible.value = false
  }
}

provide('isMobile', isMobile)
provide('toggleSidebar', toggleSidebar)

onMounted(() => {
  checkMobile()
  window.addEventListener('resize', checkMobile)

  installService.status(data => {
    if (!data) {
      router.push('/install')
    }
  })
})

onUnmounted(() => {
  window.removeEventListener('resize', checkMobile)
})
</script>

<style>
[v-cloak] {
  display: none !important;
}

.app-container {
  overflow: hidden;
}

.desktop-sidebar {
  flex-shrink: 0;
}

.main-container {
  flex: 1;
  min-width: 0;
  overflow: hidden;
}

.app-header {
  padding: 0;
  background-color: #fff;
  border-bottom: 1px solid var(--border-color, #e4e7ed);
  box-shadow: 0 1px 4px rgba(0, 0, 0, 0.08);
  position: relative;
  z-index: 100;
}

.app-main {
  background-color: var(--bg-color, #f6f7fb);
}

#main-container {
  height: 100%;
}

#main-container .el-container {
  height: 100%;
  background-color: transparent;
}

#main-container .el-main {
  height: auto;
  overflow-y: auto;
}

.mobile-sidebar-drawer {
  background-color: #304156;
}

.mobile-sidebar-drawer :deep(.el-drawer__body) {
  padding: 0;
  background-color: #304156;
}

.mobile-sidebar-drawer :deep(.el-overlay) {
  background-color: rgba(0, 0, 0, 0.5);
}

@media screen and (max-width: 768px) {
  .app-header {
    height: 56px !important;
  }
}

@media (prefers-color-scheme: dark) {
  .app-header {
    background-color: #1f1f1f;
    border-bottom-color: #333;
  }
}
</style>
