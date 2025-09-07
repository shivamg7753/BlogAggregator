import { useEffect, useState } from 'react'
import { listPosts } from '../api.js'

export default function Posts() {
  const [posts, setPosts] = useState([])
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    listPosts().then((data) => setPosts(data)).finally(() => setLoading(false))
  }, [])

  if (loading) return <div className="p-6"><span className="loading loading-spinner loading-md"></span></div>

  return (
    <div className="container mx-auto p-4">
      <div className="card bg-base-100 shadow-sm">
        <div className="card-body">
          <h2 className="card-title">Latest Posts</h2>
          <ul className="menu">
          {posts.map((p) => (
            <li key={p.id} className="">
              <a href={p.link} target="_blank" rel="noreferrer" className="text-primary font-semibold">{p.title}</a>
              <div className="text-xs opacity-70">{new Date(p.published).toLocaleString()}</div>
              <p className="mt-1">{p.content?.slice(0, 200)}</p>
            </li>
          ))}
          </ul>
        </div>
      </div>
    </div>
  )
}


