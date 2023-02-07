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
  <div class="model-group-list">
    <div
      class="model-group-item"
      v-for="(modelGroup, groupIndex) in groups"
      :key="modelGroup.bk_classification_id"
      :class="{
        'is-collapse':
          modelGroupCollapseStates[modelGroup.bk_classification_id]
      }"
    >
      <div
        class="model-group-header"
        @click="toggleModelGroupCollapse(modelGroup.bk_classification_id)"
      >
        <bk-icon class="model-group-collapse-icon" type="down-shape" />
        {{ modelGroup.bk_classification_name }}（
        {{ modelGroup.bk_objects.length }} ）
        <slot name="group-header-append" v-bind="{ modelGroup }"></slot>
      </div>
      <bk-transition name="collapse" duration-type="ease">
        <ul
          class="model-list"
          v-show="!modelGroupCollapseStates[modelGroup.bk_classification_id]"
        >
          <li
            class="model-item"
            :class="{
              'is-active': selectedModelId === model.bk_obj_id
            }"
            v-for="(model, modelIndex) in modelGroup.bk_objects"
            @click="selectModel(model)"
            :key="model.bk_obj_id"
          >
            <model-summary :data="model"></model-summary>
            <slot name="model-append"
              v-bind="{
                modelGroup,
                groupIndex,
                model,
                modelIndex
              }">
            </slot>
          </li>
        </ul>
      </bk-transition>
    </div>
  </div>
</template>

<script>
  import { defineComponent, reactive, toRef } from 'vue'
  import ModelSummary from './model-summary.vue'

  export default defineComponent({
    name: 'ModelGroupList',
    components: {
      ModelSummary
    },
    props: {
      // 模型分组数据
      groups: {
        type: Array,
        required: true,
        default: () => []
      },
      // 选中模型的 bk_obj_id
      selectedModelId: {
        type: String,
        default: ''
      }
    },
    setup(props, { emit }) {
      const modelGroupCollapseStates = reactive({})

      const toggleModelGroupCollapse = (groupId) => {
        const collapse = toRef(modelGroupCollapseStates, groupId)

        collapse.value = !collapse.value
      }

      const selectModel = (model) => {
        emit('update:selectedModelId', model.bk_obj_id)
        emit('model-select', model)
      }

      return {
        selectModel,
        modelGroupCollapseStates,
        toggleModelGroupCollapse
      }
    },
  })
</script>

<style lang="scss" scoped>
.model-group {
  &-list {
    font-size: 12px;
  }

  &-header {
    font-weight: 600;
    height: 40px;
    display: flex;
    align-items: center;
    cursor: pointer;
  }

  &-collapse-icon {
    margin-left: 12px;
    margin-right: 6px;
    transition: transform 200ms ease;

    @at-root .is-collapse & {
      transform: rotate(-90deg);
    }
  }

  &-item {
    border-bottom: 1px solid $borderColor;
  }
}

.model {
  &-item {
    position: relative;
    z-index: 1;
    display: flex;
    align-items: center;
    height: 40px;
    cursor: pointer;
    width: calc(100% - 1px);
    transition: background-color 200ms ease;

    &::before {
      content: "";
      display: block;
      position: absolute;
      top: 0;
      right: 0;
      left: 30px;
      height: 1px;
      background-color: $borderColor;
    }

    .model-summary {
      margin-left: 45px;
    }

    &:hover {
      background-color: #fff;
    }

    &.is-active {
      background-color: #fff;
      width: 100%;
    }

  }
}
</style>
