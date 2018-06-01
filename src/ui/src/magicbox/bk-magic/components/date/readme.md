<script>
    export default {
        data() {
            return {
                date: ''
            };
        },
        methods: {
            selected (value) {
                // alert(value)
            },
            change (oldValue, newValue) {
                console.log('改变前的值为：' + oldValue)
                console.log('改变后的值为：' + newValue)
                //   todo:...
            }
        }
    }
</script>

## Date 日期选择器

### 基础用法

:::demo 基础设置包括是否开始时间设置`timer`、默认日期`init-date`、可选起始时间范围设置`start-date` `end-date`、设置是否禁用`disabled`以及选择事件回调`@date-selected`

```html
<template>
    <bk-datepicker
        :timer="false"
        :init-date="date"
        :start-date="'2017-05-05'"
        :end-date="'2019-09-09'"
        :disabled="false"
        @date-selected="selected">
    </bk-datepicker>
</template>
<script>
    export default {
        data () {
            return {
               date: '2017-10-24' 
            }
        },
        methods: {
            selected (value) {
                alert(value)
            }
        }
    }
</script>
```
:::

### 开启时间设置

:::demo 可以使用`timer`属性来定义是否开启时间设置

```html
<template>
    <bk-datepicker
        :timer="true"
        :init-date="'2017-07-07'"
        :start-date="'2017-05-05'"
        :end-date="'2019-12-09'"
        :disabled="false"
        @date-selected="selected">
    </bk-datepicker>
</template>
<script>
    export default {
        data () {
            return {

            }
        },
        methods: {
            selected (value) {
                alert(value)
            }
        }
    }
</script>
```
:::

### 禁用状态

:::demo 可以使用`disabled`属性来定义日期选择器是否禁用，它接受一个`Boolean`值

```html
<template>
    <bk-datepicker
        :init-date="'2017-07-07'"
        :disabled="true">
    </bk-datepicker>
</template>
```
:::

### 属性
| 参数         | 说明             | 类型      | 可选值         | 默认值   |
|----------   |--------          |---------- |-------------  |-------- |
|timer        |是否开启时间设置    |String     |true / false  |true   |
|init-date    |默认日期           |String     |——             |—— |
|start-date   |可选范围开始日期    |String     |——            | —— |
|end-date     |可选范围结束日期    |String     |——              | —— |
|disabled     |是否禁用           |Boolean    |true / false    | false |

### 事件
| 事件名称 | 说明 | 回调参数 |
|---------|------|---------|
| date-selected | 选择日期时的回调函数 | 当前选择的日期(value) |
| change | 日期改变时的回调函数 | 改变之前的日期和改变之后的日期(oldValue, newValue) |
