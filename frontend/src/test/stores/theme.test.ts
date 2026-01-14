import { describe, it, expect, beforeEach } from 'vitest';
import { setActivePinia, createPinia } from 'pinia';
import { useThemeStore } from '@/stores/theme';

describe('Theme Store', () => {
  beforeEach(() => {
    // Crear una nueva instancia de Pinia antes de cada test
    setActivePinia(createPinia());
    // Limpiar localStorage
    localStorage.clear();
  });

  it('initializes with light theme by default', () => {
    const store = useThemeStore();
    expect(store.theme).toBe('light');
  });

  it('loads theme from localStorage if available', () => {
    localStorage.setItem('theme', 'dark');
    const store = useThemeStore();
    expect(store.theme).toBe('dark');
  });

  it('sets theme correctly', () => {
    const store = useThemeStore();
    store.setTheme('dark');
    expect(store.theme).toBe('dark');
    expect(localStorage.getItem('theme')).toBe('dark');
  });

  it('toggles theme correctly', () => {
    const store = useThemeStore();
    
    // Inicialmente light
    expect(store.theme).toBe('light');
    
    // Toggle a dark
    store.toggleTheme();
    expect(store.theme).toBe('dark');
    
    // Toggle de vuelta a light
    store.toggleTheme();
    expect(store.theme).toBe('light');
  });

  it('applies dark class to document when theme is dark', () => {
    const store = useThemeStore();
    store.setTheme('dark');
    
    expect(document.documentElement.classList.contains('dark')).toBe(true);
  });

  it('removes dark class from document when theme is light', () => {
    const store = useThemeStore();
    store.setTheme('dark');
    store.setTheme('light');
    
    expect(document.documentElement.classList.contains('dark')).toBe(false);
  });
});
