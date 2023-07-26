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
import { UNIUQE_TYPES } from '@/dictionary/model-constants'

// eslint-disable-next-line no-unused-vars
const props = defineProps({
  field: {
    type: Object,
  },
  fieldIndex: {
    type: Number,
  },
  sortable: {
    type: Boolean,
    default: true,
  },
  deletable: {
    type: Boolean,
    default: true,
  },
  fieldUnique: {
    type: Object,
    default: () => ({
      list: [],
      type: '',
    }),
  },
  isTemplate: {
    type: Boolean,
    default: false,
  },
  removeDisabled: {
    type: Boolean,
    default: false,
  },
  removeDisabledTips: {
    type: String,
    default: '',
  },
})

const emit = defineEmits(['click-field'])

const getUniqueRuleContent = uniqueList => {
  const content = []
  uniqueList.forEach(item => content.push(item.names.join(', ')))
  return content.join('<br />')
}

const handleClickRemove = field => {
  emit('remove-field', field)
}

const handleClickField = field => {
  emit('click-field', field)
}
</script>

<template>
  <div :class="['field-card', { sortable }]" @click="handleClickField(field)">
    <span v-if="sortable" class="drag-icon"></span>
    <div class="field-icon">
      <i :class="`icon-cc-field-${field.bk_property_type}`"></i>
    </div>
    <div class="field-info">
      <div class="field-name-area">
        <span class="field-name" :title="field.bk_property_name">{{
          field.bk_property_name
        }}</span>
        <span
          v-if="field.isrequired?.value ?? field.isrequired"
          class="field-required"
          >*</span
        >
        <slot name="flag-append"></slot>
      </div>
      <div class="field-id-area">
        <span class="field-id">{{ field.bk_property_id }}</span>
      </div>
    </div>
    <div
      v-if="isTemplate || $slots['tag-append'] || fieldUnique.list.length"
      class="tags">
      <span
        v-if="
          fieldUnique.type === UNIUQE_TYPES.SINGLE && fieldUnique.list.length
        "
        class="tag unique">
        <em class="tag-text">{{ $t('单独唯一') }}</em>
      </span>
      <span
        v-else-if="fieldUnique.type === UNIUQE_TYPES.UNION"
        v-bk-tooltips="{
          content: getUniqueRuleContent(fieldUnique.list),
        }"
        class="tag unique union">
        <em class="tag-text">{{ $t('联合唯一') }}</em>
      </span>
      <span v-if="isTemplate" class="tag template"
        ><em class="tag-text">{{ $t('模板') }}</em></span
      >
      <slot name="tag-append"></slot>
    </div>
    <div class="field-action" @click.stop>
      <bk-button
        v-if="deletable"
        class="field-button"
        :text="true"
        :disabled="fieldUnique.list.length > 0 || removeDisabled"
        @click.stop="handleClickRemove(field, fieldIndex)">
        <bk-icon
          v-bk-tooltips="{
            disabled: !fieldUnique.list.length && !removeDisabled,
            content:
              fieldUnique.list.length > 0
                ? $t('不允许删除在唯一校验中的字段')
                : removeDisabledTips,
          }"
          class="field-button-icon"
          type="delete" />
      </bk-button>
      <slot name="action-append" v-bind="{ field, fieldIndex }"></slot>
    </div>
  </div>
</template>

<style lang="scss" scoped>
.field-card {
  display: flex;
  align-items: center;
  position: relative;
  height: 60px;
  border: 1px solid transparent;
  background: #fff;
  box-shadow: 0 2px 4px 0 #1919290d;
  border-radius: 2px;
  padding: 0 12px;
  user-select: none;
  cursor: pointer;

  &.sortable {
    padding: 0 12px 0 0;
  }

  .field-icon {
    color: #989ca8;
    width: 24px;
    height: 24px;
    font-size: 22px;
    text-align: center;
    line-height: 24px;
  }

  .field-info {
    flex: 1;
    margin-left: 12px;
    width: 0;
  }

  .field-name-area {
    display: flex;
    gap: 4px;
    position: relative;
    top: 4px;
  }

  .field-name {
    font-size: 12px;

    @include ellipsis;
  }

  .field-id {
    font-size: 12px;
    color: #c4c6cc;
  }

  .field-required {
    font-size: 12px;
    color: #ff5656;
  }

  .tags {
    display: flex;
    position: absolute;
    gap: 2px;
    right: 0;
    top: 0;

    .tag {
      background: #e4faf0;
      border-radius: 2px;
      padding: 1px 4px;
      height: 16px;
      line-height: 16px;
      white-space: nowrap;
      display: flex;
      align-items: center;

      .tag-text {
        display: block;
        font-size: 12px;
        font-style: normal;
        transform: scale(0.875);
      }

      &.template {
        color: #14a568;
        background: #e4faf0;
        border-radius: 2px;
      }

      &.unique {
        color: #3a84ff;
        background: #f0f5ff;
        border-radius: 2px;

        &.union {
          cursor: pointer;

          &:hover {
            background: #e1ecff;
          }
        }
      }
    }
  }

  .drag-icon {
    @include dragIcon;

    visibility: hidden;
    margin: 0 4px;
  }

  .field-button {
    visibility: hidden;
    color: #63656e;

    &:hover {
      color: #3a84ff;
    }

    .field-button-icon {
      font-size: 14px;
    }

    &.is-disabled {
      color: #c4c6cc;
    }
  }

  &:hover {
    background: #f0f5ff;
    border: 1px solid #3a84ff;

    .drag-icon {
      visibility: visible;
    }

    .field-button {
      visibility: visible;
    }
  }
}
</style>
