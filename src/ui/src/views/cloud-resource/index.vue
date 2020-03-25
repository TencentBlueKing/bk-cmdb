<template>
    <div class="cloud-account-layout">
        <cmdb-tips class="cloud-account-tips">提示语</cmdb-tips>
        <div class="cloud-account-options">
            <bk-button theme="primary" @click="handleCreate">{{$t('新建')}}</bk-button>
        </div>
        <bk-table class="cloud-account-table"
            v-bkloading="{ isLoading: $loading(Object.values(request)) }"
            :data="list"
            :pagination="pagination"
            @page-change="handlePageChange"
            @page-limit-change="handleLimitChange"
            @cell-click="handleCellClick">
            <bk-table-column :label="$t('任务名称')" prop="bk_task_name" class-name="is-highlight" show-overflow-tooltip></bk-table-column>
            <bk-table-column :label="$t('资源')" prop="bk_resource_type" :formatter="resourceTypeFormatter"></bk-table-column>
            <bk-table-column :label="$t('账户名称')" prop="bk_account_name" show-overflow-tooltip>
                <task-account-selector slot-scope="{ row }"
                    display="info"
                    :value="row.bk_account_id">
                </task-account-selector>
            </bk-table-column>
            <bk-table-column :label="$t('账户类型')" prop="bk_cloud_vendor" :formatter="vendorFormatter"></bk-table-column>
            <bk-table-column :label="$t('最近同步状态')" prop="bk_sync_status">
                <div class="row-status"
                    slot-scope="{ row }"
                    v-bk-tooltips.right="{
                        disabled: !(row.bk_status_description && row.bk_status_description.error_info),
                        content: row.bk_status_description && row.bk_status_description.error_info
                    }">
                    <template v-if="row.bk_sync_status">
                        <i :class="['status', { 'is-error': row.bk_sync_status !== 'cloud_sync_success' }]"></i>
                        {{row.bk_sync_status !== 'cloud_sync_success' ? $t('失败') : $t('成功')}}
                    </template>
                    <template v-else>--</template>
                </div>
            </bk-table-column>
            <bk-table-column :label="$t('最近同步时间')" prop="bk_last_sync_time" show-overflow-tooltip>
                <template slot-scope="{ row }">{{row.bk_last_sync_time | formatter('time')}}</template>
            </bk-table-column>
            <bk-table-column :label="$t('最近编辑人')" prop="bk_last_editor">
                <template slot-scope="{ row }">{{row.bk_last_editor | formatter('singlechar')}}</template>
            </bk-table-column>
            <bk-table-column :label="$t('操作')" fixed="right">
                <template slot-scope="{ row }">
                    <link-button class="mr10" @click="handleEdit(row)">{{$t('编辑')}}</link-button>
                    <bk-popconfirm
                        trigger="click"
                        :title="$t('确定删除该任务')"
                        @confirm="handleDelete(row)">
                        <link-button>{{$t('删除')}}</link-button>
                    </bk-popconfirm>
                </template>
            </bk-table-column>
        </bk-table>
        <task-sideslider ref="taskSideslider"
            @request-refresh="getData">
        </task-sideslider>
    </div>
</template>

<script>
    import TaskSideslider from './children/task-sideslider.vue'
    import { formatter as resourceTypeFormatter } from '@/dictionary/cloud-resource-type'
    import { formatter as vendorFormatter } from '@/dictionary/cloud-vendor'
    import TaskAccountSelector from './children/task-account-selector.vue'
    import TaskForm from './children/task-form.vue'
    import Bus from '@/utils/bus.js'
    export default {
        components: {
            TaskSideslider,
            TaskAccountSelector
        },
        data () {
            return {
                list: [],
                pagination: this.$tools.getDefaultPaginationConfig(),
                request: {
                    findTask: Symbol('findTask'),
                    findAccount: Symbol('findAccount')
                }
            }
        },
        created () {
            Bus.$on('request-refresh', this.getData)
            this.getData()
        },
        beforeDestroy () {
            Bus.$off('request-refresh', this.getData)
        },
        methods: {
            handleCreate () {
                this.$refs.taskSideslider.show({
                    mode: 'create',
                    title: this.$t('新建发现任务'),
                    props: {
                        type: 'create'
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
                if (column.property === 'bk_task_name') {
                    this.handleView(row)
                }
            },
            handleView (row) {
                this.$refs.taskSideslider.show({
                    mode: 'details',
                    title: `${this.$t('任务详情')} 【${row.bk_task_name}】`,
                    props: {
                        id: row.bk_task_id
                    }
                })
            },
            handleEdit (row) {
                this.$refs.taskSideslider.show({
                    mode: 'details',
                    title: `${this.$t('任务详情')} 【${row.bk_task_name}】`,
                    props: {
                        id: row.bk_task_id,
                        defaultComponent: TaskForm.name
                    }
                })
            },
            async handleDelete (row) {
                try {
                    await this.$store.dispatch('cloud/resource/deleteTask', {
                        id: row.bk_task_id
                    })
                    this.$success('删除成功')
                    this.getData()
                } catch (e) {
                    console.error(e)
                }
            },
            async getData () {
                try {
                    const data = await this.$store.dispatch('cloud/resource/findTask', {
                        params: {
                            fields: [],
                            condition: {},
                            exact: false,
                            page: this.$tools.getPageParams(this.pagination)
                        },
                        config: {
                            requestId: this.request.findTask
                        }
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
            },
            resourceTypeFormatter (row, column) {
                return this.$t(resourceTypeFormatter(row[column.property]))
            },
            vendorFormatter (row, column) {
                return this.$t(vendorFormatter(row[column.property]))
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
            &.is-error {
                background-color: $dangerColor;
            }
        }
    }
</style>
