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
const MockJS = require('mockjs')
const bodyParser = require('body-parser')
const { pathToRegexp } = require('path-to-regexp')

const mock = require('../../mock/index')

const { HOST } = process.env
const PORT = process.env.PORT && Number(process.env.PORT)

module.exports = config => ({
  setupMiddlewares(middlewares, devServer) {
    if (!devServer) {
      throw new Error('webpack-dev-server is not defined')
    }

    if (config.dev.useMock) {
      // parse application/x-www-form-urlencoded
      devServer.app.use(bodyParser.urlencoded({ extended: true }))

      // parse application/json
      devServer.app.use(bodyParser.json())
    }

    const launchMiddleware = require('launch-editor-middleware')
    devServer.app.use('/__open-in-editor', launchMiddleware())

    devServer.app.use(/^\/mock/, (req, res, next) => {
      const mockDefs = mock.getDefs()
      let def = mockDefs[req.path]

      // 完全匹配未找到，尝试使用路径正则匹配
      if (!def) {
        for (const [path, value] of Object.entries(mockDefs)) {
          const reg = pathToRegexp(path)
          const result = reg.exec(req.path)
          if (result) {
            def = value
            req.pathRegResult = result
            break
          }
        }
      }

      // 未找到mock定义则退出，交由proxy服务处理
      if (!def) {
        return next()
      }

      let data
      if (def.data) {
        if (typeof def.data === 'function') {
          data = def.data(req)
        } else {
          data = def.data
        }
      } else if (def.path) {
        delete require.cache[def.fullpath]
        data = require(def.fullpath)
      }

      res.json(MockJS.mock(data))
    })

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
