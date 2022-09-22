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

const webpack = require('webpack')
const ESLintPlugin = require('eslint-webpack-plugin')
const { VueLoaderPlugin } = require('vue-loader')
const MiniCssExtractPlugin = require('mini-css-extract-plugin')
const HtmlWebpackPlugin = require('html-webpack-plugin')
const CopyPlugin = require('copy-webpack-plugin')
const ProgressBarPlugin = require('progress-bar-webpack-plugin')
const { BundleAnalyzerPlugin } = require('webpack-bundle-analyzer')
const chalk = require('chalk')


const { isProd, resolveBase, modeValue } = require('../utils')
const devEnv = require('../config/dev.env')
const prodEnv = require('../config/prod.env')

const getCommonPlugins = config => ([
  new ESLintPlugin({
    extensions: ['js', 'vue', 'ts', 'tsx'],
    files: ['src'],
    failOnWarning: true,
    formatter: require('eslint-friendly-formatter')
  }),

  new webpack.DefinePlugin({
    'process.env': modeValue(prodEnv, devEnv)
  }),

  new webpack.ProvidePlugin({
    process: 'process/browser',
    Buffer: ['buffer', 'Buffer'],
  }),

  new VueLoaderPlugin(),

  new HtmlWebpackPlugin({
    filename: 'index.html', // dest, relative output.path
    template: 'index.html',
    config: modeValue(config.build.config, config.dev.config),
    templateParameters: modeValue(prodEnv, devEnv),
    excludeChunks: ['login'],
    minify: isProd // default eq webpack mode
  }),
  new HtmlWebpackPlugin({
    filename: 'login.html', // dest, relative output.path
    template: 'login.html',
    config: modeValue(config.build.config, config.dev.config),
    templateParameters: modeValue(prodEnv, devEnv),
    excludeChunks: ['app']
  }),

  new CopyPlugin({
    patterns: [
      {
        from: resolveBase('static'),
        to: modeValue(config.build.assetsSubDirectory, config.dev.assetsSubDirectory),
        globOptions: { dot: true, ignore: ['.*'] }
      }
    ],
    options: {
      concurrency: 300
    }
  }),

  new webpack.ContextReplacementPlugin(/moment[/\\]locale$/, /zh-cn|en/),

  new ProgressBarPlugin({
    format: `  build [:bar] ${chalk.green.bold(':percent')} (:elapsed seconds)`,
    clear: false
  })
])

const getProdPlugins = config => ([
  new MiniCssExtractPlugin({
    filename: isProd ? 'css/[name][contenthash:7].css' : '[name].css',
    ignoreOrder: true
  })
].concat((process.env.ANALYZER || config.build.bundleAnalyzerReport) ? [
  new BundleAnalyzerPlugin()
] : []))

module.exports = (config) => {
  const commonPlugins = getCommonPlugins(config)
  const prodPlugins = getProdPlugins(config)
  return isProd ? [...commonPlugins, ...prodPlugins] : commonPlugins
}
