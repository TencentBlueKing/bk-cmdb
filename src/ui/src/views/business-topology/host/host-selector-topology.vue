<template>
    <div>
        <div class="tree-wrapper">
            <bk-big-tree ref="tree" class="tree"
                display-matched-node-descendants
                :lazy-method="loadHost"
                :lazy-disabled="isLazyDisabled"
                :show-checkbox="shouldShowCheckbox"
                :selectable="false"
                :options="{
                    idKey: getNodeId,
                    nameKey: 'bk_inst_name',
                    childrenKey: 'child'
                }"
                :style="{
                    'min-width': `calc(100% + ${deepestExpandedLevel * 60}px)`
                }"
                :height="270"
                :node-height="36"
                :filter-method="filterMethod"
                :before-check="beforeCheck"
                @node-click="handleNodeClick"
                @check-change="handleCheckedChange"
                @expand-change="handleExpandChange">
                <div class="node-info clearfix" slot-scope="{ node, data }">
                    <template v-if="data.bk_obj_id !== 'host'">
                        <i class="internal-node-icon fl"
                            v-if="data.default !== 0"
                            :class="getInternalNodeClass(node, data)">
                        </i>
                        <i v-else
                            :class="['node-icon fl', { 'is-selected': node.selected, 'is-template': isTemplate(node) }]">
                            {{data.bk_obj_name[0]}}
                        </i>
                    </template>
                    <span class="node-count fr" v-if="data.bk_obj_id !== 'host'">
                        {{getNodeCount(data)}}
                    </span>
                    <span class="node-name">{{node.name}}</span>
                </div>
            </bk-big-tree>
        </div>
        <bk-input class="filter"
            clearable
            left-icon="icon-search"
            :placeholder="$t('筛选')"
            v-model.trim="filter">
        </bk-input>
    </div>
</template>

