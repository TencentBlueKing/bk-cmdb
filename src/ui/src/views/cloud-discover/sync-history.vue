<template>
    <div class="sync-history-layout">
        <div class="sync-history-options clearfix">
            <cmdb-form-date-range class="sync-options-filter" v-model="dateRange"></cmdb-form-date-range>
        </div>
        <bk-table
            v-bkloading="{ isLoading: $loading('getSyncHistory') }"
            :data="table.list"
            :pagination="table.pagination"
            :max-height="$APP.height - 220"
            @page-change="handlePageChange"
            @page-limit-change="handleSizeChange"
            @sort-change="handleSortChange">
            <bk-table-column prop="bk_obj_id" :label="$t('模型')" width="80"></bk-table-column>
            <bk-table-column prop="bk_status" :label="$t('状态')" width="80">
                <template slot-scope="{ row }">
                    <span class="sync-success" v-if="row.bk_status === 'success'">
                        {{$t('成功')}}
                    </span>
                    <span class="sync-fail" v-else-if="row.bk_status === 'fail'">
                        {{$t('失败')}}
                    </span>
                </template>
            </bk-table-column>
            <bk-table-column prop="bk_time_consume" :label="$t('处理耗时')" width="115"></bk-table-column>
            <bk-table-column prop="details" :label="$t('详情')">
                <template slot-scope="{ row }">
                    <span v-if="row.fail_reason === 'AuthFailure'">
                        {{ $t('ID和Key认证失败') }}
                    </span>
                    <span v-else-if="row.fail_reason === 'else'">
                        {{ $t('服务器错误') }}
                    </span>
                    <span v-else>
                        {{$t('新增')}} ({{row.new_add}}) / {{$t('变更update')}} ({{row.attr_changed}})
                    </span>
                </template>
            </bk-table-column>
            <bk-table-column prop="bk_start_time" :label="$t('启动时间')"></bk-table-column>
        </bk-table>
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
                    list: [],
                    allList: [],
                    pagination: {
                        current: 1,
                        count: 0,
                        limit: 10
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
                    start: (pagination.current - 1) * pagination.limit,
                    limit: pagination.limit,
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
                    data['bk_obj_id'] = this.$t('主机')
                    return data
                })
                pagination.count = res.count
            },
            handlePageChange (page) {
                this.table.pagination.current = page
                this.getTableData()
            },
            handleSizeChange (size) {
                this.table.pagination.limit = size
                this.handlePageChange(1)
            },
            handleSortChange (sort) {
                this.table.sort = this.$tools.getSort(sort)
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
