<!--
 * Tencent
  components: { gridItem },is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2017-2022 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
-->

<script setup>
  import { computed, reactive, watch, getCurrentInstance, inject } from 'vue'
  import { Validator } from 'vee-validate'
  import { t } from '@/i18n'
  import { PROPERTY_TYPES, PROPERTY_TYPE_LIST } from '@/dictionary/property-constants'
  import { getValidateRules } from '@/utils/tools'
  import GridLayout from '@/components/ui/other/grid-layout.vue'
  import GridItem from '@/components/ui/other/grid-item.vue'
  import TheFieldChar from '../char'
  import TheFieldInt from '../int'
  import TheFieldFloat from '../float'
  import TheFieldEnum from '../enum'
  import TheFieldBool from '../bool'

  Validator.extend('maxLongcahr', {
    validate: (value) => {
      if (value !== PROPERTY_TYPES.LONGCHAR) {
        return true
      }
      return existingTypes.value.filter(type => type === PROPERTY_TYPES.LONGCHAR).length < 2
    },
    getMessage: () => t('最多只能添加2个长字符类型')
  })

  const props = defineProps({
    value: {
      type: Boolean
    },
    formData: {
      type: Object
    },
    // 内部表头编辑
    isEdit: Boolean,
    // 外部编辑整个字段
    isEditField: {
      type: Boolean,
      default: false
    }
  })

  const emit = defineEmits(['input', 'save', 'add'])
  const instance = getCurrentInstance().proxy

  const headers = inject('headers')

  const defaultSettings = () => ({
    bk_property_id: '',
    bk_property_name: '',
    bk_property_type: '',
    unit: '',
    placeholder: '',
    bk_property_type: PROPERTY_TYPES.SINGLECHAR,
    editable: true,
    isrequired: false,
    ismultiple: false,
    option: '',
    default: ''
  })

  const settings = reactive(defaultSettings())

  const isShow = computed({
    get() {
      return props.value
    },
    set(val) {
      emit('input', val)
    }
  })

  const typeList = computed(() => {
    const availableTypes = [
      PROPERTY_TYPES.SINGLECHAR,
      PROPERTY_TYPES.LONGCHAR,
      PROPERTY_TYPES.INT,
      PROPERTY_TYPES.FLOAT,
      PROPERTY_TYPES.ENUMMULTI,
      PROPERTY_TYPES.BOOL
    ]
    return PROPERTY_TYPE_LIST.filter(item => availableTypes.includes(item.id))
  })

  const optionComp = computed(() => {
    const comps = {
      [PROPERTY_TYPES.SINGLECHAR]: TheFieldChar,
      [PROPERTY_TYPES.LONGCHAR]: TheFieldChar,
      [PROPERTY_TYPES.INT]: TheFieldInt,
      [PROPERTY_TYPES.FLOAT]: TheFieldFloat,
      [PROPERTY_TYPES.ENUMMULTI]: TheFieldEnum,
      [PROPERTY_TYPES.BOOL]: TheFieldBool
    }
    return comps[settings.bk_property_type]
  })

  const isRequiredShow = computed(() => (
    ![PROPERTY_TYPES.BOOL, PROPERTY_TYPES.ENUMMULTI].includes(settings.bk_property_type)
  ))

  const isDefaultShow = computed(() => ![PROPERTY_TYPES.ENUMMULTI, PROPERTY_TYPES.BOOL]
    .includes(settings.bk_property_type))
  const isUnitShow = computed(() => [PROPERTY_TYPES.INT, PROPERTY_TYPES.FLOAT].includes(settings.bk_property_type))
  const isMultipleShow = computed(() => [PROPERTY_TYPES.ENUMMULTI].includes(settings.bk_property_type))

  const existingIds = computed(() => {
    if (props.isEdit) {
      // 编辑时排除当前id
      return headers.value
        .filter(item => item.bk_property_id !== props.formData.bk_property_id)
        .map(item => item.bk_property_id)
        .join(',')
    }
    return headers.value.map(item => item.bk_property_id).join(',')
  })
  const existingTypes = computed(() => {
    if (props.isEdit) {
      return headers.value
        .filter(item => item.bk_property_type !== props.formData.bk_property_type)
        .map(item => item.bk_property_type)
    }
    return headers.value.map(item => item.bk_property_type)
  })

  watch(isShow, (isShow) => {
    if (isShow) {
      const defaultData = {
        ...defaultSettings(),
        ...props.formData,
      }
      Object.keys(defaultData).forEach((key) => {
        settings[key] = defaultData[key]
      })
    }
  })

  watch(() => settings.bk_property_type, (type) => {
    if (!props.isEditField || !props.isEdit) {
      switch (type) {
        case PROPERTY_TYPES.INT:
        case PROPERTY_TYPES.FLOAT:
          settings.option = {
            min: '',
            max: ''
          }
          settings.default = ''
          break
        case PROPERTY_TYPES.ENUMMULTI:
          settings.option = []
          settings.ismultiple = true
          break
        default:
          settings.default = ''
          settings.option = ''
          settings.ismultiple = false
      }
    }
  })

  const validateValue = async () => {
    const validate = [
      instance.$validator.validateAll()
    ]
    if (instance.$refs.componentOption) {
      validate.push(instance.$refs.componentOption.$validator.validateAll())
    }
    const results = await Promise.all(validate)
    return results
  }

  const getDefaultValueValidateRules = (property) => {
    const rules = getValidateRules(property)
    Reflect.deleteProperty(rules, 'required')
    return rules
  }

  const handleConfirm = async () => {
    const results = await validateValue()
    if (!results.every(result => result)) {
      const comps = [
        instance.$refs.componentBase,
        instance.$refs.componentOption,
      ]
      const index = results.findIndex(item => item === false)
      comps[index]?.$el?.scrollIntoView()
      // 确定是否为默认值校验不通过
      if (index === 0) {
        if (!await instance.$validator.validate('defalut')) {
          instance.$refs.componentDefault?.$el?.scrollIntoView()
        }
      }
      return
    }
    emit(props.isEdit ? 'save' : 'add', { ...settings })
  }
