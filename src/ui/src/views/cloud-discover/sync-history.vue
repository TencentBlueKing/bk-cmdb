<template>
    <div class="audit-history-layout">
        <div class="history-options clearfix">
            <label class="options-label">{{$t("HostResourcePool['时间范围']")}}</label>
            <cmdb-form-date-range class="options-filter" v-model="dateRange" style="width: 240px"></cmdb-form-date-range>
        </div>
        <cmdb-table class="audit-table" ref="table"
            :loading="$loading('getSyncHistory')"
            :header="table.header"
            :list="table.list"
            :pagination.sync="table.pagination"
            :wrapperMinusHeight="220"
            @handlePageChange="handlePageChange"
            @handleSizeChange="handleSizeChange"
            @handleSortChange="handleSortChange">
                <template slot="details" slot-scope="{ item }">
                    {{$t('Cloud[\'新增\']')}} ({{item.new_add}}) / {{$t('Cloud[\'变更\']')}} ({{item.attr_changed}})
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
                    header: [ {
                        id: 'bk_obj_id',
                        name: this.$t('Nav["模型"]')
                    }, {
                        id: 'bk_status',
                        name: this.$t('ProcessManagement["状态"]')
                    }, {
                        id: 'bk_time_consume',
                        name: this.$t('Cloud["处理耗时"]')
                    }, {
                        id: 'details',
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
                    defaultSort: '-bk_task_id',
                    sort: '-bk_task_id'
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
                let params = {}
                let innerParams = {}
                let pagination = this.table.pagination
                innerParams['$gte'] = this.filterRange[0]
                innerParams['$lte'] = this.filterRange[1]
                params['bk_start_time'] = innerParams
                params['bk_task_id'] = this.curPush.bk_task_id
                let res = await this.searchCloudHistory({params, config: {requestID: 'getSyncHistory'}})
                this.table.list = res.info.map(data => {
                    data['bk_start_time'] = this.$tools.formatTime(data['bk_start_time'], 'YYYY-MM-DD HH:mm:ss')
                    data['bk_obj_id'] = this.$t('Hosts["主机"]')
                    if (data['bk_status'] === 'waiting_confirm') {
                        data['bk_status'] = '等待确认'
                    } else if (data['bk_status'] === 'success') {
                        data['bk_status'] = '成功'
                    } else {
                        data['bk_status'] = '失败'
                    }
                    return data
                })
                pagination.count = res.count
            },
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
        },
        watch: {
            'filterRange' () {
                this.getTableData()
            }
        }
    }
</script>

<style lang="scss" scoped>
    .audit-history-layout {
        position: relative;
        height: 100%;
        padding: 0 20px;
    }
    .history-options {
        padding: 20px 0;
    }
</style>
