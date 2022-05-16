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
  <bk-table
    class="instance-table"
    ref="instanceTable"
    v-test-id.businessHostAndService="'svrInstList'"
    v-bkloading="{ isLoading: $loading(request.getList) }"
    :row-class-name="getRowClassName"
    :data="list"
    :pagination="pagination"
    :max-height="$APP.height - 250"
    @page-change="handlePageChange"
    @page-limit-change="handlePageLimitChange"
    @expand-change="handleExpandChange"
    @selection-change="handleSelectionChange"
    @row-click="handleRowClick"
  >
    <bk-table-column type="selection" prop="id"></bk-table-column>
    <bk-table-column
      type="expand"
      prop="expand"
      width="15"
      :before-expand-change="beforeExpandChange"
    >
      <div slot-scope="{ row }" v-bkloading="{ isLoading: row.pending }">
        <expand-list
          :readonly="true"
          :list-request="processListRequest"
          :service-instance="row"
          @update-list="handleExpandResolved(row, ...arguments)"
        >
        </expand-list>
      </div>
    </bk-table-column>
    <bk-table-column :label="$t('实例名称')" prop="name" width="340">
      <template #default="{ row }">
        <span style="font-weight:bold">{{ row.name }}</span>
      </template>
    </bk-table-column>
    <bk-table-column :label="$t('进程数量')" prop="process_count" width="240" :resizable="false">
      <template #default="{ row }">
        {{row.process_count}}
      </template>
    </bk-table-column>
    <bk-table-column :label="$t('标签')" prop="tag" min-width="150">
      <list-cell-tag
        :readonly="true"
        slot-scope="{ row }"
        :row="row">
      </list-cell-tag>
    </bk-table-column>
  </bk-table>
</template>

<script>
  import { mapGetters, mapState } from 'vuex'
  import ExpandList from '@/views/business-topology/service-instance/instance/expand-list'
  import ListCellTag from '@/views/business-topology/service-instance/instance/list-cell-tag'
  import RouterQuery from '@/router/query'
  import Bus from '@/views/business-topology/service-instance/common/bus'
  import has from 'has'
  import { ServiceInstanceService } from '@/service/business-set/service-instance.js'
  import { ProcessInstanceService } from '@/service/business-set/process-instance.js'

  export default {
    components: {
      ExpandList,
      ListCellTag
    },
    data() {
      return {
        list: [],
        selection: [],
        pagination: this.$tools.getDefaultPaginationConfig(),
        filters: [],
        request: {
          getList: Symbol('getList'),
        },
      }
    },
    computed: {
      ...mapGetters(['supplierAccount']),
      ...mapState('bizSet', ['bizSetName', 'bizSetId', 'bizId']),
      ...mapGetters('businessHost', ['selectedNode']),
      searchKey() {
        const nameFilter = this.filters.find(data => data.id === 'name')
        return nameFilter ? nameFilter.values[0].name : ''
      },
      searchTag() {
        const tagFilters = []
        this.filters.forEach((data) => {
          if (data.id === 'name') return
          if (has(data, 'condition')) {
            tagFilters.push({
              key: data.condition.id,
              operator: 'in',
              values: data.values.map(value => value.name),
            })
          } else {
            const [{ id }] = data.values
            tagFilters.push({
              key: id,
              operator: 'exists',
              values: [],
            })
          }
        })
        return tagFilters
      },
    },
    watch: {
      selectedNode() {
        this.handlePageChange(1)
      },
    },
    created() {
      this.unwatch = RouterQuery.watch(
        ['page', 'limit', '_t', 'view'],
        ({
          page = this.pagination.current,
          limit = this.pagination.limit,
          view = 'instance',
        }) => {
          if (view !== 'instance') {
            return false
          }
          this.pagination.page = parseInt(page, 10)
          this.pagination.limit = parseInt(limit, 10)
          this.getList()
        },
        { immediate: true, throttle: true }
      )
      Bus.$on('expand-all-change', this.handleExpandAllChange)
      Bus.$on('filter-change', this.handleFilterChange)
    },
    beforeDestroy() {
      Bus.$off('expand-all-change', this.handleExpandAllChange)
      Bus.$off('filter-change', this.handleFilterChange)
      this.unwatch()
    },
    methods: {
      handleFilterChange(filters) {
        this.filters = filters
        RouterQuery.set({
          node: this.selectedNode.id,
          page: 1,
          _t: Date.now(),
        })
      },
      async getList() {
        try {
          const { count, info } = await ServiceInstanceService.findAll(
            this.bizSetId,
            {
              bk_biz_id: this.bizId,
              bk_module_id: this.selectedNode.data.bk_inst_id,
              page: this.$tools.getPageParams(this.pagination),
              search_key: this.searchKey,
              selectors: this.searchTag,
              with_name: true,
            },
            {
              requestId: this.request.getList,
              cancelPrevious: true,
            }
          )
          this.list = info.map(data => ({
            ...data,
            pending: true,
            editing: { name: false },
          }))
          this.pagination.count = count
        } catch (error) {
          this.list = []
          this.pagination.count = 0
          console.error(error)
        }
      },
      processListRequest(reqParams, reqConfig) {
        return ProcessInstanceService.findProcessByServiceInstance(this.bizSetId, reqParams, reqConfig)
      },
      getRowClassName({ row }) {
        const className = ['instance-table-row']
        if (!row.process_count) {
          className.push('disabled')
        }
        return className.join(' ')
      },
      handlePageChange(page) {
        this.pagination.current = page
        RouterQuery.set({
          node: this.selectedNode.id,
          page,
          _t: Date.now(),
        })
      },
      handlePageLimitChange(limit) {
        this.pagination.limit = limit
        this.pagination.page = 1
        RouterQuery.set({
          limit,
          page: 1,
          _t: Date.now(),
        })
      },
      beforeExpandChange({ row }) {
        return !!row.process_count
      },
      handleExpandAllChange(expanded) {
        this.list.forEach((row) => {
          row.process_count
            && this.$refs.instanceTable.toggleRowExpansion(row, expanded)
        })
      },
      async handleExpandChange(row, expandedRows) {
        row.pending = expandedRows.includes(row)
      },
      handleSelectionChange(selection) {
        this.selection = selection
        Bus.$emit('instance-selection-change', selection)
      },
      handleRowClick(row) {
        this.$refs.instanceTable.toggleRowExpansion(row)
      },
      handleExpandResolved(row, list) {
        row.pending = false
        this.handleRefreshCount(row, list.length)
        if (!row.process_count) {
          this.$refs.instanceTable.toggleRowExpansion(row, false)
        }
      },
      handleRefreshCount(row, newCount) {
        row.process_count = newCount
      }
    },
  }
</script>

<style lang="scss" scoped>
.instance-table {
  /deep/ {
    .instance-table-row {
      &:hover,
      &.expanded {
        td {
          background-color: #f0f1f5;
        }
        .tag-item {
          background-color: #dcdee5;
        }
      }
      &:hover {
        .tag-empty {
          display: none;
        }
      }
      &.disabled {
        .bk-table-expand-icon {
          display: none;
          cursor: not-allowed;
        }
      }
    }
    .bk-table-expanded-cell {
      padding-left: 80px !important;
    }
  }
}
</style>
