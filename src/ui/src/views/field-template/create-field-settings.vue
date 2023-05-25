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
  import { computed, reactive } from 'vue'
  import { useStore } from '@/store'
  import routerActions from '@/router/actions'
  import {
    MENU_MODEL_FIELD_TEMPLATE_CREATE_BASIC
  } from '@/dictionary/menu-symbol'
  import fieldTemplateService from '@/service/field-template'
  import TopSteps from './children/top-steps.vue'
  import FieldManage from './children/field-manage.vue'
  import { wrapData, normalizeFieldData, normalizeUniqueData } from './children/use-field'

  const store = useStore()

  const templateDraft = computed(() => store.getters['fieldTemplate/templateDraft'])
  const fieldData = computed(() => templateDraft.value.fieldList)
  const uniqueData = computed(() => templateDraft.value.uniqueList)
  const basicData = computed(() => templateDraft.value.basic)

  const requestIds = {
    submit: Symbol('submit')
  }

  const templateData = computed(() => ({
    basic: basicData.value,
    fieldList: settingData.fieldList ? settingData.fieldList : fieldData.value,
    uniqueList: settingData.uniqueList ? settingData.uniqueList : uniqueData.value
  }))

  const submitButtonDisabled = computed(() => !templateData.value.basic?.name?.length
    || !templateData.value.fieldList?.length)

  const settingData = reactive({
    fieldList: null,
    uniqueList: null
  })
  const handleFieldUpdate = (data) => {
    settingData.fieldList = data.map(wrapData)
  }
  const handleUniqueUpdate = (data) => {
    settingData.uniqueList = data
  }

  const handlePrevStep = () => {
    console.log(templateData.value, 'templateData.value')
    store.commit('fieldTemplate/setTemplateDraft', {
      fieldList: templateData.value.fieldList,
      uniqueList: templateData.value.uniqueList
    })
    routerActions.redirect({
      name: MENU_MODEL_FIELD_TEMPLATE_CREATE_BASIC,
      history: true
    })
  }
  const handleSubmit = async () => {
    const submitData = {
      ...templateData.value.basic,
      attributes: normalizeFieldData(templateData.value.fieldList),
      uniques: normalizeUniqueData(templateData.value.uniqueList, templateData.value.fieldList)
    }

    await fieldTemplateService.create(submitData, { requestId: requestIds.submit })
  }
  const handleCancel = () => {}
  const handlePreview = () => {}
</script>

<template>
  <cmdb-sticky-layout class="cmdb-config-sticky-layout">
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
        <cmdb-auth :auth="{ type: $OPERATION.C_FIELD_TEMPLATE }">
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
          @click="handlePreview">
          {{$t('预览')}}
        </bk-button>
        <bk-button theme="default" @click="handleCancel">{{$t('取消')}}</bk-button>
      </div>
    </template>
  </cmdb-sticky-layout>
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
</style>
