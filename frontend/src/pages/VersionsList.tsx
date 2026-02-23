import { useEffect, useState } from 'react'
import { useParams, Link } from 'react-router-dom'
import { api, type PageInfo } from '../api'
import { Layout } from '../components/Layout'

export function VersionsList() {
  const { page } = useParams<{ page: string }>()
  const filename = page ? (page.endsWith('.md') ? page : page + '.md') : ''
  const [versions, setVersions] = useState<PageInfo[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    if (!filename) return
    api.listVersions(filename)
      .then(setVersions)
      .catch(e => setError(e.message))
      .finally(() => setLoading(false))
  }, [filename])

  return (
    <Layout>
      <div className="border border-gh-border rounded-lg overflow-hidden">
        <div className="px-4 py-3 bg-[#f6f8fa] border-b border-gh-border flex items-center gap-3">
          <Link to={`/${filename}`} className="text-gh-muted hover:text-gh-link text-sm">← {filename}</Link>
          <h1 className="font-semibold">Version history</h1>
        </div>
        {loading && <p className="p-4 text-gh-muted text-sm">Loading...</p>}
        {error && <p className="p-4 text-red-600 text-sm">{error}</p>}
        {!loading && !error && versions.length === 0 && (
          <p className="p-4 text-gh-muted text-sm">No versions available.</p>
        )}
        <ul>
          {versions.map((v, i) => (
            <li key={v.filename} className={`flex items-center justify-between px-4 py-3 ${i !== 0 ? 'border-t border-gh-border' : ''}`}>
              <Link
                to={`/${filename}/version/${encodeURIComponent(v.filename)}`}
                className="text-gh-link hover:underline font-medium text-sm"
              >
                {new Date(v.modTime).toLocaleString()}
              </Link>
              <span className="text-xs text-gh-muted">{v.filename}</span>
            </li>
          ))}
        </ul>
      </div>
    </Layout>
  )
}
