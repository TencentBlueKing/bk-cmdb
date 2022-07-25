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

const path = require('path')

const { HOST } = process.env
const PORT = process.env.PORT && Number(process.env.PORT)

module.exports = config => ({
  setupMiddlewares(middlewares, devServer) {
    if (!devServer) {
      throw new Error('webpack-dev-server is not defined')
    }

    const launchMiddleware = require('launch-editor-middleware')

    devServer.app.use('/__open-in-editor', launchMiddleware())

    return middlewares
  },

  client: {
    logging: 'warn',
    progress: false,
    overlay: config.dev.errorOverlay
      ? { warnings: false, errors: true }
      : false,
  },

  historyApiFallback: {
    rewrites: [
      { from: /.*/, to: path.posix.join(config.dev.assetsPublicPath, 'index.html') },
    ],
  },

  static: false,

  hot: true, // Enabling HMR
  compress: true,

  host: HOST || config.dev.host,
  port: PORT || config.dev.port,

  open: config.dev.autoOpenBrowser,

  proxy: config.dev.proxyTable,

  devMiddleware: {
    stats: 'errors-only', // 'errors-only' | 'minimal' | 'normal' | 'verbose'
  }
})
