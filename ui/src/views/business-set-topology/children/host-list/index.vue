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
  <div class="list-layout">
    <host-list-options v-test-id></host-list-options>

    <host-filter-tag class="filter-tag" ref="filterTag"></host-filter-tag>

    <bk-table class="host-table" v-test-id.businessHostAndService="'hostList'"
      ref="table"
      v-bkloading="{ isLoading: $loading(Object.values(requestIds)) }"
      :data="table.data"
      :pagination="table.pagination"
      :max-height="$APP.height - filtersTagHeight - 250"
      @page-change="handlePageChange"
      @page-limit-change="handleLimitChange"
      @sort-change="handleSortChange"
      @selection-change="handleSelectionChange"
      @header-click="handleHeaderClick">
      <bk-table-column type="selection" width="50" align="center" fixed></bk-table-column>
      <bk-table-column v-for="column in tableHeader"
        show-overflow-tooltip
        :min-width="column.bk_property_id === 'bk_host_id' ? 80 : 120"
        :key="column.bk_property_id"
        :sortable="getColumnSortable(column)"
        :prop="column.bk_property_id"
        :fixed="column.bk_property_id === 'bk_host_id'"
        :render-header="() => renderHeader(column)">
        <template slot-scope="{ row }">
          <cmdb-property-value
            :theme="column.bk_property_id === 'bk_host_id' ? 'primary' : 'default'"
            :value="row | hostValueFilter(column.bk_obj_id, column.bk_property_id)"
            :show-unit="false"
            :property="column"
            :multiple="column.bk_obj_id !== BUILTIN_MODELS.HOST"
            @click.native.stop="handleValueClick(row, column)">
          </cmdb-property-value>
        </template>
      </bk-table-column>
      <bk-table-column type="setting"></bk-table-column>
    </bk-table>
  </div>
</template>

