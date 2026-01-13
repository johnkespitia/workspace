# Frontend - Sistema de Información de Acciones

Frontend desarrollado con Vue 3, TypeScript, Tailwind CSS y GraphQL.

## 🚀 Características

- **Vue 3** con Composition API
- **TypeScript** para type safety
- **Pinia** para gestión de estado
- **Vue Router** para navegación
- **GraphQL** con @urql/core para comunicación con el backend
- **Tailwind CSS** para estilos
- **Storybook** para documentación de componentes
- **Design System** completo con tokens, componentes y temas

## 📁 Estructura del Proyecto

```
frontend/
├── src/
│   ├── design-system/          # Design System
│   │   ├── components/         # Componentes reusables
│   │   ├── tokens/             # Design tokens (colores, espaciado, tipografía)
│   │   └── themes/             # Temas (light/dark)
│   ├── hoc/                    # Higher Order Components
│   │   ├── withLoading.ts
│   │   ├── withError.ts
│   │   ├── withPagination.ts
│   │   └── withSearch.ts
│   ├── composables/           # Composables Vue
│   │   ├── useApi.ts
│   │   ├── useStock.ts
│   │   ├── useRecommendations.ts
│   │   └── useDebounce.ts
│   ├── stores/                 # Stores de Pinia
│   │   ├── stock.ts
│   │   └── theme.ts
│   ├── views/                  # Vistas/páginas
│   │   ├── StockList.vue
│   │   ├── StockDetail.vue
│   │   └── Recommendations.vue
│   ├── router/                 # Configuración de rutas
│   │   └── index.ts
│   ├── utils/                  # Utilidades
│   │   └── api.ts              # Cliente GraphQL
│   ├── App.vue
│   └── main.ts
└── .storybook/                 # Configuración Storybook
```

## 🛠️ Instalación

```bash
# Instalar dependencias
npm install

# Iniciar servidor de desarrollo
npm run dev

# Construir para producción
npm run build

# Iniciar Storybook
npm run storybook
```

## 🎨 Design System

El proyecto incluye un Design System completo con:

- **Design Tokens**: Colores, espaciado, tipografía
- **Temas**: Light y Dark mode
- **Componentes Base**: Button, Input, Table, Card, ThemeToggle
- **HOCs**: withLoading, withError, withPagination, withSearch
- **Accesibilidad**: ARIA labels, navegación por teclado, contraste WCAG AA

## 📡 Configuración de API

El endpoint de GraphQL se configura mediante variable de entorno:

```env
VITE_GRAPHQL_ENDPOINT=http://localhost:8080/query
```

Por defecto usa `http://localhost:8080/query`.

## 🧩 Componentes Principales

### HOCs (Higher Order Components)

Los HOCs permiten agregar funcionalidad a componentes:

```typescript
import { withLoading } from '@/hoc/withLoading';
import StockList from './StockList.vue';

const StockListWithLoading = withLoading(StockList);
```

### Composables

Los composables proporcionan lógica reutilizable:

```typescript
import { useStock } from '@/composables/useStock';

const stock = useStock();
await stock.loadStocks();
```

## 🧪 Testing

Para ejecutar Storybook y ver los componentes:

```bash
npm run storybook
```

## 📝 Notas

- El frontend está completamente tipado con TypeScript
- Todos los componentes incluyen soporte de accesibilidad
- El tema dark/light se persiste en localStorage
- Las queries GraphQL incluyen cache y deduplicación de requests
