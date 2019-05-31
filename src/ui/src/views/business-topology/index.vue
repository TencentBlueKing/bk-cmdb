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
                <bk-tabpanel name="serviceInstances" :title="$t('BusinessTopology[\'服务实例\']')">
                    <cmdb-service-instances></cmdb-service-instances>
                </bk-tabpanel>
                <bk-tabpanel name="nodeInfo" :title="$t('BusinessTopology[\'节点信息\']')"></bk-tabpanel>
            </bk-tab>
        </div>
    </div>
</template>

<script>
    import cmdbTopologyTree from './children/topology-tree.vue'
    import cmdbServiceInstances from './children/service-instances.vue'
    export default {
        components: {
            cmdbTopologyTree,
            cmdbServiceInstances
        },
        data () {
            return {
                active: 'serviceInstances'
            }
        },
        beforeDestroy () {
            this.$store.commit('businessTopology/resetProperties')
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
</style>
