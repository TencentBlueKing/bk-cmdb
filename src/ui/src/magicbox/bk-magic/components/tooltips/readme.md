<script>
    export default {
        data () {
            return {
                defaultTips: {
                    top: {
                        isShow: false,
                        direction: 'top',
                        content: '提示信息'
                    },
                    right: {
                        isShow: false,
                        direction: 'right',
                        content: '提示信息'
                    },
                    bottom: {
                        isShow: false,
                        direction: 'bottom',
                        content: '提示信息'
                    },
                    left: {
                        isShow: false,
                        direction: 'left',
                        content: '提示信息'
                    }
                },
                customTips: {
                    text: {
                        isShow: false
                    },
                    img: {
                        isShow: false
                    },
                    fn: {
                        isShow: false
                    }
                },
                triggerTips: {
                    isShow: true,
                    content: '提示信息',
                    direction: 'left'
                }
            }
        },
        methods: {
            getVNodeText () {
                return this.$createElement('p', {style: {color: 'red', margin: 0}}, 'Hello World!')
            },
            getVNodeImage () {
                return this.$createElement('img', {attrs: {src:     'http://magicbox.bk.tencent.com/platform/app/styles/images/magic.png'}})
            },
            onShow () {
                this.$bkMessage({
                    theme: 'success',
                    message: 'onShow'
                })
            },
            onClose () {
                this.$bkMessage({
                    theme: 'error',
                    message: 'onClose'
                })
            },
            toggleTrigger () {
                this.triggerTips.isShow = !this.triggerTips.isShow
            }
        }
    }
</script>

## Tooltips 工具提示

当鼠标指向页面元素时给出简单的提示

### 基础用法

:::demo 通过指令`v-bktooltips`展示提示框，默认提供了四种方向的出现方式，可根据需要选择
```html
    <template>
        <bk-button class="mr30" type="primary"
            v-bktooltips="{
                isShow: defaultTips.left.isShow,
                content: defaultTips.left.content,
                direction: defaultTips.left.direction
            }">左边显示</bk-button>

        <bk-button class="mr30" type="primary"
            v-bktooltips="{
                isShow: defaultTips.top.isShow,
                content: defaultTips.top.content,
                direction: defaultTips.top.direction
            }">上方显示</bk-button>

        <bk-button class="mr30" type="primary"
            v-bktooltips="{
                isShow: defaultTips.bottom.isShow,
                content: defaultTips.bottom.content,
                direction: defaultTips.bottom.direction
            }">下方显示</bk-button>


        <bk-button class="mr30" type="primary"
            v-bktooltips="{
                isShow: defaultTips.right.isShow,
                content: defaultTips.right.content,
                direction: defaultTips.right.direction
            }">右边显示</bk-button>
    </template>

    <script>
        export default {
            data () {
                return {
                    defaultTips: {
                        top: {
                            isShow: false,
                            direction: 'top',
                            content: '提示信息'
                        },
                        right: {
                            isShow: false,
                            direction: 'right',
                            content: '提示信息'
                        },
                        bottom: {
                            isShow: false,
                            direction: 'bottom',
                            content: '提示信息'
                        },
                        left: {
                            isShow: false,
                            direction: 'left',
                            content: '提示信息'
                        }
                    }
                }
            }
        }
    </script>
```
:::

### 不同的触发方式
:::demo 配置`trigger`参数可实现不同的触发方式
```html
    <template>
        <bk-button type="primary" class="mr30"
            v-bktooltips="{
                isShow: triggerTips.isShow,
                content: triggerTips.content,
                direction: triggerTips.direction,
                trigger: 'show'
            }"
            @click="toggleTrigger"
        >
            使用变量控制显示，点击按钮切换
        </bk-button>
    </template>

    <script>
        export default {
            data () {
                return {
                    triggerTips: {
                        isShow: true,
                        content: '提示信息',
                        direction: 'left'
                    }
                }
            }
        }
    </script>
```
:::

### 自定义内容

