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
          <host-search v-if="activeName === 'host'"></host-search>
          <full-text-search-bar v-else />
        </div>
        <full-text-search-result-tab
          v-if="fullTextSearchProps.showResultList"
          :result="searchResult" />
      </div>
      <div class="search-content">
        <div class="loading-ghost" v-show="loading" v-bkloading="{ isLoading: loading }"></div>
        <full-text-search-result-list v-show="!loading" v-if="fullTextSearchProps.showResultList"
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
  import hostSearch from './children/host-search'
  import fullTextSearchBar from './children/full-text-search/search-bar.vue'
  import fullTextSearchResultList from './children/full-text-search/result-list.vue'
  import fullTextSearchResultTab from './children/full-text-search/result-tab.vue'
  import theMap from './children/map'
  import theFooter from './children/footer'
  import { mapGetters } from 'vuex'
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
      ...mapGetters(['site']),
      isFullTextSearch() {
        return this.site.fullTextSearch === 'on'
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
    methods: {
      handleChangeTab(name) {
        this.activeName = name
      },
      setLayout() {
        const { query } = this.$route
        const showFullText = has(query, 'keyword')

        this.showResultList = showFullText
        // 暂只有全文检索需要显示查询结果
        this.fullTextSearchProps.showResultList = this.showResultList
        if (this.isFullTextSearch && this.fullTextSearchProps.showResultList) {
          this.activeName = 'fullText'
        } else {
          this.activeName = 'host'
        }
      },
      handleSearchComplete(result) {
        this.searchResult = result
      }
    }
  }
</script>

<style lang="scss" scoped>
  .index-layout {
    padding: 0 0 65px;
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
        background: #FAFBFD;
        top: 0;
      }
    }

    .search-tab {
      max-width: 726px;
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
