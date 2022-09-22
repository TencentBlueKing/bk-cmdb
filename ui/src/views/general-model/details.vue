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
  <div class="details-layout">
    <bk-tab class="details-tab"
      type="unborder-card"
      :active.sync="active">
      <bk-tab-panel name="property" :label="$t('属性')">
        <cmdb-property
          v-if="propertyListActive"
          :properties="properties"
          :property-groups="propertyGroups"
          @after-update="handleAfterUpdate"
          :inst="inst">
        </cmdb-property>
      </bk-tab-panel>
      <bk-tab-panel name="association" :label="$t('关联')">
        <cmdb-relation
          v-if="active === 'association'"
          :obj-id="objId"
          :inst="inst">
        </cmdb-relation>
      </bk-tab-panel>
      <bk-tab-panel name="history" :label="$t('变更记录')">
        <cmdb-audit-history v-if="active === 'history'"
          resource-type="model_instance"
          :obj-id="objId"
          :resource-id="instId">
        </cmdb-audit-history>
      </bk-tab-panel>
    </bk-tab>
  </div>
</template>

<script>
  import { mapGetters, mapActions } from 'vuex'
  import cmdbProperty from '@/components/model-instance/property'
  import cmdbAuditHistory from '@/components/model-instance/audit-history'
  import cmdbRelation from '@/components/model-instance/relation'
  import instanceService from '@/service/instance/instance'
  export default {
    components: {
      cmdbProperty,
      cmdbAuditHistory,
      cmdbRelation
    },
    data() {
      return {
        inst: {},
        properties: [],
        propertyListActive: true,
        propertyGroups: [],
        active: this.$route.query.tab || 'property'
      }
    },
    computed: {
      ...mapGetters(['supplierAccount', 'userName']),
      ...mapGetters('objectModelClassify', ['models', 'getModelById']),
      instId() {
        return parseInt(this.$route.params.instId, 10)
      },
      objId() {
        return this.$route.params.objId
      },
      model() {
        return this.getModelById(this.objId) || {}
      }
    },
    watch: {
      instId() {
        this.getData()
      },
      objId() {
        this.getData()
      }
    },
    created() {
      this.getData()
    },
    methods: {
      ...mapActions('objectModelFieldGroup', ['searchGroup']),
      ...mapActions('objectModelProperty', ['searchObjectAttribute']),
      setBreadcrumbs(inst) {
        this.$store.commit('setTitle', `${this.model.bk_obj_name}【${inst.bk_inst_name}】`)
      },
      async getData() {
        try {
          const [inst, properties, propertyGroups] = await Promise.all([
            this.getInstInfo(),
            this.getProperties(),
            this.getPropertyGroups()
          ])

          this.inst = inst
          this.properties = properties
          this.propertyGroups = propertyGroups

          this.setBreadcrumbs(inst)
        } catch (e) {
          console.error(e)
        }
      },
      getInstInfo() {
        return instanceService.findOne({ bk_obj_id: this.objId, bk_inst_id: this.instId })
      },
      async getProperties() {
        try {
          const properties = await this.searchObjectAttribute({
            params: {
              bk_obj_id: this.objId,
              bk_supplier_account: this.supplierAccount
            },
            config: {
              requestId: `post_searchObjectAttribute_${this.objId}`,
              fromCache: false
            }
          })

          return properties
        } catch (e) {
          console.error(e)
        }
      },
      async getPropertyGroups() {
        try {
          const propertyGroups = this.searchGroup({
            objId: this.objId,
            params: {},
            config: {
              fromCache: false,
              requestId: `post_searchGroup_${this.objId}`
            }
          })

          return propertyGroups
        } catch (e) {
          console.error(e)
        }
      },
      // 单个属性变更后可能会引起模型的权限变更，需要重载组件获取新的权限，避免出现权限已变更但 UI 仍然显示旧权限的情况。
      // 因为 property 组件的重载会让数据变为初始化数据，所以需要手动重新获取实例数据，以便更新为修改后的数据。
      async handleAfterUpdate() {
        try {
          const inst = await this.getInstInfo()
          this.inst = inst
          this.propertyListActive = false
          this.$nextTick(() => {
            this.propertyListActive = true
          })
        } catch (error) {
          console.log(error)
        }
      }
    }
  }
</script>

<style lang="scss" scoped>
    .details-layout {
        overflow: hidden;
        .details-tab {
            min-height: 400px;
            /deep/ {
                .bk-tab-header {
                    padding: 0;
                    margin: 0 20px;
                }
                .bk-tab-section {
                    @include scrollbar-y;
                    padding-bottom: 10px;
                }
            }
        }
    }
</style>
