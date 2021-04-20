process.env.NODE_ENV = 'development'

const webpack = require('webpack')
const WebpackDevServer = require('webpack-dev-server')

const config = require('./config')
const devServerConfig = require('./webpack/devserver')(config)

const webpackConfig = require('./webpack')

WebpackDevServer.addDevServerEntrypoints(webpackConfig, devServerConfig)
const compiler = webpack(webpackConfig)
const server = new WebpackDevServer(compiler, devServerConfig)

server.listen(devServerConfig.port, devServerConfig.host, (err) => {
  if (err) {
    return console.error(err)
  }
})
