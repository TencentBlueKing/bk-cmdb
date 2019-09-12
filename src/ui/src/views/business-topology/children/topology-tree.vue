<template>
    <div class="topology-tree-wrapper">
        <bk-input class="filter-input"
            right-icon="bk-icon icon-search"
            :placeholder="$t('请输入关键字')"
            v-model="nodeKeyworyds"
            @enter="handleFilterNode">
        </bk-input>
        <bk-big-tree class="topology-tree"
            ref="tree"
            v-bkloading="{
                isLoading: $loading(['getTopologyData', 'getMainLine'])
            }"
            :options="{
                idKey: idGenerator,
                nameKey: 'bk_inst_name',
                childrenKey: 'child'
            }"
            :expand-on-click="false"
            selectable
            expand-icon="bk-icon icon-down-shape"
            collapse-icon="bk-icon icon-right-shape"
            @select-change="handleSelectChange">
            <div class="node-info clearfix" slot-scope="{ node, data }">
                <i :class="['node-model-icon fl', { 'is-selected': node.selected }, { 'is-template': isTemplate(node) }]">{{modelIconMap[data.bk_obj_id]}}</i>
                <!-- <span v-if="showCreate(node, data)"
                    class="fr"
                    style="display: inline-block; font-size: 0;"
                    v-cursor="{
                        active: !$isAuthorized($OPERATION.C_TOPO),
                        auth: [$OPERATION.C_TOPO]
                    }">
                    <bk-button class="node-button"
                        theme="primary"
                        :disabled="!$isAuthorized($OPERATION.C_TOPO)"
                        @click.stop="showCreateDialog(node)">
                        {{$t('新建')}}
                    </bk-button>
                </span> -->

                <cmdb-dot-menu class="operation-menu fr" ref="dotMenu" color="#3a84ff" @click.native.stop>
                    <ul class="menu-list">
                        <li class="menu-item"
                            v-for="(menu, index) in operationMenu(node, data)"
                            :key="index">
                            <span class="menu-span"
                                v-cursor="{
                                    active: !$isAuthorized($OPERATION[menu.auth]),
                                    auth: [$OPERATION[menu.auth]]
                                }">
                                <bk-button class="menu-button"
                                    :text="true"
                                    :disabled="!$isAuthorized($OPERATION[menu.auth])"
                                    @click="menu.handler(node, data)">
                                    {{menu.name}}
                                </bk-button>
                            </span>
                        </li>
                    </ul>
                </cmdb-dot-menu>

                <div class="info-content">
                    <span class="node-name">{{data.bk_inst_name}}</span>
                    <span class="instance-num">{{data.service_instance_count}}</span>
                </div>
            </div>
        </bk-big-tree>
        <bk-dialog class="bk-dialog-no-padding"
            v-model="createInfo.show"
            :show-footer="false"
            :mask-close="false"
            :width="580"
            @after-leave="handleAfterCancelCreateNode"
            @cancel="handleCancelCreateNode">
            <template v-if="createInfo.nextModelId === 'module'">
                <create-module v-if="createInfo.visible"
                    :parent-node="createInfo.parentNode"
                    @submit="handleCreateNode"
                    @cancel="handleCancelCreateNode">
                </create-module>
            </template>
            <!-- <template v-else>
                <create-node v-if="createInfo.visible"
                    :next-model-id="createInfo.nextModelId"
                    :properties="createInfo.properties"
                    :parent-node="createInfo.parentNode"
                    @submit="handleCreateNode"
                    @cancel="handleCancelCreateNode">
                </create-node>
            </template> -->
            <template v-else-if="createInfo.nextModelId === 'set'">
                <create-cluster v-if="createInfo.visible"
                    :next-model-id="createInfo.nextModelId"
                    :properties="createInfo.properties"
                    :parent-node="createInfo.parentNode"
                    @submit="handleCreateNode"
                    @cancel="handleCancelCreateNode">
                </create-cluster>
            </template>
            <template v-else>
                <create-new-node v-if="createInfo.visible"
                    :next-model-id="createInfo.nextModelId"
                    :properties="createInfo.properties"
                    :parent-node="createInfo.parentNode"
                    @submit="handleCreateNode"
                    @cancel="handleCancelCreateNode">
                </create-new-node>
            </template>
        </bk-dialog>
    </div>
</template>

