# 前端优化方案

## 🚨 紧急修复（立即处理）

### 1. 修复 httpClient.js 错误
```javascript
// 将 ElElMessage 改为 ElMessage
- ElElMessage.error()
+ ElMessage.error()
```

### 2. 完成过滤器迁移
所有组件中的过滤器需要改为方法调用

## 🎯 核心优化（推荐）

### 1. 替换 Vuex → Pinia
**原因：**
- Pinia 是 Vue 官方推荐
- 更好的 TypeScript 支持
- 更简洁的 API
- 更小的包体积

**改动：**
```bash
yarn add pinia
yarn remove vuex
```

### 2. 重构 HTTP 客户端
**改为 Promise/async-await：**
```javascript
// 之前
taskService.list(params, (data) => {})

// 之后
const data = await taskService.list(params)
```

### 3. Element Plus 按需引入
**减少 ~600KB 体积：**
```javascript
// vite.config.js
import AutoImport from 'unplugin-auto-import/vite'
import Components from 'unplugin-vue-components/vite'
import { ElementPlusResolver } from 'unplugin-vue-components/resolvers'
```

### 4. 添加 TypeScript
**提升代码质量和开发体验**

### 5. 改进存储方案
```javascript
// 使用 pinia-plugin-persistedstate
// 支持加密、过期时间、类型安全
```

## 📦 依赖更新

### 需要添加
```json
{
  "pinia": "^2.1.7",
  "pinia-plugin-persistedstate": "^3.2.1",
  "unplugin-auto-import": "^0.17.0",
  "unplugin-vue-components": "^0.26.0",
  "@vueuse/core": "^10.7.0"
}
```

### 需要移除
```json
{
  "vuex": "^4.1.0",
  "qs": "^6.11.0"  // axios 内置支持
}
```

## 🔧 工具链优化

### 1. 添加代码质量工具
```json
{
  "eslint": "^8.56.0",
  "prettier": "^3.1.1",
  "@vue/eslint-config-prettier": "^9.0.0"
}
```

### 2. 添加 Git Hooks
```bash
yarn add -D husky lint-staged
```

### 3. 添加环境变量
```javascript
// .env.development
VITE_API_BASE_URL=http://localhost:5920

// .env.production
VITE_API_BASE_URL=/api
```

## 🚀 性能优化

### 1. 路由懒加载
```javascript
const TaskList = () => import('../pages/task/list.vue')
```

### 2. 图片优化
```bash
yarn add -D vite-plugin-imagemin
```

### 3. 组件缓存
```vue
<router-view v-slot="{ Component }">
  <keep-alive>
    <component :is="Component" />
  </keep-alive>
</router-view>
```

### 4. 虚拟滚动
对于长列表使用 vue-virtual-scroller

## 🔒 安全优化

### 1. Token 存储改进
```javascript
// 使用 httpOnly cookie（需要后端配合）
// 或使用 sessionStorage + 刷新令牌机制
```

### 2. 请求签名
```javascript
// 添加请求签名防止篡改
```

### 3. XSS 防护
```javascript
// 使用 DOMPurify 清理 HTML
import DOMPurify from 'dompurify'
```

## 📊 优化效果预估

| 项目 | 优化前 | 优化后 | 提升 |
|------|--------|--------|------|
| 首屏加载 | ~1.5MB | ~600KB | 60% |
| 构建时间 | 2.7s | 1.5s | 44% |
| 代码质量 | 无检查 | ESLint+TS | - |
| 类型安全 | 无 | TypeScript | - |
| 状态管理 | Vuex | Pinia | 更简洁 |

## 🎯 实施优先级

### P0 - 立即修复
1. ✅ 修复 httpClient.js 错误
2. ✅ 完成过滤器迁移

### P1 - 本周完成
1. Element Plus 按需引入
2. 重构 HTTP 客户端为 async/await
3. 添加 ESLint + Prettier

### P2 - 本月完成
1. Vuex → Pinia
2. 添加 TypeScript
3. 路由懒加载
4. 添加单元测试

### P3 - 长期优化
1. 组件重构为 Composition API
2. 添加 E2E 测试
3. 性能监控
4. 安全加固

## 📝 注意事项

1. 每个优化都应该有对应的测试
2. 逐步迁移，不要一次性改动太大
3. 保持向后兼容
4. 做好代码审查
5. 更新文档

## 🔗 参考资源

- [Pinia 文档](https://pinia.vuejs.org/)
- [Vite 性能优化](https://vitejs.dev/guide/performance.html)
- [Vue 3 最佳实践](https://vuejs.org/guide/best-practices/)
- [Element Plus 按需引入](https://element-plus.org/zh-CN/guide/quickstart.html#%E6%8C%89%E9%9C%80%E5%AF%BC%E5%85%A5)
