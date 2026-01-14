#!/bin/bash
# Script para iniciar el servidor API automáticamente en background
# Este script se ejecuta cuando el contenedor se inicia

set +e  # No fallar en errores menores

export GOPATH=/root/go
export GOBIN=$GOPATH/bin
export PATH=$GOBIN:$PATH

# Función para iniciar el API
start_api() {
    echo "🚀 Iniciando servidor API..."
    cd /workspace/api || {
        echo "❌ Error: /workspace/api directory not found"
        return 1
    }

    # Verificar que .air.toml existe
    if [ ! -f ".air.toml" ]; then
        echo "📝 Creando .air.toml..."
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

    # Verificar que air existe, instalar si no
    AIR_BIN=$GOBIN/air
    if [ ! -f "$AIR_BIN" ]; then
        echo "📦 Instalando air..."
        go install github.com/air-verse/air@latest || {
            echo "❌ Error: Failed to install air"
            return 1
        }
    fi

    # Verificar que air se instaló correctamente
    if [ ! -f "$AIR_BIN" ]; then
        echo "❌ Error: air binary not found at $AIR_BIN"
        return 1
    fi

    # Esperar a que CockroachDB esté listo
    echo "⏳ Esperando a que CockroachDB esté listo..."
    max_attempts=30
    attempt=0
    while [ $attempt -lt $max_attempts ]; do
        if curl -f -s --connect-timeout 2 http://cockroachdb:8080/health >/dev/null 2>&1; then
            echo "✅ CockroachDB está listo!"
            break
        fi
        attempt=$((attempt + 1))
        if [ $attempt -lt $max_attempts ]; then
            sleep 2
        fi
    done

    if [ $attempt -eq $max_attempts ]; then
        echo "⚠️  Warning: CockroachDB puede no estar completamente listo, pero continuando..."
    fi

    # Iniciar air en background
    echo "▶️  Iniciando API con air (hot reload)..."
    nohup "$AIR_BIN" -c .air.toml > /tmp/api.log 2>&1 &
    API_PID=$!
    echo $API_PID > /tmp/api.pid
    echo "✅ API iniciado en background (PID: $API_PID)"
    echo "📋 Logs disponibles en: /tmp/api.log"
    echo "🌐 API disponible en: http://localhost:8080"
    
    # Esperar un momento para verificar que se inició correctamente
    sleep 5
    if ps -p $API_PID > /dev/null 2>&1; then
        echo "✅ API está corriendo correctamente"
        # Verificar que el API responda
        sleep 2
        if curl -f -s http://localhost:8080/health >/dev/null 2>&1; then
            echo "✅ API responde correctamente en /health"
        else
            echo "⚠️  API iniciado pero aún no responde (puede estar compilando...)"
        fi
    else
        echo "⚠️  Warning: El proceso del API puede haber terminado. Revisa los logs en /tmp/api.log"
        cat /tmp/api.log 2>/dev/null | tail -20 || true
    fi
}

# Si se ejecuta directamente, iniciar el API
if [ "${BASH_SOURCE[0]}" = "${0}" ]; then
    start_api
    # Mantener el contenedor corriendo y mostrar logs del API
    echo ""
    echo "📋 El contenedor se mantendrá corriendo. Para ver los logs del API:"
    echo "   docker logs -f go-react-test-api"
    echo "   o dentro del contenedor: tail -f /tmp/api.log"
    echo ""
    # Mantener el contenedor corriendo
    exec tail -f /tmp/api.log /dev/null 2>/dev/null || exec tail -f /dev/null
fi
