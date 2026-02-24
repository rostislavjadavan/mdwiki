import { useEffect, useState } from 'react'
import { useParams, useNavigate } from 'react-router-dom'
import { api, type Page } from '../api'
import { Layout } from '../components/Layout'
import CodeMirror from '@uiw/react-codemirror'
import { githubLight, githubDark } from '@uiw/codemirror-theme-github'
import { markdown } from '@codemirror/lang-markdown'

function useDarkMode() {
  const [isDark, setIsDark] = useState(() => document.documentElement.classList.contains('dark'))
  useEffect(() => {
    const obs = new MutationObserver(() => setIsDark(document.documentElement.classList.contains('dark')))
    obs.observe(document.documentElement, { attributeFilter: ['class'] })
    return () => obs.disconnect()
  }, [])
  return isDark
}

export function EditPage() {
  const { page } = useParams<{ page: string }>()
  const filename = page ? (page.endsWith('.md') ? page : page + '.md') : ''
  const [data, setData] = useState<Page | null>(null)
  const [name, setName] = useState('')
  const [content, setContent] = useState('')
  const [saving, setSaving] = useState(false)
  const [error, setError] = useState<string | null>(null)
  const navigate = useNavigate()
  const isDark = useDarkMode()

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
        <div className="px-4 py-3 bg-gh-subtle border-b border-gh-border">
          <h1 className="font-semibold">Edit page</h1>
        </div>
        <form onSubmit={handleSave} className="p-6 space-y-4">
          {error && (
            <div className="border rounded p-3 text-sm" style={{ borderColor: 'var(--gh-error-border)', background: 'var(--gh-error-bg)', color: 'var(--gh-error-text)' }}>{error}</div>
          )}
          <div>
            <label className="block text-sm font-medium text-gh-text mb-1">Filename</label>
            <input
              type="text"
              value={name}
              onChange={e => setName(e.target.value)}
              className="w-full border border-gh-border rounded px-3 py-2 text-sm bg-gh-bg text-gh-text focus:outline-none focus:border-blue-400 focus:ring-1 focus:ring-blue-400"
              required
            />
          </div>
          <div>
            <label className="block text-sm font-medium text-gh-text mb-1">Content (Markdown)</label>
            <div className="border border-gh-border rounded overflow-hidden focus-within:border-blue-400 focus-within:ring-1 focus-within:ring-blue-400">
              <CodeMirror
                value={content}
                onChange={setContent}
                extensions={[markdown()]}
                theme={isDark ? githubDark : githubLight}
                basicSetup={{ lineNumbers: false, foldGutter: false }}
                className="text-sm"
              />
            </div>
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
              className="px-4 py-2 text-sm font-medium text-gh-text bg-gh-bg rounded border border-gh-border hover:bg-gh-subtle"
            >
              Cancel
            </button>
          </div>
        </form>
      </div>
    </Layout>
  )
}
