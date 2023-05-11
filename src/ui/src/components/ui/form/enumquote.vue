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
  import { computed, ref, watch } from 'vue'
  import modelInstanceSelector from '@/components/model-instance/model-instance-selector.vue'
  import { t } from '@/i18n'
  import { BUILTIN_MODELS } from '@/dictionary/model-constants.js'
  import isEqual from 'lodash/isEqual'
  import { isEmptyPropertyValue } from '@/utils/tools'

  const props = defineProps({
    value: {
      type: [Array, String, Number],
      default: ''
    },
    multiple: {
      type: Boolean,
      default: true
    },
    autoSelect: {
      type: Boolean,
      default: true
    },
    disabled: {
      type: Boolean,
      default: false
    },
    options: {
      type: Array,
      default() {
        return []
      }
    }
  })

  const emit = defineEmits(['input', 'change', 'on-selected'])
  const instanceSelector = ref(null)

  // 如果初始值是大于1个元素的数组但multiple为false，则设置选择组件的multipl为true，满足编辑时仍然展示原始值的需求
  const initValue = props.value
  const localMultiple = computed(() => {
    if (Array.isArray(initValue) && initValue.length > 1 && !props.multiple) {
      return true
    }
    return props.multiple
  })

  const refModelId = computed(() => props.options.map(item => item.bk_obj_id)?.[0])

  const refModelInstIds = computed({
    get() {
      if (!isEmptyPropertyValue(props.value)) {
        if (!localMultiple.value) {
          return Array.isArray(props.value) ? props.value[0] : props.value
        }
        return props.value
      }

      // 自动选择时取出默认值
      if (props.autoSelect) {
        const defaultValue = props.options.map(item => item.bk_inst_id)
        return localMultiple.value ? defaultValue : defaultValue[0]
      }

      return localMultiple.value ? [] : ''
    },
    set(values) {
      emit('input', values)
      emit('change', values)
      emit('on-selected', values)
    }
  })

  const searchPlaceholder = computed(() => t('请输入xx', { name: t(refModelId.value === BUILTIN_MODELS.HOST ? 'IP' : '名称') }))

  watch(() => props.value, (val) => {
    // 将默认值同步回给上层组件
    if (!isEqual(val, refModelInstIds.value)) {
      emit('input', refModelInstIds.value)
    }
  }, { immediate: true })

  defineExpose({
    focus: () => instanceSelector?.value?.focus()
  })
</script>
<script>
  export default {
    name: 'cmdb-form-enumquote'
  }
</script>

<template>
  <model-instance-selector
    ref="instanceSelector"
    class="form-enumqoute-selector"
    :obj-id="refModelId"
    :placeholder="$t('请选择xx', { name: $t('模型实例') })"
    :search-placeholder="searchPlaceholder"
    :display-tag="true"
    :disabled="disabled"
    :multiple="localMultiple"
    v-model="refModelInstIds">
  </model-instance-selector>
</template>

<style lang="scss" scoped>
  .form-enumqoute-selector {
    width: 100%;
  }
</style>
