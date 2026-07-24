<!--
 * Tencent is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2017 Tencent. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
-->

<template>
  <span class="user-value-wrapper" v-bk-overflow-tips>
    <bk-user-display-name ref="userDisplayName" :user-id="localValue"></bk-user-display-name>
  </span>
</template>

<script>
  export default {
    name: 'cmdb-user-value',
    props: {
      value: {
        type: [String, Array],
        default: ''
      }
    },
    data() {
      return {}
    },
    computed: {
      localValue() {
        return this.value || ''
      }
    },
    methods: {
      getCopyValue() {
        return this.$refs?.userDisplayName?.textContent
      }
    }
  }
</script>

<style lang="scss" scoped>
  .user-value-wrapper {
    display: inline-block;
    max-width: 100%;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
    vertical-align: middle;
  }
</style>
