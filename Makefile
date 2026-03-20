.PHONY: dev build frontend-build frontend-install clean

# Install frontend dependencies
frontend-install:
	cd frontend && npm install

# Build React frontend
frontend-build:
	cd frontend && npm run build

# Development: run Go server + Vite dev server concurrently
dev:
	@echo "Starting Go API server on :8080 and Vite dev server on :5173"
	@(cd frontend && npm run dev) & (go run .) ; wait

# Production build: build frontend then compile Go binary
build: frontend-build
	go build -o yalex .

clean:
	rm -rf frontend/dist yalex
