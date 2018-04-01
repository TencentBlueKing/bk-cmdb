<script>
    export default {
        data () {
            return {
                isLoading: true,
                isDisabled: true
            };
        },
        methods: {
            handleClick (event) {
                console.log(event);
                alert('button clicked!');
            }
        }
    }
</script>

## Button 基础按钮
常用的操作按钮

### 基础用法

基础的按钮用法

:::demo 基础按钮提供5种主题，由`type`属性来定义，默认为`default`

```html
<template>
    <bk-button type="primary" title="主要按钮" @click="handleClick">
        主要按钮
    </bk-button>
</template>
<script>
    export default {
        data () {
            return {
                isLoading: true,
                isDisabled: true
            }
        },
        methods: {
            handleClick (event) {
                console.log(event)
                alert('button clicked!')
            }
        }
    }
</script>
```
:::

### 不同颜色

含有不同的颜色，代表成功、失败、警告等

:::demo 可以使用`type`属性来定义按钮的颜色，可接受`default` `primary` `warning` `success` `danger`

```html
<template>
    <bk-button type="default">
        default
    </bk-button>
    <bk-button type="primary">
        primary
    </bk-button>
    <bk-button type="warning">
        warning
    </bk-button>
    <bk-button type="success">
        success
    </bk-button>
    <bk-button type="danger">
        danger
    </bk-button>
</template>
```
:::

### 带图标

带图标的按钮可增强辨识度

:::demo 可以使用`icon`属性来定义按钮的图标

```html
<template>
    <bk-button type="primary" icon="bk">
        成功
    </bk-button>
</template>
```
:::

### 禁用状态

按钮不可用状态

:::demo 可以使用`disabled`属性来定义按钮是否禁用，它接受一个`Boolean`值

```html
<template>
    <bk-button type="primary" title="主要按钮" :disabled="isDisabled">
        主要按钮
    </bk-button>
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

点击按钮后如果需要进行数据异步操作，建议在按钮上显示加载状态

:::demo 可以使用`loading`属性来定义按钮是否显示加载中状态，它接受一个`Boolean`值

```html
<template>
    <bk-button type="primary" title="主要按钮" :loading="isLoading">
        主要按钮
    </bk-button>
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
<template>
    <bk-button type="primary" size="mini">
        超小按钮
    </bk-button>
    <bk-button type="primary" size="small">
        小按钮
    </bk-button>
    <bk-button type="primary">
        正常按钮
    </bk-button>
    <bk-button type="primary" size="large">
        大按钮
    </bk-button>
</template>

```
:::

### 按钮组 

:::demo 可以把按钮放在容器`div.bk-button-group`，实现按钮组效果

```html
<template>
    <div class="bk-button-group">
        <bk-button type="default">
            <span>北京</span>
        </bk-button>
        <bk-button type="default">
            <span>上海</span>
        </bk-button>
        <bk-button type="default">
            <span>广州</span>
        </bk-button>
        <bk-button type="default">
            <span>深圳</span>
        </bk-button>
        <bk-button type="default">
            <span>其它</span><span class="bk-icon icon-angle-down"></span>
        </bk-button>
    </div>
</template>

```
:::

### 属性
| 参数      | 说明    | 类型      | 可选值       | 默认值   |
|---------- |-------- |---------- |-------------  |-------- |
|type       |颜色     |String     |`default` `primary` `success` `warning` `danger` |default|
|size       |尺寸     |String     |`mini` `small` `normal` `large` |normal|
|icon       |图标     |String     | —— | —— |
|disabled   |是否禁用 |Boolean    |——            | false |
|loading    |是否加载中 |Boolean    |——            | false |

