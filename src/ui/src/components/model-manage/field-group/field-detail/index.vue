<!--
 * Tencent is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2017-2022 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
-->

<template>
  <cmdb-sticky-layout class="model-slider-content">
    <div class="slider-main" ref="sliderMain">
      <label class="form-label">
        <span class="label-text">
          {{$t('唯一标识')}}
          <span class="color-danger">*</span>
        </span>
        <div :class="['cmdb-form-item', { 'is-error': errors.has('fieldId') }]"
          v-bk-tooltips.top="{
            disabled: isEditField && !isCreateMode,
            theme: 'light',
            trigger: 'click',
            content: $t('模型字段唯一标识提示语')
          }">
          <bk-input type="text" class="cmdb-form-input"
            name="fieldId"
            v-model.trim="fieldInfo.bk_property_id"
            :placeholder="$t('请输入唯一标识')"
            :disabled="isEditField && !isCreateMode"
            v-validate="isEditField ? null : 'required|fieldId|reservedWord|length:128'">
          </bk-input>
          <p class="form-error" :title="errors.first('fieldId')">{{errors.first('fieldId')}}</p>
        </div>
      </label>
      <label class="form-label">
        <span class="label-text">
          {{$t('名称')}}
          <span class="color-danger">*</span>
        </span>
        <div class="cmdb-form-item" :class="{ 'is-error': errors.has('fieldName') }">
          <bk-input type="text" class="cmdb-form-input"
            name="fieldName"
            :placeholder="isSystemCreate ? $t('请输入名称，国际化时将会自动翻译，不可修改') : $t('请输入字段名称')"
            v-model.trim="fieldInfo.bk_property_name"
            :disabled="isReadOnly || isSystemCreate || field.ispre || isFromTemplateField"
            v-validate="'required|length:128'">
          </bk-input>
          <p class="form-error">{{errors.first('fieldName')}}</p>
        </div>
      </label>
      <div class="form-label" v-if="!isSettingScene">
        <span class="label-text">
          {{$t('字段分组')}}
          <span class="color-danger">*</span>
        </span>
        <div class="cmdb-form-item">
          <bk-select
            class="bk-select-full-width"
            searchable
            :clearable="false"
            v-model="fieldInfo.bk_property_group">
            <bk-option v-for="(option, index) in groups"
              :key="index"
              :id="option.bk_group_id"
              :name="option.bk_group_name">
            </bk-option>
          </bk-select>
        </div>
      </div>
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
            :disabled="isEditField && !isCreateMode"
            :popover-options="{
              a11y: false
            }">
            <bk-option v-for="(option, index) in availableFieldTypeList"
              :key="index"
              :id="option.id"
              :name="option.name">
            </bk-option>
          </bk-select>
        </div>
      </div>
      <div class="field-detail" v-if="!['foreignkey'].includes(fieldType)">
        <!-- 添加key防止复用组件时内部状态错误 -->
        <component
          class="cmdb-form-item"
          :key="fieldInfo.bk_property_type"
          v-if="isSettingComponentShow"
          :is-read-only="isReadOnly || field.ispre || isFromTemplateField"
          :is="`the-field-${fieldType}`"
          :multiple.sync="fieldInfo.ismultiple"
          v-model="fieldInfo.option"
          :is-edit-field="isEditField"
          :type="fieldInfo.bk_property_type"
          :default-value.sync="fieldInfo.default"
          ref="component">
        </component>
        <label class="form-label" v-if="isDefaultComponentShow">
          <span class="label-text">
            {{$t('默认值')}}
          </span>
          <div class="cmdb-form-item" :class="{ 'is-error': errors.has('defalut') }">
            <div class="form-item-row">
              <component
                :class="`form-item-${fieldInfo.bk_property_type}`"
                name="defalut"
                :key="fieldInfo.bk_property_type"
                :is-read-only="isReadOnly || field.ispre"
                :is="`cmdb-form-${fieldInfo.bk_property_type}`"
                :multiple="fieldInfo.ismultiple"
                :options="fieldInfo.option || []"
                :disabled="isReadOnly || isSystemCreate || field.ispre || isFromTemplateField"
                v-model="fieldInfo.default"
                v-validate="getValidateRules(fieldInfo)"
                ref="component">
              </component>
              <bk-checkbox
                v-if="isMultipleShow"
                class="checkbox"
                v-model="fieldInfo.ismultiple"
                :disabled="isReadOnly || field.ispre || isFromTemplateField">
                <span>{{$t('可多选')}}</span>
              </bk-checkbox>
            </div>
            <p class="form-error">{{errors.first('defalut')}}</p>
          </div>
        </label>
      </div>
      <the-config
        :type="fieldInfo.bk_property_type"
        :is-read-only="isReadOnly"
        :is-main-line-model="isMainLineModel"
        :ispre="isEditField && field.ispre"
        :is-edit-field="isEditField"
        :editable.sync="fieldInfo.editable"
        :isrequired.sync="fieldInfo.isrequired"
        :isrequired-lock.sync="fieldSettingExtra.lock.isrequired"
        :editable-lock.sync="fieldSettingExtra.lock.editable">
      </the-config>
      <label class="form-label" v-show="['int', 'float'].includes(fieldType)">
        <span class="label-text">
          {{$t('单位')}}
        </span>
        <div class="cmdb-form-item">
          <bk-input type="text" class="cmdb-form-input"
            v-model.trim="fieldInfo['unit']"
            :disabled="isReadOnly || isFromTemplateField"
            :placeholder="$t('请输入单位')">
          </bk-input>
        </div>
      </label>
      <div class="form-label">
        <span class="label-text">{{$t('用户提示')}}</span>
        <div class="cmdb-form-item" :class="{ 'is-error': errors.has('placeholder') }">
          <div class="form-component">
            <bk-input
              class="raw"
              :rows="3"
              :maxlength="2000"
              name="placeholder"
              :type="'textarea'"
              v-model.trim="fieldInfo['placeholder']"
              :disabled="isReadOnly || (isFromTemplateField && fieldSettingExtra.lock.placeholder)"
              v-validate="'length:2000'">
            </bk-input>
          </div>
          <p class="form-error" v-if="errors.has('placeholder')">{{errors.first('placeholder')}}</p>
        </div>
      </div>
      <slot name="append-unique"
        v-if="isUniqueSettingShow"
        v-bind="{
          disabled: uniqueDisabled,
          fieldInfo: { id: field.id, ...fieldInfo }
        }">
      </slot>
      <slot name="append-lock"
        v-bind="{
          fieldInfo: { id: field.id, ...fieldInfo }
        }">
      </slot>
    </div>
    <template slot="footer" slot-scope="{ sticky }">
      <div class="btn-group" :class="{ 'is-sticky': sticky }">
        <bk-button theme="primary"
          :loading="$loading(['updateObjectAttribute', 'createObjectAttribute'])"
          :disabled="saveDisabled"
          @click="saveField">
          <span v-if="!isSettingScene">{{isEditField ? $t('保存') : $t('提交')}}</span>
          <span v-else>{{ $t('确定') }}</span>
        </bk-button>
        <bk-button theme="default" @click="cancel">
          {{$t('取消')}}
        </bk-button>
      </div>
    </template>
  </cmdb-sticky-layout>
