#!/bin/bash
# Script para iniciar tanto Vite como Storybook en paralelo

# NO usar set -e aquí porque queremos capturar errores y mostrarlos
set +e

# Función para logging
log() {
    echo "[$(date +'%Y-%m-%d %H:%M:%S')] $*" >&2
}

# Función para manejar errores
handle_error() {
    log "❌ Error en línea $1: $2"
    log "📋 Estado actual:"
    log "   Directorio: $(pwd)"
    log "   Usuario: $(whoami)"
    log "   PID: $$"
    exit 1
}

trap 'handle_error $LINENO "$BASH_COMMAND"' ERR

log "🚀 Iniciando script de frontend..."
log "   Usuario: $(whoami)"
log "   PID: $$"

# Cambiar al directorio del frontend
cd /workspace/frontend || {
    log "❌ Error: No se pudo acceder al directorio /workspace/frontend"
    exit 1
}

log "📁 Directorio actual: $(pwd)"
log "📋 Contenido del directorio:"
ls -la | head -10

# Instalar dependencias si no están instaladas
if [ ! -d "node_modules" ]; then
    log "📦 Instalando dependencias..."
    npm install
    if [ $? -ne 0 ]; then
        log "❌ Error: No se pudieron instalar las dependencias"
        exit 1
    fi
    log "✅ Dependencias instaladas"
else
    log "✅ Dependencias ya instaladas"
fi

# Verificar que concurrently esté instalado
if [ ! -f "node_modules/.bin/concurrently" ]; then
    log "📦 Instalando concurrently..."
    npm install concurrently --save-dev
    if [ $? -ne 0 ]; then
        log "❌ Error: No se pudo instalar concurrently"
        exit 1
    fi
    log "✅ concurrently instalado"
fi

# Verificar que los scripts existan en package.json
if ! grep -q '"dev"' package.json; then
    log "❌ Error: Script 'dev' no encontrado en package.json"
    exit 1
fi

if ! grep -q '"storybook"' package.json; then
    log "❌ Error: Script 'storybook' no encontrado en package.json"
    exit 1
fi

# Verificar que los binarios existan
if [ ! -f "node_modules/.bin/vite" ]; then
    log "❌ Error: Vite no está instalado en node_modules/.bin/vite"
    exit 1
fi

log "🚀 Iniciando Vite y Storybook en paralelo..."
log "   - Vite: http://0.0.0.0:3000 (puerto 3001 en el host)"
log "   - Storybook: http://0.0.0.0:6006"
log ""

# Mostrar qué comandos se van a ejecutar
log "📋 Comandos a ejecutar:"
log "   1. npm run dev"
log "   2. npm run storybook"
log ""

# Ejecutar ambos servicios con concurrently
# Usar exec para que concurrently sea el proceso principal y mantenga el contenedor vivo
log "▶️  Ejecutando concurrently..."
exec npx concurrently \
    --names "VITE,STORYBOOK" \
    --prefix-colors "cyan,magenta" \
    --kill-others-on-fail \
    --raw \
    "npm run dev" \
    "npm run storybook"

# Si llegamos aquí, algo salió mal
log "❌ concurrently terminó inesperadamente"
exit 1
