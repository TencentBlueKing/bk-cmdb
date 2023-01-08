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
  import { t } from '@/i18n'
  import { BUILTIN_MODELS } from '@/dictionary/model-constants.js'
  import GridLayout from '@/components/ui/other/grid-layout.vue'
  import GridItem from '@/components/ui/other/grid-item.vue'
  import ModelSelector from '@/components/model-instance/model-selector.vue'
  import ModelInstanceSelector from '@/components/model-instance/model-instance-selector.vue'

  const props = defineProps({
    value: {
      type: [Array, String],
      default: ''
    },
    multiple: {
      type: Boolean,
      default: true
    }
  })

  const emit = defineEmits(['input'])

  const excludeModelIds = [BUILTIN_MODELS.SET, BUILTIN_MODELS.MODULE]

  const refModelId = ref('')
  const refModelInstIds = ref(props.multiple ? [] : '')

  const searchPlaceholder = computed(() => t('请输入xx', { name: t(refModelId.value === BUILTIN_MODELS.HOST ? 'IP' : '名称') }))

  watchEffect(() => {
    if (props.multiple) {
      if (props.value?.length) {
        refModelId.value = props.value.map(item => item.bk_obj_id)?.[0]
        refModelInstIds.value = props.value.map(item => item.bk_inst_id)
      } else {
        refModelInstIds.value = []
      }
    } else {
      if (props.value?.length) {
        refModelId.value = props.value.map(item => item.bk_obj_id)?.[0]
        refModelInstIds.value = props.value.map(item => item.bk_inst_id)?.[0]
      } else {
        refModelInstIds.value = ''
      }
    }
  })

  const handleModelInstChange = (modelInstIds) => {
    const instIds = Array.isArray(modelInstIds) ? modelInstIds : [modelInstIds]
    const option = instIds.map(instId => ({
      bk_obj_id: refModelId.value,
      bk_inst_id: instId,
      type: 'int'
    }))

    emit('input', option)
  }
</script>

<template>
  <div>
    <grid-layout mode="form" :gap="36" :font-size="'14px'" :max-columns="1">
      <grid-item
        direction="column"
        required
        :class="['cmdb-form-item', 'form-item', { 'is-error': errors.has('refModel') }]"
        :label="$t('引用模型')">
        <model-selector
          class="model-selector"
          searchable
          name="refModel"
          v-validate="'required'"
          :exclude="excludeModelIds"
          :placeholder="$t('请选择xx', { name: $t('模型') })"
          v-model="refModelId">
        </model-selector>
        <template #append>
          <div class="form-error" v-if="errors.has('refModel')">{{errors.first('refModel')}}</div>
          <div class="tips" v-else>{{$t('默认以实例名称作为枚举选项')}}</div>
        </template>
      </grid-item>
      <grid-item
        direction="column"
        required
        :class="['cmdb-form-item', 'form-item', { 'is-error': errors.has('refModelInst') }]"
        :label="$t('默认值')">
        <model-instance-selector
          :key="String(multiple)"
          class="model-instance-selector"
          name="refModelInst"
          v-validate="'required'"
          :obj-id="refModelId"
          :placeholder="$t('请选择xx', { name: $t('模型实例') })"
          :search-placeholder="searchPlaceholder"
          :display-tag="true"
          :multiple="multiple"
          v-model="refModelInstIds"
          @change="handleModelInstChange">
        </model-instance-selector>
        <template #append>
          <div class="form-error" v-if="errors.has('refModelInst')">{{errors.first('refModelInst')}}</div>
        </template>
      </grid-item>
    </grid-layout>
  </div>
</template>

<style lang="scss" scoped>
  .model-selector,
  .model-instance-selector {
    width: 100%;
  }

  .form-item {
    position: relative;
    .tips {
      position: absolute;
      top: 100%;
      left: 0;
      font-size: 12px;
      margin-top: 4px;
    }
  }
</style>
