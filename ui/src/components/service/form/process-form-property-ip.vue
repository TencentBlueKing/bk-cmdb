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
  <cmdb-input-select
    name="ip"
    :placeholder="$t('请选择或输入IP')"
    :options="IPList"
    v-bind="$attrs"
    v-model="localValue">
  </cmdb-input-select>
</template>

<script>
  export default {
    props: {
      value: {
        type: String,
        default: ''
      }
    },
    inject: ['form'],
    data() {
      return {
        IPList: []
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
      requestId() {
        return `getInstanceIpByHost_${this.form.hostId}`
      }
    },
    created() {
      this.getBindIPList()
    },
    beforeDestroy() {
      this.$http.cancel(this.requestId)
    },
    methods: {
      async getBindIPList() {
        try {
          const { options } = await this.$store.dispatch('serviceInstance/getInstanceIpByHost', {
            hostId: this.form.hostId,
            config: {
              requestId: this.requestId,
              fromCache: true
            }
          })
          this.IPList = options.map(ip => ({ id: ip, name: ip }))
        } catch (error) {
          this.IPList = []
          console.error(error)
        }
      }
    }
  }
</script>
