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
    v-bind="$attrs"
    :popover-min-width="120"
    :clearable="false">
    <bk-option v-for="(name, id) in IPMap" :key="id" :id="id" :name="name"></bk-option>
  </bk-select>
</template>

<script>
  import { PROCESS_BIND_IP_ALL_MAP, PROCESS_BIND_IPV4_MAP, PROCESS_BIND_IPV6_MAP } from '@/dictionary/process-bind-ip.js'

  export default {
    props: {
      value: {
        type: String,
        default: ''
      },
      type: {
        type: String,
        default: ''
      }
    },
    computed: {
      IPMap() {
        if (this.type === 'v4') {
          return PROCESS_BIND_IPV4_MAP
        }
        if (this.type === 'v6') {
          return PROCESS_BIND_IPV6_MAP
        }
        return PROCESS_BIND_IP_ALL_MAP
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
