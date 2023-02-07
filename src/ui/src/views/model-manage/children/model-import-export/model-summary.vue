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

<template>
  <div class="model-summary">
    <span class="model-icon-container" :class="{ 'is-builtin': isBuiltin }">
      <i
        class="model-icon"
        :class="modelIconClass"
      ></i>
    </span>
    <span class="model-name">{{ modelName }}</span>
    <span class="model-id">({{ modelId }})</span>
  </div>
</template>
<script>
  import { defineComponent, computed } from 'vue'

  export default defineComponent({
    name: 'ModelSummary',
    props: {
      /**
       * 模型数据
       * @property {Object} data
       * @property {Boolean} data.ispre 是否为内置模型
       * @property {String} data.bk_obj_icon 模型图标
       * @property {String} data.bk_obj_id 模型 ID
       * @property {String} data.bk_obj_name 模型名称
       */
      data: {
        type: Object,
        required: true,
        default: null
      }
    },
    setup({ data }) {
      const isBuiltin = computed(() => data.ispre)
      const modelIconClass = computed(() => data.bk_obj_icon || '')
      const modelName = computed(() => data.bk_obj_name || '')
      const modelId = computed(() => data.bk_obj_id || '')

      return {
        isBuiltin,
        modelIconClass,
        modelName,
        modelId
      }
    }
  })
</script>

<style lang="scss" scoped>
.model {
  &-summary {
    display: flex;
    font-size: 12px;
    align-items: center;
  }

  &-icon-container {
    width: 26px;
    height: 26px;
    background: #e1ecff;
    border-radius: 50%;
    display: flex;
    align-items: center;
    justify-content: center;

    @at-root .is-builtin & {
      background-color: #f0f1f5;
    }
  }

  &-icon {
    color: #3a84ff;

    @at-root .is-builtin & {
      color: #798aad;
    }
  }

  &-name {
    margin-left: 12px;
  }

  &-id {
    color: #979ba5;
    margin-left: 5px;
  }
}
</style>
