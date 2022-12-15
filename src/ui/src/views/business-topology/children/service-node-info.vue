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
  <div class="node-info"
    v-bkloading="{
      isLoading: $loading([
        'getModelProperties',
        'getModelPropertyGroups',
        'getNodeInstance',
        'updateNodeInstance',
        'deleteNodeInstance',
        'removeServiceTemplate'
      ])
    }"
  >
    <cmdb-permission v-if="permission" class="permission-tips" :permission="permission"></cmdb-permission>
    <cmdb-details class="topology-details"
      v-else-if="type === 'details'"
      :class="{ pt10: !isSetNode && !isModuleNode }"
      :properties="properties"
      :property-groups="propertyGroups"
      :inst="instance"
      :show-copy="true"
      :show-options="modelId !== 'biz' && editable">
      <node-extra-info slot="prepend" :instance="instance"></node-extra-info>
      <template slot="details-options">
        <cmdb-auth :auth="{ type: $OPERATION.U_TOPO, relation: [business] }">
          <template slot-scope="{ disabled }">
            <bk-button class="button-edit" v-test-id="'edit'"
              theme="primary"
              :disabled="disabled"
              @click="handleEdit">
              {{$t('编辑')}}
            </bk-button>
          </template>
        </cmdb-auth>
        <cmdb-auth :auth="{ type: $OPERATION.D_TOPO, relation: [business] }">
          <template slot-scope="{ disabled }">
            <span class="inline-block-middle" v-if="moduleFromSetTemplate"
              v-bk-tooltips="$t('由集群模板创建的模块无法删除')">
              <bk-button class="btn-delete" hover-theme="danger" disabled>
                {{$t('删除节点')}}
              </bk-button>
            </span>
            <bk-button class="btn-delete" v-else v-test-id="'del'"
              hover-theme="danger"
              :disabled="disabled"
              @click="handleDelete">
              {{$t('删除节点')}}
            </bk-button>
          </template>
        </cmdb-auth>
      </template>
    </cmdb-details>
    <template v-else-if="type === 'update'">
      <cmdb-form class="topology-form"
        ref="form"
        :properties="properties"
        :property-groups="propertyGroups"
        :disabled-properties="disabledProperties"
        :inst="instance"
        :type="type"
        :save-auth="{ type: $OPERATION.U_TOPO, relation: [business] }"
        :is-main-line="true"
        @on-submit="handleSubmit"
        @on-cancel="handleCancel">
        <form-service-category slot="prepend"
          v-if="!withTemplate && isModuleNode"
          :instance="instance"
          @change="handleServiceCategoryChange">
        </form-service-category>

        <!-- 配置字段不允许修改，提示跳转到模板 -->
        <template v-for="config in propertyConfigList" #[config.property.bk_property_id]>
          <bk-popover
            class="config-form-tips"
            ext-cls="config-uneditable-popover"
            :key="config.property.bk_property_id">
            <component
              :is="`cmdb-form-${config.property.bk_property_type}`"
              :class="['form-element-item', config.property.bk_property_type]"
              :unit="config.property.unit"
              :options="config.property.option || []"
              :placeholder="getPlaceholder(config.property)"
              :auto-check="false"
              :disabled="true"
              v-model.trim="config.instance_value">
            </component>
            <template slot="content">
              <i18n path="模板定义不可修改提示">
                <template #link>
                  <bk-link theme="primary" @click="handleRedirect" class="goto-link">
                    {{$t(withTemplate ? '跳转服务模板' : '跳转集群模板')}}
                  </bk-link>
                </template>
              </i18n>
            </template>
          </bk-popover>
        </template>
      </cmdb-form>
    </template>
  </div>
</template>

