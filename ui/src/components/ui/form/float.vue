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
  <bk-input class="cmdb-form-float" :type="inputType" ref="input"
    :placeholder="placeholder || $t('请输入浮点数')"
    :value="value"
    :disabled="disabled"
    :precision="precision"
    v-bind="$attrs"
    @blur="handleInput"
    @enter="handleEnter"
    @change="handleChange">
    <template slot="append" v-if="unit">
      <div class="unit" :title="unit">{{unit}}</div>
    </template>
  </bk-input>
</template>

<script>
  export default {
    name: 'cmdb-form-float',
    props: {
      value: {
        default: null,
        validator(val) {
          return ['string', 'number'].includes(typeof val) || val === null
        }
      },
      disabled: {
        type: Boolean,
        default: false
      },
      placeholder: {
        type: String,
        default: ''
      },
      unit: {
        type: String,
        default: ''
      },
      autoCheck: {
        type: Boolean,
        default: true
      },
      inputType: {
        type: String,
        default: 'text'
      },
      precision: {
        type: Number,
        default: 5
      }
    },
    computed: {
      localValue: {
        get() {
          return this.value === null ? '' : this.value
        },
        set(value) {
          const emitValue = value === '' ? null : value
          this.$emit('input', emitValue)
          this.$emit('change', emitValue)
          this.$emit('on-change', emitValue)
        }
      }
    },
    methods: {
      handleInput(value, event) {
        const originalValue = String(event.target.value).trim()
        const floatValue = originalValue.length ? Number(event.target.value.trim()) : null
        if (isNaN(floatValue)) {
          value = this.autoCheck ? null : value
        } else {
          value = floatValue
        }
        this.$refs.input.curValue = value
        this.localValue = value
      },
      handleChange() {
        this.$emit('on-change', this.localValue)
      },
      handleEnter() {
        this.$emit('enter', this.localValue)
      },
      focus() {
        this.$el.querySelector('input').focus()
      }
    }
  }
</script>

<style lang="scss" scoped>
.cmdb-form-float {
  .unit {
    max-width: 120px;
    font-size: 12px;
    @include ellipsis;
    padding: 0 10px;
    height: 30px;
    line-height: 30px;
    background: #f2f4f8;
    color: #63656e;
  }

  &[size="small"] {
    .unit {
      height: 24px;
      line-height: 24px;
    }
  }
}
</style>
