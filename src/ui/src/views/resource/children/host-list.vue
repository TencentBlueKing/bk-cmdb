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
  <div class="resource-layout">
    <host-list-options></host-list-options>
    <host-filter-tag class="filter-tag" ref="filterTag"></host-filter-tag>
    <bk-table class="hosts-table"
      ref="tableRef"
      v-bkloading="{ isLoading: $loading(Object.values(request)) }"
      :data="table.list"
      :pagination="table.pagination"
      :max-height="$APP.height - filtersTagHeight - 230"
      @selection-change="handleSelectionChange"
      @sort-change="handleSortChange"
      @page-change="handlePageChange"
      @page-limit-change="handleSizeChange"
      @header-click="handleHeaderClick">
      <bk-table-column type="selection" width="60" align="center" fixed
        class-name="bk-table-selection">
      </bk-table-column>
      <bk-table-column v-for="property in tableHeader"
        :show-overflow-tooltip="property.bk_property_type !== 'topology'"
        :min-width="getColumnMinWidth(property)"
        :key="property.id"
        :sortable="isPropertySortable(property) ? 'custom' : false"
        :prop="property.bk_property_id"
        :fixed="['bk_host_id'].includes(property.bk_property_id)"
        :class-name="['bk_host_id'].includes(property.bk_property_id) ? 'is-highlight' : ''"
        :render-header="() => renderHeader(property)">
        <template slot-scope="{ row }">
          <cmdb-host-topo-path
            v-if="property.bk_property_type === 'topology'"
            :host="row"
            :is-resource-assigned="isResourceAssigned"
            :is-container-search-mode="isContainerSearchMode"
            @path-ready="handlePathReady(row, ...arguments)">
          </cmdb-host-topo-path>
          <cmdb-property-value
            v-else
            :theme="['bk_host_id'].includes(property.bk_property_id) ? 'primary' : 'default'"
            :value="row | hostValueFilter(property.bk_obj_id, property.bk_property_id)"
            :show-unit="false"
            :property="property"
            :multiple="property.bk_obj_id !== 'host'"
            @click.native.stop="handleValueClick(row, property)">
          </cmdb-property-value>
        </template>
      </bk-table-column>
      <bk-table-column type="setting"></bk-table-column>
      <cmdb-table-empty slot="empty" :stuff="table.stuff"></cmdb-table-empty>
    </bk-table>
  </div>
</template>

