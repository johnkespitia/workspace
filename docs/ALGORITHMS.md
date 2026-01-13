# Algoritmos y Optimizaciones - Sistema de Información de Acciones

## 📊 Algoritmo de Recomendación de Stocks

### Objetivo
Identificar las mejores acciones para invertir basándose en:
- Cambio en precio objetivo (target)
- Rating de la acción
- Consistencia del rating

### Complejidad Temporal
- **O(n log n)** donde n = número de stocks
  - Filtrado: O(n)
  - Cálculo de scores: O(n)
  - Ordenamiento: O(n log n)
  - Selección top N: O(1)

### Complejidad Espacial
- **O(n)** para almacenar stocks y scores

### Pseudocódigo

```
ALGORITMO: GetRecommendations(stocks, limit)
ENTRADA: stocks[] (lista de acciones), limit (número de recomendaciones)
SALIDA: recommendations[] (top acciones recomendadas)

1. FILTRAR stocks con rating positivo:
   - Buy
   - Strong Buy
   - Speculative Buy
   - Market Perform (si target aumentó)

2. PARA CADA stock en stocks_filtrados:
   a. Calcular cambio porcentual:
      priceChange = (target_to - target_from) / target_from * 100
   
   b. Calcular score de rating:
      SI rating_to == "Strong Buy":
         ratingScore = 5
      SINO SI rating_to == "Buy":
         ratingScore = 3
      SINO SI rating_to == "Speculative Buy":
         ratingScore = 2
      SINO:
         ratingScore = 1
      
      SI rating cambió de menor a mayor:
         ratingScore = ratingScore + 2 (bonus)
   
   c. Calcular score de acción:
      SI action == "target raised by":
         actionScore = 3
      SINO SI action == "target lowered by":
         actionScore = -2
      SINO:
         actionScore = 0
   
   d. Calcular score final:
      finalScore = (priceChange * 0.5) + (ratingScore * 0.3) + (actionScore * 0.2)

3. ORDENAR stocks por finalScore DESCENDENTE

4. RETORNAR primeros 'limit' stocks
```

### Implementación Go

```go
type RecommendationScore struct {
    Stock       *Stock
    Score       float64
    PriceChange float64
    RatingScore float64
    ActionScore float64
}

func (s *RecommendationService) CalculateRecommendations(
    ctx context.Context,
    stocks []*domain.Stock,
    limit int,
) ([]*RecommendationScore, error) {
    // Paso 1: Filtrar stocks con rating positivo
    filteredStocks := s.filterPositiveRatings(stocks)
    
    // Paso 2: Calcular scores
    recommendations := make([]*RecommendationScore, 0, len(filteredStocks))
    for _, stock := range filteredStocks {
        score := s.calculateScore(stock)
        recommendations = append(recommendations, score)
    }
    
    // Paso 3: Ordenar por score descendente
    sort.Slice(recommendations, func(i, j int) bool {
        return recommendations[i].Score > recommendations[j].Score
    })
    
    // Paso 4: Retornar top N
    if len(recommendations) > limit {
        recommendations = recommendations[:limit]
    }
    
    return recommendations, nil
}

func (s *RecommendationService) calculateScore(stock *domain.Stock) *RecommendationScore {
    // Calcular cambio porcentual
    priceChange := 0.0
    if stock.TargetFrom.Value() > 0 {
        priceChange = ((stock.TargetTo.Value() - stock.TargetFrom.Value()) / 
                      stock.TargetFrom.Value()) * 100
    }
    
    // Calcular score de rating
    ratingScore := s.getRatingScore(stock.RatingTo, stock.RatingFrom)
    
    // Calcular score de acción
    actionScore := s.getActionScore(stock.Action)
    
    // Score final (pesos: 50% price change, 30% rating, 20% action)
    finalScore := (priceChange * 0.5) + (ratingScore * 0.3) + (actionScore * 0.2)
    
    return &RecommendationScore{
        Stock:       stock,
        Score:       finalScore,
        PriceChange: priceChange,
        RatingScore: ratingScore,
        ActionScore: actionScore,
    }
}
```

### Ejemplo de Cálculo

