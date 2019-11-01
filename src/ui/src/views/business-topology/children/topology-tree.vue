<template>
    <div class="topology-tree-wrapper">
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
            <div class="node-info clearfix" :class="{ 'is-selected': node.selected }" slot-scope="{ node, data }">
                <i class="node-model-icon fl"
                    :class="{
                        'is-selected': node.selected,
                        'is-template': isTemplate(node),
                        'is-leaf-icon': node.isLeaf
                    }">
                    {{modelIconMap[data.bk_obj_id]}}
                </i>
                <cmdb-auth v-if="showCreate(node, data)"
                    class="info-create-trigger fr"
                    :auth="$authResources({ type: $OPERATION.C_TOPO })">
                    <template slot-scope="{ disabled }">
                        <i v-if="data.set_template_id"
                            class="node-button set-template-button"
                            v-bk-tooltips="{ content: $t('模板集群添加模块提示'), placement: 'top' }">
                            {{$t('新建')}}
                        </i>
                        <bk-button v-else class="node-button"
                            theme="primary"
                            :disabled="disabled"
                            @click.stop="showCreateDialog(node)">
                            {{$t('新建')}}
                        </bk-button>
                    </template>
                </cmdb-auth>
                <span class="instance-num fr">{{data.service_instance_count}}</span>
                <div class="info-content">
                    <span class="node-name">{{data.bk_inst_name}}</span>
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
            <template v-else-if="createInfo.nextModelId === 'set'">
                <create-set v-if="createInfo.visible"
                    :parent-node="createInfo.parentNode"
                    @submit="handleCreateSetNode"
                    @cancel="handleCancelCreateNode">
                </create-set>
            </template>
            <template v-else>
                <create-node v-if="createInfo.visible"
                    :next-model-id="createInfo.nextModelId"
                    :properties="createInfo.properties"
                    :parent-node="createInfo.parentNode"
                    @submit="handleCreateNode"
                    @cancel="handleCancelCreateNode">
                </create-node>
            </template>
        </bk-dialog>
    </div>
</template>

<script>
    import { mapGetters } from 'vuex'
    import createNode from './create-node.vue'
    import createSet from './create-set.vue'
    import createModule from './create-module.vue'
    export default {
        components: {
            createNode,
            createSet,
            createModule
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
                }
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
                        datum.service_instance_count = num > 999 ? '999+' : num || 0
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
                return !isModule && !this.isBlueKing
            },
            isTemplate (node) {
                return node.data.service_template_id || node.data.set_template_id
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
                    if (!['set', 'module'].includes(nextModelId)) {
                        this.$store.commit('businessTopology/setProperties', {
                            id: nextModelId,
                            properties: properties
                        })
                    }
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
            async handleCreateSetNode (value) {
                try {
                    const parentNode = this.createInfo.parentNode
                    const nextModel = this.mainLineModels.find(model => model.bk_obj_id === 'set')
                    const formData = (value.sets || []).map(set => {
                        return this.$injectMetadata({
                            ...set,
                            'bk_biz_id': this.business,
                            'bk_parent_id': parentNode.data.bk_inst_id
                        })
                    })
                    const data = await this.createSet(formData)
                    data && data.forEach(set => {
                        if (set.data) {
                            const nodeData = {
                                default: 0,
                                child: [],
                                bk_obj_name: nextModel.bk_obj_name,
                                bk_obj_id: nextModel.bk_obj_id,
                                service_instance_count: 0,
                                service_template_id: value.service_template_id,
                                bk_inst_id: set.data.bk_set_id,
                                bk_inst_name: set.data.bk_set_name,
                                set_template_id: value.set_template_id
                            }
                            this.$refs.tree.addNode(nodeData, parentNode.id, 0)
                            if (value.set_template_id) {
                                this.addModulesInSetTemplate(nodeData, set.data.bk_set_id)
                            }
                        } else {
                            this.$error(set.error_message)
                        }
                    })
                    this.$success(this.$t('新建成功'))
                    this.handleCancelCreateNode()
                } catch (e) {
                    console.error(e)
                }
            },
            async addModulesInSetTemplate (parentNodeData, id) {
                const modules = await this.$store.dispatch('objectModule/searchModule', {
                    bizId: this.business,
                    setId: id,
                    params: this.$injectMetadata(),
                    config: {
                        requestId: 'searchModule'
                    }
                })
                const parentNodeId = this.idGenerator(parentNodeData)
                const nextModel = this.mainLineModels.find(model => model.bk_obj_id === 'module')
                modules.info && modules.info.forEach(_module => {
                    const nodeData = {
                        default: 0,
                        child: [],
                        bk_obj_name: nextModel.bk_obj_name,
                        bk_obj_id: nextModel.bk_obj_id,
                        service_instance_count: 0,
                        service_template_id: _module.service_template_id,
                        bk_inst_id: _module.bk_module_id,
                        bk_inst_name: _module.bk_module_name
                    }
                    this.$refs.tree.addNode(nodeData, parentNodeId, 0)
                })
            },
            async createSet (value) {
                const data = await this.$store.dispatch('objectSet/createSetBatch', {
                    bizId: this.business,
                    params: {
                        sets: value.map(set => {
                            return {
                                ...set,
                                bk_supplier_account: this.supplierAccount
                            }
                        })
                    }
                })
                return data || []
            },
            async createModule (value) {
                const data = await this.$store.dispatch('objectModule/createModule', {
                    bizId: this.business,
                    setId: this.createInfo.parentNode.data.bk_inst_id,
                    params: this.$injectMetadata({
                        ...value,
                        bk_supplier_account: this.supplierAccount
                    })
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
            }
        }
    }
</script>

<style lang="scss" scoped>
    .topology-tree-wrapper {
        height: 100%;
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
        &:hover,
        &.is-selected {
            .info-create-trigger {
                display: inline-block;
                & ~ .instance-num {
                    display: none;
                }
            }
        }
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
            &.is-leaf-icon {
                margin-left: 2px;
            }
        }
        .info-create-trigger {
            display: none;
            font-size: 0;
        }
        .node-button {
            height: 24px;
            padding: 0 6px;
            margin: 0 20px 0 4px;
            line-height: 22px;
            border-radius: 4px;
            font-size: 12px;
            min-width: auto;
            &.set-template-button {
                @include inlineBlock;
                font-style: normal;
                background-color: #dcdee5;
                color: #ffffff;
                outline: none;
                cursor: not-allowed;
            }
        }
        .instance-num {
            margin: 9px 20px 9px 5px;
            padding: 0 5px;
            height: 18px;
            line-height: 17px;
            border-radius: 2px;
            background-color: #f0f1f5;
            color: #979ba5;
            font-size: 12px;
            text-align: center;
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
        }
    }
    .topology-tree {
        height: 100%;
        .bk-big-tree-node.is-selected {
            .instance-num {
                background-color: #a2c5fd;
                color: #fff;
            }
        }
    }
</style>
