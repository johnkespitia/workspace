# Sistema de Información de Acciones

Sistema completo para recuperar, almacenar y visualizar información de acciones desde una API externa, con recomendaciones inteligentes de inversión.

**Stack**: Go + Vue 3 + GraphQL + CockroachDB + DDD

## 📚 Documentación

La documentación está organizada para facilitar la transferencia de conocimiento, desde la visión general hasta los detalles técnicos.

### 🎯 Para Empezar

Si eres nuevo en el proyecto, comienza por estos documentos en orden:

1. **[Resumen Ejecutivo](./docs/RESUMEN_EJECUTIVO.md)** - Visión general del proyecto, objetivos, stack tecnológico y características principales
2. **[Manual del Desarrollador](./docs/DEVELOPER_MANUAL.md)** - Guía completa para configurar el entorno, entender la estructura del proyecto y comenzar a desarrollar

### 🏗️ Arquitectura y Diseño

Para entender cómo está construido el sistema:

3. **[Arquitectura](./docs/ARCHITECTURE.md)** - Arquitectura DDD, capas, estructura y principios de diseño
4. **[Algoritmos](./docs/ALGORITHMS.md)** - Algoritmos de recomendación, búsqueda y optimizaciones implementadas

### 💻 Desarrollo Práctico

Documentación para el día a día del desarrollo:

5. **[Frontend](./docs/FRONTEND.md)** - Guía completa del frontend: componentes, HOCs, composables, stores y más
6. **[GraphQL API Reference](./docs/GRAPHQL_API_REFERENCE.md)** - Referencia completa de la API GraphQL: queries, mutations, filtros y ejemplos
7. **[Testing](./docs/TESTING.md)** - Estrategia de testing, cómo escribir tests y ejecutarlos (backend y frontend)

### ⚙️ Configuración e Infraestructura

Para configurar y entender el entorno de desarrollo:

8. **[Dev Container](./docs/DEVCONTAINER.md)** - Guía del dev container: configuración, scripts y solución de problemas
9. **[Infraestructura](./docs/INFRASTRUCTURE.md)** - Servicios, Docker, base de datos y configuración de red
10. **[Makefile](./docs/MAKEFILE.md)** - Todos los comandos disponibles y cómo usarlos

### 📋 Referencia Histórica

11. **[Plan de Acción](./docs/PLAN_DE_ACCION.md)** - Plan detallado por fases de implementación (referencia histórica)

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

#### Inicializar desde el Host (Makefile)

Puedes inicializar el dev container directamente desde el host usando Make:

```bash
# Ver todos los comandos disponibles
make help

# Inicializar el dev container (construye e inicia servicios)
make dev-init

# Iniciar servicios
make dev-up

# Ver estado de los servicios
make dev-status

# Ver logs
make dev-logs

# Detener servicios
make dev-down
```

**Comandos disponibles:**

- `make dev-init` - Inicializa el dev container (construye e inicia servicios) y abre el IDE
- `make dev-up` - Inicia los servicios
- `make dev-down` - Detiene los servicios
- `make dev-rebuild` - Reconstruye las imágenes y reinicia
- `make dev-logs` - Muestra logs de todos los servicios
- `make dev-logs-api` - Logs solo del API
- `make dev-logs-frontend` - Logs solo del Frontend
- `make dev-status` - Estado de los servicios
- `make dev-shell` - Abre shell en el contenedor API
- `make dev-restart` - Reinicia todos los servicios
- `make dev-health` - Verifica el estado de salud
- `make dev-open` - Abre el IDE (Cursor o VS Code) con el devcontainer

#### Configurar IDE Preferido

El comando `make dev-init` abrirá automáticamente tu IDE con el devcontainer. Puedes configurar tu IDE preferido de las siguientes maneras:

**Opción 1: Variable de entorno (recomendado)**

```bash
# Para usar Cursor
export IDE=cursor
make dev-init

# Para usar VS Code
export IDE=code
make dev-init

# Para auto-detectar (usa el primero disponible)
export IDE=auto
make dev-init
```

**Opción 2: Pregunta interactiva**

Si no defines la variable `IDE`, el script te preguntará qué IDE deseas usar:

```bash
make dev-init
# Te mostrará un menú para seleccionar entre Cursor, VS Code o auto-detectar
```

**Opción 3: Abrir IDE manualmente después**

Si prefieres abrir el IDE manualmente después de `dev-init`:

```bash
make dev-init
# ... espera a que termine ...
make dev-open  # Abre el IDE
```

**Nota:** Asegúrate de tener instalada la extensión "Dev Containers" en tu IDE. El IDE detectará automáticamente el devcontainer y te preguntará si deseas abrirlo.

#### Abrir en IDE con Dev Container

**Automáticamente (recomendado):**

- Ejecuta `make dev-init` y el IDE se abrirá automáticamente
- El IDE detectará el devcontainer y te preguntará si deseas abrirlo

**Manualmente:**

1. **Abrir en Dev Container:**

   - Abre Cursor o VS Code con la extensión "Dev Containers" instalada
   - Abre la carpeta del proyecto
   - Ejecuta el comando: `Dev Containers: Reopen in Container` (Cmd+Shift+P / Ctrl+Shift+P)

