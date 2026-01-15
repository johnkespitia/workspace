# Instrucciones para reiniciar el frontend

El contenedor del frontend necesita ser reiniciado para aplicar los cambios.

## Opción 1: Reiniciar solo el frontend (rápido)

Desde el host, ejecuta:

```bash
make dev-restart-frontend
```

## Opción 2: Ver logs del frontend primero

Para ver qué está pasando con el frontend:

```bash
make dev-logs-frontend
```

## Opción 3: Reconstruir completamente

Si el reinicio no funciona, reconstruye el contenedor:

```bash
make dev-rebuild
```

## Verificar que esté funcionando

Después de reiniciar, verifica el estado:

```bash
make dev-health
```

O visita directamente:
- Frontend: http://localhost:3001
- Storybook: http://localhost:6006
