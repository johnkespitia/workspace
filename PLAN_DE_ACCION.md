# Plan de Acción - Sistema de Información de Acciones

## 📋 Resumen del Proyecto

Sistema completo para recuperar, almacenar y visualizar información de acciones desde una API externa, con recomendaciones inteligentes de inversión.

## 🎯 Objetivos

1. **Conectar a API externa** y almacenar datos en CockroachDB
2. **Crear API GraphQL** y UI intuitiva con búsqueda, ordenamiento y visualización
3. **Algoritmo de recomendación** para mejores acciones de inversión
4. **Tests unitarios** (opcional pero recomendado)

---

## 🏗️ Arquitectura del Sistema

### Backend (Go + GraphQL + DDD)

```
api/
├── cmd/
│   └── main.go                    # Punto de entrada
├── internal/
│   ├── domain/                    # Capa de Dominio (DDD)
│   │   ├── stock/
│   │   │   ├── entity.go          # Entidad Stock
│   │   │   ├── repository.go      # Interfaz del repositorio
│   │   │   └── service.go         # Lógica de negocio
│   │   └── recommendation/
│   │       ├── entity.go          # Entidad Recommendation
│   │       └── algorithm.go       # Algoritmo de recomendación
│   ├── infrastructure/            # Capa de Infraestructura
│   │   ├── database/
│   │   │   ├── cockroach.go       # Conexión a CockroachDB
│   │   │   └── migrations/        # Migraciones de BD
│   │   ├── external/
│   │   │   └── karenai_api.go     # Cliente HTTP para API externa
│   │   └── repository/
│   │       └── stock_repository.go # Implementación del repositorio
│   ├── application/               # Capa de Aplicación
│   │   ├── handlers/
│   │   │   ├── graphql_handler.go # Handler GraphQL
│   │   │   └── http_handler.go    # Handler HTTP (Swagger)
│   │   ├── services/
│   │   │   ├── stock_service.go   # Servicio de acciones
│   │   │   └── sync_service.go    # Servicio de sincronización
│   │   └── graphql/
│   │       ├── schema.graphql      # Schema GraphQL
│   │       ├── resolvers.go        # Resolvers GraphQL
│   │       └── types.go            # Tipos GraphQL
│   └── config/
│       └── config.go               # Configuración
├── docs/
│   ├── api/
│   │   ├── swagger.yaml            # Documentación Swagger
│   │   └── graphql.md              # Documentación GraphQL
│   └── architecture.md             # Documentación de arquitectura
└── go.mod
```

### Frontend (Vue 3 + TypeScript + Design System)

```
frontend/
├── src/
│   ├── main.ts                     # Punto de entrada
│   ├── App.vue                     # Componente raíz
│   ├── router/                      # Vue Router
│   │   └── index.ts
│   ├── stores/                      # Pinia stores
│   │   ├── stock.ts                 # Store de acciones
│   │   └── theme.ts                 # Store de temas
│   ├── composables/                 # Composables Vue
│   │   ├── useStock.ts              # Lógica de acciones
│   │   └── useApi.ts                # Cliente API
│   ├── hoc/                         # Higher Order Components
│   │   ├── withLoading.ts           # HOC para loading
│   │   ├── withError.ts             # HOC para errores
│   │   └── withPagination.ts        # HOC para paginación
│   ├── design-system/               # Design System
│   │   ├── components/              # Componentes reusables
│   │   │   ├── Button/
│   │   │   │   ├── Button.vue
│   │   │   │   └── Button.stories.ts
│   │   │   ├── Table/
│   │   │   │   ├── Table.vue
│   │   │   │   └── Table.stories.ts
│   │   │   ├── Input/
│   │   │   │   ├── Input.vue
│   │   │   │   └── Input.stories.ts
│   │   │   ├── Card/
│   │   │   │   ├── Card.vue
│   │   │   │   └── Card.stories.ts
│   │   │   └── ThemeToggle/
│   │   │       ├── ThemeToggle.vue
│   │   │       └── ThemeToggle.stories.ts
│   │   ├── tokens/                  # Design tokens
│   │   │   ├── colors.ts
│   │   │   ├── spacing.ts
│   │   │   └── typography.ts
│   │   └── themes/                  # Temas
│   │       ├── light.ts
│   │       └── dark.ts
│   ├── views/                       # Vistas/páginas
│   │   ├── StockList.vue            # Lista de acciones
│   │   ├── StockDetail.vue          # Detalle de acción
│   │   └── Recommendations.vue     # Recomendaciones
│   └── utils/                       # Utilidades
│       ├── api.ts                   # Cliente GraphQL
│       ├── debounce.ts              # Debounce para búsqueda
│       └── accessibility.ts        # Utilidades de accesibilidad
├── .storybook/                      # Configuración Storybook
│   ├── main.ts
│   └── preview.ts
└── package.json
```

