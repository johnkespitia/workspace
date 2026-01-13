# Arquitectura del Sistema - Sistema de Información de Acciones

## 🏛️ Principios de Arquitectura

### Domain-Driven Design (DDD)

El sistema está estructurado siguiendo los principios de DDD para garantizar:
- **Desacoplamiento**: Cada capa tiene responsabilidades claras
- **Testabilidad**: Fácil de testear con mocks
- **Mantenibilidad**: Código organizado y fácil de entender
- **Escalabilidad**: Preparado para crecer

### Capas de la Arquitectura

```
┌─────────────────────────────────────────┐
│         Presentation Layer              │
│  (GraphQL Handlers, HTTP Handlers)      │
└─────────────────┬───────────────────────┘
                  │
┌─────────────────▼───────────────────────┐
│      Application Layer                   │
│  (Services, Use Cases, DTOs)            │
└─────────────────┬───────────────────────┘
                  │
┌─────────────────▼───────────────────────┐
│         Domain Layer                     │
│  (Entities, Value Objects, Interfaces)  │
└─────────────────┬───────────────────────┘
                  │
┌─────────────────▼───────────────────────┐
│    Infrastructure Layer                 │
│  (Repositories, External APIs, DB)      │
└─────────────────────────────────────────┘
```

---

## 📦 Estructura de Capas

### 1. Domain Layer (Capa de Dominio)

**Responsabilidad**: Contiene la lógica de negocio pura, sin dependencias externas.

#### Entidades

```go
// domain/stock/entity.go
type Stock struct {
    ID          uuid.UUID
    Ticker      string
    CompanyName string
    Brokerage   string
    Action      string
    RatingFrom  Rating
    RatingTo    Rating
    TargetFrom  Price
    TargetTo    Price
    CreatedAt   time.Time
    UpdatedAt   time.Time
}

// domain/stock/value_objects.go
type Rating string
type Price decimal.Decimal
```

#### Interfaces de Repositorio

```go
// domain/stock/repository.go
type Repository interface {
    Save(ctx context.Context, stock *Stock) error
    FindByID(ctx context.Context, id uuid.UUID) (*Stock, error)
    FindByTicker(ctx context.Context, ticker string) (*Stock, error)
    FindAll(ctx context.Context, filter Filter, sort Sort) ([]*Stock, error)
    Count(ctx context.Context, filter Filter) (int, error)
}
```

#### Servicios de Dominio

```go
// domain/stock/service.go
type Service interface {
    CalculatePriceChange(stock *Stock) decimal.Decimal
    IsRatingUpgrade(stock *Stock) bool
    CalculateRecommendationScore(stock *Stock) float64
}
```

---

### 2. Application Layer (Capa de Aplicación)

**Responsabilidad**: Orquesta los casos de uso y coordina entre capas.

#### Servicios de Aplicación

```go
// application/services/stock_service.go
type StockService struct {
    repo     domain.StockRepository
    domainSvc domain.StockService
}

func (s *StockService) GetStocks(ctx context.Context, filter Filter, sort Sort) ([]*StockDTO, error)
func (s *StockService) GetStock(ctx context.Context, ticker string) (*StockDTO, error)
func (s *StockService) SyncStocks(ctx context.Context) error
```

#### DTOs (Data Transfer Objects)

```go
// application/dto/stock.go
type StockDTO struct {
    ID          string
    Ticker      string
    CompanyName string
    Rating      string
    TargetPrice float64
    PriceChange float64
}
```

---

### 3. Infrastructure Layer (Capa de Infraestructura)

**Responsabilidad**: Implementaciones concretas de interfaces externas.

#### Repositorio

```go
// infrastructure/repository/stock_repository.go
type CockroachStockRepository struct {
    db *sql.DB
}

func (r *CockroachStockRepository) Save(ctx context.Context, stock *domain.Stock) error {
    // Implementación con SQL
}
```

#### Cliente API Externa

```go
// infrastructure/external/karenai_api.go
type KarenAIClient struct {
    httpClient *http.Client
    baseURL    string
    apiKey     string
}

func (c *KarenAIClient) FetchStocks(ctx context.Context, nextPage string) (*APIResponse, error) {
    // Implementación HTTP
}
```

---

### 4. Presentation Layer (Capa de Presentación)

