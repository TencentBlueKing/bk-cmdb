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

const TerserPlugin = require('terser-webpack-plugin')
const CssMinimizerPlugin = require('css-minimizer-webpack-plugin')

const { isProd } = require('../utils')

module.exports = () => ({
  // built-in optimizations works, more: https://webpack.js.org/configuration/mode/
  minimize: isProd, // is default follow mode setting
  minimizer: [
    '...',
    new CssMinimizerPlugin({
      parallel: true
    }),
    new TerserPlugin({
      exclude: /\.min\.js$/,
      parallel: true
    })
  ],
  runtimeChunk: 'single', // shared for all generated chunks
  splitChunks: {
    minChunks: 1, // default
    cacheGroups: {
      bkMagixbox: {
        test: /[\\/]bk-magic/,
        name: 'bk-magicbox',
        chunks: 'all',
        priority: 20,
        reuseExistingChunk: true,  // default
      },
      vendors: {
        test: /[\\/]node_modules[\\/]/,
        name: 'vendors',
        chunks: 'initial',
        priority: 10,
      },
      commons: {
        chunks: 'initial',
        name: 'commons',
        minChunks: 2
      }
    }
  }
})
