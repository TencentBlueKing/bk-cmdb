<template>
    <div class="audit-history-layout">
        <div class="history-options clearfix">
            <label class="options-label">{{$t("HostResourcePool['时间范围']")}}</label>
            <cmdb-form-date-range class="options-filter" v-model="dateRange" style="width: 240px"></cmdb-form-date-range>
        </div>
        <cmdb-table class="audit-table" ref="table"
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
    import moment from 'moment'
    import { mapActions } from 'vuex'
    export default {
        props: {
            curPush: {
            }
        },
        data () {
            return {
                dateRange: [],
                operator: '',
                table: {
                    header: [ {
                        id: 'bk_obj_id',
                        name: '模型'
                    }, {
                        id: 'bk_status',
                        name: '状态'
                    }, {
                        id: 'bk_time_consume',
                        name: '处理耗时'
                    }, {
                        id: 'bk_sync_detail',
                        name: '新增/改变'
                    }, {
                        id: 'bk_start_time',
                        name: '启动时间'
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
        beforeDestroy () {
            this.$http.cancel('getOperationLog')
        },
        methods: {
            ...mapActions('operationAudit', ['getOperationLog']),
            ...mapActions('cloudDiscover', ['searchCloudHistory']),
            initDateRange () {
                const start = this.$tools.formatTime(moment().subtract(14, 'days'), 'YYYY-MM-DD')
                const end = this.$tools.formatTime(moment(), 'YYYY-MM-DD')
                this.dateRange = [start, end]
            },
            async getTableData () {
                let pagination = this.table.pagination
                let taskID = this.curPush.bk_task_id
                let res = await this.searchCloudHistory({taskID})
                this.table.list = res.info.map(data => {
                    data['bk_start_time'] = this.$tools.formatTime(data['bk_start_time'], 'YYYY-MM-DD HH:mm:ss')
                    if (data['bk_obj_id'] === 'host') {
                        data['bk_obj_id'] = '主机'
                    } else {
                        data['bk_obj_id'] = '交换机'
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

        },
        components: {
        }
    }
</script>

<style lang="scss" scoped>
    .audit-history-layout{
        position: relative;
        height: 100%;
        padding: 0 20px;
    }
    .history-options{
        padding: 20px 0;
        .options-group{
            white-space: nowrap;
            margin-right: 20px;
            .options-label{
                display: inline-block;
                vertical-align: middle;
            }
            .options-filter{
                display: inline-block;
                vertical-align: middle;
                width: 280px;
            }
        }
    }
    .history-details{
        position: absolute;
        top: 20px;
        left: 30px;
        width: 709px;
        height: 577px;
        background-color: #ffffff;
        box-shadow: 0px 2px 9px 0px rgba(0, 0, 0, 0.4);
        z-index: 1;
        .details-title{
            position: relative;
            margin: 15px 0;
            line-height: 26px;
            color: #333948;
            padding: 0 40px;
            font-weight: bold;
            .icon-close{
                font-size: 14px;
                position: absolute;
                right: 12px;
                top: 0;
                cursor: pointer;
            }
        }
        .details-content{
            padding: 0 40px;
        }
    }
</style>
