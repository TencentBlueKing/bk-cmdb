<template>
    <div class="model-slider-content">
        <div class="slider-main" ref="sliderMain">
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
                        v-validate="isEditField ? null : 'required|fieldId|reservedWord|length:128'">
                    </bk-input>
                    <p class="form-error" :title="errors.first('fieldId')">{{errors.first('fieldId')}}</p>
                </div>
                <i class="icon-cc-exclamation-tips" tabindex="-1" v-bk-tooltips="$t('模型字段唯一标识提示语')"></i>
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
                        :disabled="isReadOnly || isSystemCreate"
                        v-validate="'required|length:128'">
                    </bk-input>
                    <p class="form-error">{{errors.first('fieldName')}}</p>
                </div>
                <i class="icon-cc-exclamation-tips" v-if="isSystemCreate" tabindex="-1" v-bk-tooltips="$t('国际化配置翻译，不可修改')"></i>
            </label>
            <div class="form-label">
                <span class="label-text">
                    {{$t('字段类型')}}
                    <span class="color-danger">*</span>
                </span>
                <div class="cmdb-form-item">
                    <bk-select
                        class="bk-select-full-width"
                        searchable
                        :clearable="false"
                        v-model="fieldInfo.bk_property_type"
                        :disabled="isEditField"
                        :popover-options="{
                            a11y: false
                        }">
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
                    :ispre="isEditField && field.ispre"
                    :editable.sync="fieldInfo['editable']"
                    :isrequired.sync="fieldInfo['isrequired']"
                ></the-config>
                <component
                    v-if="isComponentShow"
                    :is-read-only="isReadOnly || field.ispre"
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
                <div class="cmdb-form-item" :class="{ 'is-error': errors.has('placeholder') }">
                    <textarea
                        class="raw"
                        name="placeholder"
                        v-model.trim="fieldInfo['placeholder']"
                        :disabled="isReadOnly"
                        v-validate="'length:2000'">
                    </textarea>
                    <p class="form-error" v-if="errors.has('placeholder')">{{errors.first('placeholder')}}</p>
                </div>
            </div>
            <div class="btn-group" :class="{ 'sticky-layout': scrollbar }">
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
    </div>
</template>

