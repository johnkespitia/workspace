# Configuración de Variables de Entorno

Este proyecto requiere variables de entorno para el backend (Go) y frontend (Vue).

## 📋 Archivos de Ejemplo

Se han creado archivos de ejemplo que puedes copiar:

- `env.development.example` - Variables para desarrollo (raíz del proyecto)
- `api/env.development.example` - Variables del backend
- `frontend/env.development.example` - Variables del frontend

## 🚀 Configuración Rápida

### Backend (Go)

1. Copia el archivo de ejemplo:
```bash
cp api/env.development.example api/.env.development
```

2. Edita `api/.env.development` y actualiza los valores, especialmente:
   - `API_KEY` - Tu API key de KarenAI
   - `DB_PASSWORD` - Contraseña de CockroachDB (si aplica)

3. Para cargar las variables en Go, puedes usar `godotenv`:
```bash
go get github.com/joho/godotenv
```

Luego en `cmd/main.go`:
```go
import _ "github.com/joho/godotenv/autoload"
```

O exportar manualmente:
```bash
export $(cat api/.env.development | grep -v '^#' | xargs)
go run ./api/cmd/main.go
```

### Frontend (Vue)

1. Copia el archivo de ejemplo:
```bash
cp frontend/env.development.example frontend/.env.development
```

2. Vite carga automáticamente `.env.development` en modo desarrollo

3. Las variables deben empezar con `VITE_` para ser expuestas al cliente

4. Acceso en código:
```typescript
const endpoint = import.meta.env.VITE_GRAPHQL_ENDPOINT
```

## 📝 Variables Requeridas

### Backend

| Variable | Descripción | Valor por Defecto |
|----------|-------------|-------------------|
| `DB_HOST` | Host de CockroachDB | `localhost` |
| `DB_PORT` | Puerto de CockroachDB | `26257` |
| `DB_USER` | Usuario de la base de datos | `root` |
| `DB_PASSWORD` | Contraseña de la base de datos | (vacío) |
| `DB_NAME` | Nombre de la base de datos | `stocks` |
| `DB_SSLMODE` | Modo SSL | `disable` |
| `API_BASE_URL` | URL base de la API externa | `https://api.karenai.click` |
| `API_KEY` | API key de KarenAI | **(requerido)** |
| `PORT` | Puerto del servidor backend | `8080` |

### Frontend

| Variable | Descripción | Valor por Defecto |
|----------|-------------|-------------------|
| `VITE_GRAPHQL_ENDPOINT` | Endpoint GraphQL del backend | `http://localhost:8080/query` |

## 🔒 Seguridad

- **NUNCA** commitees archivos `.env` o `.env.development` con valores reales
- Los archivos `.env*` están en `.gitignore`
- Usa `.env.example` o `env.development.example` para documentar las variables necesarias
- En producción, usa variables de entorno del sistema o un gestor de secretos

## 🛠️ Desarrollo

Para desarrollo local, puedes:

1. **Usar direnv** (recomendado):
```bash
# Instalar direnv
brew install direnv  # macOS
# o
sudo apt install direnv  # Linux

# Configurar en tu shell
echo 'eval "$(direnv hook bash)"' >> ~/.bashrc

# Crear .envrc en la raíz del proyecto
echo "dotenv api/.env.development" > .envrc
direnv allow
```

2. **Exportar manualmente**:
```bash
export $(cat api/.env.development | grep -v '^#' | xargs)
```

3. **Usar godotenv en Go** (ya implementado en el código si lo agregas)
