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

<script setup>
  import { ref, nextTick, computed, reactive } from 'vue'
  import { useStore } from '@/store'
  import routerActions from '@/router/actions'
  import {
    MENU_MODEL_FIELD_TEMPLATE_CREATE_BASIC,
    MENU_MODEL_FIELD_TEMPLATE_BIND,
    MENU_MODEL_FIELD_TEMPLATE,
  } from '@/dictionary/menu-symbol'
  import fieldTemplateService from '@/service/field-template'
  import LeaveConfirm from '@/components/ui/dialog/leave-confirm'
  import TopSteps from './children/top-steps.vue'
  import FieldManage from './children/field-manage.vue'
  import CreateSuccess from './children/create-success.vue'
  import { wrapData, normalizeFieldData, normalizeUniqueData } from './children/use-field'
  import FieldPreview from './children/field-preview-drawer.vue'

  const store = useStore()

  const templateDraft = computed(() => store.getters['fieldTemplate/templateDraft'])
  const fieldData = computed(() => templateDraft.value.fieldList || [])
  const uniqueData = computed(() => templateDraft.value.uniqueList || [])
  const basicData = computed(() => templateDraft.value.basic)

  const requestIds = {
    submit: Symbol('submit')
  }

  const isCreateSuccess = ref(false)
  const newTemplateId = ref(null)
  const previewShow = ref(false)

  const leaveConfirmConfig = reactive({
    id: 'createFlowField',
    active: true
  })

  const previewFieldList = computed(() => templateData.value.fieldList)

  const templateData = computed(() => ({
    basic: basicData.value,
    fieldList: settingData.fieldList ?? fieldData.value,
    uniqueList: settingData.uniqueList ?? uniqueData.value
  }))

  const isDraftValid = computed(() => templateData.value.basic?.name?.length > 0)
  const submitButtonDisabled = computed(() => !isDraftValid.value || !templateData.value.fieldList?.length)

  const settingData = reactive({
    fieldList: null,
    uniqueList: null
  })

  const clearTemplateDraft = () => {
    store.commit('fieldTemplate/clearTemplateDraft')
  }

  const handleFieldUpdate = (data) => {
    settingData.fieldList = data.map(wrapData)
  }
  const handleUniqueUpdate = (data) => {
    settingData.uniqueList = data
  }

  const handlePrevStep = () => {
    leaveConfirmConfig.active = false

    store.commit('fieldTemplate/setTemplateDraft', {
      fieldList: templateData.value.fieldList,
      uniqueList: templateData.value.uniqueList
    })

    nextTick(() => {
      routerActions.redirect({
        name: MENU_MODEL_FIELD_TEMPLATE_CREATE_BASIC,
        history: true
      })
    })
  }
  const handleSubmit = async () => {
    const submitData = {
      ...templateData.value.basic,
      attributes: normalizeFieldData(templateData.value.fieldList),
      uniques: normalizeUniqueData(templateData.value.uniqueList, templateData.value.fieldList)
    }

    try {
      const result = await fieldTemplateService.create(submitData, { requestId: requestIds.submit })
      newTemplateId.value = result.id
      isCreateSuccess.value = true
      leaveConfirmConfig.active = false

      clearTemplateDraft()
    } catch (err) {
      isCreateSuccess.value = false
      console.error(err)
    }
  }
  const handleCancel = () => {
    leaveConfirmConfig.active = false
    nextTick(() => {
      routerActions.redirect({
        name: MENU_MODEL_FIELD_TEMPLATE
      })
    })
  }
  const handlePreview = () => {
    previewShow.value = true
  }

  const handleSuccessAction = (action) => {
    leaveConfirmConfig.active = false

    nextTick(() => {
      if (action === 'bind') {
        routerActions.redirect({
          name: MENU_MODEL_FIELD_TEMPLATE_BIND,
          params: {
            id: newTemplateId.value
          }
        })
      } else if (action === 'back') {
        routerActions.redirect({
          name: MENU_MODEL_FIELD_TEMPLATE
        })
      }
    })
  }

  const handleLeave = () => {
    clearTemplateDraft()
  }

  defineExpose({
    leaveConfirmConfig,
    clearTemplateDraft
  })
</script>
<script>
  export default {
    beforeRouteLeave(to, from, next) {
      if (![MENU_MODEL_FIELD_TEMPLATE_CREATE_BASIC].includes(to.name)) {
        if (!this.leaveConfirmConfig.active) {
          this.clearTemplateDraft()
        }
      }
      next()
    }
  }
</script>

<template>
  <div class="create-field-settings">
    <cmdb-sticky-layout class="cmdb-config-sticky-layout" v-if="!isCreateSuccess">
      <template #header="{ sticky }">
        <top-steps width="360px" :current="2" :class="{ 'is-sticky': sticky }"></top-steps>
      </template>
      <field-manage
        :field-list="fieldData"
        :unique-list="uniqueData"
        :is-create-mode="true"
        @update-field="handleFieldUpdate"
        @update-unique="handleUniqueUpdate">
      </field-manage>
      <template #footer="{ sticky }">
        <div :class="['layout-footer', { 'is-sticky': sticky }]">
          <cmdb-auth :auth="{ type: $OPERATION.C_FIELD_TEMPLATE }">
            <template #default="{ disabled }">
              <bk-button
                :disabled="disabled"
                @click="handlePrevStep">
                {{$t('上一步')}}
              </bk-button>
            </template>
          </cmdb-auth>
          <cmdb-auth :auth="{ type: $OPERATION.C_FIELD_TEMPLATE }" v-bk-tooltips="{
            disabled: isDraftValid,
            content: $t('未填写基础信息')
          }">
            <template #default="{ disabled }">
              <bk-button
                theme="primary"
                :disabled="submitButtonDisabled || disabled"
                :loading="$loading(requestIds.submit)"
                @click="handleSubmit">
                {{$t('提交')}}
              </bk-button>
            </template>
          </cmdb-auth>
          <bk-button
            :disabled="submitButtonDisabled"
            @click="handlePreview">
            {{$t('预览')}}
          </bk-button>
          <bk-button theme="default" @click="handleCancel">{{$t('取消')}}</bk-button>
        </div>
      </template>
    </cmdb-sticky-layout>
    <div class="create-success-container" v-else>
      <create-success @action="handleSuccessAction"></create-success>
    </div>
    <field-preview
      :preview-show.sync="previewShow"
      :properties="previewFieldList">
    </field-preview>
    <leave-confirm
      v-bind="leaveConfirmConfig"
      :reverse="true"
      :title="$t('是否退出')"
      :content="$t('新建步骤未完成，退出将撤销当前操作')"
      :ok-text="$t('退出')"
      :cancel-text="$t('取消')"
      @leave="handleLeave">
    </leave-confirm>
  </div>
</template>

<style lang="scss" scoped>
  .layout-footer {
    padding: 0;
    margin-left: 108px;

    &.is-sticky {
      margin-left: 0;
      justify-content: center;
    }
  }

  .create-success-container {
    padding: 20px 24px;
  }
</style>
