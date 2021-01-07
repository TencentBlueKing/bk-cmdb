<template>
    <div class="process-bind-info-value">
        <bk-popover v-bind="popoverOptoins" :disabled="popoverList.length < 2">
            <table-value
                ref="table"
                :value="localValue"
                :show-on="'cell'"
                :format-cell-value="formatCellValue"
                :property="property">
            </table-value>
            <ul slot="content">
                <li v-for="(item, index) in popoverList" :key="index">{{item}}</li>
            </ul>
        </bk-popover>
    </div>
</template>

<script>
    import TableValue from '@/components/ui/other/table-value'
    export default {
        components: {
            TableValue
        },
        props: {
            value: {
                type: Array,
                default: () => ([])
            },
            property: {
                type: Object,
                default: () => ({})
            },
            popoverOptoins: {
                type: Object,
                default: () => ({})
            }
        },
        data () {
            return {
                popoverList: [],
                localValue: []
            }
        },
        watch: {
            value: {
                handler (value) {
                    this.localValue = value || []
                    this.setPopoverList()
                },
                immediate: true
            }
        },
        methods: {
            ipText (value) {
                const map = {
                    '1': '127.0.0.1',
                    '2': '0.0.0.0',
                    '3': this.$t('第一内网IP'),
                    '4': this.$t('第一外网IP')
                }
                return map[value] || value || '--'
            },
            setPopoverList () {
                this.$nextTick(() => {
                    const list = this.$refs.table.cellValue
                    this.popoverList = list.map(this.getRowValue)
                })
            },
            getRowValue (row) {
                const ip = this.ipText(row.ip)
                return `${row.protocol} ${ip}:${row.port}`
            },
            formatCellValue (list) {
                if (!list.length) {
                    return '--'
                }
                const newList = list.map(this.getRowValue)
                const total = list.length
                const showCount = total > 1
                return (
                    <div class={`bind-info-value${showCount ? ' show-count' : ''}`}>
                        <span>{newList.join(', ')}</span>
                        {showCount ? <span class="count">{total}</span> : ''}
                    </div>
                )
            }
        }
    }
</script>

<style lang="scss" scoped>
    .bind-info-value {
        position: relative;
        display: inline-block;
        vertical-align: middle;
        padding-right: 0;
        max-width: 100%;
        overflow: hidden;
        text-overflow: ellipsis;
        white-space: nowrap;

        &.show-count {
            padding-right: 26px;
        }

        .count {
            position: absolute;
            display: inline-block;
            right: 2px;
            top: 0;
            color: #979ba5;
            background-color: #f0f1f5;
            border-radius: 3px;
            padding: 0 2px;
        }
    }

    /deep/.process-bind-info-value {
        .bk-tooltip,
        .bk-tooltip-ref {
            display: block;
        }
    }
</style>
