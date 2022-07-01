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
  <span v-if="isLimited"></span>
  <bk-popover class="process-form-tips" ext-cls="process-lock-popover" v-else>
    <i class="icon-cc-lock-fill"></i>
    <template slot="content">
      <i18n path="进程表单锁定提示">
        <template #link>
          <bk-link theme="primary" @click="handleRedirect" class="goto-link">{{$t('跳转服务模板')}}</bk-link>
        </template>
      </i18n>
    </template>
  </bk-popover>
</template>

<script>
  import Tippy from 'bk-magic-vue/lib/utils/tippy'
  import { MENU_BUSINESS_SERVICE_TEMPLATE_DETAILS } from '@/dictionary/menu-symbol'

  export default {
    props: {
      serviceTemplateId: Number,
      bizId: Number,
      property: {
        type: Object,
        default: () => ({})
      }
    },
    computed: {
      isLimited() {
        return ['bk_func_name', 'bk_process_name'].includes(this.property.bk_property_id)
      }
    },
    mounted() {
      if (this.isLimited) {
        this.setupTips()
      } else {
        this.hackRadius()
      }
    },
    methods: {
      setupTips() {
        const DOM = this.$el.previousElementSibling
        // eslint-disable-next-line new-cap
        Tippy(DOM, {
          theme: 'dark process-uneditable-tips',
          content: this.$t('系统限定不可修改'),
          arrow: true,
          placement: 'top'
        })
      },
      hackRadius() {
        const hackDOM = this.$el.parentElement.querySelectorAll('.bk-form-input,.bk-form-textarea,.bk-textarea-wrapper,.bk-select')
        Array.prototype.forEach.call(hackDOM, (dom) => {
          dom.style.borderTopRightRadius = 0
          dom.style.borderBottomRightRadius = 0
        })
      },
      handleRedirect() {
        this.$routerActions.open({
          name: MENU_BUSINESS_SERVICE_TEMPLATE_DETAILS,
          params: {
            bizId: this.bizId,
            templateId: this.serviceTemplateId
          }
        })
      }
    }
  }
</script>

<style lang="scss" scoped>
    .process-form-tips {
        width: 24px;
        display: inline-flex;
        align-items: center;
        justify-content: center;
        border: 1px solid #dcdee5;
        border-left: none;
        background-color: #fafbfd;
        font-size: 14px;
        overflow: hidden;
        cursor: pointer;
        /deep/ .bk-tooltip-ref {
            height: 100%;
            display: flex;
            align-items: center;
            justify-content: center;
        }
    }
</style>
<style lang="scss">
  .process-lock-popover {
    .goto-link {
      .bk-link-text {
        font-size: 12px;
      }
    }
  }
  .process-uneditable-tips-theme {
    font-size: 12px;
  }
</style>
