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
  <div>
    <span v-if="value === 'enable'" class="project-status">
      <i class="status-circle enable"></i>
      {{$t('启用中')}}
    </span>
    <span v-else-if="value === 'disabled'" class="project-status">
      <i class="status-circle disabled"></i>
      {{$t('已停用')}}
    </span>
  </div>
</template>

<script>
  export default {
    name: 'InstanceStatusColumn',
    props: {
      value: {
        type: String,
        default: ''
      }
    }
  }
</script>

<style lang="scss" scoped>
.project-status {
    .status-circle {
        display: inline-block;
        width: 8px;
        height: 8px;
        margin-right: 4px;
        border-radius: 50%;
        &.disabled {
          background-color: rgb(240, 241, 245);
          border: 1px solid rgb(198,200,204);
        }
        &.enable {
          border: 1px solid rgb(112,192,167);
          background-color: rgb(232, 245, 235);
        }
    }
}
</style>
