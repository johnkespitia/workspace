# Configuración de Documentación Web

## 📍 Ubicación de Archivos

Los archivos de documentación deben estar en las siguientes ubicaciones:

- `api/docs/API_DOCUMENTATION.md` - Documentación completa de la API
- `api/docs/USER_GUIDE.md` - Guía de usuario
- `docs/GRAPHQL_API_REFERENCE.md` - Referencia GraphQL (en workspace root)
- `api/docs/openapi.yaml` - Especificación OpenAPI

## 🔗 URLs Disponibles

Cuando el servidor está corriendo, puedes acceder a:

| URL                                            | Descripción                            |
| ---------------------------------------------- | -------------------------------------- |
| `http://localhost:8080/docs`                   | Página principal de documentación      |
| `http://localhost:8080/docs/api`               | Documentación de API (Markdown → HTML) |
| `http://localhost:8080/docs/guide`             | Guía de Usuario (Markdown → HTML)      |
| `http://localhost:8080/docs/graphql-reference` | Referencia GraphQL (Markdown → HTML)   |
| `http://localhost:8080/docs/swagger`           | Swagger UI                             |
| `http://localhost:8080/docs/openapi.yaml`      | Especificación OpenAPI                 |
| `http://localhost:8080/playground`             | GraphQL Playground                     |

## 🔧 Cómo Funciona

1. **MarkdownDocHandler**: Convierte archivos Markdown a HTML usando `blackfriday`
2. **Búsqueda de archivos**: Busca en múltiples ubicaciones para encontrar los archivos
3. **Renderizado**: Crea una página HTML completa con estilos CSS

## 🐛 Troubleshooting

### Error 404 en enlaces de documentación

**Causa**: Los archivos Markdown no se encuentran en las rutas esperadas.

**Solución**:

1. Verifica que los archivos existan:

   ```bash
   ls -la api/docs/API_DOCUMENTATION.md
   ls -la api/docs/USER_GUIDE.md
   ls -la docs/GRAPHQL_API_REFERENCE.md
   ```

2. Si los archivos están en otra ubicación, el handler los buscará automáticamente en:

   - `./docs/`
   - `../docs/`
   - `../../docs/`
   - `./api/docs/`
   - `../api/docs/`
   - Y desde el directorio de trabajo actual

3. Verifica los logs del servidor para ver qué rutas está intentando

### Archivo no encontrado

Si ves "Documentation file not found", verifica:

- Que el archivo existe
- Que tiene permisos de lectura
- Que el nombre del archivo coincide exactamente (case-sensitive)

## 📝 Notas

- Los archivos Markdown se convierten a HTML en tiempo real
- Los estilos CSS están embebidos en el HTML generado
- Los enlaces internos en Markdown funcionan correctamente
- El código se resalta con estilos básicos