**Responsabilidad**: Maneja las peticiones HTTP/GraphQL y las convierte en llamadas a servicios.

#### GraphQL Handler

```go
// application/handlers/graphql_handler.go
type GraphQLHandler struct {
    stockService *services.StockService
    syncService  *services.SyncService
}

func (h *GraphQLHandler) ResolveStocks(params graphql.ResolveParams) (interface{}, error) {
    // Llamada a servicio de aplicación
}
```

---

## 🔄 Flujo de Datos

### Sincronización de Stocks

```
1. HTTP Request → GraphQL Mutation (syncStocks)
2. Handler → SyncService
3. SyncService → KarenAIClient (API Externa)
4. KarenAIClient → API Externa (paginación)
5. SyncService → StockRepository (guardar en BD)
6. Repository → CockroachDB
7. Response → GraphQL Response
```

### Consulta de Stocks

```
1. GraphQL Query → Handler
2. Handler → StockService
3. StockService → StockRepository
4. Repository → CockroachDB
5. Domain Entities → DTOs
6. DTOs → GraphQL Types
7. Response → Cliente
```

---

## 🧮 Algoritmos

### Algoritmo de Recomendación

**Complejidad**: O(n log n) donde n = número de stocks

**Pseudocódigo**:
```
1. Filtrar stocks con rating positivo (Buy, Strong Buy, Speculative Buy)
2. Para cada stock:
   a. Calcular cambio porcentual: (target_to - target_from) / target_from
   b. Calcular score de rating:
      - Buy → 3 puntos
      - Strong Buy → 5 puntos
      - Speculative Buy → 2 puntos
      - Rating upgrade → +2 puntos bonus
   c. Calcular score final: (cambio_porcentual * 0.6) + (rating_score * 0.4)
3. Ordenar por score descendente (O(n log n))
4. Retornar top 10
```

**Implementación Go**:
```go
func (s *RecommendationService) GetRecommendations(ctx context.Context, limit int) ([]*Recommendation, error) {
    stocks, err := s.repo.FindAll(ctx, Filter{Rating: []string{"Buy", "Strong Buy"}}, Sort{})
    if err != nil {
        return nil, err
    }
    
    recommendations := make([]*Recommendation, 0, len(stocks))
    for _, stock := range stocks {
        score := s.calculateScore(stock)
        recommendations = append(recommendations, &Recommendation{
            Stock: stock,
            Score: score,
        })
    }
    
    sort.Slice(recommendations, func(i, j int) bool {
        return recommendations[i].Score > recommendations[j].Score
    })
    
    if len(recommendations) > limit {
        recommendations = recommendations[:limit]
    }
    
    return recommendations, nil
}
```

### Algoritmo de Búsqueda

**Complejidad**: O(n) donde n = número de stocks en memoria (con índices DB: O(log n))

**Estrategia**:
- Búsqueda por índice en base de datos (ticker, company_name)
- Filtrado en memoria para múltiples criterios
- Cache de resultados para búsquedas frecuentes

---

## 🔌 Integración con API Externa

### Estrategia de Paginación

```go
func (c *KarenAIClient) FetchAllStocks(ctx context.Context) ([]*Stock, error) {
    var allStocks []*Stock
    nextPage := ""
    
    for {
        response, err := c.FetchStocks(ctx, nextPage)
        if err != nil {
            return nil, err
        }
        
        allStocks = append(allStocks, response.Stocks...)
        
        if response.NextPage == "" {
            break
        }
        nextPage = response.NextPage
    }
    
    return allStocks, nil
}
```

### Manejo de Errores

- **Retry Logic**: 3 intentos con backoff exponencial
- **Rate Limiting**: Respetar límites de la API
- **Timeout**: 30 segundos por request
- **Circuit Breaker**: Prevenir cascading failures

---

## 🗄️ Modelo de Datos

### Esquema de Base de Datos

