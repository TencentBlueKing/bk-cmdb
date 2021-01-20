<template>
    <div class="host-selector-topology">
        <cmdb-resize-layout class="tree-layout"
            direction="right"
            :handler-offset="3"
            :min="281"
            :max="420">
            <div class="topo-wrapper">
                <bk-input class="search-input" ref="filterInput" clearable
                    right-icon="icon-search"
                    :placeholder="$t('搜索拓扑节点')"
                    v-model.trim="filter.keyword"
                    @click.native="handleClickFilterInput">
                </bk-input>
                <div ref="topoSearchResult" class="topo-search-result" v-show="filter.show">
                    <div class="search-result-head">
                        <span class="title">{{$t('搜索结果')}}</span>
                    </div>
                    <div class="search-result-body" v-if="filter.list.length">
                        <div class="search-result-item" v-for="(module, index) in filter.list" :key="index" @click="handleClickFilterItem(module)">
                            <div class="path-name">
                                <p class="name">{{module.bk_inst_name}}</p>
                                <p class="path" :title="module.path.join(' / ')">{{module.path.join(' / ')}}</p>
                            </div>
                            <div class="checkbox" @click.stop="handleCheckFilterItem(module)">
                                <bk-checkbox :value="filter.checked[module.bk_inst_id]"></bk-checkbox>
                            </div>
                        </div>
                    </div>
                    <div class="search-result-empty" v-else>
                        <span class="text">{{$t('暂无数据')}}</span>
                    </div>
                </div>
                <div class="tree-wrapper">
                    <bk-big-tree ref="tree" class="tree"
                        :selectable="true"
                        :options="{
                            idKey: getNodeId,
                            nameKey: 'bk_inst_name',
                            childrenKey: 'child'
                        }"
                        :node-height="36"
                        :before-select="beforeSelect"
                        @select-change="handleModuleSelectChange">
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
                            <span class="node-name" :title="node.name">{{node.name}}</span>
                        </div>
                    </bk-big-tree>
                </div>
            </div>
        </cmdb-resize-layout>
        <div class="table-wrapper" v-bkloading="{ isLoading: $loading(Object.values(request)) }">
            <host-table :list="hostList" :selected="selected" @select-change="handleHostSelectChange" />
        </div>
    </div>
</template>

