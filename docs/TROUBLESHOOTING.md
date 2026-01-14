# Troubleshooting - Problemas Comunes

## 🔴 Problema: syncStocks no retorna resultados

### Síntomas

- Ejecutas la mutation `syncStocks` pero no ves la respuesta
- Solo ves el query que enviaste, no la respuesta
- La respuesta está vacía o no se muestra

### Posibles Causas y Soluciones

#### 1. API_KEY no configurada o inválida

**Verificar**:

```bash
# En el directorio api/
cat .env.development | grep API_KEY
```

**Solución**:

```bash
# Editar .env.development
API_KEY=tu_api_key_valida_aqui
```

**Verificar en el servidor**: Si falta API_KEY, el servidor mostrará error al iniciar:

```
Failed to load config: API_KEY environment variable is required
```

#### 2. Error en la conversión de datos

**Síntoma**: La API responde pero no se convierten los stocks

**Verificar**: Revisa los logs del servidor cuando ejecutas syncStocks

**Posibles problemas**:

- Formato de datos inesperado de la API
- Ratings inválidos
- Precios negativos o inválidos

**Solución**: Revisa el código de `convertToDomainEntity` en `karenai_api.go`

#### 3. Problema de red/conectividad

**Verificar**:

```bash
# Probar conectividad a la API
curl -H "Authorization: Bearer $API_KEY" https://api.karenai.click/swechallenge/list
```

**Solución**: Verifica tu conexión a internet y que la API esté disponible

#### 4. Base de datos no conectada

**Verificar**: Revisa los logs del servidor:

```
Database connected successfully
```

**Solución**:

- Verifica que CockroachDB esté corriendo
- Verifica las variables de entorno de la base de datos

#### 5. Respuesta GraphQL no se muestra en Playground

**Síntoma**: El playground muestra el query pero no la respuesta

**Posibles causas**:

- Error en el handler GraphQL
- Problema con CORS
- Error silencioso

**Solución**:

1. Abre las herramientas de desarrollador del navegador (F12)
2. Ve a la pestaña "Network"
3. Ejecuta la mutation
4. Busca la request a `/query`
5. Revisa la respuesta en la pestaña "Response"

### Debugging Paso a Paso

#### Paso 1: Verificar que el servidor está corriendo

```bash
# En otra terminal
curl http://localhost:8080/health
# Debe retornar: OK
```

#### Paso 2: Probar el endpoint GraphQL directamente

```bash
curl -X POST http://localhost:8080/query \
  -H "Content-Type: application/json" \
  -d '{
    "query": "mutation { syncStocks { success message stocksSynced } }"
  }'
```

Esto te mostrará la respuesta completa, incluyendo errores.

#### Paso 3: Verificar logs del servidor

Cuando ejecutas `syncStocks`, el servidor debería mostrar errores si los hay. Revisa la terminal donde está corriendo el servidor.

#### Paso 4: Probar la API externa directamente

```bash
# Reemplaza YOUR_API_KEY con tu API key real
curl -H "Authorization: Bearer YOUR_API_KEY" \
  https://api.karenai.click/swechallenge/list
```

Esto te mostrará si la API está respondiendo correctamente.

### Respuesta Esperada

**Éxito**:

```json
{
  "data": {
    "syncStocks": {
      "success": true,
      "message": "Stocks synchronized successfully",
      "stocksSynced": 150
    }
  }
}
```

**Error**:

```json
{
  "data": {
    "syncStocks": {
      "success": false,
      "message": "failed to fetch stocks from API: ...",
      "stocksSynced": 0
    }
  }
}
```

**Error de GraphQL**:

```json
{
  "errors": [
    {
      "message": "error message here",
      "locations": [...],
      "path": [...]
    }
  ],
  "data": null
}
```

### Comandos Útiles

```bash
# Verificar variables de entorno cargadas
cd api
go run ./cmd/main.go
# Revisa los logs al iniciar

# Probar conexión a base de datos
# (si tienes acceso a CockroachDB CLI)
cockroach sql --insecure --host=localhost:26257

# Ver stocks en la base de datos
SELECT COUNT(*) FROM stocks;
SELECT * FROM stocks LIMIT 5;
```

### Logs a Revisar

Cuando ejecutas `syncStocks`, busca estos mensajes en los logs:

- ✅ `Database connected successfully` - Base de datos OK
- ✅ `Database already initialized` - Migraciones OK
- ❌ `Failed to load config` - Problema con variables de entorno
- ❌ `Failed to connect to database` - Problema de conexión a BD
- ❌ `failed to fetch stocks from API` - Problema con API externa
- ❌ `failed to save stocks to database` - Problema al guardar
