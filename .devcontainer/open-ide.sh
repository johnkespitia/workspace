#!/usr/bin/env bash
# Script para abrir el IDE (Cursor o VS Code) con el devcontainer

set -e

# Colores
GREEN='\033[0;32m'
YELLOW='\033[0;33m'
RED='\033[0;31m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Obtener el directorio del proyecto (raíz del workspace)
PROJECT_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

# Función para detectar qué IDE está disponible
detect_ide() {
    if command -v cursor &> /dev/null; then
        echo "cursor"
    elif command -v code &> /dev/null; then
        echo "code"
    else
        echo ""
    fi
}

# Función para verificar que los contenedores estén corriendo
wait_for_containers() {
    echo -e "${BLUE}⏳ Esperando a que los contenedores estén listos...${NC}"
    local max_attempts=30
    local attempt=0
    
    while [ $attempt -lt $max_attempts ]; do
        if docker ps --filter "name=go-react-test-api" --format "{{.Status}}" | grep -q "Up" 2>/dev/null; then
            echo -e "${GREEN}✅ Contenedores listos${NC}"
            # Esperar un poco más para que el API se inicie
            sleep 3
            return 0
        fi
        attempt=$((attempt + 1))
        sleep 1
    done
    
    echo -e "${YELLOW}⚠️  Los contenedores pueden no estar completamente listos${NC}"
    return 0
}

# Función para abrir con Cursor
open_cursor() {
    echo -e "${BLUE}🚀 Abriendo Cursor...${NC}"
    
    # Esperar a que los contenedores estén listos
    wait_for_containers
    
    cd "$PROJECT_ROOT"
    
    # Intentar conectar directamente al devcontainer usando el CLI
    echo -e "${YELLOW}📦 Intentando conectar al devcontainer...${NC}"
    
    # Verificar si el contenedor del API está corriendo
    if ! docker ps --filter "name=go-react-test-api" --format "{{.Names}}" | grep -q "go-react-test-api" 2>/dev/null; then
        echo -e "${RED}❌ El contenedor del API no está corriendo${NC}"
        echo -e "${YELLOW}Ejecuta 'make dev-up' primero${NC}"
        return 1
    fi
    
    # Abrir Cursor en el directorio del proyecto
    # Cursor detectará automáticamente el devcontainer.json y preguntará si deseas conectarte
    echo -e "${GREEN}✅ Abriendo Cursor en el directorio del proyecto...${NC}"
    cursor . 2>/dev/null || {
        echo -e "${YELLOW}⚠️  Cursor no se pudo abrir automáticamente${NC}"
        echo -e "${YELLOW}Por favor, abre Cursor manualmente en: ${PROJECT_ROOT}${NC}"
        return 1
    }
    
    echo -e "${GREEN}✅ Cursor abierto${NC}"
    echo ""
    echo -e "${YELLOW}📋 Si no se conectó automáticamente:${NC}"
    echo -e "${YELLOW}   1. Cursor debería detectar automáticamente el devcontainer${NC}"
    echo -e "${YELLOW}   2. Si aparece una notificación, haz clic en 'Reopen in Container'${NC}"
    echo -e "${YELLOW}   3. O manualmente: Command Palette (Cmd+Shift+P) > 'Dev Containers: Reopen in Container'${NC}"
    echo ""
}

