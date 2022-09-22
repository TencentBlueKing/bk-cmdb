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

process.env.NODE_ENV = 'development'
const chalk = require('chalk')

const webpack = require('webpack')
const WebpackDevServer = require('webpack-dev-server')

const config = require('./config')
const devServerConfig = require('./webpack/devserver')(config)

const webpackConfig = require('./webpack')

const compiler = webpack(webpackConfig)
const server = new WebpackDevServer(devServerConfig, compiler)

compiler.hooks.done.tapAsync('done', (stats, callback) => {
  if (!stats.hasErrors()) {
    console.clear()
    console.log(chalk.cyan(`\n  App running at: http://${devServerConfig.host}:${devServerConfig.port}\n`))
  }
  callback()
})

server.startCallback(() => {
  console.log('Running')
})
