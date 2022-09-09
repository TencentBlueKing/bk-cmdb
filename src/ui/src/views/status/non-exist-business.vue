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
          <h2>{{$t('业务不存在或无权限')}}</h2>
        </div>
        <div class="btns">
          <bk-button theme="primary" @click="handleApplyPermission" :loading="$loading('getSkipUrl')">
            {{$t('申请业务访问权限')}}
          </bk-button>
          <bk-button theme="primary" @click="handleCreate">
            {{$t('创建业务')}}
          </bk-button>
        </div>
      </bk-exception>
    </div>
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
