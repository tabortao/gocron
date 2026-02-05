<template>
  <el-container style="height: 100vh">
    <app-sidebar v-if="userStore.isLogin"></app-sidebar>
    <el-container style="flex-direction: column">
      <el-header v-if="userStore.isLogin" height="60px">
        <app-header></app-header>
      </el-header>
      <el-main style="padding: 0; overflow-y: auto">
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
import { computed, onMounted } from 'vue'
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
  overflow: hidden;
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
  padding: 20px;
  margin: 0;
  background-color: #f5f7fa;
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
