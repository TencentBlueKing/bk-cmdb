<template>
    <div class="topology-layout clearfix">
        <cmdb-resize-layout class="tree-layout fl"
            v-bkloading="{isLoading: $loading(['getInstTopo', 'getInternalTopo'])}"
            direction="right"
            :min="200"
            :max="480">
            <cmdb-business-selector class="business-selector" v-model="business">
            </cmdb-business-selector>
            <cmdb-tree ref="topoTree" class="topo-tree"
                id-key="bk_inst_id"
                label-key="bk_inst_name"
                children-key="child"
                :id-generator="getTopoNodeId"
                :tree="tree.data"
                @on-selected="handleNodeSelected">
                <div class="tree-node clearfix" slot-scope="{node, state}">
                    <template v-if="[1, 2].includes(node.default)">
                        <i class='topo-node-icon topo-node-icon-internal icon-cc-host-free-pool' v-if="node.default === 1"></i>
                        <i class='topo-node-icon topo-node-icon-internal icon-cc-host-breakdown' v-else></i>
                    </template>
                    <i class="topo-node-icon topo-node-icon-text" v-else>{{node['bk_obj_name'][0]}}</i>
                    <span class="topo-node-text">{{node['bk_inst_name']}}</span>
                    <bk-button type="primary" class="topo-node-btn-create fr"
                        v-if="showCreate(node, state)"
                        @click.stop="handleCreate">
                        {{$t('Common[\'新增\']')}}
                    </bk-button>
                </div>
            </cmdb-tree>
        </cmdb-resize-layout>
        <div class="hosts-layout">
            <bk-tab :active-name.sync="tab.active" @tab-changed="handleTabChanged">
                <bk-tabpanel class="topo-tabpanel" name="hosts" :title="$t('BusinessTopology[\'主机调配\']')">
                    <bk-button class="topo-table-btn-refresh" type="primary"
                        :disabled="$loading()"
                        @click="handleRefresh">
                        {{$t("HostResourcePool['刷新查询']")}}
                    </bk-button>
                    <cmdb-hosts-table class="topo-table" ref="topoTable"
                        :columns-config-key="table.columnsConfigKey"
                        :columns-config-properties="columnsConfigProperties">
                    </cmdb-hosts-table>
                </bk-tabpanel>
                <bk-tabpanel name="attribute" :title="$t('BusinessTopology[\'节点属性\']')"
                    v-bkloading="{isLoading: $loading()}"
                    :show="showAttributePanel">
                    <cmdb-topo-node-details 
                        v-if="isNodeDetailsActive"
                        :properties="tab.properties"
                        :property-groups="tab.propertyGroups"
                        :inst="tree.flatternedSelectedNodeInst"
                        @on-edit="handleEdit">
                    </cmdb-topo-node-details>
                    <cmdb-topo-node-form v-else-if="['update', 'create'].includes(tab.type)"
                        :properties="tab.properties"
                        :property-groups="tab.propertyGroups"
                        :inst="tree.selectedNodeInst"
                        :type="tab.type"
                        @on-submit="handleSubmit"
                        @on-cancel="handleCancel"
                        @on-delete="handleDelete">
                    </cmdb-topo-node-form>
                </bk-tabpanel>
                <bk-tabpanel name="process" :title="$t('ProcessManagement[\'进程信息\']')"
                    :show="showProcessPanel">
                </bk-tabpanel>
            </bk-tab>
        </div>
    </div>
</template>

