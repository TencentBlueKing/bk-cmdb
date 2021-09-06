<template>
  <div class="host-search-layout">
    <div class="search-bar">
      <bk-input class="search-input"
        ref="searchInput"
        type="textarea"
        :placeholder="$t('首页主机搜索提示语')"
        :rows="rows"
        v-model="searchContent"
        @focus="handleFocus"
        @blur="handleBlur"
        @keydown="handleKeydown">
      </bk-input>
      <bk-popover v-bind="popoverProps" ref="popover">
        <bk-button theme="primary" class="search-btn"
          :loading="$loading(request.search)"
          @click="handleSearch()">
          <i class="bk-icon icon-search"></i>
          {{$t('搜索')}}
        </bk-button>
        <div class="picking-popover-content" slot="content">
          <i18n tag="p" path="检测到输入框中包含非标准IP格式字符串，请选择以XXX自动解析">
            <span place="c1">&lt;{{$t('IP')}}&gt;</span>
            <span place="c2">&lt;{{$t('固资编号')}}&gt;</span>
          </i18n>
          <div class="buttons">
            <bk-button theme="primary" size="small" outline
              @click="handleSearch('ip')">
              {{$t('IP')}}
            </bk-button>
            <bk-button theme="primary" size="small" outline
              @click="handleSearch('asset')">
              {{$t('固资编号')}}
            </bk-button>
          </div>
        </div>
      </bk-popover>
      <bk-link theme="primary" class="advanced-link" @click="handleClickAdvancedSearch">{{$t('高级筛选')}}</bk-link>
    </div>
  </div>
</template>

