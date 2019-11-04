<template>
    <section class="tree-layout" v-bkloading="{ isLoading: $loading(Object.values(request)) }">
        <bk-input class="tree-search" v-model="filter"></bk-input>
        <bk-big-tree ref="tree" class="topology-tree"
            selectable
            :expand-on-click="false"
            :options="{
                idKey: getNodeId,
                nameKey: 'bk_inst_name',
                childrenKey: 'child'
            }"
            @select-change="handleSelectChange">
            <template slot-scope="{ node, data }">
                <i class="internal-node-icon fl"
                    v-if="data.default !== 0"
                    :class="getInternalNodeClass(node, data)">
                </i>
                <i v-else
                    :class="['node-icon fl', { 'is-selected': node.selected }]">
                    {{data.bk_obj_name[0]}}
                </i>
                <span :class="['node-count fr', { 'is-selected': node.selected }]">
                    {{getNodeCount(data)}}
                </span>
                <span class="node-name">{{node.name}}</span>
            </template>
        </bk-big-tree>
    </section>
</template>

<script>
    import { mapGetters } from 'vuex'
    import debounce from 'lodash.debounce'
    export default {
        data () {
            return {
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
                    internal: Symbol('internal')
                }
            }
        },
        computed: {
            ...mapGetters('objectBiz', ['bizId'])
        },
        watch: {
            filter (value) {
                this.handleFilter()
            }
        },
        created () {
            this.handleFilter = debounce(() => {
                this.$refs.tree.filter(this.filter)
            }, 300)
            this.initTopology()
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
                return this.$store.dispatch('objectMainLineModule/getInstTopo', {
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
                return count
            },
            handleSelectChange (node) {
                this.$store.commit('businessHost/setCurrentNode', node)
            }
        }
    }
</script>

<style lang="scss" scoped>
    .tree-search {
        display: block;
        width: auto;
        margin: 0 20px;
    }
    .topology-tree {
        width: 100%;
        height: calc(100vh - 180px);
        padding: 10px 0;
        @include scrollbar-y;
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
            &.is-selected {
                background-color: #3A84FF;
            }
        }
        .node-name {
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
</style>
