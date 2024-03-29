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
  <div class="non-exist-business">
    <!-- 找不到当前业务 -->
    <bk-exception class="exception" type="404" v-if="isBizNotFound">
      <h2 class="title">{{$t('业务不存在')}}</h2>
      <h3 class="subtitle">
        {{$t('你访问的业务不存在xxx')}}
      </h3>
      <div class="action">
        <bk-button theme="primary" @click="handleCreate">
          {{$t('创建业务')}}
        </bk-button>
      </div>
    </bk-exception>

    <!-- 找不到一个有权限的业务 -->
    <bk-exception class="exception" type="403" v-else-if="isNoBiz">
      <h2 class="title">{{$t('无任何业务权限')}}</h2>
      <h3 class="subtitle">
        {{$t('你没有任何业务的访问权限xxx')}}
      </h3>
      <div class="action">
        <bk-button theme="primary" @click="handleApplyPermission" :loading="$loading('getSkipUrl')">
          {{$t('去申请')}}
        </bk-button>
        <bk-button @click="handleCreate">
          {{$t('立即创建')}}
        </bk-button>
      </div>
    </bk-exception>

    <!-- 没有当前业务的权限 -->
    <bk-exception class="exception" type="403" v-else-if="isBizUnauthed">
      <h2 class="title">{{$t('无当前业务权限')}}</h2>
      <h3 class="subtitle">
        {{$t('你没有当前业务的访问权限xxx')}}
      </h3>
      <div class="action">
        <bk-button theme="primary" @click="handleApplyPermission" :loading="$loading('getSkipUrl')">
          {{$t('申请权限')}}
        </bk-button>
      </div>
    </bk-exception>
  </div>
</template>
<script>
  import { translateAuth } from '@/setup/permission'
  import { MENU_RESOURCE_BUSINESS } from '@/dictionary/menu-symbol'
  import businessService from '@/service/business/search.js'

  export default {
    computed: {
      bizId() {
        return this.$route.params.bizId
      },
      isNoBiz() {
        return this.$route.name === 'no-business'
      },
      isBizNotFound() {
        return this.$route.meta?.extra?.isNotFound
      },
      isBizUnauthed() {
        return this.$route.meta?.extra?.isUnauthed
      }
    },
    methods: {
      async handleApplyPermission() {
        const allBusinessList = await businessService.findAll()
        const availableBusiness = allBusinessList.some(business => business.bk_biz_id === this.bizId)
        try {
          const permission = translateAuth({
            type: this.$OPERATION.R_BIZ_RESOURCE,
            relation: availableBusiness ? [this.bizId] : []
          })
          const skipUrl = await this.$store.dispatch('auth/getSkipUrl', {
            params: permission,
            config: {
              requestId: 'getSkipUrl'
            }
          })
          window.open(skipUrl)
        } catch (e) {
          console.error(e)
        }
      },
      handleCreate() {
        this.$routerActions.redirect({ name: MENU_RESOURCE_BUSINESS })
      }
    }
  }
</script>

<style lang="scss" scoped>
.non-exist-business {
  display: flex;
  justify-content: center;

  .exception {
    justify-content: center;

    .title {
      font-weight: 400;
      color: #63656E;
      font-size: 24px;
      margin: 8px 0;
    }
    .subtitle {
      font-weight: 400;
      font-size: 14px;
      color: #979BA5;
      margin: 16px 0;
    }
    .action {
      margin-top: 24px;
      .bk-button {
        min-width: 88px;
      }
    }
  }
}
</style>
