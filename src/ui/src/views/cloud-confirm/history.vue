<template>
    <div class="confirm-history-layout">
        <div class="confirm-history-options clearfix">
            <bk-button class="fl" type="primary" @click="back">{{$t('Common["返回"]')}}</bk-button>
            <cmdb-form-date-range class="confirm-filter" v-model="dateRange" position="left"></cmdb-form-date-range>
        </div>
        <cmdb-table ref="table"
            :loading="$loading('getConfirHistory')"
            :header="table.header"
            :list="table.list"
            :pagination.sync="table.pagination"
            :wrapper-minus-height="220"
            @handlePageChange="handlePageChange"
            @handleSizeChange="handleSizeChange"
            @handleSortChange="handleSortChange">
            <template slot="bk_resource_type" slot-scope="{ item }">
                <span class="change-span" v-if="item.bk_resource_type === 'change'">
                    {{$t('Cloud["变更"]')}}
                </span>
                <span class="new-add-span" v-else>
                    {{$t('Cloud["新增"]')}}
                </span>
            </template>
            <template slot="bk_account_type">
                <span>{{$t('Cloud["腾讯云"]')}}</span>
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
                        name: this.$t('Cloud["资源名称"]')
                    }, {
                        id: 'bk_resource_type',
                        name: this.$t('Cloud["资源类型"]')
                    }, {
                        id: 'bk_obj_id',
                        name: this.$t('Cloud["模型"]')
                    }, {
                        id: 'bk_task_name',
                        name: this.$t('Cloud["任务名称"]')
                    }, {
                        id: 'bk_account_type',
                        name: this.$t('Cloud["账号类型"]')
                    }, {
                        id: 'bk_account_admin',
                        sortable: false,
                        name: this.$t('Cloud["任务维护人"]')
                    }, {
                        id: 'create_time',
                        width: 180,
                        name: this.$t('Cloud["发现时间"]')
                    }, {
                        id: 'confirm_time',
                        width: 180,
                        name: this.$t('Cloud["确认时间"]')
                    }],
                    list: [],
                    allList: [],
                    pagination: {
                        current: 1,
                        count: 0,
                        size: 10
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
            this.$store.commit('setHeaderTitle', this.$t('Cloud["确认记录"]'))
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
                    start: (pagination.current - 1) * pagination.size,
                    limit: pagination.size,
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
            back () {
                this.$router.go(-1)
            },
            handleSortChange (sort) {
                this.table.sort = sort
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
