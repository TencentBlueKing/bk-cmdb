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
            :show-link-line="false"
            :default-expanded-nodes="defaultExpandedNodes"
            expand-icon="bk-icon icon-down-shape"
            collapse-icon="bk-icon icon-right-shape">
            <div class="node-info clearfix" slot-scope="{ node, data }">
                <i :class="['node-model-icon fl', modelIconMap[data.bk_obj_id]]"></i>
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
                @on-submit="handleCreateNode"
                @on-cancel="handleCancelCreateNode">
            </create-node>
        </bk-dialog>
    </div>
</template>

<script>
    import createNode from './create-node.vue'
    export default {
        components: {
            createNode
        },
        data () {
            return {
                defaultExpandedNodes: [],
                treeData: [],
                mainLine: [],
                createInfo: {
                    show: false,
                    visible: false,
                    properties: [],
                    parentNode: null
                }
            }
        },
        computed: {
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
            const [data, mainLine] = await Promise.all([
                this.getTopologyData(),
                this.getMainLine()
            ])
            this.treeData = data
            this.mainLine = mainLine
            this.defaultExpandedNodes = [this.idGenerator(data[0])]
            this.$nextTick(() => {
                this.$refs.tree.setData(data)
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
                this.createInfo.parentNode = node
                this.createInfo.show = true
                this.createInfo.visible = true
                const nodeModel = this.mainLine.find(data => data.bk_obj_id === node.data.bk_obj_id)
                const nextModelId = nodeModel.bk_next_obj
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
            },
            handleCancelCreateNode () {
                this.createInfo.show = false
            },
            handleCreateNode () {}
        }
    }
</script>

<style lang="scss" scoped>
    .node-info {
        .node-model-icon {
            width: 20px;
            line-height: 32px;
            font-size: 18px;
            margin-right: 4px;
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
