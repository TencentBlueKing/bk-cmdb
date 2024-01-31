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

const config = require('../config')

const moduleConfig = require('./module')
const pluginsConfig = require('./plugins')
const optimizationConfig = require('./optimization')

const { appDir, resolveBase, modeValue } = require('../utils')

module.exports = {
  mode: modeValue('production', 'development'),

  target: 'web', // default

  context: appDir,

  entry: {
    app: './src/main.js',
    login: './src/login/index.js'
  },

  output: {
    path: config.build.assetsRoot,
    filename: 'js/[name].[chunkhash].js',
    chunkFilename: modeValue('js/[name].[chunkhash].js', 'js/[name].js'),
    publicPath: modeValue(config.build.assetsPublicPath, config.dev.assetsPublicPath),
    clean: true // 5.20.0+
  },

  resolve: {
    modules: [resolveBase('src'), 'node_modules'],
    extensions: ['.js', '.vue', '.json', '.tsx', '.ts'],
    alias: {
      vue$: 'vue/dist/vue.esm.js',
      '@': resolveBase('src')
    },
    fallback: {
      fs: false,
      assert: false,
      buffer: require.resolve('buffer/'), // lookup npm module
      path: require.resolve('path-browserify'),
      crypto: require.resolve('crypto-browserify'),
      os: require.resolve('os-browserify/browser'),
      stream: require.resolve('stream-browserify'),
      zlib: require.resolve('browserify-zlib'), // for unzip
      util: require.resolve('util/'), // for unzip
      'process/browser': require.resolve('process/browser')
    }
  },

  module: {
    ...moduleConfig(config)
  },

  devtool: modeValue(undefined, 'eval'),

  cache: {
    type: 'filesystem',
    buildDependencies: {
      config: [__filename]
    }
  },

  optimization: {
    ...optimizationConfig(config)
  },

  plugins: [
    ...pluginsConfig(config)
  ]
}
