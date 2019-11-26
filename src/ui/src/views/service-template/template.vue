<template>
    <div class="template-layout">
        <bk-tab class="template-tab"
            type="unborder-card"
            :show-header="isUpdate"
            :active.sync="active">
            <bk-tab-panel :label="$t('服务模板配置')" render-directive="if" name="config">
                <service-template-config :style="{ padding: isUpdate ? '20px 18px' : '0 20px' }"></service-template-config>
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
    export default {
        components: {
            ServiceTemplateConfig,
            ServiceTemplateInstance
        },
        data () {
            return {
                active: 'config'
            }
        },
        computed: {
            isUpdate () {
                return this.$route.params.templateId !== undefined
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
        /deep/ .bk-tab-header {
            padding: 0;
            margin: 0 20px;
        }
    }
</style>
