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
  import { computed, defineProps, ref, watch, onMounted } from 'vue'
  import debounce from 'lodash.debounce'
  import { getModelInstanceOptions } from '@/service/instance/common'
  import { t } from '@/i18n'
  import store from '@/store'
  import { OPERATION } from '@/dictionary/iam-auth'

  const props = defineProps({
    value: {
      type: [Array, String, Number],
      default: ''
    },
    objId: String,
    multiple: {
      type: Boolean,
      default: true
    }
  })
  const emit = defineEmits(['input', 'toggle'])

  const getInitValue = () => (props.multiple ? (props.value || []) : (props.value || ''))
  const resetValue = () => (props.multiple ? [] : '')

  const list = ref([])
  const loading = ref(false)
  const selector = ref(null)
  const placeholder = ref('')
  const auth = ref({})
  const disabled = ref(false)

  const search = async (keyword) => {
    loading.value = true
    const results = await getModelInstanceOptions(
      props.objId, keyword, props.value,
      { page: { limit: 50 } },
      { globalPermission: false }
    )
    if (results.length === 0) {
      placeholder.value = t('该字段暂无权限配置，点击申请权限')
      localValue.value = props.multiple ? [] : ''
      disabled.value = true
    } else {
      placeholder.value = t('请选择模型实例')
      localValue.value = getInitValue()
      disabled.value = false
    }
    list.value = results
    loading.value = false
  }

  const remoteSearch = debounce(search, 200)

  const localValue = computed({
    get() {
      return getInitValue()
    },
    set(values) {
      emit('input', values)
      emit('change', values)
    }
  })

  const isActive = ref(false)

  onMounted(() => {
    setTimeout(() => {
      selector?.value?.$refs.bkSelectTag.calcOverflow()
    }, 100)
    auth.value = getAuth(props.objId)
  })

  watch(() => props.objId, (cur, prev) => {
    if (cur && cur !== prev) {
      auth.value = getAuth(cur)
      search()
    }

    localValue.value = resetValue()
  })

  if (props.objId) {
    search()
  }

  const handleToggle = (active) => {
    isActive.value = active
    emit('toggle', active)
  }

  const getAuth = (objId) => {
    const relationModel = store.getters['objectModelClassify/getModelById'](objId)
    return { type: OPERATION.R_INST, relation: [relationModel.id] }
  }

  defineExpose({
    focus: () => selector?.value?.show?.()
  })
</script>

<template>
  <cmdb-auth-mask :auth="auth" :authorized="!disabled">
    <div class="model-instance-selector">
      <bk-select
        :class="['selector', { 'active': isActive }]"
        ref="selector"
        v-bind="$attrs"
        v-model="localValue"
        searchable
        :multiple="multiple"
        :placeholder="placeholder"
        :disabled="disabled"
        font-size="normal"
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
    </div>
  </cmdb-auth-mask>

</template>

<style lang="scss" scoped>
    .model-instance-selector {
        position: relative;
        width: 100%;
        height: 32px;
        .selector {
            width: 100%;
            &.active {
                position: absolute;
                z-index: 2;
            }
          }
    }
</style>
