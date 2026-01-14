# Optimización de Bundle

Guía para analizar y optimizar el tamaño del bundle.

## 📊 Análizar Bundle

### Usar Bundle Analyzer

```bash
# Analizar bundle después del build
npm run build:analyze

# O usar el visualizer directamente
npm run analyze
```

Esto generará un reporte HTML en `dist/stats.html` mostrando:

- Tamaño de cada chunk
- Dependencias
- Tamaño gzip y brotli

## 🎯 Optimizaciones Implementadas

### 1. Code Splitting

Las rutas usan lazy loading:

```typescript
const StockList = () => import("@/views/StockList.vue");
```

### 2. Vendor Chunks

Separación de vendors en `vite.config.ts`:

- `vue-vendor`: Vue, Vue Router, Pinia
- `graphql-vendor`: @urql/core, @urql/vue, graphql

### 3. Tree Shaking

Vite automáticamente hace tree shaking de imports no usados.

## 📈 Métricas Objetivo

- Bundle inicial < 200KB (gzipped)
- Chunks individuales < 100KB
- Tiempo de carga inicial < 2s

## 🔍 Identificar Problemas

### Chunks Grandes

Si un chunk es muy grande:

1. Verificar imports innecesarios
2. Considerar lazy loading adicional
3. Separar en chunks más pequeños

### Dependencias Duplicadas

Verificar duplicados:

```bash
npm run build:analyze
```

Buscar la misma librería en múltiples chunks.

## 🛠️ Mejoras Futuras

- [ ] Lazy load de componentes pesados
- [ ] Preload de rutas críticas
- [ ] Optimización de imágenes (si se agregan)
- [ ] Service Worker para cache
