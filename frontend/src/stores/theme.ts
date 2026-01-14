import { defineStore } from "pinia";
import { ref, watch, nextTick } from "vue";

export type Theme = "light" | "dark";

export const useThemeStore = defineStore("theme", () => {
  // Obtener tema del localStorage o usar 'light' por defecto
  const getStoredTheme = (): Theme => {
    if (typeof window === "undefined") return "light";
    const stored = localStorage.getItem("theme");
    return (stored === "dark" || stored === "light" ? stored : "light") as Theme;
  };

  const theme = ref<Theme>(getStoredTheme());

  const applyTheme = (themeToApply: Theme) => {
    if (typeof window === "undefined") return;
    
    const root = document.documentElement;
    if (themeToApply === "dark") {
      root.classList.add("dark");
    } else {
      root.classList.remove("dark");
    }
  };

  const setTheme = (newTheme: Theme) => {
    theme.value = newTheme;
    if (typeof window !== "undefined") {
      localStorage.setItem("theme", newTheme);
    }
    // Aplicar tema inmediatamente
    nextTick(() => {
      applyTheme(newTheme);
    });
  };

  const toggleTheme = () => {
    const newTheme = theme.value === "light" ? "dark" : "light";
    setTheme(newTheme);
  };

  // Observar cambios en el tema y aplicarlos
  watch(
    theme,
    (newTheme) => {
      nextTick(() => {
        applyTheme(newTheme);
      });
    },
    { immediate: true }
  );

  // Aplicar tema inicial
  if (typeof window !== "undefined") {
    nextTick(() => {
      applyTheme(theme.value);
    });
  }

  return {
    theme,
    setTheme,
    toggleTheme,
  };
});
