<template>
    <div>
        <bk-button class="create-btn" type="primary">
            {{$t('ModelManagement["新建关联关系"]')}}
        </bk-button>
        <cmdb-table
            class="relation-table"
            :loading="$loading('getOperationLog')"
            :header="table.header"
            :list="table.list"
            :pagination.sync="table.pagination"
            :wrapperMinusHeight="220"
            @handlePageChange="handlePageChange"
            @handleSizeChange="handleSizeChange"
            @handleSortChange="handleSortChange">
        </cmdb-table>
    </div>
</template>

<script>
    export default {
        data () {
            return {
                table: {
                    header: [{
                        id: 'name',
                        name: this.$t('ModelManagement["唯一标识"]')
                    }, {
                        id: 'name',
                        name: this.$t('ModelManagement["关联类型"]')
                    }, {
                        id: 'name',
                        name: this.$t('ModelManagement["约束"]')
                    }, {
                        id: 'name',
                        name: this.$t('ModelManagement["源模型"]')
                    }, {
                        id: 'last_time',
                        name: this.$t('ModelManagement["目标模型"]')
                    }, {
                        id: 'operation',
                        name: this.$t('ModelManagement["操作"]')
                    }],
                    list: [],
                    pagination: {
                        count: 0,
                        current: 1,
                        size: 10
                    },
                    defaultSort: '-op_time',
                    sort: '-op_time'
                }
            }
        },
        methods: {
            handlePageChange (current) {
                this.pagination.current = current
                this.refresh()
            },
            handleSizeChange (size) {
                this.pagination.size = size
                this.handlePageChange(1)
            },
            handleSortChange (sort) {
                this.sort = sort
                this.refresh()
            }
        }
    }
</script>

<style lang="scss" scoped>
    .create-btn {
        margin: 10px 0;
    }
</style>
