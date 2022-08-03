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
  <div class="models-layout">
    <div class="models-options clearfix">
      <div class="options-button clearfix fl">
        <cmdb-auth class="fl mr10" :auth="{ type: $OPERATION.C_INST, relation: [model.id] }">
          <bk-button slot-scope="{ disabled }"
            theme="primary"
            :disabled="disabled"
            @click="handleCreate">
            {{$t('新建')}}
          </bk-button>
        </cmdb-auth>
        <cmdb-auth class="fl mr10"
          :auth="[
            { type: $OPERATION.C_INST, relation: [model.id] }
          ]">
          <bk-button slot-scope="{ disabled }"
            class="models-button"
            :disabled="disabled"
            @click="handleImport">
            {{$t('导入')}}
          </bk-button>
        </cmdb-auth>
        <bk-button class="models-button" theme="default"
          :disabled="!table.checked.length"
          @click="handleExport">
          {{$t('导出')}}
        </bk-button>
        <cmdb-auth class="fl mr10" :auth="batchUpdateAuth">
          <template #default="{ disabled }">
            <bk-button class="models-button"
              :disabled="!table.checked.length || disabled"
              @click="handleMultipleEdit">
              {{$t('批量更新')}}
            </bk-button>
          </template>
        </cmdb-auth>
        <cmdb-auth class="fl mr10" :auth="batchDeleteAuth">
          <template #default="{ disabled }">
            <bk-button class="models-button button-delete"
              hover-theme="danger"
              :disabled="!table.checked.length || disabled"
              @click="handleMultipleDelete">
              {{$t('删除')}}
            </bk-button>
          </template>
        </cmdb-auth>
      </div>
      <div class="options-button fr">
        <icon-button class="option-filter ml5" icon="icon-cc-funnel"
          v-bk-tooltips.top="$t('高级筛选')"
          @click="handleSetFilters">
        </icon-button>
        <icon-button class="ml5"
          v-bk-tooltips="$t('查看删除历史')"
          icon="icon-cc-history"
          @click="routeToHistory">
        </icon-button>
        <icon-button class="ml5"
          v-bk-tooltips="$t('列表显示属性配置')"
          icon="icon-cc-setting"
          @click="columnsConfig.show = true">
        </icon-button>
      </div>
      <div class="options-filter clearfix fr">
        <cmdb-property-selector class="filter-selector"
          v-model="filter.field"
          :properties="properties"
          :loading="$loading([request.properties, request.groups])">
        </cmdb-property-selector>
        <component class="filter-value"
          :is="`cmdb-search-${filterType}`"
          :placeholder="filterPlaceholder"
          :class="filterType"
          :fuzzy="filter.fuzzyQuery"
          v-bind="filterComponentProps"
          v-model="filter.value"
          @change="handleFilterValueChange"
          @enter="handleFilterValueEnter"
          @clear="handleFilterValueEnter">
        </component>
        <bk-checkbox class="filter-exact" size="small"
          v-if="allowFuzzyQuery"
          v-model="filter.fuzzyQuery">
          {{$t('模糊')}}
        </bk-checkbox>
      </div>
    </div>
    <general-model-filter-tag
      class="filter-tag"
      ref="filterTag"
      :filter-selected="filterSelected"
      :filter-condition="filterCondition">
    </general-model-filter-tag>
    <bk-table class="models-table" ref="table"
      v-bkloading="{ isLoading: $loading(request.list) }"
      :data="table.list"
      :pagination="table.pagination"
      :max-height="$APP.height - filterTagHeight - 190"
      @sort-change="handleSortChange"
      @page-limit-change="handleSizeChange"
      @page-change="handlePageChange"
      @selection-change="handleSelectChange">
      <bk-table-column type="selection" width="60" align="center" fixed
        class-name="bk-table-selection">
      </bk-table-column>
      <bk-table-column v-for="column in table.header"
        sortable="custom"
        min-width="80"
        :key="column.id"
        :prop="column.id"
        :label="column.name"
        show-overflow-tooltip>
        <template slot-scope="{ row }">
          <cmdb-property-value
            :theme="column.id === 'bk_inst_id' ? 'primary' : 'default'"
            :show-unit="false"
            :value="row[column.id]"
            :property="column.property"
            @click.native.stop="handleValueClick(row, column)">
          </cmdb-property-value>
        </template>
      </bk-table-column>
      <bk-table-column fixed="right" :label="$t('操作')">
        <template slot-scope="{ row }">
          <cmdb-auth :auth="{ type: $OPERATION.D_INST, relation: [model.id, row.bk_inst_id] }">
            <template slot-scope="{ disabled }">
              <bk-button theme="primary" text :disabled="disabled" @click.stop="handleDelete(row)">
                {{$t('删除')}}
              </bk-button>
            </template>
          </cmdb-auth>
        </template>
      </bk-table-column>
      <cmdb-table-empty
        slot="empty"
        :auth="{ type: $OPERATION.C_INST, relation: [model.id] }"
        :stuff="table.stuff"
        @create="handleCreate">
      </cmdb-table-empty>
    </bk-table>
    <bk-sideslider
      v-transfer-dom
      :is-show.sync="slider.show"
      :title="slider.title"
      :width="800"
      :before-close="handleSliderBeforeClose">
      <template slot="content" v-if="slider.contentShow">
        <cmdb-form v-if="attribute.type === 'create'"
          ref="form"
          :properties="properties"
          :property-groups="propertyGroups"
          :inst="attribute.inst.edit"
          :is-main-line="isMainLineModel"
          :type="attribute.type"
          :save-auth="{ type: $OPERATION.C_INST, relation: [model.id] }"
          @on-submit="handleSave"
          @on-cancel="handleCancel">
        </cmdb-form>
        <cmdb-form-multiple v-else-if="attribute.type === 'multiple'"
          ref="multipleForm"
          :uneditable-properties="['bk_inst_name']"
          :properties="properties"
          :property-groups="propertyGroups"
          :save-auth="saveAuth"
          @on-submit="handleMultipleSave"
          @on-cancel="handleMultipleCancel">
        </cmdb-form-multiple>
      </template>
    </bk-sideslider>
    <bk-sideslider v-transfer-dom :is-show.sync="columnsConfig.show" :width="600" :title="$t('列表显示属性配置')">
      <cmdb-columns-config slot="content"
        v-if="columnsConfig.show"
        :properties="properties"
        :selected="columnsConfig.selected"
        :disabled-columns="columnsConfig.disabledColumns"
        @on-apply="handleApplyColumnsConfig"
        @on-cancel="columnsConfig.show = false"
        @on-reset="handleResetColumnsConfig">
      </cmdb-columns-config>
    </bk-sideslider>
    <bk-sideslider v-transfer-dom
      :show-mask="false"
      :is-show.sync="advancedFilterShow"
      :width="400"
      :title="$t('高级筛选')">
      <general-model-filter-form slot="content"
        v-if="advancedFilterShow"
        :obj-id="objId"
        :filter-selected="filterSelected"
        :filter-condition="filterCondition"
        :properties="properties"
        :property-groups="propertyGroups"
        @close="advancedFilterShow = false">
      </general-model-filter-form>
    </bk-sideslider>
    <router-subview></router-subview>
  </div>