# Función para abrir con VS Code
open_vscode() {
    echo -e "${BLUE}🚀 Abriendo VS Code...${NC}"
    
    # Esperar a que los contenedores estén listos
    wait_for_containers
    
    cd "$PROJECT_ROOT"
    
    # Intentar conectar directamente al devcontainer usando el CLI
    echo -e "${YELLOW}📦 Intentando conectar al devcontainer...${NC}"
    
    # Verificar si el contenedor del API está corriendo
    if ! docker ps --filter "name=go-react-test-api" --format "{{.Names}}" | grep -q "go-react-test-api" 2>/dev/null; then
        echo -e "${RED}❌ El contenedor del API no está corriendo${NC}"
        echo -e "${YELLOW}Ejecuta 'make dev-up' primero${NC}"
        return 1
    fi
    
    # Abrir VS Code en el directorio del proyecto
    # VS Code detectará automáticamente el devcontainer.json y preguntará si deseas conectarte
    echo -e "${GREEN}✅ Abriendo VS Code en el directorio del proyecto...${NC}"
    code . 2>/dev/null || {
        echo -e "${YELLOW}⚠️  VS Code no se pudo abrir automáticamente${NC}"
        echo -e "${YELLOW}Por favor, abre VS Code manualmente en: ${PROJECT_ROOT}${NC}"
        return 1
    }
    
    echo -e "${GREEN}✅ VS Code abierto${NC}"
    echo ""
    echo -e "${YELLOW}📋 Si no se conectó automáticamente:${NC}"
    echo -e "${YELLOW}   1. VS Code debería detectar automáticamente el devcontainer${NC}"
    echo -e "${YELLOW}   2. Si aparece una notificación, haz clic en 'Reopen in Container'${NC}"
    echo -e "${YELLOW}   3. O manualmente: Command Palette (Cmd+Shift+P) > 'Dev Containers: Reopen in Container'${NC}"
    echo ""
}

# Función para preguntar al usuario
ask_user() {
    local available_ide=$(detect_ide)
    
    if [ -z "$available_ide" ]; then
        echo -e "${RED}❌ No se encontró Cursor ni VS Code instalado${NC}"
        echo -e "${YELLOW}Por favor, instala Cursor o VS Code y asegúrate de que esté en tu PATH${NC}"
        return 1
    fi
    
    echo -e "${GREEN}IDE disponible: ${available_ide}${NC}"
    echo ""
    echo -e "${YELLOW}¿Qué IDE deseas usar?${NC}"
    echo "  1) Cursor (si está disponible)"
    echo "  2) VS Code (si está disponible)"
    echo "  3) Auto-detectar (usar el primero disponible)"
    echo "  4) Cancelar"
    echo ""
    read -p "Selecciona una opción (1-4): " choice
    
    case $choice in
        1)
            if command -v cursor &> /dev/null; then
                open_cursor
            else
                echo -e "${RED}❌ Cursor no está disponible${NC}"
                return 1
            fi
            ;;
        2)
            if command -v code &> /dev/null; then
                open_vscode
            else
                echo -e "${RED}❌ VS Code no está disponible${NC}"
                return 1
            fi
            ;;
        3)
            if [ "$available_ide" = "cursor" ]; then
                open_cursor
            else
                open_vscode
            fi
            ;;
        4)
            echo -e "${YELLOW}Operación cancelada${NC}"
            return 0
            ;;
        *)
            echo -e "${RED}❌ Opción inválida${NC}"
            return 1
            ;;
    esac
}

# Lógica principal
main() {
    # Verificar variable de entorno IDE
    if [ -n "$IDE" ]; then
        case "$IDE" in
            cursor)
                if command -v cursor &> /dev/null; then
                    open_cursor
                else
                    echo -e "${RED}❌ Cursor no está disponible${NC}"
                    echo -e "${YELLOW}Cambiando a auto-detección...${NC}"
                    ask_user
                fi
                ;;
            code|vscode)
                if command -v code &> /dev/null; then
                    open_vscode
                else
                    echo -e "${RED}❌ VS Code no está disponible${NC}"
                    echo -e "${YELLOW}Cambiando a auto-detección...${NC}"
                    ask_user
                fi
                ;;
            auto)
                local ide=$(detect_ide)
                if [ -z "$ide" ]; then
                    echo -e "${RED}❌ No se encontró ningún IDE${NC}"
                    return 1
                fi
                if [ "$ide" = "cursor" ]; then
                    open_cursor
                else
                    open_vscode
                fi
                ;;
            *)
                echo -e "${YELLOW}⚠️  Valor inválido para IDE: $IDE${NC}"
                echo -e "${YELLOW}Valores válidos: cursor, code, vscode, auto${NC}"
                ask_user
                ;;
        esac
    else
        # Si no hay variable de entorno, preguntar al usuario
        ask_user
    fi
}

main "$@"
