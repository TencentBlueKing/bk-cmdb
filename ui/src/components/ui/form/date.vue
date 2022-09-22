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
  <bk-date-picker class="cmdb-date"
    v-bind="otherAttrs"
    v-model="date"
    transfer
    editable
    :clearable="clearable"
    :placeholder="placeholder"
    :disabled="disabled">
  </bk-date-picker>
</template>

<script>
  export default {
    name: 'cmdb-form-date',
    props: {
      value: {
        type: String,
        default: ''
      },
      disabled: {
        type: Boolean,
        default: false
      },
      clearable: {
        type: Boolean,
        default: true
      },
      placeholder: {
        type: String,
        default: ''
      }
    },
    computed: {
      date: {
        get() {
          if (!this.value) {
            return ''
          }
          return new Date(this.value)
        },
        set(value) {
          const previousValue = this.value
          const currentValue = this.$tools.formatTime(value, 'YYYY-MM-DD')
          this.$emit('input', currentValue)
          if (currentValue !== previousValue) {
            this.$emit('change', currentValue, previousValue)
          }
        },
      },
      otherAttrs() {
        // 排除options属性，因与date-picker组件props类型冲突，不能直接用
        const { options, ...otherAttrs } = this.$attrs
        return otherAttrs
      }
    }
  }
</script>

<style lang="scss" scoped>
  .cmdb-date {
    width: 100%;

    &[size="small"] {
      ::v-deep {
        .bk-date-picker-rel {
          .icon-wrapper {
            height: 26px;
            line-height: 26px;
          }
          .bk-date-picker-editor {
            height: 26px;
            line-height: 26px;
          }
        }
      }
    }
  }
</style>
