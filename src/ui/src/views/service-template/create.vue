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
  import router from '@/router/index.js'
  import { t } from '@/i18n'
  import { $bkInfo } from '@/magicbox/index.js'
  import routerActions from '@/router/actions'
  import ManagementForm from './children/management-form.vue'
  import store from '@/store'
  import serviceTemplateService from '@/services/service-template'
  import {
    MENU_BUSINESS_SERVICE_TEMPLATE,
    MENU_BUSINESS_SERVICE_TEMPLATE_CREATE,
    MENU_BUSINESS_SERVICE_TEMPLATE_DETAILS,
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
      const cloneTemplateId = computed(() => parseInt(router.app.$route.query.clone, 10))
      const isClone = computed(() => cloneTemplateId.value > 0)

      const loading = ref(true)
      const submitDisabled = ref(false)
      const successDialogShow = ref(false)
      let newId = null

      const requestIds = {
        create: Symbol('create')
      }

      const createTemplate = async (formData) => {
        const data = { bk_biz_id: bizId.value, ...formData }
        const { id } = await serviceTemplateService.create(data, { requestId: requestIds.create })
        newId = id
        successDialogShow.value = true
      }

      const handleSubmit = async () => {
        const valid = await managementForm.value.validateAll()
        if (!valid) {
          return
        }

        const formData = managementForm.value.getData()

        if (!formData.processes.length) {
          $bkInfo({
            title: t('确认提交'),
            subTitle: t('服务模板创建没进程提示'),
            extCls: 'bk-dialog-sub-header-center',
            confirmFn: async () => {
              createTemplate(formData)
            },
            confirmLoading: true
          })
          return
        }

        createTemplate(formData)
      }

      const handleToSetTemplate = () => {
        routerActions.redirect({ name: MENU_BUSINESS_SET_TEMPLATE })
      }
      const handleToBusinessTopo = () => {
        routerActions.redirect({ name: MENU_BUSINESS_HOST_AND_SERVICE })
      }
      const handleContinueCreating = () => {
        successDialogShow.value = false
        routerActions.redirect({
          name: MENU_BUSINESS_SERVICE_TEMPLATE_CREATE,
          reload: true
        })
      }
      const handleSuccessDialogClose = () => {
        successDialogShow.value = false
        routerActions.redirect({
          name: MENU_BUSINESS_SERVICE_TEMPLATE_DETAILS,
          params: {
            templateId: newId
          },
          reload: true
        })
      }

      const handleCancel = () => {
        routerActions.redirect({ name: MENU_BUSINESS_SERVICE_TEMPLATE })
      }

      const handleDataLoaded = () => {
        loading.value = false
      }

      return {
        bizId,
        cloneTemplateId,
        isClone,
        loading,
        managementForm,
        submitDisabled,
        requestIds,
        successDialogShow,
        handleSubmit,
        handleCancel,
        handleToSetTemplate,
        handleToBusinessTopo,
        handleContinueCreating,
        handleSuccessDialogClose,
        handleDataLoaded
      }
    }
  })
</script>

<template>
  <cmdb-sticky-layout class="create-layout" v-bkloading="{ isLoading: loading }">
    <div class="layout-main">
      <management-form
        ref="managementForm"
        :data-id="cloneTemplateId"
        :is-clone="isClone"
        :submit-disabled.sync="submitDisabled"
        @data-loaded="handleDataLoaded">
      </management-form>
    </div>
    <template #footer="{ sticky }">
      <div :class="['layout-footer', { 'is-sticky': sticky }]">
        <cmdb-auth :auth="{ type: $OPERATION.C_SERVICE_TEMPLATE, relation: [bizId] }">
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

    <bk-dialog v-model="successDialogShow"
      width="520"
      :esc-close="false"
      :mask-close="false"
      :show-footer="false"
      :close-icon="false">
      <div class="update-alert-layout">
        <i class="bk-icon icon-check-1"></i>
        <h3 class="title">{{$t('创建成功')}}</h3>
        <i18n path="服务模板创建成功，您可以在XXX" tag="p" class="next-content">
          <template #link1>
            <bk-link theme="primary" @click="handleToSetTemplate">{{$t('集群模板')}}</bk-link>
          </template>
          <template #link2>
            <bk-link theme="primary" @click="handleToBusinessTopo">{{$t('业务拓扑')}}</bk-link>
          </template>
        </i18n>
        <div class="btns">
          <bk-button class="btn mr10" theme="primary" @click="handleContinueCreating">{{$t('继续创建')}}</bk-button>
          <bk-button class="btn" theme="default" @click="handleSuccessDialogClose">{{$t('关闭')}}</bk-button>
        </div>
      </div>
    </bk-dialog>
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

.update-alert-layout {
  text-align: center;
  .bk-icon {
    width: 58px;
    height: 58px;
    line-height: 58px;
    font-size: 30px;
    color: #fff;
    border-radius: 50%;
    background-color: #2dcb56;
    margin: 8px 0 15px;
  }
  .title {
    font-size: 24px;
    color: #313238;
    font-weight: normal;
    padding-bottom: 16px;
  }
  .btns {
    font-size: 0;
    padding-bottom: 20px;
    .btn {
      min-width: 86px;
    }
  }
  .next-content {
    padding-bottom: 24px;
    .bk-link {
      vertical-align: baseline;
    }
  }
}
</style>
