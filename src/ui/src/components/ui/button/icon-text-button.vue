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
  const props = defineProps({
    icon: {
      type: String,
      default: 'plus'
    },
    text: String,
    disabled: Boolean,
    disabledTips: String
  })
  const emit = defineEmits(['click'])

  const handleClick = () => {
    if (props.disabled) {
      return false
    }
    emit('click')
  }
</script>

<template>
  <div :class="['icon-text-button', { 'disabled': props.disabled }]"
    @click="handleClick"
    v-bk-tooltips="{
      disabled: !props.disabled || !props.disabledTips,
      content: props.disabledTips,
      allowHtml: true
    }">
    <bk-icon :type="props.icon" />{{ props.text }}
  </div>
</template>

<style lang="scss" scoped>
  .icon-text-button {
    cursor: pointer;
    color: $primaryColor;
    display: inline-flex;
    align-items: center;
    .icon-plus {
      font-size: 20px !important;
    }

    &.disabled {
      color: $textDisabledColor;
    }
  }
</style>
