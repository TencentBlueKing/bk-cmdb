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
  <div class="embed-permission">
    <div class="content">
      <bk-exception type="403">
        <div>
          <i18n path="抱歉您没有查看权限">
            <template #link>
              <bk-button class="text-btn"
                text
                theme="primary"
                @click="handleApplyPermission">
                {{$t('去申请')}}
              </bk-button>
            </template>
          </i18n>
        </div>
      </bk-exception>
    </div>
  </div>
</template>
<script>
  import permissionMixins from '@/mixins/permission'
  export default {
    name: 'cmdb-permission',
    mixins: [permissionMixins],
    props: {
      permission: {
        type: Array,
        default: () => []
      }
    }
  }
</script>

<style lang="scss" scoped>
    .embed-permission {
        color: #63656e;
        font-size: 14px;
        text-align: center;
        .text-btn {
            font-size: 14px;
            height: auto;
        }
    }
</style>
