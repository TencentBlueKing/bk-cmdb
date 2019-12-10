<template>
    <div class="model-slider-content">
        <label class="form-label">
            <span class="label-text">
                {{$t('唯一标识')}}
                <span class="color-danger">*</span>
            </span>
            <div class="cmdb-form-item" :class="{ 'is-error': errors.has('fieldId') }">
                <bk-input type="text" class="cmdb-form-input"
                    name="fieldId"
                    v-model.trim="fieldInfo['bk_property_id']"
                    :placeholder="$t('请输入唯一标识')"
                    :disabled="isEditField"
                    v-validate="'required|fieldId'">
                </bk-input>
                <p class="form-error">{{$t('唯一标识必须为英文字母、数字和下划线组成')}}</p>
            </div>
            <i class="icon-cc-exclamation-tips" v-bk-tooltips="$t('下划线/数字/字母')"></i>
        </label>
        <label class="form-label">
            <span class="label-text">
                {{$t('名称')}}
                <span class="color-danger">*</span>
            </span>
            <div class="cmdb-form-item" :class="{ 'is-error': errors.has('fieldName') }">
                <bk-input type="text" class="cmdb-form-input"
                    name="fieldName"
                    :placeholder="$t('请输入字段名称')"
                    v-model.trim="fieldInfo['bk_property_name']"
                    :disabled="isReadOnly"
                    v-validate="'required|enumName'">
                </bk-input>
                <p class="form-error">{{errors.first('fieldName')}}</p>
            </div>
        </label>
        <div class="form-label">
            <span class="label-text">
                {{$t('字段类型')}}
                <span class="color-danger">*</span>
            </span>
            <div class="cmdb-form-item">
                <bk-select
                    class="bk-select-full-width"
                    :clearable="false"
                    v-model="fieldInfo.bk_property_type"
                    :disabled="isEditField">
                    <bk-option v-for="(option, index) in fieldTypeList"
                        :key="index"
                        :id="option.id"
                        :name="option.name">
                    </bk-option>
                </bk-select>
            </div>
        </div>
        <div class="field-detail">
            <the-config
                :type="fieldInfo['bk_property_type']"
                :is-read-only="isReadOnly"
                :is-main-line-model="isMainLineModel"
                :editable.sync="fieldInfo['editable']"
                :isrequired.sync="fieldInfo['isrequired']"
            ></the-config>
            <component
                v-if="isComponentShow"
                :is-read-only="isReadOnly"
                :is="`the-field-${fieldType}`"
                v-model="fieldInfo.option"
                ref="component"
            ></component>
        </div>
        <label class="form-label" v-show="['int', 'float'].includes(fieldType)">
            <span class="label-text">
                {{$t('单位')}}
            </span>
            <div class="cmdb-form-item">
                <bk-input type="text" class="cmdb-form-input"
                    v-model.trim="fieldInfo['unit']"
                    :disabled="isReadOnly"
                    :placeholder="$t('请输入单位')">
                </bk-input>
            </div>
        </label>
        <div class="form-label">
            <span class="label-text">{{$t('用户提示')}}</span>
            <textarea style="width: 94%;" v-model.trim="fieldInfo['placeholder']" :disabled="isReadOnly"></textarea>
        </div>
        <div class="btn-group">
            <bk-button theme="primary"
                :loading="$loading(['updateObjectAttribute', 'createObjectAttribute'])"
                @click="saveField">
                {{isEditField ? $t('保存') : $t('提交')}}
            </bk-button>
            <bk-button theme="default" @click="cancel">
                {{$t('取消')}}
            </bk-button>
        </div>
    </div>
</template>

