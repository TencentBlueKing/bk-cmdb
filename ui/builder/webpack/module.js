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

const { isProd, resolveBase } = require('../utils')
const MiniCssExtractPlugin = require('mini-css-extract-plugin')

const baseStyleLoaders = [
  isProd ? { loader: MiniCssExtractPlugin.loader } : 'vue-style-loader',
  {
    loader: 'css-loader',
    options: {
      esModule: false,
      sourceMap: !isProd
    }
  }
]

module.exports = () => ({
  noParse: [
    /^(vue|vue-router|vuex)$/,
    /^(axios|moment|plotly.js|cytoscape|bk-magic-vue)$/
  ],
  rules: [
    {
      test: /\.vue$/,
      loader: 'vue-loader',
    },

    {
      test: /\.tsx?$/,
      use: [
        {
          loader: 'ts-loader',
          options: {
            transpileOnly: true,
            appendTsSuffixTo: [/\.vue$/],
          }
        }
      ],
      include: [resolveBase('src')],
      exclude: [resolveBase('node_modules')]
    },

    {
      test: /\.js$/,
      use: [
        {
          loader: 'thread-loader'
        },
        {
          loader: 'babel-loader',
          options: {
            cacheDirectory: true // node_modules/.cache/babel-loader
          }
        }
      ],
      include: [resolveBase('src')],
      exclude: [resolveBase('node_modules')]
    },

    {
      test: /\.(png|jpe?g|gif|svg|webp)(\?.*)?$/,
      type: 'asset',
      generator: {
        filename: 'img/[name].[hash:7].[ext]'
      },
      parser: {
        dataUrlCondition: {
          maxSize: 8 * 1024 // 8kb, defaults
        }
      }
    },

    {
      test: /\.(mp4|webm|ogg|mp3|wav|flac|aac)(\?.*)?$/,
      type: 'asset',
      generator: {
        filename: 'media/[name].[hash:7].[ext]'
      }
    },

    {
      test: /\.(woff2?|eot|ttf|otf)(\?.*)?$/,
      type: 'asset',
      generator: {
        filename: 'fonts/[name].[hash:7].[ext]'
      }
    },

    {
      test: /\.css$/,
      use: [
        ...baseStyleLoaders,
        {
          loader: 'postcss-loader'
        }
      ]
    },

    {
      test: /\.s[ac]ss$/,
      use: [
        ...baseStyleLoaders,
        {
          loader: 'sass-loader',
          options: {
            sassOptions: {
              includePaths: [
                resolveBase('src/assets'),
                resolveBase('src/magicbox')
              ],
            }
          }
        },
        {
          loader: 'postcss-loader'
        },
        {
          loader: 'sass-resources-loader',
          options: {
            resources: [
              resolveBase('src/assets/scss/_vars.scss'),
              resolveBase('src/assets/scss/_mixins.scss'),
            ]
          }
        }
      ]
    }
  ]
})
