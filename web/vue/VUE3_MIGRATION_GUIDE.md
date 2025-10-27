# Vue 3 迁移指南

## 已完成的迁移

1. ✅ 构建工具从 Webpack 迁移到 Vite
2. ✅ Vue 2 升级到 Vue 3
3. ✅ Vue Router 3 升级到 Vue Router 4
4. ✅ Vuex 3 升级到 Vuex 4
5. ✅ Element UI 升级到 Element Plus
6. ✅ 更新 main.js、router、store、App.vue
7. ✅ 更新 httpClient.js

## 需要手动更新的组件

所有 `.vue` 组件文件需要进行以下更新：

### 1. Element UI 组件名称变更

需要全局替换：
- `el-button` → 保持不变
- `el-form` → 保持不变
- `el-table` → 保持不变
- `el-dialog` → 保持不变
- `el-input` → 保持不变
- `el-select` → 保持不变
- `el-option` → 保持不变

但需要注意：
- `Message` → `ElMessage`
- `MessageBox` → `ElMessageBox`
- `Notification` → `ElNotification`
- `Loading` → `ElLoading`

### 2. 过滤器 (Filters) 迁移

Vue 3 移除了过滤器，需要改为方法或计算属性：

**之前 (Vue 2):**
```vue
<template>
  <div>{{ time | formatTime }}</div>
</template>
```

**之后 (Vue 3):**
```vue
<template>
  <div>{{ $filters.formatTime(time) }}</div>
</template>
```

### 3. $listeners 和 $attrs 变更

Vue 3 中 `$listeners` 已被移除，所有监听器都在 `$attrs` 中。

### 4. v-model 变更

自定义组件的 v-model：
- `model` 选项被移除
- `.sync` 修饰符被移除，统一使用 `v-model:propName`

### 5. 事件 API 变更

- `$on`、`$off`、`$once` 被移除
- 使用 mitt 或其他事件总线库替代

### 6. 异步组件

**之前:**
```js
const AsyncComponent = () => import('./MyComponent.vue')
```

**之后:**
```js
import { defineAsyncComponent } from 'vue'
const AsyncComponent = defineAsyncComponent(() => import('./MyComponent.vue'))
```

### 7. 函数式组件

函数式组件语法完全改变，需要重写为普通组件或使用新的函数式组件语法。

## 推荐的组件迁移步骤

对于每个 `.vue` 组件：

1. 检查是否使用了过滤器，改为 `$filters.xxx()` 方法调用
2. 检查 Element UI 组件的 API 变更（特别是 Table、Form、Dialog）
3. 如果使用了 `$listeners`，改为 `$attrs`
4. 检查自定义 v-model 和 .sync 的使用
5. 测试组件功能

## 运行项目

### 开发模式
```bash
cd web/vue
yarn dev
```

### 构建生产版本
```bash
cd web/vue
yarn build
```

或使用 Makefile：
```bash
make build-vue
```

## 注意事项

1. Element Plus 的一些组件 API 有变化，需要查阅文档
2. Vue 3 的响应式系统有变化，但对于 Options API 影响较小
3. 如果要使用 Composition API，可以逐步重构组件
4. Vite 开发服务器比 Webpack 快很多
5. 生产构建也会更快

## 参考文档

- [Vue 3 迁移指南](https://v3-migration.vuejs.org/)
- [Element Plus 文档](https://element-plus.org/)
- [Vite 文档](https://vitejs.dev/)
- [Vue Router 4 迁移](https://router.vuejs.org/guide/migration/)
- [Vuex 4 迁移](https://vuex.vuejs.org/guide/migrating-to-4-0-from-3-x.html)