2. **Servicios incluidos (se inician automáticamente):**

   - **API (Go)**: Puerto 8080 con hot reload (Air)
     - Se inicia automáticamente al abrir el devcontainer
     - Usa Air para recarga automática al modificar archivos Go
     - Disponible en: `http://localhost:8080`
   - **Frontend (Vue 3)**: Puerto 3001 con hot reload (Vite)
     - Se inicia automáticamente al abrir el devcontainer
     - Usa Vite para recarga automática al modificar archivos Vue/TypeScript
     - Disponible en: `http://localhost:3001`
   - **Storybook**: Puerto 6006
     - Se inicia automáticamente junto con el frontend
     - Documentación y pruebas de componentes
     - Disponible en: `http://localhost:6006`
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

El frontend estará disponible en `http://localhost:3001`

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
├── docs/                         # Documentación completa del proyecto
│   ├── RESUMEN_EJECUTIVO.md      # Visión general (empezar aquí)
│   ├── DEVELOPER_MANUAL.md       # Manual del desarrollador
│   ├── ARCHITECTURE.md           # Arquitectura DDD
│   ├── ALGORITHMS.md             # Algoritmos y optimizaciones
│   ├── FRONTEND.md               # Guía del frontend
│   ├── GRAPHQL_API_REFERENCE.md  # Referencia de la API
│   ├── TESTING.md                # Guía de testing
│   ├── DEVCONTAINER.md           # Configuración del dev container
│   ├── INFRASTRUCTURE.md         # Infraestructura y servicios
│   ├── MAKEFILE.md               # Comandos del Makefile
│   └── PLAN_DE_ACCION.md         # Plan de acción (referencia histórica)
├── .devcontainer/                # Configuración del devcontainer
│   ├── devcontainer.json
│   ├── docker-compose.yml
│   └── Dockerfile.*
└── README.md                     # Este archivo
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

- **Frontend**: Puerto 3001 - `http://localhost:3001`
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
# Abre: http://localhost:3001
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
curl http://localhost:3001          # Frontend

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

### Para Nuevos Desarrolladores

Sigue estos pasos para comenzar a trabajar en el proyecto:

1. **Leer la documentación inicial**:

   - Comienza con el [Resumen Ejecutivo](./docs/RESUMEN_EJECUTIVO.md) para entender el proyecto
   - Revisa el [Manual del Desarrollador](./docs/DEVELOPER_MANUAL.md) para configurar tu entorno

2. **Configurar el entorno de desarrollo**:

   ```bash
   # Inicializar el dev container (recomendado)
   make dev-init

   # O seguir la guía en docs/DEVCONTAINER.md
   ```

3. **Entender la arquitectura**:

   - Revisa [Arquitectura](./docs/ARCHITECTURE.md) para entender la estructura DDD
   - Consulta [Algoritmos](./docs/ALGORITHMS.md) para entender las optimizaciones

4. **Comenzar a desarrollar**:
   - Frontend: Consulta [Frontend](./docs/FRONTEND.md) para guías de desarrollo
   - Backend: Consulta [GraphQL API Reference](./docs/GRAPHQL_API_REFERENCE.md) para la API
   - Testing: Revisa [Testing](./docs/TESTING.md) para escribir y ejecutar tests

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

## 🆘 Solución de Problemas

Para problemas comunes, consulta la documentación específica:

- **Problemas con el dev container**: Ver [Dev Container](./docs/DEVCONTAINER.md#solución-de-problemas)
- **Problemas con servicios**: Ver [Infraestructura](./docs/INFRASTRUCTURE.md#mantenimiento)
- **Problemas con comandos**: Ver [Makefile](./docs/MAKEFILE.md#ejemplos-de-uso)

### Comandos de Diagnóstico Rápido

```bash
# Verificar estado de servicios
make dev-status

# Verificar salud de servicios
make dev-health

# Diagnóstico completo
make dev-diagnose

# Ver logs
make dev-logs
```

## 🛠️ Herramientas y Extensiones

### Extensiones Preconfiguradas en Dev Container

Las siguientes extensiones están preconfiguradas automáticamente:

- **Go** (golang.Go) - Soporte para Go
- **Vue Language Features** (Vue.volar) - Soporte para Vue 3
- **TypeScript Vue Plugin** (Vue.vscode-typescript-vue-plugin) - TypeScript en Vue
- **ESLint** - Linting de código
- **Prettier** - Formateo de código
- **Tailwind CSS IntelliSense** - Autocompletado de Tailwind

### Comandos Útiles del Makefile

Para ver todos los comandos disponibles:

```bash
make help
```

**Comandos más usados**:

- `make dev-init` - Inicializar proyecto
- `make dev-status` - Ver estado de servicios
- `make dev-health` - Verificar salud de servicios
- `make dev-logs` - Ver logs en tiempo real

Ver [Makefile](./docs/MAKEFILE.md) para documentación completa de todos los comandos.

---

## 📖 Guía de Lectura Recomendada

**Para nuevos desarrolladores**, sigue este orden de lectura:

1. **[Resumen Ejecutivo](./docs/RESUMEN_EJECUTIVO.md)** (10 min) - Entender qué es el proyecto
2. **[Manual del Desarrollador](./docs/DEVELOPER_MANUAL.md)** (30 min) - Configurar y comenzar
3. **[Arquitectura](./docs/ARCHITECTURE.md)** (20 min) - Entender la estructura
4. **[Frontend](./docs/FRONTEND.md)** o **[GraphQL API Reference](./docs/GRAPHQL_API_REFERENCE.md)** - Según tu área de trabajo

**Para referencia rápida**:

- Comandos: [Makefile](./docs/MAKEFILE.md)
- API: [GraphQL API Reference](./docs/GRAPHQL_API_REFERENCE.md)
- Configuración: [Dev Container](./docs/DEVCONTAINER.md) y [Infraestructura](./docs/INFRASTRUCTURE.md)
