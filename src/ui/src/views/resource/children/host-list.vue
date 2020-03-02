<template>
    <div class="resource-layout">
        <host-list-options></host-list-options>
        <bk-table class="hosts-table"
            v-bkloading="{ isLoading: $loading(['searchHosts', 'batchSearchProperties', 'post_searchGroup_host']) }"
            :data="table.list"
            :pagination="table.pagination"
            :row-style="{ cursor: 'pointer' }"
            :max-height="$APP.height - 230"
            @selection-change="handleSelectionChange"
            @row-click="handleRowClick"
            @sort-change="handleSortChange"
            @page-change="handlePageChange"
            @page-limit-change="handleSizeChange">
            <bk-table-column type="selection" width="60" align="center" fixed class-name="bk-table-selection"></bk-table-column>
            <bk-table-column v-for="column in table.header"
                :key="column.id"
                :label="column.name"
                :sortable="column.sortable ? 'custom' : false"
                :prop="column.id"
                :fixed="column.id === 'bk_host_innerip'"
                :class-name="column.id === 'bk_host_innerip' ? 'is-highlight' : ''">
                <template slot-scope="{ row }">
                    {{ row | hostValueFilter(column.objId, column.id) | formatter(column.type, getPropertyValue(column.objId, column.id, 'option'))}}
                </template>
            </bk-table-column>
            <cmdb-table-empty slot="empty" :stuff="table.stuff"></cmdb-table-empty>
        </bk-table>
    </div>
</template>

