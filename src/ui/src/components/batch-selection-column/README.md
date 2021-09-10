# 批量全选选择列

让表格支持批量跨页全选，如果没有跨页全选的需求，建议使用 bk-table 自带的全选列。

## 使用方法

### 基础用法

```html
<bk-table :data="tableData">
  <batch-selection-column :data="tableData" :selected-value.sync="selectedItems" :all-selected="isAllSelected">
</bk-table>

<script>
  export default {
    data: {
      isAllSelected: false,
      selectedItems: [],
      tableData: [
        {
          id: 1,
        }
      ]
    },
  }
</script>
```

### 使用 `selection-change` 事件

```html
<bk-table :data="tableData">
  <batch-selection-column :data="tableData" @selection-change="handleSelectionChange">
</bk-table>

<script>
  export default {
    data: {
      isSelectAll: false,
      selectedItems: [],
      tableData: [
        {
          id: 1,
        }
      ]
    },
    methods: {
      handleSelectionChange(selectedItems, allSelected) {
        this.selectedItems = selectedItems // Array
        this.isSelectAll = allSelected // true or false
      }
    }
  }
</script>
```

### 使用 `reserve-selection` 记住已选择的选项

```html
<bk-input @change="handleSearchValueChange" v-model="searchValue"></bk-input>
<bk-table :data="tableData">
  <batch-selection-column ref="batchSelectionColumn" :data="tableData" row-key="id" reserve-selection @selection-change="handleSelectionChange">
</bk-table>

<script>
  export default {
    data: {
      searchValue: '',
      isSelectAll: false,
      selectedItems: [],
      tableData: [
        {
          id: 1,
        }
      ]
    },
    methods: {
      handleSelectionChange(selectedItems, allSelected) {
        this.selectedItems = selectedItems // Array
        this.isSelectAll = allSelected // true or false
      },
      handleSearchValueChange (){
        // 搜索时清除掉所有选择状态
        this.$refs.batchSelectionColumn.clearSelection()
      },
    }
  }
</script>
```

## 属性

|属性名|描述|类型|默认值|必须|
|-|-|-|-|-|
|data|表格数据|Array|[]|必须|
|selected-value|已选择数据，支持 .sync 修饰符|Array|[]|非必须|
|all-selected|是否跨页全选，支持 .sync 修饰符。跨页全选时，已选择的数据会清空|Boolean|false|非必须|
|selectable|可选状态函数，依据返回的布尔值决定项目是否可选|Function|null|非必须|
|row-key|项目的主键，配合 reserveSelection 使用|String|''|非必须|
|reserve-selection|保存已选数据，跨页仍会保存下来|Boolean|false|非必须|
|cancel-tooltip|取消跨页全选后提示用户|Boolean|true|非必须
|cancel-tooltip-text|取消跨页全选后提示用户的文本|String|已取消跨页全选|非必须|
|page-selection-disabled|全选当页禁用开关|Boolean|false|非必须|
|all-selection-disabled|跨页全选禁用开关|Boolean|false|非必须|

## 事件

|事件名|描述|参数
|-|-|-|
|selection-change|选择状态变化事件|`已选择的项目`,`是否跨页全选`

## 方法

|方法名|描述|
|-|-|
|clearSelection|清除所有选择状态|
