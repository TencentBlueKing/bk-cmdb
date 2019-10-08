<template>
    <div class="template-info-layout">
        <bk-tab class="info-tab"
            type="unborder-card"
            :active.sync="active">
            <bk-tab-panel name="setting" :label="$t('模板配置')"></bk-tab-panel>
            <bk-tab-panel name="instance" :label="$t('模板实例')">
                <template-instance :template-id="templateId"></template-instance>
            </bk-tab-panel>
        </bk-tab>
    </div>
</template>

<script>
    import templateInstance from './children/template-intance.vue'
    export default {
        components: {
            templateInstance
        },
        data () {
            return {
                active: this.$route.query.tab || 'setting'
            }
        },
        computed: {
            templateId () {
                return this.$route.params.templateId || null
            }
        }
    }
</script>

<style lang="scss" scoped>
    .info-tab {
        /deep/ .bk-tab-header {
            padding: 0;
            margin: 0 20px;
        }
    }
</style>
