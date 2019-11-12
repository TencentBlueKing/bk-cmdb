<template>
    <div class="host-apply-confirm-wrapper">
        <feature-tips
            :feature-name="'hostApply'"
            :show-tips="showFeatureTips"
            :desc="$t('冲突的属性若不解决，在应用时将被忽略，确定应用后，新转入该模块的主机将自动应用配置的属性')"
            @close-tips="showFeatureTips = false">
        </feature-tips>
        <property-confirm-table :confirm-data="confirmData"></property-confirm-table>
        <div class="bottom-actionbar">
            <div class="actionbar-inner">
                <bk-button theme="primary" :disabled="applyButtonDisabled" @click="handleApply">应用</bk-button>
                <bk-button theme="default">上一步</bk-button>
                <bk-button theme="default">取消</bk-button>
            </div>
        </div>
    </div>
</template>

<script>
    import { mapGetters, mapActions } from 'vuex'
    import featureTips from '@/components/feature-tips/index'
    import propertyConfirmTable from './children/property-confirm-table'
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
            this.showFeatureTips = this.featureTipsParams['hostApplyConfirm']
        },
        methods: {
            ...mapActions('objectModelProperty', [
                'searchObjectAttribute'
            ]),
            handleApply () {
            }
        }
    }
</script>

<style lang="scss" scoped>
    .host-apply-confirm-wrapper {
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
        }
    }
</style>