<script>
    import { mapGetters, mapActions } from 'vuex'
    import cmdbHostsTable from '@/components/hosts/table'
    import cmdbTopoNodeDetails from './children/_node-details'
    import cmdbTopoNodeForm from './children/_node-form'
    export default {
        components: {
            cmdbHostsTable,
            cmdbTopoNodeDetails,
            cmdbTopoNodeForm
        },
        data () {
            return {
                properties: {
                    biz: [],
                    host: [],
                    set: [],
                    module: []
                },
                business: '',
                businessResolver: null,
                businessTopo: [],
                topoModel: [],
                tree: {
                    data: [],
                    selectedNode: null,
                    selectedNodeState: null,
                    selectedNodeInst: {},
                    flatternedSelectedNodeInst: {},
                    internalIdleId: null,
                    internalFaultId: null
                },
                tab: {
                    active: 'hosts',
                    type: 'details',
                    properties: [],
                    propertyGroups: []
                },
                table: {
                    params: null,
                    columnsConfigKey: 'topology_table_columns'
                }
            }
        },
        computed: {
            ...mapGetters(['supplierAccount']),
            columnsConfigProperties () {
                const setProperties = this.properties.set.filter(property => ['bk_set_name'].includes(property['bk_property_id']))
                const moduleProperties = this.properties.module.filter(property => ['bk_module_name'].includes(property['bk_property_id']))
                const hostProperties = this.properties.host
                return [...setProperties, ...moduleProperties, ...hostProperties]
            },
            showAttributePanel () {
                const selectedNode = this.tree.selectedNode
                if (selectedNode) {
                    const isBusinessNode = selectedNode['bk_obj_id'] === 'biz'
                    const isDefault = !!selectedNode.default
                    const isCreate = this.tab.type === 'create'
                    const isTreeLoaded = !!this.tree.data.length
                    return isTreeLoaded && ((!isBusinessNode && !isDefault) || isCreate)
                }
                return false
            },
            showProcessPanel () {
                const selectedNode = this.tree.selectedNode
                if (selectedNode) {
                    const isDefault = !!selectedNode.default
                    const isModule = selectedNode['bk_obj_id'] === 'module'
                    return !isDefault && isModule
                }
                return false
            },
            isNodeDetailsActive () {
                const isAttributePanel = this.tab.active === 'attribute'
                const isDetails = this.tab.type === 'details'
                return isAttributePanel && isDetails
            }
        },
        watch: {
            business (business) {
                if (this.businessResolver) {
                    this.businessResolver()
                } else {
                    this.getBusinessTopo()
                }
            },
            showAttributePanel (val) {
                if (!val) {
                    this.tab.active = 'hosts'
                }
            },
            showProcessPanel (val) {
                if (!val) {
                    this.tab.active = 'hosts'
                }
            }
        },
        async created () {
            try {
                await Promise.all([
                    this.getBusiness(),
                    this.getProperties()
                ])
                await this.getMainlineModel()
                await this.getBusinessTopo()
            } catch (e) {
                console.log(e)
            }
        },
        methods: {
            ...mapActions('objectModelProperty', [
                'searchObjectAttribute',
                'batchSearchObjectAttribute'
            ]),
            ...mapActions('objectModelFieldGroup', ['searchGroup']),
            ...mapActions('objectSet', ['searchSet']),
            ...mapActions('objectModule', ['searchModule']),
            ...mapActions('objectCommonInst', ['searchInst']),
            ...mapActions('objectMainLineModule', [
                'searchMainlineObject',
                'getInstTopo',
                'getInternalTopo'
            ]),
            getBusiness () {
                return new Promise((resolve, reject) => {
                    this.businessResolver = () => {
                        this.businessResolver = null
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
                        requestId: 'hostsAttribute',
                        fromCache: true
                    }
                }).then(result => {
                    Object.keys(this.properties).forEach(objId => {
                        this.properties[objId] = result[objId]
                    })
                    return result
                })
            },
            getCommonProperties (objId) {
                if (this.properties.hasOwnProperty(objId)) {
                    this.tab.properties = this.properties[objId]
                    return Promise.resolve(this.properties[objId])
                }
                this.properties[objId] = []
                return this.searchObjectAttribute({
                    params: {
                        'bk_obj_id': objId,
                        'bk_supplier_account': this.supplierAccount
                    },
                    config: {
                        requestId: `${objId}Attribute`,
                        fromCache: true
                    }
                }).then(properties => {
                    this.properties[objId] = properties
                    this.tab.properties = properties
                    return properties
                })
            },
            getPropertyGroups (objId) {
                this.tab.propertyGroups = []
                this.searchGroup({
                    objId,
                    config: {
                        fromCache: true,
                        requestId: `${objId}AttributeGroup`
                    }
                }).then(groups => {
                    this.tab.propertyGroups = groups
                    return groups
                })
            },
            getNodeObjPropertyInfo (objId) {
                return Promise.all([
                    this.getPropertyGroups(objId),
                    this.getCommonProperties(objId)
                ])
            },
            getNodeInst () {
                const selectedNode = this.tree.selectedNode
                const objId = selectedNode['bk_obj_id']
                const instId = selectedNode['bk_inst_id']
                const requestParams = {
                    page: {start: 0, limit: 1},
                    fields: [],
                    condition: {}
                }
                const requestConfig = {
                    cancelPrevious: true
                }
                let promise
                if (objId === 'set') {
                    requestParams.condition['bk_set_id'] = instId
                    promise = this.searchSet({
                        bizId: this.business,
                        params: requestParams,
                        config: requestConfig
                    })
                } else if (objId === 'module') {
                    requestParams.condition['bk_module_id'] = instId
                    requestParams.condition['bk_supplier_account'] = this.supplierAccount
                    promise = this.searchModule({
                        bizId: this.business,
                        setId: this.tree.selectedNodeState.parent.node['bk_inst_id'],
                        params: requestParams,
                        config: requestConfig
                    })
                } else {
                    requestParams.fields = {}
                    requestParams.condition[objId] = [{
                        field: 'bk_inst_id',
                        operator: '$eq',
                        value: instId
                    }]
                    promise = this.searchInst({
                        objId,
                        params: requestParams,
                        config: requestConfig
                    })
                }
                promise.then(data => {
                    this.tree.selectedNodeInst = data.info[0]
                    this.tree.flatternedSelectedNodeInst = this.$tools.flatternList(this.tab.properties, data.info)[0]
                })
            },
            getMainlineModel () {
                return this.searchMainlineObject({fromCache: true}).then(topoModel => {
                    this.topoModel = topoModel
                    return topoModel
                })
            },
            getBusinessTopo () {
                return Promise.all([
                    this.getInstTopo({
                        bizId: this.business,
                        config: {
                            requestId: 'getInstTopo',
                            cancelPrevious: true
                        }
                    }),
                    this.getInternalTopo({
                        bizId: this.business,
                        config: {
                            requestId: 'getInternalTopo',
                            cancelPrevious: true
                        }
                    })
                ]).then(([instTopo, internalTopo]) => {
                    const moduleModel = this.getModelByObjId('module')
                    const internalModule = internalTopo.module.map(module => {
                        return {
                            'default': ['空闲机', 'idle machine'].includes(module['bk_module_name']) ? 1 : 2,
                            'bk_obj_id': 'module',
                            'bk_obj_name': moduleModel['bk_obj_name'],
                            'bk_inst_id': module['bk_module_id'],
                            'bk_inst_name': module['bk_module_name']
                        }
                    })
                    internalModule.forEach(node => {
                        if (node.default === 1) {
                            this.tree.internalIdleId = node['bk_inst_id']
                        } else {
                            this.tree.internalFaultId = node['bk_inst_id']
                        }
                    })
                    instTopo[0] = {
                        selected: true,
                        expanded: true,
                        ...instTopo[0],
                        child: [...internalModule, ...instTopo[0].child]
                    }
                    this.tree.data = instTopo
                })
            },
            getModelByObjId (id) {
                return this.topoModel.find(model => model['bk_obj_id'] === id)
            },
            getTopoNodeId (node) {
                return `${node['bk_obj_id']}-${node['bk_inst_id']}`
            },
            handleNodeSelected (node, state) {
                this.tree.selectedNode = node
                this.tree.selectedNodeState = state
                const activeTab = this.tab.active
                if (activeTab === 'attribute') {
                    this.tab.type = 'details'
                    this.handleTabChanged(activeTab)
                } else if (activeTab === 'hosts') {
                    this.setSearchParams()
                    this.handleRefresh()
                }
            },
            handleRefresh () {
                this.$refs.topoTable.search(this.business, this.table.params)
            },
            setSearchParams () {
                const necessaryObj = Object.keys(this.properties)
                const condition = necessaryObj.map(objId => {
                    return {
                        'bk_obj_id': objId,
                        condition: [],
                        fields: []
                    }
                })
                const params = {
                    ip: {
                        data: [],
                        exact: 0,
                        flag: 'bk_host_innerip|bk_host_outerip'
                    },
                    condition
                }
                const selectedNodeObjId = this.tree.selectedNode['bk_obj_id']
                if (['module', 'set'].includes(selectedNodeObjId)) {
                    const objCondition = condition.find(meta => meta['bk_obj_id'] === selectedNodeObjId)
                    objCondition.condition.push({
                        field: `bk_${selectedNodeObjId}_id`,
                        operator: '$eq',
                        value: this.tree.selectedNode['bk_inst_id']
                    })
                } else if (!necessaryObj.includes(selectedNodeObjId)) {
                    condition.push({
                        'bk_obj_id': 'object',
                        condition: [{
                            field: 'bk_inst_id',
                            operator: '$eq',
                            value: this.tree.selectedNode['bk_inst_id']
                        }],
                        fields: []
                    })
                }
                this.table.params = params
            },
            async handleTabChanged (active) {
                const selectedNode = this.tree.selectedNode
                if (this.showAttributePanel && active === 'attribute') {
                    if (this.tab.type === 'details') {
                        await this.getNodeObjPropertyInfo(selectedNode['bk_obj_id'])
                        await this.getNodeInst()
                    } else {
                        const model = this.topoModel.find(model => model['bk_obj_id'] === selectedNode['bk_obj_id'])
                        await this.getNodeObjPropertyInfo(model['bk_next_obj'])
                        this.tree.selectedNodeInst = {}
                        this.tree.flatternedSelectedNodeInst = {}
                    }
                } else {
                    this.tab.type = 'details'
                    this.tree.selectedNodeInst = {}
                    this.tree.flatternedSelectedNodeInst = {}
                }
            },
            handleEdit () {
                this.tab.type = 'update'
            },
            handleCreate () {
                this.tab.type = 'create'
                if (this.tab.active === 'attribute') {
                    this.handleTabChanged(this.tab.active)
                } else {
                    this.tab.active = 'attribute'
                }
            },
            handleSubmit (value, changedValue, originalInst, type) {},
            handleCancel () {
                if (this.tab.type === 'update') {
                    this.tab.type = 'details'
                } else {
                    this.tab.active = 'hosts'
                    this.tab.type = 'update'
                }
            },
            handleDelete () {},
            showCreate (node, state) {
                const selected = state.selected
                const isBlueKing = this.tree.data[0]['bk_inst_name'] === '蓝鲸'
                const isModule = node['bk_obj_id'] === 'module'
                return selected && !isBlueKing && !isModule
            }
        }
    }
