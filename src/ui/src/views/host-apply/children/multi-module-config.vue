<template>
    <div class="multi-module-config">
        <div class="config-bd">
            <div class="config-item">
                <div class="item-label">已选择 {{moduleIds.length}} 个模块：</div>
                <div class="item-content">
                    <div :class="['module-list', { 'show-more': showMore.isMoreModuleShowed }]" ref="moduleList">
                        <div class="module-item" :title="getModulePath(id)" v-for="(id, index) in moduleIds" :key="index">
                            <span class="module-icon">{{$i18n.locale === 'en' ? 'M' : '模'}}</span>
                            {{$parent.getModuleName(id)}}
                        </div>
                        <div
                            :class="['module-item', 'more', { 'opened': showMore.isMoreModuleShowed }]"
                            :style="{ left: `${showMore.linkLeft}px` }"
                            v-show="showMore.showLink" @click="handleShowMore"
                        >
                            {{showMore.isMoreModuleShowed ? '收起' : '展开更多'}}<i class="bk-cc-icon icon-cc-arrow-down"></i>
                        </div>
                    </div>
                </div>
            </div>
            <div class="config-item">
                <div class="item-label">
                    {{$t(isDel ? '请勾选要删除的字段：' : '已配置的字段：')}}
                </div>
                <div class="item-content">
                    <div class="choose-toolbar">
                        <bk-button theme="default" icon="plus" @click="handleChooseField" v-if="!isDel">选择字段</bk-button>
                        <span class="tips"><i class="bk-cc-icon icon-cc-tips"></i><span>此功能可以批量设置字段的自动应用，不需要批量变更的字段需点击“删除”从列表中移除</span></span>
                    </div>
                    <div class="config-table" v-show="checkedPropertyIdList.length">
                        <property-config-table
                            ref="configEditTable"
                            :multiple="true"
                            :readonly="isDel"
                            :deletable="isDel"
                            :checked-property-id-list.sync="checkedPropertyIdList"
                            :rule-list="initRuleList"
                            :module-id-list="moduleIds"
                            @property-value-change="handlePropertyValueChange"
                            @selection-change="handlePropertySelectionChange"
                        >
                        </property-config-table>
                    </div>
                </div>
            </div>
        </div>
        <div class="config-ft">
            <bk-button theme="primary" :disabled="nextButtonDisabled" @click="handleNextStep" v-if="!isDel">下一步</bk-button>
            <bk-button theme="primary" :disabled="delButtonDisabled" @click="handleDel" v-else>确定删除</bk-button>
            <bk-button theme="default" @click="handleCancel">取消</bk-button>
        </div>

        <host-property-modal
            :visible.sync="propertyModalVisible"
            :checked-list.sync="checkedPropertyIdList"
        >
        </host-property-modal>
        <leave-confirm
            v-bind="leaveConfirmConfig"
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
        name: 'multi-module-config',
        components: {
            leaveConfirm,
            hostPropertyModal,
            propertyConfigTable
        },
        props: {
            moduleIds: {
                type: Array,
                default: () => ([])
            },
            action: {
                type: String,
                default: ''
            }
        },
        data () {
            return {
                initRuleList: [],
                checkedPropertyIdList: [],
                showMore: {
                    moduleListMaxRow: 2,
                    showLink: false,
                    isMoreModuleShowed: false,
                    linkLeft: 0
                },
                selectedPropertyRow: [],
                propertyModalVisible: false,
                nextButtonDisabled: false,
                delButtonDisabled: true,
                leaveConfirmConfig: {
                    id: 'multiModule',
                    active: true
                }
            }
        },
        computed: {
            ...mapGetters('objectBiz', ['bizId']),
            ...mapState('hostApply', ['ruleDraft']),
            isDel () {
                return this.action === 'batch-del'
            },
            hasRuleDraft () {
                return Object.keys(this.ruleDraft).length > 0
            }
        },
        watch: {
            checkedPropertyIdList (val) {
                this.$nextTick(() => {
                    this.toggleNextButtonDisabled()
                })
            }
        },
        created () {
            this.initData()
            this.leaveConfirmConfig.active = !this.isDel
        },
        mounted () {
            this.setShowMoreLinkStatus()
            window.addEventListener('resize', this.setShowMoreLinkStatus)
        },
        beforeDestroy () {
            window.removeEventListener('resize', this.setShowMoreLinkStatus)
        },
        methods: {
            async initData () {
                try {
                    const ruleData = await this.getRules()
                    this.initRuleList = ruleData.info
                    const attrIds = this.initRuleList.map(item => item.bk_attribute_id)
                    const checkedPropertyIdList = [...new Set(attrIds)]
                    this.checkedPropertyIdList = this.hasRuleDraft ? [...new Set([...this.checkedPropertyIdList])] : checkedPropertyIdList
                } catch (e) {
                    console.log(e)
                }
            },
            getRules () {
                return this.$store.dispatch('hostApply/getRules', {
                    bizId: this.bizId,
                    params: {
                        bk_module_ids: this.moduleIds
                    },
                    config: {
                        requestId: 'getHostApplyConfigs'
                    }
                })
            },
            getModulePath (id) {
                return this.$parent.getModulePath(id)
            },
            setShowMoreLinkStatus () {
                const moduleList = this.$refs.moduleList
                const moduleItemEl = moduleList.getElementsByClassName('module-item')[0]
                const moduleItemStyle = getComputedStyle(moduleItemEl)
                const moduleItemWidth = moduleItemEl.offsetWidth + parseInt(moduleItemStyle.marginLeft, 10) + parseInt(moduleItemStyle.marginRight, 10)
                const moduleListWidth = moduleList.clientWidth
                const maxCountInRow = Math.floor(moduleListWidth / moduleItemWidth)
                const rowCount = Math.ceil(this.moduleIds.length / maxCountInRow)
                this.showMore.showLink = rowCount > this.showMore.moduleListMaxRow
                this.showMore.linkLeft = moduleItemWidth * (maxCountInRow - 1)
            },
            toggleNextButtonDisabled () {
                const modulePropertyList = this.$refs.configEditTable.modulePropertyList
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
            },
            goBack () {
                // 删除离开不用确认
                this.leaveConfirmConfig.active = !this.isDel
                this.$nextTick(function () {
                    // 回到入口页
                    this.$router.push({
                        name: MENU_BUSINESS_HOST_APPLY
                    })
                })
            },
            handleNextStep () {
                const { modulePropertyList, ignoreRuleIds } = this.$refs.configEditTable
                const additionalRules = []
                this.moduleIds.forEach(moduleId => {
                    modulePropertyList.forEach(property => {
                        additionalRules.push({
                            bk_attribute_id: property.id,
                            bk_module_id: moduleId,
                            bk_property_value: property.__extra__.value
                        })
                    })
                })

                const savePropertyConfig = {
                    // 模块列表
                    bk_module_ids: this.moduleIds,
                    // 附加的规则
                    additional_rules: additionalRules,
                    // 删除的规则，来源于编辑表格删除
                    ignore_rule_ids: ignoreRuleIds
                }
                this.$store.commit('hostApply/setPropertyConfig', savePropertyConfig)
                this.$store.commit('hostApply/setRuleDraft', {
                    moduleIds: this.moduleIds,
                    rules: modulePropertyList
                })

                this.leaveConfirmConfig.active = false
                this.$nextTick(function () {
                    this.$router.push({
                        name: 'hostApplyConfirm',
                        query: {
                            batch: 1
                        }
                    })
                })
            },
            handleDel () {
                this.$bkInfo({
                    title: this.$t('确认删除自动应用字段？'),
                    subTitle: this.$t('删除后，将会移除字段在对应模块中的配置'),
                    confirmFn: async () => {
                        const ruleIds = this.selectedPropertyRow.reduce((acc, cur) => acc.concat(cur.__extra__.ruleList.map(item => item.id)), [])
                        try {
                            await this.$store.dispatch('hostApply/deleteRules', {
                                bizId: this.bizId,
                                params: {
                                    data: {
                                        host_apply_rule_ids: ruleIds
                                    }
                                }
                            })

                            this.goBack()
                        } catch (e) {
                            console.log(e)
                        }
                    }
                })
            },
            handleCancel () {
                this.$store.commit('hostApply/clearRuleDraft')
                this.goBack()
            },
            handlePropertySelectionChange (value) {
                this.selectedPropertyRow = value
                this.delButtonDisabled = this.selectedPropertyRow.length <= 0
            },
            handlePropertyValueChange () {
                this.toggleNextButtonDisabled()
            },
            handleChooseField () {
                this.propertyModalVisible = true
            },
            handleShowMore () {
                this.showMore.isMoreModuleShowed = !this.showMore.isMoreModuleShowed
            }
        }
    }
