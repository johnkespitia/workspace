# Debug de Query GraphQL

## Problema Reportado

Al ejecutar la query con filtros y ordenamiento, se obtiene un error.

## Query que causa el problema

```graphql
query GetStocks($filter: StockFilter, $sort: StockSort, $limit: Int, $offset: Int) {
  stocks(filter: $filter, sort: $sort, limit: $limit, offset: $offset) {
    stocks {
      id
      ticker
      companyName
      ...
    }
    totalCount
    pageInfo {
      hasNextPage
      hasPreviousPage
    }
  }
}
```

Variables:

```json
{
  "filter": {
    "rating": ["Speculative Buy"]
  },
  "limit": 50,
  "offset": 0,
  "sort": {
    "direction": "asc",
    "field": "ticker"
  }
}
```

## Problemas Identificados y Corregidos

### 1. Filtro "rating" vs "ratings"

**Problema**: El schema espera `ratings` (plural) pero se está usando `rating` (singular).

**Solución**: Agregado soporte para ambos:

- `ratings` (plural) - formato correcto
- `rating` (singular) - compatibilidad

### 2. Sort Field como Enum

**Problema**: El `field` viene como enum (`TICKER`, `COMPANY_NAME`, etc.) pero el código esperaba un string directo.

**Solución**: Agregadas funciones de mapeo:

- `mapEnumFieldToDBField()` - convierte `TICKER` → `ticker`
- `mapEnumDirectionToDBDirection()` - convierte `ASC` → `asc`

## Cómo Probar

### Opción 1: Usar "ratings" (plural) - Recomendado

```json
{
  "filter": {
    "ratings": ["Speculative Buy"]
  },
  "sort": {
    "field": "TICKER",
    "direction": "ASC"
  }
}
```

### Opción 2: Usar "rating" (singular) - Compatibilidad

```json
{
  "filter": {
    "rating": ["Speculative Buy"]
  },
  "sort": {
    "field": "TICKER",
    "direction": "ASC"
  }
}
```

## Valores de Enum Válidos

### StockSortField

- `TICKER`
- `COMPANY_NAME`
- `RATING_TO`
- `TARGET_TO`
- `CREATED_AT`

### SortDirection

- `ASC`
- `DESC`

## Debugging

Si aún hay errores:

1. **Verificar el error exacto**: Revisa la respuesta completa en el playground
2. **Revisar logs del servidor**: Puede haber errores de SQL o validación
3. **Probar sin filtros**: Primero prueba sin filtros para verificar que la query básica funciona

```graphql
query {
  stocks(limit: 10) {
    stocks {
      ticker
      companyName
    }
    totalCount
  }
}
```

4. **Probar con filtro simple**:

```graphql
query {
  stocks(filter: { ratings: ["Buy"] }, limit: 10) {
    stocks {
      ticker
      companyName
      ratingTo
    }
    totalCount
  }
}
```
