<template>
    <div class="history">
        <div class="history-filter">
            <cmdb-form-date-range class="filter-item filter-range"
                v-model="dateRange"
                @input="handlePageChange(1)">
            </cmdb-form-date-range>
            <cmdb-form-objuser class="filter-item filter-user"
                v-model="operator"
                :exclude="false"
                :multiple="false"
                :palceholder="$t('操作账号')"
                @input="handlePageChange(1)">
            </cmdb-form-objuser>
        </div>
        <bk-table class="history-table"
            v-bkloading="{ isLoading: $loading('getHostAuditLog') }"
            :data="history"
            :pagination="pagination"
            :max-height="$APP.height - 325"
            :row-style="{ cursor: 'pointer' }"
            @page-change="handlePageChange"
            @page-limit-change="handleSizeChange"
            @sort-change="handleSortChange"
            @row-click="handleRowClick">
            <bk-table-column :label="$t('操作描述')" :formatter="getFormatterDesc"></bk-table-column>
            <bk-table-column prop="user" :label="$t('操作账号')" sortable="custom"></bk-table-column>
            <bk-table-column prop="operation_time" :label="$t('操作时间')" sortable="custom">
                <template slot-scope="{ row }">
                    {{$tools.formatTime(row['operation_time'])}}
                </template>
            </bk-table-column>
            <cmdb-table-empty slot="empty" :stuff="table.stuff"></cmdb-table-empty>
        </bk-table>
        <bk-sideslider
            v-transfer-dom
            :is-show.sync="details.show"
            :width="800"
            :title="$t('操作详情')">
            <cmdb-host-history-details :details="details.data" slot="content" v-if="details.show"></cmdb-host-history-details>
        </bk-sideslider>
    </div>
</template>

<script>
    import { mapActions } from 'vuex'
    import cmdbHostHistoryDetails from '@/components/audit-history/details'
    export default {
        components: {
            cmdbHostHistoryDetails
        },
        props: {
            target: {
                type: String,
                default: ''
            },
            instId: {
                type: Number
            },
            resourceType: {
                type: String,
                default: ''
            }
        },
        data () {
            return {
                dateRange: [],
                operator: '',
                history: [],
                pagination: {
                    count: 0,
                    current: 1,
                    limit: 10
                },
                table: {
                    stuff: {
                        type: 'default',
                        payload: {
                            emptyText: this.$t('bk.table.emptyText')
                        }
                    }
                },
                sort: '-operation_time',
                details: {
                    show: false,
                    data: null
                }
            }
        },
        computed: {
        },
        created () {
            this.getHistory()
        },
        methods: {
            ...mapActions('operationAudit', ['getUserOperationLog']),
            getFormatterDesc (row) {
                const funcActions = this.$store.state.operationAudit.funcActions
                const modules = [...funcActions.business, ...funcActions.resource].filter(item => [this.resourceType, 'instance_association'].includes(item.id))
                const operations = modules.reduce((acc, item) => acc.concat(item.operations), [])
                const actionSet = {}
                operations.forEach(operation => {
                    actionSet[operation.id] = this.$t(operation.name)
                })
                let action = ''
                if (row.label) {
                    const label = Object.keys(row.label)[0]
                    action = actionSet[`${row.resource_type}-${row.action}-${label}`]
                } else {
                    action = actionSet[`${row.resource_type}-${row.action}`]
                }
                let name = ''
                const data = row.operation_detail
                if (['assign_host', 'unassign_host', 'transfer_host_module'].includes(row.action)) {
                    name = data.bk_host_innerip
                } else if (['instance_association'].includes(row.resource_type)) {
                    name = data.target_instance_name
                } else {
                    name = data.basic_detail && data.basic_detail.resource_name
                }
                return `${action}"${name}"`
            },
            async getHistory (event) {
                try {
                    const condition = {
                        resource_id: Number(this.instId)
                    }
                    if (this.dateRange.length) {
                        condition.operation_time = [this.dateRange[0] + ' 00:00:00', this.dateRange[1] + ' 23:59:59']
                    }
                    if (this.operator) {
                        condition.user = this.operator
                    }

                    const data = await this.getUserOperationLog({
                        objId: this.target,
                        params: {
                            condition,
                            limit: this.pagination.limit,
                            sort: this.sort,
                            start: (this.pagination.current - 1) * this.pagination.limit
                        },
                        config: {
                            cancelPrevious: true,
                            requestId: 'getUserOperationLog'
                        }
                    })

                    this.history = data.info
                    this.pagination.count = data.count

                    if (event) {
                        this.table.stuff.type = 'search'
                    }
                } catch (e) {
                    console.log(e)
                    this.history = []
                    this.pagination.count = 0
                }
            },
            handlePageChange (page) {
                this.pagination.current = page
                this.getHistory(true)
            },
            handleSizeChange (size) {
                this.pagination.limit = size
                this.pagination.current = 1
                this.getHistory()
            },
            handleSortChange (sort) {
                this.sort = this.$tools.getSort(sort)
                this.getHistory()
            },
            handleRowClick (item) {
                this.details.data = item
                this.details.show = true
            }
        }
    }
</script>

<style lang="scss" scoped>
    .history {
        height: 100%;
    }
    .history-filter {
        padding: 14px 0;
        .filter-item {
            display: inline-block;
            vertical-align: middle;
            &.filter-range {
                width: 300px !important;
                margin: 0 5px 0 0;
            }
            &.filter-user {
                width: 240px;
                height: 32px;
            }
        }
    }
</style>