<script>
    import { mapGetters } from 'vuex'
    // import createNode from './create-node.vue'
    import createModule from './create-module.vue'
    import createCluster from './create-cluster.vue'
    import createNewNode from './create-node.new.vue'
    export default {
        components: {
            // createNode,
            createModule,
            createCluster,
            createNewNode
        },
        data () {
            return {
                treeData: [],
                mainLine: [],
                createInfo: {
                    show: false,
                    visible: false,
                    properties: [],
                    parentNode: null,
                    nextModelId: null
                },
                nodeKeyworyds: ''
            }
        },
        computed: {
            ...mapGetters(['supplierAccount', 'isAdminView']),
            business () {
                return this.$store.state.objectBiz.bizId
            },
            propertyMap () {
                return this.$store.state.businessTopology.propertyMap
            },
            mainLineModels () {
                const models = this.$store.getters['objectModelClassify/models']
                return this.mainLine.map(data => models.find(model => model.bk_obj_id === data.bk_obj_id))
            },
            modelIconMap () {
                const map = {}
                this.mainLineModels.forEach(model => {
                    map[model.bk_obj_id] = model.bk_obj_name[0]
                })
                return map
            },
            isBlueKing () {
                return (this.treeData[0] || {}).bk_inst_name === '蓝鲸'
            }
        },
        async created () {
            const [data, mainLine] = await Promise.all([
                this.getTopologyData(),
                this.getMainLine()
            ])
            this.getTopologyInstanceNum()
            this.treeData = data
            this.mainLine = mainLine
            this.$nextTick(() => {
                this.setDefaultState(data)
            })
        },
        methods: {
            setDefaultState (data) {
                this.$refs.tree.setData(data)
                const businessData = data[0]
                const businessNodeId = this.idGenerator(businessData)
                const queryModule = parseInt(this.$route.query.module)
                let defaultNodeId = businessNodeId
                if (!isNaN(queryModule)) {
                    const nodeId = `module_${queryModule}`
                    const node = this.$refs.tree.getNodeById(nodeId)
                    if (node) {
                        defaultNodeId = nodeId
                    }
                } else if (Array.isArray(businessData.child) && businessData.child.length) {
                    defaultNodeId = this.idGenerator(businessData.child[0])
                }
                this.$refs.tree.setSelected(defaultNodeId, { emitEvent: true })
                this.$refs.tree.setExpanded(defaultNodeId)
            },
            getTopologyData () {
                return this.$store.dispatch('objectMainLineModule/getInstTopo', {
                    bizId: this.business,
                    config: {
                        requestId: 'getTopologyData'
                    }
                })
            },
            getMainLine () {
                return this.$store.dispatch('objectMainLineModule/searchMainlineObject', {
                    config: {
                        requestId: 'getMainLine'
                    }
                })
            },
            getTopologyInstanceNum () {
                this.$store.dispatch('objectMainLineModule/getInstTopoInstanceNum', {
                    bizId: this.business,
                    config: {
                        requestId: 'getTopologyInstanceNum'
                    }
                }).then(data => {
                    this.setNodeNum(data)
                })
            },
            setNodeNum (data) {
                data.forEach((datum, index) => {
                    const id = this.idGenerator(datum)
                    const node = this.$refs.tree.getNodeById(id)
                    if (node) {
                        const num = datum.service_instance_count
                        datum.service_instance_count = num > 999 ? '999+' : num
                        node.data = datum
                    }
                    const child = datum.child
                    if (Array.isArray(child) && child.length) {
                        this.setNodeNum(child)
                    }
                })
            },
            idGenerator (data) {
                return `${data.bk_obj_id}_${data.bk_inst_id}`
            },
            showCreate (node, data) {
                const isModule = data.bk_obj_id === 'module'
                return node.selected && !isModule && !this.isBlueKing
            },
            operationMenu (node, data) {
                const operation = [{
                    name: this.$t('新增模块'),
                    handler: this.showCreateDialog,
                    disabled: '',
                    auth: 'C_TOPO'
                }, {
                    name: this.$t('新增子节点'),
                    handler: this.showCreateDialog,
                    disabled: '',
                    auth: 'C_TOPO'
                }, {
                    name: this.$t('新增集群'),
                    handler: this.showCreateDialog,
                    disabled: '',
                    auth: 'C_TOPO'
                }, {
                    name: this.$t('删除'),
                    handler: this.handleDelete,
                    disabled: '',
                    auth: 'D_TOPO'
                }]
                return operation
            },
            isTemplate (node) {
                return node.data.service_template_id
            },
            async showCreateDialog (node) {
                const nodeModel = this.mainLine.find(data => data.bk_obj_id === node.data.bk_obj_id)
                const nextModelId = nodeModel.bk_next_obj
                this.createInfo.nextModelId = nextModelId
                this.createInfo.parentNode = node
                this.createInfo.show = true
                this.createInfo.visible = true
                if (this.propertyMap.hasOwnProperty(nextModelId)) {
                    this.createInfo.properties = this.propertyMap[nextModelId]
                } else {
                    const action = 'objectModelProperty/searchObjectAttribute'
                    const properties = await this.$store.dispatch(action, {
                        params: this.$injectMetadata({
                            bk_obj_id: nextModelId,
                            bk_supplier_account: this.$store.getters.supplierAccount
                        }),
                        config: {
                            requestId: 'getNextModelProperties'
                        }
                    })
                    this.$store.commit('businessTopology/setProperties', {
                        id: nextModelId,
                        properties: properties
                    })
                    this.createInfo.properties = properties
                }
            },
            handleAfterCancelCreateNode () {
                this.createInfo.visible = false
                this.createInfo.properties = []
                this.createInfo.parentNode = null
                this.createInfo.nextModelId = null
            },
            handleCancelCreateNode () {
                this.createInfo.show = false
            },
            async handleCreateNode (value) {
                try {
                    const parentNode = this.createInfo.parentNode
                    const formData = this.$injectMetadata({
                        ...value,
                        'bk_biz_id': this.business,
                        'bk_parent_id': parentNode.data.bk_inst_id
                    })
                    const nextModelId = this.createInfo.nextModelId
                    const nextModel = this.mainLineModels.find(model => model.bk_obj_id === nextModelId)
                    const handlerMap = {
                        set: this.createSet,
                        module: this.createModule
                    }
                    const data = await (handlerMap[nextModelId] || this.createCommonInstance)(formData)
                    const nodeData = {
                        default: 0,
                        child: [],
                        bk_obj_name: nextModel.bk_obj_name,
                        bk_obj_id: nextModel.bk_obj_id,
                        service_instance_count: 0,
                        service_template_id: value.service_template_id,
                        ...data
                    }
                    this.$refs.tree.addNode(nodeData, parentNode.id, 0)
                    this.$success(this.$t('新建成功'))
                    this.handleCancelCreateNode()
                } catch (e) {
                    console.error(e)
                }
            },
            async createSet (value) {
                const data = await this.$store.dispatch('objectSet/createSet', {
                    bizId: this.business,
                    params: {
                        ...value,
                        bk_supplier_account: this.supplierAccount
                    }
                })
                return {
                    bk_inst_id: data.bk_set_id,
                    bk_inst_name: data.bk_set_name
                }
            },
            async createModule (value) {
                const data = await this.$store.dispatch('objectModule/createModule', {
                    bizId: this.business,
                    setId: this.createInfo.parentNode.data.bk_inst_id,
                    params: {
                        ...value,
                        bk_supplier_account: this.supplierAccount
                    }
                })
                return {
                    bk_inst_id: data.bk_module_id,
                    bk_inst_name: data.bk_module_name
                }
            },
            async createCommonInstance (value) {
                const data = await this.$store.dispatch('objectCommonInst/createInst', {
                    objId: this.createInfo.nextModelId,
                    params: value
                })
                return {
                    bk_inst_id: data.bk_inst_id,
                    bk_inst_name: data.bk_inst_name
                }
            },
            handleSelectChange (node) {
                this.$store.commit('businessTopology/setSelectedNode', node)
            },
            handleFilterNode () {
                this.$refs.tree.filter(this.nodeKeyworyds)
            },
            handleDelete (curNode, data) {
                const instId = data.bk_inst_id
                const otherId = ['set', 'module'].includes(data.bk_obj_id) ? curNode.parent.data.bk_inst_id : data.bk_obj_id
                this.$bkInfo({
                    title: `${this.$t('确定删除')} ${curNode.name}?`,
                    subTitle: data.bk_obj_id === 'module'
                        ? this.$t('删除模块提示')
                        : this.$t('下属层级都会被删除，请先转移其下所有的主机'),
                    extCls: 'bk-dialog-sub-header-center',
                    confirmFn: async () => {
                        const promiseMap = {
                            set: this.deleteSetInstance,
                            module: this.deleteModuleInstance
                        }
                        try {
                            await (promiseMap[data.bk_obj_id] || this.deleteCustomInstance)(instId, otherId)
                            const tree = curNode.tree
                            const parentId = curNode.parent.id
                            const nodeId = curNode.id
                            tree.setSelected(parentId, {
                                emitEvent: true
                            })
                            tree.removeNode(nodeId)
                            this.$success(this.$t('删除成功'))
                        } catch (e) {
                            console.error(e)
                        }
                    }
                })
            },
            deleteSetInstance (id) {
                return this.$store.dispatch('objectSet/deleteSet', {
                    bizId: this.business,
                    setId: id,
                    config: {
                        requestId: 'deleteNodeInstance',
                        data: this.$injectMetadata({})
                    }
                })
            },
            deleteModuleInstance (id, setId) {
                return this.$store.dispatch('objectModule/deleteModule', {
                    bizId: this.business,
                    setId: setId,
                    moduleId: id,
                    config: {
                        requestId: 'deleteNodeInstance',
                        data: this.$injectMetadata({})
                    }
                })
            }
        },
        deleteCustomInstance (id, modelId) {
            return this.$store.dispatch('objectCommonInst/deleteInst', {
                objId: modelId,
                instId: id,
                config: {
                    requestId: 'deleteNodeInstance',
                    data: this.$injectMetadata()
                }
            })
        }
    }
