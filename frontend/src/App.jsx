import { BrowserRouter, Routes, Route, Link, Navigate } from 'react-router-dom'
import Login from './pages/Login.jsx'
import Register from './pages/Register.jsx'
import Feeds from './pages/Feeds.jsx'
import Posts from './pages/Posts.jsx'
import UserFeed from './pages/UserFeed.jsx'
import { useAuth } from './context/AuthContext.jsx'
import './App.css'

const ProtectedRoute = ({ children }) => {
  // use location state if needed to redirect back
  const { isAuthenticated } = useAuth()
  if (!isAuthenticated) return <Navigate to="/login" replace />
  return children
}

export default function App() {
  const { isAuthenticated, user, logout } = useAuth()
  return (
    <BrowserRouter>
      <div className="navbar bg-base-100 shadow-sm">
        <div className="flex-1">
          <a className="btn btn-ghost text-xl" href="/">Blog Aggregator</a>
        </div>
        <div className="flex-none gap-2">
          <ul className="menu menu-horizontal px-1">
            <li><Link to="/">Posts</Link></li>
            <li><Link to="/feeds">Feeds</Link></li>
            <li><Link to="/me">My Feed</Link></li>
          </ul>
          {!isAuthenticated ? (
            <div className="join">
              <Link className="btn btn-outline btn-sm join-item" to="/login">Login</Link>
              <Link className="btn btn-primary btn-sm join-item" to="/register">Register</Link>
            </div>
          ) : (
            <div className="dropdown dropdown-end">
              <div tabIndex={0} role="button" className="btn btn-ghost">
                <span className="badge badge-primary badge-lg">{user?.username}</span>
              </div>
              <ul tabIndex={0} className="dropdown-content menu bg-base-100 rounded-box z-[1] w-52 p-2 shadow">
                <li><button onClick={logout}>Logout</button></li>
              </ul>
            </div>
          )}
        </div>
      </div>
      <Routes>
        <Route path="/" element={<Posts />} />
        <Route path="/login" element={<Login />} />
        <Route path="/register" element={<Register />} />
        <Route path="/feeds" element={<Feeds />} />
        <Route path="/me" element={<ProtectedRoute><UserFeed /></ProtectedRoute>} />
      </Routes>
    </BrowserRouter>
  )
}
