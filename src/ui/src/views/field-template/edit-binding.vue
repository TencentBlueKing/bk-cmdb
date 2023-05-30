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
  import { ref, computed, watchEffect } from 'vue'
  import { t } from '@/i18n'
  import { useRoute } from '@/router/index'
  import { useStore } from '@/store'
  import routerActions from '@/router/actions'
  import {
    MENU_MODEL_FIELD_TEMPLATE,
    MENU_MODEL_FIELD_TEMPLATE_EDIT_FIELD_SETTINGS
  } from '@/dictionary/menu-symbol'
  import TopSteps from './children/top-steps.vue'
  import BindModel from './children/bind-model.vue'
  import SyncResults from './children/sync-results.vue'
  import useField, { normalizeFieldData, normalizeUniqueData, unwrapData } from './children/use-field'
  import fieldTemplateService from '@/service/field-template'

  const route = useRoute()
  const store = useStore()
  const steps = [
    { title: '基础信息', icon: 1 },
    { title: '字段设置', icon: 2 },
    { title: '模型信息确认', icon: 3 }
  ]
  const requestIds = {
    submit: Symbol('submit')
  }

  const templateId = computed(() => Number(route.params.id))

  const isEditSuccessed = ref(false)

  // 草稿数据
  const templateDraft = computed(() => store.getters['fieldTemplate/templateDraft'])
  const fieldLocalList = computed(() => templateDraft.value.fieldList.map(unwrapData))

  // 接口数据
  const bindModelData = ref([])
  const fieldData = ref([])
  const uniqueData = ref([])

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
      bk_template_id: templateId.value,
      // filter: {}
    })
    console.log(modelList, template, 'modelList, template')
    bindModelData.value = modelList
  })

  // 状态数据实时再算一次，当中如果模板数据在其它地方被意外的修改，可能会出现非预期的数据不一致问题
  // 如果不再算一次，则需要依赖上一步中的fieldStatus数据，当直接进入到此页时并不存在上一页的数据
  const { fieldStatus } = useField(fieldData, fieldLocalList)

  const isDraftValid = computed(() => !templateDraft.value.basic.name)
  // 存在冲突或者没有编辑中的草稿数据
  const submitButtonDisabled = computed(() => isDraftValid.value)

  const handleSubmit = async () => {
    console.log(templateDraft)
    const submitData = {
      id: templateId.value,
      ...templateDraft.value.basic,
      attributes: normalizeFieldData(templateDraft.value.fieldList, false, fieldStatus),
      uniques: normalizeUniqueData(templateDraft.value.uniqueList, templateDraft.value.fieldList, false)
    }

    console.log(submitData, fieldStatus)

    try {
      await fieldTemplateService.update(submitData, { requestId: requestIds.submit })
      isEditSuccessed.value = true
    } catch (err) {
      console.error(err)
    }
  }

  const handlePrevStep = () => {
    routerActions.redirect({
      name: MENU_MODEL_FIELD_TEMPLATE_EDIT_FIELD_SETTINGS,
      history: true
    })
  }
  const handleCancel = () => {
    routerActions.redirect({
      name: MENU_MODEL_FIELD_TEMPLATE
    })
  }
</script>
<script>
  export default {
    beforeRouteLeave(to, from, next) {
      if (![MENU_MODEL_FIELD_TEMPLATE_EDIT_FIELD_SETTINGS].includes(to.name)) {
        this.$store.commit('fieldTemplate/clearTemplateDraft')
      }
      next()
    }
  }
</script>

<template>
  <div class="edit-binging">
    <template v-if="!isEditSuccessed">
      <top-steps :steps="steps" width="632px" :current="3"></top-steps>
      <bind-model :height="`${$APP.height - 161 - 52}px`" :model-list="bindModelData"></bind-model>
      <div class="edit-binging-footer">
        <bk-button
          @click="handlePrevStep">
          {{$t('上一步')}}
        </bk-button>
        <cmdb-auth :auth="{ type: $OPERATION.U_FIELD_TEMPLATE, relation: [templateId] }">
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
    <div class="sync-container" v-else>
      <sync-results></sync-results>
    </div>
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

  .sync-container {
    padding: 20px 24px;
  }
</style>