<script>
    import { mapGetters } from 'vuex'
    import HostTable from './host-table.vue'
    import debounce from 'lodash.debounce'
    export default {
        components: {
            HostTable
        },
        props: {
            selected: {
                type: Array,
                default: () => ([])
            }
        },
        data () {
            return {
                hostList: [],
                topoModuleList: [],
                filter: {
                    show: false,
                    list: [],
                    keyword: '',
                    checked: {},
                    popover: null
                },
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
            ...mapGetters('businessHost', ['getDefaultSearchCondition'])
        },
        watch: {
            'filter.keyword' () {
                this.handleFilter()
            }
        },
        created () {
            this.initTopology()
            this.handleFilter = debounce(this.searchTopology, 300)
        },
        activated () {
            this.filter.keyword = ''
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
                    this.topoModuleList = this.getTopoModuleList(topology)
                    const defaultNodeId = this.getNodeId(topology[0])
                    this.$refs.tree.setExpanded(defaultNodeId)
                } catch (e) {
                    console.error(e)
                }
            },
            isModule (node) {
                return node.data.bk_obj_id === 'module'
            },
            async beforeSelect (node) {
                return this.isModule(node)
            },
            getInstanceTopology () {
                return this.$store.dispatch('objectMainLineModule/getInstTopoInstanceNum', {
                    bizId: this.bizId
                })
            },
            getInternalTopology () {
                return this.$store.dispatch('objectMainLineModule/getInternalTopo', {
                    bizId: this.bizId
                })
            },
            getNodeId (data) {
                return `${data.bk_obj_id}-${data.bk_inst_id}`
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
                        requestId: this.request.host
                    }
                })
            },
            searchTopology () {
                const keyword = this.filter.keyword.toLowerCase()
                if (!keyword) {
                    this.filter.list = []
                    this.filter.show = false
                    this.getFilterPopover().hide()
                    return
                }
                const result = this.topoModuleList.filter(mod => {
                    return mod.path.findIndex(path => path.toLowerCase().indexOf(keyword) !== -1) !== -1
                })
                this.filter.list = result
                this.showFilterPopover()
            },
            getTopoModuleList (treeData) {
                const modules = []
                const findModuleNode = function (data, parent) {
                    data.forEach(item => {
                        item.path = parent ? [...parent.path, item.bk_inst_name] : [item.bk_inst_name]
                        if (item.bk_obj_id === 'module') {
                            modules.push(item)
                        }
                        if (item.child) {
                            findModuleNode(item.child, item)
                        }
                    })
                }
                findModuleNode(treeData)

                return modules
            },
            getFilterPopover () {
                if (this.filter.popover) {
                    return this.filter.popover
                }
                this.filter.popover = this.$bkPopover(this.$refs.filterInput.$el, {
                    content: this.$refs.topoSearchResult,
                    allowHTML: true,
                    delay: 300,
                    trigger: 'manual',
                    boundary: 'window',
                    placement: 'bottom-start',
                    theme: 'light host-selector-toposearch-popover',
                    distance: 6,
                    interactive: true
                })
                return this.filter.popover
            },
            async handleModuleSelectChange (node) {
                const result = await this.searchHost(node)
                this.hostList = result.info
            },
            handleClickFilterInput () {
                if (!this.filter.keyword) {
                    return
                }
                this.showFilterPopover()
            },
            handleClickFilterItem (module) {
                const nodeId = this.getNodeId(module)
                this.$refs.tree.setSelected(nodeId, { emitEvent: true })
                this.$refs.tree.setExpanded(nodeId)
                this.hideFilterPopover()
            },
            async handleCheckFilterItem (module) {
                const moduleId = module.bk_inst_id
                const checked = !this.filter.checked[moduleId]
                this.$set(this.filter.checked, moduleId, checked)
                const result = await this.searchHost(this.$refs.tree.getNodeById(this.getNodeId(module)))
                if (checked) {
                    this.handleHostSelectChange({ removed: [], selected: result.info })
                } else {
                    this.handleHostSelectChange({ removed: result.info, selected: [] })
                }
                this.hideFilterPopover()
            },
            handleHostSelectChange (data) {
                this.$emit('select-change', data)
            },
            showFilterPopover () {
                this.filter.show = true
                this.$nextTick(() => {
                    this.getFilterPopover().show()
                })
            },
            hideFilterPopover () {
                this.filter.show = false
                this.getFilterPopover().hide()
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
            }
        }
    }
</script>

<style lang="scss" scoped>
    .host-selector-topology {
        display: flex;
        height: 100%;

        .tree-layout {
            width: 281px;
            height: 100%;
            border-right: 1px solid #dcdee5;
        }
        .topo-wrapper {
            margin-top: 24px;
            height: calc(100% - 24px);
        }
        .table-wrapper {
            flex: auto;
            margin-left: 20px;
            margin-top: 24px;
            width: 0; // 使table宽度自适应
        }
        .search-input {
            display: block;
            width: auto;
            margin-right: 20px;
        }
    }

    .tree-wrapper {
        height: calc(100% - 32px - 24px);
        margin: 12px 0;
        @include scrollbar;
    }
    .tree {
        padding: 0;
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

    .topo-search-result {
        width: 380px;
        background: #fff;
        padding-bottom: 20px;

        .search-result-head {
            display: flex;
            justify-content: space-between;
            align-items: center;
            height: 32px;
            line-height: 32px;
            padding: 0 20px;

            .title {
                font-size: 12px;
                color: #c4c6cc;
            }
            .check-all {
                /deep/ .bk-checkbox-text {
                    font-size: 12px;
                }
            }
        }
        .search-result-body {
            max-height: 300px;
            @include scrollbar-y;

            .search-result-item {
                display: flex;
                align-items: center;
                justify-content: space-between;
                height: 58px;
                padding: 0 20px;
                cursor: pointer;

                &:hover {
                    background: #e1ecff;
                }

                .path-name {
                    width: 240px;
                    .name {
                        font-weight: 700;
                        color: #63656e;
                        line-height: 16px;
                        overflow: hidden;
                        @include ellipsis;
                    }
                    .path {
                        color: #979ba5;
                        line-height: 16px;
                        overflow: hidden;
                        @include ellipsis;
                    }
                }
                .checkbox {
                    padding: 4px;
                }
            }
        }

        .search-result-empty {
            padding: 20px 0;
            text-align: center;
            .text {
                font-size: 12px;
                color: #63656e;
            }
        }
    }
</style>
<style lang="scss">
    .host-selector-toposearch-popover-theme {
        padding: 0;
    }
</style>
