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
  import { computed, ref, watchEffect } from 'vue'
  import store from '@/store'
  import Loading from '@/components/loading/index.vue'
  import FlexTag from '@/components/ui/flex-tag'
  import { isEmptyPropertyValue } from '@/utils/tools'

  const props = defineProps({
    value: {
      type: [Array, String, Number],
      default: () => ([])
    },
    property: {
      type: Object,
      default: () => ({})
    },
    showOn: String
  })

  defineExpose({
    getCopyValue: () => tagList.value.join('\n') || '--'
  })

  const list = ref([])

  const tagList = computed(() => list.value.map(item => item.full_name.split('/').join(' / ')))

  const requestId = computed(() => `get_department_id_${Array.isArray(props.value) ? props.value.join('_') : String(props.value)}`)

  const isTextStyle = computed(() => props.showOn === 'search')

  const getOrganization = async (value) => {
    const res = await store.dispatch('organization/getDepartment', {
      params: {
        lookup_field: 'id',
        exact_lookups: Array.isArray(value) ? value.join(',') : value
      },
      fromCache: true,
      requestId: requestId.value
    })

    return res?.results ?? []
  }

  watchEffect(async () => {
    if (isEmptyPropertyValue(props.value)) {
      list.value = []
      return
    }
    try {
      const result = await getOrganization(props.value)
      list.value = result
    } catch (error) {
      list.value = []
    }
  })
</script>

<template>
  <div class="org-value">
    <loading :loading="$loading(requestId)">
      <div class="empty" v-if="!list.length">--</div>
      <flex-tag v-else :list="tagList" :is-text-style="isTextStyle"></flex-tag>
    </loading>
  </div>
</template>

<style lang="scss" scoped>
.org-value {
  ::v-deep .loading {
    margin-top: 2px;
  }
}
</style>
