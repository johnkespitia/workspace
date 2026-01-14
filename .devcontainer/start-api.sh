#!/bin/bash
# Script para iniciar el servidor API con hot reload
# Este script se ejecuta cuando el contenedor se inicia como servicio
# Para devcontainer, el contenedor se mantiene corriendo con bash

set +e  # No fallar en errores menores

export GOPATH=/root/go
export GOBIN=$GOPATH/bin
export PATH=$GOBIN:$PATH

# Si estamos en modo devcontainer (variable de entorno), mantener el contenedor corriendo
if [ -n "$REMOTE_CONTAINERS" ] || [ -n "$VSCODE_INJECTION" ]; then
  echo "Devcontainer mode detected. Container will stay running."
  echo "To start the API server, run: cd /workspace/api && air -c .air.toml"
  # Mantener el contenedor corriendo
  exec tail -f /dev/null
fi

# Modo servicio: iniciar el servidor automáticamente
echo "Starting API server..."
cd /workspace/api || {
  echo "Error: /workspace/api directory not found"
  exit 1
}

# Verificar que .air.toml existe
if [ ! -f ".air.toml" ]; then
  echo "Warning: .air.toml not found. Creating default configuration..."
  # Crear un .air.toml básico si no existe
  cat > .air.toml << 'EOF'
root = "."
testdata_dir = "testdata"
tmp_dir = "tmp"

[build]
  args_bin = []
  bin = "./tmp/main"
  cmd = "go build -o ./tmp/main ./cmd/main.go"
  delay = 1000
  exclude_dir = ["assets", "tmp", "vendor", "testdata"]
  exclude_file = []
  exclude_regex = ["_test.go"]
  exclude_unchanged = false
  follow_symlink = false
  full_bin = ""
  include_dir = []
  include_ext = ["go", "tpl", "tmpl", "html"]
  include_file = []
  kill_delay = "0s"
  log = "build-errors.log"
  poll = false
  poll_interval = 0
  rerun = false
  rerun_delay = 500
  send_interrupt = false
  stop_on_error = false

[color]
  app = ""
  build = "yellow"
  main = "magenta"
  runner = "green"
  watcher = "cyan"

[log]
  main_only = false
  time = false

[misc]
  clean_on_exit = false

[screen]
  clear_on_rebuild = false
  keep_scroll = true
EOF
fi

# Use full path to air
AIR_BIN=$GOBIN/air

# Verify air exists, install if not
if [ ! -f "$AIR_BIN" ]; then
  echo "Installing air to $AIR_BIN..."
  go install github.com/air-verse/air@latest || {
    echo "Error: Failed to install air"
    exit 1
  }
fi

# Verificar que air se instaló correctamente
if [ ! -f "$AIR_BIN" ]; then
  echo "Error: air binary not found at $AIR_BIN"
  exit 1
fi

echo "Running air from $AIR_BIN..."
exec "$AIR_BIN" -c .air.toml
