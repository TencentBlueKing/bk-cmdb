<template>
    <div class="host-apply-wrapper">
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
                @module-selected="handleSelectModule"
            ></sidebar>
        </cmdb-resize-layout>
        <div class="main-layout">
            <template>
                <!-- <single-module></single-module> -->
                <component :is="currentModeComponent" :data="selectedModule"></component>
            </template>
        </div>
    </div>
</template>

<script>
    import { mapGetters, mapActions } from 'vuex'
    import featureTips from '@/components/feature-tips/index'
    import sidebar from './children/sidebar.vue'
    import singleModule from './children/single-module'
    import multiModule from './children/multi-module'
    export default {
        components: {
            sidebar,
            featureTips,
            singleModule,
            multiModule
        },
        data () {
            return {
                actionMode: '',
                selectedModule: {},
                showFeatureTips: false,
                currentModeComponent: singleModule
            }
        },
        computed: {
            ...mapGetters(['featureTipsParams', 'supplierAccount'])
        },
        watch: {
        },
        created () {
            this.showFeatureTips = this.featureTipsParams['hostApply']
            this.getHostPropertyList()
        },
        beforeDestroy () {
            this.$store.commit('businessTopology/resetState')
        },
        methods: {
            ...mapActions('objectModelProperty', [
                'searchObjectAttribute'
            ]),
            async getHostPropertyList () {
                try {
                    const data = await this.searchObjectAttribute({
                        params: this.$injectMetadata({
                            bk_obj_id: 'host',
                            bk_supplier_account: this.supplierAccount
                        }),
                        config: {
                            requestId: 'getHostPropertyList',
                            fromCache: true
                        }
                    })

                    this.$store.commit('hosts/setPropertyList', data)
                } catch (e) {
                    console.error(e)
                }
            },
            handleAddNode () {
                this.$refs.topologyTree.showCreateDialog(this.selectedNode)
            },
            handleSelectModule (data) {
                this.selectedModule = data
                console.dir(data)
            }
        }
    }
</script>

<style lang="scss" scoped>
    .host-apply-wrapper {
        padding: 0 20px;
    }
    .tree-layout {
        width: 280px;
        height: 100%;
        padding: 10px 0;
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
