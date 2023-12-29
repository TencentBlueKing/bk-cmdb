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
  <div class="options-layout clearfix">
    <div class="options-left fl">
      <template v-if="scope === 1">
        <cmdb-auth class="mr10"
          :auth="[
            { type: $OPERATION.C_RESOURCE_HOST, relation: [directoryId] },
            { type: $OPERATION.U_RESOURCE_HOST, relation: [directoryId] }
          ]">
          <bk-button slot-scope="{ disabled }"
            theme="primary"
            style="margin-left: 0"
            :disabled="disabled"
            @click="handleNewImportInst">
            {{$t('导入主机')}}
          </bk-button>
        </cmdb-auth>
      </template>
      <span style="display: inline-block;"
        v-bk-tooltips="{ content: this.$t('仅主机池主机可进行此操作'), disabled: isAllResourceHost }">
        <bk-select
          class="assign-selector mr10"
          font-size="medium"
          :popover-width="180"
          :disabled="!table.checked.length || !isAllResourceHost"
          :clearable="false"
          :placeholder="$t('分配到')"
          v-model="assign.curSelected"
          @selected="handleAssignHosts">
          <bk-option id="-1" :name="$t('分配到')" hidden></bk-option>
          <bk-option id="toBusiness" :name="$t('业务空闲机', { idleSet: $store.state.globalConfig.config.set })"></bk-option>
          <bk-option id="toDirs" :name="$t('主机池其他目录')"></bk-option>
        </bk-select>
      </span>
      <cmdb-transfer-menu class="mr10" v-if="scope !== 1" />
      <cmdb-button-group
        class="mr10"
        :trigger-text="$t('编辑')"
        :buttons="editButtonGroup"
        :expand="false">
      </cmdb-button-group>
      <cmdb-clipboard-selector class="mr10"
        label-key="bk_property_name"
        :list="clipboardList"
        :disabled="!table.checked.length"
        @on-copy="handleCopy">
      </cmdb-clipboard-selector>
      <cmdb-button-group
        class="mr10"
        :buttons="buttons"
        :expand="false">
      </cmdb-button-group>
    </div>
    <div class="options-right">
      <filter-fast-search class="option-fast-search" @search="searchFilter"></filter-fast-search>
      <icon-button class="option-filter ml10" icon="icon-cc-funnel"
        v-bk-tooltips.top="$t('高级筛选')"
        @click="handleSetFilters">
      </icon-button>
      <icon-button class="ml10"
        v-bk-tooltips="$t('查看删除历史')"
        icon="icon icon-cc-history"
        @click="routeToHistory">
      </icon-button>
    </div>

    <bk-sideslider
      v-transfer-dom
      :is-show.sync="slider.show"
      :title="slider.title"
      :width="800"
      :before-close="handleSliderBeforeClose">
      <div slot="content" style="height: 100%" v-bkloading="{ isLoading: $loading(Object.values(slider.request)) }">
        <component :is="slider.component"
          ref="multipleForm"
          :properties="properties.host"
          :property-groups="propertyGroups"
          :save-auth="saveAuth"
          @on-submit="handleMultipleSave"
          @on-cancel="handleSliderBeforeClose">
        </component>
      </div>
    </bk-sideslider>

    <bk-dialog
      class="assign-dialog"
      v-model="assign.show"
      header-position="left"
      :width="480"
      :mask-close="false"
      :esc-close="false"
      :close-icon="false"
      :title="assign.title"
      @cancel="closeAssignDialog">
      <div class="assign-content" v-if="assign.show">
        <i18n class="assign-count" tag="div" path="已选择主机">
          <template #count><span>{{table.checked.length}}</span></template>
        </i18n>
        <div class="assign-seleted">
          <p>{{assign.label}}</p>
          <bk-select
            font-size="normal"
            :searchable="true"
            :clearable="false"
            :disabled="$loading(assign.requestId)"
            :loading="$loading(authRequestId)"
            :placeholder="assign.placeholder"
            v-model="assign.id">
            <bk-option v-for="option in assignOptions"
              :key="option.id"
              :id="option.id"
              :name="option.name"
              :disabled="option.disabled">
              <cmdb-auth style="display: block;"
                :auth="option.auth"
                @update-auth="handleUpdateAssignAuth(option, ...arguments)">
                {{option.name}}
              </cmdb-auth>
            </bk-option>
          </bk-select>
        </div>
      </div>
      <div class="assign-footer" slot="footer">
        <bk-button
          class="mr10"
          theme="primary"
          :disabled="assign.id === ''"
          :loading="$loading(assign.requestId)"
          @click="handleConfirmAssign">
          {{$t('确定')}}
        </bk-button>
        <bk-button theme="default"
          :disabled="$loading(assign.requestId)"
          @click="closeAssignDialog">
          {{$t('取消')}}
        </bk-button>
      </div>
    </bk-dialog>
  </div>
