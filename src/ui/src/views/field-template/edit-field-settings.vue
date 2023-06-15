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
  import { computed, ref, watchEffect, reactive } from 'vue'
  import { t } from '@/i18n'
  import { useRoute } from '@/router/index'
  import { useStore } from '@/store'
  import routerActions from '@/router/actions'
  import {
    MENU_MODEL_FIELD_TEMPLATE,
    MENU_MODEL_FIELD_TEMPLATE_EDIT_BASIC,
    MENU_MODEL_FIELD_TEMPLATE_EDIT_BINDING
  } from '@/dictionary/menu-symbol'
  import TopSteps from './children/top-steps.vue'
  import FieldManage from './children/field-manage.vue'
  import { wrapData } from './children/use-field'
  import fieldTemplateService from '@/service/field-template'
  import FieldPreview from './children/field-preview-drawer.vue'

  const route = useRoute()
  const store = useStore()
  const steps = [
    { title: '基础信息', icon: 1 },
    { title: '字段设置', icon: 2 },
    { title: '模型信息确认', icon: 3 }
  ]

  const templateId = computed(() => Number(route.params.id))

  // 模板初始数据
  const fieldData = ref([])
  const uniqueData = ref([])
  const basicData = ref({
    name: '',
    description: ''
  })
  // 编辑前最初的数据，接口原始数据
  const beforeFieldList = ref([])
  const beforeUniqueList = ref([])

  const previewShow = ref(false)

  const nextButtonDisabled = computed(() => !fieldData.value.length)

  const templateDraft = computed(() => store.getters['fieldTemplate/templateDraft'])

  watchEffect(async () => {
    const [template, templateFieldList, templateUniqueList] = await Promise.all([
      fieldTemplateService.findById(templateId.value),
      fieldTemplateService.getFieldList({ bk_template_id: templateId.value }),
      fieldTemplateService.getUniqueList({ bk_template_id: templateId.value })
    ])

    beforeFieldList.value = templateFieldList?.info || []
    beforeUniqueList.value = templateUniqueList?.info || []

    // 如果存在草稿，优先使用
    fieldData.value = templateDraft.value.fieldList ?? (templateFieldList?.info || [])
    uniqueData.value = templateDraft.value.uniqueList ?? (templateUniqueList?.info || [])

    basicData.value.name = templateDraft.value.basic.name ?? template.name
    basicData.value.description = templateDraft.value.basic.description ?? template.description

    store.commit('setTitle', `${t('编辑字段组合模板')}【${template.name}】`)
  })

  // 模板最终数据，编辑后的数据优先否则为初始数据
  const templateData = computed(() => ({
    basic: basicData.value,
    fieldList: settingData.fieldList ?? fieldData.value,
    uniqueList: settingData.uniqueList ?? uniqueData.value
  }))

  const previewFieldList = computed(() => templateData.value.fieldList)

  // 编辑后的数据
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
  const saveDraft = () => {
    // 数据写入草稿
    store.commit('fieldTemplate/setTemplateDraft', {
      basic: templateData.value.basic,
      fieldList: templateData.value.fieldList,
      uniqueList: templateData.value.uniqueList
    })
  }

  const handlePrevStep = () => {
    saveDraft()
    routerActions.redirect({
      name: MENU_MODEL_FIELD_TEMPLATE_EDIT_BASIC,
      history: true
    })
  }
  const handleNextStep = () => {
    saveDraft()
    routerActions.redirect({
      name: MENU_MODEL_FIELD_TEMPLATE_EDIT_BINDING,
      history: true
    })
  }
  const handleCancel = () => {
    routerActions.redirect({
      name: MENU_MODEL_FIELD_TEMPLATE
    })
  }
  const handlePreview = () => {
    previewShow.value = true
  }
</script>
<script>
  export default {
    beforeRouteLeave(to, from, next) {
      if (![MENU_MODEL_FIELD_TEMPLATE_EDIT_BASIC, MENU_MODEL_FIELD_TEMPLATE_EDIT_BINDING].includes(to.name)) {
        this.$store.commit('fieldTemplate/clearTemplateDraft')
      }
      next()
    }
  }
</script>

<template>
  <cmdb-sticky-layout class="cmdb-config-sticky-layout">
    <template #header="{ sticky }">
      <top-steps :steps="steps" width="632px" :current="2" :class="{ 'is-sticky': sticky }"></top-steps>
    </template>
    <field-manage
      :field-list="fieldData"
      :unique-list="uniqueData"
      :before-field-list="beforeFieldList"
      :before-unique-list="beforeUniqueList"
      :is-create-mode="false"
      @update-field="handleFieldUpdate"
      @update-unique="handleUniqueUpdate">
    </field-manage>
    <template #footer="{ sticky }">
      <div :class="['layout-footer', { 'is-sticky': sticky }]">
        <bk-button
          @click="handlePrevStep">
          {{$t('上一步')}}
        </bk-button>
        <bk-button
          theme="primary"
          :disabled="nextButtonDisabled"
          @click="handleNextStep">
          {{$t('下一步')}}
        </bk-button>
        <bk-button
          @click="handlePreview">
          {{$t('预览')}}
        </bk-button>
        <bk-button theme="default" @click="handleCancel">{{$t('取消')}}</bk-button>
      </div>
    </template>
    <field-preview
      :preview-show.sync="previewShow"
      :properties="previewFieldList">
    </field-preview>
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
