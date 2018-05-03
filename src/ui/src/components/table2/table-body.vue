<template>
    <table class="cc-table-body">
        <colgroup>
            <col v-for="(width, index) in layout.colgroup" :key="index" :width="getColWidth(width, index)">
        </colgroup>
        <tbody>
            <tr v-for="(item, index) in table.list" :key="index">
                <template v-for="(head, index) in table.header">
                    <td v-if="head.type === 'checkbox'" class="body-checkbox">
                        <label :for="getCheckboxId(head)" class="bk-form-checkbox bk-checkbox-small">
                            <input type="checkbox" :id="getCheckboxId(head)">
                        </label>
                    </td>
                    <td v-else>
                        <slot :name="head[table.valueKey]" :item="item">
                            <div>{{item[head[table.valueKey]]}}</div>
                        </slot>
                    </td>
                </template>
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
            this.$slots = this.table.$slots
            this.$scopedSlots = this.table.$scopedSlots
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
                padding: 0 20px;
                border: 1px solid $tableBorderColor;
                border-left: none;
                border-top: none;
                @include ellipsis;
                &.body-checkbox{
                    padding: 0;
                }
                &:last-child{
                    border-right: none;
                }
            }
        }
        tr:last-child{
            td {
                border-bottom: none;
            }
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
</style>