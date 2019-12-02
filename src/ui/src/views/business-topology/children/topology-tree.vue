<template>
    <section class="tree-layout" v-bkloading="{ isLoading: $loading(Object.values(request)) }">
        <bk-input class="tree-search"
            clearable
            right-icon="bk-icon icon-search"
            :placeholder="$t('请输入关键词')"
            v-model="filter">
        </bk-input>
        <bk-big-tree ref="tree" class="topology-tree"
            selectable
            :expand-on-click="false"
            :style="{
                height: $APP.height - 160 + 'px'
            }"
            :options="{
                idKey: getNodeId,
                nameKey: 'bk_inst_name',
                childrenKey: 'child'
            }"
            @select-change="handleSelectChange">
            <div :class="['node-info clearfix', { 'is-selected': node.selected }]" slot-scope="{ node, data }">
                <i class="internal-node-icon fl"
                    v-if="data.default !== 0"
                    :class="getInternalNodeClass(node, data)">
                </i>
                <i v-else
                    :class="['node-icon fl', { 'is-selected': node.selected, 'is-template': isTemplate(node) }]">
                    {{data.bk_obj_name[0]}}
                </i>
                <cmdb-auth v-if="showCreate(node, data)"
                    class="info-create-trigger fr"
                    :auth="$authResources({ type: $OPERATION.C_TOPO })">
                    <template slot-scope="{ disabled }">
                        <i v-if="isBlueKing"
                            class="node-button disabled-node-button"
                            v-bk-tooltips="{ content: $t('蓝鲸业务拓扑节点提示'), placement: 'top' }">
                            {{$t('新建')}}
                        </i>
                        <i v-else-if="data.set_template_id"
                            class="node-button disabled-node-button"
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
                <span :class="['node-count fr', { 'is-selected': node.selected }]">
                    {{getNodeCount(data)}}
                </span>
                <span class="node-name">{{node.name}}</span>
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
    </section>
</template>