<script>
  import has from 'has'
  import HostListOptions from './options.vue'
  import hostValueFilter from '@/filters/host'
  import { mapGetters, mapState } from 'vuex'
  import {
    MENU_BUSINESS_SET_TOPOLOGY,
    MENU_BUSINESS_SET_HOST_DETAILS,
  } from '@/dictionary/menu-symbol'
  import RouterQuery from '@/router/query'
  import HostFilterTag from '@/components/filters/filter-tag'
  import FilterStore, { setupFilterStore } from '@/components/filters/store'
  import ColumnsConfig from '@/components/columns-config/columns-config.js'
  import { HostService } from '@/service/business-set/host.js'
  import { BUILTIN_MODELS } from '@/dictionary/model-constants.js'

  export default {
    name: 'HostList',
    components: {
      HostListOptions,
      HostFilterTag,
    },
    filters: {
      hostValueFilter
    },
    props: {
      active: Boolean
    },
    data() {
      this.BUILTIN_MODELS = BUILTIN_MODELS

      return {
        table: {
          data: [],
          selection: [],
          sort: 'bk_host_id',
          pagination: this.$tools.getDefaultPaginationConfig()
        },
        requestIds: {
          table: Symbol('table')
        },
        filtersTagHeight: 0
      }
    },
    computed: {
      ...mapState('bizSet', ['bizSetId', 'bizId']),
      ...mapGetters('objectModelClassify', ['getModelById']),
      ...mapGetters('businessHost', ['selectedNode']),
      tableHeader() {
        return FilterStore.header
      }
    },
    created() {
      this.initFilterStore()
      this.unwatchRouter = RouterQuery.watch(
        '*', ({
          tab = 'hostList',
          node,
          page = 1,
          limit = this.table.pagination.limit
        }) => {
          if (this.$route.name !== MENU_BUSINESS_SET_TOPOLOGY) {
            return false
          }

          this.table.pagination.current = parseInt(page, 10)
          this.table.pagination.limit = parseInt(limit, 10)

          if (tab === 'hostList' && node && this.selectedNode) {
            this.getHostList()
          }
        },
        { throttle: 16, ignore: ['keyword'] }
      )
    },
    mounted() {
      this.unwatchFilter = this.$watch(() => [FilterStore.condition, FilterStore.IP], () => {
        const el = this.$refs.filterTag.$el
        if (el.getBoundingClientRect) {
          this.filtersTagHeight = el.getBoundingClientRect().height
        } else {
          this.filtersTagHeight = 0
        }
      }, { immediate: true, deep: true })
      this.disabledTableSettingDefaultBehavior()
    },
    beforeDestroy() {
      this.unwatchRouter()
      this.unwatchFilter()
    },
    methods: {
      disabledTableSettingDefaultBehavior() {
        setTimeout(() => {
          const settingReference = this.$refs.table.$el.querySelector('.bk-table-column-setting .bk-tooltip-ref')
          // eslint-disable-next-line no-underscore-dangle
          settingReference?._tippy?.disable()
        }, 1000)
      },
      initFilterStore() {
        if (!FilterStore.hasCondition) {
          setupFilterStore({
            bk_biz_id: () => this.bizId,
            header: {
              custom: this.$route.meta.customInstanceColumn,
              global: 'host_global_custom_table_columns'
            }
          })
        }
      },
      getColumnSortable(column) {
        const isHostProperty = column.bk_obj_id === BUILTIN_MODELS.HOST
        const isForeignKey = column.bk_property_type === 'foreignkey'
        return (isHostProperty && !isForeignKey) ? 'custom' : false
      },
      renderHeader(property) {
        const content = [this.$tools.getHeaderPropertyName(property)]
        const modelId = property.bk_obj_id
        if (modelId !== BUILTIN_MODELS.HOST) {
          const model = this.getModelById(modelId)
          const suffix = this.$createElement('span', { style: { color: '#979BA5', marginLeft: '4px' } }, [`(${model.bk_obj_name})`])
          content.push(suffix)
        }
        return this.$createElement('span', {}, content)
      },
      handlePageChange(current = 1) {
        RouterQuery.set({
          page: current,
          _t: Date.now()
        })
      },
      handleLimitChange(limit) {
        RouterQuery.set({
          limit,
          page: 1,
          _t: Date.now()
        })
      },
      handleSortChange(sort) {
        this.table.sort = this.$tools.getSort(sort)
        RouterQuery.set('_t', Date.now())
      },
      handleValueClick(row, column) {
        if (column.bk_obj_id !== BUILTIN_MODELS.HOST || column.bk_property_id !== 'bk_host_id') {
          return
        }
        this.$routerActions.redirect({
          name: MENU_BUSINESS_SET_HOST_DETAILS,
          params: {
            bizId: this.bizId,
            id: row.host.bk_host_id
          },
          history: true
        })
      },
      handleSelectionChange(selection) {
        this.table.selection = selection
      },
      handleHeaderClick(column) {
        if (column.type !== 'setting') {
          return false
        }
        ColumnsConfig.open({
          props: {
            properties: FilterStore.properties.filter(property => property.bk_obj_id === BUILTIN_MODELS.HOST
              || (property.bk_obj_id === BUILTIN_MODELS.MODULE && property.bk_property_id === 'bk_module_name')
              || (property.bk_obj_id === BUILTIN_MODELS.SET && property.bk_property_id === 'bk_set_name')),
            selected: FilterStore.defaultHeader.map(property => property.bk_property_id),
            disabledColumns: ['bk_host_id', 'bk_host_innerip', 'bk_cloud_id']
          },
          handler: {
            apply: async (properties) => {
              await this.handleApplyColumnsConfig(properties)
              FilterStore.setHeader(properties)
              FilterStore.dispatchSearch()
            },
            reset: async () => {
              await this.handleApplyColumnsConfig()
              FilterStore.setHeader(FilterStore.defaultHeader)
              FilterStore.dispatchSearch()
            }
          }
        })
      },
      handleApplyColumnsConfig(properties = []) {
        return this.$store.dispatch('userCustom/saveUsercustom', {
          [this.$route.meta.customInstanceColumn]: properties.map(property => property.bk_property_id)
        })
      },
      async getHostList() {
        try {
          const params = {
            ...this.getParams()
          }
          const result = await HostService.findAll(this.bizSetId, params, {
            requestId: this.requestIds.table,
            cancelPrevious: true
          })

          this.table.data = result.info
          this.table.pagination.count = result.count
        } catch (e) {
          console.error(e)
          this.table.data = []
          this.table.pagination.count = 0
        }
      },
      getParams() {
        const params = {
          ...FilterStore.getSearchParams(),
          page: {
            ...this.$tools.getPageParams(this.table.pagination),
            sort: this.table.sort
          }
        }
        const instance = this.selectedNode.data
        const fieldMap = {
          biz: 'bk_biz_id',
          set: 'bk_set_id',
          module: 'bk_module_id'
        }
        const instanceCondition = {
          field: fieldMap[instance.bk_obj_id] || 'bk_inst_id',
          operator: '$eq',
          value: instance.bk_inst_id
        }
        const modelConditionId = has(fieldMap, instance.bk_obj_id) ? instance.bk_obj_id : 'object'
        const modelCondition = params.condition.find(modelCondition => modelCondition.bk_obj_id === modelConditionId)

        modelCondition.condition.push(instanceCondition)

        return params
      },
      doLayoutTable() {
        this.$refs.table.doLayout()
      }
    }
  }
</script>

<style lang="scss" scoped>
    .list-layout {
        overflow: hidden;
    }
    .filter-tag ~ .host-table {
        margin-top: 0;
    }
    .host-table {
        margin-top: 10px;
    }
</style>