<script>
    import theFieldChar from './char'
    import theFieldInt from './int'
    import theFieldFloat from './float'
    import theFieldEnum from './enum'
    import theFieldList from './list'
    import theConfig from './config'
    import { mapGetters, mapActions } from 'vuex'
    export default {
        components: {
            theFieldChar,
            theFieldInt,
            theFieldFloat,
            theFieldEnum,
            theFieldList,
            theConfig
        },
        props: {
            field: {
                type: Object
            },
            group: {
                type: Object
            },
            isReadOnly: {
                type: Boolean,
                default: false
            },
            isEditField: {
                type: Boolean,
                default: false
            },
            propertyIndex: {
                type: Number,
                default: 0
            },
            isMainLineModel: {
                type: Boolean,
                default: false
            }
        },
        data () {
            return {
                fieldTypeList: [{
                    id: 'singlechar',
                    name: this.$t('短字符')
                }, {
                    id: 'int',
                    name: this.$t('数字')
                }, {
                    id: 'float',
                    name: this.$t('浮点')
                }, {
                    id: 'enum',
                    name: this.$t('枚举')
                }, {
                    id: 'date',
                    name: this.$t('日期')
                }, {
                    id: 'time',
                    name: this.$t('时间')
                }, {
                    id: 'longchar',
                    name: this.$t('长字符')
                }, {
                    id: 'objuser',
                    name: this.$t('用户')
                }, {
                    id: 'timezone',
                    name: this.$t('时区')
                }, {
                    id: 'bool',
                    name: 'bool'
                }, {
                    id: 'list',
                    name: this.$t('列表')
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
                originalFieldInfo: {},
                charMap: ['singlechar', 'longchar']
            }
        },
        computed: {
            ...mapGetters(['supplierAccount', 'userName', 'isAdminView']),
            ...mapGetters('objectModel', [
                'activeModel',
                'isPublicModel',
                'isInjectable'
            ]),
            fieldType () {
                const {
                    bk_property_type: type
                } = this.fieldInfo
                if (this.charMap.indexOf(type) !== -1) {
                    return 'char'
                }
                return type
            },
            isComponentShow () {
                return ['singlechar', 'longchar', 'enum', 'int', 'float', 'list'].indexOf(this.fieldInfo['bk_property_type']) !== -1
            },
            changedValues () {
                const changedValues = {}
                for (const propertyId in this.fieldInfo) {
                    if (JSON.stringify(this.fieldInfo[propertyId]) !== JSON.stringify(this.originalFieldInfo[propertyId])) {
                        changedValues[propertyId] = this.fieldInfo[propertyId]
                    }
                }
                return changedValues
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
            this.originalFieldInfo = this.$tools.clone(this.fieldInfo)
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
                for (const key in this.fieldInfo) {
                    this.fieldInfo[key] = this.$tools.clone(this.field[key])
                }
                this.originalFieldInfo = this.$tools.clone(this.fieldInfo)
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
                        params: this.$injectMetadata(this.fieldInfo, {
                            clone: true, inject: this.isInjectable
                        }),
                        config: {
                            requestId: 'updateObjectAttribute'
                        }
                    }).then(() => {
                        this.$http.cancel(`post_searchObjectAttribute_${this.activeModel['bk_obj_id']}`)
                        this.$http.cancelCache('getHostPropertyList')
                    })
                } else {
                    const groupId = (this.isPublicModel && !this.isAdminView) ? 'bizdefault' : 'default'
                    const otherParams = {
                        creator: this.userName,
                        bk_property_group: this.group.bk_group_id || groupId,
                        bk_property_index: this.propertyIndex || 0,
                        bk_obj_id: this.group.bk_obj_id,
                        bk_supplier_account: this.supplierAccount
                    }
                    await this.createObjectAttribute({
                        params: this.$injectMetadata({
                            ...this.fieldInfo,
                            ...otherParams
                        }, {
                            inject: this.isInjectable
                        }),
                        config: {
                            requestId: 'createObjectAttribute'
                        }
                    }).then(() => {
                        this.$http.cancel(`post_searchObjectAttribute_${this.activeModel['bk_obj_id']}`)
                        this.$http.cancelCache('getHostPropertyList')
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
    .slider-content {
        /deep/ textarea[disabled] {
            background-color: #fafbfd!important;
            cursor: not-allowed;
        }
    }
    .icon-info-circle {
        font-size: 18px;
        color: $cmdbBorderColor;
        padding-left: 5px;
    }
    .field-detail {
        width: 94%;
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
    .cmdb-form-item {
        width: 94% !important;
    }
    .icon-cc-exclamation-tips {
        font-size: 18px;
        color: #979ba5;
        margin-left: 10px;
    }
</style>
