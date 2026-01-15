# Guía de Testing

## 📋 Resumen

Este documento describe la estrategia de testing del proyecto, incluyendo tests del backend (Go) y frontend (Vue 3 + TypeScript).

---

## 🎯 Estrategia de Testing

### Cobertura Objetivo

- **Backend**: > 70% de cobertura
- **Frontend**: > 70% de cobertura

### Tipos de Tests

1. **Unit Tests**: Tests de unidades individuales (funciones, métodos)
2. **Integration Tests**: Tests de integración entre componentes
3. **E2E Tests**: Tests end-to-end (futuro)

---

## 🔧 Testing Backend (Go)

### Framework

- **Testing estándar de Go**: `testing` package
- **Assertions**: `github.com/stretchr/testify`
- **Mocks**: `github.com/DATA-DOG/go-sqlmock` para base de datos

### Estructura de Tests

Los tests están organizados por capa:

```
api/internal/
├── domain/
│   ├── stock/
│   │   ├── service.go
│   │   └── service_test.go        # Tests de servicios de dominio
│   └── recommendation/
│       ├── algorithm.go
│       └── algorithm_test.go      # Tests de algoritmo
├── infrastructure/
│   └── repository/
│       ├── stock_repository.go
│       └── stock_repository_test.go # Tests de repositorios
└── application/
    └── graphql/
        ├── resolvers.go
        └── resolvers_test.go       # Tests de resolvers
```

### Ejecutar Tests

```bash
# Todos los tests
cd api
go test ./internal/... -v

# Tests con cobertura
go test ./internal/... -cover

# Tests específicos
go test ./internal/domain/stock -v

# Tests con cobertura detallada
go test ./internal/... -coverprofile=coverage.out
go tool cover -html=coverage.out
```

### Ejemplos de Tests

#### Test de Servicio de Dominio

```go
// domain/stock/service_test.go
package stock

import (
    "context"
    "testing"
    "github.com/stretchr/testify/assert"
)

func TestDomainService_CalculatePriceChange(t *testing.T) {
    service := NewService()
    
    stock := &Stock{
        TargetFrom: NewPrice(100.0),
        TargetTo:   NewPrice(120.0),
    }
    
    change := service.CalculatePriceChange(stock)
    
    assert.Equal(t, 20.0, change)
}
```

#### Test de Repositorio con Mock

```go
// infrastructure/repository/stock_repository_test.go
package repository

import (
    "context"
    "testing"
    "github.com/DATA-DOG/go-sqlmock"
    "github.com/stretchr/testify/assert"
)

func TestCockroachStockRepository_Save(t *testing.T) {
    db, mock, err := sqlmock.New()
    assert.NoError(t, err)
    defer db.Close()
    
    repo := NewCockroachStockRepository(db)
    
    mock.ExpectExec("INSERT INTO stocks").
        WithArgs("AAPL", "Apple Inc.", ...).
        WillReturnResult(sqlmock.NewResult(1, 1))
    
    stock := &domain.Stock{
        Ticker:      "AAPL",
        CompanyName: "Apple Inc.",
    }
    
    err = repo.Save(context.Background(), stock)
    assert.NoError(t, err)
    assert.NoError(t, mock.ExpectationsWereMet())
}
```

### Cobertura Actual

Según `api/TEST_SUMMARY.md`:

- **Domain Layer**: ~70-80% cobertura
- **Infrastructure Layer**: ~52% cobertura
- **Application Layer**: Tests básicos de parsing
- **Cobertura Total**: ~49-52%

---

## 🎨 Testing Frontend (Vue 3 + TypeScript)

### Framework

- **Vitest**: Framework de testing
- **Vue Test Utils**: Utilidades para testing de componentes Vue
- **Testing Library**: Utilidades adicionales para testing

### Estructura de Tests

```
frontend/src/test/
├── setup.ts                    # Configuración global
├── components/                 # Tests de componentes
│   └── Button.test.ts
├── hoc/                        # Tests de HOCs
│   └── withLoading.test.ts
├── stores/                     # Tests de stores
│   └── theme.test.ts
└── composables/                # Tests de composables
    └── useDebounce.test.ts
```

### Ejecutar Tests

