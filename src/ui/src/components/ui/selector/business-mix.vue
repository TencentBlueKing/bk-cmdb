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
  <bk-select style="text-align: left;"
    :loading="$loading([requestIds.biz, requestIds.bizset])"
    v-model="localSelected"
    :popover-width="320"
    :searchable="true"
    :search-with-pinyin="true"
    :clearable="false"
    :placeholder="$t('请选择业务')"
    :disabled="disabled"
    :popover-options="popoverOptions"
    font-size="normal"
    ext-popover-cls="cmdb-business-mix-selector-dropdown-content"
    @toggle="handleSelectToggle">
    <bk-option v-for="option in sortedList"
      :key="option.id"
      :id="option.id"
      :name="option.name"
      :disabled="!option.authorized">
      <auth-mask tag="div" :auth="option.auth" :authorized="option.authorized">
        <div
          :class="['option-item-content', { disabled: !option.authorized }]"
          :title="option.name">
          <div class="text">
            <span class="item-name">{{option.rawName}}</span>
            <span class="item-id">({{option.rawId}})</span>
          </div>
          <i class="icon icon-cc-business-set" v-if="option.isBizSet"></i>
          <i :class="['icon', 'bk-icon', 'collection', isCollected(option) ? 'icon-star-shape' : 'icon-star']"
            @click.prevent.stop="handleCollect(option)">
          </i>
        </div>
      </auth-mask>
    </bk-option>
    <div class="business-extension" slot="extension">
      <cmdb-auth :auth="{ type: $OPERATION.C_BUSINESS }" tag="div" class="extension-link"
        @click="handleCreateBusiness">
        <i class="bk-icon icon-plus-circle"></i>
        {{$t('新建业务')}}
      </cmdb-auth>
      <cmdb-auth :auth="{ type: $OPERATION.C_BUSINESS_SET }" tag="div" class="extension-link"
        @click="handleCreateBusinessSet">
        <i class="bk-icon icon-plus-circle"></i>
        {{$t('新建业务集')}}
      </cmdb-auth>
    </div>
  </bk-select>
</template>

