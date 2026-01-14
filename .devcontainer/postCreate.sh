#!/usr/bin/env bash
# No usar set -e para permitir que el script continúe incluso si hay errores menores
# Los errores se manejan individualmente con || echo
set +e
echo "Running post-create tasks..."

# Inicializar variables si no están definidas
export GOPATH="${GOPATH:-/root/go}"
export GOBIN="${GOBIN:-$GOPATH/bin}"
export PATH="$GOBIN:$PATH"

# Ensure GOPATH exists
mkdir -p "$GOPATH/bin" || true

echo "Installing Go language tools..."
if command -v go >/dev/null 2>&1; then
  echo "  Installing gopls..."
  go install golang.org/x/tools/gopls@latest 2>&1 || echo "    Warning: Failed to install gopls"
  echo "  Installing delve..."
  go install github.com/go-delve/delve/cmd/dlv@latest 2>&1 || echo "    Warning: Failed to install delve"
  echo "  Installing air (live-reload)..."
  go install github.com/air-verse/air@latest 2>&1 || echo "    Warning: Failed to install air"
else
  echo "  Warning: Go not found, skipping Go tools installation"
fi

echo "Enabling corepack (yarn support)..."
if command -v corepack >/dev/null 2>&1; then
  corepack enable 2>&1 || echo "  Warning: Failed to enable corepack"
else
  echo "  Warning: corepack not found, skipping"
fi

if [ -d frontend ]; then
  echo "Installing frontend dependencies in ./frontend (Vue 3 + TypeScript)..."
  if [ -f frontend/package-lock.json ] || [ -f frontend/yarn.lock ] || [ -f frontend/package.json ]; then
    (cd frontend && if [ -f yarn.lock ]; then yarn install 2>&1 || echo "    Warning: Failed to install with yarn"; else npm install 2>&1 || echo "    Warning: Failed to install with npm"; fi)
  else
    echo "  Warning: No package.json found in frontend directory"
  fi
else
  echo "  Warning: frontend directory not found"
fi

if [ -d api ]; then
  echo "Downloading Go module dependencies for ./api"
  if command -v go >/dev/null 2>&1; then
    (cd api && go mod download 2>&1 || echo "    Warning: Failed to download Go modules")
  else
    echo "  Warning: Go not found, skipping module download"
  fi
else
  echo "  Warning: api directory not found"
fi

echo "Checking CockroachDB availability..."
# Wait for CockroachDB to be ready (opcional, no fallar si no está disponible)
if command -v curl >/dev/null 2>&1; then
  max_attempts=10
  attempt=0
  cockroach_ready=false
  
  while [ $attempt -lt $max_attempts ]; do
    if curl -f -s --connect-timeout 2 http://cockroachdb:8080/health >/dev/null 2>&1; then
      echo "✅ CockroachDB is ready!"
      cockroach_ready=true
      break
    fi
    attempt=$((attempt + 1))
    if [ $attempt -lt $max_attempts ]; then
      echo "  Waiting for CockroachDB... (attempt $attempt/$max_attempts)"
      sleep 2
    fi
  done
  
  if [ "$cockroach_ready" = false ]; then
    echo "⚠️  Warning: CockroachDB may not be ready yet. It will be available once the service starts."
    echo "   Connection string: postgresql://root@cockroachdb:26257/defaultdb?sslmode=disable"
    echo "   Web UI: http://localhost:8081"
  else
    echo "✅ CockroachDB is running and ready to use!"
    echo "   Connection string: postgresql://root@cockroachdb:26257/defaultdb?sslmode=disable"
    echo "   Web UI: http://localhost:8081"
  fi
else
  echo "⚠️  Warning: curl not available, skipping CockroachDB health check"
  echo "   CockroachDB should be available at: postgresql://root@cockroachdb:26257/defaultdb?sslmode=disable"
fi

echo ""
echo "Verifying Docker access..."
if command -v docker >/dev/null 2>&1; then
  if docker ps >/dev/null 2>&1; then
    echo "✅ Docker access verified"
    echo "   Docker version: $(docker --version 2>/dev/null || echo 'unknown')"
    if command -v docker-compose >/dev/null 2>&1 || docker compose version >/dev/null 2>&1; then
      echo "✅ Docker Compose available"
      echo "   Compose version: $(docker-compose --version 2>/dev/null || docker compose version 2>/dev/null || echo 'unknown')"
    else
      echo "⚠️  Docker Compose not found, but Docker is available"
    fi
  else
    echo "⚠️  Docker command found but cannot access Docker daemon"
    echo "   This is normal if Docker socket is not mounted"
  fi
else
  echo "⚠️  Docker CLI not found in container"
fi

echo ""
echo "✅ Post-create script completed successfully!"
echo ""
echo "📝 Next steps:"
echo "   - API is configured to start automatically with air (hot reload)"
echo "   - Frontend is configured to run on port 3001"
echo "   - CockroachDB is available on port 26257 (SQL) and 8081 (Web UI)"
echo ""
echo "💡 Docker access:"
echo "   - To access Docker from inside the container, the Docker socket must be mounted"
echo "   - Run 'make dev-rebuild' from the host to rebuild with Docker access"
echo ""
exit 0
