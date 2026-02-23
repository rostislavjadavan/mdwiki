import { useEffect, useState } from 'react'
import { useParams, useNavigate } from 'react-router-dom'
import { api, type Page } from '../api'
import { Layout } from '../components/Layout'

export function EditPage() {
  const { page } = useParams<{ page: string }>()
  const filename = page ? (page.endsWith('.md') ? page : page + '.md') : ''
  const [data, setData] = useState<Page | null>(null)
  const [name, setName] = useState('')
  const [content, setContent] = useState('')
  const [saving, setSaving] = useState(false)
  const [error, setError] = useState<string | null>(null)
  const navigate = useNavigate()

  useEffect(() => {
    api.getPage(filename)
      .then(p => {
        setData(p)
        setName(p.filename)
        setContent(p.rawContent)
      })
      .catch(e => setError(e.message))
  }, [filename])

  async function handleSave(e: React.FormEvent) {
    e.preventDefault()
    if (!data) return
    setSaving(true)
    setError(null)
    try {
      const updated = await api.updatePage(data.filename, { filename: name, content })
      navigate(`/${updated.filename}`)
    } catch (e: any) {
      setError(e.message)
      setSaving(false)
    }
  }

  if (!data && !error) return <Layout><p className="text-gh-muted">Loading...</p></Layout>

  return (
    <Layout>
      <div className="border border-gh-border rounded-lg overflow-hidden">
        <div className="px-4 py-3 bg-[#f6f8fa] border-b border-gh-border">
          <h1 className="font-semibold">Edit page</h1>
        </div>
        <form onSubmit={handleSave} className="p-6 space-y-4">
          {error && (
            <div className="border border-red-300 bg-red-50 rounded p-3 text-red-700 text-sm">{error}</div>
          )}
          <div>
            <label className="block text-sm font-medium text-gh-text mb-1">Filename</label>
            <input
              type="text"
              value={name}
              onChange={e => setName(e.target.value)}
              className="w-full border border-gh-border rounded px-3 py-2 text-sm focus:outline-none focus:border-blue-400 focus:ring-1 focus:ring-blue-400"
              required
            />
          </div>
          <div>
            <label className="block text-sm font-medium text-gh-text mb-1">Content (Markdown)</label>
            <textarea
              value={content}
              onChange={e => setContent(e.target.value)}
              rows={20}
              className="w-full border border-gh-border rounded px-3 py-2 text-sm font-mono focus:outline-none focus:border-blue-400 focus:ring-1 focus:ring-blue-400 resize-y"
            />
          </div>
          <div className="flex gap-3">
            <button
              type="submit"
              disabled={saving}
              className="px-4 py-2 text-sm font-medium text-white bg-gh-btn-primary rounded border border-transparent hover:opacity-90 disabled:opacity-60"
            >
              {saving ? 'Saving...' : 'Save'}
            </button>
            <button
              type="button"
              onClick={() => navigate(-1)}
              className="px-4 py-2 text-sm font-medium text-gh-text bg-white rounded border border-gh-border hover:bg-[#f6f8fa]"
            >
              Cancel
            </button>
          </div>
        </form>
      </div>
    </Layout>
  )
}
