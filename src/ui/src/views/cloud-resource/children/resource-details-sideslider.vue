<template>
    <bk-sideslider v-transfer-dom
        :is-show.sync="isShow"
        :title="title"
        :width="695"
        @hidden="handleHidden">
        <template slot="content">
            <bk-tab class="details-tab" :active.sync="active" type="unborder-card" slot="content">
                <bk-tab-panel name="details" :label="$t('任务详情')">
                </bk-tab-panel>
                <bk-tab-panel name="history" :label="$t('录入历史')">
                </bk-tab-panel>
            </bk-tab>
            <keep-alive>
                <component class="details-component" :is="component"
                    :container="this"
                    v-bind="componentProps">
                </component>
            </keep-alive>
        </template>
    </bk-sideslider>
</template>

<script>
    import CloudResourceDetailsInfo from './resource-details-info.vue'
    import CloudResourceDetailsHistory from './resource-details-history.vue'
    import CloudResourceForm from './resource-form.vue'
    export default {
        name: 'cloud-resource-details',
        components: {
            [CloudResourceDetailsHistory.name]: CloudResourceDetailsHistory,
            [CloudResourceDetailsInfo.name]: CloudResourceDetailsInfo,
            [CloudResourceForm.name]: CloudResourceForm
        },
        data () {
            return {
                title: '',
                isShow: false,
                active: 'details',
                detailsComponent: null,
                componentProps: {}
            }
        },
        computed: {
            component () {
                if (this.active === 'details') {
                    return this.detailsComponent
                }
                return CloudResourceDetailsHistory.name
            }
        },
        methods: {
            show (options) {
                this.title = options.title || this.title
                this.componentProps = options.props
                this.detailsComponent = options.detailsComponent || CloudResourceDetailsInfo.name
                this.isShow = true
            },
            hide (eventType) {
                this.isShow = false
                eventType && this.$emit(eventType)
            },
            handleHidden () {
                this.active = 'details'
                this.detailsComponent = null
                this.componentProps = {}
            }
        }
    }
</script>

<style lang="scss" scoped>
    .details-tab {
        height: auto;
        /deep/ {
            .bk-tab-header {
                padding: 0;
                margin: 0 30px;
            }
            .bk-tab-section {
                height: 0;
            }
        }
    }
    .details-component {
        height: calc(100% - 58px);
        @include scrollbar-y;
    }
</style>
