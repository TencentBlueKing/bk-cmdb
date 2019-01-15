<template>
    <div class="slider-content">
        <label class="form-label">
            <span class="label-text">
                {{$t('ModelManagement["唯一标识"]')}}
                <span class="color-danger">*</span>
            </span>
            <div class="cmdb-form-item" :class="{'is-error': errors.has('fieldId')}">
                <input type="text" class="cmdb-form-input"
                name="fieldId"
                :placeholder="$t('ModelManagement[\'下划线/数字/字母\']')"
                v-model.trim="fieldInfo['bk_property_id']"
                :disabled="isEditField"
                v-validate="'required|fieldId'">
                <p class="form-error">{{errors.first('fieldId')}}</p>
            </div>
        </label>
        <label class="form-label">
            <span class="label-text">
                {{$t('ModelManagement["名称"]')}}
                <span class="color-danger">*</span>
            </span>
            <div class="cmdb-form-item" :class="{'is-error': errors.has('fieldName')}">
                <input type="text" class="cmdb-form-input"
                name="fieldName"
                :placeholder="$t('ModelManagement[\'请输入字段名称\']')"
                v-model.trim="fieldInfo['bk_property_name']"
                :disabled="isReadOnly"
                v-validate="'required|enumName'">
                <p class="form-error">{{errors.first('fieldName')}}</p>
            </div>
        </label>
        <div class="form-label">
            <span class="label-text">
                {{$t('ModelManagement["字段类型"]')}}
                <span class="color-danger">*</span>
            </span>
            <div class="cmdb-form-item">
                <bk-selector
                    :disabled="isEditField"
                    :list="fieldTypeList"
                    :selected.sync="fieldInfo['bk_property_type']"
                ></bk-selector>
            </div>
        </div>
        <div class="field-detail">
            <the-config
                :type="fieldInfo['bk_property_type']"
                :isReadOnly="isReadOnly"
                :editable.sync="fieldInfo['editable']"
                :isrequired.sync="fieldInfo['isrequired']"
            ></the-config>
            <component 
                v-if="isComponentShow"
                :isReadOnly="isReadOnly"
                :is="`the-field-${fieldType}`"
                v-model="fieldInfo.option"
                ref="component"
            ></component>
        </div>
        <label class="form-label">
            <span class="label-text">
                {{$t('ModelManagement["单位"]')}}
            </span>
            <div class="cmdb-form-item">
                <input type="text" class="cmdb-form-input"
                v-model.trim="fieldInfo['unit']"
                :disabled="isReadOnly"
                :placeholder="$t('ModelManagement[\'请输入单位\']')">
            </div>
        </label>
        <div class="form-label">
            <span class="label-text">{{$t('ModelManagement["用户提示"]')}}</span>
            <textarea v-model.trim="fieldInfo['placeholder']" :disabled="isReadOnly"></textarea>
        </div>
        <div class="btn-group">
            <bk-button type="primary" :loading="$loading(['updateObjectAttribute', 'createObjectAttribute'])" @click="saveField">
                {{$t('Common["确定"]')}}
            </bk-button>
            <bk-button type="default" @click="cancel">
                {{$t('Common["取消"]')}}
            </bk-button>
        </div>
    </div>
</template>

<script>
    import theFieldChar from './char'
    import theFieldInt from './int'
    import theFieldFloat from './float'
    import theFieldEnum from './enum'
    import theConfig from './config'
    import { mapGetters, mapActions } from 'vuex'
    export default {
        components: {
            theFieldChar,
            theFieldInt,
            theFieldFloat,
            theFieldEnum,
            theConfig
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
                    id: 'float',
                    name: this.$t('ModelManagement["浮点"]')
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
                    option: ''
                },
                charMap: ['singlechar', 'longchar']
            }
        },
        computed: {
            ...mapGetters(['supplierAccount', 'userName']),
            ...mapGetters('objectModel', [
                'activeModel'
            ]),
            fieldType () {
                let {
                    bk_property_type: type
                } = this.fieldInfo
                if (this.charMap.indexOf(type) !== -1) {
                    return 'char'
                }
                return type
            },
            isComponentShow () {
                return ['singlechar', 'longchar', 'enum', 'int', 'float'].indexOf(this.fieldInfo['bk_property_type']) !== -1
            }
        },
        watch: {
            'fieldInfo.bk_property_type' (type) {
                if (!this.isEditField) {
                    switch (type) {
                        case 'int':
                        case 'float':
                            this.fieldInfo.option = {
                                min: '',
                                max: ''
                            }
                            break
                        default:
                            this.fieldInfo.option = ''
                    }
                }
            }
        },
        created () {
            if (this.isEditField) {
                this.initData()
            }
        },
        methods: {
            ...mapActions('objectModelProperty', [
                'createObjectAttribute',
                'updateObjectAttribute'
            ]),
            initData () {
                for (let key in this.fieldInfo) {
                    this.fieldInfo[key] = this.$tools.clone(this.field[key])
                }
            },
            async validateValue () {
                if (!await this.$validator.validateAll()) {
                    return false
                }
                if (this.$refs.component && this.$refs.component.hasOwnProperty('validate')) {
                    if (!await this.$refs.component.validate()) {
                        return false
                    }
                }
                return true
            },
            async saveField () {
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
                this.$emit('save')
            },
            cancel () {
                this.$emit('cancel')
            }
        }
    }
</script>

<style lang="scss" scoped>
    .icon-info-circle {
        font-size: 18px;
        color: $cmdbBorderColor;
        padding-left: 5px;
    }
    .field-detail {
        margin-bottom: 20px;
        padding: 20px;
        background: #f3f8ff;
        .form-label:last-child {
            margin: 0;
        }
        .label-text {
            vertical-align: top;
        }
        .cmdb-form-checkbox {
            width: 90px;
            line-height: 22px;
            vertical-align: middle;
        }
    }
</style>
