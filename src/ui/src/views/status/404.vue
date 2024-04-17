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
  <bk-exception class="exception-wrap-item" type="404">
    <div class="title">{{$t('页面不存在')}}</div>
    <div class="text-subtitle">{{$t('您访问的页面不存在xxx')}}</div>
    <div class="text-wrap">
      <bk-button theme="primary" @click="handleRedirect">{{$t('返回首页')}}</bk-button>
    </div>
  </bk-exception>
</template>

<script>
  import { MENU_INDEX } from '@/dictionary/menu-symbol'
  export default {
    methods: {
      handleRedirect() {
        this.$routerActions.redirect({
          name: MENU_INDEX
        })
      }
    }
  }
</script>

<style lang="scss" scoped>
.exception-wrap-item {
  justify-content: center;

  .title {
    margin: 8px 0;
  }
  .text-subtitle {
    color: #979BA5;
    font-size: 14px;
    text-align: center;
    margin-top: 14px;
  }
  .text-wrap {
    display: flex;
    align-items: center;
    justify-content: center;
    color: #3A84FF;
    font-size: 14px;
    margin-top: 24px;
  }
}
</style>
