import { defineConfig } from "vite";
import vue from "@vitejs/plugin-vue";
import { fileURLToPath, URL } from "node:url";
import { visualizer } from "rollup-plugin-visualizer";

export default defineConfig(({ mode }) => {
  const plugins = [vue()];

  // Agregar bundle analyzer solo en modo analyze
  if (mode === "analyze") {
    plugins.push(
      visualizer({
        open: true,
        filename: "dist/stats.html",
        gzipSize: true,
        brotliSize: true,
      })
    );
  }

  return {
    plugins,
    resolve: {
      alias: {
        "@": fileURLToPath(new URL("./src", import.meta.url)),
      },
    },
    build: {
      // Optimizaciones de build
      rollupOptions: {
        output: {
          manualChunks: {
            // Separar vendor chunks
            "vue-vendor": ["vue", "vue-router", "pinia"],
            "graphql-vendor": ["@urql/core", "@urql/vue", "graphql"],
          },
        },
      },
      // Chunk size warnings
      chunkSizeWarningLimit: 1000,
    },
    server: {
      host: "0.0.0.0", // Escuchar en todas las interfaces de red
      port: 3000,
      strictPort: true, // Fallar si el puerto está ocupado (para detectar problemas)
      hmr: {
        host: "localhost", // HMR desde el host usa localhost
        clientPort: 3000, // Puerto para HMR (Hot Module Replacement)
      },
      watch: {
        usePolling: true, // Usar polling para detectar cambios en Docker
      },
    },
  };
});