```bash
# Tests en modo watch
cd frontend
npm run test

# Tests con UI interactiva
npm run test:ui

# Tests una vez (sin watch)
npm run test:run

# Tests con cobertura
npm run test:coverage
```

### Ejemplos de Tests

#### Test de Componente

```typescript
// test/components/Button.test.ts
import { describe, it, expect } from 'vitest';
import { mount } from '@vue/test-utils';
import Button from '@/design-system/components/Button.vue';

describe('Button', () => {
  it('renders correctly', () => {
    const wrapper = mount(Button, {
      props: { label: 'Click me' },
    });
    expect(wrapper.text()).toContain('Click me');
  });

  it('emits click event', async () => {
    const wrapper = mount(Button);
    await wrapper.trigger('click');
    expect(wrapper.emitted('click')).toBeTruthy();
  });
});
```

#### Test de Store (Pinia)

```typescript
// test/stores/theme.test.ts
import { describe, it, expect, beforeEach } from 'vitest';
import { setActivePinia, createPinia } from 'pinia';
import { useThemeStore } from '@/stores/theme';

describe('ThemeStore', () => {
  beforeEach(() => {
    setActivePinia(createPinia());
  });

  it('initializes with light theme', () => {
    const store = useThemeStore();
    expect(store.theme).toBe('light');
  });

  it('toggles theme', () => {
    const store = useThemeStore();
    store.toggleTheme();
    expect(store.theme).toBe('dark');
  });
});
```

#### Test de Composable

```typescript
// test/composables/useDebounce.test.ts
import { describe, it, expect, vi } from 'vitest';
import { ref } from 'vue';
import { useDebounce } from '@/composables/useDebounce';

describe('useDebounce', () => {
  it('debounces value changes', async () => {
    vi.useFakeTimers();
    const value = ref('');
    const debounced = useDebounce(value, 300);

    value.value = 'a';
    value.value = 'ab';
    value.value = 'abc';

    vi.advanceTimersByTime(300);
    expect(debounced.value).toBe('abc');
    vi.useRealTimers();
  });
});
```

### Configuración de Vitest

```typescript
// vitest.config.ts
import { defineConfig } from 'vitest/config';
import vue from '@vitejs/plugin-vue';

export default defineConfig({
  plugins: [vue()],
  test: {
    globals: true,
    environment: 'jsdom',
    setupFiles: ['./src/test/setup.ts'],
    coverage: {
      provider: 'v8',
      reporter: ['text', 'json', 'html'],
      exclude: [
        'node_modules/',
        'src/test/',
        '**/*.d.ts',
        '**/*.config.*',
      ],
    },
  },
});
```

---

## 📊 Métricas de Cobertura

### Backend

Ejecutar con cobertura:

```bash
cd api
go test ./internal/... -coverprofile=coverage.out
go tool cover -func=coverage.out
```

### Frontend

```bash
cd frontend
npm run test:coverage
```

El reporte se genera en `frontend/coverage/`.

---

## 🎯 Mejores Prácticas

### Backend

1. **Tests independientes**: Cada test debe ser independiente
2. **Mocks apropiados**: Usar mocks para dependencias externas
3. **Nombres descriptivos**: `TestService_Method_Scenario`
4. **Arrange-Act-Assert**: Estructura clara de tests

### Frontend

1. **Tests de comportamiento**: Testear comportamiento, no implementación
2. **Accesibilidad**: Incluir tests de accesibilidad cuando sea relevante
3. **Snapshots con cuidado**: Usar snapshots solo cuando sea necesario
4. **Mock de APIs**: Mockear llamadas a APIs en tests

---

## 🚀 CI/CD (Futuro)

### Pipeline Sugerido

1. **Lint**: Verificar código
2. **Unit Tests**: Ejecutar tests unitarios
3. **Coverage**: Verificar cobertura mínima
4. **Build**: Verificar que el build funciona
5. **E2E Tests**: Tests end-to-end (futuro)

---

## 📚 Recursos

- [Go Testing Package](https://pkg.go.dev/testing)
- [Vitest Documentation](https://vitest.dev/)
- [Vue Test Utils](https://test-utils.vuejs.org/)
- [Testing Library](https://testing-library.com/)

---

**Última actualización**: 2026-01-15
