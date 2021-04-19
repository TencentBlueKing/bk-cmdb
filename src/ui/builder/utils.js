const path = require('path')

const isProd = process.env.NODE_ENV === 'production'

const resolveBase = paths => path.resolve(__dirname, paths)

const modeValue = (truthy, falsy) => (isProd ? truthy : falsy)

module.exports = {
  isProd,
  resolveBase,
  modeValue
}
