import { createApp } from "vue";
import { createPinia } from "pinia";
import App from "./App.vue";
import router from "./router";
import "./style.css";

const app = createApp(App);
const pinia = createPinia();

app.use(pinia);
app.use(router);

// Inicializar el tema antes de montar la app
import { useThemeStore } from "./stores/theme";
const themeStore = useThemeStore();
// Asegurar que el tema se aplique al cargar
if (typeof window !== "undefined") {
  const root = document.documentElement;
  if (themeStore.theme === "dark") {
    root.classList.add("dark");
  } else {
    root.classList.remove("dark");
  }
}

app.mount("#app");
