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

<script setup>
  import { reactive, ref, watchEffect, set, nextTick, computed } from 'vue'
  import cloneDeep from 'lodash/cloneDeep'
  import { v4 as uuidv4 } from 'uuid'
  import { t } from '@/i18n'
  import { $bkInfo } from '@/magicbox'
  import { swapItem } from '@/utils/util'
  import GridLayout from '@/components/ui/other/grid-layout.vue'
  import GridItem from '@/components/ui/other/grid-item.vue'
  import FieldGrid from '@/components/model-manage/field-grid.vue'
  import FieldCard from '@/components/model-manage/field-card.vue'
  import FieldSettingForm from '@/components/model-manage/field-group/field-detail/index.vue'
  import { PROPERTY_TYPES } from '@/dictionary/property-constants'
  import { UNIUQE_TYPES } from '@/dictionary/model-constants'
  import UniqueManage from './unique-manage.vue'
  import ModelFieldSelector from './model-field-selector.vue'
  import useField, { unwrapData, excludeFieldType } from './use-field'
  import useUnique from './use-unique'

  const props = defineProps({
    fieldList: {
      type: Array,
      default: () => ([])
    },
    uniqueList: {
      type: Array,
      default: () => ([])
    },
    isCreateMode: {
      type: Boolean,
      default: false
    }
  })

  const emit = defineEmits(['update-field', 'update-unique'])

  const sliderViews = {
    SETTING_FORM: 'settingForm',
    MODEL_FIELD_SELECTOR: 'modelFieldSelector',
  }

  const fieldLocalList = ref([])
  const uniqueLocalList = ref([])

  const settingFormComp = ref(null)
  const uniqueManageComp = ref(null)
  const uniqueTypeComp = ref(null)

  // use方法中参数默认必须是Ref类型
  const { fieldStatus } = useField(props.fieldList, fieldLocalList)
  const { getUniqueByField, clearUniqueByField } =  useUnique(uniqueLocalList)

  const slider = reactive({
    title: '',
    uniqueEnabled: false,
    uniqueType: UNIUQE_TYPES.SINGLE,
    isShow: false,
    beforeClose: null,
    view: '',
    isEditField: false
  })

  watchEffect(() => {
    const fieldList = cloneDeep(props.fieldList || [])
    if (props.isCreateMode && !fieldList.length) {
      fieldList.push({
        id: uuidv4(),
        bk_property_id: 'demo',
        bk_property_name: '示例字段',
        bk_property_type: PROPERTY_TYPES.SINGLECHAR
      })
    }
    fieldLocalList.value = fieldList.map(unwrapData)
    uniqueLocalList.value = cloneDeep(props.uniqueList || [])
  })

  // 只有字段属性的列表
  const pureFieldList = computed(() => fieldLocalList.value.map(item => item.field))

  const handleAddField = () => {
    slider.title = t('新建字段')
    slider.isEditField = false
    slider.isCreateMode = true

    slider.curField = {}
    slider.curFieldSetting = {}

    slider.uniqueEnabled = false
    slider.uniqueType = UNIUQE_TYPES.SINGLE

    slider.isShow = true
    slider.beforeClose = handleSettingSliderBeforeClose
    slider.view = sliderViews.SETTING_FORM
  }
  const handleEditField = (field) => {
    const data = fieldLocalList.value.find(item => item.field.id === field.id)
    if (!data) {
      console.error('error data!')
      return
    }
    slider.title = t('编辑字段')
    slider.isEditField = true
    slider.isCreateMode = fieldStatus.value[field.id].new

    slider.curField = data.field
    slider.curFieldSetting = data.extra

    const { list: fieldUniqueList, type: fieldUniqueType } = getUniqueByField(data.field)

    slider.uniqueEnabled = fieldUniqueList.length > 0
    slider.uniqueType = fieldUniqueType

    slider.isShow = true
    slider.beforeClose = handleSettingSliderBeforeClose
    slider.view = 'settingForm'
  }
  const handleClickImport = () => {
    slider.title = t('从模型导入')
    slider.isShow = true
    slider.beforeClose = handleImportSliderBeforeClose
    slider.view = sliderViews.MODEL_FIELD_SELECTOR
  }

  const syncField = () => {
    emit('update-field', fieldLocalList.value)
  }
  const syncUnique = () => {
    emit('update-unique', uniqueLocalList.value)
  }
  const appendField = (fieldData, extraData = {}) => {
    const data = {
      field: {
        // 在页面中创建的数据，此id键与后台数据有意保持一致为了简化在更新查找时的逻辑
        // 如果是添加导入的字段，此id应该被模型中字段id覆盖
        id: uuidv4(),
        ...fieldData
      },
      extra: extraData
    }
    fieldLocalList.value.push(data)

    return data
  }
  const updateField = (id, fieldData, extraData) => {
    const fieldIndex = fieldLocalList.value.findIndex(item => item.field.id === id)
    const data = {
      field: { id, ...fieldData },
      extra: extraData
    }
    if (~fieldIndex) {
      set(fieldLocalList.value, fieldIndex, data)
    }

    return data
  }
  const setUnique = (uniqueList, currentField) => {
    // TODO: 删除的情况
    uniqueList.forEach((unique) => {
      const uniqueIndex = uniqueLocalList.value.findIndex(item => item.id === unique.id)
      if (~uniqueIndex) {
        set(uniqueLocalList.value, uniqueIndex, {
          ...uniqueLocalList.value[uniqueIndex],
          keys: unique.keys
        })
      } else {
        unique.keys.forEach((key, index) => {
          if (key === -1) {
            unique.keys[index] = currentField.field.id
          }
        })
        uniqueLocalList.value.push({
          id: unique.id,
          keys: unique.keys
        })
      }
    })
  }

  const getFieldUnique = (field) => {
    const { list: fieldUniqueList, type: fieldUniqueType } = getUniqueByField(field)
    const fieldUniqueWithNameList = fieldUniqueList.map(item => ({
      ...item,
      names: item.keys.map(key => pureFieldList.value.find(field => field.id === key)?.bk_property_name)
    }))
    return {
      list: fieldUniqueWithNameList,
      type: fieldUniqueType
    }
  }

  const handleFieldSave = (id, fieldData, extraData) => {
    // 校验“唯一校验”
    const validateUniqueResult = uniqueManageComp.value?.validateAll?.()

    if (!validateUniqueResult) {
      uniqueManageComp.value?.$el?.scrollIntoView?.()
      return
    }

    let currentField = null
    if (id) {
      currentField = updateField(id, fieldData, extraData)
    } else {
      currentField = appendField(fieldData, extraData)
    }

    if (slider.uniqueEnabled) {
      if (slider.uniqueType === UNIUQE_TYPES.UNION) {
        // 当前字段的唯一检验数据
        const fieldUniqueList = uniqueManageComp.value?.uniqueList
        setUnique(fieldUniqueList, currentField)
      } else {
        // 默认情况下，单独唯一同样可以使用隐藏的唯一校验组件，此处将得到一条默认的唯一校验
        const fieldUniqueList = uniqueManageComp.value?.uniqueList

        const { list, type } = getUniqueByField(currentField.field)
        // 当前字段保存前无任何唯一检验数据
        if (!list.length) {
          setUnique(fieldUniqueList, currentField)
        } else if (type === UNIUQE_TYPES.UNION) {
          // 从联合切为单独，清除相关并重新添加一条
          clearUniqueByField(currentField.field)
          setUnique([{
            id: uuidv4(),
            keys: [currentField.field.id]
          }], currentField)
        }
      }
    } else {
      // 关闭唯一校验，清除当前字段相关的唯一校验
      clearUniqueByField(currentField.field)
    }

    syncField()
    syncUnique()

    sliderClose()
  }

  const handleImportSave = (fields) => {
    fields.forEach(field => appendField(field))

    syncField()
    sliderClose()
  }

  const handleUniqueTypeChange = (type) => {
    if (type === UNIUQE_TYPES.UNION) {
      nextTick(() => {
        uniqueManageComp.value?.$el?.scrollIntoView?.()
      })
    }
  }
  const handleUniqueEnabledChange = (enabled) => {
    if (enabled) {
      nextTick(() => {
        uniqueTypeComp.value?.$el?.scrollIntoView?.()
      })
    }
  }

  const handleRemoveField = (field, index) => {
    fieldLocalList.value.splice(index, 1)

    syncField()
  }

  const handleSortChange = (event) => {
    if (!event.moved) {
      return
    }

    const { newIndex, oldIndex } = event.moved
    swapItem(fieldLocalList.value, oldIndex, newIndex)

    syncField()
  }

  let promiseResolver = null
  let promiseRejecter = null
  const uniqueEnabledTogglePreCheck = () => new Promise((resolve, reject) => {
    promiseResolver = resolve
    promiseRejecter = reject
  })
  const handleToggleUniqueEnabledConfirm = () => promiseResolver?.()
  const handleToggleUniqueEnabledCancel = () => promiseRejecter?.()
  const uniqueEnabledClickOutSideMiddleware = event => (!event.target.closest('.tippy-popper'))

  const sliderClose = () => {
    slider.isShow = false
    slider.curField = {}
    slider.beforeClose = null

    slider.uniqueEnabled = false
    slider.uniqueType = UNIUQE_TYPES.SINGLE
  }
  const handleSettingSliderBeforeClose = () => {
    const hasChanged = Object.keys(settingFormComp.value.changedValues).length
    if (hasChanged) {
      return new Promise((resolve) => {
        $bkInfo({
          title: t('确认退出'),
          subTitle: t('退出会导致未保存信息丢失'),
          extCls: 'bk-dialog-sub-header-center',
          confirmFn: () => {
            sliderClose()
            resolve(true)
          },
          cancelFn: () => {
            resolve(false)
          }
        })
      })
    }
    sliderClose()
    return true
  }
  const handleImportSliderBeforeClose = () => {
    sliderClose()
    return true
  }
  const handleSliderHidden = () => {
    sliderClose()
  }
