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
  <bk-sideslider v-transfer-dom
    :is-show.sync="isShow"
    :title="title"
    :width="540"
    @hidden="handleHidden">
    <component slot="content" ref="component"
      :is="component"
      :container="this"
      v-bind="componentProps">
    </component>
  </bk-sideslider>
</template>

<script>
  import AccountForm from './account-form.vue'
  import AccountDetails from './account-details.vue'
  export default {
    components: {
      [AccountForm.name]: AccountForm,
      [AccountDetails.name]: AccountDetails
    },
    data() {
      return {
        isShow: false,
        title: '',
        component: null,
        componentProps: {}
      }
    },
    methods: {
      show(options) {
        this.componentProps = options.props || {}
        if (options.type === 'form') {
          this.component = AccountForm.name
        } else if (options.type === 'details') {
          this.component = AccountDetails.name
        }
        this.title = options.title
        this.isShow = true
      },
      hide() {
        this.isShow = false
      },
      handleHidden() {
        this.component = null
        this.componentProps = {}
      }
    }
  }
</script>
