import { useState, useEffect, useRef } from 'react'
import { Link, useLocation } from 'react-router-dom'
import { api, type SearchResult } from '../api'

type Theme = 'light' | 'dark' | 'system'

interface Props {
  theme: Theme
  setTheme: (t: Theme) => void
}

function ThemeIcon({ theme }: { theme: Theme }) {
  if (theme === 'dark') {
    return (
      <svg className="w-4 h-4" fill="currentColor" viewBox="0 0 20 20">
        <path d="M17.293 13.293A8 8 0 016.707 2.707a8.001 8.001 0 1010.586 10.586z" />
      </svg>
    )
  }
  if (theme === 'light') {
    return (
      <svg className="w-4 h-4" fill="currentColor" viewBox="0 0 20 20">
        <path fillRule="evenodd" d="M10 2a1 1 0 011 1v1a1 1 0 11-2 0V3a1 1 0 011-1zm4 8a4 4 0 11-8 0 4 4 0 018 0zm-.464 4.95l.707.707a1 1 0 001.414-1.414l-.707-.707a1 1 0 00-1.414 1.414zm2.12-10.607a1 1 0 010 1.414l-.706.707a1 1 0 11-1.414-1.414l.707-.707a1 1 0 011.414 0zM17 11a1 1 0 100-2h-1a1 1 0 100 2h1zm-7 4a1 1 0 011 1v1a1 1 0 11-2 0v-1a1 1 0 011-1zM5.05 6.464A1 1 0 106.465 5.05l-.708-.707a1 1 0 00-1.414 1.414l.707.707zm1.414 8.486l-.707.707a1 1 0 01-1.414-1.414l.707-.707a1 1 0 011.414 1.414zM4 11a1 1 0 100-2H3a1 1 0 000 2h1z" clipRule="evenodd" />
      </svg>
    )
  }
  // system
  return (
    <svg className="w-4 h-4" fill="currentColor" viewBox="0 0 20 20">
      <path fillRule="evenodd" d="M3 5a2 2 0 012-2h10a2 2 0 012 2v8a2 2 0 01-2 2h-2.22l.123.489.804.804A1 1 0 0113 18H7a1 1 0 01-.707-1.707l.804-.804L7.22 15H5a2 2 0 01-2-2V5zm5.771 7H5V5h10v7H8.771z" clipRule="evenodd" />
    </svg>
  )
}

function SearchDropdown({ result, loading, query, onSelect }: {
  result: SearchResult | null
  loading: boolean
  query: string
  onSelect: () => void
}) {
  if (loading && !result) {
    return (
      <div className="absolute top-full mt-1 w-[400px] max-w-[calc(100vw-2rem)] bg-gh-bg border border-gh-border rounded-lg shadow-lg z-50 p-3">
        <p className="text-gh-muted text-sm">Searching...</p>
      </div>
    )
  }

  if (!result) return null

  const empty = result.filenames.length === 0 && result.pageContent.length === 0

  return (
    <div className="absolute top-full mt-1 w-[400px] max-w-[calc(100vw-2rem)] bg-gh-bg border border-gh-border rounded-lg shadow-lg z-50 max-h-[70vh] overflow-y-auto">
      {empty && (
        <p className="text-gh-muted text-sm p-3">No results for &ldquo;{query}&rdquo;</p>
      )}

      {result.filenames.length > 0 && (
        <div>
          <div className="px-3 py-2 bg-gh-subtle border-b border-gh-border text-xs font-semibold text-gh-muted">
            Filename matches
          </div>
          {result.filenames.map((r) => (
            <Link
              key={r.filename}
              to={`/${r.filename}`}
              onClick={onSelect}
              className="block px-3 py-2 hover:bg-gh-subtle border-b border-gh-border last:border-b-0"
            >
              <span className="text-gh-link text-sm font-medium">{r.filename}</span>
              {r.preview && (
                <p className="text-xs text-gh-muted mt-0.5 truncate" dangerouslySetInnerHTML={{ __html: r.preview }} />
              )}
            </Link>
          ))}
        </div>
      )}

      {result.pageContent.length > 0 && (
        <div>
          <div className="px-3 py-2 bg-gh-subtle border-b border-gh-border text-xs font-semibold text-gh-muted">
            Content matches
          </div>
          {result.pageContent.map((r) => (
            <Link
              key={r.filename}
              to={`/${r.filename}`}
              onClick={onSelect}
              className="block px-3 py-2 hover:bg-gh-subtle border-b border-gh-border last:border-b-0"
            >
              <span className="text-gh-link text-sm font-medium">{r.filename}</span>
              {r.preview && (
                <p className="text-xs text-gh-muted mt-0.5 line-clamp-2" dangerouslySetInnerHTML={{ __html: r.preview }} />
              )}
            </Link>
          ))}
        </div>
      )}
    </div>
  )
}

