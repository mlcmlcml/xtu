import subjectApi from '@/api/subject'

const state = () => ({
  subjects: []
})

const getters = {
  subjectEnumFormat: (state) => (key) => {
    for (let item of state.subjects) {
      if (item.id === key) {
        return item.name + ' ( ' + item.levelName + ' )'
      }
    }
    return null
  }
}

const actions = {
  initSubject ({ commit }) {
    subjectApi.list().then(re => {
      commit('setSubjects', re.response)
    })
  }
}

const mutations = {
  setSubjects: (state, subjects) => {
    state.subjects = subjects
  }
}

export default {
  namespaced: true,
  state,
  getters,
  actions,
  mutations
}