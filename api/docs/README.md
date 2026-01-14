# Documentación de la API

## 📚 Índice de Documentación

### Para Desarrolladores

1. **[API_DOCUMENTATION.md](./API_DOCUMENTATION.md)**

   - Documentación completa de endpoints HTTP y GraphQL
   - Ejemplos de requests/responses
   - Códigos de error
   - Guía de integración

2. **[openapi.yaml](./openapi.yaml)**

   - Especificación OpenAPI 3.0
   - Compatible con Swagger UI
   - Documentación de endpoints REST

3. **[USER_GUIDE.md](./USER_GUIDE.md)**
   - Guía de uso completa para usuarios
   - Casos de uso comunes
   - Ejemplos paso a paso
   - Solución de problemas

### Referencias Técnicas

4. **[../docs/GRAPHQL_API_REFERENCE.md](../../docs/GRAPHQL_API_REFERENCE.md)**

   - Referencia completa del schema GraphQL
   - Tipos, queries, mutations
   - Campos correctos para el frontend

5. **[../docs/GRAPHQL_EXAMPLES.md](../../docs/GRAPHQL_EXAMPLES.md)**
   - Ejemplos prácticos de queries
   - Ejemplos de mutations
   - Queries complejas

### Documentación de Arquitectura

6. **[../docs/ARCHITECTURE.md](../../docs/ARCHITECTURE.md)**

   - Arquitectura del sistema
   - Estructura DDD
   - Diagramas y flujos

7. **[../docs/ALGORITHMS.md](../../docs/ALGORITHMS.md)**
   - Algoritmo de recomendación
   - Complejidad y estrategias

---

## 🚀 Inicio Rápido

### Ver Documentación Interactiva

1. **GraphQL Playground**: `http://localhost:8080/playground`
2. **Swagger UI** (si se configura): Importar `openapi.yaml`

### Leer Documentación

- **Nuevo en la API?** → Empieza con [USER_GUIDE.md](./USER_GUIDE.md)
- **Integrando la API?** → Lee [API_DOCUMENTATION.md](./API_DOCUMENTATION.md)
- **Desarrollando frontend?** → Consulta [GRAPHQL_API_REFERENCE.md](../../docs/GRAPHQL_API_REFERENCE.md)

---

## 📖 Estructura de Documentos

```
api/
├── docs/
│   ├── README.md (este archivo)
│   ├── API_DOCUMENTATION.md (documentación completa)
│   ├── USER_GUIDE.md (guía de usuario)
│   └── openapi.yaml (especificación OpenAPI)
└── ...
docs/
├── GRAPHQL_API_REFERENCE.md (referencia GraphQL)
├── GRAPHQL_EXAMPLES.md (ejemplos GraphQL)
├── ARCHITECTURE.md (arquitectura)
└── ...
```

---

## 🔍 Búsqueda Rápida

### ¿Cómo...?

- **...obtener stocks?** → [USER_GUIDE.md - Caso 1](./USER_GUIDE.md#caso-1-ver-lista-de-stocks)
- **...filtrar por rating?** → [USER_GUIDE.md - Caso 2](./USER_GUIDE.md#caso-2-buscar-stocks-por-rating)
- **...sincronizar stocks?** → [USER_GUIDE.md - Caso 6](./USER_GUIDE.md#caso-6-sincronizar-stocks-desde-api-externa)
- **...obtener recomendaciones?** → [USER_GUIDE.md - Caso 5](./USER_GUIDE.md#caso-5-obtener-recomendaciones-de-inversión)
- **...integrar con React?** → [API_DOCUMENTATION.md - Integración](./API_DOCUMENTATION.md#integración-con-frontend)
- **...manejar errores?** → [API_DOCUMENTATION.md - Códigos de Error](./API_DOCUMENTATION.md#códigos-de-error)

---

## 📝 Notas

- Todos los documentos están en Markdown para fácil lectura
- La especificación OpenAPI está en YAML (compatible con Swagger)
- Los ejemplos están probados y funcionan con el código actual

---

**Última actualización**: 2024-01-15
