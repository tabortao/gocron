<template>
  <el-dropdown @command="handleCommand">
    <span class="language-switcher"> 🌐 {{ currentLanguage }} </span>
    <template #dropdown>
      <el-dropdown-menu>
        <el-dropdown-item
          v-for="lang in availableLanguages"
          :key="lang.value"
          :command="lang.value"
          :disabled="locale === lang.value"
        >
          {{ lang.label }}
        </el-dropdown-item>
      </el-dropdown-menu>
    </template>
  </el-dropdown>
</template>

<script setup>
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'
import { availableLanguages } from '@/const/index'

const { locale } = useI18n()

const currentLanguage = computed(() => {
  return availableLanguages[locale.value] || availableLanguages.zhCN.label
})

const handleCommand = command => {
  locale.value = command
  localStorage.setItem('locale', command)
}
</script>

<style scoped>
.language-switcher {
  cursor: pointer;
  display: flex;
  align-items: center;
  gap: 4px;
  padding: 0 12px;
  font-size: 14px;
  white-space: nowrap;
}

.language-switcher:hover {
  color: #409eff;
}
</style>
