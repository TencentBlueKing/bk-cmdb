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
  import { computed, reactive, watch, getCurrentInstance } from 'vue'
  import { PROPERTY_TYPES, PROPERTY_TYPE_LIST } from '@/dictionary/property-constants'
  import GridLayout from '@/components/ui/other/grid-layout.vue'
  import GridItem from '@/components/ui/other/grid-item.vue'
  import TheFieldChar from '../char'
  import TheFieldInt from '../int'
  import TheFieldFloat from '../float'
  import TheFieldEnum from '../enum'
  import TheFieldBool from '../bool'

  const props = defineProps({
    value: {
      type: Boolean
    },
    formData: {
      type: Object
    },
    isEdit: Boolean
  })

  const emit = defineEmits(['input', 'save', 'add'])
  const instance = getCurrentInstance().proxy

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
      PROPERTY_TYPES.ENUM,
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
      [PROPERTY_TYPES.ENUM]: TheFieldEnum,
      [PROPERTY_TYPES.BOOL]: TheFieldBool
    }
    return comps[settings.bk_property_type]
  })

  const isRequiredShow = computed(() => (
    ![PROPERTY_TYPES.BOOL, PROPERTY_TYPES.ENUM].includes(settings.bk_property_type)
  ))

  const isDefaultShow = computed(() => ![PROPERTY_TYPES.ENUM, PROPERTY_TYPES.BOOL].includes(settings.bk_property_type))
  const isUnitShow = computed(() => [PROPERTY_TYPES.INT, PROPERTY_TYPES.FLOAT].includes(settings.bk_property_type))

  watch(isShow, (isShow) => {
    if (isShow) {
      const defaultData = {
        ...defaultSettings(),
        ...props.formData,
      }
      Object.keys(defaultData).forEach((key) => {
        settings[key] = defaultData[key]
      })
      console.log(settings, '---watch after settings', props.formData)
    }
  })

  watch(() => settings.bk_property_type, (type) => {
    switch (type) {
      case PROPERTY_TYPES.INT:
      case PROPERTY_TYPES.FLOAT:
        settings.option = {
          min: '',
          max: ''
        }
        settings.default = ''
        break
      default:
        settings.default = ''
        settings.option = ''
        settings.ismultiple = false
    }
  })

  const validateValue = async () => {
    const validate = [
      instance.$validator.validateAll()
    ]
    if (instance.$refs.component) {
      validate.push(instance.$refs.component.$validator.validateAll())
    }
    console.log(validate, '--xxxx')
    const results = await Promise.all(validate)
    return results.every(result => result)
  }

  const handleConfirm = async () => {
    if (!await validateValue()) {
      return
    }
    emit(props.isEdit ? 'save' : 'add', { ...settings })
  }
</script>

<template>
  <bk-dialog
    v-model="isShow"
    width="640"
    :title="$t('字段设置')"
    header-position="left"
    :mask-close="false"
    :auto-close="false"
    @confirm="handleConfirm">
    <div class="content-layout">
      <grid-layout mode="form" :gap="36" :font-size="'14px'" :max-columns="2">
        <grid-item
          direction="column"
          required
          :class="['cmdb-form-item', 'form-item', { 'is-error': errors.has('propertyId') }]"
          :label="$t('字段ID')">
          <bk-input
            name="propertyId"
            v-validate="'required'"
            :disabled="props.isEdit"
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
          :label="$t('字段名称')">
          <bk-input
            name="propertyName"
            v-validate="'required'"
            v-model="settings.bk_property_name">
          </bk-input>
          <template #append>
            <div class="form-error" v-if="errors.has('propertyName')">{{errors.first('propertyName')}}</div>
          </template>
        </grid-item>
        <grid-item
          direction="column"
          required
          :class="['cmdb-form-item', 'form-item']"
          :label="$t('字段类型')">
          <bk-select
            class="bk-select-full-width"
            searchable
            :clearable="false"
            v-model="settings.bk_property_type"
            :disabled="props.isEdit">
            <bk-option v-for="(option, index) in typeList"
              :key="index"
              :id="option.id"
              :name="option.name">
            </bk-option>
          </bk-select>
        </grid-item>
        <grid-item
          direction="column"
          :class="['cmdb-form-item', 'form-item']"
          :label="$t('字段属性')">
          <bk-checkbox
            v-model="settings.editable">
            {{$t('可编辑')}}
          </bk-checkbox>
          <bk-checkbox
            class="ml10"
            v-if="isRequiredShow"
            v-model="settings.isrequired">
            {{$t('必填')}}
          </bk-checkbox>
        </grid-item>
      </grid-layout>

      <grid-layout class="field-option-container" mode="form" :gap="36" :font-size="'14px'" :max-columns="1">
        <component
          :key="settings.bk_property_type"
          :is="optionComp"
          :multiple="settings.ismultiple"
          v-model="settings.option"
          :type="settings.bk_property_type"
          ref="component">
        </component>
      </grid-layout>

      <grid-layout mode="form" :gap="24" :font-size="'14px'" :max-columns="1">
        <grid-item
          v-if="isDefaultShow"
          direction="column"
          :class="['cmdb-form-item', 'form-item', { 'is-error': errors.has('defalut') }]"
          :label="$t('默认值')">
          <component
            name="defalut"
            :key="settings.bk_property_type"
            :is="`cmdb-form-${settings.bk_property_type}`"
            :multiple="settings.ismultiple"
            :options="settings.option || []"
            v-model="settings.default"
            v-validate="$tools.getValidateRules(settings)"
            ref="component">
          </component>
          <template #append>
            <div class="form-error" v-if="errors.has('defalut')">{{errors.first('defalut')}}</div>
          </template>
        </grid-item>
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
        <grid-item
          direction="column"
          :class="['cmdb-form-item', 'form-item', { 'is-error': errors.has('placeholder') }]"
          :label="$t('用户提示')">
          <bk-input
            class="raw"
            :rows="3"
            :maxlength="100"
            name="placeholder"
            :type="'textarea'"
            v-model.trim="settings.placeholder"
            v-validate="'length:2000'">
          </bk-input>
          <template #append>
            <div class="form-error" v-if="errors.has('placeholder')">{{errors.first('placeholder')}}</div>
          </template>
        </grid-item>
      </grid-layout>
    </div>
  </bk-dialog>
</template>

<style lang="scss" scoped>
  .content-layout {
    max-height: 340px;
    padding: 0 24px;
    @include scrollbar-y;
  }

  .field-option-container {
    width: 100%;
    margin-bottom: 20px;
    margin-top: 20px;
    padding: 16px;
    background: #F5F7FB;

    :deep(.label-text) {
      position: relative;
      display: block;
      padding-right: 10px;
      line-height: 36px;
      font-size: 14px;
      @include ellipsis;
    }
  }
</style>
