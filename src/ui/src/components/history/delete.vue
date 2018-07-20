<template>
    <div class="filing-wrapper" v-show="isShow">
        <div class="title-contain clearfix">
            <span class="title" v-if="objId !== 'biz'">{{$t('Common["已删除历史"]')}}</span>
            <span class="title" v-else>{{$t('Common["归档历史"]')}}</span>
            <div class="fr operation-group">
                <bk-daterangepicker v-if="objId !== 'biz'"
                    ref="dateRangePicker"
                    class="datepicker"
                    :ranges="ranges"
                    :range-separator="'-'"
                    :quick-select="true"
                    :start-date="startDate"
                    :end-date="endDate"
                    @change="setFilterTime">
                </bk-daterangepicker>
                <bk-button type="primary" @click="closeFiling">{{$t('Common["返回"]')}}</bk-button>
            </div>
        </div>
        <div class="table-content">
            <v-table
                :header="tableHeader"
                :list="tableList"
                :pagination.sync="pagination"
                :loading="isLoading"
                :sortable="false"
                :wrapperMinusHeight="150"
                @handlePageChange="setCurrentPage"
                @handleSizeChange="setCurrentSize"
                @handleSortChange="setCurrentSort"
                @handleRowClick="showDetails">
                <template slot="$recovery" slot-scope="{ item }">
                    <bk-button type="primary" size="mini" @click="recoveryBizConfirm(item)">{{$t('Inst["恢复业务"]')}}</bk-button>
                </template>
            </v-table>
        </div>
        <v-sideslider :isShow.sync="details.isShow" :title="{text: $t('OperationAudit[\'操作详情\']')}">
            <v-history-details :details="details.data" :isShow="details.isShow" slot="content"></v-history-details>
        </v-sideslider>
    </div>
</template>

