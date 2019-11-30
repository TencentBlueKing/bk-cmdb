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
            <div v-else>
                {{$t('请先选择模块')}}
            </div>
        </div>
    </div>
</template>

<script>
    import { mapGetters } from 'vuex'
    import featureTips from '@/components/feature-tips/index'
    import sidebar from './children/sidebar.vue'
    import singleModuleConfig from './children/single-module-config'
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
                editing: false
            }
        },
        computed: {
            ...mapGetters(['featureTipsParams'])
        },
        created () {
            this.showFeatureTips = this.featureTipsParams['hostApply']
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
    }
</style>
