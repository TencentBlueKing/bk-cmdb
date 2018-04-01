<script>
    export default {
        data () {
            return {
                demoConf1: {
                    value: 10,
                    min: 0,
                    max: 20
                },
                demoConf2: {
                    min: 0,
                    max: 10,
                    disabled: false
                },
                demoConf3: {
                    disabled: false
                },
                demoConf4: {
                    steps: 1,
                    disabled: true
                },
            }
        },
        methods: {
            change (val) {
                console.log(val)
            }
        }
    }
</script>

## NumberInput 纯数字输入框


### 基础用法

:::demo NumberInput 纯数字输入框展示基础用法

```html
<template>
    <bk-number-input
        :value.sync="demoConf1.value"
        :ex-style="{'width': '180px'}"
        :placeholder="'数字'"
        :min="demoConf1.min"
        :max="demoConf1.max"
        @change="change">
    </bk-number-input>
</template>
<script>
    export default {
        data () {
            return {
                demoConf1: {
                    value: 10,
                    min: 0,
                    max: 20
                }
            }
        },
        methods: {
            change (val) {
                console.log(val)
            }
        }
    }
</script>
```
:::

### 尺寸

可根据不同场景选择合适的尺寸

:::demo 可以使用`size`属性来定义开关的尺寸，可接受`small`

```html
<template>
    <bk-number-input
        size="small">
    </bk-number-input>
</template>
```
:::


### 禁用状态

:::demo 可以使用`disabled`属性来定义开关是否禁用，它接受一个`Boolean`值

```html
<template>
    <bk-number-input
        :value="10"
        :disabled="demoConf4.disabled">
    </bk-number-input>
</template>
<script>
    export default {
        data () {
            return {
                demoConf4: {
                    disabled: true
                },
            }
        }
    }
</script>
```
:::


### 属性
| 参数      | 说明    | 类型      | 可选值       | 默认值   |
|---------- |-------- |---------- |-------------  |-------- |
|placeholder       |空白时显示 |String         |  ——           | ——   |
|size       |尺寸 |String         |`large` `small`           | normal |
|min   |限制输入框的最小值 |Number    |——            | -infinity |
|max |是否显示的文本 |Number    |——            | infinity |
|disabled   |是否禁用 |Boolean    |——            | false |
|ex-style   |扩展样式 | ——   |——            | —— |
|hide-operation   |是否隐藏加减 | ——   |——            | false |



### 事件
| 参数      | 说明    | 回调参数  |
|---------- |-------- |---------- |
|change   |状态发生变化时回调函数 |新状态值（Boolean）    |
