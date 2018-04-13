<template>
    <div class="filing-wrapper" v-show="isShow">
        <div class="title-contain clearfix">
            <span class="title">已删除历史</span>
            <div class="fr operation-group">
                <bk-daterangepicker
                    class="datepicker"
                    :range-separator="'-'"
                    :quick-select="true"
                    :start-date="startDate"
                    :end-date="endDate"
                    @change="setFilterTime">
                </bk-daterangepicker>
                <bk-button type="primary" @click="closeFiling">返回</bk-button>
            </div>
        </div>
        <div class="table-content">
            <v-table
                :tableHeader="tableHeader"
                :tableList="tableList"
                :pagination="pagination"
                :isLoading="isLoading"
                @handlePageTurning="setCurrentPage"
                @handlePageSizeChange="setCurrentSize"
                @handleTableSortClick="setCurrentSort"
            ></v-table>
        </div>
    </div>
</template>

<script>
    import moment from 'moment'
    import vTable from '@/components/table/table'
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
                tableList: [],
                isLoading: false,
                opTime: [],
                sort: '-op_time'
            }
        },
        computed: {
            tableHeader () {
                let header = this.$deepClone(this.objTableHeader)
                if (header.length && header[0].hasOwnProperty('type') && header[0].type === 'checkbox') {
                    header.shift()
                }
                header.push({
                    id: 'op_time',
                    name: '更新时间'
                })
                header.unshift({
                    id: 'id',
                    name: 'ID'
                })
                return header
            },
            /* 开始时间 */
            startDate () {
                return this.$formatTime(moment().subtract(1, 'days'), 'YYYY-MM-DD')
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
                    this.getTableList()
                }
            },
            '$route.path' () {
                this.$emit('update:isShow', false)
            }
        },
        methods: {
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
                    this.$alertMsg(e.data['bk_error_msg'])
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
                            } else {
                                item['id'] = item.content['pre_data']['bk_inst_id']
                            }
                        } else if (hIndex === (this.tableHeader.length - 1)) {
                            item['op_time'] = this.$formatTime(moment(item['op_time']))
                        } else if (list.property['bk_property_type'] === 'singleasst' || list.property['bk_property_type'] === 'multiasst') {
                            item[list.id] = item.content['pre_data'][list.id][0]['bk_inst_name']
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
