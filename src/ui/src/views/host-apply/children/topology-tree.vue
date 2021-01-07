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
            :height="$APP.height - 160"
            :node-height="36"
            :check-on-click="true"
            :before-select="beforeSelect"
            :filter-method="filterMethod"
            @select-change="handleSelectChange"
            @check-change="handleCheckChange"
        >
            <div
                class="node-info clearfix"
                :title="(data.host_apply_rule_count === 0 && isDel) ? $t('暂无策略') : ''"
                :class="{ 'is-selected': node.selected }"
                slot-scope="{ node, data }"
            >
                <i class="internal-node-icon fl"
                    v-if="data.default !== 0"
                    :class="getInternalNodeClass(node, data)">
                </i>
                <i class="node-icon fl" v-else
                    :class="{
                        'is-selected': node.selected,
                        'is-template': isTemplate(node),
                        'is-leaf-icon': node.isLeaf
                    }">
                    {{modelIconMap[data.bk_obj_id]}}
                </i>
                <span v-show="applyEnabled(node)" class="config-icon fr"><i class="bk-cc-icon icon-cc-selected"></i></span>
                <div class="info-content">
                    <span class="node-name">{{data.bk_inst_name}}</span>
                </div>
            </div>
            <div slot="empty" class="empty">
                <span>{{$t('bk.bigTree.emptyText')}}</span>
            </div>
        </bk-big-tree>
    </div>
</template>