```sql
CREATE TABLE stocks (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    ticker VARCHAR(10) NOT NULL UNIQUE,
    company_name VARCHAR(255) NOT NULL,
    brokerage VARCHAR(255),
    action VARCHAR(50),
    rating_from VARCHAR(50),
    rating_to VARCHAR(50),
    target_from DECIMAL(10,2),
    target_to DECIMAL(10,2),
    created_at TIMESTAMP DEFAULT now(),
    updated_at TIMESTAMP DEFAULT now()
);

-- Índices para optimización
CREATE INDEX idx_stocks_ticker ON stocks(ticker);
CREATE INDEX idx_stocks_rating_to ON stocks(rating_to);
CREATE INDEX idx_stocks_target_to ON stocks(target_to);
CREATE INDEX idx_stocks_company_name ON stocks(company_name);

-- Índice compuesto para búsquedas frecuentes
CREATE INDEX idx_stocks_rating_target ON stocks(rating_to, target_to);
```

---

## 🎨 Frontend Architecture

### Higher Order Components (HOCs) en Vue 3

**Importante**: Los HOCs (Higher Order Components) son funciones que toman un componente y retornan un nuevo componente con funcionalidad adicional. **NO son hooks/composables**.

En Vue 3, los HOCs se implementan de la siguiente manera:

```typescript
// Ejemplo: withLoading HOC
import { defineComponent, h, Component } from 'vue';

export function withLoading<T extends Component>(WrappedComponent: T) {
  return defineComponent({
    name: `withLoading(${WrappedComponent.name || 'Component'})`,
    props: {
      loading: {
        type: Boolean,
        default: false
      }
    },
    setup(props, { slots, attrs }) {
      return () => {
        if (props.loading) {
          return h('div', { class: 'loading-container' }, [
            h('div', { class: 'spinner' }, 'Cargando...')
          ]);
        }
        return h(WrappedComponent, attrs, slots);
      };
    }
  });
}

// Uso:
const StockListWithLoading = withLoading(StockList);
```

**Diferencia con Composables**:
- **HOCs**: Envuelven componentes, modifican su estructura/renderizado
- **Composables**: Proporcionan lógica reutilizable (similar a hooks de React), se usan dentro de `setup()`

**HOCs a implementar**:
- `withLoading`: Muestra spinner mientras carga
- `withError`: Maneja y muestra errores
- `withPagination`: Agrega lógica de paginación
- `withSearch`: Agrega funcionalidad de búsqueda con debounce

### Component Hierarchy

```
App
├── ThemeProvider
├── Router
│   ├── StockList (withLoading, withError, withPagination)
│   │   ├── SearchBar (withSearch HOC)
│   │   ├── StockTable (Design System)
│   │   └── Pagination (Design System)
│   ├── StockDetail
│   │   └── StockCard (Design System)
│   └── Recommendations
│       └── RecommendationList
└── ThemeToggle
```

### State Management (Pinia)

```typescript
// stores/stock.ts
export const useStockStore = defineStore('stock', {
  state: () => ({
    stocks: [] as Stock[],
    loading: false,
    error: null as string | null,
    cache: new Map<string, Stock[]>(),
  }),
  
  actions: {
    async fetchStocks(filter: Filter) {
      // Con cache y deduplicación
    }
  }
})
```

---

## 🚀 Optimizaciones

### Backend

1. **Database**:
   - Índices en columnas frecuentemente consultadas
   - Connection pooling
   - Query optimization

2. **API**:
   - DataLoader para evitar N+1 queries
   - Response caching
   - Pagination en GraphQL

3. **Sincronización**:
   - Batch inserts
   - Upsert en lugar de insert + update
   - Background jobs para sync periódico

### Frontend

1. **Rendering**:
   - Virtual scrolling para listas grandes
   - Lazy loading de componentes
   - Code splitting por ruta

2. **API Calls**:
   - Request deduplication
   - Cache en memoria con TTL
   - Debounce en búsquedas

3. **Performance**:
   - Memoization de componentes pesados
   - Web Workers para cálculos complejos
   - Image optimization

---

## 🔒 Seguridad

1. **API Key**: Almacenada en variables de entorno
2. **SQL Injection**: Usar prepared statements
3. **CORS**: Configurado apropiadamente
4. **Rate Limiting**: Implementar en API
5. **Input Validation**: Validar todos los inputs

---

## 📊 Monitoreo y Logging

1. **Logging**: Structured logging con niveles
2. **Metrics**: Tiempo de respuesta, errores, etc.
3. **Health Checks**: Endpoint `/health`
4. **Error Tracking**: Captura y reporte de errores

---

**Última actualización**: [Fecha]
