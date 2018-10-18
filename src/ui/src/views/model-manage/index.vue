<template>
    <bk-tab :active-name.sync="activeTabName" class="tab-wrapper model-manage-wrapper">
        <bk-tabpanel name="group" :title="$t('ModelManagement[\'模型分组\']')">
            <v-group></v-group>
        </bk-tabpanel>
        <bk-tabpanel name="topo" :title="$t('ModelManagement[\'模型拓扑\']')">
            <v-topo></v-topo>
        </bk-tabpanel>
        <bk-tabpanel name="relation" :title="$t('ModelManagement[\'关联类型\']')">
            <v-relation></v-relation>
        </bk-tabpanel>
    </bk-tab>
</template>

<script>
    import vGroup from './group'
    import vTopo from './topo'
    import vRelation from './relation'
    export default {
        components: {
            vGroup,
            vTopo,
            vRelation
        },
        data () {
            return {
                activeTabName: 'group'
            }
        }
    }
</script>

<style lang="scss" scoped>
    .tab-wrapper {
        padding-top: 0;
    }
</style>

<style lang="scss">
    @import '@/assets/scss/model-manage.scss';
</style>
