# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## What This Is

mdwiki is a self-contained markdown wiki: a Go backend with an embedded React SPA frontend. No database — pages are plain `.md` files on disk. Single binary deployment.

## Commands

### Build
```bash
make          # Build frontend (Vite) + compile Go binary → ./mdwiki
make clean    # Remove mdwiki binary, frontend/dist/, tmp/
```

### Development
```bash
make dev          # Run both: Go backend (air hot reload) + Vite dev server
make dev-backend  # Go backend only (air, port 8080)
make dev-frontend # Vite only (port 5173)
```
In dev mode, set `VITE_DEV_URL` to proxy frontend requests from Go to Vite.

### Testing
```bash
make test         # go test ./...
```
Go tests only — no frontend test setup exists.

### Frontend (run from `frontend/`)
```bash
npm run lint      # ESLint
npm run build     # tsc -b && vite build
```

## Architecture

### Backend (`src/`)
- `application.go` — Entry point: loads config, initializes storage, sets up Echo, registers routes
- `src/api/rest.go` — All REST endpoints
- `src/storage/` — File-based persistence (pages, versions, trash as `.md` files under `.storage/`)
- `src/search/` — Fuzzy search over filenames + content using `sahilm/fuzzy`
- `src/handlers/spa.go` — Serves embedded frontend in production; proxies to Vite in dev

**Storage layout at runtime:**
```
.storage/
  pages/      # Active pages (filename.md)
  trash/      # Deleted pages
  versions/   # Snapshots named filename.md__<unix_timestamp>
```

Every save creates a version snapshot. Deleting moves to trash (filesystem rename).

### Frontend (`frontend/src/`)
- `App.tsx` — React Router routes
- `api.ts` — Typed fetch wrapper for all API calls
- `useTheme.ts` — Light/dark/system theme (persisted to localStorage, Tailwind `class` strategy)
- `components/` — `Navbar` (search + theme toggle), `Layout`, `MarkdownContent`
- `pages/` — One component per route (PageView, EditPage, VersionsList, TrashList, etc.)

Editor uses CodeMirror 6 with GitHub themes. No global state manager — local `useState`/`useEffect` only.

### Frontend ↔ Backend
- Production: frontend built to `frontend/dist/`, embedded in Go binary via `//go:embed`, served at `/static/app/`
- Development: Vite runs on :5173; Go proxies SPA routes to it

### API Routes
```
GET/POST       /api/pages
GET/PUT/DELETE /api/pages/:page
GET            /api/pages/:page/versions
GET/POST       /api/versions/:ver  (restore)
GET            /api/search?query=
GET/DELETE     /api/trash
GET/POST       /api/trash/:page  (restore)
GET            /api/settings
```

## Key Constraints
- Pages are identified by filename (without `.md`). Renaming is done via `PUT /api/pages/:page` with a new name in the body.
- The `home` page (`home.md`) is auto-created on first run and is the app's root route.
- No authentication — designed for single-user/trusted-network use.
