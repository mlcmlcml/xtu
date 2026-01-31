import { routes } from '@/router/index'

function hasPermission(roles, route) {
  if (route.meta && route.meta.roles) {
    return roles.some(role => route.meta.roles.indexOf(role) >= 0)
  } else {
    return true
  }
}

const state = () => ({
  routers: routes,
  addRouters: []
})

const mutations = {
  SET_ROUTERS: (state, routers) => {
    state.addRouters = routers
    // 注意：这里假设你的 router/index.js 导出了 routes 数组作为基础路由
    state.routers = routes.concat(routers) 
  },
  RESET_ROUTERS: (state) => {
    state.addRouters = []
    state.routers = routes
  }
}

const actions = {
  GenerateRoutes({ commit }, data) {
    return new Promise(resolve => {
      const { roles } = data
      // 这里的 asyncRoutes 需要根据你项目实际导出的异步路由表来修改
      // 如果没有定义，这里会报错。暂时保持逻辑，但请确保变量存在。
      const accessedRouters = [] // 建议此处传入你的异步路由配置
      
      commit('SET_ROUTERS', accessedRouters)
      resolve(accessedRouters)
    })
  }
}

export default {
  namespaced: true,
  state,
  mutations,
  actions
}