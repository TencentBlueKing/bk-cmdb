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
  <div :class="['data-empty', type]">
    <div class="content" v-if="type === 'search'">
      <bk-exception type="search-empty" scene="part">
        <p>{{ $t('搜索结果为空') }}</p>
        <div class="data-tips">
          <i18n class="operation-text" path="搜索为空提示语">
            <template #filter><span class="search-text">{{$t('调整关键词')}}</span></template>
            <template #clear>
              <bk-button class="text-btn" theme="primary" text style="margin-left: 3px" @click="$emit('clear')">
                {{$t('清空筛选条件')}}
              </bk-button>
            </template>
          </i18n>
        </div>
      </bk-exception>
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
      <div>
        <template v-if="$slots.default">
          <slot></slot>
        </template>
        <template v-else>
          <div class="content" v-if="type === 'default'">
            <bk-exception type="empty" scene="part">
              <div class="data-tips"></div>
              <i18n :path="path" tag="div" v-if="!emptyText">
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
                      {{action}}
                    </bk-button>
                  </cmdb-auth>
                </template>
                <template #empty-link>
                  <a class="empty-link" href="javascript:void(0)" @click="$emit('empty-link')">{{emptyLinkText}}</a>
                </template>
              </i18n>
            </bk-exception>
          </div>
          <div class="content" v-else>
            <slot name="permission">
              <bk-exception type="empty" scene="part">
                {{ defaultText }}
              </bk-exception>
            </slot>
          </div>
        </template>
      </div>
    </div>
  </div>
</template>

<script>
  import permissionMixins from '@/mixins/permission'
  export default {
    name: 'cmdb-data-empty',
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
        permission: this.stuff.payload?.permission || ''
      }
    },
    computed: {
      type() {
        return this.stuff?.type || ''
      },
      action() {
        return this.stuff.payload.action || this.$t('创建')
      },
      resource() {
        return this.stuff.payload?.resource || ''
      },
      emptyText() {
        return this.stuff.payload?.emptyText || ''
      },
      payload() {
        return this.stuff?.payload || ''
      },
      defaultText() {
        return this.stuff.payload?.defaultText || ''
      },
      path() {
        return this.stuff.payload?.path || ''
      },
      emptyLinkText() {
        return this.stuff.payload?.emptyLinkText || ''
      }
    },
    watch: {
      stuff: {
        handler(value) {
          this.permission = value.payload?.permission
        },
        deep: true
      }
    }
  }
</script>

<style lang="scss" scoped>
.data-empty {
  color: #63656e;
  font-size: 14px;
  .img-empty {
    width: 90px;
  }
  .text-btn {
    font-size: 14px;
    height: auto;
  }
  .search-text{
    margin: 0 3px;
  }
  .data-tips{
    margin-top: 15px;
  }
}
.empty-link {
  color: #3A84FF;
}
</style>
