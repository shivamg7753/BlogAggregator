import { useEffect, useState } from 'react'
import { userFeed, listFeeds, subscribe, createFeed } from '../api.js'
import { useAuth } from '../context/AuthContext.jsx'

export default function UserFeed() {
  const [page, setPage] = useState(1)
  const [limit] = useState(10)
  const [posts, setPosts] = useState([])
  const [feeds, setFeeds] = useState([])
  const [title, setTitle] = useState('')
  const [url, setUrl] = useState('')
  const [message, setMessage] = useState('')
  const { isAuthenticated, user } = useAuth()
  const [subscribedIds, setSubscribedIds] = useState(() => new Set())
  const [loading, setLoading] = useState(true)

  const load = async () => {
    setLoading(true)
    const data = await userFeed({ user_id: user?.id, page, limit })
    const list = data.posts || []
    setPosts(list)
    const ids = new Set(list.map(p => p.feed_id))
    setSubscribedIds(ids)
    setLoading(false)
  }

  const loadFeeds = async () => {
    const data = await listFeeds()
    setFeeds(data)
  }

  useEffect(() => { load() }, [page])
  useEffect(() => { loadFeeds() }, [])

  const onSubscribe = async (feedId) => {
    try {
      setMessage('')
      if (!isAuthenticated || !user) {
        setMessage('You are not logged in')
        return
      }
      await subscribe({ user_id: user.id, feed_id: feedId })
      setMessage('Subscribed')
      setSubscribedIds(prev => new Set([...prev, feedId]))
      await load()
    } catch (e) {
      const msg = (e?.response?.data?.error || '').toString().toLowerCase()
      if (msg.includes('unique') || msg.includes('duplicate') || e?.response?.status === 409) {
        setMessage('Already subscribed')
        setSubscribedIds(prev => new Set([...prev, feedId]))
      } else if (e?.response?.status === 401 || e?.response?.status === 400) {
        setMessage('You are not logged in')
      } else {
        setMessage(e?.response?.data?.error || 'Failed to subscribe')
      }
    }
  }

  const onAddAndSubscribe = async (e) => {
    e.preventDefault()
    try {
      setMessage('')
      const created = await createFeed({ title, url })
      setTitle(''); setUrl('')
      await loadFeeds()
      await onSubscribe(created.id)
    } catch (e) {
      setMessage(e?.response?.data?.error || 'Failed to add feed')
    }
  }

  return (
    <div className="container mx-auto p-4">
      <div className="card bg-base-100 shadow-sm">
        <div className="card-body">
          <h2 className="card-title">My Feed</h2>
          {message && <div className="alert alert-info py-2 px-3">{message}</div>}
          <div className="card bg-base-100 border mb-4">
            <div className="card-body">
              <h3 className="font-semibold">Add a feed and subscribe</h3>
              <form onSubmit={onAddAndSubscribe} className="flex gap-2">
                <input className="input input-bordered" placeholder="Title" value={title} onChange={(e) => setTitle(e.target.value)} />
                <input className="input input-bordered flex-1" placeholder="URL" value={url} onChange={(e) => setUrl(e.target.value)} />
                <button className="btn btn-primary" type="submit">Add + Subscribe</button>
              </form>
            </div>
          </div>
          <div className="card bg-base-100 border mb-4">
            <div className="card-body">
              <h3 className="font-semibold">All Feeds</h3>
              <ul className="menu">
                {feeds.map((f) => (
                  <li key={f.id}>
                    <div className="flex items-center gap-2">
                      <div className="min-w-0">
                        <div className="font-semibold truncate">{f.title}</div>
                        <a className="text-xs text-primary" href={f.url} target="_blank">{f.url}</a>
                      </div>
                      <div className="ml-auto flex gap-2">
                        <button className="btn btn-primary btn-sm" onClick={() => onSubscribe(f.id)} disabled={subscribedIds.has(f.id)}>
                          {subscribedIds.has(f.id) ? 'Subscribed' : 'Subscribe'}
                        </button>
                      </div>
                    </div>
                  </li>
                ))}
              </ul>
            </div>
          </div>
          {loading ? <div className="p-6"><span className="loading loading-spinner loading-md"></span></div> : (
          <>
            <ul className="menu">
              {posts.map((p) => (
                <li key={p.id}>
                  <a href={p.link} target="_blank" rel="noreferrer" className="text-primary font-semibold">{p.title}</a>
                  <div className="text-xs opacity-70">{new Date(p.published).toLocaleString()}</div>
                  <p className="mt-1">{p.content?.slice(0, 200)}</p>
                </li>
              ))}
            </ul>
            <div className="flex items-center gap-2 mt-2">
              <button className="btn btn-sm" onClick={() => setPage((p) => Math.max(1, p - 1))} disabled={page === 1}>Prev</button>
              <span className="text-sm">Page {page}</span>
              <button className="btn btn-sm" onClick={() => setPage((p) => p + 1)}>Next</button>
            </div>
          </>
          )}
        </div>
      </div>
    </div>
  )
}


