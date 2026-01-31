// src/store/modules/enumItem.js

const state = () => ({
  items: [], // 存储枚举项
})

const mutations = {
  SET_ITEMS(state, items) {
    state.items = items; // 更新枚举项
  },
}

const actions = {
  fetchItems({ commit }) {
    // 模拟异步获取枚举项
    const items = [
      { id: 1, name: '选项1' },
      { id: 2, name: '选项2' },
    ]; 
    commit('SET_ITEMS', items); // 提交变更
  },
}

const getters = {
  items: (state) => state.items, // 获取枚举项
}

export default {
  namespaced: true, // 补齐了命名空间，防止与其他模块冲突
  state,
  mutations,
  actions,
  getters,
}