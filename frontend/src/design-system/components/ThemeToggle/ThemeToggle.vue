<template>
  <button
    :class="buttonClasses"
    :aria-label="ariaLabel"
    @click="toggleTheme"
  >
    <span v-if="isDark" class="text-xl" aria-hidden="true">☀️</span>
    <span v-else class="text-xl" aria-hidden="true">🌙</span>
    <span class="sr-only">{{ ariaLabel }}</span>
  </button>
</template>

<script setup lang="ts">
import { computed } from 'vue';
import { useThemeStore } from '@/stores/theme';

const themeStore = useThemeStore();

const isDark = computed(() => themeStore.theme === 'dark');

const ariaLabel = computed(() => 
  isDark.value ? 'Cambiar a tema claro' : 'Cambiar a tema oscuro'
);

const buttonClasses = computed(() => {
  return 'inline-flex items-center justify-center w-10 h-10 rounded-lg bg-gray-100 dark:bg-gray-800 hover:bg-gray-200 dark:hover:bg-gray-700 transition-colors duration-200 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:ring-offset-2';
});

const toggleTheme = () => {
  themeStore.toggleTheme();
};
</script>

<style scoped>
.sr-only {
  position: absolute;
  width: 1px;
  height: 1px;
  padding: 0;
  margin: -1px;
  overflow: hidden;
  clip: rect(0, 0, 0, 0);
  white-space: nowrap;
  border-width: 0;
}
</style>
