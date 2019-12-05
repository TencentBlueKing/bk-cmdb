<template>
    <div class="single-module-config">
        <div class="config-head">
            <h2 class="config-title">
                <span class="module-name">{{module.bk_inst_name}}</span>
                <small class="last-edit-time" v-if="hasRule">( {{$t('上次编辑时间')}}：{{ruleLastEditTime}} )</small>
            </h2>
        </div>
        <div class="config-body">
            <template v-if="applyEnabled || isEdit">
                <div class="view-field" v-if="hasRule && !isEdit">
                    <div class="view-bd">
                        <div class="field-list">
                            <div class="field-list-table">
                                <property-config-table
                                    ref="propertyConfigTable"
                                    :readonly="true"
                                    :checked-property-id-list.sync="checkedPropertyIdList"
                                    :rule-list="initRuleList"
                                >
                                </property-config-table>
                            </div>
                        </div>
                    </div>
                    <div class="view-ft">
                        <bk-button theme="primary" @click="handleEdit">{{$t('编辑')}}</bk-button>
                        <bk-button theme="default" :disabled="!hasConflict" @click="handleViewConflict">
                            <span v-bk-tooltips="{ content: $t('无冲突需处理') }" v-if="!hasConflict">
                                {{$t('查看冲突')}}<em class="conflict-num">{{conflictNum}}</em>
                            </span>
                            <span v-else>
                                {{$t('查看冲突')}}<em class="conflict-num">{{conflictNum}}</em>
                            </span>
                        </bk-button>
                        <bk-button theme="default" @click="handleCloseApply">{{$t('关闭自动应用')}}</bk-button>
                    </div>
                </div>
                <div :class="['choose-field', { 'not-choose': !checkedPropertyIdList.length }]" v-else>
                    <div class="choose-hd">
                        <span class="label">{{$t('自动应用字段：')}}</span>
                        <bk-button theme="default" icon="plus" @click="handleChooseField">选择字段</bk-button>
                    </div>
                    <div class="choose-bd" v-show="checkedPropertyIdList.length">
                        <property-config-table
                            ref="propertyConfigTable"
                            :checked-property-id-list.sync="checkedPropertyIdList"
                            :rule-list="initRuleList"
                            @property-value-change="handlePropertyValueChange"
                        >
                        </property-config-table>
                    </div>
                    <div class="choose-ft">
                        <bk-button theme="primary" :disabled="nextButtonDisabled" @click="handleNextStep">下一步</bk-button>
                        <bk-button theme="default" @click="handleCancel">取消</bk-button>
                    </div>
                </div>
            </template>
            <template v-else>
                <div class="empty" v-if="!hasRule">
                    <div class="desc">
                        <i class="bk-cc-icon icon-cc-tips"></i>
                        <span>{{$t('当前模块未启用自动应用策略')}}</span>
                    </div>
                    <div class="action">
                        <bk-button theme="primary" :outline="true" @click="handleEdit">立即启用</bk-button>
                    </div>
                </div>
                <div class="view-field" v-else>
                    <div class="view-bd">
                        <div class="field-list">
                            <div class="field-list-table disabled">
                                <property-config-table
                                    ref="propertyConfigTable"
                                    :readonly="true"
                                    :checked-property-id-list.sync="checkedPropertyIdList"
                                    :rule-list="initRuleList"
                                >
                                </property-config-table>
                            </div>
                            <div class="closed-mask">
                                <div class="empty">
                                    <div class="desc">
                                        <i class="bk-cc-icon icon-cc-tips"></i>
                                        <span>{{$t('该模块已关闭属性自动应用')}}</span>
                                    </div>
                                    <div class="action">
                                        <bk-button theme="primary" :outline="true" @click="handleEdit">{{$t('重新启用')}}</bk-button>
                                    </div>
                                </div>
                            </div>
                        </div>
                    </div>
                </div>
            </template>
        </div>
        <host-property-modal
            :visible.sync="propertyModalVisible"
            :checked-list.sync="checkedPropertyIdList"
        >
        </host-property-modal>
        <leave-confirm
            :id="leaveConfirm.id"
            :active="leaveConfirm.active"
            title="是否放弃？"
            content="启用步骤未完成，是否放弃当前配置"
            ok-text="留在当前页"
            cancel-text="确认放弃"
        >
        </leave-confirm>
    </div>
</template>

