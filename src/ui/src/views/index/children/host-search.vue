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
      <div class="search-input-close">
        <div class="search-input"
          v-test-id
          ref="searchInput"
          :placeholder="$t('首页主机搜索提示语')"
          :focusTip="$t('首页主机搜索聚焦提示')"
          :blurTip="$t('首页主机搜索失焦焦提示')"
          contenteditable="plaintext-only"
          @blur="handleBlur"
          @keydown="handleKeydown"
          @focus="handleFocus"
          @input="handleInput"
          @paste="handlePaste">
        </div>
        <i class="search-close bk-icon icon-close-circle-shape" @mousedown="handleClear" v-if="searchContent"></i>
      </div>
      <bk-popover v-bind="popoverProps" ref="popover">
        <bk-button theme="primary" class="search-btn" v-test-id="'search'"
          :loading="$loading(request.search)"
          @click="handleSearch()">
          <i class="bk-icon icon-search"></i>
          {{$t('搜索')}}
        </bk-button>
        <div class="picking-popover-content" slot="content">
          <p>{{$t('检测到输入框包含多种格式数据，请选择以哪个字段进行搜索：')}}</p>
          <div class="buttons">
            <bk-button theme="primary" size="small" outline v-test-id="'ipSearch'" v-if="searchFlag.ip"
              @click="handleSearch('ip')">
              {{$t('IP')}}
            </bk-button>
            <bk-button theme="primary" size="small" outline v-test-id="'assetSearch'" v-if="searchFlag.asset"
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
  import FilterUtils from '@/components/filters/utils.js'
  import { HOME_HOST_SEARCH_CONTENT_STORE_KEY } from '@/dictionary/storage-keys.js'
  import { IP_SEARCH_MAX_CLOUD, IP_SEARCH_MAX_COUNT } from '@/setup/validate'
  import { ALL_IP_REGEXP, LT_REGEXP } from '@/dictionary/regexp'

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
        pasteData: {
          cursor: 0,
          length: 0,
          input: false
        },
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
        searchFlag: {
          ip: false,
          asset: false
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
      this.textareaDom = this.$refs.searchInput
      this.setInputHtml(this.searchContent)
    },
    methods: {
      // 获取光标位置
      getCursorPosition() {
        const selection = window.getSelection()
        const element = this.textareaDom
        let caretOffset = 0
        // false表示进行了范围选择
        const { isCollapsed } = selection
        // 选中的区域
        if (selection.rangeCount > 0) {
          const range = selection.getRangeAt(0)
          // 克隆一个选中区域
          const preCaretRange = range.cloneRange()
          // 设置选中区域的节点内容为当前节点
          preCaretRange.selectNodeContents(element)
          // 重置选中区域的结束位置
          preCaretRange.setEnd(range.endContainer, range.endOffset)
          const { length } = preCaretRange.toString()
          caretOffset = isCollapsed ? length : length - selection.toString().length
        }
        return caretOffset
      },
      // 设置光标位置
      setCursorPostion() {
        const selection = window.getSelection()
        const element = document.getElementsByClassName('search-input')[0].getElementsByClassName('new-paste')[0]
        // 创建一个选中区域
        const range = document.createRange()
        // 选中节点的内容
        range.selectNodeContents(element)
        if (element.innerHTML.length > 0) {
          // 设置光标起始为指定位置
          range.setStart(element, element.childNodes.length)
        }
        // 设置选中区域为一个点
        range.collapse(true)
        // 移除所有的选中范围
        selection.removeAllRanges()
        // 添加新建的范围
        selection.addRange(range)
      },
      getSearchList() {
        // 使用切割IP的方法分割内容，方法在此处完全适用且能与高级搜索的IP分割保持一致
        return FilterUtils.splitIP(this.searchContent)
      },
      setRows() {
        const rows = this.searchContent.split('\n').length || 1
        this.rows = Math.min(10, rows)
      },
      // 清除一些结构不太完整的ip， 比如‘2.168.1.5]’
      deleteIp(searchList) {
        const ipList = Array.from(searchList)
        const IPs = FilterUtils.parseIP(ipList)
        const { assetList } = IPs
        assetList?.map(ip => searchList.delete(ip))
        return searchList
      },
      parseIp() {
        if (!this.pasteData.input) return
        this.pasteData = {
          cursor: 0,
          length: 0,
          input: false
        }
        const ipList = new Set(this.searchContent.match(ALL_IP_REGEXP))
        const newHtml = (Array.from(this.deleteIp(ipList))?.map(ip => ip) || []).join('\n')
        this.searchContent = newHtml
        this.setInputHtml(newHtml)
      },
      handleClear(event) {
        this.pasteData = {
          cursor: 0,
          length: 0,
          input: false
        }
        this.searchContent = ''
        this.setInputHtml()
        event.preventDefault()
      },
      handleFocus() {
        this.$emit('focus', true)
        this.setRows()
      },
      handleBlur() {
        this.parseIp()
        this.textareaDom && this.textareaDom.blur()
        this.$emit('focus', false)
      },
      handleKeydown(event) {
        const agent = window.navigator.userAgent.toLowerCase()
        const isMac = /macintosh|mac os x/i.test(agent)
        const modifierKey = isMac ? event.metaKey : event.ctrlKey
        if (modifierKey && event.code.toLowerCase() === 'enter') {
          this.handleSearch()
        }
      },
      handlePaste(event) {
        const val = event?.clipboardData?.getData('text')?.replace(/\r/g, '')
        this.pasteData = {
          cursor: this.getCursorPosition(),
          length: val.length
        }
      },
      handleInput(event) {
        const { inputType } = event
        const val = this.$refs.searchInput.innerText
        this.searchContent = val
        this.pasteData.input = true
        // 如果是粘贴，需要处理实时高光逻辑
        if (inputType === 'insertFromPaste') {
          const { cursor, length } = this.pasteData
          this.setHighLight(cursor, length)
        }
      },
      // 处理粘贴数据高光
      setHighLight(cursor, length) {
        const content = this.searchContent
        const start = content.substring(0, cursor).replace(LT_REGEXP, '&lt')
        const end = content.substring(cursor + length).replace(LT_REGEXP, '&lt')
        const paste = `<span class="new-paste">${content.substring(cursor, cursor + length).replace(LT_REGEXP, '&lt')}</span>`
        const newHtml = start + paste.replace(ALL_IP_REGEXP, ' <span class="high-light">$<ip></span> ') + end
        this.setInputHtml(newHtml)
        this.setCursorPostion()
      },
      setInputHtml(html = '') {
        this.textareaDom.innerHTML = html
      },
      async handleSearch(force = '') {
        const searchList = this.getSearchList()
        if (searchList.length > IP_SEARCH_MAX_COUNT) {
          this.$warn(this.$t('最多支持搜索10000条数据'))
          return
        }

        // 保存本次搜索内容
        window.sessionStorage.setItem(HOME_HOST_SEARCH_CONTENT_STORE_KEY, JSON.stringify(this.searchContent))

        if (searchList.length) {
          const IPs = FilterUtils.parseIP(searchList)
          const { IPv4List, IPv4WithCloudList, IPv6List, IPv6WithCloudList, assetList, cloudIdSet } = IPs

          this.searchFlag.ip = IPv4List.length
            || IPv4WithCloudList.length
            || IPv6List.length
            || IPv6WithCloudList.length
          this.searchFlag.asset = assetList.length > 0
          const isMixSearch = Object.values(this.searchFlag).filter(x => x).length > 1

          // 判断是否存在IP、IPv6、固资编号混合搜索
          if (!force && isMixSearch) {
            this.$refs.popover.showHandler()
            return
          }

          const assetSearch = () => this.handleAssetSearch(assetList)

          const ipSearch = () => {
            // 不同管控区域+IP的混合搜索
            if (cloudIdSet.size > IP_SEARCH_MAX_CLOUD) {
              return this.$warn(this.$t('最多支持50个不同管控区域的混合搜索'))
            }

            this.handleIPSearch(IPs)
          }

          // 优先使用混合搜索下的选择
          if (force === 'asset') {
            return assetSearch()
          }
          if (force === 'ip') {
            return ipSearch()
          }

          // 非混合搜索
          if (this.searchFlag.asset) {
            return assetSearch()
          }
          if (this.searchFlag.ip) {
            return ipSearch()
          }
        } else {
          this.searchContent = ''
          this.textareaDom && this.textareaDom.focus()
        }
      },
      handleIPSearch(IPs) {
        const IPList = [...IPs.IPv4List, ...IPs.IPv6List]
        IPs.IPv4WithCloudList.forEach(([cloud, ip]) => IPList.push(`${cloud}:[${ip}]`))
        IPs.IPv6WithCloudList.forEach(([cloud, ip]) => IPList.push(`${cloud}:[${ip}]`))

        const ip = Object.assign(FilterUtils.getDefaultIP(), { text: IPList.join('\n') })

        this.$routerActions.redirect({
          name: MENU_RESOURCE_HOST,
          query: {
            scope: 'all',
            ip: QS.stringify(ip, { encode: false }),
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
    @mixin tip {
        color: #C4C6CC;
        position: sticky;
        left: 0px;
        bottom: -5px;
        font-size: 12px;
        line-height: 17px;
        display: block;
        width: 100%;
        background: white;
        padding-bottom: 5px;
    }
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
    .search-input-close {
      flex: 1;
      max-width: 646px;
      position: relative;
    }
    .search-input[contenteditable]:empty::before {
        content: attr(placeholder);
        color: #C4C6CC;
        cursor: text;
        font-size: 12px;
        position: absolute;
        left: 16px;
    }
    .search-input[contenteditable]:focus {
        border-color: #3A84FF;
    }
    .search-input[contenteditable]:not(:empty):focus::after {
        content: attr(focusTip);
        @include tip;
    }
    .search-input[contenteditable]:not(:empty):not(:focus)::after {
        content: attr(blurTip);
        @include tip;
    }
    .search-input {
        flex: 1;
        max-width: 646px;
        background: white;
        min-height: 100%;
        height: max-content;
        padding: 5px 32px 5px 16px;
        border: 1px solid #c4c6cc;
        border-radius: 0 0 0 2px;
        font-size: 14px;
        line-height: 30px;
        max-height: 400px;
        position: relative;
        @include scrollbar-y;
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
                padding: 5px 32px 5px 16px;
                border-radius: 0 0 0 2px;
            }
            .right-icon {
              right: 20px !important;
            }
            .high-light {
              display: inline-block;
              background: #FFE8C3;
              border-radius: 2px;
              cursor: pointer;
              &:hover {
                background: #FFD695;
              }
            }
        }
    }
    .search-close {
      position: absolute;
      right: 10px;
      top: 13px;
      cursor: pointer;
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
