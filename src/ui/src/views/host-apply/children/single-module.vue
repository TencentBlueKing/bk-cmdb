<template>
    <div class="single-module-panel">
        <div class="panel-head">
            <h2 class="panel-title">
                <i18n path="配置XXX模块主机属性">
                    <span class="module-name" place="module">{{data.bk_inst_name}}</span>
                </i18n>
                <small class="last-edit-time">( {{$t('上次编辑时间')}}：2017.09.13 )</small>
            </h2>
        </div>
        <div class="panel-body">
            <template v-if="hasConfig && !isEdit">
                <div class="view-field">
                    <div class="view-bd">
                        <div class="field-list">
                            <div class="field-list-table">
                                <bk-table
                                    :data="modulePropertyList"
                                >
                                    <bk-table-column width="200" :label="$t('字段名称')" prop="bk_property_name"></bk-table-column>
                                    <bk-table-column :label="$t('值')" class-name="table-cell-form-element">
                                        <template slot-scope="{ row }">
                                            <span>{{getDisplayValue(row)}}</span>
                                        </template>
                                    </bk-table-column>
                                </bk-table>
                            </div>
                            <div class="closed-mask" v-if="!isConfigEnabled">
                                <div class="empty">
                                    <div class="desc">
                                        <i class="bk-cc-icon icon-cc-tips"></i>
                                        <span>{{$t('该模块已关闭属性自动应用')}}</span>
                                    </div>
                                    <div class="action">
                                        <bk-button theme="primary" :outline="true" @click="handleEdit">重新启用</bk-button>
                                    </div>
                                </div>
                            </div>
                        </div>
                    </div>
                    <div class="view-ft" v-if="isConfigEnabled">
                        <bk-button theme="primary" @click="handleEdit">编辑</bk-button>
                        <bk-button theme="default">查看冲突</bk-button>
                        <bk-button theme="default">关闭自动应用</bk-button>
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
                <div :class="['choose-field', { 'not-choose': !modulePropertyList.length }]" v-else>
                    <div class="choose-hd">
                        <span class="label">{{$t('请添加自动应用的字段')}}</span>
                        <bk-button theme="default" icon="plus" @click="handleChooseField">添加字段</bk-button>
                    </div>
                    <div class="choose-bd" v-if="modulePropertyList.length">
                        <bk-table
                            :data="modulePropertyList"
                        >
                            <bk-table-column width="200" :label="$t('字段名称')" prop="bk_property_name"></bk-table-column>
                            <bk-table-column :label="$t('值')" class-name="table-cell-form-element">
                                <template slot-scope="{ row }">
                                    <div class="form-element-content">
                                        <property-form-element :property="row" @value-change="handlePropertyValueChange"></property-form-element>
                                    </div>
                                </template>
                            </bk-table-column>
                            <bk-table-column width="200" :label="$t('操作')">
                                <template slot-scope="{ row }">
                                    <bk-button theme="primary" text @click="handlePropertyRowDel(row.bk_property_id)">删除</bk-button>
                                </template>
                            </bk-table-column>
                        </bk-table>
                    </div>
                    <div class="choose-ft">
                        <bk-button theme="primary" :disabled="nextDisabled" @click="handleNextStep">下一步</bk-button>
                        <bk-button theme="default">取消</bk-button>
                    </div>
                </div>
            </template>
        </div>
        <div class="panel-foot">

        </div>
        <host-property-modal
            :property-list="configPropertyList"
            :visiable.sync="propertyModalVisiable"
            :checked-list.sync="checkedPropertyIdList"
        >
        </host-property-modal>
    </div>
</template>

<script>
    import { mapGetters } from 'vuex'
    import hostPropertyModal from './host-property-modal'
    import propertyFormElement from './property-form-element'
    export default {
        components: {
            hostPropertyModal,
            propertyFormElement
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
                nextDisabled: true,
                propertyModalVisiable: false,
                initConfigData: [],
                modulePropertyList: [],
                checkedPropertyIdList: []
            }
        },
        computed: {
            ...mapGetters('hosts', ['configPropertyList']),
            hasConfig () {
                return true || this.data.host_config_id
            },
            isConfigEnabled () {
                return false || this.data.enabled
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
                this.setModulePropertyList()
            },
            modulePropertyList () {
                this.toggleNexButtonDisabled()
            }
        },
        created () {
        },
        methods: {
            async initData () {
                // 请求接口获取配置数据
                this.initConfigData = await Promise.resolve([
                    {
                        bk_property_id: 'bk_asset_id',
                        value: 'a123'
                    },
                    {
                        bk_property_id: 'bk_state_name',
                        value: 'BR'
                    }
                ])

                this.checkedPropertyIdList = this.initConfigData.map(item => item.bk_property_id)
            },
            setModulePropertyList () {
                // 当前模块属性列表中不存在，则添加
                this.checkedPropertyIdList.forEach(propertyId => {
                    if (this.modulePropertyList.findIndex(property => propertyId === property.bk_property_id) === -1) {
                        // 使用原始主机属性对象
                        const findProperty = this.configPropertyList.find(item => propertyId === item.bk_property_id)
                        if (findProperty) {
                            // 初始化值
                            findProperty.__extra__.value = (this.initConfigData.find(item => item.bk_property_id === findProperty.bk_property_id) || {}).value
                            this.modulePropertyList.push(this.$tools.clone(findProperty))
                        }
                    }
                })

                // 删除或取消选择的，则去除
                this.modulePropertyList = this.modulePropertyList.filter(property => this.checkedPropertyIdList.includes(property.bk_property_id))
            },
            getDisplayValue (property) {
                const value = property.__extra__.value
                let displayValue = value
                switch (property.bk_property_type) {
                    case 'enum':
                        displayValue = (property.option.find(item => item.id === value) || {}).name
                        break
                    case 'bool':
                        displayValue = ['否', '是'][+value]
                        break
                }

                return displayValue
            },
            toggleNexButtonDisabled () {
                this.nextDisabled = !this.modulePropertyList.length || !this.modulePropertyList.every(property => property.__extra__.value)
            },
            handleNextStep () {
                const configData = this.modulePropertyList.map(property => ({
                    key: property.bk_property_id,
                    value: property.__extra__.value
                }))

                this.$router.push({
                    name: 'hostApplyConfirm'
                })
                console.log(configData)
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
            handlePropertyRowDel (id) {
                const checkedIndex = this.checkedPropertyIdList.findIndex(propertyId => propertyId === id)
                this.checkedPropertyIdList.splice(checkedIndex, 1)
            }
        }
    }
</script>

<style lang="scss" scoped>
    .single-module-panel {
        display: flex;
        flex-direction: column;
        height: 100%;

        .panel-head,
        .panel-foot {
            flex: none;
        }
        .panel-body {
            flex: auto;
        }
    }
    .panel-title {
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
                opacity: 0.2;
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
        }
    }
</style>

<style lang="scss">
    .table-cell-form-element {
        .cell {
            overflow: unset;
            display: block;
        }
        .search-input-wrapper {
            position: relative;
        }
        .form-objuser .suggestion-list {
            z-index: 1000;
        }
    }
</style>
