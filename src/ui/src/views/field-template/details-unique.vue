<!--
 * Tencent is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2017 Tencent. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
-->

<script setup>
  import { ref, computed, watchEffect } from 'vue'
  import { getUniqueName } from './children/use-unique'
  import fieldTemplateService from '@/service/field-template'
  import { DETAILS_FIELDLIST_REQUEST_ID } from './children/use-field'
  import { t } from '@/i18n'

  const props = defineProps({
    templateId: {
      type: Number
    },
    uniqueList: {
      type: Array,
      default: () => ([])
    }
  })

  const filterWord = ref('')

  const fieldList = ref([])
  const stuff = ref({ type: 'default', payload: { emptyText: t('bk.table.emptyText') } })

  watchEffect(async () => {
    const fieldData = await fieldTemplateService.getFieldList({
      bk_template_id: props.templateId
    }, { requestId: DETAILS_FIELDLIST_REQUEST_ID, fromCache: true })
    fieldList.value = fieldData?.info || []
  })

  const displayUniqueList = computed(() => {
    const list = props.uniqueList.map(unique => ({
      ...unique,
      name: getUniqueName(unique, fieldList.value, true, 'id')
    }))
    if (filterWord.value) {
      stuff.value.type = 'search'
      const reg = new RegExp(filterWord.value, 'i')
      return list.filter(unique => reg.test(unique.name))
    }
    stuff.value.type = 'default'
    return list
  })

  const handleClearFilter = () => {
    filterWord.value = ''
    stuff.value.type = 'default'
  }
</script>

<template>
  <div class="details-unique">
    <bk-input
      class="search-input"
      v-model="filterWord"
      :placeholder="$t('请输入字段名称')"
      :right-icon="'bk-icon icon-search'" />
    <bk-table class="data-table"
      :data="displayUniqueList"
      :max-height="$APP.height - 229">
      <bk-table-column
        :label="$t('唯一校验')"
        prop="name"
        show-overflow-tooltip>
      </bk-table-column>
      <bk-table-column
        :label="$t('创建时间')"
        prop="create_time"
        show-overflow-tooltip>
        <template #default="{ row }">
          <span>{{$tools.formatTime(row.create_time)}}</span>
        </template>
      </bk-table-column>
      <cmdb-table-empty
        slot="empty"
        :stuff="stuff"
        @clear="handleClearFilter"
      ></cmdb-table-empty>
    </bk-table>
  </div>
</template>

<style lang="scss" scoped>
.details-unique {
  padding: 16px 8px;

  .search-input {
    width: 320px;
    margin-bottom: 16px;
  }
}
</style>
