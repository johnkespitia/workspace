# GraphQL API Reference - Campos Correctos para el Frontend

## Query: `stocks`

Obtiene una lista de stocks con filtros, ordenamiento y paginación.

### Query GraphQL

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
      brokerage
      action
      ratingFrom
      ratingTo
      targetFrom
      targetTo
      createdAt
      updatedAt
    }
    totalCount
    pageInfo {
      hasNextPage
      hasPreviousPage
    }
  }
}
```

## Variables Correctas

### StockFilter

**IMPORTANTE**: Todos los campos son opcionales. Usa `ratings` (plural), NO `rating` (singular).

```typescript
interface StockFilter {
  ticker?: string;           // Búsqueda exacta por ticker
  companyName?: string;       // Búsqueda parcial por nombre de empresa
  ratings?: string[];        // ✅ CORRECTO: Array de ratings (PLURAL)
  action?: string;           // Búsqueda exacta por acción
}
```

### StockSort

**IMPORTANTE**: El `sort` es opcional. Si no se envía, se usa `created_at DESC` por defecto.

```typescript
interface StockSort {
  field?: StockSortField;      // Campo por el cual ordenar
  direction?: SortDirection;   // Dirección del ordenamiento
}

enum StockSortField {
  TICKER = "TICKER",           // ✅ Usar mayúsculas
  COMPANY_NAME = "COMPANY_NAME",
  RATING_TO = "RATING_TO",
  TARGET_TO = "TARGET_TO",
  CREATED_AT = "CREATED_AT"
}

enum SortDirection {
  ASC = "ASC",                 // ✅ Usar mayúsculas
  DESC = "DESC"
}
```

**NOTA**: También se aceptan valores en minúsculas (`ticker`, `asc`), pero se recomienda usar mayúsculas para consistencia.

### Paginación

```typescript
interface Pagination {
  limit?: number;   // Por defecto: 50
  offset?: number;  // Por defecto: 0
}
```

## Ejemplos de Variables Correctas

### 1. Sin filtros ni ordenamiento (valores por defecto)

```json
{
  "filter": null,
  "sort": null,
  "limit": 50,
  "offset": 0
}
```

O simplemente:

```json
{}
```

### 2. Filtrar por ratings (CORRECTO)

```json
{
  "filter": {
    "ratings": ["Speculative Buy"]
  },
  "limit": 50,
  "offset": 0
}
```

**❌ INCORRECTO** (no usar):
```json
{
  "filter": {
    "rating": ["Speculative Buy"]  // ❌ Campo incorrecto
  }
}
```

### 3. Filtrar por múltiples ratings

```json
{
  "filter": {
    "ratings": ["Buy", "Strong Buy", "Speculative Buy"]
  }
}
```

### 4. Filtrar por ticker

```json
{
  "filter": {
    "ticker": "AAPL"
  }
}
```

### 5. Filtrar por nombre de empresa (búsqueda parcial)

```json
{
  "filter": {
    "companyName": "Apple"
  }
}
```

### 6. Filtrar por acción

```json
{
  "filter": {
    "action": "target raised by"
  }
}
```

### 7. Combinar múltiples filtros

```json
{
  "filter": {
    "ratings": ["Buy", "Strong Buy"],
    "action": "target raised by"
  }
}
```

### 8. Ordenar por ticker ascendente

```json
{
  "sort": {
    "field": "TICKER",
    "direction": "ASC"
  }
}
```

### 9. Ordenar por nombre de empresa descendente

```json
{
  "sort": {
    "field": "COMPANY_NAME",
    "direction": "DESC"
  }
}
```

### 10. Filtrar y ordenar juntos

```json
{
  "filter": {
    "ratings": ["Buy"]
  },
  "sort": {
    "field": "TICKER",
    "direction": "ASC"
  },
  "limit": 50,
  "offset": 0
}
```

### 11. Solo ordenar (sin filtros)

```json
{
  "sort": {
    "field": "CREATED_AT",
    "direction": "DESC"
  }
}
```

### 12. Solo paginación (sin filtros ni sort)

```json
{
  "limit": 10,
  "offset": 20
}
```

## Ejemplo Completo para el Frontend (TypeScript/Vue)

```typescript
// Tipos TypeScript
interface StockFilter {
  ticker?: string;
  companyName?: string;
  ratings?: string[];  // ✅ PLURAL
  action?: string;
}

