<template>
    <div class="template-layout">
        <bk-tab class="template-tab"
            type="unborder-card"
            :class="{
                'no-header': !isUpdate
            }"
            :show-header="isUpdate"
            :active.sync="active">
            <bk-tab-panel :label="$t('服务模板配置')" name="config">
                <service-template-config></service-template-config>
            </bk-tab-panel>
            <bk-tab-panel :label="$t('服务模板实例')" name="instance" v-if="isUpdate">
                <service-template-instance :active="active === 'instance'"></service-template-instance>
            </bk-tab-panel>
        </bk-tab>
    </div>
</template>

<script>
    import ServiceTemplateConfig from './children/operational'
    import ServiceTemplateInstance from './children/template-instance'
    import Bus from '@/utils/bus'
    import RouterQuery from '@/router/query'
    export default {
        components: {
            ServiceTemplateConfig,
            ServiceTemplateInstance
        },
        data () {
            return {
                active: RouterQuery.get('tab', 'config')
            }
        },
        computed: {
            isUpdate () {
                return this.$route.params.templateId !== undefined
            }
        },
        watch: {
            active: {
                immediate: true,
                handler (active) {
                    RouterQuery.set({
                        tab: active
                    })
                }
            }
        },
        created () {
            Bus.$on('active-change', active => {
                this.active = active
            })
        },
        beforeDestroy () {
            Bus.$off('active-change')
        }
    }
</script>

<style lang="scss" scoped>
    .template-tab {
        /deep/ {
            .bk-tab-header {
                padding: 0;
                margin: 0 20px;
            }
            .bk-tab-section {
                padding: 0;
            }
        }
        &.no-header {
            /deep/ .bk-tab-section {
                height: 100%;
            }
        }
    }
</style>
