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
  import { computed, ref } from 'vue'
  import { useStore } from '@/store'
  import routerActions from '@/router/actions'
  import {
    MENU_MODEL_FIELD_TEMPLATE,
    MENU_MODEL_FIELD_TEMPLATE_CREATE_FIELD_SETTINGS
  } from '@/dictionary/menu-symbol'
  import TopSteps from './children/top-steps.vue'
  import BasicForm from './children/basic-form.vue'

  const store = useStore()

  const nextButtonDisabled = ref(false)
  const basicFormRef = ref(null)

  const basicDefaultData = {
    name: '',
    description: ''
  }

  const templateDraft = computed(() => store.getters['fieldTemplate/templateDraft'])
  const basicData = computed(() => ({ ...basicDefaultData, ...templateDraft.value.basic }))

  const handleNextStep = async () => {
    if (!await basicFormRef.value.$validator.validateAll()) {
      return
    }
    const { formData } = basicFormRef.value
    store.commit('fieldTemplate/setTemplateDraft', { basic: formData })
    routerActions.redirect({
      name: MENU_MODEL_FIELD_TEMPLATE_CREATE_FIELD_SETTINGS,
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
      if (![MENU_MODEL_FIELD_TEMPLATE_CREATE_FIELD_SETTINGS].includes(to.name)) {
        this.$store.commit('fieldTemplate/clearTemplateDraft')
      }
      next()
    }
  }
</script>

<template>
  <cmdb-sticky-layout class="cmdb-config-sticky-layout">
    <template #header="{ sticky }">
      <top-steps width="360px" :class="{ 'is-sticky': sticky }"></top-steps>
    </template>
    <basic-form :data="basicData" ref="basicFormRef"></basic-form>
    <template #footer="{ sticky }">
      <div :class="['layout-footer', { 'is-sticky': sticky }]">
        <cmdb-auth :auth="{ type: $OPERATION.C_FIELD_TEMPLATE }">
          <template #default="{ disabled }">
            <bk-button
              theme="primary"
              :disabled="nextButtonDisabled || disabled"
              @click="handleNextStep">
              {{$t('下一步')}}
            </bk-button>
          </template>
        </cmdb-auth>
        <bk-button theme="default" @click="handleCancel">{{$t('取消')}}</bk-button>
      </div>
    </template>
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