<script>
    import { mapGetters, mapState } from 'vuex'
    import ConfirmStore from '@/components/ui/dialog/confirm-store.js'
    import leaveConfirm from '@/components/ui/dialog/leave-confirm'
    import hostPropertyModal from './host-property-modal'
    import propertyConfigTable from './property-config-table'
    export default {
        components: {
            leaveConfirm,
            hostPropertyModal,
            propertyConfigTable
        },
        props: {
            module: {
                type: Object,
                default: () => ({})
            },
            editing: {
                type: Boolean
            }
        },
        data () {
            return {
                initRuleList: [],
                checkedPropertyIdList: [],
                // 用于取消编辑时的还原
                checkedPropertyIdListCopy: [],
                conflictNum: 0,
                hasRule: false,
                isEdit: this.editing,
                nextButtonDisabled: true,
                propertyModalVisible: false,
                clearRules: false,
                leaveConfirm: {
                    id: 'singleModule',
                    active: false
                }
            }
        },
        computed: {
            ...mapGetters('objectBiz', ['bizId']),
            ...mapGetters('hostApply', ['configPropertyList']),
            ...mapState('hostApply', ['ruleDraft']),
            moduleId () {
                return this.module.bk_inst_id
            },
            applyEnabled () {
                return this.module.host_apply_enabled
            },
            hasConflict () {
                return this.conflictNum > 0
            },
            ruleLastEditTime () {
                const lastTimeList = this.initRuleList.map(rule => new Date(rule.last_time).getTime())
                const latestTime = Math.max(...lastTimeList)
                return this.$tools.formatTime(latestTime, 'YYYY-MM-DD HH:mm:ss')
            },
            hasRuleDraft () {
                return Object.keys(this.ruleDraft).length > 0
            }
        },
        watch: {
            module () {
                this.getConfigData()
                // 切换模块时将草稿数据清空
                this.$store.commit('hostApply/clearRuleDraft')
            },
            editing (value) {
                this.isEdit = value
            },
            checkedPropertyIdList () {
                this.toggleNextButtonDisabled()
            },
            isEdit (value) {
                this.$emit('update:editing', value)
                this.leaveConfirm.active = value
                this.toggleNextButtonDisabled()
            }
        },
        created () {
            this.getConfigData()
            this.isEdit = this.hasRuleDraft
        },
        methods: {
            async getConfigData () {
                try {
                    const ruleData = await this.getRules()

                    // 重置配置表格数据
                    if (this.$refs.propertyConfigTable) {
                        this.$refs.propertyConfigTable.reset()
                    }

                    this.initRuleList = ruleData.info || []
                    this.hasRule = ruleData.count > 0
                    const checkedPropertyIdList = this.initRuleList.map(item => item.bk_attribute_id)
                    this.checkedPropertyIdList = this.hasRuleDraft ? [...this.checkedPropertyIdList, ...checkedPropertyIdList] : checkedPropertyIdList
                    this.checkedPropertyIdListCopy = [...this.checkedPropertyIdList]

                    if (this.hasRule && this.applyEnabled) {
                        const previewData = await this.getApplyPreview()
                        this.conflictNum = previewData.unresolved_conflict_count
                    }
                } catch (e) {
                    console.log(e)
                }
            },
            getRules () {
                return this.$store.dispatch('hostApply/getRules', {
                    bizId: this.bizId,
                    params: {
                        bk_module_ids: [this.moduleId]
                    },
                    config: {
                        requestId: `getHostApplyRules`
                    }
                })
            },
            getApplyPreview () {
                return this.$store.dispatch('hostApply/getApplyPreview', {
                    bizId: this.bizId,
                    params: {
                        bk_module_ids: [this.moduleId]
                    },
                    config: {
                        requestId: `getHostApplyPreview`
                    }
                })
            },
            toggleNextButtonDisabled () {
                this.$nextTick(() => {
                    if (this.$refs.propertyConfigTable) {
                        const { modulePropertyList } = this.$refs.propertyConfigTable
                        const everyTruthy = modulePropertyList.every(property => {
                            const validTruthy = property.__extra__.valid !== false
                            let valueTruthy = property.__extra__.value
                            if (property.bk_property_type === 'bool') {
                                valueTruthy = true
                            } else if (property.bk_property_type === 'int') {
                                valueTruthy = valueTruthy !== null && String(valueTruthy)
                            }
                            return valueTruthy && validTruthy
                        })
                        this.nextButtonDisabled = !this.checkedPropertyIdList.length || !everyTruthy
                    }
                })
            },
            emptyRules () {
                this.checkedPropertyIdList = []
                this.hasRule = false
            },
            async handleNextStep () {
                const { modulePropertyList, removeRuleIds } = this.$refs.propertyConfigTable
                const additionalRules = modulePropertyList.map(property => ({
                    bk_attribute_id: property.id,
                    bk_module_id: this.moduleId,
                    bk_property_value: property.__extra__.value
                }))

                const savePropertyConfig = {
                    // 模块列表
                    bk_module_ids: [this.moduleId],
                    // 附加的规则
                    additional_rules: additionalRules,
                    // 删除的规则，来源于编辑表格删除
                    remove_rule_ids: removeRuleIds
                }

                this.$store.commit('hostApply/setPropertyConfig', savePropertyConfig)
                this.$store.commit('hostApply/setRuleDraft', {
                    moduleIds: [this.moduleId],
                    rules: modulePropertyList
                })

                // 使离开确认失活
                this.leaveConfirm.active = false
                this.$nextTick(function () {
                    this.$router.push({
                        name: 'hostApplyConfirm'
                    })
                })
            },
            handleViewConflict () {
                this.$router.push({
                    name: 'hostApplyConflict',
                    query: {
                        mid: this.moduleId
                    }
                })
            },
            handleCloseApply () {
                const h = this.$createElement
                this.$bkInfo({
                    title: this.$t('确认关闭？'),
                    extCls: 'close-apply-confirm-modal',
                    subHeader: h('div', { class: 'content' }, [
                        h('p', { class: 'tips' }, this.$t('关闭后转入模块的主机属性不再自动被应用')),
                        h('bk-checkbox', {
                            props: {
                                checked: true,
                                trueValue: true,
                                falseFalue: false
                            },
                            on: {
                                change: (value) => (this.clearRules = !value)
                            }
                        }, '保留当前自动应用策略')
                    ]),
                    confirmFn: async () => {
                        try {
                            await this.$store.dispatch('hostApply/setEnableStatus', {
                                bizId: this.bizId,
                                moduleId: this.moduleId,
                                params: {
                                    host_apply_enabled: false,
                                    clear_rules: this.clearRules
                                }
                            })

                            this.$success(this.$t('关闭成功'))
                            if (this.clearRules) {
                                this.emptyRules()
                            }
                            this.$emit('after-close', this.moduleId, this.clearRules)
                        } catch (e) {
                            console.log(e)
                        }
                    }
                })
            },
            handlePropertyValueChange () {
                this.toggleNextButtonDisabled()
            },
            handleEdit () {
                this.isEdit = true
            },
            handleChooseField () {
                this.propertyModalVisible = true
            },
            async handleCancel (id) {
                const leaveConfirmResult = await ConfirmStore.popup(this.leaveConfirm.id)
                if (leaveConfirmResult) {
                    this.checkedPropertyIdList = [...this.checkedPropertyIdListCopy]
                    this.isEdit = false
                }
            }
        }
    }
