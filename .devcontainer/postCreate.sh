#!/usr/bin/env bash
set -euo pipefail
echo "Running post-create tasks..."

# Ensure GOPATH exists
mkdir -p "$GOPATH/bin" || true

echo "Installing Go language tools..."
if command -v go >/dev/null 2>&1; then
  export GOBIN="$GOBIN"
  go install golang.org/x/tools/gopls@latest || true
  go install github.com/go-delve/delve/cmd/dlv@latest || true
  echo "Installing air (live-reload)"
  go install github.com/cosmtrek/air@latest || true
fi

echo "Enabling corepack (yarn support)..."
corepack enable || true

if [ -d frontend ]; then
  echo "Installing frontend dependencies in ./frontend (Vue 3 + TypeScript)..."
  if [ -f frontend/package-lock.json ] || [ -f frontend/yarn.lock ] || [ -f frontend/package.json ]; then
    (cd frontend && if [ -f yarn.lock ]; then yarn install || true; else npm install || true; fi)
  fi
fi

if [ -d api ]; then
  echo "Downloading Go module dependencies for ./api"
  if command -v go >/dev/null 2>&1; then
    (cd api && go mod download || true)
  fi
fi

echo "Waiting for CockroachDB to be ready..."
# Wait for CockroachDB to be ready
max_attempts=30
attempt=0
while [ $attempt -lt $max_attempts ]; do
  if curl -f http://cockroachdb:8080/health >/dev/null 2>&1; then
    echo "CockroachDB is ready!"
    break
  fi
  attempt=$((attempt + 1))
  echo "Waiting for CockroachDB... (attempt $attempt/$max_attempts)"
  sleep 2
done

if [ $attempt -eq $max_attempts ]; then
  echo "Warning: CockroachDB may not be ready yet. You may need to wait a bit longer."
else
  echo "Initializing CockroachDB database..."
  # Initialize database if needed (optional - CockroachDB creates defaultdb automatically)
  echo "CockroachDB is running and ready to use!"
  echo ""
  echo "Connection string: postgresql://root@cockroachdb:26257/defaultdb?sslmode=disable"
  echo "Web UI available at: http://localhost:8081"
fi

echo ""
echo "To start the API with live reload, run: cd /workspace/api && air -c .air.toml"
echo "Frontend is already running on port 3000"
echo ""
echo "Post-create script finished."
