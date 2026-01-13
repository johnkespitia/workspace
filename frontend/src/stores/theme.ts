import { defineStore } from 'pinia';
import { ref, watch } from 'vue';

export type Theme = 'light' | 'dark';

export const useThemeStore = defineStore('theme', () => {
  const theme = ref<Theme>((localStorage.getItem('theme') as Theme) || 'light');

  const setTheme = (newTheme: Theme) => {
    theme.value = newTheme;
    localStorage.setItem('theme', newTheme);
    applyTheme(newTheme);
  };

  const toggleTheme = () => {
    setTheme(theme.value === 'light' ? 'dark' : 'light');
  };

  const applyTheme = (themeToApply: Theme) => {
    const root = document.documentElement;
    if (themeToApply === 'dark') {
      root.classList.add('dark');
    } else {
      root.classList.remove('dark');
    }
  };

  // Aplicar tema al inicializar
  applyTheme(theme.value);

  // Observar cambios en el tema
  watch(theme, (newTheme) => {
    applyTheme(newTheme);
  });

  return {
    theme,
    setTheme,
    toggleTheme,
  };
});