</template>

<script>
  import { mapState, mapGetters, mapActions } from 'vuex'
  import QS from 'qs'
  import cmdbColumnsConfig from '@/components/columns-config/columns-config.vue'
  import generalModelFilterForm from '@/components/filters/general-model-filter-form.vue'
  import generalModelFilterTag from '@/components/filters/general-model-filter-tag.vue'
  import cmdbImport from '@/components/import/import'
  import { MENU_RESOURCE_INSTANCE, MENU_RESOURCE_INSTANCE_DETAILS } from '@/dictionary/menu-symbol'
  import cmdbPropertySelector from '@/components/property-selector'
  import RouterQuery from '@/router/query'
  import Utils from '@/components/filters/utils'
  import throttle from  'lodash.throttle'
  import instanceImportService from '@/service/instance/import'
  import instanceService from '@/service/instance/instance'
  import { resetConditionValue } from '@/components/filters/general-model-filter.js'

  const defaultFastSearch = () => ({
    field: 'bk_inst_name',
    value: [],
    operator: '$in',
    fuzzyQuery: false
  })

  export default {
    components: {
      cmdbColumnsConfig,
      cmdbImport,
      cmdbPropertySelector,
      generalModelFilterForm,
      generalModelFilterTag
    },
    data() {
      return {
        properties: [],
        propertyGroups: [],
        table: {
          checked: [],
          header: [],
          list: [],
          pagination: {
            count: 0,
            current: 1,
            ...this.$tools.getDefaultPaginationConfig()
          },
          defaultSort: 'bk_inst_id',
          sort: 'bk_inst_id',
          stuff: {
            type: 'default',
            payload: {}
          }
        },
        filter: defaultFastSearch(),
        filterSelected: [],
        filterCondition: {},
        advancedFilterShow: false,
        slider: {
          show: false,
          contentShow: false,
          title: ''
        },
        tab: {
          active: 'attribute'
        },
        attribute: {
          type: null,
          inst: {
            details: {},
            edit: {}
          }
        },
        columnsConfig: {
          show: false,
          selected: [],
          disabledColumns: ['bk_inst_id', 'bk_inst_name']
        },
        request: {
          properties: Symbol('properties'),
          groups: Symbol('groups'),
          list: Symbol('list')
        },
        filterTagHeight: 0
      }
    },
    computed: {
      ...mapState('userCustom', ['globalUsercustom']),
      ...mapGetters(['supplierAccount', 'userName']),
      ...mapGetters('userCustom', ['usercustom']),
      ...mapGetters('objectBiz', ['bizId']),
      ...mapGetters('objectModelClassify', ['models', 'getModelById']),
      ...mapGetters('objectMainLineModule', ['isMainLine']),
      objId() {
        return this.$route.params.objId
      },
      model() {
        return this.getModelById(this.objId) || {}
      },
      customConfigKey() {
        return `${this.objId}_custom_table_columns`
      },
      customColumns() {
        return this.usercustom[this.customConfigKey] || []
      },
      globalCustomColumns() {
        return this.globalUsercustom[`${this.objId}_global_custom_table_columns`] || []
      },
      parentLayers() {
        return [{
          resource_id: this.model.id,
          resource_type: 'model'
        }]
      },
      filterProperty() {
        const property = this.properties.find(property => property.bk_property_id === this.filter.field)
        return property || null
      },
      filterType() {
        if (this.filterProperty) {
          return this.filterProperty.bk_property_type
        }
        return 'singlechar'
      },
      filterPlaceholder() {
        return Utils.getPlaceholder(this.filterProperty)
      },
      filterComponentProps() {
        return Utils.getBindProps(this.filterProperty)
      },
      allowFuzzyQuery() {
        return ['singlechar', 'longchar'].includes(this.filterType)
      },
      saveAuth() {
        return this.table.checked.map(instId => ({
          type: this.$OPERATION.U_INST,
          relation: [this.model.id, parseInt(instId, 10)]
        }))
      },
      batchDeleteAuth() {
        return this.table.checked.map(instId => ({
          type: this.$OPERATION.D_INST,
          relation: [this.model.id, parseInt(instId, 10)]
        }))
      },
      batchUpdateAuth() {
        return this.table.checked.map(instId => ({
          type: this.$OPERATION.U_INST,
          relation: [this.model.id, parseInt(instId, 10)]
        }))
      },
      isMainLineModel() {
        return this.isMainLine(this.model)
      }
    },
    watch: {
      '$route.query'() {
        if (this.$route.name !== MENU_RESOURCE_INSTANCE) {
          return
        }
        this.setupFilter()
        this.setDynamicBreadcrumbs()
        this.throttleGetTableData()
        this.updateFilterTagHeight()
      },
      'filter.field'() {
        // 模糊搜索
        if (this.allowFuzzyQuery && this.filter.fuzzyQuery) {
          this.filter.value = ''
          this.filter.operator = '$regex'
          return
        }
        const defaultData = Utils.getDefaultData(this.filterProperty)
        this.filter.value = defaultData.value
        this.filter.operator = defaultData.operator
      },
      'filter.fuzzyQuery'(fuzzy) {
        if (!this.allowFuzzyQuery) return

        if (fuzzy) {
          this.filter.value = ''
          this.filter.operator = '$regex'
          return
        }

        const defaultData = Utils.getDefaultData(this.filterProperty)
        this.filter.value = defaultData.value
        this.filter.operator = defaultData.operator
      },
      'slider.show'(show) {
        if (!show) {
          this.tab.active = 'attribute'
        }
        this.$nextTick(() => {
          this.slider.contentShow = show
        })
      },
      customColumns() {
        this.setTableHeader()
      },
      objId() {
        // 切换模型需要重新获取当前模型的数据
        this.fetchData()
      }
    },
    async created() {
      this.throttleGetTableData = throttle(this.getTableData, 300, { leading: false, trailing: true })
      await this.fetchData()
      this.setupFilter()
      this.setDynamicBreadcrumbs()
      this.throttleGetTableData()
      this.getMainLine()
    },
    mounted() {
      this.updateFilterTagHeight()
    },
    beforeRouteUpdate(to, from, next) {
      this.setDynamicBreadcrumbs()
      next()
    },
    methods: {
      ...mapActions('objectModelFieldGroup', ['searchGroup']),
      ...mapActions('objectMainLineModule', ['searchMainlineObject']),
      ...mapActions('objectModelProperty', ['searchObjectAttribute']),
      ...mapActions('objectCommonInst', [
        'createInst',
        'updateInst',
        'batchUpdateInst',
        'deleteInst',
        'batchDeleteInst'
      ]),
      setDynamicBreadcrumbs() {
        this.$store.commit('setTitle', this.model.bk_obj_name)
      },
      async fetchData() {
        try {
          // 先重置当前数据
          this.resetData()

          const [groups, properties] = await Promise.all([
            this.getPropertyGroups(),
            this.getProperties()
          ])
          this.propertyGroups = groups
          this.properties = properties

          // 确定数据更新后更新表头
          this.setTableHeader()
        } catch (e) {
          console.error(e)
        }
      },
      setupFilter() {
        const {
          page = 1,
          limit = this.table.pagination.limit,
          field = 'bk_inst_name',
          operator,
          filter: value = '',
          fuzzy = 0,
          s: searchType = 'fast'
        } = this.$route.query

        const advQuery = QS.parse(RouterQuery.get('filter_adv'))
        const fastQuery = { field, operator, value }

        // 设置快速搜索的字段与类型
        this.filter.field = fastQuery.field
        this.filter.fuzzyQuery = Boolean(Number(fuzzy))

        // 更新表格分页
        this.table.pagination.current = parseInt(page, 10)
        this.table.pagination.limit = parseInt(limit, 10)

        // 每次先将条件项的值重置，避免在清空后或切换模型时产生遗留数据
        this.filterCondition = resetConditionValue(this.filterCondition, this.filterSelected)

        // 默认以高级搜索的查询query填充条件项与选择项
        try {
          Object.keys(advQuery).forEach((key) => {
            const [id, operator] = key.split('.')
            const property = Utils.findProperty(id, this.properties)
            const value = advQuery[key].toString().split(',')
            if (property && operator && value.length) {
              this.$set(this.filterCondition, property.id, {
                operator: `$${operator}`,
                value: Utils.convertValue(value, `$${operator}`, property)
              })
              // eslint-disable-next-line max-len
              const exist = this.filterSelected.findIndex(item => item.bk_property_id === property.bk_property_id) !== -1
              if (!exist) {
                this.filterSelected.push(property)
              }
            }
          })
        } catch (error) {
          this.$warn(this.$t('解析查询链接出错提示'))
        }

        // 合并快速搜索栏条件
        if (searchType === 'fast') {
          const fastSearchProperty = this.properties.find(property => property.bk_property_id === fastQuery.field)
          // eslint-disable-next-line max-len
          const exist = this.filterSelected.findIndex(item => item.bk_property_id === fastSearchProperty.bk_property_id) !== -1

          // 不存在则添加进选择项列表中
          if (!exist) {
            this.filterSelected.push(fastSearchProperty)
          }

          // 更新或初始化条件项
          const defaultData = Utils.getDefaultData(fastSearchProperty)
          const conditionValue = {
            operator: operator || defaultData.operator,
            value: this.formatFilterValue({ value: fastQuery.value, operator: fastQuery.operator }, defaultData.value)
          }
          this.$set(this.filterCondition, fastSearchProperty.id, conditionValue)

          // 给快速搜索框初始化一个正确的值
          this.filter.value = conditionValue.value
        }

        // 高级搜索时重置快速搜索的数据值，因高级搜索是快速搜索的超集在快速搜索中不能复原
        if (searchType === 'adv') {
          this.resetFastSearch()
        }
      },
      formatFilterValue({ value: currentValue, operator }, defaultValue) {
        let value = currentValue.toString().length ? currentValue : defaultValue
        const isNumber = ['int', 'float'].includes(this.filterType)
        if (isNumber && value) {
          value = parseFloat(value, 10)
        } else if (operator === '$in') {
          // eslint-disable-next-line no-nested-ternary
          value = Array.isArray(value) ? value : !!value ? [value] : []
        } else if (operator === '$regex') {
          value = Array.isArray(value) ? (value[0] || '') : value
        } else if (Array.isArray(value)) {
          value = value.filter(value => !!value)
        }
        return value
      },
      handleFilterValueChange() {
        const hasEnterEvnet = ['float', 'int', 'longchar', 'singlechar']
        if (hasEnterEvnet.includes(this.filterType)) return
        this.handleFilterValueEnter()
      },
      handleFilterValueEnter() {
        const query = {
          _t: Date.now(),
          s: 'fast',
          page: 1,
          field: this.filter.field,
          filter: this.filter.value,
          operator: this.filter.operator,
        }
        if (this.allowFuzzyQuery) {
          query.fuzzy = this.filter.fuzzyQuery ? 1 : 0
        }
        RouterQuery.set(query)
      },
      resetData() {
        this.table = {
          checked: [],
          header: [],
          list: [],
          pagination: {
            count: 0,
            current: 1,
            ...this.$tools.getDefaultPaginationConfig()
          },
          defaultSort: 'bk_inst_id',
          sort: 'bk_inst_id',
          stuff: {
            type: 'default',
            payload: {
              resource: this.model.bk_obj_name
            }
          }
        }

        // 重置筛选项与条件
        this.filterSelected = []
        this.filterCondition = {}
      },
      getProperties() {
        return this.searchObjectAttribute({
          injectId: this.objId,
          params: {
            bk_obj_id: this.objId,
            bk_supplier_account: this.supplierAccount
          },
          config: {
            requestId: this.request.properties,
          }
        })
      },
      getPropertyGroups() {
        return this.searchGroup({
          objId: this.objId,
          params: {},
          config: { requestId: this.request.groups }
        })
      },
      getMainLine() {
        return this.searchMainlineObject({})
      },
      setTableHeader() {
        return new Promise((resolve) => {
          const customColumns = this.customColumns.length ? this.customColumns : this.globalCustomColumns
          // eslint-disable-next-line max-len
          const headerProperties = this.$tools.getHeaderProperties(this.properties, customColumns, this.columnsConfig.disabledColumns)
          resolve(headerProperties)
        }).then((properties) => {
          this.updateTableHeader(properties)

          // 同步更新筛选项与条件
          this.updateFilter(properties)

          this.columnsConfig.selected = properties.map(property => property.bk_property_id)
        })
      },
      updateTableHeader(properties) {
        // 将搜索项追加到表头
        const headerIds = properties.map(item => item.bk_property_id)
        let finalProperties = []

        if (this.filterSelected.length >= properties.length && this.isColumnApply) {
          // 列表项减少
          finalProperties = properties
          this.isColumnApply = false
        } else {
          // 列表项增加
          const newHeaders = this.filterSelected.filter(item => !headerIds.includes(item.bk_property_id))
          finalProperties = properties.concat(newHeaders)
        }

        this.table.header = finalProperties.map(property => ({
          id: property.bk_property_id,
          name: this.$tools.getHeaderPropertyName(property),
          property
        }))
      },
      updateFilter(properties = []) {
        const availableProperties = properties.filter(property => property.bk_obj_id === this.objId)
        availableProperties.forEach((property) => {
          // eslint-disable-next-line max-len
          const exist = this.filterSelected.findIndex(item => item.bk_property_id === property.bk_property_id) !== -1
          if (!exist) {
            const defaultData = Utils.getDefaultData(property)
            this.filterSelected.push(property)
            this.$set(this.filterCondition, property.id, {
              operator: defaultData.operator,
              value: defaultData.value
            })
          }
        })
      },
      handleValueClick(item, column) {
        if (column.id !== 'bk_inst_id') {
          return false
        }
        this.$routerActions.redirect({
          name: MENU_RESOURCE_INSTANCE_DETAILS,
          params: {
            objId: this.objId,
            instId: item.bk_inst_id
          },
          history: true
        })
      },
      handleSortChange(sort) {
        this.table.sort = this.$tools.getSort(sort)
        RouterQuery.refresh()
      },
      handleSizeChange(size) {
        RouterQuery.set({
          limit: size,
          page: 1,
          _t: Date.now()
        })
      },
      handlePageChange(page) {
        RouterQuery.set({
          page,
          _t: Date.now()
        })
      },
      handleSelectChange(selection) {
        this.table.checked = selection.map(row => row.bk_inst_id)
      },
      getInstList(config = { cancelPrevious: true }) {
        return instanceService.find({
          bk_obj_id: this.objId,
          params: {
            ...this.getSearchParams(),
            ...Utils.transformGeneralModelCondition(this.filterCondition, this.filterSelected) || {}
          },
          config: Object.assign({ requestId: this.request.list }, config)
        })
      },
      async getTableData() {
        // 防止切换到子路由产生预期外的请求
        if (this.$route.name !== MENU_RESOURCE_INSTANCE) {
          return
        }

        try {
          const { count, info } = await this.getInstList({ cancelPrevious: true, globalPermission: false })
          if (count && !info.length) {
            RouterQuery.set({
              page: this.table.pagination.current - 1,
              _t: Date.now()
            })
          }
          this.table.list = info
          this.table.pagination.count = count
          this.table.stuff.type = this.$route.query?.s?.length ? 'search' : 'default'
        } catch (err) {
          console.error(err)
          if (err.permission) {
            this.table.stuff = {
              type: 'permission',
              payload: { permission: err.permission }
            }
          }
        }
      },
      getSearchParams() {
        const params = {
          fields: [],
          page: {
            start: this.table.pagination.limit * (this.table.pagination.current - 1),
            limit: this.table.pagination.limit,
            sort: this.table.sort
          }
        }
        return params
      },
      handleCreate() {
        this.attribute.type = 'create'
        this.attribute.inst.edit = {}
        this.slider.show = true
        this.slider.title = `${this.$t('创建')} ${this.model.bk_obj_name}`
      },
      handleDelete(inst) {
        this.$bkInfo({
          title: this.$t('确认要删除', { name: inst.bk_inst_name }),
          confirmFn: () => {
            this.deleteInst({
              objId: this.objId,
              instId: inst.bk_inst_id
            }).then(() => {
              this.slider.show = false
              this.$success(this.$t('删除成功'))
              RouterQuery.refresh()
            })
          }
        })
      },
      handleSave(values, changedValues, originalValues, type) {
        if (type === 'update') {
          this.updateInst({
            objId: this.objId,
            instId: originalValues.bk_inst_id,
            params: values
          }).then(() => {
            this.attribute.inst.details = Object.assign({}, originalValues, values)
            this.handleCancel()
            this.$success(this.$t('修改成功'))
            RouterQuery.refresh()
          })
        } else {
          delete values.bk_inst_id // properties中注入了前端自定义的bk_inst_id属性
          this.createInst({
            params: values,
            objId: this.objId
          }).then(() => {
            RouterQuery.set({
              _t: Date.now(),
              page: 1
            })
            this.handleCancel()
            this.$success(this.$t('创建成功'))
          })
        }
      },
      handleCancel() {
        if (this.attribute.type === 'create') {
          this.slider.show = false
        }
      },
      handleMultipleEdit() {
        this.attribute.type = 'multiple'
        this.slider.title = this.$t('批量更新')
        this.slider.show = true
      },
      handleMultipleSave(values) {
        this.batchUpdateInst({
          objId: this.objId,
          params: {
            update: this.table.checked.map(instId => ({
              datas: values,
              inst_id: instId
            }))
          },
          config: {
            requestId: `${this.objId}BatchUpdate`
          }
        }).then(() => {
          this.$success(this.$t('修改成功'))
          this.slider.show = false
          RouterQuery.set({
            _t: Date.now(),
            page: 1
          })
        })
      },
      handleMultipleCancel() {
        this.slider.show = false
      },
      handleMultipleDelete() {
        this.$bkInfo({
          title: this.$t('确定删除选中的实例'),
          confirmFn: () => {
            this.doBatchDeleteInst()
          }
        })
      },
      doBatchDeleteInst() {
        this.batchDeleteInst({
          objId: this.objId,
          config: {
            data: {
              delete: {
                inst_ids: this.table.checked
              }
            }
          }
        }).then(() => {
          this.$success(this.$t('删除成功'))
          this.table.checked = []
          RouterQuery.refresh()
        })
      },
      handleApplyColumnsConfig(properties) {
        this.isColumnApply = true
        this.$store.dispatch('userCustom/saveUsercustom', {
          [this.customConfigKey]: properties.map(property => property.bk_property_id)
        })
        this.columnsConfig.show = false
      },
      handleResetColumnsConfig() {
        this.isColumnApply = true
        this.$store.dispatch('userCustom/saveUsercustom', {
          [this.customConfigKey]: []
        })
        this.columnsConfig.show = false
      },
      routeToHistory() {
        this.$routerActions.redirect({
          name: 'instanceHistory',
          params: {
            objId: this.objId
          },
          history: true
        })
      },
      handleSliderBeforeClose() {
        if (this.tab.active === 'attribute') {
          const $form = this.attribute.type === 'multiple' ? this.$refs.multipleForm : this.$refs.form
          if ($form.hasChange) {
            return new Promise((resolve) => {
              this.$bkInfo({
                title: this.$t('确认退出'),
                subTitle: this.$t('退出会导致未保存信息丢失'),
                extCls: 'bk-dialog-sub-header-center',
                confirmFn: () => {
                  resolve(true)
                },
                cancelFn: () => {
                  resolve(false)
                }
              })
            })
          }
          return true
        }
        return true
      },
      async handleImport() {
        const useImport = await import('@/components/import-file')
        const [, { show: showImport, setState: setImportState }] = useImport.default()
        setImportState({
          title: this.$t('批量导入'),
          bk_obj_id: this.objId,
          template: `${window.API_HOST}importtemplate/${this.objId}`,
          submit: (options) => {
            const params = {
              op: options.step
            }
            if (options.importRelation) {
              params.object_unique_id = options.object_unique_id
              params.association_condition = options.relations
            }
            return instanceImportService.update({
              file: options.file,
              params,
              config: options.config,
              bk_obj_id: this.objId
            })
          },
          success: () => RouterQuery.set({ _t: Date.now() })
        })
        showImport()
      },
      async handleExport() {
        const useExport = await import('@/components/export-file')
        useExport.default({
          title: this.$t('导出选中'),
          bk_obj_id: this.objId,
          defaultSelectedFields: this.table.header.map(item => item.id),
          count: this.table.checked.length,
          submit: (state, task) => {
            const { fields, exportRelation  } = state
            const params = {
              export_custom_fields: fields.value.map(property => property.bk_property_id),
              bk_inst_ids: this.table.checked
            }
            if (exportRelation.value) {
              params.object_unique_id = state.object_unique_id.value
              params.association_condition = state.relations.value
            }
            return this.$http.download({
              url: `${window.API_HOST}insts/object/${this.objId}/export`,
              method: 'post',
              name: task.current.value.name,
              data: params
            })
          }
        }).show()
      },
      handleSetFilters() {
        this.advancedFilterShow = true
      },
      updateFilterTagHeight() {
        setTimeout(() => {
          const el = this.$refs.filterTag.$el
          if (el?.getBoundingClientRect) {
            this.filterTagHeight = el.getBoundingClientRect().height
          } else {
            this.filterTagHeight = 0
          }
        }, 300)
      },
      resetFastSearch() {
        this.filter = defaultFastSearch()
      }
    }
  }
</script>

<style lang="scss" scoped>
    .models-layout {
        padding: 15px 20px 0;
    }
    .options-filter{
        position: relative;
        margin-right: 5px;
        display: flex;
        align-items: center;
        width: 440px;
        .filter-selector{
            width: 120px;
            border-radius: 2px 0 0 2px;
            margin-right: -1px;
        }
        .filter-value{
            flex: 1;
            border-radius: 0 2px 2px 0;
            &.singlechar,
            &.longchar {
              border-radius: unset;
              /deep/ .bk-tag-input {
                border-radius: unset;
              }
            }
            /deep/ .bk-form-input {
                line-height: 32px;
            }
        }
        .filter-exact {
          display: inline-flex;
          align-items: center;
          padding: 0 5px;
          height: 32px;
          border: 1px solid #c4c6cc;
          border-radius: 0 2px 2px 0;
          border-left: none;
        }
    }
    .models-button{
        display: inline-block;
        position: relative;
        &:hover{
            z-index: 1;
        }
    }
    .models-table{
        margin-top: 14px;
    }
    .filter-tag ~ .models-table {
        margin-top: 0;
    }
</style>
