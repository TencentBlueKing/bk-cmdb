<template>
    <div class="custom-fields-layout">
        <feature-tips class="ml20 mr20"
            :feature-name="'customFields'"
            :show-tips="showFeatureTips"
            :desc="$t('自定义字段功能提示')"
            @close-tips="showFeatureTips = false">
        </feature-tips>
        <bk-tab class="tab-layout"
            :style="`--subHeight: ${showFeatureTips ? '42px' : 0}`"
            type="unborder-card"
            @tab-change="handleTabChange">
            <bk-tab-panel v-for="model in mainLine"
                render-directive="if"
                :key="model.bk_obj_id"
                :name="model.bk_obj_id"
                :label="model.bk_obj_name">
                <field-group class="model-detail-wrapper"
                    :class="{
                        'has-tips': showFeatureTips
                    }"
                    :custom-obj-id="model.bk_obj_id">
                </field-group>
            </bk-tab-panel>
        </bk-tab>
    </div>
</template>

<script>
    import featureTips from '@/components/feature-tips/index'
    import fieldGroup from '@/components/model-manage/field-group'
    import { mapGetters } from 'vuex'
    export default {
        components: {
            fieldGroup,
            featureTips
        },
        data () {
            return {
                mainLine: [],
                showFeatureTips: false
            }
        },
        computed: {
            ...mapGetters(['featureTipsParams'])
        },
        async created () {
            try {
                this.showFeatureTips = this.featureTipsParams['customFields']
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
    .custom-fields-layout {
        padding: 0;
    }
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
