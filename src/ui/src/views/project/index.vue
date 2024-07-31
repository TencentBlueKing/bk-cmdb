<template>
  <div class="project-layout">
    <div class="project-options clearfix">
      <cmdb-auth class="fl" :auth="{ type: $OPERATION.C_PROJECT }">
        <bk-button slot-scope="{ disabled }"
          class="fl"
          theme="primary"
          :disabled="disabled"
          @click="handleCreate">
          {{$t('新建')}}
        </bk-button>
      </cmdb-auth>
      <cmdb-auth :auth="batchUpdateAuth">
        <template #default="{ disabled }">
          <bk-button
            class="ml10"
            :disabled="selectedRows.length === 0 || disabled"
            @click="handleBatchEdit">
            {{ $t("批量编辑") }}
          </bk-button>
        </template>
      </cmdb-auth>
      <cmdb-button-group
        class="mr10"
        :buttons="buttons"
        :expand="false">
      </cmdb-button-group>
      <div class="options-button fr">
        <icon-button
          icon="icon-cc-setting"
          v-bk-tooltips.top="$t('列表显示属性配置')"
          @click="columnsConfig.show = true"
        >
        </icon-button>
      </div>
      <div class="options-filter clearfix fr">
        <cmdb-property-selector
          class="filter-selector fl"
          v-model="filter.field"
          :properties="fastSearchProperties">
        </cmdb-property-selector>
        <component class="filter-value fl r0"
          :is="`cmdb-search-${filterType}`"
          :placeholder="filterPlaceholder"
          :property="filterProperty"
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
    <bk-table class="project-table"
      v-bkloading="{ isLoading: $loading(requestId.searchProject) }"
      :data="table.visibleList"
      :pagination="table.pagination"
      :max-height="$APP.height - 200"
      @sort-change="handleSortChange"
      @page-limit-change="handleSizeChange"
      @page-change="handlePageChange">
      <batch-selection-column
        width="60px"
        row-key="id"
        ref="batchSelectionColumn"
        indeterminate
        :cross-page="table.visibleList.length >= table.pagination.limit"
        :selected-rows.sync="selectedRows"
        :data="table.visibleList"
        :full-data="table.list"
      >
      </batch-selection-column>
      <bk-table-column v-for="column in table.header"
        sortable="custom"
        :key="column.id"
        :prop="column.id"
        :label="column.name"
        :min-width="$tools.getHeaderPropertyMinWidth(column.property, { hasSort: true })"
        :show-overflow-tooltip="$tools.isShowOverflowTips(column.property)">
        <template slot-scope="{ row }">
          <cmdb-property-value
            v-if="customizeContent(column.id)"
            :theme="column.id === 'id' ? 'primary' : 'default'"
            :value="row[column.id]"
            :show-unit="false"
            :property="column.property"
            @click.native.stop="handleValueClick(row, column)">
          </cmdb-property-value>
          <InstanceStatusColumn :value="row[column.id]" v-else-if="column.id === 'bk_status'"></InstanceStatusColumn>
        </template>
      </bk-table-column>
      <bk-table-column :label="$t('操作')" fixed="right">
        <template slot-scope="{ row }">
          <cmdb-auth @click.native.stop :auth="{ type: $OPERATION.U_PROJECT, relation: [row.id] }">
            <template slot-scope="{ disabled }">
              <bk-popconfirm
                v-if="row.bk_status === 'disabled'"
                :content="$t('启用操作提示语')"
                width="288"
                trigger="click"
                @confirm="handleConfirm('enable', row)">
                <bk-button
                  theme="primary"
                  :disabled="disabled"
                  :text="true">
                  {{$t('启用')}}
                </bk-button>
              </bk-popconfirm>
              <bk-popconfirm
                v-else
                :content="$t('停用操作提示语')"
                width="288"
                trigger="click"
                @confirm="handleConfirm('disabled', row)">
                <bk-button
                  theme="primary"
                  :disabled="disabled"
                  :text="true">
                  {{$t('停用')}}
                </bk-button>
              </bk-popconfirm>
            </template>
          </cmdb-auth>
        </template>
      </bk-table-column>

      <cmdb-table-empty
        slot="empty"
        :auth="{ type: $OPERATION.C_PROJECT, relation: [model.id] }"
        :stuff="table.stuff"
        @create="handleCreate"
        @clear="handleClearFilter">
      </cmdb-table-empty>
    </bk-table>

    <bk-sideslider
      v-transfer-dom
      :is-show.sync="slider.show"
      :title="slider.title"
      :width="800"
      :before-close="handleSliderBeforeClose">
      <bk-tab :active.sync="tab.active" type="unborder-card" slot="content" v-if="slider.show">
        <bk-tab-panel name="attribute" :label="$t('属性')" style="width: calc(100% + 40px);margin: 0 -20px;">
          <cmdb-form
            ref="form"
            :properties="properties"
            :property-groups="propertyGroups"
            :inst="attribute.inst.edit"
            :is-main-line="true"
            :type="attribute.type"
            :save-auth="saveAuth"
            @on-submit="handleSave"
            @on-cancel="handleSliderBeforeClose">
          </cmdb-form>
        </bk-tab-panel>
      </bk-tab>
    </bk-sideslider>

    <bk-sideslider
      v-transfer-dom
      :is-show.sync="batchUpdateSlider.show"
      :title="$t('批量修改项目')"
      :width="800"
      :before-close="handleBatchUpdateSliderBeforeClose"
    >
      <bk-tab
        :active.sync="tab.active"
        type="unborder-card"
        slot="content"
        v-if="batchUpdateSlider.show"
      >
        <bk-tab-panel
          name="attribute"
          :label="$t('属性')"
          style="width: calc(100% + 40px); margin: 0 -20px"
        >
          <cmdb-form-multiple
            ref="batchUpdateForm"
            :properties="properties"
            :property-groups="propertyGroups"
            :save-auth="saveAuth"
            :show-default-value="true"
            @on-submit="handleMultipleSave"
            :loading="batchUpdateSlider.loading"
            @on-cancel="handleBatchUpdateSliderBeforeClose"
          >
          </cmdb-form-multiple>
        </bk-tab-panel>
      </bk-tab>
    </bk-sideslider>

    <bk-sideslider
      v-transfer-dom
      :is-show.sync="columnsConfig.show"
      :width="600"
      :title="$t('列表显示属性配置')"
      :before-close="handleColumnsConfigSliderBeforeClose"
    >
      <cmdb-columns-config
        slot="content"
        v-if="columnsConfig.show"
        ref="cmdbColumnsConfig"
        :properties="properties"
        :selected="columnsConfig.selected"
        :disabled-columns="columnsConfig.disabledColumns"
        @on-apply="handleApplayColumnsConfig"
        @on-cancel="handleColumnsConfigSliderBeforeClose"
        @on-reset="handleResetColumnsConfig">
      </cmdb-columns-config>
    </bk-sideslider>
  </div>
