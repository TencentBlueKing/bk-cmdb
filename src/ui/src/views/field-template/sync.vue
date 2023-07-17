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
    MENU_MODEL_FIELD_TEMPLATE
  } from '@/dictionary/menu-symbol'
  import { normalizeFieldData, normalizeUniqueData } from './children/use-field'
  import BindModel from './children/bind-model.vue'
  import SyncResults from './children/sync-results.vue'
  import fieldTemplateService from '@/service/field-template'

  const route = useRoute()
  const store = useStore()

  const templateId = computed(() => Number(route.params.id))
  const modelId = computed(() => Number(route.params.modelId))
  const fieldData = ref([])
  const uniqueData = ref([])

  const bindModelData = ref([])

  const isDiffDone = ref(false)
  const hasDiffError = ref(false)
  const hasDiffConflict = ref(false)

  const bindModelRef = ref(null)
  const requestIds = {
    bind: Symbol('bind')
  }
  const isBindSuccess = ref(false)
  const bindedModelIdList = computed(() => [modelId.value])

  watchEffect(async () => {
    try {
      const [template, templateFieldList, templateUniqueList] = await Promise.all([
        fieldTemplateService.findById(templateId.value),
        fieldTemplateService.getFieldList({ bk_template_id: templateId.value }),
        fieldTemplateService.getUniqueList({ bk_template_id: templateId.value })
      ])

      fieldData.value = templateFieldList?.info || []
      uniqueData.value = templateUniqueList?.info || []

      store.commit('setTitle', `${t('同步模板')}【${template.name}】`)

      const modelList = await fieldTemplateService.getBindModel({
        bk_template_id: templateId.value
      })
      bindModelData.value = modelList.filter(model => model.id === modelId.value)
    } catch (err) {
      console.error(err)
    }
  })

  const finalFieldList = computed(() => normalizeFieldData(fieldData.value, false))
  const finalUniqueList = computed(() => normalizeUniqueData(uniqueData.value, fieldData.value, false))

  const handleDiffUpdate = (hasError, hasConflict) => {
    isDiffDone.value = true
    hasDiffError.value = hasError
    hasDiffConflict.value = hasConflict
  }

  const handleSubmit = async () => {
    isBindSuccess.value = true
  }

  const handleCancel = () => {
    routerActions.redirect({
      name: MENU_MODEL_FIELD_TEMPLATE
    })
  }
</script>

<template>
  <div class="bind">
    <template v-if="!isBindSuccess">
      <bind-model
        ref="bindModelRef"
        :height="`${$APP.height - 111 - 52}px`"
        :readonly="true"
        :template-id="templateId"
        :model-list="bindModelData"
        :binded-model-list="bindModelData"
        :field-list="finalFieldList"
        :unique-list="finalUniqueList"
        @update-diffs="handleDiffUpdate">
      </bind-model>
      <div class="bind-footer">
        <cmdb-auth :auth="{ type: $OPERATION.U_FIELD_TEMPLATE, relation: [templateId] }" v-bk-tooltips="{
          disabled: !hasDiffConflict,
          content: $t('模型存在冲突，无法提交')
        }">
          <template #default="{ disabled }">
            <bk-button
              theme="primary"
              :disabled="disabled || !isDiffDone || hasDiffError || hasDiffConflict"
              :loading="$loading(requestIds.bind)"
              @click="handleSubmit">
              {{$t('提交')}}
            </bk-button>
          </template>
        </cmdb-auth>
        <bk-button theme="default" @click="handleCancel">{{$t('取消')}}</bk-button>
      </div>
    </template>
    <sync-results v-else
      scene="bind"
      :template-id="templateId"
      :model-ids="bindedModelIdList">
    </sync-results>
  </div>
</template>

<style lang="scss" scoped>
  .bind-footer {
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
