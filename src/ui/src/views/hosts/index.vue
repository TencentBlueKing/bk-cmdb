<template>
    <div class="hosts-layout clearfix">
        <cmdb-hosts-filter class="hosts-filter fr"
            :filter-config-key="filter.filterConfigKey"
            :collection-content="{business: filter.business}"
            @on-refresh="handleRefresh">
            <div class="filter-group" slot="business">
                <label class="filter-label">{{$t('Hosts[\'选择业务\']')}}</label>
                <cmdb-business-selector class="filter-field" v-model="filter.business"></cmdb-business-selector>
            </div>
        </cmdb-hosts-filter>
        <cmdb-hosts-table class="hosts-main" ref="hostsTable"
            :columns-config-key="table.columnsConfigKey"
            :columns-config-properties="columnsConfigProperties">
        </cmdb-hosts-table>
    </div>
</template>

<script>
    import { mapGetters, mapActions } from 'vuex'
    import cmdbHostsFilter from '@/components/hosts/filter'
    import cmdbHostsTable from '@/components/hosts/table'
    export default {
        components: {
            cmdbHostsFilter,
            cmdbHostsTable
        },
        data () {
            return {
                properties: {
                    biz: [],
                    host: [],
                    set: [],
                    module: []
                },
                table: {
                    columnsConfigKey: 'hosts_table_columns'
                },
                filter: {
                    filterConfigKey: 'hosts_filter_fields',
                    business: null,
                    businessResolver: null,
                    params: null,
                    paramsResolver: null
                }
            }
        },
        computed: {
            ...mapGetters(['supplierAccount']),
            ...mapGetters('hostFavorites', ['applyingInfo']),
            columnsConfigProperties () {
                const setProperties = this.properties.set.filter(property => ['bk_set_name'].includes(property['bk_property_id']))
                const moduleProperties = this.properties.module.filter(property => ['bk_module_name'].includes(property['bk_property_id']))
                const hostProperties = this.properties.host
                return [...setProperties, ...moduleProperties, ...hostProperties]
            }
        },
        watch: {
            'filter.business' (business) {
                if (this.filter.businessResolver) {
                    this.filter.businessResolver()
                } else {
                    this.table.checked = []
                    this.getHostList()
                }
            },
            applyingInfo (info) {
                if (info) {
                    this.filter.business = info['bk_biz_id']
                }
            }
        },
        async created () {
            try {
                await Promise.all([
                    this.getBusiness(),
                    this.getParams(),
                    this.getProperties()
                ])
                this.getHostList()
            } catch (e) {
                console.log(e)
            }
        },
        beforeRouteUpdate (to, from, next) {
            this.$store.commit('hostFavorites/setApplying', null)
            next()
        },
        beforeRouteLeave (to, from, next) {
            this.$store.commit('hostFavorites/setApplying', null)
            next()
        },
        methods: {
            ...mapActions('objectModelProperty', ['batchSearchObjectAttribute']),
            getBusiness () {
                const query = this.$route.query
                if (query.hasOwnProperty('business')) {
                    this.filter.business = parseInt(query.business)
                    return Promise.resolve()
                }
                return new Promise((resolve, reject) => {
                    this.filter.businessResolver = () => {
                        this.filter.businessResolver = null
                        resolve()
                    }
                })
            },
            getParams () {
                return new Promise((resolve, reject) => {
                    this.filter.paramsResolver = () => {
                        this.filter.paramsResolver = null
                        resolve()
                    }
                })
            },
            getProperties () {
                return this.batchSearchObjectAttribute({
                    params: {
                        bk_obj_id: {'$in': Object.keys(this.properties)},
                        bk_supplier_account: this.supplierAccount
                    },
                    config: {
                        requestId: `post_batchSearchObjectAttribute_${Object.keys(this.properties).join('_')}`,
                        requestGroup: Object.keys(this.properties).map(id => `post_searchObjectAttribute_${id}`),
                        fromCache: true
                    }
                }).then(result => {
                    Object.keys(this.properties).forEach(objId => {
                        this.properties[objId] = result[objId]
                    })
                    return result
                })
            },
            handleRefresh (params) {
                this.filter.params = params
                if (this.filter.paramsResolver) {
                    this.filter.paramsResolver()
                } else {
                    this.getHostList()
                }
            },
            getHostList () {
                this.$refs.hostsTable.search(this.filter.business, this.filter.params)
            }
        }
    }
</script>

<style lang="scss" scoped>
    .hosts-layout{
        height: 100%;
        padding: 0;
        overflow: hidden;
        .hosts-main{
            height: 100%;
            padding: 20px;
            overflow: hidden;
        }
        .hosts-filter{
            height: 100%;
        }
    }
    .hosts-options{
        font-size: 0;
        .options-button{
            position: relative;
            display: inline-block;
            vertical-align: middle;
            border-radius: 0;
            font-size: 14px;
            margin-left: -1px;
            &:hover{
                z-index: 1;
            }
        }
    }
    .hosts-table{
        margin-top: 20px;
    }
</style>