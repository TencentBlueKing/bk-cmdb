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
  <bk-select
    searchable
    v-model="localValue"
    v-bind="$attrs"
    :multiple="multiple"
    @clear="() => $emit('clear')"
    @toggle="handleToggle">
    <bk-option v-for="option in options"
      :key="option.id"
      :id="option.id"
      :name="option.name">
    </bk-option>
  </bk-select>
</template>

<script>
  import activeMixin from './mixins/active'
  export default {
    name: 'cmdb-search-enum',
    mixins: [activeMixin],
    props: {
      value: {
        type: [String, Array],
        default: () => ([])
      },
      options: {
        type: Array,
        default: () => ([])
      }
    },
    computed: {
      multiple() {
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
    }
  }
</script>
