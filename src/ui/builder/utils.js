const path = require('path')
const fs = require('fs')

const isProd = process.env.NODE_ENV === 'production'

const appDir = fs.realpathSync(process.cwd())
const resolveBase = relativePath => path.resolve(appDir, relativePath)

const modeValue = (truthy, falsy) => (isProd ? truthy : falsy)

module.exports = {
  isProd,
  appDir,
  resolveBase,
  modeValue
}
