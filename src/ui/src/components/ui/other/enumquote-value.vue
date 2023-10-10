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
  import routerActions from '@/router/actions'
  import { getModelInstanceByIds, getModelInstanceDetailRoute } from '@/service/instance/common'
  import Loading from '@/components/loading/index.vue'
  import FlexTag from '@/components/ui/flex-tag'

  const props = defineProps({
    value: {
      type: [Array, String],
      default: () => ([])
    },
    property: {
      type: Object,
      default: () => ({})
    }
  })

  defineExpose({
    getCopyValue: () => tagList.value.join('\n') || '--'
  })

  const options = computed(() => props.property.option || [])
  const instIds = computed(() => (props.value || []).map(id => Number(id)))
  const modelId = computed(() => options.value?.[0]?.bk_obj_id)

  const list = ref([])
  const requestId = Symbol()
  watchEffect(async () => {
    if (!instIds.value.length) {
      list.value = []
      return
    }

    try {
      const result = await getModelInstanceByIds(modelId.value, instIds.value, { requestId })
      list.value = result
    } catch (error) {
      list.value = []
    }
  })

  const tagList = computed(() => list.value.map(item => item.name))

  const handleGoDetail = (index) => {
    const item = list.value[index]
    const route = getModelInstanceDetailRoute(item.modelId, item.id, item)
    routerActions.open(route)
  }
</script>

<template>
  <div class="enumquote-value">
    <loading :loading="$loading(requestId)">
      <div class="empty" v-if="!list.length">--</div>
      <flex-tag v-else :is-link-style="true" :list="tagList" @click="handleGoDetail"></flex-tag>
    </loading>
  </div>
</template>

<style lang="scss" scoped>
.enumquote-value {
  ::v-deep .loading {
    margin-top: 2px;
  }
}
</style>