</template>

<script>
  import { mapState, mapActions, mapGetters } from 'vuex'
  import RouterQuery from '@/router/query'
  import throttle from 'lodash.throttle'
  import Utils from '@/components/filters/utils'
  import BatchSelectionColumn from '@/components/batch-selection-column'
  import cmdbPropertySelector from '@/components/property-selector'
  import cmdbColumnsConfig from '@/components/columns-config/columns-config.vue'
  import {  MENU_RESOURCE_PROJECT_DETAILS  } from '@/dictionary/menu-symbol'
  import InstanceStatusColumn from './children/instance-status-column.vue'
  import { BUILTIN_MODELS } from '@/dictionary/model-constants.js'
  import projectService from '@/service/project/index.js'
  import { PROPERTY_TYPES } from '@/dictionary/property-constants'
  import cmdbButtonGroup from '@/components/ui/other/button-group'

  export default {
    components: {
      cmdbColumnsConfig,
      cmdbPropertySelector,
      BatchSelectionColumn,
      InstanceStatusColumn,
      cmdbButtonGroup
    },
    data() {
      return {
        table: {
          header: [],
          list: [],
          visibleList: [],
          pagination: {
            count: 0,
            current: 1,
            ...this.$tools.getDefaultPaginationConfig()
          },
          defaultSort: 'id',
          sort: 'id',
          stuff: {
            type: 'default',
            payload: {
              resource: this.$t('项目')
            }
          }
        },
        selectedRows: [],
        properties: [],
        propertyGroups: [],
        filter: {
          field: 'bk_project_name',
          value: '',
          operator: ''
        },
        columnsConfig: {
          show: false,
          selected: [],
          disabledColumns: ['id', 'bk_project_name']
        },
        columnsConfigKey: 'pro_custom_table_columns',
        attribute: {
          type: null,
          inst: {
            edit: {},
            details: {}
          }
        },
        tab: {
          active: 'attribute'
        },
        slider: {
          show: false,
          title: ''
        },
        batchUpdateSlider: {
          show: false,
          loading: false,
        },
        requestId: {
          searchObjectAttribute: Symbol(),
          searchProject: Symbol(),
          searchGroup: Symbol()
        }
      }
    },
    computed: {
      ...mapState('userCustom', ['globalUsercustom']),
      ...mapGetters('userCustom', ['usercustom']),
      ...mapGetters('objectModelClassify', ['getModelById']),
      filterProperty() {
        const property = this.properties.find(property => property.bk_property_id === this.filter.field)
        return property || null
      },
      filterPlaceholder() {
        return Utils.getPlaceholder(this.filterProperty)
      },
      customProjectColumns() {
        return this.usercustom[this.columnsConfigKey] || []
      },
      filterType() {
        if (this.filterProperty) {
          return this.filterProperty.bk_property_type
        }
        return 'singlechar'
      },
      filterComponentProps() {
        return Utils.getBindProps(this.filterProperty)
      },
      model() {
        return this.getModelById(BUILTIN_MODELS.PROJECT) || {}
      },
      saveAuth() {
        const { type } = this.attribute
        if (type === 'create') {
          return { type: this.$OPERATION.C_PROJECT }
        }
        return null
      },
      fastSearchProperties() {
        return this.properties.filter(item => item.bk_property_type !== PROPERTY_TYPES.INNER_TABLE)
      },
      batchUpdateAuth() {
        if (!this.selectedRows.length) {
          return null
        }
        return this.selectedRows.map(item => ({
          type: this.$OPERATION.U_PROJECT,
          relation: [parseInt(item.id, 10)]
        }))
      },
      buttons() {
        const buttonConfig = [{
          id: 'export',
          text: this.$t('导出选中'),
          handler: this.exportField,
          disabled: !this.selectedRows.length
        }, {
          id: 'batchExport',
          text: this.$t('导出全部'),
          handler: () => this.exportField('all'),
          disabled: !this.table.pagination.count
        }]
        return buttonConfig
      },
    },
    watch: {
      'filter.field'() {
        this.genFilterCondition()
      },
      'slider.show'(show) {
        if (!show) {
          this.tab.active = 'attribute'
        }
      },
      customProjectColumns() {
        this.setTableHeader()
      }
    },
    async created() {
      try {
        this.properties = await this.searchObjectAttribute({
          injectId: BUILTIN_MODELS.PROJECT,
          params: {
            bk_obj_id: BUILTIN_MODELS.PROJECT,
            bk_supplier_account: this.supplierAccount
          },
          config: {
            requestId: this.requestId.searchObjectAttribute,
            fromCache: true
          }
        })
        await Promise.all([
          this.getPropertyGroups(),
          this.setTableHeader()
        ])
        this.throttleGetTableData = throttle(this.getTableData, 300, { leading: false, trailing: true })
        this.unwatch = RouterQuery.watch('*', async ({
          page = 1,
          limit = this.table.pagination.limit,
          filter = '',
          operator = '',
          field = 'bk_project_name'
        }) => {
          this.filter.field = field
          this.table.pagination.current = parseInt(page, 10)
          this.table.pagination.limit = parseInt(limit, 10)
          await this.$nextTick()
          this.genFilterCondition(filter, operator)
          this.throttleGetTableData()
        }, { immediate: true })
      } catch (e) {
        // ignore
      }
      if (this.$route.query.create) {
        this.handleCreate()
      }
      this.properties = this.properties.filter(item => item.bk_property_id !== 'bk_project_icon')
    },
    beforeDestroy() {
      this.unwatch()
    },
    methods: {
      ...mapActions('objectModelFieldGroup', ['searchGroup']),
      ...mapActions('objectModelProperty', ['searchObjectAttribute']),

      async exportField(type = 'select') {
        const useExport = await import('@/components/export-file')
        const title = type === 'select' ? '导出选中' : '导出全部'
        const count = type === 'select' ? this.selectedRows.length : this.table.pagination.count

        useExport.default({
          title: this.$t(title),
          bk_obj_id: BUILTIN_MODELS.PROJECT,
          defaultSelectedFields: this.table.header.map(item => item.id),
          count,
          steps: [{ title: this.$t('选择字段'), icon: 1 }],
          submit: (state, task) => {
            const { fields } = state
            const params = {
              export_custom_fields: fields.value.map(property => property.bk_property_id),
            }
            if (type === 'select') {
              const selected = this.selectedRows.map(e => e.id)
              params.ids = selected
              params.export_condition = {
                page: {
                  start: 0,
                  limit: selected.length,
                  sort: this.table.sort
                }
              }
            }
            if (type === 'all') {
              const {
                conditions,
                time_condition: timeCondition
              } = this.getCondition()
              params.export_condition = {
                filter: conditions,
                time_condition: timeCondition,
                page: {
                  ...task.current.value.page,
                  sort: this.table.sort
                }
              }
            }

            return this.$http.download({
              url: `${window.API_HOST}project/export`,
              method: 'post',
              name: task.current.value.name,
              data: params
            })
          }
        }).show()
      },
      async getTableData() {
        try {
          const [{ count }, { info }] = await Promise.all([
            this.getProjectList('count', { cancelPrevious: true, globalPermission: false }),
            this.getProjectList('filed', { cancelPrevious: true, globalPermission: false })
          ])
          this.table.pagination.count = count
          this.table.list = info
          this.table.stuff.type = this.filter.value.toString().length ? 'search' : 'default'
          this.renderVisibleList()
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
      getCondition() {
        // 这里先直接复用转换通用模型实例查询条件的方法
        const condition = {
          [this.filterProperty.id]: {
            value: this.filter.value,
            operator: this.filter.operator
          }
        }
        const {
          conditions,
          time_condition
        } = Utils.transformGeneralModelCondition(condition, this.properties)
        return { conditions, time_condition }
      },
      getProjectList(type, config = { cancelPrevious: true }) {
        const {
          conditions,
          time_condition: timeCondition
        } = this.getCondition()

        const params = this.getSearchParams(type)
        if (conditions) {
          params.filter = conditions
        }
        if (timeCondition) {
          params.time_condition = timeCondition
        }

        return projectService.find({
          params,
          config: Object.assign({ requestId: this.requestId.searchProject }, config)
        })
      },
      getSearchParams(type) {
        const searchParams = {
          page: {
            start: 0,
            limit: 200,
            sort: this.table.sort,
            enable_count: false
          }
        }
        const countParams = {
          page: {
            enable_count: true
          }
        }
        const params = type === 'filed' ? searchParams : countParams
        return params
      },
      setTableHeader() {
        return new Promise((resolve) => {
          // eslint-disable-next-line max-len
          const customColumns = this.customProjectColumns.length ? this.customProjectColumns : this.globalCustomColumns
          // eslint-disable-next-line max-len
          const headerProperties = this.$tools.getHeaderProperties(this.properties, customColumns, this.columnsConfig.disabledColumns)
          resolve(headerProperties)
        }).then((properties) => {
          this.updateTableHeader(properties)
          this.columnsConfig.selected = properties.map(property => property.bk_property_id)
        })
      },
      updateTableHeader(properties) {
        this.table.header = properties.map(property => ({
          id: property.bk_property_id,
          name: this.$tools.getHeaderPropertyName(property),
          property
        }))
      },
      renderVisibleList() {
        const { limit, current } = this.table.pagination
        this.table.visibleList = this.table.list.slice((current - 1) * limit, current * limit)
      },
      handleSortChange(sort) {
        this.table.sort = this.$tools.getSort(sort)
        this.handlePageChange(1)
        this.getTableData()
      },
      handleSizeChange(size) {
        this.table.pagination.limit = size
        this.handlePageChange(1)
        this.renderVisibleList()
      },
      handlePageChange(page) {
        this.table.pagination.current = page
        this.renderVisibleList()
      },
      handleApplayColumnsConfig(properties) {
        this.$store.dispatch('userCustom/saveUsercustom', {
          [this.columnsConfigKey]: properties.map(property => property.bk_property_id)
        })
        this.columnsConfig.show = false
      },
      handleResetColumnsConfig() {
        this.$store.dispatch('userCustom/saveUsercustom', {
          [this.columnsConfigKey]: []
        })
        this.columnsConfig.show = false
      },
      handleFilterValueChange() {
        const hasEnterEvnet = ['float', 'int', 'longchar', 'singlechar']
        if (hasEnterEvnet.includes(this.filterType)) return
        this.handleFilterValueEnter()
      },
      handleFilterValueEnter() {
        this.$refs.batchSelectionColumn.clearSelection()
        RouterQuery.set({
          _t: Date.now(),
          page: 1,
          field: this.filter.field,
          filter: this.filter.value,
        })
      },
      genFilterCondition(filter = '', operator = '') {
        const defaultData = Utils.getDefaultData(this.filterProperty)
        const isProName = ['singlechar', 'longchar'].includes(this.filterType)
        if (isProName) {
          this.filter.operator = '$regex'
          this.filter.value = filter || ''
        } else {
          this.filter.operator = operator || defaultData.operator
          this.filter.value = this.formatFilterValue(
            { value: filter, operator: this.filter.operator },
            defaultData.value
          )
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
      handleValueClick(row, column) {
        if (column.id !== 'id') {
          return false
        }
        this.$routerActions.redirect({
          name: MENU_RESOURCE_PROJECT_DETAILS,
          params: {
            projId: row.id
          },
          history: true
        })
      },
      getPropertyGroups() {
        return this.searchGroup({
          objId: BUILTIN_MODELS.PROJECT,
          params: {},
          config: {
            fromCache: true,
            requestId: this.requestId.searchGroup
          }
        }).then((groups) => {
          this.propertyGroups = groups
          return groups
        })
      },
      handleCreate() {
        this.attribute.type = 'create'
        this.attribute.inst.edit = {}
        this.slider.show = true
        this.slider.title = `${this.$t('创建')} ${this.model.bk_obj_name}`
      },
      handleBatchEdit() {
        this.batchUpdateSlider.show = true
      },
      closeCreateSlider() {
        if (this.attribute.type === 'create') {
          this.slider.show = false
        }
      },
      handleBatchUpdateSliderBeforeClose() {
        this.addDoubleConfirm(this.$refs.batchUpdateForm, () => {
          this.batchUpdateSlider.show = false
        })
      },
      handleSliderBeforeClose() {
        this.addDoubleConfirm(this.$refs.form, this.closeCreateSlider)
      },
      addDoubleConfirm(componentRef, confirmCallback) {
        const { changedValues } = componentRef
        if (this.tab.active === 'attribute') {
          if (Object.keys(changedValues).length) {
            componentRef.setChanged(true)
            return componentRef.beforeClose(confirmCallback)
          }

          confirmCallback && confirmCallback()

          return true
        }

        confirmCallback && confirmCallback()

        return true
      },
      handleColumnsConfigSliderBeforeClose() {
        const refColumns = this.$refs.cmdbColumnsConfig
        if (!refColumns) {
          return
        }
        const { columnsChangedValues } = refColumns
        if (columnsChangedValues?.()) {
          refColumns.setChanged(true)
          return refColumns.beforeClose(() => {
            this.columnsConfig.show = false
          })
        }
        this.columnsConfig.show = false
      },
      handleSave(values) {
        const data = {
          data: [values]
        }
        projectService.create(data).then(() => {
          this.getTableData()
          this.closeCreateSlider()
          this.$success(this.$t('创建成功'))
          this.$http.cancel('post_searchrProject_$ne_disabled')
        })
      },
      async handleMultipleSave(changedValues) {
        const includeProjectIds = this.selectedRows.map(r => r.id)
        const params = {
          ids: includeProjectIds,
          data: changedValues
        }
        this.batchUpdateSlider.loading = true
        projectService.update(params)
          .then(() => {
            this.$refs.batchSelectionColumn.clearSelection()
            this.batchUpdateSlider.show = false
            RouterQuery.set({
              _t: Date.now(),
            })
          })
          .catch((err) => {
            console.log(err)
          })
          .finally(() => {
            this.batchUpdateSlider.loading = false
          })
      },
      handleConfirm(status, row) {
        const params = {
          ids: [row.id],
          data: {
            bk_status: status
          }
        }
        projectService.update(params).then(() => {
          this.$bkMessage({
            theme: 'success',
            message: this.$t('操作成功')
          })
          RouterQuery.set({
            _t: Date.now(),
          })
        })
          .catch((err) => {
            console.log(err)
          })
      },
      // 项目状态是否启用
      customizeContent(id) {
        return !['bk_status'].includes(id)
      },
      handleClearFilter() {
        RouterQuery.clear()
      }
    }
  }
</script>

<style lang="scss" scoped>
  .project-layout {
    padding: 15px 20px 0;
  }
  .options-filter {
    position: relative;
    margin-right: 10px;
    width: 439px;

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
  .project-table {
    margin-top: 14px;
  }
  .project-icon{
    width: 32px;
    height: 100%;
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
