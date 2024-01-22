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
  <span class="search-input-wrapper" v-if="multiple">
    <bk-input :min="min" :max="max"
      class="search-input" type="number" v-model="start" :precision="precision"
      :allow-number-paste="allowPaste" v-on="listeners">
    </bk-input>
    <span class="search-input-grep">-</span>
    <bk-input :min="min" :max="max"
      class="search-input" type="number" v-model="end" :precision="precision"
      :allow-number-paste="allowPaste" v-on="listeners"></bk-input>
  </span>
  <bk-input v-else
    class="search-input"
    type="number"
    :min="min"
    :max="max"
    v-model="localValue"
    :precision="precision"
    :allow-number-paste="allowPaste"
    v-on="listeners">
  </bk-input>
</template>

<script>
  import activeMixin from './mixins/active'
  import numberFormTypeMixin from '@/mixins/number-form-type'
  export default {
    name: 'cmdb-search-float',
    mixins: [activeMixin, numberFormTypeMixin],
    props: {
      value: {
        type: [Number, String, Array],
        default: ''
      },
      precision: {
        type: Number,
        default: 5
      },
      allowPaste: {
        type: Boolean,
        default: true
      }
    },
    data() {
      return {
        listeners: {
          focus: () => this.handleToggle(true),
          blur: () => this.handleToggle(false),
          clear: () => this.$emit('clear')
        }
      }
    },
    computed: {
      multiple() {
        return Array.isArray(this.value)
      },
      localValue: {
        get() {
          return String(this.value) === 'NaN' ? '' : this.value
        },
        set(value) {
          let newValue
          if (Array.isArray(value)) {
            newValue = value.map(number => (number.toString().length ? Number(number) : number))
          } else {
            newValue = value.toString().length ? Number(value) : value
          }
          this.$emit('input', newValue)
          this.$emit('change', newValue)
        }
      },
      start: {
        get() {
          const [start] = this.value
          return start
        },
        set(start) {
          const [, end = ''] = this.value
          this.localValue = [start, end]
        }
      },
      end: {
        get() {
          const [, end] = this.value
          return end
        },
        set(end) {
          const [start = ''] = this.value
          this.localValue = [start, end]
        }
      }
    }
  }
</script>

<style lang="scss" scoped>
    .search-input-wrapper {
        display: inline-flex;
        .search-input-grep {
            flex: 20px 0 0;
            text-align: center;
        }
        .search-input {
            flex: 1;
        }
    }
</style>
