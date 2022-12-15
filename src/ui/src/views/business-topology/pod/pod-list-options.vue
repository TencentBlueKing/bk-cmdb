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

<script>
  import { computed, defineComponent, getCurrentInstance } from 'vue'
  import { t } from '@/i18n'
  import { $success, $error } from '@/magicbox/index.js'
  import RouterQuery from '@/router/query'
  import { getPropertyCopyValue } from '@/utils/tools.js'

  export default defineComponent({
    props: {
      tableHeader: {
        type: Array,
        default: () => ([])
      },
      tableSelection: {
        type: Array,
        default: () => ([])
      },
      filter: {
        type: Object,
        default: () => ({})
      }
    },
    setup(props) {
      const $this = getCurrentInstance()

      const clipboardList = computed(() => props.tableHeader.slice())

      const hasSelection = computed(() => !!props.tableSelection?.length)

      const searchHander = (value) => {
        RouterQuery.set({
          _t: Date.now(),
          page: 1,
          field: props.filter.field,
          value,
          operator: props.filter.operator
        })
      }

      const handleCopy = (column) => {
        const copyText = props.tableSelection.map((row) => {
          if (column.id === 'ref') {
            return row?.[column.id]?.name
          }
          return getPropertyCopyValue(row[column.id], column.property)
        })
        $this.proxy.$copyText(copyText.join('\n')).then(() => {
          $success(t('复制成功'))
        }, () => {
          $error(t('复制失败'))
        })
      }

      const handleSearch = async (value) => {
        searchHander(value)
      }

      const handlePaste = (value, event) => {
        event.preventDefault()
        const text = event.clipboardData.getData('text').trim()
        searchHander(text)
      }

      return {
        clipboardList,
        hasSelection,
        handleCopy,
        handleSearch,
        handlePaste
      }
    }
  })
</script>

<template>
  <div class="pod-list-options">
    <div class="options options-left">
      <cmdb-clipboard-selector class="options-clipboard" v-test-id
        :list="clipboardList"
        :disabled="!hasSelection"
        @on-copy="handleCopy">
      </cmdb-clipboard-selector>
    </div>
    <div class="options options-right">
      <bk-input class="filter-fast-search"
        v-model.trim="filter.value"
        :placeholder="$t('请输入名称')"
        @enter="handleSearch"
        @paste="handlePaste">
      </bk-input>
    </div>
  </div>
</template>

<style lang="scss" scoped>
.pod-list-options {
  display: flex;
  justify-content: space-between;
  margin-top: 12px;
}
.options {
  display: flex;
  align-items: center;
  &.options-right {
    overflow: hidden;
    justify-content: flex-end;
  }
}
.filter-fast-search {
  width: 300px;
}
</style>
