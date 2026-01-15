# Guía de Dev Container

## 📋 Resumen

Esta guía explica cómo usar y configurar el Dev Container del proyecto.

---

## 🚀 Inicio Rápido

### Inicialización Automática

```bash
# Inicializar y abrir IDE
make dev-init
```

Este comando:
1. Construye las imágenes Docker
2. Inicia todos los servicios
3. Abre automáticamente el IDE con el devcontainer

### Verificar Instalación

```bash
# Ver estado de servicios
make dev-status

# Verificar salud
make dev-health
```

---

## 📁 Estructura de Archivos

```
.devcontainer/
├── devcontainer.json          # Configuración principal
├── docker-compose.yml         # Orquestación de servicios
├── Dockerfile.api             # Imagen del backend
├── Dockerfile.frontend        # Imagen del frontend
├── postCreate.sh              # Script post-creación
├── start-api-auto.sh          # Script de inicio del API
├── start-frontend.sh          # Script de inicio del frontend
└── diagnose-frontend.sh       # Script de diagnóstico
```

---

## ⚙️ Configuración

### devcontainer.json

**Servicios principales**:
- `api`: Contenedor principal (backend Go)
- `frontend`: Servidor de desarrollo Vue 3
- `cockroachdb`: Base de datos

**Puertos forwardeados**:
- `3001`: Frontend
- `6006`: Storybook
- `8080`: Backend API
- `26257`: CockroachDB SQL

**Extensiones preinstaladas**:
- Go (golang.Go)
- Vue Language Features (Vue.volar)
- TypeScript
- ESLint
- Prettier
- Tailwind CSS IntelliSense

### docker-compose.yml

Define tres servicios:

1. **api**: Backend Go con hot reload
2. **frontend**: Frontend Vue 3 con Vite y Storybook
3. **cockroachdb**: Base de datos CockroachDB

---

## 🔄 Scripts de Inicio

### postCreate.sh

Se ejecuta automáticamente después de crear el contenedor:

**Funcionalidades**:
- Instala herramientas de Go (gopls, delve, air)
- Instala dependencias del frontend (`npm install`)
- Descarga módulos de Go (`go mod download`)
- Configura permisos de Docker
- Espera a que CockroachDB esté listo

### start-api-auto.sh

Inicia el servidor API automáticamente:

**Funcionalidades**:
- Verifica que CockroachDB esté listo
- Inicia `air` para hot reload
- Redirige logs a `/tmp/api.log`

### start-frontend.sh

Inicia Vite y Storybook en paralelo:

**Funcionalidades**:
- Verifica dependencias instaladas
- Instala `concurrently` si es necesario
- Inicia `npm run dev` y `npm run storybook` en paralelo
- Maneja señales correctamente

---

## 🛠️ Comandos Útiles

### Gestión de Servicios

```bash
# Iniciar servicios
make dev-up

# Detener servicios
make dev-down

# Reiniciar servicios
make dev-restart

# Reconstruir imágenes
make dev-rebuild
```

### Logs

```bash
# Todos los servicios
make dev-logs

# Servicio específico
make dev-logs-api
make dev-logs-frontend
make dev-logs-db
```

### Shell

```bash
# Shell en contenedor API
make dev-shell

# Shell en contenedor Frontend
make dev-shell-frontend
```

### Diagnóstico

```bash
# Diagnóstico completo
make dev-diagnose

# Diagnóstico del frontend
make dev-diagnose-frontend

# Estado de salud
make dev-health
```

---

## 🔧 Configuración Avanzada

### Variables de Entorno

Las variables de entorno se configuran en `devcontainer.json`:

```json
"remoteEnv": {
  "DATABASE_URL": "postgresql://root@cockroachdb:26257/defaultdb?sslmode=disable",
  "COCKROACH_HOST": "cockroachdb",
  "COCKROACH_PORT": "26257",
  "COCKROACH_USER": "root",
  "COCKROACH_DB": "defaultdb"
}
```

### Personalizar Extensiones

Editar `.devcontainer/devcontainer.json`:

```json
"extensions": [
  "golang.Go",
  "Vue.volar",
  // Agregar más extensiones aquí
]
```

### Personalizar Configuración de VS Code

Editar `.devcontainer/devcontainer.json`:

```json
"settings": {
  "editor.formatOnSave": true,
  "editor.defaultFormatter": "esbenp.prettier-vscode",
  // Agregar más configuraciones aquí
}
```

---

## 🐛 Solución de Problemas

### Servicios no inician

```bash
# Reconstruir contenedor
make dev-rebuild

# Ver logs
make dev-logs
```

### Puerto en uso

```bash
# Verificar qué proceso usa el puerto
lsof -i :3001
lsof -i :8080

# Detener proceso
kill -9 <PID>
```

### Permisos de Docker

Si tienes problemas con permisos de Docker dentro del contenedor:

```bash
# Verificar grupo docker
groups

# El script postCreate.sh debería configurar esto automáticamente
```

### Frontend no accesible

```bash
# Diagnóstico del frontend
make dev-diagnose-frontend

# Verificar puertos
make dev-status
```

---

## 📚 Recursos

- [Dev Containers Documentation](https://code.visualstudio.com/docs/devcontainers/containers)
- [Docker Compose Documentation](https://docs.docker.com/compose/)
- [Makefile Documentation](./MAKEFILE.md)

---

**Última actualización**: 2026-01-15
