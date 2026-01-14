.PHONY: help dev-init dev-up dev-down dev-rebuild dev-logs dev-clean dev-shell dev-status dev-open dev-diagnose dev-install-frontend

# Variables
COMPOSE_FILE := .devcontainer/docker-compose.yml
DEV_CONTAINER_DIR := .devcontainer
PROJECT_NAME := go-react-test-devcontainer
# IDE preferido: cursor, code, vscode, auto (o dejar vacío para preguntar)
IDE ?= $(shell echo $$IDE)

# Detectar si usar $(DOCKER_COMPOSE) (V1) o docker compose (V2)
DOCKER_COMPOSE := $(shell command -v $(DOCKER_COMPOSE) 2>/dev/null || echo "docker compose")

# Colores para output
GREEN := \033[0;32m
YELLOW := \033[0;33m
RED := \033[0;31m
BLUE := \033[0;34m
NC := \033[0m # No Color

help: ## Muestra esta ayuda
	@echo "$(GREEN)Comandos disponibles para el Dev Container:$(NC)"
	@echo ""
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  $(YELLOW)%-15s$(NC) %s\n", $$1, $$2}'
	@echo ""

dev-init: ## Inicializa el dev container (construye e inicia servicios)
	@echo "$(GREEN)🚀 Inicializando dev container...$(NC)"
	@cd $(DEV_CONTAINER_DIR) && $(DOCKER_COMPOSE) -p $(PROJECT_NAME) -f docker-compose.yml build
	@cd $(DEV_CONTAINER_DIR) && $(DOCKER_COMPOSE) -p $(PROJECT_NAME) -f docker-compose.yml up -d
	@echo "$(GREEN)✅ Dev container inicializado$(NC)"
	@echo "$(YELLOW)📝 Ejecuta 'make dev-logs' para ver los logs$(NC)"
	@echo "$(YELLOW)📝 Ejecuta 'make dev-status' para ver el estado$(NC)"
	@echo ""
	@echo "$(BLUE)💻 Abriendo IDE con devcontainer...$(NC)"
	@bash $(DEV_CONTAINER_DIR)/open-ide.sh || echo "$(YELLOW)⚠️  No se pudo abrir el IDE automáticamente$(NC)"

dev-up: ## Inicia los servicios del dev container
	@echo "$(GREEN)▶️  Iniciando servicios...$(NC)"
	@cd $(DEV_CONTAINER_DIR) && $(DOCKER_COMPOSE) -p $(PROJECT_NAME) -f docker-compose.yml up -d
	@echo "$(GREEN)✅ Servicios iniciados$(NC)"
	@echo "$(YELLOW)Frontend: http://localhost:3001$(NC)"
	@echo "$(YELLOW)Backend: http://localhost:8080$(NC)"
	@echo "$(YELLOW)CockroachDB UI: http://localhost:8081$(NC)"

dev-down: ## Detiene los servicios del dev container
	@echo "$(YELLOW)⏹️  Deteniendo servicios...$(NC)"
	@cd $(DEV_CONTAINER_DIR) && $(DOCKER_COMPOSE) -p $(PROJECT_NAME) -f docker-compose.yml down
	@echo "$(GREEN)✅ Servicios detenidos$(NC)"

dev-stop: dev-down ## Alias para dev-down

dev-rebuild: ## Reconstruye las imágenes y reinicia los servicios
	@echo "$(YELLOW)🔨 Reconstruyendo imágenes...$(NC)"
	@cd $(DEV_CONTAINER_DIR) && $(DOCKER_COMPOSE) -p $(PROJECT_NAME) -f docker-compose.yml build --no-cache
	@echo "$(GREEN)▶️  Reiniciando servicios...$(NC)"
	@cd $(DEV_CONTAINER_DIR) && $(DOCKER_COMPOSE) -p $(PROJECT_NAME) -f docker-compose.yml up -d
	@echo "$(GREEN)✅ Dev container reconstruido y reiniciado$(NC)"

dev-logs: ## Muestra los logs de todos los servicios
	@cd $(DEV_CONTAINER_DIR) && $(DOCKER_COMPOSE) -p $(PROJECT_NAME) -f docker-compose.yml logs -f

dev-logs-api: ## Muestra los logs del servicio API
	@cd $(DEV_CONTAINER_DIR) && $(DOCKER_COMPOSE) -p $(PROJECT_NAME) -f docker-compose.yml logs -f api

dev-logs-frontend: ## Muestra los logs del servicio Frontend
	@cd $(DEV_CONTAINER_DIR) && $(DOCKER_COMPOSE) -p $(PROJECT_NAME) -f docker-compose.yml logs -f frontend

dev-logs-db: ## Muestra los logs de CockroachDB
	@cd $(DEV_CONTAINER_DIR) && $(DOCKER_COMPOSE) -p $(PROJECT_NAME) -f docker-compose.yml logs -f cockroachdb

dev-status: ## Muestra el estado de los servicios
	@echo "$(GREEN)📊 Estado de los servicios:$(NC)"
	@cd $(DEV_CONTAINER_DIR) && $(DOCKER_COMPOSE) -p $(PROJECT_NAME) -f docker-compose.yml ps
	@echo ""
	@echo "$(YELLOW)Puertos expuestos:$(NC)"
	@echo "  Frontend:    http://localhost:3001"
	@echo "  Backend:     http://localhost:8080"
	@echo "  CockroachDB: http://localhost:8081 (Web UI)"
	@echo "  CockroachDB: localhost:26257 (SQL)"

dev-shell: ## Abre una shell en el contenedor API
	@cd $(DEV_CONTAINER_DIR) && $(DOCKER_COMPOSE) -p $(PROJECT_NAME) -f docker-compose.yml exec api bash

dev-shell-frontend: ## Abre una shell en el contenedor Frontend
	@cd $(DEV_CONTAINER_DIR) && $(DOCKER_COMPOSE) -p $(PROJECT_NAME) -f docker-compose.yml exec frontend bash

dev-clean: ## Detiene servicios y elimina volúmenes (⚠️  elimina datos de la BD)
	@echo "$(RED)⚠️  Esto eliminará los volúmenes y datos de CockroachDB$(NC)"
	@echo "$(YELLOW)Ejecutando limpieza...$(NC)"
	@cd $(DEV_CONTAINER_DIR) && $(DOCKER_COMPOSE) -p $(PROJECT_NAME) -f docker-compose.yml down -v
	@echo "$(GREEN)✅ Limpieza completada$(NC)"

dev-restart: ## Reinicia todos los servicios
	@echo "$(YELLOW)🔄 Reiniciando servicios...$(NC)"
	@cd $(DEV_CONTAINER_DIR) && $(DOCKER_COMPOSE) -p $(PROJECT_NAME) -f docker-compose.yml restart
	@echo "$(GREEN)✅ Servicios reiniciados$(NC)"

dev-restart-api: ## Reinicia solo el servicio API
	@cd $(DEV_CONTAINER_DIR) && $(DOCKER_COMPOSE) -p $(PROJECT_NAME) -f docker-compose.yml restart api
	@echo "$(GREEN)✅ API reiniciado$(NC)"

dev-restart-frontend: ## Reinicia solo el servicio Frontend
	@cd $(DEV_CONTAINER_DIR) && $(DOCKER_COMPOSE) -p $(PROJECT_NAME) -f docker-compose.yml restart frontend
	@echo "$(GREEN)✅ Frontend reiniciado$(NC)"

dev-health: ## Verifica el estado de salud de los servicios
	@echo "$(GREEN)🏥 Verificando salud de los servicios...$(NC)"
	@echo ""
	@echo "$(YELLOW)Backend:$(NC)"
	@if curl -s -f --max-time 5 http://localhost:8080/health >/dev/null 2>&1; then \
		echo "$(GREEN)✅ Backend responde correctamente$(NC)"; \
	else \
		echo "$(RED)❌ Backend no responde$(NC)"; \
	fi
	@echo ""
	@echo "$(YELLOW)Frontend:$(NC)"
	@if curl -s -f --max-time 5 -o /dev/null http://localhost:3001 >/dev/null 2>&1; then \
		HTTP_CODE=$$(curl -s -o /dev/null -w "%{http_code}" --max-time 5 http://localhost:3001 2>/dev/null || echo "000"); \
		if [ "$$HTTP_CODE" != "000" ] && [ "$$HTTP_CODE" != "" ]; then \
			echo "$(GREEN)✅ Frontend responde (HTTP $$HTTP_CODE)$(NC)"; \
		else \
			echo "$(RED)❌ Frontend no responde$(NC)"; \
		fi; \
	else \
		echo "$(RED)❌ Frontend no responde$(NC)"; \
	fi
	@echo ""
	@echo "$(YELLOW)CockroachDB:$(NC)"
	@if cd $(DEV_CONTAINER_DIR) && $(DOCKER_COMPOSE) -p $(PROJECT_NAME) -f docker-compose.yml exec -T cockroachdb curl -s -f --max-time 5 http://localhost:8080/health >/dev/null 2>&1; then \
		echo "$(GREEN)✅ CockroachDB responde correctamente$(NC)"; \
	else \
		echo "$(RED)❌ CockroachDB no responde$(NC)"; \
	fi

dev-open: ## Abre el IDE (Cursor o VS Code) con el devcontainer
	@echo "$(BLUE)💻 Abriendo IDE...$(NC)"
	@bash $(DEV_CONTAINER_DIR)/open-ide.sh

dev-install-frontend: ## Instala las dependencias del frontend
	@echo "$(GREEN)📦 Instalando dependencias del frontend...$(NC)"
	@cd $(DEV_CONTAINER_DIR) && $(DOCKER_COMPOSE) -p $(PROJECT_NAME) -f docker-compose.yml exec frontend npm install
	@echo "$(GREEN)✅ Dependencias instaladas$(NC)"

dev-diagnose: ## Diagnóstico detallado de los servicios
	@echo "$(GREEN)🔍 Diagnóstico de servicios...$(NC)"
	@echo ""
	@echo "$(YELLOW)=== Estado de contenedores ===$(NC)"
	@cd $(DEV_CONTAINER_DIR) && $(DOCKER_COMPOSE) -p $(PROJECT_NAME) -f docker-compose.yml ps
	@echo ""
	@echo "$(YELLOW)=== Últimos logs del Frontend ===$(NC)"
	@cd $(DEV_CONTAINER_DIR) && $(DOCKER_COMPOSE) -p $(PROJECT_NAME) -f docker-compose.yml logs --tail=20 frontend || echo "$(RED)No se pudieron obtener logs$(NC)"
	@echo ""
	@echo "$(YELLOW)=== Verificando puertos ===$(NC)"
	@echo "Puerto 3001 (Frontend):"
	@lsof -i :3001 2>/dev/null || netstat -an | grep :3001 2>/dev/null || echo "  No se encontró proceso escuchando en puerto 3001"
	@echo "Puerto 8080 (Backend):"
	@lsof -i :8080 2>/dev/null || netstat -an | grep :8080 2>/dev/null || echo "  No se encontró proceso escuchando en puerto 8080"
	@echo ""
	@echo "$(YELLOW)=== Verificando conectividad ===$(NC)"
	@echo "Frontend (localhost:3001):"
	@curl -s -o /dev/null -w "  HTTP Status: %{http_code}\n" --max-time 3 http://localhost:3001 2>/dev/null || echo "  ❌ No responde"
	@echo "Backend (localhost:8080):"
	@curl -s -o /dev/null -w "  HTTP Status: %{http_code}\n" --max-time 3 http://localhost:8080/health 2>/dev/null || echo "  ❌ No responde"
