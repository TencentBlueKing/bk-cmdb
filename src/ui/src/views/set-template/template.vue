<template>
    <div class="template-info-layout">
        <template v-if="mode === 'view'">
            <bk-tab class="info-tab"
                type="unborder-card"
                :active.sync="active">
                <bk-tab-panel name="setting" :label="$t('集群模板配置')">
                    <set-template-config class="template-config-view"></set-template-config>
                </bk-tab-panel>
                <bk-tab-panel name="instance" render-directive="if" :label="$t('集群模板实例')">
                    <set-template-instance :template-id="templateId"></set-template-instance>
                </bk-tab-panel>
            </bk-tab>
        </template>
        <template v-else>
            <set-template-config></set-template-config>
        </template>
    </div>
</template>

<script>
    import setTemplateInstance from './children/template-intance.vue'
    import setTemplateConfig from './children/template-config.vue'
    export default {
        components: {
            setTemplateInstance,
            setTemplateConfig
        },
        data () {
            return {
                active: 'setting'
            }
        },
        computed: {
            mode () {
                return this.$route.params.mode
            },
            templateId () {
                return this.$route.params.templateId
            }
        },
        watch: {
            $route: {
                immediate: true,
                handler ($route) {
                    if ($route.query.tab) {
                        this.active = $route.query.tab
                    }
                }
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
    .template-config-view {
        padding: 24px 0 0 0;
    }
</style>
