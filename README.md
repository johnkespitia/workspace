# Sistema de Información de Acciones

Sistema completo para recuperar, almacenar y visualizar información de acciones desde una API externa, con recomendaciones inteligentes de inversión.

**Stack**: Go + Vue 3 + GraphQL + CockroachDB + DDD

## 📚 Documentación

- **[Plan de Acción](./PLAN_DE_ACCION.md)**: Plan detallado por fases de implementación
- **[Arquitectura](./docs/ARCHITECTURE.md)**: Arquitectura DDD, capas y flujos de datos
- **[Algoritmos](./docs/ALGORITHMS.md)**: Algoritmos de recomendación, búsqueda y optimizaciones
- **[HOCs en Vue 3](./docs/HOCS_VUE3.md)**: Guía completa de Higher Order Components (NO hooks)
- **[Flujos del Sistema](./docs/FLUJOS.md)**: Diagramas de flujo de todos los procesos
- **[Resumen Ejecutivo](./docs/RESUMEN_EJECUTIVO.md)**: Visión general del proyecto

---

## 🎯 Características Principales

- ✅ **Sincronización de Datos**: Conexión a API externa con almacenamiento en CockroachDB
- ✅ **API GraphQL**: Consultas, búsqueda, filtrado y recomendaciones
- ✅ **Interfaz Moderna**: Design System con componentes reusables documentados en Storybook
- ✅ **Temas**: Soporte para light/dark mode
- ✅ **Accesibilidad**: WCAG AA compliance, navegación por teclado
- ✅ **Optimización**: Cache, debounce, request deduplication
- ✅ **Algoritmo de Recomendación**: O(n log n) para identificar mejores acciones

---

# Go + Vue 3 - Stack Tecnológico

Proyecto full-stack con backend en Golang y frontend en Vue 3.

## Stack Tecnológico

### Backend
- **Golang** - Lenguaje de programación del servidor

### Frontend
- **Vue 3** - Framework JavaScript progresivo
- **TypeScript** - Superset tipado de JavaScript
- **Pinia** - Store de estado para Vue
- **Tailwind CSS** - Framework CSS utility-first
- **Vite** - Build tool y dev server

### Base de Datos
- **CockroachDB** - Base de datos distribuida

## Configuración del Entorno

### Opción 1: Usando Dev Container (Recomendado)

El proyecto incluye un devcontainer configurado con todos los servicios necesarios.

1. **Abrir en Dev Container:**
   - Usa VS Code con la extensión "Dev Containers"
   - Abre la carpeta del proyecto
   - Ejecuta el comando: `Remote-Containers: Reopen in Container`

2. **Servicios incluidos (se inician automáticamente):**
   - **API (Go)**: Puerto 8080 con hot reload (Air)
     - Se inicia automáticamente al abrir el devcontainer
     - Usa Air para recarga automática al modificar archivos Go
     - Disponible en: `http://localhost:8080`
   - **Frontend (Vue 3)**: Puerto 3000 con hot reload (Vite)
     - Se inicia automáticamente al abrir el devcontainer
     - Usa Vite para recarga automática al modificar archivos Vue/TypeScript
     - Disponible en: `http://localhost:3000`
   - **CockroachDB**: Puerto 26257 (SQL) y 8081 (Web UI)
     - Se inicia automáticamente al abrir el devcontainer
     - Base de datos: `defaultdb`
     - Web UI disponible en: `http://localhost:8081`

3. **Inicio automático de servicios:**
   - Todos los servicios se inician automáticamente cuando se abre el devcontainer
   - No es necesario ejecutar comandos manuales para iniciar los servicios
   - El script `postCreate.sh` se ejecuta automáticamente para instalar dependencias y configurar el entorno

4. **Variables de entorno:**
   - Todas las variables están preconfiguradas en el devcontainer
   - Ver `.devcontainer/ENV_VARIABLES.md` para más detalles

### Opción 2: Instalación Local

#### Prerrequisitos
- Node.js 18+ y npm
- Go 1.21+
- CockroachDB instalado localmente

#### Instalación

1. **Instalar dependencias del frontend:**
```bash
cd frontend
npm install
```

