import { useEffect, useState } from 'react'
import { useParams, useNavigate, Link } from 'react-router-dom'
import { api, type Page } from '../api'
import { Layout } from '../components/Layout'
import { MarkdownContent } from '../components/MarkdownContent'

export function PageView() {
  const { page } = useParams()
  const filename = page ? (page.endsWith('.md') ? page : page + '.md') : 'home.md'
  const [data, setData] = useState<Page | null>(null)
  const [error, setError] = useState<string | null>(null)
  const navigate = useNavigate()

  useEffect(() => {
    api.getPage(filename)
      .then(setData)
      .catch(e => setError(e.message))
  }, [filename])

  async function handleDelete() {
    if (!data) return
    if (!confirm(`Delete "${data.filename}"?`)) return
    try {
      await api.deletePage(data.filename)
      navigate('/list')
    } catch (e: any) {
      alert(e.message)
    }
  }

  if (error) return (
    <Layout>
      <div className="rounded-lg p-6" style={{ borderColor: 'var(--gh-error-border)', background: 'var(--gh-error-bg)', color: 'var(--gh-error-text)', border: '1px solid var(--gh-error-border)' }}>
        <h1 className="text-lg font-semibold mb-2">Page not found</h1>
        <p>{error}</p>
        <Link to="/create" className="text-gh-link mt-4 inline-block">Create a new page</Link>
      </div>
    </Layout>
  )

  if (!data) return <Layout><p className="text-gh-muted">Loading...</p></Layout>

  return (
    <Layout>
      <div className="border border-gh-border rounded-lg overflow-hidden">
        <div className="flex items-center justify-between px-4 py-3 bg-gh-subtle border-b border-gh-border">
          <h1 className="font-semibold text-gh-text">{data.filename}</h1>
          <div className="flex items-center gap-2">
            <Link
              to={`/${data.filename}/version`}
              className="text-xs text-gh-muted hover:text-gh-link border border-gh-border rounded px-2 py-1"
            >
              History
            </Link>
            <Link
              to={`/edit/${data.filename}`}
              className="text-xs text-white bg-gh-btn-primary hover:opacity-90 border border-transparent rounded px-3 py-1 font-medium"
            >
              Edit
            </Link>
            <button
              onClick={handleDelete}
              className="text-xs text-white bg-gh-btn-danger hover:opacity-90 border border-transparent rounded px-3 py-1 font-medium"
            >
              Delete
            </button>
          </div>
        </div>
        <div className="p-6">
          <MarkdownContent html={data.content} />
        </div>
      </div>
    </Layout>
  )
}
