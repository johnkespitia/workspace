# Higher Order Components (HOCs) en Vue 3

## 📚 Concepto

**Higher Order Components (HOCs)** son funciones que toman un componente y retornan un nuevo componente con funcionalidad adicional. 

**IMPORTANTE**: Los HOCs **NO son hooks/composables**. Son un patrón de diseño para componer componentes.

### Diferencia entre HOCs y Composables

| Aspecto | HOCs | Composables |
|---------|------|-------------|
| **Propósito** | Envolver componentes, modificar estructura | Proporcionar lógica reutilizable |
| **Uso** | Se aplican al componente | Se llaman dentro de `setup()` |
| **Retorno** | Nuevo componente | Valores reactivos, funciones |
| **Ejemplo** | `withLoading(MyComponent)` | `const { data, loading } = useApi()` |

---

## 🏗️ Implementación en Vue 3

### Estructura Básica de un HOC

```typescript
import { defineComponent, h, Component } from 'vue';

export function withFeature<T extends Component>(
  WrappedComponent: T,
  options?: FeatureOptions
) {
  return defineComponent({
    name: `withFeature(${WrappedComponent.name || 'Component'})`,
    props: {
      // Props del componente envuelto
      ...WrappedComponent.props,
      // Props adicionales del HOC
      featureProp: {
        type: Boolean,
        default: false
      }
    },
    setup(props, { slots, attrs, emit }) {
      // Lógica del HOC
      const featureState = ref(false);
      
      return () => {
        // Renderizado condicional o modificado
        if (props.featureProp) {
          return h('div', { class: 'feature-wrapper' }, [
            h(WrappedComponent, { ...attrs, ...props }, slots)
          ]);
        }
        return h(WrappedComponent, { ...attrs, ...props }, slots);
      };
    }
  });
}
```

---

## 🎯 HOCs a Implementar

### 1. withLoading

**Propósito**: Mostrar un spinner o estado de carga mientras se cargan datos.

```typescript
// hoc/withLoading.ts
import { defineComponent, h, Component } from 'vue';

export function withLoading<T extends Component>(WrappedComponent: T) {
  return defineComponent({
    name: `withLoading(${WrappedComponent.name || 'Component'})`,
    props: {
      loading: {
        type: Boolean,
        default: false
      },
      loadingText: {
        type: String,
        default: 'Cargando...'
      }
    },
    setup(props, { slots, attrs }) {
      return () => {
        if (props.loading) {
          return h('div', { 
            class: 'loading-container',
            'aria-live': 'polite',
            'aria-busy': 'true'
          }, [
            h('div', { class: 'spinner' }, props.loadingText)
          ]);
        }
        return h(WrappedComponent, attrs, slots);
      };
    }
  });
}
```

**Uso**:
```vue
<script setup lang="ts">
import StockList from './StockList.vue';
import { withLoading } from '@/hoc/withLoading';

const StockListWithLoading = withLoading(StockList);
const loading = ref(false);
</script>

<template>
  <StockListWithLoading :loading="loading" />
</template>
```

---

### 2. withError

**Propósito**: Manejar y mostrar errores de forma consistente.

```typescript
// hoc/withError.ts
import { defineComponent, h, Component } from 'vue';

export function withError<T extends Component>(WrappedComponent: T) {
  return defineComponent({
    name: `withError(${WrappedComponent.name || 'Component'})`,
    props: {
      error: {
        type: [String, Error, null] as PropType<string | Error | null>,
        default: null
      },
      errorMessage: {
        type: String,
        default: 'Ha ocurrido un error'
      }
    },
    setup(props, { slots, attrs }) {
      const getErrorMessage = () => {
        if (!props.error) return null;
        if (typeof props.error === 'string') return props.error;
        return props.error.message || props.errorMessage;
      };

      return () => {
        const errorMsg = getErrorMessage();
        
        if (errorMsg) {
          return h('div', { 
            class: 'error-container',
            role: 'alert',
            'aria-live': 'assertive'
          }, [
            h('div', { class: 'error-icon' }, '⚠️'),
            h('div', { class: 'error-message' }, errorMsg),
            h('button', {
              class: 'error-retry',
              onClick: () => {
                // Emitir evento para reintentar
                emit('retry');
              }
            }, 'Reintentar')
          ]);
        }
        
        return h(WrappedComponent, attrs, slots);
      };
    }
  });
}
```

**Uso**:
```vue
<script setup lang="ts">
import StockList from './StockList.vue';
import { withError } from '@/hoc/withError';

const StockListWithError = withError(StockList);
const error = ref<string | null>(null);
</script>

<template>
  <StockListWithError :error="error" @retry="handleRetry" />
</template>
```

---

### 3. withPagination

**Propósito**: Agregar lógica de paginación a componentes de lista.