</template>

<script>
  import { mapGetters } from 'vuex'
  import { AuthRequestId } from '@/components/ui/auth/auth-queue.js'
  import cmdbImport from '@/components/import/import'
  import cmdbButtonGroup from '@/components/ui/other/button-group'
  import Bus from '@/utils/bus.js'
  import RouterQuery from '@/router/query'
  import HostStore from '../transfer/host-store'
  import cmdbTransferMenu from '../transfer/transfer-menu'
  import FilterForm from '@/components/filters/filter-form.js'
  import FilterFastSearch from '@/components/filters/filter-fast-search'
  import FilterStore from '@/components/filters/store'
  import FilterUtils from '@/components/filters/utils'
  import hostImportService from '@/service/host/import'
  import { isUseComplexValueType, isEmptyPropertyValue } from '@/utils/tools'
  import isEqual from 'lodash/isEqual'
  const CUSTOM_STICKY_KEY = 'sticky-directory'

  export default {
    components: {
      cmdbImport,
      cmdbButtonGroup,
      cmdbTransferMenu,
      FilterFastSearch
    },
    data() {
      return {
        scope: '',
        businessList: [],
        slider: {
          show: false,
          title: '',
          component: null,
          request: {}
        },
        assign: {
          show: false,
          id: '',
          curSelected: '-1',
          placeholder: this.$t('请选择xx', { name: this.$t('业务') }),
          label: this.$t('业务列表'),
          title: this.$t('分配到业务空闲机', { idleSet: this.$store.state.globalConfig.config.set }),
          requestId: Symbol('assignHosts')
        },
        assignOptions: [],
        IPWithCloudSymbol: Symbol('IPWithCloud'),
        IPv6WithCloudSymbol: Symbol('IPv6WithCloud'),
        IPv46WithCloudSymbol: Symbol('IPv46WithCloud'),
        IPv64WithCloudSymbol: Symbol('IPv64WithCloud'),
        authRequestId: AuthRequestId,
        assignOptionAuthedMap: {
          biz: {},
          dir: {}
        },
        currentDirectoryIdList: []
      }
    },
    computed: {
      ...mapGetters('userCustom', ['usercustom']),
      ...mapGetters('resourceHost', [
        'activeDirectory',
        'defaultDirectory',
        'directorySortedList'
      ]),
      stickyDirectory() {
        return this.usercustom[CUSTOM_STICKY_KEY] || []
      },
      directoryId() {
        if (this.activeDirectory) {
          return this.activeDirectory.bk_module_id
        }
        // 默认会导入到空闲机目录
        return this.defaultDirectory?.bk_module_id
      },
      table() {
        return this.$parent.table
      },
      tableHeaderPropertyIdList() {
        return this.$parent.tableHeader.map(item => item.bk_property_id)
      },
      isAllResourceHost() {
        return this.table.selection.every(({ biz }) => {
          const [currentBiz] = biz
          return currentBiz.default === 1
        })
      },
      clipboardList() {
        const IPWithCloudFields = {
          [this.IPWithCloudSymbol]: `${this.$t('管控区域')}ID:IPv4`,
          [this.IPv6WithCloudSymbol]: `${this.$t('管控区域')}ID:IPv6`,
          [this.IPv46WithCloudSymbol]: `${this.$t('管控区域')}ID:IP(${this.$t('IPv4优先')})`,
          [this.IPv64WithCloudSymbol]: `${this.$t('管控区域')}ID:IP(${this.$t('IPv6优先')})`
        }
        const IPWithClouds = Object.getOwnPropertySymbols(IPWithCloudFields).map(key => FilterUtils.defineProperty({
          id: key,
          bk_obj_id: 'host',
          bk_property_id: key,
          bk_property_name: IPWithCloudFields[key],
          bk_property_type: 'singlechar'
        }))
        const clipboardList = this.$parent.tableHeader.slice()
        clipboardList.splice(1, 0, ...IPWithClouds)
        return clipboardList
      },
      properties() {
        return FilterStore.modelPropertyMap
      },
      propertyGroups() {
        return FilterStore.propertyGroups
      },
      buttons() {
        const buttonConfig = [{
          id: 'delete',
          text: this.$t('删除'),
          handler: this.handleMultipleDelete,
          disabled: !this.table.checked.length
        }, {
          id: 'export',
          text: this.$t('导出选中'),
          handler: this.exportField,
          disabled: !this.table.checked.length
        }, {
          id: 'batchExport',
          text: this.$t('导出全部'),
          handler: this.batchExportField,
          disabled: !this.table.pagination.count
        }]

        if (this.scope !== 1) {
          buttonConfig.splice(0, 1)
        }

        // 已分配并且是容器搜索模式
        if (this?.$parent?.isResourceAssigned && this?.$parent?.isContainerSearchMode) {
          buttonConfig.splice(buttonConfig.length - 1, 1)
        }

        return buttonConfig
      },
      saveAuth() {
        return this.table.selection.map(({ host, module, biz }) => {
          if (biz[0].default === 0) {
            return {
              type: this.$OPERATION.U_HOST,
              relation: [biz[0].bk_biz_id, host.bk_host_id]
            }
          }
          return {
            type: this.$OPERATION.U_RESOURCE_HOST,
            relation: [module[0].bk_module_id, host.bk_host_id]
          }
        })
      },
      editButtonGroup() {
        const buttonConfig = [{
          id: 'batch-edit',
          text: this.$t('批量编辑'),
          handler: this.handleMultipleEdit,
          disabled: !this.table.checked.length,
          tooltips: { content: this.$t('请先勾选需要编辑的主机'), disabled: this.table.checked.length > 0 }
        }, {
          id: 'import-edit',
          text: this.$t('导入编辑'),
          handler: this.handleEditImportInst
        }]
        if (this.scope !== 1) {
          buttonConfig.pop()
        }
        return buttonConfig
      },
      importInstError() {
        const importInstError = this.importInst.error || {}
        if (importInstError.bk_error_msg) {
          return importInstError.bk_error_msg
        }
        return importInstError.message || ''
      }
    },
    watch: {
      'importInst.show'(show) {
        if (!show) {
          this.importInst.type = 'new'
          this.importInst.active = 'import'
          this.importInst.error = null
          this.importInst.showTips = false
        } else {
          this.importInst.directory = this.directoryId
        }
      },
      'importInst.directory'(directory) {
        if (this.importInst.type === 'new') {
          this.importInst.payload = {
            bk_module_id: directory
          }
        }
      },
      'importInst.active'() {
        this.importInst.error = null
        this.importInst.showTips = false
      },
      '$route.query'(query, prev) {
        // 高级筛选自动打开面板，首次进入时无_t参数
        // eslint-disable-next-line no-underscore-dangle
        if (!prev._t && query.adv) {
          FilterForm.show()
        }
      },
      assignOptionAuthedMap: {
        handler() {
          this.sortAssignOptionByAuth()
        },
        deep: true
      }
    },
    async created() {
      try {
        this.unwatchScope = RouterQuery.watch('scope', (scope = 1) => {
          this.scope = isNaN(scope) ? 'all' : parseInt(scope, 10)
        }, { immediate: true })
        await this.getFullAmountBusiness()
      } catch (e) {
        console.error(e.message)
      }
    },
    beforeDestroy() {
      this.unwatchScope()
    },
    methods: {
      async getFullAmountBusiness() {
        try {
          const data = await this.$http.get('biz/simplify?sort=bk_biz_id')
          this.businessList = data.info || []
        } catch (e) {
          console.error(e)
          this.businessList = []
        } finally {
          HostStore.setBusinessList(this.businessList)
        }
      },
      sortAssignOptionByAuth() {
        // 暂时只排序业务，目录因为有置顶逻辑且数量不多
        if (this.assign.curSelected === 'toBusiness') {
          const list = this.businessList.map(item => ({
            ...item,
            is_pass: this.assignOptionAuthedMap.biz[item.bk_biz_id]
          }))
          list.sort((itemA, itemB) => itemB?.is_pass - itemA?.is_pass)
          this.businessList = list

          // 使用排序后的业务列表更新列表选项
          this.setAssignOptions(this.currentDirectoryIdList)
        }
      },
      openAgentApp() {
        const { agent } = window.Site
        if (agent) {
          const topWindow = window.top
          const isPaasConsole = topWindow !== window
          if (isPaasConsole) {
            topWindow.postMessage(JSON.stringify({
              action: 'open_other_app',
              app_code: 'bk_nodeman'
            }), '*')
          } else {
            window.open(agent)
          }
        } else {
          this.$warn(this.$t('未配置Agent安装APP地址'))
        }
      },
      handleAssignHosts(id) {
        this.assign.show = true

        this.assignOptionAuthedMap = {
          biz: {},
          dir: {}
        }

        const hosts = HostStore.getSelected()
        const directoryIds = hosts.map((host) => {
          const [module] = host.module
          return module.bk_module_id
        })
        const directoryIdList = [...new Set(directoryIds)]

        this.currentDirectoryIdList = directoryIdList

        if (id === 'toBusiness') {
          this.assign.placeholder = this.$t('请选择xx', { name: this.$t('业务') })
          this.assign.label = this.$t('业务列表')
          this.assign.title = this.$t('分配到业务空闲机', { idleSet: this.$store.state.globalConfig.config.set })
        } else {
          this.assign.placeholder = this.$t('请选择xx', { name: this.$t('目录') })
          this.assign.label = this.$t('目录列表')
          this.assign.title = this.$t('分配到主机池其他目录')
        }

        this.setAssignOptions(directoryIdList)
      },
      setAssignOptions(directoryIdList) {
        if (this.assign.curSelected === 'toBusiness') {
          this.assignOptions = this.businessList.map(item => ({
            id: item.bk_biz_id,
            name: `[${item.bk_biz_id}] ${item.bk_biz_name}`,
            disabled: !item?.is_pass ?? true,
            auth: directoryIdList.map(directoryId => ({
              type: this.$OPERATION.TRANSFER_HOST_TO_BIZ,
              relation: [[[directoryId], [item.bk_biz_id]]]
            }))
          }))
        } else {
          const directoryList = this.directorySortedList(this.stickyDirectory)
          const directorySortedList = directoryIdList?.length === 1
            ? directoryList.filter(item => item.bk_module_id !== directoryIdList[0])
            : directoryList

          this.assignOptions = directorySortedList.map(item => ({
            id: item.bk_module_id,
            name: item.bk_module_name,
            disabled: true,
            auth: directoryIdList.map(directoryId => ({
              type: this.$OPERATION.TRANSFER_HOST_TO_DIRECTORY,
              relation: [[[directoryId], [item.bk_module_id]]]
            }))
          }))
        }
      },
      handleUpdateAssignAuth(option, authorized) {
        option.disabled = !authorized
        if (this.assign.curSelected === 'toBusiness') {
          this.$set(this.assignOptionAuthedMap.biz, option.id, authorized)
        } else {
          this.$set(this.assignOptionAuthedMap.dir, option.id, authorized)
        }
      },
      closeAssignDialog() {
        this.assign.id = ''
        this.assign.show = false
        this.assign.curSelected = '-1'
      },
      handleConfirmAssign() {
        this.assign.curSelected === 'toBusiness' ? this.assignHostsToBusiness() : this.changeHostsDir()
      },
      async assignHostsToBusiness() {
        await this.$store.dispatch('resourceDirectory/assignHostsToBusiness', {
          params: {
            bk_biz_id: this.assign.id,
            bk_host_id: this.table.checked
          },
          config: {
            requestId: this.assign.requestId
          }
        }).then(() => {
          Bus.$emit('refresh-dir-count')
          this.$success(this.$t('分配成功'))
          this.closeAssignDialog()
          RouterQuery.set({
            page: 1,
            _t: Date.now()
          })
        })
          .catch((e) => {
            console.error(e)
          })
      },
      async changeHostsDir() {
        try {
          await this.$store.dispatch('resource/host/transfer/directory', {
            params: {
              bk_module_id: this.assign.id,
              bk_host_id: this.table.checked
            },
            config: {
              requestId: this.assign.requestId
            }
          })
          Bus.$emit('refresh-dir-count')
          this.$success(this.$t('转移成功'))
          this.closeAssignDialog()
          RouterQuery.set({
            page: 1,
            _t: Date.now()
          })
        } catch (e) {
          console.error(e)
        }
      },
      handleCopy(property) {
        const copyText = this.table.selection.map((data, index) => {
          const modelId = property.bk_obj_id
          const modelData = data[modelId]

          if (isUseComplexValueType(property)) {
            const value = this.$parent?.$refs?.[`table-cell-property-value-${property.bk_property_id}`]?.[index]?.getCopyValue()
            return value
          }

          const IPWithCloudKeys = [
            this.IPWithCloudSymbol,
            this.IPv6WithCloudSymbol,
            this.IPv46WithCloudSymbol,
            this.IPv64WithCloudSymbol
          ]
          if (IPWithCloudKeys.includes(property.id)) {
            const cloud = this.$tools.getPropertyCopyValue(modelData.bk_cloud_id, 'foreignkey')
            const ip = this.$tools.getPropertyCopyValue(modelData.bk_host_innerip, 'singlechar')
            const ipv6 = this.$tools.getPropertyCopyValue(modelData.bk_host_innerip_v6, 'singlechar')
            const isEmptyIPv4Value = isEmptyPropertyValue(modelData.bk_host_innerip)
            const isEmptyIPv6Value = isEmptyPropertyValue(modelData.bk_host_innerip_v6)
            if (property.id === this.IPWithCloudSymbol) {
              return `${cloud}:${ip}`
            }
            if (property.id === this.IPv6WithCloudSymbol) {
              return `${cloud}:[${ipv6}]`
            }
            if (property.id === this.IPv46WithCloudSymbol) {
              if (!isEmptyIPv4Value || isEmptyIPv6Value) {
                return `${cloud}:${ip}`
              }
              return `${cloud}:[${ipv6}]`
            }
            if (property.id === this.IPv64WithCloudSymbol) {
              if (isEmptyIPv4Value || !isEmptyIPv6Value) {
                return `${cloud}:[${ipv6}]`
              }
              return `${cloud}:${ip}`
            }
          }
          if (property.bk_property_type === 'topology') {
            // eslint-disable-next-line no-underscore-dangle
            return data.__bk_host_topology__.join(',')
          }
          const propertyId = property.bk_property_id
          const copyValueOptions = {}
          if (propertyId === 'bk_cloud_id') {
            copyValueOptions.isFullCloud = true
          }
          if (Array.isArray(modelData)) {
            const value = modelData
              .map(item => this.$tools.getPropertyCopyValue(item[propertyId], property, copyValueOptions))
            return value.join(',')
          }
          return this.$tools.getPropertyCopyValue(modelData[propertyId], property, copyValueOptions)
        })
        this.$copyText(copyText.join('\n')).then(() => {
          this.$success(this.$t('复制成功'))
        }, () => {
          this.$error(this.$t('复制失败'))
        })
      },
      async handleMultipleEdit() {
        this.slider.title = this.$t('主机属性')
        this.slider.show = true
        setTimeout(() => {
          this.slider.component = 'cmdb-form-multiple'
        }, 200)
      },
      async handleMultipleSave(changedValues) {
        await this.$store.dispatch('hostUpdate/updateHost', {
          params: {
            ...changedValues,
            bk_host_id: this.table.checked.join(',')
          }
        })
        this.slider.show = false
        RouterQuery.set({
          _t: Date.now()
        })
      },
      handleMultipleDelete() {
        this.$bkInfo({
          title: `${this.$t('确定删除选中的主机')}？`,
          confirmFn: () => {
            this.$store.dispatch('hostDelete/deleteHost', {
              params: {
                data: {
                  bk_host_id: this.table.checked.join(','),
                  bk_supplier_account: this.supplierAccount
                }
              }
            }).then(() => {
              this.$success(this.$t('成功删除选中的主机'))
              RouterQuery.set({
                page: 1,
                _t: Date.now()
              })
              Bus.$emit('refresh-dir-count')
            })
          }
        })
      },
      handleSliderBeforeClose() {
        const $form = this.$refs.multipleForm
        const { values, refrenceValues } = $form
        const changedValues =  !isEqual(values, refrenceValues)
        if (changedValues) {
          $form.setChanged(true)
          return $form.beforeClose(() => {
            this.slider.component = null
            this.slider.show = false
          })
        }
        this.slider.show = false
        this.slider.component = null
      },
      async exportField() {
        const useExport = await import('@/components/export-file')
        useExport.default({
          title: this.$t('导出选中'),
          bk_obj_id: 'host',
          presetFields: ['bk_cloud_id', 'bk_host_innerip'],
          defaultSelectedFields: this.tableHeaderPropertyIdList,
          count: this.table.selection.length,
          submit: (state, task) => {
            const { fields, exportRelation  } = state
            const params = {
              export_custom_fields: fields.value.map(property => property.bk_property_id),
              bk_host_ids: this.table.selection.map(({ host }) => host.bk_host_id),
              export_condition: {
                page: {
                  start: 0,
                  limit: this.table.selection.length
                }
              }
            }
            if (exportRelation.value) {
              params.object_unique_id = state.object_unique_id.value
              params.association_condition = state.relations.value
            }
            return this.$http.download({
              url: `${window.API_HOST}hosts/export`,
              method: 'post',
              name: task.current.value.name,
              data: params
            })
          }
        }).show()
      },
      async batchExportField() {
        const useExport = await import('@/components/export-file')
        useExport.default({
          title: this.$t('导出全部'),
          bk_biz_id: this.bizId,
          bk_obj_id: 'host',
          presetFields: ['bk_cloud_id', 'bk_host_innerip'],
          defaultSelectedFields: this.tableHeaderPropertyIdList,
          count: this.table.pagination.count,
          submit: (state, task) => {
            const { fields, exportRelation  } = state
            const exportCondition = this.$parent.getParams()
            const params = {
              export_custom_fields: fields.value.map(property => property.bk_property_id),
              bk_biz_id: this.bizId,
              export_condition: {
                ...exportCondition,
                page: {
                  ...task.current.value.page,
                  sort: 'bk_host_id'
                }
              }
            }
            if (exportRelation.value) {
              params.object_unique_id = state.object_unique_id.value
              params.association_condition = state.relations.value
            }
            return this.$http.download({
              url: `${window.API_HOST}hosts/export`,
              method: 'post',
              name: task.current.value.name,
              data: params
            })
          }
        }).show()
      },
      routeToHistory() {
        this.$routerActions.redirect({
          name: 'hostHistory',
          history: true
        })
      },
      handleSetFilters() {
        FilterForm.show()
      },
      async handleNewImportInst() {
        const useImport = await import('@/components/import-file')
        const [, { show: showImport, setState: setImportState }] = useImport.default()
        setImportState({
          title: this.$t('导入主机'),
          bk_obj_id: 'host',
          template: `${window.API_HOST}importtemplate/host`,
          fileTips: this.$t('导入文件大小提示'),
          submit: (options) => {
            const params = {
              bk_module_id: this.directoryId,
              op: options.step
            }
            if (options.importRelation) {
              params.object_unique_id = options.object_unique_id
              params.association_condition = options.relations
            }
            return hostImportService.create({ file: options.file, params, config: options.config })
          },
          success: () => {
            RouterQuery.set({ _t: Date.now() })
            Bus.$emit('refresh-dir-count')
          }
        })
        showImport()
      },
      async handleEditImportInst() {
        const useImport = await import('@/components/import-file')
        const [, { show: showImport, setState: setImportState }] = useImport.default()
        setImportState({
          title: this.$t('导入编辑'),
          bk_obj_id: 'host',
          fileTips: this.$t('导入文件大小提示'),
          submit: (options) => {
            const params = {
              op: options.step
            }
            if (options.importRelation) {
              params.object_unique_id = options.object_unique_id
              params.association_condition = options.relations
            }
            return hostImportService.update({ file: options.file, params, config: options.config })
          },
          success: () => RouterQuery.set({ _t: Date.now() })
        })
        showImport()
      },
      searchFilter(value) {
        this.$emit('search', value)
      }
    }
  }
