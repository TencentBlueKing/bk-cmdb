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
  import { reactive, ref, watchEffect, set, getCurrentInstance, nextTick, provide } from 'vue'
  import { t } from '@/i18n'
  import { swapItem } from '@/utils/util'
  import { isEmptyPropertyValue } from '@/utils/tools'
  import { PROPERTY_TYPES, PROPERTY_TYPE_NAMES } from '@/dictionary/property-constants'
  import GridLayout from '@/components/ui/other/grid-layout.vue'
  import GridItem from '@/components/ui/other/grid-item.vue'
  import IconTextButton from '@/components/ui/button/icon-text-button.vue'
  import TableDefaultSettings from '../../table-default-settings.vue'
  import FieldSettingsModel from './field-settings-model.vue'

  const props = defineProps({
    value: {
      type: [Object, String],
      default: () => ({})
    },
    isEditField: {
      type: Boolean,
      default: false
    }
  })

  const emit = defineEmits(['input'])

  const currentInstance = getCurrentInstance().proxy

  const settingsModelState = reactive({
    isShow: false,
    isEdit: false,
    editDataIndex: null,
    formData: {}
  })

  const columns = [
    {
      id: 'bk_property_id',
      label: t('列ID'),
    },
    {
      id: 'bk_property_name',
      label: t('列名称'),
    },
    {
      id: 'bk_property_type',
      label: t('列类型'),
    }
  ]

  const headers = ref([])
  const defaults = ref([])

  watchEffect(() => {
    // 这里不用处理数据引用问题，因为在上层已经clone过
    headers.value = props.value?.header || []
    defaults.value = props.value?.default || []
  })

  const updateValue = () => {
    headers.value.forEach((header) => {
      // 转换成数字格式
      if ([PROPERTY_TYPES.INT, PROPERTY_TYPES.FLOAT].includes(header.bk_property_type)) {
        header.option.min = isEmptyPropertyValue(header.option.min) ? '' : Number(header.option.min)
        header.option.max = isEmptyPropertyValue(header.option.max) ? '' : Number(header.option.max)
      }

      // default字段处理
      if (props.isEditField) {
        header.default = isEmptyPropertyValue(header.default) ? null : header.default
      } else if (isEmptyPropertyValue(header.default)) {
        Reflect.deleteProperty(header, 'default')
      }
    })
    emit('input', {
      header: headers.value,
      default: defaults.value
    })

    nextTick(() => {
      currentInstance.$validator.validateAll()
    })
  }

  const handleClickAddField = () => {
    settingsModelState.isShow = true
    settingsModelState.isEdit = false
    settingsModelState.formData = {}
  }
  const handleClickEditField = (index) => {
    settingsModelState.isShow = true
    settingsModelState.isEdit = true
    settingsModelState.editDataIndex = index
    settingsModelState.formData = headers.value[index]
  }
  const handleClickRemoveField = (index) => {
    headers.value.splice(index, 1)

    updateValue()
  }
  const handleClickUpField = (index) => {
    if (index === 0) {
      return
    }
    swapItem(headers.value, index, index - 1)

    updateValue()
  }
  const handleClickDownField = (index) => {
    if (index === headers.value.length - 1) {
      return
    }
    swapItem(headers.value, index, index + 1)

    updateValue()
  }
  const handleAddField = (data) => {
    headers.value.push(data)
    settingsModelState.isShow = false

    updateValue()
  }
  const handleSaveField = (data) => {
    const newData = { ...headers.value[settingsModelState.editDataIndex], ...data }
    set(headers.value, settingsModelState.editDataIndex, newData)
    settingsModelState.isShow = false

    updateValue()
  }
  const handleUpdateDefaults = (data) => {
    emit('input', {
      header: headers.value,
      default: data
    })
  }

  provide('headers', headers)
</script>

