<template>
    <div class="archived-layout">
        <div class="archived-filter">
            <div class="filter-item">
                <bk-input v-model="filter.name"
                    clearable
                    :placeholder="$t('请输入xx', { name: $t('业务') })"
                    right-icon="bk-icon icon-search"
                    @enter="handlePageChange(1, $event)">
                </bk-input>
            </div>
        </div>
        <bk-table class="archived-table"
            :pagination="pagination"
            :data="list"
            :max-height="$APP.height - 190"
            @page-change="handlePageChange"
            @page-limit-change="handleSizeChange">
            <bk-table-column v-for="column in header"
                :key="column.id"
                :prop="column.id"
                :label="column.name">
            </bk-table-column>
            <bk-table-column :label="$t('操作')" fixed="right">
                <template slot-scope="{ row }">
                    <cmdb-auth class="inline-block-middle" :auth="$authResources({ type: $OPERATION.BUSINESS_ARCHIVE })">
                        <bk-button slot-scope="{ disabled }"
                            theme="primary"
                            size="small"
                            :disabled="disabled"
                            @click="handleRecovery(row)">
                            {{$t('恢复业务')}}
                        </bk-button>
                    </cmdb-auth>
                </template>
            </bk-table-column>
            <cmdb-table-empty slot="empty" :stuff="table.stuff"></cmdb-table-empty>
        </bk-table>
    </div>
</template>

<script>
    import { mapGetters, mapActions } from 'vuex'
    import { MENU_RESOURCE_BUSINESS, MENU_RESOURCE_MANAGEMENT } from '@/dictionary/menu-symbol'
    export default {
        data () {
            return {
                properties: [],
                header: [],
                list: [],
                filter: {
                    range: [],
                    name: ''
                },
                pagination: {
                    current: 1,
                    count: 0,
                    ...this.$tools.getDefaultPaginationConfig()
                },
                table: {
                    stuff: {
                        type: 'default',
                        payload: {
                            emptyText: this.$t('bk.table.emptyText')
                        }
                    }
                }
            }
        },
        computed: {
            ...mapGetters(['supplierAccount', 'isAdminView', 'userName']),
            ...mapGetters('userCustom', ['usercustom']),
            ...mapGetters('objectBiz', ['bizId']),
            customBusinessColumns () {
                return this.usercustom[`${this.userName}_biz_${this.isAdminView ? 'adminView' : this.bizId}_table_columns`]
            }
        },
        async created () {
            try {
                this.setDynamicBreadcrumbs()
                this.properties = await this.searchObjectAttribute({
                    params: this.$injectMetadata({
                        bk_obj_id: 'biz',
                        bk_supplier_account: this.supplierAccount
                    }),
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
            setDynamicBreadcrumbs () {
                this.$store.commit('setBreadcrumbs', [{
                    label: this.$t('资源目录'),
                    route: {
                        name: MENU_RESOURCE_MANAGEMENT
                    }
                }, {
                    label: this.$t('业务'),
                    route: {
                        name: MENU_RESOURCE_BUSINESS
                    }
                }, {
                    label: this.$t('已归档业务')
                }])
            },
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
                    name: this.$t('更新时间')
                }])
            },
            getTableData (event) {
                this.searchBusiness({
                    params: this.getSearchParams(),
                    config: {
                        globalPermission: false,
                        cancelPrevious: true,
                        requestId: 'searchArchivedBusiness'
                    }
                }).then(business => {
                    if (business.count && !business.info.length) {
                        this.pagination.current -= 1
                        this.getTableData()
                    }
                    this.pagination.count = business.count
                    this.list = this.$tools.flattenList(this.properties, business.info.map(biz => {
                        biz['last_time'] = this.$tools.formatTime(biz['last_time'], 'YYYY-MM-DD HH:mm:ss')
                        return biz
                    }))

                    if (event) {
                        this.table.stuff.type = 'search'
                    }
                }).catch(({ permission }) => {
                    if (permission) {
                        this.table.stuff = {
                            type: 'permission',
                            payload: { permission }
                        }
                    }
                })
            },
            getSearchParams () {
                const params = {
                    condition: {
                        'bk_data_status': 'disabled'
                    },
                    fields: [],
                    page: {
                        start: (this.pagination.current - 1) * this.pagination.limit,
                        limit: this.pagination.limit,
                        sort: '-bk_biz_id'
                    }
                }
                if (this.filter.range.length) {
                    params.condition.last_time = {
                        '$gte': this.filter.range[0],
                        '$lte': this.filter.range[1]
                    }
                }
                if (this.filter.name) {
                    params.condition.bk_biz_name = { '$regex': this.filter.name }
                }
                return params
            },
            handleRecovery (biz) {
                this.$bkInfo({
                    title: this.$t('是否确认恢复业务？'),
                    subTitle: this.$t('恢复业务提示', { bizName: biz['bk_biz_name'] }),
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
                    this.$success(this.$t('恢复业务成功'))
                    this.getTableData()
                })
            },
            handleSizeChange (size) {
                this.pagination.limit = size
                this.handlePageChange(1)
            },
            handlePageChange (current, event) {
                this.pagination.current = current
                this.getTableData(event)
            }
        }
    }
</script>

<style lang="scss" scoped>
    .archived-layout{
        padding: 0 20px;
    }
    .archived-filter {
        padding: 0 0 15px 0;
        .filter-item {
            width: 220px;
            margin-right: 5px;
            @include inlineBlock;
        }
    }
</style>