<script>
  import { mapGetters } from 'vuex'
  import hostListOptions from './host-options.vue'
  import hostValueFilter from '@/filters/host'
  import {
    MENU_RESOURCE_HOST,
    MENU_RESOURCE_HOST_DETAILS,
    MENU_RESOURCE_BUSINESS_HOST_DETAILS
  } from '@/dictionary/menu-symbol'
  import RouterQuery from '@/router/query'
  import tableMixin from '@/mixins/table'
  import CmdbHostTopoPath from '@/components/host-topo-path/host-topo-path.vue'
  import HostStore from '../transfer/host-store'
  import HostFilterTag from '@/components/filters/filter-tag'
  import FilterStore, { setupFilterStore } from '@/components/filters/store'
  import ColumnsConfig from '@/components/columns-config/columns-config.js'
  import containerHostService from '@/service/container/host.js'

  export default {
    components: {
      hostListOptions,
      CmdbHostTopoPath,
      HostFilterTag
    },
    filters: {
      hostValueFilter
    },
    mixins: [tableMixin],
    data() {
      return {
        directory: null,
        scope: 1,
        table: {
          checked: [],
          selection: [],
          list: [],
          pagination: {
            current: 1,
            count: 0,
            ...this.$tools.getDefaultPaginationConfig()
          },
          sort: 'bk_host_id',
          exportUrl: `${window.API_HOST}hosts/export`,
          stuff: {
            type: 'default',
            payload: {
              emptyText: this.$t('bk.table.emptyText')
            }
          }
        },
        request: {
          list: Symbol('list')
        },
        filtersTagHeight: 0,
        tableHeader: []
      }
    },
    computed: {
      ...mapGetters(['userName']),
      ...mapGetters('resourceHost', ['activeDirectory']),
      ...mapGetters('objectModelClassify', ['getModelById']),
      moduleProperties() {
        return FilterStore.getModelProperties('module')
      },
      customInstanceColumnKey() {
        if (this.isContainerSearchMode) {
          return this.$route.meta.customContainerInstanceColumn
        }
        return this.$route.meta.customInstanceColumn
      },
      isContainerSearchMode() {
        return FilterStore.isContainerSearchMode
      },
      isResourceAssigned() {
        return FilterStore.isResourceAssigned
      },
      searchMode() {
        return FilterStore.searchMode
      }
    },
    watch: {
      scope() {
        this.setModuleNamePropertyState()
        this.tableHeader = FilterStore.getHeader()
        // 重置selection防止因数据结构不同导致获取数据错误
        this.table.selection = []
      },
      $route() {
        this.initFilterStore()
      },
      searchMode() {
        // 传统与容器模式的表头不一样，需要重新设置
        this.tableHeader = FilterStore.getHeader()
      }
    },
    async created() {
      try {
        this.initFilterStore()
        this.setModuleNamePropertyState()
        this.unwatchRouter = RouterQuery.watch('*', ({
          scope = 1,
          page = 1,
          sort = 'bk_host_id',
          limit = this.table.pagination.limit,
          directory = null
        }) => {
          if (this.$route.name !== MENU_RESOURCE_HOST) {
            return false
          }
          this.table.pagination.current = parseInt(page, 10)
          this.table.pagination.limit = parseInt(limit, 10)
          this.table.sort = sort
          this.directory = parseInt(directory, 10) || null

          this.scope = isNaN(scope) ? 'all' : parseInt(scope, 10)

          FilterStore.setResourceScope(scope)

          this.getHostList()
        }, { throttle: 100 })
        this.unwatchScopeAndDirectory = RouterQuery.watch(['scope', 'directory'], FilterStore.resetAll)
      } catch (error) {
        console.error(error)
      }
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
      this.unwatchScopeAndDirectory()
    },
    methods: {
      async initFilterStore() {
        const currentRouteName = this.$route.name
        if (this.storageRouteName === currentRouteName) return
        this.storageRouteName = currentRouteName
        await setupFilterStore({
          header: {
            custom: this.$route.meta.customInstanceColumn,
            customContainer: this.$route.meta.customContainerInstanceColumn,
            global: 'host_global_custom_table_columns'
          }
        })

        this.tableHeader = FilterStore.getHeader()
      },
      setModuleNamePropertyState() {
        const property = this.moduleProperties.find(property => property.bk_property_id === 'bk_module_name')
        if (property) {
          const normalName = this.$t('模块名')
          const directoryName = this.$t('目录名')
          const scopeModuleName = {
            0: normalName,
            1: directoryName,
            all: `${directoryName}/${normalName}`
          }
          property.bk_property_name = scopeModuleName[this.scope]
        }
      },
      getColumnMinWidth(property) {
        let name = this.$tools.getHeaderPropertyName(property)
        const modelId = property.bk_obj_id
        if (modelId !== 'host') {
          const model = this.getModelById(modelId)
          name = `${name}(${model.bk_obj_name})`
        }

        const preset = {}
        if (property.bk_property_type === 'topology') {
          preset[property.bk_property_id] = 200
        }

        return this.$tools.getHeaderPropertyMinWidth(property, {
          name,
          hasSort: this.isPropertySortable(property) ? 'custom' : false,
          preset
        })
      },
      isPropertySortable(property) {
        return property.bk_obj_id === 'host' && !['foreignkey', 'topology'].includes(property.bk_property_type)
      },
      renderHeader(property) {
        const content = [this.$tools.getHeaderPropertyName(property)]
        const modelId = property.bk_obj_id
        if (modelId !== 'host') {
          const model = this.getModelById(modelId)
          const suffix = this.$createElement('span', { style: { color: '#979BA5', marginLeft: '4px' } }, [`(${model.bk_obj_name})`])
          content.push(suffix)
        }
        return this.$createElement('span', {}, content)
      },
      async getHostList(event) {
        try {
          const { count, info } = await this.getSearchRequest()

          // 容器主机时为每条记录添加biz，统一数据结构供后续使用
          if (this.isContainerSearchMode) {
            info.forEach((item, index) => {
              const bizId = item?.node?.[index]?.bk_biz_id
              item.biz = [{
                bk_biz_id: bizId,
                default: 0
              }]
            })
          }

          this.table.pagination.count = count
          this.table.list = info
          this.table.stuff.type = event ? 'search' : 'default'
        } catch (error) {
          this.table.pagination.count = 0
          this.table.checked = []
          this.table.list = []
          console.error(error)
        }
      },
      getSearchRequest() {
        const params = this.getParams()
        const config = {
          requestId: this.request.list,
          cancelPrevious: true
        }

        if (this.isContainerSearchMode) {
          return containerHostService.findAll(params, config)
        }

        return this.$store.dispatch('hostSearch/searchHost', { params, config })
      },
      getParams() {
        const params = {
          ...FilterStore.getSearchParams(),
          page: {
            ...this.$tools.getPageParams(this.table.pagination),
            sort: this.table.sort
          }
        }

        if (!this.isContainerSearchMode) {
          this.injectScope(params)
          this.scope === 1 && this.injectDirectory(params)
        }

        return params
      },
      injectScope(params) {
        const biz = params.condition.find(condition => condition.bk_obj_id === 'biz')
        if (this.scope === 'all') {
          biz.condition = biz.condition.filter(condition => condition.field !== 'default')
        } else {
          const newMeta = {
            field: 'default',
            operator: '$eq',
            value: this.scope
          }
          // eslint-disable-next-line max-len
          const existMeta = biz.condition.find(({ field, operator }) => field === newMeta.field && operator === newMeta.operator)
          if (existMeta) {
            existMeta.value = newMeta.value
          } else {
            biz.condition.push(newMeta)
          }
        }
        return params
      },
      injectDirectory(params) {
        if (!this.directory) {
          return false
        }
        const moduleCondition = params.condition.find(condition => condition.bk_obj_id === 'module')
        const directoryMeta = {
          field: 'bk_module_id',
          operator: '$eq',
          value: this.directory
        }
        // eslint-disable-next-line max-len
        const existMeta = moduleCondition.condition.find(({ field, operator }) => field === directoryMeta.field && operator === directoryMeta.operator)
        if (existMeta) {
          existMeta.value = directoryMeta.value
        } else {
          moduleCondition.condition.push(directoryMeta)
        }
      },
      handleSelectionChange(selection) {
        this.table.selection = selection
        this.table.checked = selection.map(item => item.host.bk_host_id)
        HostStore.setSelected(selection)
      },
      handleValueClick(item, property) {
        if (property.bk_obj_id !== 'host' || property.bk_property_id !== 'bk_host_id') {
          return
        }
        // eslint-disable-next-line prefer-destructuring
        const business = item.biz[0]
        if (business.default) {
          this.$routerActions.redirect({
            name: MENU_RESOURCE_HOST_DETAILS,
            params: {
              id: item.host.bk_host_id
            },
            query: {
              from: 'resource'
            },
            history: true
          })
        } else {
          this.$routerActions.redirect({
            name: MENU_RESOURCE_BUSINESS_HOST_DETAILS,
            params: {
              business: business.bk_biz_id,
              id: item.host.bk_host_id
            },
            query: {
              from: 'resource'
            },
            history: true
          })
        }
      },
      handlePageChange(current) {
        RouterQuery.set({
          page: current,
          _t: Date.now()
        })
      },
      handleSizeChange(limit) {
        RouterQuery.set({
          limit,
          page: 1,
          _t: Date.now()
        })
      },
      handleSortChange(sort) {
        RouterQuery.set({
          sort: this.$tools.getSort(sort),
          _t: Date.now()
        })
      },
      // 拓扑路径写入数据中，用于复制
      handlePathReady(row, paths) {
        // eslint-disable-next-line no-underscore-dangle
        row.__bk_host_topology__ = paths
      },
      handleHeaderClick(column) {
        if (column.type !== 'setting') {
          return false
        }
        ColumnsConfig.open({
          props: {
            properties: FilterStore.columnConfigProperties,
            selected: FilterStore.defaultHeader.map(property => property.bk_property_id),
            disabledColumns: FilterStore.fixedPropertyIds
          },
          handler: {
            apply: async (properties) => {
              // 先清空表头，防止更新排序后未重新渲染
              this.tableHeader = []

              await this.handleApplyColumnsConfig(properties)

              // 获取最新的表头，内部会读取到上方保存的配置
              this.tableHeader = FilterStore.getHeader()

              FilterStore.dispatchSearch()
            },
            reset: async () => {
              await this.handleApplyColumnsConfig()
              this.tableHeader = FilterStore.getHeader()
              FilterStore.dispatchSearch()
            }
          }
        })
      },
      handleApplyColumnsConfig(properties = []) {
        return this.$store.dispatch('userCustom/saveUsercustom', {
          [this.customInstanceColumnKey]: properties.map(property => property.bk_property_id)
        })
      }
    }
  }
</script>

<style lang="scss" scoped>
    .filter-tag ~ .hosts-table {
        margin-top: 0;
    }
    .hosts-table {
        margin-top: 10px;
    }
</style>
