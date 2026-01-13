.PHONY: help dev-init dev-up dev-down dev-rebuild dev-logs dev-clean dev-shell dev-status

# Variables
COMPOSE_FILE := .devcontainer/docker-compose.yml
DEV_CONTAINER_DIR := .devcontainer

# Colores para output
GREEN := \033[0;32m
YELLOW := \033[0;33m
RED := \033[0;31m
NC := \033[0m # No Color

help: ## Muestra esta ayuda
	@echo "$(GREEN)Comandos disponibles para el Dev Container:$(NC)"
	@echo ""
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  $(YELLOW)%-15s$(NC) %s\n", $$1, $$2}'
	@echo ""

dev-init: ## Inicializa el dev container (construye e inicia servicios)
	@echo "$(GREEN)🚀 Inicializando dev container...$(NC)"
	@cd $(DEV_CONTAINER_DIR) && docker-compose -f docker-compose.yml build
	@cd $(DEV_CONTAINER_DIR) && docker-compose -f docker-compose.yml up -d
	@echo "$(GREEN)✅ Dev container inicializado$(NC)"
	@echo "$(YELLOW)📝 Ejecuta 'make dev-logs' para ver los logs$(NC)"
	@echo "$(YELLOW)📝 Ejecuta 'make dev-status' para ver el estado$(NC)"

dev-up: ## Inicia los servicios del dev container
	@echo "$(GREEN)▶️  Iniciando servicios...$(NC)"
	@cd $(DEV_CONTAINER_DIR) && docker-compose -f docker-compose.yml up -d
	@echo "$(GREEN)✅ Servicios iniciados$(NC)"
	@echo "$(YELLOW)Frontend: http://localhost:3000$(NC)"
	@echo "$(YELLOW)Backend: http://localhost:8080$(NC)"
	@echo "$(YELLOW)CockroachDB UI: http://localhost:8081$(NC)"

dev-down: ## Detiene los servicios del dev container
	@echo "$(YELLOW)⏹️  Deteniendo servicios...$(NC)"
	@cd $(DEV_CONTAINER_DIR) && docker-compose -f docker-compose.yml down
	@echo "$(GREEN)✅ Servicios detenidos$(NC)"

dev-stop: dev-down ## Alias para dev-down

dev-rebuild: ## Reconstruye las imágenes y reinicia los servicios
	@echo "$(YELLOW)🔨 Reconstruyendo imágenes...$(NC)"
	@cd $(DEV_CONTAINER_DIR) && docker-compose -f docker-compose.yml build --no-cache
	@echo "$(GREEN)▶️  Reiniciando servicios...$(NC)"
	@cd $(DEV_CONTAINER_DIR) && docker-compose -f docker-compose.yml up -d
	@echo "$(GREEN)✅ Dev container reconstruido y reiniciado$(NC)"

dev-logs: ## Muestra los logs de todos los servicios
	@cd $(DEV_CONTAINER_DIR) && docker-compose -f docker-compose.yml logs -f

dev-logs-api: ## Muestra los logs del servicio API
	@cd $(DEV_CONTAINER_DIR) && docker-compose -f docker-compose.yml logs -f api

dev-logs-frontend: ## Muestra los logs del servicio Frontend
	@cd $(DEV_CONTAINER_DIR) && docker-compose -f docker-compose.yml logs -f frontend

dev-logs-db: ## Muestra los logs de CockroachDB
	@cd $(DEV_CONTAINER_DIR) && docker-compose -f docker-compose.yml logs -f cockroachdb

dev-status: ## Muestra el estado de los servicios
	@echo "$(GREEN)📊 Estado de los servicios:$(NC)"
	@cd $(DEV_CONTAINER_DIR) && docker-compose -f docker-compose.yml ps
	@echo ""
	@echo "$(YELLOW)Puertos expuestos:$(NC)"
	@echo "  Frontend:    http://localhost:3000"
	@echo "  Backend:     http://localhost:8080"
	@echo "  CockroachDB: http://localhost:8081 (Web UI)"
	@echo "  CockroachDB: localhost:26257 (SQL)"

dev-shell: ## Abre una shell en el contenedor API
	@cd $(DEV_CONTAINER_DIR) && docker-compose -f docker-compose.yml exec api bash

dev-shell-frontend: ## Abre una shell en el contenedor Frontend
	@cd $(DEV_CONTAINER_DIR) && docker-compose -f docker-compose.yml exec frontend bash

dev-clean: ## Detiene servicios y elimina volúmenes (⚠️  elimina datos de la BD)
	@echo "$(RED)⚠️  Esto eliminará los volúmenes y datos de CockroachDB$(NC)"
	@echo "$(YELLOW)Ejecutando limpieza...$(NC)"
	@cd $(DEV_CONTAINER_DIR) && docker-compose -f docker-compose.yml down -v
	@echo "$(GREEN)✅ Limpieza completada$(NC)"

dev-restart: ## Reinicia todos los servicios
	@echo "$(YELLOW)🔄 Reiniciando servicios...$(NC)"
	@cd $(DEV_CONTAINER_DIR) && docker-compose -f docker-compose.yml restart
	@echo "$(GREEN)✅ Servicios reiniciados$(NC)"

dev-restart-api: ## Reinicia solo el servicio API
	@cd $(DEV_CONTAINER_DIR) && docker-compose -f docker-compose.yml restart api
	@echo "$(GREEN)✅ API reiniciado$(NC)"

dev-restart-frontend: ## Reinicia solo el servicio Frontend
	@cd $(DEV_CONTAINER_DIR) && docker-compose -f docker-compose.yml restart frontend
	@echo "$(GREEN)✅ Frontend reiniciado$(NC)"

dev-health: ## Verifica el estado de salud de los servicios
	@echo "$(GREEN)🏥 Verificando salud de los servicios...$(NC)"
	@echo ""
	@echo "$(YELLOW)Backend:$(NC)"
	@curl -s http://localhost:8080/health || echo "$(RED)❌ Backend no responde$(NC)"
	@echo ""
	@echo "$(YELLOW)Frontend:$(NC)"
	@curl -s -o /dev/null -w "HTTP Status: %{http_code}\n" http://localhost:3000 || echo "$(RED)❌ Frontend no responde$(NC)"
	@echo ""
	@echo "$(YELLOW)CockroachDB:$(NC)"
	@cd $(DEV_CONTAINER_DIR) && docker-compose -f docker-compose.yml exec -T cockroachdb curl -s http://localhost:8080/health || echo "$(RED)❌ CockroachDB no responde$(NC)"
