<template>
    <div class="model-field-wrapper">
        <div class="form-content clearfix">
            <h3>{{$t('ModelManagement["字段配置"]')}}</h3>
            <div class="form-item has-right-content">
                <label class="form-label">{{$t('ModelManagement["中文名称"]')}}<span class="color-danger">*</span></label>
                <div class="input-box">
                    <input type="text" class="cmdb-form-input"
                        v-model.trim="fieldInfo['bk_property_name']"
                        v-validate="'required|enumName'"
                        :disabled="isReadOnly"
                        name="fieldName">
                    <div v-show="errors.has('fieldName')" class="error-msg color-danger">{{ errors.first('fieldName') }}</div>
                </div>
            </div>
            <div class="form-item has-right-content">
                <label class="form-label">{{$t('ModelManagement["英文名称"]')}}<span class="color-danger">*</span></label>
                <div class="input-box">
                    <input type="text" class="cmdb-form-input"
                        v-model.trim="fieldInfo['bk_property_id']"
                        v-validate="'required|fieldId'"
                        :disabled="isEditField"
                        name="fieldId">
                    <div v-show="errors.has('fieldId')" class="error-msg color-danger">{{ errors.first('fieldId') }}</div>
                </div>
            </div>
            <div class="form-item">
                <label class="form-label unit">{{$t('ModelManagement["单位"]')}}</label>
                <input type="text" class="cmdb-form-input" v-model.trim="fieldInfo['unit']" :disabled="isReadOnly">
            </div>
            <div class="form-item block">
                <label class="form-label">{{$t('ModelManagement["提示语"]')}}</label>
                <input type="text" class="cmdb-form-input" v-model.trim="fieldInfo['placeholder']" :disabled="isReadOnly">
            </div>
        </div>
        <div class="form-content" :class="{'details': fieldInfo['bk_property_type'] === 'enum'}">
            <h3>{{$t('ModelManagement["选项"]')}}</h3>
            <div class="clearfix">
                <div class="form-item has-right-content">
                    <label class="form-label">{{$t('ModelManagement["类型"]')}}</label>
                    <bk-selector
                        class="form-selector"
                        :list="fieldTypeList"
                        :content-max-height="200"
                        :selected.sync="fieldInfo['bk_property_type']"
                        :disabled="isEditField"
                    ></bk-selector>
                </div>
                <v-config :type="fieldInfo['bk_property_type']"
                    :isReadOnly="isReadOnly"
                    :editable.sync="fieldInfo['editable']"
                    :isrequired.sync="fieldInfo['isrequired']"
                    :isonly.sync="fieldInfo['isonly']"></v-config>
            </div>
            <div class="field-config clearfix" v-if="isComponentShow">
                <component v-if="fieldType !== 'asst'"
                    :is="`model-field-${fieldType}`"
                    v-model="fieldInfo.option"
                    :isReadOnly="isReadOnly"
                    ref="component"
                ></component>
                <component v-else
                    :is="`model-field-${fieldType}`"
                    :isEditField="isEditField"
                    v-model="fieldInfo['bk_asst_obj_id']"
                    :isReadOnly="isReadOnly"
                    ref="component"
                ></component>
            </div>
            <div class="btn-wrapper" v-if="!isReadOnly">
                <bk-button type="primary" :loading="$loading(['createObjectAttribute', 'updateObjectAttribute'])" @click="save">
                    {{$t('Common["保存"]')}}
                </bk-button>
                <bk-button type="default" @click="cancel">
                    {{$t('Common["取消"]')}}
                </bk-button>
            </div>
        </div>
    </div>
</template>

