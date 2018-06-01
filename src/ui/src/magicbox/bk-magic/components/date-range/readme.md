<script>
    import moment from 'moment'
    export default {
        data() {
            return {
                startDate: '2017-11-11',
                endDate: '2017-12-12',
                ranges: {
                    '昨天': [moment().subtract(1, 'days'), moment()],
                    '最近一周': [moment().subtract(7, 'days'), moment()],
                    '最近一个月': [moment().subtract(1, 'month'), moment()],
                    '最近三个月': [moment().subtract(3, 'month'), moment()]
                }
            };
        },
        methods: {
            change (oldValue, newValue) {
                console.log('改变前的值为：' + oldValue)
                console.log('改变后的值为：' + newValue)
            },
            setDate () {
                this.startDate = '2017-11-13'
                this.endDate = '2017-12-14'
            }
        }
    }
</script>

## DateRange 日期范围选择

### 基础用法

:::demo 日期范围选择器默认配置；支持`change`事件回调

```html
<template>
    <bk-daterangepicker
        @change="change"
        :range-separator="'-'"
        :quick-select="false"
        :disabled="false"
        :ranges="ranges">
    </bk-daterangepicker>
</template>
<script>
    export default {
        data () {
            return {
                ranges: {
                    '昨天': [moment().subtract(1, 'days'), moment()],
                    '最近一周': [moment().subtract(7, 'days'), moment()],
                    '最近一个月': [moment().subtract(1, 'month'), moment()],
                    '最近三个月': [moment().subtract(3, 'month'), moment()]
                }
            }
        },
        methods: {
            selected (value) {
                alert(value)
            },
            change (oldValue, newValue) {
                // Todolist
            }
        }
    }
</script>
```
:::

### 精确到时分秒

:::demo 通过配置`timer`，精确到时分秒

```html
<template>
    <bk-daterangepicker
        :range-separator="'-'"
        :quick-select="false"
        :timer="true">
</template>
```
:::

### 配置快捷选择栏

包含开关设置以及快捷列表配置

:::demo 设置`quick-select`为true开启快捷选择，通过`ranges`对象字段来自定义快捷选择配置

```html
<template>
    <bk-daterangepicker
        :quick-select="true"
        :ranges="ranges">
    </bk-daterangepicker>
</template>
<script>
    export default {
        data () {
            return {
                ranges: {
                    '昨天': [moment().subtract(1, 'days'), moment()],
                    '最近一周': [moment().subtract(7, 'days'), moment()],
                    '最近一个月': [moment().subtract(1, 'month'), moment()],
                    '最近三个月': [moment().subtract(3, 'month'), moment()]
                }
            }
        },
        methods: {
            change (oldValue, newValue) {
                // Todolist
            }
        }
    }
</script>
```
:::

### 设置默认值

设置默认值初始化日期选择器

:::demo 可以使用`icon`属性来定义按钮的图标

```html
<template>
    <bk-daterangepicker
        :quick-select="true"
        :is-show-confirm="true"
        :start-date="startDate"
        :end-date="endDate"
        :min-date="'2017-11-10'"
        :max-date="'2017-12-25'"
        :ranges="ranges">
    </bk-daterangepicker>
    <button class="bk-button bk-primary mt10" @click="setDate">改变日期</button>
</template>
<script>
    export default {
        data () {
            return {
                ranges: {
                    '昨天': [moment().subtract(1, 'days'), moment()],
                    '最近一周': [moment().subtract(7, 'days'), moment()],
                    '最近一个月': [moment().subtract(1, 'month'), moment()],
                    '最近三个月': [moment().subtract(3, 'month'), moment()]
                }
            }
        },
        methods: {
            change (oldValue, newValue) {
                // Todolist
            }
        }
    }
</script>
```
:::

### 禁用状态

禁用日期范围选择器

:::demo 使用`disabled`属性来定义日期选择器是否禁用，它接受一个`Boolean`值

```html
<template>
    <bk-daterangepicker
        :quick-select="true"
        :start-date="'2017-07-07'"
        :end-date="'2017-09-09'"
        :disabled="true"
        :ranges="ranges">
    </bk-daterangepicker>
</template>
<script>
    export default {
        data () {
            return {
                ranges: {
                    '昨天': [moment().subtract(1, 'days'), moment()],
                    '最近一周': [moment().subtract(7, 'days'), moment()],
                    '最近一个月': [moment().subtract(1, 'month'), moment()],
                    '最近三个月': [moment().subtract(3, 'month'), moment()]
                }
            }
        }
    }
</script>
```
:::

### 属性
| 参数             | 说明    | 类型      | 可选值       | 默认值   |
|----------        |-------- |---------- |-------------  |-------- |
|placeholder       |占位内容     |String     |——   |请选择日期|
|disabled          |是否禁用     |Boolean     |true / false |false|
|align             |对齐方式     |String     |left / center / right | left |
|range-separator   |日期分隔符   |String    |——           | - |
|start-date        |默认开始日期 |String    |——            | 上一个月 |
|end-date          |默认结束日期 |String    |——            | 本月 |
|quick-select      |是否开启快捷选择 |Boolean    |true / false            | false |
|ranges            |快捷菜单设置 |Object    |——            | —— |
|timer            |精确到时分秒 |Boolean    |——            | false |

### 事件
| 事件名称 | 说明 | 回调参数 |
|---------|------|---------|
| change | 日期范围改变时的回调函数 | 改变之前的日期和改变之后的日期(oldValue, newValue) |
