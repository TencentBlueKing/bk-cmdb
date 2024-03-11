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
  <div class="expand">
    <bk-select class="form-enummulti-selector"
      v-model="selected"
      :clearable="allowClear"
      :searchable="searchable"
      :disabled="disabled"
      display-tag
      :multiple="localMultiple"
      :placeholder="placeholder"
      :font-size="fontSize"
      :popover-options="{
        boundary: 'window'
      }"
      v-bind="$attrs"
      ref="selector">
      <bk-option
        v-for="(option, index) in options"
        :key="index"
        :id="option.id"
        :name="option.name">
      </bk-option>
    </bk-select>
  </div>

</template>

<script>
  import isEqual from 'lodash/isEqual'
  import { isEmptyPropertyValue } from '@/utils/tools'

  export default {
    name: 'cmdb-form-enummulti',
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
      multiple: {
        type: Boolean,
        default: true
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
        default() {
          return []
        }
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
        initValue: this.value
      }
    },
    computed: {
      searchable() {
        return this.options.length > 7
      },
      localMultiple() {
        if (Array.isArray(this.initValue) && this.initValue.length > 1 && !this.multiple) {
          return true
        }
        return this.multiple
      },
      selected: {
        get() {
          if (isEmptyPropertyValue(this.value)) {
            return this.getDefaultValue()
          }

          if (!this.localMultiple) {
            return Array.isArray(this.value) ? this.value[0] : this.value
          }

          // 多选时需要过滤掉不在选项列表中的值
          const vals = !Array.isArray(this.value) ? [this.value] : this.value
          return vals.filter(val => this.options?.some?.(option => option.id === val))
        },
        set(value) {
          this.$emit('input', value)
          this.$emit('on-selected', value)
          this.$emit('change', value)
        }
      }
    },
    watch: {
      value: {
        immediate: true,
        handler() {
          this.checkSelected()
        }
      }
    },
    methods: {
      isEmpty(value) {
        return !value?.length
      },
      getDefaultValue() {
        if (this.autoSelect) {
          const defaultOptions = (this.options || []).filter(option => option.is_default)
          const defaultValue = defaultOptions.map(option => option.id)
          return this.localMultiple ? defaultValue : defaultValue[0]
        }

        return this.localMultiple ? [] : ''
      },
      checkSelected() {
        const { selected } = this
        if (!isEqual(this.value, selected)) {
          this.selected = selected
        }
      },
      focus() {
        this.$refs.selector.show()
      }
    }
  }
</script>

<style lang="scss" scoped>
    .form-enummulti-selector {
        width: 100%;
    }
</style>
