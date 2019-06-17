<template>
    <div class="layout">
        <cmdb-resize-layout class="tree-layout fl"
            direction="right"
            :handler-offset="3"
            :min="200"
            :max="480">
            <cmdb-topology-tree></cmdb-topology-tree>
        </cmdb-resize-layout>
        <div class="tab-layout">
            <bk-tab :active-name="active">
                <bk-tabpanel name="serviceInstances"
                    :title="$t('BusinessTopology[\'服务实例\']')"
                    :show="isModuleNode">
                    <cmdb-service-instances></cmdb-service-instances>
                </bk-tabpanel>
                <bk-tabpanel name="nodeInfo" :title="$t('BusinessTopology[\'节点信息\']')">
                    <cmdb-service-node-info></cmdb-service-node-info>
                </bk-tabpanel>
                <div slot="setting" v-show="withTemplate">
                    <a href="javascript:void(0)" class="text-primary sync-template-link" @click="handleSyncTemplate">
                        <i class="bk-icon icon-refresh"></i>
                        {{$t('BusinessTopology["同步模板"]')}}
                    </a>
                </div>
            </bk-tab>
        </div>
    </div>
</template>

<script>
    import cmdbTopologyTree from './children/topology-tree.vue'
    import cmdbServiceInstances from './children/service-instances.vue'
    import cmdbServiceNodeInfo from './children/service-node-info.vue'
    export default {
        components: {
            cmdbTopologyTree,
            cmdbServiceInstances,
            cmdbServiceNodeInfo
        },
        data () {
            return {
                active: 'nodeInfo'
            }
        },
        computed: {
            selectedNode () {
                return this.$store.state.businessTopology.selectedNode
            },
            isModuleNode () {
                return this.selectedNode && this.selectedNode.data.bk_obj_id === 'module'
            },
            withTemplate () {
                const nodeInstance = this.$store.state.businessTopology.selectedNodeInstance
                return this.isModuleNode && nodeInstance && nodeInstance.service_template_id
            }
        },
        watch: {
            isModuleNode (isModuleNode) {
                if (isModuleNode) {
                    this.active = 'serviceInstances'
                } else {
                    this.active = 'nodeInfo'
                }
            }
        },
        beforeDestroy () {
            this.$store.commit('businessTopology/resetState')
        },
        methods: {
            handleSyncTemplate () {
                this.$router.push({
                    name: 'synchronous',
                    params: {
                        moduleId: this.selectedNode.data.bk_inst_id,
                        setId: this.selectedNode.parent.data.bk_inst_id
                    },
                    query: {
                        path: [...this.selectedNode.parents, this.selectedNode].map(node => node.name).join(' / '),
                        from: `${this.$route.path}?module=${this.selectedNode.data.bk_inst_id}`
                    }
                })
            }
        }
    }
</script>

<style lang="scss" scoped>
    .layout {
        height: 100%;
        padding: 0;
    }
    .tree-layout {
        width: 280px;
        height: 100%;
        padding: 20px 0 20px 20px;
        border-right: 1px solid $cmdbBorderColor;
    }
    .tab-layout {
        height: 100%;
        overflow: hidden;
    }
    .sync-template-link {
        margin: 0 -20px 0 0;
        font-size: 14px;
        .icon-refresh {
            @include inlineBlock;
            font-size: 12px;
        }
    }
</style>
