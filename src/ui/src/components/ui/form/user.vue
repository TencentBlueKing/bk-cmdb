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
  <div class="cmdb-form-user">
    <div class="prepend" v-if="$slots.prepend">
      <slot name="prepend" />
    </div>
    <bk-user-selector
      class="selector"
      v-model="localValue"
      :multiple="multiple"
      :placeholder="localPlaceholder"
      :tenant-id="tenantId"
      :data-error-handler="errorHandler"
      v-bind="selectorMoreProps">
    </bk-user-selector>
  </div>
</template>

<script>
  import BkUserSelector from '@blueking/bk-user-selector/vue2'
  import { mapGetters } from 'vuex'
  import { showLoginModal } from '@/utils/login-helper'
  import '@blueking/bk-user-selector/vue2/vue2.css'

  export default {
    name: 'cmdb-form-objuser',
    components: {
      BkUserSelector
    },
    props: {
      value: {
        type: String,
        default: ''
      },
      fastSelect: Boolean,
      multiple: {
        type: Boolean,
        default: true
      },
      placeholder: {
        type: String,
        default: ''
      }
    },
    computed: {
      ...mapGetters(['userName']),
      api() {
        return window.Site.userManageUrl
      },
      tenantId() {
        return window.Site.tenantId
      },
      localValue: {
        get() {
          if (this.multiple) {
            return this.value?.split(',') ?? []
          }
          return this.value
        },
        set(val) {
          this.$emit('input', val.toString())
          this.$emit('change', val.toString, this.value)
        }
      },
      selectorMoreProps() {
        const moreProps = { ...this.$attrs }

        if (this.api) {
          moreProps.apiBaseUrl = this.api
        } else {
          // TODO: 未配置userManageUrl，使用自定义api
        }

        if (this.fastSelect) {
          moreProps.currentUserId = this.userName
        }

        return moreProps
      },
      localPlaceholder() {
        return this.placeholder || this.$t('请输入用户')
      }
    },
    methods: {
      errorHandler(res) {
        if (res.code === 1306000) {
          showLoginModal()
        }
      }
    }
  }
</script>

<style lang="scss" scoped>
.cmdb-form-user {
  display: flex;

  .prepend {
    margin-right: -1px;
  }

  .selector {
    width: 100%;

    &[size="small"] {
      height: 26px !important;

      :deep(.user-selector-container:not(.focus)) {
        height: 26px !important;

        &.placeholder:after {
          line-height: 24px;
        }

        .user-selector-input {
          margin-top: 2px;
          height: 20px;
          line-height: 20px;
        }

        .user-selector-selected,
        .user-selector-overflow-tag {
          margin: 2px 0 2px 6px;
          line-height: 20px;
        }
      }
    }

    :deep(.user-selector-selected:has(.non-existent)) {
      background: #FFDEDE;
      .user-selector-selected-value {
        color: #EB3333;
      }
      .user-selector-selected-clear {
        color: #F15656
      }
    }
  }
}
</style>
