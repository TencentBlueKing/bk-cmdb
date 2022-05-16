<!--
 * Tencent is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2017-2022 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
-->

<template>
  <bk-table-column
    ref="batchSelectionColumn"
    :class-name="`batch-selection-column ${crossPage ? 'is-cross-page' : ''}`"
    :render-header="columnHeader"
    v-bind="$attrs"
    v-on="$listeners"
  >
    <template #default="{ $index }">
      <bk-checkbox
        v-if="rows[$index]"
        :disabled="rows[$index].disabled"
        @change="handleRowSelectionChange(rows[$index])"
        v-model="rows[$index].checked"
        @click.native.stop
      ></bk-checkbox>
    </template>
  </bk-table-column>
</template>
<script>
  import cloneDeep from 'lodash/cloneDeep'
  import safeGet from 'lodash/get'

  export default {
    name: 'BatchSelectionColumn',
    props: {
      /**
       * 是否开启跨页全选，默认开启
       */
      crossPage: {
        type: Boolean,
        default: true,
      },
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
      selectedRows: {
        type: Array,
        default: () => [],
      },
      /**
       * 反选数据
       */
      unselectedRows: {
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
       * 支持记住选择状态
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
        default: true
      },
      /**
       * 全量数据，用于前端分页时选择所有数据，传入 fullData 则会只是用此数据作为全量数据
       */
      fullData: {
        type: Array,
        default: null
      },
    },
    data() {
      return {
        rows: [],
        reservedSelectedRows: [], // 记住的已选项
        reservedUnselectedRows: [], // 全选时，记住的未选项
        isPageSelected: false, // 全选当页
        onCrossPageMode: false, // 是否在跨页全选模式
        isAllSelected: false, // 全选所有
        pageSelectionIndeterminate: false,
        allSelectionIndeterminate: false,
      }
    },
    computed: {
      innerSelectedRows() {
        if (!this.rows?.length) return []
        return this.rows.filter(i => i.checked)
      },
      selectableRows() {
        if (!this.rows?.length) return []
        return this.rows.filter(i => !i.disabled)
      },
      fullDataAllDisabled() {
        if (this.fullData && this.selectable) {
          return this.fullData.every((i, index) => {
            if (this.selectable) return !this.selectable(i, index)
            return true
          })
        }
        return false
      },
    },
    watch: {
      data: {
        immediate: true,
        deep: true,
        handler() {
          this.initRows()
        },
      },
      rows: {
        deep: true,
        handler() {
          this.emitRows()
        },
      },
      selectedRows: {
        immediate: true,
        handler(val) {
          if (this.reserveSelection) {
            this.reservedSelectedRows = cloneDeep(val)
          }
        }
      }
    },
    methods: {
      initRows() {
        this.rows = cloneDeep(this.data)

        if (this.reserveSelection && this.rowKey) {
          this.generateRowSelection()
          this.generatePageSelection()
        } else {
          this.clearSelection()
        }
      },
      emitRows() {
        const { onCrossPageMode, reserveSelection, reservedSelectedRows, reservedUnselectedRows } = this
        let { innerSelectedRows } = this

        if (onCrossPageMode) {
          if (this.fullData) {
            const unselectedRowKeys = reservedUnselectedRows.map(r => safeGet(r, this.rowKey))
            innerSelectedRows = this.fullData
              .filter((i, index) => {
                const isSelected =  !unselectedRowKeys.includes(safeGet(i, this.rowKey))
                if (this.selectable) {
                  return this.selectable(i, index) && isSelected
                }
                return isSelected
              })
          } else {
            innerSelectedRows = []
          }
        } else if (reserveSelection) {
          innerSelectedRows = reservedSelectedRows
        }

        this.$emit('selection-change', innerSelectedRows, reservedUnselectedRows, onCrossPageMode)
        this.$emit('update:selectedRows', innerSelectedRows)
        this.$emit('update:unselectedRows', reservedUnselectedRows)
        this.$emit('update:allSelected', onCrossPageMode)
      },
      setReservedRow(arr, row, checked) {
        if (!this.reserveSelection) return false

        const findRowIndex = (row) => {
          let rowIndex = -1

          arr.forEach((i, index) => {
            if (safeGet(i, this.rowKey) === safeGet(row, this.rowKey)) {
              rowIndex = index
            }
          })

          return rowIndex
        }

        const rowIndex = findRowIndex(row)

        if (rowIndex === -1 && row.checked === checked) {
          arr.push(row)
        } else if (row.checked !== checked && !row.disabled) {
          arr.splice(rowIndex, 1)
        }
      },
      handleRowSelectionChange(row) {
        if (this.onCrossPageMode) {
          this.setReservedRow(this.reservedUnselectedRows, row, false)
        } else {
          this.setReservedRow(this.reservedSelectedRows, row, true)
        }
        this.generatePageSelection()
      },
      handlePageSelectionChange(isSelected) {
        if (this.indeterminate && this.pageSelectionIndeterminate) {
          this.generateRowSelection(true)
        } else {
          this.generateRowSelection(isSelected)
        }

        this.rows.forEach((row) => {
          if (this.onCrossPageMode) {
            this.setReservedRow(this.reservedUnselectedRows, row, false)
          } else {
            this.setReservedRow(this.reservedSelectedRows, row, true)
          }
        })

        this.$nextTick(() => {
          this.generatePageSelection()
        })
      },
      handleAllSelectionChange(isSelected) {
        this.generateRowSelection(isSelected)

        if (this.indeterminate && this.allSelectionIndeterminate) {
          this.onCrossPageMode = true
        } else {
          this.onCrossPageMode = isSelected
        }

        this.reservedSelectedRows = []
        this.reservedUnselectedRows = []

        this.$nextTick(() => {
          this.generatePageSelection()
        })
      },
      columnHeader() {
        const pageSelection = <bk-checkbox
                indeterminate={this.pageSelectionIndeterminate}
                class={[{ 'is-total-selected': this.onCrossPageMode }, 'page-select-checkbox']}
                disabled={this.pageSelectionDisabled || this.selectableRows.length === 0}
                vModel={this.isPageSelected}
                onChange={this.handlePageSelectionChange}
              ></bk-checkbox>

        if (!this.crossPage) {
          return  <div class="batch-selection-label">{pageSelection}</div>
        }

        return (
        <div class="batch-selection-label">
          <bk-popover
            placement="right"
            theme="light"
            arrow={false}
            size="regular"
            tippy-options={{
              distance: -10,
              duration: 100
            }}
          >
            {pageSelection}
            <template slot="content">
              <bk-checkbox
                indeterminate={this.allSelectionIndeterminate}
                disabled={this.allSelectionDisabled || this.fullDataAllDisabled}
                class={[{ 'is-total-selected': this.onCrossPageMode }, 'all-select-checkbox']}
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
        const selectabeRowslLen = this.selectableRows.length
        const selectedRowsLen = this.innerSelectedRows.length

        this.isPageSelected = selectedRowsLen === selectabeRowslLen && selectabeRowslLen > 0

        if (this.indeterminate) {
          this.isAllSelected = this.reservedUnselectedRows.length === 0 && this.onCrossPageMode
          this.pageSelectionIndeterminate = selectedRowsLen > 0 && !this.isPageSelected
          this.allSelectionIndeterminate = this.reservedUnselectedRows.length > 0 && !this.isAllSelected
        } else {
          this.isAllSelected = this.onCrossPageMode
        }
      },
      // 生成单个项目选择状态
      generateRowSelection(currentChecked) {
        this.rows = this.rows.map((i, index) => {
          // 如果传入了当前选中值，则让所有选项变为当前选中值
          if (currentChecked !== undefined && typeof currentChecked === 'boolean') {
            return { ...i, checked: currentChecked && !i.disabled }
          }

          // 当跨页全选时，默认为选中状态，非跨页全选时则默认为未选中状态
          let checked = this.onCrossPageMode

          // 不可选的选项一律置灰
          if (this.selectable) {
            const disabled = !this.selectable(i, index)
            if (disabled) {
              return { ...i, checked: false, disabled }
            }
          }

          // 记住选择状态
          if (this.reserveSelection) {
            if (this.onCrossPageMode) {
              // 跨页全选时，记住没有选择的项目
              checked = !this.reservedUnselectedRows
                .some(unselectedRow => safeGet(unselectedRow, this.rowKey) === safeGet(i, this.rowKey))
            } else {
              // 非跨页全选时，记住已选择的项目
              checked = this.reservedSelectedRows
                .some(selectedRow => safeGet(selectedRow, this.rowKey) === safeGet(i, this.rowKey))
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

        this.reservedSelectedRows = []
        this.reservedUnselectedRows = []
        this.generateRowSelection(false)
      },
      toggleRowSelection(row) {
        const currentRow = this.rows.find(r => safeGet(row, this.rowKey) === safeGet(r, this.rowKey))
        currentRow.checked = !currentRow.checked
        this.handleRowSelectionChange(currentRow)
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

  th.batch-selection-column.is-cross-page .cell {
    padding-right: 0;
    padding-left: 0;

    .bk-tooltip-ref{
      padding-right: 15px;
      padding-left: 15px;
    }
  }
</style>

