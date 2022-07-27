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
  <div :class="['table-stuff', type]">
    <div class="content" v-if="type === 'search'">
      <i class="bk-cc-icon icon-cc-tips"></i>
      <span class="text">{{$t('搜索内容为空')}}</span>
    </div>
    <div class="content" v-else-if="type === 'permission'">
      <slot name="permission">
        <bk-exception type="403" scene="part">
          <i18n path="抱歉您没有查看权限">
            <template #link>
              <bk-button class="text-btn"
                text
                theme="primary"
                @click="handleApplyPermission">
                {{$t('去申请')}}
              </bk-button>
            </template>
          </i18n>
        </bk-exception>
      </slot>
    </div>
    <div class="content" v-else>
      <img class="img-empty" src="../../../assets/images/empty-content.png" alt="">
      <div>
        <template v-if="$slots.default">
          <slot></slot>
        </template>
        <template v-else>
          <i18n path="您还未XXX" tag="div" v-if="!emptyText">
            <template #action><span>{{action}}</span></template>
            <template #resource><span>{{resource}}</span></template>
            <template #link>
              <cmdb-auth :auth="auth">
                <bk-button class="text-btn"
                  text
                  theme="primary"
                  slot-scope="{ disabled }"
                  :disabled="disabled"
                  @click="$emit('create')">
                  {{$i18n.locale === 'en' ? `${action} now` : `立即${action}`}}
                </bk-button>
              </cmdb-auth>
            </template>
          </i18n>
          <span v-else>
            {{emptyText}}
          </span>
        </template>
      </div>
    </div>
  </div>
</template>

<script>
  import permissionMixins from '@/mixins/permission'
  export default {
    name: 'cmdb-table-empty',
    mixins: [permissionMixins],
    props: {
      stuff: {
        type: Object,
        default: () => ({
          type: 'default',
          payload: {}
        })
      },
      auth: {
        type: Object,
        default: null
      }
    },
    data() {
      return {
        permission: this.stuff.payload.permission
      }
    },
    computed: {
      type() {
        return this.stuff.type
      },
      action() {
        return this.stuff.payload.action || this.$t('创建')
      },
      resource() {
        return this.stuff.payload.resource
      },
      emptyText() {
        return this.stuff.payload.emptyText
      },
      payload() {
        return this.stuff.payload
      }
    },
    watch: {
      stuff: {
        handler(value) {
          this.permission = value.payload.permission
        },
        deep: true
      }
    },
    mounted() {
    },
    methods: {
    }
  }
</script>

<style lang="scss" scoped>
    .table-stuff {
        color: #63656e;
        font-size: 14px;
        .img-empty {
            width: 90px;
        }
        .text-btn {
            font-size: 14px;
            height: auto;
        }
    }
</style>
