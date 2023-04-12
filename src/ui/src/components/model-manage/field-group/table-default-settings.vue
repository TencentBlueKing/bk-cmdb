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
  import { nextTick, reactive, ref, watch, set, watchEffect } from 'vue'
  import has from 'has'
  import cloneDeep from 'lodash/cloneDeep'
  import IconTextButton from '@/components/ui/button/icon-text-button.vue'
  import PropertyFormElement from '@/components/ui/form/property-form-element.vue'
  import { getPropertyDefaultValue, isShowOverflowTips, isEmptyPropertyValue } from '@/utils/tools.js'

  const props = defineProps({
    defaults: {
      type: Array,
      default: () => []
    },
    headers: {
      type: Array,
      default: () => []
    },
    readonly: {
      type: Boolean,
      default: false
    },
    preview: {
      type: Boolean,
      default: false
    }
  })

  const emit = defineEmits(['update'])

  const list = ref([])

  const editState = reactive({
    index: null,
    row: {}
  })

  const propertyFormEl = ref(null)
  const tableRef = ref(null)
  const refreshKey = ref(Date.now())

  watchEffect(() => {
    list.value = cloneDeep(props.defaults || [])
  })

  const updateValue = () => {
    emit('update', list.value)
  }

  watch(() => props.headers, (headers) => {
    refreshKey.value = Date.now()

    headers.forEach((header) => {
      list.value.forEach((row) => {
        // 不存在的字段先补齐（先创建了行后添加的表头）
        if (!has(row, header.bk_property_id)) {
          row[header.bk_property_id] = getPropertyDefaultValue(header)
        }

        // 刷新默认值（开始默认值为空后来变更了）
        for (const [key, value] of Object.entries(row)) {
          const prop = headers.find(header => header.bk_property_id === key)
          if (prop) {
            row[key] = getPropertyDefaultValue(prop, isEmptyPropertyValue(value) ? undefined : value)
          }
        }
      })
    })

    updateValue()
  }, { deep: true })

  const newRowData = () => {
    const data = {}
    props.headers.forEach((prop) => {
      data[prop.bk_property_id] = getPropertyDefaultValue(prop)
    })
    return data
  }

  const validateAll = async () => {
    // 获得每一个表单元素的校验方法
    const validates = (propertyFormEl.value || [])
      .map(formElement => formElement.$validator.validateAll())

    if (validates.length) {
      const results = await Promise.all(validates)
      return results.every(valid => valid)
    }

    return true
  }

  const exitEdit = () => {
    editState.index = null
    editState.row = {}
  }
  const enterEdit = (index) => {
    editState.index = index
    editState.row = cloneDeep(list.value[index])
    nextTick(() => {
      const component = propertyFormEl.value?.[0].$refs[`component-${props.headers[0].bk_property_id}`]
      component?.focus?.()
    })
  }
  const updateRow = (index) => {
    set(list.value, index, cloneDeep(editState.row))
  }

  const handleClickEdit = (index) => {
    enterEdit(index)
  }
  const handleClickRemove = (index) => {
    list.value.splice(index, 1)
    updateValue()
  }
  const handleClickAdd = async () => {
    if (!await validateAll()) {
      return
    }
    const length = list.value.push(newRowData())

    enterEdit(length - 1)
  }

  const handleConfirm = async (index) => {
    if (!await validateAll()) {
      return
    }
    updateRow(index)
    exitEdit()

    updateValue()
  }
  const handleCancel = () => {
    exitEdit()
  }

  const clickOutsideMiddleware = (event) => {
    const path = event.composedPath ? event.composedPath() : event.path
    return !path?.some?.(node => node.className === 'tippy-popper')
  }
  const handleClickOutside = () => {
    exitEdit()
  }
</script>

<template>
  <div :class="['table-default-settings', { preview: props.preview }]" v-click-outside="{
    handler: handleClickOutside,
    middleware: clickOutsideMiddleware
  }">
    <bk-table
      :key="refreshKey"
      ref="tableRef"
      class="settings-table"
      :data="list"
      :outer-border="props.preview"
      :header-border="false"
      :row-auto-height="true"
      :max-height="344">
      <bk-table-column
        v-for="prop in props.headers"
        :key="prop.bk_property_id"
        :label="prop.bk_property_name"
        :prop="prop.bk_property_id"
        :min-width="$tools.getHeaderPropertyMinWidth(prop, { min: 120 })"
        :show-overflow-tooltip="true">
        <template #default="{ row, $index }">
          <cmdb-property-value
            v-if="$index !== editState.index"
            :is-show-overflow-tips="isShowOverflowTips(prop)"
            :class="'property-value'"
            :value="row[prop.bk_property_id]"
            :property="prop">
          </cmdb-property-value>
          <property-form-element
            v-else
            class="detault-form-el"
            ref="propertyFormEl"
            :property="prop"
            :size="'small'"
            :font-size="'normal'"
            :row="1"
            error-display-type="tooltips"
            v-model="editState.row[prop.bk_property_id]">
          </property-form-element>
        </template>
      </bk-table-column>
      <bk-table-column :label="$t('操作')" width="90" fixed="right" v-if="!props.readonly">
        <template #default="{ $index }">
          <div class="operation-cell">
            <template v-if="$index !== editState.index">
              <bk-button text theme="primary"
                class="action-button"
                @click="handleClickEdit($index)">{{$t('编辑')}}</bk-button>
              <bk-button text theme="primary"
                class="action-button"
                @click="handleClickRemove($index)">{{$t('删除')}}</bk-button>
            </template>
            <template v-else>
              <bk-button text theme="primary"
                class="action-button"
                @click="handleConfirm($index)">确定</bk-button>
              <bk-button text theme="primary"
                class="action-button"
                @click="handleCancel">取消</bk-button>
            </template>
          </div>
        </template>
      </bk-table-column>
      <template #empty v-if="!props.readonly"><icon-text-button :text="$t('新增')" @click="handleClickAdd" /></template>
    </bk-table>
    <div class="table-append" v-if="list.length > 0 && !props.readonly">
      <icon-text-button
        :text="$t('新增')"
        @click="handleClickAdd"
        :disabled="list.length === 50"
        :disabled-tips="$t('最多添加50行')" />
    </div>
  </div>
</template>

<style lang="scss" scoped>
  .table-default-settings {
    &.preview {
      .table-append {
        border: 1px solid #dfe0e5;
        border-top: none;
      }
    }
  }

  .settings-table {
    .operation-cell {
      .action-button {
        & + .action-button {
          margin-left: 4px;
        }
      }
    }

    .detault-form-el {
      :deep(.form-error) {
        position: static;
      }
    }
  }
  .table-append {
    padding: 10px;
    background: #fff;
    font-size: 12px;
  }
</style>
