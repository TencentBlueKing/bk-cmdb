<template>
    <div class="history-layout">
        <div class="history-options clearfix">
            <label class="fl">{{$t('已删除历史')}}</label>
            <bk-button class="fr ml10" theme="primary" @click="back">{{$t('返回')}}</bk-button>
            <bk-date-picker class="history-date-range fr"
                placement="bottom-end"
                type="daterange"
                :shortcuts="ranges"
                :clearable="false"
                :start-date="startDate"
                v-model="defaultDate"
                @change="setFilterTime">
            </bk-date-picker>
        </div>
        <bk-table class="history-table"
            v-bkloading="{ isLoading: $loading() }"
            :pagination="pagination"
            :data="list"
            @page-change="handlePageChange"
            @page-limit-change="handleSizeChange">
            <bk-table-column v-for="column in header"
                :key="column.id"
                :prop="column.id"
                :label="column.name">
            </bk-table-column>
        </bk-table>
    </div>
</template>

<script>
    import { mapGetters, mapActions } from 'vuex'
    import moment from 'moment'
    export default {
        data () {
            const startDate = moment().subtract(1, 'month').toDate()
            const endDate = moment().toDate()
            const opSatrtTime = this.$tools.formatTime(startDate, 'YYYY-MM-DD') + ' 00:00:00'
            const opEndTime = this.$tools.formatTime(endDate, 'YYYY-MM-DD') + ' 23:59:59'
            return {
                properties: [],
                header: [],
                list: [],
                pagination: {
                    current: 1,
                    limit: 10,
                    count: 0
                },
                opTime: [opSatrtTime, opEndTime],
                startDate,
                endDate,
                opTimeResolver: null,
                defaultDate: [startDate, endDate]
            }
        },
        computed: {
            ...mapGetters(['supplierAccount', 'userName', 'isAdminView']),
            ...mapGetters('userCustom', ['usercustom']),
            ...mapGetters('objectBiz', ['bizId']),
            customColumns () {
                const customKeyMap = {
                    [this.objId]: `${this.userName}_${this.objId}_${this.isAdminView ? 'adminView' : this.bizId}_table_columns`,
                    'host': `${this.userName}_$resource_${this.isAdminView ? 'adminView' : this.bizId}_table_columns`
                }
                return this.usercustom[customKeyMap[this.objId]] || []
            },
            objId () {
                return this.$route.params.objId
            },
            ranges () {
                const language = this.$i18n.locale
                if (language === 'en') {
                    return [{
                        text: 'Yesterday',
                        value () {
                            return [moment().subtract(1, 'days').toDate(), moment().toDate()]
                        }
                    }, {
                        text: 'Last Week',
                        value () {
                            return [moment().subtract(7, 'days').toDate(), moment().toDate()]
                        }
                    }, {
                        text: 'Last Month',
                        value () {
                            return [moment().subtract(1, 'month').toDate(), moment().toDate()]
                        }
                    }, {
                        text: 'Last Three Month',
                        value () {
                            return [moment().subtract(3, 'month').toDate(), moment().toDate()]
                        }
                    }]
                }
                return [{
                    text: '昨天',
                    value () {
                        return [moment().subtract(1, 'days').toDate(), moment().toDate()]
                    }
                }, {
                    text: '最近一周',
                    value () {
                        return [moment().subtract(7, 'days').toDate(), moment().toDate()]
                    }
                }, {
                    text: '最近一个月',
                    value () {
                        return [moment().subtract(1, 'month').toDate(), moment().toDate()]
                    }
                }, {
                    text: '最近三个月',
                    value () {
                        return [moment().subtract(3, 'month').toDate(), moment().toDate()]
                    }
                }]
            }
        },
        watch: {
            opTime (opTime) {
                if (this.opTimeResolver) {
                    this.opTimeResolver()
                } else {
                    this.handlePageChange(1)
                }
            }
        },
        async created () {
            try {
                // await this.setTimeResolver()
                this.properties = await this.searchObjectAttribute({
                    params: this.$injectMetadata({
                        bk_obj_id: this.objId,
                        bk_supplier_account: this.supplierAccount
                    }),
                    config: {
                        requestId: `post_searchObjectAttribute_${this.objId}`,
                        fromCache: false
                    }
                })
                await this.setTableHeader()
                this.getTableData()
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
            setTimeResolver () {
                return new Promise((resolve, reject) => {
                    this.opTimeResolver = () => {
                        this.opTimeResolver = null
                        resolve()
                    }
                })
            },
            setTableHeader () {
                const idMap = {
                    'host': 'bk_host_id',
                    'set': 'bk_set_id',
                    'module': 'bk_module_id',
                    'biz': 'bk_biz_id',
                    'plat': 'bk_plat_id'
                }
                const fixedPropertyMap = {
                    'host': ['bk_host_innerip', 'bk_cloud_id'],
                    'set': ['bk_set_name'],
                    'module': ['bk_module_name'],
                    'biz': ['bk_biz_name'],
                    'plat': ['bk_plat_name']
                }
                const headerProperties = this.$tools.getHeaderProperties(this.properties, this.customColumns, fixedPropertyMap[this.objId] || ['bk_inst_name'])
                this.header = [{
                    id: idMap[this.objId] || 'bk_inst_id',
                    name: 'ID'
                }].concat(headerProperties.map(property => {
                    return {
                        id: property['bk_property_id'],
                        name: property['bk_property_name']
                    }
                })).concat([{
                    id: 'op_time',
                    width: 180,
                    name: this.$t('更新时间')
                }])
                return Promise.resolve(this.header)
            },
            setFilterTime (daterange) {
                this.opTime = daterange.map((date, index) => {
                    return index === 0 ? (date + ' 00:00:00') : (date + ' 23:59:59')
                })
            },
            getTableData () {
                this.getOperationLog({
                    params: this.getSearchParams(),
                    config: {
                        cancelPrevious: true,
                        requestId: `search${this.objId}OperationLog`
                    }
                }).then(log => {
                    try {
                        this.pagination.count = log.count
                        const list = log.info.map(data => {
                            return {
                                ...(data.content['cur_data'] ? data.content['cur_data'] : data.content['pre_data']),
                                'op_time': this.$tools.formatTime(data['op_time'])
                            }
                        })
                        this.list = this.$tools.flattenList(this.properties, list)
                    } catch (e) {
                        this.list = []
                        this.$error(e.message)
                    }
                })
            },
            getSearchParams () {
                return {
                    condition: {
                        'op_type': 3,
                        'op_time': this.opTime,
                        'op_target': this.objId
                    },
                    start: (this.pagination.current - 1) * this.pagination.limit,
                    limit: this.pagination.limit,
                    sort: '-op_time'
                }
            },
            handleSizeChange (size) {
                this.pagination.limit = size
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