**Stock A:**
- Target: $50 → $60 (aumento 20%)
- Rating: Buy → Buy (consistente)
- Action: "target raised by"
- Score: (20 * 0.5) + (3 * 0.3) + (3 * 0.2) = **13.5**

**Stock B:**
- Target: $100 → $110 (aumento 10%)
- Rating: Neutral → Buy (upgrade)
- Action: "target raised by"
- Score: (10 * 0.5) + (5 * 0.3) + (3 * 0.2) = **8.6**

**Resultado**: Stock A tiene mayor score y se recomienda primero.

---

## 🔍 Algoritmo de Búsqueda

### Objetivo
Buscar stocks por ticker, nombre de compañía, rating, etc.

### Complejidad Temporal
- **Con índices DB**: O(log n) para búsqueda por índice
- **Sin índices**: O(n) para búsqueda lineal
- **Búsqueda múltiple**: O(n) donde n = resultados filtrados

### Estrategia de Implementación

```go
// Búsqueda optimizada con índices
func (r *CockroachStockRepository) Search(
    ctx context.Context,
    query string,
) ([]*domain.Stock, error) {
    // Si query es ticker (formato corto), usar índice de ticker
    if len(query) <= 10 && strings.ToUpper(query) == query {
        return r.findByTicker(ctx, query)
    }
    
    // Si no, búsqueda por nombre de compañía (con índice)
    return r.findByCompanyName(ctx, query)
}
```

### Búsqueda con Filtros Múltiples

```go
// Complejidad: O(n) donde n = número de stocks en resultado
func (r *CockroachStockRepository) FindWithFilters(
    ctx context.Context,
    filters Filter,
) ([]*domain.Stock, error) {
    // Construir query SQL con WHERE clauses
    // Usar índices compuestos cuando sea posible
    query := "SELECT * FROM stocks WHERE 1=1"
    args := []interface{}{}
    
    if filters.Ticker != "" {
        query += " AND ticker = $1"
        args = append(args, filters.Ticker)
    }
    
    if len(filters.Ratings) > 0 {
        query += " AND rating_to = ANY($2)"
        args = append(args, filters.Ratings)
    }
    
    // Ejecutar query con índices
    // ...
}
```

---

## 🔄 Algoritmo de Sincronización

### Objetivo
Sincronizar datos desde API externa a base de datos local.

### Complejidad Temporal
- **O(n)** donde n = número total de registros
- **Por página**: O(m) donde m = registros por página

### Estrategia de Paginación

```go
func (s *SyncService) SyncAllStocks(ctx context.Context) error {
    var allStocks []*domain.Stock
    nextPage := ""
    pageCount := 0
    
    // Fetch todas las páginas
    for {
        response, err := s.apiClient.FetchStocks(ctx, nextPage)
        if err != nil {
            return fmt.Errorf("error fetching page %d: %w", pageCount, err)
        }
        
        // Convertir DTOs a entidades de dominio
        stocks := s.convertToDomainEntities(response.Stocks)
        allStocks = append(allStocks, stocks...)
        
        // Verificar si hay más páginas
        if response.NextPage == "" {
            break
        }
        nextPage = response.NextPage
        pageCount++
    }
    
    // Batch insert/update en base de datos
    return s.batchUpsert(ctx, allStocks)
}
```

### Batch Upsert (Optimización)

```go
// Complejidad: O(n) con batch size b
// Mejor que O(n) inserts individuales
func (r *CockroachStockRepository) BatchUpsert(
    ctx context.Context,
    stocks []*domain.Stock,
) error {
    batchSize := 100
    for i := 0; i < len(stocks); i += batchSize {
        end := i + batchSize
        if end > len(stocks) {
            end = len(stocks)
        }
        
        batch := stocks[i:end]
        if err := r.upsertBatch(ctx, batch); err != nil {
            return err
        }
    }
    return nil
}

func (r *CockroachStockRepository) upsertBatch(
    ctx context.Context,
    stocks []*domain.Stock,
) error {
    // Usar INSERT ... ON CONFLICT (UPSERT) de PostgreSQL/CockroachDB
    query := `
        INSERT INTO stocks (ticker, company_name, ...)
        VALUES (unnest($1::text[]), unnest($2::text[]), ...)
        ON CONFLICT (ticker) 
        DO UPDATE SET 
            company_name = EXCLUDED.company_name,
            ...
    `
    // Ejecutar batch insert
    // ...
}
```

