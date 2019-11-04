<template>
    <div>
        <div class="tree-wrapper">
            <bk-big-tree ref="tree" class="tree"
                :show-checkbox="shouldShowCheckbox"
                :selectable="false"
                :show-link-line="true"
                :options="{
                    idKey: getNodeId,
                    nameKey: 'bk_inst_name',
                    childrenKey: 'child'
                }"
                :style="{
                    'min-width': `calc(100% + ${deepestExpandedLevel * 30}px)`
                }"
                @node-click="handleNodeClick"
                @check-change="handleCheckedChange"
                @expand-change="handleExpandChange">
            </bk-big-tree>
        </div>
        <bk-input class="filter"
            left-icon="icon-search"
            :placeholder="$t('筛选')"
            v-model="filter">
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
                expandedNodes: [],
                handleFilter: () => ({}),
                hostMap: {},
                request: {
                    host: Symbol('host')
                }
            }
        },
        computed: {
            ...mapGetters('objectBiz', ['bizId']),
            ...mapGetters('businessHost', ['getDefaultSearchCondition']),
            deepestExpandedLevel () {
                const maxLevel = Math.max.apply(null, [0, ...this.expandedNodes.map(node => node.level)])
                return Math.max(2, maxLevel) - 2
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
                this.$refs.tree.filter(this.filter)
            }, 300)
            this.initTopology()
        },
        activated () {
            this.filter = ''
        },
        methods: {
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
                this.$refs.tree.setChecked(nodes.map(node => node.id), { checked })
            },
            async initTopology () {
                try {
                    const [topology, internal, allHost] = await Promise.all([
                        this.getInstanceTopology(),
                        this.getInternalTopology(),
                        this.getAllHost()
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
                    const defaultNodeId = this.getNodeId(topology[0])
                    this.$refs.tree.setData(topology)
                    this.$refs.tree.setExpanded(defaultNodeId)
                    this.appendHostNode(allHost)
                } catch (e) {
                    console.error(e)
                }
            },
            appendHostNode (hostList) {
                const dataMap = {}
                hostList.forEach(item => {
                    item.module.map(module => {
                        const moduleId = module.bk_module_id
                        if (dataMap.hasOwnProperty(moduleId)) {
                            dataMap[moduleId].push(item)
                        } else {
                            dataMap[moduleId] = [item]
                        }
                    })
                })
                Object.keys(dataMap).forEach(moduleId => {
                    const moduleNodeId = this.getNodeId({ bk_obj_id: 'module', bk_inst_id: moduleId })
                    this.$refs.tree.addNode(dataMap[moduleId].map(item => ({
                        bk_obj_id: 'host',
                        bk_inst_id: `${moduleNodeId}-${item.host.bk_host_id}`,
                        bk_inst_name: item.host.bk_host_innerip,
                        bk_host_id: item.host.bk_host_id,
                        item: item
                    })), {
                        parentId: moduleNodeId,
                        expandParent: false
                    })
                })
            },
            async getAllHost () {
                const params = {
                    bk_biz_id: this.bizId,
                    ip: { data: [], exact: 0, flag: 'bk_host_innerip|bk_host_outerip' },
                    page: {},
                    condition: this.getDefaultSearchCondition()
                }
                const result = await this.$store.dispatch('hostSearch/searchHost', {
                    params: params,
                    config: {
                        requestId: this.$parent.request.host
                    }
                })
                return result.info
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
                return data.bk_obj_id !== 'biz' && data.host_count > 0
            },
            async loadHost (node) {
                try {
                    const result = await this.searchHost(node)
                    const data = []
                    const leaf = []
                    result.info.forEach(item => {
                        const nodeData = {
                            bk_obj_id: 'host',
                            bk_inst_id: `${node.id}-${item.host.bk_host_id}`, // 额外加上父节点id，防止不同模块下的主机id重复
                            bk_inst_name: item.host.bk_host_innerip,
                            bk_host_id: item.host.bk_host_id
                        }
                        data.push(nodeData)
                        leaf.push(this.getNodeId(nodeData))
                        this.$set(this.hostMap, item.host.bk_host_id, item)
                    })
                    this.$set(node.data, 'child', result.info)
                    setTimeout(() => {
                        data.forEach(nodeData => {
                            const isSelected = this.$parent.selected.some(item => item.host.bk_host_id === nodeData.bk_host_id)
                            if (isSelected) {
                                this.$refs.tree.setChecked(this.getNodeId(nodeData))
                            }
                        })
                    }, 0)
                    return { data, leaf }
                } catch (e) {
                    console.log(e)
                    return { data: [], leaf: [] }
                }
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
                    this.$refs.tree.setChecked(node.id, { checked: !node.checked, emitEvent: true })
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
