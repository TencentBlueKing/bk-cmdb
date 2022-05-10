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
  <div
    class="node-info"
    v-bkloading="{
      isLoading: $loading([
        'getModelProperties',
        'getModelPropertyGroups',
        'getNodeInstance'
      ])
    }"
  >
    <cmdb-permission
      v-if="permission"
      class="permission-tips"
      :permission="permission"
    ></cmdb-permission>

    <cmdb-details
      class="topology-details"
      v-else
      :class="{ pt10: !isSetNode && !isModuleNode }"
      :properties="properties"
      :property-groups="propertyGroups"
      :inst="instance"
      :show-options="false"
    >
      <node-extra-info slot="prepend" :instance="instance"></node-extra-info>
    </cmdb-details>
  </div>
</template>

<script>
  import { mapState } from 'vuex'
  import has from 'has'
  import debounce from 'lodash.debounce'
  import NodeExtraInfo from '@/views/business-topology/children/node-extra-info.vue'
  import instanceService from '@/service/instance/instance'
  import { BusinessService } from '@/service/business-set/business.js'
  import { SetService } from '@/service/business-set/set.js'
  import { ModuleService } from '@/service/business-set/module.js'
  import to from 'await-to-js'
  import { BUILTIN_MODELS } from '@/dictionary/model-constants.js'

  export default {
    name: 'nodeInfo',
    components: {
      NodeExtraInfo,
    },
    props: {
      active: Boolean,
    },
    data() {
      return {
        properties: [],
        disabledProperties: [],
        propertyGroups: [],
        instance: {},
        refresh: null,
        permission: null,
      }
    },
    computed: {
      ...mapState('bizSet', ['bizSetId', 'bizId']),
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
        return this.selectedNode?.data?.bk_obj_id === BUILTIN_MODELS.MODULE
      },
      isSetNode() {
        return this.selectedNode?.data?.bk_obj_id === BUILTIN_MODELS.SET
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
            module: 'bk_module_name',
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
            module: 'bk_module_id',
          }
          return map[this.modelId] || 'bk_inst_id'
        }
        return null
      },
      withTemplate() {
        return this.isModuleNode && !!this.instance.service_template_id
      },
    },
    watch: {
      modelId: {
        immediate: true,
        handler(modelId) {
          if (modelId && this.active) {
            this.initProperties()
          }
        },
      },
      selectedNode: {
        immediate: true,
        handler(node) {
          if (node && this.active) {
            this.getNodeDetails()
          }
        },
      },
      active: {
        immediate: true,
        handler(active) {
          if (active) {
            this.refresh && this.refresh()
          }
        },
      },
    },
    created() {
      this.refresh = debounce(() => {
        this.initProperties()
        this.getNodeDetails()
      }, 10)
    },
    methods: {
      async initProperties() {
        try {
          const [properties, groups] = await Promise.all([
            this.getProperties(),
            this.getPropertyGroups(),
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
        await this.getInstance()
        this.disabledProperties = this.selectedNode.data.bk_obj_id === BUILTIN_MODELS.MODULE && this.withTemplate
          ? ['bk_module_name']
          : []
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
              bk_biz_id: this.bizId,
              bk_obj_id: modelId,
              bk_supplier_account: this.$store.getters.supplierAccount,
            },
            config: {
              requestId: 'getModelProperties',
            },
          })
          this.$store.commit('businessHost/setProperties', {
            id: modelId,
            properties,
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
            params: { bk_biz_id: this.bizId },
            config: {
              requestId: 'getModelPropertyGroups',
            },
          })
          this.$store.commit('businessHost/setPropertyGroups', {
            id: modelId,
            groups,
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
            module: this.getModuleInstance,
          }
          this.instance = await (promiseMap[modelId] || this.getCustomInstance)()
          this.$store.commit(
            'businessHost/setSelectedNodeInstance',
            this.instance
          )
        } catch (e) {
          console.error(e)
          this.instance = {}
        }
      },
      async getBizInstance() {
        const [err, { info: [instance] }] = await to(BusinessService.findOne(
          {
            bizSetId: this.bizSetId,
            bizId: this.bizId
          },
          {
            requestId: 'getNodeInstance',
            cancelPrevious: true,
            globalPermission: false,
          }
        ))

        if (err?.permission) {
          this.permission = err?.permission
          return {}
        }

        return instance
      },
      async getSetInstance() {
        const setId = this.selectedNode.data.bk_inst_id
        const [, { info: [instance] }] = await to(SetService.findOne(
          {
            bizSetId: this.bizSetId,
            bizId: this.bizId
          },
          {
            page: { start: 0, limit: 1 },
            fields: [],
            condition: {
              bk_set_id: setId,
            },
          },
          {
            requestId: 'getNodeInstance',
            cancelPrevious: true,
          }
        ))
        return instance
      },
      insertNodeIdProperty(properties) {
        if (
          properties.find(property => property.bk_property_id === this.nodePrimaryPropertyId)
        ) {
          return
        }
        const nodeNameProperty = properties.find(property => property.bk_property_id === this.nodeNamePropertyId)
        properties.unshift({
          ...nodeNameProperty,
          bk_property_id: this.nodePrimaryPropertyId,
          bk_property_name: this.$t('ID'),
        })
      },
      async getModuleInstance() {
        const setId = this.selectedNode.parent.data.bk_inst_id
        const moduleId = this.selectedNode.data.bk_inst_id
        const [, { info: [instance] }] = await to(ModuleService.findOne(
          {
            bizSetId: this.bizSetId,
            bizId: this.bizId,
            setId,
          },
          {
            page: { start: 0, limit: 1 },
            fields: [],
            condition: {
              bk_module_id: moduleId,
              bk_supplier_account: this.$store.getters.supplierAccount,
            },
          },
          {
            requestId: 'getNodeInstance',
            cancelPrevious: true,
          }
        ))
        return instance
      },
      async getCustomInstance() {
        return instanceService.findOne({
          bk_obj_id: this.modelId,
          bk_inst_id: this.selectedNode.data.bk_inst_id,
        })
      }
    },
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
      padding: 10px 0 0 36px;
    }
  }
}
</style>
