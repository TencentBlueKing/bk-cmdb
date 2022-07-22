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
  <bk-dialog v-model="isShow"
    theme="primary"
    :width="840"
    :mask-close="false"
    :show-footer="false"
    header-position="left"
    :title="$t('业务集预览')">
    <div class="content" v-bkloading="{ isLoading: $loading(requestId) }">
      <div class="content-head">
        <i18n path="共N个业务">
          <template #count><em class="count">{{total}}</em></template>
        </i18n>
        <bk-input class="search-input" clearable
          :value="keyword"
          right-icon="icon-search"
          :placeholder="$t('业务名')"
          @enter="handleSearch"
          @clear="handleSearch">
        </bk-input>
      </div>
      <div class="content-main">
        <ul class="business-list" v-if="total > 0">
          <li v-for="(item, index) in businessList" :key="index" class="business-item">
            <bk-link :title="item.bk_biz_name" @click="handleClickName(item)">{{item.bk_biz_name}}</bk-link>
          </li>
        </ul>
        <bk-exception :type="keyword ? 'search-empty' : 'empty'" scene="part" v-else></bk-exception>
      </div>
      <div class="content-foot">
        <bk-pagination small
          v-bind="pagination"
          :count="total"
          :current.sync="pagination.current"
          :limit.sync="pagination.limit"
          @limit-change="pagination.current = 1" />
      </div>
    </div>
  </bk-dialog>
</template>

<script>
  import { computed, defineComponent, reactive, ref, toRefs, watchEffect, watch } from '@vue/composition-api'
  import { MENU_BUSINESS } from '@/dictionary/menu-symbol'
  import businessSetService from '@/service/business-set/index.js'
  import routerActions from '@/router/actions'

  export default defineComponent({
    props: {
      show: {
        type: Boolean,
        default: false
      },
      mode: {
        type: String,
        required: true,
        defalut: 'before',
        validator: value => ['before', 'after'].indexOf(value) !== -1
      },
      payload: {
        type: Object,
        required: true,
        defalut: {}
      }
    },
    setup(props, { emit }) {
      const { show, mode, payload } = toRefs(props)

      const isShow = computed({
        get: () => show.value,
        set: value => emit('update:show', value)
      })
      const requestId = Symbol()

      const keyword = ref('')
      const pagination = reactive({
        current: 1,
        limit: 15,
        'limit-list': [15, 30, 60, 120]
      })

      const searcher = computed(() => {
        const actions = {
          before: businessSetService.previewOfBeforeCreate,
          after: businessSetService.previewOfAfterCreate
        }
        const params = {
          page: {
            start: pagination.limit * (pagination.current - 1),
            limit: pagination.limit
          }
        }

        // 业务名模糊搜索
        if (keyword.value) {
          params.filter = {
            condition: 'AND',
            rules: [{
              field: 'bk_biz_name',
              operator: 'contains',
              value: keyword.value,
            }]
          }
        }

        const { bk_scope: scope, bk_biz_set_id: bizSetId } = payload.value
        if (mode.value === 'before') {
          params.bk_scope = scope.filter
        } else if (mode.value === 'after') {
          params.bk_biz_set_id = bizSetId
        }

        const searchMethod = actions[mode.value]

        return config => searchMethod(params, config)
      })

      const businessList = ref([])
      const total = ref(0)

      // searcher更新时会隐式触发此方法
      const getList = async () => {
        const { list, count } = await searcher.value({ requestId })
        businessList.value = list
        total.value = count
      }

      watchEffect(async () => {
        // dialog组件显示状态再触发数据查询（if渲染有点问题）
        if (!isShow.value) return

        getList()
      })

      const handleSearch = (value) => {
        keyword.value = value
      }

      // 隐藏时重置值
      watch(isShow, (val) => {
        if (!val) {
          keyword.value = ''
          pagination.current = 1
        }
      })

      const handleClickName = (item) => {
        routerActions.open({
          name: MENU_BUSINESS,
          params: {
            bizId: item.bk_biz_id
          }
        })
      }

      return {
        isShow,
        requestId,
        keyword,
        total,
        pagination,
        businessList,
        handleSearch,
        handleClickName
      }
    }
  })
</script>

<style lang="scss" scoped>
  .content {
    .content-head {
      display: flex;
      justify-content: space-between;

      .search-input {
        width: 320px;
      }
      .count {
        font-weight: 700;
        font-style: normal;
        margin: 0 2px;
      }
    }

    .content-main {
      height: 220px;
      margin: 24px 0;
      @include scrollbar-y;
    }

    .content-foot {
      margin: 12px 0;
    }
  }

  .business-list {
    display: flex;
    flex-wrap: wrap;

    .business-item {
      margin: 0 12px 10px 0;
      display: flex;
      flex: none;
      align-items: center;
      height: 32px;
      background: #F5F7FA;
      width: calc(33.333% - 12px);
      padding-left: 4px;

      &:nth-of-type(3n + 3) {
        margin-right: 0;
        width: 33.333%
      }

      ::v-deep .bk-link {
        width: 100%;
        justify-content: flex-start;

        .bk-link-text {
          font-size: 12px;
          @include ellipsis;
        }
      }
    }
  }
</style>
