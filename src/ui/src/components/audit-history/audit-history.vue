<template>
    <div class="audit-history-layout">
        <div class="history-options clearfix">
            <div class="options-group fl">
                <label class="options-label">{{$t('时间范围')}}</label>
                <cmdb-form-date-range class="options-filter" :clearable="false" v-model="dateRange"></cmdb-form-date-range>
            </div>
            <div class="options-group fl" style="margin: 0">
                <label class="options-label">{{$t('操作账号')}}</label>
                <cmdb-form-objuser class="options-filter"
                    v-model="operator"
                    :exclude="false"
                    :multiple="false"
                    :palceholder="$t('操作账号')">
                </cmdb-form-objuser>
            </div>
            <bk-button class="fl ml10" theme="primary" @click="refresh(true)">{{$t('查询')}}</bk-button>
        </div>
        <bk-table
            v-bkloading="{ isLoading: $loading('getUserOperationLog') }"
            :data="list"
            :pagination="pagination"
            :max-height="$APP.height - 220"
            :row-style="{ cursor: 'pointer' }"
            @page-change="handlePageChange"
            @page-limit-change="handleSizeChange"
            @sort-change="handleSortChange"
            @row-click="handleRowClick">
            <bk-table-column :label="$t('变更内容')" prop="op_desc" show-overflow-tooltip></bk-table-column>
            <bk-table-column :label="$t('操作账号')" prop="operator" show-overflow-tooltip></bk-table-column>
            <bk-table-column :label="$t('操作时间')" prop="op_time" show-overflow-tooltip>
                <template slot-scope="{ row }">
                    {{row.op_time | formatter('time')}}
                </template>
            </bk-table-column>
        </bk-table>
        <div class="history-details" v-if="details.isShow" v-click-outside="closeDetails">
            <p class="details-title">
                <span>{{$t('操作详情')}}</span>
                <i class="bk-icon icon-close" @click="closeDetails"></i>
            </p>
            <v-details class="details-content"
                :is-show="details.isShow"
                :details="details.data"
                :height="342"
                :width="635"></v-details>
        </div>
    </div>
</template>

<script>
    import moment from 'moment'
    import vDetails from './details'
    import { mapActions } from 'vuex'
    export default {
        components: {
            vDetails
        },
        props: {
            extKey: {
                type: Object,
                default () {
                    return null
                }
            },
            target: {
                type: String,
                default: ''
            },
            instId: {
                type: Number
            }
        },
        data () {
            return {
                dateRange: [],
                operator: '',
                sendOperator: '',
                list: [],
                pagination: {
                    count: 0,
                    current: 1,
                    limit: 10,
                    size: 'small'
                },
                defaultSort: '-op_time',
                sort: '-op_time',
                details: {
                    isShow: false,
                    data: null,
                    clickoutside: true
                }
            }
        },
        computed: {
            filterRange () {
                const range = [
                    this.dateRange[0] ? this.dateRange[0] + ' 00:00:00' : '',
                    this.dateRange[1] ? this.dateRange[1] + ' 23:59:59' : ''
                ]
                return range.filter(date => !!date)
            }
        },
        created () {
            this.initDateRange()
            this.refresh()
        },
        beforeDestroy () {
            this.$http.cancel('getUserOperationLog')
        },
        methods: {
            ...mapActions('operationAudit', ['getUserOperationLog']),
            initDateRange () {
                const start = this.$tools.formatTime(moment().subtract(14, 'days'), 'YYYY-MM-DD')
                const end = this.$tools.formatTime(moment(), 'YYYY-MM-DD')
                this.dateRange = [start, end]
            },
            closeDetails () {
                if (!this.details.clickoutside) {
                    this.details.isShow = false
                    this.details.data = null
                }
            },
            refresh (isClickSearch) {
                if (isClickSearch) {
                    this.pagination.current = 1
                    this.sendOperator = this.operator
                }
                this.getUserOperationLog({
                    objId: this.target,
                    params: this.getParams(),
                    config: {
                        cancelPrevious: true,
                        requestId: 'getUserOperationLog'
                    }
                }).then(data => {
                    this.list = data.info
                    this.pagination.count = data.count
                })
            },
            getParams () {
                const condition = {
                    'op_target': this.target,
                    'op_time': this.filterRange
                }
                if (this.extKey) {
                    condition['ext_key'] = this.extKey
                }
                if (!isNaN(this.instId)) {
                    condition['inst_id'] = this.instId
                }
                if (this.sendOperator) {
                    condition.operator = this.sendOperator
                }
                return {
                    condition,
                    limit: this.pagination.limit,
                    sort: this.sort,
                    start: (this.pagination.current - 1) * this.pagination.limit
                }
            },
            handlePageChange (current) {
                this.pagination.current = current
                this.refresh()
            },
            handleSizeChange (size) {
                this.pagination.limit = size
                this.handlePageChange(1)
            },
            handleSortChange (sort) {
                this.sort = this.$tools.getSort(sort)
                this.refresh()
            },
            handleRowClick (item) {
                this.details.isShow = true
                this.details.clickoutside = true
                this.details.data = item
                this.$nextTick(() => {
                    this.details.clickoutside = false
                })
            }
        }
    }
</script>

<style lang="scss" scoped>
    .audit-history-layout{
        position: relative;
        height: 100%;
    }
    .history-options{
        padding: 20px 0 14px;
        font-size: 14px;
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
                width: 240px !important;
                height: 32px;
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
                font-size: 20px;
                position: absolute;
                right: 12px;
                top: 3px;
                cursor: pointer;
            }
        }
        .details-content{
            padding: 0 40px;
        }
    }
</style>
