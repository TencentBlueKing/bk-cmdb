<script>
    export default {
        data () {
            return {
                isLoading: true,
                isDisabled: true
            };
        },
        methods: {
            handleClick(event) {
                console.log(event);
                alert('button clicked!');
            }
        },
        mounted () {}
    }
</script>

## Button 图标按钮
常用的操作按钮

### 基础用法

基础的按钮用法

:::demo 图标按钮提供5种主题，由`type`属性来定义，默认为`default`

```html
<template>
    <bk-icon-button type="primary" title="主要按钮" icon="bk">
        主要按钮
    </bk-icon-button>
</template>
```
:::

### 不同颜色

不同的颜色，代表成功、失败、警告等

:::demo 可以使用`type`属性来定义按钮的颜色，可接受`default` `primary` `warning` `success` `danger`

```html
<template>
    <bk-icon-button icon="bk" type="default">
        default
    </bk-icon-button>
    <bk-icon-button icon="bk" type="primary">
        primary
    </bk-icon-button>
    <bk-icon-button icon="bk" type="warning">
        warning
    </bk-icon-button>
    <bk-icon-button icon="bk" type="success">
        success
    </bk-icon-button>
    <bk-icon-button icon="bk" type="danger">
        danger
    </bk-icon-button>
</template>
```
:::

### 隐藏文本

可以隐藏文本

:::demo 可以使用`hide-text`属性来定义按钮是否隐藏文本，它接受一个`Boolean`值

```html
<template>
    <bk-icon-button type="primary" icon="bk" hide-text="true">
        成功
    </bk-icon-button>
</template>
```
:::

### 禁用状态

按钮不可用状态

:::demo 可以使用`disabled`属性来定义按钮是否禁用，它接受一个`Boolean`值

```html
<template>
    <bk-icon-button type="primary" title="按钮" :disabled="isDisabled">
        按钮
    </bk-icon-button>
</template>
<script>
    export default {
        data () {
            return {
                isDisabled: true
            }
        }
    }
</script>
```
:::

### 加载中状态

点击按钮后进行数据加载操作，在按钮上显示加载状态

:::demo 可以使用`loading`属性来定义按钮是否显示加载中状态，它接受一个`Boolean`值

```html
<template>
    <bk-icon-button type="primary" title="按钮" :loading="isLoading">
        按钮
    </bk-icon-button>
</template>
<script>
    export default {
        data () {
            return {
                isLoading: true
            }
        }
    }
</script>
```
:::

### 尺寸  

除了默认尺寸还提供另外三种尺寸，可根据不同场景选择合适的按钮

:::demo 可以使用`size`属性来定义按钮的尺寸，可接受`mini` `small` `large`

```html
<bk-icon-button type="primary" icon="bk" size="mini">
    超小按钮
</bk-icon-button>
<bk-icon-button type="primary" icon="bk" size="small">
    小按钮
</bk-icon-button>
<bk-icon-button icon="bk" type="primary">
    正常按钮
</bk-icon-button>
<bk-icon-button type="primary" icon="bk" size="large">
    大按钮
</bk-icon-button>
```
:::

### 属性
| 参数      | 说明    | 类型      | 可选值       | 默认值   |
|---------- |-------- |---------- |-------------  |-------- |
|type       |颜色     |String     |`default` `primary` `success` `warning` `danger` |default|
|size       |尺寸     |String     |`mini` `small` `normal` `large` |normal|
|icon       |图标     |String     |—— |—— |
|disabled   |是否禁用 |Boolean    |——            | false |
|hide-text   |是否隐藏文本 |Boolean    |——            | false |
|loading    |是否加载中 |Boolean    |——            | false |

