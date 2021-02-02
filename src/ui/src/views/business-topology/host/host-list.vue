<template>
    <div class="list-layout">
        <host-list-options @transfer="handleTransfer"></host-list-options>
        <host-filter-tag class="filter-tag" ref="filterTag"></host-filter-tag>
        <bk-table class="host-table"
            ref="table"
            v-bkloading="{ isLoading: $loading(Object.values(request)) }"
            :data="table.data"
            :pagination="table.pagination"
            :max-height="$APP.height - filtersTagHeight - 250"
            @page-change="handlePageChange"
            @page-limit-change="handleLimitChange"
            @sort-change="handleSortChange"
            @selection-change="handleSelectionChange"
            @header-click="handleHeaderClick">
            <bk-table-column type="selection" width="50" align="center" fixed></bk-table-column>
            <bk-table-column v-for="column in tableHeader"
                show-overflow-tooltip
                :min-width="column.bk_property_id === 'bk_host_id' ? 80 : 120"
                :key="column.bk_property_id"
                :sortable="getColumnSortable(column)"
                :prop="column.bk_property_id"
                :fixed="column.bk_property_id === 'bk_host_id'"
                :render-header="() => renderHeader(column)">
                <template slot-scope="{ row }">
                    <cmdb-property-value
                        :theme="column.bk_property_id === 'bk_host_id' ? 'primary' : 'default'"
                        :value="row | hostValueFilter(column.bk_obj_id, column.bk_property_id)"
                        :show-unit="false"
                        :property="column"
                        @click.native.stop="handleValueClick(row, column)">
                    </cmdb-property-value>
                </template>
            </bk-table-column>
            <bk-table-column type="setting"></bk-table-column>
        </bk-table>
        <cmdb-dialog v-model="dialog.show" :width="dialog.width" :height="dialog.height">
            <component
                :is="dialog.component"
                v-bind="dialog.props"
                @cancel="handleDialogCancel"
                @confirm="handleDialogConfirm">
            </component>
        </cmdb-dialog>
    </div>
</template>

