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
  import { ref, computed, toRefs, watchEffect } from 'vue'
  import fieldTemplateService from '@/service/field-template'
  import FieldCard from '@/components/model-manage/field-card.vue'
  import useUnique from './children/use-unique'
  import { DETAILS_FIELDLIST_REQUEST_ID } from './children/use-field'

  const props = defineProps({
    templateId: {
      type: Number
    },
    uniqueList: {
      type: Array,
      default: () => ([])
    }
  })

  const emit = defineEmits(['updated'])

  const { uniqueList: templateUniqueList } = toRefs(props)

  const filterWord = ref('')

  const { getFieldUniqueWithName } =  useUnique(templateUniqueList, templateUniqueList)

  const fieldList = ref([])

  watchEffect(async () => {
    const fieldData = await fieldTemplateService.getFieldList({
      bk_template_id: props.templateId
    }, { requestId: DETAILS_FIELDLIST_REQUEST_ID })
    fieldList.value = fieldData?.info || []

    emit('updated', fieldList.value)
  })

  const displayFieldList = computed(() => {
    if (filterWord.value) {
      const reg = new RegExp(filterWord.value, 'i')
      return fieldList.value.filter(field => reg.test(field.bk_property_name))
    }
    return fieldList.value
  })

  const handleClearFilter = () => {
    filterWord.value = ''
  }
</script>

<template>
  <div class="details-field" v-bkloading="{ isLoading: $loading(DETAILS_FIELDLIST_REQUEST_ID) }">
    <bk-input
      class="search-input"
      v-model="filterWord"
      :placeholder="$t('请输入字段名称')"
      :right-icon="'bk-icon icon-search'" />
    <div class="field-list" v-if="displayFieldList.length">
      <field-card
        v-for="(field, index) in displayFieldList"
        class="details-field-card"
        :key="index"
        :field="field"
        :field-unique="getFieldUniqueWithName(field, fieldList)"
        :sortable="false"
        :deletable="false">
      </field-card>
    </div>
    <cmdb-table-empty
      v-else
      slot="empty"
      :stuff="{ type: 'search' }"
      @clear="handleClearFilter"
    ></cmdb-table-empty>
  </div>
</template>

<style lang="scss" scoped>
.details-field {
  padding: 16px 0;

  .search-input {
    width: 320px;
    margin-bottom: 18px;
    margin-left: 8px;
  }
  .field-list {
    display: grid;
    gap: 16px;
    grid-template-columns: repeat(auto-fill, minmax(276px, 1fr));
    width: 100%;
    padding: 0 8px;
    align-content: flex-start;
    max-height: calc(100vh - 278px - 50px);
    @include scrollbar-y;

    .details-field-card {
      border: 1px solid #DCDEE5;
      box-shadow: none;
    }
  }
}
</style>
