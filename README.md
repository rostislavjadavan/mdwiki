![](logo/mdwiki_logo_small.png)

# mdwiki

Simple, self-contained markdown wiki written in Go. Single binary, no database, no user accounts.

## Features

- Single binary deployment (frontend embedded in Go binary)
- Pages stored as plain markdown files
- Page versioning (every save creates a new version)
- Trash with restore support
- Fuzzy search (instant results in navbar dropdown as you type)
- React SPA frontend with Tailwind CSS

## Tech Stack

- **Backend:** Go, [Echo](https://echo.labstack.com/) framework, `embed` for static assets
- **Frontend:** React, TypeScript, Vite, Tailwind CSS
- **Storage:** Plain markdown files on disk

## How to Build

Requires Go 1.16+ and Node.js.

```
make
```

This will install frontend dependencies, build the React SPA, and compile the Go binary with the frontend embedded.

## Development

```
make dev
```

This starts both the Vite dev server (with hot reload) and the Go backend (with [air](https://github.com/air-verse/air) for hot reload). Air is installed automatically if not present.

Open the **Vite dev server** address (default `http://localhost:5173`) in your browser. API requests are proxied to the Go backend automatically.

## Production

```
make
./mdwiki
```

The built binary serves everything — open the **Go backend** address (default `http://localhost:8080`, configured in `config.yml`).

Other targets:

```
make clean          # Remove build artifacts
make test           # Run Go tests
make dev-frontend   # Run Vite dev server only
make dev-backend    # Run Go server with air only
```

## Configuration

Application loads `config.yml` on startup.

```yaml
host: localhost
port: 8080
storage: .storage
```

- `host` & `port` - web server bind address
- `storage` - directory where pages are stored

## Running as `systemd` Service

### Service file example

```ini
[Unit]
Description=mdwiki
After=network.target
Wants=network-online.target

[Service]
Restart=always
Type=simple
WorkingDirectory=/opt
ExecStart=/opt/mdwiki

[Install]
WantedBy=multi-user.target
```

- Update `ExecStart` and `WorkingDirectory` based on installation path
- Working directory must contain `config.yml`
- Place the file in `/etc/systemd/system/mdwiki.service`

### Running the service

- Start: `systemctl start mdwiki`
- Stop: `systemctl stop mdwiki`
- Status: `systemctl status mdwiki`
- Enable on boot: `systemctl enable mdwiki`

## Logo

Made using https://excalidraw.com