<script>
    import modelFieldChar from './char'
    import modelFieldInt from './int'
    import modelFieldEnum from './enum'
    import modelFieldAsst from './asst'
    import vConfig from './config'
    import { mapGetters, mapActions } from 'vuex'
    export default {
        components: {
            modelFieldChar,
            modelFieldInt,
            modelFieldEnum,
            modelFieldAsst,
            vConfig
        },
        props: {
            field: {
                type: Object
            },
            isReadOnly: {
                type: Boolean,
                default: false
            },
            isEditField: {
                type: Boolean,
                default: false
            }
        },
        data () {
            return {
                fieldTypeList: [{
                    id: 'singlechar',
                    name: this.$t('ModelManagement["短字符"]')
                }, {
                    id: 'int',
                    name: this.$t('ModelManagement["数字"]')
                }, {
                    id: 'enum',
                    name: this.$t('ModelManagement["枚举"]')
                }, {
                    id: 'date',
                    name: this.$t('ModelManagement["日期"]')
                }, {
                    id: 'time',
                    name: this.$t('ModelManagement["时间"]')
                }, {
                    id: 'longchar',
                    name: this.$t('ModelManagement["长字符"]')
                }, {
                    id: 'singleasst',
                    name: this.$t('ModelManagement["单关联"]')
                }, {
                    id: 'multiasst',
                    name: this.$t('ModelManagement["多关联"]')
                }, {
                    id: 'objuser',
                    name: this.$t('ModelManagement["用户"]')
                }, {
                    id: 'timezone',
                    name: this.$t('ModelManagement["时区"]')
                }, {
                    id: 'bool',
                    name: 'bool'
                }],
                fieldInfo: {
                    bk_property_name: '',
                    bk_property_id: '',
                    unit: '',
                    placeholder: '',
                    bk_property_type: 'singlechar',
                    editable: true,
                    isrequired: false,
                    isonly: false,
                    option: '',
                    bk_asst_obj_id: ''
                },
                charMap: ['singlechar', 'longchar'],
                asstMap: ['singleasst', 'multiasst']
            }
        },
        computed: {
            ...mapGetters([
                'userName',
                'supplierAccount'
            ]),
            ...mapGetters('objectModel', [
                'activeModel'
            ]),
            isComponentShow () {
                return ['singlechar', 'longchar', 'multichar', 'singleasst', 'multiasst', 'enum', 'int'].indexOf(this.fieldInfo['bk_property_type']) !== -1
            },
            fieldType () {
                let {
                    bk_property_type: type
                } = this.fieldInfo
                if (this.charMap.indexOf(type) !== -1) {
                    return 'char'
                } else if (this.asstMap.indexOf(type) !== -1) {
                    return 'asst'
                }
                return type
            }
        },
        watch: {
            'fieldInfo.bk_property_type' (type) {
                if (!this.isEditField) {
                    switch (type) {
                        case 'int':
                            this.fieldInfo.option = {
                                min: '',
                                max: ''
                            }
                            this.fieldInfo.bk_asst_obj_id = ''
                            break
                        default:
                            this.fieldInfo.option = ''
                            this.fieldInfo.bk_asst_obj_id = ''
                    }
                }
            }
        },
        methods: {
            ...mapActions('objectModelProperty', [
                'createObjectAttribute',
                'updateObjectAttribute'
            ]),
            isCloseConfirmShow () {
                let result = false
                for (let key in this.fieldInfo) {
                    if (this.fieldInfo[key] !== this.field[key]) {
                        result = true
                        break
                    }
                }
                return result
            },
            initData () {
                for (let key in this.fieldInfo) {
                    this.fieldInfo[key] = this.$tools.clone(this.field[key])
                }
            },
            async validateValue () {
                let res = await this.$validator.validateAll()
                if (!res) {
                    return false
                }
                if (this.$refs.component && this.$refs.component.hasOwnProperty('validate')) {
                    res = await this.$refs.component.validate()
                    if (!res) {
                        return false
                    }
                }
                return true
            },
            async save () {
                if (!await this.validateValue()) {
                    return
                }
                
                if (this.isEditField) {
                    await this.updateObjectAttribute({
                        id: this.field.id,
                        params: this.fieldInfo,
                        config: {
                            requestId: 'updateObjectAttribute'
                        }
                    }).then(() => {
                        this.$http.cancel(`post_searchObjectAttribute_${this.activeModel['bk_obj_id']}`)
                    })
                } else {
                    let otherParams = {
                        creator: this.userName,
                        bk_property_group: 'default',
                        bk_obj_id: this.activeModel['bk_obj_id'],
                        bk_supplier_account: this.supplierAccount
                    }
                    await this.createObjectAttribute({
                        params: {...this.fieldInfo, ...otherParams},
                        config: {
                            requestId: 'createObjectAttribute'
                        }
                    }).then(() => {
                        this.$http.cancel(`post_searchObjectAttribute_${this.activeModel['bk_obj_id']}`)
                    })
                }
                this.$emit('save', this.fieldInfo['bk_property_type'])
            },
            cancel () {
                this.$emit('cancel')
            }
        },
        mounted () {
            if (this.isEditField) {
                this.initData()
            }
        }
    }
</script>


<style lang="scss" scoped>
    .model-field-wrapper {
        height: 100%;
        padding: 20px 20px 10px;
        .form-content {
            margin-bottom: 20px;
            &.details {
                height: calc(100% - 134px);
                .field-config {
                    max-height: calc(100% - 154px);
                    @include scrollbar;
                }
            }
            h3 {
                margin-bottom: 10px;
                font-size: 14px;
                line-height: 1;
            }
            .form-item {
                &.block {
                    margin-top: 10px;
                    .cmdb-form-input {
                        width: 625px;
                    }
                }
                .input-box {
                    display: inline-block;
                }
                .unit {
                    width: 28px;
                }
            }
            .field-config {
                margin-top: 10px;
            }
            .btn-wrapper {
                margin-top: 20px;
                padding-left: 80px;
                font-size: 0px;
                button {
                    &:first-child {
                        margin-right: 10px;
                    }
                }
            }
        }
    }
</style>
