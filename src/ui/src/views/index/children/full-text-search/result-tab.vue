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
  <div class="result-tab">
    <div :class="['categories', { expanded }]"
      :style="{
        '--rows': rows,
        '--itemHeight': `${itemHeight}px`,
        '--itemMarginRight': `${itemMarginRight}px`,
        '--itemMarginBottom': `${itemMarginBottom}px`
      }">
      <span :class="['category-item', { 'category-active': !currentCategory }]"
        @click="handleSelectCategory()">
        {{$t('全部结果')}}（{{total}}）
      </span>
      <span v-for="(category, index) in categories"
        :key="index"
        :class="['category-item', { 'category-active': category.id === currentCategory }]"
        @click="handleSelectCategory(category)">
        {{category.name}}（{{category.count}}）
      </span>
      <div class="toggle-anchor" v-show="showMore" @click="handleToggle">
        <div class="anchor-button" v-show="!expanded">
          <span>+{{categories.length}}</span>
          <bk-icon type="angle-double-down" />
        </div>
        <bk-link v-show="expanded" theme="primary">收起</bk-link>
      </div>
    </div>
  </div>
</template>

<script>
  import { computed, defineComponent, onMounted, onUpdated, toRefs } from 'vue'
  import routerActions from '@/router/actions'
  import RouterQuery from '@/router/query'
  import useTab, { sizes } from './use-tab.js'

  export default defineComponent({
    props: {
      result: {
        type: Object,
        default: () => ({})
      }
    },
    setup(props) {
      const route = computed(() => RouterQuery.route)

      const { result } = toRefs(props)
      const currentCategory = computed(() => route.value.query.c)

      // 分类标签
      const aggregations = computed(() => result.value.aggregations || [])
      const { categories, calculateSizes } = useTab(aggregations)

      const total = computed(() => (result.value.total > 999 ? '999+' : result.value.total))

      const handleSelectCategory = (category) => {
        routerActions.redirect({
          name: route.value.name,
          query: {
            ...route.value.query,
            c: category?.id,
            k: category?.kind,
            p: undefined,
            ps: undefined
          }
        })
      }

      onMounted(() => {
        calculateSizes()
      })

      onUpdated(() => {
        calculateSizes()
      })

      const handleToggle = () => {
        sizes.expanded = !sizes.expanded
      }

      return {
        currentCategory,
        categories,
        total,
        handleSelectCategory,
        handleToggle,
        ...toRefs(sizes)
      }
    }
  })
</script>

<style lang="scss" scoped>
  .result-tab {
    width: 1280px;
    margin: 38px auto 0;
  }

  .categories {
    display: flex;
    flex-wrap: wrap;
    height: var(--itemHeight);
    overflow: hidden;
    color: $cmdbTextColor;
    background-color: #FFF;
    font-size: 14px;
    transition: all .125s ease-out;

    .category-item {
      height: 20px;
      margin-right: var(--itemMarginRight);
      margin-bottom: var(--itemMarginBottom);
      cursor: pointer;
      &.category-active {
        color: #3a84ff;
        font-weight: bold;
      }
      &:hover {
        color: #3a84ff;
      }
    }

    .anchor-button {
      height: 20px;
      background: #f0f1f5;
      border-radius: 8px;
      text-align: center;
      min-width: 66px;
      cursor: pointer;
    }

    &.expanded {
      height: calc(var(--rows) * var(--itemHeight));
    }
  }
</style>
