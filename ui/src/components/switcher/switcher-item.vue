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
  <span class="switcher-item"
    v-bk-tooltips="{
      boundary: 'window',
      placement: 'top',
      interactive: false,
      disabled: !tips,
      content: tips
    }"
    :class="{
      'is-active': group && group.active === name
    }"
    @click="handleClick">
    <slot></slot>
  </span>
</template>

<script>
  export default {
    name: 'cmdb-switcher-item',
    props: {
      name: {
        type: [String, Number],
        default: ''
      },
      tips: {
        type: String,
        default: ''
      }
    },
    computed: {
      group() {
        let group = this.$parent
        while (group && !group.isSwitcherGroup) {
          group = group.$parent
        }
        return group
      }
    },
    methods: {
      handleClick() {
        if (this.group) {
          this.group.setActive(this.name)
        }
      }
    }
  }
</script>

<style lang="scss" scoped>
    .switcher-item {
        display: inline-flex;
        justify-content: center;
        align-items: center;
        height: 32px;
        padding: 0 10px;
        border: 1px solid #c4c6cc;
        font-size: 14px;
        color: #c4c6cc;
        cursor: pointer;
        z-index: 1;
        & ~ .switcher-item {
            margin-left: -1px;
        }
        &:first-child {
            border-radius: 2px 0px 0px 2px;
        }
        &:last-child {
            border-radius: 0px 2px 2px 0px;
        }
        &.is-active {
            color: $primaryColor;
            border-color: $primaryColor;
            z-index: 2;
        }
    }
</style>
