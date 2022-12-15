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
    v-model="localValue"
    v-bind="$attrs"
    :trigger="trigger"
    :list="list"
    @removeAll="() => $emit('clear')"
    @click.native="handleToggle(true)"
    @blur="handleToggle(false, ...arguments)">
  </bk-tag-input>
</template>

<script>
  import activeMixin from './mixins/active'
  export default {
    name: 'cmdb-search-map',
    mixins: [activeMixin],
    props: {
      value: {
        type: Array,
        default: () => ([])
      },
      options: {
        type: Array,
        default: () => ([])
      },
      idKey: {
        type: String,
        default: 'key'
      },
      nameKey: {
        type: String,
        default: 'val'
      },
      trigger: {
        type: String,
        default: 'focus'
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
      list() {
        return this.options.map(opt => ({
          id: `${opt[this.idKey]}=${opt[this.nameKey]}`,
          name: `${opt[this.idKey]}=${opt[this.nameKey]}`
        }))
      }
    },
    methods: {
    }
  }
</script>
