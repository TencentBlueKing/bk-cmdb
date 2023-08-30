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
  import { computed } from 'vue'
  import FlexTag from '@/components/ui/flex-tag'

  const props = defineProps({
    value: {
      type: [String, Object, Array],
      default: () => ({})
    }
  })

  const tags = computed(() => {
    if (!props.value) {
      return []
    }

    let list = props.value
    if (!Array.isArray(props.value)) {
      list = [props.value]
    }

    const labels = []
    list.filter(item => item).forEach((item) => {
      labels.push(...Object.keys(item).map(key => `${key}: ${item[key]}`))
    })
    return labels
  })

  defineExpose({
    getCopyValue: () => tags.value.join('\n') || '--'
  })
</script>

<template>
  <div class="mapstring-value">
    <div class="empty" v-if="!tags.length">--</div>
    <flex-tag v-else :list="tags"></flex-tag>
  </div>
</template>
