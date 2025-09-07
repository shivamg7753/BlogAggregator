import { useState } from 'react'
import { login } from '../api.js'
import { useAuth } from '../context/AuthContext.jsx'
import { useNavigate, Link } from 'react-router-dom'

export default function Login() {
  const [username, setUsername] = useState('')
  const [password, setPassword] = useState('')
  const [error, setError] = useState('')
  const navigate = useNavigate()
  const { login: setSession } = useAuth()

  const onSubmit = async (e) => {
    e.preventDefault()
    setError('')
    try {
      const data = await login(username, password)
      setSession({ token: data.token, user: data.user })
      navigate('/me')
    } catch (err) {
      setError(err?.response?.data?.error || 'Login failed')
    }
  }

  return (
    <div style={{ maxWidth: 420, margin: '24px auto' }}>
      <h2>Login</h2>
      <form onSubmit={onSubmit} style={{ display: 'grid', gap: 12 }}>
        <input placeholder="Username" value={username} onChange={(e) => setUsername(e.target.value)} />
        <input type="password" placeholder="Password" value={password} onChange={(e) => setPassword(e.target.value)} />
        <button type="submit">Login</button>
      </form>
      {error && <p style={{ color: 'red' }}>{error}</p>}
      <p>Need an account? <Link to="/register">Register</Link></p>
    </div>
  )
}