</script>

<template>
  <div class="field-manage">
    <div class="toolbar">
      <cmdb-auth :auth="{ type: $OPERATION.C_FIELD_TEMPLATE }">
        <template #default="{ disabled }">
          <bk-button
            theme="primary"
            :disabled="disabled"
            @click="handleAddField">
            {{$t('添加字段')}}
          </bk-button>
        </template>
      </cmdb-auth>
      <cmdb-auth :auth="{ type: $OPERATION.C_FIELD_TEMPLATE }">
        <template #default="{ disabled }">
          <bk-button
            :disabled="disabled"
            @click="handleClickImport">
            {{$t('从模型导入')}}
          </bk-button>
        </template>
      </cmdb-auth>
      <div class="filter">
        <bk-input
          class="search-input"
          :placeholder="$t('请输入字段名称')"
          :right-icon="'bk-icon icon-search'" />
        <bk-button theme="default" class="unique-button">
          {{$t('唯一校验')}}
          <em class="num">0</em>
        </bk-button>
      </div>
    </div>
    <field-grid
      :field-list="fieldLocalList"
      @sort-change="handleSortChange">
      <template #field-card="{ field, itemClass }">
        <field-card
          :class="itemClass"
          :field="field.field"
          :field-unique="getFieldUnique(field.field)"
          @click-field="handleEditField"
          @remove-field="handleRemoveField">
        </field-card>
      </template>
    </field-grid>

    <bk-sideslider
      ref="sidesliderComp"
      v-transfer-dom
      :width="640"
      :title="slider.title"
      :is-show.sync="slider.isShow"
      :before-close="slider.beforeClose"
      @hidden="handleSliderHidden">
      <template v-if="slider.isShow" #content>
        <!-- 添加/编辑字段 -->
        <field-setting-form v-if="slider.view === sliderViews.SETTING_FORM"
          ref="settingFormComp"
          scene="setting"
          :is-main-line-model="false"
          :is-read-only="false"
          :is-edit-field="slider.isEditField"
          :is-create-mode="slider.isCreateMode"
          :field="slider.curField"
          :field-setting="slider.curFieldSetting"
          :exclude-type="excludeFieldType"
          @confirm="handleFieldSave"
          @cancel="handleSettingSliderBeforeClose">
          <template #append-unique="{ disabled, fieldInfo }">
            <grid-layout class="mt20" mode="form" :gap="24" :font-size="'14px'" :max-columns="1">
              <grid-item :label="$t('设置为唯一校验')">
                <bk-popconfirm
                  v-if="slider.uniqueEnabled"
                  :content="$t('取消字段唯一校验确认提示语')"
                  width="260"
                  trigger="click"
                  :confirm-text="$t('继续')"
                  @confirm="handleToggleUniqueEnabledConfirm"
                  @cancel="handleToggleUniqueEnabledCancel">
                  <bk-switcher
                    v-click-outside="{
                      handler: handleToggleUniqueEnabledCancel,
                      middleware: uniqueEnabledClickOutSideMiddleware
                    }"
                    v-model="slider.uniqueEnabled"
                    theme="primary"
                    :pre-check="uniqueEnabledTogglePreCheck"
                    @change="handleUniqueEnabledChange">
                  </bk-switcher>
                </bk-popconfirm>
                <bk-switcher
                  v-else
                  v-bk-tooltips="{
                    disabled: !disabled,
                    content: $t('需先设置唯一标识和名称，再进行唯一校验的设置')
                  }"
                  :disabled="disabled"
                  v-model="slider.uniqueEnabled"
                  theme="primary"
                  @change="handleUniqueEnabledChange">
                </bk-switcher>
              </grid-item>
              <grid-item required :label="$t('校验类型')" ref="uniqueTypeComp" v-if="slider.uniqueEnabled">
                <bk-radio-group class="full-width-radio" v-model="slider.uniqueType" @change="handleUniqueTypeChange">
                  <bk-radio-button :value="UNIUQE_TYPES.SINGLE">{{$t('单独唯一')}}</bk-radio-button>
                  <bk-radio-button :value="UNIUQE_TYPES.UNION">{{$t('联合唯一')}}</bk-radio-button>
                </bk-radio-group>
              </grid-item>
            </grid-layout>
            <grid-layout v-show="slider.uniqueEnabled && slider.uniqueType === UNIUQE_TYPES.UNION"
              class="unique-option-container"
              mode="form"
              :gap="0"
              :font-size="'14px'"
              :max-columns="1">
              <unique-manage
                type="union"
                :field="fieldInfo"
                :field-list="pureFieldList"
                :unique-list="uniqueLocalList"
                ref="uniqueManageComp">
              </unique-manage>
            </grid-layout>
          </template>
        </field-setting-form>
        <model-field-selector
          v-else-if="slider.view === sliderViews.MODEL_FIELD_SELECTOR"
          :template-field-list="pureFieldList"
          @confirm="handleImportSave"
          @cancel="handleImportSliderBeforeClose">
        </model-field-selector>
      </template>
    </bk-sideslider>
  </div>
</template>

<style lang="scss" scoped>
  .field-manage {
    margin: 24px 108px;

    .toolbar {
      display: flex;
      margin-bottom: 30px;
      .auth-box {
        & + .auth-box {
          margin-left: 8px;
        }
      }

      .filter {
        margin-left: auto;
        display: flex;
        gap: 8px;
      }

      .search-input {
        width: 430px;
      }
      .unique-button {
        .num {
          font-style: normal;
          font-size: 12px;
          background: #F0F1F5;
          border-radius: 2px;
          padding: 0 .5em;
          color: #979BA5;
          margin-left: 2px;
        }
      }
    }
  }

  .unique-option-container {
    width: 100%;
    margin-bottom: 20px;
    margin-top: 20px;
    padding: 16px;
    background: #F5F7FB;
  }

  .full-width-radio {
    display: flex;
    .bk-form-radio-button {
      flex: 1;
      :deep(.bk-radio-button-text) {
        width: 100%;
      }
    }
  }
</style>
