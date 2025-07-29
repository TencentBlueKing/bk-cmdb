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
  <blueking-user-selector ref="userSelector" type="info"
    v-if="localValue.length"
    :api="api"
    :value="localValue"
    v-bk-overflow-tips>
  </blueking-user-selector>
  <span v-else>--</span>
</template>

<script>
  import BluekingUserSelector from '@blueking/user-selector'
  export default {
    components: {
      BluekingUserSelector
    },
    props: {
      value: {
        type: String,
        default: ''
      }
    },
    data() {
      return {
        api: window.ESB.userManage
      }
    },
    computed: {
      localValue: {
        get() {
          if (this.value) {
            return this.value.split(',')
          }
          return []
        }
      }
    },
    methods: {
      getCopyValue() {
        return this.$refs?.userSelector?.userInfo || '--'
      }
    }
  }
</script>

<style lang="scss" scoped>
    .bk-table{
        .user-selector{
            width: 100%;
            overflow: hidden;
            display: block;
            white-space: nowrap;
            text-overflow: ellipsis;
        }
    }

</style>
