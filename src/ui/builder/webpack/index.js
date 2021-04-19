const config = require('../config')

const moduleConfig = require('./module')
const pluginsConfig = require('./plugins')
const optimizationConfig = require('./optimization')

const { resolveBase, modeValue } = require('../utils')

module.exports = {
  mode: modeValue('production', 'development'),

  entry: {
    app: './src/main.js',
    login: './src/login/index.js'
  },

  output: {
    path: config.build.assetsRoot,
    filename: 'js/[name].[chunkhash].js',
    // webpack5中可能受optimization.chunkIds配置影响，暂不配置使用默认行为
    // chunkFilename: 'js/[id].[chunkhash].js'
    publicPath: modeValue(config.build.assetsPublicPath, config.dev.assetsPublicPath),
    clean: true // 5.20.0+
  },

  resolve: {
    modules: [resolveBase('src'), 'node_modules'],
    extensions: ['.js', '.vue', '.json'],
    alias: {
      vue$: 'vue/dist/vue.esm.js',
      '@': resolveBase('src'),
    }
  },

  module: {
    ...moduleConfig(config)
  },

  devtool: modeValue(false, 'eval-cheap-module-source-map'),

  target: 'browserslist', // default

  optimization: {
    ...optimizationConfig(config)
  },

  plugins: [
    ...pluginsConfig(config)
  ]
}
