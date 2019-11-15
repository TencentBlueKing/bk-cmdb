<template>
    <div class="conflict-list">
        <feature-tips
            :feature-name="'hostApply'"
            :show-tips="showFeatureTips"
            :desc="$t('主机属于多个模块，且同一属性的自动应用配置有差异，若不处理，将维持主机转移前的该项属性值')"
            @close-tips="showFeatureTips = false">
        </feature-tips>
        <property-confirm-table :confirm-data="confirmData"></property-confirm-table>
        <div class="bottom-actionbar">
            <div class="actionbar-inner">
                <bk-button theme="primary" :disabled="applyButtonDisabled">应用</bk-button>
                <bk-button theme="default">返回</bk-button>
            </div>
        </div>
    </div>
</template>

<script>
    import { mapGetters, mapActions } from 'vuex'
    import featureTips from '@/components/feature-tips/index'
    import propertyConfirmTable from './children/property-confirm-table'
    import { MENU_BUSINESS_HOST_APPLY } from '@/dictionary/menu-symbol'
    export default {
        components: {
            featureTips,
            propertyConfirmTable
        },
        data () {
            return {
                showFeatureTips: false,
                applyButtonDisabled: true,
                confirmData: [
                    {
                        id: 1,
                        bk_host_innerip: '120.23.3.3',
                        bk_cloud_id: '直连区域',
                        bk_asset_id: 'No90378',
                        bk_host_name: 'nginx1_lol',
                        diff_value: 'xxx'
                    }
                ]
            }
        },
        computed: {
            ...mapGetters(['featureTipsParams', 'supplierAccount'])
        },
        watch: {

        },
        created () {
            this.showFeatureTips = this.featureTipsParams['hostApplyConflict']
            this.setBreadcrumbs()
        },
        methods: {
            ...mapActions('objectModelProperty', [
                'searchObjectAttribute'
            ]),
            setBreadcrumbs () {
                this.$store.commit('setTitle', this.$t('冲突列表'))
                this.$store.commit('setBreadcrumbs', [{
                    label: this.$t('主机属性自动应用'),
                    route: {
                        name: MENU_BUSINESS_HOST_APPLY
                    }
                }, {
                    label: this.$t('冲突列表')
                }])
            }
        }
    }
</script>

<style lang="scss" scoped>
    .conflict-list {
        padding: 0 20px;
    }

    .bottom-actionbar {
        position: absolute;
        width: 100%;
        height: 50px;
        border-top: 1px solid #dcdee5;
        bottom: 0;
        left: 0;

        .actionbar-inner {
            padding: 8px 0 0 20px;

            .bk-button {
                min-width: 86px;
            }
        }
    }
</style>
