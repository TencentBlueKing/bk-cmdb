<template>
    <div class="custom-fields-layout">
        <bk-tab class="tab-layout"
            type="unborder-card"
            @tab-change="handleTabChange">
            <bk-tab-panel v-for="model in mainLine"
                render-directive="if"
                :key="model.bk_obj_id"
                :name="model.bk_obj_id"
                :label="model.bk_obj_name">
                <field-group class="model-detail-wrapper"
                    :custom-obj-id="model.bk_obj_id">
                </field-group>
            </bk-tab-panel>
        </bk-tab>
    </div>
</template>

<script>
    import fieldGroup from '@/components/model-manage/field-group'
    export default {
        components: {
            fieldGroup
        },
        data () {
            return {
                mainLine: []
            }
        },
        async created () {
            this.mainLine = await this.getMainLine()
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
        height: 100%;
    }
    .tab-layout {
        /deep/ .bk-tab-header {
            padding: 0;
            margin: 0 20px;
        }
    }
</style>

<style lang="scss">
    @import '@/assets/scss/model-manage.scss';
</style>
