<!--
 * Tencent is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2017 Tencent. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
-->

<template>
  <div class="g-expand" v-if="multiple">
    <bk-tag-input ref="tagInput"
      allow-create
      allow-auto-match
      :collapse-tags="true"
      :create-tag-validator="tagValidator"
      v-model="localValue"
      v-bind="$attrs"
      :list="[]"
      :has-delete-icon="true"
      :paste-fn="handlePasteFn"
      @removeAll="() => $emit('clear')"
      @click.native="handleToggle(true)"
      @blur="handleToggle(false, ...arguments)"
      @inputchange="handleInputChange">
    </bk-tag-input>
  </div>
  <bk-input v-else
    v-model.trim="localValue"
    v-bind="$attrs"
    @clear="() => $emit('clear')"
    @focus="handleToggle(true, ...arguments)"
    @blur="handleToggle(false, ...arguments)">
  </bk-input>
</template>

<script>
  import activeMixin from './mixins/active'
  import { isNumeric } from '@/utils/util'
  export default {
    name: 'cmdb-search-singlechar',
    mixins: [activeMixin],
    props: {
      value: {
        type: [Array, String, Number],
        default: () => ([])
      },
      /**
       * value 为外部输入，用 value 的数据类型来控制匹配模式不可靠，所以增加 fuzzy 属性来确定匹配模式。如果传入 fuzzy 则优先使用 fuzzy 来进行模式的切换，否则使用 value
       */
      fuzzy: {
        type: Boolean,
        default: undefined
      },
      // 只可输入数字
      onlyNumber: {
        type: Boolean,
        default: false
      },
      isPasteSplit: {
        type: Boolean,
        default: false
      }
    },
    computed: {
      multiple() {
        if (typeof this.fuzzy === 'boolean') {
          return !this.fuzzy
        }
        return Array.isArray(this.value)
      },
      tagValidator() {
        if (this.onlyNumber) {
          return isNumeric
        }
        return null
      },
      localValue: {
        get() {
          return this.value
        },
        set(value) {
          let newValue = this.$tools.clone(value)
          if (this.onlyNumber) {
            newValue = value.map(val => +val)
          }
          this.$emit('input', newValue)
          this.$emit('change', newValue)
        }
      }
    },
    methods: {
      handlePasteFn(value) {
        if (!value) return
        if (this.onlyNumber) {
          // 如果当前需要分割功能，在只能输入number的时候，需要将value去掉分隔符
          const detectionValue = this.isPasteSplit ? value.replace(/,|;|\n|\s/g, '') : value
          if (!isNumeric(detectionValue)) return
        }
        let val = [value]
        if (this.isPasteSplit && this.multiple) {
          val = (value.split(/,|;|\n|\s/)).map(value => value.trim())
            .filter(value => value.length)
        }
        this.localValue = [...new Set([...this.localValue, ...val])]
        return this.localValue
      },
      handleInputChange(value) {
        this.$emit('inputchange', value)
      }
    }
  }
</script>
