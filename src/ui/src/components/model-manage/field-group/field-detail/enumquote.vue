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
  import { computed, nextTick, ref, watch, watchEffect, inject } from 'vue'
  import { t } from '@/i18n'
  import router from '@/router/index.js'
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
    },
    isReadOnly: {
      type: Boolean,
      default: false
    }
  })

  const emit = defineEmits(['input', 'update:multiple'])

  const customObjId = inject('customObjId')

  const defaultValueSelectEl = ref(null)

  const refModelId = ref('')
  const refModelInstIds = ref(props.multiple ? [] : '')

  const searchPlaceholder = computed(() => t('请输入xx', { name: t(refModelId.value === BUILTIN_MODELS.HOST ? 'IP' : '名称') }))
  const excludeModelIds = computed(() => ([
    BUILTIN_MODELS.SET,
    BUILTIN_MODELS.MODULE,
    router.app.$route.params.modelId ?? customObjId
  ]))

  const isMultiple = computed({
    get() {
      return props.multiple
    },
    set(val) {
      emit('update:multiple', val)
    }
  })

  watchEffect(() => {
    if (props.value?.length) {
      refModelId.value = props.value.map(item => item.bk_obj_id)?.[0]
      refModelInstIds.value = props.value.map(item => item.bk_inst_id)
    } else {
      refModelInstIds.value = []
    }
  })

  watch(() => props.multiple, () => nextTick(async () => defaultValueSelectEl.value.$validator.validate('refModelInst')))

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
        <div class="form-item-row">
          <model-instance-selector
            ref="defaultValueSelectEl"
            class="model-instance-selector"
            name="refModelInst"
            data-vv-validate-on="change"
            v-validate="`required|maxSelectLength:${ multiple ? -1 : 1 }`"
            :obj-id="refModelId"
            :placeholder="$t('请选择xx', { name: $t('模型实例') })"
            :search-placeholder="searchPlaceholder"
            :display-tag="true"
            :multiple="true"
            v-model="refModelInstIds"
            @change="handleModelInstChange">
          </model-instance-selector>
          <bk-checkbox
            class="checkbox"
            v-model="isMultiple"
            :disabled="isReadOnly">
            <span>{{$t('可多选')}}</span>
          </bk-checkbox>
        </div>
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
  .form-item-row {
    display: flex;
    align-items: center;
    gap: 12px;
    .checkbox {
      flex: none;
    }
  }
</style>