</script>
<style lang="scss" scoped>
    .multi-module-config {
        // width: 1066px;
        --labelWidth: 180px;
        .config-item {
            display: flex;
            margin: 8px 0;

            .item-label {
                flex: none;
                width: var(--labelWidth);
                font-size: 14px;
                font-weight: bold;
                color: #63656e;
                text-align: right;
                margin-right: 12px;
            }
            .item-content {
                flex: auto;
            }

            .choose-toolbar {
                margin-bottom: 18px;;
                .tips {
                    font-size: 12px;
                    margin-left: 8px;
                    .icon-cc-tips {
                        margin-right: 8px;
                    }
                }
            }
        }
        .config-ft {
            margin: 20px 0 20px calc(var(--labelWidth) + 12px);
            .bk-button {
                min-width: 86px;
            }
        }
    }

    .module-list {
        position: relative;
        max-height: 72px;
        overflow: hidden;
        transition: all .2s ease-out;

        &.show-more {
            max-height: 100%;
        }
    }
    .module-item {
        position: relative;
        display: inline-block;
        vertical-align: middle;
        height: 26px;
        width: 120px;
        margin: 0 10px 10px 0;
        line-height: 24px;
        padding: 0 20px 0 25px;
        border: 1px solid #c4c6cc;
        border-radius: 13px;
        color: $textColor;
        font-size: 12px;
        cursor: default;
        @include ellipsis;
        .module-icon {
            position: absolute;
            left: 2px;
            top: 2px;
            width: 20px;
            height: 20px;
            border-radius: 50%;
            line-height: 20px;
            text-align: center;
            color: #fff;
            font-size: 12px;
            background-color: #c4c6cc;
        }

        &.more {
            position: absolute;
            left: 0;
            bottom: 0;
            background: #fafbfd;
            border: 0 none;
            border-radius: unset;
            cursor: pointer;
            color: #3a84ff;
            font-size: 14px;
            text-align: center;
            padding: 0;
            .bk-cc-icon {
                font-size: 22px;
            }

            &.opened {
                position: static;
                .bk-cc-icon {
                    transform: rotate(180deg);
                }
            }
        }
    }
</style>
