<template>
    <div class="topology-layout clearfix">
        <cmdb-resize-layout class="tree-layout fl" direction="right" :min="200" :max="480">
            <cmdb-business-selector class="business-selector" v-model="business">
            </cmdb-business-selector>
            <cmdb-tree ref="topoTree" class="topo-tree"
                id-key="bk_inst_id"
                label-key="bk_inst_name"
                children-key="child"
                :id-generator="getTopoNodeId"
                :tree="tree.data">
                <div class="tree-node" slot-scope="{node}">
                    <template v-if="[1, 2].includes(node.default)">
                        <i class='topo-node-icon topo-node-icon-internal icon-cc-host-free-pool' v-if="node.default === 1"></i>
                        <i class='topo-node-icon topo-node-icon-internal icon-cc-host-breakdown' v-else></i>
                    </template>
                    <i class="topo-node-icon topo-node-icon-text" v-else>{{node['bk_obj_name'][0]}}</i>
                    <span class="topo-node-text">{{node['bk_inst_name']}}</span>
                </div>
            </cmdb-tree>
        </cmdb-resize-layout>
        <div class="hosts-layout"></div>
    </div>
</template>

<script>
    import { mapActions } from 'vuex'
    export default {
        data () {
            return {
                business: '',
                businessResolver: null,
                businessTopo: [],
                topoModel: [],
                tree: {
                    data: [],
                    internalIdleId: null,
                    internalFaultId: null
                }
            }
        },
        watch: {
            business (business) {
                if (this.businessResolver) {
                    this.businessResolver()
                } else {
                    this.getBusinessTopo()
                    this.getHostList()
                }
            }
        },
        async created () {
            try {
                await this.getBusiness()
                await this.getMainlineModel()
                await this.getBusinessTopo()
            } catch (e) {
                console.log(e)
            }
        },
        methods: {
            ...mapActions('objectMainLineModule', [
                'searchMainlineObject',
                'getInstTopo',
                'getInternalTopo'
            ]),
            getBusiness () {
                return new Promise((resolve, reject) => {
                    this.businessResolver = () => {
                        this.businessResolver = null
                        resolve()
                    }
                })
            },
            getMainlineModel () {
                return this.searchMainlineObject({fromCache: true}).then(topoModel => {
                    this.topoModel = topoModel
                    return topoModel
                })
            },
            getBusinessTopo () {
                return Promise.all([
                    this.getInstTopo({
                        bizId: this.business
                    }),
                    this.getInternalTopo({
                        bizId: this.business
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
                    internalModule.forEach(node => {
                        if (node.default === 1) {
                            this.tree.internalIdleId = node['bk_inst_id']
                        } else {
                            this.tree.internalFaultId = node['bk_inst_id']
                        }
                    })
                    instTopo[0] = {
                        selected: true,
                        expanded: true,
                        ...instTopo[0],
                        child: [...internalModule, ...instTopo[0].child]
                    }
                    this.tree.data = instTopo
                })
            },
            getModelByObjId (id) {
                return this.topoModel.find(model => model['bk_obj_id'] === id)
            },
            getTopoNodeId (node) {
                return `${node['bk_obj_id']}-${node['bk_inst_id']}`
            }
        }
    }
</script>

<style lang="scss" scoped>
    .topology-layout{
        padding: 0;
        height: 100%;
    }
    .tree-layout{
        width: 280px;
        height: 100%;
        border-right: 1px solid $cmdbBorderColor;
        background-color: #fafbfd;
        .business-selector{
            display: block;
            width: auto;
            margin: 20px;
        }
        .topo-tree{
            padding: 0 0 0 20px;
            height: calc(100% - 76px);
            @include scrollbar-y;
            .tree-node {
                font-size: 0;
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
    }
    .hosts-layout{
        overflow: hidden;
        height: 100%;
    }
</style>