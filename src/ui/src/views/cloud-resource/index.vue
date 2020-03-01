<template>
    <div class="cloud-account-layout">
        <cmdb-tips class="cloud-account-tips">提示语</cmdb-tips>
        <div class="cloud-account-options">
            <bk-button theme="primary" @click="handleCreate">{{$t('新建')}}</bk-button>
        </div>
        <bk-table class="cloud-account-table"
            :data="list"
            :pagination="pagination"
            @page-change="handlePageChange"
            @page-limit-change="handleLimitChange"
            @cell-click="handleCellClick">
            <bk-table-column :label="$t('任务名称')" prop="name" class-name="is-highlight"></bk-table-column>
            <bk-table-column :label="$t('资源')" prop="resource"></bk-table-column>
            <bk-table-column :label="$t('账户名称')" prop="account_name"></bk-table-column>
            <bk-table-column :label="$t('账户类型')" prop="account_type"></bk-table-column>
            <bk-table-column :label="$t('最近同步状态')" prop="status">
                <div class="row-status"
                    slot-scope="{ row }"
                    v-bk-tooltips.right="{
                        disabled: !row.error,
                        content: '异常原因'
                    }">
                    <i :class="['status', { 'is-error': row.error }]"></i>
                    异常
                </div>
            </bk-table-column>
            <bk-table-column :label="$t('最近同步时间')" prop="last_time"></bk-table-column>
            <bk-table-column :label="$t('最近编辑人')" prop="account_type"></bk-table-column>
            <bk-table-column :label="$t('操作')">
                <template slot-scope="{ row }">
                    <link-button class="mr10" @click="handleEdit(row)">{{$t('编辑')}}</link-button>
                    <link-button @click="handleDelete(row)">{{$t('删除')}}</link-button>
                </template>
            </bk-table-column>
        </bk-table>
        <resource-create-sideslider ref="resourceCreateSideslider"
            @request-refresh="getData">
        </resource-create-sideslider>
        <resource-details-sideslider ref="resourceDetailsSideslider"
            @request-refresh="getData">
        </resource-details-sideslider>
    </div>
</template>

<script>
    import ResourceCreateSideslider from './children/resource-sideslider.vue'
    import ResourceDetailsSideslider from './children/resource-details-sideslider.vue'
    export default {
        components: {
            ResourceCreateSideslider,
            ResourceDetailsSideslider
        },
        data () {
            return {
                list: [],
                pagination: {
                    ...this.$tools.getDefaultPaginationConfig(),
                    count: 2
                }
            }
        },
        created () {
            this.getData()
        },
        methods: {
            handleCreate () {
                this.$refs.resourceCreateSideslider.show({
                    type: 'form',
                    title: this.$t('新建发现任务'),
                    props: {
                        mode: 'create'
                    }
                })
            },
            handlePageChange (page) {
                this.pagination.current = page
                this.getData()
            },
            handleLimitChange (limit) {
                this.pagination.limit = limit
                this.pagination.current = 1
                this.getData()
            },
            handleCellClick (row, column) {
                if (column.property === 'name') {
                    this.handleView(row)
                }
            },
            handleView (row) {
                this.$refs.resourceDetailsSideslider.show({
                    type: 'details',
                    title: `${this.$t('任务详情')} 【${row.name}】`,
                    props: {
                        id: row.id
                    }
                })
            },
            handleEdit (row) {},
            async handleDelete (row) {
                try {
                    await Promise.resolve()
                    this.$success('删除成功')
                    this.getData()
                } catch (e) {
                    console.error(e)
                }
            },
            async getData () {
                try {
                    const data = await Promise.resolve({
                        count: 2,
                        info: [{ id: 1 }, { id: 2 }]
                    })
                    if (data.count && !data.info.length) {
                        this.handlePageChange(this.pagination.current - 1)
                        return
                    }
                    this.list = data.info
                    this.pagination.count = data.count
                } catch (e) {
                    console.error(e)
                    this.list = []
                    this.pagination.count = 0
                }
            }
        }
    }
</script>

<style lang="scss" scoped>
    .cloud-account-layout {
        padding: 0 20px;
    }
    .cloud-account-tips {
        margin-top: 10px;
    }
    .cloud-account-options {
        margin-top: 10px;
    }
    .cloud-account-table {
        margin-top: 10px;
    }
    .row-status {
        display: inline-block;
        .status {
            display: inline-block;
            margin-right: 4px;
            width: 7px;
            height: 7px;
            border-radius: 50%;
            background-color: $successColor;
        }
    }
</style>
