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
  <bk-input class="filter-fast-search"
    v-model.trim="value"
    :placeholder="$t('请输入IP或固资编号')"
    @enter="handleSearch"
    @paste="handlePaste">
  </bk-input>
</template>

<script>
  import FilterStore from './store'
  import { splitIP, parseIP, getDefaultIP } from '@/components/filters/utils'

  export default {
    data() {
      return {
        value: ''
      }
    },
    methods: {
      async handleSearch() {
        this.dispatchFilter(this.value)
      },
      handlePaste(value, event) {
        event.preventDefault()
        const text = event.clipboardData.getData('text').trim()
        this.dispatchFilter(text)
      },
      dispatchFilter(currentValue) {
        const IPs = parseIP(splitIP(currentValue))
        const IPList = [...IPs.IPv4List, ...IPs.IPv6List]
        IPs.IPv4WithCloudList.forEach(([, ip]) => IPList.push(ip))
        IPs.IPv6WithCloudList.forEach(([, ip]) => IPList.push(ip))
        const ip = Object.assign(getDefaultIP(), { text: currentValue })

        FilterStore.resetPage(true)
        if (IPList.length) {
          FilterStore.updateIP(ip)
        } else if (IPs.assetList.length) {
          FilterStore.createOrUpdateCondition([{
            field: 'bk_asset_id',
            model: 'host',
            operator: '$in',
            value: IPs.assetList
          }])
        }
        this.value = ''
      }
    }
  }
</script>

<style lang="scss" scoped>
    .filter-fast-search {
        display: inline-flex;
    }
</style>
