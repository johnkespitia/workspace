# Diagnóstico del Frontend

Si el frontend no está funcionando, sigue estos pasos:

## 1. Verificar el estado del contenedor

```bash
make dev-status
```

Esto mostrará si el contenedor `frontend` está corriendo.

## 2. Ver los logs del frontend

```bash
make dev-logs-frontend
```

Esto mostrará los logs en tiempo real. Busca errores como:
- Errores de dependencias faltantes
- Errores de puertos ocupados
- Errores de permisos

## 3. Ejecutar diagnóstico completo

```bash
make dev-diagnose-frontend
```

Este comando:
- Verifica el estado del contenedor
- Muestra los últimos logs
- Verifica si los puertos están en uso en el host
- Ejecuta diagnóstico dentro del contenedor

## 4. Verificar salud de los servicios

```bash
make dev-health
```

Esto verifica si los servicios responden a peticiones HTTP.

## 5. Reiniciar el frontend

Si encuentras problemas, intenta reiniciar:

```bash
make dev-restart-frontend
```

## 6. Reconstruir el contenedor

Si el reinicio no funciona, reconstruye:

```bash
make dev-rebuild
```

## Problemas comunes

### El contenedor se cierra inmediatamente

**Síntoma**: El contenedor aparece como "Exited" en `make dev-status`

**Solución**: 
1. Verifica los logs: `make dev-logs-frontend`
2. Verifica que las dependencias estén instaladas
3. Verifica que el script `start-frontend.sh` tenga permisos de ejecución

### Los puertos no están accesibles

**Síntoma**: `ERR_EMPTY_RESPONSE` en el navegador

**Solución**:
1. Verifica que el contenedor esté corriendo: `make dev-status`
2. Verifica que los puertos estén mapeados: `make dev-diagnose-frontend`
3. Verifica que los servicios estén escuchando en `0.0.0.0` (no solo `localhost`)

### Dependencias faltantes

**Síntoma**: Errores sobre módulos no encontrados en los logs

**Solución**:
```bash
make dev-install-frontend
```

## Acceso directo al contenedor

Si necesitas depurar manualmente:

```bash
make dev-shell-frontend
```

Dentro del contenedor puedes:
- Verificar procesos: `ps aux | grep -E "(vite|storybook)"`
- Verificar puertos: `netstat -tlnp | grep -E "(3000|6006)"`
- Ejecutar comandos manualmente: `npm run dev` o `npm run storybook`
