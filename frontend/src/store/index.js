// src/store/index.js
import { createStore } from 'vuex' // 改为 Vue 3 适配的 createStore
import getters from './getters'
import enumItem from './modules/enumItem'
import user from './modules/user'
import dialog from './modules/dialog'

// 注意：在 Vue 3 中不再需要 Vue.use(Vuex)

export default createStore({
  modules: {
    // 注意：namespaced 应该在每个具体的子模块（如 user.js）里定义，
    // 这里保持你原来的结构，但 Vue 3 主要是通过各模块导出的。
    user,
    enumItem,
    dialog
  },
  getters
})