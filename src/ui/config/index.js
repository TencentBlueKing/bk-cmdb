'use strict'
// Template version: 1.3.1
// see http://vuejs-templates.github.io/webpack for documentation.

const path = require('path')

const config = {
    'BUILD_VERSION': '',
    'BUILD_TITLE': '配置平台 | 蓝鲸智云企业版',
    'BUILD_OUTPUT': '../bin/enterprise/cmdb'
}

process.argv.slice(2).forEach(str => {
    const argv = str.split('=')
    config[argv[0]] = argv.slice(1).join('=')
})
process.CMDB_CONFIG = config

module.exports = {
  dev: {
    // custom config
    config: Object.assign({}, config, {
        'API_URL': '""',
        'API_VERSION': '"v3"',
        'API_LOGIN': '""',
        'AGENT_URL': 'null',
        'BUILD_VERSION': 'dev',
        'USER_ROLE': '"1"',
        'USER_NAME': '"admin"'
    }),

    // Paths
    assetsSubDirectory: '',
    assetsPublicPath: '/static/',
    proxyTable: {},

    // Various Dev Server settings
    host: 'localhost', // can be overwritten by process.env.HOST
    port: 9090, // can be overwritten by process.env.PORT, if port is in use, a free one will be determined
    autoOpenBrowser: true,
    errorOverlay: true,
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

    cssSourceMap: true
  },

  build: {
    // custom config
    config: Object.assign({}, config, {
        'API_URL': '{{.site}}',
        'API_VERSION': '{{.version}}',
        'API_LOGIN': '{{.curl}}',
        'AGENT_URL': '{{.agentAppUrl}}',
        'USER_ROLE': '{{.role}}',
        'USER_NAME': '{{.userName}}'
    }),

    // Template for index.html
    index: `${path.resolve(config.BUILD_OUTPUT)}/web/index.html`,

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
