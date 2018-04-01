<script>
    export default {
        data () {
            return {
                list: [{
                    id: 1,
                    name: '爬山',
                    alias: 'climb'
                }, 
                {
                    id: 2,
                    name: '跑步',
                    alias: 'running'
                }, 
                {
                    id: 3,
                    name: '打球',
                    alias: 'play-ball'
                }, 
                {
                    id: 4,
                    name: '跳舞',
                    alias: 'dancing'
                }, 
                {
                    id: 5,
                    name: '健身',
                    alias: 'gym'
                }, 
                {
                    id: 6,
                    name: '骑车',
                    alias: 'cycling'
                }],
                defaultDemo: {
                    selected: 0
                },
                disabledDemo: {
                    selected: 1,
                    disabled: true
                },
                multiDemo: {
                    selected: [0, 1],
                    multiSelect: true
                },
                searchDemo: {
                    selected: -1,
                    searchable: true,
                    hasCreateItem: false,
                    searchKey: 'alias'
                },
                customDemo: {
                    selected: 0,
                    hasCreateItem: false,
                    tools: {
                        edit: true,
                        del: true
                    },
                    allowClear: true
                }
            }
        },
        methods: {
            selected (index, data) {
                console.log(data.name)
            },
            customDemoEdit () {
                console.log('edit')
            },
            customDemoDelete () {
                console.log('delete')
            },
            defaultDemoCreate () {
                console.log('create')
            },
            customSelected (value, item) {
                console.log(value, item)
            },
            test (v, k) {
                console.log(v, k)
            }
        }
    }
</script>

## Dropdown 下拉选框

将动作或菜单折叠到下拉菜单中，支持单选和多选

### 基础用法

:::demo 用默认配置初始化组件，参数list和selected为必传

```html
<template>
    <bk-dropdown
        :list="list"
        :selected.sync="defaultDemo.selected"
        :setting-key="'alias'"
        @item-selected="selected">
    </bk-dropdown>
</template>

<script>
    export default {
        data() {
            return {
                list: [{
                    id: 1,
                    name: '爬山',
                    alias: 'climb'
                }, 
                {
                    id: 2,
                    name: '跑步',
                    alias: 'running'
                }, 
                {
                    id: 3,
                    name: '打球',
                    alias: 'play-ball'
                }, 
                {
                    id: 4,
                    name: '跳舞',
                    alias: 'dancing'
                }, 
                {
                    id: 5,
                    name: '健身',
                    alias: 'gym'
                }, 
                {
                    id: 6,
                    name: '骑车',
                    alias: 'cycling'
                }],
                defaultDemo: {
                    selected: 0
                }
            }
        },
        methods: {
            selected (index, data) {
                console.log(data.name)
            }
        }
    }
</script>
```
:::

### 禁用状态

组件的禁用状态

:::demo 配置禁用状态的参数`disabled`为`true`，组件不可用，默认选中第二项

```html
<template>
    <bk-dropdown
        :list="list"
        :selected.sync="disabledDemo.selected"
        :disabled="disabledDemo.disabled">
    </bk-dropdown>
</template>

<script>
    export default {
        data () {
            return {
                list: [
                    {
                        id: 1,
                        name: '爬山',
                        alias: 'climb'
                    },
                    {
                        id: 2,
                        name: '跑步',
                        alias: 'running'
                    },
                    {
                        id: 3,
                        name: '打球',
                        alias: 'play-ball'
                    },
                    {
                        id: 4,
                        name: '跳舞',
                        alias: 'dancing'
                    },
                    {
                        id: 5,
                        name: '健身',
                        alias: 'gym'
                    },
                    {
                        id: 6,
                        name: '骑车',
                        alias: 'cycling'
                    }
                ],
                disabledDemo: {
                    selected: 1,
                    disabled: true
                }
            }
        }
    }
</script>
```
:::

### 多选

可从下拉框中选择多个选项

:::demo 配置`multi-select`为`true`，此时`selected`为选中项在`list`中的索引值
```html
<template>
    <bk-dropdown
        :list="list"
        :selected.sync="multiDemo.selected"
        :multi-select="multiDemo.multiSelect">
    </bk-dropdown>
</template>

<script>
    export default {
        data () {
            return {
                list: [
                    {
                        id: 1,
                        name: '爬山',
                         alias: 'climb'
                    },
                    {
                        id: 2,
                        name: '跑步',
                        alias: 'running'
                    },
                    {
                        id: 3,
                        name: '打球',
                        alias: 'play-ball'
                    },
                    {
                        id: 4,
                        name: '跳舞',
                        alias: 'dancing'
                    },
                    {
                        id: 5,
                        name: '健身',
                        alias: 'gym'
                    },
                    {
                        id: 6,
                        name: '骑车',
                        alias: 'cycling'
                    }
                ],
                multiDemo: {
                    selected: [0, 1],
                    multiSelect: true
                }
            }
        }
    }
</script>
```
:::

