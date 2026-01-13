import type { StorybookConfig } from '@storybook/vue3';
import { mergeConfig } from 'vite';
import path from 'path';

const config: StorybookConfig = {
  stories: ['../src/**/*.stories.@(js|jsx|ts|tsx|mdx)'],
  addons: [
    '@storybook/addon-links',
    '@storybook/addon-essentials',
    '@storybook/addon-interactions',
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
    });
  },
};

export default config;
