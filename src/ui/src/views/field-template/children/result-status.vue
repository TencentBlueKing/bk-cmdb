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
  // eslint-disable-next-line no-unused-vars
  const props = defineProps({
    status: {
      type: String,
      default: ''
    }
  })
</script>

<template>
  <div class="result-status">
    <img class="status-icon icon-loading" src="@/assets/images/icon/loading.svg" v-if="status === 'loading'" />
    <bk-icon type="close-circle-shape" class="status-icon icon-error"
      v-else-if="status === 'apierror' || status === 'error'" />
    <bk-icon type="check-circle-shape" class="status-icon icon-success" v-else-if="status === 'success'" />
    <bk-icon type="exclamation-circle-shape" class="status-icon icon-abnormal" v-else-if="status === 'abnormal'" />

    <div class="title">
      <slot name="title"></slot>
    </div>
    <div class="summary">
      <slot name="summary"></slot>
    </div>
    <div class="actions">
      <slot name="actions"></slot>
    </div>
  </div>
</template>

<style lang="scss" scoped>
  .result-status {
    display: flex;
    align-items: center;
    justify-content: center;
    flex-direction: column;
    gap: 16px;
    min-height: 420px;
    background: #fff;

    .status-icon {
      margin-bottom: 18px;
      &:not(.icon-loading) {
        font-size: 64px !important;
      }
      &.icon-loading {
        width: 56px;
        height: 56px;
      }
      &.icon-error {
        color: #EA3636;
      }
      &.icon-success {
        color: #2DCB56;
      }
      &.icon-fail {
        color: #FF9C01;
      }
      &.icon-abnormal {
        color: #FF9C01;
      }
    }

    .title {
      color: #313238;
      font-size: 24px;
      .count {
        font-weight: 700;
        font-style: normal;
        padding: 0 .2em;
      }
    }

    .summary {
      font-size: 14px;
      color: #63656E;
      margin-bottom: 8px;
    }

    .actions {
      display: flex;
      gap: 8px;
    }
  }
</style>
