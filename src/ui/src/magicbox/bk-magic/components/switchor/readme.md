<script>
    export default {
        data () {
            return {
                isSelected: true,
                isDisabled: true,
                isOutline: true,
                isSquare: true,
                showText: true
            }
        },
        methods: {
            change (val) {
                alert(val)
            }
        }
    }
</script>
<style>
    .demo-box .bk-switcher{
        margin-right: 20px;
    }
</style>

## Switcher 开关

在两种状态之间的切换

### 基础用法

:::demo 可以通过`selected`属性来定义开关是否为打开状态，同时可以通过`show-text`属性来定义开关是否显示文案

```html
<template>
    <bk-switcher 
        :selected="isSelected">
    </bk-switcher>
    <bk-switcher 
        :selected="isSelected" 
        on-text="打开"
        off-text="关闭"
        :show-text="showText">
    </bk-switcher>
</template>
<script>
    export default {
        data () {
            return {
                isSelected: true,
                showText: false
            }
        }
    }
</script>
```
:::

### 类型

提供不同的展示效果

:::demo 可以通过`is-outline`属性来定义开关是否为描边效果，同时可以通过`is-square`属性来定义开关是否为方形效果

```html
<template>
    <bk-switcher 
        :selected="isSelected"
        :is-outline="isOutline">
    </bk-switcher>
    <bk-switcher 
        :selected="isSelected"
        :is-square="isSquare">
    </bk-switcher>
    <bk-switcher 
        :selected="isSelected"
        :is-outline="isOutline"
        :is-square="isSquare">
    </bk-switcher>
</template>
<script>
    export default {
        data () {
            return {
                isSelected: true,
                isOutline: true,
                isSquare: true,
                showText: false
            }
        }
    }
</script>
```
:::

### 禁用状态

不可用状态

:::demo 可以使用`disabled`属性来定义开关是否禁用，它接受一个`Boolean`值

```html
<template>
    <bk-switcher 
        :selected="isSelected" 
        :disabled="isDisabled" 
        :show-text="showText">
    </bk-switcher>
</template>
<script>
    export default {
        data () {
            return {
                isSelected: true,
                isDisabled: true,
                showText: true
            }
        }
    }
</script>
```
:::


### 尺寸  

除了默认尺寸还提供另外一种尺寸，可根据不同场景选择合适的尺寸

:::demo 可以使用`size`属性来定义开关的尺寸，可接受`small`

```html
<template>
    <bk-switcher 
        :selected="isSelected" 
        size="small">
    </bk-switcher>
    <bk-switcher 
        :selected="isSelected">
    </bk-switcher>
</template>

<script>
    export default {
        data () {
            return {
                isSelected: true
            }
        }
    }
</script>
```
:::


### 属性
| 参数      | 说明    | 类型      | 可选值       | 默认值   |
|---------- |-------- |---------- |-------------  |-------- |
|selected   |是否打开 |Boolean    |——            | false |
|disabled   |是否禁用 |Boolean    |——            | false |
|show-text |是否显示的文本 |Boolean    |——            | false |
|on-text |打开状态显示的文本 |String    |建议1到3个字符，过长显示不全         | ON |
|off-text |关闭状态显示文本 |String    |建议1到3个字符，过长显示不全           | OFF |
|size       |尺寸 |String         |`normal` `small`           | normal |
|is-outline       |是否为描边效果 |Boolean         |——            | false |
|is-square       |是否为方形效果 |Boolean         |——            | false |


### 事件
| 参数      | 说明    | 回调参数  |
|---------- |-------- |---------- |
|change   |状态发生变化时回调函数 |新状态值（Boolean）    |