### 可筛选选项

:::demo 配置`searchable`为`true`
```html
<template>
    <bk-dropdown
        :list="list"
        :selected.sync="searchDemo.selected"
        :searchable="searchDemo.searchable"
        :has-create-item="searchDemo.hasCreateItem">
    </bk-dropdown>
</template>

<script>
  export default {
    data () {
      return {
        list: [
          {
            id: 1,
            name: '爬山',
            alias: 'climb'
          },
          {
            id: 2,
            name: '跑步',
            alias: 'running'
          },
          {
            id: 3,
            name: '打球',
            alias: 'play-ball'
          },
          {
            id: 4,
            name: '跳舞',
            alias: 'dancing'
          },
          {
            id: 5,
            name: '健身',
            alias: 'gym'
          },
          {
            id: 6,
            name: '骑车',
            alias: 'cycling'
          }
        ],
        searchDemo: {
          selected: -1,
          searchable: true,
          hasCreateItem: false,
          searchKey: 'alias'
        }
      }
    }
  }
</script>
```
:::

### 更多自定义配置

:::demo 根据业务需求配置其它参数
```html
<template>
    <bk-dropdown
        :list="list"
        :selected.sync="customDemo.selected"
        :has-create-item="customDemo.hasCreateItem"
        :tools="customDemo.tools"
        :allow-clear="customDemo.allowClear"
        :setting-key="'alias'"
        @edit="customDemoEdit"
        @del="customDemoDelete"
        @item-selected="customSelected">
    </bk-dropdown>
</template>

<script>
  export default {
        data () {
            return {
                list: [
                    {
                        id: 1,
                        name: '爬山',
                        alias: 'climb'
                    },
                    {
                        id: 2,
                        name: '跑步',
                        alias: 'running'
                    },
                    {
                        id: 3,
                        name: '打球',
                        alias: 'play-ball'
                    },
                    {
                        id: 4,
                        name: '跳舞',
                        alias: 'dancing'
                    },
                    {
                        id: 5,
                        name: '健身',
                        alias: 'gym'
                    },
                    {
                        id: 6,
                        name: '骑车',
                        alias: 'cycling'
                    }
                ],
                customDemo: {
                    selected: 0,
                    hasCreateItem: false,
                    tools: {
                        edit: true,
                        del: true
                    },
                    allowClear: true
                }
            }
        },
        methods: {
            customDemoEdit () {
                console.log('edit')
            },
            customDemoDelete () {
                console.log('delete')
            },
            customSelected (value, item) {
                console.log(value, item)
            }
        }
    }
</script>
```
:::

### 属性
|参数           | 说明    | 类型      | 可选值       | 默认值   |
|---------------|-------- |---------- |-------------  |-------- |
|list           | 下拉菜单所需的数据列表，必选 | Array | —— | —— |
|selected       | 选中的项的index值，支持.sync修饰符，必选 | Number/Array | —— | —— |
|extCls         | 自定义样式的类名，会被加在最外层的DOM节点上 | String | —— | —— |
|has-create-item  | 拉菜单中是否有新增项 | Boolean | —— | false |
|create-text     | 下拉菜单中新增项的文字 | String | —— | —— |
|tools          | 待选项右侧的工具按钮，有两个可配置的key：edit和del，默认为两者都显示。| Object/Boolean | —— | [{edit: true, del: true}] |
|placeholder    | 未选选项时，组件显示的文字 | String | —— | '请选择' |
|displayKey     | 循环list时，显示字段的key值 | String | —— | 'name' |
|disabled       | 是否禁用组件 | Boolean | —— | false |
|multi-select    | 是否支持多选 | Boolean | —— | false |
|searchable     | 是否支持筛选 | Boolean | —— | false |
|search-key      | 筛选时，搜索的key值 | String | —— | 'name' |
|allow-clear     | 是否可以清除单选时选中的项 | Boolean | —— | false |
|setting-key     | 单选选中某一项时，根据当前字段确定`item-selected`事件传入的第一个参数的值 | String | —— | 'id' |

### 事件
| 事件名称 | 说明 | 回调参数 |
|---------|------|---------|
| item-selected | 选中选项时触发 | 选中的选项的index值 |
| typing | 输入筛选条件时触发 | 当前输入的值 |
| clear | 单选模式下清空当前选择时触发 | 当前清除的项的索引值 |
| visible-toggle | 下拉框打开/关闭时触发 | 打开为true，关闭为false |
| edit | 编辑当前项 | 当前项的索引值 |
| del | 删除当前项 | 当前项的索引值 |
| create | 新增选项 | —— |
