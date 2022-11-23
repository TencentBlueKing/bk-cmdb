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
  import { computed, defineComponent, ref } from 'vue'
  import router from '@/router/index.js'
  import routerActions from '@/router/actions'
  import ManagementForm from './children/management-form.vue'
  import store from '@/store'
  import setTemplateService from '@/service/set-template'
  import { getTemplateSyncStatus } from './children/use-template-data.js'
  import {
    MENU_BUSINESS_HOST_AND_SERVICE,
    MENU_BUSINESS_SET_TEMPLATE_DETAILS,
  } from '@/dictionary/menu-symbol'

  export default defineComponent({
    components: {
      ManagementForm
    },
    setup() {
      const managementForm = ref(null)

      const bizId = computed(() => store.getters['objectBiz/bizId'])

      const templateId = computed(() => parseInt(router.app.$route.params.templateId, 10))

      const loading = ref(true)
      const submitDisabled = ref(false)
      const successDialogShow = ref(false)
      const needSync = ref(false)

      const requestIds = {
        update: Symbol('update')
      }

      const updateTemplate = async (formData) => {
        const data = {
          id: templateId.value,
          bk_biz_id: bizId.value,
          ...formData
        }

        await setTemplateService.update(data, { requestId: requestIds.update })

        // 更新模板后，获取模板的同步状态
        const syncStatus = await getTemplateSyncStatus(bizId.value, templateId.value)
        needSync.value = syncStatus

        successDialogShow.value = true
      }

      const handleSubmit = async () => {
        const valid = await managementForm.value.validateAll()
        if (!valid) {
          return
        }

        const formData = managementForm.value.getData()
        updateTemplate(formData)
      }

      const handleCancel = () => {
        routerActions.back()
      }

      const handleDataLoaded = () => {
        loading.value = false
      }


      const handleToCreateInstance = () => {
        routerActions.redirect({ name: MENU_BUSINESS_HOST_AND_SERVICE })
      }
      const handleToSyncInstance = () => {
        successDialogShow.value = false
        routerActions.redirect({
          name: MENU_BUSINESS_SET_TEMPLATE_DETAILS,
          params: {
            templateId: templateId.value
          },
          query: {
            tab: 'instance'
          }
        })
      }
      const handleCloseSuccessDialog = () => {
        successDialogShow.value = false
        routerActions.redirect({
          name: MENU_BUSINESS_SET_TEMPLATE_DETAILS,
          params: {
            templateId: templateId.value
          },
          history: false
        })
      }

      return {
        bizId,
        templateId,
        loading,
        managementForm,
        submitDisabled,
        requestIds,
        successDialogShow,
        needSync,
        handleSubmit,
        handleDataLoaded,
        handleCancel,
        handleCloseSuccessDialog,
        handleToCreateInstance,
        handleToSyncInstance
      }
    }
  })
</script>

<template>
  <cmdb-sticky-layout class="edit-layout" v-bkloading="{ isLoading: loading }">
    <div class="layout-main">
      <management-form
        ref="managementForm"
        :data-id="templateId"
        :submit-disabled.sync="submitDisabled"
        @data-loaded="handleDataLoaded">
      </management-form>
    </div>
    <template #footer="{ sticky }">
      <div :class="['layout-footer', { 'is-sticky': sticky }]">
        <cmdb-auth :auth="{ type: $OPERATION.U_SET_TEMPLATE, relation: [bizId, templateId] }">
          <bk-button
            theme="primary"
            slot-scope="{ disabled }"
            :disabled="disabled || submitDisabled"
            :loading="$loading(requestIds.update)"
            @click="handleSubmit">
            {{$t('保存')}}
          </bk-button>
        </cmdb-auth>
        <bk-button theme="default" @click="handleCancel">{{$t('取消')}}</bk-button>
      </div>
    </template>

    <bk-dialog v-model="successDialogShow"
      width="480"
      :esc-close="false"
      :mask-close="true"
      :show-footer="false"
      :close-icon="false">
      <div class="update-alert-layout">
        <i class="bk-icon icon-check-1"></i>
        <h3>{{$t('修改成功')}}</h3>
        <p class="update-success-tips">{{$t(needSync ? '集群模板修改成功并需要同步提示' : '集群模板修改成功提示')}}</p>
        <div class="btns">
          <bk-button class="mr10" theme="primary"
            v-if="needSync"
            @click="handleToSyncInstance">
            {{$t('同步集群')}}
          </bk-button>
          <bk-button class="mr10"
            :theme="needSync ? 'default' : 'primary'"
            @click="handleToCreateInstance">
            {{$t('创建集群')}}
          </bk-button>
          <bk-button theme="default" @click="handleCloseSuccessDialog">{{$t('关闭按钮')}}</bk-button>
        </div>
      </div>
    </bk-dialog>
  </cmdb-sticky-layout>
</template>

<style lang="scss" scoped>
.edit-layout {
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
    font-size: 50px;
    color: #fff;
    border-radius: 50%;
    background-color: #2dcb56;
    margin: 8px 0 15px;
  }
  h3 {
    font-size: 24px;
    color: #313238;
    font-weight: normal;
    padding-bottom: 16px;
  }
  .btns {
    font-size: 0;
    padding-bottom: 20px;
    .bk-button {
      min-width: 86px;
    }
  }
  .update-success-tips {
    padding-bottom: 24px;
  }
}
</style>