export function Navbar({ theme, setTheme }: Props) {
  const [query, setQuery] = useState('')
  const [menuOpen, setMenuOpen] = useState(false)
  const [result, setResult] = useState<SearchResult | null>(null)
  const [loading, setLoading] = useState(false)
  const [showResults, setShowResults] = useState(false)
  const abortRef = useRef<AbortController | null>(null)
  const wrapperRef = useRef<HTMLDivElement>(null)
  const location = useLocation()

  // Close dropdown on route change
  useEffect(() => {
    setShowResults(false)
    setQuery('')
    setResult(null)
  }, [location.pathname])

  // Debounced search
  useEffect(() => {
    if (!query.trim()) {
      setResult(null)
      setShowResults(false)
      return
    }

    setShowResults(true)
    const timer = setTimeout(() => {
      abortRef.current?.abort()
      const controller = new AbortController()
      abortRef.current = controller

      setLoading(true)
      api.search(query, controller.signal)
        .then((r) => {
          setResult(r)
          setShowResults(true)
        })
        .catch(err => {
          if (err.name !== 'AbortError') console.error(err)
        })
        .finally(() => {
          if (!controller.signal.aborted) setLoading(false)
        })
    }, 250)

    return () => clearTimeout(timer)
  }, [query])

  // Close on click outside
  useEffect(() => {
    function handleClick(e: MouseEvent) {
      if (wrapperRef.current && !wrapperRef.current.contains(e.target as Node)) {
        setShowResults(false)
      }
    }
    document.addEventListener('mousedown', handleClick)
    return () => document.removeEventListener('mousedown', handleClick)
  }, [])

  function handleKeyDown(e: React.KeyboardEvent) {
    if (e.key === 'Escape') {
      setShowResults(false)
    }
  }

  function handleSelect() {
    setShowResults(false)
    setQuery('')
    setResult(null)
    setMenuOpen(false)
  }

  function cycleTheme() {
    const order: Theme[] = ['light', 'dark', 'system']
    const next = order[(order.indexOf(theme) + 1) % order.length]
    setTheme(next)
  }

  return (
    <header className="bg-gh-header text-white sticky top-0 z-50">
      <div className="max-w-[1012px] mx-auto px-4 flex items-center gap-4 h-14">
        <Link to="/" className="shrink-0">
          <img src="/logo.png" alt="mdwiki" width={91} height={30} />
        </Link>

        <div ref={wrapperRef} className="relative flex-1 hidden sm:flex">
          <input
            type="text"
            placeholder="Search..."
            value={query}
            onChange={e => setQuery(e.target.value)}
            onFocus={() => { if (query.trim() && result) setShowResults(true) }}
            onKeyDown={handleKeyDown}
            className="w-full max-w-xs px-3 py-1 rounded bg-white/10 border border-white/20 text-white placeholder-white/50 text-sm focus:outline-none focus:bg-white/20"
          />
          {showResults && (
            <SearchDropdown result={result} loading={loading} query={query} onSelect={handleSelect} />
          )}
        </div>

        <nav className="hidden sm:flex items-center gap-4 text-sm text-white/80 shrink-0">
          <Link to="/list" className="hover:text-white">Pages</Link>
          <Link to="/create" className="hover:text-white">New Page</Link>
          <Link to="/trash" className="hover:text-white">Trash</Link>
          <button
            onClick={cycleTheme}
            className="hover:text-white flex items-center gap-1"
            title={`Theme: ${theme}`}
          >
            <ThemeIcon theme={theme} />
          </button>
        </nav>

        <button
          className="sm:hidden ml-auto text-white/80 hover:text-white"
          onClick={() => setMenuOpen(!menuOpen)}
          aria-label="Toggle menu"
        >
          <svg className="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2}
              d={menuOpen ? 'M6 18L18 6M6 6l12 12' : 'M4 6h16M4 12h16M4 18h16'} />
          </svg>
        </button>
      </div>

      {menuOpen && (
        <div className="sm:hidden bg-gh-header border-t border-white/10 px-4 pb-4 flex flex-col gap-3 text-sm text-white/80">
          <div className="relative pt-3">
            <input
              type="text"
              placeholder="Search..."
              value={query}
              onChange={e => setQuery(e.target.value)}
              onKeyDown={handleKeyDown}
              className="w-full px-3 py-1 rounded bg-white/10 border border-white/20 text-white placeholder-white/50 text-sm focus:outline-none"
            />
            {showResults && (
              <SearchDropdown result={result} loading={loading} query={query} onSelect={handleSelect} />
            )}
          </div>
          <Link to="/list" onClick={() => setMenuOpen(false)} className="hover:text-white">Pages</Link>
          <Link to="/create" onClick={() => setMenuOpen(false)} className="hover:text-white">New Page</Link>
          <Link to="/trash" onClick={() => setMenuOpen(false)} className="hover:text-white">Trash</Link>
          <button
            onClick={cycleTheme}
            className="hover:text-white flex items-center gap-1 text-left"
          >
            <ThemeIcon theme={theme} />
            <span className="capitalize">{theme}</span>
          </button>
        </div>
      )}
    </header>
  )
}
