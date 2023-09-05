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
  import { computed, ref, watchEffect, watch, nextTick } from 'vue'
  import cloneDeep from 'lodash/cloneDeep'
  import { useRoute } from '@/router/index'
  import { t } from '@/i18n'
  import { $success } from '@/magicbox/index.js'
  import CmdbTab from '@/components/ui/tab/index.vue'
  import EditableField from '@/components/ui/details/editable-field.vue'
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

  const emit = defineEmits(['close', 'bind-change', 'clone-done', 'delete-done', 'update-template'])

  const route = useRoute()

  const templateId = computed(() => props.template?.id)
  const queryTab = computed(() => route.query.tab)

  const templateList = computed(() => [cloneDeep(props.template)])
  const templateName = computed({
    get() {
      return templateList.value?.[0]?.name
    },
    set(val) {
      emit('update-template', templateId.value, val, 'name')
    }
  })
  const templateDesc = computed({
    get() {
      return templateList.value?.[0]?.description
    },
    set(val) {
      emit('update-template', templateId.value, val, 'description')
    }
  })
  const { handleDelete: handleDeleteTemplate } = useTemplate(templateList)

  const fieldCount = ref('')
  const modelCount = ref('')
  const uniqueList = ref([])
  const previewFieldList = ref([])
  const previewShow = ref(false)

  const isNameEditing = ref(false)
  const isDescEditing = ref(false)
  watch([isNameEditing, isDescEditing], ([nameEditing, descEditing]) => {
    nextTick(() => {
      if (nameEditing) {
        isDescEditing.value = false
      }
      if (descEditing) {
        isNameEditing.value = false
      }
    })
  })

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

  const handleEdit = (id, routeKey) => {
    const editRoutes = {
      [tabIds.field]: MENU_MODEL_FIELD_TEMPLATE_EDIT_FIELD_SETTINGS,
      [tabIds.unique]: MENU_MODEL_FIELD_TEMPLATE_EDIT_FIELD_SETTINGS,
      [tabIds.model]: MENU_MODEL_FIELD_TEMPLATE_EDIT_BINDING,
    }
    const query = {}
    if (!routeKey && tabActive.value === tabIds.unique) {
      query.action = 'openUnqiueDrawer'
    }
    routerActions.redirect({
      name: editRoutes[routeKey || tabActive.value],
      params: {
        id: id ?? templateId.value
      },
      query
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
    isDescEditing.value = false
    isNameEditing.value = false
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

  const handleCloneSuccess = (res) => {
    const { id } = res
    if (!id) return
    handleEdit(id, tabIds.field)
    $success(t('克隆成功'))
  }
  const handleCloneDialogToggle = (val) => {
    isShowCloneDialog.value = val
  }
  const handleClose = () => {
    emit('close')
  }

  const handleSaveTemplate = async ({ value, confirm, stop }, dataKey) => {
    try {
      await fieldTemplateService.updateBaseInfo({
        id: templateId.value,
        name: templateName.value,
        [dataKey]: value
      })
      confirm()
    } catch (err) {
      stop()
      console.log(err)
    }
  }
  const handleSaveName = (arg) => {
    handleSaveTemplate(arg, 'name')
  }
  const handleSaveDesc = async (arg) => {
    handleSaveTemplate(arg, 'description')
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
          <div class="data-value title">
            <editable-field
              class="editable-field-name"
              v-model="templateName"
              :editing.sync="isNameEditing"
              font-size="12px"
              validate="required|length:256"
              :placeholder="$t('请输入模板名称')"
              :auth="{ type: $OPERATION.U_FIELD_TEMPLATE, relation: [templateId] }"
              @confirm="handleSaveName">
            </editable-field>
          </div>
        </div>
        <div class="data-row">
          <div class="data-key">{{$t('描述：')}}</div>
          <div class="data-value desc">
            <editable-field
              class="editable-field-desc"
              v-model="templateDesc"
              :editing.sync="isDescEditing"
              type="longchar"
              :rows="4"
              font-size="12px"
              validate="length:2000"
              :placeholder="$t('请输入模板描述')"
              :auth="{ type: $OPERATION.U_FIELD_TEMPLATE, relation: [templateId] }"
              @confirm="handleSaveDesc">
            </editable-field>
          </div>
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
          @unbound="handleModelUnbound"
          @close="handleClose">
        </details-model>
      </div>
      <template slot="footer" slot-scope="{ sticky }">
        <div class="action-bar" :class="{ 'is-sticky': sticky }">
          <cmdb-auth class="mr10" :auth="{ type: $OPERATION.U_FIELD_TEMPLATE, relation: [templateId] }">
            <template #default="{ disabled }">
              <bk-button theme="primary" @click="() => handleEdit()" :disabled="disabled">
                {{ $t('进入编辑') }}
              </bk-button>
            </template>
          </cmdb-auth>
          <bk-button theme="default" @click="handlePreviewField">
            {{$t('预览字段')}}
          </bk-button>
          <cmdb-auth class="mr10" :auth="[
            { type: $OPERATION.C_FIELD_TEMPLATE },
            { type: $OPERATION.U_FIELD_TEMPLATE, relation: [templateId] }
          ]">
            <template #default="{ disabled }">
              <bk-button theme="default" @click="handleClone" :disabled="disabled">
                {{$t('克隆')}}
              </bk-button>
            </template>
          </cmdb-auth>
          <cmdb-auth :auth="{ type: $OPERATION.D_FIELD_TEMPLATE, relation: [templateId] }">
            <template #default="{ disabled }">
              <bk-button theme="default" @click="handleDelete" :disabled="disabled || modelCount > 0">
                {{$t('删除')}}
              </bk-button>
            </template>
          </cmdb-auth>
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
        margin-bottom: 4px;
        line-height: 34px;

        .data-key {
          color: #63656E;
          flex: none;
          align-self: start;
        }
        .data-value {
          color: #313238;
        }
        .title {
          width: 368px;
          font-weight: 700 !important;
          font-size: 14px;
        }
        .desc {
          width: 368px;
        }

        .editable-field-name,
        .editable-field-desc {
          width: 100%;
          vertical-align: initial;
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
