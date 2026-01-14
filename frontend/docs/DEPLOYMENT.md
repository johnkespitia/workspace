# Guía de Deployment

Esta guía explica cómo hacer deploy del frontend en diferentes entornos.

## 📋 Prerrequisitos

- Node.js 18+ y npm
- Acceso al servidor/plataforma de deployment
- Variables de entorno configuradas

## 🔧 Variables de Entorno

Crear archivo `.env.production`:

```env
# Endpoint de GraphQL API
VITE_GRAPHQL_ENDPOINT=https://api.tudominio.com/query

# Otros variables si es necesario
VITE_APP_TITLE=Stock Info
```

### Variables Disponibles

- `VITE_GRAPHQL_ENDPOINT`: URL del endpoint GraphQL (requerido)
- `VITE_APP_TITLE`: Título de la aplicación (opcional)

## 🏗️ Build para Producción

### Build Local

```bash
# Instalar dependencias
npm install

# Build para producción
npm run build

# El build estará en la carpeta dist/
```

### Verificar Build

```bash
# Preview del build local
npm run preview
```

## 🐳 Docker

### Dockerfile

```dockerfile
# Build stage
FROM node:18-alpine AS builder

WORKDIR /app

# Copiar archivos de dependencias
COPY package*.json ./
RUN npm ci

# Copiar código fuente
COPY . .

# Build
RUN npm run build

# Production stage
FROM nginx:alpine

# Copiar build
COPY --from=builder /app/dist /usr/share/nginx/html

# Copiar configuración de nginx
COPY nginx.conf /etc/nginx/conf.d/default.conf

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]
```

### docker-compose.yml

```yaml
version: "3.8"

services:
  frontend:
    build:
      context: .
      dockerfile: Dockerfile
    ports:
      - "3000:80"
    environment:
      - VITE_GRAPHQL_ENDPOINT=http://api:8080/query
    depends_on:
      - api
```

### Build y Run

```bash
# Build imagen
docker build -t stock-frontend .

# Run contenedor
docker run -p 3000:80 \
  -e VITE_GRAPHQL_ENDPOINT=http://api:8080/query \
  stock-frontend
```

## 🌐 Nginx

### Configuración Básica

Crear `nginx.conf`:

```nginx
server {
    listen 80;
    server_name _;
    root /usr/share/nginx/html;
    index index.html;

    # Gzip compression
    gzip on;
    gzip_vary on;
    gzip_min_length 1024;
    gzip_types text/plain text/css text/xml text/javascript
               application/x-javascript application/xml+rss
               application/javascript application/json;

    # Cache estáticos
    location ~* \.(jpg|jpeg|png|gif|ico|css|js|svg|woff|woff2)$ {
        expires 1y;
        add_header Cache-Control "public, immutable";
    }

    # SPA routing
    location / {
        try_files $uri $uri/ /index.html;
    }

    # Security headers
    add_header X-Frame-Options "SAMEORIGIN" always;
    add_header X-Content-Type-Options "nosniff" always;
    add_header X-XSS-Protection "1; mode=block" always;
}
```

## ☁️ Plataformas Cloud

### Vercel

1. Conectar repositorio
2. Configurar variables de entorno
3. Build command: `npm run build`
4. Output directory: `dist`

### Netlify

1. Conectar repositorio
2. Configurar:
   - Build command: `npm run build`
   - Publish directory: `dist`
3. Agregar variables de entorno

### AWS S3 + CloudFront

```bash
# Build
npm run build

# Subir a S3
aws s3 sync dist/ s3://tu-bucket-name --delete

# Invalidar CloudFront
aws cloudfront create-invalidation \
  --distribution-id E1234567890 \
  --paths "/*"
```

### Google Cloud Platform

```bash
# Build
npm run build

# Deploy a Cloud Storage
gsutil -m rsync -r -d dist/ gs://tu-bucket-name

# O usar App Engine
gcloud app deploy app.yaml
```

## 🔒 Seguridad

### Headers de Seguridad

Agregar en nginx o servidor:

```
X-Frame-Options: SAMEORIGIN
X-Content-Type-Options: nosniff
X-XSS-Protection: 1; mode=block
Content-Security-Policy: default-src 'self'
```

### HTTPS

Siempre usar HTTPS en producción:

```nginx
server {
    listen 443 ssl http2;
    ssl_certificate /path/to/cert.pem;
    ssl_certificate_key /path/to/key.pem;
    # ...
}
```

## 📊 Monitoreo

### Performance

- Lighthouse CI para métricas automáticas
- Google Analytics para tracking
- Sentry para error tracking

### Health Check

Endpoint simple para verificar que la app está funcionando:

```typescript
// En el servidor, servir /health
// Retorna 200 OK si la app está funcionando
```

## 🔄 CI/CD

### GitHub Actions

Crear `.github/workflows/deploy.yml`:

```yaml
name: Deploy

on:
  push:
    branches: [main]

jobs:
  build-and-deploy:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-node@v3
        with:
          node-version: "18"
      - run: npm ci
      - run: npm run build
      - run: npm run test:run
      - name: Deploy
        # Agregar paso de deployment
```

## 🐛 Troubleshooting

### Build Falla

1. Verificar Node.js version: `node --version` (debe ser 18+)
2. Limpiar cache: `rm -rf node_modules package-lock.json && npm install`
3. Verificar variables de entorno

### Errores de CORS

Configurar CORS en el backend para permitir el dominio del frontend.

### Assets No Cargan

Verificar que la ruta base esté correcta en `vite.config.ts` si se deploya en subdirectorio.

### GraphQL No Conecta

1. Verificar `VITE_GRAPHQL_ENDPOINT`
2. Verificar que el backend esté corriendo
3. Verificar CORS

## 📈 Optimizaciones Post-Deployment

### CDN

Usar CDN para servir assets estáticos:

- CloudFlare
- AWS CloudFront
- Google Cloud CDN

### Caching

- Cache de assets estáticos (1 año)
- Cache de HTML (no cache o corto)
- Service Worker para cache offline (futuro)

### Compression

- Gzip o Brotli en servidor
- Ya configurado en nginx.conf de ejemplo

## 📝 Checklist Pre-Deployment

- [ ] Variables de entorno configuradas
- [ ] Build exitoso sin errores
- [ ] Tests pasando
- [ ] Linter sin errores
- [ ] Bundle size verificado
- [ ] HTTPS configurado
- [ ] Headers de seguridad configurados
- [ ] CORS configurado en backend
- [ ] Monitoreo configurado
- [ ] Documentación actualizada

## 🚀 Deployment Rápido

### Opción 1: Vercel (Más Fácil)

```bash
npm install -g vercel
vercel
```

### Opción 2: Netlify

```bash
npm install -g netlify-cli
netlify deploy --prod
```

### Opción 3: Docker

```bash
docker build -t frontend .
docker run -p 3000:80 frontend
```

---

**¿Problemas?** Revisa la sección de Troubleshooting o crea un issue.
