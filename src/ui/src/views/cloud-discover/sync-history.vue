<template>
    <div class="sync-history-layout">
        <div class="sync-history-options clearfix">
            <cmdb-form-date-range class="sync-options-filter" v-model="dateRange"></cmdb-form-date-range>
        </div>
        <cmdb-table ref="table"
            :loading="$loading('getSyncHistory')"
            :header="table.header"
            :list="table.list"
            :pagination.sync="table.pagination"
            :wrapper-minus-height="220"
            @handlePageChange="handlePageChange"
            @handleSizeChange="handleSizeChange"
            @handleSortChange="handleSortChange">
            <template slot="bk_status" slot-scope="{ item }">
                <span class="sync-success" v-if="item.bk_status === 'success'">
                    {{$t('Inst["成功"]')}}
                </span>
                <span class="sync-fail" v-else-if="item.bk_status === 'fail'">
                    {{$t('EventPush["失败"]')}}
                </span>
            </template>
            <template slot="details" slot-scope="{ item }">
                <span v-if="item.fail_reason === 'AuthFailure'">
                    {{ $t('Cloud["ID和Key认证失败"]') }}
                </span>
                <span v-else-if="item.fail_reason === 'else'">
                    {{ $t('Cloud["服务器错误"]') }}
                </span>
                <span v-else>
                    {{$t('Cloud[\'新增\']')}} ({{item.new_add}}) / {{$t('Cloud[\'变更\']')}} ({{item.attr_changed}})
                </span>
            </template>
        </cmdb-table>
    </div>
</template>

<script>
    import moment from 'moment'
    import { mapActions } from 'vuex'
    export default {
        props: {
            curPush: {
                type: Object
            },
            type: {
                type: String,
                default: 'create'
            }
        },
        data () {
            return {
                dateRange: [],
                operator: '',
                table: {
                    header: [{
                        id: 'bk_obj_id',
                        sortable: false,
                        width: 80,
                        name: this.$t('Nav["模型"]')
                    }, {
                        id: 'bk_status',
                        width: 80,
                        sortable: false,
                        name: this.$t('ProcessManagement["状态"]')
                    }, {
                        id: 'bk_time_consume',
                        width: 115,
                        name: this.$t('Cloud["处理耗时"]')
                    }, {
                        id: 'details',
                        sortable: false,
                        name: this.$t('Cloud["详情"]')
                    }, {
                        id: 'bk_start_time',
                        name: this.$t('HostResourcePool["启动时间"]')
                    }],
                    list: [],
                    allList: [],
                    pagination: {
                        current: 1,
                        count: 0,
                        size: 10
                    },
                    checked: [],
                    defaultSort: '-bk_start_time',
                    sort: '-bk_start_time'
                }
            }
        },
        computed: {
            filterRange () {
                return [
                    this.dateRange[0] ? this.dateRange[0] + ' 00:00:00' : '',
                    this.dateRange[1] ? this.dateRange[1] + ' 23:59:59' : ''
                ]
            }
        },
        watch: {
            'filterRange' () {
                this.getTableData()
            }
        },
        created () {
            this.initDateRange()
            this.getTableData()
        },
        methods: {
            ...mapActions('cloudDiscover', ['searchCloudHistory']),
            initDateRange () {
                const start = this.$tools.formatTime(moment().subtract(14, 'days'), 'YYYY-MM-DD')
                const end = this.$tools.formatTime(moment(), 'YYYY-MM-DD')
                this.dateRange = [start, end]
            },
            async getTableData () {
                const params = {}
                const innerParams = {}
                const pagination = this.table.pagination
                const page = {
                    start: (pagination.current - 1) * pagination.size,
                    limit: pagination.size,
                    sort: this.table.sort
                }
                innerParams['$gte'] = this.filterRange[0]
                innerParams['$lte'] = this.filterRange[1]
                params['bk_start_time'] = innerParams
                params['bk_task_id'] = this.curPush.bk_task_id
                params['page'] = page
                const res = await this.searchCloudHistory({ params, config: { requestID: 'getSyncHistory' } })
                this.table.list = res.info.map(data => {
                    data['start_time'] = this.$tools.formatTime(data['start_time'], 'YYYY-MM-DD HH:mm:ss')
                    data['bk_obj_id'] = this.$t('Hosts["主机"]')
                    return data
                })
                pagination.count = res.count
            },
            handlePageChange (page) {
                this.table.pagination.current = page
                this.getTableData()
            },
            handleSizeChange (size) {
                this.table.pagination.size = size
                this.handlePageChange(1)
            },
            handleSortChange (sort) {
                this.table.sort = sort
                this.getTableData()
            }
        }
    }
</script>

<style lang="scss" scoped>
    .sync-history-layout {
        position: relative;
        height: 100px;
        padding: 0 20px;
    }
    .sync-history-options {
        padding: 20px 0;
        .sync-options-filter {
            display: inline-block;
            vertical-align: middle;
            width: 240px;
        }
    }
    .sync-success {
        color: #2cc545;
    }
    .sync-fail {
        color: #fc2e2e;
    }
</style>
