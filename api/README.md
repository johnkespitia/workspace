# API Backend - Sistema de Información de Acciones

## 📋 Descripción

Backend desarrollado en Go siguiendo arquitectura Domain-Driven Design (DDD) para gestionar información de acciones desde una API externa.

## 🏗️ Arquitectura

El proyecto sigue una arquitectura DDD con las siguientes capas:

```
api/
├── cmd/                    # Punto de entrada
│   └── main.go
├── internal/
│   ├── domain/            # Capa de Dominio (DDD)
│   │   ├── stock/         # Entidades y lógica de negocio de stocks
│   │   └── recommendation/ # Algoritmo de recomendación
│   ├── infrastructure/    # Capa de Infraestructura
│   │   ├── database/      # Conexión a CockroachDB
│   │   ├── external/      # Cliente API externa
│   │   └── repository/    # Implementación de repositorios
│   ├── application/       # Capa de Aplicación
│   │   └── services/      # Servicios de aplicación
│   └── config/            # Configuración
└── go.mod
```

## 🚀 Configuración

### Variables de Entorno

Crea un archivo `.env` o configura las siguientes variables:

```bash
# Base de datos
DB_HOST=localhost
DB_PORT=26257
DB_USER=root
DB_PASSWORD=
DB_NAME=stocks
DB_SSLMODE=disable

# API Externa
API_BASE_URL=https://api.karenai.click
API_KEY=tu_api_key_aqui

# Servidor
PORT=8080
```

### Instalación de Dependencias

```bash
go mod download
```

## 🗄️ Base de Datos

### CockroachDB

El proyecto usa CockroachDB (compatible con PostgreSQL).

**Nota**: Las migraciones se verifican automáticamente al iniciar la aplicación. Si la base de datos no está inicializada, se ejecutarán las migraciones automáticamente.

### Gestión de Migraciones

Las migraciones se encuentran en `internal/infrastructure/database/migrations/` y se ejecutan en orden alfabético. El sistema:

- Lee automáticamente todos los archivos `.sql` de la carpeta migrations
- Los ejecuta en orden alfabético
- Verifica si las tablas ya existen antes de ejecutar

### Esquema de Tabla

```sql
CREATE TABLE stocks (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    ticker VARCHAR(10) NOT NULL UNIQUE,
    company_name VARCHAR(255) NOT NULL,
    brokerage VARCHAR(255),
    action VARCHAR(50),
    rating_from VARCHAR(50) NOT NULL,
    rating_to VARCHAR(50) NOT NULL,
    target_from DECIMAL(10,2) NOT NULL,
    target_to DECIMAL(10,2) NOT NULL,
    created_at TIMESTAMP DEFAULT now(),
    updated_at TIMESTAMP DEFAULT now()
);
```

## 🔧 Compilación y Ejecución

### Compilar

```bash
# Compilar servidor principal
go build ./cmd/main.go

# Compilar herramienta de migraciones
go build ./cmd/migrate
```

### Ejecutar Servidor

```bash
./main
# o
go run ./cmd/main.go
```

### Comandos de Migración

La herramienta de migraciones permite gestionar el esquema de la base de datos:

#### Verificar estado de migraciones

```bash
go run ./cmd/migrate -check
```

#### Ejecutar migraciones pendientes

```bash
go run ./cmd/migrate -up
```

#### Reiniciar base de datos (⚠️ elimina todos los datos)

```bash
go run ./cmd/migrate -reset
```

**Nota**: El comando `-reset` requiere confirmación y eliminará todas las tablas antes de recrearlas.

## 📦 Dependencias Principales

- `github.com/lib/pq` - Driver PostgreSQL para CockroachDB
- `github.com/google/uuid` - Generación de UUIDs
- `github.com/shopspring/decimal` - Manejo preciso de decimales para precios

## 🎯 Funcionalidades Implementadas (FASE 1)

