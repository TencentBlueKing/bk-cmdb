<script>
    export default {
        data () {
            return {
                list: [
                    {
                        value: 1,
                        label: '爬山',
                        alias: 'pashan',
                        disabled: true
                    },
                    {
                        value: 2,
                        label: '跑步',
                        alias: 'paobu'
                    },
                    {
                        value: 3,
                        label: '跳舞',
                        alias: 'tiaowu'
                    },
                    {
                        value: 4,
                        label: '健身',
                        alias: 'jianshen'
                    },
                    {
                        value: 5,
                        label: '打球',
                        alias: 'daqiu'
                    }
                ],
                optionDisabledList: [
                    {
                        value: 'climbing',
                        label: '爬山',
                        alias: 'pashan'
                    },
                    {
                        value: 'running',
                        label: '跑步',
                        alias: 'paobu',
                        disabled: true
                    },
                    {
                        value: 'dancing',
                        label: '跳舞',
                        alias: 'tiaowu'
                    },
                    {
                        value: 'gym',
                        label: '健身',
                        alias: 'jianshen'
                    },
                    {
                        value: 'play-ball',
                        label: '打球',
                        alias: 'daqiu'
                    }
                ],
                groupList: [
                    {
                        label: '运动',
                        children: [
                            {
                                label: '骑车',
                                value: 'cycling'
                            },
                            {
                                label: '骑马',
                                value: 'horsing'
                            },
                            {
                                label: '跳舞',
                                value: 'dancing'
                            }
                        ]
                    },
                    {
                        label: '棋牌',
                        children: [
                            {
                                label: '三国杀',
                                value: 'shanguosha'
                            },
                            {
                                label: '围棋',
                                value: 'weiqi'
                            }
                        ]
                    }
                ],
                multiple: true,
                filterable: true,
                defaultMultiple: '',
                disabled: true,
                defaultModel: 3,
                disabledModel: '',
                disabledOptionModel: '',
                multipleModel: '',
                groupModel: '',
                filterModel: '',
                extModel: ''
            }
        }
    }
</script>

## Select 选择器

模拟原生select选择器

### 基础用法

:::demo 用默认配置初始化组件，`bk-select`中的`selected`和`bk-select-option`中的`value`为必选项

```html
<template>
    <bk-select
        :selected.sync="defaultModel"
        :list="list">
        <bk-select-option
            v-for="option of list"
            :key="option.value"
            :value="option.value"
            :label="option.label">
        </bk-select-option>
    </bk-select>
</template>

<script>
    export default {
        data () {
            return {
                list: [
                    {
                        value: 'climbing',
                        label: '爬山',
                        alias: 'pashan'
                    },
                    {
                        value: 'running',
                        label: '跑步',
                        alias: 'paobu'
                    },
                    {
                        value: 'dancing',
                        label: '跳舞',
                        alias: 'tiaowu'
                    },
                    {
                        value: 'gym',
                        label: '健身',
                        alias: 'jianshen'
                    },
                    {
                        value: 'play-ball',
                        label: '打球',
                        alias: 'daqiu'
                    }
                ],
                defaultModel: ''
            }
        }
    }
</script>
```
:::

### 禁用状态

组件不可用

:::demo 在`bk-select`中配置`disabled`参数

```html
<template>
    <bk-select
        :selected="disabledModel"
        :disabled="disabled">
        <bk-select-option
            v-for="option of list"
            :key="option.value"
            :value="option.value"
            :label="option.label">
        </bk-select-option>
    </bk-select>
</template>

<script>
    export default {
        data () {
            return {
                list: [
                    {
                        value: 'climbing',
                        label: '爬山',
                        alias: 'pashan'
                    },
                    {
                        value: 'running',
                        label: '跑步',
                        alias: 'paobu'
                    },
                    {
                        value: 'dancing',
                        label: '跳舞',
                        alias: 'tiaowu'
                    },
                    {
                        value: 'gym',
                        label: '健身',
                        alias: 'jianshen'
                    },
                    {
                        value: 'play-ball',
                        label: '打球',
                        alias: 'daqiu'
                    }
                ],
                defaultModel: ''
            }
        }
    }
</script>
```
:::

### 有禁用选项

:::demo 在`bk-select-option`中给特定的项配置`disabled`参数。若之后因业务需要改变禁用选项的状态，请在选项中预先配置`{disabled: false}`以提供数据绑定。

```html
<template>
    <bk-select
        :selected="disabledOptionModel">
        <bk-select-option
            v-for="option of optionDisabledList"
            :key="option.value"
            :value="option.value"
            :label="option.label"
            :disabled="option.disabled">
        </bk-select-option>
    </bk-select>
</template>

<script>
    export default {
        data () {
            return {
                optionDisabledList: [
                    {
                        value: 'climbing',
                        label: '爬山',
                        alias: 'pashan'
                    },
                    {
                        value: 'running',
                        label: '跑步',
                        alias: 'paobu',
                        disabled: true
                    },
                    {
                        value: 'dancing',
                        label: '跳舞',
                        alias: 'tiaowu'
                    },
                    {
                        value: 'gym',
                        label: '健身',
                        alias: 'jianshen'
                    },
                    {
                        value: 'play-ball',
                        label: '打球',
                        alias: 'daqiu'
                    }
                ]
            }
        }
    }
</script>
```
:::

