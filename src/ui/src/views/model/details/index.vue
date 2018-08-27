<template>
    <div class="details-wrapper">
        <bk-tab :active-name.sync="curTabName">
            <bk-tabpanel name="host" :title="$t('ModelManagement[\'模型配置\']')">
                <v-field></v-field>
            </bk-tabpanel>
            <!-- <bk-tabpanel name="host" :title="$t('ModelManagement[\'字段分组\']')">
                <v-layout></v-layout>
            </bk-tabpanel>
            <bk-tabpanel name="host" :title="$t('ModelManagement[\'其他操作\']')">
                <v-other></v-other>
            </bk-tabpanel> -->
        </bk-tab>
    </div>
</template>

<script>
    import vField from './field'
    export default {
        components: {
            vField
        },
        data () {
            return {
                curTabName: ''
            }
        },
        methods: {
        }
    }
</script>

<style lang="scss" scoped>
    .details-wrapper{
        padding: 8px 20px 20px;
        height: 100%;
    }
</style>
