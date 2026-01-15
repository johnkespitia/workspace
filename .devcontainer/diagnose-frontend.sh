#!/bin/bash
# Script de diagnóstico para el frontend

echo "🔍 Diagnóstico del Frontend"
echo "============================"
echo ""

echo "📁 Directorio actual: $(pwd)"
echo ""

echo "📦 Verificando dependencias..."
if [ -d "node_modules" ]; then
    echo "✅ node_modules existe"
    if [ -f "node_modules/.bin/vite" ]; then
        echo "✅ Vite está instalado"
    else
        echo "❌ Vite NO está instalado"
    fi
    if [ -f "node_modules/.bin/storybook" ]; then
        echo "✅ Storybook está instalado"
    else
        echo "❌ Storybook NO está instalado"
    fi
    if [ -f "node_modules/.bin/concurrently" ]; then
        echo "✅ concurrently está instalado"
    else
        echo "❌ concurrently NO está instalado"
    fi
else
    echo "❌ node_modules NO existe"
fi
echo ""

echo "🌐 Verificando puertos..."
if command -v netstat >/dev/null 2>&1; then
    netstat -tlnp 2>/dev/null | grep -E "(3000|6006)" || echo "⚠️  No se encontraron puertos escuchando"
elif command -v ss >/dev/null 2>&1; then
    ss -tlnp 2>/dev/null | grep -E "(3000|6006)" || echo "⚠️  No se encontraron puertos escuchando"
else
    echo "⚠️  No se encontraron herramientas para verificar puertos"
fi
echo ""

echo "🔧 Verificando procesos..."
ps aux | grep -E "(vite|storybook|node.*dev|node.*storybook)" | grep -v grep || echo "⚠️  No se encontraron procesos de Vite o Storybook"
echo ""

echo "📋 Verificando package.json..."
if [ -f "package.json" ]; then
    echo "✅ package.json existe"
    if grep -q '"dev"' package.json; then
        echo "✅ Script 'dev' encontrado"
        echo "   $(grep '"dev"' package.json)"
    else
        echo "❌ Script 'dev' NO encontrado"
    fi
    if grep -q '"storybook"' package.json; then
        echo "✅ Script 'storybook' encontrado"
        echo "   $(grep '"storybook"' package.json)"
    else
        echo "❌ Script 'storybook' NO encontrado"
    fi
else
    echo "❌ package.json NO existe"
fi
echo ""

echo "🌍 Verificando configuración de red..."
if [ -f "vite.config.ts" ]; then
    echo "✅ vite.config.ts existe"
    if grep -q "host.*0.0.0.0" vite.config.ts; then
        echo "✅ Vite configurado para escuchar en 0.0.0.0"
    else
        echo "⚠️  Vite puede no estar configurado para escuchar en todas las interfaces"
    fi
else
    echo "❌ vite.config.ts NO existe"
fi
echo ""

echo "📊 Variables de entorno relevantes:"
echo "   HOSTNAME: ${HOSTNAME:-no definido}"
echo "   PORT: ${PORT:-no definido}"
echo ""

echo "✅ Diagnóstico completado"
