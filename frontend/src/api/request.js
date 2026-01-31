import axios from 'axios'
import { ElMessage } from 'element-plus' // 1. 修改引入包名和组件名
import store from '@/store'

const service = axios.create({
  baseURL: 'http://127.0.0.1:8001', 
  timeout: 5000 
})

service.interceptors.request.use(
  config => {
    if (store.getters.token) {
      config.headers['X-Token'] = store.getters.token || ''
    }
    return config
  },
  error => {
    return Promise.reject(error)
  }
)

service.interceptors.response.use(
  response => {
    const res = response.data
    if (res.code !== 20000) {
      let message = res.message || 'Error'
      
      if (res.code == 20002) {
        store.dispatch("dialog/setlogin", true);
      }

      // 2. 如果你需要开启错误提示，解除下方注释并使用 ElMessage
      /*
      ElMessage({
        message: message,
        type: 'error',
        duration: 3.6 * 1000
      })
      */

      return Promise.reject(new Error(res.message || 'Error'))
    } else {
      return res
    }
  },
  error => {
    // 3. 修改拦截器中的错误提示
    ElMessage({
      message: error.message,
      type: 'error',
      duration: 5 * 1000
    })
    return Promise.reject(error)
  }
)

export default service