<script>
  import { MENU_RESOURCE_HOST } from '@/dictionary/menu-symbol'
  import QS from 'qs'
  import isIP from 'validator/es/lib/isIP'
  import isInt from 'validator/es/lib/isInt'
  export default {
    data() {
      return {
        rows: 1,
        searchText: '',
        searchContent: '',
        textarea: '',
        textareaDom: null,
        popoverProps: {
          width: 280,
          trigger: 'manual',
          distance: 12,
          theme: 'light',
          placement: 'bottom',
          tippyOptions: {
            hideOnClick: true
          }
        },
        request: {
          search: Symbol('search')
        }
      }
    },
    watch: {
      searchContent() {
        this.$nextTick(this.setRows)
      }
    },
    mounted() {
      this.textareaDom = this.$refs.searchInput && this.$refs.searchInput.$refs.textarea
    },
    methods: {
      getSearchList() {
        const searchList = []
        this.searchContent.split('\n').forEach((text) => {
          const trimText = text.trim()
          if (trimText.length) {
            searchList.push(trimText)
          }
        })
        return searchList
      },
      setRows() {
        const rows = this.searchContent.split('\n').length || 1
        this.rows = Math.min(10, rows)
      },
      handleFocus() {
        this.$emit('focus', true)
        this.setRows()
      },
      handleBlur() {
        if (!this.searchContent.trim().length) {
          this.searchContent = ''
        }
        this.textareaDom && this.textareaDom.blur()
        this.$emit('focus', false)
      },
      handleKeydown(content, event) {
        const agent = window.navigator.userAgent.toLowerCase()
        const isMac = /macintosh|mac os x/i.test(agent)
        const modifierKey = isMac ? event.metaKey : event.ctrlKey
        if (modifierKey && event.code.toLowerCase() === 'enter') {
          this.handleSearch()
        }
      },
      async handleSearch(force = '') {
        const searchList = this.getSearchList()
        if (searchList.length > 10000) {
          this.$warn(this.$t('最多支持搜索10000条数据'))
          return
        }

        if (searchList.length) {
          const IPList = []
          const IPWithCloudList = []
          const assetList = []
          const cloudIdSet = new Set()
          searchList.forEach((text) => {
            if (isIP(text, 4)) {
              IPList.push(text)
            } else {
              const splitData = text.split(':')
              const [cloudId, ip] = splitData
              if (splitData.length === 2 && isInt(cloudId) && isIP(ip)) {
                IPWithCloudList.push(text)
                cloudIdSet.add(parseInt(cloudId, 10))
              } else {
                assetList.push(text)
              }
            }
          })
          console.log(IPList, IPWithCloudList, assetList, cloudIdSet, force)
          // 判断是否存在IP、固资编号混合搜索
          if (!force && (IPList.length || IPWithCloudList.length) && assetList.length) {
            this.$refs.popover.showHandler()
            return
          }

          const assetSearch = () => this.handleAssetSearch(assetList)

          const ipSearch = () => {
            // 无云区域与有云区域的混合搜索
            if (IPList.length && IPWithCloudList.length) {
              return this.$warn(this.$t('暂不支持不同云区域的混合搜索'))
            }
            // 纯IP搜索
            if (IPList.length) {
              return this.handleIPSearch(IPList)
            }
            // 不同云区域+IP的混合搜索
            if (cloudIdSet.size > 1) {
              return this.$warn(this.$t('暂不支持不同云区域的混合搜索'))
            }
            this.handleIPWithCloudSearch(IPWithCloudList, cloudIdSet)
          }

          // 优先使用混合搜索下的选择
          if (force === 'asset') {
            return assetSearch()
          }
          if (force === 'ip') {
            return ipSearch()
          }

          // 纯固资编号搜索
          if (assetList.length) {
            return assetSearch()
          }
          // IP系列搜索
          ipSearch()
        } else {
          this.searchContent = ''
          this.textareaDom && this.textareaDom.focus()
        }
      },
      handleIPSearch(list) {
        const ip = {
          text: list.join('\n'),
          inner: true,
          outer: true,
          exact: true
        }
        this.$routerActions.redirect({
          name: MENU_RESOURCE_HOST,
          query: {
            scope: 'all',
            ip: QS.stringify(ip, { encode: false })
          },
          history: true
        })
      },
      handleIPWithCloudSearch(list, cloudSet) {
        const IPList = list.map((text) => {
          const [, ip] = text.split(':')
          return ip
        })
        const ip = {
          text: IPList.join('\n'),
          inner: true,
          outer: true,
          exact: true
        }
        const filter = {
          'bk_cloud_id.in': [cloudSet.values().next().value].join(',')
        }
        this.$routerActions.redirect({
          name: MENU_RESOURCE_HOST,
          query: {
            scope: 'all',
            ip: QS.stringify(ip, { encode: false }),
            filter: QS.stringify(filter, { encode: false })
          },
          history: true
        })
      },
      async handleAssetSearch(list) {
        try {
          const filter = {
            'bk_asset_id.in': list.join(',')
          }
          this.$routerActions.redirect({
            name: MENU_RESOURCE_HOST,
            query: {
              scope: 'all',
              filter: QS.stringify(filter, { encode: false })
            },
            history: true
          })
        } catch (error) {
          console.error(true)
        }
      },
      handleClickAdvancedSearch() {
        this.$routerActions.redirect({
          name: MENU_RESOURCE_HOST,
          query: {
            adv: 1,
            scope: 'all'
          },
          history: false
        })
      }
    }
  }
</script>

<style lang="scss" scoped>
    .host-search-layout {
        position: relative;
        width: 100%;
        max-width: 726px;
        height: 42px;
        margin: 0 auto;
    }
    .search-bar {
        position: absolute;
        width: 100%;
        height: 42px;
        z-index: 999;
        display: flex;
    }
    .search-input {
        flex: 1;
        /deep/ {
            .bk-textarea-wrapper {
                border: 0;
                border-radius: 0 0 0 2px;
            }
            .bk-form-textarea {
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
    .search-text {
        position: absolute;
        left: 0;
        top: 0;
        width: 100%;
        max-width: 640px;
        height: 42px;
        line-height: 30px;
        font-size: 14px;
        color: #63656E;
        border: 1px solid #C4C6CC;
        background-color: #FFFFFF;
        padding: 5px 16px;
        z-index: 1;
        cursor: text;
        @include ellipsis;
    }
    .advanced-link {
      margin-left: 8px;
      /deep/ .bk-link-text {
        font-size: 12px;
      }
    }
    .picking-popover-content {
      padding: 6px;
      .buttons {
        margin-top: 12px;
        text-align: right;
      }
    }
</style>
