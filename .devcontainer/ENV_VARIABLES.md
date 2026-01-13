# Variables de Entorno

Este documento describe las variables de entorno disponibles en el devcontainer.

## Base de Datos (CockroachDB)

```bash
DATABASE_URL=postgresql://root@cockroachdb:26257/defaultdb?sslmode=disable
COCKROACH_HOST=cockroachdb
COCKROACH_PORT=26257
COCKROACH_USER=root
COCKROACH_DB=defaultdb
```

## API (Go)

```bash
PORT=8080
```

## Frontend (Vue 3)

```bash
VITE_API_URL=http://localhost:8080
```

## Uso en Go

Para usar estas variables en tu código Go:

```go
import (
    "os"
    "database/sql"
    _ "github.com/lib/pq"
)

func main() {
    databaseURL := os.Getenv("DATABASE_URL")
    if databaseURL == "" {
        databaseURL = "postgresql://root@cockroachdb:26257/defaultdb?sslmode=disable"
    }
    
    db, err := sql.Open("postgres", databaseURL)
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()
}
```

## Uso en Vue 3

Para usar variables de entorno en Vue 3 con Vite, crea un archivo `.env` en la carpeta `frontend`:

```bash
VITE_API_URL=http://localhost:8080
```

Luego accede a ellas en tu código:

```typescript
const apiUrl = import.meta.env.VITE_API_URL
```
