import { defineConfig } from 'eslint/config';
import js from '@eslint/js';
import vuePlugin from 'eslint-plugin-vue';
import prettierPlugin from 'eslint-plugin-prettier';
import tsParser from '@typescript-eslint/parser';
import globals from 'globals';
import path from 'node:path';
import { fileURLToPath } from 'node:url';

const __filename = fileURLToPath(import.meta.url);
const __dirname = path.dirname(__filename);

export default defineConfig([
  js.configs.recommended,
  {
    files: ['**/*.{ts,tsx,vue}'],
    languageOptions: {
      parser: tsParser,
      parserOptions: {
        project: [path.join(__dirname, 'tsconfig.json'), path.join(__dirname, 'tsconfig.eslint.json')],
        tsconfigRootDir: __dirname,
        extraFileExtensions: ['.vue'],
        sourceType: 'module',
      },
      globals: {
        ...globals.browser,
        ...globals.node,
      },
    },
  },
  {
    files: ['**/*.vue'],
    plugins: { vue: vuePlugin },
    extends: vuePlugin.configs['flat/recommended'],
    rules: {
      'vue/html-indent': ['warn', 2],
      'vue/max-attributes-per-line': ['warn', { singleline: 3 }],
      'vue/singleline-html-element-content-newline': 'off',
      'vue/multiline-html-element-content-newline': 'off',
    },
  },
  {
    files: ['**/*.{js,ts,vue}'],
    plugins: { prettier: prettierPlugin },
    rules: { 'prettier/prettier': 'warn' },
  },
  {
    ignores: ['dist/**', 'node_modules/**'],
  },
  {
    files: ['vite.config.ts', 'uno.config.ts'],
    languageOptions: {
      parser: tsParser,
      parserOptions: { project: undefined },
    },
  },
]);
