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
          resource-type="project"
          :properties="properties"
          :property-groups="propertyGroups"
          :obj-id="objId"
          :inst="inst"
          :show-default-value="true">
        </cmdb-property>
      </bk-tab-panel>
      <bk-tab-panel name="association" :label="$t('关联')">
        <cmdb-relation
          v-if="active === 'association'"
          resource-type="project"
          :obj-id="objId"
          :inst="inst">
        </cmdb-relation>
      </bk-tab-panel>
      <bk-tab-panel name="history" :label="$t('变更记录')">
        <cmdb-audit-history v-if="active === 'history'"
          resource-type="project"
          :obj-id="objId"
          :biz-id="projId"
          :resource-id="projId">
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
  import { BUILTIN_MODELS } from '@/dictionary/model-constants.js'
  import projectService from '@/service/project/index.js'

  export default {
    components: {
      cmdbProperty,
      cmdbAuditHistory,
      cmdbRelation
    },
    data() {
      return {
        inst: {},
        objId: BUILTIN_MODELS.PROJECT,
        properties: [],
        propertyGroups: [],
        active: this.$route.query.tab || 'property'
      }
    },
    computed: {
      ...mapGetters(['supplierAccount', 'userName']),
      ...mapGetters('objectModelClassify', ['getModelById']),
      projId() {
        return parseInt(this.$route.params.projId, 10)
      },
      model() {
        return this.getModelById(BUILTIN_MODELS.PROJECT) || {}
      }
    },
    watch: {
      projId() {
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
        this.$store.commit('setTitle', `${this.model.bk_obj_name}【${inst.bk_project_name}】`)
      },
      async getData() {
        const id = Number(this.$route.params.projId)
        try {
          const config = { requestId: `post_searchProjectById_${this.projId}`, cancelPrevious: true }
          const [inst, properties, propertyGroups] = await Promise.all([
            projectService.findOne({ id }, config),
            this.getProperties(),
            this.getPropertyGroups()
          ])
          this.inst = inst
          this.properties = properties.filter(item => item.bk_property_id !== 'bk_project_icon')
          this.propertyGroups = propertyGroups
          this.setBreadcrumbs(inst)
        } catch (e) {
          console.error(e)
        }
      },
      async getProperties() {
        try {
          const properties = await this.searchObjectAttribute({
            injectId: BUILTIN_MODELS.PROJECT,
            params: {
              bk_obj_id: this.objId
            },
            config: {
              requestId: 'post_searchObjectAttribute_bk_project',
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
              requestId: 'post_searchGroup_bk_project'
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
