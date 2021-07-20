process.env.NODE_ENV = 'production'

const webpack = require('webpack')
const webpackConfig = require('./webpack')

webpack(webpackConfig, (err, stats) => {
  if (err) {
    console.error(err)
    return
  }

  if (stats.hasErrors()) {
    stats.compilation.errors.forEach((e) => {
      console.error(e.message)
    })
    return
  }

  console.log(stats.toString({ colors: true }))
})