---

## ⚡ Optimizaciones Frontend

### Debounce para Búsqueda

**Complejidad**: O(1) por llamada (cancelación de timeout anterior)

```typescript
function useDebounce<T>(value: T, delay: number): T {
  const [debouncedValue, setDebouncedValue] = useState<T>(value);

  useEffect(() => {
    const handler = setTimeout(() => {
      setDebouncedValue(value);
    }, delay);

    return () => {
      clearTimeout(handler); // Cancelar timeout anterior
    };
  }, [value, delay]);

  return debouncedValue;
}
```

### Request Deduplication

**Complejidad**: O(1) con Map lookup

```typescript
class ApiClient {
  private pendingRequests = new Map<string, Promise<any>>();

  async fetchStocks(filter: Filter): Promise<Stock[]> {
    const key = JSON.stringify(filter);
    
    // Si ya hay una request pendiente con el mismo filtro, reutilizarla
    if (this.pendingRequests.has(key)) {
      return this.pendingRequests.get(key)!;
    }
    
    // Crear nueva request
    const promise = this.executeRequest(filter);
    this.pendingRequests.set(key, promise);
    
    // Limpiar después de completar
    promise.finally(() => {
      this.pendingRequests.delete(key);
    });
    
    return promise;
  }
}
```

### Cache con TTL

**Complejidad**: O(1) para get/set

```typescript
interface CacheEntry<T> {
  data: T;
  timestamp: number;
  ttl: number; // Time to live en ms
}

class Cache<T> {
  private cache = new Map<string, CacheEntry<T>>();

  get(key: string): T | null {
    const entry = this.cache.get(key);
    if (!entry) return null;
    
    // Verificar si expiró
    if (Date.now() - entry.timestamp > entry.ttl) {
      this.cache.delete(key);
      return null;
    }
    
    return entry.data;
  }

  set(key: string, data: T, ttl: number = 60000): void {
    this.cache.set(key, {
      data,
      timestamp: Date.now(),
      ttl,
    });
  }
}
```

### Virtual Scrolling (Para listas grandes)

**Complejidad**: O(visible_items) en lugar de O(all_items)

```typescript
// Solo renderizar items visibles en viewport
function VirtualList<T>({ items, itemHeight }: Props<T>) {
  const [scrollTop, setScrollTop] = useState(0);
  const containerHeight = 600; // altura del contenedor
  
  const startIndex = Math.floor(scrollTop / itemHeight);
  const endIndex = Math.min(
    startIndex + Math.ceil(containerHeight / itemHeight),
    items.length
  );
  
  const visibleItems = items.slice(startIndex, endIndex);
  
  return (
    <div onScroll={(e) => setScrollTop(e.target.scrollTop)}>
      <div style={{ height: items.length * itemHeight }}>
        {visibleItems.map((item, index) => (
          <Item key={startIndex + index} item={item} />
        ))}
      </div>
    </div>
  );
}
```

---

## 📈 Métricas de Performance Esperadas

### Backend
- **Query de stocks**: < 100ms (con índices)
- **Sincronización completa**: < 5s (depende de número de páginas)
- **Cálculo de recomendaciones**: < 50ms (para 1000 stocks)

### Frontend
- **Tiempo de carga inicial**: < 2s
- **Búsqueda con debounce**: < 300ms
- **Renderizado de tabla**: < 100ms (con virtual scrolling)
- **Cambio de tema**: < 50ms

---

## 🔧 Mejoras Futuras

1. **Caching distribuido**: Redis para cache compartido
2. **Background jobs**: Sincronización periódica automática
3. **WebSockets**: Actualizaciones en tiempo real
4. **Machine Learning**: Mejorar algoritmo de recomendación
5. **CDN**: Servir assets estáticos desde CDN

---

**Última actualización**: [Fecha]
