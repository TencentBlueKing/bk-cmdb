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
  <div class="business-set-layout">
    <div class="business-set-options clearfix">
      <cmdb-auth class="fl" :auth="{ type: $OPERATION.C_BUSINESS_SET }">
        <bk-button slot-scope="{ disabled }"
          class="fl"
          theme="primary"
          :disabled="disabled"
          @click="handleCreate">
          {{$t('新建')}}
        </bk-button>
      </cmdb-auth>
      <div class="options-button fr">
        <icon-button
          icon="icon-cc-setting"
          v-bk-tooltips.top="$t('列表显示属性配置')"
          @click="columnsConfigShow = true">
        </icon-button>
      </div>
      <div class="options-filter clearfix fr">
        <cmdb-property-selector
          class="filter-selector fl"
          v-model="filter.field"
          :properties="properties"
          @change="handleFilterFieldChange">
        </cmdb-property-selector>
        <component class="filter-value fl"
          :is="`cmdb-search-${filterType}`"
          :placeholder="filterPlaceholder"
          :class="filterType"
          :fuzzy="true"
          v-bind="filterComponentProps"
          v-model="filter.value"
          @change="handleFilterValueChange"
          @enter="handleFilterValueEnter"
          @clear="handleFilterValueEnter">
        </component>
      </div>
    </div>
    <bk-table class="business-table"
      v-bkloading="{ isLoading: $loading(requestId) }"
      :data="table.list"
      :pagination="table.pagination"
      :max-height="$APP.height - 200"
      @sort-change="handleSortChange"
      @page-limit-change="handleSizeChange"
      @page-change="handlePageChange">
      <bk-table-column v-for="column in table.header"
        sortable="custom"
        :key="column.id"
        :prop="column.id"
        :label="column.name"
        :min-width="$tools.getHeaderPropertyMinWidth(column.property, { hasSort: true })"
        show-overflow-tooltip>
        <template slot-scope="{ row }">
          <cmdb-property-value
            :theme="column.id === 'bk_biz_set_id' ? 'primary' : 'default'"
            :value="row[column.id]"
            :show-unit="false"
            :property="column.property"
            @click.native.stop="handleValueClick(row, column)">
          </cmdb-property-value>
        </template>
      </bk-table-column>
      <bk-table-column :label="$t('操作')" fixed="right">
        <template slot-scope="{ row }">
          <cmdb-auth class="mr10" :auth="{ type: $OPERATION.R_BIZ_SET_RESOURCE, relation: [row.bk_biz_set_id] }">
            <template slot-scope="{ disabled }">
              <bk-button
                theme="primary"
                :disabled="disabled"
                :text="true"
                @click.stop="handlePreview(row)">
                {{$t('预览')}}
              </bk-button>
            </template>
          </cmdb-auth>
          <cmdb-auth class="mr10" :auth="{ type: $OPERATION.U_BUSINESS_SET, relation: [row.bk_biz_set_id] }">
            <template slot-scope="{ disabled }">
              <bk-button
                theme="primary"
                :disabled="disabled"
                :text="true"
                @click.stop="handleEdit(row)">
                {{$t('编辑')}}
              </bk-button>
            </template>
          </cmdb-auth>
          <cmdb-auth
            :auth="{ type: $OPERATION.D_BUSINESS_SET, relation: [row.bk_biz_set_id] }"
            v-bk-tooltips.top="{ content: $t('内置业务集不可删除'), disabled: !isBuiltin(row) }">
            <template slot-scope="{ disabled }">
              <bk-button
                theme="primary"
                :disabled="isBuiltin(row) || disabled"
                :text="true"
                @click.stop="handleDelete(row)">
                {{$t('删除')}}
              </bk-button>
            </template>
          </cmdb-auth>
        </template>
      </bk-table-column>
      <cmdb-table-empty
        slot="empty"
        :stuff="table.stuff">
        <i18n path="业务集列表提示语" class="table-empty-tips">
          <template #auth><bk-link theme="primary" @click="handleApplyPermission">{{$t('申请查看权限')}}</bk-link></template>
          <template #create>
            <cmdb-auth :auth="{ type: $OPERATION.C_BUSINESS_SET }">
              <bk-button slot-scope="{ disabled }" text
                theme="primary"
                class="text-btn"
                :disabled="disabled"
                @click="handleCreate">
                {{$t('立即创建')}}
              </bk-button>
            </cmdb-auth>
          </template>
        </i18n>
      </cmdb-table-empty>
    </bk-table>

    <management-form
      :properties="properties"
      :property-groups="propertyGroups"
      :show.sync="managementFormState.show"
      :data="managementFormState.data"
      @save-success="handleSaveSuccess" />

    <columns-config
      :show.sync="columnsConfigShow"
      :properties="properties"
      @update-header="handleUpdateHeader" />

    <business-scope-preview v-bind="previewProps" :show.sync="previewProps.show" />

  </div>
</template>

