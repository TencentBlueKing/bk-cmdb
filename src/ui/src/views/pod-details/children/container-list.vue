<script>
  import { computed, defineComponent, reactive, ref, watch, watchEffect, getCurrentInstance } from 'vue'
  import store from '@/store'
  import { t } from '@/i18n'
  import { $success, $error } from '@/magicbox/index.js'
  import routerActions from '@/router/actions'
  import RouterQuery from '@/router/query'
  import tableMixin from '@/mixins/table'
  import { getDefaultPaginationConfig, getSort, getHeaderProperties, getHeaderPropertyName, getPropertyCopyValue } from '@/utils/tools.js'
  import { transformGeneralModelCondition } from '@/components/filters/utils.js'
  import ColumnsConfig from '@/components/columns-config/columns-config.js'
  import { CONTAINER_OBJECTS, CONTAINER_OBJECT_INST_KEYS } from '@/dictionary/container'
  import { MENU_POD_CONTAINER_DETAILS } from '@/dictionary/menu-symbol'
  import containerPropertyService from '@/service/container/property.js'
  import containerConService from '@/service/container/container.js'

  export default defineComponent({
    mixins: [tableMixin],
    mounted() {
      this.disabledTableSettingDefaultBehavior()
    },
    setup() {
      const $this = getCurrentInstance()

      const tableRef = ref(null)

      const requestIds = {
        property: Symbol(),
        list: Symbol()
      }

      const route = computed(() => RouterQuery.route)

      const MODEL_ID_KEY = CONTAINER_OBJECT_INST_KEYS[CONTAINER_OBJECTS.CONTAINER].ID
      const MODEL_NAME_KEY = CONTAINER_OBJECT_INST_KEYS[CONTAINER_OBJECTS.CONTAINER].NAME
      const MODEL_FULL_NAME_KEY = CONTAINER_OBJECT_INST_KEYS[CONTAINER_OBJECTS.CONTAINER].FULL_NAME

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
          MODEL_NAME_KEY,
          'container_uid'
        ]
      })

      const properties = ref([])

      const columnsConfigKey = 'pod_container_custom_table_columns'

      const bizId = computed(() => store.getters['objectBiz/bizId'])
      const podId = computed(() => parseInt(route.value.params.podId, 10))

      const filter = reactive({
        field: MODEL_FULL_NAME_KEY,
        value: '',
        operator: '$regex'
      })

      // 在模型管理中可以配置展示列（全局的）
      const globalUsercustom = computed(() => store.getters['userCustom/globalUsercustom'])

      const usercustom = computed(() => store.getters['userCustom/usercustom'])
      const customColumns = computed(() => usercustom.value[columnsConfigKey] || [])
      const globalCustomColumns = computed(() =>  globalUsercustom.value?.pod_global_custom_table_columns || [])

      const clipboardList = computed(() => table.header.slice())

      const hasSelection = computed(() => !!table.selection?.length)

      const getList = async () => {
        const params = getSearchParams()
        if (!params.fields?.length) {
          return
        }

        try {
          const { list, count } = await containerConService.find(params, {
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
        const conProperties = await containerPropertyService.getMany({
          objId: CONTAINER_OBJECTS.CONTAINER
        }, {
          requestId: requestIds.property,
          fromCache: true
        })
        properties.value = conProperties

        setTableHeader()

        getList()
      })

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
            rules: [
              {
                field: 'bk_pod_id',
                operator: 'equal',
                value: podId.value
              }
            ]
          }
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

      const handlePageChange = (current = 1) => {
        table.pagination.current = current
        getList()
      }

      const handleLimitChange = (limit) => {
        table.pagination.current = 1
        table.pagination.limit = limit
        getList()
      }

      const handleSortChange = (sort) => {
        table.sort = getSort(sort)
        getList()
      }

      const handleValueClick = (row, column) => {
        if (column.id !== 'id') {
          return
        }
        routerActions.redirect({
          name: MENU_POD_CONTAINER_DETAILS,
          params: {
            bizId: bizId.value,
            podId: podId.value,
            containerId: row.id
          },
          query: {
            tab: 'property'
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

      const handleSearch = async (value) => {
        table.pagination.current = 1
        filter.value = value
        getList()
      }

      const handleCopy = (column) => {
        const copyText = table.selection.map(row => getPropertyCopyValue(row[column.id], column.property))
        $this.proxy.$copyText(copyText.join('\n')).then(() => {
          $success(t('复制成功'))
        }, () => {
          $error(t('复制失败'))
        })
      }

      return {
        requestIds,
        tableRef,
        table,
        filter,
        clipboardList,
        hasSelection,
        handlePageChange,
        handleLimitChange,
        handleSortChange,
        handleValueClick,
        handleSelectionChange,
        handleHeaderClick,
        handleSearch,
        handleCopy
      }
    }
  })
</script>

<template>
  <div class="pod-list">
    <div class="list-options">
      <div class="option">
        <cmdb-clipboard-selector class="options-clipboard"
          :list="clipboardList"
          :disabled="!hasSelection"
          @on-copy="handleCopy">
        </cmdb-clipboard-selector>
      </div>
      <div class="option">
        <bk-input class="filter-fast-search"
          v-model.trim="filter.value"
          :placeholder="$t('请输入名称')"
          @enter="handleSearch">
        </bk-input>
      </div>
    </div>

    <bk-table class="list-table"
      v-bkloading="{ isLoading: $loading(Object.values(requestIds)) }"
      ref="tableRef"
      :data="table.data"
      :pagination="table.pagination"
      :max-height="$APP.height - 325"
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
  height: 100%;
  overflow: auto;
  @include scrollbar-y;

  .list-options {
    display: flex;
    justify-content: space-between;
    margin-top: 14px;

    .option {
      display: flex;
      align-items: center;

      .filter-fast-search {
        width: 300px;
      }
    }
  }

  .list-table {
    margin-top: 14px;
  }
}
</style>
