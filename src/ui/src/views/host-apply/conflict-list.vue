<template>
    <div class="conflict-list">
        <feature-tips
            :feature-name="'hostApply'"
            :show-tips="showFeatureTips"
            :desc="$t('因为以下主机复用模块的自动应用策略不一致，导致策略失效，需要手动维护不一致的属性。要彻底解决此问题，可以修改复用模块的策略为一致或移除模块的策略配置')"
            @close-tips="showFeatureTips = false">
        </feature-tips>
        <property-confirm-table
            ref="propertyConfirmTable"
            :list="table.list"
            :total="table.total"
        >
        </property-confirm-table>
        <div class="bottom-actionbar">
            <div class="actionbar-inner">
                <bk-button theme="primary" :disabled="applyButtonDisabled" @click="handleApply">应用</bk-button>
                <bk-button theme="default" @click="handleBack">返回</bk-button>
            </div>
        </div>
        <apply-status-modal
            ref="applyStatusModal"
            :request="applyRequest"
            @return="handleStatusModalBack"
            @view-host="handleViewHost"
            @view-failed="handleViewFailed"
        >
        </apply-status-modal>
    </div>
</template>

<script>
    import { mapGetters, mapActions } from 'vuex'
    import featureTips from '@/components/feature-tips/index'
    import propertyConfirmTable from './children/property-confirm-table'
    import applyStatusModal from './children/apply-status'
    import { MENU_BUSINESS_HOST_AND_SERVICE, MENU_BUSINESS_HOST_APPLY } from '@/dictionary/menu-symbol'
    export default {
        components: {
            featureTips,
            applyStatusModal,
            propertyConfirmTable
        },
        data () {
            return {
                showFeatureTips: false,
                table: {
                    list: [],
                    total: 0
                },
                applyRequest: null,
                applyButtonDisabled: false
            }
        },
        computed: {
            ...mapGetters('objectBiz', ['bizId']),
            ...mapGetters(['featureTipsParams', 'supplierAccount']),
            moduleIds () {
                const mid = this.$route.query.mid
                let moduleIds = []
                if (mid) {
                    moduleIds = String(mid).split(',').map(id => Number(id))
                }
                return moduleIds
            }
        },
        watch: {
        },
        created () {
            this.showFeatureTips = this.featureTipsParams['hostApplyConflict']
            this.setBreadcrumbs()
            this.initData()
        },
        methods: {
            ...mapActions('hostApply', [
                'getApplyPreview',
                'runApply'
            ]),
            async initData () {
                try {
                    const previewData = await this.getApplyPreview({
                        bizId: this.bizId,
                        params: { bk_module_ids: this.moduleIds },
                        config: {
                            requestId: 'getHostApplyPreview'
                        }
                    })
                    const conflictList = (previewData.plans || []).filter(item => item.unresolved_conflict_count > 0)
                    this.table.list = conflictList
                    this.table.total = previewData.unresolved_conflict_count
                } catch (e) {
                    this.applyButtonDisabled = true
                    console.error(e)
                }
            },
            setBreadcrumbs () {
                this.$store.commit('setTitle', this.$t('策略失效列表'))
                this.$store.commit('setBreadcrumbs', [{
                    label: this.$t('主机属性自动应用'),
                    route: {
                        name: MENU_BUSINESS_HOST_APPLY
                    }
                }, {
                    label: this.$t('策略失效列表')
                }])
            },
            goBack () {
                this.$router.push({
                    name: MENU_BUSINESS_HOST_APPLY
                })
            },
            postApply () {
                const conflictResolveResult = this.$refs.propertyConfirmTable.conflictResolveResult
                const conflictResolvers = []
                Object.keys(conflictResolveResult).forEach(key => {
                    const propertyList = conflictResolveResult[key]
                    propertyList.forEach(property => {
                        conflictResolvers.push({
                            bk_host_id: Number(key),
                            bk_attribute_id: property.id,
                            bk_property_value: property.__extra__.value
                        })
                    })
                })

                this.applyRequest = this.runApply({
                    bizId: this.bizId,
                    params: {
                        bk_module_ids: this.moduleIds,
                        conflict_resolvers: conflictResolvers
                    },
                    config: {
                        requestId: 'runHostApply'
                    }
                })
                this.$refs.applyStatusModal.show()
            },
            async handleApply () {
                const allResolved = this.$refs.propertyConfirmTable.list.every(item => item.unresolved_conflict_count === 0)
                if (allResolved) {
                    this.postApply()
                } else {
                    this.$bkInfo({
                        title: this.$t('确认应用'),
                        subTitle: this.$t('您还有无法自动应用的主机属性需确认，是要保留主机原有属性值不做修改吗？'),
                        confirmFn: () => {
                            this.postApply()
                        }
                    })
                }
            },
            handleBack () {
                this.$router.back()
            },
            handleStatusModalBack () {
                this.goBack()
            },
            handleViewHost () {
                this.$router.push({
                    name: MENU_BUSINESS_HOST_AND_SERVICE
                })
            },
            handleViewFailed () {
                this.$router.push({
                    name: 'hostApplyFailed'
                })
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
