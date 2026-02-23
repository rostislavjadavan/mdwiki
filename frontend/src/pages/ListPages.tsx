import { useEffect, useState } from 'react'
import { Link } from 'react-router-dom'
import { api, type PageInfo } from '../api'
import { Layout } from '../components/Layout'

export function ListPages() {
  const [pages, setPages] = useState<PageInfo[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    api.listPages()
      .then(setPages)
      .catch(e => setError(e.message))
      .finally(() => setLoading(false))
  }, [])

  return (
    <Layout>
      <div className="border border-gh-border rounded-lg overflow-hidden">
        <div className="px-4 py-3 bg-gh-subtle border-b border-gh-border flex items-center justify-between">
          <h1 className="font-semibold">All pages</h1>
          <Link to="/create" className="text-xs text-white bg-gh-btn-primary rounded px-3 py-1 font-medium hover:opacity-90">
            New page
          </Link>
        </div>
        {loading && <p className="p-4 text-gh-muted text-sm">Loading...</p>}
        {error && <p className="p-4 text-sm" style={{ color: 'var(--gh-error-text)' }}>{error}</p>}
        {!loading && !error && pages.length === 0 && (
          <p className="p-4 text-gh-muted text-sm">No pages yet.</p>
        )}
        <ul>
          {pages.map((p, i) => (
            <li key={p.filename} className={`flex items-center justify-between px-4 py-3 ${i !== 0 ? 'border-t border-gh-border' : ''}`}>
              <Link to={`/${p.filename}`} className="text-gh-link hover:underline font-medium">
                {p.filename}
              </Link>
              <span className="text-xs text-gh-muted">
                {new Date(p.modTime).toLocaleDateString()}
              </span>
            </li>
          ))}
        </ul>
      </div>
    </Layout>
  )
}
