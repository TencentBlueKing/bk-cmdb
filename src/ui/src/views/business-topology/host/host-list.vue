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
    <host-list-options @transfer="handleTransfer" v-test-id></host-list-options>
    <host-filter-tag class="filter-tag" ref="filterTag"></host-filter-tag>
    <bk-table class="host-table" v-test-id.businessHostAndService="'hostList'"
      ref="tableRef"
      v-bkloading="{ isLoading: $loading(Object.values(request)) }"
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
        :show-overflow-tooltip="column.bk_property_type !== 'map'"
        :min-width="getColumnMinWidth(column)"
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
            :multiple="column.bk_obj_id !== 'host'"
            @click.native.stop="handleValueClick(row, column)">
          </cmdb-property-value>
        </template>
      </bk-table-column>
      <bk-table-column type="setting"></bk-table-column>
    </bk-table>
    <cmdb-dialog v-model="dialog.show" :width="dialog.width" :height="dialog.height">
      <component
        :is="dialog.component"
        v-bind="dialog.props"
        @cancel="handleDialogCancel"
        @confirm="handleDialogConfirm">
      </component>
    </cmdb-dialog>
  </div>
</template>

<script>
  import has from 'has'
  import HostListOptions from './host-list-options.vue'
  import ModuleSelector from './module-selector.vue'
  import AcrossBusinessConfirm from './across-business-confirm.vue'
  import AcrossBusinessModuleSelector from './across-business-module-selector.vue'
  import MoveToResourceConfirm from './move-to-resource-confirm.vue'
  import hostValueFilter from '@/filters/host'
  import tableMixin from '@/mixins/table'
  import { mapGetters } from 'vuex'
  import {
    MENU_BUSINESS_HOST_AND_SERVICE,
    MENU_BUSINESS_HOST_DETAILS,
    MENU_BUSINESS_TRANSFER_HOST
  } from '@/dictionary/menu-symbol'
  import Bus from '@/utils/bus.js'
  import RouterQuery from '@/router/query'
  import HostFilterTag from '@/components/filters/filter-tag'
  import FilterUtils from '@/components/filters/utils'
  import FilterStore, { setupFilterStore } from '@/components/filters/store'
  import ColumnsConfig from '@/components/columns-config/columns-config.js'
  import { ONE_TO_ONE } from '@/dictionary/host-transfer-type.js'
  import { BUILTIN_MODELS, BUILTIN_MODEL_PROPERTY_KEYS } from '@/dictionary/model-constants.js'
  import { CONTAINER_OBJECTS, CONTAINER_OBJECT_PROPERTY_KEYS, TOPO_MODE_KEYS } from '@/dictionary/container.js'
  import containerHostService from '@/service/container/host.js'
  import { getContainerNodeType } from '@/service/container/common.js'

  export default {
    components: {
      HostListOptions,
      HostFilterTag,
      [ModuleSelector.name]: ModuleSelector,
      [AcrossBusinessConfirm.name]: AcrossBusinessConfirm,
      [AcrossBusinessModuleSelector.name]: AcrossBusinessModuleSelector,
      [MoveToResourceConfirm.name]: MoveToResourceConfirm
    },
    filters: {
      hostValueFilter
    },
    mixins: [tableMixin],
    props: {
      active: Boolean
    },
    data() {
      return {
        commonRequestFinished: false,
        table: {
          data: [],
          selection: [],
          sort: 'bk_host_id',
          pagination: this.$tools.getDefaultPaginationConfig()
        },
        dialog: {
          width: 830,
          height: 600,
          show: false,
          component: null,
          props: {}
        },
        request: {
          table: Symbol('table'),
          moveToResource: Symbol('moveToResource'),
          moveToIdleModule: Symbol('moveToIdleModule')
        },
        filtersTagHeight: 0,
        tableHeader: []
      }
    },
    computed: {
      ...mapGetters('objectBiz', ['bizId', 'currentBusiness']),
      ...mapGetters('objectModelClassify', ['getModelById']),
      ...mapGetters('businessHost', [
        'columnsConfigProperties',
        'selectedNode',
        'commonRequest'
      ]),
      isContainerNode() {
        return !!this.selectedNode?.data?.is_container
      },
      isBizNode() {
        return this.selectedNode?.data?.bk_obj_id === BUILTIN_MODELS.BUSINESS
      },
      topoMode() {
        if (this.isContainerNode) {
          return TOPO_MODE_KEYS.CONTAINER
        }
        if (this.isBizNode) {
          return TOPO_MODE_KEYS.BIZ_NODE
        }
        return TOPO_MODE_KEYS.NORMAL
      },
      customInstanceColumnKey() {
        if (this.isContainerHost) {
          return this.$route.meta.customContainerInstanceColumn
        }
        return this.$route.meta.customInstanceColumn
      },
      isContainerSearchMode() {
        return FilterStore.isContainerSearchMode
      },
      searchMode() {
        return FilterStore.searchMode
      },
      isContainerHost() {
        return this.isContainerSearchMode || this.isContainerNode
      }
    },
    watch: {
      $route() {
        this.initFilterStore()
      },
      topoMode(mode) {
        FilterStore.setTopoMode(mode)

        this.tableHeader = FilterStore.getHeader()
        // 重置selection防止因数据结构不同导致获取数据错误
        this.table.selection = []
      },
      searchMode() {
        this.tableHeader = FilterStore.getHeader()
      }
    },
    created() {
      FilterStore.setTopoMode(this.topoMode)

      this.initFilterStore()

      this.unwatchRouter = RouterQuery.watch('*', ({
        tab = 'hostList',
        node,
        page = 1,
        limit = this.table.pagination.limit
      }) => {
        if (this.$route.name !== MENU_BUSINESS_HOST_AND_SERVICE) {
          return false
        }
        this.table.pagination.current = parseInt(page, 10)
        this.table.pagination.limit = parseInt(limit, 10)

        if (tab === 'hostList' && node && this.selectedNode) {
          this.getHostList()
        }
      }, { throttle: 16, ignore: ['keyword'] })
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
      async initFilterStore() {
        const currentRouteName = this.$route.name
        if (this.storageRouteName === currentRouteName || currentRouteName !== MENU_BUSINESS_HOST_AND_SERVICE) {
          return
        }
        this.storageRouteName = currentRouteName

        await setupFilterStore({
          bk_biz_id: this.bizId,
          header: {
            custom: this.$route.meta.customInstanceColumn,
            customContainer: this.$route.meta.customContainerInstanceColumn,
            global: 'host_global_custom_table_columns'
          }
        })

        this.tableHeader = FilterStore.getHeader()
      },
      getColumnSortable(column) {
        const isHostProperty = column.bk_obj_id === 'host'
        const isForeignKey = column.bk_property_type === 'foreignkey'
        return (isHostProperty && !isForeignKey) ? 'custom' : false
      },
      renderHeader(property) {
        const content = [this.$tools.getHeaderPropertyName(property)]
        const modelId = property.bk_obj_id
        if (modelId !== 'host' && modelId !== CONTAINER_OBJECTS.NODE) {
          const model = this.getModelById(modelId)
          const suffix = this.$createElement('span', { style: { color: '#979BA5', marginLeft: '4px' } }, [`(${model.bk_obj_name})`])
          content.push(suffix)
        }
        return this.$createElement('span', {}, content)
      },
      getColumnMinWidth(property) {
        let name = this.$tools.getHeaderPropertyName(property)
        const modelId = property.bk_obj_id
        if (modelId !== 'host' && modelId !== CONTAINER_OBJECTS.NODE) {
          const model = this.getModelById(modelId)
          name = `${name}(${model.bk_obj_name})`
        }
        return this.$tools.getHeaderPropertyMinWidth(property, { name, hasSort: this.getColumnSortable(property) })
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
        if (column.bk_obj_id !== 'host' || column.bk_property_id !== 'bk_host_id') {
          return
        }
        this.$routerActions.redirect({
          name: MENU_BUSINESS_HOST_DETAILS,
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
      },
      async getHostList() {
        try {
          await this.commonRequest
          this.commonRequestFinished = true

          const result = await this.getSearchRequest()

          this.table.data = result.info || []
          this.table.pagination.count = result.count
        } catch (e) {
          console.error(e)
          this.table.data = []
          this.table.pagination.count = 0
        }
      },
      getSearchRequest() {
        const params = this.getParams()
        const config = {
          requestId: this.request.table,
          cancelPrevious: true
        }

        if (this.isContainerHost) {
          return containerHostService.findAll(params, config)
        }

        return this.$store.dispatch('hostSearch/searchHost', { params, config })
      },
      getParams() {
        const type = this.topoMode
        const paramsMap = {
          normal: this.getNormalParams,
          container: this.getContainerParams,
          bizNode: this.getBizNodeParams
        }

        return paramsMap[type]()
      },
      getNormalParams() {
        const params = {
          ...FilterStore.getSearchParams(),
          page: {
            ...this.$tools.getPageParams(this.table.pagination),
            sort: this.table.sort
          }
        }
        const topoNodeData = this.selectedNode.data
        const fieldMap = {
          biz: 'bk_biz_id',
          set: 'bk_set_id',
          module: 'bk_module_id'
        }
        const topoCondition = {
          field: fieldMap[topoNodeData.bk_obj_id] || 'bk_inst_id',
          operator: '$eq',
          value: topoNodeData.bk_inst_id
        }
        const modelConditionId = has(fieldMap, topoNodeData.bk_obj_id) ? topoNodeData.bk_obj_id : 'object'
        const modelCondition = params.condition.find(modelCondition => modelCondition.bk_obj_id === modelConditionId)
        modelCondition.condition.push(topoCondition)
        return params
      },
      getContainerParams() {
        const params = {
          ...FilterStore.getSearchParams(),
          page: {
            ...this.$tools.getPageParams(this.table.pagination),
            sort: this.table.sort
          }
        }

        const selectedNodeData = this.selectedNode.data

        // 容器节点的属性ID
        const fieldMap = {
          [CONTAINER_OBJECTS.CLUSTER]: CONTAINER_OBJECT_PROPERTY_KEYS[CONTAINER_OBJECTS.CLUSTER].ID,
          [CONTAINER_OBJECTS.FOLDER]: CONTAINER_OBJECT_PROPERTY_KEYS[CONTAINER_OBJECTS.FOLDER].ID,
          [CONTAINER_OBJECTS.NAMESPACE]: CONTAINER_OBJECT_PROPERTY_KEYS[CONTAINER_OBJECTS.NAMESPACE].ID,
          [CONTAINER_OBJECTS.WORKLOAD]: CONTAINER_OBJECT_PROPERTY_KEYS[CONTAINER_OBJECTS.WORKLOAD].ID,
          [BUILTIN_MODELS.BUSINESS]: BUILTIN_MODEL_PROPERTY_KEYS[BUILTIN_MODELS.BUSINESS].ID,
        }
        const nodeType = getContainerNodeType(selectedNodeData.bk_obj_id)

        // folder节点参数特殊处理
        if (nodeType === CONTAINER_OBJECTS.FOLDER) {
          params.folder = true
          // folder父节点为cluster节点
          params[fieldMap[CONTAINER_OBJECTS.CLUSTER]] = this.selectedNode.parent.data.bk_inst_id
        } else {
          // 添加节点的属性ID参数，如 bk_namespace_id
          params[fieldMap[nodeType]] = selectedNodeData.bk_inst_id
        }

        // 节点的类型值，workload节点时为具体的类型，如 daemonSet
        params.kind = selectedNodeData.bk_obj_id

        return params
      },
      getBizNodeParams() {
        if (this.isContainerSearchMode) {
          return this.getContainerParams()
        }
        return this.getNormalParams()
      },
      handleTransfer(type) {
        const actionMap = {
          idle: this.openModuleSelector,
          business: this.openModuleSelector,
          acrossBusiness: this.openAcrossBusiness,
          resource: this.openResourceConfirm,
          increment: this.openModuleSelector
        }
        actionMap[type] && actionMap[type](type)
      },
      openModuleSelector(type) {
        const props = {
          moduleType: type === 'increment' ? 'business' : type,
          transferType: type,
          business: this.currentBusiness
        }
        if (type === 'idle') {
          props.title = this.$t('转移主机到空闲模块', { idleSet: this.$store.state.globalConfig.config.set })
        } else if (type === 'increment') {
          props.title = this.$t('追加主机到业务模块')
        } else {
          props.title = this.$t('转移主机到业务模块')
          const { selection } = this.table
          const firstSelectionModules = selection[0].module.map(module => module.bk_module_id).sort()
          const firstSelectionModulesStr = firstSelectionModules.join(',')
          const allSame = selection.slice(1).every((item) => {
            const modules = item.module.map(module => module.bk_module_id).sort()
              .join(',')
            return modules === firstSelectionModulesStr
          })
          if (allSame) {
            props.previousModules = firstSelectionModules
          }
        }
        this.dialog.props = props
        this.dialog.width = 830
        this.dialog.height = 600
        this.dialog.component = ModuleSelector.name
        this.dialog.show = true
      },
      openResourceConfirm() {
        const invalidList = this.validteIdleHost()
        if (!invalidList) return
        this.dialog.props = {
          count: this.table.selection.length,
          bizId: this.bizId,
          invalidList
        }
        const hasInvalid = !!invalidList.length
        this.dialog.width = hasInvalid ? 640 : 460
        this.dialog.height = hasInvalid ? undefined : 250
        this.dialog.component = MoveToResourceConfirm.name
        this.dialog.show = true
      },
      openAcrossBusiness() {
        const invalidList = this.validteIdleHost()
        if (!invalidList) return
        if (invalidList.length) {
          this.openAcrossBusinessConfirm(invalidList)
        } else {
          this.openAcrossBusinessModuleSelector()
        }
      },
      openAcrossBusinessConfirm(invalidList) {
        this.dialog.props = {
          invalidList,
          count: this.table.selection.length
        }
        this.dialog.width = 640
        this.dialog.height = undefined
        this.dialog.component = AcrossBusinessConfirm.name
        this.dialog.show = true
      },
      openAcrossBusinessModuleSelector() {
        this.dialog.props = {
          title: this.$t('转移主机到其他业务'),
          business: this.currentBusiness,
          type: ONE_TO_ONE
        }
        this.dialog.width = 830
        this.dialog.height = 600
        this.dialog.component = AcrossBusinessModuleSelector.name
        this.dialog.show = true
      },
      validteIdleHost() {
        const invalidList = this.table.selection.filter((item) => {
          const [module] = item.module
          // 非空闲机池
          return module.default === 0
        }).map(item => item.host.bk_host_innerip)

        if (invalidList.length === this.table.selection.length) {
          this.$warn(this.$t('主机不属于空闲机池提示', { idleSet: this.$store.state.globalConfig.config.set }))
          return false
        }

        return invalidList
      },
      handleDialogCancel() {
        this.dialog.show = false
      },
      handleDialogConfirm() {
        this.dialog.show = false
        const type = this.dialog.props.transferType || this.dialog.props.moduleType
        if (this.dialog.component === ModuleSelector.name) {
          if (type === 'idle') {
            const isAllIdleSetHost = this.table.selection.every((data) => {
              const modules = data.module
              return modules.every(module => module.default !== 0)
            })
            if (isAllIdleSetHost) {
              // eslint-disable-next-line prefer-rest-params
              this.transferDirectly(...arguments)
              return
            }
            // eslint-disable-next-line prefer-rest-params
            this.gotoTransferPage(...arguments)
            return
          }
          // eslint-disable-next-line prefer-rest-params
          this.gotoTransferPage(...arguments)
          return
        }
        if (this.dialog.component === MoveToResourceConfirm.name) {
          // eslint-disable-next-line prefer-rest-params
          return this.moveHostToResource(...arguments)
        }
        if (this.dialog.component === AcrossBusinessModuleSelector.name) {
          // eslint-disable-next-line prefer-rest-params
          return this.moveHostToOtherBusiness(...arguments)
        }
        if (this.dialog.component === AcrossBusinessConfirm.name) {
          return this.openAcrossBusinessModuleSelector()
        }
      },
      async transferDirectly(modules) {
        try {
          // eslint-disable-next-line prefer-destructuring
          const internalModule = modules[0]
          await this.$http.post(`host/transfer_with_auto_clear_service_instance/bk_biz_id/${this.bizId}`, {
            bk_host_ids: FilterUtils.getSelectedHostIds(this.table.selection),
            default_internal_module: internalModule.data.bk_inst_id,
            is_remove_from_all: true
          }, {
            requestId: this.request.moveToIdleModule
          })
          Bus.$emit('refresh-count', {
            hosts: [...this.table.selection],
            target: internalModule
          })
          this.table.selection = []
          this.$success('转移成功')
          RouterQuery.set({
            _t: Date.now(),
            page: 1
          })
        } catch (e) {
          console.error(e)
        }
      },
      gotoTransferPage(modules) {
        const query = {
          sourceModel: this.selectedNode.data.bk_obj_id,
          sourceId: this.selectedNode.data.bk_inst_id,
          targetModules: modules.map(node => node.data.bk_inst_id).join(','),
          resources: FilterUtils.getSelectedHostIds(this.table.selection)?.join(','),
          node: this.selectedNode.id
        }
        this.$routerActions.redirect({
          name: MENU_BUSINESS_TRANSFER_HOST,
          params: {
            type: this.dialog.props.transferType || this.dialog.props.moduleType
          },
          query,
          history: true
        })
      },
      async moveHostToResource(directoryId) {
        try {
          const validList = this.table.selection.filter((item) => {
            const [module] = item.module
            return module.default >= 1
          })
          await this.$store.dispatch('hostRelation/transferHostToResourceModule', {
            params: {
              bk_biz_id: this.bizId,
              bk_host_id: validList.map(item => item.host.bk_host_id),
              bk_module_id: directoryId
            },
            config: {
              requestId: this.request.moveToResource
            }
          })
          this.refreshHost()
        } catch (e) {
          console.error(e)
        }
      },
      async moveHostToOtherBusiness(modules, targetBizId) {
        try {
          const [targetModule] = modules
          const validList = this.table.selection.filter((item) => {
            const [module] = item.module
            return module.default >= 1
          })
          await this.$http.post('hosts/modules/across/biz', {
            src_bk_biz_id: this.bizId,
            dst_bk_biz_id: targetBizId,
            bk_host_id: validList.map(({ host }) => host.bk_host_id),
            bk_module_id: targetModule.data.bk_inst_id
          })
          this.refreshHost()
        } catch (error) {
          console.error(error)
        }
      },
      refreshHost() {
        Bus.$emit('refresh-count', {
          hosts: [...this.table.selection]
        })
        this.table.selection = []
        this.$success('转移成功')
        RouterQuery.set({
          _t: Date.now(),
          page: 1
        })
      },
      doLayoutTable() {
        this.$refs?.table?.doLayout()
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
