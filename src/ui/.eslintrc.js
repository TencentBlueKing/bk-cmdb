/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2017-2022 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
 */

module.exports = {
  extends: [
    'eslint:recommended',
    'plugin:vue/recommended',
    'plugin:import/recommended',
    // 覆盖禁用冲突的规则
    'prettier',
  ],
  parser: 'vue-eslint-parser',
  parserOptions: {
    parser: '@typescript-eslint/parser',
    ecmaVersion: 2018,
    sourceType: 'module',
    tsconfigRootDir: __dirname,
    extraFileExtensions: ['.vue'],
    ecmaFeatures: {
      jsx: true,
      modules: true,
    },
  },
  env: {
    browser: true,
    node: true,
    commonjs: true,
    es6: true,
  },
  plugins: ['vue'],
  root: true,
  rules: {
    'no-async-promise-executor': 'off',
    'no-unused-vars': [
      'error',
      { vars: 'all', args: 'after-used', ignoreRestSiblings: true },
    ],
    'no-empty': ['error', { allowEmptyCatch: true }],
    'prefer-const': ['error'],
    'no-var': ['error'],

    'vue/component-definition-name-casing': 'off',
    'vue/no-mutating-props': 'off',
    'vue/this-in-template': 'off',
    'vue/multi-word-component-names': 'off',
    'vue/no-setup-props-destructure': 'off',
    'vue/require-default-prop': 'off',

    'import/no-named-as-default-member': 'off',
    'import/no-named-as-default': 'off',

    'import/order': [
      'error',
      {
        groups: [
          'builtin',
          'external',
          'internal',
          'parent',
          'sibling',
          'index',
          'object',
          'type',
        ],
        'newlines-between': 'always',
        pathGroups: [
          {
            pattern: 'vue',
            group: 'builtin',
            position: 'before',
          },
          {
            pattern: 'vuex',
            group: 'builtin',
            position: 'before',
          },
          {
            pattern: '@blueking/**',
            group: 'builtin',
            position: 'before',
          },
          {
            pattern: '@/*',
            group: 'internal',
            position: 'before',
          },
          {
            pattern: '@/api/**',
            group: 'internal',
            position: 'before',
          },
          {
            pattern: '@/magicbox/**',
            group: 'internal',
            position: 'before',
          },
          {
            pattern: '@/utils/**',
            group: 'internal',
            position: 'before',
          },
          {
            pattern: '@/dictionary/**',
            group: 'internal',
            position: 'before',
          },
          {
            pattern: '@/router/**',
            group: 'internal',
            position: 'before',
          },
          {
            pattern: '@/service/**',
            group: 'internal',
            position: 'before',
          },
          {
            pattern: '@/views/**',
            group: 'internal',
            position: 'before',
          },
          {
            pattern: '@/components/**',
            group: 'internal',
            position: 'before',
          },
          {
            pattern: '@/filters/**',
            group: 'internal',
            position: 'before',
          },
          {
            pattern: '@/mixins/**',
            group: 'internal',
            position: 'before',
          },
        ],
        distinctGroup: false,
      },
    ],
  },
  settings: {
    'import/resolver': {
      webpack: {
        config: './builder/webpack/index.js',
      },
      node: true,
    },
  },
  overrides: [
    {
      files: ['*.ts', '*.tsx'],
      parser: '@typescript-eslint/parser',
      extends: ['plugin:@typescript-eslint/recommended'],
      plugins: ['@typescript-eslint'],
      parserOptions: {
        ecmaFeatures: {
          jsx: true,
        },
      },
      rules: {
        '@typescript-eslint/semi': ['error', 'never'],
      },
    },
  ],
}
