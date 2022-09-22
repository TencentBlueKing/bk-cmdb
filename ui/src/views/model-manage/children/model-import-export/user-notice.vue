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
  <div class="user-notice">
    <div class="notice-icon">
      <img alt="Notice" src="@/assets/images/notice-pane.png" />
    </div>
    <h2 class="notice-title">
      {{ t("用户须知") }}
    </h2>
    <div class="notice-content">
      <slot></slot>
    </div>
  </div>
</template>

<script>
  import { defineComponent } from '@vue/composition-api'
  import { t } from '@/i18n'

  export default defineComponent({
    name: 'UserNotice',
    setup() {
      return {
        t
      }
    }
  })
</script>

<style lang="scss" scoped>
.notice {
  &-icon {
    margin: 30px auto 0;
    width: 158px;
    height: 106px;

    img {
      width: 100%;
      height: 100%;
    }
  }

  &-title {
    text-align: center;
    margin-top: 20px;
    font-weight: 400;
    font-size: 24px;
  }

  &-content {
    width: 680px;
    margin: 20px auto 0;
    padding: 16px 16px 26px;
    font-size: 12px;
    line-height: 20px;
    background: #fafbfd;
    border: 1px solid #dcdee5;
    border-radius: 2px;
  }
}
</style>
