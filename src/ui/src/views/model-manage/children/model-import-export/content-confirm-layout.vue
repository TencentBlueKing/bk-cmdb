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
  <div class="content-confirm-layout">
    <div class="content-confirm-aside">
      <slot name="aside"></slot>
    </div>
    <div class="content-confirm-main">
      <slot name="main"></slot>
    </div>
  </div>
</template>

<script>
  import { defineComponent } from 'vue'

  export default defineComponent({
    name: 'ContentConfirmLayout'
  })
</script>

<style lang="scss" scoped>
$borderColor: #dcdee5;
.content-confirm {
  &-layout {
    display: flex;
    height: 100%;
  }

  &-aside {
    position: relative;
    flex: 0 0 360px;
    width: 360px;
    overflow-y: auto;
    overflow-x: visible;
    background: #fafbfd;

    &::after {
      content:"";
      display: block;
      position: absolute;
      right: 0;
      top: 0;
      bottom: 0;
      width: 1px;
      background-color: $borderColor;
    }
  }

  &-main {
    flex: 1 1 50%;
    margin-left: 25px;
    height: 100%;
    overflow-y: auto;
  }
}
</style>
