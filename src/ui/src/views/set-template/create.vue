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

<script lang="ts">
  import { computed, defineComponent, ref } from '@vue/composition-api'
  import { t } from '@/i18n'
  import { $bkInfo } from '@/magicbox/index.js'
  import routerActions from '@/router/actions'
  import ManagementForm from './children/management-form.vue'
  import store from '@/store'
  import setTemplateService from '@/services/set-template'
  import {
    MENU_BUSINESS_HOST_AND_SERVICE,
    MENU_BUSINESS_SET_TEMPLATE
  } from '@/dictionary/menu-symbol'

  export default defineComponent({
    components: {
      ManagementForm
    },
    setup() {
      const managementForm = ref(null)

      const bizId = computed(() => store.getters['objectBiz/bizId'])

      const loading = ref(true)
      const submitDisabled = ref(false)

      const requestIds = {
        create: Symbol('create')
      }

      const createTemplate = async (formData) => {
        const data = { bk_biz_id: bizId.value, ...formData }

        await setTemplateService.create(data, { requestId: requestIds.create })

        $bkInfo({
          type: 'success',
          width: 480,
          title: t('创建成功'),
          subTitle: t('创建集群模板成功提示'),
          okText: t('创建集群'),
          cancelText: t('返回列表'),
          confirmFn: () => {
            routerActions.redirect({ name: MENU_BUSINESS_HOST_AND_SERVICE })
          },
          cancelFn: () => {
            routerActions.redirect({ name: MENU_BUSINESS_SET_TEMPLATE })
          }
        })
      }

      const handleSubmit = async () => {
        const valid = await managementForm.value.validateAll()
        if (!valid) {
          return
        }

        const formData = managementForm.value.getData()
        createTemplate(formData)
      }

      const handleCancel = () => {
        routerActions.redirect({ name: MENU_BUSINESS_SET_TEMPLATE })
      }

      const handleDataLoaded = () => {
        loading.value = false
      }

      return {
        bizId,
        loading,
        managementForm,
        submitDisabled,
        requestIds,
        handleSubmit,
        handleDataLoaded,
        handleCancel
      }
    }
  })
</script>

<template>
  <cmdb-sticky-layout class="create-layout" v-bkloading="{ isLoading: loading }">
    <div class="layout-main">
      <management-form
        ref="managementForm"
        :submit-disabled.sync="submitDisabled"
        @data-loaded="handleDataLoaded">
      </management-form>
    </div>
    <template #footer="{ sticky }">
      <div :class="['layout-footer', { 'is-sticky': sticky }]">
        <cmdb-auth :auth="{ type: $OPERATION.C_SET_TEMPLATE, relation: [bizId] }">
          <bk-button
            theme="primary"
            slot-scope="{ disabled }"
            :disabled="disabled || submitDisabled"
            :loading="$loading(requestIds.create)"
            @click="handleSubmit">
            {{$t('提交')}}
          </bk-button>
        </cmdb-auth>
        <bk-button theme="default" @click="handleCancel">{{$t('取消')}}</bk-button>
      </div>
    </template>
  </cmdb-sticky-layout>
</template>

<style lang="scss" scoped>
.create-layout {
  .layout-main {
    padding: 15px 20px 0 20px;
  }

  .layout-footer {
    display: flex;
    align-items: center;
    height: 52px;
    padding: 0 20px;
    margin-top: 8px;
    .bk-button {
      min-width: 86px;

      & + .bk-button {
        margin-left: 8px;
      }
    }
    .auth-box + .bk-button {
      margin-left: 8px;
    }
    &.is-sticky {
      background-color: #fff;
      border-top: 1px solid $borderColor;
    }
  }
}
</style>
