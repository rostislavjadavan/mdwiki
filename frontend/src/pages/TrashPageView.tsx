import { useEffect, useState } from 'react'
import { useParams, useNavigate, Link } from 'react-router-dom'
import { api, type Page } from '../api'
import { Layout } from '../components/Layout'
import { MarkdownContent } from '../components/MarkdownContent'

export function TrashPageView() {
  const { page } = useParams<{ page: string }>()
  const [data, setData] = useState<Page | null>(null)
  const [error, setError] = useState<string | null>(null)
  const navigate = useNavigate()

  useEffect(() => {
    if (!page) return
    api.getTrashPage(page)
      .then(setData)
      .catch(e => setError(e.message))
  }, [page])

  async function handleRestore() {
    if (!data) return
    try {
      await api.restoreTrash(data.filename)
      navigate('/trash')
    } catch (e: any) {
      alert(e.message)
    }
  }

  if (error) return (
    <Layout>
      <div className="rounded-lg p-6" style={{ border: '1px solid var(--gh-error-border)', background: 'var(--gh-error-bg)', color: 'var(--gh-error-text)' }}>
        <p>{error}</p>
        <Link to="/trash" className="text-gh-link mt-2 inline-block">Back to trash</Link>
      </div>
    </Layout>
  )

  if (!data) return <Layout><p className="text-gh-muted">Loading...</p></Layout>

  return (
    <Layout>
      <div className="border border-gh-border rounded-lg overflow-hidden">
        <div className="flex items-center justify-between px-4 py-3 bg-gh-subtle border-b border-gh-border">
          <div className="flex items-center gap-3">
            <Link to="/trash" className="text-gh-muted hover:text-gh-link text-sm">← Trash</Link>
            <h1 className="font-semibold">{data.filename}</h1>
            <span className="text-xs rounded px-2 py-0.5" style={{ background: 'var(--gh-badge-yellow-bg)', color: 'var(--gh-badge-yellow-text)', border: '1px solid var(--gh-badge-yellow-border)' }}>deleted</span>
          </div>
          <button
            onClick={handleRestore}
            className="text-xs text-white bg-gh-btn-primary rounded px-3 py-1 font-medium hover:opacity-90"
          >
            Restore
          </button>
        </div>
        <div className="p-6">
          <MarkdownContent html={data.content} />
        </div>
      </div>
    </Layout>
  )
}
