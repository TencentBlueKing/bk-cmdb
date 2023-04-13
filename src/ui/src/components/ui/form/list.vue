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
  <bk-select class="form-list-selector"
    v-model="selected"
    :clearable="allowClear"
    :searchable="searchable"
    :disabled="disabled"
    :multiple="multiple"
    :placeholder="placeholder"
    :font-size="fontSize"
    :popover-options="{
      boundary: 'window'
    }"
    v-bind="$attrs"
    ref="selector">
    <bk-option v-for="option in options"
      :key="option.value"
      :id="option.value"
      :name="option.name">
    </bk-option>
  </bk-select>
</template>

<script>
  export default {
    name: 'cmdb-form-list',
    props: {
      value: {
        type: [Array, String],
        default: ''
      },
      disabled: {
        type: Boolean,
        default: false
      },
      multiple: {
        type: Boolean,
        default: false
      },
      allowClear: {
        type: Boolean,
        default: false
      },
      autoSelect: {
        type: Boolean,
        default: true
      },
      options: {
        type: Array,
        default: () => []
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
        selected: this.multiple ? [] : ''
      }
    },
    computed: {
      searchable() {
        return this.options.length > 7
      }
    },
    watch: {
      value(value) {
        this.selected = this.options.find(item => item.name === value).value
      },
      selected(selected) {
        const selectValue = this.options.find(item => item.value === selected).name
        this.$emit('input', selectValue)
        this.$emit('on-selected', selectValue)
      }
    },
    created() {
      this.initValue()
    },
    methods: {
      initValue() {
        try {
          if (this.autoSelect && (!this.value || (this.multiple && !this.value.length))) {
            this.selected = this.multiple ? [this.options[0].value] : (this.options[0].value || '')
          } else {
            this.selected = this.value
          }
        } catch (error) {
          this.selected = this.multiple ? [] : ''
        }
      },
      focus() {
        this.$refs.selector.show()
      }
    }
  }
</script>

<style lang="scss" scoped>
    .form-list-selector {
        width: 100%;
    }
</style>
