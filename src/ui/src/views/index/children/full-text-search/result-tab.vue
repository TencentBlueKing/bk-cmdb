<template>
  <div class="result-tab">
    <div class="categories">
      <span class="category-item"
        :class="['category-item', { 'category-active': !currentCategory }]"
        @click="toggleCategory()">
        {{$t('全部结果')}}（{{total}}）
      </span>
      <span class="category-item"
        v-for="(category, index) in categories"
        :key="index"
        :class="['category-item', { 'category-active': category.id === currentCategory }]"
        @click="toggleCategory(category.id)">
        {{category.name}}（{{category.count}}）
      </span>
    </div>
  </div>
</template>

<script>
  import { computed, defineComponent, toRefs } from '@vue/composition-api'
  import useRoute from './use-route.js'

  export default defineComponent({
    props: {
      result: {
        type: Object,
        default: () => ({})
      }
    },
    setup(props, { root }) {
      const { $store, $route, $routerActions } = root
      const { route } = useRoute(root)

      const { result } = toRefs(props)
      const currentCategory = computed(() => route.value.query.c)

      const getModelById = $store.getters['objectModelClassify/getModelById']

      // 分类标签
      const categories = computed(() => {
        const aggregations = result.value.aggregations || []
        const categories =  aggregations.sort((prev, next) => next.count - prev.count).map(({ key, kind, count }) => {
          let { bk_obj_id: id, bk_obj_name: name } = getModelById(key) || {}
          if (kind === 'model') {
            id = 'model'
            name = root.$t('模型')
          }
          return { id, name, kind, count: count > 999 ? '999+' : count }
        })
        categories.sort(item => (item.kind === 'model' ? -1 : 0))
        return categories
      })

      const toggleCategory = (objId) => {
        $routerActions.redirect({
          name: $route.name,
          query: {
            ...route.value.query,
            c: objId
          }
        })
      }

      const total = computed(() => (result.value.total > 999 ? '999+' : result.value.total))

      return {
        currentCategory,
        categories,
        total,
        toggleCategory
      }
    }
  })
</script>

<style lang="scss" scoped>
  .result-tab {
    width: 90%;
    margin: 38px auto 0;
  }

  .categories {
    color: $cmdbTextColor;
    background-color: #FAFBFD;
    font-size: 14px;
    border-bottom: 1px solid #dde4eb;
    .category-item {
      display: inline-block;
      margin-right: 20px;
      margin-bottom: 12px;
      cursor: pointer;
      &.category-active {
        color: #3a84ff;
        font-weight: bold;
      }
      &:hover {
        color: #3a84ff;
      }
    }
  }
</style>
