import { useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { api } from '../api'
import { Layout } from '../components/Layout'

export function CreatePage() {
  const [filename, setFilename] = useState('')
  const [creating, setCreating] = useState(false)
  const [error, setError] = useState<string | null>(null)
  const navigate = useNavigate()

  async function handleCreate(e: React.FormEvent) {
    e.preventDefault()
    setCreating(true)
    setError(null)
    try {
      const page = await api.createPage(filename)
      navigate(`/edit/${page.filename}`)
    } catch (e: any) {
      setError(e.message)
      setCreating(false)
    }
  }

  return (
    <Layout>
      <div className="border border-gh-border rounded-lg overflow-hidden max-w-lg">
        <div className="px-4 py-3 bg-gh-subtle border-b border-gh-border">
          <h1 className="font-semibold">Create new page</h1>
        </div>
        <form onSubmit={handleCreate} className="p-6 space-y-4">
          {error && (
            <div className="border rounded p-3 text-sm" style={{ borderColor: 'var(--gh-error-border)', background: 'var(--gh-error-bg)', color: 'var(--gh-error-text)' }}>{error}</div>
          )}
          <div>
            <label className="block text-sm font-medium text-gh-text mb-1">
              Filename <span className="text-gh-muted font-normal">(e.g. my-page.md)</span>
            </label>
            <input
              type="text"
              value={filename}
              onChange={e => setFilename(e.target.value)}
              placeholder="my-page.md"
              className="w-full border border-gh-border rounded px-3 py-2 text-sm bg-gh-bg text-gh-text focus:outline-none focus:border-blue-400 focus:ring-1 focus:ring-blue-400"
              required
            />
          </div>
          <button
            type="submit"
            disabled={creating}
            className="px-4 py-2 text-sm font-medium text-white bg-gh-btn-primary rounded border border-transparent hover:opacity-90 disabled:opacity-60"
          >
            {creating ? 'Creating...' : 'Create page'}
          </button>
        </form>
      </div>
    </Layout>
  )
}
