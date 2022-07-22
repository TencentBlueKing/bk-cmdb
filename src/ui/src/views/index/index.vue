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
  <div class="index-layout" :style="{ '--defaultPaddingTop': `${paddingTop}px` }">
    <div :class="['search-layout', { sticky: showResultList }]">
      <div class="search-top">
        <div class="search-tab" v-if="isFullTextSearch">
          <span :class="['tab-item', { 'active': activeName === 'host' }]"
            @click="handleChangeTab('host')">
            {{$t('主机搜索')}}
          </span>
          <span :class="['tab-item', { 'active': activeName === 'fullText' }]"
            @click="handleChangeTab('fullText')">
            {{$t('全文检索')}}
          </span>
        </div>
        <div class="tab-content">
          <host-search v-if="activeName === 'host'" v-test-id></host-search>
          <full-text-search-bar v-else v-test-id />
        </div>
        <full-text-search-result-tab v-test-id
          v-if="fullTextSearchProps.showResultList"
          :result="searchResult" />
      </div>
      <div class="search-content">
        <div class="loading-ghost" v-show="loading" v-bkloading="{ isLoading: loading }"></div>
        <full-text-search-result-list v-show="!loading" v-if="fullTextSearchProps.showResultList" v-test-id
          @complete="handleSearchComplete"
          :fetching.sync="loading" />
      </div>
    </div>
    <the-map style="user-select: none;"></the-map>
    <the-footer></the-footer>
  </div>
</template>

<script>
  import has from 'has'
  import { MENU_RESOURCE_HOST } from '@/dictionary/menu-symbol'
  import { HOME_HOST_SEARCH_CONTENT_STORE_KEY } from '@/dictionary/storage-keys.js'
  import hostSearch from './children/host-search'
  import fullTextSearchBar from './children/full-text-search/search-bar.vue'
  import fullTextSearchResultList from './children/full-text-search/result-list.vue'
  import fullTextSearchResultTab from './children/full-text-search/result-tab.vue'
  import theMap from './children/map'
  import theFooter from './children/footer'
  export default {
    name: 'index',
    components: {
      hostSearch,
      fullTextSearchBar,
      fullTextSearchResultList,
      fullTextSearchResultTab,
      theMap,
      theFooter
    },
    data() {
      return {
        activeName: 'host',
        showResultList: false,
        fullTextSearchProps: {},
        searchResult: {},
        loading: false
      }
    },
    computed: {
      isFullTextSearch() {
        return this.$Site.fullTextSearch === 'on'
      },
      paddingTop() {
        return parseInt((this.$APP.height - 58) / 3, 10)
      }
    },
    watch: {
      '$route.query': {
        handler: 'setLayout',
        immediate: true
      }
    },
    beforeRouteEnter(to, from, next) {
      // 来源不是搜索落地页则清除搜索内容记录
      // eslint-disable-next-line no-underscore-dangle
      if (from?.name !== MENU_RESOURCE_HOST || !from?.query?._f) {
        window.sessionStorage.removeItem(HOME_HOST_SEARCH_CONTENT_STORE_KEY)
      }

      next()
    },
    methods: {
      handleChangeTab(name) {
        this.$router.push({
          query: {
            ...this.$route.query,
            tab: name
          }
        })
      },
      setLayout() {
        const { query } = this.$route
        this.activeName = query.tab || 'host'

        this.showResultList = has(query, 'keyword')
        // 暂只有全文检索需要显示查询结果
        this.fullTextSearchProps.showResultList = this.showResultList
      },
      handleSearchComplete(result) {
        this.searchResult = result
      }
    }
  }
</script>

<style lang="scss" scoped>
  .index-layout {
    padding: 0 0 52px;
    background-color: #F5F6FA;
    position: relative;
    z-index: 1;
  }

  .search-layout {
    height: 100%;
    transition: all 0.4s;
    padding-top: var(--defaultPaddingTop);
    overflow-y: auto;

    &.sticky {
      padding-top: 0;
      background: #FAFBFD;

      .search-top {
        padding-top: 50px;
        position: sticky;
        background: #FFF;
        border-bottom: 1px solid #dde4eb;
        top: 0;
      }
    }

    .search-tab {
      max-width: 806px;
      margin: 0 auto;
      font-size: 0;
      .tab-item {
        @include inlineBlock;
        position: relative;
        height: 30px;
        line-height: 30px;
        text-align: center;
        padding: 0 14px;
        margin: 0 4px -1px 0;
        font-size: 14px;
        color: #63656E;
        background-color: #DCDEE5;
        border: 1px solid #C4C6CC;
        border-radius: 6px 6px 0 0;
        transition: all 0.2s;
        cursor: pointer;
        &.active {
          background-color: #FFFFFF;
          border-bottom-color: #FFFFFF !important;
          z-index: 1000;
        }
      }
    }

    &:focus-within {
       .search-tab {
         .tab-item.active {
          border-color: #3A84FF;
        }
       }
    }

    .tab-content {
      display: block;
    }
  }

  .search-content {
    padding: 24px 0;
    .loading-ghost {
      min-height: 360px;
      /deep/ .bk-loading {
        background-color: #FAFBFD !important;
        z-index: 3 !important;
      }
    }
  }
</style>
