# Infraestructura

## 📋 Resumen

Este documento describe la infraestructura del proyecto, incluyendo servicios, configuración de red, base de datos y deployment.

---

## 🏗️ Arquitectura de Servicios

### Servicios Principales

```
┌─────────────────────────────────────────┐
│         Docker Compose                   │
├─────────────────────────────────────────┤
│                                         │
│  ┌──────────────┐  ┌──────────────┐   │
│  │   Frontend   │  │     API      │   │
│  │   (Vue 3)    │  │    (Go)      │   │
│  │   :3001      │  │    :8080     │   │
│  └──────┬───────┘  └──────┬───────┘   │
│         │                 │            │
│         │                 │            │
│         └────────┬────────┘            │
│                  │                      │
│         ┌────────▼────────┐            │
│         │   CockroachDB   │            │
│         │     :26257      │            │
│         └─────────────────┘            │
│                                         │
└─────────────────────────────────────────┘
```

### Puertos

| Servicio        | Puerto Interno | Puerto Host | Descripción                  |
| --------------- | -------------- | ----------- | ---------------------------- |
| Frontend        | 3000           | 3001        | Servidor de desarrollo Vite  |
| Storybook       | 6006           | 6006        | Documentación de componentes |
| Backend API     | 8080           | 8080        | Servidor GraphQL             |
| CockroachDB SQL | 26257          | 26257       | Puerto SQL                   |
| CockroachDB UI  | 8080           | 8081        | Interfaz web                 |

---

## 🐳 Docker y Dev Container

### Configuración

El proyecto usa **Dev Containers** para un entorno de desarrollo consistente.

**Archivos principales**:

- `.devcontainer/devcontainer.json` - Configuración del dev container
- `.devcontainer/docker-compose.yml` - Orquestación de servicios
- `.devcontainer/Dockerfile.api` - Imagen del backend
- `.devcontainer/Dockerfile.frontend` - Imagen del frontend

### Servicios Docker

#### 1. API (Backend Go)

```yaml
api:
  build:
    context: ..
    dockerfile: .devcontainer/Dockerfile.api
  container_name: go-react-test-api
  volumes:
    - ..:/workspace:cached
  ports:
    - "8080:8080"
  environment:
    - DATABASE_URL=postgresql://root@cockroachdb:26257/defaultdb?sslmode=disable
  command: ["bash", "/workspace/.devcontainer/start-api-auto.sh"]
```

**Características**:

- Hot reload con `air`
- Acceso a Docker socket para comandos Docker
- Variables de entorno preconfiguradas

#### 2. Frontend (Vue 3)

```yaml
frontend:
  build:
    context: ..
    dockerfile: .devcontainer/Dockerfile.frontend
  container_name: go-react-test-frontend
  volumes:
    - ..:/workspace:cached
  ports:
    - "127.0.0.1:3001:3000"
    - "127.0.0.1:6006:6006"
  command: ["bash", "/workspace/.devcontainer/start-frontend.sh"]
```

**Características**:

- Hot reload con Vite
- Storybook integrado
- Mapeo de puertos a localhost

#### 3. CockroachDB

```yaml
cockroachdb:
  image: cockroachdb/cockroach:latest-v23.1
  container_name: go-react-test-cockroachdb
  ports:
    - "26257:26257"
    - "8081:8080"
  command: start-single-node --insecure
  healthcheck:
    test: ["CMD", "curl", "-f", "http://localhost:8080/health"]
```

**Características**:

- Modo single-node para desarrollo
- Sin SSL (insecure mode)
- Health check configurado

---

## 🗄️ Base de Datos

### CockroachDB

**Versión**: v23.1 (latest)

**Configuración**:

- **Modo**: Single-node (desarrollo)
- **SSL**: Deshabilitado (insecure)
- **Base de datos**: `defaultdb`
- **Usuario**: `root`
- **Puerto SQL**: `26257`
- **Puerto Web UI**: `8081`

### Conexión

**Desde el contenedor**:

```
postgresql://root@cockroachdb:26257/defaultdb?sslmode=disable
```

**Desde el host**:

