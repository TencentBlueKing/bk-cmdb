<template>
    <div class="conflict-list">
        <property-confirm-table
            ref="propertyConfirmTable"
            :list="table.list"
            :total="table.total"
        >
        </property-confirm-table>
        <div :class="['bottom-actionbar', { 'is-sticky': hasScrollbar }]">
            <div class="actionbar-inner">
                <cmdb-auth :auth="{ type: $OPERATION.U_HOST_APPLY, relation: [bizId] }">
                    <bk-button
                        theme="primary"
                        slot-scope="{ disabled }"
                        :disabled="applyButtonDisabled || disabled"
                        @click="handleApply"
                    >
                        {{$t('应用')}}
                    </bk-button>
                </cmdb-auth>
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
    import propertyConfirmTable from '@/components/host-apply/property-confirm-table'
    import applyStatusModal from './children/apply-status'
    import {
        MENU_BUSINESS_HOST_AND_SERVICE,
        MENU_BUSINESS_HOST_APPLY,
        MENU_BUSINESS_HOST_APPLY_FAILED
    } from '@/dictionary/menu-symbol'
    import { addResizeListener, removeResizeListener } from '@/utils/resize-events'
    export default {
        components: {
            applyStatusModal,
            propertyConfirmTable
        },
        data () {
            return {
                table: {
                    list: [],
                    total: 0
                },
                applyRequest: null,
                applyButtonDisabled: false,
                hasScrollbar: false
            }
        },
        computed: {
            ...mapGetters('objectBiz', ['bizId']),
            moduleIds () {
                const mid = this.$route.query.mid
                let moduleIds = []
                if (mid) {
                    moduleIds = String(mid).split(',').map(id => Number(id))
                }
                return moduleIds
            }
        },
        mounted () {
            addResizeListener(this.$refs.propertyConfirmTable.$el, this.resizeHandler)
        },
        beforeDestroy () {
            removeResizeListener(this.$refs.propertyConfirmTable.$el, this.resizeHandler)
        },
        created () {
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
                this.$store.commit('setTitle', this.$t('策略失效主机'))
            },
            goBack () {
                this.$routerActions.redirect({
                    name: MENU_BUSINESS_HOST_APPLY
                })
            },
            resizeHandler (a, b, c) {
                this.$nextTick(() => {
                    const scroller = this.$refs.propertyConfirmTable.$el.querySelector('.bk-table-body-wrapper')
                    this.hasScrollbar = scroller.scrollHeight > scroller.offsetHeight
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
                        subTitle: this.$t('您还有无法自动应用的主机属性需确认'),
                        confirmFn: () => {
                            this.postApply()
                        }
                    })
                }
            },
            handleStatusModalBack () {
                this.goBack()
            },
            handleViewHost () {
                this.$routerActions.redirect({
                    name: MENU_BUSINESS_HOST_AND_SERVICE,
                    query: {
                        node: `module-${this.moduleIds[0]}`
                    },
                    history: true
                })
            },
            handleViewFailed () {
                this.$routerActions.redirect({
                    name: MENU_BUSINESS_HOST_APPLY_FAILED,
                    query: this.$route.query,
                    history: true
                })
            }
        }
    }
</script>

<style lang="scss" scoped>
    .conflict-list {
        padding: 15px 20px 0;
    }

    .bottom-actionbar {
        width: 100%;
        height: 50px;
        bottom: 0;
        left: 0;
        z-index: 100;

        .actionbar-inner {
            padding: 20px 0 0 0;
            .bk-button {
                min-width: 86px;
            }
        }

        &.is-sticky {
            position: absolute;
            background: #fff;
            border-top: 1px solid #dcdee5;

            .actionbar-inner {
                padding: 8px 0 0 20px;
            }
        }
    }
</style>
