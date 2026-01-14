# Documentación Completa de la API

## 📋 Índice

1. [Endpoints HTTP](#endpoints-http)
2. [GraphQL API](#graphql-api)
3. [Ejemplos de Uso](#ejemplos-de-uso)
4. [Guía de Integración](#guía-de-integración)
5. [Códigos de Error](#códigos-de-error)
6. [Rate Limiting](#rate-limiting)
7. [Mejores Prácticas](#mejores-prácticas)

---

## 🌐 Endpoints HTTP

### GET /docs

**Descripción**: Página principal de documentación con enlaces a todas las herramientas de documentación.

**Request**:

```http
GET /docs HTTP/1.1
Host: localhost:8080
```

**Response** (200 OK):

- HTML con página de documentación que incluye:
  - Enlace a GraphQL Playground
  - Enlace a Swagger UI
  - Enlaces a documentación completa
  - Lista de todos los endpoints disponibles

**Uso**: Punto de entrada principal para acceder a toda la documentación de la API.

---

### GET /docs/swagger

**Descripción**: Interfaz Swagger UI para explorar la especificación OpenAPI de forma interactiva.

**Request**:

```http
GET /docs/swagger HTTP/1.1
Host: localhost:8080
```

**Response** (200 OK):

- HTML con Swagger UI cargando automáticamente `/docs/openapi.yaml`

**Características**:

- Documentación interactiva de endpoints REST
- Pruebas de endpoints directamente desde el navegador
- Esquemas de request/response
- Ejemplos de código

---

### GET /docs/openapi.yaml

**Descripción**: Especificación OpenAPI 3.0 en formato YAML.

**Request**:

```http
GET /docs/openapi.yaml HTTP/1.1
Host: localhost:8080
```

**Response** (200 OK):

- Archivo YAML con la especificación OpenAPI completa

**Uso**:

- Importar en herramientas como Swagger UI, Postman, Insomnia
- Generar clientes automáticamente
- Integrar con herramientas de documentación

---

### GET /health

**Descripción**: Verifica el estado del servidor.

**Request**:

```http
GET /health HTTP/1.1
Host: localhost:8080
```

**Response** (200 OK):

```
OK
```

**Uso**:

- Monitoreo de salud del servicio
- Health checks de load balancers
- Verificación rápida de disponibilidad

---

### POST /query

**Descripción**: Endpoint principal para ejecutar queries y mutations GraphQL.

**Request Headers**:

```http
Content-Type: application/json
```

**Request Body**:

```json
{
  "query": "query { stocks(limit: 10) { stocks { ticker } } }",
  "variables": {},
  "operationName": "GetStocks"
}
```

**Response** (200 OK):

```json
{
  "data": {
    "stocks": {
      "stocks": [
        {
          "ticker": "AAPL"
        }
      ],
      "totalCount": 1
    }
  }
}
```

**Códigos de Estado**:

- `200 OK`: Request procesado (puede contener errores en el body)
- `400 Bad Request`: Request inválido
- `405 Method Not Allowed`: Método no permitido (solo POST)
- `408 Request Timeout`: Operación muy larga

---

### GET /playground

**Descripción**: Interfaz visual interactiva para explorar el schema GraphQL.

**Request**:

```http
GET /playground HTTP/1.1
Host: localhost:8080
```

**Response** (200 OK):

- HTML con interfaz del GraphQL Playground

**Características**:

- Explorador de schema interactivo
- Autocompletado de queries
- Validación en tiempo real
- Historial de queries

---

## 🔷 GraphQL API

### Schema Completo

Ver documentación detallada en: [GRAPHQL_API_REFERENCE.md](../docs/GRAPHQL_API_REFERENCE.md)

### Tipos Principales

#### Stock

```graphql
type Stock {
  id: ID!
  ticker: String!
  companyName: String!
  brokerage: String
  action: String
  ratingFrom: String!
  ratingTo: String!
  targetFrom: Float!
  targetTo: Float!
  createdAt: DateTime!
  updatedAt: DateTime!
}
```

#### StockConnection

```graphql
type StockConnection {
  stocks: [Stock!]!
  totalCount: Int!
  pageInfo: PageInfo!
}
```

#### Recommendation

```graphql
type Recommendation {
  stock: Stock!
  score: Float!
  priceChange: Float!
  ratingScore: Float!
  actionScore: Float!
}
```

### Queries Disponibles

#### stocks

Obtiene una lista de stocks con filtros, ordenamiento y paginación.

```graphql
query GetStocks(
  $filter: StockFilter
  $sort: StockSort
  $limit: Int
  $offset: Int
) {
  stocks(filter: $filter, sort: $sort, limit: $limit, offset: $offset) {
    stocks {
      id
      ticker
      companyName
      ratingTo
      targetTo
    }
    totalCount
    pageInfo {
      hasNextPage
      hasPreviousPage
    }
  }
}
```

#### stock

Obtiene un stock específico por ticker.

```graphql
query GetStock($ticker: String!) {
  stock(ticker: $ticker) {
    id
    ticker
    companyName
    brokerage
    action
    ratingFrom
    ratingTo
    targetFrom
    targetTo
    createdAt
    updatedAt
  }
}
```

#### recommendations

Obtiene recomendaciones de inversión basadas en el algoritmo.

```graphql
query GetRecommendations($limit: Int) {
  recommendations(limit: $limit) {
    stock {
      ticker
      companyName
      ratingTo
      targetTo
    }
    score
    priceChange
    ratingScore
    actionScore
  }
}
```

### Mutations Disponibles

#### syncStocks

Sincroniza stocks desde la API externa.

```graphql
mutation SyncStocks {
  syncStocks {
    success
    message
    stocksSynced
  }
}
```

---

## 📝 Ejemplos de Uso

### Ejemplo 1: Obtener stocks con filtros

**Query**:

```graphql
query GetFilteredStocks {
  stocks(
    filter: { ratings: ["Buy", "Strong Buy"], action: "target raised by" }
    sort: { field: TICKER, direction: ASC }
    limit: 50
    offset: 0
  ) {
    stocks {
      ticker
      companyName
      ratingTo
      targetTo
    }
    totalCount
    pageInfo {
      hasNextPage
      hasPreviousPage
    }
  }
}
```

**cURL**:

```bash
curl -X POST http://localhost:8080/query \
  -H "Content-Type: application/json" \
  -d '{
    "query": "query { stocks(filter: { ratings: [\"Buy\"] }, limit: 10) { stocks { ticker companyName } totalCount } }"
  }'
```

**JavaScript (fetch)**:

```javascript
const response = await fetch("http://localhost:8080/query", {
  method: "POST",
  headers: {
    "Content-Type": "application/json",
  },
  body: JSON.stringify({
    query: `
      query GetStocks {
        stocks(filter: { ratings: ["Buy"] }, limit: 10) {
          stocks {
            ticker
            companyName
            ratingTo
            targetTo
          }
          totalCount
        }
      }
    `,
  }),
});

const data = await response.json();
console.log(data);
```

---

### Ejemplo 2: Sincronizar stocks

**Mutation**:

```graphql
mutation {
  syncStocks {
    success
    message
    stocksSynced
  }
}
```

**cURL**:

```bash
curl -X POST http://localhost:8080/query \
  -H "Content-Type: application/json" \
  -d '{
    "query": "mutation { syncStocks { success message stocksSynced } }"
  }'
```

**Nota**: Esta operación puede tardar varios segundos dependiendo de la cantidad de stocks a sincronizar.

---

### Ejemplo 3: Obtener recomendaciones

**Query**:

```graphql
query GetRecommendations {
  recommendations(limit: 10) {
    stock {
      ticker
      companyName
      ratingTo
      targetTo
    }
    score
    priceChange
    ratingScore
    actionScore
  }
}
```

**cURL**:

```bash
curl -X POST http://localhost:8080/query \
  -H "Content-Type: application/json" \
  -d '{
    "query": "query { recommendations(limit: 10) { stock { ticker companyName } score } }"
  }'
```

---

### Ejemplo 4: Búsqueda por ticker

**Query**:

```graphql
query GetStockByTicker {
  stock(ticker: "AAPL") {
    id
    ticker
    companyName
    brokerage
    action
    ratingFrom
    ratingTo
    targetFrom
    targetTo
    createdAt
    updatedAt
  }
}
```

**cURL**:

```bash
curl -X POST http://localhost:8080/query \
  -H "Content-Type: application/json" \
  -d '{
    "query": "query { stock(ticker: \"AAPL\") { ticker companyName ratingTo targetTo } }"
  }'
```

---

## 🔌 Guía de Integración

### Configuración Inicial

1. **Configurar variables de entorno**:

```bash
DB_HOST=localhost
DB_PORT=26257
DB_USER=root
DB_PASSWORD=
DB_NAME=stocks
DB_SSLMODE=disable
API_BASE_URL=https://api.karenai.click
API_KEY=tu_api_key
PORT=8080
```

2. **Iniciar el servidor**:

```bash
go run ./cmd/main.go
```

3. **Verificar que funciona**:

```bash
curl http://localhost:8080/health
```

### Cliente GraphQL en JavaScript/TypeScript

```typescript
// Cliente GraphQL simple
async function graphqlQuery(query: string, variables?: Record<string, any>) {
  const response = await fetch("http://localhost:8080/query", {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify({
      query,
      variables,
    }),
  });

  if (!response.ok) {
    throw new Error(`HTTP error! status: ${response.status}`);
  }

  const data = await response.json();

  if (data.errors) {
    throw new Error(`GraphQL errors: ${JSON.stringify(data.errors)}`);
  }

  return data.data;
}

// Uso
const result = await graphqlQuery(`
  query {
    stocks(limit: 10) {
      stocks {
        ticker
        companyName
      }
      totalCount
    }
  }
`);
```

### Cliente GraphQL en Python

```python
import requests

def graphql_query(query, variables=None):
    response = requests.post(
        'http://localhost:8080/query',
        json={
            'query': query,
            'variables': variables or {}
        },
        headers={'Content-Type': 'application/json'}
    )
    response.raise_for_status()
    data = response.json()

    if 'errors' in data:
        raise Exception(f"GraphQL errors: {data['errors']}")

    return data['data']

# Uso
result = graphql_query("""
    query {
        stocks(limit: 10) {
            stocks {
                ticker
                companyName
            }
            totalCount
        }
    }
""")
```

---

## ⚠️ Códigos de Error

### HTTP Status Codes

| Código | Descripción           | Solución                                         |
| ------ | --------------------- | ------------------------------------------------ |
| 200    | OK                    | Request procesado correctamente                  |
| 400    | Bad Request           | Verificar formato del request                    |
| 405    | Method Not Allowed    | Usar POST para /query                            |
| 408    | Request Timeout       | Operación muy larga, considerar aumentar timeout |
| 500    | Internal Server Error | Error del servidor, revisar logs                 |

### GraphQL Errors

Los errores GraphQL se retornan en el campo `errors` del response:

```json
{
  "errors": [
    {
      "message": "Cannot query field 'invalidField' on type 'Stock'",
      "locations": [
        {
          "line": 2,
          "column": 5
        }
      ]
    }
  ]
}
```

**Errores comunes**:

- `Cannot query field 'X' on type 'Y'`: Campo no existe en el tipo
- `Variable "$X" got invalid value`: Variable con tipo incorrecto
- `Unknown field 'X' in input`: Campo no válido en input
- `Expected type X, found Y`: Tipo incorrecto

---

## 🚦 Rate Limiting

Actualmente no hay rate limiting implementado en los endpoints HTTP/GraphQL.

**Nota**: El cliente de API externa tiene rate limiting configurado (10 requests/segundo) para evitar sobrecargar la API externa.

---

## 💡 Mejores Prácticas

### 1. Usar Variables en Queries

✅ **Bueno**:

```graphql
query GetStocks($limit: Int, $ratings: [String!]) {
  stocks(filter: { ratings: $ratings }, limit: $limit) {
    stocks {
      ticker
    }
  }
}
```

❌ **Evitar**:

```graphql
query {
  stocks(filter: { ratings: ["Buy"] }, limit: 10) {
    stocks {
      ticker
    }
  }
}
```

### 2. Paginación

Siempre usar paginación para queries grandes:

```graphql
query {
  stocks(limit: 50, offset: 0) {
    stocks {
      ticker
    }
    pageInfo {
      hasNextPage
    }
  }
}
```

### 3. Manejo de Errores

Siempre verificar el campo `errors` en la respuesta:

```typescript
const response = await fetch('/query', { ... });
const data = await response.json();

if (data.errors) {
  console.error('GraphQL errors:', data.errors);
  // Manejar errores
}

// Usar data.data
```

### 4. Timeouts

Configurar timeouts apropiados para operaciones largas:

```typescript
const controller = new AbortController();
const timeoutId = setTimeout(() => controller.abort(), 60000); // 60s

try {
  const response = await fetch("/query", {
    signal: controller.signal,
    // ...
  });
} catch (error) {
  if (error.name === "AbortError") {
    console.error("Request timeout");
  }
} finally {
  clearTimeout(timeoutId);
}
```

### 5. Caching

Implementar caching en el cliente para queries frecuentes:

```typescript
const cache = new Map();

async function cachedQuery(query: string, ttl: number = 60000) {
  const key = query;
  const cached = cache.get(key);

  if (cached && Date.now() - cached.timestamp < ttl) {
    return cached.data;
  }

  const data = await graphqlQuery(query);
  cache.set(key, { data, timestamp: Date.now() });
  return data;
}
```

---

## 📚 Recursos Adicionales

- [GraphQL API Reference](../docs/GRAPHQL_API_REFERENCE.md) - Referencia completa del schema
- [GraphQL Examples](../docs/GRAPHQL_EXAMPLES.md) - Más ejemplos de queries
- [OpenAPI Specification](./openapi.yaml) - Especificación OpenAPI 3.0
- [Architecture Documentation](../docs/ARCHITECTURE.md) - Arquitectura del sistema

---

## 🔄 Changelog

### v1.0.0 (2024-01-15)

- ✅ Implementación inicial de GraphQL API
- ✅ Endpoints HTTP básicos
- ✅ GraphQL Playground
- ✅ Health check endpoint
- ✅ Documentación completa

---

**Última actualización**: 2024-01-15
