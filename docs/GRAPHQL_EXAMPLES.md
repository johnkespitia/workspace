# Ejemplos de Queries GraphQL

Este documento contiene ejemplos de queries y mutations para probar la API GraphQL.

## 🚀 Iniciar el Servidor

```bash
cd api
go run ./cmd/main.go
```

El servidor se iniciará en `http://localhost:8080` (o el puerto configurado en `PORT`).

## 🎮 Acceder al Playground

Abre tu navegador y ve a:

```
http://localhost:8080/playground
```

## 📝 Ejemplos de Queries

### 1. Obtener lista de stocks (básico)

```graphql
query {
  stocks(limit: 10) {
    stocks {
      id
      ticker
      companyName
      ratingTo
      targetTo
      targetFrom
    }
    totalCount
    pageInfo {
      hasNextPage
      hasPreviousPage
    }
  }
}
```

### 2. Obtener stocks con filtros

```graphql
query {
  stocks(filter: { ratings: ["Buy", "Strong Buy"] }, limit: 20) {
    stocks {
      ticker
      companyName
      ratingFrom
      ratingTo
      targetFrom
      targetTo
    }
    totalCount
  }
}
```

### 3. Buscar por ticker

```graphql
query {
  stocks(filter: { ticker: "AAPL" }) {
    stocks {
      ticker
      companyName
      brokerage
      action
      ratingTo
      targetTo
    }
  }
}
```

### 4. Buscar por nombre de compañía

```graphql
query {
  stocks(filter: { companyName: "Apple" }) {
    stocks {
      ticker
      companyName
      ratingTo
      targetTo
    }
  }
}
```

### 5. Ordenar stocks

```graphql
query {
  stocks(sort: { field: TARGET_TO, direction: DESC }, limit: 10) {
    stocks {
      ticker
      companyName
      targetTo
      ratingTo
    }
  }
}
```

### 6. Obtener un stock específico

```graphql
query {
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

### 7. Obtener recomendaciones

```graphql
query {
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

### 8. Paginación

```graphql
query {
  stocks(limit: 5, offset: 0) {
    stocks {
      ticker
      companyName
    }
    totalCount
    pageInfo {
      hasNextPage
      hasPreviousPage
    }
  }
}
```

## 🔄 Ejemplos de Mutations

### 1. Sincronizar stocks desde API externa

```graphql
mutation {
  syncStocks {
    success
    message
    stocksSynced
  }
}
```

**Nota**: Esta mutation puede tardar varios segundos dependiendo de cuántos stocks haya en la API externa.

## 🎯 Queries Complejas

### Combinar filtros y ordenamiento

```graphql
query {
  stocks(
    filter: { ratings: ["Buy", "Strong Buy"], action: "target raised by" }
    sort: { field: TARGET_TO, direction: DESC }
    limit: 20
    offset: 0
  ) {
    stocks {
      ticker
      companyName
      action
      ratingTo
      targetFrom
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

### Obtener recomendaciones con detalles completos

```graphql
query {
  recommendations(limit: 5) {
    stock {
      id
      ticker
      companyName
      brokerage
      action
      ratingFrom
      ratingTo
      targetFrom
      targetTo
    }
    score
    priceChange
    ratingScore
    actionScore
  }
}
```

## ⚠️ Notas Importantes

1. **Primera vez**: Si la base de datos está vacía, primero ejecuta la mutation `syncStocks` para cargar datos.

2. **Variables**: Puedes usar variables en las queries:

```graphql
query GetStocks($limit: Int, $ratings: [String!]) {
  stocks(filter: { ratings: $ratings }, limit: $limit) {
    stocks {
      ticker
      companyName
    }
    totalCount
  }
}
```

Variables JSON:

```json
{
  "limit": 10,
  "ratings": ["Buy", "Strong Buy"]
}
```

3. **Errores**: Si hay errores, aparecerán en el campo `errors` de la respuesta, pero el status HTTP será 200 (esto es estándar en GraphQL).

## 🔍 Campos Disponibles

### Stock

- `id` (ID!)
- `ticker` (String!)
- `companyName` (String!)
- `brokerage` (String)
- `action` (String)
- `ratingFrom` (String!)
- `ratingTo` (String!)
- `targetFrom` (Float!)
- `targetTo` (Float!)
- `createdAt` (DateTime!)
- `updatedAt` (DateTime!)

### Recommendation

- `stock` (Stock!)
- `score` (Float!)
- `priceChange` (Float!)
- `ratingScore` (Float!)
- `actionScore` (Float!)

### StockFilter

- `ticker` (String)
- `companyName` (String)
- `ratings` ([String!])
- `action` (String)

### StockSort

- `field` (StockSortField!): TICKER, COMPANY_NAME, RATING_TO, TARGET_TO, CREATED_AT
- `direction` (SortDirection!): ASC, DESC