</script>

<style lang="scss" scoped>
    .options-left {
        font-size: 0;
        .assign-selector {
            min-width: 80px;
        }
    }
    .options-right {
        display: flex;
        justify-content: flex-end;
        overflow: hidden;
        .option-fast-search {
            flex: 1;
            margin-left: 10px;
            max-width: 300px;
        }
        .option-filter:hover {
            color: $primaryColor;
        }
    }
    .import-prepend {
        margin: 20px 29px -10px 33px;
        /deep/ {
            .bk-form-item {
                display: flex;
                margin-bottom: 8px;
            }
            .bk-label {
                width: auto !important;
            }
            .bk-form-content {
                flex: 1;
                margin-left: auto !important;
            }
        }
        .form-item-tips {
            position: absolute;
            color: #EA3636;
            font-size: 12px;
        }
    }
    .automatic-import{
        padding: 44px 30px 0 30px;
        display: flex;
        flex-direction: column;
        align-items: center;
        .agent-install-tips1 {
            font-size: 14px;
            margin: 10px 0;
        }
        .agent-install-tips2 {
            font-size: 12px;
            color: #979BA5;
        }
        .agent-install-button {
            margin: 14px 0;
        }
    }
    .assign-dialog {
        color: #63656E;
        .assign-count {
            font-size: 12px;
            padding-bottom: 20px;
            span {
                font-weight: bold;
            }
        }
        .assign-seleted {
            .bk-select {
                width: 100%;
                margin-top: 10px;
            }
        }
        .assign-footer {
            font-size: 0;
        }
    }
    .apply-others {
        display: inline-block;
        width: 60%;
        font-size: 12px;
        color: #63656E;
        line-height: 32px;
        cursor: pointer;
        &:hover {
            color: #3a84ff;
        }
        .bk-icon {
            font-size: 14px;
            display: inline-block;
            vertical-align: -2px;
        }
    }

    .edit-import-panel {
        .alert {
            margin: 24px 29px 0 33px
        }
        /deep/ {
            .up-file {
                margin-top: 20px;
            }
        }
    }

    .upload-error-message {
        margin: 8px 0;
        color: $dangerColor;
    }
</style>
