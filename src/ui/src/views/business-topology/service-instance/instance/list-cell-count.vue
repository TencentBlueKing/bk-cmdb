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
  <span v-if="row.process_count">{{row.process_count}}</span>
  <span class="process-count-tips" v-else-if="row.service_template_id">
    <i class="tips-icon bk-icon icon-exclamation-circle"></i>
    <i18n class="tips-content" path="模板服务实例无进程提示">
      <template #link>
        <cmdb-auth class="tips-link"
          :auth="{ type: $OPERATION.U_SERVICE_INSTANCE, relation: [bizId] }"
          @click.native.stop
          @click="redirectToTemplate">
          {{$t('跳转添加并同步')}}
        </cmdb-auth>
      </template>
    </i18n>
  </span>
  <span class="process-count-tips" v-else>
    <i class="tips-icon bk-icon icon-exclamation-circle"></i>
    <i18n class="tips-content" path="普通服务实例无进程提示">
      <template #link>
        <cmdb-auth class="tips-link"
          :auth="{ type: $OPERATION.U_SERVICE_INSTANCE, relation: [bizId] }"
          @click.native.stop
          @click="handleAddProcess">
          {{$t('立即添加')}}
        </cmdb-auth>
      </template>
    </i18n>
  </span>
</template>

<script>
  import { mapGetters } from 'vuex'
  import createProcessMixin from './create-process-mixin'
  import { MENU_BUSINESS_SERVICE_TEMPLATE_DETAILS } from '@/dictionary/menu-symbol'

  export default {
    name: 'list-cell-count',
    mixins: [createProcessMixin],
    props: {
      row: Object
    },
    computed: {
      ...mapGetters('objectBiz', ['bizId'])
    },
    methods: {
      redirectToTemplate() {
        this.$routerActions.redirect({
          name: MENU_BUSINESS_SERVICE_TEMPLATE_DETAILS,
          params: {
            bizId: this.bizId,
            templateId: this.row.service_template_id
          },
          history: true
        })
      }
    }
  }
</script>

<style lang="scss" scoped>
    .process-count-tips {
        display: flex;
        align-items: center;
        .tips-icon {
            color: $warningColor;
            font-size: 14px;
        }
        .tips-content {
            padding: 0 4px;
            color: $textDisabledColor;
            .tips-link {
                color: $primaryColor;
                cursor: pointer;
                &.disabled {
                    color: $textDisabledColor;
                }
            }
        }
    }
</style>
