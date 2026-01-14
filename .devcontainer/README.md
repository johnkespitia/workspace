# Dev Container Configuration

Este devcontainer está configurado para trabajar con el siguiente stack:

## Stack Tecnológico

- **Backend**: Golang
- **Frontend**: Vue 3 + TypeScript + Pinia + Tailwind CSS
- **Base de Datos**: CockroachDB
- **Build Tool**: Vite

## Servicios

### API (Go)
- Puerto: 8080
- Hot reload con Air
- Ubicación: `/workspace/api`

### Frontend (Vue 3)
- Puerto: 3000
- Hot reload con Vite
- Ubicación: `/workspace/frontend`

### CockroachDB
- Puerto SQL: 26257
- Puerto Web UI: 8081
- Base de datos inicial: `defaultdb`
- Modo: Insecure (solo para desarrollo)

## Conexión a CockroachDB

### Desde Go
```go
import (
    "database/sql"
    _ "github.com/lib/pq"
)

db, err := sql.Open("postgres", "postgresql://root@cockroachdb:26257/defaultdb?sslmode=disable")
```

### Variables de Entorno
Las siguientes variables están disponibles en el contenedor:
- `DATABASE_URL`: URL completa de conexión
- `COCKROACH_HOST`: Host de CockroachDB (cockroachdb)
- `COCKROACH_PORT`: Puerto (26257)
- `COCKROACH_USER`: Usuario (root)
- `COCKROACH_DB`: Base de datos (defaultdb)

## Acceso a CockroachDB

### SQL Shell
```bash
docker exec -it go-react-test-cockroachdb ./cockroach sql --insecure
```

### Web UI
Abre en tu navegador: http://localhost:8081

## Extensiones de VS Code

- Go (golang.Go)
- Vue Language Features (Vue.volar)
- TypeScript Vue Plugin (Vue.vscode-typescript-vue-plugin)
- ESLint
- Prettier
- Tailwind CSS IntelliSense
- TypeScript

## Comandos Útiles

### Iniciar API con hot reload
```bash
cd /workspace/api
air -c .air.toml
```

### Instalar dependencias del frontend
```bash
cd /workspace/frontend
npm install
```

### Ver logs de CockroachDB
```bash
docker logs go-react-test-cockroachdb
```

### Reiniciar CockroachDB
```bash
cd .devcontainer && docker-compose -p go-react-test-devcontainer restart cockroachdb
```

## Abrir IDE con Dev Container

El proyecto incluye un script para abrir automáticamente Cursor o VS Code con el devcontainer:

### Desde el Makefile (recomendado)
```bash
# Inicializar y abrir IDE automáticamente
make dev-init

# Solo abrir el IDE (si ya está inicializado)
make dev-open
```

### Configurar IDE Preferido

Puedes configurar tu IDE preferido mediante variable de entorno:

```bash
# Usar Cursor
export IDE=cursor
make dev-init

# Usar VS Code
export IDE=code
make dev-init

# Auto-detectar (usa el primero disponible)
export IDE=auto
make dev-init
```

Si no defines `IDE`, el script te preguntará interactivamente qué IDE deseas usar.

### Manualmente

Si prefieres abrir el IDE manualmente:

1. Abre Cursor o VS Code
2. Abre la carpeta del proyecto
3. Ejecuta: `Dev Containers: Reopen in Container` (Cmd+Shift+P / Ctrl+Shift+P)

**Nota:** Asegúrate de tener instalada la extensión "Dev Containers" en tu IDE.
