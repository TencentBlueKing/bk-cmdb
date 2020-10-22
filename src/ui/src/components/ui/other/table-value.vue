<template>
    <bk-table v-if="displayType === 'table'" class="table-value" :data="list">
        <bk-table-column v-for="col in header"
            :key="col.bk_property_id"
            :prop="col.bk_property_id"
            :label="col.bk_property_name"
            :width="col.bk_property_type === 'bool' ? '90px' : ''"
            show-overflow-tooltip>
            <template slot-scope="{ row }">
                <cmdb-property-value
                    v-bk-overflow-tips
                    :value="row[col['bk_property_id']]"
                    :property="col">
                </cmdb-property-value>
            </template>
        </bk-table-column>
        <div slot="empty">
            <span>{{$t('暂无数据')}}</span>
        </div>
    </bk-table>
    <div class="table-cell-value" v-else>
        <vnodes :vnode="getCellValue()"></vnodes>
    </div>
</template>

<script>
    export default {
        components: {
            vnodes: {
                functional: true,
                render: (h, ctx) => ctx.props.vnode
            }
        },
        props: {
            value: {
                type: [Array, String], // String是为了兼容后台数据未给默认值的情况
                default: () => ([])
            },
            property: {
                type: Object,
                default: () => ({})
            },
            showOn: {
                type: String,
                default: 'default',
                validator (value) {
                    return ['default', 'cell'].includes(value)
                }
            },
            formatCellValue: Function
        },
        data () {
            return {
                list: []
            }
        },
        computed: {
            header () {
                return (this.property.option || []).map(option => option)
            },
            displayType () {
                if (this.header.length) {
                    return this.showOn === 'default' ? 'table' : 'info'
                }
                return 'info'
            },
            cellValue () {
                const list = this.list.map(item => {
                    const values = {}
                    Object.keys(item).forEach(key => {
                        const value = item[key]
                        const options = this.property.option
                        const property = options.find(property => property.bk_property_id === key)
                        if (property) {
                            const displayValue = this.$options.filters['formatter'](value, property, property.option)
                            values[key] = displayValue
                        }
                    })
                    return values
                })
                return list
            }
        },
        watch: {
            value: {
                handler (value) {
                    const formattedValue = (value || []).map(item => {
                        const row = { ...item }
                        Object.keys(row).forEach(key => {
                            const field = row[key]
                            if (field !== null && typeof field === 'object') {
                                row[key] = field.value
                            } else {
                                row[key] = field
                            }
                        })
                        return row
                    })
                    this.list = formattedValue
                },
                immediate: true
            }
        },
        methods: {
            getCellValue () {
                if (this.formatCellValue) {
                    return (<span>{this.formatCellValue(this.cellValue)}</span>)
                }
                return (<span >{this.cellValue.map(item => (Object.values(item).join(' '))).join(',') || '--'}</span>)
            }
        }
    }
</script>

<style lang="scss" scoped>
    .table-value {
        &.property-value {
            width: 100% !important;
            padding: 0 !important;
        }
    }
</style>
