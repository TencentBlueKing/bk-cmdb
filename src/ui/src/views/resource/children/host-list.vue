<template>
    <div class="resource-layout">
        <host-list-options></host-list-options>
        <host-filter-tag class="filter-tag" ref="filterTag"></host-filter-tag>
        <bk-table class="hosts-table"
            ref="table"
            v-bkloading="{ isLoading: $loading(Object.values(request)) }"
            :data="table.list"
            :pagination="table.pagination"
            :max-height="$APP.height - filtersTagHeight - 230"
            @selection-change="handleSelectionChange"
            @sort-change="handleSortChange"
            @page-change="handlePageChange"
            @page-limit-change="handleSizeChange"
            @header-click="handleHeaderClick">
            <bk-table-column type="selection" width="60" align="center" fixed class-name="bk-table-selection"></bk-table-column>
            <bk-table-column v-for="property in tableHeader"
                :show-overflow-tooltip="property.bk_property_type !== 'topology'"
                :min-width="getColumnMinWidth(property)"
                :key="property.id"
                :sortable="isPropertySortable(property) ? 'custom' : false"
                :prop="property.bk_property_id"
                :fixed="['bk_host_id'].includes(property.bk_property_id)"
                :class-name="['bk_host_id'].includes(property.bk_property_id) ? 'is-highlight' : ''"
                :render-header="() => renderHeader(property)">
                <template slot-scope="{ row }">
                    <cmdb-host-topo-path
                        v-if="property.bk_property_type === 'topology'"
                        :host="row"
                        @path-ready="handlePathReady(row, ...arguments)">
                    </cmdb-host-topo-path>
                    <cmdb-property-value
                        v-else
                        :theme="['bk_host_id'].includes(property.bk_property_id) ? 'primary' : 'default'"
                        :value="row | hostValueFilter(property.bk_obj_id, property.bk_property_id)"
                        :show-unit="false"
                        :property="property"
                        @click.native.stop="handleValueClick(row, property)">
                    </cmdb-property-value>
                </template>
            </bk-table-column>
            <bk-table-column type="setting"></bk-table-column>
            <cmdb-table-empty slot="empty" :stuff="table.stuff"></cmdb-table-empty>
        </bk-table>
    </div>
</template>

