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
        <div v-bk-tooltips.top.light.click="$t('模型字段唯一标识提示语')"
          class="cmdb-form-item" :class="{ 'is-error': errors.has('fieldId') }">
          <bk-input type="text" class="cmdb-form-input"
            name="fieldId"
            v-model.trim="fieldInfo.bk_property_id"
            :placeholder="$t('请输入唯一标识')"
            :disabled="isEditField"
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
            :disabled="isReadOnly || isSystemCreate || field.ispre"
            v-validate="'required|length:128'">
          </bk-input>
          <p class="form-error">{{errors.first('fieldName')}}</p>
        </div>
      </label>
      <div class="form-label">
        <span class="label-text">
          {{$t('字段分组')}}
          <span class="color-danger">*</span>
        </span>
        <div class="cmdb-form-item">
          <bk-select
            class="bk-select-full-width"
            searchable
            :clearable="false"
            v-model="fieldInfo.bk_property_group"
            :disabled="isEditField">
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
      <div class="field-detail" v-show="!['foreignkey'].includes(fieldType)">
        <the-config
          :type="fieldInfo.bk_property_type"
          :is-read-only="isReadOnly"
          :is-main-line-model="isMainLineModel"
          :ispre="isEditField && field.ispre"
          :is-edit-field="isEditField"
          :editable.sync="fieldInfo.editable"
          :isrequired.sync="fieldInfo.isrequired"
          :multiple.sync="fieldInfo.ismultiple"
        ></the-config>
        <!-- 添加key防止复用组件时内部状态错误 -->
        <component
          class="cmdb-form-item"
          :key="fieldInfo.bk_property_type"
          v-if="isSettingComponentShow"
          :is-read-only="isReadOnly || field.ispre"
          :is="`the-field-${fieldType}`"
          :multiple="fieldInfo.ismultiple"
          v-model="fieldInfo.option"
          :is-edit-field="isEditField"
          :type="fieldInfo.bk_property_type"
          :default-value.sync="fieldInfo.default"
          ref="component"
        ></component>
        <label class="form-label" v-if="isDefaultComponentShow">
          <span class="label-text">
            {{$t('默认值')}}
          </span>
          <div class="cmdb-form-item" :class="{ 'is-error': errors.has('defalut') }">
            <component
              name="defalut"
              :key="fieldInfo.bk_property_type"
              :is-read-only="isReadOnly || field.ispre"
              :is="`cmdb-form-${fieldInfo.bk_property_type}`"
              :multiple="fieldInfo.ismultiple"
              :options="fieldInfo.option || []"
              :disabled="isReadOnly || isSystemCreate || field.ispre"
              v-model="fieldInfo.default"
              v-validate="getValidateRules(fieldInfo)"
              ref="component"
            ></component>
            <p class="form-error">{{errors.first('defalut')}}</p>
          </div>
        </label>
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
          <bk-input
            class="raw"
            :rows="3"
            :maxlength="100"
            name="placeholder"
            :type="'textarea'"
            v-model.trim="fieldInfo['placeholder']"
            :disabled="isReadOnly"
            v-validate="'length:2000'">
          </bk-input>
          <p class="form-error" v-if="errors.has('placeholder')">{{errors.first('placeholder')}}</p>
        </div>
      </div>
    </div>
    <template slot="footer" slot-scope="{ sticky }">
      <div class="btn-group" :class="{ 'is-sticky': sticky }">
        <bk-button theme="primary"
          :loading="$loading(['updateObjectAttribute', 'createObjectAttribute'])"
          @click="saveField">
          {{isEditField ? $t('保存') : $t('提交')}}
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
        required: true
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
        charMap: [PROPERTY_TYPES.SINGLECHAR, PROPERTY_TYPES.LONGCHAR]
      }
    },
    provide() {
      return {
        customObjId: this.customObjId
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

      const { beforeClose } = useSideslider(this.fieldInfo)
      this.beforeClose = beforeClose
    },
    methods: {
      ...mapActions('objectModelProperty', [
        'createBizObjectAttribute',
        'createObjectAttribute',
        'updateObjectAttribute',
        'updateBizObjectAttribute'
      ]),
      initData() {
        Object.keys(this.fieldInfo).forEach((key) => {
          this.fieldInfo[key] = this.$tools.clone(this.field[key] ?? '')
        })
        this.originalFieldInfo = this.$tools.clone(this.fieldInfo)
      },
      async validateValue() {
        const validate = [
          this.$validator.validateAll()
        ]
        if (this.$refs.component) {
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
        if (this.isEditField) {
          const action = this.customObjId ? 'updateBizObjectAttribute' : 'updateObjectAttribute'
          const params = this.field.ispre ? this.getPreFieldUpdateParams() : this.fieldInfo
          if (!this.isGlobalView) {
            params.bk_biz_id = this.bizId
          }
          if (isEmptyPropertyValue(this.fieldInfo.default)) {
            params.default = null
          }
          await this[action]({
            bizId: this.bizId,
            id: this.field.id,
            params,
            config: {
              requestId: 'updateObjectAttribute'
            }
          }).then(() => {
            fieldId = this.fieldInfo.bk_property_id
            this.$http.cancel(`post_searchObjectAttribute_${this.activeModel.bk_obj_id}`)
            this.$http.cancelCache('getHostPropertyList')
            this.$success(this.$t('修改成功'))
          })
        } else {
          if (isEmptyPropertyValue(this.fieldInfo.default)) {
            Reflect.deleteProperty(this.fieldInfo, 'default')
          }
          const groupId = this.isGlobalView ? 'default' : 'bizdefault'
          const selectedGroup = this.groups.find(group => group.bk_group_id === this.fieldInfo.bk_property_group)
          const otherParams = {
            creator: this.userName,
            bk_property_group: this.fieldInfo.bk_property_group || this.group.bk_group_id || groupId,
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
            this.$http.cancel(`post_searchObjectAttribute_${this.activeModel.bk_obj_id}`)
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
