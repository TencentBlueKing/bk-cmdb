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
  <component href="javascript:void(0)"
    :class="['link-button', { 'disabled': disabled }]"
    :is="tag"
    @click="handleClick">
    <slot></slot>
  </component>
</template>

<script>
  export default {
    name: 'link-button',
    props: {
      disabled: Boolean,
      tag: {
        type: String,
        default: 'a'
      },
      href: {
        type: String,
        default: ''
      }
    },
    methods: {
      handleClick(event) {
        if (this.disabled) {
          return false
        }
        this.$emit('click', event)
      }
    }
  }
</script>

<style lang="scss" scoped>
    .link-button {
        display: inline-block;
        vertical-align: middle;
        line-height: 32px;
        color: $primaryColor;
        &.disabled {
            color: $textDisabledColor;
            cursor: not-allowed;
        }
    }
</style>
