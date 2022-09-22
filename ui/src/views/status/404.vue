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
    <p>{{$t('页面不存在')}}</p>
    <div class="text-wrap">
      <bk-link theme="primary" @click="handleRedirect">{{$t('返回首页')}}</bk-link>
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
    }
</style>
