<template>
    <div class="layout">
        <cmdb-resize-layout class="tree-layout fl"
            direction="right"
            :handler-offset="3"
            :min="200"
            :max="480">
            <cmdb-topology-tree></cmdb-topology-tree>
        </cmdb-resize-layout>
        <div class="tab-layout"></div>
    </div>
</template>

<script>
    import cmdbTopologyTree from './children/topology-tree.vue'
    export default {
        components: {
            cmdbTopologyTree
        },
        data () {
            return {}
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
