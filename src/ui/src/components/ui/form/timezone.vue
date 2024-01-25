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
  <bk-select class="form-timezone-selector"
    searchable
    v-bind="$attrs"
    v-model="selected"
    :clearable="false"
    :disabled="disabled"
    :multiple="multiple"
    :display-tag="multiple"
    :selected-style="getSelectedStyle"
    :placeholder="placeholder"
    ref="selector">
    <bk-option
      v-for="(option, index) in timezoneList"
      :key="index"
      :id="option.id"
      :name="option.name">
    </bk-option>
  </bk-select>
</template>

<script>
  import TIMEZONE from './timezone.json'
  export default {
    name: 'cmdb-form-timezone',
    props: {
      value: {
        type: [Array, String, Number],
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
      placeholder: {
        type: String,
        default: ''
      }
    },
    data() {
      const timezoneList = TIMEZONE.map(timezone => ({
        id: timezone,
        name: timezone
      }))
      return {
        timezoneList,
        selected: this.multiple ? [] : ''
      }
    },
    computed: {
      getSelectedStyle() {
        return this.multiple ? 'checkbox' : 'check'
      },
    },
    watch: {
      value(value) {
        if (value !== this.selected) {
          this.selected = value
        }
      },
      selected(selected) {
        this.$emit('input', selected)
        this.$emit('on-selected', selected)
      },
      disabled(disabled) {
        if (!disabled) {
          this.selected = this.getDefaultValue()
        }
      }
    },
    created() {
      this.selected = this.getDefaultValue()
    },
    methods: {
      getDefaultValue() {
        let value = this.value || ''
        if (this.multiple && !value.length) {
          value = ['Asia/Shanghai']
        } else {
          value = value || 'Asia/Shanghai'
        }
        return value
      },
      focus() {
        this.$refs.selector.show()
      }
    }
  }
</script>

<style lang="scss" scoped>
    .form-timezone-selector{
      width: 100%;
    }
</style>
