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
  import { ref, computed, reactive, nextTick, watchEffect, set } from 'vue'
  import { t } from '@/i18n'
  import { useRoute } from '@/router/index'
  import { useStore } from '@/store'
  import routerActions from '@/router/actions'
  import {
    MENU_MODEL_FIELD_TEMPLATE,
    MENU_MODEL_FIELD_TEMPLATE_EDIT_FIELD_SETTINGS
  } from '@/dictionary/menu-symbol'
  import LeaveConfirm from '@/components/ui/dialog/leave-confirm'
  import TopSteps from './children/top-steps.vue'
  import BindModel from './children/bind-model.vue'
  import SyncResults from './children/sync-results.vue'
  import useField, { normalizeFieldData, normalizeUniqueData, unwrapData } from './children/use-field'
  import useUnique from './children/use-unique'
  import fieldTemplateService from '@/service/field-template'

  const route = useRoute()
  const store = useStore()
  const steps = [
    { title: t('基础信息'), icon: 1 },
    { title: t('字段设置'), icon: 2 },
    { title: t('模型信息确认'), icon: 3 }
  ]
  const requestIds = {
    submit: Symbol('submit')
  }

  const templateId = computed(() => Number(route.params.id))

  const isEditSuccess = ref(false)

  const leaveConfirmConfig = reactive({
    id: 'editFlowBind',
    active: true
  })

  // 草稿数据
  const templateDraft = computed(() => store.getters['fieldTemplate/templateDraft'])
  const fieldLocalList = computed(() => templateDraft.value.fieldList?.map(unwrapData) ?? [])
  const uniqueLocalList = computed(() => templateDraft.value.uniqueList ?? [])

  // 接口数据
  const bindModelData = ref([])
  const basicData = ref({
    name: '',
    description: ''
  })
  const fieldData = ref([])
  const uniqueData = ref([])

  // 对比的状态
  const isDiffDone = ref(false)
  const hasDiffError = ref(false)
  const hasDiffConflict = ref(false)

  const modelIdList = ref([])

  const modelEditAuths = ref({})

  watchEffect(async () => {
    const [template, templateFieldList, templateUniqueList] = await Promise.all([
      fieldTemplateService.findById(templateId.value),
      fieldTemplateService.getFieldList({ bk_template_id: templateId.value }),
      fieldTemplateService.getUniqueList({ bk_template_id: templateId.value })
    ])
    store.commit('setTitle', `${t('编辑字段组合模板')}【${template.name}】`)

    fieldData.value = templateFieldList?.info || []
    uniqueData.value = templateUniqueList?.info || []

    const modelList = await fieldTemplateService.getBindModel({
      bk_template_id: templateId.value
    })

    basicData.value.name = templateDraft.value.basic.name ?? template.name
    basicData.value.description = templateDraft.value.basic.description ?? template.description
    bindModelData.value = modelList ?? []
    modelIdList.value = modelList.filter(model => !model.bk_ispaused).map(model => model.id)

    // 无绑定的模型diffDone为true
    isDiffDone.value = !modelIdList.value?.length
  })

  // 状态数据实时再算一次，当中如果模板数据在其它地方被意外的修改，可能会出现非预期的数据不一致问题
  // 如果不再算一次，则需要依赖上一步中的fieldStatus数据，当直接进入到此页时并不存在上一页的数据
  const { fieldStatus, removedFieldList } = useField(fieldData, fieldLocalList)
  const { uniqueStatus, removedUniqueList } = useUnique(uniqueData, uniqueLocalList)

  // 模板最终数据，草稿数据优先否则为接口数据
  const templateData = computed(() => ({
    basic: basicData.value,
    fieldList: templateDraft.value.fieldList ?? fieldData.value,
    uniqueList: templateDraft.value.uniqueList ?? uniqueData.value
  }))

  const isDraftValid = computed(() => !basicData.value.name)

  const hasModelEditAuth = computed(() => {
    if (!modelIdList.value?.length) {
      return true
    }
    return Object.values(modelEditAuths.value).every(isPass => isPass)
  })

  // 存在冲突或者没有编辑中的草稿数据
  const submitButtonDisabled = computed(() => isDraftValid.value
    || !isDiffDone.value
    || hasDiffError.value
    || hasDiffConflict.value
    || !hasModelEditAuth.value)

  const clearTemplateDraft = () => {
    store.commit('fieldTemplate/clearTemplateDraft')
  }

  const handleDiffUpdate = (hasError, hasConflict) => {
    isDiffDone.value = true
    hasDiffError.value = hasError
    hasDiffConflict.value = hasConflict
  }

  const finalFieldList = computed(() => normalizeFieldData(templateData.value.fieldList, false, fieldStatus))
  // eslint-disable-next-line max-len
  const finalUniqueList = computed(() => normalizeUniqueData(templateData.value.uniqueList, templateData.value.fieldList, false, uniqueStatus))

  const handleSubmit = async () => {
    const submitData = {
      id: templateId.value,
      ...templateData.value.basic,
      attributes: finalFieldList.value,
      uniques: finalUniqueList.value
    }

    try {
      await fieldTemplateService.update(submitData, { requestId: requestIds.submit })
      leaveConfirmConfig.active = false
      isEditSuccess.value = true

      clearTemplateDraft()
    } catch (err) {
      console.error(err)
    }
  }

  const handlePrevStep = () => {
    leaveConfirmConfig.active = false

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

  const handleModelAuthUpdate = (model, isPass) => {
    set(modelEditAuths.value, model.id, isPass)
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
  <div class="edit-binging">
    <template v-if="!isEditSuccess">
      <top-steps :steps="steps" width="632px" :current="3"></top-steps>
      <bind-model
        :height="`${$APP.height - 161 - 52}px`"
        :readonly="true"
        :template-id="templateId"
        :model-list="bindModelData"
        :field-list="finalFieldList"
        :unique-list="finalUniqueList"
        :removed-field-list="removedFieldList"
        :removed-unique-list="removedUniqueList"
        @update-diffs="handleDiffUpdate"
        @update-model-auth="handleModelAuthUpdate">
      </bind-model>
      <div class="edit-binging-footer">
        <bk-button
          @click="handlePrevStep">
          {{$t('上一步')}}
        </bk-button>
        <cmdb-auth :auth="{ type: $OPERATION.U_FIELD_TEMPLATE, relation: [templateId] }" v-bk-tooltips="{
          disabled: hasModelEditAuth && !hasDiffConflict,
          content: !hasModelEditAuth ? $t('暂无对应模型的编辑权限，请先申请模型的编辑权限') : $t('模型存在冲突，无法提交')
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
        <bk-button theme="default" @click="handleCancel">{{$t('取消')}}</bk-button>
      </div>
    </template>
    <sync-results v-else
      scene="edit"
      :template-id="templateId"
      :model-ids="modelIdList">
    </sync-results>
    <leave-confirm
      v-bind="leaveConfirmConfig"
      :reverse="true"
      :title="$t('是否退出')"
      :content="$t('编辑步骤未完成，退出将撤销当前操作')"
      :ok-text="$t('退出')"
      :cancel-text="$t('取消')"
      @leave="handleLeave">
    </leave-confirm>
  </div>
</template>

<style lang="scss" scoped>
  .edit-binging-footer {
    display: flex;
    align-items: center;
    justify-content: center;
    height: 52px;
    padding: 0 20px;
    background-color: #fff;
    border-top: 1px solid $borderColor;

    .bk-button {
      min-width: 86px;

      & + .bk-button {
        margin-left: 8px;
      }
      & + .auth-box {
          margin-left: 8px;
      }
    }
    .auth-box {
      & + .bk-button,
      & + .auth-box {
          margin-left: 8px;
      }
    }
  }
</style>