2. **Ejecutar el frontend en modo desarrollo:**
```bash
npm run dev
```
El frontend estará disponible en `http://localhost:3000`

3. **Ejecutar el backend:**
```bash
cd api
go run cmd/main.go
```
El backend estará disponible en `http://localhost:8080`

4. **Iniciar CockroachDB:**
```bash
# Instalar CockroachDB según tu sistema operativo
# Luego iniciar en modo desarrollo:
cockroach start-single-node --insecure --http-addr=localhost:8081
```

## Estructura del Proyecto

```
.
├── api/                          # Backend en Golang (DDD)
│   ├── cmd/                      # Punto de entrada
│   ├── internal/
│   │   ├── domain/               # Capa de Dominio (DDD)
│   │   │   ├── stock/            # Entidades y servicios de dominio
│   │   │   └── recommendation/   # Algoritmo de recomendación
│   │   ├── application/          # Capa de Aplicación
│   │   │   ├── handlers/         # GraphQL handlers
│   │   │   ├── services/         # Servicios de aplicación
│   │   │   └── graphql/          # Schema y resolvers
│   │   └── infrastructure/       # Capa de Infraestructura
│   │       ├── database/         # Conexión a CockroachDB
│   │       ├── external/         # Cliente API externa
│   │       └── repository/       # Implementación de repositorios
│   └── docs/                     # Documentación API (Swagger)
├── frontend/                     # Frontend en Vue 3
│   ├── src/
│   │   ├── design-system/        # Design System con componentes
│   │   │   ├── components/       # Componentes reusables
│   │   │   ├── tokens/           # Design tokens
│   │   │   └── themes/           # Temas (light/dark)
│   │   ├── hoc/                  # Higher Order Components
│   │   ├── views/                # Vistas/páginas
│   │   ├── stores/               # Stores de Pinia
│   │   └── composables/          # Composables Vue
│   └── .storybook/               # Configuración Storybook
├── docs/                         # Documentación general
│   ├── ARCHITECTURE.md           # Arquitectura DDD
│   ├── ALGORITHMS.md             # Algoritmos y optimizaciones
│   └── RESUMEN_EJECUTIVO.md      # Resumen ejecutivo
├── .devcontainer/                # Configuración del devcontainer
│   ├── devcontainer.json
│   ├── docker-compose.yml
│   └── README.md
├── PLAN_DE_ACCION.md             # Plan de acción detallado
└── README.md
```

## Desarrollo

### Inicio Automático de Servicios

Al abrir el devcontainer, todos los servicios se inician automáticamente:

1. **CockroachDB**: Se inicia primero (el API depende de él)
2. **Frontend (Vue 3)**: Se inicia automáticamente con `npm run dev`
3. **Backend (Go API)**: Se inicia automáticamente con Air para hot-reload

El script `postCreate.sh` se ejecuta automáticamente para:
- Instalar herramientas de Go (gopls, delve, air)
- Instalar dependencias del frontend (`npm install`)
- Descargar módulos de Go (`go mod download`)
- Esperar a que CockroachDB esté listo

### Puertos
- **Frontend**: Puerto 3000 - `http://localhost:3000`
- **Backend**: Puerto 8080 - `http://localhost:8080`
- **CockroachDB SQL**: Puerto 26257
- **CockroachDB Web UI**: Puerto 8081 - `http://localhost:8081`

### Verificar que los servicios están corriendo

```bash
# Verificar backend
curl http://localhost:8080/health
# Debería responder: OK

# Verificar endpoint de prueba del backend
curl http://localhost:8080/hello
# Debería responder: {"message":"Hello from Go API"}

# Verificar frontend (desde el navegador)
# Abre: http://localhost:3000
```

### Conexión a CockroachDB

**Desde el devcontainer:**
```
postgresql://root@cockroachdb:26257/defaultdb?sslmode=disable
```

**Desde localhost:**
```
postgresql://root@localhost:26257/defaultdb?sslmode=disable
```

### Comandos Útiles

**En el devcontainer:**

> **Nota:** Los servicios (API, Frontend y CockroachDB) se inician automáticamente al abrir el devcontainer. No es necesario ejecutarlos manualmente.

