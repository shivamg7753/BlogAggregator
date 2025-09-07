import { createContext, useContext, useEffect, useMemo, useState } from 'react'
import { getToken, setToken as persistToken, getUser as readUser, setUser as persistUser } from '../auth.js'

const AuthContext = createContext(null)

export function AuthProvider({ children }) {
  const [token, setToken] = useState(() => getToken() || null)
  const [user, setUser] = useState(() => readUser())

  useEffect(() => { persistToken(token || '') }, [token])
  useEffect(() => { persistUser(user || null) }, [user])

  const value = useMemo(() => ({
    token,
    user,
    isAuthenticated: Boolean(token),
    login: ({ token: t, user: u }) => { setToken(t); setUser(u) },
    logout: () => { setToken(null); setUser(null) },
  }), [token, user])

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>
}

export function useAuth() {
  const ctx = useContext(AuthContext)
  if (!ctx) throw new Error('useAuth must be used within AuthProvider')
  return ctx
}


