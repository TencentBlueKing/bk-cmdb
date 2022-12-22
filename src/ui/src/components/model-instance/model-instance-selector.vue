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
  import { computed, defineProps, ref, watch, watchEffect } from 'vue'
  import debounce from 'lodash.debounce'
  import { getModelInstanceOptions } from '@/service/instance/common'

  const props = defineProps({
    value: {
      type: [Array, String],
      default: ''
    },
    objId: String,
    multiple: {
      type: Boolean,
      default: true
    }
  })
  const emit = defineEmits(['input', 'toggle'])

  watch(() => props.objId, () => {
    localValue.value = resetValue()
  })

  const list = ref([])
  const loading = ref(false)
  const selector = ref(null)

  const search = async (keyword) => {
    loading.value = true
    const results = await getModelInstanceOptions(props.objId, keyword, props.value, { page: { limit: 50 } })
    list.value = results
    loading.value = false
  }

  const remoteSearch = debounce(search, 200)

  const getInitValue = () => (props.multiple ? (props.value || []) : (props.value || ''))
  const resetValue = () => (props.multiple ? [] : '')
  const focus = () => selector?.value?.show?.()

  const localValue = computed({
    get() {
      return getInitValue()
    },
    set(values) {
      emit('input', values)
      emit('change', values)
    }
  })

  watchEffect(() => {
    if (props.objId) {
      search()
    }
  })

  const handleToggle = active => emit('toggle', active)

  defineExpose({
    focus
  })
</script>

<template>
  <bk-select
    ref="selector"
    v-bind="$attrs"
    v-model="localValue"
    searchable
    :multiple="multiple"
    :loading="loading"
    :is-tag-width-limit="true"
    :remote-method="remoteSearch"
    @toggle="handleToggle">
    <bk-option v-for="option in list"
      :key="option.id"
      :id="option.id"
      :name="option.name">
    </bk-option>
  </bk-select>
</template>
