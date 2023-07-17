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
  import { ref, computed, reactive, nextTick, watchEffect } from 'vue'
  import { t } from '@/i18n'
  import { useRoute } from '@/router/index'
  import { useStore } from '@/store'
  import routerActions from '@/router/actions'
  import {
    MENU_MODEL_FIELD_TEMPLATE,
    MENU_MODEL_FIELD_TEMPLATE_EDIT_FIELD_SETTINGS
  } from '@/dictionary/menu-symbol'
  import { $success } from '@/magicbox/index.js'
  import LeaveConfirm from '@/components/ui/dialog/leave-confirm'
  import TopSteps from './children/top-steps.vue'
  import BasicForm from './children/basic-form.vue'
  import fieldTemplateService from '@/service/field-template'

  const route = useRoute()
  const store = useStore()

  const requestIds = {
    save: Symbol('save')
  }
  const steps = [
    { title: t('基础信息'), icon: 1 },
    { title: t('字段设置'), icon: 2 },
    { title: t('模型信息确认'), icon: 3 }
  ]

  const leaveConfirmConfig = reactive({
    id: 'editFlowBasic',
    active: false
  })

  const templateId = computed(() => Number(route.params.id))
  const basicFormRef = ref(null)

  const templateDraft = computed(() => store.getters['fieldTemplate/templateDraft'])

  const basicData = ref({
    name: '',
    description: ''
  })

  watchEffect(async () => {
    const template = await fieldTemplateService.findById(templateId.value)
    basicData.value.name = templateDraft.value.basic.name ?? template.name
    basicData.value.description = templateDraft.value.basic.description ?? template.description
  })

  const clearTemplateDraft = () => {
    store.commit('fieldTemplate/clearTemplateDraft')
  }

  const handleFormDataChange = () => {
    leaveConfirmConfig.active = true
  }

  const handleSave = async (type = 'next') => {
    const { formData } = basicFormRef.value
    const { name, description } = basicData.value // 一开始的名称和描述
    const { name: newName, description: newDescription }  = formData  // 新的名称和描述
    if (!await basicFormRef.value.validateAll()) {
      return
    }
    const saveData = {
      id: templateId.value,
      ...formData
    }

    try {
      if (name !== newName || description !== newDescription) {
        await fieldTemplateService.updateBaseInfo(saveData, { requestId: requestIds.save })
        $success(t('保存成功'))
      }
      if (type === 'next') {
        return true
      }
      handleCancel()
    } catch (err) {
      console.error(err)
    }
    return false
  }
  const handleNextStep = async () => {
    if (!await handleSave()) return
    leaveConfirmConfig.active = false

    const { formData } = basicFormRef.value
    store.commit('fieldTemplate/setTemplateDraft', { basic: formData })

    nextTick(() => {
      routerActions.redirect({
        name: MENU_MODEL_FIELD_TEMPLATE_EDIT_FIELD_SETTINGS,
        history: false
      })
    })
  }
  const handleCancel = () => {
    leaveConfirmConfig.active = false
    nextTick(() => {
      routerActions.redirect({
        name: MENU_MODEL_FIELD_TEMPLATE
      })
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
      if (![MENU_MODEL_FIELD_TEMPLATE_EDIT_FIELD_SETTINGS].includes(to.name)) {
        if (!this.leaveConfirmConfig.active) {
          this.clearTemplateDraft()
        }
      }
      next()
    }
  }
</script>

<template>
  <cmdb-sticky-layout class="cmdb-config-sticky-layout">
    <template #header="{ sticky }">
      <top-steps :steps="steps" width="45%" :class="{ 'is-sticky': sticky }"></top-steps>
    </template>
    <basic-form :data="basicData" ref="basicFormRef" @change="handleFormDataChange"></basic-form>
    <template #footer="{ sticky }">
      <div :class="['layout-footer', { 'is-sticky': sticky }]">
        <cmdb-auth :auth="{ type: $OPERATION.U_FIELD_TEMPLATE, relation: [templateId] }">
          <template #default="{ disabled }">
            <bk-button
              theme="primary"
              :disabled="disabled"
              :loading="$loading(requestIds.save)"
              @click="handleSave">
              {{$t('保存&退出')}}
            </bk-button>
          </template>
        </cmdb-auth>
        <cmdb-auth :auth="{ type: $OPERATION.U_FIELD_TEMPLATE, relation: [templateId] }">
          <template #default="{ disabled }">
            <bk-button
              :disabled="disabled"
              :loading="$loading(requestIds.save)"
              @click="handleNextStep">
              {{$t('保存&下一步')}}
            </bk-button>
          </template>
        </cmdb-auth>
        <bk-button theme="default" @click="handleCancel">{{$t('取消')}}</bk-button>
      </div>
    </template>
    <leave-confirm
      v-bind="leaveConfirmConfig"
      :reverse="true"
      :title="$t('是否退出')"
      :content="$t('编辑步骤未完成，退出将撤销当前操作')"
      :ok-text="$t('退出')"
      :cancel-text="$t('取消')"
      @leave="handleLeave">
    </leave-confirm>
  </cmdb-sticky-layout>
</template>

<style lang="scss" scoped>
  .layout-footer {
    width: 628px;
    margin: 0 auto;
    padding: 0 0 0 140px;
    position: relative;
    left: -36px;

    &.is-sticky {
      width: 100%;
      left: 0;
      padding: 0;
      justify-content: center;
    }
  }
</style>
