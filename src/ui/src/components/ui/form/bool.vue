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
  <bk-switcher
    class="form-bool"
    size="small"
    theme="primary"
    v-model="localValue"
    :disabled="disabled">
  </bk-switcher>
</template>

<script>
  export default {
    name: 'cmdb-form-bool',
    props: {
      value: {
        type: [String, Boolean], // String是为了兼容初始数据为空的情况
        default: false
      },
      disabled: {
        type: Boolean,
        default: false
      }
    },
    computed: {
      localValue: {
        get() {
          if (typeof this.value === 'boolean') {
            return this.value
          }
          if (this.value === 'true') {
            return true
          }
          return false
        },
        set(value) {
          this.$emit('input', value)
          this.$emit('change', value)
        }
      }
    }
  }
</script>

<style lang="scss" scoped>
    .form-bool{
        display: inline-block;
        vertical-align: middle;
    }
</style>
