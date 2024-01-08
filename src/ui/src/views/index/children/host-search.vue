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
      <editable-block
        v-test-id
        ref="ipEditableBlock"
        :placeholder="$t('首页主机搜索提示语')"
        :focus-tip="$t('首页主机搜索聚焦提示')"
        :blur-tip="$t('首页主机搜索失焦焦提示')"
        :value="searchContent"
        @search="handleSearch"
        @blur="handleBlur"
        @focus="handleFocus">
      </editable-block>
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
  import EditableBlock from '@/components/editable-block/index.vue'

  export default {
    components: {
      EditableBlock
    },
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
        ipEditableBlock: null,
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
      },
      'ipEditableBlock.searchContent': {
        handler() {
          this.$nextTick(this.setRows)
        }
      },
    },
    mounted() {
      this.ipEditableBlock = this.$refs.ipEditableBlock
    },
    methods: {
      getAuthMaskProps() {
        const auth = { type: this.$OPERATION.R_RESOURCE_HOST }
        return {
          auth,
          tag: 'div',
          authorized: this.isViewAuthed(auth)
        }
      },
      getSearchList() {
        // 使用切割IP的方法分割内容，方法在此处完全适用且能与高级搜索的IP分割保持一致
        return FilterUtils.splitIP(this.ipEditableBlock.searchContent)
      },
      setRows() {
        const rows = this.ipEditableBlock.searchContent.split('\n').length || 1
        this.rows = Math.min(10, rows)
      },
      handleFocus() {
        this.$emit('focus', true)
        this.setRows()
      },
      handleBlur() {
        this.$emit('focus', false)
      },
      async handleSearch(force = '') {
        const { hasIP, searchContent } = this.ipEditableBlock
        if (!hasIP || !searchContent.length) {
          this.ipEditableBlock?.focus()
          return
        }
        const searchList = this.getSearchList()
        if (searchList.length > IP_SEARCH_MAX_COUNT) {
          this.$warn(this.$t('最多支持搜索10000条数据'))
          return
        }

        // 保存本次搜索内容
        window.sessionStorage.setItem(
          HOME_HOST_SEARCH_CONTENT_STORE_KEY,
          JSON.stringify(this.ipEditableBlock.searchContent)
        )

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
      },
      handleIPSearch(IPs) {
        const IPList = [...IPs.IPv4List, ...IPs.IPv6List]
        IPs.IPv4WithCloudList.forEach(([cloud, ip, type]) => {
          let val = `${cloud}:[${ip}]`
          if (type === 0) {
            val = val.replace(/\[|\]/g, '')
          }
          IPList.push(val)
        })
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
    .host-search-layout {
        position: relative;
        width: 100%;
        max-width: 806px;
        height: 42px;
        margin: 0 auto;

        .auth-mask {
          height: 100%;
        }
    }
    .search-bar {
        position: absolute;
        width: 100%;
        height: 42px;
        z-index: 999;
        display: flex;
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
