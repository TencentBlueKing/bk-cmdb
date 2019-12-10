<template>
    <div class="single-module-config">
        <div class="config-body">
            <div :class="['choose-field', { 'not-choose': !checkedPropertyIdList.length }]">
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
        </div>
        <host-property-modal
            :visible.sync="propertyModalVisible"
            :checked-list.sync="checkedPropertyIdList"
        >
        </host-property-modal>
        <leave-confirm
            :id="leaveConfirmConfig.id"
            :active="leaveConfirmConfig.active"
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
    import leaveConfirm from '@/components/ui/dialog/leave-confirm'
    import hostPropertyModal from './host-property-modal'
    import propertyConfigTable from './property-config-table'
    import { MENU_BUSINESS_HOST_APPLY } from '@/dictionary/menu-symbol'
    export default {
        name: 'single-module-config',
        components: {
            leaveConfirm,
            hostPropertyModal,
            propertyConfigTable
        },
        props: {
            moduleIds: {
                type: Array,
                default: () => ([])
            }
        },
        data () {
            return {
                initRuleList: [],
                checkedPropertyIdList: [],
                nextButtonDisabled: true,
                propertyModalVisible: false,
                leaveConfirmConfig: {
                    id: 'singleModule',
                    active: true
                }
            }
        },
        computed: {
            ...mapGetters('objectBiz', ['bizId']),
            ...mapState('hostApply', ['ruleDraft']),
            moduleId () {
                return this.moduleIds[0]
            },
            hasRuleDraft () {
                return Object.keys(this.ruleDraft).length > 0
            }
        },
        watch: {
            checkedPropertyIdList () {
                this.toggleNextButtonDisabled()
            }
        },
        created () {
            this.initData()
        },
        methods: {
            async initData () {
                try {
                    const ruleData = await this.getRules()
                    this.initRuleList = ruleData.info || []
                    const checkedPropertyIdList = this.initRuleList.map(item => item.bk_attribute_id)
                    this.checkedPropertyIdList = this.hasRuleDraft ? [...this.checkedPropertyIdList] : checkedPropertyIdList
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
                this.leaveConfirmConfig.active = false
                this.$nextTick(function () {
                    this.$router.push({
                        name: 'hostApplyConfirm'
                    })
                })
            },
            handlePropertyValueChange () {
                this.toggleNextButtonDisabled()
            },
            handleChooseField () {
                this.propertyModalVisible = true
            },
            handleCancel () {
                this.$router.push({
                    name: MENU_BUSINESS_HOST_APPLY,
                    query: {
                        module: this.moduleId
                    }
                })
            }
        }
    }
</script>

<style lang="scss" scoped>
    .single-module-config {
        display: flex;
        flex-direction: column;
        width: 1066px;
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
</style>
