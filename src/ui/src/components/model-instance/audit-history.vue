<!--
 * Tencent is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2017 Tencent. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
-->

<template>
  <div class="history">
    <div class="history-filter">
      <cmdb-form-date-range class="filter-item filter-range"
        v-model="condition.operation_time"
        :clearable="false"
        @input="handlePageChange(1)">
      </cmdb-form-date-range>
      <cmdb-form-objuser class="filter-item filter-user"
        v-model="condition.user"
        :exclude="false"
        :multiple="false"
        :placeholder="$t('操作账号')"
        @input="handlePageChange(1)">
      </cmdb-form-objuser>
    </div>
    <bk-table class="history-table"
      v-bkloading="{ isLoading: $loading(requestId) }"
      :data="history"
      :pagination="pagination"
      :max-height="$APP.height - 325"
      :row-style="{ cursor: 'pointer' }"
      @page-change="handlePageChange"
      @page-limit-change="handleSizeChange"
      @row-click="handleRowClick">
      <bk-table-column :label="$t('操作描述')" :formatter="getFormatterDesc"></bk-table-column>
      <bk-table-column prop="user" :label="$t('操作账号')">
        <template #default="{ row }">
          <cmdb-user-value :value="row.user" />
        </template>
      </bk-table-column>
      <bk-table-column prop="operation_time" :label="$t('操作时间11')">
        <template slot-scope="{ row }">
          {{$tools.formatTime(row['operation_time'], 'YYYY-MM-DD HH:mm:ssZZ')}}
        </template>
      </bk-table-column>
      <cmdb-table-empty slot="empty" :stuff="table.stuff" @clear="handleClearFilter"></cmdb-table-empty>
    </bk-table>
  </div>
</template>

<script>
  import AuditDetails from '@/components/audit-history/details.js'
  import tools from '@/utils/tools'
  import resourcePoolService from '@/service/resource-pool/index'

  const today = tools.formatTime(new Date(), 'YYYY-MM-DD')
  const formatValue = () => ({
    operation_time: [today, today],
    user: ''
  })

  export default {
    props: {
      objId: {
        type: String,
        required: true
      },
      bizId: {
        type: Number
      },
      resourceId: {
        type: [Number, String]
      },
      resourceType: {
        type: String,
        default: ''
      }
    },
    data() {
      return {
        history: [],
        dictionary: [],
        condition: {
          ...formatValue(),
          resource_id: this.resourceId,
          resource_type: this.resourceType
        },
        table: {
          stuff: {
            type: 'default',
            payload: {
              emptyText: this.$t('bk.table.emptyText')
            }
          }
        },
        pagination: {
          count: 0,
          current: 1,
          limit: 10
        },
        requestId: Symbol('getList'),
        resourcePoolBizId: undefined
      }
    },
    async created() {
      const resourcePoolData = await resourcePoolService.getBiz()
      this.resourcePoolBizId = resourcePoolData?.bk_biz_id

      this.getAuditDictionary()
      this.getHistory()
    },
    methods: {
      async getAuditDictionary() {
        try {
          this.dictionary = await this.$store.dispatch('audit/getDictionary', {
            fromCache: true,
            globalPermission: false
          })
        } catch (error) {
          this.dictionary = []
        }
      },
      async getHistory() {
        try {
          const { info, count } = await this.$store.dispatch('audit/getInstList', {
            params: {
              condition: this.getUsefulCondition(),
              page: {
                ...this.$tools.getPageParams(this.pagination),
                sort: '-operation_time'
              }
            },
            config: {
              requestId: this.requestId,
              globalPermission: false
            }
          })
          this.table.stuff.type = this.condition.user ? 'search' : 'default'
          this.pagination.count = count
          this.history = info
        } catch ({ permission }) {
          if (permission) {
            this.table.stuff = {
              type: 'permission',
              payload: { permission }
            }
          }
          this.history = []
        }
      },
      getUsefulCondition() {
        const usefuleCondition = {
          bk_obj_id: this.objId
        }

        // 通用模型实例不传，未分配业务主机传1
        if (!isNaN(this.bizId)) {
          usefuleCondition.bk_biz_id = this.bizId > 0 ? this.bizId : this.resourcePoolBizId
        }

        Object.keys(this.condition).forEach((key) => {
          const value = this.condition[key]
          if (String(value).length) {
            usefuleCondition[key] = value
          }
        })
        if (usefuleCondition.operation_time) {
          const [start, end] = usefuleCondition.operation_time
          usefuleCondition.operation_time = {
            start: `${start} 00:00:00`,
            end: `${end} 23:59:59`
          }
        }
        return usefuleCondition
      },
      getFormatterDesc(row) {
        const type = this.dictionary.find(type => type.id === row.resource_type) || {}
        const action = (type.operations || []).find(action => action.id === row.action) || {}
        return `${action.name || row.action}${type.name || row.resource_type}`
      },
      handlePageChange(page) {
        this.pagination.current = page
        this.getHistory(true)
      },
      handleSizeChange(size) {
        this.pagination.limit = size
        this.pagination.current = 1
        this.getHistory()
      },
      handleSortChange(sort) {
        this.sort = this.$tools.getSort(sort)
        this.getHistory()
      },
      handleRowClick(item) {
        AuditDetails.show({
          aduitTarget: 'instance',
          id: item.id,
          resourceType: this.resourceType,
          bizId: this.bizId,
          objId: this.objId
        })
      },
      handleClearFilter() {
        this.condition = {
          ...formatValue(),
          resource_id: this.resourceId,
          resource_type: this.resourceType
        }
        this.getHistory()
      }
    }
  }
</script>

<style lang="scss" scoped>
    .history {
        height: 100%;
    }
    .history-filter {
        padding: 14px 0;
        .filter-item {
            display: inline-block;
            vertical-align: middle;
            &.filter-range {
                width: 300px !important;
                margin: 0 5px 0 0;
            }
            &.filter-user {
                width: 240px;
                height: 32px;
            }
        }
    }
</style>
