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
  <bk-table v-if="!serviceInstance.pending"
    :data="list"
    :outer-border="false"
    :header-cell-style="{ backgroundColor: '#fff' }"
    v-bind="dynamicProps"
    v-bkloading="{ isLoading: $loading(request.list) }">
    <bk-table-column v-for="property in header"
      :key="property.bk_property_id"
      :label="property.bk_property_name"
      :prop="property.bk_property_id"
      :show-overflow-tooltip="property.bk_property_id !== 'bind_info'">
      <template slot-scope="{ row }">
        <cmdb-property-value v-if="property.bk_property_id !== 'bind_info'"
          :theme="property.bk_property_id === 'bk_func_name' ? 'primary' : 'default'"
          :value="row.property[property.bk_property_id]"
          :show-unit="false"
          :property="property"
          @click.native="handleView(row)">
        </cmdb-property-value>
        <process-bind-info-value v-else
          :value="row.property[property.bk_property_id]"
          :property="property">
        </process-bind-info-value>
      </template>
    </bk-table-column>
    <bk-table-column v-if="!readonly" width="150" :resizable="false">
      <div class="options-wrapper" slot-scope="{ row }">
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
      </div>
    </bk-table-column>
  </bk-table>
</template>

<script>
  import { processPropertyRequestId } from '@/components/service/form/symbol'
  import { processTableHeader } from '@/dictionary/table-header'
  import { mapGetters } from 'vuex'
  import Form from '@/components/service/form/form.js'
  import ProcessBindInfoValue from '@/components/service/process-bind-info-value'
  import Bus from '../common/bus'
  export default {
    components: {
      ProcessBindInfoValue
    },
    props: {
      serviceInstance: Object,
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
        properties: [],
        header: [],
        list: [],
        request: {
          list: Symbol('getList'),
          delete: Symbol('delete')
        }
      }
    },
    computed: {
      ...mapGetters(['supplierAccount']),
      bizId() {
        const  { objectBiz, bizSet } = this.$store.state
        return objectBiz.bizId || bizSet.bizId
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
      Bus.$on('refresh-process-list', this.handleRefresh)
    },
    beforeDestroy() {
      Bus.$off('refresh-process-list', this.handleRefresh)
    },
    methods: {
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
          const property = this.properties.find(property => property.bk_property_id === id)
          if (property) {
            header.push(property)
          }
        })
        this.header = header
      },
      handleRefresh(target) {
        if (target !== this.serviceInstance) {
          return
        }
        this.getList()
      },
      async getList() {
        try {
          const reqParams = {
            bk_biz_id: this.bizId,
            service_instance_id: this.serviceInstance.id
          }
          const reqConfig = {
            requestId: this.request.list,
            cancelPrevious: true,
            cancelWhenRouteChange: true
          }

          this.list = await this.listRequest(reqParams, reqConfig)
        } catch (error) {
          console.error(error)
        } finally {
          this.$emit('update-list', this.list)
        }
      },
      handleView(row) {
        Form.show({
          type: 'view',
          title: this.$t('查看进程'),
          instance: row.property,
          hostId: row.relation.bk_host_id,
          bizId: this.bizId,
          serviceTemplateId: this.serviceInstance.service_template_id,
          processTemplateId: row.relation.process_template_id,
          submitHandler: this.editSubmitHandler,
          showOptions: !this.readonly,
        })
      },
      handleEdit(row) {
        Form.show({
          type: 'update',
          title: this.$t('编辑进程'),
          instance: row.property,
          hostId: row.relation.bk_host_id,
          bizId: this.bizId,
          serviceTemplateId: this.serviceInstance.service_template_id,
          processTemplateId: row.relation.process_template_id,
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
          await this.$store.dispatch('processInstance/deleteServiceInstanceProcess', {
            config: {
              data: {
                bk_biz_id: this.bizId,
                process_instance_ids: [row.property.bk_process_id]
              },
              requestId: this.request.delete
            }
          })
          if (this.list.length === 1) {
            Bus.$emit('delete-complete')
          } else {
            this.getList()
          }
        } catch (error) {
          console.error(error)
        }
      }
    }
  }
</script>
