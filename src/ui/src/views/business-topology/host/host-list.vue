<template>
    <div class="list-layout">
        <host-list-options @transfer="handleTransfer"></host-list-options>
        <bk-table class="host-table"
            v-bkloading="{ isLoading: $loading(Object.values(request)) || !commonRequestFinished }"
            :data="table.data"
            :pagination="table.pagination"
            :row-style="{ cursor: 'pointer' }"
            :max-height="$APP.height - 250"
            @page-change="refresh"
            @page-limit-change="handleLimitChange"
            @sort-change="handleSortChange"
            @row-click="handleRowClick"
            @selection-change="handleSelectionChange">
            <bk-table-column type="selection" width="50"></bk-table-column>
            <bk-table-column v-for="column in table.header"
                :key="column.bk_property_id"
                :label="column.bk_property_name"
                :sortable="getColumnSortable(column)"
                :prop="column.bk_property_id"
                :width="column.bk_property_id === 'bk_host_innerip' ? 130 : 'auto'"
                :fixed="column.bk_property_id === 'bk_host_innerip'"
                :class-name="column.bk_property_id === 'bk_host_innerip' ? 'is-highlight' : ''">
                <div slot-scope="{ row }"
                    :title="row | hostValueFilter(column.bk_obj_id, column.bk_property_id) | formatter(column) | unit(column.unit)">
                    {{ row | hostValueFilter(column.bk_obj_id, column.bk_property_id) | formatter(column) | unit(column.unit) }}
                </div>
            </bk-table-column>
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
    import MoveToResourceConfirm from './move-to-resource-confirm.vue'
    import hostValueFilter from '@/filters/host'
    import debounce from 'lodash.debounce'
    import { mapGetters, mapState } from 'vuex'
    import {
        MENU_BUSINESS_HOST_DETAILS,
        MENU_BUSINESS_TRANSFER_HOST
    } from '@/dictionary/menu-symbol'
    import Bus from '@/utils/bus.js'
    export default {
        components: {
            HostListOptions,
            [ModuleSelector.name]: ModuleSelector,
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
                    header: [],
                    selection: [],
                    sort: 'bk_host_id',
                    pagination: this.$tools.getDefaultPaginationConfig()
                },
                columnsConfig: {
                    selected: [],
                    fixedColumns: ['bk_host_innerip', 'bk_cloud_id', 'bk_module_name', 'bk_set_name']
                },
                dialog: {
                    width: 720,
                    height: 460,
                    show: false,
                    component: null,
                    props: {}
                },
                request: {
                    table: Symbol('table'),
                    moveToResource: Symbol('moveToResource'),
                    moveToIdleModule: Symbol('moveToIdleModule')
                },
                refresh: null
            }
        },
        computed: {
            ...mapGetters('userCustom', ['usercustom']),
            ...mapGetters('objectBiz', ['bizId']),
            ...mapGetters('businessHost', [
                'columnsConfigProperties',
                'selectedNode',
                'getDefaultSearchCondition',
                'commonRequest'
            ]),
            ...mapState('hosts', ['filterParams']),
            customColumns () {
                const customColumnKey = this.$route.meta.customInstanceColumn
                return this.usercustom[customColumnKey] || []
            }
        },
        watch: {
            customColumns () {
                this.setTableHeader()
            },
            columnsConfigProperties () {
                this.setTableHeader()
            },
            selectedNode (node) {
                node && this.active && this.refresh(1)
            },
            filterParams () {
                this.selectedNode && this.refresh(1)
            }
        },
        created () {
            this.refresh = debounce(current => {
                this.handlePageChange(current)
            }, 10)
            Bus.$on('refresh-list', this.handlePageChange)
        },
        beforeDestroy () {
            Bus.$off('refresh-list', this.handlePageChange)
        },
        methods: {
            setTableHeader () {
                const properties = this.$tools.getHeaderProperties(
                    this.columnsConfigProperties,
                    this.customColumns,
                    this.columnsConfig.fixedColumns
                )
                this.table.header = properties
            },
            getColumnSortable (column) {
                const isHostProperty = column.bk_obj_id === 'host'
                const isForeignKey = column.bk_property_type === 'foreignkey'
                return (isHostProperty && !isForeignKey) ? 'custom' : false
            },
            handlePageChange (current) {
                this.table.pagination.current = current
                this.getHostList()
            },
            handleLimitChange (limit) {
                this.table.pagination.limit = limit
                this.refresh(1)
            },
            handleSortChange (sort) {
                this.table.sort = this.$tools.getSort(sort)
                this.refresh(1)
            },
            handleRowClick (row, event, column) {
                if (column.type === 'selection') {
                    return false
                }
                this.$router.push({
                    name: MENU_BUSINESS_HOST_DETAILS,
                    params: {
                        business: this.bizId,
                        id: row.host.bk_host_id
                    },
                    query: {
                        from: 'business',
                        node: this.selectedNode.id
                    }
                })
            },
            handleSelectionChange (selection) {
                this.table.selection = selection
            },
            async getHostList () {
                try {
                    await this.commonRequest
                    this.commonRequestFinished = true
                    const result = await this.$store.dispatch('hostSearch/searchHost', {
                        params: this.getParams(),
                        config: {
                            requestId: this.request.table
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
                    bk_biz_id: this.bizId,
                    ip: this.filterParams.ip,
                    page: {
                        ...this.$tools.getPageParams(this.table.pagination),
                        sort: this.table.sort
                    },
                    condition: this.getDefaultSearchCondition()
                }
                params.condition.forEach(modleCondition => {
                    modleCondition.condition = this.filterParams[modleCondition.bk_obj_id] || []
                })
                const idMap = {
                    host: 'bk_host_id',
                    set: 'bk_set_id',
                    module: 'bk_module_id',
                    biz: 'bk_biz_id',
                    object: 'bk_inst_id'
                }
                const nodeData = this.selectedNode.data
                const conditionObjectId = Object.keys(idMap).includes(nodeData.bk_obj_id) ? nodeData.bk_obj_id : 'object'
                const selectedNodeCondition = params.condition.find(target => target.bk_obj_id === conditionObjectId)
                selectedNodeCondition.condition.push({
                    field: idMap[conditionObjectId],
                    operator: '$eq',
                    value: nodeData.bk_inst_id
                })
                return params
            },
            handleTransfer (type) {
                if (['idle', 'business'].includes(type)) {
                    this.dialog.props = {
                        moduleType: type,
                        title: type === 'idle' ? this.$t('转移主机到空闲模块') : this.$t('转移主机到业务模块')
                    }
                    this.dialog.width = 720
                    this.dialog.height = 460
                    this.dialog.component = ModuleSelector.name
                } else {
                    this.dialog.props = {
                        count: this.table.selection.length
                    }
                    this.dialog.width = 400
                    this.dialog.height = 231
                    this.dialog.component = MoveToResourceConfirm.name
                }
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
                    this.moveHostToResource()
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
                    this.refresh(1)
                    this.table.selection = []
                    this.$success('转移成功')
                } catch (e) {
                    console.error(e)
                }
            },
            gotoTransferPage (modules) {
                this.$router.push({
                    name: MENU_BUSINESS_TRANSFER_HOST,
                    params: {
                        type: this.dialog.props.moduleType
                    },
                    query: {
                        sourceModel: this.selectedNode.data.bk_obj_id,
                        sourceId: this.selectedNode.data.bk_inst_id,
                        targetModules: modules.map(node => node.data.bk_inst_id).join(','),
                        resources: this.table.selection.map(item => item.host.bk_host_id).join(',')
                    }
                })
            },
            async moveHostToResource () {
                try {
                    await this.$store.dispatch('hostRelation/transferHostToResourceModule', {
                        params: {
                            bk_biz_id: this.bizId,
                            bk_host_id: this.table.selection.map(item => item.host.bk_host_id)
                        },
                        config: {
                            requestId: this.request.moveToResource
                        }
                    })
                    Bus.$emit('refresh-count', {
                        type: 'host_count',
                        hosts: [...this.table.selection]
                    })
                    this.table.selection = []
                    this.refresh(1)
                    this.$success('转移成功')
                } catch (e) {
                    console.error(e)
                }
            }
        }
    }
</script>

<style lang="scss" scoped>
    .list-layout {
        overflow: hidden;
    }
    .host-table {
        margin-top: 12px;
    }
</style>