<script>
  import { mapGetters } from 'vuex'
  import businessSetService from '@/service/business-set/index.js'
  import { verifyAuth } from '@/services/auth.js'
  import { TRANSFORM_TO_INTERNAL } from '@/dictionary/iam-auth'
  import AuthMask from '@/components/ui/auth/auth-mask.vue'
  import {
    BUSINESS_SELECTOR_COLLECTION,
    MENU_RESOURCE_BUSINESS,
    MENU_RESOURCE_BUSINESS_SET
  } from '@/dictionary/menu-symbol'

  const MAX_COLLECT_COUNT = 8

  export default {
    name: 'cmdb-business-mix-selector',
    components: {
      AuthMask
    },
    props: {
      value: {
        type: String
      },
      disabled: {
        type: Boolean,
        default: false
      },
      popoverOptions: {
        type: Object,
        default() {
          return {}
        }
      },
      showApplyCreate: Boolean
    },
    data() {
      return {
        useIAM: this.$Site.authscheme === 'iam',
        normalizationList: [],
        sortedList: [],
        requestIds: {
          biz: Symbol('biz'),
          bizset: Symbol('bizset'),
          collection: Symbol()
        }
      }
    },
    computed: {
      ...mapGetters('objectBiz', ['bizId']),
      ...mapGetters('userCustom', ['usercustom']),
      collection() {
        return this.usercustom[BUSINESS_SELECTOR_COLLECTION] || []
      },
      localSelected: {
        get() {
          return this.value
        },
        set(value) {
          const [id, type] = value.split('-')
          this.$emit('input', value)
          this.$emit('select', value, Number(id), type === 'bizset')
        }
      }
    },
    async created() {
      this.getData()
    },
    methods: {
      async getData() {
        const [{ info: bizList = [] }, { info: bizsetList = [] }] = await Promise.all([
          this.$http.get('biz/simplify', {
            requestId: this.requestIds.biz,
            fromCache: false
          }),
          businessSetService.getAll({
            requestId: this.requestIds.bizset,
            fromCache: false
          })
        ])

        const allList = [...bizList, ...bizsetList]
        const normalizationList = []
        const authList = []

        allList.forEach((item) => {
          const isBizSet = Boolean(item.bk_biz_set_id)
          const rawId = isBizSet ? item.bk_biz_set_id : item.bk_biz_id
          const rawName = isBizSet ? item.bk_biz_set_name : item.bk_biz_name
          normalizationList.push({
            isBizSet,
            rawId,
            rawName,
            authType: isBizSet ? 'bizSet' : 'business',
            authorized: true,
            // id值加后缀标明类型
            id: isBizSet ? `${rawId}-bizset` : `${rawId}-biz`,
            name: `${rawName} (${rawId})`
          })

          if (this.useIAM) {
            authList.push({
              type: isBizSet ? this.$OPERATION.R_BIZ_SET_RESOURCE : this.$OPERATION.R_BIZ_RESOURCE,
              relation: [rawId]
            })
          }
        })

        if (this.useIAM) {
          const authResult = await verifyAuth(TRANSFORM_TO_INTERNAL(authList)) || []
          authResult.forEach(({ resource_id: id, resource_type: type, is_pass: isPass }, index) => {
            const matched = normalizationList.find(item => item.authType === type && item.rawId === id)
            matched.authorized = isPass
            matched.auth = authList[index]
          })
        }

        this.normalizationList = normalizationList
        this.sortedList = this.normalizationList
      },
      isCollected(option) {
        return this.collection.includes(option.id)
      },
      sortList(list) {
        return list.slice().sort((a, b) => {
          const isACollected = this.isCollected(a)
          const isBCollected = this.isCollected(b)

          // 收藏的优先级最高
          if (isACollected > isBCollected) {
            return -1
          }

          // 同为收藏，先收藏的在前
          if (isACollected && isBCollected) {
            if (this.collection.indexOf(a.id) < this.collection.indexOf(b.id)) {
              return -1
            }
            return 1
          }

          if (!isACollected && !isBCollected) {
            // 有权限排前面
            if (a.authorized > b.authorized) {
              return -1
            }

            if (!a.authorized && !b.authorized) {
              // 同为业务，id正序
              if (!a.isBizSet && !b.isBizSet && a.rawId < b.rawId) {
                return -1
              }

              // 同为业务集，id正序
              if (a.isBizSet && b.isBizSet && a.rawId < b.rawId) {
                return -1
              }
            }
          }

          return 0
        })
      },
      async handleCollect(option) {
        if (this.$loading(this.requestIds.collection)) {
          return
        }

        let newCollection = []
        const isAdd = !this.collection.some(item => item === option.id)

        if (isAdd && this.collection.length >= MAX_COLLECT_COUNT) {
          this.$warn(this.$t('限制收藏个数提示', { max: MAX_COLLECT_COUNT }))
          return
        }

        if (isAdd) {
          newCollection = this.collection.concat(option.id)
        } else {
          newCollection = this.collection.filter(item => item !== option.id)
        }

        try {
          await this.$store.dispatch('userCustom/saveUsercustom', {
            [BUSINESS_SELECTOR_COLLECTION]: newCollection
          }, { requestId: this.requestIds.collection })
          this.$success(this.$t(isAdd ? '收藏成功' : '取消收藏成功'))
        } catch (err) {
          this.$error(this.$t(isAdd ? '收藏失败' : '取消收藏失败'))
        }
      },
      handleSelectToggle(isOpen) {
        // 每次下拉展开时重新排序数据，操作收藏时不排序防止顺序跳动
        if (isOpen) {
          this.sortedList = this.sortList(this.normalizationList)
        }
      },
      handleCreateBusiness() {
        this.$routerActions.redirect({
          name: MENU_RESOURCE_BUSINESS,
          query: {
            create: 1
          }
        })
      },
      handleCreateBusinessSet() {
        this.$routerActions.redirect({
          name: MENU_RESOURCE_BUSINESS_SET,
          query: {
            create: 1
          }
        })
      }
    }
  }
</script>

<style lang="scss" scoped>
  .option-item-content {
    color: #63656E;
    font-size: 12px;
    display: flex;
    align-items: center;
    justify-content: space-between;

    &.disabled {
      color: #c4c6cc;
    }

    .text {
      flex: 1;
      overflow: hidden;
      text-overflow: ellipsis;
    }
    .icon {
      margin-left: 8px;
      &.collection {
        padding: 2px;

        &:not(.icon-star-shape) {
          display: none;
        }
      }

      &.icon-star-shape {
        color: #FFB400;
      }

      &.icon-cc-business-set {
        color: #979ba5;
        font-size: 14px;
      }
    }

    .item-id {
      color: #C4C6CC;
    }

    &:hover {
      .icon {
        &.collection {
          display: block;
        }
      }
    }
  }

  .business-extension {
    display: flex;
    width: 100%;
    background-color: #FAFBFD;
    .extension-link {
      display: flex;
      flex: 1;
      align-items: center;
      justify-content: center;
      position: relative;
      font-size: 12px;
      color: #63656E;
      cursor: pointer;

      &:hover {
        opacity: .85;
      }

      .bk-icon {
        font-size: 16px;
        color: #979BA5;
        margin-right: 4px;
      }

      &.disabled {
        color: $textDisabledColor;
        .bk-icon {
          color: $textDisabledColor;
        }
      }

      & + .extension-link {
        &::before {
          position: absolute;
          content: "";
          left: 0;
          top: 10px;
          height: 12px;
          width: 1px;
          background: #c4c6cc;
        }
      }
    }
  }
</style>
<style lang="scss">
  .cmdb-business-mix-selector-dropdown-content {
    .bk-select-extension {
      padding: 0;
    }
  }
</style>