<template>
  <grid-layout mode="form" :gap="36" :font-size="'14px'" :max-columns="1">
    <grid-item
      direction="column"
      required
      :class="['cmdb-form-item', 'form-item', { 'is-error': errors.has('headers') }]">
      <template #label>
        <div class="label-inner">
          <span class="label-text">{{ $t('表格列设置') }}</span>
          <i18n path="共N列" class="count" v-show="headers.length > 0">
            <template #count><em class="num">{{headers.length}}</em></template>
          </i18n>
        </div>
        <input
          v-validate="'required|min_value:1'"
          v-model="headers.length"
          name="headers"
          data-vv-validate-on="change"
          type="hidden">
      </template>
      <bk-table
        class="table-header-settings"
        :data="headers"
        :outer-border="false"
        :header-border="false">
        <bk-table-column
          v-for="column in columns"
          :key="column.id"
          :label="column.label"
          :prop="column.id"
          :width="column.width"
          :show-overflow-tooltip="true">
          <template #default="{ row }">
            <span v-if="column.id === 'bk_property_type'">{{ PROPERTY_TYPE_NAMES[row[column.id]] }}</span>
            <span v-else>{{ row[column.id] }}</span>
          </template>
        </bk-table-column>
        <bk-table-column :label="$t('操作')" min-width="90">
          <template #default="{ $index }">
            <div class="operation-cell">
              <i :title="$t('编辑')"
                class="icon-cc-edit-shape action-button edit-button"
                @click="handleClickEditField($index)"></i>
              <bk-icon :title="$t('移除')" type="delete"
                :class="['action-button', 'del-button']"
                @click="handleClickRemoveField($index)" />
              <bk-icon :title="$t('上移')" type="arrows-up"
                :class="['action-button', 'up-button', { disabled: $index === 0 }]"
                @click="handleClickUpField($index)" />
              <bk-icon :title="$t('下移')" type="arrows-down"
                :class="['action-button', 'down-button', { disabled: $index === headers.length - 1 }]"
                @click="handleClickDownField($index)" />
            </div>
          </template>
        </bk-table-column>
        <template #empty><icon-text-button :text="$t('添加新列')" @click="handleClickAddField" /></template>
        <template #append v-if="headers.length > 0">
          <div class="table-append">
            <icon-text-button
              :text="$t('添加新列')"
              @click="handleClickAddField"
              :disabled="headers.length === 8"
              :disabled-tips="$t('最多添加8列')" />
          </div>
        </template>
      </bk-table>
      <template #append>
        <div class="form-error" v-if="errors.has('headers')">{{$t('请设置表头字段')}}</div>
      </template>
    </grid-item>
    <grid-item
      direction="column"
      :class="['cmdb-form-item', 'form-item', { 'is-error': errors.has('refModelInst') }]"
      :label="$t('默认值')">
      <table-default-settings
        v-show="headers.length > 0"
        :headers="headers"
        :defaults="defaults"
        @update="handleUpdateDefaults" />
      <template #append>
        <span v-show="!headers.length" class="header-empty-tips">{{ $t('请先添加表格列') }}</span>
      </template>
    </grid-item>

    <field-settings-model
      v-model="settingsModelState.isShow"
      :is-edit="settingsModelState.isEdit"
      :is-edit-field="isEditField"
      :form-data="settingsModelState.formData"
      @save="handleSaveField"
      @add="handleAddField">
    </field-settings-model>
  </grid-layout>
</template>

<style lang="scss" scoped>
  .label-inner {
    .count {
      margin-left: 12px;
      font-size: 12px;
      color: $grayColor;
      .num {
        font-style: normal;
      }
    }
  }

  .table-header-settings {
    .operation-cell {
      .up-button,
      .down-button {
        font-size: 20px !important;
        margin-left: 0px !important;
      }
      .up-button {
        margin-left: 8px !important;
      }
      .action-button {
        cursor: pointer;

        &:hover {
          color: $primaryColor;
        }
        & + .action-button {
          margin-left: 8px;
        }

        &.disabled {
          color: $textDisabledColor;
          cursor: not-allowed;
        }
      }
    }

    .table-append {
      padding: 12px;
    }
  }

  .header-empty-tips {
    font-size: 12px;
    color: $grayColor;
  }
</style>
