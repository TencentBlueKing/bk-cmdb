<script>
    export default {
        data () {
            return {
                isDropdownShow: false
            }
        },
        methods: {
            dropdownShow () {
                this.isDropdownShow = true
            },
            dropdownHide () {
                this.isDropdownShow = false
            },
            triggerHandler () {
                // do
                this.$refs.dropdown.hide();
            }
        }
    }
</script>

## DropdownMenu 下拉菜单

### 基础用法

:::demo `slot[name=dropdown-trigger]`配置触发对象，`slot[name=dropdown-content]`配置下拉菜单

```html
<template>
    <bk-dropdown-menu 
        @show="dropdownShow" 
        @hide="dropdownHide" 
        ref="dropdown">
        <button class="bk-button bk-primary" slot="dropdown-trigger">
            <span>更多操作</span>
            <i :class="['bk-icon icon-angle-down',{'icon-flip': isDropdownShow}]"></i>
        </button>
        <ul class="bk-dropdown-list" slot="dropdown-content">
            <li>
                <a href="javascript:;" @click="triggerHandler">生产环境</a>
            </li>
            <li>
                <a href="javascript:;" @click="triggerHandler">预发布环境</a>
            </li>
            <li>
                <a href="javascript:;" @click="triggerHandler">测试环境</a>
            </li>
            <li>
                <a href="javascript:;" @click="triggerHandler">正式环境</a>
            </li>
        </ul>
    </bk-dropdown-menu>
</template>
<script>
    export default {
        data () {
            return {
                isDropdownShow: false
            }
        },
        methods: {
            dropdownShow () {
                this.isDropdownShow = true
            },
            dropdownHide () {
                this.isDropdownShow = false
            },
            triggerHandler () {
                // do
                this.$refs.dropdown.hide();
            }
        }
    }
</script>
```
:::

### 对齐方式
:::demo 通过配置参数`align`可以让下拉菜单以触发对象左对齐或右对齐，默认为左对齐

```html
<template>
    <bk-dropdown-menu :align="'right'">
        <button class="bk-button bk-primary is-icon" slot="dropdown-trigger">
            <i class="bk-icon icon-cog-shape"></i>
        </button>
        <ul class="bk-dropdown-list" slot="dropdown-content">
            <li>
                <a href="javascript:;">生产环境</a>
            </li>
            <li>
                <a href="javascript:;">预发布环境</a>
            </li>
            <li>
                <a href="javascript:;">测试环境</a>
            </li>
            <li>
                <a href="javascript:;">正式环境</a>
            </li>
        </ul>
    </bk-dropdown-menu>
</template>
```
:::

### 其它应用

:::demo  通过`slot[name=dropdown-trigger]`来配置触发对象

```html
<template>
    <bk-dropdown-menu :align="'center'">
        <button class="bk-button bk-primary" slot="dropdown-trigger">
            <i class="bk-icon icon-area-chart"></i>
            <span>预览</span>
        </button>
        <div slot="dropdown-content">
            <img src="" style="vertical-align: middle; padding: 5px; height: 168px;">
        </div>
    </bk-dropdown-menu>
</template>
```
:::


### 属性
|参数           | 说明    | 类型      | 可选值       | 默认值   |
|---------------|-------- |---------- |-------------  |-------- |
|align           | 与触发对象的对齐方式，包括左对齐，居中，右对齐三种方式| `left`、`center`、`right`| —— | `left` |

### 事件
| 事件名称 | 说明 | 回调参数 |
|---------|------|---------|
| show | 显示时触发此回调函数 | —— |
| hide | 隐藏时触发此回调函数 | —— |