---

## 📝 Fases de Implementación

### **FASE 1: Backend - Infraestructura y Dominio** ⏱️ Estimado: 2-3 días

#### 1.1 Configuración de Base de Datos
- [ ] Crear esquema de base de datos para acciones
- [ ] Implementar migraciones con estructura DDD
- [ ] Configurar conexión a CockroachDB
- [ ] Crear índices para optimización de consultas

**Estructura de Tabla:**
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
    updated_at TIMESTAMP DEFAULT now(),
    INDEX idx_ticker (ticker),
    INDEX idx_rating_to (rating_to),
    INDEX idx_target_to (target_to)
);
```

#### 1.2 Capa de Dominio (DDD)
- [ ] Crear entidad `Stock` con validaciones
- [ ] Definir interfaces de repositorio
- [ ] Crear servicios de dominio
- [ ] Implementar value objects (Rating, TargetPrice)

#### 1.3 Cliente API Externa
- [ ] Implementar cliente HTTP para `api.karenai.click`
- [ ] Manejo de paginación con `next_page`
- [ ] Manejo de errores y retry logic
- [ ] Rate limiting y caching

**Algoritmo de Sincronización:**
- **Complejidad**: O(n) donde n = número de registros por página
- **Estrategia**: 
  - Fetch paginado de la API externa
  - Upsert en base de datos (evitar duplicados)
  - Procesamiento en batch para optimizar escrituras

---

### **FASE 2: Backend - GraphQL API** ⏱️ Estimado: 2-3 días

#### 2.1 Schema GraphQL
- [ ] Definir tipos: `Stock`, `StockConnection`, `Recommendation`
- [ ] Queries: `stocks`, `stock`, `recommendations`
- [ ] Mutations: `syncStocks` (para sincronizar desde API externa)
- [ ] Inputs: `StockFilter`, `StockSort`

#### 2.2 Resolvers
- [ ] Implementar resolvers con inyección de dependencias
- [ ] Implementar DataLoader para evitar N+1 queries
- [ ] Manejo de errores GraphQL

#### 2.3 Servicios de Aplicación
- [ ] `StockService`: Lógica de negocio para acciones
- [ ] `SyncService`: Sincronización con API externa
- [ ] `RecommendationService`: Algoritmo de recomendación

**Algoritmo de Recomendación:**
- **Complejidad**: O(n log n) para ordenamiento
- **Estrategia**:
  1. Filtrar acciones con rating positivo (Buy, Strong Buy)
  2. Calcular score basado en:
     - Cambio porcentual en target: `(target_to - target_from) / target_from`
     - Rating upgrade/downgrade
     - Consistencia del rating
  3. Ordenar por score descendente
  4. Retornar top N recomendaciones

---

### **FASE 3: Backend - Documentación y Tests** ⏱️ Estimado: 1-2 días

#### 3.1 Documentación Swagger
- [ ] Configurar Swagger/OpenAPI
- [ ] Documentar endpoints HTTP (si los hay)
- [ ] Ejemplos de requests/responses

#### 3.2 Documentación GraphQL
- [ ] Documentar schema GraphQL
- [ ] Ejemplos de queries y mutations
- [ ] Guía de uso

#### 3.3 Tests Unitarios
- [ ] Tests para servicios de dominio
- [ ] Tests para repositorios (con mocks)
- [ ] Tests para resolvers GraphQL
- [ ] Tests para algoritmo de recomendación

---

### **FASE 4: Frontend - Design System** ⏱️ Estimado: 3-4 días

#### 4.1 Configuración Storybook
- [ ] Instalar y configurar Storybook
- [ ] Configurar temas (light/dark)
- [ ] Configurar accesibilidad addon

#### 4.2 Componentes Base
- [ ] **Button**: Variantes, estados, accesibilidad
- [ ] **Input**: Búsqueda, validación, accesibilidad
- [ ] **Table**: Ordenamiento, paginación, accesibilidad
- [ ] **Card**: Variantes, estados
- [ ] **ThemeToggle**: Cambio de tema

#### 4.3 Design Tokens
- [ ] Colores (light/dark themes)
- [ ] Espaciado
- [ ] Tipografía
- [ ] Breakpoints

#### 4.4 Accesibilidad
- [ ] ARIA labels en todos los componentes
- [ ] Navegación por teclado
- [ ] Contraste de colores (WCAG AA)
- [ ] Screen reader support

---

### **FASE 5: Frontend - HOCs y Lógica** ⏱️ Estimado: 2-3 días

#### 5.1 Higher Order Components (HOCs)
**Nota**: HOCs son Higher Order Components (funciones que toman un componente y retornan uno nuevo), NO hooks/composables.

- [ ] `withLoading`: HOC que envuelve componentes y muestra spinner durante carga
- [ ] `withError`: HOC que maneja y muestra errores en componentes
- [ ] `withPagination`: HOC que agrega lógica de paginación a componentes de lista
- [ ] `withSearch`: HOC que agrega funcionalidad de búsqueda con debounce

**Implementación en Vue 3**:
```typescript
// Patrón: Función que retorna un componente usando defineComponent y h()
export function withLoading<T extends Component>(WrappedComponent: T) {
  return defineComponent({
    setup(props, { slots, attrs }) {
      return () => props.loading 
        ? h('div', 'Loading...') 
        : h(WrappedComponent, attrs, slots);
    }
  });
}
```

**Optimización de Búsqueda:**
- **Algoritmo**: Debounce con complejidad O(1) por llamada
- **Implementación**: 
  - Debounce de 300ms para búsqueda
  - Cache de resultados en memoria
  - Cancelación de requests anteriores

#### 5.2 Composables
- [ ] `useStock`: Lógica de acciones con cache
- [ ] `useApi`: Cliente GraphQL con cache
- [ ] `useRecommendations`: Lógica de recomendaciones

**Optimización de API Calls:**
- **Estrategia**:
  - Cache en memoria con TTL
  - Request deduplication
  - Lazy loading de datos
  - Paginación en frontend

---

### **FASE 6: Frontend - Vistas y Integración** ⏱️ Estimado: 2-3 días

#### 6.1 Vista de Lista de Acciones
- [ ] Tabla con todas las acciones
- [ ] Búsqueda por ticker/company
- [ ] Ordenamiento por columnas
- [ ] Paginación
- [ ] Filtros (rating, action)

#### 6.2 Vista de Detalle
- [ ] Información completa de la acción
- [ ] Historial de cambios
- [ ] Gráficos (opcional)

#### 6.3 Vista de Recomendaciones
- [ ] Lista de mejores acciones
- [ ] Score de recomendación
- [ ] Explicación del algoritmo

#### 6.4 Integración
- [ ] Conectar con GraphQL API
- [ ] Manejo de estados globales (Pinia)
- [ ] Manejo de errores y loading
- [ ] Optimización de rendimiento

---

### **FASE 7: Optimización y Pulido** ⏱️ Estimado: 1-2 días

#### 7.1 Optimización de Rendimiento
- [ ] Code splitting
- [ ] Lazy loading de rutas
- [ ] Optimización de imágenes
- [ ] Bundle size optimization

#### 7.2 Testing Frontend
- [ ] Tests unitarios de componentes
- [ ] Tests de HOCs
- [ ] Tests de stores (Pinia)

#### 7.3 Documentación Final
- [ ] README actualizado
- [ ] Guía de desarrollo
- [ ] Guía de deployment

---

## 🔧 Stack Tecnológico Detallado

### Backend
- **Go 1.21+**
- **GraphQL**: `github.com/graphql-go/graphql` o `github.com/99designs/gqlgen`
- **Database**: CockroachDB (driver: `github.com/lib/pq`)
- **HTTP Client**: `net/http` estándar o `github.com/go-resty/resty`
- **Swagger**: `github.com/swaggo/swag`
- **Testing**: `testing` package + `github.com/stretchr/testify`

### Frontend
- **Vue 3** (Composition API)
- **TypeScript**
- **Pinia** (State management)
- **Vue Router** (Routing)
- **Apollo Client** o **urql** (GraphQL client)
- **Tailwind CSS** (Styling)
- **Storybook** (Design System docs)
- **Vitest** (Testing)

---

## 📊 Métricas de Éxito

1. **Performance**:
   - Tiempo de carga inicial < 2s
   - Búsqueda con respuesta < 300ms
   - API response time < 500ms

2. **Accesibilidad**:
   - Score Lighthouse > 90
   - WCAG AA compliance
   - Keyboard navigation completa

3. **Código**:
   - Cobertura de tests > 70%
   - Documentación completa
   - Código desacoplado y mantenible

---

## 🚀 Próximos Pasos Inmediatos

1. **Crear estructura de carpetas** según arquitectura DDD
2. **Configurar base de datos** y migraciones
3. **Implementar cliente API externa**
4. **Crear entidades de dominio**
5. **Configurar GraphQL**

---

## 📚 Referencias y Notas

- **API Externa**: `https://api.karenai.click/swechallenge/list`
- **Auth Token**: Incluido en headers
- **Paginación**: Usar `next_page` query parameter
- **Sample Data**: Ver imagen proporcionada (TICKER, COMPANY, RATING, TARGET, etc.)

---

**Última actualización**: [Fecha]
**Estado**: 🟡 En Planificación
