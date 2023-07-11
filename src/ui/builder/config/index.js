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

'use strict'
// Template version: 1.3.1
// see http://vuejs-templates.github.io/webpack for documentation.

const path = require('path')
const fs = require('fs')
const parseArgs = require('minimist')
const { fixRequestBody } = require('http-proxy-middleware')


const config = {
  BUILD_TITLE: '',
  BUILD_OUTPUT: '../bin/enterprise/cmdb'
}

const argv = parseArgs(process.argv.slice(2))

process.argv.slice(2).forEach((str) => {
  const arg = str.split('=')
  if (Object.prototype.hasOwnProperty.call(config, arg[0])) {
    config[arg[0]] = arg.slice(1).join('=')
  }
})
process.CMDB_CONFIG = config
const dev = {
  // custom config
  config: Object.assign({}, config, {
    API_URL: JSON.stringify('http://{host}:{port}/proxy/'),
    API_VERSION: JSON.stringify('v3'),
    API_LOGIN: JSON.stringify(''),
    AGENT_URL: JSON.stringify(''),
    AUTH_SCHEME: JSON.stringify('internal'),
    AUTH_CENTER: JSON.stringify({}),
    BUILD_VERSION: JSON.stringify('dev'),
    USER_ROLE: JSON.stringify(1),
    USER_NAME: JSON.stringify('admin'),
    FULL_TEXT_SEARCH: JSON.stringify('off'),
    USER_MANAGE: JSON.stringify(''),
    HELP_DOC_URL: JSON.stringify(''),
    DISABLE_OPERATION_STATISTIC: false,
    COOKIE_DOMAIN: JSON.stringify(''),
    DESKTOP_URL: JSON.stringify('')
  }),

  // Paths
  assetsSubDirectory: '',
  assetsPublicPath: '/static/',
  proxyTable: {
    '/proxy': {
      logLevel: 'info',
      changeOrigin: true,
      target: 'http://{webserver地址}/',
      pathRewrite: {
        '^/proxy': ''
      }
    }
  },
  // Various Dev Server settings
  host: 'localhost', // can be overwritten by process.env.HOST
  port: 9090, // can be overwritten by process.env.PORT, if port is in use, a free one will be determined
  autoOpenBrowser: true,
  errorOverlay: false,
  notifyOnErrors: true,
  poll: false, // https://webpack.js.org/configuration/dev-server/#devserver-watchoptions-

  // Use Eslint Loader?
  // If true, your code will be linted during bundling and
  // linting errors and warnings will be shown in the console.
  useEslint: true,
  // If true, eslint errors and warnings will also be shown in the error overlay
  // in the browser.
  showEslintErrorsInOverlay: true,

  /**
     * Source Maps
     */

  // https://webpack.js.org/configuration/devtool/#development
  devtool: 'cheap-module-eval-source-map',

  // If you have problems debugging vue-files in devtools,
  // set this to false - it *may* help
  // https://vue-loader.vuejs.org/en/options.html#cachebusting
  cacheBusting: true,

  cssSourceMap: true,

  useMock: false
}

const customDevConfigPath = path.resolve(__dirname, `index.dev.${argv.env || 'ee'}.js`)
const isCustomDevConfigExist = fs.existsSync(customDevConfigPath)
if (isCustomDevConfigExist) {
  const customDevConfig = require(customDevConfigPath)
  Object.assign(dev, customDevConfig)
}

if (argv.mock) {
  // 将所有请求修改为/mock下
  dev.config.API_URL = dev.config.API_URL.replace('/proxy/', '/mock/')

  // 当devserver中的/mock未匹配时会使用此代理，此代理将/mock的请求代理回默认的/proxy
  dev.proxyTable['/mock'] = {
    // 使用默认proxy配置
    ...dev.proxyTable['/proxy'],
    // 此时地址都是/mock前缀，同样需要重写为''
    pathRewrite: {
      '^/mock': ''
    },
    // fix proxied POST requests when bodyParser is applied before this middleware
    onProxyReq: fixRequestBody
  }

  dev.useMock = true
}

module.exports = {
  dev,

  build: {
    // custom config
    config: Object.assign({}, config, {
      API_URL: '{{.site}}',
      API_VERSION: '{{.version}}',
      BUILD_VERSION: '{{.ccversion}}',
      API_LOGIN: '{{.curl}}',
      AGENT_URL: '{{.agentAppUrl}}',
      AUTH_SCHEME: '{{.authscheme}}',
      AUTH_CENTER: '{{.authCenter}}',
      USER_ROLE: '{{.role}}',
      USER_NAME: '{{.userName}}',
      FULL_TEXT_SEARCH: '{{.fullTextSearch}}',
      USER_MANAGE: '{{.userManage}}',
      HELP_DOC_URL: '{{.helpDocUrl}}',
      DISABLE_OPERATION_STATISTIC: '{{.disableOperationStatistic}}',
      COOKIE_DOMAIN: '{{.cookieDomain}}',
      DESKTOP_URL: '{{.bkDesktopUrl}}'
    }),

    // Template for index.html
    index: `${path.resolve(config.BUILD_OUTPUT)}/web/index.html`,

    // Template for login.html
    login: `${path.resolve(config.BUILD_OUTPUT)}/web/login.html`,

    // Paths
    assetsRoot: `${path.resolve(config.BUILD_OUTPUT)}/web`,

    assetsSubDirectory: '',
    assetsPublicPath: '/static/',

    /**
         * Source Maps
         */

    productionSourceMap: true,
    // https://webpack.js.org/configuration/devtool/#production
    devtool: '#source-map',

    // Gzip off by default as many popular static hosts such as
    // Surge or Netlify already gzip all static assets for you.
    // Before setting to `true`, make sure to:
    // npm install --save-dev compression-webpack-plugin
    productionGzip: false,
    productionGzipExtensions: ['js', 'css'],

    // Run the build command with an extra argument to
    // View the bundle analyzer report after build finishes:
    // `npm run build --report`
    // Set to `true` or `false` to always turn it on or off
    bundleAnalyzerReport: process.env.npm_config_report
  }
}