<script>
    import HostListOptions from './host-list-options.vue'
    import ModuleSelector from './module-selector.vue'
    import AcrossBusinessModuleSelector from './across-business-module-selector.vue'
    import MoveToResourceConfirm from './move-to-resource-confirm.vue'
    import hostValueFilter from '@/filters/host'
    import { mapGetters } from 'vuex'
    import {
        MENU_BUSINESS_HOST_AND_SERVICE,
        MENU_BUSINESS_HOST_DETAILS,
        MENU_BUSINESS_TRANSFER_HOST
    } from '@/dictionary/menu-symbol'
    import Bus from '@/utils/bus.js'
    import RouterQuery from '@/router/query'
    import HostFilterTag from '@/components/filters/filter-tag'
    import FilterStore, { setupFilterStore } from '@/components/filters/store'
    import ColumnsConfig from '@/components/columns-config/columns-config.js'
    export default {
        components: {
            HostListOptions,
            HostFilterTag,
            [ModuleSelector.name]: ModuleSelector,
            [AcrossBusinessModuleSelector.name]: AcrossBusinessModuleSelector,
            [MoveToResourceConfirm.name]: MoveToResourceConfirm
        },
        filters: {
            hostValueFilter
        },
        props: {
            active: Boolean
        },
        data () {
            return {
                commonRequestFinished: false,
                table: {
                    data: [],
                    selection: [],
                    sort: 'bk_host_id',
                    pagination: this.$tools.getDefaultPaginationConfig()
                },
                dialog: {
                    width: 830,
                    height: 600,
                    show: false,
                    component: null,
                    props: {}
                },
                request: {
                    table: Symbol('table'),
                    moveToResource: Symbol('moveToResource'),
                    moveToIdleModule: Symbol('moveToIdleModule')
                },
                filtersTagHeight: 0
            }
        },
        computed: {
            ...mapGetters('objectBiz', ['bizId', 'currentBusiness']),
            ...mapGetters('objectModelClassify', ['getModelById']),
            ...mapGetters('businessHost', [
                'columnsConfigProperties',
                'selectedNode',
                'commonRequest'
            ]),
            tableHeader () {
                return FilterStore.header
            }
        },
        created () {
            setupFilterStore({
                bk_biz_id: this.bizId,
                header: {
                    custom: this.$route.meta.customInstanceColumn,
                    global: 'host_global_custom_table_columns'
                }
            })
            this.unwatchRouter = RouterQuery.watch('*', ({
                tab = 'hostList',
                node,
                page = 1,
                limit = this.table.pagination.limit
            }) => {
                if (this.$route.name !== MENU_BUSINESS_HOST_AND_SERVICE) {
                    return false
                }
                this.table.pagination.current = parseInt(page)
                this.table.pagination.limit = parseInt(limit)
                tab === 'hostList' && node && this.selectedNode && this.getHostList()
            }, { throttle: 16, ignore: ['keyword'] })
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
        },
        methods: {
            disabledTableSettingDefaultBehavior () {
                setTimeout(() => {
                    const settingReference = this.$refs.table.$el.querySelector('.bk-table-column-setting .bk-tooltip-ref')
                    settingReference && settingReference._tippy && settingReference._tippy.disable()
                }, 1000)
            },
            getColumnSortable (column) {
                const isHostProperty = column.bk_obj_id === 'host'
                const isForeignKey = column.bk_property_type === 'foreignkey'
                return (isHostProperty && !isForeignKey) ? 'custom' : false
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
            handlePageChange (current = 1) {
                RouterQuery.set({
                    page: current,
                    _t: Date.now()
                })
            },
            handleLimitChange (limit) {
                RouterQuery.set({
                    limit: limit,
                    page: 1,
                    _t: Date.now()
                })
            },
            handleSortChange (sort) {
                this.table.sort = this.$tools.getSort(sort)
                RouterQuery.set('_t', Date.now())
            },
            handleValueClick (row, column) {
                if (column.bk_obj_id !== 'host' || column.bk_property_id !== 'bk_host_id') {
                    return
                }
                this.$routerActions.redirect({
                    name: MENU_BUSINESS_HOST_DETAILS,
                    params: {
                        bizId: this.bizId,
                        id: row.host.bk_host_id
                    },
                    history: true
                })
            },
            handleSelectionChange (selection) {
                this.table.selection = selection
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
            },
            async getHostList () {
                try {
                    await this.commonRequest
                    this.commonRequestFinished = true
                    const result = await this.$store.dispatch('hostSearch/searchHost', {
                        params: this.getParams(),
                        config: {
                            requestId: this.request.table,
                            cancelPrevious: true
                        }
                    })
                    this.table.data = result.info
                    this.table.pagination.count = result.count
                } catch (e) {
                    console.error(e)
                    this.table.data = []
                    this.table.pagination.count = 0
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
                const topoNodeData = this.selectedNode.data
                const fieldMap = {
                    'biz': 'bk_biz_id',
                    'set': 'bk_set_id',
                    'module': 'bk_module_id'
                }
                const topoCondition = {
                    field: fieldMap[topoNodeData.bk_obj_id] || 'bk_inst_id',
                    operator: '$eq',
                    value: topoNodeData.bk_inst_id
                }
                const modelConditionId = fieldMap.hasOwnProperty(topoNodeData.bk_obj_id) ? topoNodeData.bk_obj_id : 'object'
                const modelCondition = params.condition.find(modelCondition => modelCondition.bk_obj_id === modelConditionId)
                modelCondition.condition.push(topoCondition)
                return params
            },
            handleTransfer (type) {
                const actionMap = {
                    idle: this.openModuleSelector,
                    business: this.openModuleSelector,
                    acrossBusiness: this.openAcrollBusinessModuleSelector,
                    resource: this.openResourceConfirm
                }
                actionMap[type] && actionMap[type](type)
            },
            openModuleSelector (type) {
                const props = {
                    moduleType: type,
                    business: this.currentBusiness
                }
                if (type === 'idle') {
                    props.title = this.$t('转移主机到空闲模块')
                } else {
                    props.title = this.$t('转移主机到业务模块')
                    const selection = this.table.selection
                    const firstSelectionModules = selection[0].module.map(module => module.bk_module_id).sort()
                    const firstSelectionModulesStr = firstSelectionModules.join(',')
                    const allSame = selection.slice(1).every(item => {
                        const modules = item.module.map(module => module.bk_module_id).sort().join(',')
                        return modules === firstSelectionModulesStr
                    })
                    if (allSame) {
                        props.previousModules = firstSelectionModules
                    }
                }
                this.dialog.props = props
                this.dialog.width = 830
                this.dialog.height = 600
                this.dialog.component = ModuleSelector.name
                this.dialog.show = true
            },
            openResourceConfirm () {
                this.dialog.props = {
                    count: this.table.selection.length,
                    bizId: this.bizId
                }
                this.dialog.width = 400
                this.dialog.height = 250
                this.dialog.component = MoveToResourceConfirm.name
                this.dialog.show = true
            },
            openAcrollBusinessModuleSelector () {
                this.dialog.props = {
                    title: this.$t('转移主机到其他业务'),
                    business: this.currentBusiness
                }
                this.dialog.width = 830
                this.dialog.height = 600
                this.dialog.component = AcrossBusinessModuleSelector.name
                this.dialog.show = true
            },
            handleDialogCancel () {
                this.dialog.show = false
            },
            handleDialogConfirm () {
                this.dialog.show = false
                if (this.dialog.component === ModuleSelector.name) {
                    if (this.dialog.props.moduleType === 'idle') {
                        const isAllIdleSetHost = this.table.selection.every(data => {
                            const modules = data.module
                            return modules.every(module => module.default !== 0)
                        })
                        if (isAllIdleSetHost) {
                            this.transferDirectly(...arguments)
                        } else {
                            this.gotoTransferPage(...arguments)
                        }
                    } else {
                        this.gotoTransferPage(...arguments)
                    }
                } else if (this.dialog.component === MoveToResourceConfirm.name) {
                    this.moveHostToResource(...arguments)
                } else if (this.dialog.component === AcrossBusinessModuleSelector.name) {
                    this.moveHostToOtherBusiness(...arguments)
                }
            },
            async transferDirectly (modules) {
                try {
                    const internalModule = modules[0]
                    const selectedNode = this.selectedNode
                    await this.$http.post(
                        `host/transfer_with_auto_clear_service_instance/bk_biz_id/${this.bizId}`, {
                            bk_host_ids: this.table.selection.map(data => data.host.bk_host_id),
                            default_internal_module: internalModule.data.bk_inst_id,
                            remove_from_node: {
                                bk_inst_id: selectedNode.data.bk_inst_id,
                                bk_obj_id: selectedNode.data.bk_obj_id
                            }
                        }, {
                            requestId: this.request.moveToIdleModule
                        }
                    )
                    Bus.$emit('refresh-count', {
                        type: 'host_count',
                        hosts: [...this.table.selection],
                        target: internalModule
                    })
                    this.table.selection = []
                    this.$success('转移成功')
                    RouterQuery.set({
                        _t: Date.now(),
                        page: 1
                    })
                } catch (e) {
                    console.error(e)
                }
            },
            gotoTransferPage (modules) {
                const query = {
                    sourceModel: this.selectedNode.data.bk_obj_id,
                    sourceId: this.selectedNode.data.bk_inst_id,
                    targetModules: modules.map(node => node.data.bk_inst_id).join(','),
                    resources: this.table.selection.map(item => item.host.bk_host_id).join(','),
                    node: this.selectedNode.id
                }
                this.$routerActions.redirect({
                    name: MENU_BUSINESS_TRANSFER_HOST,
                    params: {
                        type: this.dialog.props.moduleType
                    },
                    query: query,
                    history: true
                })
            },
            async moveHostToResource (directoryId) {
                try {
                    await this.$store.dispatch('hostRelation/transferHostToResourceModule', {
                        params: {
                            bk_biz_id: this.bizId,
                            bk_host_id: this.table.selection.map(item => item.host.bk_host_id),
                            bk_module_id: directoryId
                        },
                        config: {
                            requestId: this.request.moveToResource
                        }
                    })
                    this.refreshHost()
                } catch (e) {
                    console.error(e)
                }
            },
            async moveHostToOtherBusiness (modules, targetBizId) {
                try {
                    const [targetModule] = modules
                    await this.$http.post('hosts/modules/across/biz', {
                        src_bk_biz_id: this.bizId,
                        dst_bk_biz_id: targetBizId,
                        bk_host_id: this.table.selection.map(({ host }) => host.bk_host_id),
                        bk_module_id: targetModule.data.bk_inst_id
                    })
                    this.refreshHost()
                } catch (error) {
                    console.error(error)
                }
            },
            refreshHost () {
                Bus.$emit('refresh-count', {
                    type: 'host_count',
                    hosts: [...this.table.selection]
                })
                this.table.selection = []
                this.$success('转移成功')
                RouterQuery.set({
                    _t: Date.now(),
                    page: 1
                })
            },
            doLayoutTable () {
                this.$refs.table.doLayout()
            }
        }
    }
</script>

<style lang="scss" scoped>
    .list-layout {
        overflow: hidden;
    }
    .filter-tag ~ .host-table {
        margin-top: 0;
    }
    .host-table {
        margin-top: 10px;
    }
</style>
