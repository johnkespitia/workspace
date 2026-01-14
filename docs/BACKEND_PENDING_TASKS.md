# Tareas Pendientes del Backend

## 📊 Estado Actual

### ✅ FASE 1: Backend - Infraestructura y Dominio (90% Completado)

#### ✅ Completado:

- [x] Crear esquema de base de datos para acciones
- [x] Implementar migraciones con estructura DDD
- [x] Configurar conexión a CockroachDB
- [x] Crear índices para optimización de consultas
- [x] Crear entidad `Stock` con validaciones
- [x] Definir interfaces de repositorio
- [x] Crear servicios de dominio
- [x] Implementar value objects (Rating, TargetPrice)
- [x] Implementar cliente HTTP para `api.karenai.click`
- [x] Manejo de paginación con `next_page`
- [x] Procesamiento en batch para optimizar escrituras

#### ⚠️ Pendiente:

- [x] **Manejo de errores y retry logic** en cliente API externa ✅
- [x] **Rate limiting** en cliente API externa ✅
- [x] **Caching** en cliente API externa ✅

---

### ✅ FASE 2: Backend - GraphQL API (95% Completado)

#### ✅ Completado:

- [x] Definir tipos: `Stock`, `StockConnection`, `Recommendation`
- [x] Queries: `stocks`, `stock`, `recommendations`
- [x] Mutations: `syncStocks`
- [x] Inputs: `StockFilter`, `StockSort`
- [x] Implementar resolvers con inyección de dependencias
- [x] Manejo de errores GraphQL
- [x] `StockService`: Lógica de negocio para acciones
- [x] `SyncService`: Sincronización con API externa
- [x] `RecommendationService`: Algoritmo de recomendación

#### ⚠️ Pendiente:

- [x] **Implementar DataLoader para evitar N+1 queries** ✅

---

### ❌ FASE 3: Backend - Documentación y Tests (0% Completado)

#### ✅ Completado:

- [x] **Tests para servicios de dominio** ✅
- [x] **Tests para repositorios** (con mocks) ✅
- [x] **Tests para resolvers GraphQL** ✅
- [x] **Tests para algoritmo de recomendación** ✅

#### ✅ Completado (Opcional):

- [x] **Configurar Swagger/OpenAPI** ✅ (openapi.yaml creado)
- [x] **Documentar endpoints HTTP** ✅ (API_DOCUMENTATION.md)
- [x] **Ejemplos de requests/responses** ✅ (API_DOCUMENTATION.md)
- [x] **Documentar schema GraphQL** ✅ (GRAPHQL_API_REFERENCE.md mejorado)
- [x] **Ejemplos de queries y mutations** ✅ (GRAPHQL_EXAMPLES.md y API_DOCUMENTATION.md)
- [x] **Guía de uso completa** ✅ (USER_GUIDE.md)

---

## 🎯 Tareas Prioritarias

### 1. **Implementar DataLoader para N+1 Queries** (Alta Prioridad)

**Problema**: Actualmente, si una query GraphQL solicita múltiples stocks con relaciones, puede generar múltiples queries a la base de datos (N+1 problem).

**Solución**: Implementar DataLoader para batch loading.

**Archivo a crear**: `api/internal/application/graphql/dataloader.go`

**Implementación sugerida**:

```go
package graphql

import (
    "context"
    "github.com/graph-gophers/dataloader/v7"
)

type StockLoader struct {
    loader *dataloader.Loader[string, *stock.Stock]
}

func NewStockLoader(stockService *services.StockService) *StockLoader {
    return &StockLoader{
        loader: dataloader.NewBatchedLoader(
            func(ctx context.Context, keys []string) []*dataloader.Result[*stock.Stock] {
                // Batch fetch stocks
                results := make([]*dataloader.Result[*stock.Stock], len(keys))
                stocks, err := stockService.GetStocksByTickers(ctx, keys)
                // ... implementar lógica
                return results
            },
        ),
    }
}
```

**Dependencia necesaria**: `github.com/graph-gophers/dataloader/v7`

---

### 2. **Implementar Retry Logic en Cliente API Externa** (Media Prioridad)

**Archivo a modificar**: `api/internal/infrastructure/external/karenai_api.go`

**Implementación sugerida**:

```go
func (c *KarenAIClient) fetchWithRetry(ctx context.Context, url string, maxRetries int) (*APIResponse, error) {
    var lastErr error
    for i := 0; i < maxRetries; i++ {
        resp, err := c.fetchPage(ctx, url)
        if err == nil {
            return resp, nil
        }
        lastErr = err

        // Exponential backoff
        backoff := time.Duration(i+1) * time.Second
        time.Sleep(backoff)
    }
    return nil, fmt.Errorf("failed after %d retries: %w", maxRetries, lastErr)
}
```

---

### 3. **Implementar Rate Limiting** (Media Prioridad)

**Archivo a modificar**: `api/internal/infrastructure/external/karenai_api.go`

**Implementación sugerida**:

