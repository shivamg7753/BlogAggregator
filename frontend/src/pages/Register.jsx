import { useState } from 'react'
import { register } from '../api.js'
import { Link, useNavigate } from 'react-router-dom'

export default function Register() {
  const [username, setUsername] = useState('')
  const [email, setEmail] = useState('')
  const [password, setPassword] = useState('')
  const [error, setError] = useState('')
  const [ok, setOk] = useState('')
  const navigate = useNavigate()

  const onSubmit = async (e) => {
    e.preventDefault()
    setError('')
    setOk('')
    try {
      await register({ username, email, password })
      setOk('Registered! You can now login.')
      setTimeout(() => navigate('/login'), 600)
    } catch (err) {
      setError(err?.response?.data?.error || 'Registration failed')
    }
  }

  return (
    <div style={{ maxWidth: 420, margin: '24px auto' }}>
      <h2>Register</h2>
      <form onSubmit={onSubmit} style={{ display: 'grid', gap: 12 }}>
        <input placeholder="Username" value={username} onChange={(e) => setUsername(e.target.value)} />
        <input placeholder="Email" value={email} onChange={(e) => setEmail(e.target.value)} />
        <input type="password" placeholder="Password" value={password} onChange={(e) => setPassword(e.target.value)} />
        <button type="submit">Create Account</button>
      </form>
      {ok && <p style={{ color: 'green' }}>{ok}</p>}
      {error && <p style={{ color: 'red' }}>{error}</p>}
      <p>Have an account? <Link to="/login">Login</Link></p>
    </div>
  )
}


