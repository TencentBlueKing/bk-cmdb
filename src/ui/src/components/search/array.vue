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
      v-model="localValue"
      v-bind="$attrs"
      :list="[]"
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
  export default {
    name: 'cmdb-search-array',
    mixins: [activeMixin],
    props: {
      value: {
        type: [Array, String],
        default: () => ([])
      },
      fuzzy: {
        type: Boolean,
        default: undefined
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
      localValue: {
        get() {
          return this.value
        },
        set(value) {
          this.$emit('input', value)
          this.$emit('change', value)
        }
      }
    },
    methods: {
      handlePasteFn(value) {
        if (!value) return
        let val = [value]
        if (this.isPasteSplit && this.multiple) {
          val = (value.split(/,|;|\n/)).map(value => value.trim())
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