</script>

<template>
  <bk-dialog
    v-model="isShow"
    width="670"
    render-directive="if"
    :title="$t('表格列设置')"
    header-position="left"
    :mask-close="false"
    :auto-close="false"
    @confirm="handleConfirm">
    <div class="content-layout">
      <grid-layout mode="form" :gap="24" :font-size="'14px'" :max-columns="2" ref="componentBase">
        <grid-item
          direction="column"
          required
          :class="['cmdb-form-item', 'form-item', { 'is-error': errors.has('propertyId') }]"
          :label="$t('列ID')">
          <bk-input
            name="propertyId"
            :data-vv-as="$t('列ID')"
            v-validate="`required|fieldId|reservedWord|length:128|excluded:${existingIds}`"
            :disabled="props.isEditField && props.isEdit"
            v-model="settings.bk_property_id">
          </bk-input>
          <template #append>
            <div class="form-error" v-if="errors.has('propertyId')">{{errors.first('propertyId')}}</div>
          </template>
        </grid-item>
        <grid-item
          direction="column"
          required
          :class="['cmdb-form-item', 'form-item', { 'is-error': errors.has('propertyName') }]"
          :label="$t('列名称')">
          <bk-input
            name="propertyName"
            v-validate="'required|length:128'"
            v-model="settings.bk_property_name">
          </bk-input>
          <template #append>
            <div class="form-error" v-if="errors.has('propertyName')">{{errors.first('propertyName')}}</div>
          </template>
        </grid-item>
        <grid-item
          direction="column"
          required
          :class="['cmdb-form-item', 'form-item', { 'is-error': errors.has('propertyType') }]"
          :label="$t('列类型')">
          <bk-select
            name="propertyType"
            class="bk-select-full-width"
            searchable
            :clearable="false"
            v-model="settings.bk_property_type"
            v-validate="'maxLongcahr'"
            :popover-options="{
              appendTo: 'parent'
            }"
            :disabled="props.isEdit">
            <bk-option v-for="(option, index) in typeList"
              :key="index"
              :id="option.id"
              :name="option.name">
            </bk-option>
          </bk-select>
          <template #append>
            <div class="form-error" v-if="errors.has('propertyType')">{{errors.first('propertyType')}}</div>
          </template>
        </grid-item>
        <grid-item
          direction="column"
          :class="['cmdb-form-item', 'form-item', 'checkbox-options']"
          :label="$t('列属性')">
          <bk-checkbox
            class="checkbox"
            v-model="settings.editable">
            <span class="g-has-dashed-tooltips" v-bk-tooltips="$t('字段设置可编辑提示语')">
              {{$t('当前列允许编辑')}}
            </span>
          </bk-checkbox>
          <bk-checkbox
            class="ml10 checkbox"
            v-if="isRequiredShow"
            v-model="settings.isrequired">
            {{$t('当前列的值必填')}}
          </bk-checkbox>
          <bk-checkbox
            class="ml10 checkbox"
            v-if="isMultipleShow"
            v-model="settings.ismultiple">
            {{$t('可多选')}}
          </bk-checkbox>
        </grid-item>
      </grid-layout>

      <grid-layout class="field-option-container" mode="form" :gap="0" :font-size="'14px'" :max-columns="1">
        <component
          :key="settings.bk_property_type"
          :is="optionComp"
          :multiple="settings.ismultiple"
          v-model="settings.option"
          :type="settings.bk_property_type"
          ref="componentOption">
        </component>
        <div class="form-label" v-if="isDefaultShow">
          <div class="label-text">
            <span class="g-has-dashed-tooltips" v-bk-tooltips="$t('表格字段列默认值提示语')">{{$t('列默认值')}}</span>
          </div>
          <div :class="['cmdb-form-item', 'form-item', { 'is-error': errors.has('defalut') }]">
            <component
              name="defalut"
              :key="settings.bk_property_type"
              :is="`cmdb-form-${settings.bk_property_type}`"
              :multiple="settings.ismultiple"
              :options="settings.option || []"
              v-model="settings.default"
              v-validate="getDefaultValueValidateRules(settings)"
              ref="componentDefault">
            </component>
            <div class="form-error" v-if="errors.has('defalut')">{{errors.first('defalut')}}</div>
          </div>
        </div>
      </grid-layout>

      <grid-layout mode="form" :gap="24" :font-size="'14px'" :max-columns="1">
        <grid-item
          v-show="isUnitShow"
          direction="column"
          :class="['cmdb-form-item', 'form-item']"
          :label="$t('单位')">
          <bk-input type="text" class="cmdb-form-input"
            v-model.trim="settings.unit"
            :placeholder="$t('请输入单位')">
          </bk-input>
        </grid-item>
      </grid-layout>
    </div>
  </bk-dialog>
</template>

<style lang="scss" scoped>
  .content-layout {
    height: 340px;
    padding: 0 12px;
    @include scrollbar-y;
  }

  .field-option-container {
    width: 100%;
    margin-bottom: 20px;
    margin-top: 20px;
    padding: 16px;
    background: #F5F7FB;

    :deep(.form-label) {
      margin-bottom: 15px;
    }

    :deep(.label-text) {
      position: relative;
      display: block;
      padding-right: 10px;
      line-height: 36px;
      font-size: 14px;
      @include ellipsis;
    }
  }

  .checkbox-options {
    margin-bottom: 10px;
    .checkbox {
      height: 24px;
      line-height: 24px;
    }
  }
</style>
