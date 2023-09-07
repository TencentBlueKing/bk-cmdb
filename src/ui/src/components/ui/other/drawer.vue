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
    open: {
      type: Boolean,
      default: false
    },
    title: {
      type: String,
      default: ''
    },
    width: {
      type: String,
      default: '210px'
    },
    placement: {
      type: String,
      default: 'right'
    }
  })

  const emit = defineEmits(['close'])

  // watch(() => props.open, (open) => {

  // })

  const handleClose = () => emit('close')
</script>

<template>
  <div :class="[
    'drawer',
    `drawer-${placement}`,
    { 'drawer-open': open },
    { 'drawer-hidden': !open }
  ]">
    <div class="drawer-head">
      <slot name="title">
        <div class="drawer-title">{{ title }}</div>
      </slot>
      <bk-icon type="close" class="icon-close" @click="handleClose" />
    </div>
    <div class="drawer-body">
      <slot name="content"></slot>
    </div>
  </div>
</template>

<style lang="scss" scoped>
.drawer {
  width: v-bind(width);
  position: absolute;
  transition: all 0.3s;
  background: #fff;

  &.drawer-right {
    right: 0;
    top: 0;
  }

  &.drawer-hidden {
    display: none;
  }

  .drawer-head {
    display: flex;
    align-items: center;
    justify-content: space-between;

    .drawer-title {
      font-size: 16px;
      color: #313238;
    }

    .icon-close {
      color: #979BA5;
      &:hover {
        color: $primaryColor;
      }
    }
  }
  .drawer-body {
    padding: 16px 12px;
  }
}
</style>