<script>
    import { mapState, mapGetters } from 'vuex'
    import hostListOptions from './host-options.vue'
    import hostValueFilter from '@/filters/host'
    import {
        MENU_RESOURCE_HOST_DETAILS,
        MENU_RESOURCE_BUSINESS_HOST_DETAILS
    } from '@/dictionary/menu-symbol'
    import Bus from '@/utils/bus.js'
    export default {
        components: {
            hostListOptions
        },
        filters: {
            hostValueFilter
        },
        data () {
            return {
                properties: {
                    biz: [],
                    host: [],
                    set: [],
                    module: []
                },
                propertyGroups: [],
                table: {
                    checked: [],
                    header: Array(8).fill({}),
                    list: [],
                    pagination: {
                        current: 1,
                        count: 0,
                        ...this.$tools.getDefaultPaginationConfig()
                    },
                    defaultSort: 'bk_host_id',
                    sort: 'bk_host_id',
                    exportUrl: `${window.API_HOST}hosts/export`,
                    stuff: {
                        type: 'default',
                        payload: {
                            emptyText: this.$t('bk.table.emptyText')
                        }
                    }
                },
                columnsConfig: {
                    selected: []
                },
                columnsConfigDisabledColumns: ['bk_host_innerip', 'bk_cloud_id', 'bk_biz_name', 'bk_module_name'],
                ready: false,
                scope: 1
            }
        },
        computed: {
            ...mapState('hosts', ['filterParams']),
            ...mapGetters(['userName']),
            ...mapGetters('userCustom', ['usercustom']),
            ...mapGetters('resourceHost', ['activeDirectory']),
            columnsConfigKey () {
                return `${this.userName}_$resource_adminView_table_columns`
            },
            customColumns () {
                return this.usercustom[this.columnsConfigKey] || []
            },
            columnsConfigProperties () {
                const setProperties = this.properties.set.filter(property => ['bk_set_name'].includes(property['bk_property_id']))
                const moduleProperties = this.properties.module.filter(property => ['bk_module_name'].includes(property['bk_property_id']))
                const businessProperties = this.properties.biz.filter(property => ['bk_biz_name'].includes(property['bk_property_id']))
                const hostProperties = this.properties.host
                return [...setProperties, ...moduleProperties, ...businessProperties, ...hostProperties]
            },
            filterCondition () {
                const defaultModel = ['biz', 'set', 'module', 'host']
                const params = {
                    bk_biz_id: -1,
                    ip: this.filterParams.ip,
                    condition: defaultModel.map(model => {
                        return {
                            bk_obj_id: model,
                            condition: this.filterParams[model] || [],
                            fields: []
                        }
                    })
                }
                if (this.activeDirectory && this.scope === 1) {
                    const moduleCondition = params.condition.find(target => target.bk_obj_id === 'module')
                    moduleCondition.condition.push({
                        field: 'bk_module_id',
                        operator: '$eq',
                        value: this.activeDirectory.bk_inst_id
                    })
                }
                return params
            }
        },
        watch: {
            customColumns () {
                this.setTableHeader()
            },
            columnsConfigProperties () {
                this.setTableHeader()
            },
            filterParams () {
                this.ready && this.getHostList(true)
            }
        },
        async created () {
            try {
                Bus.$on('refresh-resource-list', this.handlePageChange)
                await Promise.all([
                    this.getProperties(),
                    this.getHostPropertyGroups()
                ])
                this.getHostList()
                this.ready = true
            } catch (e) {
                console.log(e)
            }
        },
        beforeDestroy () {
            this.ready = false
            Bus.$off('refresh-resource-list', this.handlePageChange)
        },
        methods: {
            getPropertyValue (modelId, propertyId, field) {
                const model = this.properties[modelId]
                if (!model) {
                    return ''
                }
                const curProperty = model.find(property => property.bk_property_id === propertyId)
                return curProperty ? curProperty[field] : ''
            },
            getProperties () {
                return this.$store.dispatch('objectModelProperty/batchSearchObjectAttribute', {
                    params: this.$injectMetadata({
                        bk_obj_id: { '$in': Object.keys(this.properties) },
                        bk_supplier_account: this.supplierAccount
                    }, { inject: this.$route.name !== 'resource' }),
                    config: {
                        requestId: 'batchSearchProperties'
                    }
                }).then(result => {
                    Object.keys(this.properties).forEach(objId => {
                        this.properties[objId] = result[objId]
                    })
                    return result
                })
            },
            getHostPropertyGroups () {
                return this.$store.dispatch('objectModelFieldGroup/searchGroup', {
                    objId: 'host',
                    params: this.$injectMetadata(),
                    config: {
                        requestId: 'post_searchGroup_host'
                    }
                }).then(groups => {
                    this.propertyGroups = groups
                    return groups
                })
            },
            setTableHeader () {
                const customColumns = this.customColumns.length ? this.customColumns : this.globalCustomColumns
                const properties = this.$tools.getHeaderProperties(this.columnsConfigProperties, customColumns, this.columnsConfigDisabledColumns)
                this.table.header = properties.map(property => {
                    return {
                        id: property.bk_property_id,
                        name: this.$tools.getHeaderPropertyName(property),
                        type: property.bk_property_type,
                        objId: property.bk_obj_id,
                        sortable: property.bk_obj_id === 'host' && !['foreignkey'].includes(property.bk_property_type)
                    }
                })
                this.columnsConfig.selected = properties.map(property => property['bk_property_id'])
            },
            getHostList (event) {
                this.$store.dispatch('hostSearch/searchHost', {
                    params: this.injectScope({
                        ...this.filterCondition,
                        page: {
                            start: (this.table.pagination.current - 1) * this.table.pagination.limit,
                            limit: this.table.pagination.limit,
                            sort: this.table.sort
                        }
                    }),
                    config: {
                        requestId: 'searchHosts',
                        cancelPrevious: true
                    }
                }).then(data => {
                    this.table.pagination.count = data.count
                    this.table.list = data.info

                    if (event) {
                        this.table.stuff.type = 'search'
                    }

                    return data
                }).catch(e => {
                    this.table.checked = []
                    this.table.list = []
                    this.table.pagination.count = 0
                })
            },
            injectScope (params) {
                const biz = params.condition.find(condition => condition.bk_obj_id === 'biz')
                if (this.scope === 'all') {
                    biz.condition = biz.condition.filter(condition => condition.field !== 'default')
                } else {
                    const newCondition = {
                        field: 'default',
                        operator: '$eq',
                        value: this.scope
                    }
                    const existCondition = biz.condition.find(condition => condition.field === 'default')
                    if (existCondition) {
                        Object.assign(existCondition, newCondition)
                    } else {
                        biz.condition.push(newCondition)
                    }
                }
                return params
            },
            handleSelectionChange (selection) {
                this.table.checked = selection.map(item => item.host.bk_host_id)
            },
            handleRowClick (item) {
                const business = item.biz[0]
                if (business.default) {
                    this.$router.push({
                        name: MENU_RESOURCE_HOST_DETAILS,
                        params: {
                            id: item.host.bk_host_id
                        },
                        query: {
                            from: 'resource'
                        }
                    })
                } else {
                    this.$router.push({
                        name: MENU_RESOURCE_BUSINESS_HOST_DETAILS,
                        params: {
                            business: business.bk_biz_id,
                            id: item.host.bk_host_id
                        },
                        query: {
                            from: 'resource'
                        }
                    })
                }
            },
            handlePageChange (current, event) {
                this.table.pagination.current = current
                this.getHostList(event)
            },
            handleSizeChange (limit) {
                this.table.pagination.limit = limit
                this.handlePageChange(1)
            },
            handleSortChange (sort) {
                this.table.sort = this.$tools.getSort(sort)
                this.handlePageChange(1)
            }
        }
    }
</script>

<style lang="scss" scoped>
</style>