interface StockSort {
  field?: 'TICKER' | 'COMPANY_NAME' | 'RATING_TO' | 'TARGET_TO' | 'CREATED_AT';
  direction?: 'ASC' | 'DESC';
}

interface GetStocksVariables {
  filter?: StockFilter;
  sort?: StockSort;
  limit?: number;
  offset?: number;
}

// Query GraphQL
const GET_STOCKS_QUERY = `
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
        brokerage
        action
        ratingFrom
        ratingTo
        targetFrom
        targetTo
        createdAt
        updatedAt
      }
      totalCount
      pageInfo {
        hasNextPage
        hasPreviousPage
      }
    }
  }
`;

// Ejemplo de uso
const variables: GetStocksVariables = {
  filter: {
    ratings: ['Buy', 'Strong Buy'],  // ✅ PLURAL
  },
  sort: {
    field: 'TICKER',                  // ✅ Mayúsculas
    direction: 'ASC'                  // ✅ Mayúsculas
  },
  limit: 50,
  offset: 0
};
```

## Valores de Rating Válidos

Los ratings válidos que puedes usar en el array `ratings`:

- `"Buy"`
- `"Strong Buy"`
- `"Speculative Buy"`
- `"Neutral"`
- `"Market Perform"`
- `"Sell"`
- `"Strong Sell"`

## Valores de Action Válidos

Algunos ejemplos de acciones válidas:

- `"target raised by"`
- `"target lowered by"`
- `"target raised"`
- `"target lowered"`
- `"initiated coverage"`
- `"speculative buy"`

## Resumen de Campos Correctos

| Campo | Tipo | Requerido | Valores | Notas |
|-------|------|-----------|---------|-------|
| `filter.ticker` | `string` | No | Cualquier string | Búsqueda exacta |
| `filter.companyName` | `string` | No | Cualquier string | Búsqueda parcial (ILIKE) |
| `filter.ratings` | `string[]` | No | Array de ratings | ✅ **PLURAL** - Campo correcto |
| `filter.action` | `string` | No | Cualquier string | Búsqueda exacta |
| `sort.field` | `StockSortField` | No | `TICKER`, `COMPANY_NAME`, `RATING_TO`, `TARGET_TO`, `CREATED_AT` | Recomendado: mayúsculas |
| `sort.direction` | `SortDirection` | No | `ASC`, `DESC` | Recomendado: mayúsculas |
| `limit` | `number` | No | Entero positivo | Por defecto: 50 |
| `offset` | `number` | No | Entero positivo | Por defecto: 0 |

## Errores Comunes

### ❌ Error: "Unknown field 'rating'"
**Causa**: Usar `rating` (singular) en lugar de `ratings` (plural)
**Solución**: Cambiar a `ratings`

```json
// ❌ Incorrecto
{ "filter": { "rating": ["Buy"] } }

// ✅ Correcto
{ "filter": { "ratings": ["Buy"] } }
```

### ❌ Error: "Expected type StockSortField, found ticker"
**Causa**: Usar minúsculas en el enum
**Solución**: Usar mayúsculas (aunque también se aceptan minúsculas)

```json
// ✅ Recomendado
{ "sort": { "field": "TICKER", "direction": "ASC" } }

// ✅ También funciona
{ "sort": { "field": "ticker", "direction": "asc" } }
```

### ❌ Error: "Expected type SortDirection, found asc"
**Causa**: Usar minúsculas en el enum
**Solución**: Usar mayúsculas (aunque también se aceptan minúsculas)

```json
// ✅ Recomendado
{ "sort": { "direction": "ASC" } }

// ✅ También funciona
{ "sort": { "direction": "asc" } }
```

## Testing en GraphQL Playground

1. Abre el playground en: `http://localhost:8080/playground`
2. Usa la query `GetStocks` de arriba
3. En "Query Variables", pega:

```json
{
  "filter": {
    "ratings": ["Speculative Buy"]
  },
  "sort": {
    "field": "TICKER",
    "direction": "ASC"
  },
  "limit": 50,
  "offset": 0
}
```

4. Ejecuta la query