<script>
    import { mapGetters, mapState } from 'vuex'
    import Bus from '@/utils/bus'
    import { addResizeListener, removeResizeListener } from '@/utils/resize-events'
    export default {
        props: {
            treeOptions: {
                type: Object,
                default: () => ({})
            },
            action: {
                type: String,
                default: 'batch-edit'
            },
            checked: {
                type: Array,
                default: () => ([])
            }
        },
        data () {
            return {
                treeData: [],
                treeStat: {},
                mainLine: [],
                createInfo: {
                    show: false,
                    visible: false,
                    properties: [],
                    parentNode: null,
                    nextModelId: null
                },
                request: {
                    searchNode: Symbol('searchNode')
                },
                nodeIconMap: {
                    1: 'icon-cc-host-free-pool',
                    2: 'icon-cc-host-breakdown',
                    default: 'icon-cc-host-free-pool'
                }
            }
        },
        computed: {
            ...mapGetters(['supplierAccount']),
            ...mapState('hostApply', ['ruleDraft']),
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
            isDel () {
                return this.action === 'batch-del'
            }
        },
        watch: {
            action () {
                this.setNodeDisabled()
            }
        },
        async created () {
            Bus.$on('topology-search', this.handleSearch)
            const [data, mainLine] = await Promise.all([
                this.getTopologyData(),
                this.getMainLine()
            ])

            // 将空闲机池放到顶部
            const root = data[0] || {}
            const children = root.child || []
            const idleIndex = children.findIndex(item => item.default === 1)
            if (idleIndex !== -1) {
                const idlePool = children[idleIndex]
                idlePool.child.sort((a, b) => a.bk_inst_id - b.bk_inst_id)
                children.splice(idleIndex, 1)
                children.unshift(idlePool)
            }

            this.treeData = data
            this.mainLine = mainLine
            this.treeStat = this.getTreeStat()
            this.$nextTick(() => {
                this.setDefaultState(data)
            })
        },
        mounted () {
            addResizeListener(this.$el, this.handleResize)
        },
        beforeDestroy () {
            Bus.$off('topology-search', this.handleSearch)
            removeResizeListener(this.$el, this.handleResize)
        },
        methods: {
            async handleSearch (params) {
                try {
                    if (params.query_filter.rules.length) {
                        const data = await this.$store.dispatch('hostApply/searchNode', {
                            bizId: this.business,
                            params: params,
                            config: {
                                requestId: this.request.searchNode
                            }
                        })
                        this.$refs.tree.filter(data)
                    } else {
                        this.$refs.tree.filter()
                    }
                    if (this.checked.length) {
                        this.$refs.tree.setChecked(this.checked)
                    }
                } catch (e) {
                    console.error(e)
                }
            },
            filterMethod (remoteData, node) {
                return remoteData.some(datum => datum.bk_inst_id === node.data.bk_inst_id && datum.bk_obj_id === node.data.bk_obj_id)
            },
            setDefaultState (data) {
                this.$refs.tree.setData(data)
                let defaultNodeId
                const queryModule = parseInt(this.$route.query.module)
                const firstModule = this.treeStat.firstModule
                if (!isNaN(queryModule)) {
                    defaultNodeId = `module_${queryModule}`
                } else if (this.ruleDraft.moduleIds) {
                    defaultNodeId = `module_${this.ruleDraft.moduleIds[0]}`
                } else if (firstModule) {
                    defaultNodeId = this.idGenerator(firstModule)
                }
                if (defaultNodeId) {
                    this.$refs.tree.setSelected(defaultNodeId, { emitEvent: true })
                    this.$refs.tree.setExpanded(defaultNodeId)
                }
            },
            getTreeStat () {
                const stat = {
                    firstModule: null,
                    levels: {},
                    noRuleIds: []
                }
                const findModule = function (data, parent) {
                    for (const item of data) {
                        stat.levels[item.bk_inst_id] = parent ? (stat.levels[parent.bk_inst_id] + 1) : 0
                        if (item.bk_obj_id === 'module') {
                            if (item.host_apply_rule_count === 0) {
                                stat.noRuleIds.push(item.bk_inst_id)
                            }
                            if (!stat.firstModule) {
                                stat.firstModule = item
                            }
                        } else if (item.child.length) {
                            const match = findModule(item.child, item)
                            if (match && !stat.firstModule) {
                                stat.firstModule = item
                            }
                        }
                    }
                }
                findModule(this.treeData)
                return stat
            },
            setNodeDisabled () {
                const nodeIds = this.treeStat.noRuleIds.map(id => `module_${id}`)
                this.$refs.tree.setDisabled(nodeIds, { emitEvent: true, disabled: this.isDel })
            },
            updateNodeStatus (id, isClear) {
                const nodeData = this.$refs.tree.getNodeById(`module_${id}`).data
                nodeData.host_apply_enabled = false
                if (isClear) {
                    nodeData.host_apply_rule_count = 0
                }
                this.treeStat = this.getTreeStat()
            },
            getTopologyData () {
                return this.$store.dispatch('objectMainLineModule/getInstTopoInstanceNum', {
                    bizId: this.business,
                    config: {
                        params: {
                            with_default: 1
                        },
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
            applyEnabled (node) {
                return this.isModule(node) && node.data.host_apply_enabled
            },
            isTemplate (node) {
                return node.data.service_template_id || node.data.set_template_id
            },
            isModule (node) {
                return node.data.bk_obj_id === 'module'
            },
            async beforeSelect (node) {
                return this.isModule(node)
            },
            getInternalNodeClass (node, data) {
                const clazz = []
                clazz.push(this.nodeIconMap[data.default] || this.nodeIconMap.default)
                if (node.selected) {
                    clazz.push('is-selected')
                }
                return clazz
            },
            handleSelectChange (node) {
                this.$emit('selected', node)
            },
            handleCheckChange (id, checked) {
                this.$emit('checked', id, checked)
            },
            handleResize () {
                this.$refs.tree.resize()
            }
        }
    }
</script>

<style lang="scss" scoped>
    .topology-tree {
        padding: 10px 0;
        margin-right: 2px;
        .node-info {
            .node-icon {
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
            .config-icon {
                position: relative;
                top: 6px;
                right: 4px;
                padding: 0 5px;
                height: 18px;
                line-height: 17px;
                color: #979ba5;
                font-size: 26px;
                text-align: center;
                color: #2dcb56;
            }
            .internal-node-icon{
                width: 20px;
                height: 20px;
                line-height: 20px;
                text-align: center;
                margin: 8px 4px 8px 0;
                &.is-selected {
                    color: #ffb400;
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
            }
        }

        .empty {
            position: absolute;
            display: flex;
            height: calc(100% - 30px);
            width: 100%;
            justify-content: center;
            align-items: center;
        }
    }
</style>
