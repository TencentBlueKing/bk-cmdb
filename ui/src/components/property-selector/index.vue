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
    v-model="localValue"
    :searchable="searchable"
    font-size="medium"
    :clearable="false"
    v-bind="$attrs">
    <bk-option v-for="option in options"
      :key="option.bk_property_id"
      :id="option.bk_property_id"
      :name="option.bk_property_name">
    </bk-option>
  </bk-select>
</template>

<script>
  export default {
    name: 'cmdb-property-selector',
    props: {
      properties: {
        type: Array,
        default: () => ([])
      },
      value: {
        type: [String, Number],
        default: ''
      },
      searchable: {
        type: Boolean,
        default: true
      }
    },
    computed: {
      localValue: {
        get() {
          return this.value
        },
        set(value) {
          this.$emit('input', value)
          this.$emit('change', value)
        }
      },
      options() {
        return this.properties.filter(property => !!property.id)
      }
    }
  }
</script>
