<template>
    <div class="single-module-config">
        <div class="config-head">
            <h2 class="config-title">
                <i18n path="配置XXX模块主机属性">
                    <span class="module-name" place="module">{{data.bk_inst_name}}</span>
                </i18n>
                <small class="last-edit-time" v-if="hasConfig">( {{$t('上次编辑时间')}}：2017.09.13 )</small>
            </h2>
        </div>
        <div class="config-body">
            <template v-if="hasConfig && !isEdit">
                <div class="view-field">
                    <div class="view-bd">
                        <div class="field-list">
                            <div :class="['field-list-table', { disabled: !isConfigEnabled }]">
                                <property-config-table
                                    :readonly="true"
                                    :checked-property-id-list="checkedPropertyIdList"
                                    :init-data="initConfigData.config_list"
                                >
                                </property-config-table>
                            </div>
                            <div class="closed-mask" v-if="!isConfigEnabled">
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
                    <div class="view-ft" v-if="isConfigEnabled">
                        <bk-button theme="primary" @click="handleEdit">{{$t('编辑')}}</bk-button>
                        <bk-button theme="default" :disabled="!hasConflict" @click="handleViewConflict">
                            <span v-bk-tooltips="{ content: $t('无冲突需处理') }" v-if="!hasConflict">
                                {{$t('查看冲突')}}<em class="conflict-num">{{initConfigData.conflict_num}}</em>
                            </span>
                            <span v-else>
                                {{$t('查看冲突')}}<em class="conflict-num">{{initConfigData.conflict_num}}</em>
                            </span>
                        </bk-button>
                        <bk-button theme="default" @click="handleCloseApply">{{$t('关闭自动应用')}}</bk-button>
                    </div>
                </div>
            </template>
            <template v-else>
                <div class="empty" v-if="!isEdit">
                    <div class="desc">
                        <i class="bk-cc-icon icon-cc-tips"></i>
                        <span>{{$t('尚未配置模块的属性自动应用')}}</span>
                    </div>
                    <div class="action">
                        <bk-button theme="primary" :outline="true" @click="handleEdit">立即启用</bk-button>
                    </div>
                </div>
                <div :class="['choose-field', { 'not-choose': !checkedPropertyIdList.length }]" v-else>
                    <div class="choose-hd">
                        <span class="label">{{$t('请添加自动应用的字段')}}</span>
                        <bk-button theme="default" icon="plus" @click="handleChooseField">添加字段</bk-button>
                    </div>
                    <div class="choose-bd" v-if="checkedPropertyIdList.length">
                        <property-config-table
                            ref="configEditTable"
                            :checked-property-id-list="checkedPropertyIdList"
                            :init-data="initConfigData.config_list"
                            @property-value-change="handlePropertyValueChange"
                        >
                        </property-config-table>
                    </div>
                    <div class="choose-ft">
                        <bk-button theme="primary" :disabled="nextDisabled" @click="handleNextStep">下一步</bk-button>
                        <bk-button theme="default" @click="handleCancel">取消</bk-button>
                    </div>
                </div>
            </template>
        </div>
        <host-property-modal
            :visiable.sync="propertyModalVisiable"
            :checked-list.sync="checkedPropertyIdList"
        >
        </host-property-modal>
    </div>
</template>

<script>
    import { mapGetters } from 'vuex'
    import hostPropertyModal from './host-property-modal'
    import propertyConfigTable from './property-config-table'
    export default {
        components: {
            hostPropertyModal,
            propertyConfigTable
        },
        props: {
            data: {
                type: Object,
                default: () => ({})
            }
        },
        data () {
            return {
                isEdit: false,
                nextDisabled: false,
                propertyModalVisiable: false,
                initConfigData: {},
                checkedPropertyIdList: []
            }
        },
        computed: {
            ...mapGetters('hosts', ['configPropertyList']),
            hasConfig () {
                return true || this.data.host_config_id
            },
            isConfigEnabled () {
                return true || this.data.enabled
            },
            hasConflict () {
                return this.initConfigData.conflict_num > 0
            }
        },
        watch: {
            configPropertyList: {
                handler () {
                    this.initData()
                },
                immediate: true
            },
            checkedPropertyIdList () {
                this.$nextTick(() => {
                    this.isEidt && this.toggleNexButtonDisabled()
                })
            }
        },
        created () {
        },
        methods: {
            async initData () {
                // 请求接口获取配置数据
                if (this.hasConfig) {
                    this.initConfigData = await Promise.resolve({
                        config_list: [
                            {
                                bk_property_id: 'bk_asset_id',
                                value: 'a12356'
                            },
                            {
                                bk_property_id: 'bk_state_name',
                                value: 'BR'
                            }
                        ],
                        conflict_num: 12
                    })

                    this.checkedPropertyIdList = this.initConfigData.config_list.map(item => item.bk_property_id)
                }
            },
            toggleNexButtonDisabled () {
                const modulePropertyList = this.$refs.configEditTable.modulePropertyList
                this.nextDisabled = !this.checkedPropertyIdList.length || !modulePropertyList.every(property => property.__extra__.value)
            },
            handleNextStep () {
                const modulePropertyList = this.$refs.configEditTable.modulePropertyList
                const configData = modulePropertyList.map(property => ({
                    key: property.bk_property_id,
                    value: property.__extra__.value
                }))

                // TODO 调用接口保存配置

                this.$router.push({
                    name: 'hostApplyConfirm',
                    query: {
                        batch: 0,
                        cid: 1
                    }
                })
                console.log(configData)
            },
            handleViewConflict () {
                this.$router.push({
                    name: 'hostApplyConflict',
                    query: {
                        cid: 1
                    }
                })
            },
            handleCloseApply () {
                this.$bkInfo({
                    title: this.$t('确认关闭自动应用？'),
                    subTitle: this.$t('关闭后，新转入的的主机将不会自动应用模块的主机属性'),
                    confirmFn: () => {
                    }
                })
            },
            handlePropertyValueChange () {
                this.toggleNexButtonDisabled()
            },
            handleEdit () {
                this.isEdit = true
            },
            handleChooseField () {
                this.propertyModalVisiable = true
            },
            handleCancel (id) {
                this.isEdit = false
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

        .module-name {
            padding: 0 .2em;
        }
        .last-edit-time {
            font-size: 12px;
            font-weight: 400;
            color: #979ba5;
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
                margin-left: 152px;
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
</style>
