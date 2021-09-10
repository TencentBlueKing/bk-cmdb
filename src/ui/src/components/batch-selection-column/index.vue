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
        default: ''
      },
      /**
       * 支持记住上次跨页全选状态
       */
      reserveSelection: {
        type: Boolean,
        default: false
      },
      /**
       * 取消跨页全选是否提示用户
       */
      cancelTooltip: {
        type: Boolean,
        default: true
      },
      /**
       * 取消跨页全选时的提示
       */
      cancelTooltipText: {
        type: String,
        default: '已取消跨页全选'
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
    },
    data() {
      return {
        items: [],
        reservedItems: [], // 记住的选项
        isPageSelected: false, // 全选当页
        isAllSelected: false, // 全选所有
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
          if (this.isAllSelected) {
            this.generatePageSelection(true)
          } else {
            this.generatePageSelection()
          }

          const selectabeItemslLen = this.selectableItems.length
          this.isPageSelected = this.selectedItems.length === selectabeItemslLen && selectabeItemslLen > 0
        } else {
          this.clearSelection()
        }
      },
      emitItems() {
        const { isAllSelected, reserveSelection, reservedItems } = this
        let { selectedItems } = this

        if (isAllSelected) {
          selectedItems = []
        } else if (reserveSelection) {
          selectedItems = reservedItems
        }

        this.$emit('selection-change', selectedItems, isAllSelected)
        this.$emit('update:selectedValue', selectedItems)
        this.$emit('update:allSelected', isAllSelected)
      },
      toggleItemSelection() {
        this.isPageSelected = this.selectedItems.length === this.selectableItems.length
        if (!this.isPageSelected && this.isAllSelected) {
          this.clearSelection()
          this.$bkMessage({
            message: this.cancelTooltipText,
          })
        }
      },
      setReservedItem(item) {
        const findItemIndex = (item) => {
          let itemIndex = -1
          this.reservedItems.forEach((i, index) => {
            if (i[this.rowKey] === item[this.rowKey]) {
              itemIndex = index
            }
          })
          return itemIndex
        }

        const itemIndex = findItemIndex(item)

        if (itemIndex === -1 && item.checked) {
          this.reservedItems.push(item)
        } else if (!item.checked) {
          this.reservedItems.splice(itemIndex, 1)
        }
      },
      handleItemSelectionChange(item) {
        this.toggleItemSelection()
        this.setReservedItem(item)
      },
      generatePageSelection(val) {
        this.items = this.items.map((i, index) => {
          let disabled = false
          let checked = false

          if (this.selectable && typeof this.selectable === 'function') {
            disabled = !this.selectable(i, index)
          }

          if (disabled) {
            return { ...i, checked: false, disabled }
          }

          if (this.reserveSelection) {
            checked = this.reservedItems.some(reservedItem => reservedItem[this.rowKey] === i[this.rowKey])
          }


          if (val !== undefined && typeof val === 'boolean') {
            return { ...i, checked: val }
          }

          return { ...i, checked }
        })
      },
      handlePageSelectionChange(isSelected) {
        this.generatePageSelection(isSelected)

        this.items.forEach((item) => {
          this.setReservedItem(item)
        })

        if (!isSelected && this.isAllSelected) {
          this.clearSelection()
          this.$bkMessage({
            message: this.cancelTooltipText,
          })
        }
      },
      handleAllSelectionChange(isSelected) {
        this.generatePageSelection(isSelected)
        this.isPageSelected = isSelected
        this.reservedItems = []
      },
      clearSelection() {
        this.isPageSelected = false
        this.isAllSelected = false
        this.reservedItems = []
        this.generatePageSelection(false)
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
                disabled={this.pageSelectionDisabled}
                class={{ 'is-total-selected': this.isAllSelected }}
                vModel={this.isPageSelected}
                onChange={this.handlePageSelectionChange}
              ></bk-checkbox>
            </div>
            <template slot="content">
              <bk-checkbox
                disabled={this.allSelectionDisabled}
                class={{ 'is-total-selected': this.isAllSelected }}
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
    },
  }
</script>

<style lang="scss">
.bk-form-checkbox.is-checked.is-total-selected .bk-checkbox{
  background-color: #2dcb56;
  border-color: #2dcb56;
}
</style>