</script>

<style lang="scss" scoped>
    .topology-layout{
        padding: 0;
        height: 100%;
    }
    .tree-layout{
        width: 280px;
        height: 100%;
        border-right: 1px solid $cmdbBorderColor;
        background-color: #fafbfd;
        .business-selector{
            display: block;
            width: auto;
            margin: 20px;
        }
        .topo-tree{
            padding: 0 0 0 20px;
            height: calc(100% - 76px);
            @include scrollbar-y;
            .tree-node {
                font-size: 0;
            }
            .topo-node-icon{
                display: inline-block;
                vertical-align: middle;
                width: 16px;
                height: 16px;
                line-height: 16px;
                font-size: 12px;
                text-align: center;
                color: #fff;
                font-style: normal;
                background-color: #c3cdd7;
                &.topo-node-icon-internal{
                    font-size: 16px;
                    color: $cmdbTextColor;
                    background-color: transparent;
                }
            }
            .topo-node-text{
                display: inline-block;
                vertical-align: middle;
                padding: 0 0 0 8px;
                font-size: 14px;
            }
            .topo-node-btn-create{
                width: auto;
                height: 18px;
                padding: 0 6px;
                margin: 3px 4px;
                line-height: 16px;
                border-radius: 4px;
                font-size: 12px;
            }
        }
    }
    .hosts-layout{
        overflow: hidden;
        padding: 0 20px;
        height: 100%;
        .topo-tabpanel{
            padding: 20px 0 0 0;
            position: relative;
            .topo-table-btn-refresh{
                position: absolute;
                top: 20px;
                right: 0;
            }
        }
        .options{
            margin: 30px 0 0 150px;
        }
    }
</style>