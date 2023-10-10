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
    template: {
      type: Object,
      default: () => ({})
    }
  })

  const emit = defineEmits(['action'])

  const handleAction = (actionName) => {
    emit('action', actionName)
  }
</script>

<template>
  <div class="create-success">
    <i class="bk-icon icon-success icon-check-circle-shape"></i>
    <div class="title">{{ $t('字段模板创建成功') }}</div>
    <div class="summary">{{ $t('接下来，您可以绑定需要应用的模型') }}</div>
    <div class="actions">
      <bk-button theme="primary" @click="handleAction('bind')">{{ $t('立即绑定') }}</bk-button>
      <bk-button @click="handleAction('back')">{{ $t('返回列表') }}</bk-button>
    </div>
  </div>
</template>

<style lang="scss" scoped>
  .create-success {
    display: flex;
    align-items: center;
    justify-content: center;
    flex-direction: column;
    gap: 16px;
    min-height: 420px;
    background: #fff;

    .icon-success {
      color: #2DCB56;
      font-size: 64px;
      margin-bottom: 18px;
    }

    .title {
      color: #313238;
      font-size: 24px;
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
