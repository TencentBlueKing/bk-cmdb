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
  <bk-input
    v-model="localValue"
    v-bind="$attrs"
    :placeholder="placeholder || $t('请输入长字符')"
    :disabled="disabled"
    :type="'textarea'"
    :rows="row"
    :maxlength="maxlength"
    :clearable="!disabled"
    @blur="handleBlur"
    @enter="handleEnter"
    @on-change="handleChange">
  </bk-input>
</template>

<script>
  export default {
    name: 'cmdb-form-longchar',
    props: {
      value: {
        type: [String, Number],
        default: ''
      },
      disabled: {
        type: Boolean,
        default: false
      },
      maxlength: {
        type: Number,
        default: 2000
      },
      minlength: {
        type: Number,
        default: 2000
      },
      row: {
        type: Number,
        default: 3
      },
      placeholder: {
        type: String,
        default: ''
      }
    },
    computed: {
      localValue: {
        get() {
          return (this.value === null || this.value === undefined) ? '' : this.value
        },
        set(value) {
          this.$emit('input', value)
        }
      }
    },
    methods: {
      handleChange(value) {
        this.$emit('on-change', value)
      },
      handleEnter(value) {
        this.$emit('enter', value)
      },
      handleBlur(value) {
        this.$emit('blur', value)
      },
      focus() {
        this.$el.querySelector('textarea').focus()
      }
    }
  }
</script>

<style lang="scss" scoped>
    .bk-form-control {
        /deep/ .bk-textarea-wrapper {
            .bk-form-textarea {
                min-height: auto !important;
                padding: 5px 10px 8px;
                @include scrollbar-y;
                &.textarea-maxlength {
                    margin-bottom: 0 !important;
                }
            }
        }
        /deep/ .bk-limit-box {
            display: none !important
        }
    }
</style>
