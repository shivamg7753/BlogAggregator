import axios from 'axios'

const api = axios.create({
  baseURL: '/api',
})

api.interceptors.request.use((config) => {
  const token = localStorage.getItem('token')
  if (token) {
    config.headers['Authorization'] = `Bearer ${token}`
  }
  return config
})

export const login = async (username, password) => {
  const { data } = await api.post('/login', { username, password })
  return data
}

export const register = async ({ username, email, password }) => {
  const { data } = await api.post('/users/register', { username, email, password })
  return data
}

export const listFeeds = async () => {
  const { data } = await api.get('/feeds')
  return data
}

export const createFeed = async ({ title, url }) => {
  const { data } = await api.post('/feeds', { title, url })
  return data
}

export const refreshFeed = async (feed_id) => {
  const { data } = await api.post('/feeds/refresh', { feed_id })
  return data
}

export const listPosts = async () => {
  const { data } = await api.get('/posts')
  return data
}

export const subscribe = async ({ user_id, feed_id }) => {
  const { data } = await api.post('/subscriptions', { user_id, feed_id })
  return data
}

export const userFeed = async ({ user_id, page = 1, limit = 10 }) => {
  const { data } = await api.get(`/users/${user_id}/feed?page=${page}&limit=${limit}`)
  return data
}

export default api


