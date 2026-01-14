# Guía de Uso Completa - Stock Information API

## 📖 Introducción

Esta guía te ayudará a entender y usar la API de información de acciones. La API proporciona acceso a datos de stocks, recomendaciones de inversión y sincronización con fuentes externas.

---

## 🚀 Inicio Rápido

### 1. Configuración Inicial

```bash
# 1. Clonar el repositorio
git clone <repository-url>
cd api

# 2. Configurar variables de entorno
cp env.development.example .env.development
# Editar .env.development con tus credenciales

# 3. Instalar dependencias
go mod download

# 4. Ejecutar migraciones
go run ./cmd/migrate -up

# 5. Iniciar servidor
go run ./cmd/main.go
```

### 2. Verificar que Funciona

```bash
# Health check
curl http://localhost:8080/health

# Debería retornar: OK
```

### 3. Acceder a la Documentación

Tienes varias opciones para ver la documentación:

1. **Página Principal de Documentación**:

   - URL: `http://localhost:8080/docs`
   - Página HTML con enlaces a todas las herramientas

2. **GraphQL Playground**:

   - URL: `http://localhost:8080/playground`
   - Interfaz visual para probar queries GraphQL

3. **Swagger UI**:
   - URL: `http://localhost:8080/docs/swagger`
   - Documentación interactiva de la API REST

---

## 📚 Conceptos Básicos

### ¿Qué es GraphQL?

GraphQL es un lenguaje de consulta para APIs que te permite:

- Solicitar exactamente los datos que necesitas
- Obtener múltiples recursos en una sola petición
- Tener un schema tipado y autodocumentado

### Estructura de una Query

```graphql
query NombreDeLaQuery {
  campo {
    subcampo
  }
}
```

### Estructura de una Mutation

```graphql
mutation NombreDeLaMutation {
  operacion {
    resultado
  }
}
```

---

## 🔍 Casos de Uso Comunes

### Caso 1: Ver Lista de Stocks

**Objetivo**: Ver todos los stocks disponibles con información básica.

**Query**:

