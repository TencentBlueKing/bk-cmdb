<template>
    <div :class="['host-apply', { 'show-feature-tips': showFeatureTips }]">
        <feature-tips
            class-name="top-tips"
            feature-name="hostApply"
            :show-tips="showFeatureTips"
            :desc="$t('主机属性自动应用功能提示')"
            @close-tips="showFeatureTips = false"
        >
        </feature-tips>
        <div class="main-wrapper">
            <cmdb-resize-layout class="tree-layout fl"
                direction="right"
                :handler-offset="3"
                :min="300"
                :max="480"
            >
                <sidebar ref="sidebar" @module-selected="handleSelectModule" @action-change="handleActionChange"></sidebar>
            </cmdb-resize-layout>
            <div class="main-content">
                <template v-if="moduleId">
                    <div class="config-panel" v-show="!batchAction">
                        <div class="config-head">
                            <h2 class="config-title">
                                <span class="module-name">{{currentModule.bk_inst_name}}</span>
                                <small class="last-edit-time" v-if="hasRule">( {{$t('上次编辑时间')}}：{{ruleLastEditTime}} )</small>
                            </h2>
                        </div>
                        <div class="config-body">
                            <template v-if="applyEnabled">
                                <div class="view-field">
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
                                            <span v-bk-tooltips="{ content: $t('无失效需处理') }" v-if="!hasConflict">
                                                {{$t('失效列表')}}<em class="conflict-num">{{conflictNum}}</em>
                                            </span>
                                            <span v-else>
                                                {{$t('失效列表')}}<em class="conflict-num">{{conflictNum}}</em>
                                            </span>
                                        </bk-button>
                                        <bk-button theme="default" @click="handleCloseApply">{{$t('关闭自动应用')}}</bk-button>
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
                    </div>
                </template>
                <div class="empty" v-else>
                    <div class="desc">
                        <i class="bk-cc-icon icon-cc-tips"></i>
                        <span>{{$t('暂无业务模块，无法启用属性自动应用，可跳转业务拓扑创建')}}</span>
                    </div>
                    <div class="action">
                        <bk-button
                            outline
                            theme="primary"
                            @click="$router.push({ name: hostAndServiceRouteName })"
                        >
                            {{$t('跳转创建')}}
                        </bk-button>
                    </div>
                </div>
            </div>
        </div>
    </div>
</template>

<script>
    import { mapGetters } from 'vuex'
    import featureTips from '@/components/feature-tips/index'
    import sidebar from './children/sidebar.vue'
    import propertyConfigTable from './children/property-config-table'
    import { MENU_BUSINESS_HOST_AND_SERVICE } from '@/dictionary/menu-symbol'
    export default {
        components: {
            sidebar,
            featureTips,
            propertyConfigTable
        },
        data () {
            return {
                currentModule: {},
                initRuleList: [],
                checkedPropertyIdList: [],
                conflictNum: 0,
                clearRules: false,
                hasRule: false,
                showFeatureTips: false,
                batchAction: false,
                hostAndServiceRouteName: MENU_BUSINESS_HOST_AND_SERVICE
            }
        },
        computed: {
            ...mapGetters('objectBiz', ['bizId']),
            ...mapGetters(['featureTipsParams']),
            applyEnabled () {
                return this.currentModule.host_apply_enabled
            },
            moduleId () {
                return this.currentModule.bk_inst_id
            },
            ruleLastEditTime () {
                const lastTimeList = this.initRuleList.map(rule => new Date(rule.last_time).getTime())
                const latestTime = Math.max(...lastTimeList)
                return this.$tools.formatTime(latestTime, 'YYYY-MM-DD HH:mm:ss')
            },
            hasConflict () {
                return this.conflictNum > 0
            }
        },
        created () {
            this.showFeatureTips = this.featureTipsParams['hostApply']
            this.getHostPropertyList()
        },
        methods: {
            async getData () {
                try {
                    const ruleData = await this.getRules()

                    // 重置配置表格数据
                    if (this.$refs.propertyConfigTable) {
                        this.$refs.propertyConfigTable.reset()
                    }

                    this.initRuleList = ruleData.info || []
                    this.hasRule = ruleData.count > 0
                    this.checkedPropertyIdList = this.initRuleList.map(item => item.bk_attribute_id)

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
            async getHostPropertyList () {
                try {
                    const properties = await this.$store.dispatch('hostApply/getProperties', {
                        requestId: 'getHostPropertyList',
                        fromCache: true
                    })
                    this.$store.commit('hostApply/setPropertyList', properties)
                } catch (e) {
                    console.error(e)
                }
            },
            emptyRules () {
                this.checkedPropertyIdList = []
                this.hasRule = false
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
                            this.$refs.sidebar.setApplyClosed(this.moduleId, this.clearRules)
                        } catch (e) {
                            console.log(e)
                        }
                    }
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
            handleEdit () {
                this.$router.push({
                    name: 'hostApplyEdit',
                    query: {
                        mid: this.moduleId
                    }
                })
            },
            handleSelectModule (data) {
                this.currentModule = data
                this.getData()
            },
            handleActionChange (action) {
                this.batchAction = action
            }
        }
    }
</script>

<style lang="scss" scoped>
    .host-apply {
        .top-tips {
            margin: 0 20px 10px;
        }

        &.show-feature-tips {
            .main-wrapper {
                height: calc(100% - 42px);
            }
        }
    }
    .main-wrapper {
        height: 100%;
        border-top: 1px solid $cmdbLayoutBorderColor;

    }
    .tree-layout {
        width: 300px;
        height: 100%;
        border-right: 1px solid $cmdbLayoutBorderColor;
    }
    .main-content {
        @include scrollbar-y;
        height: 100%;
        padding: 0 20px;

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
    }

    .config-panel {
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

        .view-field {
            .field-list {
                position: relative;

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
                }
            }
            .view-bd,
            .view-ft {
                margin: 20px 0;
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
