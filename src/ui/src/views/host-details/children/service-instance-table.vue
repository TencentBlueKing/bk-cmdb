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
  <div class="table-layout" v-show="show">
    <div class="table-title" @click="localExpanded = !localExpanded"
      @mouseenter="handleShowDotMenu"
      @mouseleave="handleHideDotMenu">
      <bk-checkbox v-if="!readonly" class="title-checkbox"
        :size="16"
        v-model="checked"
        @click.native.stop>
      </bk-checkbox>
      <i class="title-icon bk-icon icon-down-shape" v-if="localExpanded"></i>
      <i class="title-icon bk-icon icon-right-shape" v-else></i>
      <template v-if="!instance.editing.name">
        <span class="title-label" v-bk-overflow-tips>{{instance.name}}</span>
        <cmdb-dot-menu v-if="!readonly" class="instance-menu" ref="dotMenu" @click.native.stop>
          <ul class="menu-list"
            @mouseenter="handleShowDotMenu"
            @mouseleave="handleHideDotMenu">
            <li class="menu-item"
              v-for="(menu, index) in instanceMenu"
              :key="index">
              <cmdb-auth class="menu-span" :auth="HOST_AUTH[menu.auth]">
                <bk-button slot-scope="{ disabled }"
                  class="menu-button"
                  :text="true"
                  :disabled="disabled"
                  @click="menu.handler">
                  {{menu.name}}
                </bk-button>
              </cmdb-auth>
            </li>
          </ul>
        </cmdb-dot-menu>
      </template>
      <div class="edit-form" v-else>
        <service-instance-name-edit-form ref="nameEditForm"
          :value="instance.name"
          @click.native.stop
          @confirm="handleConfirmEditName"
          @cancel="handleCancelEditName" />
      </div>
      <div class="right-content">
        <div class="instance-label clearfix" @click.stop>
          <div class="label-list fl">
            <div class="label-item" :title="`${label.key}：${label.value}`" :key="index"
              v-for="(label, index) in labelShowList">
              <span>{{label.key}}</span>
              <span>:</span>
              <span>{{label.value}}</span>
            </div>
            <bk-popover class="label-item label-tips"
              v-if="labelTipsList.length"
              theme="light label-tips"
              :width="290"
              placement="bottom-end">
              <span>...</span>
              <div class="tips-label-list" slot="content">
                <span class="label-item" :title="`${label.key}：${label.value}`" :key="index"
                  v-for="(label, index) in labelTipsList">
                  <span>{{label.key}}</span>
                  <span>:</span>
                  <span>{{label.value}}</span>
                </span>
              </div>
            </bk-popover>
          </div>
        </div>
        <span class="topology-path" v-bk-overflow-tips
          @click.stop="goTopologyInstance">
          {{topologyPath}}
        </span>
      </div>
    </div>
    <bk-table class="process-list-table"
      v-show="localExpanded"
      v-bkloading="{ isLoading: $loading(Object.values(requestId)) }"
      :data="list"
      :row-style="{ cursor: 'pointer' }"
      @row-click="showProcessDetails">
      <bk-table-column v-for="(column, index) in header"
        :class-name="index === 0 ? 'is-highlight' : ''"
        :key="column.id"
        :prop="column.id"
        :label="column.name"
        :show-overflow-tooltip="column.id !== 'bind_info'">
        <template slot-scope="{ row }">
          <cmdb-property-value v-if="column.id !== 'bind_info'"
            :value="(row.property || {})[column.id]"
            :show-unit="false"
            :property="column.property">
          </cmdb-property-value>
          <process-bind-info-value v-else
            :value="(row.property || {})[column.id]"
            :property="column.property"
            :popover-optoins="{ placement: 'bottom' }">
          </process-bind-info-value>
        </template>
      </bk-table-column>
      <bk-table-column v-if="!readonly" width="150" :resizable="false" :label="$t('操作')">
        <template slot-scope="{ row }">
          <cmdb-auth class="mr10" :auth="{ type: $OPERATION.U_SERVICE_INSTANCE, relation: [bizId] }">
            <bk-button slot-scope="{ disabled }"
              theme="primary" text
              :disabled="disabled"
              @click.native.stop
              @click="handleEdit(row)">
              {{$t('编辑')}}
            </bk-button>
          </cmdb-auth>
          <cmdb-auth :auth="{ type: $OPERATION.U_SERVICE_INSTANCE, relation: [bizId] }"
            v-if="!instance.service_template_id">
            <bk-button slot-scope="{ disabled }"
              theme="primary" text
              :disabled="disabled"
              @click.native.stop
              @click="handleDelete(row)">
              {{$t('删除')}}
            </bk-button>
          </cmdb-auth>
        </template>
      </bk-table-column>
      <template slot="empty">
        <span class="process-count-tips" v-if="instance.service_template_id">
          <i class="tips-icon bk-icon icon-exclamation-circle"></i>
          <i18n class="tips-content" path="模板服务实例无进程提示">
            <template #link>
              <cmdb-auth class="tips-link"
                :auth="{ type: $OPERATION.U_SERVICE_INSTANCE, relation: [bizId] }"
                @click="redirectToTemplate">
                {{$t('跳转添加并同步')}}
              </cmdb-auth>
            </template>
          </i18n>
        </span>
        <span class="process-count-tips" v-else>
          <i class="tips-icon bk-icon icon-exclamation-circle"></i>
          <i18n class="tips-content" path="普通服务实例无进程提示">
            <template #link>
              <cmdb-auth class="tips-link"
                :auth="{ type: $OPERATION.U_SERVICE_INSTANCE, relation: [bizId] }"
                @click="handleAddProcess">
                {{$t('立即添加')}}
              </cmdb-auth>
            </template>
          </i18n>
        </span>
      </template>
    </bk-table>
  </div>
