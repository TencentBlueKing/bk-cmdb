<template>
    <div class="filing-wrapper" v-show="isShow">
        <div class="title-contain clearfix">
            <span class="title">{{$t('Common["已删除历史"]')}}</span>
            <div class="fr operation-group">
                <bk-daterangepicker
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
                :tableHeader="tableHeader"
                :tableList="tableList"
                :pagination="pagination"
                :isLoading="isLoading"
                :sortable="false"
                @handlePageTurning="setCurrentPage"
                @handlePageSizeChange="setCurrentSize"
                @handleTableSortClick="setCurrentSort"
            >
            </v-table>
        </div>
    </div>
</template>

<script>
    import moment from 'moment'
    import vTable from '@/components/table/table'
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
                'language'
            ]),
            tableHeader () {
                let header = this.$deepClone(this.objTableHeader)
                // 为业务时删除第一列的ID
                if (this.objId === 'biz') {
                    header = header.slice(1)
                }
                if (header.length && header[0].hasOwnProperty('type') && header[0].type === 'checkbox') {
                    header.shift()
                }
                header.push({
                    id: 'op_time',
                    name: this.$t('EventPush["更新时间"]')
                })
                header.unshift({
                    id: 'id',
                    name: 'ID'
                })
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
            searchParams () {
                if (!this.opTime.length) {
                    this.setFilterTime(null, `${this.startDate} - ${this.endDate}`)
                }
                let params = {
                    condition: {
                        op_type: 'delete',
                        op_time: this.opTime,
                        op_target: this.objId
                    },
                    start: (this.pagination.current - 1) * this.pagination.size,
                    limit: this.pagination.size,
                    sort: this.sort
                }
                return params
            }
        },
        watch: {
            isShow (isShow) {
                if (isShow) {
                    this.setFilterTime(null, `${this.startDate} - ${this.endDate}`)
                    this.resetDateRangePicker()
                }
            },
            '$route.path' () {
                this.$emit('update:isShow', false)
            }
        },
        methods: {
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
                if (this.opTime.length === 2) {
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
                    let res = await this.$axios.post('audit/search/', this.searchParams)
                    this.initTableList(res.data.info)
                    this.pagination.count = res.data.count
                } catch (e) {
                    this.$alertMsg(e.message || e.statusText || e.data['bk_error_msg'])
                } finally {
                    this.isLoading = false
                }
            },
            initTableList (list) {
                list.forEach((item, index) => {
                    this.tableHeader.map((list, hIndex) => {
                        if (hIndex === 0) {
                            if (this.objId === 'host') {
                                item['id'] = item.content['pre_data']['bk_host_id']
                            } else if (this.objId === 'biz') {
                                item['id'] = item.content['pre_data']['bk_biz_id']
                            } else {
                                item['id'] = item.content['pre_data']['bk_inst_id']
                            }
                        } else if (hIndex === (this.tableHeader.length - 1)) {
                            item['op_time'] = this.$formatTime(moment(item['op_time']))
                        } else if (list.property['bk_property_type'] === 'singleasst' || list.property['bk_property_type'] === 'multiasst') {
                            let name = []
                            if (item.content['pre_data'].hasOwnProperty(list.id)) {
                                if (item.content['pre_data'][list.id]) {
                                    item.content['pre_data'][list.id].map(({bk_inst_name: bkInstName}) => {
                                        name.push(bkInstName)
                                    })
                                } else {
                                    name.push('')
                                }
                            }
                            item[list.id] = name.join(',')
                        } else if (list.property['bk_property_type'] === 'enum') {
                            let option = (list.property.option || []).find(({id}) => id === item.content['pre_data'][list.id])
                            item[list.id] = option ? option.name : ''
                        } else {
                            item[list.id] = item.content['pre_data'][list.id]
                        }
                    })
                })
                this.tableList = list
            },
            closeFiling () {
                this.$emit('update:isShow', false)
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
            vTable
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
