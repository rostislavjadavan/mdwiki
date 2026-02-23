import { useState } from 'react'
import { Link, useNavigate } from 'react-router-dom'

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

export function Navbar({ theme, setTheme }: Props) {
  const [query, setQuery] = useState('')
  const [menuOpen, setMenuOpen] = useState(false)
  const navigate = useNavigate()

  function handleSearch(e: React.FormEvent) {
    e.preventDefault()
    if (query.trim()) {
      navigate(`/search?query=${encodeURIComponent(query.trim())}`)
    }
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

        <form onSubmit={handleSearch} className="flex-1 hidden sm:flex">
          <input
            type="text"
            placeholder="Search..."
            value={query}
            onChange={e => setQuery(e.target.value)}
            className="w-full max-w-xs px-3 py-1 rounded bg-white/10 border border-white/20 text-white placeholder-white/50 text-sm focus:outline-none focus:bg-white/20"
          />
        </form>

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
          <form onSubmit={handleSearch} className="pt-3">
            <input
              type="text"
              placeholder="Search..."
              value={query}
              onChange={e => setQuery(e.target.value)}
              className="w-full px-3 py-1 rounded bg-white/10 border border-white/20 text-white placeholder-white/50 text-sm focus:outline-none"
            />
          </form>
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
