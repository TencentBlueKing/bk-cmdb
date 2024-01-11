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
  <bk-tag-input ref="tagInput"
    allow-create
    allow-auto-match
    v-if="multiple"
    v-model="localValue"
    v-bind="$attrs"
    :list="[]"
    :has-delete-icon="true"
    @removeAll="() => $emit('clear')"
    @click.native="handleToggle(true)"
    @blur="handleToggle(false, ...arguments)">
  </bk-tag-input>
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
    name: 'cmdb-search-singlechar',
    mixins: [activeMixin],
    props: {
      value: {
        type: [Array, String],
        default: () => ([])
      },
      /**
       * value 为外部输入，用 value 的数据类型来控制匹配模式不可靠，所以增加 fuzzy 属性来确定匹配模式。如果传入 fuzzy 则优先使用 fuzzy 来进行模式的切换，否则使用 value
       */
      fuzzy: {
        type: Boolean,
        default: undefined
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
    watch: {
      multiple: {
        immediate: true,
        handler(multiple) {
          multiple ? this.addPasteEvent() : this.removePasteEvent()
        }
      }
    },
    beforeDestroy() {
      this.removePasteEvent()
    },
    methods: {
      async addPasteEvent() {
        await this.$nextTick()
        const { tagInput } = this.$refs
        if (!tagInput) return
        tagInput.$refs.input?.addEventListener('paste', this.handlePaste)
      },
      async removePasteEvent() {
        await this.$nextTick()
        const { tagInput } = this.$refs
        if (!tagInput) return
        tagInput.$refs.input?.removeEventListener('paste', this.handlePaste)
      },
      handlePaste(event) {
        const text = event.clipboardData.getData('text')
        const values = text.split(/,|;|\n/).map(value => value.trim())
          .filter(value => value.length)
        const value = [...new Set([...this.localValue, ...values])]
        this.localValue = value
      }
    }
  }
</script>
