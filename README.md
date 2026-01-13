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

2. **Servicios incluidos:**
   - **API (Go)**: Puerto 8080 con hot reload (Air)
   - **Frontend (Vue 3)**: Puerto 3000 con hot reload (Vite)
   - **CockroachDB**: Puerto 26257 (SQL) y 8081 (Web UI)

3. **Variables de entorno:**
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
├── api/                 # Backend en Golang
│   ├── cmd/
│   ├── app/
│   └── repositories/
├── frontend/            # Frontend en Vue 3
│   ├── src/
│   │   ├── stores/      # Stores de Pinia
│   │   └── App.vue      # Componente principal
│   └── package.json
├── .devcontainer/       # Configuración del devcontainer
│   ├── devcontainer.json
│   ├── docker-compose.yml
│   └── README.md
└── README.md
```

## Desarrollo

### Puertos
- **Frontend**: Puerto 3000
- **Backend**: Puerto 8080
- **CockroachDB SQL**: Puerto 26257
- **CockroachDB Web UI**: Puerto 8081

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
```bash
# Iniciar API con hot reload
cd /workspace/api && air -c .air.toml

# Acceder a CockroachDB SQL shell
docker exec -it cockroachdb ./cockroach sql --insecure

# Ver logs de CockroachDB
docker logs cockroachdb
```

## Características

El proyecto incluye una página demo que muestra:
- Integración entre el frontend Vue 3 y el backend Go
- Manejo de estado con Pinia
- Diseño moderno con Tailwind CSS
- Conexión a CockroachDB
- Hot reload para desarrollo rápido

## Extensiones de VS Code Recomendadas

- Go (golang.Go)
- Vue Language Features (Vue.volar)
- TypeScript Vue Plugin (Vue.vscode-typescript-vue-plugin)
- ESLint
- Prettier
- Tailwind CSS IntelliSense

Todas estas extensiones están preconfiguradas en el devcontainer.
