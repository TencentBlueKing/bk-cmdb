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
  import { useRoute } from '@/router/index'
  import BindModel from './children/bind-model.vue'

  const route = useRoute()

  const templateId = computed(() => Number(route.params.id))

  const handleCancel = () => {}
  const handleSubmit = () => {}
</script>

<template>
  <div class="bind">
    <bind-model :height="`${$APP.height - 111 - 52}px`"></bind-model>
    <div class="bind-footer">
      <cmdb-auth :auth="{ type: $OPERATION.U_FIELD_TEMPLATE, relation: [templateId] }">
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
  </div>
</template>

<style lang="scss" scoped>
  .bind-footer {
    display: flex;
    align-items: center;
    justify-content: center;
    height: 52px;
    padding: 0 20px;
    background-color: #fff;
    border-top: 1px solid $borderColor;

    .bk-button {
      min-width: 86px;

      & + .bk-button {
        margin-left: 8px;
      }
      & + .auth-box {
          margin-left: 8px;
      }
    }
    .auth-box {
      & + .bk-button,
      & + .auth-box {
          margin-left: 8px;
      }
    }
  }
</style>
