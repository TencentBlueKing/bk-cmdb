<template>
    <table :class="['cc-table-body', {'row-border': table.rowBorder, 'col-border': table.colBorder, 'stripe': table.stripe}]">
        <colgroup>
            <col v-for="(width, index) in layout.colgroup" :key="index" :width="width">
        </colgroup>
        <tbody v-show="table.list.length">
            <tr v-for="(item, rowIndex) in table.list" 
                :key="rowIndex"
                @click="handleRowClick(item, rowIndex)"
                @mouseover="handleRowMouseover($event, item, rowIndex)"
                @mouseout="handleRowMouseout($event, item, rowIndex)">
                <template v-for="(head, colIndex) in table.header">
                    <td v-if="head.type === 'checkbox' && !table.$scopedSlots[head[table.valueKey]]" class="data-content checkbox-content" @click.stop :key="colIndex">
                        <label class="bk-form-checkbox bk-checkbox-small" :for="getCheckboxId(head, rowIndex)">
                            <input type="checkbox"
                                :id="getCheckboxId(head, rowIndex)"
                                :checked="checked.indexOf(item[head[table.valueKey]]) !== -1"
                                @change="handleRowCheck(item, item[head[table.valueKey]], rowIndex)">
                        </label>
                    </td>
                    <td is="data-content" :class="['data-content', {'checkbox-content': head.type === 'checkbox'}]" v-else
                        :key="colIndex"
                        :item="item"
                        :head="head"
                        :layout="layout"
                        :rowIndex="rowIndex"
                        :colIndex="colIndex">
                    </td>
                </template>
            </tr>
        </tbody>
        <tbody v-if="!table.list.length">
            <tr>
                <td is="data-empty" class="data-empty" align="center"
                    :colspan="table.header.length"
                    :style="{height: emptyHeight}" 
                    :layout="layout">
                </td>
            </tr>
        </tbody>
    </table>
</template>

<script>
    export default {
        props: {
            layout: Object
        },
        data () {
            return {}
        },
        computed: {
            table () {
                return this.layout.table
            },
            checked () {
                return this.table.checked
            },
            emptyHeight () {
                return this.table.emptyHeight + 'px'
            }
        },
        methods: {
            getCheckboxId (head, rowIndex) {
                return `table-${this.layout.id}-body-${head[this.table.valueKey]}-checkbox-${rowIndex}`
            },
            getColWidth (width, index) {
                let total = this.layout.colgroup.length
                if ((index === total - 1) && this.layout.scrollY) {
                    return width - this.table.gutterWidth
                }
                return width
            },
            handleRowCheck (item, value, rowIndex) {
                let checked = [...this.checked]
                const index = checked.indexOf(value)
                if (this.table.multipleCheck) {
                    if (index === -1) {
                        checked.push(value)
                    } else {
                        checked.splice(index, 1)
                    }
                } else {
                    checked = index === -1 ? [value] : []
                }
                this.table.$emit('update:checked', checked)
                this.table.$emit('handleRowCheck', value, rowIndex)
            },
            handleRowClick (item, rowIndex) {
                this.table.$emit('handleRowClick', item, rowIndex)
            },
            handleRowMouseover (event, item, rowIndex) {
                const rowHoverColor = this.table.rowHoverColor
                if (rowHoverColor) {
                    event.currentTarget.style.backgroundColor = rowHoverColor
                }
            },
            handleRowMouseout (event, item, rowIndex) {
                const rowHoverColor = this.table.rowHoverColor
                if (rowHoverColor) {
                    event.currentTarget.style.backgroundColor = 'inherit'
                }
            }
        },
        components: {
            'data-content': {
                props: ['head', 'item', 'layout', 'rowIndex'],
                render (h) {
                    const table = this.layout.table
                    const column = this.head[table.valueKey]
                    const defaultConfig = {
                        on: {
                            click: this.handleCellClick
                        }
                    }
                    if (typeof table.renderCell === 'function') {
                        return h('td', defaultConfig, table.renderCell(this.item, this.head, this.layout))
                    } else if (table.$scopedSlots[column]) {
                        return h('td', defaultConfig, table.$scopedSlots[column]({item: this.item, rowIndex: this.rowIndex, colIndex: this.colIndex, layout: this.layout}))
                    } else {
                        return h('td', Object.assign({}, defaultConfig, {attrs: {title: this.item[this.head[table.valueKey]]}}), this.item[this.head[table.valueKey]])
                    }
                },
                methods: {
                    handleCellClick () {
                        this.layout.table.$emit('handleCellClick', this.item, this.rowIndex, this.colIndex)
                    }
                }
            },
            'data-empty': {
                props: ['layout'],
                render (h) {
                    const dataEmptySlot = this.layout.table.$slots['data-empty']
                    if (dataEmptySlot) {
                        return h('td', {}, dataEmptySlot)
                    } else {
                        return h('td', {}, this.$t("Common['暂时没有数据']"))
                    }
                }
            }
        }
    }
</script>

<style lang="scss" scoped>
    .cc-table-body{
        color: $textColor;
        text-align: left;
        border-collapse: separate;
        border-spacing: 0;
        table-layout: fixed;
        tr {
            td {
                height: 40px;
                cursor: pointer;
                @include ellipsis;
            }
        }
    }
    .cc-table-body.row-border {
        tr {
            td {
                border-bottom: 1px solid $tableBorderColor;
            }
        }
        tr:last-child{
            td {
                border-bottom: none;
            }
        }
    }
    .cc-table-body.col-border{
        tr {
            td {
                border-right: 1px solid $tableBorderColor;
                &:last-child{
                    border-right: none;
                }
            }
        }
    }
    .cc-table-body.stripe {
        tr:nth-child(2n) {
            background-color: #f1f7ff;
        }
    }
    .bk-form-checkbox{
        display: block;
        width: 100%;
        height: 100%;
        text-align: center;
        padding: 0;
        cursor: pointer;
        &:before{
            content: '';
            width: 0;
            height: 100%;
            display: inline-block;
            vertical-align: middle;
        }
        input[type='checkbox'] {
            display: inline-block;
            vertical-align: middle;
        }
    }
    .data-content{
        padding: 0 16px;
        font-size: 12px;
        @include ellipsis;
        &.checkbox-content{
            padding: 0;
            height: 100%;
        }
    }
    .data-empty{
        font-size: 12px;
    }
</style>