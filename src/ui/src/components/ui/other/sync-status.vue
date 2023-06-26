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
  import { t } from '@/i18n'

  const props = defineProps({
    status: String,
    tips: {
      type: String,
      default: ''
    },
    mini: Boolean
  })

  const statusText = computed(() => {
    const texts = {
      syncing: t('同步中'),
      need_sync: t('待同步'),
      finished: t('已同步'),
      failure: t('同步失败')
    }
    return texts[props.status] || '--'
  })
</script>

<template>
  <span :class="['sync-status', { 'has-tips': tips.length > 0 }]" v-bk-tooltips="{
    disabled: !tips,
    content: tips
  }">
    <template v-if="status === 'syncing'">
      <img class="svg-icon" src="@/assets/images/icon/loading.svg">
    </template>
    <template v-if="status === 'need_sync'">
      <i class="status-circle waiting"></i>
    </template>
    <template v-if="status === 'finished'">
      <i class="status-circle success"></i>
    </template>
    <template v-if="status === 'failure'">
      <i class="status-circle fail"></i>
    </template>
    <slot v-if="!mini">{{ statusText }}</slot>
  </span>
</template>

<style lang="scss" scoped>
.sync-status {
  display: inline-flex;
  align-items: center;
  font-size: 12px;
  color: #63656E;
  gap: 4px;
  .status-circle {
    display: inline-block;
    width: 8px;
    height: 8px;
    border-radius: 50%;
    &.waiting {
      background: #FFF;
      border: 1px solid #3A84FF;
    }
    &.success {
      background: #E5F6EA;
      border: 1px solid #3FC06D;
    }
    &.fail {
      background-color: #FFE6E6;
      border: 1px solid #EA3636;
    }
  }
  .svg-icon {
    @include inlineBlock;
    margin-top: -4px;
    width: 16px;
  }

  &.has-tips {
    cursor: pointer;
  }
}
</style>
