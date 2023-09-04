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
    <span class="label-text">{{$t('字段设置')}}</span>
    <div class="checkbox-options">
      <label v-if="isEditableShow">
        <bk-checkbox
          class="checkbox"
          v-model="localValue.editable"
          :disabled="isReadOnly || ispre">
          <span class="g-has-dashed-tooltips" v-bk-tooltips="$t('字段设置可编辑提示语')">
            {{$t('可编辑')}}
          </span>
          <i class="bk-cc-icon icon-cc-tips disabled-tips"
            v-if="modelId === 'host'"
            v-bk-tooltips="$t('主机属性设置为不可编辑状态后提示')"></i>
        </bk-checkbox>
      </label>
      <label class="ml30" v-if="isRequiredShow && !isMainLineModel">
        <bk-checkbox
          class="checkbox"
          v-model="localValue.isrequired"
          :disabled="isReadOnly || ispre">
          <span>{{$t('必填')}}</span>
        </bk-checkbox>
      </label>
      <label class="ml30" v-if="isMultipleShow">
        <bk-checkbox
          class="checkbox"
          v-model="localValue.multiple"
          :disabled="isReadOnly || ispre">
          <span>{{$t('可多选')}}</span>
        </bk-checkbox>
      </label>
    </div>
  </div>
</template>

<script>
  import { PROPERTY_TYPES } from '@/dictionary/property-constants'

  export default {
    props: {
      isReadOnly: {
        type: Boolean,
        default: false
      },
      type: {
        type: String,
        required: true
      },
      editable: {
        type: Boolean,
        default: true
      },
      isrequired: {
        type: Boolean,
        default: false
      },
      multiple: {
        type: Boolean,
        default: false
      },
      isMainLineModel: {
        type: Boolean,
        default: false
      },
      ispre: Boolean,
      isEditField: {
        type: Boolean,
        default: false
      }
    },
    data() {
      return {
        editableMap: [
          PROPERTY_TYPES.SINGLECHAR,
          PROPERTY_TYPES.INT,
          PROPERTY_TYPES.FLOAT,
          PROPERTY_TYPES.ENUM,
          PROPERTY_TYPES.DATE,
          PROPERTY_TYPES.TIME,
          PROPERTY_TYPES.LONGCHAR,
          PROPERTY_TYPES.OBJUSER,
          PROPERTY_TYPES.TIMEZONE,
          PROPERTY_TYPES.BOOL,
          PROPERTY_TYPES.LIST,
          PROPERTY_TYPES.ORGANIZATION,
          PROPERTY_TYPES.ENUMMULTI,
          PROPERTY_TYPES.ENUMQUOTE,
          PROPERTY_TYPES.INNER_TABLE
        ],
        isRequiredMap: [
          PROPERTY_TYPES.SINGLECHAR,
          PROPERTY_TYPES.INT,
          PROPERTY_TYPES.FLOAT,
          PROPERTY_TYPES.DATE,
          PROPERTY_TYPES.TIME,
          PROPERTY_TYPES.LONGCHAR,
          PROPERTY_TYPES.OBJUSER,
          PROPERTY_TYPES.TIMEZONE,
          PROPERTY_TYPES.LIST,
          PROPERTY_TYPES.ORGANIZATION,
          PROPERTY_TYPES.INNER_TABLE
        ],
        isMultipleMap: [
          PROPERTY_TYPES.ORGANIZATION,
          PROPERTY_TYPES.ENUMQUOTE,
          PROPERTY_TYPES.ENUMMULTI
        ],
        localValue: {
          editable: this.editable,
          isrequired: this.isrequired,
          multiple: this.multiple
        }
      }
    },
    inject: ['customObjId'], // 来源于自定义字段编辑
    computed: {
      isEditableShow() {
        return this.editableMap.indexOf(this.type) !== -1
      },
      isRequiredShow() {
        return this.isRequiredMap.indexOf(this.type) !== -1
      },
      isMultipleShow() {
        return this.isMultipleMap.indexOf(this.type) !== -1
      },
      modelId() {
        return this.$route.params.modelId ?? this.customObjId
      }
    },
    watch: {
      editable(editable) {
        this.localValue.editable = editable
      },
      isrequired(isrequired) {
        this.localValue.isrequired = isrequired
      },
      multiple(multiple) {
        this.localValue.multiple = multiple
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
      'localValue.multiple'(multiple) {
        this.$emit('update:multiple', multiple)
      }
    }
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