```
postgresql://root@localhost:26257/defaultdb?sslmode=disable
```

### Migraciones

Las migraciones están en `api/internal/infrastructure/database/migrations/`.

**Ejecutar migraciones**:

```bash
cd api
go run cmd/migrate/main.go
```

**Estructura de migraciones**:

```
migrations/
├── 001_create_stocks_table.sql
└── ...
```

### Esquema de Base de Datos

```sql
CREATE TABLE stocks (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    ticker VARCHAR(10) NOT NULL UNIQUE,
    company_name VARCHAR(255) NOT NULL,
    brokerage VARCHAR(255),
    action VARCHAR(50),
    rating_from VARCHAR(50),
    rating_to VARCHAR(50),
    target_from DECIMAL(10,2),
    target_to DECIMAL(10,2),
    created_at TIMESTAMP DEFAULT now(),
    updated_at TIMESTAMP DEFAULT now()
);

-- Índices
CREATE INDEX idx_stocks_ticker ON stocks(ticker);
CREATE INDEX idx_stocks_rating_to ON stocks(rating_to);
CREATE INDEX idx_stocks_target_to ON stocks(target_to);
CREATE INDEX idx_stocks_company_name ON stocks(company_name);
```

---

## 🌐 Red y Conectividad

### Red Docker

Los servicios están en la misma red Docker (`app-network`) y se comunican usando los nombres de servicio:

- `api` → `cockroachdb:26257`
- `frontend` → `api:8080`

### Mapeo de Puertos

Los puertos se mapean explícitamente a `127.0.0.1` para evitar conflictos:

```yaml
ports:
  - "127.0.0.1:3001:3000" # Frontend
  - "127.0.0.1:6006:6006" # Storybook
  - "8080:8080" # Backend
  - "26257:26257" # CockroachDB SQL
  - "8081:8080" # CockroachDB UI
```

---

## 🔐 Variables de Entorno

### Backend

```bash
# Base de datos
DATABASE_URL=postgresql://root@cockroachdb:26257/defaultdb?sslmode=disable
COCKROACH_HOST=cockroachdb
COCKROACH_PORT=26257
COCKROACH_USER=root
COCKROACH_DB=defaultdb

# API Externa
API_BASE_URL=https://api.karenai.click
API_KEY=tu_api_key_aqui

# Servidor
PORT=8080
```

### Frontend

```bash
# GraphQL Endpoint
VITE_GRAPHQL_ENDPOINT=http://localhost:8080/query
```

---

## 📦 Dependencias

### Backend (Go)

**Herramientas principales**:

- Go 1.21+
- Air (hot reload)
- gqlgen (GraphQL)
- CockroachDB driver

**Instalación automática**:
Las herramientas se instalan automáticamente en `postCreate.sh`.

### Frontend (Node.js)

**Herramientas principales**:

- Node.js 18+
- npm/yarn
- Vite
- Vue 3
- Storybook

**Instalación automática**:
Las dependencias se instalan automáticamente en `postCreate.sh`.

---

## 🚀 Deployment (Futuro)

### Consideraciones

1. **Producción**:

   - CockroachDB multi-node
   - SSL habilitado
   - Variables de entorno seguras
   - Health checks

2. **CI/CD**:

   - Tests automáticos
   - Build de imágenes
   - Deployment automático

3. **Monitoreo**:
   - Logs centralizados
   - Métricas de performance
   - Alertas

---

## 🔧 Mantenimiento

### Limpiar Volúmenes

```bash
# Eliminar volúmenes (⚠️ elimina datos)
make dev-clean
```

### Reconstruir Imágenes

```bash
# Reconstruir sin caché
make dev-rebuild
```

### Ver Logs

```bash
# Todos los servicios
make dev-logs

# Servicio específico
make dev-logs-api
make dev-logs-frontend
make dev-logs-db
```

---

## 📚 Recursos

- [Docker Compose Documentation](https://docs.docker.com/compose/)
- [Dev Containers](https://code.visualstudio.com/docs/devcontainers/containers)
- [CockroachDB Documentation](https://www.cockroachlabs.com/docs/)

---

**Última actualización**: 2026-01-15