```typescript
// hoc/withPagination.ts
import { defineComponent, h, Component, computed } from 'vue';

interface PaginationProps {
  currentPage: number;
  totalPages: number;
  pageSize: number;
  totalItems: number;
}

export function withPagination<T extends Component>(WrappedComponent: T) {
  return defineComponent({
    name: `withPagination(${WrappedComponent.name || 'Component'})`,
    props: {
      currentPage: {
        type: Number,
        default: 1
      },
      totalPages: {
        type: Number,
        required: true
      },
      pageSize: {
        type: Number,
        default: 10
      },
      totalItems: {
        type: Number,
        required: true
      }
    },
    emits: ['page-change'],
    setup(props, { slots, attrs, emit }) {
      const paginationInfo = computed(() => ({
        currentPage: props.currentPage,
        totalPages: props.totalPages,
        pageSize: props.pageSize,
        totalItems: props.totalItems,
        startItem: (props.currentPage - 1) * props.pageSize + 1,
        endItem: Math.min(props.currentPage * props.pageSize, props.totalItems)
      }));

      const goToPage = (page: number) => {
        if (page >= 1 && page <= props.totalPages) {
          emit('page-change', page);
        }
      };

      return () => {
        return h('div', { class: 'pagination-wrapper' }, [
          h(WrappedComponent, {
            ...attrs,
            pagination: paginationInfo.value
          }, slots),
          h('div', { class: 'pagination-controls' }, [
            h('button', {
              disabled: props.currentPage === 1,
              onClick: () => goToPage(props.currentPage - 1),
              'aria-label': 'Página anterior'
            }, '← Anterior'),
            h('span', { class: 'page-info' }, 
              `Página ${props.currentPage} de ${props.totalPages}`
            ),
            h('button', {
              disabled: props.currentPage === props.totalPages,
              onClick: () => goToPage(props.currentPage + 1),
              'aria-label': 'Página siguiente'
            }, 'Siguiente →')
          ])
        ]);
      };
    }
  });
}
```

**Uso**:
```vue
<script setup lang="ts">
import StockList from './StockList.vue';
import { withPagination } from '@/hoc/withPagination';

const StockListWithPagination = withPagination(StockList);
const currentPage = ref(1);
const totalPages = computed(() => Math.ceil(stocks.value.length / 10));
</script>

<template>
  <StockListWithPagination
    :current-page="currentPage"
    :total-pages="totalPages"
    :total-items="stocks.length"
    @page-change="currentPage = $event"
  />
</template>
```

---

### 4. withSearch

**Propósito**: Agregar funcionalidad de búsqueda con debounce.

```typescript
// hoc/withSearch.ts
import { defineComponent, h, Component, ref, watch } from 'vue';
import { useDebounce } from '@/composables/useDebounce';

export function withSearch<T extends Component>(WrappedComponent: T) {
  return defineComponent({
    name: `withSearch(${WrappedComponent.name || 'Component'})`,
    props: {
      searchPlaceholder: {
        type: String,
        default: 'Buscar...'
      },
      debounceMs: {
        type: Number,
        default: 300
      }
    },
    emits: ['search'],
    setup(props, { slots, attrs, emit }) {
      const searchQuery = ref('');
      const debouncedQuery = useDebounce(searchQuery, props.debounceMs);

      watch(debouncedQuery, (newQuery) => {
        emit('search', newQuery);
      });

      return () => {
        return h('div', { class: 'search-wrapper' }, [
          h('div', { class: 'search-input-container' }, [
            h('input', {
              type: 'text',
              class: 'search-input',
              placeholder: props.searchPlaceholder,
              value: searchQuery.value,
              onInput: (e: Event) => {
                searchQuery.value = (e.target as HTMLInputElement).value;
              },
              'aria-label': 'Campo de búsqueda'
            }),
            h('span', { class: 'search-icon' }, '🔍')
          ]),
          h(WrappedComponent, {
            ...attrs,
            searchQuery: debouncedQuery.value
          }, slots)
        ]);
      };
    }
  });
}
```

**Uso**:
```vue
<script setup lang="ts">
import StockList from './StockList.vue';
import { withSearch } from '@/hoc/withSearch';

const StockListWithSearch = withSearch(StockList);

const handleSearch = (query: string) => {
  // Filtrar stocks basado en query
  filteredStocks.value = stocks.value.filter(s => 
    s.ticker.toLowerCase().includes(query.toLowerCase()) ||
    s.companyName.toLowerCase().includes(query.toLowerCase())
  );
};
</script>

<template>
  <StockListWithSearch
    search-placeholder="Buscar por ticker o compañía..."
    @search="handleSearch"
  />
</template>
```

---

## 🔄 Composición de Múltiples HOCs

Puedes componer múltiples HOCs para agregar varias funcionalidades:

```typescript
// Componer HOCs
const StockListEnhanced = withLoading(
  withError(
    withPagination(
      withSearch(StockList)
    )
  )
);
```

O usando una función helper:

```typescript
// utils/composeHOCs.ts
export function composeHOCs(...hocs: Function[]) {
  return (component: Component) => {
    return hocs.reduceRight(
      (acc, hoc) => hoc(acc),
      component
    );
  };
}

// Uso:
const enhanceStockList = composeHOCs(
  withLoading,
  withError,
  withPagination,
  withSearch
);

const StockListEnhanced = enhanceStockList(StockList);
```

---

## ✅ Ventajas de HOCs en Vue 3

1. **Separación de responsabilidades**: Lógica reutilizable separada de componentes
2. **Composición**: Fácil combinar múltiples HOCs
3. **Testabilidad**: HOCs se pueden testear independientemente
4. **Flexibilidad**: Mismo componente con diferentes funcionalidades

## ⚠️ Consideraciones

1. **Props forwarding**: Asegúrate de pasar todas las props al componente envuelto
2. **Slots forwarding**: Pasa los slots correctamente
3. **Nombres de componentes**: Usa nombres descriptivos para debugging
4. **Accesibilidad**: Mantén atributos ARIA cuando sea necesario

---

**Última actualización**: [Fecha]
