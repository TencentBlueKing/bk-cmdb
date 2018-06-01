<script>
    export default {
        data () {
            return {
                testId: 0
            }
        },
        methods: {
            handleDefault () {
                this.$bkMessage('消息')
            },
            handleDanger () {
                this.$bkMessage({
                    message: '错误类',
                    theme: 'error'
                })
            },
            handleWarning () {
                this.$bkMessage({
                    message: '提示类',
                    theme: 'warning'
                })
            },
            handleSuccess () {
                this.$bkMessage({
                    message: '成功',
                    theme: 'success'
                })
            },
            handlePrimary () {
                this.$bkMessage({
                    message: '消息',
                    theme: 'primary'
                })
            },
            handleDelay () {
                this.$bkMessage({
                    message: '自定义关闭时间，1s后关闭',
                    delay: 1000
                })
            },
            handleOnClose () {
                this.$bkMessage({
                    message: '自定义关闭回调函数',
                    onClose: () => {
                        this.$bkMessage({
                            message: '这是关闭的回调函数',
                            theme: 'success',
                            delay: 1000
                        })
                    }
                })
            },
            handleManualClose () {
                this.$bkMessage({
                    message: '可手动关闭的组件',
                    hasCloseIcon: true
                })
            },
            handleCreateElement () {
                let h = this.$createElement

                this.$bkMessage({
                    message: h('p', {
                        style: {
                            color: 'red',
                            margin: 0
                        }
                    }, 'Hello World!')
                })
            }
        },
        mounted () {}
    }
</script>

## Message 消息提示
用户操作后的消息提示，用于成功、失败、警告等消息提醒

### 基础用法

:::demo 使用默认配置的消息提示
```html
<template>
    <bk-button type="primary" @click="handleDefault">默认配置的消息提示</bk-button>
</template>

<script>
    export default {
        methods: {
            handleDefault () {
                this.$bkMessage('这是一条提示信息')
            }
        }
    }
</script>
```
:::


### 内置状态

:::demo 配置`theme`字段，使用组件的内置状态
```html
<template>
    <bk-button type="danger" @click="handleDanger">失败</bk-button>
    <bk-button type="warning" @click="handleWarning">警告</bk-button>
    <bk-button type="success" @click="handleSuccess">成功</bk-button>
    <bk-button type="primary" @click="handlePrimary">消息</bk-button>
</template>

<script>
    export default {
        methods: {
            handleDanger () {
                this.$bkMessage({
                    message: '错误类',
                    theme: 'error'
                })
            },
            handleWarning () {
                this.$bkMessage({
                    message: '提示类',
                    theme: 'warning'
                })
            },
            handleSuccess () {
                this.$bkMessage({
                    message: '成功',
                    theme: 'success'
                })
            },
            handlePrimary () {
                this.$bkMessage({
                    message: '消息',
                    theme: 'primary'
                })
            }
        }
    }
</script>
```
:::

### 更多自定义配置

:::demo 配置其它自定义字段，可实现更多类型的组件
```html
<template>
    <bk-button type="primary" @click="handleDelay">自定义关闭时间</bk-button>
    <bk-button type="primary" @click="handleOnClose">自定义关闭回调函数</bk-button>
    <bk-button type="primary" @click="handleManualClose">可手动关闭组件</bk-button>
    <bk-button type="primary" @click="handleCreateElement">自定义显示内容样式</bk-button>
</template>

<script>
    export default {
        methods: {
            handleDelay () {
                this.$bkMessage({
                    message: '自定义关闭时间，1s后关闭',
                    delay: 1000
                })
            },
            handleOnClose () {
                this.$bkMessage({
                    message: '自定义关闭回调函数',
                    onClose: () => {
                        this.$bkMessage({
                            message: '这是关闭的回调函数',
                            theme: 'success',
                            delay: 1000
                        })
                    }
                })
            },
            handleManualClose () {
                this.$bkMessage({
                    message: '可手动关闭的组件',
                    hasCloseIcon: true
                })
            },
            handleCreateElement () {
                let h = this.$createElement

                this.$bkMessage({
                    message: h('p', {
                        style: {
                            color: 'red',
                            margin: 0
                        }
                    }, 'Hello World!')
                })
            }
        }
    }
</script>
```
:::

### 属性
| 参数 | 说明 | 类型 | 可选值 | 默认值 |
|-----|------|------|-------|--------|
| theme | 组件主题色 | String | success/error/warning/primary | primary |
| icon | 组件左侧图标 | String | —— | dialogue-shape |
| message | 组件显示的文字内容，支持字符串或用`this.$createElement`生成的DOM片段 | String/DOM | —— | —— |
| delay | 组件是否延时自动关闭，值为0时需要手动关闭 | Number | —— | 3000 |
| hasCloseIcon | 是否显示关闭按钮 | Boolean | —— | false |
| onClose | 关闭组件时的回调函数 | Function | —— | —— |
| onShow | 显示组件时的回调函数 | Function | —— | —— |

### 方法
手动关闭组件时，可以调用`this.$bkMessage`返回的实例的`close`方法
| 方法 | 说明 |
| ---- | ---- |
| close | 关闭当前的组件 |
