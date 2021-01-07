<template>
    <div class="cloud-account-layout">
        <cmdb-tips class="cloud-account-tips" tips-key="cloud-account-tips">
            {{$t('云资源发现提示语')}}
        </cmdb-tips>
        <div class="cloud-account-options">
            <div class="options-left">
                <cmdb-auth :auth="{ type: $OPERATION.C_CLOUD_RESOURCE_TASK }">
                    <bk-button theme="primary" slot-scope="{ disabled }"
                        :disabled="disabled"
                        @click="handleCreate">
                        {{$t('新建')}}
                    </bk-button>
                </cmdb-auth>
            </div>
            <div class="options-right">
                <bk-input class="options-filter" clearable
                    v-model.trim="filter"
                    right-icon="icon-search"
                    :placeholder="$t('请输入xx', { name: $t('任务名称') })">
                </bk-input>
            </div>
        </div>
        <bk-table class="cloud-account-table"
            v-bkloading="{ isLoading: $loading(Object.values(request)) }"
            :data="list"
            :pagination="pagination"
            :max-height="$APP.height - 220"
            @sort-change="handleSortChange"
            @page-change="handlePageChange"
            @page-limit-change="handleLimitChange"
            @cell-click="handleCellClick">
            <bk-table-column
                :label="$t('任务名称')"
                sortable="custom"
                prop="bk_task_name"
                class-name="is-highlight"
                show-overflow-tooltip>
            </bk-table-column>
            <bk-table-column :label="$t('资源')" prop="bk_resource_type" :formatter="resourceTypeFormatter"></bk-table-column>
            <bk-table-column :label="$t('账户名称')" prop="bk_account_name" show-overflow-tooltip sortable="custom">
                <task-account-selector slot-scope="{ row }"
                    display="info"
                    :value="row.bk_account_id">
                </task-account-selector>
            </bk-table-column>
            <bk-table-column :label="$t('账户类型')" prop="bk_cloud_vendor" sortable="custom">
                <cmdb-vendor slot-scope="{ row }" :type="row.bk_cloud_vendor"></cmdb-vendor>
            </bk-table-column>
            <bk-table-column :label="$t('最近同步状态')" prop="bk_sync_status" sortable="custom">
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
            <bk-table-column :label="$t('最近同步时间')" prop="bk_last_sync_time" show-overflow-tooltip sortable="custom">
                <template slot-scope="{ row }">{{row.bk_last_sync_time | formatter('time')}}</template>
            </bk-table-column>
            <bk-table-column :label="$t('编辑人')" prop="bk_last_editor">
                <template slot-scope="{ row }">{{row.bk_last_editor | formatter('singlechar')}}</template>
            </bk-table-column>
            <bk-table-column :label="$t('操作')" fixed="right">
                <template slot-scope="{ row }">
                    <cmdb-auth class="mr10" :auth="{ type: $OPERATION.U_CLOUD_RESOURCE_TASK, relation: [row.bk_task_id] }">
                        <link-button slot-scope="{ disabled }"
                            :disabled="disabled"
                            @click="handleEdit(row)">
                            {{$t('编辑')}}
                        </link-button>
                    </cmdb-auth>
                    <cmdb-auth class="mr10" :auth="{ type: $OPERATION.D_CLOUD_RESOURCE_TASK, relation: [row.bk_task_id] }">
                        <link-button slot-scope="{ disabled }"
                            :disabled="disabled"
                            @click="handleDelete(row)">
                            {{$t('删除')}}
                        </link-button>
                    </cmdb-auth>
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
    import TaskAccountSelector from './children/task-account-selector.vue'
    import TaskForm from './children/task-form.vue'
    import Bus from '@/utils/bus.js'
    import symbols from './common/symbol'
    import CmdbVendor from '@/components/ui/other/vendor'
    import throttle from 'lodash.throttle'
    import RouterQuery from '@/router/query'
    export default {
        components: {
            TaskSideslider,
            TaskAccountSelector,
            CmdbVendor
        },
        data () {
            return {
                filter: '',
                list: [],
                pagination: this.$tools.getDefaultPaginationConfig(),
                sort: 'bk_task_id',
                request: {
                    findTask: Symbol('findTask'),
                    findAccount: Symbol('findAccount')
                },
                scheduleSearch: throttle(this.handlePageChange, 800, { leading: false, trailing: true })
            }
        },
        watch: {
            filter () {
                this.scheduleSearch()
            }
        },
        created () {
            Bus.$on('request-refresh', this.getData)
            this.getData()
            this.unwatch = RouterQuery.watch('_t', () => {
                this.handlePageChange(1)
            })
        },
        beforeDestroy () {
            Bus.$off('request-refresh', this.getData)
            this.$http.cancelCache(symbols.all)
            this.unwatch && this.unwatch()
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
            handleSortChange (sort) {
                this.sort = this.$tools.getSort(sort, { prop: 'bk_task_id' })
                this.getData()
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
                const infoInstance = this.$bkInfo({
                    title: this.$t('确认删除xx', { instance: row.bk_task_name }),
                    closeIcon: false,
                    confirmFn: async () => {
                        try {
                            await this.$store.dispatch('cloud/resource/deleteTask', {
                                id: row.bk_task_id
                            })
                            this.$success('删除成功')
                            this.getData()
                            return true
                        } catch (e) {
                            console.error(e)
                            return false
                        } finally {
                            infoInstance.buttonLoading = false
                        }
                    }
                })
            },
            async getData () {
                try {
                    const params = {
                        fields: [],
                        condition: {},
                        exact: false,
                        page: {
                            ...this.$tools.getPageParams(this.pagination),
                            sort: this.sort
                        }
                    }
                    if (this.filter) {
                        params.condition.bk_task_name = this.filter
                        params.is_fuzzy = true
                    }
                    const data = await this.$store.dispatch('cloud/resource/findTask', {
                        params: params,
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
        display: flex;
        align-items: center;
        justify-content: space-between;
        .options-filter {
            width: 260px;
        }
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