</script>

<style lang="scss" scoped>
    .topology-tree-wrapper {
        height: 100%;
        .filter-input {
            width: auto;
            margin: 10px 20px 20px;
        }
        /deep/ .bk-big-tree-node {
            .node-options {
                .bk-icon {
                    font-size: 16px;
                    margin: 0;
                    line-height: 38px;
                    color: #c4c6cc;
                }
            }
            &.is-selected .node-options .bk-icon{
                color: #3a84ff;
            }
        }
    }
    .node-info {
        .node-model-icon {
            width: 22px;
            height: 22px;
            line-height: 21px;
            text-align: center;
            font-style: normal;
            font-size: 12px;
            margin: 8px 8px 0 6px;
            border-radius: 50%;
            background-color: #c4c6cc;
            color: #fff;
            &.is-template {
                background-color: #97aed6;
            }
            &.is-selected {
                background-color: #3a84ff;
            }
        }
        .node-button {
            height: 24px;
            padding: 0 6px;
            margin: 0 10px 0 4px;
            line-height: 22px;
            border-radius: 4px;
            font-size: 12px;
            min-width: auto;
        }
        .operation-menu {
            width: 26px;
            height: 26px;
            margin: 5px 14px 0 0;
            // display: none;
            /deep/ {
                .bk-tooltip-ref {
                    height: 100%;
                    border-radius: 50%;
                    background-color: #d0e0fd;
                }
                .menu-trigger {
                    padding: 3px 0;
                }
            }
        }
        .info-content {
            display: flex;
            align-items: center;
            line-height: 36px;
            font-size: 14px;
            .node-name {
                @include ellipsis;
                margin-right: 8px;
            }
            .instance-num {
                margin-right: 5px;
                padding: 0 5px;
                height: 18px;
                line-height: 17px;
                border-radius: 2px;
                background-color: #f0f1f5;
                color: #979ba5;
                font-size: 12px;
                text-align: center;
            }
        }
    }
    .topology-tree {
        height: 100%;
        .bk-big-tree-node {
            &:hover {
                .operation-menu {
                    display: block;
                }
            }
            &.is-selected {
                .instance-num {
                    background-color: #a2c5fd;
                    color: #fff;
                }
            }
        }
    }
    .menu-list {
        padding: 6px 0;
        font-size: 14px;
        .menu-item {
            height: 32px;
            line-height: 32px;
            &:hover {
                .menu-button:not(.is-disabled) {
                    color: #3a84ff;
                    background-color: #e1ecff;
                }
            }
        }
        .menu-span {
            display: block;
        }
        .menu-button {
            width: 100%;
            height: 32px;
            line-height: 32px;
            text-align: left;
            padding: 0 14px;
            color: #63656e;
            &.is-disabled {
                color: #c4c6cc;
            }
        }
    }
</style>
