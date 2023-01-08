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
  import modelInstanceSelector from '@/components/model-instance/model-instance-selector.vue'
  import { t } from '@/i18n'
  import { BUILTIN_MODELS } from '@/dictionary/model-constants.js'
  import { getModelInstanceByIds } from '@/service/instance/common'

  const props = defineProps({
    value: {
      type: [Array, String, Number],
      default: () => ([])
    },
    options: {
      type: Array,
      default() {
        return []
      }
    },
    displayType: {
      type: String,
      default: 'selector',
      validator(type) {
        return ['selector', 'info'].includes(type)
      }
    }
  })

  const emit = defineEmits(['input', 'change'])

  const refModelId = computed(() => props.options.map?.(item => item.bk_obj_id)?.[0])
  const refModelInstIds = computed({
    get() {
      return props.value?.map?.(val => Number(val)) ?? []
    },
    set(values) {
      emit('input', values)
      emit('change', values)
    }
  })
  const isInfoType = computed(() => props.displayType === 'info')

  const searchPlaceholder = computed(() => t('请输入xx', { name: t(refModelId.value === BUILTIN_MODELS.HOST ? 'IP' : '名称') }))

  const infoValue = ref('')
  watchEffect(async () => {
    if (!refModelInstIds.value.length || !isInfoType.value) {
      return
    }
    const result = await getModelInstanceByIds(refModelId.value, refModelInstIds.value)
    infoValue.value = result?.map?.(item => item.name)?.join(' | ') || '--'
  })
</script>
<script>
  export default {
    name: 'cmdb-search-enumquote'
  }
</script>

<template>
  <div v-if="isInfoType">
    <slot name="info-prepend"></slot>
    {{infoValue}}
  </div>
  <model-instance-selector
    v-else
    :obj-id="refModelId"
    :placeholder="$t('请选择xx', { name: $t('模型实例') })"
    :search-placeholder="searchPlaceholder"
    :display-tag="false"
    :multiple="true"
    v-model="refModelInstIds">
  </model-instance-selector>
</template>
