// src/utils/axios-config.ts
import axios from 'axios'

// 设置基础URL
axios.defaults.baseURL = import.meta.env.VITE_API_BASE_URL

// 可选：添加请求拦截器
axios.interceptors.request.use(
    config => {
        // 在发送请求前做些什么
        return config
    },
    error => {
        return Promise.reject(error)
    }
)

// 可选：添加响应拦截器
axios.interceptors.response.use(
    response => {
        // 对响应数据做点什么
        return response
    },
    error => {
        // 处理错误响应
        return Promise.reject(error)
    }
)

export default axios