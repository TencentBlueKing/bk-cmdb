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
  defineProps({
    text: String,
    theme: String,
    hoverBgColor: {
      type: String,
      default: ''
    }
  })
</script>

<template>
  <span :class="['tag', theme, { nohover: !hoverBgColor.length }]">
    <em class="tag-text">
      <slot>{{text}}</slot>
    </em>
  </span>
</template>

<style lang="scss" scoped>
.tag {
  background: #E4FAF0;
  color: #14A568;
  border-radius: 2px;
  padding: 1px 4px;
  height: 16px;
  line-height: 16px;
  white-space: nowrap;
  display: flex;
  align-items: center;

  .tag-text {
    display: block;
    font-size: 12px;
    font-style: normal;
    transform: scale(.875);
  }

  &.new {
    color: #14A568;
    background: #E4FAF0;
  }
  &.changed {
    color: #FF9C01;
    background: #FFF3E1;
  }
  &.removed {
    color: #EA3636;
    background: #FCE9E8;
  }
  &.paused {
    color: #63656E;
  }

  &:not(.nohover):hover {
    background: v-bind(hoverBgColor) !important;
  }
}
</style>