</script>

<style lang="scss" scoped>
    .single-module-config {
        display: flex;
        flex-direction: column;
        height: 100%;

        .config-head,
        .config-foot {
            flex: none;
        }
        .config-body {
            flex: auto;
        }
    }
    .config-title {
        height: 32px;
        line-height: 32px;
        font-size: 14px;
        color: #313238;
        font-weight: 700;

        .last-edit-time {
            font-size: 12px;
            font-weight: 400;
            color: #979ba5;
            margin-left: .2em;
        }
    }
    .empty {
        display: flex;
        flex-direction: column;
        align-items: center;
        justify-content: center;
        height: 80%;

        .desc {
            font-size: 14px;
            color: #63656e;
        }
        .action {
            margin-top: 18px;
        }
    }

    .choose-field {
        padding: 16px 2px;
        .choose-hd {
            .label {
                font-size: 14px;
                color: #63656e;
                margin-right: 8px;
            }
        }
        .choose-bd {
            margin-top: 20px;

            .form-element-content {
                padding: 4px 0;
            }
        }
        .choose-ft {
            margin-top: 20px;
            .bk-button {
                min-width: 86px;
            }
        }

        &.not-choose {
            .choose-ft {
                margin-left: 111px;
            }
        }
    }

    .view-field {
        .field-list {
            position: relative;
            &:hover {
                .closed-mask {
                    display: block;
                }
            }

            .field-list-table {
                &.disabled {
                    opacity: 0.2;
                }
            }
            .closed-mask {
                position: absolute;
                width: 100%;
                height: 100%;
                min-height: 210px;
                left: 0;
                top: 0;
                display: none;
            }
        }
        .view-bd,
        .view-ft {
            margin-top: 20px;
            .bk-button {
                min-width: 86px;
            }
        }
    }

    .conflict-num {
        font-size: 12px;
        color: #fff;
        background: #c4c6cc;
        border-radius: 8px;
        font-style: normal;
        padding: 0px 4px;
        font-family: arial;
        margin-left: 4px;
    }

    .close-apply-confirm-modal {
        .content {
            font-size: 14px;
        }
        .tips {
            margin: 12px 0;
        }
    }
</style>

<style lang="scss">
    .close-apply-confirm-modal {
        .bk-dialog-sub-header {
            padding-left: 32px !important;
            padding-right: 32px !important;
        }
    }
</style>
