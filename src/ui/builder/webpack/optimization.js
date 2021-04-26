const TerserPlugin = require('terser-webpack-plugin')
const CssMinimizerPlugin = require('css-minimizer-webpack-plugin')

const { isProd } = require('../utils')

module.exports = () => ({
  // built-in optimizations works, more: https://webpack.js.org/configuration/mode/
  minimize: isProd, // is default follow mode setting
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
  runtimeChunk: 'single', // shared for all generated chunks
  splitChunks: {
    minChunks: 1, // default
    cacheGroups: {
      bkMagixbox: {
        test: /[\\/]bk-magic/,
        name: 'bk-magicbox',
        chunks: 'all',
        priority: 20,
        reuseExistingChunk: true,  // default
      },
      vendors: {
        test: /[\\/]node_modules[\\/]/,
        name: 'vendors',
        chunks: 'initial',
        priority: 10,
      },
      commons: {
        chunks: 'initial',
        name: 'commons',
        minChunks: 2
      }
    }
  }
})
