<template>
    <div class="history-layout">
        <div class="history-options clearfix">
            <label class="fl">{{$t('Common["已删除历史"]')}}</label>
            <bk-button class="fr" type="primary" @click="back">{{$t('Common["返回"]')}}</bk-button>
            <bk-date-range class="history-date-range fr"
                :ranges="ranges"
                :range-separator="'-'"
                :quick-select="true"
                :start-date="startDate"
                :end-date="endDate"
                position="bottom-left"
                @change="setFilterTime">
            </bk-date-range>
        </div>
        <cmdb-table class="history-table"
            :sortable="false"
            :pagination.sync="pagination"
            :list="list"
            :header="header"
            :wrapperMinusHeight="157"
            @handlePageChange="handlePageChange"
            @handleSizeChange="handleSizeChange">
        </cmdb-table>
    </div>
</template>

<script>
    import { mapGetters, mapActions } from 'vuex'
    import moment from 'moment'
    export default {
        data () {
            return {
                properties: [],
                header: [],
                list: [],
                pagination: {
                    current: 1,
                    size: 10,
                    count: 0
                },
                ranges: {},
                opTime: []
            }
        },
        computed: {
            ...mapGetters(['supplierAccount']),
            ...mapGetters('userCustom', ['usercustom']),
            customColumns () {
                const customKeyMap = {
                    [this.objId]: `${this.objId}_table_columns`,
                    'host': 'resource_table_columns'
                }
                return this.usercustom[customKeyMap[this.objId]] || []
            },
            objId () {
                return this.$route.params.objId
            },
            startDate () {
                return this.$tools.formatTime(moment().subtract(1, 'month'), 'YYYY-MM-DD')
            },
            endDate () {
                return this.$tools.formatTime(moment(), 'YYYY-MM-DD')
            }
        },
        async created () {
            try {
                this.properties = await this.searchObjectAttribute({
                    params: {
                        bk_obj_id: this.objId,
                        bk_supplier_account: this.supplierAccount
                    },
                    config: {
                        requestId: `${this.objId}Attribute`,
                        fromCache: true
                    }
                })
                await this.setTableHeader()
                await this.initFilterTime()
            } catch (e) {
                // ignore
            }
        },
        methods: {
            ...mapActions('objectModelProperty', ['searchObjectAttribute']),
            ...mapActions('operationAudit', ['getOperationLog']),
            back () {
                this.$router.go(-1)
            },
            initFilterTime () {
                const language = this.$i18n.locale
                if (language === 'en') {
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
                this.opTime = [this.startDate + ' 00:00:00', this.endDate + ' 23:59:59']
                return Promise.resolve()
            },
            setTableHeader () {
                const idMap = {
                    'host': 'bk_host_id',
                    'set': 'bk_set_id',
                    'module': 'bk_module_id',
                    'biz': 'bk_biz_id',
                    'plat': 'bk_plat_id',
                    [this.objId]: 'bk_inst_id'
                }
                const fixedPropertyMap = {
                    'host': 'bk_host_innerip',
                    'set': 'bk_set_name',
                    'module': 'bk_module_name',
                    'biz': 'bk_biz_name',
                    'plat': 'bk_plat_name',
                    [this.objId]: 'bk_inst_name'
                }
                const headerProperties = this.$tools.getHeaderProperties(this.properties, this.customColumns, [fixedPropertyMap[this.objId]])
                this.header = [{
                    id: idMap[this.objId],
                    name: 'ID'
                }].concat(headerProperties.map(property => {
                    return {
                        id: property['bk_property_id'],
                        name: property['bk_property_name']
                    }
                })).concat([{
                    id: 'last_time',
                    name: this.$t('Common["更新时间"]')
                }])
                return Promise.resolve(this.header)
            },
            setFilterTime (oldVal, newVal) {
                this.opTime = newVal.split(' - ').map((date, index) => {
                    return index === 0 ? (date + ' 00:00:00') : (date + ' 23:59:59')
                })
                this.handlePageChange(1)
            },
            getTableData () {
                this.getOperationLog({
                    params: this.getSearchParams(),
                    config: {
                        cancelPrevious: true,
                        requestId: `search${this.objId}OperationLog`
                    }
                }).then(log => {
                    this.pagination.count = log.count
                    this.list = []
                })
            },
            getSearchParams () {
                return {
                    condition: {
                        'op_type': 3,
                        'op_time': this.opTime,
                        'op_target': this.objId
                    },
                    start: (this.pagination.current - 1) * this.pagination.size,
                    limit: this.pagination.size,
                    sort: '-op_time'
                }
            },
            handleSizeChange (size) {
                this.pagination.size = size
                this.handlePageChange(1)
            },
            handlePageChange (current) {
                this.pagination.current = current
                this.getTableData()
            }
        }
    }
</script>

<style lang="scss" scoped>
    .history-layout{
        padding: 20px;
    }
    .history-options{
        height: 36px;
        line-height: 36px;
        font-size: 14px;
    }
    .history-table{
        margin-top: 20px;
    }
</style>

<style lang="scss">
    .history-date-range{
        .range-action{
            display: none;
        }
    }
</style>