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
  <div class="result-list">
    <template v-if="!fetching && list.length">
      <div class="data-list">
        <component v-for="(item, index) in list" :key="index"
          :is="`item-${item.comp || item.type}`"
          :property-map="propertyMap"
          :data="item" />
      </div>

      <div class="pagination">
        <span class="mr10">{{$tc('共计N条', pagination.total, { N: pagination.total })}}</span>
        <bk-pagination
          size="small"
          align="right"
          :type="'compact'"
          :current.sync="pagination.current"
          :limit="pagination.limit"
          :count="pagination.total"
          @limit-change="handleLimitChange"
          @change="handlePageChange">
        </bk-pagination>
      </div>
    </template>
    <no-search-results v-else-if="fetching !== -1" :text="$t('搜不到相关内容')" />
  </div>
</template>

<script>
  import { computed, defineComponent, reactive, ref, watch } from 'vue'
  import store from '@/store'
  import routerActions from '@/router/actions'
  import RouterQuery from '@/router/query'
  import NoSearchResults from '@/views/status/no-search-results.vue'
  import ItemBiz from './item-biz.vue'
  import ItemBizSet from './item-bizset.vue'
  import ItemModel from './item-model.vue'
  import ItemInstance from './item-instance.vue'
  import ItemHost from './item-host.vue'
  import ItemSet from './item-set.vue'
  import ItemModule from './item-module.vue'
  import useResult from './use-result'
  import useItem from './use-item'
  import { categories } from './use-tab.js'

  export default defineComponent({
    components: {
      NoSearchResults,
      [ItemBiz.name]: ItemBiz,
      [ItemBizSet.name]: ItemBizSet,
      [ItemBiz.name]: ItemBiz,
      [ItemModel.name]: ItemModel,
      [ItemInstance.name]: ItemInstance,
      [ItemHost.name]: ItemHost,
      [ItemSet.name]: ItemSet,
      [ItemModule.name]: ItemModule,
    },
    setup(props, { emit }) {
      const route = computed(() => RouterQuery.route)
      const { result, fetching, getSearchResult } = useResult({ route })

      const pagination = reactive({
        limit: 10,
        current: 1,
        total: 0
      })

      // 依赖query参数启动与响应
      watch(() => route.value.query, (query) => {
        const { ps: limit = 10, p: page = 1, tab } = query
        pagination.limit = Number(limit)
        pagination.current = Number(page)
        if (tab === 'fullText') {
          getSearchResult()
        }
      }, { immediate: true })

      // 结果列表
      const hitList = computed(() => result.value.hits || [])
      const { normalizationList: list } = useItem(hitList)

      // 根据当前分类设置分页总数
      watch(categories, (categories) => {
        const { c: objId, k: kind } = route.value.query
        const current = categories.find((item) => {
          if (kind === 'model') {
            return item.kind === kind
          }
          return item.kind === kind && item.id === objId
        })

        if (current) {
          pagination.total = current.total
        } else {
          pagination.total = result.value.total
        }
      })

      // 统一查询对象属性
      const propertyMap = ref({})
      watch(result, async (result) => {
        emit('complete', result)

        const hits = result.hits || []
        const modelIds = hits.map(item => item.key)
        if (!modelIds.length) {
          return
        }

        propertyMap.value = await store.dispatch('objectModelProperty/batchSearchObjectAttribute', {
          params: {
            bk_obj_id: { $in: [...new Set(modelIds)] },
            bk_supplier_account: store.getters.supplierAccount
          }
        })
      })

      watch(fetching, fetching => emit('update:fetching', fetching))

      const handleLimitChange = (limit) => {
        pagination.limit = limit
        routerActions.redirect({
          name: route.value.name,
          query: {
            ...route.value.query,
            ps: limit
          }
        })
      }
      const handlePageChange = (page) => {
        pagination.current = page
        routerActions.redirect({
          name: route.value.name,
          query: {
            ...route.value.query,
            p: page
          }
        })
      }

      return {
        list,
        pagination,
        fetching,
        propertyMap,
        handleLimitChange,
        handlePageChange
      }
    }
  })
</script>

<style lang="scss" scoped>
  .result-list {
    width: 1280px;
    margin: 0 auto;
  }

  .data-list {
    padding-top: 14px;
    color: $cmdbTextColor;
    .result-item {
      padding-bottom: 35px;
      color: #63656e;
      /deep/ {
        .hl {
          em {
            color: #3a84ff !important;
            font-style: normal !important;
            word-break: break-all;
          }
        }
        .result-title {
          display: inline-block;
          font-size: 18px;
          font-weight: bold;
          margin-bottom: 4px;
          cursor: pointer;
          &:hover {
            span {
              color: #3a84ff;
              text-decoration: underline;
            }
          }
          .tag-disabled {
            height: 18px;
            line-height: 16px;
            padding: 0 4px;
            font-style: normal;
            font-size: 12px;
            color: #979BA5;
            border: 1px solid #C4C6CC;
            background-color: #FAFBFD;
            border-radius: 2px;
            margin-left: 4px;
            text-decoration: none;
          }
        }
        .result-desc {
          display: flex;
          flex-wrap: wrap;
          font-size: 14px;
          .desc-item {
            flex: none;
            max-width: 100%;
            word-wrap: break-word;
            word-break: break-all;
            margin-bottom: 6px;
            margin-right: 16px;
          }
          &:hover {
            color: #313238;
            cursor: pointer;
          }
        }
      }
    }
  }

  .pagination {
    display: flex;
    align-items: center;
    font-size: 12px;
    color: #737987;
    .bk-page {
      flex: 1;
    }
  }
</style>
