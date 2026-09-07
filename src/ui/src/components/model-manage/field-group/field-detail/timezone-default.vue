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
    class="field-timezone-default"
    searchable
    :clearable="!disabled"
    :disabled="disabled"
    :multiple="multiple"
    :display-tag="multiple"
    v-model="localValue">
    <bk-option
      v-for="timezone in timezoneList"
      :key="timezone"
      :id="timezone"
      :name="timezone">
    </bk-option>
  </bk-select>
</template>

<script>
  import TIMEZONE from '@/components/ui/form/timezone.json'

  export default {
    name: 'the-field-timezone-default',
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
      }
    },
    data() {
      return {
        timezoneList: TIMEZONE
      }
    },
    computed: {
      localValue: {
        get() {
          if (this.multiple) {
            return Array.isArray(this.value) ? this.value : []
          }
          return this.value || ''
        },
        set(val) {
          this.$emit('input', val)
        }
      }
    }
  }
</script>

<style lang="scss" scoped>
  .field-timezone-default {
    width: 100%;
  }
</style>