:::demo 可使用`$createElement`生成DOM节点作为提示信息的内容
```html
    <template>
        <bk-button type="primary" class="mr30"
            v-bktooltips="{
                maxWidth: 300,
                isShow: customTips.text.isShow,
                content: getVNodeText()
            }">自定义文字内容</bk-button>

        <bk-button type="primary" class="mr30"
            v-bktooltips="{
                isShow: customTips.img.isShow,
                content: getVNodeImage()
            }">自定义图片内容</bk-button>

        <bk-button type="primary" class="mr30"
            v-bktooltips="{
                isShow: customTips.fn.isShow,
                content: '提示信息',
                onShow,
                onClose
            }">自定义回调函数</bk-button>
    </template>

    <script>
        export default {
            data () {
                return {
                    customTips: {
                        text: {
                            isShow: false
                        },
                        img: {
                            isShow: false
                        },
                        fn: {
                            isShow: false
                        }
                    }
                }
            },
            methods: {
                getVNodeText () {
                    return this.$createElement('p', {style: {color: 'red', margin: 0}}, 'Hello World!')
                },
                getVNodeImage () {
                    return this.$createElement('img', {attrs: {src:     'http://magicbox.bk.tencent.com/platform/app/styles/images/magic.png'}})
                },
                onShow () {
                    this.$bkMessage({
                        theme: 'success',
                        message: 'onShow'
                    })
                },
                onClose () {
                    this.$bkMessage({
                        theme: 'error',
                        message: 'onClose'
                    })
                }
            }
        }
    </script>
```
:::

### 配置按钮

:::demo 可使用按钮作为提示信息的内容
```html
<template>
    <bk-button class="mr30" type="primary"
        v-bktooltips="{
            isShow: defaultTips.top.isShow,
            content: [
                {
                    text: '添加',
                    onclick: () => {
                        //添加回调事件
                    }
                },
                {
                    text: '删除',
                    onclick:() => {
                        //添加回调事件
                    }
                },
                {
                    text: '收藏',
                    onclick:() => {
                        //添加回调事件
                    }
                }
            ],
            direction: defaultTips.top.direction
        }">
        上边显示
    </bk-button>

    <bk-button class="mr30" type="primary"
        v-bktooltips="{
            isShow: defaultTips.bottom.isShow,
            content: [
                {
                    text: '添加',
                    onclick: () => {
                        //添加回调事件
                    }
                },
                {
                    text: '删除',
                    onclick:() => {
                        //添加回调事件
                    }
                },
                {
                    text: '收藏',
                    onclick:() => {
                        //添加回调事件
                    }
                }
            ],
            direction: defaultTips.bottom.direction
        }">下边显示</bk-button>

</template>

<script>
    export default {
        data () {
            return {
                defaultTips: {
                    top: {
                        isShow: false,
                        direction: 'top',
                        content: '提示信息'
                    },
                    bottom: {
                        isShow: false,
                        direction: 'bottom',
                        content: '提示信息'
                    },
                }
            }
        }
    }
</script>
```
:::

### 属性
| 参数 | 说明    | 类型      | 可选值       | 默认值   |
| ---- | ------ | --------- | ----------- | -------- |
| isShow | 组件是否显示 | Boolean | —— | —— |
| maxWidth | 组件最大宽度 | Number | —— | —— |
| content | 提示信息内容，接受字符串或使用`$createElement`生成的DOM节点或者对象 | String/VNode/Object | —— | 提示信息 |
| direction | 组件出现时相对于触发元素的方向 | String | top/right/bottom/left | right |
| trigger | 组件显示的方式，配置为`show`时，可通过改变参数`isShow`的值控制是否显示 | String | show | show |

### 事件
| 事件名称 | 说明 | 回调参数 |
|---------|------|---------|
| onShow | 组件显示时的回调函数 | 当前组件的实例 |
| onClose | 组件关闭时的回调函数 | 当前组件的实例 |
| onclick | 组件点击时的回调函数 | 当前组件的实例 |
