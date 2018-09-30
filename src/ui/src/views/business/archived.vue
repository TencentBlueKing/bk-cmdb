<template>
    <div class="archived-layout">
        <div class="archived-options clearfix">
            <label class="fl">{{$t('Common["归档历史"]')}}</label>
            <bk-button class="fr" type="primary" @click="back">{{$t('Common["返回"]')}}</bk-button>
        </div>
        <cmdb-table class="archived-table"
            rowCursor="default"
            :sortable="false"
            :pagination.sync="pagination"
            :list="list"
            :header="header"
            :wrapperMinusHeight="157"
            @handlePageChange="handlePageChange"
            @handleSizeChange="handleSizeChange">
            <template slot="options" slot-scope="{ item }">
                <bk-button type="primary" size="mini" @click="handleRecovery(item)">{{$t('Inst["恢复业务"]')}}</bk-button>
            </template>
        </cmdb-table>
    </div>
</template>

<script>
    import { mapGetters, mapActions } from 'vuex'
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
                }
            }
        },
        computed: {
            ...mapGetters(['supplierAccount']),
            ...mapGetters('userCustom', ['usercustom']),
            customBusinessColumns () {
                return this.usercustom['biz_table_columns']
            }
        },
        async created () {
            try {
                this.properties = await this.searchObjectAttribute({
                    params: {
                        bk_obj_id: 'biz',
                        bk_supplier_account: this.supplierAccount
                    },
                    config: {
                        requestId: 'post_searchObjectAttribute_biz',
                        fromCache: true
                    }
                })
                this.setTableHeader()
                this.getTableData()
            } catch (e) {
                // ignore
            }
        },
        methods: {
            ...mapActions('objectModelProperty', ['searchObjectAttribute']),
            ...mapActions('objectBiz', ['searchBusiness', 'recoveryBusiness']),
            back () {
                this.$router.go(-1)
            },
            setTableHeader () {
                const headerProperties = this.$tools.getHeaderProperties(this.properties, this.customBusinessColumns, ['bk_biz_name'])
                this.header = [{
                    id: 'bk_biz_id',
                    name: 'ID'
                }].concat(headerProperties.map(property => {
                    return {
                        id: property['bk_property_id'],
                        name: property['bk_property_name']
                    }
                })).concat([{
                    id: 'last_time',
                    name: this.$t('Common["更新时间"]')
                }, {
                    id: 'options',
                    name: this.$t('Common["操作"]')
                }])
            },
            getTableData () {
                this.searchBusiness({
                    params: this.getSearchParams(),
                    config: {
                        cancelPrevious: true,
                        requestId: 'searchArchivedBusiness'
                    }
                }).then(business => {
                    this.pagination.count = business.count
                    this.list = this.$tools.flatternList(this.properties, business.info.map(biz => {
                        biz['last_time'] = this.$tools.formatTime(biz['last_time'], 'YYYY-MM-DD HH:mm:ss')
                        return biz
                    }))
                })
            },
            getSearchParams () {
                return {
                    condition: {
                        'bk_data_status': 'disabled'
                    },
                    fields: [],
                    page: {
                        start: (this.pagination.current - 1) * this.pagination.size,
                        limit: this.pagination.size,
                        sort: '-bk_biz_id'
                    }
                }
            },
            handleRecovery (biz) {
                this.$bkInfo({
                    title: this.$t('Inst["是否确认恢复业务？"]'),
                    content: this.$t('Inst["恢复业务提示"]', {bizName: biz['bk_biz_name']}),
                    confirmFn: () => {
                        this.recoveryBiz(biz)
                    }
                })
            },
            recoveryBiz (biz) {
                this.recoveryBusiness({
                    params: {
                        'bk_biz_id': biz['bk_biz_id']
                    },
                    config: {
                        cancelWhenRouteChange: false
                    }
                }).then(() => {
                    this.$http.cancel('post_searchBusiness_$ne_disabled')
                    this.$success(this.$t('Inst["恢复业务成功"]'))
                    this.handlePageChange(1)
                })
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
    .archived-layout{
        padding: 20px;
    }
    .archived-options{
        height: 36px;
        line-height: 36px;
        font-size: 14px;
    }
    .archived-table{
        margin-top: 20px;
    }
</style>