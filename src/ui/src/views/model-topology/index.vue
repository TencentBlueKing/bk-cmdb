<template>
    <div class="topology-layout clearfix">
        <cmdb-graphics-toolbar class="toolbar"></cmdb-graphics-toolbar>
        <cmdb-graphics-models class="models fl" ref="models"></cmdb-graphics-models>
        <cmdb-topology-graphics class="graphics" ref="graphics"></cmdb-topology-graphics>
        <cmdb-graphics-config ref="config"></cmdb-graphics-config>
        <cmdb-graphics-association ref="association"></cmdb-graphics-association>
    </div>
</template>
<script>
    import cmdbGraphicsToolbar from './children/graphics-toolbar'
    import cmdbTopologyGraphics from './children/topology-graphics'
    import cmdbGraphicsConfig from './children/graphics-config'
    import cmdbGraphicsModels from './children/graphics-models'
    import cmdbGraphicsAssociation from './children/graphics-association'
    export default {
        components: {
            cmdbGraphicsToolbar,
            cmdbGraphicsConfig,
            cmdbTopologyGraphics,
            cmdbGraphicsModels,
            cmdbGraphicsAssociation
        },
        provide () {
            return {
                parentRefs: this.$refs
            }
        }
    }
</script>

<style lang="scss" scoped>
    .topology-layout {
        height: 100%;
        padding: 0;
        .toolbar {
            height: 50px;
        }
        .models {
            height: calc(100% - 50px);
            @include scrollbar-y;
        }
        .graphics {
            overflow: hidden;
            height: calc(100% - 50px);
        }
    }
</style>
