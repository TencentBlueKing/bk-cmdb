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
  import { v4 as uuidv4 } from 'uuid'
  import { t } from '@/i18n'
  import { computed, ref } from 'vue'
  import previewField from '@/components/model-manage/field-group/preview-field.vue'
  import { PROPERTY_TYPES } from '@/dictionary/property-constants'
  import { defaultFieldData } from './use-field'

  const props = defineProps({
    previewShow: {
      type: Boolean,
      default: false
    },
    properties: {
      type: Array,
      default: () => []
    }
  })

  const groups = ref([{
    bk_biz_id: 0,
    bk_group_id: 'default',
    bk_group_index: 0,
    bk_group_name: 'Default',
    bk_isdefault: true,
    bk_obj_id: '',
  }])

  const defaultData = ref([
    {
      ...defaultFieldData(),
      id: uuidv4(),
      bk_property_id: 'demo',
      bk_property_name: t('示例字段'),
      bk_property_type: PROPERTY_TYPES.SINGLECHAR,
      bk_property_group: 'default',
      bk_property_group_name: 'Default',
      placeholder: '',
      isrequired: false
    }
  ])

  const perviewData = computed(() => {
    if (props.properties?.length) {
      return props.properties.map(item => ({
        ...item,
        bk_property_group: 'default',
        bk_property_group_name: 'Default',
        placeholder: item.placeholder.value,
        isrequired: item.isrequired.value
      }))
    }
    return defaultData.value
  })

  const emit = defineEmits(['update:previewShow'])

  const isPreviewShow = computed({
    get() {
      return props.previewShow
    },
    set(val) {
      emit('update:previewShow', val)
    }
  })

  const handleCloseSilder = () => {
    emit('update:previewShow', false)
  }

</script>

<template>
  <bk-sideslider
    ref="sidesliderComp"
    v-transfer-dom
    :width="676"
    :title="$t('字段预览')"
    :is-show="isPreviewShow"
    :before-close="handleCloseSilder">
    <preview-field v-if="isPreviewShow"
      slot="content"
      :properties="perviewData"
      :property-groups="groups">
    </preview-field>
  </bk-sideslider>
</template>
