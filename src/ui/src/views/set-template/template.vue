<template>
    <div class="template-info-layout" :style="{ padding: mode === 'view' ? '0' : '15px 0 0 0' }">
        <template v-if="mode === 'view'">
            <bk-tab class="info-tab"
                type="unborder-card"
                ref="infoTab"
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
    import RouterQuery from '@/router/query'
    export default {
        components: {
            setTemplateInstance,
            setTemplateConfig
        },
        data () {
            return {
                active: RouterQuery.get('tab', 'setting')
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
            },
            active: {
                immediate: true,
                handler (active) {
                    if (this.mode === 'view') {
                        this.checkSyncStatus()
                    }
                    RouterQuery.set({
                        tab: active
                    })
                }
            }
        },
        methods: {
            async checkSyncStatus () {
                try {
                    const data = await this.$store.dispatch('setTemplate/getSetTemplateStatus', {
                        bizId: this.$store.getters['objectBiz/bizId'],
                        params: {
                            set_template_ids: [Number(this.templateId)]
                        }
                    })
                    const needSync = this.$tools.getValue(data, '0.need_sync')
                    const tabHeader = this.$refs.infoTab.$el.querySelector('.bk-tab-label-item.is-last')
                    if (needSync) {
                        tabHeader.classList.add('has-tips')
                    } else {
                        tabHeader.classList.remove('has-tips')
                    }
                } catch (e) {
                    console.error(e)
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
        /deep/ .bk-tab-label-item.has-tips:before {
            content: "";
            position: absolute;
            top: 18px;
            right: 12px;
            width: 6px;
            height: 6px;
            border-radius: 50%;
            background-color: $dangerColor;
        }
    }
    .template-config-view {
        padding: 24px 0 0 0;
    }
</style>
