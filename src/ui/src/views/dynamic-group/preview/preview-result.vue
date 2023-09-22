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
    <div class="header">
      <div class="title">
        {{$t('结果预览')}}
        <span class="date" v-if="hasCondition">
          (
          {{ $t('生成时间xx', { date: time(now) }) }}
          )
        </span>
      </div>
      <bk-button
        theme="default"
        icon="refresh"
        class="mr10 refresh"
        size="small"
        v-if="hasCondition"
        @click="handleRefresh">
        {{$t('刷新')}}
      </bk-button>
    </div>
    <bk-table v-if="hasCondition" class="host-table" v-test-id.businessHostAndService="'hostList'"
      ref="tableRef"
      v-bkloading="{ isLoading: $loading(Object.values(request)) }"
      :data="table.data"
      :pagination="table.pagination"
      :max-height="$APP.height - filtersTagHeight - 250"
      @page-change="handlePageChange"
      @page-limit-change="handleLimitChange"
      @sort-change="handleSortChange"
      @header-click="handleHeaderClick">
      <bk-table-column v-for="column in tableHeader"
        :show-overflow-tooltip="$tools.isShowOverflowTips(column)"
        :min-width="getColumnMinWidth(column)"
        :key="column.bk_property_id"
        :sortable="getColumnSortable(column)"
        :prop="column.bk_property_id"
        :fixed="column.bk_property_id === 'bk_host_id'"
        :render-header="() => renderHeader(column)">
        <template slot-scope="{ row }">
          <cmdb-property-value
            :ref="getTableCellPropertyValueRefId(column)"
            :theme="column.bk_property_id === 'bk_host_id' ? 'primary' : 'default'"
            :value="row | hostValueFilter(column.bk_obj_id, column.bk_property_id)"
            :show-unit="false"
            :property="column"
            :multiple="column.bk_obj_id !== 'host'"
            :instance="row"
            show-on="cell"
            @click.native.stop="handleValueClick(row, column)">
          </cmdb-property-value>
        </template>
      </bk-table-column>
      <bk-table-column type="setting" :tippy-options="{ zIndex: -1 }"></bk-table-column>
      <cmdb-table-empty
        slot="empty"
        :stuff="table.stuff">
      </cmdb-table-empty>
    </bk-table>
    <div class="no-condition" v-else>
      <cmdb-data-empty
        slot="empty"
        :stuff="dataEmpty">
      </cmdb-data-empty>
    </div>
  </div>
</template>

