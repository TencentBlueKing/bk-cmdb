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
  import { reactive, ref, watch } from 'vue'
  import { $bkPopover } from '@/magicbox/index.js'
  import DropdownOptionButton from './dropdown-option-button.vue'

  const props = defineProps({
    commands: {
        type: Array,
        default: () => ([])
      },
      show: Boolean,
      target: [Element, null]
  })

  const emit = defineEmits(['hide'])

  const content = ref(null)

  const operateModelPopover = reactive({
    instance: null,
  })

  const showModelOperation = () => {
    const { target } = props
    operateModelPopover.instance?.destroy?.()
    operateModelPopover.instance = $bkPopover(target, {
      content: content.value,
      delay: 0,
      hideOnClick: true,
      placement: 'bottom',
      animateFill: false,
      sticky: true,
      theme: 'light bound-model-popover',
      boundary: 'window',
      trigger: 'manual', // 'manual mouseenter',
      arrow: true,
      onHidden: () => {
        emit('hide')
      }
    })

    operateModelPopover.instance.show()
  }

  watch(()=> props.show, (val)=> {
    if(val) showModelOperation()
  })
</script>

<template>
  <ul class="bk-dropdown-list" ref="content">
    <li v-for="(cmd, cmdIndex) in props.commands" :key="cmdIndex">
      <!-- 带权限的按钮 -->
      <cmdb-auth
        v-if="cmd.auth"
        :auth="cmd.auth"
        @update-auth="cmd.onUpdateAuth"
      >
        <template #default="{ disabled: authDisabled }">
          <dropdown-option-button
            v-bk-tooltips="{
              content: cmd.disabledTooltips,
              disabled: !cmd.disabled && cmd.disabledTooltips,
              placement: 'right'
            }"
            :disabled="authDisabled || cmd.disabled"
            @click="cmd.handler"
          >
            {{ cmd.text }}
          </dropdown-option-button>
        </template>
      </cmdb-auth>
      <!-- 无权限按钮 -->
      <dropdown-option-button
        v-else
        :disabled="cmd.disabled"
        @click="cmd.handler"
      >
        {{ cmd.text }}
      </dropdown-option-button>
    </li>
  </ul>
</template>
