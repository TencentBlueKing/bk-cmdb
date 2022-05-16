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
  <bk-table class="process-table" ref="processTable" v-test-id.businessHostAndService="'processList'"
    v-bkloading="{ isLoading: $loading(request.getProcessList) }"
    row-class-name="process-table-row"
    :data="list"
    :pagination="pagination"
    :max-height="$APP.height - 250"
    @page-change="handlePageChange"
    @page-limit-change="handlePageLimitChange"
    @expand-change="handleExpandChange"
    @row-click="handleRowClick">
    <bk-table-column type="expand" width="28">
      <div slot-scope="{ row }" v-bkloading="{ isLoading: row.pending }">
        <expand-list
          :readonly="true"
          :process="row"
          :list-request="instanceListRequest"
          @resolved="handleExpandResolved(row, ...arguments)">
        </expand-list>
      </div>
    </bk-table-column>
    <bk-table-column :label="$t('进程别名')" prop="bk_process_name" width="300" show-overflow-tooltip>
      <span class="process-name" slot-scope="{ row }">{{row.bk_process_name}}</span>
    </bk-table-column>
    <bk-table-column :label="$t('实例数量')">
      <template slot-scope="{ row }">{{row.process_ids.length}}</template>
    </bk-table-column>
  </bk-table>
</template>

<script>
  import { mapGetters, mapState } from 'vuex'
  import RouterQuery from '@/router/query'
  import Bus from '@/views/business-topology/service-instance/common/bus.js'
  import ExpandList from '@/views/business-topology/service-instance/process/expand-list.vue'
  import { ProcessInstanceService } from '@/service/business-set/process-instance.js'
  export default {
    name: 'ProcessList',
    components: {
      ExpandList
    },
    data() {
      return {
        filter: '',
        list: [],
        pagination: this.$tools.getDefaultPaginationConfig(),
        request: {
          getProcessList: Symbol('getProcessList')
        }
      }
    },
    computed: {
      ...mapState('bizSet', ['bizSetId', 'bizId']),
      ...mapGetters('businessHost', ['selectedNode'])
    },
    watch: {
      selectedNode() {
        this.handlePageChange(1)
      }
    },
    created() {
      this.unwatch = RouterQuery.watch(['page', 'limit', '_t'], ({
        page = 1,
        limit = this.pagination.limit
      }) => {
        this.pagination.current = parseInt(page, 10)
        this.pagination.limit = parseInt(limit, 10)
        this.getProcessList()
      }, { immediate: true })
      Bus.$on('expand-all-change', this.handleExpandAllChange)
      Bus.$on('update-reserve-selection', this.handleReserveSelectionChange)
      Bus.$on('filter-list', this.handleFilterList)
    },
    beforeDestroy() {
      this.unwatch()
      Bus.$off('expand-all-change', this.handleExpandAllChange)
      Bus.$off('update-reserve-selection', this.handleReserveSelectionChange)
      Bus.$off('filter-list', this.handleFilterList)
    },
    methods: {
      handleFilterList(value) {
        this.filter = value
        RouterQuery.set({
          page: 1,
          _t: Date.now()
        })
      },
      async getProcessList() {
        try {
          const { count, info } = await ProcessInstanceService.findAll(
            this.bizSetId,
            {
              bk_module_id: this.selectedNode.data.bk_inst_id,
              bk_biz_id: this.bizId,
              process_name: this.filter,
              page: this.$tools.getPageParams(this.pagination)
            },
            {
              requestId: this.request.getProcessList,
              cancelPrevious: true
            }
          )
          this.list = info.map(item => ({ ...item, pending: true, reserved: [] }))
          this.pagination.count = count
        } catch (error) {
          this.list = []
          this.pagination.count = 0
          console.error(error)
        } finally {
          Bus.$emit('process-list-change')
        }
      },
      instanceListRequest(reqParams, reqConfig) {
        return ProcessInstanceService
          .findServiceInstanceByProcess(this.bizSetId, reqParams, reqConfig).then(({ info }) => info)
      },
      handlePageChange(page) {
        RouterQuery.set({
          node: this.selectedNode.id,
          page,
          _t: Date.now()
        })
      },
      handlePageLimitChange(limit) {
        RouterQuery.set({
          limit,
          page: 1,
          _t: Date.now()
        })
      },
      handleExpandChange(row, expandedRows) {
        row.pending = expandedRows.includes(row)
      },
      handleExpandAllChange(expand) {
        this.list.forEach((row) => {
          this.$refs.processTable.toggleRowExpansion(row, expand)
        })
      },
      handleRowClick(row) {
        this.$refs.processTable.toggleRowExpansion(row)
      },
      handleReserveSelectionChange(process, selection) {
        this.list.forEach((row) => {
          row.reserved = row === process ? selection : []
        })
      },
      handleExpandResolved(row, list) {
        row.pending = false
        row.process_ids = list.map(process => process.process_id)
      }
    }
  }
</script>

<style lang="scss" scoped>
    .process-table {
        .process-name {
            font-weight: bold;
        }
    }
    /deep/ {
        .process-table-row {
            &:hover,
            &.expanded {
                td {
                    background-color: #f0f1f5;
                }
            }
            td {
                position: sticky;
                top: 0;
                z-index: 100;
                background-color: #fff;
            }
        }
        .bk-table-expand-icon {
            text-align: right !important;
            justify-content: flex-end !important;
            .bk-icon {
                position: static;
                margin: 0;
            }
        }
    }
</style>
