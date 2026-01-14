import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import { fileURLToPath, URL } from 'node:url'

export default defineConfig({
  plugins: [vue()],
  resolve: {
    alias: {
      '@': fileURLToPath(new URL('./src', import.meta.url))
    }
  },
  server: {
    host: '0.0.0.0', // Escuchar en todas las interfaces de red
    port: 3000,
    strictPort: true, // Fallar si el puerto está ocupado (para detectar problemas)
    hmr: {
      host: 'localhost', // HMR desde el host usa localhost
      clientPort: 3000 // Puerto para HMR (Hot Module Replacement)
    },
    watch: {
      usePolling: true // Usar polling para detectar cambios en Docker
    }
  }
})
