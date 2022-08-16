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
  import { computed, defineComponent, reactive, ref, watch, watchEffect } from '@vue/composition-api'
  import store from '@/store'
  import RouterQuery from '@/router/query'
  import routerActions from '@/router/actions'
  import { getDefaultPaginationConfig, getSort, getHeaderProperties, getHeaderPropertyName } from '@/utils/tools.js'
  import { transformGeneralModelCondition, getDefaultData } from '@/components/filters/utils.js'
  import ColumnsConfig from '@/components/columns-config/columns-config.js'
  import PodListOptions from './pod-list-options.vue'
  import {
    MENU_POD_DETAILS,
  } from '@/dictionary/menu-symbol'
  import { CONTAINER_OBJECTS } from '@/dictionary/container'
  import containerPropertyService from '@/service/container/property.js'
  import podService from '@/service/pod/index.js'

  export default defineComponent({
    components: {
      PodListOptions
    },
    setup() {
      const requestIds = {
        property: Symbol(),
        list: Symbol()
      }

      const MODEL_NAME_KEY = 'pod_name'

      const table = reactive({
        data: [],
        header: [],
        selection: [],
        sort: 'id',
        pagination: getDefaultPaginationConfig()
      })

      const columnsConfig = reactive({
        disabledColumns: ['name', 'namespace'],
        selected: []
      })

      const properties = ref([])

      const columnsConfigKey = 'pod_custom_table_columns'

      const bizId = computed(() => store.getters['objectBiz/bizId'])

      const query = computed(() => RouterQuery.getAll())

      const filter = reactive({
        field: query.value.field || MODEL_NAME_KEY,
        value: '',
        operator: '$regex'
      })

      // 在模型管理中可以配置展示列（全局的）
      const globalUsercustom = computed(() => store.getters['userCustom/globalUsercustom'])

      const usercustom = computed(() => store.getters['userCustom/usercustom'])
      const customColumns = computed(() => usercustom.value[columnsConfigKey] || [])
      const globalCustomColumns = computed(() =>  globalUsercustom.value?.pod_global_custom_table_columns || [])

      watch(customColumns, () => setTableHeader())

      watchEffect(async () => {
        const podProperties = await containerPropertyService.getAll({
          objId: CONTAINER_OBJECTS.POD
        }, {
          requestId: requestIds.property,
          fromCache: true
        })
        console.log(podProperties)
        properties.value = podProperties

        setTableHeader()
        getList()
      })

      // 监听查询参数触发查询
      watch(
        query,
        async (query) => {
          const {
            page = 1,
            limit = table.pagination.limit,
            pod_name: name = '',
            operator = '',
            field = MODEL_NAME_KEY
          } = query
          updateFilter(field, name, operator)

          table.pagination.current = parseInt(page, 10)
          table.pagination.limit = parseInt(limit, 10)

          getList()
        }
      )

      const setTableHeader = () => {
        console.log(usercustom.value[columnsConfigKey] || [], 'usercustom.value[columnsConfigKey]', customColumns)
        const customColumns = customColumns?.value.length ? customColumns.value : globalCustomColumns.value
        const headerProperties = getHeaderProperties(properties.value, customColumns, columnsConfig.disabledColumns)

        table.header = headerProperties.map(property => ({
          id: property.bk_property_id,
          name: getHeaderPropertyName(property),
          property
        }))

        columnsConfig.selected = properties.value.map(property => property.bk_property_id)
      }


      // 查询条件组件相关属性数据
      const filterProperty = computed(() => properties.value.find(property => property.id === filter.field))
      const filterType = computed(() => filterProperty.value?.bk_property_type ?? 'singlechar')

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
      const searchParams = computed(() => {
        const params = {
          bk_biz_id: bizId.value,
          fields: table.header.map(item => item.id),
          page: {
            start: table.pagination.limit * (table.pagination.current - 1),
            limit: table.pagination.limit,
            sort: table.sort
          }
        }

        // 这里先直接复用转换通用模型实例查询条件的方法
        const condition = {
          [filter.field]: {
            value: filter.value,
            operator: filter.operator
          }
        }
        console.log(condition, 'conditionconditioncondition')
        const { conditions } = transformGeneralModelCondition(condition, properties.value)

        console.log(conditions, 'conditionsconditions++')

        if (conditions) {
          params.filter = {
            condition: 'AND',
            rules: conditions.rules
          }
        }

        return params
      })

      const getList = async () => {
        try {
          const { list, count } = await podService.find(searchParams.value, {
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
        if (column.bk_property_id !== 'id') {
          return
        }
        routerActions.redirect({
          name: MENU_POD_DETAILS,
          params: {
            bizId: bizId.value,
            id: row.id
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
            properties: [],
            selected: [],
            disabledColumns: []
          },
          handler: {
            apply: async (properties) => {
              console.log(properties)
            },
            reset: async () => {
            }
          }
        })
      }

      return {
        requestIds,
        table,
        filter,
        handlePageChange,
        handleLimitChange,
        handleSortChange,
        handleValueClick,
        handleSelectionChange,
        handleHeaderClick
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
      ref="table"
      v-bkloading="{ isLoading: $loading(Object.values(requestIds)) }"
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
        show-overflow-tooltip
        :min-width="column.id === 'id' ? 80 : 120"
        :key="column.id"
        :sortable="'custom'"
        :prop="column.id"
        :label="column.name"
        :fixed="column.id === 'id'">
        <template slot-scope="{ row }">
          <cmdb-property-value
            :theme="column.id === 'id' ? 'primary' : 'default'"
            :value="row[column.id]"
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
