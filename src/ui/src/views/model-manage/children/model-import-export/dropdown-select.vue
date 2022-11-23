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
  <bk-dropdown-menu
    align="center"
    :disabled="disabled"
    trigger="click"
    @click.native.stop
    class="dropdown-select"
  >
    <div class="dropdown-select-trigger" slot="dropdown-trigger">
      <span>{{ selectedValue.label || selectedValue.value }}</span>
      <bk-icon type="down-shape" />
    </div>
    <div class="bk-dropdown-list" slot="dropdown-content">
      <li :key="optionIndex" v-for="(option, optionIndex) in options">
        <dropdown-option-button
          :disabled="option.disabled"
          :active="option.value === selectedValue.value"
          @click="handleOptionClick(option)"
        >{{ option.label || option.value }}</dropdown-option-button
        >
      </li>
    </div>
  </bk-dropdown-menu>
</template>

<script>
  import { defineComponent, ref, watch } from 'vue'
  import DropdownOptionButton from '../dropdown-option-button.vue'

  export default defineComponent({
    name: 'DropdownSelect',
    components: { DropdownOptionButton },
    model: {
      prop: 'value',
      event: 'value-change'
    },
    props: {
      // 当前选中的值，支持 v-model
      value: {
        type: [String, Number, Boolean, Object],
        default: ''
      },
      // 选项列表
      options: {
        type: Array,
        default: () => []
      },
      // 是否禁用
      disabled: {
        type: Boolean,
        default: false
      }
    },
    setup({ value, options }, { emit }) {
      const selectedValue = ref({})
      const valueRef = ref(value)

      watch(
        valueRef,
        () => {
          selectedValue.value = options.find(opt => opt.value === valueRef.value) || {}
        },
        {
          immediate: true,
          deep: true
        }
      )

      const handleOptionClick = (option) => {
        selectedValue.value = option
        emit('value-change', option.value)
      }

      return {
        selectedValue,
        handleOptionClick
      }
    }
  })
</script>

<style lang="scss" scoped>
.dropdown-select {
  font-weight: normal;

  &-trigger {
    display: flex;
    align-items: center;
    background-color: #f0f1f5;
    line-height: 22px;
    padding: 0 8px;
    border-radius: 2px;
    user-select: none;

    .bk-icon {
      margin-left: auto;
    }
  }

  .dropdown-option-button {
    white-space: nowrap;
  }
}
</style>
