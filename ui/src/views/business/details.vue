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
          resource-type="business"
          :properties="properties"
          :property-groups="propertyGroups"
          :inst="inst">
        </cmdb-property>
      </bk-tab-panel>
      <bk-tab-panel name="association" :label="$t('关联')">
        <cmdb-relation
          v-if="active === 'association'"
          resource-type="business"
          :obj-id="objId"
          :inst="inst">
        </cmdb-relation>
      </bk-tab-panel>
      <bk-tab-panel name="history" :label="$t('变更记录')">
        <cmdb-audit-history v-if="active === 'history'"
          resource-type="business"
          :obj-id="objId"
          :biz-id="bizId"
          :resource-id="bizId">
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
  export default {
    components: {
      cmdbProperty,
      cmdbAuditHistory,
      cmdbRelation
    },
    data() {
      return {
        inst: {},
        objId: 'biz',
        properties: [],
        propertyGroups: [],
        active: this.$route.query.tab || 'property'
      }
    },
    computed: {
      ...mapGetters(['supplierAccount', 'userName']),
      ...mapGetters('objectModelClassify', ['getModelById']),
      bizId() {
        return parseInt(this.$route.params.bizId, 10)
      },
      model() {
        return this.getModelById('biz') || {}
      }
    },
    watch: {
      bizId() {
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
      ...mapActions('objectBiz', [
        'searchBusinessById'
      ]),
      setBreadcrumbs(inst) {
        this.$store.commit('setTitle', `${this.model.bk_obj_name}【${inst.bk_biz_name}】`)
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
      async getInstInfo() {
        try {
          const inst = await this.searchBusinessById({
            bizId: this.bizId,
            config: { requestId: `post_searchBusinessById_${this.bizId}`, cancelPrevious: true }
          })
          return inst
        } catch (e) {
          console.error(e)
        }
      },
      async getProperties() {
        try {
          const properties = await this.searchObjectAttribute({
            params: {
              bk_obj_id: this.objId,
              bk_supplier_account: this.supplierAccount
            },
            config: {
              requestId: 'post_searchObjectAttribute_biz',
              fromCache: true
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
              fromCache: true,
              requestId: 'post_searchGroup_biz'
            }
          })

          return propertyGroups
        } catch (e) {
          console.error(e)
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
