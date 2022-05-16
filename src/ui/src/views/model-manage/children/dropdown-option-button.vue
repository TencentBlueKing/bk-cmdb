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
  <div
    :class="['dropdown-option-button', { disabled: disabled }, { 'is-active': active }]"
    @click="handleClick"
  >
    <slot></slot>
  </div>
</template>

<script>
  export default {
    name: 'DropdownOptionButton',
    props: {
      // 禁用选项
      disabled: {
        type: Boolean,
        default: false
      },
      // 激活选项
      active: {
        type: Boolean,
        default: false
      },
    },
    methods: {
      handleClick() {
        if (!this.disabled) {
          this.$emit('click')
        }
      },
    },
  }
</script>

<style lang="scss" scoped>
.dropdown-option-button {
  height: 32px;
  line-height: 32px;
  white-space: nowrap;
  color: #63656E;
  padding: 0 12px;
  font-size: 12px;
  user-select: none;

  &:hover,
  &.is-active {
    background-color: #e1ecff;
    color: #3A84FF;
  }

  &.disabled {
    color: $textDisabledColor;
    cursor: not-allowed;
  }
}
</style>
