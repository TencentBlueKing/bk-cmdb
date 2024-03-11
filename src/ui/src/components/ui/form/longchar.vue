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
    ref="autoSize"
    :placeholder="placeholder || $t('请输入长字符')"
    :disabled="disabled"
    :type="'textarea'"
    :rows="row"
    :maxlength="maxlength"
    :clearable="!disabled"
    @blur="handleBlur"
    @enter="handleEnter"
    @change="handleChange"
    @on-change="handleChange">
  </bk-input>
</template>

<script>
  import { calcTextareaHeight } from '@/utils/util.js'
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
    data() {
      return {
        height: 0,
        minHeight: 0
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
    mounted() {
      this.init()
    },
    methods: {
      init() {
        const { autoSize } = this.$refs
        const { textarea } = autoSize.$refs
        const parent = autoSize.$el.querySelector('.bk-textarea-wrapper')
        const { height, minHeight } = calcTextareaHeight(textarea, this.row)

        this.minHeight = minHeight + 4
        this.height = height

        textarea.style.height = `${minHeight}px`
        autoSize.$el.style.height = `${this.minHeight}px`
        parent.style.height = `${this.minHeight}px`
        parent.style.minHeight = `${this.minHeight}px`
      },
      handleChange(value) {
        this.$emit('on-change', value)
        const { autoSize } = this.$refs
        const { textarea } = autoSize.$refs
        this.$nextTick(() => {
          const { height } = calcTextareaHeight(textarea, this.row)
          textarea.style.height = `${height}px`
        })
      },
      handleEnter(value) {
        this.$emit('enter', value)
      },
      handleBlur(value) {
        this.$refs.autoSize.$el.querySelector('.bk-textarea-wrapper').style.height = `${this.minHeight}px`
        this.$emit('blur', value)
      },
      focus() {
        const { autoSize } = this.$refs
        const { textarea } = autoSize.$refs
        textarea.style.height = `${this.height}px`
        this.$el.querySelector('textarea').focus()
      }
    }
  }
</script>

<style lang="scss" scoped>
    .bk-form-control {
        position: relative;
        height: 32px;
        &.control-active {
          :deep(.bk-textarea-wrapper) {
            height: auto !important;
          }
        }

        :deep(.control-icon) {
            z-index: 3;
        }

        /deep/ .bk-textarea-wrapper {
            position: absolute;
            width: 100%;
            height: 32px;
            min-height: 32px;
            z-index: 2;
            @include scrollbar-y;
            &:hover {
              height: auto !important;
            }

            .bk-form-textarea {
                min-height: 28px;
                max-height: 400px;
                padding: 5px 10px;
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
