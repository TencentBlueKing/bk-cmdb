<template>
    <div class="display-wrapper">
        <div class="display-box">
            <div class="display-setting">
                <label class="cmdb-form-checkbox cmdb-checkbox-small">
                    <input type="checkbox" v-model="isShowName">
                    <span class="cmdb-checkbox-text">{{$t('ModelManagement["显示模型名称"]')}}</span>
                </label>
                <label class="cmdb-form-checkbox cmdb-checkbox-small">
                    <input type="checkbox" v-model="isShowAsst">
                    <span class="cmdb-checkbox-text">{{$t('ModelManagement["显示关联名称"]')}}</span>
                </label>
            </div>
            <ul class="display-list">
                <li class="group-item" v-for="(group, groupIndex) in displayModelGroups" :key="groupIndex">
                    <p class="group-name">
                        <span class="group-name-text">{{group['bk_classification_name']}}</span>
                        <label class="group-name-checkbox cmdb-form-checkbox cmdb-checkbox-small">
                            <input type="checkbox"
                                :checked="isGroupChecked(group)"
                                @change="handleToggleGroup($event, group)">
                        </label>
                    </p>
                    <ul class="clearfix">
                        <li class="model-item" :class="{'active': model['bk_obj_id'] === activePop}" v-for="(model, modelIndex) in group['bk_objects']" :key="modelIndex">
                            <label class="cmdb-form-checkbox checkbox cmdb-checkbox-small">
                                <input type="checkbox" v-model="localModelConfig[model.bk_obj_id]">
                                <span class="cmdb-checkbox-text">{{model['bk_obj_name']}}</span>
                            </label>
                        </li>
                    </ul>
                </li>
            </ul>
        </div>
        <div class="button-group">
            <bk-button type="primary" @click="saveDisplay">
                {{$t('Common["确定"]')}}
            </bk-button>
            <bk-button type="default" @click="reset">
                {{$t('ModelManagement["重置"]')}}
            </bk-button>
        </div>
    </div>
</template>

<script>
    import { mapGetters } from 'vuex'
    export default {
        props: {
            associationList: {
                type: Array
            },
            topoModelList: {
                type: Array
            },
            isShowModelName: {
                type: Boolean
            },
            isShowModelAsst: {
                type: Boolean
            }
        },
        data () {
            return {
                localTopoModelList: [],
                activePop: '',
                isShowName: this.isShowModelName,
                isShowAsst: this.isShowModelAsst,
                localModelConfig: {}
            }
        },
        computed: {
            ...mapGetters('objectModelClassify', [
                'classifications'
            ]),
            ...mapGetters('globalModels', [
                'modelConfig',
                'displayModelGroups'
            ])
        },
        created () {
            this.initLocalTopoModelList()
            this.initModelConfig()
        },
        methods: {
            initModelConfig () {
                this.localModelConfig = this.$tools.clone(this.modelConfig)
            },
            isGroupChecked (group) {
                return !group.bk_objects.some(model => {
                    return !this.localModelConfig[model.bk_obj_id]
                })
            },
            handleToggleGroup (event, group) {
                const isChecked = event.target.checked
                group.bk_objects.forEach(model => {
                    this.localModelConfig[model.bk_obj_id] = isChecked
                })
            },
            hidePop () {
                this.activePop = ''
            },
            toggleActivePop (objId) {
                this.activePop = this.activePop === objId ? '' : objId
            },
            initLocalTopoModelList () {
                this.localTopoModelList = this.$tools.clone(this.topoModelList)
            },
            findCurrentModelAsst (model) {
                return this.localTopoModelList.find(obj => obj['bk_obj_id'] === model['bk_obj_id'] && obj.hasOwnProperty('assts') && obj.assts.length).assts
            },
            isChecked (model) {
                let modelAsst = this.findCurrentModelAsst(model)
                return !modelAsst.some(asst => !asst.checked)
            },
            checkAll (model) {
                let modelAsst = this.findCurrentModelAsst(model)
                modelAsst.forEach(asst => {
                    asst.checked = event.target.checked
                })
                this.$forceUpdate()
            },
            asstLabel (model, asst) {
                let asstModel = this.$allModels.find(model => {
                    return model['bk_obj_id'] === asst['bk_obj_id']
                })
                if (asstModel) {
                    let association = this.associationList.find(({id}) => id === asst['bk_asst_inst_id'])
                    if (association) {
                        if (association['bk_asst_name'].length) {
                            return `${association['bk_asst_name']}->${asstModel['bk_obj_name']}`
                        }
                        return `${association['bk_asst_id']}->${asstModel['bk_obj_name']}`
                    }
                }
            },
            saveDisplay () {
                const modelChecked = Object.values(this.localModelConfig).some(value => value)
                if (!modelChecked) {
                    this.$warn(this.$t('ModelManagement["模型配置提示"]'))
                    return false
                }
                const displayConfig = {
                    isShowModelName: this.isShowName,
                    isShowModelAsst: this.isShowAsst,
                    topoModelList: this.$tools.clone(this.localTopoModelList)
                }
                const topoModelConfig = {
                    ...this.localModelConfig,
                    isShowModelName: this.isShowName,
                    isShowModelAsst: this.isShowAsst
                }
                this.$store.dispatch('userCustom/saveUsercustom', {
                    topoModelConfig
                }).then(() => {
                    this.$emit('save', displayConfig)
                    this.$emit('cancel')
                })
            },
            reset () {
                this.isShowName = true
                this.isShowAsst = true
                Object.keys(this.localModelConfig).forEach(modelId => {
                    this.localModelConfig[modelId] = false
                })
                // this.localTopoModelList.forEach(model => {
                //     if (model.hasOwnProperty('assts') && model.assts.length) {
                //         model.assts.forEach(asst => {
                //             asst.checked = true
                //         })
                //     }
                // })
            }
        }
    }
