<template>
    <div class="host-apply">
        <feature-tips
            :feature-name="'hostApply'"
            :show-tips="showFeatureTips"
            :desc="$t('主机属性自动应用功能提示')"
            @close-tips="showFeatureTips = false">
        </feature-tips>
        <cmdb-resize-layout class="tree-layout fl"
            direction="right"
            :handler-offset="3"
            :min="200"
            :max="480"
        >
            <sidebar
                :action-mode.sync="actionMode"
                :editing="editing"
                @module-selected="handleSelectModule"
            ></sidebar>
        </cmdb-resize-layout>
        <div class="main-layout">
            <template v-if="selectedModule.bk_inst_id">
                <single-module-config :module="selectedModule" :editing.sync="editing"></single-module-config>
            </template>
            <div class="empty" v-else>
                <i18n path="您还未XXX">
                    <span place="action">{{$t('创建')}}</span>
                    <span place="resource">{{$t('模块')}}</span>
                    <span place="link">
                        <bk-button class="text-btn"
                            text
                            place="link"
                            theme="primary"
                            @click="$router.push({ name: hostAndServiceRouteName })"
                        >
                            {{`立即${$t('创建')}`}}
                        </bk-button>
                    </span>
                </i18n>
            </div>
        </div>
    </div>
</template>

<script>
    import { mapGetters } from 'vuex'
    import featureTips from '@/components/feature-tips/index'
    import sidebar from './children/sidebar.vue'
    import singleModuleConfig from './children/single-module-config'
    import { MENU_BUSINESS_HOST_AND_SERVICE } from '@/dictionary/menu-symbol'
    export default {
        components: {
            sidebar,
            featureTips,
            singleModuleConfig
        },
        data () {
            return {
                actionMode: '',
                selectedModule: {},
                showFeatureTips: false,
                editing: false,
                hostAndServiceRouteName: MENU_BUSINESS_HOST_AND_SERVICE
            }
        },
        computed: {
            ...mapGetters(['featureTipsParams'])
        },
        created () {
            this.showFeatureTips = this.featureTipsParams['hostApply']
        },
        beforeRouteLeave (to, from, next) {
            if (to.name !== 'hostApplyConfirm') {
                this.$store.commit('hostApply/clearRuleDraft')
            }
            next()
        },
        methods: {
            handleSelectModule (data) {
                this.editing = false
                this.selectedModule = data
            }
        }
    }
</script>

<style lang="scss" scoped>
    .host-apply {
        padding: 0 20px;
    }
    .tree-layout {
        width: 280px;
        height: 100%;
        border: 1px solid #dcdee5;
    }
    .main-layout {
        @include scrollbar-x;
        height: 100%;
        border: 1px solid #dcdee5;
        border-left: 0 none;
        padding: 10px 18px;

        .empty {
            display: flex;
            width: 100%;
            height: 100%;
            align-items: center;
            justify-content: center;
            font-size: 14px;
        }
    }
</style>