```bash
# Verificar que los servicios están corriendo
curl http://localhost:8080/health  # Backend
curl http://localhost:3000          # Frontend

# Si necesitas reiniciar el API manualmente (ya está corriendo automáticamente)
cd /workspace/api && air -c .air.toml

# Acceder a CockroachDB SQL shell
docker exec -it cockroachdb ./cockroach sql --insecure

# Ver logs de CockroachDB
docker logs cockroachdb

# Ver logs del contenedor de la API
docker logs <container-id-api>  # Reemplaza con el ID del contenedor

# Ver logs del contenedor del frontend
docker logs <container-id-frontend>  # Reemplaza con el ID del contenedor
```

## 🚀 Inicio Rápido

### Desarrollo del Sistema de Acciones

1. **Revisar la documentación**:
   - Leer [Plan de Acción](./PLAN_DE_ACCION.md) para entender las fases
   - Revisar [Arquitectura](./docs/ARCHITECTURE.md) para entender la estructura DDD
   - Consultar [Algoritmos](./docs/ALGORITHMS.md) para entender las optimizaciones

2. **Configurar el entorno**:
   - El devcontainer ya está configurado con todos los servicios
   - CockroachDB se inicia automáticamente
   - Backend y Frontend tienen hot reload

3. **Comenzar implementación**:
   - Seguir las fases del plan de acción
   - Empezar por Fase 1: Backend - Infraestructura

### API Externa

El sistema se conecta a:
- **Endpoint**: `https://api.karenai.click/swechallenge/list`
- **Auth**: Bearer token (ver documentación)
- **Paginación**: Usar parámetro `next_page`

## Características

### Sistema de Acciones
- Sincronización automática desde API externa
- Almacenamiento en CockroachDB
- API GraphQL con queries y mutations
- Interfaz moderna con Design System
- Algoritmo de recomendación O(n log n)
- Optimizaciones de performance (cache, debounce, etc.)

### Demo Actual
El proyecto incluye una página demo que muestra:
- Integración entre el frontend Vue 3 y el backend Go
- Manejo de estado con Pinia
- Diseño moderno con Tailwind CSS
- Conexión a CockroachDB
- Hot reload para desarrollo rápido

## Solución de Problemas

### Los servicios no se inician automáticamente

Si los servicios no se inician automáticamente al abrir el devcontainer:

1. **Reconstruir el contenedor:**
   - Presiona `Ctrl+Shift+P` (o `Cmd+Shift+P` en Mac)
   - Ejecuta: `Dev Containers: Rebuild Container`

2. **Verificar que los contenedores estén corriendo:**
   ```bash
   # Desde fuera del contenedor (en tu máquina local)
   docker-compose -f .devcontainer/docker-compose.yml ps
   ```

3. **Iniciar servicios manualmente (si es necesario):**
   ```bash
   # Backend (normalmente ya está corriendo)
   cd /workspace/api && air -c .air.toml
   
   # Frontend (normalmente ya está corriendo)
   cd /workspace/frontend && npm run dev
   ```

### Error: Puerto ya en uso

Si recibes un error de que el puerto está en uso:

1. **Verificar qué proceso está usando el puerto:**
   ```bash
   # En Linux/Mac
   lsof -i :3000  # Para frontend
   lsof -i :8080  # Para backend
   ```

2. **Detener procesos duplicados:**
   ```bash
   pkill -f "vite"      # Para frontend
   pkill -f "go run"    # Para backend
   ```

### El backend no compila

Si el backend tiene errores de compilación:

1. **Verificar que las dependencias estén instaladas:**
   ```bash
   cd /workspace/api
   go mod download
   go mod tidy
   ```

2. **Verificar errores de sintaxis:**
   ```bash
   cd /workspace/api
   go build ./cmd/main.go
   ```

## Extensiones de VS Code Recomendadas

- Go (golang.Go)
- Vue Language Features (Vue.volar)
- TypeScript Vue Plugin (Vue.vscode-typescript-vue-plugin)
- ESLint
- Prettier
- Tailwind CSS IntelliSense

Todas estas extensiones están preconfiguradas en el devcontainer.
