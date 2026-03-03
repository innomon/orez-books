import vue from '@vitejs/plugin-vue';
import path from 'path';
import { defineConfig } from 'vite';

export default defineConfig({
  plugins: [vue()],
  resolve: {
    alias: {
      vue: 'vue/dist/vue.esm-bundler.js',
      fyo: path.resolve(__dirname, './src/fyo'),
      src: path.resolve(__dirname, './src'),
      schemas: path.resolve(__dirname, './src/schemas'),
      backend: path.resolve(__dirname, './src/backend'),
      models: path.resolve(__dirname, './src/models'),
      utils: path.resolve(__dirname, './src/utils'),
      regional: path.resolve(__dirname, './src/regional'),
      reports: path.resolve(__dirname, './src/reports'),
      dummy: path.resolve(__dirname, './src/dummy'),
      fixtures: path.resolve(__dirname, './src/fixtures'),
    },
  },
});
