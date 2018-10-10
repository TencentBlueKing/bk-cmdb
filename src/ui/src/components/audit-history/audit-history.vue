<template>
    <div class="audit-history-layout">
        <div class="history-options clearfix">
            <div class="options-group fl">
                <label class="options-label">{{$t("HostResourcePool['时间范围']")}}</label>
                <cmdb-form-date-range class="options-filter" v-model="dateRange"></cmdb-form-date-range>
            </div>
            <div class="options-group fl">
                <label class="options-label">{{$t("HostResourcePool['操作账号']")}}</label>
                <cmdb-form-objuser class="options-filter" v-model="operator" :exclude="false" :multiple="false"></cmdb-form-objuser>
            </div>
            <bk-button class="fr" type="primary" @click="refresh">{{$t("Common['查询']")}}</bk-button>
        </div>
        <cmdb-table class="audit-table"
            :loading="$loading('getOperationLog')"
            :header="header"
            :list="list"
            :pagination.sync="pagination"
            :wrapperMinusHeight="220"
            @handlePageChange="handlePageChange"
            @handleSizeChange="handleSizeChange"
            @handleSortChange="handleSortChange"
            @handleRowClick="handleRowClick">
            <template slot="op_time" slot-scope="{ item }">
                {{$tools.formatTime(item['op_time'])}}
            </template>
        </cmdb-table>
        <div class="history-details" v-if="details.isShow" v-click-outside="closeDetails">
            <p class="details-title">
                <span>{{$t('OperationAudit[\'操作详情\']')}}</span>
                <i class="bk-icon icon-close" @click="closeDetails"></i>
            </p>
            <v-details class="details-content" 
                :isShow="this.details.isShow"
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
                header: [{
                    id: 'op_desc',
                    name: this.$t("HostResourcePool['变更内容']")
                }, {
                    id: 'operator',
                    name: this.$t("HostResourcePool['操作账号']")
                }, {
                    id: 'op_time',
                    name: this.$t("HostResourcePool['操作时间']")
                }],
                list: [],
                pagination: {
                    count: 0,
                    current: 1,
                    size: 10
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
                return [
                    this.dateRange[0] ? this.dateRange[0] + ' 00:00:00' : '',
                    this.dateRange[1] ? this.dateRange[1] + ' 23:59:59' : ''
                ]
            }
        },
        created () {
            this.initDateRange()
            this.refresh()
        },
        beforeDestroy () {
            this.$http.cancel('getOperationLog')
        },
        methods: {
            ...mapActions('operationAudit', ['getOperationLog']),
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
            refresh () {
                this.getOperationLog({
                    params: this.getParams(),
                    config: {
                        cancelPrevious: true,
                        requestId: 'getOperationLog'
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
                if (this.operator) {
                    condition.operator = this.operator
                }
                return {
                    condition,
                    limit: this.pagination.size,
                    sort: this.sort,
                    start: (this.pagination.current - 1) * this.pagination.size
                }
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
            },
            handleRowClick (item) {
                this.details.isShow = true
                this.details.clickoutside = true
                this.details.data = item
                this.$nextTick(() => {
                    this.details.clickoutside = false
                })
            }
        },
        components: {
            vDetails
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
                width: 240px;
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