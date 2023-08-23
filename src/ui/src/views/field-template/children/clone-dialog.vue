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
  import { computed, ref, watch } from 'vue'
  import cloneDeep from 'lodash/cloneDeep'
  import fieldTemplateService from '@/service/field-template'

  const props = defineProps({
    show: {
      type: Boolean
    },
    sourceTemplate: {
      type: Object,
      default: () => ({})
    }
  })

  const emit = defineEmits(['success', 'toggle'])

  const requestId = Symbol()

  const cloneInput = ref(null)
  const cloneForm = ref(cloneDeep(props.sourceTemplate))
  watch(() => props.sourceTemplate, (sourceTemplate) => {
    cloneForm.value = cloneDeep(sourceTemplate)
    cloneForm.value.name += '-copy'
  }, { deep: true, immediate: true })

  const isShow = computed({
    get() {
      return props.show
    },
    set(val) {
      emit('toggle', val)
    }
  })

  const handleCloneConfirm = async () => {
    try {
      const params = {
        id: cloneForm.value.id,
        name: cloneForm.value.name,
        description: cloneForm.value.description
      }
      const res = await fieldTemplateService.cloneTemplate(params, { requestId })
      emit('success', res)
    } catch (error) {
      console.error(error)
      return false
    }
  }
</script>

<template>
  <bk-dialog
    v-model="isShow"
    theme="primary"
    header-position="left"
    :mask-close="false"
    :auto-close="false"
    :loading="$loading(requestId)"
    width="670"
    :title="$t('克隆字段组合模板')"
    @confirm="handleCloneConfirm">
    <bk-form :label-width="$i18n.locale === 'en' ? 140 : 90" :model="cloneForm">
      <bk-form-item :label="$t('模板名称')" :required="true" :property="'name'"
        class="cmdb-form-item" :class="{ 'is-error': errors.has('name') }">
        <bk-input
          ref="cloneInput"
          v-model="cloneForm.name"
          v-validate="'required|length:256'"
          v-autofocus
          :placeholder="$t('请输入模板名称')"
          name="name">
        </bk-input>
        <p class="form-error" v-if="errors.has('name')">{{errors.first('name')}}</p>
      </bk-form-item>
      <bk-form-item :label="$t('描述')" :property="'description'"
        class="cmdb-form-item" :class="{ 'is-error': errors.has('description') }">
        <bk-input
          type="textarea"
          name="description"
          :rows="4"
          v-model="cloneForm.description"
          :placeholder="$t('请输入模板描述')"
          v-validate="'length:2000'">
        </bk-input>
      </bk-form-item>
    </bk-form>
  </bk-dialog>
</template>
