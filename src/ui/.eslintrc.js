/**
 * @file 蓝鲸前端代码 ESLint 规则 Vue
 * bkfe
 */

module.exports = {
    extends: [
        'standard',
        'plugin:vue/recommended'
    ],
    parserOptions: {
        parser: 'babel-eslint',
        ecmaVersion: 2018,
        sourceType: 'module',
        ecmaFeatures: {
            jsx: true,
            modules: true
        }
    },
    env: {
        browser: true,
        node: true,
        commonjs: true,
        es6: true
    },
    // required to lint *.vue files
    plugins: [
        'vue'
    ],
    root: true,
    // add your custom rules hered
    rules: {
        // https://eslint.org/docs/rules/brace-style
        'brace-style': ['error', '1tbs', {'allowSingleLine': false}],

        // https://eslint.org/docs/rules/camelcase
        'camelcase': ['error', {'properties': 'never', 'ignoreDestructuring': true}],

        // 缩进使用 4 个空格，并且 switch 语句中的 Case 需要缩进
        // https://eslint.org/docs/rules/indent
        'indent': ['error', 4, {
            'SwitchCase': 1,
            'flatTernaryExpressions': true
        }],

        // 数组的括号内的前后禁止有空格
        // https://eslint.org/docs/rules/array-bracket-spacing
        'array-bracket-spacing': ['error', 'never'],

        // https://eslint.org/docs/rules/operator-linebreak
        'operator-linebreak': ['error', 'before'],

        // 在开发阶段打开调试
        // https://eslint.org/docs/rules/no-debugger
        'no-debugger': 'off',

        // 只有一个参数时，箭头函数体可以省略圆括号
        // https://eslint.org/docs/rules/arrow-parens
        'arrow-parens': 'off',

        // 禁止空语句（可在空语句写注释避免），允许空的 catch 语句
        // https://eslint.org/docs/rules/no-empty
        'no-empty': ['error', {'allowEmptyCatch': true}],

        // 禁止在语句末尾使用分号
        // https://eslint.org/docs/rules/semi
        'semi': ['error', 'never'],

        // 禁用不必要的分号
        // https://eslint.org/docs/rules/no-extra-semi
        'no-extra-semi': 'error',

        // generator 的 * 前面禁止有空格，后面必须有空格
        // https://eslint.org/docs/rules/generator-star-spacing
        'generator-star-spacing': [
            'error',
            {
                before: false,
                after: true
            }
        ],

        // 函数圆括号之前有一个空格
        // https://eslint.org/docs/rules/space-before-function-paren
        'space-before-function-paren': ['error', {
            'anonymous': 'always', // 匿名函数表达式
            'named': 'always', // 命名的函数表达式
            'asyncArrow': 'always' // 异步的箭头函数表达式
        }],

        // 禁止行尾有空格
        // https://eslint.org/docs/rules/no-trailing-spaces
        'no-trailing-spaces': ['error', {
            'skipBlankLines': true // 允许在空行使用空白符
        }],

        // 注释的斜线或 * 后必须有空格
        // https://eslint.org/docs/rules/spaced-comment
        'spaced-comment': ['error', 'always', {
            'line': {
                'markers': ['*package', '!', '/', ',', '=']
            },
            'block': {
                // 前后空格是否平衡
                'balanced': false,
                'markers': ['*package', '!', ',', ':', '::', 'flow-include'],
                'exceptions': ['*']
            }
        }],

        // https://eslint.org/docs/rules/no-template-curly-in-string
        // 禁止在字符串中使用字符串模板。不限制
        'no-template-curly-in-string': 'off',

        // https://eslint.org/docs/rules/no-useless-escape
        // 禁止出现没必要的转义。不限制
        'no-useless-escape': 'off',

        // https://eslint.org/docs/rules/no-var
        // 禁止使用 var
        'no-var': 'error',

        // https://eslint.org/docs/rules/prefer-const
        // 如果一个变量不会被重新赋值，必须使用 `const` 进行声明。
        'prefer-const': 'error',

        // https://github.com/vuejs/eslint-plugin-vue/blob/master/docs/rules/array-bracket-spacing.md
        'vue/array-bracket-spacing': ['error', 'never'],

        // https://github.com/vuejs/eslint-plugin-vue/blob/master/docs/rules/arrow-spacing.md
        'vue/arrow-spacing': ['error', {'before': true, 'after': true}],

        // https://github.com/vuejs/eslint-plugin-vue/blob/master/docs/rules/attribute-hyphenation.md
        'vue/attribute-hyphenation': ['error', 'always'],

        // https://github.com/vuejs/eslint-plugin-vue/blob/master/docs/rules/attributes-order.md
        // 属性顺序，不限制
        'vue/attributes-order': 'off',

        // https://github.com/vuejs/eslint-plugin-vue/blob/master/docs/rules/block-spacing.md
        'vue/block-spacing': ['error', 'always'],

        // https://github.com/vuejs/eslint-plugin-vue/blob/master/docs/rules/brace-style.md
        'vue/brace-style': ['error', '1tbs', {'allowSingleLine': false}],

        // https://github.com/vuejs/eslint-plugin-vue/blob/master/docs/rules/camelcase.md
        // 后端数据字段经常不是驼峰，所以不限制 properties，也不限制解构
        'vue/camelcase': ['error', {'properties': 'never', 'ignoreDestructuring': true}],

        // https://github.com/vuejs/eslint-plugin-vue/blob/master/docs/rules/comma-dangle.md
        // 禁止使用拖尾逗号，如 {demo: 'test',}
        'vue/comma-dangle': ['error', 'never'],

        // https://github.com/vuejs/eslint-plugin-vue/blob/master/docs/rules/comment-directive.md
        // vue 文件 template 中允许 eslint-disable eslint-enable eslint-disable-line eslint-disable-next-line
        // 行内注释启用/禁用某些规则，配置为 1 即允许
        'vue/comment-directive': 1,

        // https://github.com/vuejs/eslint-plugin-vue/blob/master/docs/rules/component-name-in-template-casing.md
        // 组件 html 标签的形式，连字符形式，所有 html 标签均会检测，如引入第三方不可避免，可通过 ignores 配置，支持正则，不限制
        'vue/component-name-in-template-casing': 'off',

        // https://github.com/vuejs/eslint-plugin-vue/blob/master/docs/rules/dot-location.md
        // 等待 https://github.com/vuejs/eslint-plugin-vue/pull/794 合入
        // 'vue/dot-location': ['error', 'property'],

        // https://github.com/vuejs/eslint-plugin-vue/blob/master/docs/rules/eqeqeq.md
        'vue/eqeqeq': ['error', 'always', {'null': 'ignore'}],

        // https://github.com/vuejs/eslint-plugin-vue/blob/master/docs/rules/html-closing-bracket-newline.md
        // 单行写法不需要换行，多行需要，不限制
        'vue/html-closing-bracket-newline': 'off',

        // https://github.com/vuejs/eslint-plugin-vue/blob/master/docs/rules/html-closing-bracket-spacing.md
        'vue/html-closing-bracket-spacing': ['error', {
            'startTag': 'never',
            'endTag': 'never',
            'selfClosingTag': 'always'
        }],

        // https://github.com/vuejs/eslint-plugin-vue/blob/master/docs/rules/html-end-tags.md
        'vue/html-end-tags': 'error',

        // https://github.com/vuejs/eslint-plugin-vue/blob/master/docs/rules/html-indent.md
        'vue/html-indent': ['error', 4, {
            'attribute': 1,
            'baseIndent': 1,
            'closeBracket': 0,
            'alignAttributesVertically': false,
            'ignores': []
        }],

        // https://github.com/vuejs/eslint-plugin-vue/blob/master/docs/rules/html-quotes.md
        'vue/html-quotes': ['error', 'double'],

        // https://github.com/vuejs/eslint-plugin-vue/blob/master/docs/rules/html-self-closing.md
        // html tag 是否自闭和，不限制
        'vue/html-self-closing': 'off',

        // https://github.com/vuejs/eslint-plugin-vue/blob/master/docs/rules/jsx-uses-vars.md
        // 当变量在 `JSX` 中被使用了，那么 eslint 就不会报出 `no-unused-vars` 的错误。需要开启 eslint no-unused-vars 规则才适用
        'vue/jsx-uses-vars': 1,

        // https://github.com/vuejs/eslint-plugin-vue/blob/master/docs/rules/key-spacing.md
        'vue/key-spacing': ['error', {'beforeColon': false, 'afterColon': true}],

        // 关键字周围空格一致性，在关键字前后保留空格，如 if () else {}
        // https://github.com/vuejs/eslint-plugin-vue/blob/master/docs/rules/keyword-spacing.md
        // 等待 https://github.com/vuejs/eslint-plugin-vue/pull/795 合入
        // 'vue/keyword-spacing': ['error', {'before': true, 'after': true}],

        // https://github.com/vuejs/eslint-plugin-vue/blob/master/docs/rules/match-component-file-name.md
        // 组件名称属性与其文件名匹配，不限制
        'vue/match-component-file-name': 'off',

        // https://github.com/vuejs/eslint-plugin-vue/blob/master/docs/rules/max-attributes-per-line.md
        // 每行属性的最大个数，不限制
        'vue/max-attributes-per-line': 'off',

        // https://github.com/vuejs/eslint-plugin-vue/blob/master/docs/rules/multiline-html-element-content-newline.md
        // 在多行元素的内容前后需要换行符，不限制
        'vue/multiline-html-element-content-newline': 'off',

        // https://github.com/vuejs/eslint-plugin-vue/blob/master/docs/rules/mustache-interpolation-spacing.md
        // template 中 {{var}}，不限制
        'vue/mustache-interpolation-spacing': 'off',

        // https://github.com/vuejs/eslint-plugin-vue/blob/master/docs/rules/name-property-casing.md
        'vue/name-property-casing': 'off',

        // https://github.com/vuejs/eslint-plugin-vue/blob/master/docs/rules/no-async-in-computed-properties.md
        'vue/no-async-in-computed-properties': 'error',

        // https://github.com/vuejs/eslint-plugin-vue/blob/master/docs/rules/no-boolean-default.md
        'vue/no-boolean-default': 'off',

        // https://github.com/vuejs/eslint-plugin-vue/blob/master/docs/rules/no-confusing-v-for-v-if.md
        'vue/no-confusing-v-for-v-if': 'error',

        // https://github.com/vuejs/eslint-plugin-vue/blob/master/docs/rules/no-dupe-keys.md
        // 二级属性名禁止重复
        'vue/no-dupe-keys': 'error',

        // https://github.com/vuejs/eslint-plugin-vue/blob/master/docs/rules/no-duplicate-attributes.md
        // 禁止 html 元素中出现重复的属性
        'vue/no-duplicate-attributes': 'error',

        // https://github.com/vuejs/eslint-plugin-vue/blob/master/docs/rules/no-empty-pattern.md
        // 等待 https://github.com/vuejs/eslint-plugin-vue/pull/798 合入
        // 禁止解构中出现空 {} 或 []
        // 'vue/no-empty-pattern': 'error',

        // https://github.com/vuejs/eslint-plugin-vue/blob/master/docs/rules/no-multi-spaces.md
        // 删除 html 标签中连续多个不用于缩进的空格
        'vue/no-multi-spaces': 'error',

        // https://github.com/vuejs/eslint-plugin-vue/blob/master/docs/rules/no-parsing-error.md
        // 禁止语法错误
        'vue/no-parsing-error': 'error',

        // https://github.com/vuejs/eslint-plugin-vue/blob/master/docs/rules/no-reserved-keys.md
        // 禁止使用保留字
        'vue/no-reserved-keys': 'error',

        // https://github.com/vuejs/eslint-plugin-vue/blob/master/docs/rules/no-restricted-syntax.md
        // 禁止使用特定的语法
        'vue/no-restricted-syntax': 'off',

        // https://github.com/vuejs/eslint-plugin-vue/blob/master/docs/rules/no-shared-component-data.md
        // data 属性必须是函数
        'vue/no-shared-component-data': 'error',

        // https://github.com/vuejs/eslint-plugin-vue/blob/master/docs/rules/no-side-effects-in-computed-properties.md
        // 禁止在计算属性对属性进行修改，不限制
        'vue/no-side-effects-in-computed-properties': 'off',

        // https://github.com/vuejs/eslint-plugin-vue/blob/master/docs/rules/no-spaces-around-equal-signs-in-attribute.md
        'vue/no-spaces-around-equal-signs-in-attribute': 'error',

        // https://github.com/vuejs/eslint-plugin-vue/blob/master/docs/rules/no-template-key.md
        // 禁止在 <template> 中使用 key 属性，不限制
        'vue/no-template-key': 'off',

        // https://github.com/vuejs/eslint-plugin-vue/blob/master/docs/rules/no-template-shadow.md
        'vue/no-template-shadow': 'error',

        // https://github.com/vuejs/eslint-plugin-vue/blob/master/docs/rules/no-textarea-mustache.md
        'vue/no-textarea-mustache': 'error',

        // https://github.com/vuejs/eslint-plugin-vue/blob/master/docs/rules/no-unused-components.md
        // 禁止 components 中声明的组件在 template 中没有使用
        'vue/no-unused-components': 'error',

        // https://github.com/vuejs/eslint-plugin-vue/blob/master/docs/rules/no-unused-vars.md
        'vue/no-unused-vars': 'error',

        // https://github.com/vuejs/eslint-plugin-vue/blob/master/docs/rules/no-use-v-if-with-v-for.md
        // 禁止 v-for 和 v-if 在同一元素上使用，不限制
        'vue/no-use-v-if-with-v-for': 'off',

        // https://github.com/vuejs/eslint-plugin-vue/blob/master/docs/rules/no-v-html.md
        // 禁止使用 v-html，防止 xss
        'vue/no-v-html': 'off',

        // https://github.com/vuejs/eslint-plugin-vue/blob/master/docs/rules/object-curly-spacing.md
        // 对象写在一行时，大括号里需要空格
        'vue/object-curly-spacing': ['error', 'always'],

        // https://github.com/vuejs/eslint-plugin-vue/blob/master/docs/rules/order-in-components.md
        // 官方推荐顺序
        'vue/order-in-components': ['error', {
            'order': [
                'el',
                'name',
                'parent',
                'functional',
                ['delimiters', 'comments'],
                ['components', 'directives', 'filters'],
                'extends',
                'mixins',
                'inheritAttrs',
                'model',
                ['props', 'propsData'],
                'data',
                'computed',
                'watch',
                // LIFECYCLE_HOOKS: ['beforeCreate', 'created', 'beforeMount', 'mounted', 'beforeUpdate', 'updated', 'activated', 'deactivated', 'beforeDestroy', 'destroyed']
                'LIFECYCLE_HOOKS',
                'methods',
                ['template', 'render'],
                'renderError'
            ]
        }],

        // https://github.com/vuejs/eslint-plugin-vue/blob/master/docs/rules/prop-name-casing.md
        // 组件 props 属性名驼峰命名
        'vue/prop-name-casing': ['error', 'camelCase'],

        // https://github.com/vuejs/eslint-plugin-vue/blob/master/docs/rules/require-component-is.md
        // <component> 元素必须要有 :is 属性
        'vue/require-component-is': 'error',

        // https://github.com/vuejs/eslint-plugin-vue/blob/master/docs/rules/require-default-prop.md
        // props 必须要有默认值，不限制
        'vue/require-default-prop': 'off',

        // https://github.com/vuejs/eslint-plugin-vue/blob/master/docs/rules/require-direct-export.md
        // 组件必须要直接被 export。不限制
        'vue/require-direct-export': 'off',

        // https://github.com/vuejs/eslint-plugin-vue/blob/master/docs/rules/require-prop-type-constructor.md
        // props 的 type 必须为构造函数，不能为字符串。
        'vue/require-prop-type-constructor': 'error',

        // https://github.com/vuejs/eslint-plugin-vue/blob/master/docs/rules/require-prop-types.md
        // props 必须要有 type。
        'vue/require-prop-types': 'error',

        // https://github.com/vuejs/eslint-plugin-vue/blob/master/docs/rules/require-render-return.md
        // render 函数必须要有返回值
        'vue/require-render-return': 'error',

        // https://github.com/vuejs/eslint-plugin-vue/blob/master/docs/rules/require-v-for-key.md
        // v-for 指令必须要有 key 属性
        'vue/require-v-for-key': 'error',

        // https://github.com/vuejs/eslint-plugin-vue/blob/master/docs/rules/require-valid-default-prop.md
        // props 默认值必须有效。不限制
        'vue/require-valid-default-prop': 'off',

        // https://github.com/vuejs/eslint-plugin-vue/blob/master/docs/rules/return-in-computed-property.md
        // 计算属性必须要有返回值
        'vue/return-in-computed-property': 'error',

        // https://github.com/vuejs/eslint-plugin-vue/blob/master/docs/rules/script-indent.md
        'vue/script-indent': ['error', 4, {
            'baseIndent': 1,
            'switchCase': 1
        }],

        // https://github.com/vuejs/eslint-plugin-vue/blob/master/docs/rules/singleline-html-element-content-newline.md
        // 单行 html 元素后面必须换行。不限制
        'vue/singleline-html-element-content-newline': 'off',

        // https://github.com/vuejs/eslint-plugin-vue/blob/master/docs/rules/space-infix-ops.md
        // 二元操作符两边要有空格
        'vue/space-infix-ops': 'error',

        // https://github.com/vuejs/eslint-plugin-vue/blob/master/docs/rules/space-unary-ops.md
        // new, delete, typeof, void, yield 等后面必须有空格，一元操作符 -, +, --, ++, !, !! 禁止有空格
        'vue/space-unary-ops': ['error', {'words': true, 'nonwords': false}],

        // https://github.com/vuejs/eslint-plugin-vue/blob/master/docs/rules/this-in-template.md
        // 不允许在 template 中使用 this
        'vue/this-in-template': ['error', 'never'],

        // https://github.com/vuejs/eslint-plugin-vue/blob/master/docs/rules/use-v-on-exact.md
        // 强制使用精确修饰词。不限制
        'vue/use-v-on-exact': 'off',

        // https://github.com/vuejs/eslint-plugin-vue/blob/master/docs/rules/v-bind-style.md
        // v-bind 指令的写法。不限制
        'vue/v-bind-style': 'off',

        // https://github.com/vuejs/eslint-plugin-vue/blob/master/docs/rules/v-on-function-call.md
        // 强制或禁止在 v-on 指令中不带参数的方法调用后使用括号。不限制
        'vue/v-on-function-call': 'off',

        // https://github.com/vuejs/eslint-plugin-vue/blob/master/docs/rules/v-on-style.md
        // v-on 指令的写法。限制简写
        'vue/v-on-style': ['error', 'shorthand'],

        // https://github.com/vuejs/eslint-plugin-vue/blob/master/docs/rules/valid-template-root.md
        // 根节点必须合法
        'vue/valid-template-root': 'error',

        // https://github.com/vuejs/eslint-plugin-vue/blob/master/docs/rules/valid-v-bind.md
        // v-bind 指令必须合法
        'vue/valid-v-bind': 'error',

        // https://github.com/vuejs/eslint-plugin-vue/blob/master/docs/rules/valid-v-cloak.md
        // v-cloak 指令必须合法
        'vue/valid-v-cloak': 'error',

        // https://github.com/vuejs/eslint-plugin-vue/blob/master/docs/rules/valid-v-else-if.md
        // v-else-if 指令必须合法
        'vue/valid-v-else-if': 'error',

        // https://github.com/vuejs/eslint-plugin-vue/blob/master/docs/rules/valid-v-else.md
        // v-else 指令必须合法
        'vue/valid-v-else': 'error',

        // https://github.com/vuejs/eslint-plugin-vue/blob/master/docs/rules/valid-v-for.md
        // valid-v-for 指令必须合法
        'vue/valid-v-for': 'error',

        // https://github.com/vuejs/eslint-plugin-vue/blob/master/docs/rules/valid-v-html.md
        // v-html 指令必须合法
        'vue/valid-v-html': 'error',

        // https://github.com/vuejs/eslint-plugin-vue/blob/master/docs/rules/valid-v-if.md
        // v-if 指令必须合法
        'vue/valid-v-if': 'error',

        // https://github.com/vuejs/eslint-plugin-vue/blob/master/docs/rules/valid-v-model.md
        // v-model 指令必须合法
        'vue/valid-v-model': 'error',

        // https://github.com/vuejs/eslint-plugin-vue/blob/master/docs/rules/valid-v-on.md
        // v-on 指令必须合法
        'vue/valid-v-on': 'error',

        // https://github.com/vuejs/eslint-plugin-vue/blob/master/docs/rules/valid-v-once.md
        // v-once 指令必须合法
        'vue/valid-v-once': 'error',

        // https://github.com/vuejs/eslint-plugin-vue/blob/master/docs/rules/valid-v-pre.md
        // v-pre 指令必须合法
        'vue/valid-v-pre': 'error',

        // https://github.com/vuejs/eslint-plugin-vue/blob/master/docs/rules/valid-v-show.md
        // v-show 指令必须合法
        'vue/valid-v-show': 'error',

        // https://github.com/vuejs/eslint-plugin-vue/blob/master/docs/rules/valid-v-text.md
        // v-text 指令必须合法
        'vue/valid-v-text': 'error'
    },
    overrides: [
        {
            files: ['*.vue'],
            rules: {
                indent: 'off'
            }
        }
    ]
}