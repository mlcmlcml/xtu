const state = () => ({ // 建议使用函数返回 state，防止 SSR 内存泄露
  userInfo: JSON.parse(localStorage.getItem('user')) || {
    id: '',
    stuId: '',
    userHead: '',
    nickName: '',
    userName: '',  
    userEmail: ''
  }
})

const mutations = {
  SET_USER(state, userInfo) {
    state.userInfo = { ...state.userInfo, ...userInfo };
    localStorage.setItem('user', JSON.stringify(state.userInfo));
  },
  CLEAR_USER(state) {
    state.userInfo = {
      id: '',
      stuId: '',
      userHead: '',
      nickName: '',
      userName: '', 
      userEmail: ''
    };
    localStorage.removeItem('user');
  }
}

const actions = {
  login({ commit }, userInfo) {
    commit('SET_USER', userInfo);
    return Promise.resolve();
  },
  updateUser({ commit }, updatedInfo) {
    commit('SET_USER', updatedInfo);
    return Promise.resolve();
  },
  logout({ commit }) {
    commit('CLEAR_USER');
    return Promise.resolve();
  }
}

const getters = {
  userInfo: (state) => state.userInfo,
};

export default {
  namespaced: true,
  state,
  mutations,
  actions,
  getters,
};