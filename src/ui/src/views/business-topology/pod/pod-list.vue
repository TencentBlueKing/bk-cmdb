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
  import { computed, defineComponent, reactive, ref, watch, watchEffect } from 'vue'
  import store from '@/store'
  import RouterQuery from '@/router/query'
  import routerActions from '@/router/actions'
  import tableMixin from '@/mixins/table'
  import { getDefaultPaginationConfig, getSort, getHeaderProperties, getHeaderPropertyName } from '@/utils/tools.js'
  import { transformGeneralModelCondition, getDefaultData } from '@/components/filters/utils.js'
  import podValueFilter from '@/filters/pod'
  import ColumnsConfig from '@/components/columns-config/columns-config.js'
  import PodListOptions from './pod-list-options.vue'
  import { MENU_POD_DETAILS } from '@/dictionary/menu-symbol'
  import { CONTAINER_OBJECTS, CONTAINER_OBJECT_PROPERTY_KEYS, CONTAINER_OBJECT_INST_KEYS } from '@/dictionary/container'
  import containerPropertyService, { getPodTopoNodeProps } from '@/service/container/property.js'
  import containerPodService from '@/service/container/pod.js'
  import { getContainerNodeType } from '@/service/container/common.js'

  export default defineComponent({
    components: {
      PodListOptions
    },
    filters: {
      podValueFilter
    },
    mixins: [tableMixin],
    mounted() {
      this.disabledTableSettingDefaultBehavior()
    },
    setup() {
      const requestIds = {
        property: Symbol(),
        list: Symbol()
      }

      const MODEL_ID_KEY = CONTAINER_OBJECT_INST_KEYS[CONTAINER_OBJECTS.POD].ID
      const MODEL_NAME_KEY = CONTAINER_OBJECT_INST_KEYS[CONTAINER_OBJECTS.POD].NAME
      const MODEL_FULL_NAME_KEY = CONTAINER_OBJECT_INST_KEYS[CONTAINER_OBJECTS.POD].FULL_NAME

      const tableRef = ref(null)

      const table = reactive({
        data: [],
        header: [],
        selection: [],
        sort: MODEL_ID_KEY,
        pagination: getDefaultPaginationConfig()
      })

      const columnsConfig = reactive({
        disabledColumns: [
          MODEL_ID_KEY,
          MODEL_NAME_KEY
        ]
      })

      const properties = ref([])

      const columnsConfigKey = 'pod_custom_table_columns'

      const bizId = computed(() => store.getters['objectBiz/bizId'])

      const selectedNode = computed(() => store.getters['businessHost/selectedNode'])

      const query = computed(() => RouterQuery.getAll())

      const filter = reactive({
        field: query.value.field || MODEL_FULL_NAME_KEY,
        value: '',
        operator: '$regex'
      })

      // 在模型管理中可以配置展示列（全局的）
      const globalUsercustom = computed(() => store.getters['userCustom/globalUsercustom'])

      const usercustom = computed(() => store.getters['userCustom/usercustom'])
      const customColumns = computed(() => usercustom.value[columnsConfigKey] || [])
      const globalCustomColumns = computed(() =>  globalUsercustom.value?.pod_global_custom_table_columns || [])

      // 查询条件组件相关属性数据
      const filterProperty = computed(() => properties.value.find(property => property.id === filter.field))
      const filterType = computed(() => filterProperty.value?.bk_property_type ?? 'singlechar')

      const getList = async () => {
        if (!selectedNode.value) {
          return
        }

        if (selectedNode.value?.data?.is_folder) {
          table.data = []
          table.pagination.count = 0

          return
        }

        const params = getSearchParams()
        if (!params.fields?.length) {
          return
        }

        try {
          const { list, count } = await containerPodService.find(params, {
            requestId: requestIds.list,
            cancelPrevious: true,
            globalPermission: false
          })

          table.data = list
          table.pagination.count = count
        } catch (error) {
          console.error(error)
        }
      }

      const saveColumnsConfig = (properties = []) => store.dispatch('userCustom/saveUsercustom', {
        [columnsConfigKey]: properties.map(property => property.bk_property_id)
      })

      // 自定义表头变更后，更新table.header
      watch(customColumns, () => setTableHeader())

      watchEffect(async () => {
        const podProperties = await containerPropertyService.getMany({
          objId: CONTAINER_OBJECTS.POD
        }, {
          requestId: requestIds.property,
          fromCache: true
        })

        properties.value = [...podProperties, ...getPodTopoNodeProps()]

        setTableHeader()

        getList()
      })

      // 监听查询参数触发查询
      watch(
        query,
        async (query) => {
          const {
            tab = 'podList',
            node,
            page = 1,
            limit = table.pagination.limit,
            value = '',
            operator = '',
            field = MODEL_FULL_NAME_KEY
          } = query
          updateFilter(field, value, operator)

          table.pagination.current = parseInt(page, 10)
          table.pagination.limit = parseInt(limit, 10)

          if (tab === 'podList' && node && selectedNode.value) {
            getList()
          }
        }
      )

      const setTableHeader = () => {
        const configColumns = customColumns.value.length ? customColumns.value : globalCustomColumns.value
        const headerProperties = getHeaderProperties(properties.value, configColumns, columnsConfig.disabledColumns)

        table.header = headerProperties.map(property => ({
          id: property.bk_property_id,
          name: getHeaderPropertyName(property),
          property
        }))
      }

      const clearTableHeader = () => {
        table.header = []
      }

      // 更新filter数据，无值状态时则使用默认数据初始化
      const updateFilter = (field, value = '', operator = '') => {
        if (field) {
          filter.field = field
        }

        if (!filterProperty.value) return

        // 业务集中的singlechar类型统一使用$regex
        const options = filterType.value === 'singlechar' ? { operator: '$regex', value: '' } : {}
        const defaultData = { ...getDefaultData(filterProperty.value), ...options }

        filter.operator = operator || defaultData.operator
        filter.value = value || defaultData.value
      }

      // 计算查询条件参数
      const getSearchParams = () => {
        const params = {
          bk_biz_id: bizId.value,
          fields: table.header.map(item => item.id),
          page: {
            start: table.pagination.limit * (table.pagination.current - 1),
            limit: table.pagination.limit,
            sort: table.sort
          },
          filter: {
            condition: 'AND',
            rules: []
          }
        }

        const selectedNodeData = selectedNode.value.data

        // 容器节点的属性ID
        const fieldMap = {
          [CONTAINER_OBJECTS.CLUSTER]: CONTAINER_OBJECT_PROPERTY_KEYS[CONTAINER_OBJECTS.CLUSTER].ID,
          [CONTAINER_OBJECTS.NAMESPACE]: CONTAINER_OBJECT_PROPERTY_KEYS[CONTAINER_OBJECTS.NAMESPACE].ID
        }
        const nodeType = getContainerNodeType(selectedNodeData.bk_obj_id)

        if (nodeType === CONTAINER_OBJECTS.WORKLOAD) {
          params.filter.rules.push({
            field: 'ref',
            operator: 'filter_object',
            value: {
              condition: 'AND',
              rules: [
                {
                  field: 'id',
                  operator: 'equal',
                  value: selectedNodeData.bk_inst_id
                },
                {
                  field: 'kind',
                  operator: 'equal',
                  value: selectedNodeData.bk_obj_id
                }
              ]
            }
          })
        } else if (nodeType !== CONTAINER_OBJECTS.FOLDER) {
          // 添加节点的属性ID参数，如 bk_namespace_id
          params.filter.rules.push({
            field: fieldMap[nodeType],
            operator: 'equal',
            value: selectedNodeData.bk_inst_id
          })
        }

        const condition = {
          [filter.field]: {
            value: filter.value,
            operator: filter.operator
          }
        }

        const { conditions } = transformGeneralModelCondition(condition, properties.value)

        if (conditions) {
          params.filter.rules.push(...conditions.rules)
        }

        return params
      }

      const getColumnSortable = (column) => {
        const topoNodePropIds = getPodTopoNodeProps()?.map(prop => prop.bk_property_id) ?? []
        return !topoNodePropIds.includes(column.id) ? 'custom' : false
      }

      const handlePageChange = (current = 1) => {
        RouterQuery.set({
          page: current,
          _t: Date.now()
        })
      }

      const handleLimitChange = (limit) => {
        RouterQuery.set({
          limit,
          page: 1,
          _t: Date.now()
        })
      }

      const handleSortChange = (sort) => {
        table.sort = getSort(sort)
        RouterQuery.set('_t', Date.now())
      }

      const handleValueClick = (row, column) => {
        if (column.id !== 'id') {
          return
        }
        routerActions.redirect({
          name: MENU_POD_DETAILS,
          params: {
            bizId: bizId.value,
            podId: row.id
          },
          history: true
        })
      }

      const handleSelectionChange = (selection) => {
        table.selection = selection
      }

      const handleHeaderClick = (column) => {
        if (column.type !== 'setting') {
          return false
        }
        ColumnsConfig.open({
          props: {
            properties: properties.value,
            selected: table.header.map(item => item.id),
            disabledColumns: columnsConfig.disabledColumns
          },
          handler: {
            apply: async (properties) => {
              // 先清空表头，防止更新排序后未重新渲染
              clearTableHeader()
              await saveColumnsConfig(properties)
              getList()
            },
            reset: async () => {
              clearTableHeader()
              await saveColumnsConfig([])
              getList()
            }
          }
        })
      }

      return {
        requestIds,
        tableRef,
        table,
        filter,
        handlePageChange,
        handleLimitChange,
        handleSortChange,
        handleValueClick,
        handleSelectionChange,
        handleHeaderClick,
        getColumnSortable
      }
    }
  })
