<script>
    export default {
        data () {
            return {
                iconSteps: ['edit2', 'image', 'check-1'],
                iconCurStep: 2,
                objectSteps: [{
                    title: '订单详情',
                    icon: 'order'
                },
                {
                    title: '收件人信息',
                    icon: 'user'
                },
                {
                    title: '支付',
                    icon: 'id'
                },
                {
                    title: '完成',
                    icon: 'check-1'
                }],
                controllableSteps: {
                    controllable: true,
                    steps: [{
                        title: '基本信息',
                        icon: 1
                    },
                    {
                        title: '详细信息',
                        icon: 2
                    },
                    {
                        title: '实名认证',
                        icon: 3
                    }],
                    curStep: 1
                }
            }
        },
        methods: {
            stepChanged (index) {
                console.log(`当前步骤index：${index}`)
            }
        }
    }
</script>

## Steps 步骤

引导用户按步骤完成流程的组件

### 基础用法

:::demo 使用默认配置的组件
```html
<template>
    <bk-steps></bk-steps>
</template>
```
:::

### 自定义步骤内容

:::demo 配置`steps`参数，具体内容参考下方属性表格
```html
<template>
    <bk-steps 
        :steps="iconSteps" 
        :cur-step.sync="iconCurStep" 
        class="mb40">
    </bk-steps>
    <bk-steps :steps="objectSteps"></bk-steps>
</template>

<script>
    export default {
        data () {
            return {
                iconSteps: ['edit2', 'image', 'check-1'],
                iconCurStep: 2,
                objectSteps: [{
                    title: '订单详情',
                    icon: 'order'
                },
                {
                    title: '收件人信息',
                    icon: 'user'
                },
                {
                    title: '支付',
                    icon: 'id'
                },
                {
                    title: '完成',
                    icon: 'check-1'
                }]
            }
        }
    }
</script>
```
:::

### 可控制当前步骤

:::demo 配置`controllable`参数，可实现点击步骤任意跳转
```html
<template>
    <bk-steps 
        :controllable="controllableSteps.controllable" 
        :steps="controllableSteps.steps" 
        :cur-step.sync="controllableSteps.curStep" 
        @step-changed="stepChanged">
    </bk-steps>
</template>

<script>
    export default {
        data () {
            return {
                controllableSteps: {
                    controllable: true,
                    steps: [{
                        title: '基本信息',
                        icon: 1
                    },
                    {
                        title: '详细信息',
                        icon: 2
                    },
                    {
                        title: '实名认证',
                        icon: 3
                    }],
                    curStep: 1
                }
            }
        },
        stepChanged (index) {
            console.log(`当前步骤index：${index}`)
        }
    }
</script>
```
:::

### 属性
| 参数 | 说明    | 类型      | 可选值       | 默认值   |
| ---- | ------ | --------- | ----------- | -------- |
| steps | 组件步骤内容，数组中的元素可以是对象，可以是数字，可以是字符串，也可以是三者混合；当元素是对象时，有两个可选的key：title和icon；当元素是数字或字符串时，组件会将其解析成对象形式下的icon的值；当元素是字符串时，使用蓝鲸icon | Array | —— | —— |
| cur-step | 当前步骤的索引值，从1开始；支持.sync修饰符 | Number | —— | 1 |
| controllable | 步骤可否被控制前后跳转 | Boolean | —— | false |

### 事件
| 事件名称 | 说明 | 回调参数 |
|---------|------|---------|
| step-changed | 当前步骤变化时的回调 | 变化后的步骤index |
