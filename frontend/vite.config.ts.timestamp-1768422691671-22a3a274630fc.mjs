// vite.config.ts
import { defineConfig } from "file:///workspace/frontend/node_modules/vite/dist/node/index.js";
import vue from "file:///workspace/frontend/node_modules/@vitejs/plugin-vue/dist/index.mjs";
import { fileURLToPath, URL } from "node:url";
import { visualizer } from "file:///workspace/frontend/node_modules/rollup-plugin-visualizer/dist/plugin/index.js";
var __vite_injected_original_import_meta_url = "file:///workspace/frontend/vite.config.ts";
var vite_config_default = defineConfig(({ mode }) => {
  const plugins = [vue()];
  if (mode === "analyze") {
    plugins.push(
      visualizer({
        open: true,
        filename: "dist/stats.html",
        gzipSize: true,
        brotliSize: true
      })
    );
  }
  return {
    plugins,
    resolve: {
      alias: {
        "@": fileURLToPath(new URL("./src", __vite_injected_original_import_meta_url))
      }
    },
    build: {
      // Optimizaciones de build
      rollupOptions: {
        output: {
          manualChunks: {
            // Separar vendor chunks
            "vue-vendor": ["vue", "vue-router", "pinia"],
            "graphql-vendor": ["@urql/core", "@urql/vue", "graphql"]
          }
        }
      },
      // Chunk size warnings
      chunkSizeWarningLimit: 1e3
    },
    server: {
      host: "0.0.0.0",
      // Escuchar en todas las interfaces de red
      port: 3e3,
      strictPort: true,
      // Fallar si el puerto está ocupado (para detectar problemas)
      hmr: {
        host: "localhost",
        // HMR desde el host usa localhost
        clientPort: 3e3
        // Puerto para HMR (Hot Module Replacement)
      },
      watch: {
        usePolling: true
        // Usar polling para detectar cambios en Docker
      }
    }
  };
});
export {
  vite_config_default as default
};
//# sourceMappingURL=data:application/json;base64,ewogICJ2ZXJzaW9uIjogMywKICAic291cmNlcyI6IFsidml0ZS5jb25maWcudHMiXSwKICAic291cmNlc0NvbnRlbnQiOiBbImNvbnN0IF9fdml0ZV9pbmplY3RlZF9vcmlnaW5hbF9kaXJuYW1lID0gXCIvd29ya3NwYWNlL2Zyb250ZW5kXCI7Y29uc3QgX192aXRlX2luamVjdGVkX29yaWdpbmFsX2ZpbGVuYW1lID0gXCIvd29ya3NwYWNlL2Zyb250ZW5kL3ZpdGUuY29uZmlnLnRzXCI7Y29uc3QgX192aXRlX2luamVjdGVkX29yaWdpbmFsX2ltcG9ydF9tZXRhX3VybCA9IFwiZmlsZTovLy93b3Jrc3BhY2UvZnJvbnRlbmQvdml0ZS5jb25maWcudHNcIjtpbXBvcnQgeyBkZWZpbmVDb25maWcgfSBmcm9tIFwidml0ZVwiO1xuaW1wb3J0IHZ1ZSBmcm9tIFwiQHZpdGVqcy9wbHVnaW4tdnVlXCI7XG5pbXBvcnQgeyBmaWxlVVJMVG9QYXRoLCBVUkwgfSBmcm9tIFwibm9kZTp1cmxcIjtcbmltcG9ydCB7IHZpc3VhbGl6ZXIgfSBmcm9tIFwicm9sbHVwLXBsdWdpbi12aXN1YWxpemVyXCI7XG5cbmV4cG9ydCBkZWZhdWx0IGRlZmluZUNvbmZpZygoeyBtb2RlIH0pID0+IHtcbiAgY29uc3QgcGx1Z2lucyA9IFt2dWUoKV07XG5cbiAgLy8gQWdyZWdhciBidW5kbGUgYW5hbHl6ZXIgc29sbyBlbiBtb2RvIGFuYWx5emVcbiAgaWYgKG1vZGUgPT09IFwiYW5hbHl6ZVwiKSB7XG4gICAgcGx1Z2lucy5wdXNoKFxuICAgICAgdmlzdWFsaXplcih7XG4gICAgICAgIG9wZW46IHRydWUsXG4gICAgICAgIGZpbGVuYW1lOiBcImRpc3Qvc3RhdHMuaHRtbFwiLFxuICAgICAgICBnemlwU2l6ZTogdHJ1ZSxcbiAgICAgICAgYnJvdGxpU2l6ZTogdHJ1ZSxcbiAgICAgIH0pXG4gICAgKTtcbiAgfVxuXG4gIHJldHVybiB7XG4gICAgcGx1Z2lucyxcbiAgICByZXNvbHZlOiB7XG4gICAgICBhbGlhczoge1xuICAgICAgICBcIkBcIjogZmlsZVVSTFRvUGF0aChuZXcgVVJMKFwiLi9zcmNcIiwgaW1wb3J0Lm1ldGEudXJsKSksXG4gICAgICB9LFxuICAgIH0sXG4gICAgYnVpbGQ6IHtcbiAgICAgIC8vIE9wdGltaXphY2lvbmVzIGRlIGJ1aWxkXG4gICAgICByb2xsdXBPcHRpb25zOiB7XG4gICAgICAgIG91dHB1dDoge1xuICAgICAgICAgIG1hbnVhbENodW5rczoge1xuICAgICAgICAgICAgLy8gU2VwYXJhciB2ZW5kb3IgY2h1bmtzXG4gICAgICAgICAgICBcInZ1ZS12ZW5kb3JcIjogW1widnVlXCIsIFwidnVlLXJvdXRlclwiLCBcInBpbmlhXCJdLFxuICAgICAgICAgICAgXCJncmFwaHFsLXZlbmRvclwiOiBbXCJAdXJxbC9jb3JlXCIsIFwiQHVycWwvdnVlXCIsIFwiZ3JhcGhxbFwiXSxcbiAgICAgICAgICB9LFxuICAgICAgICB9LFxuICAgICAgfSxcbiAgICAgIC8vIENodW5rIHNpemUgd2FybmluZ3NcbiAgICAgIGNodW5rU2l6ZVdhcm5pbmdMaW1pdDogMTAwMCxcbiAgICB9LFxuICAgIHNlcnZlcjoge1xuICAgICAgaG9zdDogXCIwLjAuMC4wXCIsIC8vIEVzY3VjaGFyIGVuIHRvZGFzIGxhcyBpbnRlcmZhY2VzIGRlIHJlZFxuICAgICAgcG9ydDogMzAwMCxcbiAgICAgIHN0cmljdFBvcnQ6IHRydWUsIC8vIEZhbGxhciBzaSBlbCBwdWVydG8gZXN0XHUwMEUxIG9jdXBhZG8gKHBhcmEgZGV0ZWN0YXIgcHJvYmxlbWFzKVxuICAgICAgaG1yOiB7XG4gICAgICAgIGhvc3Q6IFwibG9jYWxob3N0XCIsIC8vIEhNUiBkZXNkZSBlbCBob3N0IHVzYSBsb2NhbGhvc3RcbiAgICAgICAgY2xpZW50UG9ydDogMzAwMCwgLy8gUHVlcnRvIHBhcmEgSE1SIChIb3QgTW9kdWxlIFJlcGxhY2VtZW50KVxuICAgICAgfSxcbiAgICAgIHdhdGNoOiB7XG4gICAgICAgIHVzZVBvbGxpbmc6IHRydWUsIC8vIFVzYXIgcG9sbGluZyBwYXJhIGRldGVjdGFyIGNhbWJpb3MgZW4gRG9ja2VyXG4gICAgICB9LFxuICAgIH0sXG4gIH07XG59KTtcbiJdLAogICJtYXBwaW5ncyI6ICI7QUFBMk8sU0FBUyxvQkFBb0I7QUFDeFEsT0FBTyxTQUFTO0FBQ2hCLFNBQVMsZUFBZSxXQUFXO0FBQ25DLFNBQVMsa0JBQWtCO0FBSG1ILElBQU0sMkNBQTJDO0FBSy9MLElBQU8sc0JBQVEsYUFBYSxDQUFDLEVBQUUsS0FBSyxNQUFNO0FBQ3hDLFFBQU0sVUFBVSxDQUFDLElBQUksQ0FBQztBQUd0QixNQUFJLFNBQVMsV0FBVztBQUN0QixZQUFRO0FBQUEsTUFDTixXQUFXO0FBQUEsUUFDVCxNQUFNO0FBQUEsUUFDTixVQUFVO0FBQUEsUUFDVixVQUFVO0FBQUEsUUFDVixZQUFZO0FBQUEsTUFDZCxDQUFDO0FBQUEsSUFDSDtBQUFBLEVBQ0Y7QUFFQSxTQUFPO0FBQUEsSUFDTDtBQUFBLElBQ0EsU0FBUztBQUFBLE1BQ1AsT0FBTztBQUFBLFFBQ0wsS0FBSyxjQUFjLElBQUksSUFBSSxTQUFTLHdDQUFlLENBQUM7QUFBQSxNQUN0RDtBQUFBLElBQ0Y7QUFBQSxJQUNBLE9BQU87QUFBQTtBQUFBLE1BRUwsZUFBZTtBQUFBLFFBQ2IsUUFBUTtBQUFBLFVBQ04sY0FBYztBQUFBO0FBQUEsWUFFWixjQUFjLENBQUMsT0FBTyxjQUFjLE9BQU87QUFBQSxZQUMzQyxrQkFBa0IsQ0FBQyxjQUFjLGFBQWEsU0FBUztBQUFBLFVBQ3pEO0FBQUEsUUFDRjtBQUFBLE1BQ0Y7QUFBQTtBQUFBLE1BRUEsdUJBQXVCO0FBQUEsSUFDekI7QUFBQSxJQUNBLFFBQVE7QUFBQSxNQUNOLE1BQU07QUFBQTtBQUFBLE1BQ04sTUFBTTtBQUFBLE1BQ04sWUFBWTtBQUFBO0FBQUEsTUFDWixLQUFLO0FBQUEsUUFDSCxNQUFNO0FBQUE7QUFBQSxRQUNOLFlBQVk7QUFBQTtBQUFBLE1BQ2Q7QUFBQSxNQUNBLE9BQU87QUFBQSxRQUNMLFlBQVk7QUFBQTtBQUFBLE1BQ2Q7QUFBQSxJQUNGO0FBQUEsRUFDRjtBQUNGLENBQUM7IiwKICAibmFtZXMiOiBbXQp9Cg==
