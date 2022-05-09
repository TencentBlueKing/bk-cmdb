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
  <bk-table v-if="!process.pending"
    ref="expandListTable"
    v-bkloading="{ isLoading: $loading(Object.values(request)) }"
    :data="list"
    :outer-border="false"
    :header-cell-style="{ backgroundColor: '#fff' }"
    v-bind="dynamicProps"
    @selection-change="handleSelectionChange">
    <bk-table-column v-if="!readonly" type="selection" fixed></bk-table-column>
    <bk-table-column :label="$t('所属实例')" prop="service_instance_name" min-width="150" fixed show-overflow-tooltip>
      <template slot-scope="{ row }">
        <span class="instance-name" @click.stop="handleView(row)">{{row.service_instance_name}}</span>
      </template>
    </bk-table-column>
    <bk-table-column v-for="property in header"
      :key="property.bk_property_id"
      :label="property.bk_property_name"
      :prop="property.bk_property_id"
      :show-overflow-tooltip="property.bk_property_id !== 'bind_info'">
      <template slot-scope="{ row }">
        <cmdb-property-value v-if="property.bk_property_id !== 'bind_info'"
          :value="row.property[property.bk_property_id]"
          :show-unit="false"
          :property="property">
        </cmdb-property-value>
        <process-bind-info-value v-else
          :value="row.property[property.bk_property_id]"
          :property="property">
        </process-bind-info-value>
      </template>
    </bk-table-column>
    <bk-table-column v-if="!readonly" :label="$t('操作')" width="150" fixed="right" :resizable="false">
      <template slot-scope="{ row }">
        <cmdb-auth class="mr10" :auth="{ type: $OPERATION.U_SERVICE_INSTANCE, relation: [bizId] }">
          <bk-button slot-scope="{ disabled }" v-test-id="'edit'"
            theme="primary" text
            :disabled="disabled"
            @click="handleEdit(row)">
            {{$t('编辑')}}
          </bk-button>
        </cmdb-auth>
        <cmdb-auth :auth="{ type: $OPERATION.U_SERVICE_INSTANCE, relation: [bizId] }"
          v-if="!row.relation.process_template_id">
          <bk-popconfirm trigger="click" slot-scope="{ disabled }"
            ext-popover-cls="del-confirm"
            :content="$t('确定删除该进程')"
            confirm-loading
            @confirm="handleDelete(row)">
            <bk-button v-test-id="'del'"
              theme="primary" text
              :disabled="disabled">
              {{$t('删除')}}
            </bk-button>
          </bk-popconfirm>
        </cmdb-auth>
      </template>
    </bk-table-column>
  </bk-table>
</template>

