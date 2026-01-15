# Guía del Makefile

## 📋 Resumen

Este documento describe todos los comandos disponibles en el `Makefile` del proyecto.

---

## 🚀 Comandos Principales

### `make help`

Muestra todos los comandos disponibles con sus descripciones.

```bash
make help
```

### `make dev-init`

Inicializa el dev container (construye e inicia servicios) y abre el IDE automáticamente.

```bash
make dev-init
```

**Funcionalidades**:

- Construye las imágenes Docker
- Inicia todos los servicios
- Abre el IDE (Cursor o VS Code) con el devcontainer

**Configurar IDE**:

```bash
# Usar Cursor
export IDE=cursor
make dev-init

# Usar VS Code
export IDE=code
make dev-init

# Auto-detectar
export IDE=auto
make dev-init
```

---

## 🔄 Gestión de Servicios

### `make dev-up`

Inicia los servicios del dev container.

```bash
make dev-up
```

**Servicios iniciados**:

- Frontend: http://localhost:3001
- Storybook: http://localhost:6006
- Backend: http://localhost:8080
- CockroachDB UI: http://localhost:8081

### `make dev-down`

Detiene los servicios del dev container.

```bash
make dev-down
```

**Alias**: `make dev-stop`

### `make dev-restart`

Reinicia todos los servicios.

```bash
make dev-restart
```

### `make dev-restart-api`

Reinicia solo el servicio API.

```bash
make dev-restart-api
```

### `make dev-restart-frontend`

Reinicia solo el servicio Frontend.

```bash
make dev-restart-frontend
```

### `make dev-rebuild`

Reconstruye las imágenes Docker sin caché y reinicia los servicios.

```bash
make dev-rebuild
```

**⚠️ Nota**: Esto puede tardar varios minutos.

---

## 📊 Monitoreo

### `make dev-status`

Muestra el estado de los servicios.

```bash
make dev-status
```

**Información mostrada**:

- Estado de contenedores
- Puertos expuestos
- URLs de acceso

### `make dev-health`

Verifica el estado de salud de los servicios.

```bash
make dev-health
```

**Verificaciones**:

- Backend responde en `/health`
- Frontend accesible
- Storybook accesible
- CockroachDB responde

### `make dev-logs`

Muestra los logs de todos los servicios en tiempo real.

```bash
make dev-logs
```

**Salir**: `Ctrl+C`

### `make dev-logs-api`

Muestra los logs solo del servicio API.

```bash
make dev-logs-api
```

### `make dev-logs-frontend`

Muestra los logs solo del servicio Frontend.

```bash
make dev-logs-frontend
```

### `make dev-logs-db`

Muestra los logs solo de CockroachDB.

```bash
make dev-logs-db
```

---

## 🔍 Diagnóstico

### `make dev-diagnose`

Diagnóstico detallado de los servicios.

```bash
make dev-diagnose
```

**Información mostrada**:

- Estado de contenedores
- Últimos logs del frontend
- Puertos en uso
- Conectividad HTTP

### `make dev-diagnose-frontend`

Diagnóstico específico del frontend.

```bash
make dev-diagnose-frontend
```

**Información mostrada**:

- Estado del contenedor
- Últimos logs
- Puertos en uso
- Diagnóstico dentro del contenedor
- Conectividad desde el host
- Información de red

---

## 🛠️ Utilidades

### `make dev-shell`

Abre una shell bash en el contenedor API.

```bash
make dev-shell
```

**Uso**:

```bash
# Dentro del shell
cd /workspace/api
go run cmd/main.go
```

### `make dev-shell-frontend`

Abre una shell bash en el contenedor Frontend.

```bash
make dev-shell-frontend
```

**Uso**:

```bash
# Dentro del shell
cd /workspace/frontend
npm run dev
```

### `make dev-open`

Abre el IDE (Cursor o VS Code) con el devcontainer.

```bash
make dev-open
```

**Nota**: Requiere que el devcontainer ya esté iniciado.

### `make dev-install-frontend`

Instala las dependencias del frontend.

```bash
make dev-install-frontend
```

**Útil cuando**:

- Se agregaron nuevas dependencias
- `node_modules` está corrupto
- Reinstalación necesaria

### `make dev-clean`

Detiene servicios y elimina volúmenes.

```bash
make dev-clean
```

**⚠️ ADVERTENCIA**: Esto elimina todos los datos de CockroachDB.

**Útil cuando**:

- Necesitas empezar desde cero
- Hay problemas con volúmenes
- Limpieza completa necesaria

---

## 📝 Ejemplos de Uso

### Flujo de Desarrollo Típico

```bash
# 1. Inicializar proyecto
make dev-init

# 2. Verificar que todo está corriendo
make dev-status
make dev-health

# 3. Ver logs si hay problemas
make dev-logs-api

# 4. Reiniciar servicio si es necesario
make dev-restart-api

# 5. Al terminar, detener servicios
make dev-down
```

### Debugging

```bash
# Ver logs en tiempo real
make dev-logs

# Diagnóstico completo
make dev-diagnose

# Diagnóstico específico del frontend
make dev-diagnose-frontend

# Abrir shell para debugging
make dev-shell
```

### Reinstalación

```bash
# Reconstruir todo desde cero
make dev-clean
make dev-rebuild

# Solo reinstalar dependencias del frontend
make dev-install-frontend
```

---

## 🔧 Variables del Makefile

El Makefile detecta automáticamente:

- **Docker Compose**: Detecta `docker-compose` (V1) o `docker compose` (V2)
- **IDE**: Usa variable de entorno `IDE` o pregunta interactivamente

**Variables configurables**:

```makefile
PROJECT_NAME := go-react-test-devcontainer
DEV_CONTAINER_DIR := .devcontainer
```

---

## 🎨 Colores en Output

El Makefile usa colores para mejor legibilidad:

- 🟢 **Verde**: Éxito
- 🟡 **Amarillo**: Advertencias/Información
- 🔴 **Rojo**: Errores
- 🔵 **Azul**: Acciones importantes

---

## 📚 Recursos

- [Makefile Documentation](https://www.gnu.org/software/make/manual/)
- [Docker Compose Documentation](https://docs.docker.com/compose/)
- [Dev Container Documentation](./DEVCONTAINER.md)

---

**Última actualización**: 2026-01-15
