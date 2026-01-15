# Guía del Frontend

## 📋 Resumen

Esta guía describe la arquitectura, componentes y desarrollo del frontend del proyecto.

---

## 🏗️ Arquitectura

### Stack Tecnológico

- **Framework**: Vue 3 (Composition API)
- **Lenguaje**: TypeScript
- **Build Tool**: Vite
- **State Management**: Pinia
- **GraphQL Client**: urql
- **Styling**: Tailwind CSS
- **Componentes**: Design System documentado en Storybook

### Estructura de Carpetas

```
frontend/
├── src/
│   ├── design-system/          # Design System
│   │   ├── components/         # Componentes reusables
│   │   ├── tokens/            # Design tokens
│   │   └── themes/            # Temas (light/dark)
│   ├── hoc/                   # Higher Order Components
│   ├── views/                 # Vistas/páginas
│   ├── stores/                # Stores de Pinia
│   ├── composables/           # Composables Vue
│   ├── utils/                 # Utilidades
│   └── router/                # Configuración de rutas
├── .storybook/                # Configuración Storybook
└── public/                    # Archivos estáticos
```

---

## 🎨 Design System

### Componentes Base

Los componentes están documentados en Storybook (http://localhost:6006):

1. **Button**: Variantes, estados, accesibilidad
2. **Input**: Búsqueda, validación
3. **Table**: Ordenamiento, paginación
4. **Card**: Variantes, estados
5. **ThemeToggle**: Cambio de tema

### Uso de Componentes

```vue
<script setup lang="ts">
import { Button, Input, Table } from '@/design-system/components';
</script>

<template>
  <Button variant="primary" @click="handleClick">
    Click me
  </Button>
</template>
```

### Design Tokens

Los tokens están en `src/design-system/tokens/`:

- **Colores**: `colors.ts`
- **Espaciado**: `spacing.ts`
- **Tipografía**: `typography.ts`
- **Breakpoints**: `breakpoints.ts`

### Temas

El sistema soporta temas light y dark:

```typescript
import { useThemeStore } from '@/stores/theme';

const themeStore = useThemeStore();
themeStore.toggleTheme(); // Cambia entre light/dark
```

---

## 🔧 Higher Order Components (HOCs)

Los HOCs están en `src/hoc/`:

### withLoading

Muestra un spinner mientras carga:

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

### withError

Maneja y muestra errores:

```vue
<script setup lang="ts">
import { withError } from '@/hoc/withError';

const StockListWithError = withError(StockList);
const error = ref<string | null>(null);
</script>

<template>
  <StockListWithError :error="error" @retry="handleRetry" />
</template>
```

### withPagination

Agrega lógica de paginación:

```vue
<script setup lang="ts">
import { withPagination } from '@/hoc/withPagination';

const StockListWithPagination = withPagination(StockList);
</script>

<template>
  <StockListWithPagination
    :current-page="currentPage"
    :total-pages="totalPages"
    @page-change="handlePageChange"
  />
</template>
```

### withSearch

Agrega búsqueda con debounce:

```vue
<script setup lang="ts">
import { withSearch } from '@/hoc/withSearch';

const StockListWithSearch = withSearch(StockList);
</script>

<template>
  <StockListWithSearch @search="handleSearch" />
</template>
```

**Más información**: Ver la sección de HOCs en este documento.

---

## 🎣 Composables

Los composables están en `src/composables/`:

### useStock

Lógica de acciones con cache:

```typescript
import { useStock } from '@/composables/useStock';

const { stocks, loading, error, fetchStocks } = useStock();

await fetchStocks({ ratings: ['Buy'] });
```

### useApi

Cliente GraphQL con cache:

```typescript
import { useApi } from '@/composables/useApi';

const { query, loading, error } = useApi();

const result = await query(GET_STOCKS_QUERY, variables);
```

### useRecommendations

Lógica de recomendaciones:

```typescript
import { useRecommendations } from '@/composables/useRecommendations';

const { recommendations, loading, fetchRecommendations } = useRecommendations();

await fetchRecommendations(10);
```

### useDebounce

Debounce de valores:

```typescript
import { useDebounce } from '@/composables/useDebounce';

const searchQuery = ref('');
const debouncedQuery = useDebounce(searchQuery, 300);
```

---

## 🗄️ State Management (Pinia)

### Stores

Los stores están en `src/stores/`:

#### theme.ts

Maneja el tema (light/dark):

```typescript
import { useThemeStore } from '@/stores/theme';

const themeStore = useThemeStore();
themeStore.toggleTheme();
```

#### stock.ts

Maneja el estado de stocks:

```typescript
import { useStockStore } from '@/stores/stock';

const stockStore = useStockStore();
await stockStore.fetchStocks(filter);
```

---

## 🌐 GraphQL

### Cliente

El cliente GraphQL está configurado en `src/utils/api.ts`:

```typescript
import { graphqlClient, GET_STOCKS_QUERY } from '@/utils/api';

const result = await graphqlClient.query(GET_STOCKS_QUERY, {
  filter: { ratings: ['Buy'] },
  limit: 50,
});
```

### Queries Disponibles

Ver [GraphQL API Reference](./GRAPHQL_API_REFERENCE.md) para todas las queries disponibles.

---

## 🎨 Styling

### Tailwind CSS

El proyecto usa Tailwind CSS para estilos:

```vue
<template>
  <div class="flex items-center justify-between p-4 bg-white dark:bg-gray-800">
    <h1 class="text-2xl font-bold">Stocks</h1>
  </div>
</template>
```

### Design Tokens

Usar tokens cuando sea posible:

```typescript
import { colors, spacing } from '@/design-system/tokens';

// En lugar de valores hardcodeados
const style = {
  color: colors.primary,
  padding: spacing.md,
};
```

---

## 📚 Storybook

### Acceso

Storybook está disponible en http://localhost:6006

### Agregar Stories

Crear `ComponentName.stories.ts`:

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

---

## 🧪 Testing

### Ejecutar Tests

```bash
# Tests en modo watch
npm run test

# Tests con UI
npm run test:ui

# Tests con cobertura
npm run test:coverage
```

**Más información**: Ver [Testing](./TESTING.md)

---

## 🚀 Desarrollo

### Hot Reload

Vite proporciona hot reload automático. Los cambios se reflejan instantáneamente.

### Agregar Nueva Vista

1. Crear componente en `src/views/MyView.vue`
2. Agregar ruta en `src/router/index.ts`:

```typescript
{
  path: '/my-view',
  name: 'MyView',
  component: () => import('@/views/MyView.vue'),
}
```

### Agregar Nuevo Composable

Crear `src/composables/useMyComposable.ts`:

```typescript
import { ref, computed } from 'vue';

export function useMyComposable() {
  const data = ref(null);
  const loading = ref(false);

  const fetchData = async () => {
    loading.value = true;
    // Lógica aquí
    loading.value = false;
  };

  return {
    data,
    loading,
    fetchData,
  };
}
```

---

## ♿ Accesibilidad

### Checklist

- [ ] ARIA labels en elementos interactivos
- [ ] Navegación por teclado funcional
- [ ] Contraste WCAG AA
- [ ] Focus visible
- [ ] Texto alternativo para imágenes

### Utilidades

```typescript
import {
  checkContrast,
  handleKeyboardNavigation,
  announceToScreenReader,
} from '@/utils/accessibility';
```

---

## 📦 Build y Optimización

### Build de Producción

```bash
npm run build
```

### Análisis de Bundle

```bash
npm run build:analyze
```

Ver `dist/stats.html` para análisis detallado.

---

## 📚 Recursos

- [Vue 3 Documentation](https://vuejs.org/)
- [Vite Documentation](https://vitejs.dev/)
- [Pinia Documentation](https://pinia.vuejs.org/)
- [Tailwind CSS Documentation](https://tailwindcss.com/)
- [Storybook Documentation](https://storybook.js.org/)

---

**Última actualización**: 2026-01-15
