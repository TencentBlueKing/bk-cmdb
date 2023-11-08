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
  import { reactive, ref, watchEffect, set, nextTick, computed, toRefs } from 'vue'
  import cloneDeep from 'lodash/cloneDeep'
  import { v4 as uuidv4 } from 'uuid'
  import { t } from '@/i18n'
  import { $bkInfo, $error } from '@/magicbox'
  import { swapItem, escapeRegexChar } from '@/utils/util'
  import GridLayout from '@/components/ui/other/grid-layout.vue'
  import GridItem from '@/components/ui/other/grid-item.vue'
  import FieldGrid from '@/components/model-manage/field-grid.vue'
  import FieldCard from '@/components/model-manage/field-card.vue'
  import FieldSettingForm from '@/components/model-manage/field-group/field-detail/index.vue'
  import { UNIUQE_TYPES } from '@/dictionary/model-constants'
  import { EDITABLE_TYPES, REQUIRED_TYPES } from '@/dictionary/property-constants'
  import UniqueManage from './unique-manage.vue'
  import ModelFieldSelector from './model-field-selector.vue'
  import UniqueManageDrawer from './unique-manage-drawer.vue'
  import useField, { unwrapData, excludeFieldType, isFieldExist } from './use-field'
  import useUnique, { singleRuleTypes, unionRuleTypes } from './use-unique'
  import RouterQuery from '@/router/query'

  const props = defineProps({
    fieldList: {
      type: Array,
      default: () => ([])
    },
    uniqueList: {
      type: Array,
      default: () => ([])
    },
    beforeFieldList: {
      type: Array,
      default: () => ([])
    },
    beforeUniqueList: {
      type: Array,
      default: () => ([])
    },
    isCreateMode: {
      type: Boolean,
      default: false
    }
  })

  const emit = defineEmits(['update-field', 'update-unique'])

  const { beforeFieldList: oldFieldList, beforeUniqueList: oldUniqueList } = toRefs(props)

  const sliderViews = {
    SETTING_FORM: 'settingForm',
    MODEL_FIELD_SELECTOR: 'modelFieldSelector',
  }

  const fieldLocalList = ref([])
  const uniqueLocalList = ref([])

  const settingFormComp = ref(null)
  const modelFormComp = ref(null)
  const uniqueManageComp = ref(null)
  const uniqueTypeComp = ref(null)

  // use方法中参数默认必须是Ref类型
  const { fieldStatus, removedFieldList } = useField(oldFieldList, fieldLocalList)
  const { getUniqueByField, clearUniqueByField } =  useUnique(oldUniqueList, uniqueLocalList)

  const defaultFieldSetting = () => ({
    lock: {
      isrequired: true,
      editable: true,
      placeholder: true
    }
  })
  const slider = reactive({
    title: '',
    uniqueEnabled: false,
    uniqueType: UNIUQE_TYPES.SINGLE,
    isShow: false,
    beforeClose: null,
    view: '',
    isEditField: false,
    curFieldSetting: defaultFieldSetting()
  })

  const uniqueDrawerOpen = ref(false)
  const stuff = ref({ type: 'default' })

  const query = computed(() => RouterQuery.getAll())

  watchEffect(() => {
    const fieldList = cloneDeep(props.fieldList || [])
    fieldLocalList.value = fieldList.map(unwrapData)
    uniqueLocalList.value = cloneDeep(props.uniqueList || [])
    const { action } = query?.value ?? {}
    if (action === 'openUnqiueDrawer') {
      handleOpenUnqiueDrawer()
    }
  })

  // 只有字段属性的列表
  const pureFieldList = computed(() => fieldLocalList.value.map(item => item.field))

  const filterWord = ref('')
  const displayFieldLocalList = computed(() => {
    if (filterWord.value) {
      stuff.value.type = 'search'
      const reg = new RegExp(escapeRegexChar(filterWord.value), 'i')
      return fieldLocalList.value.filter(item => reg.test(item.field.bk_property_name))
    }
    stuff.value.type = 'default'
    return fieldLocalList.value
  })

  const handleAddField = () => {
    slider.title = t('新建字段')
    slider.isEditField = false
    slider.isCreateMode = true

    slider.curField = {}
    slider.curFieldSetting = defaultFieldSetting()

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
  const appendField = (fieldData) => {
    const data = {
      field: {
        // 在页面中创建的数据，此id键与后台数据有意保持一致为了简化在更新查找时的逻辑
        // 如果是添加导入的字段，此id应该被模型中字段id覆盖
        id: uuidv4(),
        ...fieldData
      },
      extra: slider.curFieldSetting
    }
    fieldLocalList.value.push(data)
    syncField()
    return data
  }
  const updateField = (id, fieldData) => {
    const fieldIndex = fieldLocalList.value.findIndex(item => item.field.id === id)
    const data = {
      field: { id, ...fieldData },
      extra: slider.curFieldSetting
    }
    if (~fieldIndex) {
      set(fieldLocalList.value, fieldIndex, data)
    }

    return data
  }
  const setUnique = (uniqueList, currentField) => {
    uniqueList.forEach((unique) => {
      const uniqueIndex = uniqueLocalList.value.findIndex(item => item.id === unique.id)
      if (~uniqueIndex) {
        set(uniqueLocalList.value, uniqueIndex, {
          ...uniqueLocalList.value[uniqueIndex],
          keys: unique.keys
        })
      } else {
        unique.keys.forEach((key, index) => {
          // 实时创建的字段，创建唯一校验时同时创建的字段
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

  const isUniqueTypeDisabled = (fieldInfo, type) => {
    if (type === UNIUQE_TYPES.SINGLE) {
      return !singleRuleTypes.includes(fieldInfo.bk_property_type)
    }
    if (type === UNIUQE_TYPES.UNION) {
      return !unionRuleTypes.includes(fieldInfo.bk_property_type)
    }
  }

  const handleFieldSave = async (id, fieldData) => {
    if (isFieldExist(fieldData, [...fieldLocalList.value, ...removedFieldList.value], id)) {
      $error(t('字段已在模板中存在，无法添加'))
      return
    }

    // 启用了唯一校验并且类型是联合唯一才需要校验，单独唯一不校验在后续的处理中直接对数据进行修改
    if (slider.uniqueEnabled && slider.uniqueType === UNIUQE_TYPES.UNION) {
      // 校验“唯一校验”
      const validateUniqueResult = await uniqueManageComp.value?.isValid?.()

      if (!validateUniqueResult) {
        uniqueManageComp.value?.$el?.scrollIntoView?.()
        return
      }
    }

    let currentField = null
    if (id) {
      currentField = updateField(id, fieldData)
    } else {
      currentField = appendField(fieldData)
    }

    if (slider.uniqueEnabled) {
      if (slider.uniqueType === UNIUQE_TYPES.UNION) {
        // 当前字段的唯一检验数据
        const fieldUniqueList = uniqueManageComp.value?.getUniqueList()
        const { list } = getUniqueByField(currentField.field)
        list.forEach((unique) => {
          // 找不到了
          if (!fieldUniqueList.some(item => item.id === unique.id)) {
            // 在原列表中的index
            const oldIndex = uniqueLocalList.value.findIndex(item => item.id === unique.id)
            if (~oldIndex) {
              uniqueLocalList.value.splice(oldIndex, 1)
            }
          }
        })
        setUnique(fieldUniqueList, currentField)
      } else {
        // 默认情况下，单独唯一同样可以使用隐藏的唯一校验组件，此处将得到一条默认的唯一校验
        const fieldUniqueList = uniqueManageComp.value?.getUniqueList()

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
  const handleUniqueEnabledChange = (enabled, fieldInfo) => {
    if (enabled) {
      if (singleRuleTypes.includes(fieldInfo.bk_property_type)) {
        slider.uniqueType = UNIUQE_TYPES.SINGLE
      } else if (unionRuleTypes.includes(fieldInfo.bk_property_type)) {
        slider.uniqueType = UNIUQE_TYPES.UNION
      }
      nextTick(() => {
        uniqueTypeComp.value?.$el?.scrollIntoView?.()
      })
    }
  }

  const handleChangeUnique = (uniqueList) => {
    // 先过滤掉keys为空的数据
    uniqueList.filter(unique => unique.keys.length).forEach((unique) => {
      const { id: uniqueId, keys: uniqueKeys } = unique
      const uniqueIndex = uniqueLocalList.value.findIndex(item => item.id === uniqueId)
      // 添加
      if (!~uniqueIndex) {
        uniqueLocalList.value.push({
          id: uniqueId,
          keys: uniqueKeys
        })
      } else {
        // 更新
        set(uniqueLocalList.value, uniqueIndex, {
          ...uniqueLocalList.value[uniqueIndex],
          keys: uniqueKeys
        })
      }
    })

    // 在原数据中查找编辑后的数据，找不到即为删除
    uniqueLocalList.value.forEach((unique, index) => {
      if (!uniqueList.some(item => item.id === unique.id)) {
        uniqueLocalList.value.splice(index, 1)
      }
    })

    syncUnique()
  }

  const handleRemoveField = (field) => {
    const index = fieldLocalList.value.findIndex(item => item.field.id === field.id)
    if (~index) {
      fieldLocalList.value.splice(index, 1)
      syncField()
    }
  }

  const handleSortChange = (event) => {
    if (!event.moved) {
      return
    }

    const { newIndex, oldIndex } = event.moved
    swapItem(fieldLocalList.value, oldIndex, newIndex)

    syncField()
  }

  const handleRecover = (removedField) => {
    if (isFieldExist(removedField, fieldLocalList.value)) {
      $error(t('字段已在模板中存在，无法恢复'))
      return
    }
    const oriField = oldFieldList.value.find(item => item.id === removedField.id)
    const { field, extra } = unwrapData(oriField)
    appendField(field, extra)
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
  const beforeCloseDialog = () => new Promise((resolve) => {
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

  const handleSettingSliderBeforeClose = () => {
    const hasChanged = Object.keys(settingFormComp.value.changedValues).length
    if (hasChanged) {
      return beforeCloseDialog()
    }
    sliderClose()
    return true
  }
  const handleImportSliderBeforeClose = () => {
    if (modelFormComp.value?.selectedModelId) {
      return beforeCloseDialog()
    }
    sliderClose()
    return true
  }
  const handleSliderHidden = () => {
    sliderClose()
  }
  const handleOpenUnqiueDrawer = () => {
    uniqueDrawerOpen.value = true
  }
  const handleUniqueDrawerClose = () => {
    uniqueDrawerOpen.value = false
  }
  const handleClearFilter = () => {
    stuff.value.type = 'default'
    filterWord.value = ''
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
          v-model="filterWord"
          :placeholder="$t('请输入字段名称')"
          :right-icon="'bk-icon icon-search'" />
        <bk-button theme="default" class="unique-button" @click="handleOpenUnqiueDrawer">
          {{$t('唯一校验')}}
          <em class="num">{{uniqueLocalList.length}}</em>
        </bk-button>
      </div>
    </div>
    <field-grid
      :field-list="displayFieldLocalList"
      @sort-change="handleSortChange">
      <template #field-card="{ field, itemClass }">
        <field-card
          :class="[itemClass, 'field-card-container']"
          :field="field.field"
          :field-unique="getFieldUnique(field.field)"
          :remove-disabled="fieldLocalList.length === 1"
          :remove-disabled-tips="$t('模板至少需要一个字段')"
          @click-field="handleEditField(field.field)"
          @remove-field="handleRemoveField(field.field)">
          <template #flag-append v-if="!isCreateMode">
            <div class="flag-append" v-if="fieldStatus[field.field.id].new || fieldStatus[field.field.id].changed">
              <span class="flag-tag new" v-if="fieldStatus[field.field.id].new">
                <em class="tag-text">{{$t('新增')}}</em>
              </span>
              <span class="flag-tag changed" v-else-if="fieldStatus[field.field.id].changed">
                <em class="tag-text">{{$t('更新')}}</em>
              </span>
            </div>
          </template>
        </field-card>
      </template>
    </field-grid>

    <div class="removed-container" v-if="removedFieldList.length">
      <div class="removed-title">本次删除的字段（{{removedFieldList.length}}）</div>
      <field-grid
        :field-list="removedFieldList"
        :disabled-sort="true">
        <template #field-card="{ field, itemClass }">
          <field-card
            :class="[itemClass, 'field-card-container', 'removed']"
            :field="field.field"
            :field-unique="getFieldUnique(field.field)"
            :deletable="false"
            :sortable="false">
            <template #flag-append>
              <div class="flag-append">
                <span class="flag-tag changed">
                  <em class="tag-text">{{$t('删除')}}</em>
                </span>
              </div>
            </template>
            <template #action-append="{ field: removedField }">
              <bk-link theme="primary" class="recover-link" @click="handleRecover(removedField)">{{$t('恢复')}}</bk-link>
            </template>
          </field-card>
        </template>
      </field-grid>
    </div>

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
                  v-if="slider.uniqueEnabled && slider.uniqueType === UNIUQE_TYPES.UNION"
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
                    @change="enabled => handleUniqueEnabledChange(enabled, fieldInfo)">
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
                  @change="enabled => handleUniqueEnabledChange(enabled, fieldInfo)">
                </bk-switcher>
              </grid-item>
              <grid-item required :label="$t('校验类型')" ref="uniqueTypeComp" v-if="slider.uniqueEnabled">
                <bk-radio-group class="full-width-radio" v-model="slider.uniqueType" @change="handleUniqueTypeChange">
                  <bk-radio-button
                    :disabled="isUniqueTypeDisabled(fieldInfo, UNIUQE_TYPES.SINGLE)"
                    :value="UNIUQE_TYPES.SINGLE">
                    {{$t('单独唯一')}}
                  </bk-radio-button>
                  <bk-radio-button
                    :disabled="isUniqueTypeDisabled(fieldInfo, UNIUQE_TYPES.UNION)"
                    :value="UNIUQE_TYPES.UNION">
                    {{$t('联合唯一')}}
                  </bk-radio-button>
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
                :before-unique-list="oldUniqueList"
                ref="uniqueManageComp">
              </unique-manage>
            </grid-layout>
          </template>
          <template #append-lock="{ fieldInfo }">
            <div class="lock-option-container">
              <div class="lock-option-title">
                <span class="sub-title">{{ $t('模型应用设置') }}</span>
                {{ $t('模型应用设置提示语') }}
              </div>
              <div class="option-item" v-if="EDITABLE_TYPES.includes(fieldInfo.bk_property_type)">
                <bk-checkbox
                  v-model="slider.curFieldSetting.lock.editable"
                  :true-value="false"
                  :false-value="true"
                  class="checkbox">
                  {{$t('允许在模型中修改“在实例中可编辑”的配置')}}
                </bk-checkbox>
              </div>
              <div class="option-item" v-if="REQUIRED_TYPES.includes(fieldInfo.bk_property_type)">
                <bk-checkbox
                  v-model="slider.curFieldSetting.lock.isrequired"
                  :true-value="false"
                  :false-value="true"
                  class="checkbox">
                  {{$t('允许在模型中修改“设置为必填项”的配置')}}
                </bk-checkbox>
              </div>
              <div class="option-item">
                <bk-checkbox
                  v-model="slider.curFieldSetting.lock.placeholder"
                  :true-value="false"
                  :false-value="true"
                  class="checkbox">
                  {{$t('允许在模型中修改“用户提示”的配置')}}
                </bk-checkbox>
              </div>
            </div>
          </template>
        </field-setting-form>
        <model-field-selector
          v-else-if="slider.view === sliderViews.MODEL_FIELD_SELECTOR"
          :template-field-list="pureFieldList"
          ref="modelFormComp"
          @confirm="handleImportSave"
          @cancel="handleImportSliderBeforeClose">
        </model-field-selector>
      </template>
    </bk-sideslider>

    <unique-manage-drawer
      :open="uniqueDrawerOpen"
      :field-list="pureFieldList"
      :unique-list="uniqueLocalList"
      :before-unique-list="oldUniqueList"
      @close="handleUniqueDrawerClose"
      @change-unique="handleChangeUnique">
    </unique-manage-drawer>


    <cmdb-table-empty
      v-if="!displayFieldLocalList.length"
      slot="empty"
      :stuff="stuff"
      :auth="{ type: $OPERATION.C_FIELD_TEMPLATE }"
      @clear="handleClearFilter">
      <bk-exception class="empty-set" type="empty" scene="part">
        <i18n path="尚未创建字段">
          <template #link>
            <cmdb-auth :auth="{ type: $OPERATION.C_FIELD_TEMPLATE }">
              <template #default="{ disabled }">
                <bk-button
                  text
                  theme="primary"
                  :disabled="disabled"
                  @click="handleAddField">
                  {{$t('立即创建')}}
                </bk-button>
              </template>
            </cmdb-auth>
          </template>
        </i18n>
      </bk-exception>
    </cmdb-table-empty>
  </div>
</template>

<style lang="scss" scoped>
  .field-manage {
    padding: 24px 108px;

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
    padding: 24px 16px;
    background: #F5F7FB;
  }

  .lock-option-container {
    display: flex;
    flex-direction: column;
    gap: 18px;
    border-top: 1px solid #EAEBF0;
    padding-top: 16px;
    margin-top: 20px;
    margin-bottom: 32px;
    .lock-option-title {
      font-size: 14px;
      margin-bottom: 4px;

      .sub-title {
        font-weight: 700;
      }
    }

    .option-item {
      font-size: 14px;
    }
  }

  .full-width-radio {
    display: flex;
    .bk-form-radio-button {
      flex: 1;
      :deep(.bk-radio-button-text) {
        width: 100%;
      }

      &.disabled:first-child {
        :deep(.bk-radio-button-text) {
          border-left: 1px solid currentColor;
        }
      }
    }
  }

  .field-card-container {
    .flag-append {
      margin-left: 2px;
    }
    .flag-tag {
      background: #E4FAF0;
      border-radius: 2px;
      padding: 1px 2px;
      height: 16px;
      line-height: 16px;
      white-space: nowrap;
      display: flex;
      align-items: center;
      position: relative;
      top: -2px;

      .tag-text {
        display: block;
        font-size: 12px;
        font-style: normal;
        transform: scale(.875);
      }
      &.new {
        color: #14A568;
        background: #E4FAF0;
      }
      &.changed {
        color: #FF9C01;
        background: #FFF3E1;
      }
      &.removed {
        color: #EA3636;
        background: #FCE9E8;
      }
    }

    .recover-link {
      visibility: hidden;
      :deep(.bk-link-text) {
        font-size: 12px;
      }
    }

    &:hover {
      .recover-link {
        visibility: visible;
      }
    }

    &.removed {
      opacity: 0.5;

      :deep(.field-name) {
        text-decoration: line-through;
      }

      &:hover {
        opacity: 1;
      }
    }
  }

  .removed-container {
    margin-top: 24px;
    .removed-title {
      font-size: 12px;
      margin-bottom: 16px;
    }
  }
</style>
