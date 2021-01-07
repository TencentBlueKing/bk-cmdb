<template>
    <div class="host-apply-confirm">
        <div class="caption">
            <div class="title">{{$t('请确认以下主机应用信息')}}</div>
            <div class="stat">
                <span class="conflict-item">
                    <span v-bk-tooltips="{ content: $t('当一台主机同属多个模块时，且模块配置了不同的属性'), width: 185 }">
                        <i class="bk-cc-icon icon-cc-tips"></i>
                    </span>
                    <i18n path="属性冲突N台">
                        <em place="num" class="conflict-num">{{conflictNum}}</em>
                    </i18n>
                </span>
                <i18n path="总检测N台">
                    <em place="num" class="check-num">{{table.total}}</em>
                </i18n>
            </div>
        </div>
        <property-confirm-table
            ref="propertyConfirmTable"
            :list="table.list"
            :total="table.total"
        >
        </property-confirm-table>
        <div :class="['bottom-actionbar', { 'is-sticky': hasScrollbar }]">
            <div class="actionbar-inner">
                <bk-button theme="default" @click="handlePrevStep">{{$t('上一步')}}</bk-button>
                <cmdb-auth :auth="{ type: $OPERATION.U_HOST_APPLY, relation: [bizId] }">
                    <bk-button
                        theme="primary"
                        slot-scope="{ disabled }"
                        :disabled="applyButtonDisabled || disabled"
                        @click="handleApply"
                    >
                        {{$t('保存并应用')}}
                    </bk-button>
                </cmdb-auth>
                <bk-button theme="default" @click="handleCancel">{{$t('取消')}}</bk-button>
            </div>
        </div>
        <leave-confirm
            v-bind="leaveConfirmConfig"
            reverse
            :title="$t('是否退出配置')"
            :content="$t('启用步骤未完成，退出将会丢失当前配置')"
            :ok-text="$t('退出')"
            :cancel-text="$t('取消')"
        >
        </leave-confirm>
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
    import { mapGetters, mapState, mapActions } from 'vuex'
    import leaveConfirm from '@/components/ui/dialog/leave-confirm'
    import propertyConfirmTable from '@/components/host-apply/property-confirm-table'
    import applyStatusModal from './children/apply-status'
    import {
        MENU_BUSINESS_HOST_AND_SERVICE,
        MENU_BUSINESS_HOST_APPLY,
        MENU_BUSINESS_HOST_APPLY_EDIT,
        MENU_BUSINESS_HOST_APPLY_FAILED
    } from '@/dictionary/menu-symbol'
    import { addResizeListener, removeResizeListener } from '@/utils/resize-events'
    export default {
        components: {
            leaveConfirm,
            applyStatusModal,
            propertyConfirmTable
        },
        data () {
            return {
                table: {
                    list: [],
                    total: 0
                },
                conflictNum: 0,
                leaveConfirmConfig: {
                    id: 'propertyConfirm',
                    active: true
                },
                applyRequest: null,
                applyButtonDisabled: false,
                hasScrollbar: false
            }
        },
        computed: {
            ...mapState('hostApply', ['propertyConfig']),
            ...mapGetters('objectBiz', ['bizId']),
            ...mapGetters(['supplierAccount']),
            isBatch () {
                return this.$route.query.batch === 1
            }
        },
        mounted () {
            addResizeListener(this.$refs.propertyConfirmTable.$el, this.resizeHandler)
        },
        beforeDestroy () {
            removeResizeListener(this.$refs.propertyConfirmTable.$el, this.resizeHandler)
        },
        beforeRouteLeave (to, from, next) {
            if (to.name !== MENU_BUSINESS_HOST_APPLY_EDIT) {
                this.$store.commit('hostApply/clearRuleDraft')
            }
            next()
        },
        created () {
            // 无配置数据时强制跳转至入口页
            if (!Object.keys(this.propertyConfig).length) {
                this.leaveConfirmConfig.active = false
                this.$routerActions.redirect({
                    name: MENU_BUSINESS_HOST_APPLY
                })
            } else {
                this.setBreadcrumbs()
                this.initData()
            }
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
                        params: this.propertyConfig,
                        config: {
                            requestId: 'getHostApplyPreview'
                        }
                    })

                    this.table.list = previewData.plans || []
                    this.table.total = previewData.count
                    this.conflictNum = previewData.unresolved_conflict_count
                } catch (e) {
                    this.applyButtonDisabled = true
                    console.error(e)
                }
            },
            setBreadcrumbs () {
                const title = this.isBatch ? '批量应用属性' : '应用属性'
                this.$store.commit('setTitle', this.$t(title))
            },
            goBack () {
                const query = {}
                if (!this.isBatch) {
                    query.module = this.propertyConfig.bk_module_ids[0]
                }
                this.$store.commit('hostApply/clearRuleDraft')
                this.$routerActions.redirect({
                    name: MENU_BUSINESS_HOST_APPLY,
                    query
                })
            },
            resizeHandler (a, b, c) {
                this.$nextTick(() => {
                    const scroller = this.$refs.propertyConfirmTable.$el.querySelector('.bk-table-body-wrapper')
                    this.hasScrollbar = scroller.scrollHeight > scroller.offsetHeight
                })
            },
            saveAndApply () {
                const { conflictResolveResult } = this.$refs.propertyConfirmTable
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

                // 合入冲突结果数据
                let propertyConfig = { ...this.propertyConfig, ...{ conflict_resolvers: conflictResolvers } }

                this.applyRequest = this.runApply({
                    bizId: this.bizId,
                    params: propertyConfig,
                    config: {
                        requestId: 'runHostApply'
                    }
                })
                this.$refs.applyStatusModal.show()

                this.applyRequest.then(() => {
                    // 应用请求完成则不需要离开确认
                    this.leaveConfirmConfig.active = false

                    const failHostIds = this.$refs.applyStatusModal.fail.map(item => item.bk_host_id)
                    propertyConfig = { ...propertyConfig, ...{ bk_host_ids: failHostIds } }
                    // 更新属性配置
                    this.$store.commit('hostApply/setPropertyConfig', propertyConfig)
                })
            },
            async handleApply () {
                const allResolved = this.$refs.propertyConfirmTable.list.every(item => item.unresolved_conflict_count === 0)
                if (allResolved) {
                    this.saveAndApply()
                } else {
                    this.$bkInfo({
                        title: this.$t('确认应用'),
                        subTitle: this.$t('您还有无法自动应用的主机属性需确认'),
                        confirmFn: () => {
                            this.saveAndApply()
                        }
                    })
                }
            },
            handleCancel () {
                this.goBack()
            },
            handlePrevStep () {
                this.leaveConfirmConfig.active = false
                this.$routerActions.back()
            },
            handleStatusModalBack () {
                this.goBack()
            },
            handleViewHost () {
                const query = {}
                if (!this.isBatch) {
                    query.node = `module-${this.propertyConfig.bk_module_ids[0]}`
                }
                this.$routerActions.redirect({
                    name: MENU_BUSINESS_HOST_AND_SERVICE,
                    query,
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
    .host-apply-confirm {
        padding: 15px 20px 0;

        .caption {
            display: flex;
            margin-bottom: 14px;
            justify-content: space-between;
            align-items: center;

            .title {
                color: #63656e;
                font-size: 14px;
            }

            .stat {
                color: #313238;
                font-size: 12px;
                margin-right: 8px;

                .conflict-item {
                    margin-right: 12px;
                }

                .conflict-num,
                .check-num {
                    font-style: normal;
                    font-weight: bold;
                    margin: 0 .2em;
                }
                .conflict-num {
                    color: #ff5656;
                }
                .check-num {
                    color: #2dcb56;
                }
            }
        }
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
                margin-right: 4px;
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
