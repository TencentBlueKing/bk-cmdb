const path = require('path')

const { HOST } = process.env
const PORT = process.env.PORT && Number(process.env.PORT)

module.exports = config => ({
  before(app) {
    const launchMiddleware = require('launch-editor-middleware')
    app.use('/__open-in-editor', launchMiddleware())
  },
  clientLogLevel: 'error',
  historyApiFallback: {
    rewrites: [
      { from: /.*/, to: path.posix.join(config.dev.assetsPublicPath, 'index.html') },
    ],
  },
  hot: true, // Enabling HMR
  contentBase: false, // since we use CopyWebpackPlugin.
  compress: true,
  host: HOST || config.dev.host,
  port: PORT || config.dev.port,
  open: config.dev.autoOpenBrowser,
  overlay: config.dev.errorOverlay
    ? { warnings: false, errors: true }
    : false,
  publicPath: config.dev.assetsPublicPath,
  proxy: config.dev.proxyTable,
  quiet: false, // necessary for FriendlyErrorsPlugin
  watchOptions: {
    poll: config.dev.poll,
  },
  stats: 'errors-only', // 'errors-only' | 'minimal' | 'normal' | 'verbose'
})
