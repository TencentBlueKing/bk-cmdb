<template>
    <div class="confirm-history-layout">
        <div class="confirm-history-options clearfix">
            <bk-button class="fl" theme="primary" @click="back">{{$t('返回')}}</bk-button>
            <cmdb-form-date-range class="confirm-filter" v-model="dateRange"></cmdb-form-date-range>
        </div>
        <bk-table
            v-bkloading="{ isLoading: $loading('getConfirHistory') }"
            :data="table.list"
            :pagination="table.pagination"
            :max-height="$APP.height - 220"
            @page-change="handlePageChange"
            @page-limit-change="handleSizeChange"
            @sort-change="handleSortChange">
            <bk-table-column v-for="column in table.header"
                :key="column.id"
                :prop="column.id"
                :label="column.name">
                <template slot-scope="{ row }">
                    <template v-if="column.id === 'bk_resource_type'">
                        <span class="change-span" v-if="row.bk_resource_type === 'change'">
                            {{$t('变更')}}
                        </span>
                        <span class="new-add-span" v-else>
                            {{$t('新增')}}
                        </span>
                    </template>
                    <template v-else-if="column.id === 'bk_account_type' && row.bk_account_type === 'tencent_cloud'">
                        {{$t('腾讯云')}}
                    </template>
                    <template v-else>{{row[column.id]}}</template>
                </template>
            </bk-table-column>
        </bk-table>
    </div>
</template>

<script>
    import moment from 'moment'
    import { mapActions } from 'vuex'
    export default {
        props: {
            curPush: {
                type: {
                    type: String,
                    default: 'create'
                }
            }
        },
        data () {
            return {
                dateRange: [],
                operator: '',
                table: {
                    header: [{
                        id: 'bk_host_innerip',
                        name: this.$t('资源名称')
                    }, {
                        id: 'bk_resource_type',
                        name: this.$t('资源类型')
                    }, {
                        id: 'bk_obj_id',
                        name: this.$t('模型')
                    }, {
                        id: 'bk_task_name',
                        name: this.$t('任务名称')
                    }, {
                        id: 'bk_account_type',
                        name: this.$t('账号类型')
                    }, {
                        id: 'bk_account_admin',
                        sortable: false,
                        name: this.$t('任务维护人')
                    }, {
                        id: 'create_time',
                        width: 180,
                        name: this.$t('发现时间')
                    }, {
                        id: 'confirm_time',
                        width: 180,
                        name: this.$t('确认时间')
                    }],
                    list: [],
                    allList: [],
                    pagination: {
                        current: 1,
                        count: 0,
                        limit: 10
                    },
                    checked: [],
                    defaultSort: '-confirm_history_id',
                    sort: '-confirm_history_id'
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
            ...mapActions('cloudDiscover', ['searchConfirmHistory']),
            initDateRange () {
                const start = this.$tools.formatTime(moment().subtract(14, 'days'), 'YYYY-MM-DD')
                const end = this.$tools.formatTime(moment(), 'YYYY-MM-DD')
                this.dateRange = [start, end]
            },
            async getTableData () {
                const pagination = this.table.pagination
                const params = {}
                const innerParams = {}
                const page = {
                    start: (pagination.current - 1) * pagination.limit,
                    limit: pagination.limit,
                    sort: this.table.sort
                }
                innerParams['$gte'] = this.filterRange[0]
                innerParams['$lte'] = this.filterRange[1]
                params['confirm_time'] = innerParams
                params['page'] = page
                const res = await this.searchConfirmHistory({ params, config: { requestID: 'getConfirHistory' } })
                this.table.list = res.info.map(data => {
                    data['create_time'] = this.$tools.formatTime(data['create_time'], 'YYYY-MM-DD HH:mm:ss')
                    data['confirm_time'] = this.$tools.formatTime(data['confirm_time'], 'YYYY-MM-DD HH:mm:ss')
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
            back () {
                this.$router.go(-1)
            },
            handleSortChange (sort) {
                this.table.sort = this.$tools.getSort(sort)
                this.getTableData()
            }
        }
    }
</script>

<style lang="scss" scoped>
    .confirm-history-layout {
        position: relative;
        height: 100%;
        padding: 0 20px;
    }
    .confirm-history-options {
        padding: 20px 0;
        .confirm-filter{
            float: right;
            display: inline-block;
            vertical-align: middle;
            width: 240px;
        }
    }
    .change-span {
        color: #ffb23a;
    }
    .new-add-span {
        color: #4f55f3;
    }
</style>
