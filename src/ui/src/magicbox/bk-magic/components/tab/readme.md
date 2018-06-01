<script>
    export default {
        data () {
            return {
                showTab: false
            }
        },
        methods: {
            tabChanged (name) {
                console.log(`当前tab名称：${name}`)
            },
            toggleTab () {
                this.showTab = !this.showTab
            }
        }
    }
</script>

## Tab 选项卡

展示不同类别的数据内容

### 基础用法

:::demo 默认配置的选项卡；默认显示第一个标签，通过配置`activeName`，可以默认显示指定的子面板
```html
<template>
    <bk-tab :active-name="'curveSetting'" @tab-changed="tabChanged">
        <bk-tabpanel name="curveSetting" title="项目一">
            <div class="p20">项目一</div>
        </bk-tabpanel>
        <bk-tabpanel name="exceptionDetecting" title="项目二">
            <div class="p20">项目二</div>
        </bk-tabpanel>
    </bk-tab>
</template>
```
:::

### 类型

:::demo 通过配置`type`，可以调整选项卡为背景填充效果
```html
<template>
    <bk-tab :type="'fill'" :active-name="'exceptionDetecting'" @tab-changed="tabChanged">
        <bk-tabpanel name="curveSetting" title="项目一">
            <div class="p20">项目一</div>
        </bk-tabpanel>
        <bk-tabpanel name="exceptionDetecting" title="项目二">
            <div class="p20">项目二</div>
        </bk-tabpanel>
    </bk-tab>
</template>
```
:::

### 尺寸

:::demo 通过配置`size`，可以调整选项卡尺寸
```html
<template>
    <bk-tab :size= "'small'" :active-name="'exceptionDetecting'" @tab-changed="tabChanged">
        <bk-tabpanel name="curveSetting" title="项目一">
            <div class="p15 f12">项目一</div>
        </bk-tabpanel>
        <bk-tabpanel name="exceptionDetecting" title="项目二">
            <div class="p15 f12">项目二</div>
        </bk-tabpanel>
    </bk-tab>
</template>
```
:::

### 切换显示tab

:::demo 通过配置`tabpanel`的参数`show`，可以自定义显示子面板
```html
<template>
    <section>
        <p>
            <bk-button
                :type="'primary'"
                @click="toggleTab">
                切换显示项目三
            </bk-button>
        </p>

        <bk-tab :size= "'small'" :active-name="'proj1'">
            <bk-tabpanel name="proj1" title="项目一">
                <div class="p15 f12">项目一</div>
            </bk-tabpanel>
            <bk-tabpanel name="proj2" title="项目二">
                <div class="p15 f12">项目二</div>
            </bk-tabpanel>
            <bk-tabpanel name="proj3" title="项目三" :show="showTab">
                <div class="p15 f12">项目三</div>
            </bk-tabpanel>
        </bk-tab>
    </section>
</template>
```
:::

### 自定义工具栏

:::demo 组件在选项卡右侧提供了一个工具栏的`slot`，通过传入`slot`，可实现自定义的工具栏
```html
<template>
    <bk-tab>
        <template slot="setting">
            <div class="bk-button-group">
                <button type="button" class="bk-button bk-default bk-button-small is-outline is-selected">
                    <span>天</span>
                </button>
                <button type="button" class="bk-button bk-default bk-button-small is-outline">
                    <span>周</span>
                </button>
                <button type="button" class="bk-button bk-default bk-button-small is-outline">
                    <span>月</span>
                </button>
                <button type="button" class="bk-button bk-default bk-button-small is-outline">
                    <span>年</span>
                </button>
            </div>
        </template>
        <bk-tabpanel name="curveSetting" title="项目一">
            <div class="p20">项目一</div>
        </bk-tabpanel>
        <bk-tabpanel name="exceptionDetecting" title="项目二">
            <div class="p20">项目二</div>
        </bk-tabpanel>
    </bk-tab>
</template>
```
:::

### `bk-tab`属性
| 参数 | 说明    | 类型      | 可选值       | 默认值   |
| ---- | ------ | --------- | ----------- | -------- |
| type | 选项卡类型 | String | —— | `normal` `fill` | `normal`
| size | 选项卡尺寸 | String | —— | `normal` `small` | `normal`
| active-name | 当前显示的子面板的别名，对应`bk-tabpanel`属性中的`name`字段 | String | —— | —— |

### `bk-tabpanel`属性
| 参数 | 说明    | 类型      | 可选值       | 默认值   |
| ---- | ------ | --------- | ----------- | -------- |
| name | 子面板的别名 | String | —— | —— |
| title | 子面板的名字，用作显示在选项卡的头部列表中 | String | —— | —— |
| show | 子面板是否显示 | Boolean | —— | true |

### 事件
| 事件名称 | 说明 | 回调参数 |
|---------|------|---------|
| tab-change | 切换tab的回调函数 | 切换后的tab的name |
