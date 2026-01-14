# URLs de Documentación de la API

## 🌐 URLs Disponibles

Cuando el servidor está corriendo en `http://localhost:8080`, las siguientes URLs están disponibles:

### 📚 Documentación Principal

| URL                                       | Descripción                                                                                   |
| ----------------------------------------- | --------------------------------------------------------------------------------------------- |
| `http://localhost:8080/docs`              | **Página principal de documentación** - Punto de entrada con enlaces a todas las herramientas |
| `http://localhost:8080/docs/swagger`      | **Swagger UI** - Interfaz interactiva para explorar la API REST                               |
| `http://localhost:8080/docs/openapi.yaml` | **Especificación OpenAPI** - Archivo YAML con la especificación completa                      |

### 🎮 Herramientas Interactivas

| URL                                | Descripción                                                                      |
| ---------------------------------- | -------------------------------------------------------------------------------- |
| `http://localhost:8080/playground` | **GraphQL Playground** - Interfaz visual para probar queries y mutations GraphQL |

### 🔗 Endpoints de la API

| URL                            | Método | Descripción                |
| ------------------------------ | ------ | -------------------------- |
| `http://localhost:8080/health` | GET    | Health check del servidor  |
| `http://localhost:8080/query`  | POST   | Endpoint principal GraphQL |

---

## 🚀 Inicio Rápido

1. **Iniciar el servidor**:

   ```bash
   cd api
   go run ./cmd/main.go
   ```

2. **Abrir documentación**:
   - Navega a: `http://localhost:8080/docs`
   - O directamente a: `http://localhost:8080/docs/swagger` para Swagger UI
   - O a: `http://localhost:8080/playground` para GraphQL Playground

---

## 📖 Documentación por Tipo

### Para Usuarios Finales

- **Guía de Usuario**: Ver archivo `USER_GUIDE.md` (no disponible vía web, solo archivo)
- **GraphQL Playground**: `http://localhost:8080/playground` (interfaz visual)

### Para Desarrolladores

- **Swagger UI**: `http://localhost:8080/docs/swagger` (documentación REST interactiva)
- **OpenAPI Spec**: `http://localhost:8080/docs/openapi.yaml` (especificación técnica)
- **Documentación Completa**: Ver archivo `API_DOCUMENTATION.md` (no disponible vía web, solo archivo)

### Para Integración

- **GraphQL Endpoint**: `http://localhost:8080/query` (endpoint principal)
- **GraphQL Playground**: `http://localhost:8080/playground` (para probar queries)

---

## 💡 Recomendaciones

1. **Primera vez**: Empieza en `http://localhost:8080/docs` para ver todas las opciones
2. **Probar GraphQL**: Usa `http://localhost:8080/playground` para explorar el schema
3. **Ver API REST**: Usa `http://localhost:8080/docs/swagger` para documentación interactiva
4. **Integración**: Consulta los archivos Markdown en `api/docs/` para documentación detallada

---

**Nota**: Si el servidor está corriendo en un puerto diferente, reemplaza `8080` con el puerto configurado en la variable de entorno `PORT`.