</template>

<script>
  import theFieldChar from './char'
  import theFieldInt from './int'
  import theFieldFloat from './float'
  import theFieldEnum from './enum'
  import theFieldList from './list'
  import theFieldBool from './bool'
  import theFieldEnumquote from './enumquote.vue'
  import theFieldInnertable from './inner-table/index.vue'
  import theConfig from './config'
  import { mapGetters, mapActions } from 'vuex'
  import { MENU_BUSINESS } from '@/dictionary/menu-symbol'
  import { PROPERTY_TYPES, PROPERTY_TYPE_LIST } from '@/dictionary/property-constants'
  import { isEmptyPropertyValue } from '@/utils/tools'
  import useSideslider from '@/hooks/use-sideslider'
  import { singleRuleTypes, unionRuleTypes } from '@/views/field-template/children/use-unique'
  import fieldTemplateService from '@/service/field-template'

  export default {
    components: {
      theFieldChar,
      theFieldInt,
      theFieldFloat,
      theFieldEnum,
      theFieldEnumquote,
      theFieldList,
      theFieldBool,
      theConfig,
      theFieldInnertable
    },
    props: {
      properties: {
        type: Array,
        default: () => []
      },
      field: {
        type: Object
      },
      group: {
        type: Object,
        default: () => ({})
      },
      groups: {
        type: Array,
        default: () => []
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
      },
      scene: {
        type: String,
        default: 'manage' // setting表示配置流程
      },
      excludeType: {
        type: Array,
        default: () => []
      },
      fieldSetting: {
        type: Object
      },
      // 是否为创建模式，新创建的字段二次编辑时可以修改ID等
      isCreateMode: {
        type: Boolean,
        default: false
      }
    },
    data() {
      return {
        fieldInfo: {
          bk_property_name: '',
          bk_property_id: '',
          bk_property_group: this.group.bk_group_id,
          unit: '',
          placeholder: '',
          bk_property_type: PROPERTY_TYPES.SINGLECHAR,
          editable: true,
          isrequired: false,
          ismultiple: false,
          option: '',
          default: ''
        },
        originalFieldInfo: {},
        charMap: [PROPERTY_TYPES.SINGLECHAR, PROPERTY_TYPES.LONGCHAR],
        // 字段组合模板配置额外数据
        fieldSettingExtra: {
          lock: {
            isrequired: true,
            editable: true,
            placeholder: true
          }
        },
        uniqueSetting: {
          disalbed: true
        }
      }
    },
    provide() {
      return {
        customObjId: this.customObjId,
        isSettingScene: this.isSettingScene,
        isFromTemplateField: this.isFromTemplateField
      }
    },
    computed: {
      ...mapGetters('objectBiz', ['bizId']),
      ...mapGetters(['supplierAccount', 'userName']),
      ...mapGetters('objectModel', ['activeModel']),
      isGlobalView() {
        const [topRoute] = this.$route.matched
        return topRoute ? topRoute.name !== MENU_BUSINESS : true
      },
      saveDisabled() {
        return !this.isSettingScene && this.isEditField && !this.isChanged?.value
      },
      fieldType() {
        let { bk_property_type: type } = this.fieldInfo
        switch (type) {
          case PROPERTY_TYPES.SINGLECHAR:
          case PROPERTY_TYPES.LONGCHAR:
            type = 'char'
            break
          case PROPERTY_TYPES.ENUMMULTI:
            type = 'enum'
            break
        }
        return type
      },
      isSettingComponentShow() {
        const types = [
          PROPERTY_TYPES.SINGLECHAR,
          PROPERTY_TYPES.LONGCHAR,
          PROPERTY_TYPES.ENUM,
          PROPERTY_TYPES.INT,
          PROPERTY_TYPES.FLOAT,
          PROPERTY_TYPES.LIST,
          PROPERTY_TYPES.BOOL,
          PROPERTY_TYPES.ENUMMULTI,
          PROPERTY_TYPES.ENUMQUOTE,
          PROPERTY_TYPES.INNER_TABLE
        ]
        return types.indexOf(this.fieldInfo.bk_property_type) !== -1
      },
      isDefaultComponentShow() {
        const types = [
          PROPERTY_TYPES.ENUM,
          PROPERTY_TYPES.ENUMMULTI,
          PROPERTY_TYPES.ENUMQUOTE,
          PROPERTY_TYPES.BOOL,
          PROPERTY_TYPES.INNER_TABLE,
          PROPERTY_TYPES.LIST
        ]
        return !types.includes(this.fieldInfo.bk_property_type)
      },
      changedValues() {
        const changedValues = {}
        Object.keys(this.fieldInfo).forEach((propertyId) => {
          if (JSON.stringify(this.fieldInfo[propertyId]) !== JSON.stringify(this.originalFieldInfo[propertyId])) {
            changedValues[propertyId] = this.fieldInfo[propertyId]
          }
        })
        return changedValues
      },
      isSystemCreate() {
        if (this.isEditField) {
          return this.field.creator === 'cc_system'
        }
        return false
      },
      fieldTypeList() {
        if (this.customObjId) {
          const disabledTypes = this.isEditField
            ? [PROPERTY_TYPES.INNER_TABLE]
            : [PROPERTY_TYPES.INNER_TABLE, PROPERTY_TYPES.FOREIGNKEY, PROPERTY_TYPES.ENUMQUOTE]
          return PROPERTY_TYPE_LIST.filter(item => !disabledTypes.includes(item.id))
        }
        // eslint-disable-next-line max-len
        const createFieldList = PROPERTY_TYPE_LIST.filter(item => ![PROPERTY_TYPES.ENUMQUOTE, PROPERTY_TYPES.FOREIGNKEY].includes(item.id))
        return this.isEditField ? PROPERTY_TYPE_LIST : createFieldList
      },
      availableFieldTypeList() {
        return this.fieldTypeList.filter(item => !this.excludeType.includes(item.id))
      },
      // 是否为字段组合模板配置流程
      isSettingScene() {
        return this.scene === 'setting'
      },
      isMultipleShow() {
        const types = [
          PROPERTY_TYPES.ORGANIZATION,
          PROPERTY_TYPES.ENUMQUOTE,
          PROPERTY_TYPES.ENUMMULTI
        ]
        return types.includes(this.fieldInfo.bk_property_type)
      },
      uniqueDisabled() {
        return !this.fieldInfo.bk_property_id || !this.fieldInfo.bk_property_name
      },
      isUniqueSettingShow() {
        const types = [...new Set([...singleRuleTypes, ...unionRuleTypes])]
        return types.includes(this.fieldInfo.bk_property_type) && !this.uniqueDisabled
      },
      isFromTemplateField() {
        return !this.isSettingScene && this.field.bk_template_id > 0
      }
    },
    watch: {
      'fieldInfo.bk_property_type'(type) {
        if (!this.isEditField) {
          switch (type) {
            case 'int':
            case 'float':
              this.fieldInfo.option = {
                min: '',
                max: ''
              }
              this.fieldInfo.default = ''
              break
            case PROPERTY_TYPES.ENUMMULTI:
            case PROPERTY_TYPES.ENUMQUOTE:
              this.fieldInfo.option = []
              this.fieldInfo.ismultiple = true
              break
            case PROPERTY_TYPES.OBJUSER:
            case PROPERTY_TYPES.ORGANIZATION:
              this.fieldInfo.default = ''
              this.fieldInfo.option = ''
              this.fieldInfo.ismultiple = true
              break
            case PROPERTY_TYPES.INNER_TABLE:
              this.fieldInfo.option = {
                header: [],
                default: []
              }
            default:
              this.fieldInfo.default = ''
              this.fieldInfo.option = ''
              this.fieldInfo.ismultiple = false
          }
        }
      }
    },
    created() {
      this.originalFieldInfo = this.$tools.clone(this.fieldInfo)
      if (this.isEditField) {
        this.initData()
      }
      const { beforeClose, isChanged } = useSideslider(this.fieldInfo)
      this.beforeClose = beforeClose
      this.isChanged = isChanged
    },
    methods: {
      ...mapActions('objectModelProperty', [
        'createBizObjectAttribute',
        'createObjectAttribute',
        'updateObjectAttribute',
        'updateBizObjectAttribute'
      ]),
      ...mapActions('objectModelFieldGroup', [
        'updatePropertySort'
      ]),
      async initData() {
        Object.keys(this.fieldInfo).forEach((key) => {
          this.fieldInfo[key] = this.$tools.clone(this.field[key] ?? '')
        })
        this.originalFieldInfo = this.$tools.clone(this.fieldInfo)

        if (this.isSettingScene) {
          Object.keys(this.fieldSetting).forEach((key) => {
            this.fieldSettingExtra[key] = this.$tools.clone(this.fieldSetting[key] ?? '')
          })
        } else if (this.isFromTemplateField) {
          // 模型字段编辑来自模板的字段，需要获取模板中针对字段的配置
          const templateFieldList = await fieldTemplateService.getTemplateFieldListByField(this.field)
          const templateField = templateFieldList.find(item => item.id === this.field.bk_template_id)
          this.fieldSettingExtra.lock.isrequired = templateField?.isrequired?.lock
          this.fieldSettingExtra.lock.editable = templateField?.editable?.lock
          this.fieldSettingExtra.lock.placeholder = templateField?.placeholder?.lock
        }
      },
      async validateValue() {
        const validate = [
          this.$validator.validateAll()
        ]
        if (this.$refs?.component?.$validator) {
          validate.push(this.$refs.component.$validator.validateAll())
        }
        const results = await Promise.all(validate)
        return results.every(result => result)
      },
      isNullOrUndefinedOrEmpty(value) {
        return [null, '', undefined].includes(value)
      },
      async saveField() {
        if (!await this.validateValue()) {
          return
        }

        const tableTypeCount = this.properties
          .filter(property => property.bk_property_type === PROPERTY_TYPES.INNER_TABLE).length
        const isTableType = this.fieldInfo.bk_property_type === PROPERTY_TYPES.INNER_TABLE
        if (!this.isEditField && isTableType && tableTypeCount >= 5) {
          this.$error('最多只能添加5个表格字段')
          return
        }

        let fieldId = null
        if (this.fieldInfo.bk_property_type === 'int' || this.fieldInfo.bk_property_type === 'float') {
          this.fieldInfo.option.min = this.isNullOrUndefinedOrEmpty(this.fieldInfo.option.min) ? '' : Number(this.fieldInfo.option.min)
          this.fieldInfo.option.max = this.isNullOrUndefinedOrEmpty(this.fieldInfo.option.max) ? '' : Number(this.fieldInfo.option.max)
          this.fieldInfo.default = this.isNullOrUndefinedOrEmpty(this.fieldInfo.default) ? '' : Number(this.fieldInfo.default)
        }

        // 配置流程直接抛出事件并退出，在流程中自行处理
        if (this.isSettingScene) {
          this.$emit('confirm', this.field.id, this.fieldInfo)
          return
        }
        const defaultGroupld = this.isGlobalView ? 'default' : 'bizdefault'
        const groupId =  this.fieldInfo.bk_property_group || defaultGroupld
        const activeObjId = this.$route.params.modelId || this.customObjId
        if (this.isEditField) {
          const action = this.customObjId ? 'updateBizObjectAttribute' : 'updateObjectAttribute'
          let params = this.field.ispre ? this.getPreFieldUpdateParams() : this.fieldInfo
          if (isEmptyPropertyValue(this.fieldInfo.default)) {
            params.default = null
          }
          if (this.isFromTemplateField) {
            params = this.getTemplateFieldParams()
          }
          if (!this.isGlobalView) {
            params.bk_biz_id = this.bizId
          }
          // 是否调用编辑字段接口
          const paramsLen = Object.keys(params).length
          // 是否调用修改字段分组接口
          const isChangeGroupId = groupId !== this.group.bk_group_id
          if (paramsLen) {
            await this[action]({
              bizId: this.bizId,
              id: this.field.id,
              params,
              config: {
                requestId: 'updateObjectAttribute'
              }
            }).then(() => {
              fieldId = this.fieldInfo.bk_property_id
              this.$http.cancel(`post_searchObjectAttribute_${activeObjId}`)
              this.$http.cancelCache('getHostPropertyList')
            })
          }
          if (isChangeGroupId) {
            // 每次都查询出最后一个元素的index
            const toLastIndex = this.properties.
              filter(property => property.bk_property_group === groupId)?.at(-1)?.bk_property_index
              ?? -1
            // 修改分组
            await this.updatePropertySort({
              objId: activeObjId,
              propertyId: this.field.id,
              params: {
                bk_biz_id: params.bk_biz_id,
                bk_property_group: groupId,
                bk_property_index: toLastIndex + 1
              },
              config: {
                requestId: `updatePropertySort_${activeObjId}`
              }
            }).then(() => {
              this.$success(this.$t('修改成功'))
            })
          } else {
            this.$success(this.$t('修改成功'))
          }
        } else {
          if (isEmptyPropertyValue(this.fieldInfo.default)) {
            Reflect.deleteProperty(this.fieldInfo, 'default')
          }
          const selectedGroup = this.groups.find(group => group.bk_group_id === this.fieldInfo.bk_property_group)
          const otherParams = {
            creator: this.userName,
            bk_property_group: groupId,
            bk_obj_id: selectedGroup?.bk_obj_id || this.group.bk_obj_id,
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
            params,
            config: {
              requestId: 'createObjectAttribute'
            }
          }).then(() => {
            this.$http.cancel(`post_searchObjectAttribute_${activeObjId}`)
            this.$http.cancelCache('getHostPropertyList')
            this.$success(this.$t('创建成功'))
          })
        }
        this.$emit('save', fieldId)
      },
      getPreFieldUpdateParams() {
        const allowKey = ['option', 'unit', 'placeholder']
        const params = {}
        allowKey.forEach((key) => {
          params[key] = this.fieldInfo[key]
        })
        return params
      },
      getTemplateFieldParams() {
        const allowKey = ['isrequired', 'editable', 'placeholder']
        const params = {}
        allowKey.forEach((key) => {
          if (!this.fieldSettingExtra.lock[key]) {
            params[key] = this.fieldInfo[key]
          }
        })
        return params
      },
      cancel() {
        this.$emit('cancel')
      },
      getValidateRules(fieldInfo) {
        const rules =  this.$tools.getValidateRules(fieldInfo)
        Reflect.deleteProperty(rules, 'required')
        return rules
      }
    }
  }
</script>

<style lang="scss" scoped>
    .model-slider-content {
        height: 100%;
        padding: 0;
        @include scrollbar-y;
        .slider-main {
            max-height: calc(100% - 52px);
            @include scrollbar-y;
            padding: 20px 24px;
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
            width: 100%;
            margin-bottom: 20px;
            padding: 20px;
            background: #F5F7FB;
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
            width: 100%;
            &.is-error {
                /deep/ .bk-form-input {
                    border-color: #ff5656;
                }
            }

            .form-item-row {
              display: flex;
              align-items: center;
              gap: 12px;
              min-height: 32px;
              position: relative;

              .checkbox {
                flex: none;
              }
              .form-item-objuser {
                width: 100%;
              }
            }
        }
        .icon-cc-exclamation-tips {
            font-size: 18px;
            color: #979ba5;
            margin-left: 10px;
        }
        .btn-group {
            padding: 8px 24px;
            &.is-sticky {
                border-top: 1px solid #dcdee5;
            }
            .bk-button{
                width: 88px;
                height: 32px;
            }
        }
    }
</style>
<style lang="scss">
@import '@/assets/scss/model-manage.scss';
</style>
