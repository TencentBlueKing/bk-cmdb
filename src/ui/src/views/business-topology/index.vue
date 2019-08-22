<template>
    <div class="info-layout">
        <cmdb-resize-layout class="tree-layout fl"
            direction="right"
            :handler-offset="3"
            :min="200"
            :max="480">
            <cmdb-topology-tree></cmdb-topology-tree>
        </cmdb-resize-layout>
        <div class="tab-layout">
            <bk-tab :active="active" type="unborder-card">
                <bk-tab-panel name="serviceInstances"
                    :label="$t('服务实例')"
                    :visible="isModuleNode">
                    <cmdb-service-instances></cmdb-service-instances>
                </bk-tab-panel>
                <bk-tab-panel name="nodeInfo" :label="$t('节点信息')">
                    <cmdb-service-node-info></cmdb-service-node-info>
                </bk-tab-panel>
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
        }
    }
</script>

<style lang="scss" scoped>
    .info-layout {
        min-width: 1200px;
        height: 100%;
        padding: 0;
    }
    .tree-layout {
        width: 280px;
        height: 100%;
        padding: 20px 0;
        border-right: 1px solid #dcdee5;
    }
    .tab-layout {
        height: 100%;
        overflow: hidden;
        .bk-tab {
            height: 100%;
        }
        .layout {
            height: 100%;
            padding: 0;
        }
        /deep/ .bk-tab-section {
            height: calc(100% - 42px);
            padding-top: 14px;
        }
        /deep/ .bk-tab-content {
            height: 100%;
        }
    }
</style>
