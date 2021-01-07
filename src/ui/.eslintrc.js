// https://eslint.org/docs/user-guide/configuring

module.exports = {
  root: true,
  parser: 'babel-eslint',
  parserOptions: {
    sourceType: 'module'
  },
  env: {
    browser: true,
  },
  // https://github.com/standard/standard/blob/master/docs/RULES-en.md
  extends: 'standard',
  // required to lint *.vue files
  plugins: [
    'html'
  ],
  // add your custom rules here
  'rules': {
    // // allow paren-less arrow functions
    // 'arrow-parens': 0,
    // // allow async-await
    // 'generator-star-spacing': 0,
    // // allow debugger during development
    // 'no-debugger': process.env.NODE_ENV === 'production' ? 2 : 0

    // 在开发阶段打开调试
    'no-debugger': process.env.NODE_ENV === 'production' ? 2 : 0,
    // 箭头函数体可以省略圆括号
    'arrow-parens': 0,
    // 禁止空语句（可在空语句写注释避免）
    'no-empty': 'error',
    // 禁用不必要的分号
    'no-extra-semi': 'error',
    // 禁止在语句末尾使用分号（可以根据项目配置要求使用分号 ['error', 'always']）
    'semi': ['error', 'never'],
    // 允许 async-await
    'generator-star-spacing': 0,
    // 允许出现未使用过的变量
    'no-unused-vars': 0,
    // 缩进使用4个空格，并且switch语句中的Case需要缩进
    'indent': ['error', 4, {
        'SwitchCase': 1
    }],
    // 函数圆括号之前有一个空格
    'space-before-function-paren': ['error', {
        'anonymous': 'always', // 匿名函数表达式
        'named': 'always', // 命名的函数表达式
        'asyncArrow': 'always' // 异步的箭头函数表达式
    }],
    'no-trailing-spaces': ['error', {
        'skipBlankLines': true // 允许在空行使用空白符
    }],
    // 禁止使用拖尾逗号，如{demo: 'test',}
    'comma-dangle': ['error', 'never'],
    // 对象的键和值一致性
    'key-spacing': ['error', {
        'beforeColon': false
    }],
    // 关键字周围空格一致性，在关键字前后保留空格，如if () else {}
    'keyword-spacing': ['error', {
        'before': true,
        'after': true
    }],
    // 操作符周围有空格
    'space-infix-ops': 'error',
    // 注释前有空白
    'spaced-comment': ['error', 'always']
  }
}
