/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2017-2018 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and limitations under the License.
 */

<template>
    <div class="history-wrapper">
        <div class="history-filter clearfix">
            <div class="filter-group date fl">
                <label>{{$t("HostResourcePool['时间范围']")}}</label>
                <bk-daterangepicker class="filter-field" ref="dateRangePicker"
                    @change="setFilterDate"
                    :range-separator="'-'"
                    :quick-select="false"
                    :disabled="false"
                    :start-date="initDate.start"
                    :end-date="initDate.end">
                </bk-daterangepicker>
            </div>
            <div class="filter-group user fl">
                <label>{{$t("HostResourcePool['操作账号']")}}</label>
                <v-member-selector class="filter-field" :exclude="true" :selected.sync="filter.user" :active="active" :multiple="false"></v-member-selector>
            </div>
            <div class="filter-group btn fr">
                <bk-button type="primary" :loading="$loading('auditHistory')" @click="setCurrentPage(1)">{{$t("Common['查询']")}}</bk-button>
            </div>
        </div>
        <div class="history-table">
            <v-table
                :loading="$loading('auditHistory')" 
                :header="table.header" 
                :list="table.list" 
                :pagination.sync="table.pagination" 
                :defaultSort="table.defaultSort"
                :wrapperMinusHeight="270"
                @handleSortChange="setTableSort"
                @handlePageChange="setCurrentPage"
                @handleSizeChange="setPageSize"
                @handleRowClick="showDetails">
            </v-table>
        </div>
        <div class="history-details" v-show="details.isShow" v-click-outside="closeDetails">
            <p class="details-title">
                <span>{{$t('OperationAudit[\'操作详情\']')}}</span>
                <i class="bk-icon icon-close" @click="closeDetails"></i>
            </p>
            <v-history-details style="padding: 0 40px;" 
                :details="details.data"
                :isShow="this.details.isShow" 
                slot="content" 
                :height="342" 
                :width="635"></v-history-details>
        </div>
    </div>
</template>
<script>
    import vTable from '@/components/table/table'
    import vMemberSelector from '@/components/common/selector/member'
    import vHistoryDetails from '@/components/history/details'
    import moment from 'moment'
    export default {
        props: {
            active: {
                type: Boolean,
                default: false
            },
            type: {
                type: String,
                default: 'inst' // inst | host
            },
            innerIP: String,
            instId: Number
        },
        data () {
            return {
                details: {
                    isShow: false,
                    data: null,
                    clickoutside: true
                },
                filter: {
                    date: [],
                    user: ''
                },
                table: {
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
                        current: 1,
                        count: 0,
                        size: 10
                    },
                    defaultSort: '-op_time',
                    sort: '-op_time'
                }
            }
        },
        computed: {
            searchParams () {
                let params = {
                    condition: {
                        op_time: this.filter.date,
                        op_target: this.type
                    },
                    limit: this.table.pagination.size,
                    start: (this.table.pagination.current - 1) * this.table.pagination.size,
                    sort: this.table.sort
                }
                if (this.type === 'host') {
                    params.condition['ext_key'] = {
                        '$in': [this.innerIP]
                    }
                } else {
                    params.condition['inst_id'] = this.instId
                }
                if (this.filter.user) {
                    params.condition['operator'] = this.filter.user
                }
                return params
            },
            initDate () {
                return {
                    start: this.$formatTime(moment().subtract(14, 'days'), 'YYYY-MM-DD'),
                    end: this.$formatTime(moment(), 'YYYY-MM-DD')
                }
            }
        },
        beforeMount () {
            this.filter.date = [`${this.initDate.start} 00:00:00`, `${this.initDate.end} 23:59:59`]
        },
        watch: {
            active (active) {
                if (active) {
                    this.getHistory()
                } else {
                    let $dateRangePicker = this.$refs.dateRangePicker
                    $dateRangePicker.selectedDateView = `${this.initDate.start} - ${this.initDate.end}`
                    $dateRangePicker.selectedDateRange = [this.initDate.start, this.initDate.end]
                    $dateRangePicker.selectedDateRangeTmp = [this.initDate.start, this.initDate.end]
                    this.filter.date = [`${this.initDate.start} 00:00:00`, `${this.initDate.end} 23:59:59`]
                    this.filter.user = ''
                    this.table.pagination.current = 1
                }
            }
        },
        methods: {
            showDetails (item) {
                this.details.isShow = true
                this.details.clickoutside = true
                this.details.data = item
                this.$nextTick(() => {
                    this.details.clickoutside = false
                })
            },
            closeDetails () {
                if (!this.details.clickoutside) {
                    this.details.isShow = false
                    this.details.data = null
                }
            },
            getHistory () {
                this.$axios.post('audit/search', this.searchParams, {id: 'auditHistory'}).then(res => {
                    if (res.result) {
                        res.data.info.map(history => {
                            history['op_time'] = this.$formatTime(history['op_time'])
                        })
                        this.table.list = res.data.info
                        this.table.pagination.count = res.data.count
                    } else {
                        this.$alertMsg(res['bk_error_msg'])
                    }
                })
            },
            setFilterDate (oldDate, newDate) {
                if (newDate) {
                    newDate = newDate.split(' - ')
                    newDate[0] = `${newDate[0]} 00:00:00`
                    newDate[1] = `${newDate[1]} 23:59:59`
                    this.filter.date = newDate
                }
            },
            setTableSort (sort) {
                this.table.sort = sort
                this.getHistory()
            },
            setPageSize (size) {
                this.table.pagination.size = size
                this.getHistory()
            },
            setCurrentPage (current) {
                this.table.pagination.current = current
                this.getHistory()
            }
        },
        components: {
            vTable,
            vMemberSelector,
            vHistoryDetails
        }
    }
</script>

<style lang="scss" scoped>
    .history-wrapper{
        position: relative;
    }
    .history-filter{
        padding: 20px 0;
        position: relative;
        z-index: 1;
        .filter-group{
            white-space: nowrap;
            margin: 0 20px 0 0px;
            &.btn{
                margin: 0;
            }
            &.date{
                .filter-field{
                    white-space: normal;
                }
            }
            label{
                display: inline-block;
                vertical-align: middle;
            }
            .filter-field{
                display: inline-block;
                vertical-align: middle;
                width: 240px;
            }
            .bk-button{
                width: 96px;
            }
        }
    }
    .history-details{
        position: absolute;
        top:20px;
        left: 30px;
        width: 709px;
        height: 577px;
        background-color: #ffffff;
        box-shadow: 0px 2px 9px 0px rgba(0, 0, 0, 0.4);
        z-index: 1;
        .details-title{
            position: relative;
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
    }
</style>