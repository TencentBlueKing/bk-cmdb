<template>
    <div class="transfer-layout clearfix">
        <div class="columns-layout fl">
            <div class="business-layout">
                <label class="business-label">{{$t('Common[\'业务\']')}}</label>
                <cmdb-business-selector class="business-selector" v-model="business" :disabled="true">
                </cmdb-business-selector>
            </div>
            <div class="tree-layout">
                <cmdb-tree ref="topoTree" class="topo-tree"
                    id-key="bk_inst_id"
                    label-key="bk_inst_name"
                    children-key="child"
                    :id-generator="getTopoNodeId"
                    :tree="tree.data"
                    :before-select="beforeNodeSelect"
                    @on-selected="handleNodeSelected">
                    <div class="tree-node clearfix" slot-scope="{node, state}" :class="{'tree-node-selected': state.selected}">
                        <template v-if="[1, 2].includes(node.default)">
                            <i class='topo-node-icon topo-node-icon-internal icon-cc-host-free-pool' v-if="node.default === 1"></i>
                            <i class='topo-node-icon topo-node-icon-internal icon-cc-host-breakdown' v-else></i>
                        </template>
                        <i class="topo-node-icon topo-node-icon-text" v-else>{{node['bk_obj_name'][0]}}</i>
                        <span class="topo-node-text">{{node['bk_inst_name']}}</span>
                    </div>
                </cmdb-tree>
            </div>
        </div>
        <div class="columns-layout fl">
        </div>
    </div>
</template>

<script>
    import { mapGetters, mapActions } from 'vuex'
    export default {
        data () {
            return {
                business: '',
                topoModel: [],
                tree: {
                    data: []
                },
                selectedModuleNodes: []
            }
        },
        watch: {
            async business (business) {
                if (business) {
                    await this.getMainlineModel()
                    await this.getBusinessTopo()
                }
            }
        },
        methods: {
            ...mapActions('objectMainLineModule', [
                'searchMainlineObject',
                'getInstTopo',
                'getInternalTopo'
            ]),
            getMainlineModel () {
                return this.searchMainlineObject({fromCache: true}).then(topoModel => {
                    this.topoModel = topoModel
                    return topoModel
                })
            },
            getBusinessTopo () {
                return Promise.all([
                    this.getInstTopo({
                        bizId: this.business,
                        config: {
                            requestId: 'getInstTopo',
                            fromCache: true
                        }
                    }),
                    this.getInternalTopo({
                        bizId: this.business,
                        config: {
                            requestId: 'getInternalTopo',
                            fromCache: true
                        }
                    })
                ]).then(([instTopo, internalTopo]) => {
                    const moduleModel = this.getModelByObjId('module')
                    const internalModule = internalTopo.module.map(module => {
                        return {
                            'default': ['空闲机', 'idle machine'].includes(module['bk_module_name']) ? 1 : 2,
                            'bk_obj_id': 'module',
                            'bk_obj_name': moduleModel['bk_obj_name'],
                            'bk_inst_id': module['bk_module_id'],
                            'bk_inst_name': module['bk_module_name']
                        }
                    })
                    this.tree.data = [{
                        'default': 0,
                        'bk_obj_id': 'module',
                        'bk_obj_name': this.$t('HostResourcePool["资源池"]'),
                        'bk_inst_id': 'source',
                        'bk_inst_name': this.$t('HostResourcePool["资源池"]'),
                        'child': []
                    }, {
                        expanded: true,
                        ...instTopo[0],
                        child: [...internalModule, ...instTopo[0].child]
                    }]
                })
            },
            getModelByObjId (id) {
                return this.topoModel.find(model => model['bk_obj_id'] === id)
            },
            getTopoNodeId (node) {
                return `${node['bk_obj_id']}-${node['bk_inst_id']}`
            },
            beforeNodeSelect (node, state) {
                let confirmResolver
                let confirmRejecter
                const asyncConfirm = new Promise((resolve, reject) => {
                    confirmResolver = resolve
                    confirmRejecter = reject
                })
                if (node['bk_obj_id'] !== 'module') {
                    confirmResolver(true)
                } else {
                    const isSpecialNode = !!node.default || node['bk_inst_id'] === 'source'
                    const hasNormalNode = this.selectedModuleNodes.some(node => {
                        return !node.default && node['bk_inst_id'] !== 'source'
                    })
                    const hasSpecialNode = this.selectedModuleNodes.some(node => {
                        return node.default || node['bk_inst_id'] === 'source'
                    })
                    if (isSpecialNode && hasNormalNode) {
                        this.$bkInfo({
                            title: this.$t('Common[\'转移确认\']', {target: node['bk_inst_name']}),
                            confirmFn: () => {
                                confirmResolver(true)
                            },
                            cancelFn: () => {
                                confirmRejecter(false)
                            }
                        })
                    } else {
                        if (hasSpecialNode) {
                            this.selectedModuleNodes = []
                        }
                        confirmResolver(true)
                    }
                }
                return asyncConfirm
            },
            handleNodeSelected (node, state) {
                if (node['bk_obj_id'] === 'module') {
                    this.selectedModuleNodes.push(node)
                }
            }
        }
    }
</script>

<style lang="scss" scoped>
    .transfer-layout {
        height: 540px;
        width: 720px;
        .columns-layout{
            width: 50%;
            height: 100%;
        }
    }
    .business-layout {
        border-right: 1px solid $cmdbBorderColor;
        height: 65px;
        &:before {
            display: inline-block;
            width: 0;
            height: 100%;
            content: '';
            font-size: 0;
            vertical-align: middle;
        }
        .business-label {
            display: inline-block;
            vertical-align: middle;
            padding: 0 25px;
        }
        .business-selector {
            display: inline-block;
            vertical-align: middle;

            width: 245px;
        }
    }
    .tree-layout {
        height: 415px;
        border-top: 1px solid $cmdbBorderColor;
        border-right: 1px solid $cmdbBorderColor;
    }
    .topo-tree{
        padding: 0 0 0 20px;
        height: 100%;
        @include scrollbar-y;
        .tree-node {
            font-size: 0;
            &:hover{
                .topo-node-icon.topo-node-icon-text{
                    background-color: #50abff;
                }
                .topo-node-icon.topo-node-icon-internal{
                    color: #50abff;
                }
            }
            &.tree-node-selected{
                .topo-node-icon.topo-node-icon-text{
                    background-color: #498fe0;
                }
                .topo-node-icon.topo-node-icon-internal{
                    color: #ffb400;
                }
            }
        }
        .topo-node-icon{
            display: inline-block;
            vertical-align: middle;
            width: 16px;
            height: 16px;
            line-height: 16px;
            font-size: 12px;
            text-align: center;
            color: #fff;
            font-style: normal;
            background-color: #c3cdd7;
            &.topo-node-icon-internal{
                font-size: 16px;
                color: $cmdbTextColor;
                background-color: transparent;
            }
        }
        .topo-node-text{
            display: inline-block;
            vertical-align: middle;
            padding: 0 0 0 8px;
            font-size: 14px;
        }
    }
</style>