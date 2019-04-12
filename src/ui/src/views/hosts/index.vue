<template>
    <div class="hosts-layout clearfix">
        <cmdb-hosts-filter class="hosts-filter fr"
            :filter-config-key="filterConfigKey"
            :collection-content="{business: filter.business}"
            @on-refresh="handleRefresh">
        </cmdb-hosts-filter>
        <cmdb-hosts-table class="hosts-main" ref="hostsTable"
            :columns-config-key="columnsConfigKey"
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
                filter: {
                    business: null,
                    businessResolver: null,
                    params: null,
                    paramsResolver: null
                }
            }
        },
        computed: {
            ...mapGetters(['supplierAccount', 'userName', 'isAdminView']),
            ...mapGetters('hostFavorites', ['applyingInfo']),
            ...mapGetters('objectBiz', ['bizId']),
            columnsConfigKey () {
                return `${this.userName}_host_${this.isAdminView ? 'adminView' : this.bizId}_table_columns`
            },
            filterConfigKey () {
                return `${this.userName}_host_${this.isAdminView ? 'adminView' : this.bizId}_filter_fields`
            },
            columnsConfigProperties () {
                const setProperties = this.properties.set.filter(property => ['bk_set_name'].includes(property['bk_property_id']))
                const moduleProperties = this.properties.module.filter(property => ['bk_module_name'].includes(property['bk_property_id']))
                const hostProperties = this.properties.host
                return [...setProperties, ...moduleProperties, ...hostProperties]
            }
        },
        watch: {
            bizId () {
                this.table.checked = []
                this.getHostList()
            },
            applyingInfo (info) {
                if (info) {
                    this.filter.business = info['bk_biz_id']
                }
            }
        },
        async created () {
            this.$store.commit('setHeaderTitle', this.$t('Nav["主机查询"]'))
            try {
                const res = await Promise.all([
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
                    this.$store.commit('objectBiz/setBizId', parseInt(query.business))
                    return Promise.resolve()
                }
                return Promise.resolve()
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
                    params: this.$injectMetadata({
                        bk_obj_id: {'$in': Object.keys(this.properties)},
                        bk_supplier_account: this.supplierAccount
                    }),
                    config: {
                        requestId: `post_batchSearchObjectAttribute_${Object.keys(this.properties).join('_')}`,
                        requestGroup: Object.keys(this.properties).map(id => `post_searchObjectAttribute_${id}`)
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
            getHostList (resetPage = true) {
                this.$refs.hostsTable.search(this.bizId, this.filter.params, resetPage)
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