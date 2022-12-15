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
  <div class="host-table">
    <div class="search-bar">
      <bk-input
        clearable
        right-icon="icon-search"
        @change="handleFilter"
        v-model.trim="keyword">
      </bk-input>
    </div>
    <bk-table
      ref="table"
      row-class-name="clickable-row"
      :data="displayList"
      :max-height="tableMaxHeight"
      :outer-border="false"
      :header-border="false"
      @row-click="handleRowClick"
    >
      <batch-selection-column
        width="60"
        :cross-page="false"
        ref="batchSelectionColumn"
        :selected-rows="selected"
        @selection-change="handleSelectionChange"
        :data="displayList"
        row-key="host.bk_host_id"
        reserve-selection>
      </batch-selection-column>
      <bk-table-column :label="$t('内网IP')">
        <template slot-scope="{ row }">
          {{row.host.bk_host_innerip}}
        </template>
      </bk-table-column>
      <bk-table-column :label="$t('云区域')" show-overflow-tooltip>
        <template slot-scope="{ row }">{{row.host.bk_cloud_id | foreignkey}}</template>
      </bk-table-column>
    </bk-table>
    <bk-pagination
      v-if="pagination"
      small
      :current.sync="innerPagination.current"
      :count="innerPagination.count"
      @change="handlePaginationChange"
      :show-total-count="true"
      @limit-change="handleLimitChange"
      :limit-list="[10, 20, 50,100,500]"
      :limit.sync="innerPagination.limit" />
  </div>
</template>

<script>
  import debounce from 'lodash.debounce'
  import { foreignkey } from '@/filters/formatter.js'
  import BatchSelectionColumn from '@/components/batch-selection-column'
  export default {
    components: {
      BatchSelectionColumn
    },
    filters: {
      foreignkey
    },
    props: {
      list: {
        type: Array,
        default: () => ([])
      },
      selected: {
        type: Array,
        default: () => ([])
      },
      pagination: {
        type: Object,
        default: null
      }
    },
    data() {
      return {
        keyword: '',
        displayList: [],
        innerPagination: {
          current: 1,
          limit: 500,
          count: 0
        }
      }
    },
    computed: {
      tableMaxHeight() {
        return this.pagination ? 360 : 400
      },
    },
    watch: {
      list(list) {
        this.displayList = list
      },
      'pagination.count': {
        immediate: true,
        handler(val) {
          if (val) {
            this.innerPagination.count = val
          }
        }
      },
    },
    created() {
      this.handleFilter = debounce(this.filterHost, 300)
    },
    methods: {
      handleSelectionChange(selection) {
        const ids = [...new Set(selection.map(data => data.host.bk_host_id))]
        const removed = this.displayList.filter(item => !ids.includes(item.host.bk_host_id))

        if (selection.length > 500) {
          this.$bkMessage({
            message: '批量转移主机数量不能超过 500'
          })
          return false
        }
        this.$emit('select-change', { removed, selected: selection })
      },
      handleRowClick(row) {
        this.$refs.batchSelectionColumn.toggleRowSelection(row)
      },
      filterHost() {
        if (this.keyword) {
          this.displayList = this.list.filter(item => new RegExp(this.keyword, 'i').test(item.host.bk_host_innerip))
        } else {
          this.displayList = this.list
        }
      },
      handleLimitChange() {
        this.innerPagination.current = 1
        this.handlePaginationChange()
      },
      handlePaginationChange() {
        const val = {
          start: (this.innerPagination.current - 1) * this.innerPagination.limit,
          limit: this.innerPagination.limit,
          count: this.innerPagination.count
        }

        this.$emit('update:pagination', val)
        this.$emit('pagination-change', val)
      }
    }
  }
</script>

<style lang="scss" scoped>
    .host-table {
        height: 100%;

        .search-bar {
            margin-bottom: 12px;
        }
    }

    ::v-deep .clickable-row {
      cursor: pointer;
    }

    ::v-deep .bk-page .bk-page-total-small{
      line-height: 40px;
      margin-top: 0;
    }
</style>
