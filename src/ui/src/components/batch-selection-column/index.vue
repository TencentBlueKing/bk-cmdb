<template>
  <bk-table-column
    ref="batchSelectionColumn"
    class-name="batch-selection-column"
    :render-header="columnHeader"
    v-bind="$attrs"
    v-on="$listeners"
  >
    <template #default="{ $index }">
      <bk-checkbox
        v-if="items[$index]"
        :disabled="items[$index].disabled"
        @change="handleItemSelectionChange(items[$index])"
        v-model="items[$index].checked"
      ></bk-checkbox>
    </template>
  </bk-table-column>
</template>
<script>
  import cloneDeep from 'lodash/cloneDeep'

  export default {
    name: 'BatchSelectionColumn',
    props: {
      /**
       * 可选数据
       */
      data: {
        type: Array,
        required: true,
        default: () => [],
      },
      /**
       * 已选择数据
       */
      selectedValue: {
        type: Array,
        default: () => [],
      },
      /**
       * 是否跨页全选
       */
      allSelected: {
        type: Boolean,
        default: false,
      },
      /**
       * 控制项目是否可选
       * @returns {Boolean} 可选状态
       */
      selectable: {
        type: Function,
        default: null,
      },
      /**
       * 行主键，用于记住跨页选择
       */
      rowKey: {
        type: String,
        required: true,
        default: ''
      },
      /**
       * 支持记住上次跨页全选状态
       */
      reserveSelection: {
        type: Boolean,
        default: true
      },
      /**
       * 全选当页禁用开关
       */
      pageSelectionDisabled: {
        type: Boolean,
        default: false
      },
      /**
       * 跨页全选禁用开关
       */
      allSelectionDisabled: {
        type: Boolean,
        default: false
      },
      /**
       * 是否支持半选状态
       */
      indeterminate: {
        type: Boolean,
        default: false
      }
    },
    data() {
      return {
        items: [],
        reservedSelectedItems: [], // 记住的已选项
        reservedUnselectedItems: [], // 全选时，记住的未选项
        isPageSelected: false, // 全选当页
        onCrossPageMode: false, // 是否在跨页全选模式
        isAllSelected: false, // 全选所有
        pageSelectionIndeterminate: false,
        allSelectionIndeterminate: false,
      }
    },
    computed: {
      selectedItems() {
        if (!this.items?.length) return []
        return this.items.filter(i => i.checked)
      },
      selectableItems() {
        if (!this.items?.length) return []
        return this.items.filter(i => !i.disabled)
      }
    },
    watch: {
      data: {
        immediate: true,
        handler() {
          this.initItems()
        },
      },
      items: {
        deep: true,
        handler() {
          this.emitItems()
        },
      },
    },
    methods: {
      initItems() {
        this.items = cloneDeep(this.data)

        if (this.reserveSelection && this.rowKey) {
          this.generateItemSelection()

          this.generatePageSelection()
        } else {
          this.clearSelection()
        }
      },
      emitItems() {
        const { onCrossPageMode, reserveSelection, reservedSelectedItems } = this
        let { selectedItems } = this

        if (onCrossPageMode) {
          selectedItems = []
        } else if (reserveSelection) {
          selectedItems = reservedSelectedItems
        }

        this.$emit('selection-change', selectedItems, onCrossPageMode, this.reservedUnselectedItems)
        this.$emit('update:selectedValue', selectedItems)
        this.$emit('update:unselectedValue', this.reservedUnselectedItems)
        this.$emit('update:allSelected', onCrossPageMode)
      },
      setReservedItem(arr, item, checked) {
        const findItemIndex = (item) => {
          let itemIndex = -1
          arr.forEach((i, index) => {
            if (i[this.rowKey] === item[this.rowKey]) {
              itemIndex = index
            }
          })
          return itemIndex
        }

        const itemIndex = findItemIndex(item)

        if (itemIndex === -1 && item.checked === checked) {
          arr.push(item)
        } else if (item.checked !== checked) {
          arr.splice(itemIndex, 1)
        }
      },
      handleItemSelectionChange(item) {
        if (this.onCrossPageMode) {
          this.setReservedItem(this.reservedUnselectedItems, item, false)
        } else {
          this.setReservedItem(this.reservedSelectedItems, item, true)
        }

        this.generatePageSelection()
      },
      handlePageSelectionChange(isSelected) {
        if (this.indeterminate && this.pageSelectionIndeterminate) {
          this.generateItemSelection(false)
        } else {
          this.generateItemSelection(isSelected)
        }

        this.items.forEach((item) => {
          if (this.onCrossPageMode) {
            this.setReservedItem(this.reservedUnselectedItems, item, false)
          } else {
            this.setReservedItem(this.reservedSelectedItems, item, true)
          }
        })

        this.$nextTick(() => {
          this.generatePageSelection()
        })
      },
      handleAllSelectionChange(isSelected) {
        this.generateItemSelection(isSelected)

        if (this.indeterminate && this.allSelectionIndeterminate) {
          this.onCrossPageMode = false
          this.clearSelection()
        } else {
          this.onCrossPageMode = isSelected
        }

        this.reservedSelectedItems = []
        this.reservedUnselectedItems = []

        this.$nextTick(() => {
          this.generatePageSelection()
        })
      },
      columnHeader() {
        return (
        <div class="batch-selection-label">
          <bk-popover
            placement="right"
            theme="light"
            arrow={false}
            size="regular"
          >
            <div>
              <bk-checkbox
                indeterminate={this.pageSelectionIndeterminate}
                class={{ 'is-total-selected': this.onCrossPageMode, 'page-select-checkbox': true }}
                disabled={this.pageSelectionDisabled}
                vModel={this.isPageSelected}
                onChange={this.handlePageSelectionChange}
              ></bk-checkbox>
            </div>
            <template slot="content">
              <bk-checkbox
                indeterminate={this.allSelectionIndeterminate}
                disabled={this.allSelectionDisabled}
                class={{ 'is-total-selected': this.onCrossPageMode, 'all-select-checkbox': true }}
                vModel={this.isAllSelected}
                onChange={this.handleAllSelectionChange}
              >
                跨页全选
              </bk-checkbox>
            </template>
          </bk-popover>
        </div>
        )
      },
      // 生成页面全选状态
      generatePageSelection() {
        const selectabeItemslLen = this.selectableItems.length
        const selectedItemsLen = this.selectedItems.length

        this.isPageSelected = selectedItemsLen === selectabeItemslLen && selectabeItemslLen > 0

        if (this.indeterminate) {
          this.isAllSelected = this.reservedUnselectedItems.length === 0 && this.onCrossPageMode
          this.pageSelectionIndeterminate = selectedItemsLen > 0 && !this.isPageSelected
          this.allSelectionIndeterminate = this.reservedUnselectedItems.length > 0 && !this.isAllSelected
        } else {
          this.isAllSelected = this.onCrossPageMode
        }
      },
      // 生成单个项目选择状态
      generateItemSelection(currentChecked) {
        this.items = this.items.map((i, index) => {
          // 如果传入了当前选中值，则让所有选项变为当前选中值
          if (currentChecked !== undefined && typeof currentChecked === 'boolean') {
            return { ...i, checked: currentChecked }
          }

          // 当跨页全选时，默认为选中状态，非跨页全选时则默认为未选中状态
          let checked = this.onCrossPageMode

          // 不可选的选项一律置灰
          if (this.selectable && typeof this.selectable === 'function') {
            const disabled = !this.selectable(i, index)
            if (disabled) {
              return { ...i, checked: false, disabled }
            }
          }

          // 记住选择状态
          if (this.reserveSelection) {
            if (this.onCrossPageMode) {
              // 跨页全选时，记住没有选择的项目
              checked = !this.reservedUnselectedItems
                .some(unselectedItem => unselectedItem[this.rowKey] === i[this.rowKey])
            } else {
              // 非跨页全选时，记住已选择的项目
              checked = this.reservedSelectedItems.some(selectedItem => selectedItem[this.rowKey] === i[this.rowKey])
            }
          }

          return { ...i, checked }
        })
      },
      // 清除所有选择状态
      clearSelection() {
        this.isPageSelected = false
        this.onCrossPageMode = false

        if (this.indeterminate) {
          this.pageSelectionIndeterminate = false
          this.allSelectionIndeterminate = false
        }

        this.reservedSelectedItems = []
        this.reservedUnselectedItems = []
        this.generateItemSelection(false)
      },
    },
  }
</script>
<style lang="scss">
.bk-form-checkbox.is-checked.is-total-selected .bk-checkbox,
.bk-form-checkbox.is-indeterminate.is-total-selected .bk-checkbox {
  background-color: #2dcb56;
  border-color: #2dcb56;
}
</style>
