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
const props = defineProps({
  value: {
    type: Boolean,
    default: true,
  },
  tips: String,
})
const emit = defineEmits(['input', 'click'])

const hasTips = computed(() => props.tips && props.tips.length > 0)

const handleTogglState = () => {
  emit('input', !props.value)
  emit('click', !props.value)
}
</script>

<template>
  <span
    v-bk-tooltips="{
      disabled: !hasTips,
      placement: 'top',
      interactive: false,
      content: tips,
      delay: [100, 0],
    }"
    class="lock-button"
    tabindex="-1"
    @click="handleTogglState">
    <i v-if="value" class="icon-cc-lock-fill"></i>
    <i v-else class="icon-cc-unlock-fill"></i>
  </span>
</template>

<style lang="scss" scoped>
.lock-button {
  display: inline-flex;
  width: 24px;
  align-items: center;
  justify-content: center;
  font-size: 14px;
  overflow: hidden;
  cursor: pointer;
}
</style>
