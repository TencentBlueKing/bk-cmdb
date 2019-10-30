<template>
    <div>
        <bk-big-tree ref="tree" class="tree"
            :show-checkbox="shouldShowCheckbox"
            :selectable="false"
            :show-link-line="true"
            :lazy-method="loadHost"
            :options="{
                idKey: getNodeId,
                nameKey: 'bk_inst_name',
                childrenKey: 'child'
            }"
            @node-click="handleNodeClick"
            @check-change="handleCheckedChange">
        </bk-big-tree>
        <bk-input class="filter"
            left-icon="icon-search"
            :placeholder="$t('筛选')"
            v-model="filter">
        </bk-input>
    </div>
</template>

<script>
    import { mapGetters } from 'vuex'
    export default {
        data () {
            return {
                filter: '',
                hostMap: {},
                request: {
                    host: Symbol('host')
                }
            }
        },
        computed: {
            ...mapGetters('objectBiz', ['bizId']),
            ...mapGetters('businessHost', ['getDefaultSearchCondition'])
        },
        watch: {
            '$parent.selected' (current, previous) {
                console.log(current === previous)
                this.syncState(current, previous)
            }
        },
        created () {
            this.initTopology()
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
            },
            getInstanceTopology () {
                return this.$store.dispatch('objectMainLineModule/getInstTopo', {
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
                return data.bk_obj_id === 'host'
            },
            async loadHost (node) {
                try {
                    const result = await this.searchHost(node.data.bk_inst_id)
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
                                console.log(this.$refs.tree.getNodeById(this.getNodeId(nodeData)))
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
            searchHost (moduleId) {
                const params = {
                    bk_biz_id: this.bizId,
                    ip: { data: [], exact: 0, flag: 'bk_host_innerip|bk_host_outerip' },
                    page: {
                        sort: 'bk_host_innerip'
                    },
                    condition: this.getDefaultSearchCondition()
                }
                const moduleCondition = params.condition.find(target => target.bk_obj_id === 'module')
                moduleCondition.condition.push({
                    field: 'bk_module_id',
                    operator: '$eq',
                    value: moduleId
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
            handleCheckedChange (checked, currentNode) {
                const item = this.hostMap[currentNode.data.bk_host_id]
                if (currentNode.checked) {
                    this.$parent.handleSelect(item)
                } else {
                    this.$parent.handleRemove(item)
                }
            }
        }
    }
</script>

<style lang="scss" scoped>
    .tree {
        height: calc(100% - 32px);
        padding: 0 0 0 20px;
        border-bottom: 1px solid $borderColor;
        @include scrollbar;
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
