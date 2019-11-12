<template>
    <div class="topology-tree-wrapper">
        <bk-big-tree class="topology-tree"
            ref="tree"
            v-bind="treeOptions"
            v-bkloading="{
                isLoading: $loading(['getTopologyData', 'getMainLine'])
            }"
            :options="{
                idKey: idGenerator,
                nameKey: 'bk_inst_name',
                childrenKey: 'child'
            }"
            selectable
            :before-select="beforeSelect"
            @select-change="handleSelectChange"
            @check-change="handleCheckChange"
        >
            <div class="node-info clearfix" :class="{ 'is-selected': node.selected }" slot-scope="{ node, data }">
                <i class="node-model-icon fl"
                    :class="{
                        'is-selected': node.selected,
                        'is-template': isTemplate(node),
                        'is-leaf-icon': node.isLeaf
                    }">
                    {{modelIconMap[data.bk_obj_id]}}
                </i>
                <span class="instance-num fr">√</span>
                <div class="info-content">
                    <span class="node-name">{{data.bk_inst_name}}</span>
                </div>
            </div>
        </bk-big-tree>
    </div>
</template>

<script>
    import { mapGetters } from 'vuex'
    export default {
        components: {
        },
        props: {
            treeOptions: {
                type: Object,
                default: () => {}
            }
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
            beforeSelect (node) {
                return node.data.bk_obj_id === 'module'
            },
            handleSelectChange (node) {
                this.$emit('selected', node)
                this.$store.commit('businessTopology/setSelectedNode', node)
            },
            handleCheckChange (id, checked) {
                console.log('handleCheckChange', id, checked)
                this.$emit('checked', id, checked)
            }
        }
    }
</script>

<style lang="scss" scoped>
    .topology-tree-wrapper {
        height: 100%;
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
