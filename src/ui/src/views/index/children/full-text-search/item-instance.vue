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

<template>
  <div class="result-item">
    <div class="result-title" @click="data.linkTo(data.source)">
      <span v-html="`${data.typeName} - ${data.title}`"></span>
    </div>
    <div class="result-desc" @click="data.linkTo(data.source)">
      <div class="desc-item hl" v-html="`${$t('实例ID')}：${getHighlightValue(data.source.bk_inst_id, data)}`"></div>
      <div class="desc-item" v-for="(property, childIndex) in properties" :key="childIndex">
        <div class="hl" v-if="property.bk_property_type !== PROPERTY_TYPES.INNER_TABLE"
          v-html="`${getHighlightValue(property.bk_property_name, data)}：${getText(property, data)}`">
        </div>
        <div class="table-value" v-else>
          <span>{{getHighlightValue(property.bk_property_name, data)}}：</span>
          <i class="bk-cmdb-icon icon-cc-table" />
          <span class="matched-tag" v-if="isMatchedTableValue(property)">命中</span>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
  import { defineComponent, toRefs, computed } from 'vue'
  import { getText, getHighlightValue } from './use-item.js'
  import { PROPERTY_TYPES } from '@/dictionary/property-constants'

  export default defineComponent({
    name: 'item-instance',
    props: {
      data: {
        type: Object,
        default: () => ({})
      },
      propertyMap: {
        type: Object,
        default: () => ({})
      }
    },
    setup(props) {
      const { data, propertyMap } = toRefs(props)

      const properties = computed(() => propertyMap.value[data.value.source.bk_obj_id])

      const isMatchedTableValue = property => data.value?.highlight?.tables?.includes(property.bk_property_id)

      return {
        PROPERTY_TYPES,
        properties,
        getText,
        getHighlightValue,
        isMatchedTableValue
      }
    }
  })
</script>
