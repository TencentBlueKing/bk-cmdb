<script>
    export default {
        data () {
            return {
                customBadge: {
                    visible: true
                }
            }
        },
        methods: {
            handleHideBadge () {
                this.customBadge.visible = !this.customBadge.visible
            },
            handleBadgeHover () {
                console.log('hover')
            },
            handleBadgelLeave () {
                console.log('leave')
            }
        }
    }
</script>

## Badge 标记

可以出现在任意DOM节点角上的数字或状态标记

### 基础用法

:::demo 用默认配置初始化组件
```html
<template>
    <bk-badge class="mr40" :theme="'warning'">
        <bk-button
            :type="'primary'"
            :title="'位于右上角'">
            位于右上角
        </bk-button>
    </bk-badge>
  
    <bk-badge class="mr40"
        :theme="'info'"
        :position="'bottom-right'">
        <bk-button
            :type="'primary'"
            :title="'位于右下角'">
            位于右下角
        </bk-button>
    </bk-badge>
  
    <bk-badge class="mr40"
        :theme="'success'"
        :position="'bottom-left'">
        <bk-button
            :type="'primary'"
            :title="'位于左下角'">
            位于左下角
      </bk-button>
    </bk-badge>
  
    <bk-badge class="mr40"
        :theme="'danger'"
        :position="'top-right'">
        <bk-button
            :type="'primary'"
            :title="'位于左下角'">
             位于左下角
        </bk-button>
    </bk-badge>
</template>
```
:::

### 可自定义最大值

:::demo 配置参数`val`的值
```html
<template>
    <bk-badge class="mr40"
        :theme="'danger'"
        :max="99"
        :val="150">
        <bk-button
            :type="'primary'"
            :title="'评论'">
            评论
        </bk-button>
    </bk-badge>
</template>
```
:::

### 自定义内容

:::demo 参数`val`支持数字或字符串
```html
<template>
    <bk-badge class="mr40"
        :theme="'success'"
        :val="'成功'">
        <bk-button
            :type="'primary'"
            :title="'评论'">
            评论
        </bk-button>
    </bk-badge>
  
    <bk-badge class="mr40"
        :theme="'danger'"
        :val="'失败'">
        <bk-button
            :type="'primary'"
            :title="'评论'">
            评论
        </bk-button>
    </bk-badge>
  
    <bk-badge class="mr40"
        :theme="'danger'"
        :max="99"
        :val="200"
        :visible="customBadge.visible">
        <bk-button
            :type="'primary'"
            :title="'点击可隐藏badge'"
            @click="handleHideBadge">
            点击可切换显示badge
        </bk-button>
    </bk-badge>
</template>

<script>
    export default {
        data () {
            return {
                customBadge: {
                    visible: true
                }
            }
        },
        methods: {
            handleHideBadge () {
                this.customBadge.visible = false
            },
            handleBadgeHover () {
                console.log('hover')
            },
            handleBadgelLeave () {
                console.log('leave')
            }
        }
}
</script>
```
:::

### 无内容红点

:::demo 配置参数`dot`
```html
<template>
    <bk-badge class="mr40"
        dot
        :theme="'danger'">
        <bk-button
            :type="'primary'"
            :title="'消息'">
            消息
        </bk-button>
    </bk-badge>
</template>
```
:::

### 属性
|参数           | 说明    | 类型      | 可选值       | 默认值   |
|---------------|-------- |---------- |-------------  |-------- |
| theme | 组件的主题色 | String | `primary` `info` `warning` `danger` `success`或16进制颜色值 | primary |
| val | 组件显示的值 | Number，String | —— | 1 |
| icon | 组件显示图标；当设置icon时，将忽略设置的value值 | String | —— | —— |
| max | 组件显示的最大值，当value超过max，显示数字+；仅当设置了Number类型的value值时生效 | Number | —— | —— |
| dot | 是否仅显示小圆点；当设置dot为true时，value，icon，max均会被忽略 | Boolean | —— | —— |
| visible | 是否显示组件 | Boolean | —— | —— |
| position | 组件相对于其兄弟组件的位置 | String | `top-right` `bottom-right` `bottom-left` `top-left` | top-right |

### 事件
| 事件名称 | 说明 | 回调参数 |
|---------|------|---------|
| hover | 广播给父组件mouseover事件 | —— |
| leave | 广播给父组件mouseleave事件 | —— |
