<script>
    export default {
        data () {
            return {
                defaultSettings: {
                    isShow: false
                },
                customSettings: {
                    isShow: false,
                    title: '更多参数配置'
                },
                moreSettings: {
                    isShow: false,
                    quickClose: true,
                    width: 500,
                    direction: 'left'
                }
            }
        },
        methods: {
            shownMoreSettings () {
                console.log('组件已打开')
            },
            hiddenMoreSettings () {
                console.log('组件已关闭')
            }
        }
    }
</script>

## SideSlider 侧栏

提供一个从两侧滑入的组件，供用户填写/查看更多信息

### 基础用法

:::demo 使用默认配置的组件
```html
<template>
    <div>
        <bk-button :type="'primary'" @click="defaultSettings.isShow = true">默认配置</bk-button>
        <bk-sideslider :is-show.sync="defaultSettings.isShow"></bk-sideslider>
    </div>
</template>

<script>
export default {
    data () {
        return {
            defaultSettings: {
                isShow: false
            }
        }
    }
}
</script>
```
:::

### 自定义标题和内容

:::demo 配置`title`参数和添加`slot`
```html
<template>
    <div>
        <bk-button :type="'primary'" @click="customSettings.isShow = true">自定义标题和内容</bk-button>
        <bk-sideslider :is-show.sync="customSettings.isShow" :title="customSettings.title">
            <div class="p20" slot="content">
                自定义显示的内容
            </div>
        </bk-sideslider>
    </div>
</template>

<script>
export default {
    data () {
        return {
            customSettings: {
                isShow: false,
                title: '更多参数配置'
            }
        }
    }
}
</script>
```
:::

### 更多自定义配置

:::demo 配置`quickClose`可定义点击遮罩关闭组件，配置`width`可定义组件宽度，配置`direction`可定义组件滑出的方向，可配置显示/关闭组件时的回调函数
```html
<template>
    <div>
        <bk-button :type="'primary'" @click="moreSettings.isShow = true">更多自定义配置1</bk-button>

        <bk-sideslider 
            :is-show.sync="moreSettings.isShow" 
            :quick-close="moreSettings.quickClose" 
            :width="moreSettings.width" 
            :direction="moreSettings.direction" 
            @shown="shownMoreSettings" 
            @hidden="hiddenMoreSettings">
        </bk-sideslider>
    </div>
</template>

<script>
    export default {
        data () {
            return {
                moreSettings: {
                    isShow: false,
                    quickClose: true,
                    width: 500,
                    direction: 'left'
                }
            }
        },
        methods: {
            shownMoreSettings () {
                console.log('组件已打开')
            },
            hiddenMoreSettings () {
                console.log('组件已关闭')
            }
        }
    }
</script>
```
:::

### 属性
| 参数 | 说明    | 类型      | 可选值       | 默认值   |
| ---- | ------ | --------- | ----------- | -------- |
| is-show | 是否显示组件，支持.sync修饰符 | Boolean | —— | false |
| title | 自定义组件标题 | String | —— | —— |
| quick-close | 是否支持点击遮罩关闭组件 | Boolean | —— | false |
| width | 组件的宽度 | Number | —— | 400 |
| direction | 组件滑出的方向 | String | left/right | right |

### 事件
| 事件名称 | 说明 | 回调参数 |
|---------|------|---------|
| shown | 显示组件后的回调参数 | —— |
| hidden | 关闭组件后的回调参数 | —— |