<script>
    import moment from 'moment'
    import vTable from '@/components/table/table'
    import vSideslider from '@/components/slider/sideslider'
    import vHistoryDetails from '@/components/history/details'
    import {mapGetters} from 'vuex'
    export default {
        props: {
            isShow: {
                type: Boolean,
                default: false
            },
            objId: {
                required: true
            },
            objTableHeader: {
                required: true
            }
        },
        data () {
            return {
                details: {
                    isShow: false,
                    data: null
                },
                pagination: {
                    current: 1,
                    count: 0,
                    size: 10
                },
                ranges: [],
                tableList: [],
                isLoading: false,
                opTime: [],
                sort: '-op_time'
            }
        },
        computed: {
            ...mapGetters([
                'language',
                'bkSupplierAccount'
            ]),
            tableHeader () {
                let header = this.$deepClone(this.objTableHeader)
                if (header.length && header[0].hasOwnProperty('type') && header[0].type === 'checkbox') {
                    header.shift()
                }
                header.push({
                    id: '$op_time',
                    name: this.$t('EventPush["更新时间"]')
                })
                header.unshift({
                    id: 'id',
                    name: 'ID'
                })
                // 为业务时删除第一列的ID
                if (this.objId === 'biz') {
                    header = header.slice(1)
                    header.push({
                        id: '$recovery',
                        name: '操作'
                    })
                }
                return header
            },
            /* 开始时间 */
            startDate () {
                return this.$formatTime(moment().subtract(1, 'month'), 'YYYY-MM-DD')
            },
            /* 结束时间 */
            endDate () {
                return this.$formatTime(moment(), 'YYYY-MM-DD')
            },
            axiosConfig () {
                let config = {
                    url: '',
                    params: {}
                }
                if (!this.opTime.length) {
                    this.setFilterTime(null, `${this.startDate} - ${this.endDate}`)
                }
                if (this.objId === 'biz') {
                    config.url = `biz/search/${this.bkSupplierAccount}`
                    config.params = {
                        condition: {
                            bk_data_status: 'disabled'
                            // last_time: {
                            //     '$gt': this.opTime[0],
                            //     '$lt': this.opTime[1],
                            //     'cc_time_type': 1
                            // }
                        },
                        fields: [],
                        page: {
                            start: (this.pagination.current - 1) * this.pagination.size,
                            limit: this.pagination.size,
                            sort: this.sort
                        }
                    }
                } else {
                    config.url = 'audit/search/'
                    config.params = {
                        condition: {
                            op_type: 3, // delete
                            op_time: this.opTime,
                            op_target: this.objId
                        },
                        start: (this.pagination.current - 1) * this.pagination.size,
                        limit: this.pagination.size,
                        sort: this.sort
                    }
                }
                return config
            }
        },
        watch: {
            isShow (isShow) {
                if (isShow) {
                    if (this.objId !== 'biz') {
                        this.setFilterTime(null, `${this.startDate} - ${this.endDate}`)
                        this.resetDateRangePicker()
                    } else {
                        this.getTableList()
                    }
                }
            },
            '$route.path' () {
                this.$emit('update:isShow', false)
            }
        },
        methods: {
            recoveryBizConfirm (item) {
                this.$bkInfo({
                    title: this.$t('Inst["是否确认恢复业务？"]'),
                    content: this.$t('Inst["恢复业务提示"]', {bizName: item['bk_biz_name']}),
                    confirmFn: () => {
                        this.recoveryBiz(item)
                    }
                })
            },
            async recoveryBiz (item) {
                try {
                    const res = await this.$axios.put(`biz/status/enable/${this.bkSupplierAccount}/${item['bk_biz_id']}`)
                    if (!res.result) {
                        this.$alertMsg(res['bk_error_msg'])
                    }
                    this.getTableList()
                } catch (e) {
                    this.$alertMsg(e.message || e.data['bk_error_msg'] || e.statusText)
                }
            },
            resetDateRangePicker () {
                this.$refs.dateRangePicker.selectedDateRange = [this.startDate, this.endDate]
                this.$refs.dateRangePicker.selectedDateRangeTmp = [this.startDate, this.endDate]
            },
            setFilterTime (oldVal, newVal) {
                this.opTime = newVal.split(' - ').map((time, index) => {
                    if (index === 0) {
                        return time + ' 00:00:00'
                    } else {
                        return time + ' 23:59:59'
                    }
                })
                if (this.opTime.length === 2 && this.isShow) {
                    this.getTableList()
                }
            },
            /* 表格排序事件 */
            setCurrentSort (sort) {
                this.sort = sort
                this.setCurrentPage(1)
            },
            /* 翻页事件，设置搜索参数页码 */
            setCurrentPage (current) {
                this.pagination.current = current
                this.getTableList()
            },
            /* 设置每页显示数量 */
            setCurrentSize (size) {
                this.pagination.size = size
                this.setCurrentPage(1)
            },
            /* 获取表格数据 */
            async getTableList () {
                this.isLoading = true
                try {
                    let res = await this.$axios.post(this.axiosConfig.url, this.axiosConfig.params)
                    this.initTableList(res.data.info)
                    this.pagination.count = res.data.count
                } catch (e) {
                    this.$alertMsg(e.message || e.data['bk_error_msg'] || e.statusText)
                } finally {
                    this.isLoading = false
                }
            },
            initTableList (list) {
                list.forEach((item, index) => {
                    this.tableHeader.map((list, hIndex) => {
                        let data = this.objId === 'biz' ? item : item.content['pre_data']
                        // 如果该字段为null则不展示该行
                        if (data !== null && list.id !== '$recovery') {
                            if (hIndex === 0) {
                                if (this.objId === 'host') {
                                    item['id'] = data['bk_host_id']
                                } else if (this.objId === 'biz') {
                                    item['id'] = data['bk_biz_id']
                                } else {
                                    item['id'] = data['bk_inst_id']
                                }
                            } else if (list.id === '$op_time') {
                                item['$op_time'] = this.objId === 'biz' ? this.$formatTime(moment(item['last_time'])) : this.$formatTime(moment(item['op_time']))
                            } else if (list.property['bk_property_type'] === 'singleasst' || list.property['bk_property_type'] === 'multiasst') {
                                let name = []
                                if (data.hasOwnProperty(list.id)) {
                                    if (data[list.id]) {
                                        data[list.id].map(({bk_inst_name: bkInstName}) => {
                                            name.push(bkInstName)
                                        })
                                    } else {
                                        name.push('')
                                    }
                                }
                                item[list.id] = name.join(',')
                            } else if (list.property['bk_property_type'] === 'enum') {
                                let option = (list.property.option || []).find(({id}) => id === data[list.id])
                                item[list.id] = option ? option.name : ''
                            } else if (['date', 'time'].includes(list.property['bk_property_type'])) {
                                const format = list.property['bk_property_type'] === 'date' ? 'YYYY-MM-DD' : 'YYYY-MM-DD HH:mm:ss'
                                item[list.id] = this.$formatTime(data[list.id], format)
                            } else {
                                item[list.id] = data[list.id]
                            }
                        }
                    })
                })
                this.tableList = list
            },
            closeFiling () {
                this.$emit('update:isShow', false)
            },
            showDetails (item) {
                if (this.objId !== 'biz') {
                    this.details.data = item
                    this.details.isShow = true
                }
            }
        },
        created () {
            if (this.language === 'en') {
                this.ranges = {
                    'Yesterday': [moment().subtract(1, 'days'), moment()],
                    'Last Week': [moment().subtract(7, 'days'), moment()],
                    'Last Month': [moment().subtract(1, 'month'), moment()],
                    'Last Three Month': [moment().subtract(3, 'month'), moment()]
                }
            } else {
                this.ranges = {
                    昨天: [moment().subtract(1, 'days'), moment()],
                    最近一周: [moment().subtract(7, 'days'), moment()],
                    最近一个月: [moment().subtract(1, 'month'), moment()],
                    最近三个月: [moment().subtract(3, 'month'), moment()]
                }
            }
        },
        components: {
            vTable,
            vSideslider,
            vHistoryDetails
        }
    }
</script>

<style lang="scss" scoped>
    .filing-wrapper{
        position: absolute;
        left: 0;
        top: 0;
        right: 0;
        bottom: 0;
        background: #fff;
        z-index: 1199;
        .title-contain{
            padding:20px 20px 0 20px;
            .title{
                line-height: 36px;
                font-size: 14px;
            }
        }
        .operation-group{
            font-size: 0;
            .datepicker{
                float: left;
                margin-right: 20px;
            }
        }
        .table-content{
            padding: 20px;
        }
    }
</style>