```graphql
query GetStocks {
  stocks(limit: 20) {
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

**Resultado esperado**:

```json
{
  "data": {
    "stocks": {
      "stocks": [
        {
          "ticker": "AAPL",
          "companyName": "Apple Inc.",
          "ratingTo": "Strong Buy",
          "targetTo": 200.0
        }
      ],
      "totalCount": 150
    }
  }
}
```

---

### Caso 2: Buscar Stocks por Rating

**Objetivo**: Encontrar solo stocks con rating "Buy" o "Strong Buy".

**Query**:

```graphql
query GetBuyStocks {
  stocks(filter: { ratings: ["Buy", "Strong Buy"] }, limit: 50) {
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

---

### Caso 3: Ordenar Stocks por Ticker

**Objetivo**: Ver stocks ordenados alfabéticamente por ticker.

**Query**:

```graphql
query GetSortedStocks {
  stocks(sort: { field: TICKER, direction: ASC }, limit: 50) {
    stocks {
      ticker
      companyName
    }
  }
}
```

---

### Caso 4: Buscar un Stock Específico

**Objetivo**: Obtener información completa de un stock por su ticker.

**Query**:

```graphql
query GetStock {
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

---

### Caso 5: Obtener Recomendaciones de Inversión

**Objetivo**: Ver las mejores recomendaciones basadas en el algoritmo.

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

**Interpretación del Score**:

- **score**: Score total de recomendación (más alto = mejor)
- **priceChange**: Cambio porcentual en precio objetivo
- **ratingScore**: Score basado en el rating
- **actionScore**: Score basado en la acción (target raised/lowered)

---

### Caso 6: Sincronizar Stocks desde API Externa

**Objetivo**: Actualizar la base de datos con los últimos datos de la API externa.

**Mutation**:

```graphql
mutation SyncStocks {
  syncStocks {
    success
    message
    stocksSynced
  }
}
```

**Nota**: Esta operación puede tardar varios segundos o minutos dependiendo de la cantidad de stocks.

**Resultado esperado**:

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

---

## 🎯 Filtros Avanzados

### Filtrar por Múltiples Criterios

```graphql
query GetFilteredStocks {
  stocks(
    filter: { ratings: ["Buy", "Strong Buy"], action: "target raised by" }
    sort: { field: TARGET_TO, direction: DESC }
    limit: 20
  ) {
    stocks {
      ticker
      companyName
      ratingTo
      targetTo
      action
    }
    totalCount
  }
}
```

### Búsqueda por Nombre de Empresa

```graphql
query SearchByCompany {
  stocks(filter: { companyName: "Apple" }) {
    stocks {
      ticker
      companyName
    }
  }
}
```

**Nota**: La búsqueda por `companyName` es parcial (ILIKE), por lo que "Apple" encontrará "Apple Inc.", "Apple Computer", etc.

---

## 📄 Paginación

### Navegación Básica

```graphql
query GetPage1 {
  stocks(limit: 20, offset: 0) {
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

query GetPage2 {
  stocks(limit: 20, offset: 20) {
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

### Implementación de Paginación en Cliente

```typescript
async function getStocksPage(page: number, pageSize: number = 20) {
  const offset = page * pageSize;

  const query = `
    query GetStocksPage($limit: Int, $offset: Int) {
      stocks(limit: $limit, offset: $offset) {
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
  `;

  const response = await fetch("http://localhost:8080/query", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({
      query,
      variables: { limit: pageSize, offset },
    }),
  });

  return response.json();
}
```

---

## 🔧 Integración con Frontend

### Ejemplo con React + Apollo Client

```typescript
import { ApolloClient, InMemoryCache, gql } from "@apollo/client";

const client = new ApolloClient({
  uri: "http://localhost:8080/query",
  cache: new InMemoryCache(),
});

const GET_STOCKS = gql`
  query GetStocks($filter: StockFilter, $sort: StockSort, $limit: Int) {
    stocks(filter: $filter, sort: $sort, limit: $limit) {
      stocks {
        ticker
        companyName
        ratingTo
        targetTo
      }
      totalCount
    }
  }
`;

function StocksList() {
  const { loading, error, data } = useQuery(GET_STOCKS, {
    variables: {
      filter: { ratings: ["Buy"] },
      limit: 20,
    },
  });

  if (loading) return <p>Loading...</p>;
  if (error) return <p>Error: {error.message}</p>;

  return (
    <ul>
      {data.stocks.stocks.map((stock) => (
        <li key={stock.ticker}>
          {stock.ticker} - {stock.companyName}
        </li>
      ))}
    </ul>
  );
}
```

### Ejemplo con Vue 3 + Composition API

```typescript
import { ref, onMounted } from "vue";

const stocks = ref([]);
const loading = ref(false);

async function fetchStocks() {
  loading.value = true;
  try {
    const response = await fetch("http://localhost:8080/query", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({
        query: `
          query {
            stocks(limit: 20) {
              stocks {
                ticker
                companyName
                ratingTo
                targetTo
              }
            }
          }
        `,
      }),
    });

    const { data } = await response.json();
    stocks.value = data.stocks.stocks;
  } catch (error) {
    console.error("Error fetching stocks:", error);
  } finally {
    loading.value = false;
  }
}

onMounted(() => {
  fetchStocks();
});
```

---

## 🐛 Solución de Problemas

### Problema: "Cannot query field 'X' on type 'Y'"

**Causa**: El campo no existe en el tipo especificado.

**Solución**: Verifica el schema GraphQL en el Playground o consulta la documentación.

### Problema: "Variable '$X' got invalid value"

**Causa**: El tipo o formato de la variable es incorrecto.

**Solución**: Verifica que las variables coincidan con los tipos definidos en el schema.

**Ejemplo correcto**:

```json
{
  "filter": {
    "ratings": ["Buy", "Strong Buy"]
  }
}
```

**Ejemplo incorrecto**:

```json
{
  "filter": {
    "rating": "Buy" // ❌ Debe ser "ratings" (plural) y un array
  }
}
```

### Problema: Request Timeout

**Causa**: La operación (especialmente `syncStocks`) está tomando demasiado tiempo.

**Solución**:

1. Aumentar el timeout del cliente
2. Verificar la conexión a la API externa
3. Considerar ejecutar `syncStocks` en background

### Problema: No hay stocks en la respuesta

**Causa**: Puede ser que:

1. No se hayan sincronizado stocks aún
2. Los filtros son muy restrictivos
3. Hay un error en la query

**Solución**:

1. Ejecutar `syncStocks` primero
2. Verificar los filtros
3. Revisar errores en la respuesta GraphQL

---

## 📊 Monitoreo y Logs

### Health Check

```bash
# Verificar estado del servidor
curl http://localhost:8080/health
```

### Logs del Servidor

Los logs se muestran en la consola donde se ejecuta el servidor:

```
Database connected successfully
Database already initialized
Starting server on :8080
```

### Métricas Recomendadas

- Tiempo de respuesta de queries
- Tasa de éxito de `syncStocks`
- Número de stocks sincronizados
- Errores GraphQL

---

## 🔐 Seguridad

### Estado Actual

- ✅ CORS configurado para desarrollo
- ⚠️ Sin autenticación (para producción, agregar)
- ⚠️ Sin rate limiting en endpoints (considerar para producción)

### Recomendaciones para Producción

1. **Autenticación**: Implementar JWT o API keys
2. **Rate Limiting**: Limitar requests por IP/usuario
3. **HTTPS**: Usar siempre HTTPS en producción
4. **Validación**: Validar y sanitizar todos los inputs
5. **Logging**: Implementar logging estructurado
6. **Monitoring**: Configurar alertas y monitoreo

---

## 📞 Soporte

Para más información:

- Consulta la [documentación completa de la API](./API_DOCUMENTATION.md)
- Revisa los [ejemplos de GraphQL](../docs/GRAPHQL_EXAMPLES.md)
- Consulta la [referencia de la API GraphQL](../docs/GRAPHQL_API_REFERENCE.md)

---

**Última actualización**: 2024-01-15
