import { useEffect, useState } from 'react'
import { Link } from 'react-router-dom'
import { api, type PageInfo } from '../api'
import { Layout } from '../components/Layout'

export function TrashList() {
  const [pages, setPages] = useState<PageInfo[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  async function load() {
    setLoading(true)
    try {
      setPages(await api.listTrash())
    } catch (e: any) {
      setError(e.message)
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => { load() }, [])

  async function handleEmpty() {
    if (!confirm('Permanently delete all pages in trash?')) return
    try {
      await api.emptyTrash()
      setPages([])
    } catch (e: any) {
      alert(e.message)
    }
  }

  async function handleRestore(filename: string) {
    try {
      await api.restoreTrash(filename)
      await load()
    } catch (e: any) {
      alert(e.message)
    }
  }

  return (
    <Layout>
      <div className="border border-gh-border rounded-lg overflow-hidden">
        <div className="px-4 py-3 bg-[#f6f8fa] border-b border-gh-border flex items-center justify-between">
          <h1 className="font-semibold">Trash</h1>
          {pages.length > 0 && (
            <button
              onClick={handleEmpty}
              className="text-xs text-white bg-gh-btn-danger rounded px-3 py-1 font-medium hover:opacity-90"
            >
              Empty trash
            </button>
          )}
        </div>
        {loading && <p className="p-4 text-gh-muted text-sm">Loading...</p>}
        {error && <p className="p-4 text-red-600 text-sm">{error}</p>}
        {!loading && !error && pages.length === 0 && (
          <p className="p-4 text-gh-muted text-sm">Trash is empty.</p>
        )}
        <ul>
          {pages.map((p, i) => (
            <li key={p.filename} className={`flex items-center justify-between px-4 py-3 ${i !== 0 ? 'border-t border-gh-border' : ''}`}>
              <div>
                <Link to={`/trash/${p.filename}`} className="text-gh-link hover:underline font-medium">
                  {p.filename}
                </Link>
                <span className="text-xs text-gh-muted ml-3">
                  {new Date(p.modTime).toLocaleDateString()}
                </span>
              </div>
              <button
                onClick={() => handleRestore(p.filename)}
                className="text-xs text-white bg-gh-btn-primary rounded px-3 py-1 font-medium hover:opacity-90"
              >
                Restore
              </button>
            </li>
          ))}
        </ul>
      </div>
    </Layout>
  )
}
