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
  <div class="form-label">
    <div class="checkbox-options">
      <flex-row v-if="isEditableShow" gap="2px">
        <bk-checkbox
          v-model="localValue.editable"
          class="checkbox"
          :disabled="
            isReadOnly || ispre || (isFromTemplateField && editableLockLocal)
          ">
          <span
            v-bk-tooltips="$t('字段设置可编辑提示语')"
            class="g-has-dashed-tooltips">
            {{ $t('在实例中可编辑') }}
          </span>
          <i
            v-if="modelId === 'host'"
            v-bk-tooltips="$t('主机属性设置为不可编辑状态后提示')"
            class="bk-cc-icon icon-cc-tips disabled-tips"></i>
        </bk-checkbox>
      </flex-row>
      <flex-row
        v-if="isRequiredShow && !isMainLineModel"
        gap="2px"
        class="mt20">
        <bk-checkbox
          v-model="localValue.isrequired"
          class="checkbox"
          :disabled="
            isReadOnly || ispre || (isFromTemplateField && isrequiredLockLocal)
          ">
          <span>{{ $t('设置为必填项') }}</span>
        </bk-checkbox>
      </flex-row>
    </div>
  </div>
</template>

<script>
import { EDITABLE_TYPES, REQUIRED_TYPES } from '@/dictionary/property-constants'
import FlexRow from '@/components/ui/other/flex-row.vue'

export default {
  components: {
    FlexRow,
  },
  inject: [
    // 来源于自定义字段编辑
    'customObjId',
    'isSettingScene',
    'isFromTemplateField',
  ],
  props: {
    isReadOnly: {
      type: Boolean,
      default: false,
    },
    type: {
      type: String,
      required: true,
    },
    editable: {
      type: Boolean,
      default: true,
    },
    isrequired: {
      type: Boolean,
      default: false,
    },
    isMainLineModel: {
      type: Boolean,
      default: false,
    },
    ispre: Boolean,
    isEditField: {
      type: Boolean,
      default: false,
    },
    isrequiredLock: {
      type: Boolean,
      default: false,
    },
    editableLock: {
      type: Boolean,
      default: false,
    },
  },
  data() {
    return {
      localValue: {
        editable: this.editable,
        isrequired: this.isrequired,
      },
    }
  },
  computed: {
    isEditableShow() {
      return EDITABLE_TYPES.indexOf(this.type) !== -1
    },
    isRequiredShow() {
      return REQUIRED_TYPES.indexOf(this.type) !== -1
    },
    modelId() {
      return this.$route.params.modelId ?? this.customObjId
    },
    isrequiredLockLocal: {
      get() {
        return this.isrequiredLock
      },
      set(val) {
        this.$emit('update:isrequiredLock', val)
      },
    },
    editableLockLocal: {
      get() {
        return this.editableLock
      },
      set(val) {
        this.$emit('update:editableLock', val)
      },
    },
  },
  watch: {
    editable(editable) {
      this.localValue.editable = editable
    },
    isrequired(isrequired) {
      this.localValue.isrequired = isrequired
    },
    'localValue.editable'(editable) {
      this.$emit('update:editable', editable)
    },
    'localValue.isrequired'(isrequired) {
      if (!isrequired && this.isOnlyShow) {
        this.localValue.isonly = false
      }
      this.$emit('update:isrequired', isrequired)
    },
  },
}
</script>

<style lang="scss" scoped>
.disabled-tips {
  font-size: 12px;
  margin-left: 6px;
}

.checkbox-options {
  margin-bottom: 10px;

  .checkbox {
    height: 24px;
    line-height: 24px;
  }
}
</style>