<script>
  import { computed, defineComponent, reactive, ref, watch, watchEffect } from 'vue'
  import { t } from '@/i18n'
  import { OPERATION } from '@/dictionary/iam-auth'
  import { $bkInfo, $success, $error } from '@/magicbox/index.js'
  import cmdbPropertySelector from '@/components/property-selector'
  import managementForm from './children/management-form.vue'
  import businessScopePreview from '@/components/business-scope/preview.vue'
  import columnsConfig from './children/columns-config.vue'
  import RouterQuery from '@/router/query'
  import routerActions from '@/router/actions'
  import Utils from '@/components/filters/utils'
  import { getDefaultPaginationConfig, getSort } from '@/utils/tools.js'
  import applyPermission from '@/utils/apply-permission.js'
  import businessSetService from '@/service/business-set/index.js'
  import propertyService from '@/service/property/property.js'
  import propertyGroupService from '@/service/property/group.js'
  import { MENU_RESOURCE_BUSINESS_SET_DETAILS } from '@/dictionary/menu-symbol.js'

  export default defineComponent({
    components: {
      cmdbPropertySelector,
      columnsConfig,
      managementForm,
      businessScopePreview
    },
    setup() {
      const requestId = Symbol()
      const MODEL_ID_KEY = 'bk_biz_set_id'
      const MODEL_NAME_KEY = 'bk_biz_set_name'

      // 响应式的query
      const query = computed(() => RouterQuery.getAll())

      const table = reactive({
        header: [],
        list: [],
        pagination: {
          count: 0,
          current: 1,
          ...getDefaultPaginationConfig()
        },
        sort: `-${MODEL_ID_KEY}`,
        stuff: {
          type: 'default',
          payload: {
            resource: t('业务集')
          }
        }
      })

      const filter = reactive({
        field: query.value.field || MODEL_NAME_KEY,
        value: '',
        operator: '$regex'
      })

      // 创建/编辑组件状态值
      const managementFormState = reactive({
        show: false,
        data: {}
      })

      // 列配置组件显示状态
      const columnsConfigShow = ref(false)

      const previewProps = reactive({
        show: false,
        mode: 'after',
        payload: {}
      })

      // 计算查询条件参数
      const searchParams = computed(() => {
        const params = {
          fields: [],
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
        // eslint-disable-next-line max-len
        const { conditions, time_condition: timeCondition } = Utils.transformGeneralModelCondition(condition, properties.value)

        if (timeCondition) {
          params.time_condition = timeCondition
        }

        if (conditions) {
          params.bk_biz_set_filter = {
            condition: 'AND',
            rules: conditions.rules
          }
        }

        return params
      })

      const getList = async (options = {}) => {
        try {
          const { list, count } = await businessSetService.find(searchParams.value, {
            requestId,
            cancelPrevious: true,
            globalPermission: false
          })

          if (options.isDel && count && !list?.length) {
            RouterQuery.set({
              page: table.pagination.current - 1,
              _t: Date.now()
            })
            return
          }

          table.list = list
          table.pagination.count = count

          table.stuff.type = filter.value.toString().length ? 'search' : 'default'
        } catch ({ permission }) {
          if (!permission) return
          table.stuff = {
            type: 'permission',
            payload: { permission }
          }
        }
      }

      // 获取模型属性与分组
      const properties = ref([])
      const propertyGroups = ref([])
      watchEffect(async () => {
        const [modelProperties, modelPropertyGroups] = await Promise.all([
          propertyService.findBizSet(true),
          propertyGroupService.findBizSet()
        ])
        properties.value = modelProperties
        propertyGroups.value = modelPropertyGroups

        // 默认使用查询参数更新一次filter数据
        updateFilter(query.value.field, query.value.keyword, query.value.operator)

        // 初始化查询
        getList()
      })

      // 查询条件组件相关属性数据
      const filterProperty = computed(() => properties.value.find(property => property.bk_property_id === filter.field))
      const filterType = computed(() => filterProperty.value?.bk_property_type ?? 'singlechar')
      const filterPlaceholder = computed(() => Utils.getPlaceholder(filterProperty.value))
      const filterComponentProps = computed(() => Utils.getBindProps(filterProperty.value))

      // 更新filter数据，无值状态时则使用默认数据初始化
      const updateFilter = (field, value = '', operator = '') => {
        if (field) {
          filter.field = field
        }

        if (!filterProperty.value) return

        // 业务集中的singlechar类型统一使用$regex
        const options = filterType.value === 'singlechar' ? { operator: '$regex', value: '' } : {}
        const defaultData = { ...Utils.getDefaultData(filterProperty.value), ...options }

        filter.operator = operator || defaultData.operator
        filter.value = formatFilterValue({ value, operator: filter.operator }, defaultData.value)
      }

      // 根据传入的currentValue或operator格式化值，currentValue值为空时使用defaultValue
      const formatFilterValue = ({ value: currentValue, operator }, defaultValue) => {
        let value = currentValue.toString().length ? currentValue : defaultValue
        const isNumber = ['int', 'float'].includes(filterType.value)
        if (isNumber && value) {
          value = parseFloat(value, 10)
        } else if (operator === '$in') {
          // eslint-disable-next-line no-nested-ternary
          value = Array.isArray(value) ? value : !!value ? [value] : []
          if (filterType.value === 'organization') {
            value = value.map(val => Number(val))
          }
        } else if (Array.isArray(value)) {
          value = value.filter(val => !!val)
        }
        return value
      }

      // 切换条件字段时初始化字段对应的默认值
      const handleFilterFieldChange = () => updateFilter()

      // 监听查询参数触发查询
      watch(
        query,
        async (query) => {
          const {
            page = 1,
            limit = table.pagination.limit,
            keyword = '',
            operator = '',
            field = MODEL_NAME_KEY
          } = query
          updateFilter(field, keyword, operator)

          table.pagination.current = parseInt(page, 10)
          table.pagination.limit = parseInt(limit, 10)

          getList()
        }
      )

      const handleUpdateHeader = (header) => {
        table.header = header
      }

      const handleSortChange = (sort) => {
        table.sort = getSort(sort, { prop: MODEL_ID_KEY, order: 'descending' })
        RouterQuery.refresh()
      }

      const handleSizeChange = (size) => {
        RouterQuery.set({
          limit: size,
          page: 1,
          _t: Date.now()
        })
      }

      const handlePageChange = (page) => {
        RouterQuery.set({
          page,
          _t: Date.now()
        })
      }

      const handleDelete = (inst) => {
        $bkInfo({
          title: t('确认要删除', { name: inst[MODEL_NAME_KEY] }),
          confirmLoading: true,
          confirmFn: async () => {
            try {
              await businessSetService.deleteById(inst[MODEL_ID_KEY])
              getList({ isDel: true })
              $success(t('删除成功'))
            } catch (error) {
              console.error(error)
              $error(t('删除失败'))
              return false
            }
          }
        })
      }
      const handlePreview = (inst) => {
        previewProps.show = true
        previewProps.payload = { ...inst }
      }

      const handleEdit = (inst) => {
        managementFormState.show = true
        managementFormState.data = { ...inst }
      }

      const handleCreate = () => {
        managementFormState.show = true
        managementFormState.data = {}
      }

      const handleValueClick = (row, column) => {
        if (column.id !== MODEL_ID_KEY) {
          return false
        }
        routerActions.redirect({
          name: MENU_RESOURCE_BUSINESS_SET_DETAILS,
          params: {
            bizSetId: row[MODEL_ID_KEY]
          },
          history: true
        })
      }

      const handleApplyPermission = async () => {
        try {
          await applyPermission({
            type: OPERATION.R_BUSINESS_SET,
            relation: []
          })
        } catch (e) {
          console.error(e)
        }
      }

      const handleSaveSuccess = () => {
        managementFormState.show = false
        getList()
      }

      const handleFilterValueChange = () => {
        const hasEnterEvnet = ['float', 'int', 'longchar', 'singlechar']
        if (hasEnterEvnet.includes(filterType.value)) return
        handleFilterValueEnter()
      }

      const handleFilterValueEnter = () => {
        RouterQuery.set({
          _t: Date.now(),
          page: 1,
          field: filter.field,
          keyword: filter.value,
          operator: filter.operator
        })
      }

      const isBuiltin = inst => inst?.default === 1

      if (query.value.create) {
        handleCreate()
      }

      return  {
        properties,
        propertyGroups,
        filterType,
        filterPlaceholder,
        filterComponentProps,
        table,
        filter,
        requestId,
        managementFormState,
        columnsConfigShow,
        previewProps,
        isBuiltin,
        handleCreate,
        handleFilterValueChange,
        handleFilterValueEnter,
        handleSortChange,
        handleSizeChange,
        handlePageChange,
        handleDelete,
        handlePreview,
        handleEdit,
        handleApplyPermission,
        handleValueClick,
        handleSaveSuccess,
        handleUpdateHeader,
        handleFilterFieldChange
      }
    }
  })
</script>

<style lang="scss" scoped>
    .business-set-layout {
        padding: 15px 20px 0;
    }
    .options-filter {
        position: relative;
        margin-right: 10px;
        .filter-selector {
            width: 120px;
            border-radius: 2px 0 0 2px;
            margin-right: -1px;
        }
        .filter-value {
            width: 320px;
            border-radius: 0 2px 2px 0;
            /deep/ .bk-form-input {
                border-radius: 0 2px 2px 0;
            }
        }
        .filter-search {
            position: absolute;
            right: 10px;
            top: 9px;
            cursor: pointer;
        }
    }
    .options-button {
        font-size: 0;
        .bk-button {
            width: 32px;
            padding: 0;
            /deep/ .bk-icon {
                line-height: 14px;
            }
        }
    }
    .business-table {
        margin-top: 14px;
    }
    .table-empty-tips {
        display: flex;
        align-items: center;
        justify-content: center;
        .text-btn {
            font-size: 14px;
        }
    }
</style>
