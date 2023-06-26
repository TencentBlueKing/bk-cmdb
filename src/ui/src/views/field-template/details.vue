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
  import { computed, ref, watchEffect } from 'vue'
  import { useRoute } from '@/router/index'
  import { t } from '@/i18n'
  import { $success } from '@/magicbox/index.js'
  import CmdbTab from '@/components/ui/tab/index.vue'
  import DetailsField from './details-field.vue'
  import DetailsUnique from './details-unique.vue'
  import DetailsModel from './details-model.vue'
  import FieldPreviewDrawer from './children/field-preview-drawer.vue'
  import CloneDialog from './children/clone-dialog.vue'
  import fieldTemplateService from '@/service/field-template'
  import routerActions from '@/router/actions'
  import {
    MENU_MODEL_FIELD_TEMPLATE_EDIT_FIELD_SETTINGS,
    MENU_MODEL_FIELD_TEMPLATE_EDIT_BINDING
  } from '@/dictionary/menu-symbol'
  import useTemplate from './children/use-template'

  const props = defineProps({
    open: {
      type: Boolean,
      default: false
    },
    template: {
      type: Object,
      default: () => ({})
    }
  })

  const emit = defineEmits(['close', 'bind-change', 'clone-done', 'delete-done'])

  const route = useRoute()

  const templateId = computed(() => props.template?.id)
  const queryTab = computed(() => route.query.tab)

  const templateLocal = computed(() => [props.template])
  const { handleDelete: handleDeleteTemplate } = useTemplate(templateLocal)

  const fieldCount = ref('')
  const modelCount = ref('')
  const uniqueList = ref([])
  const previewFieldList = ref([])
  const previewShow = ref(false)

  const isShowCloneDialog = ref(false)

  watchEffect(async () => {
    const [fieldCounts, templateUniqueList, modelCounts] = await Promise.all([
      fieldTemplateService.getFieldCount({ bk_template_ids: [templateId.value] }),
      fieldTemplateService.getUniqueList({ bk_template_id: templateId.value }),
      fieldTemplateService.getModelCount({ bk_template_ids: [templateId.value] })
    ])

    fieldCount.value = fieldCounts?.[0]?.count
    modelCount.value = modelCounts?.[0]?.count
    uniqueList.value = templateUniqueList?.info || []
  })

  const isShow = computed({
    get() {
      return props.open
    },
    set() {
      emit('close')
    }
  })

  const tabIds = {
    field: Symbol(),
    unique: Symbol(),
    model: Symbol(),
  }
  const tabActive = ref(tabIds[queryTab.value] || tabIds.field)
  const tabs = computed(() => ([
    {
      id: tabIds.field,
      text: t('字段配置'),
      count: fieldCount.value
    },
    {
      id: tabIds.unique,
      text: t('唯一校验'),
      count: uniqueList.value?.length
    },
    {
      id: tabIds.model,
      text: t('绑定的模型'),
      count: modelCount.value
    }
  ]))

  const handleEdit = () => {
    const editRoutes = {
      [tabIds.field]: MENU_MODEL_FIELD_TEMPLATE_EDIT_FIELD_SETTINGS,
      [tabIds.unique]: MENU_MODEL_FIELD_TEMPLATE_EDIT_FIELD_SETTINGS,
      [tabIds.model]: MENU_MODEL_FIELD_TEMPLATE_EDIT_BINDING,
    }
    routerActions.redirect({
      name: editRoutes[tabActive.value],
      params: {
        id: templateId.value
      }
    })

    emit('close')
  }

  const handlePreviewField = () => {
    previewShow.value = true
  }
  const handleClone = () => {
    isShowCloneDialog.value = true
  }
  const handleDelete = () => {
    handleDeleteTemplate(props.template, () => {
      $success(t('删除成功'))
      emit('delete-done')
      emit('close')
    })
  }

  const handleTabChange = (tab) => {
    tabActive.value = tab.id
  }

  const handleSliderHidden = () => {
    tabActive.value = tabIds.field
    emit('close')
  }

  const handleSliderBeforeClose = () => true

  const handleFieldUpdated = (list) => {
    previewFieldList.value = list
  }
  const handleModelUnbound = async () => {
    const modelCounts = await fieldTemplateService.getModelCount({ bk_template_ids: [templateId.value] })
    modelCount.value = modelCounts?.[0]?.count
    emit('bind-change', templateId.value)
  }

  const handleCloneSuccess = () => {
    $success(t('克隆成功'))
    isShowCloneDialog.value = false
    emit('clone-done')
  }
  const handleCloneDialogToggle = (val) => {
    isShowCloneDialog.value = val
  }
