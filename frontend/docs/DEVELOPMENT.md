# Guía de Desarrollo

Esta guía te ayudará a entender la estructura del proyecto y cómo contribuir.

## 📁 Estructura del Proyecto

```
frontend/
├── src/
│   ├── design-system/       # Design System
│   │   ├── components/      # Componentes reusables
│   │   ├── tokens/          # Design tokens (colores, espaciado, etc.)
│   │   └── themes/          # Temas (light/dark)
│   ├── hoc/                 # Higher Order Components
│   ├── composables/         # Composables Vue (lógica reutilizable)
│   ├── stores/              # Stores de Pinia
│   ├── views/               # Vistas/páginas
│   ├── router/              # Configuración de rutas
│   ├── utils/               # Utilidades
│   ├── components/          # Componentes específicos de la app
│   ├── test/                # Tests
│   ├── App.vue              # Componente raíz
│   └── main.ts              # Punto de entrada
├── .storybook/              # Configuración Storybook
├── docs/                    # Documentación
└── public/                  # Archivos estáticos
```

## 🚀 Inicio Rápido

### Prerrequisitos

- Node.js 18+ y npm
- Conocimientos básicos de Vue 3, TypeScript y Tailwind CSS

### Instalación

```bash
# Instalar dependencias
npm install

# Iniciar servidor de desarrollo
npm run dev

# El frontend estará disponible en http://localhost:3000
```

## 🛠️ Scripts Disponibles

```bash
# Desarrollo
npm run dev              # Servidor de desarrollo
npm run build            # Build para producción
npm run preview          # Preview del build

# Testing
npm run test             # Tests en modo watch
npm run test:ui          # Tests con UI
npm run test:coverage    # Tests con coverage
npm run test:run         # Tests una vez

# Storybook
npm run storybook        # Iniciar Storybook
npm run build-storybook  # Build de Storybook

# Análisis
npm run build:analyze    # Analizar bundle
npm run analyze          # Visualizar bundle
```

## 🏗️ Arquitectura

### Design System

El Design System está en `src/design-system/`:

- **Tokens**: Colores, espaciado, tipografía, breakpoints
- **Componentes**: Button, Input, Table, Card, ThemeToggle
- **Temas**: Light y Dark mode

### Higher Order Components (HOCs)

Los HOCs están en `src/hoc/`:

- `withLoading`: Muestra spinner durante carga
- `withError`: Maneja errores
- `withPagination`: Agrega paginación
- `withSearch`: Agrega búsqueda con debounce

**Uso**:

```typescript
import { withLoading } from "@/hoc/withLoading";
const EnhancedComponent = withLoading(MyComponent);
```

### Composables

Los composables están en `src/composables/`:

- `useApi`: Cliente GraphQL con cache
- `useStock`: Lógica de acciones
- `useRecommendations`: Lógica de recomendaciones
- `useBreakpoint`: Detección de breakpoint reactivo
- `useVirtualScroll`: Virtual scrolling para listas grandes

**Uso**:

```typescript
import { useStock } from "@/composables/useStock";
const { stocks, loadStocks } = useStock();
```

### Stores (Pinia)

Los stores están en `src/stores/`:

- `theme`: Gestión de tema light/dark
- `stock`: Estado global de acciones

## 📝 Convenciones de Código

### Nombres de Archivos

- Componentes: `PascalCase.vue` (ej: `Button.vue`)
- Composables: `camelCase.ts` con prefijo `use` (ej: `useStock.ts`)
- Utils: `camelCase.ts` (ej: `accessibility.ts`)
- Tests: `ComponentName.test.ts`

### Estructura de Componentes

```vue
<template>
  <!-- Template -->
</template>

<script setup lang="ts">
// Imports
// Props
// Emits
// Composables
// Computed
// Methods
</script>

<style scoped>
/* Estilos */
</style>
```

### TypeScript

- Usar tipos explícitos
- Evitar `any`
- Usar interfaces para props y emits
- Exportar tipos cuando sean reutilizables

### Estilos

- Usar Tailwind CSS para estilos
- Usar clases de utilidad cuando sea posible
- Estilos scoped para componentes
- Design tokens para valores consistentes

## 🧪 Testing

### Escribir Tests

Los tests están en `src/test/`:

```typescript
import { describe, it, expect } from "vitest";
import { mount } from "@vue/test-utils";
import MyComponent from "@/components/MyComponent.vue";

describe("MyComponent", () => {
  it("renders correctly", () => {
    const wrapper = mount(MyComponent);
    expect(wrapper.text()).toContain("Hello");
  });
});
```

### Cobertura

Objetivo: > 70% de cobertura

```bash
npm run test:coverage
```

## 📚 Storybook

### Agregar Stories

Crear `ComponentName.stories.ts`:

```typescript
import type { Meta, StoryObj } from "@storybook/vue3";
import MyComponent from "./MyComponent.vue";

const meta: Meta<typeof MyComponent> = {
  title: "Design System/MyComponent",
  component: MyComponent,
};

export default meta;
type Story = StoryObj<typeof MyComponent>;

export const Default: Story = {
  args: {},
};
```

## 🔍 GraphQL

### Queries

Las queries están en `src/utils/api.ts`:

```typescript
import { graphqlClient, GET_STOCKS_QUERY } from "@/utils/api";

const result = await graphqlClient.query(GET_STOCKS_QUERY, variables);
```

### Tipos

Usar los tipos exportados de `@/utils/api`:

```typescript
import type { Stock, StockFilter, StockSort } from "@/utils/api";
```

## 🎨 Accesibilidad

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
} from "@/utils/accessibility";
```

## 🐛 Debugging

### Vue DevTools

Instalar extensión de navegador para Vue DevTools.

### Console Logs

Usar `console.log` para debugging, remover antes de commit.

### Source Maps

Activados por defecto en desarrollo.

## 📦 Build y Optimización

### Análizar Bundle

```bash
npm run build:analyze
```

Ver `dist/stats.html` para análisis detallado.

### Optimizaciones

- Code splitting por ruta (ya implementado)
- Lazy loading de componentes pesados
- Vendor chunks separados

## 🔄 Git Workflow

1. Crear branch desde `main`
2. Hacer cambios
3. Escribir tests
4. Ejecutar tests y linter
5. Crear PR

### Commits

Usar mensajes descriptivos:

```
feat: agregar gráfico en vista de detalle
fix: corregir navegación por teclado en tabla
docs: actualizar guía de desarrollo
```

## 📖 Recursos

- [Vue 3 Docs](https://vuejs.org/)
- [TypeScript Handbook](https://www.typescriptlang.org/docs/)
- [Tailwind CSS](https://tailwindcss.com/)
- [Pinia](https://pinia.vuejs.org/)
- [Vitest](https://vitest.dev/)

## ❓ Preguntas Frecuentes

### ¿Cómo agrego un nuevo componente?

1. Crear en `src/design-system/components/` si es reutilizable
2. O en `src/components/` si es específico de la app
3. Agregar story en Storybook
4. Escribir tests

### ¿Cómo agrego una nueva ruta?

Editar `src/router/index.ts`:

```typescript
{
  path: '/new-route',
  name: 'NewRoute',
  component: () => import('@/views/NewRoute.vue'),
}
```

### ¿Cómo uso el tema dark/light?

```typescript
import { useThemeStore } from "@/stores/theme";
const themeStore = useThemeStore();
themeStore.toggleTheme();
```

---

**¿Necesitas ayuda?** Revisa la documentación o crea un issue.
