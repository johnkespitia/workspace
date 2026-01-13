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
docker exec -it cockroachdb ./cockroach sql --insecure
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
docker logs cockroachdb
```

### Reiniciar CockroachDB
```bash
docker-compose restart cockroachdb
```
