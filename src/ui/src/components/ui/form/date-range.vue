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
  <bk-date-picker style="width: 100%"
    v-bind="otherAttrs"
    v-model="time"
    transfer
    :font-size="fontSize"
    :placeholder="placeholder || $t('选择日期范围')"
    :clearable="clearable"
    :type="timer ? 'datetimerange' : 'daterange'"
    :disabled="disabled">
  </bk-date-picker>
</template>

<script>
  export default {
    name: 'cmdb-form-date-range',
    props: {
      value: {
        type: [Array, String],
        default() {
          return []
        }
      },
      disabled: {
        type: Boolean,
        default: false
      },
      timer: Boolean,
      clearable: {
        type: Boolean,
        default: true
      },
      placeholder: {
        type: String,
        default: ''
      },
      fontSize: {
        type: [String, Number],
        default: 'medium'
      }
    },
    data() {
      return {
        localValue: [...this.value]
      }
    },
    computed: {
      time: {
        get() {
          return this.localValue.map(date => (date ? new Date(date) : ''))
        },
        set(value) {
          const localValue = value.map(date => this.$tools.formatTime(date, this.timer ? 'YYYY-MM-DD HH:mm:ss' : 'YYYY-MM-DD'))
          this.localValue = localValue.filter(date => !!date)
        }
      },
      otherAttrs() {
        // 排除options属性，因与date-picker组件props类型冲突，不能直接用
        const { options, ...otherAttrs } = this.$attrs
        return otherAttrs
      }
    },
    watch: {
      value(value) {
        if ([...value].join('') !== this.localValue.join('')) {
          this.localValue = [...value]
        }
      },
      localValue(value, oldValue) {
        if (value.join('') !== [...this.value].join('')) {
          this.$emit('input', [...value])
          this.$emit('change', [...value], [...oldValue])
        }
      }
    }
  }
</script>