<script setup>
  import hostValueFilter from '@/filters/host'
  import tableMixin from '@/mixins/table'
  import { computed, ref, watch, reactive, h } from 'vue'
  import store from '@/store'
  import { t } from '@/i18n'
  import {
    MENU_BUSINESS_HOST_DETAILS,
  } from '@/dictionary/menu-symbol'
  import { time } from '@/filters/formatter'
  import FilterStore, { setupFilterStore } from '../store'
  import ColumnsConfig from '@/components/columns-config/columns-config.js'
  import { CONTAINER_OBJECTS } from '@/dictionary/container.js'
  import routerActions from '@/router/actions'
  import {
    getDefaultPaginationConfig,
    isPropertySortable,
    getHeaderPropertyName,
    getHeaderPropertyMinWidth,
    isUseComplexValueType,
    getSort,
    getPageParams
  } from '@/utils/tools'

  const props = defineProps({
    condition: {
      type: Object,
      default: () => ({})
    },
    mode: String
  })

  const now = ref(new Date())
  const filtersTagHeight = ref(0)
  const tableHeader = ref([])

  const dataEmpty = reactive({
    type: 'empty',
    payload: {
      defaultText: t('请先在左侧设置分组条件')
    }
  })
  const table = reactive({
    data: [],
    selection: [],
    sort: 'bk_host_id',
    pagination: getDefaultPaginationConfig(),
    stuff: {
      type: 'default',
      payload: {
        emptyText: t('bk.table.emptyText')
      }
    }
  })
  const request = reactive({
    table: Symbol('table'),
    moveToResource: Symbol('moveToResource'),
    moveToIdleModule: Symbol('moveToIdleModule')
  })

  const hasCondition = computed(() => Object.keys(props.condition)?.length > 0)
  const customInstanceColumnKey = computed(() => {
    if (props.mode === 'set') return 'dynamic_group_search_object_cluster'
    return 'business_topology_table_column_config'
  })
  const bizId = computed(() => store.getters['objectBiz/bizId'])
  const getModelById = computed(() => store.getters['objectModelClassify/getModelById'])

  const handleRefresh = (() => {
    pageCurrentChange()
    now.value = new Date()
  })
  const initFilterStore = (async () => {
    await setupFilterStore({
      bizId: bizId.value,
      mode: props.mode,
      header: {
        custom: 'business_topology_table_column_config',
        cluster: 'dynamic_group_search_object_cluster'  // 集群字段
      }
    })
  })
  const getColumnSortable = (column => (isPropertySortable(column) ? 'custom' : false))
  const renderHeader = ((property) => {
    const content = [getHeaderPropertyName(property)]
    const modelId = property.bk_obj_id
    if (modelId !== 'host' && modelId !== CONTAINER_OBJECTS.NODE) {
      const model = getModelById.value(modelId)
      const suffix = h('span', { style: { color: '#979BA5', marginLeft: '4px' } }, [`(${model.bk_obj_name})`])
      content.push(suffix)
    }
    return h('span', {}, content)
  })
  const getColumnMinWidth = ((property) => {
    let name = getHeaderPropertyName(property)
    const modelId = property.bk_obj_id
    if (modelId !== 'host' && modelId !== CONTAINER_OBJECTS.NODE) {
      const model = getModelById.value(modelId)
      name = `${name}(${model.bk_obj_name})`
    }
    return getHeaderPropertyMinWidth(property, {
      name,
      hasSort: isPropertySortable(property)
    })
  })
  const getTableCellPropertyValueRefId = (property => (isUseComplexValueType(property) ? `table-cell-property-value-${property.bk_property_id}` : null))
  const pageCurrentChange = ((current = 1) => {
    if (table.pagination.current === current) {
      getHostList()
      return
    }
    table.pagination.current = current
  })
  const handlePageChange = ((current = 1) => pageCurrentChange(current))
  const handleLimitChange = ((limit) => {
    table.pagination.limit = limit
    pageCurrentChange()
  })
  const handleSortChange = (sort => table.sort = getSort(sort))
  const handleValueClick = ((row, column) => {
    if (column.bk_obj_id !== 'host' || column.bk_property_id !== 'bk_host_id') {
      return
    }
    routerActions.open({
      name: MENU_BUSINESS_HOST_DETAILS,
      params: {
        bizId: bizId.value,
        id: row.host.bk_host_id
      }
    })
  })
  const handleHeaderClick = ((column) => {
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
          tableHeader.value = []
          await handleApplyColumnsConfig(properties)
          // 获取最新的表头，内部会读取到上方保存的配置
          tableHeader.value = FilterStore.getHeader()
          getHostList()
        },
        reset: async () => {
          await handleApplyColumnsConfig()
          tableHeader.value = FilterStore.getHeader()
          getHostList()
        }
      }
    })
  })
  const handleApplyColumnsConfig = ((properties = []) => store.dispatch('userCustom/saveUsercustom', {
    [customInstanceColumnKey.value]: properties.map(property => property.bk_property_id)
  }))
  const getSearchRequest = (() => {
    const params = getNormalParams()
    const config = {
      requestId: request.table,
      cancelPrevious: true
    }
    return store.dispatch('hostSearch/searchHost', { params, config })
  })
  const getHostList = (async () => {
    try {
      const result = await getSearchRequest()
      table.data = result.info || []
      table.pagination.count = result.count
    } catch (e) {
      console.error(e)
      table.data = []
      table.pagination.count = 0
    }
  })
  const getNormalParams = (() => {
    const params = {
      ...FilterStore.getSearchParams(),
      page: {
        ...getPageParams(table.pagination),
        sort: table.sort
      }
    }
    return params
  })

  watch(() => props.condition, (val) => {
    now.value = new Date()
    tableHeader.value = FilterStore.getHeader()
    FilterStore.setDynamicCollection(val)
    if (hasCondition.value) {
      pageCurrentChange()
    }
  }, {
    deep: true,
    immediate: true
  })
  watch(() => table.pagination.current, () => {
    getHostList()
  })

  initFilterStore()
</script>

<script>
  export default {
    filters: {
      hostValueFilter
    },
    mixins: [tableMixin],
  }
</script>

<style lang="scss" scoped>
.list-layout {
  overflow: hidden;
  height: 100%;
  position: relative;

  .no-condition {
    position: absolute;
    top: 40%;
    left: 50%;
    transform: translate(-50%);
  }
  .header {
    height: 40px;
    background: #FFFFFF;
    border-bottom: 1px solid #DCDEE5;
    border-left: 1px solid #DCDEE5;
    padding: 0 16px;
    @include space-between;

    .title {
      font-size: 14px;
      color: #313238;
      line-height: 22px;
      font-weight: bold;
    }
    .date {
      font-size: 12px;
      color: #C4C6CC;
      line-height: 16px;
      font-weight: normal;
      display: inline-block;
      margin-left: 5px;
    }
    .refresh {
      :deep(>div) {
        @include space-between;
      }
      :deep(.icon-refresh) {
        font-size: 12px;
        margin-right: 5px;
      }
    }
  }
  .host-table {
    margin: 12px 15px 0 24px;
    width: calc(100% - 39px);
  }
}

</style>
