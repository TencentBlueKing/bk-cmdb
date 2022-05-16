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
  <component class="auth-box"
    :is="tag"
    v-cursor="{
      active: !isAuthorized,
      auth: auth,
      authResults,
      onclick
    }"
    :class="{ disabled }"
    @click="handleClick">
    <slot :disabled="disabled"></slot>
  </component>
</template>

<script>
  import AuthProxy from './auth-queue'
  import deepEqual from 'deep-equal'
  export default {
    name: 'cmdb-auth',
    props: {
      ignore: Boolean,
      auth: {
        type: [Object, Array]
      },
      tag: {
        type: String,
        default: 'span'
      },
      onclick: Function
    },
    data() {
      return {
        authResults: null,
        authMetas: null,
        isAuthorized: false,
        disabled: true,
        useIAM: this.$Site.authscheme === 'iam'
      }
    },
    watch: {
      auth: {
        deep: true,
        handler(value, oldValue) {
          !deepEqual(value, oldValue) && this.setAuthProxy()
        }
      },
      ignore() {
        this.setAuthProxy()
      }
    },
    mounted() {
      this.setAuthProxy()
    },
    methods: {
      setAuthProxy() {
        if (this.useIAM && this.auth && !this.ignore) {
          AuthProxy.add({
            component: this,
            data: this.auth
          })
        } else {
          this.disabled = false
          this.isAuthorized = true
          this.$emit('update-auth', true)
        }
      },
      updateAuth(authResults, authMetas) {
        let isPass
        if (!authResults.length && authMetas.length) { // 鉴权失败
          isPass = false
        } else {
          isPass = authResults.every(result => result.is_pass)
        }
        this.authResults = authResults
        this.authMetas = authMetas
        this.isAuthorized = isPass
        this.disabled = !isPass
        this.$emit('update-auth', isPass)
      },
      handleClick() {
        if (this.disabled) {
          return
        }
        this.$emit('click')
      }
    }
  }
</script>

<style lang="scss" scoped>
    .auth-box {
        display: inline-block;
    }
</style>
