<template>
  <div class="search-bar" v-click-outside="handleClickOutside">
    <div class="input-bar">
      <bk-input class="search-input"
        ref="searchInput"
        autocomplete="off"
        maxlength="32"
        clearable
        :placeholder="$t('请输入关键字，点击或回车搜索')"
        v-model.trim="keyword"
        @focus="handleFocus"
        @keydown="handleKeydown"
        @input="handleInput"
        @enter="handleSearch"
        @clear="handleClear">
      </bk-input>
      <bk-button theme="primary" class="search-btn"
        @click="handleSearch">
        <i class="bk-icon icon-search"></i>
        {{$t('搜索')}}
      </bk-button>
    </div>

    <div class="search-popover suggestion" v-show="showSuggestion">
      <ul class="list suggestion-list">
        <li v-for="(item, index) in suggestion" :key="index"
          :title="item.title"
          class="item"
          @click="item.linkTo(item.source)">
          <span class="name">{{item.title}}</span>
          <span class="type">({{item.typeName}})</span>
          <i class="tag-disabled" v-if="item.type === 'biz' && item.source.bk_data_status === 'disabled'">
            {{$t('已归档')}}
          </i>
        </li>
      </ul>
    </div>

    <div class="search-popover history" v-show="showHistory">
      <div class="history-title clearfix">
        <span class="fl">{{$t('搜索历史')}}</span>
        <bk-button :text="true" class="clear-btn fr" @click="handlClearHistory">
          <i class="bk-icon icon-cc-delete"></i>
          {{$t('清空')}}
        </bk-button>
      </div>
      <ul class="list history-list">
        <li v-for="(history, index) in historyList"
          ref="historyItem"
          :key="index"
          :class="['item', { 'selected': selectIndex === index }]"
          @click="handleClickHistory(history)">
          {{history}}
        </li>
      </ul>
    </div>
  </div>
</template>

<script>
  import { defineComponent, ref, onUnmounted, watch } from '@vue/composition-api'
  import useResult from './use-result.js'
  import useSuggestion from './use-suggestion'
  import useHistory from './use-history'
  import useRoute from './use-route.js'

  export default defineComponent({
    setup(props, { root }) {
      const { $store, $routerActions, $route } = root
      const { route } = useRoute(root)

      const keyword = ref('')
      watch(() => route.value.query, (query) => {
        if (query.keyword) {
          keyword.value = query.keyword
        }
      }, { immediate: true })

      const focusWithin = ref(false)
      const searchInput = ref(null)
      const forceHide = ref(false)
      let maxLengthPopover = null

      const { result, getSearchResult } = useResult({ route, keyword }, root)

      const handleFocus = () => {
        focusWithin.value = true
        forceHide.value = false
      }
      const handleClickOutside = () => {
        focusWithin.value = false
      }
      const handleKeydown = (value, event) => {
        onkeydown(event)
      }
      const handleClear = () => {
        forceHide.value = true
      }

      const handleSearch = () => {
        $store.commit('fullTextSearch/setSearchHistory', keyword.value)
        forceHide.value = true
        $routerActions.redirect({
          name: $route.name,
          query: {
            keyword: keyword.value,
            t: Date.now()
          }
        })
      }

      const handleInput = (value) => {
        if (value.length <= 32) {
          maxLengthPopover && maxLengthPopover.hide()
          getSearchResult()
          return
        }
        if (!maxLengthPopover) {
          maxLengthPopover = root.$bkPopover(searchInput.value.$el, {
            theme: 'dark search-input-max-length',
            content: root.$t('最大支持搜索32个字符'),
            zIndex: 9999,
            trigger: 'manual',
            boundary: 'window',
            arrow: true
          })
        }
        maxLengthPopover.show()
      }

      const historyState = {
        keyword,
        focusWithin,
        forceHide
      }
      const {
        historyList,
        showHistory,
        selectHistory,
        selectIndex,
        handleHistorySearch,
        handlClearHistory,
        onkeydown
      } = useHistory(historyState, root)

      const suggestionState = {
        result,
        focusWithin,
        showHistory,
        selectHistory,
        forceHide,
        keyword
      }
      const { suggestion, showSuggestion } = useSuggestion(suggestionState, root, result)

      const handleClickHistory = (history) => {
        handleHistorySearch(history)
        forceHide.value = true
      }

      onUnmounted(() => {
        maxLengthPopover && maxLengthPopover.destroy()
      })

      return {
        suggestion,
        showSuggestion,

        selectIndex,
        showHistory,
        historyList,
        handleKeydown,
        handleClickHistory,

        handleFocus,
        handleClickOutside,
        handleClear,

        keyword,
        result,
        searchInput,
        handleInput,
        handleSearch,
        handlClearHistory
      }
    }
  })
</script>

<style lang="scss" scoped>
  .search-bar {
    position: relative;
    max-width: 726px;
    margin: 0 auto;

    .input-bar {
      height: 42px;
      z-index: 999;
      display: flex;
      .search-input {
        flex: 1;
        /deep/ {
          .bk-input-text {
            border: 0;
            border-radius: 0 0 0 2px;
          }
          .bk-form-input {
            min-height: 42px;
            line-height: 30px;
            font-size: 14px;
            border: 1px solid #C4C6CC;
            padding: 5px 16px;
            border-radius: 0 0 0 2px;
          }
        }
      }
      .search-btn {
        width: 86px;
        height: 42px;
        line-height: 42px;
        padding: 0;
        border-radius: 0 2px 2px 0;
        .icon-search {
          width: 18px;
          height: 18px;
          font-size: 18px;
          margin: -2px 4px 0 0;
        }
      }
    }
  }

  .search-popover {
    position: absolute;
    top: 47px;
    left: 0;
    width: calc(100% - 86px);
    background-color: #ffffff;
    box-shadow: 0px 2px 6px 0px rgba(0,0,0,0.15);
    border: 1px solid #DCDEE5;
    overflow: hidden;
    z-index: 99;

    .list {
      .item {
        color: #63656E;
        font-size: 14px;
        padding: 0 20px;
        line-height: 40px;
        cursor: pointer;
        &:hover, &.selected {
          color: #3A84FF;
          background-color: #E1ECFF;
        }
      }
    }

    .suggestion-list .item {
        display: flex;
        align-items: center;
        .name {
          max-width: 86%;
          @include ellipsis;
        }
        .type {
          padding-left: 10px;
          color: #C4C6CC;
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
        }
    }

    .history-title {
      font-size: 14px;
      line-height: 36px;
      color: #C4C6CC;
      padding: 5px 20px 0;
      &::after {
        content: '';
        display: block;
        height: 1px;
        background-color: #F0F1F5;
      }
      .clear-btn {
        color: #C4C6CC;
        &:hover {
          color: #979BA5;
        }
        .icon-cc-delete {
          margin-top: -2px;
        }
      }
    }
    .history-list {
      margin-bottom: 5px;
    }
  }
</style>

<style lang="scss">
    .search-input-max-length-theme {
        font-size: 12px;
        padding: 6px 12px;
        left: 248px !important;
    }
</style>