### 多选

:::demo 在`bk-select`上配置`multiple`参数
```html
<template>
    <bk-select
        :selected.sync="multipleModel"
        :multiple="multiple">
        <bk-select-option
            v-for="option of list"
            :disabled="option.disabled"
            :key="option.value"
            :value="option.value"
            :label="option.label">
        </bk-select-option>
    </bk-select>
</template>

<script>
    export default {
        data () {
            return {
                list: [
                    {
                        value: 'climbing',
                        label: '爬山',
                        alias: 'pashan'
                    },
                    {
                        value: 'running',
                        label: '跑步',
                        alias: 'paobu'
                    },
                    {
                        value: 'dancing',
                        label: '跳舞',
                        alias: 'tiaowu'
                    },
                    {
                        value: 'gym',
                        label: '健身',
                        alias: 'jianshen'
                    },
                    {
                        value: 'play-ball',
                        label: '打球',
                        alias: 'daqiu'
                    }
                ],
                defaultModel: '',
                multiple: true
            }
        }
    }
</script>
```
:::

### 分组

:::demo 使用`bk-option-group`作为分组

```html
<template>
    <bk-select
        :selected="groupModel">
        <bk-option-group
            v-for="group of groupList"
            :label="group.label"
            :key="group.label">
            <bk-select-option
                v-for="option of group.children"
                :key="option.value"
                :value="option.value"
                :label="option.label">
            </bk-select-option>
        </bk-option-group>
    </bk-select>
</template>

<script>
    export default {
        data () {
            return {
                groupList: [
                    {
                        label: '运动',
                        children: [
                            {
                                label: '骑车',
                                value: 'cycling'
                            },
                            {
                                label: '跳舞',
                                value: 'dancing'
                            }
                        ]
                    },
                    {
                        label: '棋牌',
                        children: [
                            {
                                label: '三国杀',
                                value: 'shanguosha'
                            },
                            {
                                label: '围棋',
                                value: 'weiqi'
                            }
                        ]
                    }
                ]
            }
        }
    }
</script>
```
:::

### 过滤列表
当前备选项过多时可以过滤

:::demo
```html
<template>
    <bk-select
        :selected="filterModel"
        :filterable="filterable">
        <bk-select-option
            v-for="option of list"
            :key="option.value"
            :value="option.value"
            :label="option.label">
        </bk-select-option>
    </bk-select>
</template>
```
:::

### 扩展

可增加自定义内容

:::demo 可在当前列表的前(配置`slot`为`pre-slot`)/后(配置`slot`为`post-slot`)增加自定义配置项
```html
<template>
    <bk-select
        :selected="extModel">
        <bk-option-group
            v-for="group of groupList"
            :label="group.label"
            :key="group.label">
            <bk-select-option
                v-for="option of group.children"
                :key="option.value"
                :value="option.value"
                :label="option.label">
            </bk-select-option>
        </bk-option-group>
        <template slot="post-ext">
            <p class="demo-select-create">
                <i class="bk-icon icon-plus-circle"></i>
                新增项
            </p>
        </template>
    </bk-select>
</template>

<script>
    export default {
        data () {
            return {
                groupList: [
                    {
                        label: '运动',
                        children: [
                            {
                                label: '骑车',
                                value: 'cycling'
                            },
                            {
                                label: '跳舞',
                                value: 'dancing'
                            }
                        ]
                    },
                    {
                        label: '棋牌',
                        children: [
                            {
                                label: '三国杀',
                                value: 'shanguosha'
                            },
                            {
                                label: '围棋',
                                value: 'weiqi'
                            }
                        ]
                    }
                ]
            }
        }
    }
</script>
```
:::

### select属性
|参数           | 说明    | 类型      | 可选值       | 默认值   |
|---------------|-------- |---------- |-------------  |-------- |
| selected | 当前选中的值，必选 | Any | —— | —— |
| value-key | 当`selected`为对象时必须传入，组件用该参数作为匹配标准 | String | —— | 'value' |
| placeholder | 未选择项时的提示语 | String | —— | '请选择' |
| disabled | 是否禁用组件 | Boolean | —— | false |
| multiple | 是否开启多选 | Boolean | —— | false |
| filterable | 是否开启列表过滤 | Boolean | —— | false |
| filter-fn | 自定义过滤方法，参数为所有列表项 | Function | —— |

### select事件
| 事件名称 | 说明 | 回调参数 |
|---------|------|---------|
| on-selected | 点击选项后触发 | 第一个参数在单选时为当前选中项的数据，多选时为当前选中的所有option的数组；第二个参数为选中项在列表中的索引值 |
| on-toggle | 下拉框打开/关闭时触发 | 打开为`true`，关闭为`false` |
| on-filter | 输入过滤条件时触发 | 第一个参数为当前输入的值，第二个参数为所有option项|

### select-option
|参数           | 说明    | 类型      | 可选值       | 默认值   |
|---------------|-------- |---------- |-------------  |-------- |
| value | 当前项在列表中的标识 | Any | —— | —— |
| label | 当前项显示的文字 | String/Number | —— | —— |
| disabled | 当前项是否禁用 | Boolean | —— | false |
