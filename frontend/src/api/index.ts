import axios from 'axios'
import { useAuthStore } from '@/stores/authStore'

export interface ApiResponse<T = any> {
  code: number
  message?: string
  data?: T
}

// 创建 axios 实例
const api = axios.create({
  baseURL: '/api/v1',
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json',
  },
})

// 请求拦截器
api.interceptors.request.use(
  (config) => {
    // 加载 API Key
    const store = useAuthStore()
    const { apiKey } = store
    if (apiKey) {
      config.headers['Authorization'] = `ApiKey ${apiKey}`
    }
    console.log('Request:', config.method?.toUpperCase(), config.url)

    return config
  },
  (error) => {
    return Promise.reject(error)
  },
)

// 响应拦截器
api.interceptors.response.use(
  (response) => {
    return response
  },
  (error) => {
    console.error('API Error:', error.response?.data || error.message)
    return Promise.reject(error)
  },
)

// Navigation API 接口
export const navApi = {
  /**
   * 获取所有导航数据
   * @returns {Promise<ApiResponse<CategoryData[]>>}
   */
  getNavs: () => {
    return api.get<ApiResponse<CategoryData[]>>('/navs')
  },

  // ===== 分类管理 =====

  /**
   * 创建分类（需要认证）
   * @param {Category} category - 分类数据对象
   * @returns {Promise<ApiResponse<CategoryData>>}
   */
  createCategory: (category: Category) => {
    return api.post<ApiResponse<CategoryData>>('/categories', category)
  },

  /**
   * 更新分类
   * @param {string} cid - 分类ID
   * @param {Category} category - 分类数据对象
   * @returns {Promise<ApiResponse>}
   */
  updateCategory: (cid: string, category: Category) => {
    return api.put<ApiResponse>(`/categories/${cid}`, category)
  },

  /**
   * 删除分类（需要认证）
   * @param {string} cid - 分类ID
   * @returns {Promise<ApiResponse>}
   */
  deleteCategory: (cid: string) => {
    return api.delete<ApiResponse>(`/categories/${cid}`)
  },

  // ===== 分组管理 =====

  /**
   * 创建分组（需要认证）
   * @param {string} cid - 分类ID
   * @param {Group} group - 分组数据对象
   * @returns {Promise<ApiResponse>}
   */
  createGroup: (cid: string, group: Group) => {
    return api.post<ApiResponse>(`/categories/${cid}/groups`, group)
  },

  /**
   * 更新分组
   * @param {string} cid - 分类ID
   * @param {string} gid - 分组ID
   * @param {Group} group - 分组数据对象
   * @returns {Promise<ApiResponse>}
   */
  updateGroup: (cid: string, gid: string, group: Group) => {
    return api.put<ApiResponse>(`/categories/${cid}/groups/${gid}`, group)
  },

  /**
   * 删除分组（需要认证）
   * @param {string} cid - 分类ID
   * @param {string} gid - 分组ID
   * @returns {Promise<ApiResponse>}
   */
  deleteGroup: (cid: string, gid: string) => {
    return api.delete<ApiResponse>(`/categories/${cid}/groups/${gid}`)
  },

  // ===== 收藏项管理 =====

  /**
   * 创建收藏项（需要认证）
   * @param {string} cid - 分类ID
   * @param {string} gid - 分组ID
   * @param {Collection} collection - 收藏项数据对象
   * @returns {Promise<ApiResponse>}
   */
  createCollection: (cid: string, gid: string, collection: Collection) => {
    return api.post<ApiResponse>(`/categories/${cid}/groups/${gid}/collections`, collection)
  },

  /**
   * 更新收藏项
   * @param {string} cid - 分类ID
   * @param {string} gid - 分组ID
   * @param {string} oldLink - 原始链接（用于定位要更新的收藏项）
   * @param {Collection} collection - 新的收藏项数据对象
   * @returns {Promise<ApiResponse>}
   */
  updateCollection: (cid: string, gid: string, oldLink: string, collection: Collection) => {
    return api.put<ApiResponse>(`/categories/${cid}/groups/${gid}/collections`, {
      old_link: oldLink,
      data: collection,
    })
  },

  /**
   * 删除收藏项（需要认证）
   * @param {string} cid - 分类ID
   * @param {string} gid - 分组ID
   * @param {string} link - 收藏项链接
   * @returns {Promise<ApiResponse>}
   */
  deleteCollection: (cid: string, gid: string, link: string) => {
    return api.delete<ApiResponse>(`/categories/${cid}/groups/${gid}/collections`, {
      params: { link },
    })
  },
}

// 默认导出
export default api