</template>

<script>
  import {
    MENU_BUSINESS_HOST_AND_SERVICE,
    MENU_BUSINESS_SERVICE_TEMPLATE_DETAILS
  } from '@/dictionary/menu-symbol'
  import { processTableHeader } from '@/dictionary/table-header'
  import ProcessBindInfoValue from '@/components/service/process-bind-info-value'
  import ServiceInstanceNameEditForm from '@/components/service/instance-name-edit-form'
  import { mapState } from 'vuex'
  import authMixin from '../mixin-auth'
  import Form from '@/components/service/form/form.js'
  import { readonlyMixin } from '../mixin-readonly'
  import { serviceInstanceProcessesProxy } from '../service-proxy'

  export default {
    components: {
      ProcessBindInfoValue,
      ServiceInstanceNameEditForm
    },
    mixins: [authMixin, readonlyMixin],
    props: {
      instance: {
        type: Object,
        required: true
      },
      expanded: Boolean
    },
    data() {
      return {
        tipsLabel: {
          show: false,
          id: null
        },
        show: true,
        checked: false,
        localExpanded: this.expanded,
        properties: [],
        header: [],
        list: []
      }
    },
    computed: {
      ...mapState('hostDetails', ['info']),
      bizId() {
        const [biz] = this.info.biz || []
        return biz ? biz.bk_biz_id : null
      },
      withTemplate() {
        return this.isModuleNode && !!this.instance.service_template_id
      },
      instanceMenu() {
        const menu = [{
          name: this.$t('删除'),
          handler: this.handleDeleteInstance,
          auth: 'D_SERVICE_INSTANCE'
        }, {
          name: this.$t('编辑'),
          handler: this.handleEditInstance,
          auth: 'U_SERVICE_INSTANCE'
        }]
        return menu
      },
      requestId() {
        return {
          processList: `get_service_instance_${this.instance.id}_processes`,
          properties: 'get_service_process_properties',
          deleteProcess: 'delete_service_process'
        }
      },
      labelList() {
        const list = []
        const { labels } = this.instance
        labels && Object.keys(labels).forEach((key) => {
          list.push({
            id: this.instance.id,
            key,
            value: labels[key]
          })
        })
        return list
      },
      labelShowList() {
        return this.labelList.slice(0, 3)
      },
      labelTipsList() {
        return this.labelList.slice(3)
      },
      topologyPath() {
        const pathArr = this.$tools.clone(this.instance.topo_path).reverse()
        const path = pathArr.map(path => path.bk_inst_name).join(' / ')
        return path
      }
    },
    watch: {
      localExpanded(expanded) {
        if (expanded) {
          this.getServiceProcessList()
        }
      },
      checked(checked) {
        this.$emit('check-change', checked, this.instance)
      }
    },
    async created() {
      await this.getProcessProperties()
      if (this.expanded) {
        this.getServiceProcessList()
      }
    },
    methods: {
      async getProcessProperties() {
        const action = 'objectModelProperty/searchObjectAttribute'
        const properties = await this.$store.dispatch(action, {
          params: {
            bk_obj_id: 'process',
            bk_supplier_account: this.$store.getters.supplierAccount
          },
          config: {
            requestId: this.requestId.properties,
            fromCache: true
          }
        })
        this.properties = properties
        this.setHeader()
      },
      getServiceProcessList() {
        serviceInstanceProcessesProxy(
          {
            service_instance_id: this.instance.id,
            bk_biz_id: this.info.biz[0].bk_biz_id
          },
          {
            requestId: this.requestId.processList
          }
        ).then((list) => {
          this.list = list
        })
          .catch(() => {
            this.list = []
          })
      },
      setHeader() {
        const header = processTableHeader.map((id) => {
          const property = this.properties.find(property => property.bk_property_id === id) || {}
          return {
            id: property.bk_property_id,
            name: this.$tools.getHeaderPropertyName(property),
            property
          }
        })
        this.header = header
      },
      handleDeleteInstance() {
        this.$bkInfo({
          title: this.$t('确定删除该服务实例'),
          confirmLoading: true,
          confirmFn: async () => {
            try {
              await this.$store.dispatch('serviceInstance/deleteServiceInstance', {
                config: {
                  data: {
                    service_instance_ids: [this.instance.id],
                    bk_biz_id: this.bizId
                  }
                }
              })
              this.$success(this.$t('删除成功'))
              this.$emit('delete-instance')
              return true
            } catch (e) {
              console.error(e)
              return false
            }
          }
        })
      },
      handleEditInstance() {
        this.$emit('edit-name')
        this.$nextTick(() => {
          this.$refs.nameEditForm.focus()
        })
      },
      async handleConfirmEditName(value) {
        try {
          await this.$store.dispatch('serviceInstance/updateServiceInstance', {
            bizId: this.bizId,
            params: {
              data: [{
                service_instance_id: this.instance.id,
                update: {
                  name: value
                }
              }]
            }
          })
          this.$emit('edit-name-success', value)
        } catch (error) {
          console.error(error)
        }
      },
      handleCancelEditName() {
        this.$emit('cancel-edit-name')
      },
      goTopologyInstance() {
        this.$routerActions.redirect({
          name: MENU_BUSINESS_HOST_AND_SERVICE,
          query: {
            tab: 'serviceInstance',
            node: `module-${this.instance.bk_module_id}`,
            instanceName: this.instance.name
          },
          history: true
        })
      },
      handleShowDotMenu() {
        if (this.$refs.dotMenu) {
          this.$refs.dotMenu.$el.style.opacity = 1
        }
      },
      handleHideDotMenu() {
        if (this.$refs.dotMenu) {
          this.$refs.dotMenu.$el.style.opacity = 0
        }
      },
      handleAddProcess() {
        Form.show({
          type: 'create',
          title: `${this.$t('添加进程')}(${this.instance.name})`,
          hostId: this.instance.bk_host_id,
          bizId: this.bizId,
          submitHandler: this.createSubmitHandler
        })
      },
      async createSubmitHandler(values) {
        try {
          await this.$store.dispatch('processInstance/createServiceInstanceProcess', {
            params: {
              bk_biz_id: this.bizId,
              service_instance_id: this.instance.id,
              processes: [{
                process_info: values
              }]
            }
          })
          this.getServiceProcessList()
          this.updateInstanceInfo()
        } catch (error) {
          console.error(error)
        }
      },
      showProcessDetails(row) {
        Form.show({
          type: 'view',
          title: this.$t('查看进程'),
          instance: row.property,
          hostId: row.relation.bk_host_id,
          bizId: this.bizId,
          serviceTemplateId: this.instance.service_template_id,
          processTemplateId: row.relation.process_template_id,
          submitHandler: this.editSubmitHandler,
          showOptions: false
        })
      },
      handleEdit(row) {
        Form.show({
          type: 'update',
          title: this.$t('编辑进程'),
          instance: row.property,
          hostId: row.relation.bk_host_id,
          bizId: this.bizId,
          serviceTemplateId: this.instance.service_template_id,
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
          this.getServiceProcessList()
          this.updateInstanceInfo()
        } catch (error) {
          console.error(error)
        }
      },
      handleDelete(row) {
        this.$bkInfo({
          title: this.$t('确定删除该进程'),
          confirmFn: async () => {
            try {
              await this.$store.dispatch('processInstance/deleteServiceInstanceProcess', {
                config: {
                  data: {
                    bk_biz_id: this.bizId,
                    process_instance_ids: [row.property.bk_process_id]
                  },
                  requestId: this.requestId.deleteProcess
                }
              })
              this.getServiceProcessList()
              this.updateInstanceInfo()
            } catch (error) {
              console.error(error)
            }
          }
        })
      },
      redirectToTemplate() {
        this.$routerActions.redirect({
          name: MENU_BUSINESS_SERVICE_TEMPLATE_DETAILS,
          params: {
            bizId: this.bizId,
            templateId: this.instance.service_template_id
          },
          history: true
        })
      },
      updateInstanceInfo() {
        this.$emit('update-instance')
      }
    }
  }
</script>

<style lang="scss" scoped>
    .table-layout {
        padding: 0 0 12px 0;
    }
    .table-title {
        display: flex;
        align-items: center;
        height: 40px;
        padding: 0 10px;
        line-height: 40px;
        border-radius: 2px 2px 0 0;
        background-color: #DCDEE5;
        overflow: hidden;
        cursor: pointer;

        .title-checkbox {
            flex: none;
            /deep/ .bk-checkbox {
                background-color: #fff;
            }
            &.is-checked {
                /deep/ .bk-checkbox {
                    background-color: #3a84ff !important;
                }
            }
        }

        .right-content {
            display: flex;
            flex-wrap: nowrap;
            align-items: center;
            max-width: 75%;
            margin-left: auto;
        }

        .title-icon {
            font-size: 14px;
            margin: 0 2px 0 6px;
            color: #63656E;
            @include inlineBlock;
        }
        .title-label {
            font-size: 14px;
            color: #313238;
            min-width: 90px;
            @include inlineBlock;
            @include ellipsis;
        }
        .topology-path {
            font-size: 12px;
            color: #63656e;
            height: 20px;
            line-height: 20px;
            padding: 0 6px;
            outline: none;
            @include inlineBlock;
            overflow: hidden;
            text-overflow: ellipsis;
            white-space: nowrap;
            direction: rtl;
            &:hover {
                color: #3a84ff;
            }
        }
        .edit-form {
            width: 420px;
            @include inlineBlock;
        }
    }
    .instance-menu {
        flex: none;
        opacity: 0;
        /deep/ .bk-tooltip-ref {
            width: 100%;
        }
    }
    .menu-list {
        min-width: 74px;
        padding: 6px 0;
        .menu-item {
            .menu-button {
                display: block;
                width: 100%;
                height: 32px;
                padding: 0 13px;
                line-height: 32px;
                outline: 0;
                border: none;
                text-align: left;
                color: #63656E;
                font-size: 12px;
                background-color: #fff;
                &:hover {
                    background-color: #e1ecff;
                    color: #3a84ff;
                }
                &[disabled] {
                    color: #dcdee5;
                }
            }
            .menu-span {
                display: inline-block;
                width: 100%;
            }
        }
    }
    .instance-label {
        flex: none;
        @include inlineBlock;
        font-size: 12px;
        .icon-cc-label {
            color: #979ba5;
            font-size: 16px;
        }
        .label-list {
            padding-left: 4px;
            line-height: 38px;
            font-size: 0;
            .label-item {
                @include inlineBlock;
                font-size: 12px;
                height: 20px;
                line-height: 20px;
                margin-right: 4px;
                padding: 0 6px;
                color: #979ba5;
                background-color: #fafbfd;
                border-radius: 2px;
                &>span {
                    @include ellipsis;
                    display: inline-block;
                    max-width: 54px;
                }
            }
            .label-tips {
                padding: 0;
                /deep/ .bk-tooltip-ref {
                    padding: 0 6px;
                    span {
                        line-height: 16px;
                        display: inline-block;
                        vertical-align: top;
                    }
                }
                &:hover {
                    background-color: #e1ecff;
                }
            }
        }
    }
    .process-list-table {
        /deep/ {
            .bk-table-empty-block {
                min-height: 42px;
                .bk-table-empty-text {
                    width: auto;
                    padding: 0;
                }
            }
        }
    }
    .process-count-tips {
        display: flex;
        align-items: center;
        .tips-icon {
            color: $warningColor;
            font-size: 14px;
        }
        .tips-content {
            padding: 0 4px;
            color: $textDisabledColor;
            .tips-link {
                color: $primaryColor;
                cursor: pointer;
                &.disabled {
                    color: $textDisabledColor;
                }
            }
        }
    }
</style>

<style lang="scss">
    .tips-label-list {
        .label-item {
            @include inlineBlock;
            font-size: 12px;
            height: 20px;
            line-height: 18px;
            margin: 5px 2px;
            padding: 0 6px;
            color: #979ba5;
            background-color: #fafbfd;
            border: 1px solid #dcdee5;
            border-radius: 2px;
            &>span {
                @include ellipsis;
                display: inline-block;
                max-width: 54px;
            }
        }
    }
    .tippy-tooltip.label-tips-theme {
        padding: 8px 6px !important;
    }
</style>