<script>
    import { addResizeListener, removeResizeListener } from '@/utils/resize-events'
    import theFieldChar from './char'
    import theFieldInt from './int'
    import theFieldFloat from './float'
    import theFieldEnum from './enum'
    import theFieldList from './list'
    import theFieldBool from './bool'
    import theConfig from './config'
    import { mapGetters, mapActions } from 'vuex'
    import { MENU_BUSINESS } from '@/dictionary/menu-symbol'
    export default {
        components: {
            theFieldChar,
            theFieldInt,
            theFieldFloat,
            theFieldEnum,
            theFieldList,
            theFieldBool,
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
            isMainLineModel: {
                type: Boolean,
                default: false
            },
            customObjId: String,
            propertyIndex: {
                type: Number,
                default: 0
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
                }, {
                    id: 'organization',
                    name: this.$t('组织')
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
                charMap: ['singlechar', 'longchar'],
                scrollbar: false
            }
        },
        computed: {
            ...mapGetters('objectBiz', ['bizId']),
            ...mapGetters(['supplierAccount', 'userName']),
            ...mapGetters('objectModel', ['activeModel']),
            isGlobalView () {
                const topRoute = this.$route.matched[0]
                return topRoute ? topRoute.name !== MENU_BUSINESS : true
            },
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
                return ['singlechar', 'longchar', 'enum', 'int', 'float', 'list', 'bool'].indexOf(this.fieldInfo['bk_property_type']) !== -1
            },
            changedValues () {
                const changedValues = {}
                for (const propertyId in this.fieldInfo) {
                    if (JSON.stringify(this.fieldInfo[propertyId]) !== JSON.stringify(this.originalFieldInfo[propertyId])) {
                        changedValues[propertyId] = this.fieldInfo[propertyId]
                    }
                }
                return changedValues
            },
            isSystemCreate () {
                if (this.isEditField) {
                    return this.field.creator === 'cc_system'
                }
                return false
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
        mounted () {
            addResizeListener(this.$refs.sliderMain, this.handleScrollbar)
        },
        beforeDestroy () {
            removeResizeListener(this.$refs.sliderMain, this.handleScrollbar)
        },
        methods: {
            ...mapActions('objectModelProperty', [
                'createBizObjectAttribute',
                'createObjectAttribute',
                'updateObjectAttribute',
                'updateBizObjectAttribute'
            ]),
            handleScrollbar () {
                const el = this.$refs.sliderMain
                this.scrollbar = el.scrollHeight !== el.offsetHeight
            },
            initData () {
                for (const key in this.fieldInfo) {
                    this.fieldInfo[key] = this.$tools.clone(this.field[key])
                }
                this.originalFieldInfo = this.$tools.clone(this.fieldInfo)
            },
            async validateValue () {
                const validate = [
                    this.$validator.validateAll()
                ]
                if (this.$refs.component) {
                    validate.push(this.$refs.component.$validator.validateAll())
                }
                const results = await Promise.all(validate)
                return results.every(result => result)
            },
            isNullOrUndefinedOrEmpty (value) {
                return [null, '', undefined].includes(value)
            },
            async saveField () {
                if (!await this.validateValue()) {
                    return
                }
                let fieldId = null
                if (this.fieldInfo.bk_property_type === 'int') {
                    this.fieldInfo.option.min = this.isNullOrUndefinedOrEmpty(this.fieldInfo.option.min) ? '' : Number(this.fieldInfo.option.min)
                    this.fieldInfo.option.max = this.isNullOrUndefinedOrEmpty(this.fieldInfo.option.max) ? '' : Number(this.fieldInfo.option.max)
                }
                if (this.isEditField) {
                    const action = this.customObjId ? 'updateBizObjectAttribute' : 'updateObjectAttribute'
                    const params = this.field.ispre ? this.getPreFieldUpdateParams() : this.fieldInfo
                    if (!this.isGlobalView) {
                        params.bk_biz_id = this.bizId
                    }
                    await this[action]({
                        bizId: this.bizId,
                        id: this.field.id,
                        params: params,
                        config: {
                            requestId: 'updateObjectAttribute'
                        }
                    }).then(() => {
                        fieldId = this.fieldInfo.bk_property_id
                        this.$http.cancel(`post_searchObjectAttribute_${this.activeModel['bk_obj_id']}`)
                        this.$http.cancelCache('getHostPropertyList')
                    })
                } else {
                    const groupId = this.isGlobalView ? 'default' : 'bizdefault'
                    const otherParams = {
                        creator: this.userName,
                        bk_property_group: this.group.bk_group_id || groupId,
                        bk_obj_id: this.group.bk_obj_id,
                        bk_supplier_account: this.supplierAccount
                    }
                    const action = this.customObjId ? 'createBizObjectAttribute' : 'createObjectAttribute'
                    const params = {
                        ...this.fieldInfo,
                        ...otherParams
                    }
                    if (!this.isGlobalView) {
                        params.bk_biz_id = this.bizId
                    }
                    await this[action]({
                        bizId: this.bizId,
                        params: params,
                        config: {
                            requestId: 'createObjectAttribute'
                        }
                    }).then(() => {
                        this.$http.cancel(`post_searchObjectAttribute_${this.activeModel['bk_obj_id']}`)
                        this.$http.cancelCache('getHostPropertyList')
                    })
                }
                this.$emit('save', fieldId)
            },
            getPreFieldUpdateParams () {
                const allowKey = ['option', 'unit', 'placeholder']
                const params = {}
                allowKey.forEach(key => {
                    params[key] = this.fieldInfo[key]
                })
                return params
            },
            cancel () {
                this.$emit('cancel')
            }
        }
    }
</script>

<style lang="scss" scoped>
    .model-slider-content {
        height: 100%;
        padding: 0;
        overflow: hidden;
        .slider-main {
            max-height: calc(100% - 52px);
            @include scrollbar-y;
            padding: 20px 20px 0;
        }
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
        .btn-group {
            padding: 10px 20px;
            &.sticky-layout {
                border-top: 1px solid #dcdee5;
            }
        }
    }
</style>
