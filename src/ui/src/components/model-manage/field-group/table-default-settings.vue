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
  import { reactive, ref } from 'vue'
  import IconTextButton from '@/components/ui/button/icon-text-button.vue'
  import PropertyFormElement from '@/components/ui/form/property-form-element.vue'
  import { getPropertyDefaultValue, isShowOverflowTips } from '@/utils/tools.js'

  const props = defineProps({
    headers: {
      type: Array,
      default: () => []
    }
  })

  const list = ref([])
  const editState = reactive({
    rowIndex: null,
  })
  const propertyFormEl = ref(null)

  const newRowData = () => {
    const data = {}
    props.headers.forEach((prop) => {
      data[prop.bk_property_id] = getPropertyDefaultValue(prop)
    })
    return data
  }

  const validateAll = async () => {
    // 获得每一个表单元素的校验方法
    const validates = propertyFormEl.value || []
      .map(formElement => formElement.$validator.validateAll())

    if (validates.length) {
      const results = await Promise.all(validates)
      return results.every(valid => valid)
    }

    return true
  }

  const handleClickEdit = (index) => {
    console.log('handleClickEdit', index)
    editState.rowIndex = index
  }
  const handleClickRemove = (index) => {
    console.log('handleClickRemove', index)
  }
  const handleClickAdd = () => {
    validateAll()
    const length = list.value.push(newRowData())
    editState.rowIndex = length - 1
  }
  const handleChange = (value, property) => {
    console.log(value, property, '--change')
    // emit('change', property, value)
  }
</script>

<template>
  <div class="table-default-settings">
    <bk-table
      class="settings-table"
      :data="list"
      :outer-border="false"
      :header-border="false"
      :max-height="344">
      <bk-table-column
        v-for="prop in props.headers"
        :key="prop.bk_property_id"
        :label="prop.bk_property_name"
        :prop="prop.bk_property_id"
        :min-width="$tools.getHeaderPropertyMinWidth(prop, { min: 90 })"
        :show-overflow-tooltip="true">
        <template #default="{ row, $index }">
          <cmdb-property-value
            v-if="$index !== editState.rowIndex"
            :is-show-overflow-tips="isShowOverflowTips(prop)"
            :class="'property-value'"
            :value="row[prop.bk_property_id]"
            :property="prop">
          </cmdb-property-value>
          <property-form-element
            v-else
            ref="propertyFormEl"
            :property="prop"
            :size="'small'"
            :font-size="'normal'"
            v-model="row[prop.bk_property_id]"
            @change="handleChange">
          </property-form-element>
        </template>
      </bk-table-column>
      <bk-table-column :label="$t('操作')" width="90" fixed="right">
        <template #default="{ $index }">
          <div class="operation-cell">
            <template v-if="$index !== editState.rowIndex">
              <i :title="$t('编辑')"
                class="icon-cc-edit-shape action-button edit-button"
                @click="handleClickEdit($index)"></i>
              <bk-icon :title="$t('移除')" type="delete"
                class="action-button del-button"
                @click="handleClickRemove($index)" />
            </template>
            <template v-else>
              <bk-button text theme="primary">确定</bk-button>
              <bk-button text theme="primary">取消</bk-button>
            </template>
          </div>
        </template>
      </bk-table-column>
      <template #empty><icon-text-button :text="$t('新增')" @click="handleClickAdd" /></template>
    </bk-table>
    <div class="table-append" v-if="list.length > 0">
      <icon-text-button :text="$t('新增')" @click="handleClickAdd" />
    </div>
  </div>
</template>

<style lang="scss" scoped>
  .settings-table {
    .operation-cell {
      .action-button {
        cursor: pointer;

        &:hover {
          color: $primaryColor;
        }
        & + .action-button {
          margin-left: 8px;
        }
      }
    }
  }
  .table-append {
    padding: 12px;
    background: #fff;
    font-size: 12px;
  }
</style>
