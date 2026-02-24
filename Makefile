.PHONY: all build frontend clean dev dev-frontend dev-backend test air

all: build

# Build frontend and then the Go binary (with embedded frontend)
build: frontend
	go build -o mdwiki .

# Build the React frontend
frontend: frontend/node_modules
	cd frontend && npm run build

# Install frontend dependencies
frontend/node_modules: frontend/package-lock.json
	cd frontend && npm ci

# Remove build artifacts
clean:
	rm -f mdwiki
	rm -rf frontend/dist tmp

# Install air if not present
air:
	@which air > /dev/null || go install github.com/air-verse/air@latest

# Run both frontend and backend dev servers
dev: air frontend/node_modules
	@mkdir -p frontend/dist
	@echo "Starting Go backend (air) and Vite dev server..."
	@trap 'kill 0' EXIT; \
		air & \
		echo "Waiting for backend on :8080..."; \
		while ! nc -z localhost 8080 2>/dev/null; do sleep 1; done; \
		echo "Backend ready, starting Vite..."; \
		cd frontend && npm run dev & \
		wait

# Run frontend dev server (with hot reload)
dev-frontend:
	cd frontend && npm run dev

# Run Go server with hot reload
dev-backend:
	air

# Run tests
test:
	go test ./...
