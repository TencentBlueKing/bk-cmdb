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
  <div class="host-search-layout">
    <div class="search-bar">
      <bk-input class="search-input" v-test-id
        ref="searchInput"
        type="textarea"
        :placeholder="$t('首页主机搜索提示语')"
        :rows="rows"
        :clearable="true"
        v-model="searchContent"
        @focus="handleFocus"
        @blur="handleBlur"
        @keydown="handleKeydown">
      </bk-input>
      <bk-popover v-bind="popoverProps" ref="popover">
        <bk-button theme="primary" class="search-btn" v-test-id="'search'"
          :loading="$loading(request.search)"
          @click="handleSearch()">
          <i class="bk-icon icon-search"></i>
          {{$t('搜索')}}
        </bk-button>
        <div class="picking-popover-content" slot="content">
          <i18n tag="p" path="检测到输入框中包含非标准IP格式字符串，请选择以XXX自动解析">
            <template #c1><span>&lt;{{$t('IP')}}&gt;</span></template>
            <template #c2><span>&lt;{{$t('固资编号')}}&gt;</span></template>
          </i18n>
          <div class="buttons">
            <bk-button theme="primary" size="small" outline v-test-id="'ipSearch'"
              @click="handleSearch('ip')">
              {{$t('IP')}}
            </bk-button>
            <bk-button theme="primary" size="small" outline v-test-id="'assetSearch'"
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
  import FilterUtils from '@/components/filters/utils.js'
  import { HOME_HOST_SEARCH_CONTENT_STORE_KEY } from '@/dictionary/storage-keys.js'

  export default {
    data() {
      const defaultSearchContent = () => {
        let content = ''
        try {
          content = JSON.parse(window.sessionStorage.getItem(HOME_HOST_SEARCH_CONTENT_STORE_KEY)) || ''
        } catch (e) {
          console.error(e)
          content = ''
        }
        return content
      }
      return {
        rows: 1,
        searchContent: defaultSearchContent(),
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
      searchContent: {
        handler() {
          this.$nextTick(this.setRows)
        },
        immediate: true
      }
    },
    mounted() {
      this.textareaDom = this.$refs.searchInput && this.$refs.searchInput.$refs.textarea
    },
    methods: {
      getSearchList() {
        // 使用切割IP的方法分割内容，方法在此处完全适用且能与高级搜索的IP分割保持一致
        return FilterUtils.splitIP(this.searchContent)
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

        // 保存本次搜索内容
        window.sessionStorage.setItem(HOME_HOST_SEARCH_CONTENT_STORE_KEY, JSON.stringify(this.searchContent))

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
          // console.log(IPList, IPWithCloudList, assetList, cloudIdSet, force)
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
        max-width: 806px;
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
        max-width: 646px;
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