<script>
  import has from 'has'
  import debounce from 'lodash.debounce'
  import Utils from '@/components/filters/utils'
  import NodeExtraInfo from './node-extra-info'
  import FormServiceCategory from './form-service-category'
  import instanceService from '@/service/instance/instance'
  import serviceTemplateService from '@/service/service-template'
  import setTemplateService from '@/service/set-template'
  import { MENU_BUSINESS_SERVICE_TEMPLATE_DETAILS, MENU_BUSINESS_SET_TEMPLATE_DETAILS } from '@/dictionary/menu-symbol'

  export default {
    components: {
      NodeExtraInfo,
      FormServiceCategory
    },
    props: {
      active: Boolean
    },
    data() {
      return {
        type: 'details',
        properties: [],
        disabledProperties: [],
        propertyGroups: [],
        instance: {},
        first: '',
        second: '',
        hasChange: false,
        templateInfo: {
          serviceTemplateName: this.$t('无'),
          serviceCategory: '',
          setTemplateName: this.$t('无')
        },
        propertyConfigList: [],
        refresh: null,
        permission: null
      }
    },
    computed: {
      business() {
        return this.$store.getters['objectBiz/bizId']
      },
      propertyMap() {
        return this.$store.state.businessHost.propertyMap
      },
      propertyGroupMap() {
        return this.$store.state.businessHost.propertyGroupMap
      },
      selectedNode() {
        return this.$store.state.businessHost.selectedNode
      },
      isModuleNode() {
        return this.selectedNode && this.selectedNode.data.bk_obj_id === 'module'
      },
      isSetNode() {
        return this.selectedNode && this.selectedNode.data.bk_obj_id === 'set'
      },
      isBlueking() {
        let rootNode = this.selectedNode || {}
        if (rootNode.parent) {
          // eslint-disable-next-line prefer-destructuring
          rootNode = rootNode.parents[0]
        }
        return rootNode.name === '蓝鲸'
      },
      modelId() {
        if (this.selectedNode) {
          return this.selectedNode.data.bk_obj_id
        }
        return null
      },
      nodeNamePropertyId() {
        if (this.modelId) {
          const map = {
            biz: 'bk_biz_name',
            set: 'bk_set_name',
            module: 'bk_module_name'
          }
          return map[this.modelId] || 'bk_inst_name'
        }
        return null
      },
      nodePrimaryPropertyId() {
        if (this.modelId) {
          const map = {
            biz: 'bk_biz_id',
            set: 'bk_set_id',
            module: 'bk_module_id'
          }
          return map[this.modelId] || 'bk_inst_id'
        }
        return null
      },
      withTemplate() {
        return this.isModuleNode && !!this.instance.service_template_id
      },
      withSetTemplate() {
        return this.isSetNode && !!this.instance.set_template_id
      },
      moduleFromSetTemplate() {
        return this.isModuleNode && !!this.selectedNode.parent.data.set_template_id
      },
      editable() {
        const editable = this.$store.state.businessHost.blueKingEditable
        return this.isBlueking ? this.isBlueking && editable : true
      },
      templateId() {
        if (this.withTemplate) {
          return this.instance.service_template_id
        }
        if (this.withSetTemplate) {
          return this.instance.set_template_id
        }
        return null
      }
    },
    watch: {
      modelId: {
        immediate: true,
        handler(modelId) {
          if (modelId && this.active) {
            this.initProperties()
          }
        }
      },
      selectedNode: {
        immediate: true,
        handler(node) {
          if (node && this.active) {
            this.getNodeDetails()
          }
        }
      },
      active: {
        immediate: true,
        handler(active) {
          if (active) {
            this.refresh && this.refresh()
          }
        }
      }
    },
    created() {
      this.refresh = debounce(() => {
        this.initProperties()
        this.getNodeDetails()
      }, 10)
    },
    methods: {
      async initProperties() {
        this.type = 'details'
        try {
          const [
            properties,
            groups
          ] = await Promise.all([
            this.getProperties(),
            this.getPropertyGroups()
          ])
          this.properties = properties
          this.propertyGroups = groups
        } catch (e) {
          console.error(e)
          this.properties = []
          this.propertyGroups = []
        }
      },
      async getNodeDetails() {
        this.type = 'details'
        await this.getInstance()

        // 获取节点的模板属性配置
        this.getPropertyConfig()

        this.disabledProperties = this.selectedNode.data.bk_obj_id === 'module' && this.withTemplate ? ['bk_module_name'] : []
      },
      async getProperties() {
        let properties = []
        const { modelId } = this
        if (has(this.propertyMap, modelId)) {
          properties = this.propertyMap[modelId]
        } else {
          const action = 'objectModelProperty/searchObjectAttribute'
          properties = await this.$store.dispatch(action, {
            params: {
              bk_biz_id: this.business,
              bk_obj_id: modelId,
              bk_supplier_account: this.$store.getters.supplierAccount
            },
            config: {
              requestId: 'getModelProperties'
            }
          })
          this.$store.commit('businessHost/setProperties', {
            id: modelId,
            properties
          })
        }
        this.insertNodeIdProperty(properties)
        return Promise.resolve(properties)
      },
      async getPropertyGroups() {
        let groups = []
        const { modelId } = this
        if (has(this.propertyGroupMap, modelId)) {
          groups = this.propertyGroupMap[modelId]
        } else {
          const action = 'objectModelFieldGroup/searchGroup'
          groups = await this.$store.dispatch(action, {
            objId: modelId,
            params: { bk_biz_id: this.business },
            config: {
              requestId: 'getModelPropertyGroups'
            }
          })
          this.$store.commit('businessHost/setPropertyGroups', {
            id: modelId,
            groups
          })
        }
        return Promise.resolve(groups)
      },
      async getInstance() {
        try {
          const { modelId } = this
          const promiseMap = {
            biz: this.getBizInstance,
            set: this.getSetInstance,
            module: this.getModuleInstance
          }
          this.instance = await (promiseMap[modelId] || this.getCustomInstance)()
          this.$store.commit('businessHost/setSelectedNodeInstance', this.instance)
        } catch (e) {
          console.error(e)
          this.instance = {}
        }
      },
      async getBizInstance() {
        try {
          const data = await this.$store.dispatch('objectBiz/searchBusiness', {
            params: {
              page: { start: 0, limit: 1 },
              fields: [],
              condition: {
                bk_biz_id: { $eq: this.selectedNode.data.bk_inst_id }
              }
            },
            config: {
              requestId: 'getNodeInstance',
              cancelPrevious: true,
              globalPermission: false
            }
          })
          return data.info[0]
        } catch ({ permission }) {
          this.permission = permission
          return {}
        }
      },
      async getSetInstance() {
        const { info: [instance] } = await this.$store.dispatch('objectSet/searchSet', {
          bizId: this.business,
          params: {
            page: { start: 0, limit: 1 },
            fields: [],
            condition: {
              bk_set_id: this.selectedNode.data.bk_inst_id
            }
          },
          config: {
            requestId: 'getNodeInstance',
            cancelPrevious: true
          }
        })
        return instance
      },
      async getPropertyConfig() {
        if (!this.withTemplate && !this.withSetTemplate) {
          return
        }

        const templateService = this.isModuleNode ? serviceTemplateService : setTemplateService

        const { attributes } = await templateService.findProperty({
          id: this.templateId,
          bk_biz_id: this.business
        })

        // 在配置数据中注入属性对象和实例对应的字段值
        const values = this.$tools.getInstFormValues(this.properties, this.instance, false)
        this.propertyConfigList = attributes.map((item) => {
          const property = this.properties.find(prop => prop.id === item.bk_attribute_id) || {}
          return {
            property,
            instance_value: values[property.bk_property_id],
            ...item
          }
        })
      },
      insertNodeIdProperty(properties) {
        if (properties.find(property => property.bk_property_id === this.nodePrimaryPropertyId)) {
          return
        }
        const nodeNameProperty = properties.find(property => property.bk_property_id === this.nodeNamePropertyId)
        properties.unshift({
          ...nodeNameProperty,
          id: Date.now(),
          bk_property_id: this.nodePrimaryPropertyId,
          bk_property_name: this.$t('ID'),
          editable: false
        })
      },
      async getModuleInstance() {
        const data = await this.$store.dispatch('objectModule/searchModule', {
          bizId: this.business,
          setId: this.selectedNode.parent.data.bk_inst_id,
          params: {
            page: { start: 0, limit: 1 },
            fields: [],
            condition: {
              bk_module_id: this.selectedNode.data.bk_inst_id,
              bk_supplier_account: this.$store.getters.supplierAccount
            }
          },
          config: {
            requestId: 'getNodeInstance',
            cancelPrevious: true
          }
        })
        // eslint-disable-next-line prefer-destructuring
        const instance = data.info[0]
        return instance
      },
      async getCustomInstance() {
        return instanceService.findOne({
          bk_obj_id: this.modelId,
          bk_inst_id: this.selectedNode.data.bk_inst_id
        })
      },
      getPlaceholder(property) {
        return Utils.getPlaceholder(property)
      },
      handleRedirect() {
        this.$routerActions.redirect({
          name: this.withTemplate ? MENU_BUSINESS_SERVICE_TEMPLATE_DETAILS : MENU_BUSINESS_SET_TEMPLATE_DETAILS,
          params: {
            bizId: this.business,
            templateId: this.templateId
          },
          history: true
        })
      },
      handleEdit() {
        this.type = 'update'
      },
      handleServiceCategoryChange(id) {
        this.$set(this.$refs.form.values, 'service_category_id', id)
      },
      async handleSubmit(value) {
        if (!await this.$validator.validateAll()) return
        const promiseMap = {
          set: this.updateSetInstance,
          module: this.updateModuleInstance
        }
        const nameMap = {
          set: 'bk_set_name',
          module: 'bk_module_name'
        }
        try {
          await (promiseMap[this.modelId] || this.updateCustomInstance)({ ...value, bk_biz_id: this.business })
          this.selectedNode.data.bk_inst_name = value[nameMap[this.modelId] || 'bk_inst_name']
          this.instance = Object.assign({}, this.instance, value)
          this.type = 'details'
          this.$success(this.$t('修改成功'))
        } catch (e) {
          console.error(e)
        }
      },
      updateSetInstance(value) {
        return this.$store.dispatch('objectSet/updateSet', {
          bizId: this.business,
          setId: this.selectedNode.data.bk_inst_id,
          params: { ...value },
          config: {
            requestId: 'updateNodeInstance'
          }
        })
      },
      updateModuleInstance(value) {
        return this.$store.dispatch('objectModule/updateModule', {
          bizId: this.business,
          setId: this.selectedNode.parent.data.bk_inst_id,
          moduleId: this.selectedNode.data.bk_inst_id,
          params: {
            bk_supplier_account: this.$store.getters.supplierAccount,
            ...value
          },
          config: {
            requestId: 'updateNodeInstance'
          }
        })
      },
      updateCustomInstance(value) {
        return this.$store.dispatch('objectCommonInst/updateInst', {
          objId: this.modelId,
          instId: this.selectedNode.data.bk_inst_id,
          params: { ...value },
          config: {
            requestId: 'updateNodeInstance'
          }
        })
      },
      handleCancel() {
        this.type = 'details'
      },
      async handleDelete() {
        const count = await this.getSelectedNodeHostCount()
        if (count) {
          this.$error(this.$t('目标包含主机，不允许删除'))
          return
        }
        this.$bkInfo({
          title: `${this.$t('确定删除')} ${this.selectedNode.name}?`,
          confirmFn: async () => {
            const promiseMap = {
              set: this.deleteSetInstance,
              module: this.deleteModuleInstance
            }
            try {
              await (promiseMap[this.modelId] || this.deleteCustomInstance)()
              const { tree } = this.selectedNode
              const parentId = this.selectedNode.parent.id
              const nodeId = this.selectedNode.id
              tree.setSelected(parentId, {
                emitEvent: true
              })
              tree.removeNode(nodeId)
              this.$success(this.$t('删除成功'))
            } catch (e) {
              console.error(e)
            }
          }
        })
      },
      deleteSetInstance() {
        return this.$store.dispatch('objectSet/deleteSet', {
          bizId: this.business,
          setId: this.selectedNode.data.bk_inst_id,
          config: {
            requestId: 'deleteNodeInstance',
            data: { bk_biz_id: this.business }
          }
        })
      },
      deleteModuleInstance() {
        return this.$store.dispatch('objectModule/deleteModule', {
          bizId: this.business,
          setId: this.selectedNode.parent.data.bk_inst_id,
          moduleId: this.selectedNode.data.bk_inst_id,
          config: {
            requestId: 'deleteNodeInstance',
            data: { bk_biz_id: this.business }
          }
        })
      },
      deleteCustomInstance() {
        return this.$store.dispatch('objectCommonInst/deleteInst', {
          objId: this.modelId,
          instId: this.selectedNode.data.bk_inst_id,
          config: {
            requestId: 'deleteNodeInstance',
            data: { bk_biz_id: this.business }
          }
        })
      },
      async getSelectedNodeHostCount() {
        const defaultModel = ['biz', 'set', 'module', 'host', 'object']
        const modelInstKey = {
          biz: 'bk_biz_id',
          set: 'bk_set_id',
          module: 'bk_module_id',
          host: 'bk_host_id',
          object: 'bk_inst_id'
        }
        const conditionParams = {
          condition: defaultModel.map(model => ({
            bk_obj_id: model,
            condition: [],
            fields: []
          }))
        }
        const { selectedNode } = this
        const selectedModel = defaultModel.includes(selectedNode.data.bk_obj_id) ? selectedNode.data.bk_obj_id : 'object'
        const selectedModelCondition = conditionParams.condition.find(model => model.bk_obj_id === selectedModel)
        selectedModelCondition.condition.push({
          field: modelInstKey[selectedModel],
          operator: '$eq',
          value: selectedNode.data.bk_inst_id
        })
        const data = await this.$store.dispatch('hostSearch/searchHost', {
          params: {
            ...conditionParams,
            bk_biz_id: this.business,
            ip: {
              flag: 'bk_host_innerip|bk_host_outer',
              exact: 0,
              data: []
            },
            page: {
              start: 0,
              limit: 1,
              sort: ''
            }
          },
          config: {
            requestId: 'searchHosts',
            cancelPrevious: true
          }
        })
        return data && data.count
      }
    }
  }