<script>
    import { mapGetters } from 'vuex'
    import debounce from 'lodash.debounce'
    import CreateNode from './create-node.vue'
    import CreateSet from './create-set.vue'
    import CreateModule from './create-module.vue'
    import Bus from '@/utils/bus'
    export default {
        components: {
            CreateNode,
            CreateSet,
            CreateModule
        },
        props: {
            active: {
                type: String,
                required: true
            }
        },
        data () {
            return {
                isBlueKing: false,
                filter: '',
                handleFilter: () => ({}),
                nodeCountType: 'host_count',
                nodeIconMap: {
                    1: 'icon-cc-host-free-pool',
                    2: 'icon-cc-host-breakdown',
                    default: 'icon-cc-host-free-pool'
                },
                request: {
                    instance: Symbol('instance'),
                    internal: Symbol('internal'),
                    property: Symbol('property')
                },
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
            ...mapGetters('objectBiz', ['bizId']),
            ...mapGetters('businessHost', ['topologyModels', 'propertyMap'])
        },
        watch: {
            filter (value) {
                this.handleFilter()
            },
            active (value) {
                const map = {
                    hostList: 'host_count',
                    serviceInstance: 'service_instance_count'
                }
                if (Object.keys(map).includes(value)) {
                    this.nodeCountType = map[value]
                }
            }
        },
        created () {
            Bus.$on('refresh-count', this.refreshCount)
            this.handleFilter = debounce(() => {
                this.$refs.tree.filter(this.filter)
            }, 300)
            this.initTopology()
        },
        beforeDestroy () {
            Bus.$off('refresh-count', this.refreshCount)
        },
        methods: {
            async initTopology () {
                try {
                    const [topology, internal] = await Promise.all([
                        this.getInstanceTopology(),
                        this.getInternalTopology()
                    ])
                    const root = topology[0] || {}
                    const children = root.child || []
                    const idlePool = {
                        bk_obj_id: 'set',
                        bk_inst_id: internal.bk_set_id,
                        bk_inst_name: internal.bk_set_name,
                        host_count: internal.host_count,
                        service_instance_count: internal.service_instance_count,
                        default: internal.default,
                        is_idle_set: true,
                        child: (internal.module || []).map(module => ({
                            bk_obj_id: 'module',
                            bk_inst_id: module.bk_module_id,
                            bk_inst_name: module.bk_module_name,
                            host_count: module.host_count,
                            service_instance_count: module.service_instance_count,
                            default: module.default
                        }))
                    }
                    children.unshift(idlePool)
                    this.isBlueKing = root.bk_inst_name === '蓝鲸'
                    this.$refs.tree.setData(topology)
                    this.setDefaultState()
                } catch (e) {
                    console.error(e)
                }
            },
            setDefaultState () {
                const businessNodeId = this.$refs.tree.nodes[0].id
                const queryNodeId = this.$route.query.node
                let defaultNodeId = businessNodeId
                if (queryNodeId) {
                    const node = this.$refs.tree.getNodeById(queryNodeId)
                    defaultNodeId = node ? queryNodeId : businessNodeId
                }
                this.$refs.tree.setExpanded(defaultNodeId)
                this.$refs.tree.setSelected(defaultNodeId, { emitEvent: true })
            },
            getInstanceTopology () {
                return this.$store.dispatch('objectMainLineModule/getInstTopoInstanceNum', {
                    bizId: this.bizId,
                    config: {
                        requestId: this.request.instance
                    }
                })
            },
            getInternalTopology () {
                return this.$store.dispatch('objectMainLineModule/getInternalTopo', {
                    bizId: this.bizId,
                    config: {
                        requestId: this.request.internal
                    }
                })
            },
            getNodeId (data) {
                return `${data.bk_obj_id}-${data.bk_inst_id}`
            },
            getInternalNodeClass (node, data) {
                const clazz = []
                clazz.push(this.nodeIconMap[data.default] || this.nodeIconMap.default)
                if (node.selected) {
                    clazz.push('is-selected')
                }
                return clazz
            },
            getNodeCount (data) {
                const count = data[this.nodeCountType]
                if (typeof count === 'number') {
                    return count > 999 ? '999+' : count
                }
                return 0
            },
            handleSelectChange (node) {
                this.$store.commit('businessHost/setSelectedNode', node)
                Bus.$emit('toggle-host-filter', false)
                if (!node.expanded) {
                    this.$refs.tree.setExpanded(node.id)
                }
            },
            showCreate (node, data) {
                const isModule = data.bk_obj_id === 'module'
                const isIdleSet = data.is_idle_set
                return !isModule && !isIdleSet
            },
            async showCreateDialog (node) {
                const nodeModel = this.topologyModels.find(data => data.bk_obj_id === node.data.bk_obj_id)
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
                            requestId: this.request.property
                        }
                    })
                    if (!['set', 'module'].includes(nextModelId)) {
                        this.$store.commit('businessHost/setProperties', {
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
                        'bk_biz_id': this.bizId,
                        'bk_parent_id': parentNode.data.bk_inst_id
                    })
                    const nextModelId = this.createInfo.nextModelId
                    const nextModel = this.topologyModels.find(model => model.bk_obj_id === nextModelId)
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
                        host_count: 0,
                        service_instance_count: 0,
                        service_template_id: value.service_template_id,
                        ...data
                    }
                    this.$refs.tree.addNode(nodeData, parentNode.id, parentNode.data.bk_obj_id === 'biz' ? 1 : 0)
                    this.$success(this.$t('新建成功'))
                    this.createInfo.show = false
                } catch (e) {
                    console.error(e)
                }
            },
            async handleCreateSetNode (value) {
                try {
                    const parentNode = this.createInfo.parentNode
                    const nextModel = this.topologyModels.find(model => model.bk_obj_id === 'set')
                    const formData = (value.sets || []).map(set => {
                        return this.$injectMetadata({
                            ...set,
                            'bk_biz_id': this.bizId,
                            'bk_parent_id': parentNode.data.bk_inst_id
                        })
                    })
                    const data = await this.createSet(formData)
                    const insertBasic = parentNode.data.bk_obj_id === 'biz' ? 1 : 0
                    data && data.forEach((set, index) => {
                        if (set.data) {
                            const nodeData = {
                                default: 0,
                                child: [],
                                bk_obj_name: nextModel.bk_obj_name,
                                bk_obj_id: nextModel.bk_obj_id,
                                host_count: 0,
                                service_instance_count: 0,
                                service_template_id: value.service_template_id,
                                bk_inst_id: set.data.bk_set_id,
                                bk_inst_name: set.data.bk_set_name,
                                set_template_id: value.set_template_id
                            }
                            this.$refs.tree.addNode(nodeData, parentNode.id, insertBasic + index)
                            if (value.set_template_id) {
                                this.addModulesInSetTemplate(nodeData, set.data.bk_set_id)
                            }
                        } else {
                            this.$error(set.error_message)
                        }
                    })
                    this.$success(this.$t('新建成功'))
                    this.createInfo.show = false
                } catch (e) {
                    console.error(e)
                }
            },
            async addModulesInSetTemplate (parentNodeData, id) {
                const modules = await this.$store.dispatch('objectModule/searchModule', {
                    bizId: this.bizId,
                    setId: id,
                    params: this.$injectMetadata(),
                    config: {
                        requestId: 'searchModule'
                    }
                })
                const parentNodeId = this.getNodeId(parentNodeData)
                const nextModel = this.topologyModels.find(model => model.bk_obj_id === 'module')
                modules.info && modules.info.forEach(_module => {
                    const nodeData = {
                        default: 0,
                        child: [],
                        bk_obj_name: nextModel.bk_obj_name,
                        bk_obj_id: nextModel.bk_obj_id,
                        host_count: 0,
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
                    bizId: this.bizId,
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
                    bizId: this.bizId,
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
            isTemplate (node) {
                return node.data.service_template_id || node.data.set_template_id
            },
            refreshCount ({ type, hosts, target }) {
                hosts.forEach(data => {
                    data.module.forEach(module => {
                        if (!target || target.data.bk_inst_id !== module.bk_module_id) {
                            const node = this.$refs.tree.getNodeById(`module-${module.bk_module_id}`)
                            const nodes = node ? [node, ...node.parents] : []
                            nodes.forEach(exist => {
                                exist.data[type]--
                            })
                        }
                    })
                })
                if (target) {
                    const targetNode = this.$refs.tree.getNodeById(`module-${target.data.bk_inst_id}`)
                    if (targetNode) {
                        [targetNode, ...targetNode.parents].forEach(exist => {
                            exist.data[type] = exist.data[type] + hosts.length
                        })
                    }
                }
            }
        }
    }
</script>

<style lang="scss" scoped>
    .tree-layout {
        overflow: hidden;
    }
    .tree-search {
        display: block;
        width: auto;
        margin: 0 20px;
    }
    .topology-tree {
        padding: 10px 0;
        margin-right: 4px;
        @include scrollbar-y(6px);
        .node-icon {
            display: block;
            width: 20px;
            height: 20px;
            margin: 8px 4px 8px 0;
            vertical-align: middle;
            border-radius: 50%;
            background-color: #C4C6CC;
            line-height: 1.666667;
            text-align: center;
            font-size: 12px;
            font-style: normal;
            color: #FFF;
            &.is-template {
                background-color: #97aed6;
            }
            &.is-selected {
                background-color: #3A84FF;
            }
        }
        .node-name {
            display: block;
            height: 36px;
            line-height: 36px;
            overflow: hidden;
            @include ellipsis;
        }
        .node-count {
            padding: 0 5px;
            margin: 9px 20px 9px 4px;
            height: 18px;
            line-height: 17px;
            border-radius: 2px;
            background-color: #f0f1f5;
            color: #979ba5;
            font-size: 12px;
            text-align: center;
            &.is-selected {
                background-color: #a2c5fd;
                color: #fff;
            }
        }
        .internal-node-icon{
            width: 20px;
            height: 20px;
            line-height: 20px;
            text-align: center;
            margin: 8px 4px 8px 0;
            &.is-selected {
                color: #FFB400;
            }
        }
    }
    .node-info {
        &:hover,
        &.is-selected {
            .info-create-trigger {
                display: inline-block;
                & ~ .node-count {
                    display: none;
                }
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
            &.disabled-node-button {
                @include inlineBlock;
                line-height: 24px;
                font-style: normal;
                background-color: #dcdee5;
                color: #ffffff;
                outline: none;
                cursor: not-allowed;
            }
        }
    }
</style>
