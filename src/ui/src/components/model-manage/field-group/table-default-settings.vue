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
  import { nextTick, reactive, ref, watch, set, watchEffect, getCurrentInstance } from 'vue'
  import has from 'has'
  import cloneDeep from 'lodash/cloneDeep'
  import IconTextButton from '@/components/ui/button/icon-text-button.vue'
  import PropertyFormElement from '@/components/ui/form/property-form-element.vue'
  import { getPropertyDefaultValue, isShowOverflowTips, getPropertyDefaultEmptyValue, formatValues } from '@/utils/tools.js'
  import { keyupCallMethod } from '@/utils/util'

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

  const instacne = getCurrentInstance().proxy

  const list = ref([])

  const editState = reactive({
    index: [],
    row: {}
  })

  const tableRef = ref(null)
  const tableEmptyAddButtonRef = ref(null)
  const refreshKey = ref(Date.now())
  const addIndex = ref(-1)

  const scrollAddButton = () => {
    if (tableEmptyAddButtonRef.value) {
      tableEmptyAddButtonRef.value.$el?.closest('.bk-table-empty-text')?.scrollIntoView?.()
    }
  }

  watchEffect(() => {
    list.value = cloneDeep(props.defaults || [])
  })

  const updateValue = () => {
    const formattedList = list.value.map(row => formatValues(row, props.headers))
    emit('update', formattedList)
  }

  watch(() => props.headers, (headers) => {
    refreshKey.value = Date.now()

    // 从数据中去掉header中已经删除的列
    list.value.forEach((row) => {
      Object.keys(row).forEach((key) => {
        if (!headers.some(header => header.bk_property_id === key)) {
          Reflect.deleteProperty(row, key)
        }
      })
    })

    headers.forEach((header) => {
      list.value.forEach((row) => {
        // 不存在的字段先补齐（先创建了行后添加的表头）
        if (!has(row, header.bk_property_id)) {
          // 这里按需求初始化的是默认空值，而非表头的默认值
          row[header.bk_property_id] = getPropertyDefaultEmptyValue(header)
        }
      })
    })

    updateValue()

    nextTick(scrollAddButton)
  }, { deep: true })

  const newRowData = () => {
    const data = {}
    props.headers.forEach((prop) => {
      data[prop.bk_property_id] = getPropertyDefaultValue(prop)
    })
    return data
  }

  const validateAll = async (index) => {
    // 获得每一个表单元素的校验方法
    const validates = (instacne.$refs[`property-form-el-${index}`] || [])
      .map(formElement => formElement.$validator.validateAll())

    if (validates.length) {
      const results = await Promise.all(validates)
      return results.every(valid => valid)
    }

    return true
  }

  const exitEdit = (index) => {
    const dataIndex = editState.index.findIndex(i => i === index)
    if (dataIndex !== -1) {
      editState.index.splice(dataIndex, 1)
      set(editState.row, index, {})
    }

    // 退出“新增一行”
    if (index === addIndex.value) {
      addIndex.value = -1
    }
  }
  const enterEdit = (index) => {
    editState.index.push(index)
    set(editState.row, index, cloneDeep(list.value[index]))
    nextTick(() => {
      const component = instacne.$refs[`property-form-el-${index}`]?.[0].$refs[`component-${props.headers[0].bk_property_id}`]
      component?.focus?.()
    })
  }
  const updateRow = (index) => {
    set(list.value, index, cloneDeep(editState.row[index]))
  }

  const handleClickEdit = (index) => {
    enterEdit(index)
  }
  const handleClickRemove = (index) => {
    list.value.splice(index, 1)
    updateValue()

    if (list.value.length === 0) {
      nextTick(scrollAddButton)
    }
  }
  const handleClickAdd = async () => {
    const length = list.value.push(newRowData())

    enterEdit(length - 1)
    addIndex.value = length - 1
  }

  const handleConfirm = async (index) => {
    if (!await validateAll(index)) {
      return
    }
    updateRow(index)
    exitEdit(index)

    updateValue()
  }
  const handleCancel = (index) => {
    // 取消的是新增的那一行
    if (index === addIndex.value) {
      list.value.splice(addIndex.value, 1)
      if (list.value.length === 0) {
        nextTick(scrollAddButton)
      }
    }
    exitEdit(index)
  }
</script>

<template>
  <div :class="['table-default-settings', { preview: props.preview }]">
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
        :label="$tools.getHeaderPropertyName(prop)"
        :prop="prop.bk_property_id"
        :min-width="$tools.getHeaderPropertyMinWidth(prop, { min: 120 })">
        <template #default="{ row, $index }">
          <cmdb-property-value
            v-if="!editState.index.includes($index)"
            :is-show-overflow-tips="isShowOverflowTips(prop)"
            :class="'property-value'"
            :value="row[prop.bk_property_id]"
            :property="prop">
          </cmdb-property-value>
          <property-form-element
            v-else
            :class="['detault-form-el', prop.bk_property_type]"
            :ref="`property-form-el-${$index}`"
            :property="prop"
            :size="'small'"
            :font-size="'normal'"
            :row="1"
            :must-required="false"
            error-display-type="tooltips"
            @keyup.native="(event) => keyupCallMethod(event, () => handleConfirm($index))"
            v-model="editState.row[$index][prop.bk_property_id]">
          </property-form-element>
        </template>
      </bk-table-column>
      <bk-table-column :label="$t('操作')" width="130" fixed="right" v-if="!props.readonly">
        <template #default="{ $index }">
          <div class="operation-cell">
            <template v-if="!editState.index.includes($index)">
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
                @click="handleConfirm($index)">{{$t('确定')}}</bk-button>
              <bk-button text theme="primary"
                class="action-button"
                @click="handleCancel($index)">{{$t('取消')}}</bk-button>
            </template>
          </div>
        </template>
      </bk-table-column>
      <template #empty v-if="!props.readonly">
        <icon-text-button
          ref="tableEmptyAddButtonRef"
          class="table-empty-add-button"
          :text="$t('新增')"
          @click="handleClickAdd" />
      </template>
    </bk-table>
    <div class="table-append" v-if="list.length > 0 && !props.readonly && addIndex === -1">
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

    .table-empty-add-button {
      // 避免遮挡
      position: relative;
      z-index: 1;
    }
  }

  .settings-table {
    &:focus-within {
      &.bk-table-scrollable-x,
      &.bk-table-scrollable-y {
        overflow: auto !important;
        :deep(.bk-table-body-wrapper) {
          overflow: auto !important;
        }
      }

      overflow: visible !important;
      :deep(.bk-table-body-wrapper) {
        overflow: visible !important;
      }
    }
    .operation-cell {
      .action-button {
        & + .action-button {
          margin-left: 4px;
        }
      }
    }

    .detault-form-el {
      &:focus-within {

        &.longchar {
          position: absolute;
          left: -1px;
          top: 2px;
          z-index: 1;
          :deep(.bk-form-textarea) {
            min-height: 90px !important;
          }
          :deep(.control-icon) {
            display: none;
          }
        }
      }
    }
  }
  .table-append {
    padding: 10px;
    background: #fff;
    font-size: 12px;
  }
</style>
