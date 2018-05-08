<template>
    <table :class="['cc-table-body', {'row-border': table.rowBorder, 'col-border': table.colBorder, 'stripe': table.stripe}]">
        <colgroup>
            <col v-for="(width, index) in layout.colgroup" :key="index" :width="getColWidth(width, index)">
        </colgroup>
        <tbody v-show="table.list.length">
            <tr v-for="(item, index) in table.list" :key="index" @click="handleRowClick(item)">
                <template v-for="(head, index) in table.header">
                    <td v-if="head.type === 'checkbox'" class="body-checkbox">
                        <label :for="getCheckboxId(head)" class="bk-form-checkbox bk-checkbox-small">
                            <input type="checkbox" :id="getCheckboxId(head)">
                        </label>
                    </td>
                    <td v-else>
                        <td-content :item="item" :head="head" :layout="layout"></td-content>
                    </td>
                </template>
            </tr>
        </tbody>
        <tbody v-show="!table.list.length">
            <tr>
                <td :colspan="table.header.length" align="center">
                    <data-empty class="data-empty" :layout="layout"></data-empty>
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
            }
        },
        created () {
            this.$slots = this.layout.table.$slots
            this.$scopedSlots = this.layout.table.$scopedSlots
        },
        methods: {
            getCheckboxId (head) {
                return `table-${this.layout.id}-body-${head[this.table.valueKey]}-checkbox`
            },
            getColWidth (width, index) {
                let total = this.layout.colgroup.length
                if ((index === total - 1) && this.layout.scrollY) {
                    return width - this.table.gutterWidth
                }
                return width
            },
            handleRowClick (item) {
                this.table.$emit('handleRowClick', item)
            }
        },
        components: {
            'td-content': {
                props: ['head', 'item', 'layout'],
                render (h) {
                    const table = this.layout.table
                    let scopedSlots = table.$scopedSlots[this.head[table.valueKey]]
                    if (scopedSlots) {
                        return scopedSlots({item: this.item})
                    } else {
                        return h('div', this.item[this.head[table.valueKey]])
                    }
                }
            },
            'data-empty': {
                props: ['layout'],
                render (h) {
                    const dataEmpty = this.layout.table.$slots['data-empty']
                    if (dataEmpty) {
                        return dataEmpty()
                    } else {
                        return h('div', this.$t("Common['暂时没有数据']"))
                    }
                },
                mounted () {
                    const bodyWrapperMaxHeight = parseInt(this.layout.table.$refs.bodyWrapper.style.maxHeight, 10)
                    this.$el.style.height = bodyWrapperMaxHeight ? bodyWrapperMaxHeight + 'px' : '220px'
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
                padding: 0 16px;
                cursor: pointer;
                @include ellipsis;
                &.body-checkbox{
                    padding: 0;
                }
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
    .data-empty{
        font-size: 12px;
        &:before{
            content: '';
            width: 0;
            height: 100%;
            display: inline-block;
            vertical-align: middle;
        }
    }
</style>