<script>
    export default {
        data () {
            return {
                list: [{
                    id: 0,
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
                list2: [{
                    id: 1,
                    name: '运动',
                    children: [{
                        id: 10,
                        name: '骑车',
                        alias: 'a'
                    },
                    {
                        id: 12,
                        name: '跳舞',
                        alias: 'b'
                    }]
                }, 
                {
                    id: 2,
                    name: '棋牌',
                    children: [{
                        id: 13,
                        name: '三国杀',
                        alias: 'c'
                    },
                    {
                        id: 14,
                        name: '围棋',
                        alias: 'd'
                    }]
                }],
                provinceList: [{
                    provinceId: 1,
                    provinceName: '广东',
                    cityList: [
                        {
                            cityId: 11,
                            cityName: '广州'
                        },
                        {
                            cityId: 12,
                            cityName: '深圳'
                        }
                    ]
                },
                {
                    provinceId: 2,
                    provinceName: '广西',
                    cityList: [
                        {
                            cityId: 21,
                            cityName: '南宁'
                        },
                        {
                            cityId: 22,
                            cityName: '桂林'
                        }
                    ]
                }],
                cityList: [],
                defaultDemo: {
                    selected: 0
                },
                childDemo: {
                    selected: [10, 12]
                },
                disabledDemo: {
                    selected: 1,
                    disabled: true
                },
                loadingDemo: {
                    selected: -1,
                    isDataLoading: true
                },
                multiDemo: {
                    selected: [2, 3],
                    multiSelect: true
                },
                searchDemo: {
                    selected: [],
                    searchable: true,
                    hasCreateItem: false,
                    searchKey: 'alias'
                },
                customDemo: {
                    selected: 3,
                    hasCreateItem: false,
                    tools: {
                        edit: true,
                        del: true
                    },
                    allowClear: true
                },
                linkDemo: {
                    provinceSelected: -1,
                    citySelected: -1
                }
            }
        },
        methods: {
            selected (id, data) {
                console.log(id)
                console.log(data)
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
            },
            addNewItem () {
                console.log('add')
            },
            selectProvince (index, data) {
                this.cityList = data.cityList
            }
        }
    }
</script>

## Selector 下拉选框

将动作或菜单折叠到下拉菜单中，支持单选和多选

### 基础用法

:::demo 用默认配置初始化组件，参数list和selected为必传

```html
<template>
    <bk-selector
        :list="list"
        :selected.sync="defaultDemo.selected"
        @item-selected="selected">
    </bk-selector>
</template>

<script>
    export default {
        data() {
            return {
                list: [{
                    id: 1,
                    name: '爬山'
                }, 
                {
                    id: 2,
                    name: '跑步'
                }, 
                {
                    id: 3,
                    name: '打球'
                }, 
                {
                    id: 4,
                    name: '跳舞'
                }, 
                {
                    id: 5,
                    name: '健身'
                }, 
                {
                    id: 6,
                    name: '骑车'
                }],
                defaultDemo: {
                    selected: 0
                }
            }
        },
        methods: {
            selected (id, data) {
                console.log(id)
                console.log(data)
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
    <bk-selector
        :list="list"
        :selected.sync="multiDemo.selected"
        :multi-select="multiDemo.multiSelect"
        @item-selected="selected">
    </bk-selector>
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
        },
        methods: {
            selected (id, data) {
                console.log(id)
                console.log(data)
            }
        }
    }
</script>
```
:::


### Loading状态

组件的加载中状态

:::demo 配置禁用状态的参数`isLoading`为`true`，可以显示loading状态，用于接口数据加载慢的情况

```html
<template>
    <bk-selector
        :list="list"
        :selected.sync="loadingDemo.selected"
        :is-loading="loadingDemo.isDataLoading">
    </bk-selector>
</template>

<script>
    export default {
        data() {
            return {
                list: [{
                    id: 1,
                    name: '爬山'
                }, 
                {
                    id: 2,
                    name: '跑步'
                }, 
                {
                    id: 3,
                    name: '打球'
                }, 
                {
                    id: 4,
                    name: '跳舞'
                }, 
                {
                    id: 5,
                    name: '健身'
                }, 
                {
                    id: 6,
                    name: '骑车'
                }],
                loadingDemo: {
                    selected: -1,
                    isDataLoading: true
                }
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
    <bk-selector
        :list="list"
        :selected.sync="disabledDemo.selected"
        :disabled="disabledDemo.disabled">
    </bk-selector>
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



### 带分组

:::demo 可以通过`list`的`children`项来配置分组

```html
<template>
    <bk-selector
        :list="list2"
        :selected.sync="childDemo.selected"
        :multi-select="true"
        @item-selected="selected">
    </bk-selector>
</template>

<script>
    export default {
        data() {
            return {
                list2: [{
                    name: '运动',
                    children: [{
                        id: 10,
                        name: '骑车'
                    },
                    {
                        id: 12,
                        name: '跳舞'
                    }]
                }, 
                {
                    name: '棋牌',
                    children: [{
                        id: 13,
                        name: '三国杀'
                    },
                    {
                        id: 14,
                        name: '围棋'
                    }]
                }],
                childDemo: {
                    selected: [10, 12]
                }
            }
        },
        methods: {
            selected (id, data) {
                console.log(data.name)
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
    <bk-selector
        :list="list"
        :selected.sync="searchDemo.selected"
        :multi-select="true"
        :searchable="searchDemo.searchable">
    </bk-selector>
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
          selected: [],
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

### 扩展

:::demo 可以通过`slot`来自定义模版，会向下拉列表追加

```html
<template>
    <bk-selector
        :list="list2"
        :selected.sync="childDemo.selected"
        :multi-select="true"
        @item-selected="selected">
        <template scope="props">
            <div class="bk-selector-create-item" @click="addNewItem">
                <i class="bk-icon icon-plus-circle"></i>
                <i class="text">新增加</i>
            </div>
        </template>
    </bk-selector>
</template>

<script>
    export default {
        data() {
            return {
                list2: [{
                    name: '运动',
                    children: [{
                        id: 10,
                        name: '骑车'
                    },
                    {
                        id: 12,
                        name: '跳舞'
                    }]
                }, 
                {
                    name: '棋牌',
                    children: [{
                        id: 13,
                        name: '三国杀'
                    },
                    {
                        id: 14,
                        name: '围棋'
                    }]
                }],
                childDemo: {
                    selected: [10, 12]
                }
            }
        },
        methods: {
            selected (id, data) {
                console.log(data.name)
            },
            addNewItem () {
                console.log('add')
            }
        }
    }
</script>
```
:::

### 联动

:::demo 可以多个下拉选框先进行联动

```html
<template>
    <bk-selector
        :style="{width: '150px', display: 'inline-block'}"
        :list="provinceList"
        :setting-key="'provinceId'"
        :display-key="'provinceName'"
        :placeholder="'选择省'"
        :selected.sync="linkDemo.provinceSelected"
        @item-selected="selectProvince">
    </bk-selector>
    <bk-selector
        :style="{width: '150px', display: 'inline-block'}"
        :list="cityList"
        :setting-key="'cityId'"
        :display-key="'cityName'"
        :placeholder="'选择市'"
        :selected.sync="linkDemo.citySelected">
    </bk-selector>
</template>

<script>
    export default {
        data () {
            return {
                provinceList: [{
                    provinceId: 1,
                    provinceName: '广东',
                    cityList: [
                        {
                            cityId: 11,
                            cityName: '广州'
                        },
                        {
                            cityId: 12,
                            cityName: '深圳'
                        }
                    ]
                },
                {
                    provinceId: 2,
                    provinceName: '广西',
                    cityList: [
                        {
                            cityId: 21,
                            cityName: '南宁'
                        },
                        {
                            cityId: 22,
                            cityName: '桂林'
                        }
                    ]
                }],
                cityList: [],
                linkDemo: {
                    provinceSelected: -1,
                    citySelected: -1
                }
            }
        },
        methods: {
            selectProvince (index, data) {
                this.cityList = data.cityList
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
|selected       | 选中的项的index值，支持.sync修饰符，必选，**如果是多选，必须为数组** | Number/Array | —— | —— |
|ext-cls         | 自定义样式的类名，会被加在最外层的DOM节点上 | String | —— | —— |
|create-text     | 下拉菜单中新增项的文字 | String | —— | —— |
|tools          | 待选项右侧的工具按钮，有两个可配置的key：edit和del，默认为两者都显示。| Object/Boolean | —— | [{edit: true, del: true}] |
|placeholder    | 未选选项时，组件显示的文字 | String | —— | '请选择' |
|display-key     | 循环list时，显示字段的key值 | String | —— | 'name' |
|disabled       | 是否禁用组件 | Boolean | —— | false |
|multi-select    | 是否支持多选 | Boolean | —— | false |
|searchable     | 是否支持筛选 | Boolean | —— | false |
|search-key      | 筛选时，搜索的key值 | String | —— | 'name' |
|allow-clear     | 是否可以清除单选时选中的项 | Boolean | —— | false |
|setting-key     | 确定列表数据项的唯一值，默认为`id` | String | —— | 'id' |

### 事件
| 事件名称 | 说明 | 回调参数 |
|---------|------|---------|
| item-selected | 选中选项时触发 | 选中的唯一键和选中的数据对象|
| typing | 输入筛选条件时触发 | 当前输入的值 |
| clear | 单选模式下清空当前选择时触发 | 当前清除的项的索引值 |
| visible-toggle | 下拉框打开/关闭时触发 | 打开为true，关闭为false |
| edit | 编辑当前项 | 当前项的索引值 |
| del | 删除当前项 | 当前项的索引值 |
