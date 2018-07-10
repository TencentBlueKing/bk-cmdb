
'use strict'
// Template version: 1.1.3
// see http://vuejs-templates.github.io/webpack for documentation.

const path = require('path')
var merge = require('webpack-merge')

var buildEnv = require('./build.env')
var env4Dev = merge(require('./dev.env'), buildEnv)
var env4Build = merge(require('./prod.env'), buildEnv)
var baseOutputPath = process.cmdb.BUILD_OUTPUT ? path.resolve(process.cmdb.BUILD_OUTPUT) : path.resolve(__dirname, '../../bin/enterprise/cmdb')
module.exports = {
  build: {
    env: env4Build,
    cmdb: process.cmdb,
    index: `${baseOutputPath}/web/index.html`,
    assetsRoot: `${baseOutputPath}/web`,
    assetsSubDirectory: '',
    assetsPublicPath: '/static/',

    productionSourceMap: buildEnv.BUILD_ENV == 'product',
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
  },
  dev: {
    env: env4Dev,
    port: process.env.PORT || 8080,
    autoOpenBrowser: true,
    assetsSubDirectory: 'static',
    assetsPublicPath: '/',
    proxyTable: {
      '**': {
          logLevel: 'debug',
          filter: function(pathname, req) {
              // 代理ajax请求
              if (req.headers["x-requested-with"] == 'XMLHttpRequest') {
                  return true;
              }
          },
          target: '127.0.0.1',
          changeOrigin: true,
          onProxyReq(proxyReq, req, res) {
            console.warn(req.originalUrl, ' --> ', req.url)
          }
      }
    },
    // CSS Sourcemaps off by default because relative paths are "buggy"
    // with this option, according to the CSS-Loader README
    // (https://github.com/webpack/css-loader#sourcemaps)
    // In our experience, they generally work as expected,
    // just be aware of this issue when enabling this option.
    cssSourceMap: true
  }
}
