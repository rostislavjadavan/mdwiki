import { useState } from 'react'
import { Link, useNavigate } from 'react-router-dom'

export function Navbar() {
  const [query, setQuery] = useState('')
  const [menuOpen, setMenuOpen] = useState(false)
  const navigate = useNavigate()

  function handleSearch(e: React.FormEvent) {
    e.preventDefault()
    if (query.trim()) {
      navigate(`/search?query=${encodeURIComponent(query.trim())}`)
    }
  }

  return (
    <header className="bg-gh-header text-white sticky top-0 z-50">
      <div className="max-w-[1012px] mx-auto px-4 flex items-center gap-4 h-14">
        <Link to="/" className="text-white font-bold text-lg shrink-0 hover:text-gray-300">
          mdwiki
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
        </div>
      )}
    </header>
  )
}
