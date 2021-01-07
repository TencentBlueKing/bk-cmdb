<template>
    <div class="topology-layout clearfix">
        <cmdb-resize-layout class="tree-layout fl"
            v-bkloading="{isLoading: $loading([`get_getInstTopo_${business}`, `get_getInternalTopo_${business}`])}"
            direction="right"
            :handler-offset="3"
            :min="200"
            :max="480">
            <cmdb-business-selector class="business-selector" v-model="business">
            </cmdb-business-selector>
            <div class="tree-simplify" v-if="false && tree.simplifyAvailable">
                <cmdb-form-bool class="tree-simplify-checkbox"
                    :size="16"
                    :true-value="true"
                    :false-value="false"
                    v-model="tree.simplify">
                    <span>
                        {{$t('BusinessTopology["精简显示"]')}}
                    </span>
                </cmdb-form-bool>
                <v-popover style="display: inline-block;" trigger="hover" :delay="200">
                    <i class="tree-simplify-tips bk-icon icon-info-circle">
                    </i>
                    <img src="@/assets/images/simplify-tips.png" slot="popover">
                </v-popover>
            </div>
            <cmdb-tree ref="topoTree" class="topo-tree"
                children-key="child"
                :id-generator="getTopoNodeId"
                :tree="tree.data"
                @on-selected="handleNodeSelected">
                <div class="tree-node clearfix" slot-scope="{node, state}"
                    :class="{
                        'tree-node-selected': state.selected,
                        'tree-node-module': node['bk_obj_id'] === 'module'
                    }">
                    <template v-if="[1, 2].includes(node.default)">
                        <i class='topo-node-icon topo-node-icon-internal icon-cc-host-free-pool' v-if="node.default === 1"></i>
                        <i class='topo-node-icon topo-node-icon-internal icon-cc-host-breakdown' v-else></i>
                    </template>
                    <i class="topo-node-icon topo-node-icon-text" v-else>{{node['bk_obj_name'][0]}}</i>
                    <span class="topo-node-text" :title="node['bk_inst_name']">{{node['bk_inst_name']}}</span>
                    <bk-button type="primary" class="topo-node-btn-create fr"
                        v-if="showCreate(node, state)"
                        @click.stop="handleCreate">
                        {{$t('Common[\'新增\']')}}
                    </bk-button>
                </div>
            </cmdb-tree>
            <bk-dialog
                :is-show.sync="tree.create.showDialog"
                :has-header="false"
                :has-footer="false"
                :padding="0"
                :quick-close="false"
                @after-transition-leave="handleAfterCancelCreateNode"
                @cancel="handleCancelCreateNode">
                <tree-node-create v-if="tree.create.active" slot="content"
                    :properties="tree.create.properties"
                    :state="tree.selectedNodeState"
                    @on-submit="handleCreateNode"
                    @on-cancel="handleCancelCreateNode">
                </tree-node-create>
            </bk-dialog>
        </cmdb-resize-layout>
        <div class="hosts-layout">
            <bk-tab :active-name.sync="tab.active" @tab-changed="handleTabChanged">
                <bk-tabpanel class="topo-tabpanel" name="hosts" :title="$t('BusinessTopology[\'主机调配\']')">
                    <bk-button class="topo-table-btn-refresh" type="primary" style="display: none;"
                        :disabled="$loading()"
                        @click="handleRefresh">
                        {{$t("HostResourcePool['刷新查询']")}}
                    </bk-button>
                    <cmdb-hosts-table class="topo-table" ref="topoTable"
                        :columns-config-key="table.columnsConfigKey"
                        :columns-config-properties="columnsConfigProperties"
                        :quick-search="true"
                        @on-quick-search="handleQuickSearch">
                    </cmdb-hosts-table>
                </bk-tabpanel>
                <bk-tabpanel name="attribute" :title="$t('BusinessTopology[\'节点属性\']')"
                    v-bkloading="{isLoading: $loading()}"
                    :show="showAttributePanel">
                    <cmdb-details class="topology-details"
                        v-if="isNodeDetailsActive"
                        :showDelete="false"
                        :properties="tab.properties"
                        :property-groups="tab.propertyGroups"
                        :inst="tree.flatternedSelectedNodeInst"
                        @on-edit="handleEdit">
                    </cmdb-details>
                    <cmdb-form class="topology-details" v-else-if="['update', 'create'].includes(tab.type)"
                        :properties="tab.properties"
                        :property-groups="tab.propertyGroups"
                        :inst="tree.selectedNodeInst"
                        :type="tab.type"
                        @on-submit="handleSubmit"
                        @on-cancel="handleCancel">
                        <template slot="extra-options">
                            <bk-button type="danger" style="margin-left: 4px" @click="handleDelete">{{$t('Common["删除"]')}}
                            </bk-button>
                        </template>
                    </cmdb-form>
                </bk-tabpanel>
                <bk-tabpanel name="process" :title="$t('ProcessManagement[\'进程信息\']')"
                    :show="showProcessPanel">
                    <cmdb-topo-node-process
                        v-if="tab.active === 'process'"
                        :business="business"
                        :module="tree.selectedNode">
                    </cmdb-topo-node-process>
                </bk-tabpanel>
            </bk-tab>
        </div>
    </div>