<script>
    import { mapGetters } from 'vuex'
    import debounce from 'lodash.debounce'
    export default {
        data () {
            return {
                filter: '',
                filterMethod: this.defaultFilterMethod,
                expandedNodes: [],
                handleFilter: () => ({}),
                hostMap: {},
                request: {
                    host: Symbol('host')
                },
                nodeIconMap: {
                    1: 'icon-cc-host-free-pool',
                    2: 'icon-cc-host-breakdown',
                    default: 'icon-cc-host-free-pool'
                }
            }
        },
        computed: {
            ...mapGetters('objectBiz', ['bizId']),
            ...mapGetters('businessHost', ['getDefaultSearchCondition']),
            deepestExpandedLevel () {
                const maxLevel = Math.max.apply(null, [0, ...this.expandedNodes.map(node => node.level)])
                return Math.max(4, maxLevel) - 4
            },
            limitDisplay () {
                return !!this.$parent.displayNodes.length
            }
        },
        watch: {
            '$parent.selected' (current, previous) {
                this.syncState(current, previous)
            },
            filter () {
                this.handleFilter()
            }
        },
        created () {
            this.handleFilter = debounce(() => {
                this.$refs.tree.filter(this.limitDisplay ? (this.filter || Symbol('any')) : this.filter)
                this.recaculateLine()
            }, 300)
            this.filterMethod = this.limitDisplay ? this.displayNodesFilterMethod : this.defaultFilterMethod
            this.initTopology()
        },
        activated () {
            this.filter = ''
        },
        methods: {
            defaultFilterMethod (keyword, node) {
                return String(node.name).toLowerCase().indexOf(keyword) > -1
            },
            displayNodesFilterMethod (keyword, node) {
                const displayNodes = this.$parent.displayNodes
                if (this.filter) {
                    return node.data.bk_obj_id === 'host' && node.name.indexOf(keyword) > -1
                }
                return displayNodes.includes(node.id) || node.data.bk_obj_id === 'host'
            },
            recaculateLine () {
                if (this.limitDisplay) {
                    const tree = this.$refs.tree
                    const displayNodes = this.$parent.displayNodes
                    tree.needsCalculateNodes.push(...displayNodes.map(id => tree.getNodeById(id)))
                }
            },
            syncState (current, previous) {
                const unselectHost = previous.filter(prev => {
                    const exist = current.some(cur => cur.host.bk_host_id === prev.host.bk_host_id)
                    return !exist
                })
                this.syncCheckedState(current, true)
                this.syncCheckedState(unselectHost, false)
            },
            syncCheckedState (list, checked) {
                const hosts = list.map(item => item.host.bk_host_id)
                const nodes = this.$refs.tree.nodes.filter(node => hosts.includes(node.data.bk_host_id))
                this.$refs.tree.setChecked(nodes.map(node => node.id), { checked, beforeCheck: false })
            },
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
                        child: this.$tools.sort((internal.module || []), 'default').map(module => ({
                            bk_obj_id: 'module',
                            bk_inst_id: module.bk_module_id,
                            bk_inst_name: module.bk_module_name,
                            host_count: module.host_count,
                            service_instance_count: module.service_instance_count,
                            default: module.default
                        }))
                    }
                    children.unshift(idlePool)
                    this.$refs.tree.setData(topology)
                    this.syncState(this.$parent.selected, [])
                    if (this.limitDisplay) {
                        this.$refs.tree.filter(Symbol('ignore'))
                        this.$refs.tree.setExpanded(this.$parent.displayNodes)
                    } else {
                        const defaultNodeId = this.getNodeId(topology[0])
                        this.$refs.tree.setExpanded(defaultNodeId)
                    }
                } catch (e) {
                    console.error(e)
                }
            },
            getInstanceTopology () {
                return this.$store.dispatch('objectMainLineModule/getInstTopoInstanceNum', {
                    bizId: this.bizId,
                    config: {
                        requestId: this.$parent.request.instance
                    }
                })
            },
            getInternalTopology () {
                return this.$store.dispatch('objectMainLineModule/getInternalTopo', {
                    bizId: this.bizId,
                    config: {
                        requestId: this.$parent.request.internal
                    }
                })
            },
            getNodeId (data) {
                return `${data.bk_obj_id}-${data.bk_inst_id}`
            },
            shouldShowCheckbox (data) {
                if (data.bk_obj_id === 'host') {
                    return true
                }
                if (data.bk_obj_id === 'module' && data.host_count > 0) {
                    return true
                }
                return false
            },
            async loadHost (node) {
                try {
                    const result = await this.searchHost(node)
                    const data = []
                    result.info.forEach(item => {
                        const nodeData = {
                            bk_obj_id: 'host',
                            bk_inst_id: `${node.id}-${item.host.bk_host_id}`, // 额外加上父节点id，防止不同模块下的主机id重复
                            bk_inst_name: item.host.bk_host_innerip,
                            bk_host_id: item.host.bk_host_id,
                            item: item
                        }
                        data.push(nodeData)
                        this.$set(this.hostMap, item.host.bk_host_id, item)
                    })
                    this.$set(node.data, 'child', result.info)
                    setTimeout(() => {
                        data.forEach(nodeData => {
                            const isSelected = this.$parent.selected.some(item => item.host.bk_host_id === nodeData.bk_host_id)
                            if (isSelected) {
                                this.$refs.tree.setChecked(this.getNodeId(nodeData), { beforeCheck: false })
                            }
                        })
                    }, 0)
                    return { data }
                } catch (e) {
                    console.log(e)
                    return { data: [] }
                }
            },
            isLazyDisabled (node) {
                return node.data.bk_obj_id === 'host' || node.data.host_count === 0
            },
            searchHost (node) {
                const params = {
                    bk_biz_id: this.bizId,
                    ip: { data: [], exact: 0, flag: 'bk_host_innerip|bk_host_outerip' },
                    page: {
                        sort: 'bk_host_innerip'
                    },
                    condition: this.getDefaultSearchCondition()
                }
                const modelId = node.data.bk_obj_id
                const fieldMap = {
                    biz: 'bk_biz_id',
                    set: 'bk_set_id',
                    module: 'bk_module_id',
                    host: 'bk_host_id'
                }
                const targetCondition = params.condition.find(target => target.bk_obj_id === modelId)
                targetCondition.condition.push({
                    field: fieldMap[modelId] || 'bk_inst_id',
                    operator: '$eq',
                    value: node.data.bk_inst_id
                })
                return this.$store.dispatch('hostSearch/searchHost', {
                    params: params,
                    config: {
                        requestId: this.$parent.request.host
                    }
                })
            },
            handleNodeClick (node) {
                if (node.data.bk_obj_id === 'host') {
                    this.$refs.tree.setChecked(node.id, { checked: !node.checked, emitEvent: true, beforeCheck: false })
                }
            },
            async handleCheckedChange (checked, selectedNode) {
                const hosts = []
                if (selectedNode.data.bk_obj_id === 'host') {
                    hosts.push(selectedNode.data.item)
                } else {
                    const descendants = selectedNode.descendants.filter(node => node.data.bk_obj_id === 'host')
                    hosts.push(...descendants.map(node => node.data.item))
                }
                if (selectedNode.checked) {
                    this.$parent.handleSelect(hosts)
                } else {
                    this.$parent.handleRemove(hosts)
                }
            },
            handleExpandChange (targetNode) {
                const nodes = [targetNode, ...targetNode.descendants]
                const ids = nodes.map(node => node.id)
                if (targetNode.expanded) {
                    this.expandedNodes.push(...nodes.filter(node => node.expanded))
                } else {
                    this.expandedNodes = this.expandedNodes.filter(exist => !ids.includes(exist.id))
                }
            },
            getInternalNodeClass (node, data) {
                return this.nodeIconMap[data.default] || this.nodeIconMap.default
            },
            isTemplate (node) {
                return node.data.service_template_id || node.data.set_template_id
            },
            getNodeCount (data) {
                const count = data.host_count
                if (typeof count === 'number') {
                    return count > 999 ? '999+' : count
                }
                return 0
            },
            async beforeCheck (node) {
                if (node.lazy) {
                    const { data } = await this.loadHost(node)
                    this.$refs.tree.addNode(data, node.id)
                    return true
                }
                return true
            }
        }
    }
</script>

<style lang="scss" scoped>
    .tree-wrapper {
        height: calc(100% - 32px);
        border-bottom: 1px solid $borderColor;
        @include scrollbar;
    }
    .tree {
        padding: 0 0 0 20px;
        @include scrollbar-x;
        .node-icon {
            display: block;
            width: 20px;
            height: 20px;
            margin: 8px 4px 8px 0;
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
            &.set-template-button {
                @include inlineBlock;
                font-style: normal;
                background-color: #dcdee5;
                color: #ffffff;
                outline: none;
                cursor: not-allowed;
            }
        }
    }
    .filter {
        /deep/ {
            .bk-form-input {
                border: none;
                border-radius: 0;
            }
            .bk-icon {
                display: inline;
                vertical-align: initial;
            }
        }
    }
</style>
