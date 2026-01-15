# Manual del Desarrollador

## 📋 Tabla de Contenidos

1. [Configuración del Entorno](#configuración-del-entorno)
2. [Estructura del Proyecto](#estructura-del-proyecto)
3. [Convenciones de Código](#convenciones-de-código)
4. [Desarrollo Backend](#desarrollo-backend)
5. [Desarrollo Frontend](#desarrollo-frontend)
6. [Testing](#testing)
7. [Debugging](#debugging)
8. [Contribución](#contribución)

---

## 🚀 Configuración del Entorno

### Prerrequisitos

- Docker Desktop instalado y corriendo
- Git instalado
- IDE con soporte para Dev Containers (VS Code o Cursor)

### Inicialización del Proyecto

```bash
# Clonar el repositorio
git clone <repository-url>
cd go-react-test

# Inicializar el dev container
make dev-init
```

El comando `make dev-init`:
- Construye las imágenes Docker necesarias
- Inicia todos los servicios (API, Frontend, CockroachDB)
- Abre automáticamente el IDE con el devcontainer

### Verificar Instalación

```bash
# Verificar estado de servicios
make dev-status

# Verificar salud de servicios
make dev-health
```

Todos los servicios deberían estar corriendo y accesibles:
- Frontend: http://localhost:3001
- Storybook: http://localhost:6006
- Backend: http://localhost:8080
- CockroachDB UI: http://localhost:8081

---

## 📁 Estructura del Proyecto

```
workspace/
├── api/                          # Backend Go
│   ├── cmd/                      # Puntos de entrada
│   │   ├── main.go              # Servidor principal
│   │   └── migrate/              # Herramienta de migraciones
│   ├── internal/
│   │   ├── domain/              # Capa de Dominio (DDD)
│   │   │   ├── stock/           # Entidades y servicios de dominio
│   │   │   └── recommendation/  # Algoritmo de recomendación
│   │   ├── application/         # Capa de Aplicación
│   │   │   ├── graphql/         # Schema y resolvers GraphQL
│   │   │   ├── handlers/        # HTTP handlers
│   │   │   └── services/        # Servicios de aplicación
│   │   └── infrastructure/      # Capa de Infraestructura
│   │       ├── database/         # Conexión y migraciones
│   │       ├── external/        # Clientes API externas
│   │       └── repository/      # Implementación de repositorios
│   └── docs/                     # Documentación API
├── frontend/                     # Frontend Vue 3
│   ├── src/
│   │   ├── design-system/        # Componentes reusables
│   │   ├── hoc/                 # Higher Order Components
│   │   ├── views/               # Vistas/páginas
│   │   ├── stores/              # Stores de Pinia
│   │   ├── composables/         # Composables Vue
│   │   └── utils/               # Utilidades
│   └── .storybook/              # Configuración Storybook
├── .devcontainer/                # Configuración Dev Container
│   ├── devcontainer.json
│   ├── docker-compose.yml
│   └── Dockerfile.*
├── docs/                         # Documentación general
└── Makefile                      # Comandos de desarrollo
```

---

## 📝 Convenciones de Código

### Backend (Go)

#### Estructura de Archivos

- **Entidades**: `entity.go`
- **Repositorios**: `repository.go` (interfaz), `*_repository.go` (implementación)
- **Servicios**: `service.go`
- **Tests**: `*_test.go`

#### Nomenclatura

- **Público**: PascalCase (`StockService`, `GetStocks`)
- **Privado**: camelCase (`calculateScore`, `fetchStocks`)
- **Constantes**: PascalCase o UPPER_CASE
- **Interfaces**: Nombre del comportamiento (`Repository`, `Service`)

#### Ejemplo

```go
// domain/stock/entity.go
type Stock struct {
    ID          uuid.UUID
    Ticker      string
    CompanyName string
    // ...
}

// domain/stock/repository.go
type Repository interface {
    Save(ctx context.Context, stock *Stock) error
    FindByTicker(ctx context.Context, ticker string) (*Stock, error)
}

// infrastructure/repository/stock_repository.go
type CockroachStockRepository struct {
    db *sql.DB
}

func (r *CockroachStockRepository) Save(ctx context.Context, stock *domain.Stock) error {
    // Implementación
}
```

### Frontend (Vue 3 + TypeScript)

#### Estructura de Componentes

```vue
<script setup lang="ts">
// 1. Imports
import { ref, computed } from 'vue';
import type { Stock } from '@/utils/api';

// 2. Props
interface Props {
  stocks: Stock[];
}
const props = defineProps<Props>();

// 3. Emits
const emit = defineEmits<{
  select: [stock: Stock];
}>();

// 4. Composables
const { data, loading } = useApi();

// 5. Estado local
const selectedStock = ref<Stock | null>(null);

// 6. Computed
const filteredStocks = computed(() => {
  return props.stocks.filter(/* ... */);
});

// 7. Métodos
function handleSelect(stock: Stock) {
  selectedStock.value = stock;
  emit('select', stock);
}
</script>

<template>
  <!-- Template -->
</template>

<style scoped>
/* Estilos */
</style>
```

#### Nomenclatura

- **Componentes**: PascalCase (`StockList.vue`, `StockCard.vue`)
- **Composables**: camelCase con prefijo `use` (`useStock.ts`, `useApi.ts`)
- **Stores**: camelCase (`useStockStore.ts`)
- **Utilidades**: camelCase (`api.ts`, `accessibility.ts`)

---

## 🔧 Desarrollo Backend

### Iniciar el Servidor

El servidor se inicia automáticamente al abrir el devcontainer. Si necesitas reiniciarlo:

```bash
# Reiniciar API
make dev-restart-api

# Ver logs
make dev-logs-api
```

### Hot Reload

El servidor usa `air` para hot reload automático. Los cambios en archivos `.go` se reflejan automáticamente.

### Migraciones de Base de Datos

```bash
# Ejecutar migraciones
cd api
go run cmd/migrate/main.go

# Ver migraciones existentes
ls api/internal/infrastructure/database/migrations/
```

### Estructura de Capas (DDD)

1. **Domain Layer**: Lógica de negocio pura, sin dependencias externas
2. **Application Layer**: Orquesta casos de uso, coordina entre capas
3. **Infrastructure Layer**: Implementaciones concretas (BD, APIs externas)
4. **Presentation Layer**: Handlers HTTP/GraphQL

### Agregar un Nuevo Endpoint GraphQL

1. **Definir en el schema** (`api/internal/application/graphql/schema.graphql`):

```graphql
type Query {
  myNewQuery(filter: MyFilter): MyResponse
}
```

2. **Regenerar código** (si usas gqlgen):

```bash
cd api
go run github.com/99designs/gqlgen generate
```

3. **Implementar resolver** (`api/internal/application/graphql/resolvers.go`):

```go
func (r *queryResolver) MyNewQuery(ctx context.Context, filter *MyFilter) (*MyResponse, error) {
    // Llamar a servicio de aplicación
    return r.myService.GetData(ctx, filter)
}
```

### Testing Backend

```bash
# Ejecutar todos los tests
cd api
go test ./internal/... -v

# Tests con cobertura
go test ./internal/... -cover

# Tests específicos
go test ./internal/domain/stock -v
```

---

## 🎨 Desarrollo Frontend

### Iniciar el Servidor de Desarrollo

El frontend se inicia automáticamente al abrir el devcontainer. Si necesitas reiniciarlo:

```bash
# Reiniciar frontend
make dev-restart-frontend

# Ver logs
make dev-logs-frontend
```

### Hot Reload

Vite proporciona hot reload automático. Los cambios se reflejan instantáneamente en el navegador.

### Storybook

Storybook está disponible en http://localhost:6006 y se inicia automáticamente junto con el frontend.

```bash
# Iniciar Storybook manualmente (si es necesario)
cd frontend
npm run storybook
```

### Agregar un Nuevo Componente

1. **Crear el componente** (`frontend/src/design-system/components/MyComponent.vue`):

```vue
<script setup lang="ts">
interface Props {
  title: string;
}
defineProps<Props>();
</script>

<template>
  <div class="my-component">
    <h2>{{ title }}</h2>
  </div>
</template>
```

2. **Crear story** (`frontend/src/design-system/components/MyComponent.stories.ts`):

```typescript
import type { Meta, StoryObj } from '@storybook/vue3';
import MyComponent from './MyComponent.vue';

const meta: Meta<typeof MyComponent> = {
  title: 'Design System/MyComponent',
  component: MyComponent,
};

export default meta;
type Story = StoryObj<typeof MyComponent>;

export const Default: Story = {
  args: {
    title: 'Mi Componente',
  },
};
```

3. **Exportar desde index** (`frontend/src/design-system/components/index.ts`):

```typescript
export { default as MyComponent } from './MyComponent.vue';
```

### Testing Frontend

```bash
# Ejecutar tests en modo watch
cd frontend
npm run test

# Tests con UI
npm run test:ui

# Tests con cobertura
npm run test:coverage
```

---

## 🧪 Testing

### Backend

Los tests están organizados por capa:

- `domain/*/service_test.go` - Tests de servicios de dominio
- `infrastructure/repository/*_test.go` - Tests de repositorios
- `application/graphql/resolvers_test.go` - Tests de resolvers

**Ejemplo de test**:

```go
func TestStockService_GetStocks(t *testing.T) {
    // Arrange
    mockRepo := &MockStockRepository{}
    service := NewStockService(mockRepo)
    
    // Act
    stocks, err := service.GetStocks(context.Background(), nil, nil)
    
    // Assert
    assert.NoError(t, err)
    assert.NotNil(t, stocks)
}
```

### Frontend

Los tests están en `frontend/src/test/`:

- `components/` - Tests de componentes
- `hoc/` - Tests de HOCs
- `stores/` - Tests de stores
- `composables/` - Tests de composables

**Ejemplo de test**:

```typescript
import { describe, it, expect } from 'vitest';
import { mount } from '@vue/test-utils';
import MyComponent from '@/components/MyComponent.vue';

describe('MyComponent', () => {
  it('renders correctly', () => {
    const wrapper = mount(MyComponent, {
      props: { title: 'Test' },
    });
    expect(wrapper.text()).toContain('Test');
  });
});
```

---

## 🐛 Debugging

### Backend

1. **Logs del servidor**: `make dev-logs-api`
2. **GraphQL Playground**: http://localhost:8080/playground
3. **Debugger**: Configurar breakpoints en VS Code/Cursor

### Frontend

1. **Vue DevTools**: Instalar extensión del navegador
2. **Console logs**: Usar `console.log` (remover antes de commit)
3. **Network tab**: Verificar requests GraphQL en DevTools

### Base de Datos

```bash
# Conectar a CockroachDB
docker exec -it go-react-test-cockroachdb ./cockroach sql --insecure

# Ver stocks
SELECT * FROM stocks LIMIT 10;
```

---

## 🤝 Contribución

### Flujo de Trabajo

1. **Crear branch**:
   ```bash
   git checkout -b feature/mi-nueva-funcionalidad
   ```

2. **Desarrollar**:
   - Seguir convenciones de código
   - Escribir tests
   - Actualizar documentación si es necesario

3. **Commit**:
   ```bash
   git add .
   git commit -m "feat: agregar nueva funcionalidad"
   ```

4. **Push y crear PR**:
   ```bash
   git push origin feature/mi-nueva-funcionalidad
   ```

### Convenciones de Commits

Usar [Conventional Commits](https://www.conventionalcommits.org/):

- `feat:` Nueva funcionalidad
- `fix:` Corrección de bug
- `docs:` Documentación
- `style:` Formato (sin cambios de código)
- `refactor:` Refactorización
- `test:` Tests
- `chore:` Tareas de mantenimiento

### Checklist antes de PR

- [ ] Código sigue las convenciones
- [ ] Tests pasan (`make test` o `npm test`)
- [ ] Documentación actualizada
- [ ] Sin warnings de linter
- [ ] Commits con mensajes descriptivos

---

## 📚 Recursos Adicionales

- [Arquitectura](./ARCHITECTURE.md) - Arquitectura DDD del proyecto
- [GraphQL API Reference](./GRAPHQL_API_REFERENCE.md) - Referencia completa de la API
- [Algoritmos](./ALGORITHMS.md) - Algoritmos y optimizaciones
- [Frontend](./FRONTEND.md) - Guía específica del frontend
- [Infraestructura](./INFRASTRUCTURE.md) - Configuración de infraestructura
- [DevContainer](./DEVCONTAINER.md) - Configuración del dev container

---

**Última actualización**: 2026-01-15
