<template>
  <div class="result-list">
    <template v-if="!fetching && list.length">
      <div class="data-list">
        <component v-for="(item, index) in list" :key="index"
          :is="`item-${item.type}`"
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
  import { computed, defineComponent, reactive, ref, watch } from '@vue/composition-api'
  import NoSearchResults from '@/views/status/no-search-results.vue'
  import ItemBiz from './item-biz.vue'
  import ItemModel from './item-model.vue'
  import ItemInstance from './item-instance.vue'
  import ItemHost from './item-host.vue'
  import useResult from './use-result'
  import useItem from './use-item'
  import useRoute from './use-route.js'

  export default defineComponent({
    components: {
      NoSearchResults,
      [ItemBiz.name]: ItemBiz,
      [ItemModel.name]: ItemModel,
      [ItemInstance.name]: ItemInstance,
      [ItemHost.name]: ItemHost
    },
    props: {
    },
    setup(props, { root, emit }) {
      const { $store, $route, $routerActions } = root

      const { route } = useRoute(root)
      const { result, fetching, getSearchResult } = useResult({ route }, root)

      const pagination = reactive({
        limit: 10,
        current: 1,
        total: 0
      })

      // 依赖query参数启动与响应
      watch(() => route.value.query, (query) => {
        const { ps: limit = 10, p: page = 1 } = query
        pagination.limit = Number(limit)
        pagination.current = Number(page)
        getSearchResult()
      }, { immediate: true })

      // 结果列表
      const hitList = computed(() => result.value.hits || [])
      const { normalizationList: list } = useItem(hitList, root)

      // 统一查询对象属性
      const propertyMap = ref({})
      watch(result, async (result) => {
        emit('complete', result)

        const aggregations = result.aggregations || []
        const hits = result.hits || []
        const objIds = aggregations.filter(item => item.kind !== 'model').map(item => item.key)
        const modelObjIds = hits.filter(item => item.type === 'model').map(item => item.source.bk_obj_id)
        const mergedObjIds = [...objIds, ...modelObjIds]
        if (!mergedObjIds.length) {
          return
        }

        propertyMap.value = await $store.dispatch('objectModelProperty/batchSearchObjectAttribute', {
          params: {
            bk_obj_id: { $in: mergedObjIds },
            bk_supplier_account: $store.getters.supplierAccount
          }
        })

        pagination.total = result.total
      })

      watch(fetching, fetching => emit('update:fetching', fetching))

      const handleLimitChange = (limit) => {
        pagination.limit = limit
        $routerActions.redirect({
          name: $route.name,
          query: {
            ...route.value.query,
            ps: limit
          }
        })
      }
      const handlePageChange = (page) => {
        pagination.current = page
        $routerActions.redirect({
          name: $route.name,
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
    width: 90%;
    margin: 0 auto;
  }

  .data-list {
    padding-top: 14px;
    color: $cmdbTextColor;
    .result-item {
      width: 65%;
      padding-bottom: 35px;
      color: #63656e;
      /deep/ {
        em.hl {
          color: #3a84ff !important;
          font-style: normal !important;
          word-break: break-all;
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
