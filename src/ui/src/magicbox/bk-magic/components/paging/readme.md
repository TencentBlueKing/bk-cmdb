<script>
    export default {
        data () {
            return {
                defaultPaging: {
                    page: 1,
                    totalPage: 5
                },
                customPaging: {
                    page1: 2,
                    totalPage1: 4,
                    page2: 7,
                    totalPage2: 100
                },
                compactPaging: {
                    page: 1,
                    type: 'compact'
                }
            }
        }
    }
</script>

## Paging 分页

分页显示数据

### 基础用法

:::demo 默认配置的组件
```html
<template>
    <bk-paging 
        :cur-page.sync="defaultPaging.page" 
        :total-page="defaultPaging.totalPage">
    </bk-paging>
</template>

<script>
export default {
    data () {
        return {
            defaultPaging: {
                page: 1,
                totalPage: 5
            }
        }
    }
}
</script>
```
:::

### 紧凑效果

:::demo 配置`type`字段
```html
<template>
    <bk-paging 
        :cur-page.sync="compactPaging.page" 
        :type="'compact'">
    </bk-paging>
</template>

<script>
    export default {
        data () {
            return {
                compactPaging: {
                    page: 1
                }
            }
        }
    }
</script>
```
:::

### 尺寸

:::demo 配置`size`字段
```html
<template>
    <bk-paging 
        :cur-page.sync="defaultPaging.page" 
        :size="'small'"
        :total-page="defaultPaging.totalPage">
    </bk-paging>
</template>

<script>
    export default {
        data () {
            return {
                defaultPaging: {
                    page: 1,
                    totalPage: 5
                }
            }
        }
    }
</script>
```
:::

### 自定义配置

:::demo 可配置总页数，当前页码等
```html
<template>
    <div class="mb20">
        <bk-paging 
            :cur-page.sync="customPaging.page1" 
            :total-page="customPaging.totalPage1">
        </bk-paging>
    </div>
    
    <div>
        <bk-paging 
            :cur-page.sync="customPaging.page2" 
            :total-page="customPaging.totalPage2" 
            :type="'compact'">
        </bk-paging>
    </div>
</template>

<script>
    export default {
        data () {
            return {
                customPaging: {
                    page1: 2,
                    totalPage1: 4,
                    page2: 7,
                    totalPage2: 100
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
| type | 组件的类型 | String | `default` `compact` | default |
| total-page | 总页数 | Number | —— | 20 |
| cur-page | 当前页码，支持`.sync`修饰符 | Number | —— | 1 |

### 事件
| 事件名称 | 说明 | 回调参数 |
|---------|------|---------|
| page-change | 当前页码变化时的回调 | 变化后的页码 |
