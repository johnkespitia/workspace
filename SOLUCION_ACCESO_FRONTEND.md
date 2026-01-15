# Solución: Frontend no accesible desde el host

## Diagnóstico

El diagnóstico muestra que:

- ✅ Contenedor está corriendo
- ✅ Procesos Vite y Storybook están activos
- ✅ Puertos están escuchando dentro del contenedor (3000 y 6006)
- ✅ Puertos están mapeados (3001:3000 y 6006:6006)
- ✅ Puertos están en uso en el host

Pero el navegador muestra `ERR_EMPTY_RESPONSE`.

## Posibles causas y soluciones

### 1. Docker Desktop en Mac - Problema de mapeo de puertos

**Síntoma**: Los puertos están mapeados pero no son accesibles desde el navegador.

**Solución**: He actualizado `docker-compose.yml` para usar `127.0.0.1:PUERTO` en lugar de solo `PUERTO`. Esto fuerza el mapeo a localhost.

**Aplicar cambios**:

```bash
make dev-restart-frontend
```

### 2. Verificar conectividad desde el host

Ejecuta estos comandos desde el host (no desde el contenedor):

```bash
# Verificar Frontend
curl -v http://localhost:3001

# Verificar Storybook
curl -v http://localhost:6006
```

Si `curl` funciona pero el navegador no, puede ser un problema de CORS o de caché del navegador.

### 3. Limpiar caché del navegador

1. Abre las herramientas de desarrollador (F12)
2. Haz clic derecho en el botón de recargar
3. Selecciona "Vaciar caché y recargar de forma forzada"

### 4. Verificar Docker Desktop

En Docker Desktop:

1. Ve a Settings > Resources > Network
2. Verifica que los puertos no estén bloqueados
3. Reinicia Docker Desktop si es necesario

### 5. Verificar firewall

En Mac:

```bash
# Verificar si hay reglas de firewall bloqueando
sudo /usr/libexec/ApplicationFirewall/socketfilterfw --listapps
```

### 6. Probar con IP directa del contenedor

Si el mapeo de puertos no funciona, puedes intentar acceder directamente a la IP del contenedor:

```bash
# Obtener la IP del contenedor
docker inspect go-react-test-frontend | grep IPAddress

# Luego acceder desde el navegador usando esa IP
# (Esto solo funciona si estás en la misma red)
```

### 7. Reiniciar completamente

Si nada funciona:

```bash
# Detener todo
make dev-down

# Reconstruir
make dev-rebuild

# Verificar
make dev-health
```

## Verificación final

Después de aplicar los cambios:

```bash
# Verificar estado
make dev-status

# Verificar salud
make dev-health

# Ver logs
make dev-logs-frontend
```

## Acceso directo

Si el mapeo de puertos sigue sin funcionar, puedes acceder directamente al contenedor:

```bash
# Abrir shell en el contenedor frontend
make dev-shell-frontend

# Dentro del contenedor, verificar que los servicios estén escuchando
netstat -tlnp | grep -E "(3000|6006)"
```

Los servicios deberían estar escuchando en `0.0.0.0:3000` y `0.0.0.0:6006` dentro del contenedor.
