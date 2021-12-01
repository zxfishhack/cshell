import axios from 'axios'
import { Message } from 'element-ui'
import { removeToken } from '@/assets/js/utils/cookie'

const service = axios.create({
  baseURL: '',
  timeout: 20000 // 请求超时时间
})

service.interceptors.response.use(
  response => {
    return response.data
  },
  error => {
    if (error.response) {
      const status = error.response.status
      // 401
      if (status === 401) {
        Message.error('认证失败，请重新登录')
        removeToken()
        // TODO
        console.log('error', error)
        location.href = '#/login'
      }
      // 500 服务器内部错误
      if (status === 500) {
        Message.error('服务器内部错误')
      }
      // 400 参数错误
      if (status === 400) {
        Message.error('参数错误')
      }
    } else {
      // 请求超时
      Message.error(error.message)
    }
    return Promise.reject(error)
  }
)

export default service
