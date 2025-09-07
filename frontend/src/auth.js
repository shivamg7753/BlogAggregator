export const getToken = () => localStorage.getItem('token')

export const setToken = (token) => {
  if (token) localStorage.setItem('token', token)
  else localStorage.removeItem('token')
}

export const isAuthenticated = () => Boolean(getToken())

export const getUser = () => {
  const raw = localStorage.getItem('user')
  try { return raw ? JSON.parse(raw) : null } catch { return null }
}

export const setUser = (user) => {
  if (user) localStorage.setItem('user', JSON.stringify(user))
  else localStorage.removeItem('user')
}


