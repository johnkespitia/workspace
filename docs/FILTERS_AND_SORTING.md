# Validación de Filtros y Ordenamiento GraphQL

## Resumen

Se han validado y corregido todos los filtros y ordenamientos para que funcionen sin errores. Las mejoras incluyen:

1. **Sort opcional**: El `sort` ahora es opcional (no requiere `NewNonNull`)
2. **Manejo robusto de tipos**: Soporte para `int`, `int32`, `int64`, `float64`
3. **Validación de campos**: Validación de campos de sort y dirección
4. **Valores por defecto**: Valores por defecto seguros cuando faltan parámetros
5. **Compatibilidad**: Soporte para `rating` (singular) y `ratings` (plural)
6. **Enums flexibles**: Acepta enums en mayúsculas (`TICKER`, `ASC`) y minúsculas (`ticker`, `asc`)

## Filtros Disponibles

### StockFilter

Todos los filtros son **opcionales**:

```graphql
input StockFilter {
  ticker: String # Búsqueda exacta por ticker
  companyName: String # Búsqueda parcial (ILIKE) por nombre de empresa
  ratings: [String!] # Filtro por ratings (array)
  action: String # Búsqueda exacta por acción
}
```

**Nota**: También se acepta `rating` (singular) por compatibilidad.

### Ejemplos de Uso

```graphql
# Filtrar por ticker
query {
  stocks(filter: { ticker: "AAPL" }) {
    stocks {
      ticker
      companyName
    }
  }
}

# Filtrar por nombre de empresa (búsqueda parcial)
query {
  stocks(filter: { companyName: "Apple" }) {
    stocks {
      ticker
      companyName
    }
  }
}

# Filtrar por ratings (plural - recomendado)
query {
  stocks(filter: { ratings: ["Buy", "Strong Buy"] }) {
    stocks {
      ticker
      ratingTo
    }
  }
}

# Filtrar por rating (singular - compatibilidad)
query {
  stocks(filter: { rating: ["Speculative Buy"] }) {
    stocks {
      ticker
      ratingTo
    }
  }
}

# Filtrar por acción
query {
  stocks(filter: { action: "target raised by" }) {
    stocks {
      ticker
      action
    }
  }
}

# Combinar múltiples filtros
query {
  stocks(filter: { ratings: ["Buy"], action: "target raised by" }) {
    stocks {
      ticker
      companyName
      ratingTo
      action
    }
  }
}
```

## Ordenamiento

### StockSort

El `sort` es **opcional**. Si no se proporciona, se usa:

- `field`: `created_at`
- `direction`: `desc`

```graphql
input StockSort {
  field: StockSortField # Campo por el cual ordenar
  direction: SortDirection # Dirección del ordenamiento
}
```

### Campos de Ordenamiento (StockSortField)

Valores válidos (acepta mayúsculas y minúsculas):

- `TICKER` o `ticker` → Ordena por ticker
- `COMPANY_NAME` o `company_name` → Ordena por nombre de empresa
- `RATING_TO` o `rating_to` → Ordena por rating final
- `TARGET_TO` o `target_to` → Ordena por precio objetivo final
- `CREATED_AT` o `created_at` → Ordena por fecha de creación (por defecto)

### Direcciones de Ordenamiento (SortDirection)

Valores válidos (acepta mayúsculas y minúsculas):

- `ASC` o `asc` → Orden ascendente
- `DESC` o `desc` → Orden descendente (por defecto)

### Ejemplos de Uso

```graphql
# Ordenar por ticker ascendente
query {
  stocks(sort: { field: TICKER, direction: ASC }) {
    stocks {
      ticker
      companyName
    }
  }
}

# Ordenar por ticker ascendente (minúsculas)
query {
  stocks(sort: { field: ticker, direction: asc }) {
    stocks {
      ticker
      companyName
    }
  }
}

# Ordenar por nombre de empresa descendente
query {
  stocks(sort: { field: COMPANY_NAME, direction: DESC }) {
    stocks {
      ticker
      companyName
    }
  }
}

# Ordenar por rating final
query {
  stocks(sort: { field: RATING_TO, direction: ASC }) {
    stocks {
      ticker
      ratingTo
    }
  }
}

# Sin sort (usa valores por defecto: created_at DESC)
query {
  stocks {
    stocks {
      ticker
      createdAt
    }
  }
}

# Solo field (direction por defecto: desc)
query {
  stocks(sort: { field: TICKER }) {
    stocks {
      ticker
    }
  }
}
```

## Paginación

```graphql
query {
  stocks(
    limit: 50 # Por defecto: 50
    offset: 0 # Por defecto: 0
  ) {
    stocks {
      ticker
    }
    totalCount
    pageInfo {
      hasNextPage
      hasPreviousPage
    }
  }
}
```

## Ejemplo Completo

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

**Variables**:

```json
{
  "filter": {
    "ratings": ["Buy", "Strong Buy"],
    "action": "target raised by"
  },
  "sort": {
    "field": "TICKER",
    "direction": "ASC"
  },
  "limit": 50,
  "offset": 0
}
```

## Validaciones Implementadas

### 1. Tipos Robustos

- Soporte para `int`, `int32`, `int64`, `float64` en `limit` y `offset`
- Manejo seguro de valores `nil`

### 2. Filtros

- Validación de strings vacíos (se ignoran)
- Soporte para `rating` (singular) y `ratings` (plural)
- Manejo seguro de arrays vacíos

### 3. Ordenamiento

- Validación de campos de sort (solo campos permitidos)
- Validación de direcciones (solo `asc` o `desc`)
- Valores por defecto seguros si el campo no es válido

### 4. Enums

- Acepta valores en mayúsculas (`TICKER`, `ASC`)
- Acepta valores en minúsculas (`ticker`, `asc`)
- Mapeo automático a valores de BD

## Errores Comunes y Soluciones

### Error: "Expected type StockSortField, found ticker"

**Solución**: Usa `TICKER` (mayúsculas) o asegúrate de que el schema tenga el valor en minúsculas (ya está implementado).

### Error: "Expected type SortDirection, found asc"

**Solución**: Usa `ASC` (mayúsculas) o asegúrate de que el schema tenga el valor en minúsculas (ya está implementado).

### Error: Variable "$sort" got invalid value

**Solución**: Verifica que los valores del enum sean correctos. Ahora acepta ambos formatos (mayúsculas y minúsculas).

### Los filtros no funcionan

**Solución**:

- Verifica que los valores no estén vacíos
- Para `ratings`, usa un array: `["Buy"]` no `"Buy"`
- Usa `ratings` (plural) en lugar de `rating` (singular) si es posible

## Testing

Para probar todos los filtros y ordenamientos:

1. **Sin filtros ni sort**:

```graphql
query {
  stocks {
    stocks {
      ticker
    }
    totalCount
  }
}
```

2. **Solo filtros**:

```graphql
query {
  stocks(filter: { ratings: ["Buy"] }) {
    stocks {
      ticker
    }
  }
}
```

3. **Solo sort**:

```graphql
query {
  stocks(sort: { field: TICKER, direction: ASC }) {
    stocks {
      ticker
    }
  }
}
```

4. **Filtros + Sort + Paginación**:

```graphql
query {
  stocks(
    filter: { ratings: ["Buy"] }
    sort: { field: TICKER, direction: ASC }
    limit: 10
    offset: 0
  ) {
    stocks {
      ticker
      companyName
      ratingTo
    }
    totalCount
  }
}
```
