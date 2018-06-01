<script>
    export default {
        data () {
            return {

            }
        },
        methods: {
            defaultInfoBox () {
                this.$bkInfo({
                    title: '确认要删除？'
                })
            },
            successInfoBox () {
                this.$bkInfo({
                    type: 'success'
                })
            },
            dangerInfoBox () {
                this.$bkInfo({
                    type: 'error'
                })

                setTimeout(() => {
                    this.$bkInfo.hide()
                }, 3000)
            },
            warningInfoBox () {
                this.$bkInfo({
                    type: 'warning',
                    confirmFn(done) {
                        done()
                    }
                })
            },
            infoInfoBox () {
                this.$bkInfo({
                    type: 'loading'
                })
            },
            vnodeInfoBox () {
                let h = this.$createElement

                this.$bkInfo({
                    title: '支持VNode的提示框',
                    content: h('p', {
                        style: {
                            color: 'red'
                        }
                    }, 'Hello World')
                })
            },
            quickCloseInfoBox () {
                this.$bkInfo({
                    quickClose: true
                })
            },
            delayCloseInfoBox () {
                this.$bkInfo({
                    title: '自定义延时关闭',
                    content: '3s后自动关闭',
                    delay: 3000
                })
            },
            colorfulInfoBox () {
                this.$bkInfo({
                    title: '自定义配色',
                    content: '配置theme字段，支持primary, info, success, warning, danger',
                    theme: 'success'
                })
            },
            customIconInfoBox () {
                this.$bkInfo({
                    type: 'success',
                    icon: 'check-circle'
                })
            },
            customStatusInfoBox () {
                this.$bkInfo({
                    type: 'success',
                    statusOpts: {
                        title: '执行任务成功',
                        subtitle: false
                    }
                })
            }
        }
    }
</script>

## Info 提示框

模态对话框组件，可用于消息提示，成功提示，错误提示，后续操作询问等

### 基础用法

:::demo 调用`$bkInfo`方法，配置title等参数
```html
<template>
    <bk-button type="primary" @click="defaultInfoBox">默认配置的提示框</bk-button>
</template>

<script>
    export default {
        methods: {
            defaultInfoBox () {
                this.$bkInfo({
                    title: '确认要删除？'
                })
            }
        }
    }
</script>
```
:::

### 各种状态

:::demo 配置type的值，实现成功，错误，警告，加载中的不同类型，配置type属性后将不后出现`确认`和`取消`按钮
```html
<template>
    <bk-button type="success" @click="successInfoBox">成功</bk-button>

    <bk-button type="danger" @click="dangerInfoBox">错误</bk-button>

    <bk-button type="warning" @click="warningInfoBox">警告</bk-button>

    <bk-button type="info" @click="infoInfoBox">加载中</bk-button>
</template>

<script>
    export default {
        methods: {
            successInfoBox () {
                this.$bkInfo({
                    type: 'success'
                })
            },
            dangerInfoBox () {
                this.$bkInfo({
                    type: 'error',
                    delay: 3000
                })
            },
            warningInfoBox () {
                this.$bkInfo({
                    type: 'warning',
                    confirmFn(done) {
                        done()
                    }
                })
            },
            infoInfoBox () {
                this.$bkInfo({
                    type: 'loading'
                })
            }
        }
    }
</script>
```
:::

### 更多自定义配置

:::demo 支持VNode，点击遮罩关闭，自定义延时关闭，配色主题，特殊状态自定义图标等
```html
<template>
    <bk-button type="primary" @click="vnodeInfoBox">支持VNode</bk-button>

    <bk-button type="primary" @click="quickCloseInfoBox">点击遮罩关闭</bk-button>

    <bk-button type="primary" @click="delayCloseInfoBox">自定义延时关闭</bk-button>

    <bk-button type="primary" @click="colorfulInfoBox">配色主题</bk-button>

    <bk-button type="primary" @click="customIconInfoBox">自定义图标</bk-button>

    <bk-button type="primary" @click="customStatusInfoBox">自定义状态内容</bk-button>
</template>

<script>
    export default {
        methods: {
            vnodeInfoBox () {
                let h = this.$createElement

                this.$bkInfo({
                    title: '支持VNode的提示框',
                    content: h('p', {
                        style: {
                            color: 'red'
                        }
                    }, 'Hello World')
                })
            },
            quickCloseInfoBox () {
                this.$bkInfo({
                    quickClose: true
                })
            },
            delayCloseInfoBox () {
                this.$bkInfo({
                    title: '自定义延时关闭',
                    content: '3s后自动关闭',
                    delay: 3000
                })
            },
            colorfulInfoBox () {
                this.$bkInfo({
                    title: '自定义配色',
                    content: '配置theme字段，支持primary, info, success, warning, danger',
                    theme: 'success'
                })
            },
            customIconInfoBox () {
                this.$bkInfo({
                    type: 'success',
                    icon: 'check-circle'
                })
            },
            customStatusInfoBox () {
                this.$bkInfo({
                    type: 'success',
                    statusOpts: {
                        title: '执行任务成功',
                        subtitle: false
                    }
                })
            }
        }
    }
</script>
```
:::

### 属性
| 参数 | 说明    | 类型      | 可选值       | 默认值   |
| ---- | ------ | --------- | ----------- | -------- |
| clsName | 自定义样式 | String | —— | —— |
| type | 消息框的类型，配置type属性后将不后出现`确认`和`取消`按钮 | String | default, success, warning, error, loading | default |
| title | 消息框的标题 | String | —— | —— |
| content | 消息框的内容，仅在type为default的时候可用；可以用vm.$createElement函数生成模版 | String/VNode | —— | —— |
| icon | 消息框状态的图标，使用蓝鲸icon | String | —— | —— |
| hasHeader | 消息框是否显示头部 | Boolean | —— | true |
| statusOpts | 消息框不同状态下（当type不为default时可用）时传入的配置项，有两个key：title和subtitle，这两个配置项的值可以是String，也可以是使用vm.$createElement函数生成的模版 | String/VNode | —— | —— |
| closeIcon | 是否显示消息框关闭按钮 | Boolean | —— | true |
| theme | 消息框的主题色 | String | primary, info, success, warning, danger | primary |
| confirm | 消息框确定按钮的文字，仅当type为default时可用 | String | —— | 确定 |
| cancel | 息框取消按钮的文字，仅当type为default时可用 | String | —— | 取消 |
| quickClose | 是否允许点击遮罩关闭消息框 | Boolean | —— | false |
| delay |可设置自动关闭的时间 | Number | —— | —— |
| confirmFn | 确认按钮的回调函数，包括一个参数done，用户可手动调用，关闭消息框 | Function | —— | —— |
| cancelFn | 取消按钮的回调函数，包括一个参数done，用户可手动调用，关闭消息框 | Function | —— | —— |
| shown | 显示组件时的回调函数 | Function | —— | —— |
| hidden | 隐藏组件时的回调函数 | Function | —— | —— |
