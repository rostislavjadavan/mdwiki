import { useEffect, useState } from 'react'
import { useParams, useNavigate, Link } from 'react-router-dom'
import { api, type Page } from '../api'
import { Layout } from '../components/Layout'
import { MarkdownContent } from '../components/MarkdownContent'

export function VersionView() {
  const { page, ver } = useParams<{ page: string; ver: string }>()
  const filename = page ? (page.endsWith('.md') ? page : page + '.md') : ''
  const [data, setData] = useState<Page | null>(null)
  const [error, setError] = useState<string | null>(null)
  const navigate = useNavigate()

  useEffect(() => {
    if (!ver) return
    api.getVersion(decodeURIComponent(ver))
      .then(setData)
      .catch(e => setError(e.message))
  }, [ver])

  async function handleRestore() {
    if (!ver) return
    try {
      await api.restoreVersion(decodeURIComponent(ver))
      navigate(`/${filename}`)
    } catch (e: any) {
      alert(e.message)
    }
  }

  if (error) return (
    <Layout>
      <div className="rounded-lg p-6" style={{ border: '1px solid var(--gh-error-border)', background: 'var(--gh-error-bg)', color: 'var(--gh-error-text)' }}>
        <p>{error}</p>
        <Link to={`/${filename}/version`} className="text-gh-link mt-2 inline-block">Back to history</Link>
      </div>
    </Layout>
  )

  if (!data) return <Layout><p className="text-gh-muted">Loading...</p></Layout>

  return (
    <Layout>
      <div className="border border-gh-border rounded-lg overflow-hidden">
        <div className="flex items-center justify-between px-4 py-3 bg-gh-subtle border-b border-gh-border">
          <div className="flex items-center gap-3">
            <Link to={`/${filename}/version`} className="text-gh-muted hover:text-gh-link text-sm">← History</Link>
            <h1 className="font-semibold text-sm">{data.filename}</h1>
            <span className="text-xs rounded px-2 py-0.5" style={{ background: 'var(--gh-badge-blue-bg)', color: 'var(--gh-badge-blue-text)', border: '1px solid var(--gh-badge-blue-border)' }}>version</span>
          </div>
          <button
            onClick={handleRestore}
            className="text-xs text-white bg-gh-btn-primary rounded px-3 py-1 font-medium hover:opacity-90"
          >
            Restore this version
          </button>
        </div>
        <div className="p-6">
          <MarkdownContent html={data.content} />
        </div>
      </div>
    </Layout>
  )
}