</script>

<style lang="scss" scoped>
    .display-wrapper {
        padding: 20px 0;
        height: 100%;
        .display-box {
            padding: 0 30px;
            max-height: calc(100% - 56px);
            @include scrollbar;
        }
    }
    .display-setting {
        .cmdb-form-checkbox {
            min-width: 154px;
        }
    }
    .display-list {
        .group-item {
            position: relative;
            margin-top: 25px;
            .group-name {
                margin-bottom: 15px;
                padding-left: 8px;
                border-left: 4px solid $cmdbBorderColor;
                line-height: 14px;
                font-size: 0;
                &:hover .group-name-checkbox {
                    display: inline-block;
                }
                .group-name-text {
                    display: inline-block;
                    vertical-align: middle;
                    color: #333948;
                    font-weight: bold;
                    font-size: 14px;
                }
                .group-name-checkbox {
                    display: none;
                    padding: 0;
                    margin: 0 0 0 4px;
                    vertical-align: middle;
                    line-height: 1;
                }
            }
        }
        .model-item {
            float: left;
            width: 175px;
            &.active {
                .count {
                    &:after {
                        position: absolute;
                        content: '';
                        border: 5px solid transparent;
                        border-bottom-color: #fff;
                        top: 17px;
                        left: calc(50% - 5px);
                        z-index: 2;
                    }
                    &:before {
                        position: absolute;
                        content: '';
                        border: 5px solid transparent;
                        border-bottom-color: $cmdbTableBorderColor;
                        top: 16px;
                        left: calc(50% - 5px);
                        z-index: 2;
                    }
                }
            }
            >.cmdb-form-checkbox {
                vertical-align: middle;
                margin-right: 0;
                &.checkbox {
                    font-size: 0;
                }
                &.text-box {
                    margin-right: 10px;
                    cursor: pointer;
                    .cmdb-checkbox-text {
                        max-width: 100px;
                    }
                }
            }
            .count {
                position: relative;
                font-size: 14px;
                vertical-align: middle;
                color: $cmdbBorderFocusColor;
                
            }
            .relation-detail {
                position: absolute;
                margin-top: 7px;
                padding: 5px 20px 10px;
                left: 0;
                width: 100%;
                background: #fff;
                border: 1px solid $cmdbTableBorderColor;
                z-index: 1;
                .detail-title {
                    font-size: 0;
                    line-height: 36px;
                    color: $cmdbBorderColor;
                    .title {
                        color: #333948;
                    }
                    .title,
                    .info {
                        font-size: 14px;
                        vertical-align: middle;
                    }
                    .cmdb-checkbox-text {
                        color: $cmdbBorderColor;
                    }
                }
            }
            .relation-list {
                font-size: 0;
                >li:nth-child(3n) {
                    .cmdb-form-checkbox {
                        margin: 0;
                    }
                }
                .cmdb-form-checkbox {
                    width: 160px;
                    margin-right: 9px;
                    .cmdb-checkbox-text {
                        max-width: 135px;
                        @include ellipsis;
                    }
                    &:nth-child(3n) {
                        margin: 0;
                    }
                }
            }
            .icon-angle-down {
                margin-left: 5px;
                font-size: 12px;
                color: $cmdbBorderColor;
            }
        }
    }
    .button-group {
        margin: 20px;
        font-size: 0;
        .bk-button {
            margin-right: 10px;
        }
    }
</style>
