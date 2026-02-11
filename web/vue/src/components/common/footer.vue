<template>
  <div class="footer-bar">
    © {{ currentYear }}
    <a href="https://github.com/tabortao/gocron" target="_blank" rel="noopener noreferrer">
      gocronx-team/gocron
    </a>
  </div>
</template>

<script>
import { computed } from 'vue'

export default {
  name: 'app-footer',
  setup() {
    const currentYear = computed(() => new Date().getFullYear())
    return {
      currentYear
    }
  }
}
</script>

<style scoped>
.footer-bar {
  padding: 12px 0;
  background-color: rgba(245, 247, 250, 0.95);
  border-top: 1px solid #e4e7ed;
  text-align: center;
  font-size: 13px;
  color: #606266;
  backdrop-filter: blur(8px);
}

.footer-bar a {
  color: #409eff;
  text-decoration: none;
  transition: color 0.3s ease;
  font-weight: 500;
  margin-left: 4px;
}

.footer-bar a:hover {
  color: #66b1ff;
  text-decoration: underline;
}
</style>
