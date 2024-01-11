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
        <template v-if="isResourceInstanceView">
          <div class="title">
            <h2>{{$t('无模型相关权限')}}</h2>
            <p>{{$t('请前往权限中心申请相关权限')}}</p>
          </div>
          <div class="btns">
            <bk-button theme="primary" @click="handleApplyPermission" :loading="$loading('getSkipUrl')">
              {{$t('去申请')}}
            </bk-button>
          </div>
        </template>
        <template v-else>
          <div class="title">
            <h2>{{$t('无功能权限')}}</h2>
            <p>{{$t('点击下方按钮申请')}}</p>
          </div>
          <div class="btns">
            <bk-button theme="primary" @click="handleApplyPermission" :loading="$loading('getSkipUrl')">
              {{$t('申请功能权限')}}
            </bk-button>
          </div>
        </template>
      </bk-exception>
    </div>
  </div>
</template>
<script>
  import { translateAuth } from '@/setup/permission'
  import {
    MENU_RESOURCE_INSTANCE
  } from '@/dictionary/menu-symbol'

  export default {
    computed: {
      isResourceInstanceView() {
        return this.$route.name === MENU_RESOURCE_INSTANCE
      }
    },
    created() {
      if (this.isResourceInstanceView) {
        const modelId = this.$route.params.objId
        const model = this.$store.getters['objectModelClassify/getModelById'](modelId)
        this.$store.commit('setTitle', model?.bk_obj_name ?? '')
      }
    },
    methods: {
      async handleApplyPermission() {
        try {
          // authKey确定用哪个权限申请
          const { authKey } = this.$route.meta

          const { permission } = this.$route.meta.auth || {}
          const view = this.$route.meta.auth[authKey]

          const viewAuth = typeof view === 'function' ? view(this.$route, this) : view
          const viewAuths = [viewAuth]

          // 如果存在superView并且未鉴权通过，则需要一起申请
          if (this.$route.meta.auth.superView) {
            const { superView } = this.$route.meta.auth
            const superAuth = typeof superView === 'function' ? superView(this.$route, this) : superView
            const authSuperViewResult = await this.$store.dispatch('auth/getViewAuth', superAuth)
            if (!authSuperViewResult) {
              viewAuths.unshift(superAuth)
            }
          }

          const skipUrl = await this.$store.dispatch('auth/getSkipUrl', {
            params: view ? translateAuth(viewAuths) : permission,
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
                min-width: 100px;
            }
        }
    }
</style>
