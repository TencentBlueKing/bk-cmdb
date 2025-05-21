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
  <bk-org-selector
    :class="[
      'cmdb-organization',
      size,
      {
        'is-focus': focused,
        'is-disabled': disabled,
        'is-unselected': unselected,
        'organization-hide': hidden
      }
    ]"
    v-model="checked"
    :api-base-url="api"
    :tenant-id="tenantId"
    :disabled="disabled"
    :z-index="9999"
    @changeResult="handleChangeResult"
    @closed="handleClose"
    @confirm="handleConfirm"
    v-bind="$attrs"
  >
  </bk-org-selector>
</template>

<script setup>
  import { computed, ref, watch } from 'vue'
  import isEqual from 'lodash/isEqual'
  import BkOrgSelector from '@blueking/bk-org-selector/vue2'
  import '@blueking/bk-org-selector/vue2/vue2.css'

  const props = defineProps({
    value: {
      type: [Array, String, Number],
      default: () => ([])
    },
    disabled: {
      type: Boolean,
      default: false
    },
    readonly: Boolean,
    multiple: {
      type: Boolean,
      default: false
    },
    clearable: Boolean,
    size: String,
    placeholder: {
      type: String,
      default: ''
    },
    zIndex: {
      type: Number,
      default: 2500
    },
    formatter: Function,
    hidden: {
      type: Boolean,
      default: false
    }
  })

  const emit = defineEmits(['on-checked', 'input', 'toggle', 'result-change', 'close', 'confirm'])

  const focused = ref(false)
  const unselected = ref(false)

  const handleChangeResult = (res) => {
    // 组织组件通过item的close按钮取消掉需要手动同步
    const oldVal = checked.value.map(item => item.id)
    const val = (res[0]?.data ?? []).map(item => item.id)
    if (!isEqual(val, oldVal)) {
      checked.value = val
    }
  }
  const handleConfirm  = (res) => {
    emit('confirm', res)
  }
  const handleClose = () => {
    emit('close')
  }

  const api = computed(() => window.Site.userManageUrl)
  const tenantId = computed(() => window.Site.tenantId)
  const checked = computed({
    get() {
      if (this.value && !Array.isArray(this.value)) {
        return [{ id: this.value, type: 'org' }]
      }
      // 需要判断是回显还是用户在机构选择器选择完成
      if (typeof this.value?.[0] === 'number') {
        // 回显
        return this.value.map(item => ({ id: item, type: 'org' })) || []
      }
      return this.value || []
    },
    set(value) {
      let val = value || null
      if (val) {
        val = Array.isArray(value) ? value : [value]
        val = val.map(item => item?.id ?? item)
      }
      emit('on-checked', val)
      emit('change', val)
      emit('input', val)
    }
  })

  watch(() => props.multiple, (isMultiple) => {
    // todo 组件支持单选后操作
  })

</script>

<script>
  export default {
    name: 'cmdb-form-organization'
  }
</script>

<style lang="scss" scoped>
.cmdb-organization {
  position: relative;
  width: 100%;
  .selector {
    width: 100%;
    &.active {
      position: absolute;
      z-index: 2;
    }
  }
}
.organization-hide {
  display: none;
}
:deep(.bk-big-tree-empty) {
  position: static;
}
</style>
