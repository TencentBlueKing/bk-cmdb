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
    multiple
    display-tag
    selected-style="checkbox"
    searchable
    v-model="localValue"
    v-bind="$attrs"
    @clear="() => $emit('clear')"
    @toggle="handleToggle">
    <bk-option v-for="timezone in timezones"
      :key="timezone"
      :id="timezone"
      :name="timezone">
    </bk-option>
  </bk-select>
</template>

<script>
  import activeMixin from './mixins/active'
  import TimeZones from '../ui/form/timezone.json'
  export default {
    name: 'cmdb-search-timezone',
    mixins: [activeMixin],
    props: {
      value: {
        type: Array,
        default: () => ([])
      }
    },
    data() {
      return {
        timezones: Object.freeze(TimeZones)
      }
    },
    computed: {
      localValue: {
        get() {
          return this.value || []
        },
        set(values) {
          this.$emit('input', values)
          this.$emit('change', values)
        }
      }
    }
  }
</script>
