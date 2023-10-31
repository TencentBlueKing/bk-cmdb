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
  import { computed, reactive } from 'vue'
  import { t } from '@/i18n'
  import { isContainerObject } from '@/service/container/common'
  import { timeFormatter } from '@/filters/formatter'
  import CodeViewer from '@/components/ui/code-viewer/index.vue'

  const props = defineProps({
    value: {
      type: [Object, Array, String],
      default: () => ({})
    },
    property: {
      type: Object,
      default: () => ({})
    },
    showOn: String,
    instance: {
      type: Object,
      default: () => ({})
    }
  })

  const jsonStringValue = computed(() => JSON.stringify(props.value, null, 2))

  const filename = computed(() => {
    if (isContainerObject(props.property.bk_obj_id)) {
      return `${props.instance.id}_${props.property.bk_property_id}_${timeFormatter(Date.now(), 'YYYYMMDDHHmmss')}.json`
    }
    return `${props.property.bk_property_id}_${timeFormatter(Date.now())}.json`
  })

  const detailsSlider = reactive({
    show: false,
    title: `【${props.property.bk_property_name}】${t('详情')}`
  })

  const handleClickView = () => {
    detailsSlider.show = true
  }

  defineExpose({
    getCopyValue: () => jsonStringValue.value
  })
</script>

<template>
  <div class="json-value">
    <div class="view-trigger" @click="handleClickView" v-if="value">
      <i class="bk-cmdb-icon icon-json" />
    </div>
    <div v-else>--</div>

    <bk-sideslider
      v-transfer-dom
      :is-show.sync="detailsSlider.show"
      :title="detailsSlider.title"
      :width="960">
      <div slot="content" class="details-content" v-if="detailsSlider.show">
        <code-viewer :filename="filename" :code="jsonStringValue" lang="json" />
      </div>
    </bk-sideslider>
  </div>
</template>

<style lang="scss" scoped>
  .view-trigger {
    cursor: pointer;
    font-size: 14px;
    color: $primaryColor;
    &:hover {
      color: #699df4;
    }
  }

  .details-content {
    padding: 20px 24px;
    height: 100%;
  }
</style>
