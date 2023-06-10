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
  import Draggable from 'vuedraggable'
  import FieldCard from './field-card.vue'

  const props = defineProps({
    fieldList: {
      type: Array,
      default: () => ([])
    },
    disabledSort: {
      type: Boolean,
      default: false
    }
  })

  const emit = defineEmits(['click-field', 'sorted', 'sort-change'])

  const isDragging = ref(false)

  const fieldListLocal = computed({
    get() {
      return props.fieldList
    },
    set(value) {
      emit('sorted', value)
    }
  })

  const handleClickField = (field) => {
    emit('click-field', field)
  }
  const handleRemoveField = (field, index) => {
    emit('remove-field', field, index)
  }

  const handleDragStart = () => {
    isDragging.value = true
  }
  const handleDragEnd = () => {
    isDragging.value = false
  }
  const handleDragChange = (event) => {
    emit('sort-change', event)
  }
</script>

<template>
  <draggable
    tag="div"
    :class="['field-grid', { dragging: isDragging }]"
    v-model="fieldListLocal"
    ghost-class="field-item-ghost"
    draggable=".field-item"
    :animation="150"
    :disabled="disabledSort"
    @start="handleDragStart"
    @end="handleDragEnd"
    @change="handleDragChange">
    <slot
      name="field-card"
      v-for="(field, index) in fieldListLocal"
      v-bind="{ field, itemClass: 'field-item', index }">
      <field-card
        class="field-item"
        :key="index"
        :field-index="index"
        :field="field"
        @click-field="handleClickField"
        @remove-field="handleRemoveField">
      </field-card>
    </slot>
  </draggable>
</template>

<style lang="scss" scoped>
.field-grid {
  display: grid;
  gap: 16px;
  grid-template-columns: repeat(auto-fill, minmax(260px, 1fr));
  width: 100%;
  align-content: flex-start;

  :deep(.field-item-ghost) {
    background-color: #f5f7fa !important;
    border: 1px dashed #dcdee5;

    &:hover {
      border-color: #dcdee5;
      background-color: #f5f7fa;
      box-shadow: none;
    }

    > * {
      display: none !important;
    }
  }

  &.dragging {
    :deep(.field-item) {
      &:hover {
        border-color: #dcdee5;
        background-color: #fff;
        .drag-icon {
          visibility: hidden;
        }
        .field-button {
          visibility: hidden;
        }
      }
    }
  }
}
</style>