### ✅ Capa de Dominio

- **Entidad Stock**: Con validaciones y lógica de negocio
- **Value Objects**: Rating y Price con validaciones
- **Interfaces de Repositorio**: Definición de contratos
- **Servicios de Dominio**: Cálculo de scores y recomendaciones

### ✅ Capa de Infraestructura

- **Conexión a CockroachDB**: Configuración y gestión de conexiones
- **Cliente API Externa**: Cliente HTTP con paginación para `api.karenai.click`
- **Repositorio CockroachDB**: Implementación completa con:
  - CRUD operations
  - Batch upsert para optimización
  - Búsqueda con filtros y ordenamiento
  - Conteo de registros

### ✅ Capa de Aplicación

- **StockService**: Servicio para gestión de stocks
- **SyncService**: Servicio para sincronización desde API externa
- **RecommendationService**: Servicio para cálculo de recomendaciones

### ✅ Algoritmo de Recomendación

Implementado en `domain/recommendation/algorithm.go`:

- **Complejidad**: O(n log n)
- **Estrategia**:
  1. Filtrar stocks con rating positivo
  2. Calcular score basado en:
     - Cambio porcentual en target (50%)
     - Rating score (30%)
     - Action score (20%)
  3. Ordenar por score descendente
  4. Retornar top N recomendaciones

## ✅ FASE 2 Completada - GraphQL API

### Implementado

- ✅ **Schema GraphQL**: Tipos, queries, mutations e inputs definidos
- ✅ **Resolvers**: Implementados con inyección de dependencias
- ✅ **Handler GraphQL**: Endpoint `/query` con soporte CORS
- ✅ **GraphQL Playground**: Endpoint `/playground` para desarrollo
- ✅ **Integración**: Conectado con servicios de aplicación

### Endpoints Disponibles

- `POST /query` - Endpoint GraphQL principal
- `GET /playground` - GraphQL Playground (interfaz visual)
- `GET /health` - Health check
- `GET /docs` - Página principal de documentación
- `GET /docs/swagger` - Swagger UI (documentación OpenAPI interactiva)
- `GET /docs/openapi.yaml` - Especificación OpenAPI en formato YAML

### Ejemplo de Query

```graphql
query {
  stocks(filter: { ratings: ["Buy", "Strong Buy"] }, limit: 10) {
    stocks {
      ticker
      companyName
      ratingTo
      targetTo
    }
    totalCount
  }
}
```

### Ejemplo de Mutation

```graphql
mutation {
  syncStocks {
    success
    message
    stocksSynced
  }
}
```

## ✅ FASE 3 Completada - Documentación y Tests

### Documentación

- ✅ **OpenAPI/Swagger**: Especificación completa en `docs/openapi.yaml`
- ✅ **Documentación de API**: Guía completa en `docs/API_DOCUMENTATION.md`
- ✅ **Guía de Usuario**: Guía paso a paso en `docs/USER_GUIDE.md`
- ✅ **Ejemplos GraphQL**: Ejemplos prácticos en `docs/GRAPHQL_EXAMPLES.md`
- ✅ **Referencia GraphQL**: Schema completo en `docs/GRAPHQL_API_REFERENCE.md`

### Tests

- ✅ **Tests unitarios**: Cobertura de servicios, repositorios, resolvers y algoritmos
- ✅ **Cobertura**: ~50% del código backend

Ver más detalles en [TEST_SUMMARY.md](./TEST_SUMMARY.md) y [docs/README.md](./docs/README.md)

## 📝 Notas

- Las migraciones se ejecutan automáticamente al iniciar la aplicación
- El servidor incluye graceful shutdown
- Health check endpoint disponible en `/health`

## 🐛 Troubleshooting

### Error de conexión a base de datos

Verifica que CockroachDB esté corriendo y las variables de entorno estén correctas.

### Error de API externa

Verifica que `API_KEY` esté configurado correctamente.
