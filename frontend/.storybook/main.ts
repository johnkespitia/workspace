import type { StorybookConfig } from '@storybook/vue3';
import { mergeConfig } from 'vite';
import path from 'path';

const config: StorybookConfig = {
  stories: ['../src/**/*.stories.@(js|jsx|ts|tsx|mdx)'],
  addons: [
    '@storybook/addon-links',
    '@storybook/addon-essentials',
    '@storybook/addon-interactions',
    '@storybook/addon-a11y', // Addon de accesibilidad
  ],
  framework: {
    name: '@storybook/vue3-vite',
    options: {},
  },
  async viteFinal(config) {
    return mergeConfig(config, {
      resolve: {
        alias: {
          '@': path.resolve(__dirname, '../src'),
        },
      },
      server: {
        host: '0.0.0.0', // Escuchar en todas las interfaces de red
        port: 6006,
        strictPort: true, // Fallar si el puerto está ocupado
        hmr: {
          host: 'localhost', // HMR desde el host usa localhost
          clientPort: 6006, // Puerto del host para HMR (mismo que el del contenedor)
        },
        watch: {
          usePolling: true, // Usar polling para detectar cambios en Docker
        },
      },
    });
  },
};

export default config;
