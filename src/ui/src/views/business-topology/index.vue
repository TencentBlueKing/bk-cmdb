<template>
    <div class="info-layout">
        <cmdb-resize-layout class="tree-layout fl"
            direction="right"
            :handler-offset="3"
            :min="200"
            :max="480">
            <cmdb-topology-tree ref="topologyTree"></cmdb-topology-tree>
        </cmdb-resize-layout>
        <div class="tab-layout">
            <bk-tab :active="active" type="unborder-card" v-if="!isBusinessNode">
                <bk-tab-panel name="serviceInstances"
                    :label="$t('服务实例')"
                    :visible="isModuleNode">
                    <cmdb-service-instances></cmdb-service-instances>
                </bk-tab-panel>
                <bk-tab-panel name="nodeInfo" :label="$t('节点信息')">
                    <cmdb-service-node-info></cmdb-service-node-info>
                </bk-tab-panel>
            </bk-tab>
            <div class="business-node-view" v-else>
                <img class="node-view-img" src="../../assets/images/add-node.png" width="103"
                    :style="{
                        'margin-top': ($APP.height - 120) * 0.2 + 'px'
                    }">
                <i18n class="node-view-handler"
                    v-if="!selectedNode.children.length"
                    path="未添加节点提示">
                    <a class="node-view-link" href="javascript:void(0)" place="link"
                        @click="handleAddNode">
                        {{$t('添加节点')}}
                    </a>
                </i18n>
                <i18n class="node-view-handler"
                    v-else
                    path="已添加节点提示">
                    <a class="node-view-link" href="javascript:void(0)" place="link" style="margin-left: -2px"
                        @click="handleAddNode">
                        {{$t('添加节点')}}
                    </a>
                </i18n>
            </div>
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
            isBusinessNode () {
                return this.selectedNode && this.selectedNode.data.bk_obj_id === 'biz'
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
            handleAddNode () {
                this.$refs.topologyTree.showCreateDialog(this.selectedNode)
            }
        }
    }
</script>

<style lang="scss" scoped>
    .info-layout {
        padding: 0;
        border-top: 1px solid $cmdbLayoutBorderColor;
    }
    .tree-layout {
        width: 280px;
        height: 100%;
        padding: 10px 0;
        border-right: 1px solid #dcdee5;
    }
    .tab-layout {
        height: 100%;
        @include scrollbar-x;
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
            min-width: 826px;
        }
        /deep/ .bk-tab-content {
            height: 100%;
        }
    }
    .business-node-view {
        text-align: center;
        font-size: 14px;
        .node-view-img {
            display: block;
            margin: 0 auto 5px;
        }
        .node-view-link {
            color: #3A84FF;
        }
    }
</style>
