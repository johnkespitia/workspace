# Proceso de Sincronización de Stocks

## 📋 Resumen

El proceso de sincronización obtiene datos de la API externa (`api.karenai.click`) y los guarda en la base de datos local (CockroachDB).

## 🔄 Flujo del Proceso

```
1. Cliente ejecuta mutation syncStocks
   ↓
2. SyncService recibe la petición
   ↓
3. KarenAIClient hace requests paginados a API externa
   ↓
4. Cada página se convierte a entidades de dominio
   ↓
5. StockRepository hace batch upsert en CockroachDB
   ↓
6. Retorna número de stocks sincronizados
```

## 🚀 Cómo Sincronizar

### Paso 1: Ejecutar la Mutation

En el GraphQL Playground (`http://localhost:8080/playground`), ejecuta:

```graphql
mutation {
  syncStocks {
    success
    message
    stocksSynced
  }
}
```

### Paso 2: Esperar la Respuesta

La sincronización puede tardar varios segundos dependiendo de:

- Número de páginas en la API externa
- Velocidad de conexión
- Tamaño de cada página

Respuesta exitosa:

```json
{
  "data": {
    "syncStocks": {
      "success": true,
      "message": "Stocks synchronized successfully",
      "stocksSynced": 150
    }
  }
}
```

Respuesta con error:

```json
{
  "data": {
    "syncStocks": {
      "success": false,
      "message": "error message here",
      "stocksSynced": 0
    }
  }
}
```

### Paso 3: Verificar los Datos

Después de sincronizar, ejecuta una query para verificar:

```graphql
query {
  stocks(limit: 10) {
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

## 🔍 Detalles Técnicos

### 1. Cliente API Externa (`KarenAIClient`)

**Archivo**: `internal/infrastructure/external/karenai_api.go`

**Proceso**:

- Hace requests HTTP GET a `https://api.karenai.click/swechallenge/list`
- Usa paginación con parámetro `next_page`
- Incluye Bearer token en header `Authorization`
- Convierte DTOs de la API a entidades de dominio

**Ejemplo de request**:

```go
GET https://api.karenai.click/swechallenge/list?next_page=abc123
Headers:
  Authorization: Bearer <API_KEY>
```

### 2. Conversión de Datos

**Mapeo de campos**:

- `TICKER` → `stock.Ticker`
- `COMPANY` → `stock.CompanyName`
- `BROKERAGE` → `stock.Brokerage`
- `ACTION` → `stock.Action`
- `RATING FROM` → `stock.RatingFrom`
- `RATING TO` → `stock.RatingTo`
- `TARGET FROM` → `stock.TargetFrom`
- `TARGET TO` → `stock.TargetTo`

### 3. Batch Upsert

**Archivo**: `internal/infrastructure/repository/stock_repository.go`

**Proceso**:

- Agrupa stocks en batches de 100
- Usa `INSERT ... ON CONFLICT (ticker) DO UPDATE` (UPSERT)
- Evita duplicados basándose en el ticker
- Actualiza registros existentes si el ticker ya existe

**Ventajas**:

- Más eficiente que inserts individuales
- Maneja duplicados automáticamente
- Transaccional (todo o nada por batch)

## ⚙️ Configuración Requerida

### Variables de Entorno

Asegúrate de tener configurado en `.env.development`:

```bash
# API Externa
API_BASE_URL=https://api.karenai.click
API_KEY=tu_api_key_aqui  # ⚠️ IMPORTANTE: Debe ser válida
```

### Verificar Configuración

```bash
# Verificar que las variables están cargadas
cd api
go run ./cmd/main.go
# Debería mostrar: "Starting server on :8080"
# Si falta API_KEY, mostrará error al iniciar
```

## 🐛 Troubleshooting

### Error: "API_KEY environment variable is required"

**Solución**: Configura `API_KEY` en tu `.env.development` o exporta la variable:

```bash
export API_KEY=tu_api_key
```

### Error: "failed to fetch stocks from API"

**Posibles causas**:

1. API_KEY inválida o expirada
2. Problemas de red/conectividad
3. API externa no disponible

**Solución**: Verifica la API_KEY y la conectividad.

### Sincronización lenta

**Normal**: La sincronización puede tardar 5-30 segundos dependiendo de:

- Número de páginas
- Velocidad de la API externa
- Tamaño de los datos

### Stocks duplicados

**No debería pasar**: El sistema usa UPSERT basado en `ticker`, por lo que:

- Si un ticker ya existe, se actualiza
- No se crean duplicados

## 📊 Monitoreo

### Ver logs del servidor

El servidor mostrará logs durante la sincronización:

```
Starting server on :8080
Database connected successfully
Database already initialized
```

Durante la sincronización (en el código del servicio):

- Se procesan páginas una por una
- Se muestran errores si ocurren

### Verificar en base de datos (opcional)

```sql
-- Conectar a CockroachDB
-- Ver número de stocks
SELECT COUNT(*) FROM stocks;

-- Ver algunos stocks
SELECT ticker, company_name, rating_to, target_to
FROM stocks
LIMIT 10;
```

## 🔄 Re-sincronización

Puedes ejecutar `syncStocks` múltiples veces:

- Los stocks existentes se actualizarán
- Los nuevos se agregarán
- No se crearán duplicados

**Recomendación**: Ejecutar periódicamente para mantener los datos actualizados.

## 📝 Ejemplo Completo

### 1. Primera sincronización

```graphql
mutation {
  syncStocks {
    success
    message
    stocksSynced
  }
}
```

**Respuesta**:

```json
{
  "data": {
    "syncStocks": {
      "success": true,
      "message": "Stocks synchronized successfully",
      "stocksSynced": 150
    }
  }
}
```

### 2. Verificar datos

```graphql
query {
  stocks(limit: 5) {
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

**Respuesta**:

```json
{
  "data": {
    "stocks": {
      "stocks": [
        {
          "ticker": "AAPL",
          "companyName": "Apple Inc.",
          "ratingTo": "Buy",
          "targetTo": 180.50
        },
        ...
      ],
      "totalCount": 150
    }
  }
}
```

### 3. Re-sincronizar (actualizar datos)

```graphql
mutation {
  syncStocks {
    success
    message
    stocksSynced
  }
}
```

Los stocks existentes se actualizarán con los datos más recientes.
