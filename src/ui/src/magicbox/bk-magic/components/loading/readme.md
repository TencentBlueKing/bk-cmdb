<script>
    export default {
        data () {
            return {
                basicLoading: true,
                textLoading: true
            }
        },
        methods: {
            showLoading () {
                let h = this.$createElement

                this.$bkLoading({
                    title: h('span', {
                        style: {
                            color: 'red'
                        }
                    }, '加载中')
                })

                setTimeout(() => {
                    this.$bkLoading.hide()
                }, 3000)
            }
        }
    }
</script>

<style>
    .test-dom{
        height: 300px;
        line-height: 300px;
        border: 1px solid #eee;
        text-align: center;
    }
</style>

## Loading 加载

覆盖正在加载数据的组件一个loading层

### 基础用法

:::demo 组件提供了自定义指令`v-bkloading`，方便对指定DOM节点添加Loading效果
```html
<template>
    <div class="test-dom" v-bkloading="{isLoading: basicLoading}">
        内容
    </div>
</template>

<script>
    export default {
        data () {
            return {
                basicLoading: true
            }
        }
    }
</script>
```
:::

### 配置文案
:::demo 传入`title`,值会被渲染到loading图标的下方
```html
<template>
    <div class="test-dom" v-bkloading="{isLoading: textLoading, title: '数据加载中'}">
        内容
    </div>
</template>

<script>
    export default {
        data () {
            return {
                textLoading: true
            }
        }
    }
</script>
```
:::

### 全屏加载
:::demo `Loading`组件还提供了绑定在`Vue`对象的原型方法
```html
<template>
    <bk-button type="primary" @click="showLoading">显示全屏加载，3s后关闭</bk-button>
</template>

<script>
    export default {
        methods: {
            showLoading () {
                let h = this.$createElement

                this.$bkLoading({
                    title: h('span', {
                        style: {
                            color: 'red'
                        }
                    }, '加载中')
                })

                setTimeout(() => {
                    this.$bkLoading.hide()
                }, 3000)
            }
        }
    }
</script>
```
:::

### 属性
| 参数 | 说明    | 类型      | 可选值       | 默认值   |
| ---- | ------ | --------- | ----------- | -------- |
| isLoading | 一个`Boolean`变量，控制组件的显示 | Boolean | —— | —— |
| title | 组件显示时的文案，支持使用`createElement`函数生成的VNode | String | —— | —— |
