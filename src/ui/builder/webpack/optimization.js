const TerserPlugin = require('terser-webpack-plugin')
const CssMinimizerPlugin = require('css-minimizer-webpack-plugin')

const { isProd } = require('../utils')

module.exports = () => ({
  // built-in optimizations works, more: https://webpack.js.org/configuration/mode/
  minimize: isProd, // is defaults, follow mode setting
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
  runtimeChunk: 'single' // shared for all generated chunks
})
