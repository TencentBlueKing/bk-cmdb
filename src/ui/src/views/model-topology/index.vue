<template>
    <div class="topology-layout">
        <cmdb-graphics-toolbar class="toolbar"></cmdb-graphics-toolbar>
        <cmdb-topology-graphics class="graphics" ref="graphics">
        </cmdb-topology-graphics>
    </div>
</template>
<script>
    import cmdbGraphicsToolbar from './children/graphics-toolbar'
    import cmdbTopologyGraphics from './children/topology-graphics'
    export default {
        components: {
            cmdbGraphicsToolbar,
            cmdbTopologyGraphics
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
        .graphics {
            height: calc(100% - 50px);
        }
    }
</style>