const config = require('../config')

const moduleConfig = require('./module')
const pluginsConfig = require('./plugins')
const optimizationConfig = require('./optimization')

const { appDir, resolveBase, modeValue } = require('../utils')

module.exports = {
  mode: modeValue('production', 'development'),

  target: 'web', // default

  context: appDir,

  entry: {
    app: './src/main.js',
    login: './src/login/index.js'
  },

  output: {
    path: config.build.assetsRoot,
    filename: 'js/[name].[chunkhash].js',
    chunkFilename: modeValue('js/[name].[chunkhash].js', 'js/[name].js'),
    publicPath: modeValue(config.build.assetsPublicPath, config.dev.assetsPublicPath),
    clean: true // 5.20.0+
  },

  resolve: {
    modules: [resolveBase('src'), 'node_modules'],
    extensions: ['.js', '.vue', '.json'],
    alias: {
      vue$: 'vue/dist/vue.esm.js',
      '@': resolveBase('src'),
    },
    fallback: {
      fs: false,
      buffer: false,
      assert: false,
      path: require.resolve('path-browserify')
    }
  },

  module: {
    ...moduleConfig(config)
  },

  devtool: modeValue(false, 'eval-cheap-module-source-map'),

  cache: {
    type: 'filesystem',
    buildDependencies: {
      config: [__filename]
    }
  },

  optimization: {
    ...optimizationConfig(config)
  },

  plugins: [
    ...pluginsConfig(config)
  ]
}
