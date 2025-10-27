<template>
  <div v-cloak>
    <el-menu
      :default-active="currentRoute"
      mode="horizontal"
      background-color="#545c64"
      text-color="#fff"
      active-text-color="#ffd04b"
      router>
      <el-row>
        <el-col :span="2">
          <el-menu-item index="/task">任务管理</el-menu-item>
        </el-col>
        <el-col :span="2">
          <el-menu-item index="/host">任务节点</el-menu-item>
        </el-col>
        <el-col :span="2">
          <el-menu-item v-if="userStore.isAdmin" index="/user">用户管理</el-menu-item>
        </el-col>
        <el-col :span="2">
          <el-menu-item v-if="userStore.isAdmin" index="/system">系统管理</el-menu-item>
        </el-col>
        <el-col :span="16"></el-col>
        <el-col :span="2" style="float:right;">
          <el-sub-menu v-if="userStore.token" index="userStatus">
            <template #title>{{ userStore.username }}</template>
            <el-menu-item index="/user/edit-my-password">修改密码</el-menu-item>
            <el-menu-item index="/user/two-factor">双因素认证</el-menu-item>
            <el-menu-item @click="logout" index="/user/logout">退出</el-menu-item>
          </el-sub-menu>
        </el-col>
      </el-row>
    </el-menu>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useUserStore } from '../../stores/user'

const route = useRoute()
const router = useRouter()
const userStore = useUserStore()

const currentRoute = computed(() => {
  if (route.path === '/') return '/task'
  const segments = route.path.split('/')
  return `/${segments[1]}`
})

const logout = () => {
  userStore.logout()
  router.push('/')
}
</script>
