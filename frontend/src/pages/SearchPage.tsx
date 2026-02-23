import { useEffect, useState } from 'react'
import { useSearchParams, Link } from 'react-router-dom'
import { api, type SearchResult } from '../api'
import { Layout } from '../components/Layout'

export function SearchPage() {
  const [searchParams, setSearchParams] = useSearchParams()
  const [query, setQuery] = useState(searchParams.get('query') || '')
  const [result, setResult] = useState<SearchResult | null>(null)
  const [loading, setLoading] = useState(false)

  useEffect(() => {
    const q = searchParams.get('query') || ''
    setQuery(q)
    if (q) {
      setLoading(true)
      api.search(q)
        .then(setResult)
        .finally(() => setLoading(false))
    }
  }, [searchParams])

  function handleSubmit(e: React.FormEvent) {
    e.preventDefault()
    setSearchParams({ query })
  }

  return (
    <Layout>
      <div className="space-y-6">
        <div className="border border-gh-border rounded-lg overflow-hidden">
          <div className="px-4 py-3 bg-gh-subtle border-b border-gh-border">
            <h1 className="font-semibold">Search</h1>
          </div>
          <form onSubmit={handleSubmit} className="p-4 flex gap-3">
            <input
              type="text"
              value={query}
              onChange={e => setQuery(e.target.value)}
              placeholder="Search pages..."
              className="flex-1 border border-gh-border rounded px-3 py-2 text-sm bg-gh-bg text-gh-text focus:outline-none focus:border-blue-400 focus:ring-1 focus:ring-blue-400"
            />
            <button
              type="submit"
              className="px-4 py-2 text-sm font-medium text-white bg-gh-btn-primary rounded border border-transparent hover:opacity-90"
            >
              Search
            </button>
          </form>
        </div>

        {loading && <p className="text-gh-muted text-sm">Searching...</p>}

        {result && (
          <>
            {result.filenames.length > 0 && (
              <div className="border border-gh-border rounded-lg overflow-hidden">
                <div className="px-4 py-3 bg-gh-subtle border-b border-gh-border">
                  <h2 className="font-semibold text-sm">Filename matches</h2>
                </div>
                <ul>
                  {result.filenames.map((r, i) => (
                    <li key={r.filename} className={`px-4 py-3 ${i !== 0 ? 'border-t border-gh-border' : ''}`}>
                      <Link to={`/${r.filename}`} className="text-gh-link hover:underline font-medium">
                        {r.filename}
                      </Link>
                      {r.preview && (
                        <p className="text-xs text-gh-muted mt-1" dangerouslySetInnerHTML={{ __html: r.preview }} />
                      )}
                    </li>
                  ))}
                </ul>
              </div>
            )}

            {result.pageContent.length > 0 && (
              <div className="border border-gh-border rounded-lg overflow-hidden">
                <div className="px-4 py-3 bg-gh-subtle border-b border-gh-border">
                  <h2 className="font-semibold text-sm">Content matches</h2>
                </div>
                <ul>
                  {result.pageContent.map((r, i) => (
                    <li key={r.filename} className={`px-4 py-3 ${i !== 0 ? 'border-t border-gh-border' : ''}`}>
                      <Link to={`/${r.filename}`} className="text-gh-link hover:underline font-medium">
                        {r.filename}
                      </Link>
                      {r.preview && (
                        <p className="text-xs text-gh-muted mt-1 whitespace-pre-wrap" dangerouslySetInnerHTML={{ __html: r.preview }} />
                      )}
                    </li>
                  ))}
                </ul>
              </div>
            )}

            {result.filenames.length === 0 && result.pageContent.length === 0 && (
              <p className="text-gh-muted text-sm">No results found for &ldquo;{result.query}&rdquo;</p>
            )}
          </>
        )}
      </div>
    </Layout>
  )
}
