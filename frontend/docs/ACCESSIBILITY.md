# Guía de Accesibilidad

Este documento describe las prácticas de accesibilidad implementadas en el frontend.

## 🎯 Objetivos

- Cumplir con **WCAG 2.1 nivel AA**
- Score Lighthouse > 90 en accesibilidad
- Navegación completa por teclado
- Soporte para screen readers

## ✅ Implementaciones

### 1. ARIA Labels

Todos los componentes incluyen atributos ARIA apropiados:

- `aria-label` para botones sin texto visible
- `aria-live` para contenido dinámico
- `aria-busy` para estados de carga
- `aria-sort` para tablas ordenables
- `role` apropiado para elementos semánticos

### 2. Navegación por Teclado

#### Tabla (Table Component)

- **Tab**: Navegar entre filas clickeables
- **Enter/Space**: Seleccionar fila
- **Arrow Up/Down**: Navegar entre filas
- **Home/End**: Ir al inicio/fin de la lista

#### Utilidades

- `handleKeyboardNavigation()`: Helper para navegación en listas
- `createFocusTrap()`: Para modales y diálogos

### 3. Contraste de Colores

Utilidad `checkContrast()` para verificar contraste WCAG AA:

```typescript
import { checkContrast } from "@/utils/accessibility";

const passes = checkContrast("#ffffff", "#000000"); // true
```

### 4. Screen Readers

- `announceToScreenReader()`: Anuncia cambios importantes
- Texto oculto con clase `.sr-only` para contexto adicional
- Estructura semántica HTML correcta

### 5. Focus Management

- Focus visible en todos los elementos interactivos
- Focus trap en modales
- Restauración de focus al cerrar modales

## 🧪 Testing de Accesibilidad

### Storybook Addon

El addon `@storybook/addon-a11y` está configurado para auditar accesibilidad:

```bash
npm run storybook
```

En Storybook, verás el panel "Accessibility" con:

- Violaciones de ARIA
- Problemas de contraste
- Recomendaciones de accesibilidad

### Lighthouse

Ejecutar auditoría de accesibilidad:

```bash
# En Chrome DevTools
# Lighthouse > Accessibility > Generate report
```

## 📋 Checklist de Componentes

Para cada componente nuevo, verificar:

- [ ] ARIA labels apropiados
- [ ] Navegación por teclado funcional
- [ ] Contraste de colores WCAG AA
- [ ] Focus visible
- [ ] Texto alternativo para imágenes
- [ ] Estructura semántica HTML

## 🔧 Utilidades Disponibles

### `@/utils/accessibility`

- `checkContrast(color1, color2)`: Verifica contraste
- `handleKeyboardNavigation()`: Navegación por teclado
- `createFocusTrap()`: Focus trap para modales
- `announceToScreenReader()`: Anuncia a screen readers

## 📚 Recursos

- [WCAG 2.1 Guidelines](https://www.w3.org/WAI/WCAG21/quickref/)
- [ARIA Authoring Practices](https://www.w3.org/WAI/ARIA/apg/)
- [WebAIM Contrast Checker](https://webaim.org/resources/contrastchecker/)
