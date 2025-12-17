<template>
  <el-container style="height: 100vh">
    <el-header v-if="userStore.isLogin">
      <app-header></app-header>
      <app-nav-menu></app-nav-menu>
    </el-header>
    <el-main style="padding: 0; display: flex; flex-direction: column; overflow: hidden">
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
</template>

<script setup>
import { onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from './stores/user'
import installService from './api/install'
import appHeader from './components/common/header.vue'
import appNavMenu from './components/common/navMenu.vue'
import { ElConfigProvider } from 'element-plus'
import zhCN from 'element-plus/es/locale/lang/zh-cn'
import en from 'element-plus/es/locale/lang/en'
import { useI18n } from 'vue-i18n'

const { locale } = useI18n()
const router = useRouter()
const userStore = useUserStore()

const activeLang = computed(() => {
  switch (locale.value) {
    case 'en-US':
      return en
    case 'zh-CN':
      return zhCN
    default:
      return zhCN
  }
})

onMounted(() => {
  installService.status(data => {
    if (!data) {
      router.push('/install')
    }
  })
})
</script>

<style>
[v-cloak] {
  display: none !important;
}
html,
body {
  margin: 0;
  padding: 0;
  height: 100%;
  overflow-x: hidden;
}
.el-header {
  padding: 0;
  margin: 0;
}
.el-container {
  padding: 0;
  margin: 0;
  width: 100%;
}
.el-main {
  padding: 0;
  margin: 0;
}
#main-container {
  flex: 1;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

#main-container .el-container {
  height: 100%;
}

#main-container .el-main {
  height: auto;
  flex: 1;
  overflow-y: auto;
  margin: 20px;
}

.el-aside .el-menu {
  height: 100%;
}
.custom-message-box {
  min-width: 420px;
}
.custom-message-box .el-message-box__message {
  font-size: 15px;
  line-height: 1.6;
}
.el-message-box__title {
  font-size: 18px;
  font-weight: 600;
}
</style>
