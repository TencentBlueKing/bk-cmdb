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
  import { BUILTIN_MODELS, BUILTIN_MODEL_PROPERTY_KEYS } from '@/dictionary/model-constants'

  const props = defineProps({
    property: {
      type: Object,
      default: () => ({})
    },
    instance: {
      type: Object,
      default: () => ({})
    },
    showOn: String
  })

  defineExpose({
    getCopyValue: () => ''
  })

  const objId = computed(() => props.property?.bk_obj_id)
  const instanceId = computed(() => {
    if (BUILTIN_MODELS.HOST === objId.value) {
      return props.instance?.[objId.value]?.[BUILTIN_MODEL_PROPERTY_KEYS[objId.value]?.ID]
    }
    return props.instance?.[BUILTIN_MODEL_PROPERTY_KEYS[objId.value]?.ID || 'bk_inst_id']
  })

  const isShowDialog = ref(false)
  const dialogTitle = computed(() => props.property.bk_property_name)

  const handleClickView = () => {
    isShowDialog.value = true
  }
</script>

<template>
  <div class="inner-table-value">
    <cmdb-form-innertable v-if="showOn === 'default'"
      :property="property"
      :obj-id="objId"
      :instance-id="instanceId"
      :immediate="false"
      :readonly="true" />
    <div class="view-trigger" v-else-if="showOn === 'cell'" @click="handleClickView">
      <i class="bk-cmdb-icon icon-cc-table" />
    </div>
    <bk-dialog
      v-model="isShowDialog"
      width="960"
      render-directive="if"
      :title="dialogTitle"
      header-position="left"
      :mask-close="true"
      :show-footer="false">
      <cmdb-form-innertable
        :mode="'info'"
        :property="property"
        :obj-id="objId"
        :instance-id="instanceId"
        :readonly="true" />
    </bk-dialog>
  </div>
</template>

<style lang="scss" scoped>
  .view-trigger {
    cursor: pointer;
    font-size: 12px;
    color: $primaryColor;
    &:hover {
      color: #699df4;
    }
  }
</style>
