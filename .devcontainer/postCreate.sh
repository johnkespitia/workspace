#!/usr/bin/env bash
# No usar set -e para permitir que el script continúe incluso si hay errores menores
# Los errores se manejan individualmente con || echo
set +e
echo "Running post-create tasks..."

# Inicializar variables si no están definidas
# El GOPATH ya está definido en el Dockerfile, pero lo verificamos
if [ -z "$GOPATH" ]; then
    export GOPATH="${HOME:-/root}/go"
fi
export GOBIN="${GOBIN:-$GOPATH/bin}"
export PATH="$GOBIN:$PATH"

# Ensure GOPATH exists con permisos correctos
# Verificar si el directorio padre existe primero
GOPATH_PARENT=$(dirname "$GOPATH")
if [ ! -d "$GOPATH_PARENT" ]; then
    mkdir -p "$GOPATH_PARENT" 2>/dev/null || {
        echo "⚠️  Warning: No se pudo crear el directorio padre de GOPATH: $GOPATH_PARENT"
        echo "   Usando /tmp/go como alternativa..."
        export GOPATH="/tmp/go"
        export GOBIN="$GOPATH/bin"
    }
fi

# Crear el directorio GOPATH si no existe
if [ ! -d "$GOPATH" ]; then
    mkdir -p "$GOPATH/bin" 2>/dev/null || {
        echo "⚠️  Warning: No se pudo crear GOPATH: $GOPATH"
        echo "   Las herramientas de Go pueden no instalarse correctamente"
    }
else
    mkdir -p "$GOPATH/bin" 2>/dev/null || true
fi

# Asegurar permisos correctos (solo si somos root o tenemos permisos)
if [ "$(id -u)" = "0" ] || [ -w "$GOPATH" ]; then
    chmod -R 755 "$GOPATH" 2>/dev/null || true
fi

echo "Installing Go language tools..."
if command -v go >/dev/null 2>&1; then
  # Verificar que GOPATH esté configurado y accesible
  if [ -z "$GOPATH" ] || [ ! -w "$GOPATH" ] 2>/dev/null; then
    echo "  ⚠️  Warning: GOPATH no está configurado o no es accesible"
    echo "     Las herramientas de Go se instalarán en el directorio de trabajo actual"
    export GOPATH="$(pwd)/.go"
    export GOBIN="$GOPATH/bin"
    mkdir -p "$GOPATH/bin" 2>/dev/null || true
  fi
  
  echo "  GOPATH: $GOPATH"
  echo "  GOBIN: $GOBIN"
  
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
echo "Configuring Docker access..."
if [ -S /var/run/docker.sock ]; then
  # Detectar el GID del socket de Docker montado
  DOCKER_SOCK_GID=$(stat -c '%g' /var/run/docker.sock 2>/dev/null || echo "")
  
  if [ -n "$DOCKER_SOCK_GID" ]; then
    echo "  Docker socket GID: $DOCKER_SOCK_GID"
    
    # Verificar si existe un grupo con ese GID
    EXISTING_GROUP=$(getent group "$DOCKER_SOCK_GID" | cut -d: -f1 2>/dev/null || echo "")
    
    if [ -n "$EXISTING_GROUP" ]; then
      echo "  Using existing group: $EXISTING_GROUP (GID: $DOCKER_SOCK_GID)"
      DOCKER_GROUP="$EXISTING_GROUP"
    else
      # Crear el grupo docker con el GID del socket
      DOCKER_GROUP="docker"
      if ! getent group "$DOCKER_GROUP" >/dev/null 2>&1; then
        groupadd -g "$DOCKER_SOCK_GID" "$DOCKER_GROUP" 2>/dev/null || {
          echo "  ⚠️  Warning: Could not create docker group with GID $DOCKER_SOCK_GID"
          echo "     Trying with default GID..."
          groupadd "$DOCKER_GROUP" 2>/dev/null || true
        }
      fi
    fi
    
    # Agregar usuarios al grupo docker
    CURRENT_USER=$(whoami)
    if [ "$CURRENT_USER" != "root" ]; then
      usermod -aG "$DOCKER_GROUP" "$CURRENT_USER" 2>/dev/null || true
    fi
    usermod -aG "$DOCKER_GROUP" root 2>/dev/null || true
    usermod -aG "$DOCKER_GROUP" vscode 2>/dev/null || true
    
    echo "  ✅ Added users to $DOCKER_GROUP group"
    echo "  ⚠️  Note: You may need to restart the container or run 'newgrp docker' for changes to take effect"
  fi
fi

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
    echo "   Socket permissions: $(ls -l /var/run/docker.sock 2>/dev/null || echo 'not found')"
    echo "   Current user groups: $(groups)"
    echo "   ⚠️  You may need to restart the container or run 'newgrp docker' for changes to take effect"
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