```go
type RateLimiter struct {
    limiter *rate.Limiter
}

func NewRateLimiter(requestsPerSecond float64) *RateLimiter {
    return &RateLimiter{
        limiter: rate.NewLimiter(rate.Limit(requestsPerSecond), 1),
    }
}

func (c *KarenAIClient) fetchPage(ctx context.Context, url string) (*APIResponse, error) {
    // Wait for rate limiter
    if err := c.rateLimiter.Wait(ctx); err != nil {
        return nil, err
    }
    // ... resto de la lógica
}
```

**Dependencia necesaria**: `golang.org/x/time/rate`

---

### 4. **Implementar Caching** (Baja Prioridad)

**Archivo a modificar**: `api/internal/infrastructure/external/karenai_api.go`

**Implementación sugerida**:

```go
type Cache interface {
    Get(key string) (interface{}, bool)
    Set(key string, value interface{}, ttl time.Duration)
}

func (c *KarenAIClient) FetchAllStocks(ctx context.Context) ([]*stock.Stock, error) {
    cacheKey := "stocks:all"
    if cached, ok := c.cache.Get(cacheKey); ok {
        return cached.([]*stock.Stock), nil
    }

    stocks, err := c.fetchAllStocks(ctx)
    if err == nil {
        c.cache.Set(cacheKey, stocks, 5*time.Minute)
    }
    return stocks, err
}
```

---

### 5. **Tests Unitarios** (Alta Prioridad)

#### 5.1 Tests para Servicios de Dominio

**Archivo a crear**: `api/internal/domain/stock/service_test.go`

```go
func TestCalculatePriceChange(t *testing.T) {
    // Test cases
}

func TestIsRatingUpgrade(t *testing.T) {
    // Test cases
}

func TestCalculateRecommendationScore(t *testing.T) {
    // Test cases
}
```

#### 5.2 Tests para Repositorios (con mocks)

**Archivo a crear**: `api/internal/infrastructure/repository/stock_repository_test.go`

**Dependencia necesaria**: `github.com/DATA-DOG/go-sqlmock` o similar

```go
func TestStockRepository_Save(t *testing.T) {
    // Mock database
    // Test Save operation
}

func TestStockRepository_FindAll(t *testing.T) {
    // Mock database
    // Test FindAll with filters
}
```

#### 5.3 Tests para Resolvers GraphQL

**Archivo a crear**: `api/internal/application/graphql/resolvers_test.go`

```go
func TestResolver_Stocks(t *testing.T) {
    // Mock services
    // Test Stocks resolver
}

func TestResolver_SyncStocks(t *testing.T) {
    // Mock services
    // Test SyncStocks mutation
}
```

#### 5.4 Tests para Algoritmo de Recomendación

**Archivo a crear**: `api/internal/domain/recommendation/algorithm_test.go`

```go
func TestRecommendationAlgorithm_CalculateScore(t *testing.T) {
    // Test cases para diferentes escenarios
}
```

---

### 6. **Documentación Swagger/OpenAPI** (Baja Prioridad)

**Nota**: Como el proyecto usa GraphQL principalmente, Swagger puede no ser necesario. Sin embargo, si se quiere documentar endpoints HTTP adicionales:

**Archivo a crear**: `api/docs/swagger.yaml`

**Dependencia necesaria**: `github.com/swaggo/swag`

---

## 📋 Checklist de Implementación

### Prioridad Alta

- [x] Implementar DataLoader para N+1 queries ✅
- [x] Tests unitarios para servicios de dominio ✅
- [x] Tests unitarios para repositorios ✅
- [x] Tests unitarios para resolvers GraphQL ✅
- [x] Tests unitarios para algoritmo de recomendación ✅

### Prioridad Media

- [x] Implementar retry logic en cliente API externa ✅
- [x] Implementar rate limiting en cliente API externa ✅

### Prioridad Baja

- [x] Implementar caching en cliente API externa ✅
- [ ] Configurar Swagger/OpenAPI (si es necesario)
- [ ] Documentación adicional

---

## 🔧 Dependencias Adicionales Necesarias

```bash
# Para DataLoader
go get github.com/graph-gophers/dataloader/v7

# Para Rate Limiting
go get golang.org/x/time/rate

# Para Testing
go get github.com/stretchr/testify
go get github.com/DATA-DOG/go-sqlmock

# Para Swagger (opcional)
go get github.com/swaggo/swag/cmd/swag
```

---

## 📝 Notas

1. **DataLoader**: Es crítico para evitar problemas de rendimiento en producción cuando hay múltiples queries GraphQL concurrentes.

2. **Tests**: Son esenciales para mantener la calidad del código y prevenir regresiones.

3. **Retry Logic y Rate Limiting**: Importantes para robustez en producción, especialmente cuando se interactúa con APIs externas.

4. **Caching**: Puede mejorar significativamente el rendimiento, pero debe implementarse con cuidado para evitar datos obsoletos.

---

## 🎯 Estimación de Tiempo

- **DataLoader**: 2-3 horas
- **Retry Logic**: 1-2 horas
- **Rate Limiting**: 1-2 horas
- **Caching**: 2-3 horas
- **Tests Unitarios**: 4-6 horas
- **Documentación Swagger**: 1-2 horas (opcional)

**Total estimado**: 11-18 horas (1.5-2.5 días)