</script>
<script>
  export default {
    name: 'template-details'
  }
</script>

<template>
  <bk-sideslider
    v-transfer-dom
    :width="960"
    :title="`${$t('字段组合模板详情')}【${template.name}】`"
    :is-show.sync="isShow"
    :quick-close="true"
    :before-close="handleSliderBeforeClose"
    @hidden="handleSliderHidden">
    <cmdb-sticky-layout slot="content" class="content">
      <div class="content-head">
        <div class="data-row">
          <div class="data-value title">{{ template.name }}</div>
        </div>
        <div class="data-row">
          <div class="data-key">{{$t('描述：')}}</div>
          <div class="data-value">{{ template.description || '--' }}</div>
        </div>
        <cmdb-tab class="details-tab" :tabs="tabs" :active="tabActive" @change="handleTabChange" />
      </div>
      <div class="content-body">
        <details-field
          v-if="tabActive === tabIds.field"
          :template-id="templateId"
          :unique-list="uniqueList"
          @updated="handleFieldUpdated">
        </details-field>
        <details-unique
          v-if="tabActive === tabIds.unique"
          :template-id="templateId"
          :unique-list="uniqueList">
        </details-unique>
        <details-model
          v-if="tabActive === tabIds.model"
          :template-id="templateId"
          @unbound="handleModelUnbound">
        </details-model>
      </div>
      <template slot="footer" slot-scope="{ sticky }">
        <div class="action-bar" :class="{ 'is-sticky': sticky }">
          <bk-button theme="primary" @click="handleEdit">
            {{ $t('进入编辑') }}
          </bk-button>
          <bk-button theme="default" @click="handlePreviewField">
            {{$t('预览字段')}}
          </bk-button>
          <bk-button theme="default" @click="handleClone">
            {{$t('克隆')}}
          </bk-button>
          <bk-button theme="default" @click="handleDelete">
            {{$t('删除')}}
          </bk-button>
        </div>
      </template>

      <field-preview-drawer
        :preview-show.sync="previewShow"
        :properties="previewFieldList">
      </field-preview-drawer>

      <clone-dialog
        :show="isShowCloneDialog"
        :source-template="template"
        @success="handleCloneSuccess"
        @toggle="handleCloneDialogToggle">
      </clone-dialog>
    </cmdb-sticky-layout>
  </bk-sideslider>
</template>

<style lang="scss" scoped>
  .content {
    height: 100%;
    // 取消外部滚动条
    overflow: hidden;

    .content-head {
      padding: 24px 40px 0 40px;
      background: #F5F7FA;

      .data-row {
        display: flex;
        align-items: center;
        font-size: 12px;
        margin-bottom: 12px;

        .data-key {
          color: #63656E;
        }
        .data-value {
          color: #313238;
        }
        .title {
          font-weight: 700;
          font-size: 14px;
        }
      }

      .details-tab {
        margin-top: 24px;
      }
    }
    .content-body {
      padding: 0 32px;
      background: #FFF;
    }
    .action-bar {
      display: flex;
      gap: 8px;
      padding: 8px 40px;
      background: #fff;
      &.is-sticky {
        border-top: 1px solid #dcdee5;
        padding: 8px 24px;
      }
      .bk-button{
        width: 88px;
        height: 32px;
      }
    }
  }
</style>
