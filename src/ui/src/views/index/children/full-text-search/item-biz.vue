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
      <i class="tag-disabled" v-if="data.source.bk_data_status === 'disabled'">{{$t('已归档')}}</i>
    </div>
    <div class="result-desc" v-if="properties" @click="data.linkTo(data.source)">
      <template v-for="(property, childIndex) in properties">
        <div class="desc-item hl"
          :key="childIndex"
          v-if="data.source[property.bk_property_id]"
          v-html="`${property.bk_property_name}：${getText(property, data)}`">
        </div>
      </template>
    </div>
  </div>
</template>

<script>
  import { defineComponent, toRefs, computed } from 'vue'
  import { getText, getHighlightValue } from './use-item.js'

  export default defineComponent({
    name: 'item-biz',
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
      const { propertyMap } = toRefs(props)

      const properties = computed(() => propertyMap.value.biz)

      return {
        properties,
        getText,
        getHighlightValue
      }
    }
  })
</script>
