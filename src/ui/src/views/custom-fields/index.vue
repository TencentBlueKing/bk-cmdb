<template>
    <div class="custom-fields-layout" :style="{ padding: featureTips ? '15px 0 0 0' : 0 }">
        <cmdb-tips class="ml20 mr20 mb10" tips-key="showCustomFields" v-model="featureTips">{{$t('自定义字段功能提示')}}</cmdb-tips>
        <bk-tab class="tab-layout"
            :style="`--subHeight: ${featureTips ? '42px' : 0}`"
            :active.sync="active"
            type="unborder-card"
            @tab-change="handleTabChange">
            <bk-tab-panel v-for="model in mainLine"
                render-directive="if"
                :key="model.bk_obj_id"
                :name="model.bk_obj_id"
                :label="model.bk_obj_name">
                <field-group class="model-detail-wrapper"
                    :class="{
                        'has-tips': featureTips
                    }"
                    :custom-obj-id="model.bk_obj_id">
                </field-group>
            </bk-tab-panel>
        </bk-tab>
    </div>
</template>

<script>
    import fieldGroup from '@/components/model-manage/field-group'
    import RouterQuery from '@/router/query'
    export default {
        components: {
            fieldGroup
        },
        data () {
            return {
                active: RouterQuery.get('tab', 'set'),
                featureTips: true,
                mainLine: []
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
        async created () {
            try {
                const data = await this.getMainLine()
                this.mainLine = data.filter(model => ['host', 'set', 'module'].includes(model.bk_obj_id))
            } catch (e) {
                this.mainLine = []
            }
        },
        methods: {
            getMainLine () {
                return this.$store.dispatch('objectMainLineModule/searchMainlineObject', {
                    config: {
                        requestId: 'getMainLine'
                    }
                })
            },
            handleTabChange (modelId) {
                const activeModel = this.mainLine.find(model => model.bk_obj_id === modelId) || {}
                this.$store.commit('objectModel/setActiveModel', activeModel)
            }
        }
    }
</script>

<style lang="scss" scoped>
    .tab-layout {
        height: calc(100% - var(--subHeight));
        /deep/ {
            .bk-tab-content {
                padding-top: 10px;
            }
            .bk-tab-header {
                padding: 0;
                margin: 0 20px;
            }
        }
    }
    .model-detail-wrapper {
        padding: 0 !important;
        &.has-tips {
            height: calc(100% - 52px);
        }
    }
</style>

<style lang="scss">
    @import '@/assets/scss/model-manage.scss';
</style>
