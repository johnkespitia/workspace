# Testing - Frontend

Este directorio contiene los tests del frontend usando Vitest y Vue Test Utils.

## Estructura

```
test/
├── setup.ts              # Configuración global de tests
├── components/           # Tests de componentes
│   └── Button.test.ts
├── hoc/                  # Tests de Higher Order Components
│   └── withLoading.test.ts
├── stores/               # Tests de stores de Pinia
│   └── theme.test.ts
└── composables/          # Tests de composables
    └── useDebounce.test.ts
```

## Ejecutar Tests

```bash
# Ejecutar tests en modo watch
npm run test

# Ejecutar tests con UI
npm run test:ui

# Ejecutar tests una vez
npm run test:run

# Ejecutar tests con coverage
npm run test:coverage
```

## Escribir Tests

### Componente Vue

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

### Store de Pinia

```typescript
import { describe, it, expect, beforeEach } from "vitest";
import { setActivePinia, createPinia } from "pinia";
import { useMyStore } from "@/stores/my";

describe("MyStore", () => {
  beforeEach(() => {
    setActivePinia(createPinia());
  });

  it("initializes correctly", () => {
    const store = useMyStore();
    expect(store.value).toBe(0);
  });
});
```

### Composable

```typescript
import { describe, it, expect } from "vitest";
import { useMyComposable } from "@/composables/useMy";

describe("useMyComposable", () => {
  it("returns expected value", () => {
    const { value } = useMyComposable();
    expect(value.value).toBeDefined();
  });
});
```

## Coverage

El coverage se genera en `coverage/` después de ejecutar `npm run test:coverage`.

Objetivo: > 70% de cobertura según métricas de éxito del plan.