</script>

<style lang="scss" scoped>
    .permission-tips {
        position: absolute;
        top: 35%;
        left: 50%;
        transform: translate(-50%, -50%);
    }
    .node-info {
        height: 100%;
        margin: 0 -20px;
    }
    .topology-details.details-layout {
        padding: 0;
        /deep/ {
            .property-group {
                padding-left: 36px;
            }
            .property-list {
                padding-left: 20px;
            }
            .details-options {
                width: 100%;
                margin: 0;
                padding-left: 56px;
            }
        }
    }
    .topology-form {
        /deep/ {
            .form-groups {
                padding: 0 20px;
            }
            .property-list {
                margin-left: 36px;
            }
            .property-item {
                max-width: 554px !important;
            }
            .form-options {
                padding: 10px 0 10px 56px;
            }
        }
    }
    .button-edit {
        min-width: 76px;
        margin-right: 4px;
    }
    .btn-delete {
        min-width: 76px;
    }

    .config-form-tips {
      display: block;
      width: 100%;
      ::v-deep .bk-tooltip-ref {
        display: block;
      }
    }
</style>
<style lang="scss">
  .config-uneditable-popover {
    .goto-link {
      vertical-align: unset;
      .bk-link-text {
        font-size: 12px;
      }
    }
  }
</style>
