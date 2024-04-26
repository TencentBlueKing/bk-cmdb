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
    <bk-table v-if="hasCondition" class="result-list"
      ref="tableRef"
      v-bkloading="{ isLoading: $loading(Object.values(request)) }"
      :data="table.data"
      :pagination="table.pagination"
      :height="$APP.height - 115"
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
            :value="getValue(row, column.bk_obj_id, column.bk_property_id)"
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

    <ul class="copy-list" ref="copyRef" v-show="showCopy">
      <li v-for="(item, index) in copyList"
        class="copy-item"
        :key="index"
        @click="(event) => handleCopy(event, item)">
        {{item.bk_property_name}}
      </li>
    </ul>
  </div>
</template>

<script setup>
  import hostValueFilter from '@/filters/host'
  import tableMixin from '@/mixins/table'
  import { computed, ref, watch, reactive, h, nextTick, getCurrentInstance } from 'vue'
  import store from '@/store'
  import { t } from '@/i18n'
  import {
    MENU_BUSINESS_HOST_DETAILS,
  } from '@/dictionary/menu-symbol'
  import { time } from '@/filters/formatter'
  import FilterStore, { setupFilterStore } from '../store'
  import hostSearchService from '@/service/host/search'
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
    getPageParams,
    getPropertyCopyValue,
    isEmptyPropertyValue,
    transformHostSearchParams
  } from '@/utils/tools'
  import { transformGeneralModelCondition } from '@/components/filters/utils.js'
  import { BUILTIN_MODELS } from '@/dictionary/model-constants'
  import FilterUtils from '@/components/filters/utils'
  import { $bkPopover, $error, $success } from '@/magicbox/index.js'
  import { rollReqUseCount } from '@/service/utils.js'

  const { proxy } = getCurrentInstance()
  const props = defineProps({
    condition: {
      type: Object,
      default: () => ({})
    },
    properties: {
      type: Array,
      default: () => ([])
    },
    mode: String
  })

  const now = ref(new Date())
  const tableHeader = ref([])
  const copyLoading = reactive({
    value: false,
    target: ''
  })
  const copyRef = ref(null)
  const showCopy = ref(false)

  const IPv4Symbol = Symbol('IPV4')
  const IPv6Symbol = Symbol('IPV6')
  const IPWithCloudSymbol = Symbol('IPWithCloud')
  const IPv6WithCloudSymbol = Symbol('IPv6WithCloud')
  const IPv46WithCloudSymbol = Symbol('IPv46WithCloud')
  const IPv64WithCloudSymbol = Symbol('IPv64WithCloud')
  let copyTip = null

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
  const copyList = computed(() => {
    const IPWithCloudFields = {
      [IPv4Symbol]: `${t('内网')}IPv4`,
      [IPv6Symbol]: `${t('内网')}IPv6`,
      [IPWithCloudSymbol]: `${t('管控区域')}:${t('内网')}IPv4`,
      [IPv6WithCloudSymbol]: `${t('管控区域')}:[${t('内网')}IPv6]`,
      [IPv46WithCloudSymbol]: `${t('管控区域')}:${t('内网')}IP(${t('IPv4优先')})`,
      [IPv64WithCloudSymbol]: `${t('管控区域')}:${t('内网')}IP(${t('IPv6优先')})`
    }
    const IPWithClouds = Object.getOwnPropertySymbols(IPWithCloudFields).map(key => FilterUtils.defineProperty({
      id: key,
      bk_obj_id: 'host',
      bk_property_id: key,
      bk_property_name: IPWithCloudFields[key],
      bk_property_type: 'singlechar'
    }))
    return IPWithClouds
  })
  const getModelById = store.getters['objectModelClassify/getModelById']

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
    const content = [<span class="property-name">{getHeaderPropertyName(property)}</span>]
    const modelId = property.bk_obj_id
    if (modelId !== BUILTIN_MODELS.HOST && modelId !== CONTAINER_OBJECTS.NODE) {
      const model = getModelById(modelId)
      const suffix = h('span', { style: { color: '#979BA5', marginLeft: '4px' } }, [`(${model.bk_obj_name})`])
      content.push(suffix)
    }
    const { value, target } = copyLoading
    const showLoading = value && target?.propertyId === property.bk_property_id
    const tooltips = {
      content: `${t('复制内容')}`,
      placement: 'top',
      disabled: value
    }
    const directive = {
      isLoading: showLoading,
      mode: 'spin',
      size: 'mini',
      zIndex: 999
    }
    const copyClass = ['icon-copy', 'bk-icon', 'copy-icon', {
      'no-sort': !getColumnSortable(property)
    }]
    const copy = <i class={ copyClass }
                    disabled={ value }
                    data-copy-loading={ showLoading }
                    on-click={ event => handleCopy(event, property)}
                    v-bk-tooltips={ tooltips }
                    v-bkloading={ directive }>
                  </i>
    content.push(copy)
    return h('span', {}, content)
  })
  const getColumnMinWidth = ((property) => {
    let name = getHeaderPropertyName(property)
    const modelId = property.bk_obj_id
    if (modelId !== BUILTIN_MODELS.HOST && modelId !== CONTAINER_OBJECTS.NODE) {
      const model = getModelById(modelId)
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
      getList()
      return
    }
    table.pagination.current = current
  })
  const showCopyList = (event) => {
    if (!copyTip) {
      copyTip = $bkPopover(event.target?.parentElement, {
        content: copyRef.value,
        allowHTML: true,
        trigger: 'mannual',
        boundary: 'window',
        placement: 'bottom-start',
        theme: 'light',
        distance: 6,
        interactive: true,
        arrow: false,
        onHidden() {
          showCopy.value = false
        },
      })
    }
    nextTick(() => {
      showCopy.value = true
      copyTip?.show()
    })
  }
  const getCopyData = async (property) => {
    const { bk_property_id: propertyId, bk_obj_id: modelId } = property
    copyLoading.value = true
    copyLoading.target = { propertyId, modelId }
    const data = await getCopyList(propertyId)
    return data
  }
  const setCopyData = (property, list = []) => {
    const copyText = list.map((data) => {
      const { bk_property_id: propertyId, bk_obj_id: modelId } = property
      const modelData = data[modelId] ?? data
      const IPWithCloudKeys = [
        IPv4Symbol,
        IPv6Symbol,
        IPWithCloudSymbol,
        IPv6WithCloudSymbol,
        IPv46WithCloudSymbol,
        IPv64WithCloudSymbol
      ]

      if (IPWithCloudKeys.includes(property.id)) {
        const cloud = getPropertyCopyValue(modelData.bk_cloud_id, 'foreignkey')
        const ip = getPropertyCopyValue(modelData.bk_host_innerip, 'singlechar')
        const ipv6 = getPropertyCopyValue(modelData.bk_host_innerip_v6, 'singlechar')
        const isEmptyIPv4Value = isEmptyPropertyValue(modelData.bk_host_innerip)
        const isEmptyIPv6Value = isEmptyPropertyValue(modelData.bk_host_innerip_v6)
        if (property.id === IPv4Symbol) {
          return `${ip}`
        }
        if (property.id === IPv6Symbol) {
          return `${ipv6}`
        }
        if (property.id === IPWithCloudSymbol) {
          return `${cloud}:${ip}`
        }
        if (property.id === IPWithCloudSymbol) {
          return `${cloud}:${ip}`
        }
        if (property.id === IPWithCloudSymbol) {
          return `${cloud}:${ip}`
        }
        if (property.id === IPv6WithCloudSymbol) {
          return `${cloud}:[${ipv6}]`
        }
        if (property.id === IPv46WithCloudSymbol) {
          if (!isEmptyIPv4Value || isEmptyIPv6Value) {
            return `${cloud}:${ip}`
          }
          return `${cloud}:[${ipv6}]`
        }
        if (property.id === IPv64WithCloudSymbol) {
          if (isEmptyIPv4Value || !isEmptyIPv6Value) {
            return `${cloud}:[${ipv6}]`
          }
          return `${cloud}:${ip}`
        }
      }

      const copyValueOptions = {}
      if (propertyId === 'bk_cloud_id') {
        copyValueOptions.isFullCloud = true
      }
      if (Array.isArray(modelData)) {
        const value = modelData
          .map(item => getPropertyCopyValue(item[propertyId], property, copyValueOptions))
        return value.join(',')
      }
      return getPropertyCopyValue(modelData[propertyId], property, copyValueOptions)
    })

    proxy.$copyText(copyText.join('\n')).then(() => {
      $success(t('复制成功'))
    }, () => {
      $error(t('复制失败'))
    })
      .finally(() => {
        copyLoading.value = false
        copyLoading.target = ''
      })
  }
  const handleCopy = async (event, property) => {
    event.stopPropagation()
    if (copyLoading.value) return

    const { bk_property_id: propertyId } = property
    if (propertyId === 'bk_host_innerip') {
      return showCopyList(event)
    }

    copyTip?.hide()
    const data = await getCopyData(property)
    return setCopyData(property, data)
  }
  const handlePageChange = ((current = 1) => pageCurrentChange(current))
  const handleLimitChange = ((limit) => {
    table.pagination.limit = limit
    pageCurrentChange()
  })
  const handleSortChange = (sort => table.sort = getSort(sort))
  const handleValueClick = ((row, column) => {
    if (column.bk_obj_id !== BUILTIN_MODELS.HOST || column.bk_property_id !== 'bk_host_id') {
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
          getList()
        },
        reset: async () => {
          await handleApplyColumnsConfig()
          tableHeader.value = FilterStore.getHeader()
          getList()
        }
      }
    })
  })
  const handleApplyColumnsConfig = ((properties = []) => store.dispatch('userCustom/saveUsercustom', {
    [customInstanceColumnKey.value]: properties.map(property => property.bk_property_id)
  }))
  const getHostList = async (searchParams, page, config) => {
    const params = {
      bk_biz_id: bizId.value,
      ...searchParams,
      page
    }
    return hostSearchService.getBizHosts({ params, config })
  }
  const getValue = (row, modelId, propertyId) => {
    const { mode } = props
    if (mode === BUILTIN_MODELS.HOST) {
      return hostValueFilter(row, modelId, propertyId)
    }
    return  row?.[propertyId]
  }
  const getSetList = async (page, config) => {
    const setParams = getSetParams()
    const { time_condition, conditions } = setParams
    return store.dispatch('objectSet/searchSet', {
      bizId: bizId.value,
      params: {
        page,
        filter: conditions,
        time_condition
      },
      config
    })
  }
  const getSetParams = () => {
    const { properties } = props
    return transformGeneralModelCondition(FilterStore.getCondition(), properties)
  }
  const getParams = (type = 'list') => {
    const searchParams = FilterStore.getSearchParams()
    const page =  {
      ...getPageParams(table.pagination),
      sort: table.sort
    }
    const config = {
      requestId: type === 'list' ? request.table : 'copy',
      cancelPrevious: true
    }
    return { searchParams, page, config }
  }
  const getFields = (propertyId) => {
    if (typeof propertyId === 'symbol') {
      return ['bk_cloud_id', 'bk_host_innerip', 'bk_host_innerip_v6']
    }
    return [propertyId]
  }
  const getCopyList = (async (propertyId) => {
    const { mode } = props
    const { searchParams, page, config } = getParams('copy')
    const fields = getFields(propertyId)

    if (mode === BUILTIN_MODELS.HOST) {
      const url = hostSearchService.getSearchUrl('biz')
      searchParams?.condition?.forEach((condition) => {
        condition.fields = fields
      })
      return rollReqUseCount(url, transformHostSearchParams({
        bk_biz_id: bizId.value,
        ...searchParams,
        page
      }), {}, config)
    }

    const setParams = getSetParams()
    const { time_condition, conditions } = setParams
    return rollReqUseCount(`set/search/${store.getters.supplierAccount}/${bizId.value}`, {
      page,
      filter: conditions,
      fields,
      time_condition
    }, {}, config)
  })
  const getList = (async () => {
    try {
      const { mode } = props
      const { searchParams, page, config } = getParams()

      let result = {}
      if (mode === BUILTIN_MODELS.HOST) {
        result = await getHostList(searchParams, page, config)
      } else {
        result = await getSetList(page, config)
      }

      table.data = result.info || []
      table.pagination.count = result.count
    } catch (e) {
      console.error(e)
      table.data = []
      table.pagination.count = 0
    }
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
    getList()
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
.copy-list {
  width: 240px;
  color: #63656E;
  letter-spacing: 0;
  cursor: pointer;

  .copy-item {
    line-height: 32px;
    margin: 0 -0.6rem;
    padding: 0 0.6rem;

    &:hover {
      background: #E1ECFF;
      color: #3A84FF;
    }
  }
}
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
  .result-list {
    margin: 12px 15px 0 24px;
    width: calc(100% - 39px);

    :deep(.bk-table-header) {
      .cell {
        &:hover {
          .copy-icon {
            opacity: 1;
          }
        }

        .property-name {
            display: inline-block;
            max-width: 100%;
            @include ellipsis;
            line-height: 100%;
        }
      }
    }

    :deep(.copy-icon) {
      position: absolute !important;
      top: 50%;
      transform: translate(18px, -50%);
      font-size: 14px;
      color: #3A84FF;
      opacity: 0;
      cursor: pointer;

      &[disabled] {
        color: #c4c6cc;
      }

      &[data-copy-loading] {
        color: transparent;
      }

      .bk-loading {
        background: transparent !important;
      }
    }

    :deep(.no-sort) {
      transform: translate(0, -50%);
      margin-left: 1px;
    }

    :deep(.bk-table-pagination-wrapper) {
      z-index: 9;
    }

    :deep(.bk-table-body-wrapper ) {
      @include scrollbar-y;
    }
  }
}

</style>
