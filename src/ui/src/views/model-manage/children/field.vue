<template>
    <div>
        <bk-button class="create-btn" type="primary">
            {{$t('ModelManagement["新建字段"]')}}
        </bk-button>
        <cmdb-table
            class="field-table"
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
                        name: this.$t('ModelManagement["字段名称"]')
                    }, {
                        id: 'name',
                        name: this.$t('ModelManagement["类型"]')
                    }, {
                        id: 'name',
                        name: this.$t('ModelManagement["唯一"]')
                    }, {
                        id: 'name',
                        name: this.$t('ModelManagement["必填"]')
                    }, {
                        id: 'name',
                        name: this.$t('ModelManagement["字段描述"]')
                    }, {
                        id: 'last_time',
                        name: this.$t('ModelManagement["创建时间"]')
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
