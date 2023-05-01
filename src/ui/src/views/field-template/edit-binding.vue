<!--
 * Tencent
  components: { topSteps },is pleased to support the open source community by making 蓝鲸 available.
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
  import { ref } from 'vue'
  import routerActions from '@/router/actions'
  import {
    MENU_MODEL_FIELD_TEMPLATE_EDIT_FIELD_SETTINGS
  } from '@/dictionary/menu-symbol'
  import TopSteps from './children/top-steps.vue'
  import BindModel from './children/bind-model.vue'

  const nextButtonDisabled = ref(false)
  const steps = [
    { title: '基础信息', icon: 1 },
    { title: '字段设置', icon: 2 },
    { title: '模型信息确认', icon: 3 }
  ]

  // 来源接口
  const bindModelData = ref([])

  const handlePrevStep = () => {
    routerActions.redirect({
      name: MENU_MODEL_FIELD_TEMPLATE_EDIT_FIELD_SETTINGS,
      history: true
    })
  }
  const handleCancel = () => {}
  const handleSubmit = () => {}
</script>

<template>
  <cmdb-sticky-layout class="cmdb-config-sticky-layout">
    <template #header="{ sticky }">
      <top-steps :steps="steps" width="45%" :current="3" :class="{ 'is-sticky': sticky }"></top-steps>
    </template>
    <bind-model :model-list="bindModelData"></bind-model>
    <template #footer="{ sticky }">
      <div :class="['layout-footer', { 'is-sticky': sticky }]">
        <cmdb-auth :auth="{ type: $OPERATION.C_FIELD_TEMPLATE }">
          <template #default="{ disabled }">
            <bk-button
              :disabled="nextButtonDisabled || disabled"
              @click="handlePrevStep">
              {{$t('上一步')}}
            </bk-button>
          </template>
        </cmdb-auth>
        <cmdb-auth :auth="{ type: $OPERATION.C_FIELD_TEMPLATE }">
          <template #default="{ disabled }">
            <bk-button
              theme="primary"
              :disabled="disabled"
              @click="handleSubmit">
              {{$t('提交')}}
            </bk-button>
          </template>
        </cmdb-auth>
        <bk-button theme="default" @click="handleCancel">{{$t('取消')}}</bk-button>
      </div>
    </template>
  </cmdb-sticky-layout>
</template>
