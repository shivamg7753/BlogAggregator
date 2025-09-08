import { useEffect, useMemo, useState } from 'react'
import { listFeeds, createFeed, refreshFeed, subscribe, unsubscribe, userFeed } from '../api.js'
import { useAuth } from '../context/AuthContext.jsx'

export default function Feeds() {
  const [feeds, setFeeds] = useState([])
  const [title, setTitle] = useState('')
  const [url, setUrl] = useState('')
  const [message, setMessage] = useState('')
  const { isAuthenticated, user } = useAuth()
  const [subscribedIds, setSubscribedIds] = useState(() => new Set())

  const load = async () => {
    const data = await listFeeds()
    setFeeds(data)
  }

  useEffect(() => { load() }, [])

  useEffect(() => {
    const loadSubscribed = async () => {
      if (!isAuthenticated || !user) return
      try {
        const data = await userFeed({ user_id: user.id, page: 1, limit: 50 })
        const ids = new Set((data.posts || []).map(p => p.feed_id))
        setSubscribedIds(ids)
      } catch {}
    }
    loadSubscribed()
  }, [isAuthenticated, user])

  const onCreate = async (e) => {
    e.preventDefault()
    setMessage('')
    try {
      await createFeed({ title, url })
      setTitle(''); setUrl('')
      await load()
      setMessage('Feed created')
    } catch (e) {
      setMessage(e?.response?.data?.error || 'Failed to create feed')
    }
  }

  const onRefresh = async (id) => {
    setMessage('')
    try {
      await refreshFeed(id)
      setMessage('Refresh requested')
    } catch (e) {
      setMessage(e?.response?.data?.error || 'Failed to refresh')
    }
  }

  const onSubscribe = async (feedId) => {
    setMessage('')
    try {
      if (!isAuthenticated || !user) {
        setMessage('You are not logged in')
        return
      }
      await subscribe({ user_id: user.id, feed_id: feedId })
      setMessage('Subscribed')
      setSubscribedIds(prev => new Set([...prev, feedId]))
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

  const onUnsubscribe = async (feedId) => {
    setMessage('')
    try {
      if (!isAuthenticated || !user) {
        setMessage('You are not logged in')
        return
      }
      await unsubscribe({ user_id: user.id, feed_id: feedId })
      setMessage('Unsubscribed')
      setSubscribedIds(prev => { const s = new Set(prev); s.delete(feedId); return s })
    } catch (e) {
      setMessage(e?.response?.data?.error || 'Failed to unsubscribe')
    }
  }

  return (
    <div className="container mx-auto p-4">
      <div className="card bg-base-100 shadow-sm">
        <div className="card-body">
          <h2 className="card-title">Feeds</h2>
          <form onSubmit={onCreate} className="flex gap-2 mb-4">
            <input className="input input-bordered" placeholder="Title" value={title} onChange={(e) => setTitle(e.target.value)} />
            <input className="input input-bordered flex-1" placeholder="URL" value={url} onChange={(e) => setUrl(e.target.value)} />
            <button className="btn btn-primary" type="submit">Add</button>
          </form>
          {message && <div className="alert alert-info py-2 px-3">{message}</div>}
          <ul className="menu">
            {feeds.map((f) => (
              <li key={f.id}>
                <div className="flex items-center gap-2">
                  <div className="min-w-0">
                    <div className="font-semibold truncate">{f.title}</div>
                    <a className="text-xs text-primary" href={f.url} target="_blank">{f.url}</a>
                  </div>
                  <div className="ml-auto flex gap-2">
                    <button className="btn btn-sm" onClick={() => onRefresh(f.id)}>Refresh</button>
                    {subscribedIds.has(f.id) ? (
                      <button className="btn btn-outline btn-sm" onClick={() => onUnsubscribe(f.id)}>Unsubscribe</button>
                    ) : (
                      <button className="btn btn-primary btn-sm" onClick={() => onSubscribe(f.id)}>Subscribe</button>
                    )}
                  </div>
                </div>
              </li>
            ))}
          </ul>
        </div>
      </div>
    </div>
  )
}


