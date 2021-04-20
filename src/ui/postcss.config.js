module.exports = {
  syntax: 'postcss-scss',
  plugins: [
    // Plugins for PostCSS
    ['postcss-deep-scopable', { sels: [] }],
    'postcss-preset-env'
  ]
}
