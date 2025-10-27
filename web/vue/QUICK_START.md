# 快速开始

## 🚀 安装和运行

```bash
# 安装依赖
cd web/vue
yarn install

# 开发模式
yarn dev
# 访问 http://localhost:8080

# 生产构建
yarn build

# 预览构建结果
yarn preview

# 代码检查
yarn lint
```

## 📦 主要变更

### 1. 状态管理（Pinia）
```javascript
// 使用 store
import { useUserStore } from '@/stores/user'

const userStore = useUserStore()
console.log(userStore.username)
userStore.setUser({ username: 'admin' })
userStore.logout()
```

### 2. HTTP 请求（推荐新方式）
```javascript
// 新方式 - async/await
import request from '@/utils/request'

const data = await request.get('/api/tasks')
const result = await request.post('/api/task', { name: 'test' })

// 旧方式 - 回调（仍然支持）
import httpClient from '@/utils/httpClient'

httpClient.get('/api/tasks', {}, (data) => {
  console.log(data)
})
```

### 3. 组件写法（推荐）
```vue
<script setup>
import { ref, computed } from 'vue'
import { useUserStore } from '@/stores/user'

const userStore = useUserStore()
const count = ref(0)
const double = computed(() => count.value * 2)

const increment = () => {
  count.value++
}
</script>
```

### 4. 路由
```javascript
import { useRouter, useRoute } from 'vue-router'

const router = useRouter()
const route = useRoute()

router.push('/task')
console.log(route.params.id)
```

## 🎯 核心优化

| 项目 | 优化 |
|------|------|
| 状态管理 | Vuex → Pinia |
| 组件库 | 全量引入 → 按需引入 |
| 路由 | 静态导入 → 懒加载 |
| 构建体积 | 1.45MB → 888KB (39% ↓) |
| 首屏 JS | 1.1MB → 300KB (73% ↓) |

## 📁 项目结构

```
web/vue/
├── src/
│   ├── api/              # API 接口
│   ├── components/       # 公共组件
│   ├── pages/            # 页面组件
│   ├── router/           # 路由配置
│   ├── stores/           # Pinia stores
│   ├── utils/            # 工具函数
│   ├── App.vue           # 根组件
│   └── main.js           # 入口文件
├── .env.development      # 开发环境变量
├── .env.production       # 生产环境变量
├── .eslintrc.cjs         # ESLint 配置
├── .prettierrc.json      # Prettier 配置
├── vite.config.js        # Vite 配置
└── package.json          # 依赖配置
```

## 🔧 配置说明

### 环境变量
```bash
# .env.development
VITE_API_BASE_URL=http://localhost:5920

# .env.production
VITE_API_BASE_URL=/api
```

### 自动导入
Vite 已配置自动导入：
- Vue API (ref, computed, watch...)
- Vue Router (useRouter, useRoute)
- Pinia (defineStore, storeToRefs)
- Element Plus 组件

无需手动导入即可使用！

## 📚 相关文档

- [OPTIMIZATION_COMPLETED.md](./OPTIMIZATION_COMPLETED.md) - 优化完成报告
- [OPTIMIZATION_PLAN.md](./OPTIMIZATION_PLAN.md) - 优化方案
- [README_VUE3.md](./README_VUE3.md) - Vue 3 迁移说明
- [MIGRATION_SUMMARY.md](./MIGRATION_SUMMARY.md) - 迁移总结

## ⚡ 性能提示

1. 使用 `<script setup>` 语法
2. 使用 Pinia 替代 Vuex
3. 使用 async/await 替代回调
4. 组件自动按需加载
5. 路由自动懒加载

## 🎉 开始开发

现在你可以开始开发了！项目已经过全面优化，享受更快的开发体验吧！
