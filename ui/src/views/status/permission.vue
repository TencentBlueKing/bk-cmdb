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
  <div class="tips-wrapper">
    <div class="content-wrapper">
      <bk-exception type="403">
        <div class="title">
          <h2>{{$t('无功能权限')}}</h2>
          <p>{{$t('点击下方按钮申请')}}</p>
        </div>
        <div class="btns">
          <bk-button theme="primary" @click="handleApplyPermission" :loading="$loading('getSkipUrl')">
            {{$t('申请功能权限')}}
          </bk-button>
        </div>
      </bk-exception>
    </div>
  </div>
</template>
<script>
  import { translateAuth } from '@/setup/permission'
  export default {
    methods: {
      async handleApplyPermission() {
        try {
          const { view, permission } = this.$route.meta.auth || {}
          const skipUrl = await this.$store.dispatch('auth/getSkipUrl', {
            params: view ? translateAuth(view) : permission,
            config: {
              requestId: 'getSkipUrl'
            }
          })
          window.open(skipUrl)
        } catch (e) {
          console.error(e)
        }
      }
    }
  }
</script>

<style lang="scss" scoped>
    .tips-wrapper {
        overflow: hidden;
    }
    .content-wrapper {
        margin-top: 100px;
        text-align: center;
        color: #63656E;
        font-size: 14px;
        .title {
            img {
                width: 136px;
            }
            h2 {
                margin-bottom: 10px;
                font-size: 22px;
                color: #313238;
                font-weight: normal;
            }
            p {
                color: #63656e;
                font-size: 14px;
            }
        }
        .btns {
            margin-top: 24px;
            .bk-button {
                border-radius: 3px;
                padding-left: 10px;
                padding-right: 10px;
            }
        }
    }
</style>
