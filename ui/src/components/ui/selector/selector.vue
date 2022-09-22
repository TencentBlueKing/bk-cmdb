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
  <bk-select v-if="hasChildren"
    ref="selector"
    v-model="selected"
    :placeholder="placeholder"
    :searchable="searchable"
    :clearable="allowClear"
    :disabled="disabled"
    :loading="loading"
    :font-size="fontSize"
    :popover-options="popoverOptions"
    :readonly="readonly"
    :size="size"
    @change="handleChange"
    @toggle="handleToggle">
    <bk-option-group v-for="(group, index) in list"
      :key="index"
      :name="group[displayKey]">
      <bk-option v-for="option in group.children || []"
        :key="option[settingKey]"
        :id="option[settingKey]"
        :name="option[displayKey]">
        <slot v-bind="option" />
      </bk-option>
    </bk-option-group>
  </bk-select>
  <bk-select v-else
    ref="selector"
    v-model="selected"
    :placeholder="placeholder"
    :searchable="searchable"
    :clearable="allowClear"
    :disabled="disabled"
    :loading="loading"
    :font-size="fontSize"
    :popover-options="popoverOptions"
    :readonly="readonly"
    :size="size"
    @change="handleChange"
    @toggle="handleToggle">
    <bk-option
      v-for="option in list"
      :key="option[settingKey]"
      :id="option[settingKey]"
      :name="option[displayKey]">
      <slot v-bind="option" />
    </bk-option>
  </bk-select>
</template>

<script>
  export default {
    name: 'cmdb-selector',
    props: {
      value: {
        type: [String, Number],
        default: ''
      },
      disabled: {
        type: Boolean,
        default: false
      },
      allowClear: {
        type: Boolean,
        default: false
      },
      list: {
        type: Array,
        default() {
          return []
        }
      },
      settingKey: {
        type: String,
        default: 'id'
      },
      displayKey: {
        type: String,
        default: 'name'
      },
      autoSelect: {
        type: Boolean,
        default: true
      },
      placeholder: {
        type: String,
        default: ''
      },
      hasChildren: {
        type: Boolean,
        default: false
      },
      emptyText: {
        type: String,
        default: ''
      },
      fontSize: {
        type: String,
        default: 'medium'
      },
      searchable: {
        type: Boolean,
        default: false
      },
      loading: Boolean,
      popoverOptions: {
        type: Object,
        default: () => ({})
      },
      readonly: Boolean,
      size: String
    },
    data() {
      return {
        selected: ''
      }
    },
    computed: {
      selectedOption() {
        return this.list.find(option => option[this.settingKey] === this.selected)
      }
    },
    watch: {
      value(value) {
        this.selected = value
      },
      selected(selected) {
        this.$emit('input', selected)
        this.$emit('on-selected', selected, this.selectedOption)
      },
      list() {
        this.setInitData()
      }
    },
    created() {
      this.setInitData()
    },
    methods: {
      setInitData() {
        let { value } = this
        if (this.autoSelect) {
          const currentOption = this.list.find(option => option[this.settingKey] === this.value)
          if (!currentOption) {
            value = this.list.length ? this.list[0][this.settingKey] : this.value
          }
        }
        this.selected = value
      },
      handleChange(newValue, oldValue) {
        this.$emit('change', newValue, oldValue)
      },
      handleToggle(isShow,) {
        this.$emit('toggle', isShow)
      },
      focus() {
        this.$refs.selector.show()
      }
    }
  }
</script>

<style lang="scss" scoped>
    .form-selector{
        display: inline-block;
        vertical-align: middle;
        width: 100%;
    }
</style>