<script>
    import { mapGetters } from 'vuex'
    import hostListOptions from './host-options.vue'
    import hostValueFilter from '@/filters/host'
    import {
        MENU_RESOURCE_HOST,
        MENU_RESOURCE_HOST_DETAILS,
        MENU_RESOURCE_BUSINESS_HOST_DETAILS
    } from '@/dictionary/menu-symbol'
    import RouterQuery from '@/router/query'
    import CmdbHostTopoPath from '@/components/host-topo-path/host-topo-path.vue'
    import HostStore from '../transfer/host-store'
    import HostFilterTag from '@/components/filters/filter-tag'
    import FilterStore, { setupFilterStore } from '@/components/filters/store'
    import ColumnsConfig from '@/components/columns-config/columns-config.js'
    export default {
        components: {
            hostListOptions,
            CmdbHostTopoPath,
            HostFilterTag
        },
        filters: {
            hostValueFilter
        },
        data () {
            return {
                directory: null,
                scope: 1,
                table: {
                    checked: [],
                    selection: [],
                    list: [],
                    pagination: {
                        current: 1,
                        count: 0,
                        ...this.$tools.getDefaultPaginationConfig()
                    },
                    sort: 'bk_host_id',
                    exportUrl: `${window.API_HOST}hosts/export`,
                    stuff: {
                        type: 'default',
                        payload: {
                            emptyText: this.$t('bk.table.emptyText')
                        }
                    }
                },
                request: {
                    list: Symbol('list')
                },
                filtersTagHeight: 0
            }
        },
        computed: {
            ...mapGetters(['userName']),
            ...mapGetters('resourceHost', ['activeDirectory']),
            ...mapGetters('objectModelClassify', ['getModelById']),
            moduleProperties () {
                return FilterStore.getModelProperties('module')
            },
            tableHeader () {
                return FilterStore.header
            }
        },
        watch: {
            scope () {
                this.setModuleNamePropertyState()
            }
        },
        async created () {
            try {
                setupFilterStore({
                    header: {
                        custom: this.$route.meta.customInstanceColumn,
                        global: 'host_global_custom_table_columns'
                    }
                })
                this.setModuleNamePropertyState()
                this.unwatchRouter = RouterQuery.watch('*', ({
                    scope = 1,
                    page = 1,
                    sort = 'bk_host_id',
                    limit = this.table.pagination.limit,
                    directory = null
                }) => {
                    if (this.$route.name !== MENU_RESOURCE_HOST) {
                        return false
                    }
                    this.table.pagination.current = parseInt(page)
                    this.table.pagination.limit = parseInt(limit)
                    this.table.sort = sort
                    this.directory = parseInt(directory) || null
                    this.scope = isNaN(scope) ? 'all' : parseInt(scope)
                    this.getHostList()
                }, { throttle: 100 })
                this.unwatchScopeAndDirectory = RouterQuery.watch(['scope', 'directory'], FilterStore.resetAll)
            } catch (error) {
                console.error(error)
            }
        },
        mounted () {
            this.unwatchFilter = this.$watch(() => {
                return [FilterStore.condition, FilterStore.IP]
            }, () => {
                const el = this.$refs.filterTag.$el
                if (el.getBoundingClientRect) {
                    this.filtersTagHeight = el.getBoundingClientRect().height
                } else {
                    this.filtersTagHeight = 0
                }
            }, { immediate: true, deep: true })
            this.disabledTableSettingDefaultBehavior()
        },
        beforeDestroy () {
            this.unwatchRouter()
            this.unwatchFilter()
            this.unwatchScopeAndDirectory()
        },
        methods: {
            disabledTableSettingDefaultBehavior () {
                setTimeout(() => {
                    const settingReference = this.$refs.table.$el.querySelector('.bk-table-column-setting .bk-tooltip-ref')
                    settingReference && settingReference._tippy && settingReference._tippy.disable()
                }, 1000)
            },
            setModuleNamePropertyState () {
                const property = this.moduleProperties.find(property => property.bk_property_id === 'bk_module_name')
                if (property) {
                    const normalName = this.$t('模块名')
                    const directoryName = this.$t('目录名')
                    const scopeModuleName = {
                        0: normalName,
                        1: directoryName,
                        all: `${directoryName}/${normalName}`
                    }
                    property.bk_property_name = scopeModuleName[this.scope]
                }
            },
            getColumnMinWidth (property) {
                if (property.bk_property_type === 'topology') {
                    return 200
                }
                return 100
            },
            isPropertySortable (property) {
                return property.bk_obj_id === 'host' && !['foreignkey', 'topology'].includes(property.bk_property_type)
            },
            renderHeader (property) {
                const content = [this.$tools.getHeaderPropertyName(property)]
                const modelId = property.bk_obj_id
                if (modelId !== 'host') {
                    const model = this.getModelById(modelId)
                    const suffix = this.$createElement('span', { style: { color: '#979BA5', marginLeft: '4px' } }, [`(${model.bk_obj_name})`])
                    content.push(suffix)
                }
                return this.$createElement('span', {}, content)
            },
            async getHostList (event) {
                try {
                    const { count, info } = await this.$store.dispatch('hostSearch/searchHost', {
                        params: this.getParams(),
                        config: {
                            requestId: this.request.list,
                            cancelPrevious: true
                        }
                    })
                    this.table.pagination.count = count
                    this.table.list = info
                    this.table.stuff.type = event ? 'search' : 'default'
                } catch (error) {
                    this.table.pagination.count = 0
                    this.table.checked = []
                    this.table.list = []
                    console.error(error)
                }
            },
            getParams () {
                const params = {
                    ...FilterStore.getSearchParams(),
                    page: {
                        ...this.$tools.getPageParams(this.table.pagination),
                        sort: this.table.sort
                    }
                }
                this.injectScope(params)
                this.scope === 1 && this.injectDirectory(params)
                return params
            },
            injectScope (params) {
                const biz = params.condition.find(condition => condition.bk_obj_id === 'biz')
                if (this.scope === 'all') {
                    biz.condition = biz.condition.filter(condition => condition.field !== 'default')
                } else {
                    const newMeta = {
                        field: 'default',
                        operator: '$eq',
                        value: this.scope
                    }
                    const existMeta = biz.condition.find(({ field, operator }) => field === newMeta.field && operator === newMeta.operator)
                    if (existMeta) {
                        existMeta.value = newMeta.value
                    } else {
                        biz.condition.push(newMeta)
                    }
                }
                return params
            },
            injectDirectory (params) {
                if (!this.directory) {
                    return false
                }
                const moduleCondition = params.condition.find(condition => condition.bk_obj_id === 'module')
                const directoryMeta = {
                    field: 'bk_module_id',
                    operator: '$eq',
                    value: this.directory
                }
                const existMeta = moduleCondition.condition.find(({ field, operator }) => field === directoryMeta.field && operator === directoryMeta.operator)
                if (existMeta) {
                    existMeta.value = directoryMeta.value
                } else {
                    moduleCondition.condition.push(directoryMeta)
                }
            },
            handleSelectionChange (selection) {
                this.table.selection = selection
                this.table.checked = selection.map(item => item.host.bk_host_id)
                HostStore.setSelected(selection)
            },
            handleValueClick (item, property) {
                if (property.bk_obj_id !== 'host' || property.bk_property_id !== 'bk_host_id') {
                    return
                }
                const business = item.biz[0]
                if (business.default) {
                    this.$routerActions.redirect({
                        name: MENU_RESOURCE_HOST_DETAILS,
                        params: {
                            id: item.host.bk_host_id
                        },
                        query: {
                            from: 'resource'
                        },
                        history: true
                    })
                } else {
                    this.$routerActions.redirect({
                        name: MENU_RESOURCE_BUSINESS_HOST_DETAILS,
                        params: {
                            business: business.bk_biz_id,
                            id: item.host.bk_host_id
                        },
                        query: {
                            from: 'resource'
                        },
                        history: true
                    })
                }
            },
            handlePageChange (current) {
                RouterQuery.set({
                    page: current,
                    _t: Date.now()
                })
            },
            handleSizeChange (limit) {
                RouterQuery.set({
                    limit: limit,
                    page: 1,
                    _t: Date.now()
                })
            },
            handleSortChange (sort) {
                RouterQuery.set({
                    sort: this.$tools.getSort(sort),
                    _t: Date.now()
                })
            },
            // 拓扑路径写入数据中，用于复制
            handlePathReady (row, paths) {
                row.__bk_host_topology__ = paths
            },
            handleHeaderClick (column) {
                if (column.type !== 'setting') {
                    return false
                }
                ColumnsConfig.open({
                    props: {
                        properties: FilterStore.properties.filter(property => {
                            return property.bk_obj_id === 'host'
                                || (property.bk_obj_id === 'module' && property.bk_property_id === 'bk_module_name')
                                || (property.bk_obj_id === 'set' && property.bk_property_id === 'bk_set_name')
                        }),
                        selected: FilterStore.defaultHeader.map(property => property.bk_property_id),
                        disabledColumns: ['bk_host_id', 'bk_host_innerip', 'bk_cloud_id']
                    },
                    handler: {
                        apply: async properties => {
                            await this.handleApplyColumnsConfig(properties)
                            FilterStore.setHeader(properties)
                            FilterStore.dispatchSearch()
                        },
                        reset: async () => {
                            await this.handleApplyColumnsConfig()
                            FilterStore.setHeader(FilterStore.defaultHeader)
                            FilterStore.dispatchSearch()
                        }
                    }
                })
            },
            handleApplyColumnsConfig (properties = []) {
                return this.$store.dispatch('userCustom/saveUsercustom', {
                    [this.$route.meta.customInstanceColumn]: properties.map(property => property.bk_property_id)
                })
            }
        }
    }
</script>

<style lang="scss" scoped>
    .filter-tag ~ .hosts-table {
        margin-top: 0;
    }
    .hosts-table {
        margin-top: 10px;
    }
</style>