</template>

<script>
    import { mapGetters, mapActions } from 'vuex'
    import cmdbHostsTable from '@/components/hosts/table'
    import cmdbTopoNodeProcess from './children/_node-process'
    import treeNodeCreate from './children/_node-create.vue'
    export default {
        components: {
            cmdbHostsTable,
            cmdbTopoNodeProcess,
            treeNodeCreate
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
                    simplify: false,
                    simplifyAvailable: false,
                    simplifyParentNode: null,
                    selectedNode: null,
                    selectedNodeState: null,
                    selectedNodeInst: {},
                    flatternedSelectedNodeInst: {},
                    internalModule: [],
                    create: {
                        showDialog: false,
                        active: false,
                        properties: []
                    }
                },
                tab: {
                    active: 'hosts',
                    type: 'details',
                    properties: [],
                    propertyGroups: []
                },
                table: {
                    params: null,
                    columnsConfigKey: 'topology_table_columns',
                    quickSearch: {
                        property: null,
                        value: ''
                    }
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
            },
            'tree.simplify' (simplify) {
                if (simplify) {
                    this.simplifyTree()
                } else {
                    this.getBusinessTopo()
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
            ...mapActions('objectSet', [
                'searchSet',
                'createSet',
                'updateSet',
                'deleteSet'
            ]),
            ...mapActions('objectModule', [
                'searchModule',
                'createModule',
                'updateModule',
                'deleteModule'
            ]),
            ...mapActions('objectCommonInst', [
                'searchInst',
                'createInst',
                'updateInst',
                'deleteInst'
            ]),
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
            getCommonProperties (objId) {
                if (this.properties.hasOwnProperty(objId)) {
                    this.tab.properties = this.properties[objId]
                    return Promise.resolve(this.properties[objId])
                }
                return this.searchObjectAttribute({
                    params: {
                        'bk_obj_id': objId,
                        'bk_supplier_account': this.supplierAccount
                    },
                    config: {
                        requestId: `post_searchObjectAttribute_${objId}`,
                        fromCache: true
                    }
                }).then(properties => {
                    this.$set(this.properties, objId, properties)
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
                        requestId: `post_searchGroup_${objId}`
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
                    this.tree.flatternedSelectedNodeInst = this.$tools.flatternItem(this.tab.properties, data.info[0])
                })
            },
            getMainlineModel () {
                return this.searchMainlineObject({
                    requestId: 'get_searchMainlineObject',
                    fromCache: true
                }).then(topoModel => {
                    this.topoModel = topoModel
                    return topoModel
                })
            },
            getBusinessTopo () {
                return Promise.all([
                    this.getInstTopo({
                        bizId: this.business,
                        config: {
                            requestId: `get_getInstTopo_${this.business}`,
                            cancelPrevious: true
                        }
                    }),
                    this.getInternalTopo({
                        bizId: this.business,
                        config: {
                            requestId: `get_getInternalTopo_${this.business}`,
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
                    this.tree.internalModule = internalModule
                    this.tree.data = [{
                        selected: true,
                        expanded: true,
                        ...instTopo[0],
                        child: [...internalModule, ...instTopo[0].child]
                    }]
                    this.setSimplifyAvailable()
                })
            },
            setSimplifyAvailable () {
                const simplifyData = this.tree.data.filter(data => !this.tree.internalModule.includes(data))
                this.tree.simplifyAvailable = simplifyData.length === 1 && simplifyData[0].child.length
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
                    this.$refs.topoTable.table.checked = []
                    this.setSearchParams()
                    this.handleRefresh()
                }
            },
            handleQuickSearch (property, value) {
                this.table.quickSearch.property = property
                this.table.quickSearch.value = value
                this.setSearchParams()
                this.handleRefresh()
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
                const quickSearch = this.table.quickSearch
                if (quickSearch.property && quickSearch.value !== null) {
                    if (['singleasst', 'multiasst'].includes(quickSearch.property['bk_property_type'])) {
                        condition.push({
                            'bk_obj_id': quickSearch.property['bk_asst_obj_id'],
                            condition: [{
                                field: 'bk_inst_name',
                                operator: '$regex',
                                value: quickSearch.value
                            }]
                        })
                    } else {
                        const hostCondition = condition.find(condition => condition['bk_obj_id'] === 'host')
                        hostCondition.condition.push({
                            field: quickSearch.property['bk_property_id'],
                            operator: '$regex',
                            value: quickSearch.value
                        })
                    }
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
                    if (active === 'hosts') {
                        this.handleRefresh()
                    }
                }
            },
            handleEdit () {
                this.tab.type = 'update'
            },
            async handleCreate () {
                this.tree.create.showDialog = true
                this.tree.create.active = true
                let targetNode = this.tree.selectedNode
                if (this.tree.simplify && targetNode['bk_obj_id'] === 'biz') {
                    targetNode = this.tree.simplifyParentNode
                }
                const model = this.topoModel.find(model => model['bk_obj_id'] === targetNode['bk_obj_id'])
                const properties = await this.getCommonProperties(model['bk_next_obj'])
                this.tree.create.properties = properties
            },
            handleCreateNode (values) {
                this.createNode(values).then(() => {
                    this.handleCancelCreateNode()
                })
            },
            handleCancelCreateNode () {
                this.tree.create.showDialog = false
            },
            handleAfterCancelCreateNode () {
                this.tree.create.active = false
                this.tree.create.properties = []
            },
            handleSubmit (value, changedValue, originalInst, type) {
                let promise = type === 'create' ? this.createNode(value) : this.updateNode(value)
                promise.then(() => {
                    this.$http.cancelCache('getInstTopo')
                })
            },
            createNode (value) {
                let selectedNode = this.tree.selectedNode
                if (this.tree.simplify && selectedNode['bk_obj_id'] === 'biz') {
                    selectedNode = this.tree.simplifyParentNode
                }
                const selectedNodeModel = this.topoModel.find(model => model['bk_obj_id'] === selectedNode['bk_obj_id'])
                const nextObjId = selectedNodeModel['bk_next_obj']
                const formData = {
                    ...value,
                    'bk_biz_id': this.business,
                    'bk_parent_id': selectedNode['bk_inst_id']
                }
                let promise
                let instIdKey
                let instNameKey
                if (nextObjId === 'set') {
                    formData['bk_supplier_account'] = this.supplierAccount
                    instIdKey = 'bk_set_id'
                    instNameKey = 'bk_set_name'
                    promise = this.createSet({
                        bizId: this.business,
                        params: formData
                    })
                } else if (nextObjId === 'module') {
                    formData['bk_supplier_account'] = this.supplierAccount
                    instIdKey = 'bk_module_id'
                    instNameKey = 'bk_module_name'
                    promise = this.createModule({
                        bizId: this.business,
                        setId: selectedNode['bk_inst_id'],
                        params: formData
                    })
                } else {
                    instIdKey = 'bk_inst_id'
                    instNameKey = 'bk_inst_name'
                    promise = this.createInst({
                        objId: nextObjId,
                        params: formData
                    })
                }
                promise.then(inst => {
                    const children = selectedNode.child || []
                    const newNode = {
                        default: 0,
                        child: [],
                        'bk_inst_id': inst[instIdKey],
                        'bk_inst_name': inst[instNameKey],
                        'bk_obj_id': nextObjId,
                        'bk_obj_name': selectedNodeModel['bk_next_name']
                    }
                    if (selectedNode['bk_obj_id'] === 'biz') {
                        children.splice(2, 0, newNode)
                    } else {
                        children.unshift(newNode)
                    }
                    this.$set(selectedNode, 'child', children)
                    this.$set(selectedNode, 'expanded', true)
                    this.tab.active = 'hosts'
                    this.tab.type = 'details'
                    this.$success(this.$t('Common[\'新建成功\']'))
                })
                return promise
            },
            updateNode (value) {
                const formData = {...value}
                const selectedNode = this.tree.selectedNode
                const objId = selectedNode['bk_obj_id']
                let promise
                if (objId === 'set') {
                    formData['bk_supplier_account'] = this.supplierAccount
                    promise = this.updateSet({
                        bizId: this.business,
                        setId: selectedNode['bk_inst_id'],
                        params: formData
                    })
                } else if (objId === 'module') {
                    formData['bk_supplier_account'] = this.supplierAccount
                    promise = this.updateModule({
                        bizId: this.business,
                        setId: value['bk_set_id'],
                        moduleId: selectedNode['bk_inst_id'],
                        params: formData
                    })
                } else {
                    promise = this.updateInst({
                        objId: objId,
                        instId: selectedNode['bk_inst_id'],
                        params: formData
                    })
                }
                promise.then(() => {
                    const instNameKey = ['set', 'module'].includes(objId) ? `bk_${objId}_name` : 'bk_inst_name'
                    selectedNode['bk_inst_name'] = formData[instNameKey]
                    this.tree.selectedNodeInst = value
                    this.tree.flatternedSelectedNodeInst = this.$tools.flatternList(this.tab.properties, value)
                    this.tab.type = 'details'
                    this.$success(this.$t('Common[\'修改成功\']'))
                })
                return promise
            },
            handleCancel () {
                if (this.tab.type === 'update') {
                    this.tab.type = 'details'
                } else {
                    this.tab.active = 'hosts'
                    this.tab.type = 'update'
                }
            },
            handleDelete () {
                const selectedNode = this.tree.selectedNode
                const parentNode = this.tree.selectedNodeState.parent.node
                const objId = selectedNode['bk_obj_id']
                const config = {requestId: 'deleteNode'}
                this.$bkInfo({
                    title: `${this.$t('Common["确定删除"]')} ${selectedNode['bk_inst_name']}?`,
                    content: objId === 'module'
                        ? this.$t('Common["请先转移其下所有的主机"]')
                        : this.$t('Common[\'下属层级都会被删除，请先转移其下所有的主机\']'),
                    confirmFn: () => {
                        let promise
                        if (objId === 'set') {
                            promise = this.deleteSet({
                                bizId: this.business,
                                setId: selectedNode['bk_inst_id'],
                                config
                            })
                        } else if (objId === 'module') {
                            promise = this.deleteModule({
                                bizId: this.business,
                                setId: parentNode['bk_inst_id'],
                                moduleId: selectedNode['bk_inst_id'],
                                config
                            })
                        } else {
                            promise = this.deleteInst({
                                objId,
                                instId: selectedNode['bk_inst_id'],
                                config
                            })
                        }
                        promise.then(() => {
                            parentNode.child = parentNode.child.filter(node => node !== selectedNode)
                            this.tab.active = 'hosts'
                            this.$refs.topoTree.selectNode(this.getTopoNodeId(this.tree.data[0]))
                            this.$success(this.$t('Common[\'删除成功\']'))
                        })
                    }
                })
            },
            showCreate (node, state) {
                const selected = state.selected
                const isBlueKing = this.tree.data[0]['bk_inst_name'] === '蓝鲸'
                const isModule = node['bk_obj_id'] === 'module'
                return selected && !isBlueKing && !isModule
            },
            simplifyTree () {
                this.$refs.topoTree.selectNode(this.getTopoNodeId(this.tree.data[0]))
                let instTopo = this.tree.data[0].child.slice(this.tree.internalModule.length)
                let simplifyParentNode
                while (instTopo.length === 1 && instTopo[0].child.length) {
                    simplifyParentNode = instTopo[0]
                    instTopo = instTopo[0].child
                }
                this.tree.simplifyParentNode = simplifyParentNode
                this.tree.data[0].child.splice(this.tree.internalModule.length, 1, ...instTopo)
                this.$refs.topoTree.$forceUpdate()
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
            margin: 20px 20px 13px;
        }
        .tree-simplify {
            padding: 0 20px;
            font-size: 14px;
            .tree-simplify-tips:hover {
                color: #0082ff;
            }
        }
        .topo-tree{
            padding: 0 0 0 20px;
            margin: 10px 0 0 0;
            height: calc(100% - 100px);
            @include scrollbar-y;
            .tree-node {
                font-size: 0;
                &:hover{
                    .topo-node-icon.topo-node-icon-text{
                        background-color: #50abff;
                    }
                    .topo-node-icon.topo-node-icon-internal{
                        color: #50abff;
                    }
                }
                &.tree-node-selected {
                    .topo-node-icon.topo-node-icon-text {
                        background-color: #498fe0;
                    }
                    .topo-node-icon.topo-node-icon-internal {
                        color: #ffb400;
                    }
                }
                &.tree-node-selected:not(.tree-node-module) {
                    .topo-node-text {
                        max-width: calc(100% - 65px);
                    }
                }
            }
            .topo-node-icon{
                display: inline-block;
                vertical-align: middle;
                width: 18px;
                height: 18px;
                line-height: 16px;
                font-size: 12px;
                text-align: center;
                color: #fff;
                font-style: normal;
                background-color: #c3cdd7;
                border-radius: 50%;
                &.topo-node-icon-internal{
                    font-size: 16px;
                    color: $cmdbTextColor;
                    background-color: transparent;
                }
            }
            .topo-node-text{
                display: inline-block;
                vertical-align: middle;
                max-width: calc(100% - 18px);
                padding: 0 0 0 8px;
                font-size: 14px;
                @include ellipsis;
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
    .topology-details {
        max-width: 700px;
    }
</style>