<script>
  import { processPropertyRequestId } from '@/components/service/form/symbol'
  import { processTableHeader } from '@/dictionary/table-header'
  import { mapGetters } from 'vuex'
  import Bus from '../common/bus'
  import RootBus from '@/utils/bus'
  import Form from '@/components/service/form/form.js'
  import ProcessBindInfoValue from '@/components/service/process-bind-info-value'
  import RouterQuery from '@/router/query'
  export default {
    components: {
      ProcessBindInfoValue
    },
    props: {
      process: Object,
      /**
       * 是否只读
       */
      readonly: {
        type: Boolean,
        default: false
      },
      /**
       * 列表请求方法
       */
      listRequest: {
        type: Function,
        default: null,
        required: true
      }
    },
    data() {
      return {
        list: [],
        properties: [],
        header: [],
        selection: [],
        isEmitByOthers: false,
        request: {
          list: Symbol('list'),
          delete: Symbol('delete')
        }
      }
    },
    computed: {
      ...mapGetters(['supplierAccount']),
      ...mapGetters('businessHost', ['selectedNode']),
      bizId() {
        const { objectBiz, bizSet } = this.$store.state
        return objectBiz.bizId || bizSet.bizId
      },
      serviceTemplateId() {
        return this.selectedNode && this.selectedNode.data.service_template_id
      },
      dynamicProps() {
        const dynamicProps = {}
        const paddingHeight = 43
        const rowHeight = 42
        if (this.list.length && this.list.length < 3) {
          // eslint-disable-next-line no-mixed-operators
          dynamicProps.height = paddingHeight + rowHeight * (this.list.length + 1)
        }
        return dynamicProps
      }
    },
    created() {
      this.getProperties()
      this.getList()
      Bus.$on('process-selection-change', this.handleProcessSelectionChange)
      Bus.$on('batch-delete', this.handeBatchDelete)
      Bus.$on('refresh-expand-list', this.handleRefreshList)
    },
    beforeDestroy() {
      Bus.$off('process-selection-change', this.handleProcessSelectionChange)
      Bus.$off('batch-delete', this.handeBatchDelete)
      Bus.$off('refresh-expand-list', this.handleRefreshList)
    },
    methods: {
      handleRefreshList(target) {
        if (target !== this.process) {
          return false
        }
        this.getList()
      },
      async getList() {
        try {
          const reqParams = {
            bk_biz_id: this.bizId,
            process_ids: this.process.process_ids,
            page: { limit: 999999999 }
          }
          const reqConfig = {
            requestId: this.request.list,
            cancelPrevious: true
          }

          this.list = await this.listRequest(reqParams, reqConfig)
        } catch (error) {
          console.error(error)
          this.list = []
        } finally {
          this.$emit('resolved', this.list)
          this.$nextTick(this.checkSelection)
        }
      },
      checkSelection() {
        const reserved = this.process.reserved.map(data => data.process_id)
        const existSelection = this.list.filter(row => reserved.includes(row.process_id))
        existSelection.forEach((row) => {
          this.$refs.expandListTable.toggleRowSelection(row, true)
        })
      },
      async getProperties() {
        try {
          this.properties = await this.$store.dispatch('objectModelProperty/searchObjectAttribute', {
            params: {
              bk_obj_id: 'process',
              bk_supplier_account: this.supplierAccount
            },
            config: {
              requestId: processPropertyRequestId,
              fromCache: true
            }
          })
          this.setHeader()
        } catch (error) {
          console.error(error)
        }
      },
      setHeader() {
        const header = []
        processTableHeader.forEach((id) => {
          if (id === 'bk_process_name') {
            return
          }
          const property = this.properties.find(property => property.bk_property_id === id)
          if (property) {
            header.push(property)
          }
        })
        this.header = header
      },
      handleProcessSelectionChange(process) {
        if (process === this.process) {
          return false
        }
        this.isEmitByOthers = true
        this.$refs.expandListTable.clearSelection()
        this.$nextTick(() => {
          this.isEmitByOthers = false
        })
      },
      handleSelectionChange(selection) {
        this.selection = selection
        if (!this.isEmitByOthers) {
          Bus.$emit('process-selection-change', this.process, selection, this.request.delete)
        }
        Bus.$emit('update-reserve-selection', this.process, selection)
      },
      handleView(row) {
        Form.show({
          type: 'view',
          title: this.$t('查看进程'),
          instance: row.property,
          serviceTemplateId: this.selectedNode.data.service_template_id,
          processTemplateId: row.relation.process_template_id,
          hostId: row.relation.bk_host_id,
          bizId: this.bizId,
          submitHandler: this.editSubmitHandler,
          showOptions: !this.readonly,
        })
      },
      handleEdit(row) {
        Form.show({
          type: 'update',
          title: this.$t('编辑进程'),
          instance: row.property,
          serviceTemplateId: this.selectedNode.data.service_template_id,
          processTemplateId: row.relation.process_template_id,
          hostId: row.relation.bk_host_id,
          bizId: this.bizId,
          submitHandler: this.editSubmitHandler
        })
      },
      async editSubmitHandler(values, changedValues, instance) {
        try {
          await this.$store.dispatch('processInstance/updateServiceInstanceProcess', {
            params: {
              bk_biz_id: this.bizId,
              processes: [{ ...instance, ...values }]
            }
          })
          this.getList()
        } catch (error) {
          console.error(error)
        }
      },
      async handleDelete(row) {
        try {
          await this.dispatchDelete([row.process_id])
          if (this.list.length === 1) {
            this.refreshParentList()
          } else {
            this.getList()
          }
        } catch (error) {
          console.error(error)
        }
      },
      async handeBatchDelete(name) {
        if (name !== this.process.bk_process_name) {
          return
        }
        try {
          await this.dispatchDelete(this.selection.map(row => row.process_id))
          if (this.selection.length === this.list.length) {
            this.refreshParentList()
          } else {
            this.getList()
          }
        } catch (error) {
          console.error(error.message)
        } finally {
          this.selection = []
        }
      },
      dispatchDelete(ids) {
        return this.$store.dispatch('processInstance/deleteServiceInstanceProcess', {
          config: {
            data: {
              bk_biz_id: this.bizId,
              process_instance_ids: ids
            },
            requestId: this.request.delete
          }
        })
      },
      refreshParentList() {
        RootBus.$emit('refresh-count-by-node')
        RouterQuery.set({
          _t: Date.now(),
          page: 1
        })
      }
    }
  }
</script>

<style lang="scss" scoped>
    .instance-name {
        cursor: pointer;
        color: $primaryColor;
    }
</style>
