<template>
    <div>
        <cmdb-tree class="topology-tree"
            ref="tree"
            v-bkloading="{
                isLoading: $loading(['getTopologyData', 'getMainLine'])
            }"
            :options="{
                idKey: idGenerator,
                nameKey: 'bk_inst_name',
                childrenKey: 'child'
            }"
            expand-icon="bk-icon icon-down-shape"
            collapse-icon="bk-icon icon-right-shape">
            <div class="node-info clearfix" slot-scope="{ node, data }">
                <i :class="['node-model-icon fl', data.bk_obj_icon || modelIconMap[data.bk_obj_id]]"></i>
                <bk-button class="node-button fr"
                    type="primary"
                    v-if="showCreate(node, data)"
                    @click.stop="showCreateDialog(node)">
                    {{$t('Common[\'新建\']')}}
                </bk-button>
                <span class="node-name">{{data.bk_inst_name}}</span>
            </div>
        </cmdb-tree>
        <bk-dialog
            :is-show.sync="createInfo.show"
            :has-header="false"
            :has-footer="false"
            :padding="0"
            :quick-close="false"
            @after-transition-leave="handleAfterCancelCreateNode"
            @cancel="handleCancelCreateNode">
            <create-node v-if="createInfo.visible" slot="content"
                :properties="createInfo.properties"
                :parent-node="createInfo.parentNode"
                @submit="handleCreateNode"
                @cancel="handleCancelCreateNode">
            </create-node>
        </bk-dialog>
    </div>
</template>

<script>
    import { mapGetters } from 'vuex'
    import createNode from './create-node.vue'
    export default {
        components: {
            createNode
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
                    map[model.bk_obj_id] = model.bk_obj_icon
                })
                return map
            }
        },
        async created () {
            const [data, internal, mainLine] = await Promise.all([
                this.getTopologyData(),
                this.getInternalTopology(),
                this.getMainLine()
            ])
            data[0].child.unshift(...internal)
            this.treeData = data
            this.mainLine = mainLine
            const defaultNode = this.idGenerator(data[0])
            this.$nextTick(() => {
                this.$refs.tree.setData(data)
                this.$refs.tree.setSelected(defaultNode)
                this.$refs.tree.setExpanded(defaultNode)
            })
        },
        methods: {
            getTopologyData () {
                return this.$store.dispatch('objectMainLineModule/getInstTopo', {
                    bizId: this.business,
                    config: {
                        requestId: 'getTopologyData'
                    }
                })
            },
            async getInternalTopology () {
                const data = await this.$store.dispatch('objectMainLineModule/getInternalTopo', {
                    bizId: this.business
                })
                const moduleModel = this.$store.getters['objectModelClassify/getModelById']('module')
                return data.module.map(module => {
                    const isIdle = ['空闲机', 'idle machine'].includes(module.bk_module_name)
                    return {
                        'default': isIdle ? 1 : 2,
                        'bk_obj_id': 'module',
                        'bk_obj_name': moduleModel.bk_obj_name,
                        'bk_inst_id': module.bk_module_id,
                        'bk_inst_name': module.bk_module_name,
                        'bk_obj_icon': isIdle ? 'icon-cc-host-free-pool' : 'icon-cc-host-breakdown'
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
            idGenerator (data) {
                return `${data.bk_obj_id}_${data.bk_inst_id}`
            },
            showCreate (node, data) {
                const isModule = data.bk_obj_id === 'module'
                const isBlueKing = this.treeData[0].bk_inst_name === '蓝鲸'
                return node.selected && !isModule && !isBlueKing
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
                    const formData = {
                        ...value,
                        'bk_biz_id': this.business,
                        'bk_parent_id': parentNode.data.bk_inst_id
                    }
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
                        ...data
                    }
                    if (parentNode.data.bk_obj_id === 'biz') {
                        this.$refs.tree.addNode(nodeData, parentNode.id, 2)
                    } else {
                        this.$refs.tree.addNode(nodeData, parentNode.id, 0)
                    }
                    this.$success(this.$t('Common[\'新建成功\']'))
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
            }
        }
    }
</script>

<style lang="scss" scoped>
    .node-info {
        .node-model-icon {
            width: 20px;
            line-height: 32px;
            font-size: 18px;
            margin: 0 4px 0 6px;
        }
        .node-button {
            height: 24px;
            padding: 0 6px;
            margin: 4px;
            line-height: 22px;
            border-radius: 4px;
            font-size: 12px;
        }
        .node-name {
            display: block;
            line-height: 32px;
            font-size: 14px;
            @include ellipsis;
        }
    }
</style>