</script>

<template>
  <div class="pod-list">
    <pod-list-options
      :table-header="table.header"
      :table-selection="table.selection"
      :filter="filter">
    </pod-list-options>
    <bk-table class="pod-table"
      v-bkloading="{ isLoading: $loading(Object.values(requestIds)) }"
      ref="tableRef"
      :data="table.data"
      :pagination="table.pagination"
      :max-height="$APP.height - 250"
      @page-change="handlePageChange"
      @page-limit-change="handleLimitChange"
      @sort-change="handleSortChange"
      @selection-change="handleSelectionChange"
      @header-click="handleHeaderClick">
      <bk-table-column type="selection" width="50" align="center" fixed></bk-table-column>
      <bk-table-column v-for="column in table.header"
        :show-overflow-tooltip="!['map'].includes(column.property.bk_property_type)"
        :min-width="$tools.getHeaderPropertyMinWidth(column.property, { hasSort: true })"
        :key="column.id"
        :sortable="getColumnSortable(column)"
        :prop="column.id"
        :label="column.name"
        :fixed="column.id === 'id'">
        <template slot-scope="{ row }">
          <cmdb-property-value
            :theme="column.id === 'id' ? 'primary' : 'default'"
            :value="row | podValueFilter(column.id)"
            :show-unit="false"
            :property="column.property"
            @click.native.stop="handleValueClick(row, column)">
          </cmdb-property-value>
        </template>
      </bk-table-column>
      <bk-table-column type="setting"></bk-table-column>
    </bk-table>
  </div>
</template>

<style lang="scss" scoped>
.pod-list {
  overflow: hidden;
}

.pod-table {
  margin-top: 10px;
}
</style>
