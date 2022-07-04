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
  <cmdb-sticky-layout
    class="service-template-sync-layout"
    v-bkloading="{ isLoading: $loading([requestIds.properties, requestIds.topopath]) }">
    <template #header="{ sticky }">
      <cmdb-tips :class="['layout-header', { 'is-sticky': sticky }]">{{$t('同步模板功能提示')}}</cmdb-tips>
    </template>

    <div class="layout-top">
      <p class="title" v-if="isSingleModule">{{$t('请确认单个实例更改信息')}}</p>
      <i18n path="请确认实例更改信息"
        tag="p"
        class="title"
        v-else>
        <template #count>
          <span>{{moduleIds.length}}</span>
        </template>
      </i18n>
      <div class="type-legend">
        <span class="legend-item">
          <i class="dot changed"></i>
          {{$t('变更')}}
        </span>
        <span class="legend-item">
          <i class="dot added"></i>
          {{$t('新增')}}
        </span>
        <span class="legend-item">
          <i class="dot removed"></i>
          {{$t('删除')}}
        </span>
      </div>
    </div>

    <div class="layout-main">
      <div class="module-instance-container" v-if="isSingleModule">
        <module-instance
          v-bkloading="{ isLoading: moduleGroup[moduleIds[0]].loading }"
          :module-id="moduleIds[0]"
          :template-id="templateId"
          :topo-path="moduleGroup[moduleIds[0]].topoPath"
          :model-property="modelProperty"
          :property-diff="moduleGroup[moduleIds[0]].propertyDiff"
          :process-diff="moduleGroup[moduleIds[0]].processDiff">
        </module-instance>
      </div>
      <div class="module-instance-group" v-else>
        <cmdb-collapse class="module-instance-container"
          v-for="moduleId in moduleIds"
          :label="moduleGroup[moduleId].topoPath"
          :collapse="moduleGroup[moduleId].collapse"
          arrow-type="filled"
          :key="moduleId"
          @collapse-change="handleModuleCollapseChange(moduleId, $event)">
          <module-instance
            v-bkloading="{ isLoading: moduleGroup[moduleId].loading }"
            class="module-instance-item"
            collapse-size="small"
            :module-id="Number(moduleId)"
            :template-id="templateId"
            :topo-path="moduleGroup[moduleId].topoPath"
            :model-property="modelProperty"
            :property-diff="moduleGroup[moduleId].propertyDiff"
            :process-diff="moduleGroup[moduleId].processDiff">
          </module-instance>
        </cmdb-collapse>
      </div>
    </div>

    <template #footer="{ sticky }">
      <div :class="['layout-footer', { 'is-sticky': sticky }]">
        <cmdb-auth :auth="{ type: $OPERATION.C_SERVICE_INSTANCE, relation: [bizId] }">
          <template #default="{ disabled }">
            <bk-button theme="primary"
              :disabled="disabled"
              :loading="confirming"
              @click="confirmAndSync">
              {{$t('确认同步')}}
            </bk-button>
          </template>
        </cmdb-auth>
        <bk-button @click="handleGoback">{{$t('取消')}}</bk-button>
      </div>
    </template>
  </cmdb-sticky-layout>
</template>

<script>
  import { mapGetters } from 'vuex'
  import to from 'await-to-js'
  import ModuleInstance from './children/module-instance.vue'

  export default {
    name: 'BusinessSynchronous',
    components: {
      ModuleInstance,
    },
    data() {
      const moduleIds = String(this.$route.params.modules).split(',')
        .map(id => Number(id))

      const moduleGroup = {}
      moduleIds.forEach((id) => {
        moduleGroup[id] = {
          collapse: true, // 是否展开
          topoPath: '', // 拓扑路径
          processDiff: [], // 进程对比数据
          propertyDiff: [], // 属性对比数据
          loaded: false, // 是否加载过
          loading: false // 是否加载中
        }
      })

      return {
        moduleIds,
        moduleGroup, // 按模块分组的数据
        modelProperty: {}, // 资源的所有属性，用来翻译
        confirming: false,
        requestIds: {
          properties: Symbol(),
          topopath: Symbol(),
          difference: Symbol()
        }
      }
    },
    computed: {
      ...mapGetters(['supplierAccount']),
      ...mapGetters('objectBiz', ['bizId']),
      templateId() {
        return Number(this.$route.params.template)
      },
      isSingleModule() {
        return this.moduleIds?.length === 1
      }
    },
    async created() {
      await to(this.loadProperties())
      await to(this.loadTopoPath())

      if (this.moduleIds?.length > 1) {
        this.$store.commit('setTitle', this.$t('批量同步模板'))
      }

      // 默认展开第1个
      this.loadDiffByModule(this.moduleIds[0])
    },
    methods: {
      /**
       * 加载进程属性，便于转换成可读中文
       */
      loadProperties() {
        return this.$store.dispatch('objectModelProperty/batchSearchObjectAttribute', {
          params: {
            bk_biz_id: this.bizId,
            bk_obj_id: { $in: ['process', 'module'] },
            bk_supplier_account: this.supplierAccount
          },
          config: {
            requestId: this.requestIds.properties,
            fromCache: true
          }
        }).then((data) => {
          this.modelProperty = data
        })
          .catch(() => {
            this.modelProperty = {}
          })
      },
      /**
       * 加载拓扑路径，用于加载涉及实例
       */
      loadTopoPath() {
        return this.$store.dispatch('objectMainLineModule/getTopoPath', {
          bizId: this.bizId,
          params: {
            topo_nodes: this.moduleIds.map(moduleId => ({ bk_obj_id: 'module', bk_inst_id: moduleId }))
          },
          config: { requestId: this.requestIds.topopath }
        }).then(({ nodes }) => {
          nodes.forEach((node) => {
            const moduleId = node.topo_node.bk_inst_id
            this.moduleGroup[moduleId].topoPath = node.topo_path.reverse().map(path => path.bk_inst_name)
              .join(' / ')
          })
        })
      },
      /**
       * 按模块获取diff信息（进程视角）
       */
      loadDiffByModule(moduleId) {
        const currentModule = this.moduleGroup[moduleId]

        currentModule.loading = true
        currentModule.collapse = false

        return this.$store.dispatch('businessSynchronous/getTplDiffs', {
          params: {
            bk_module_id: moduleId,
            bk_biz_id: this.bizId,
            service_template_id: this.templateId
          },
          config: { requestId: this.requestIds.difference }
        }).then((difference) => {
          const processDiff = []
          const processDiffTypes = ['changed', 'added', 'removed']

          // 进程变更
          Object.keys(difference).forEach((type) => {
            const diffItem = difference[type]
            if (processDiffTypes.includes(type) && diffItem) {
              diffItem.forEach(({ id, name }) => {
                processDiff.push(this.genProcessDiffItem({
                  diffType: type,
                  processId: id,
                  processName: name,
                }))
              })
            }
          })

          // 属性变更，注入原始属性对象
          if (difference.attributes) {
            currentModule.propertyDiff = difference.attributes.map((attr) => {
              const property = this.modelProperty.module.find(prop => prop.id === attr.id)
              return {
                property,
                ...attr
              }
            })
          }

          currentModule.processDiff = processDiff
        })
          .finally(() => {
            currentModule.loading = false
            currentModule.loaded = true
          })
      },
      /**
       * 生成进程变更对比项
       * @param {string} diffType 必须，变更类型
       * @param {string} processId 非必须，进程模板的变更 ID
       * @param {string} processName 非必须，进程模板的名称
       */
      genProcessDiffItem({
        diffType,
        processId,
        processName
      }) {
        return {
          type: diffType,
          process_template_id: processId,
          process_template_name: processName,
          confirmed: false
        }
      },
      confirmAndSync() {
        this.confirming = true
        this.$store.dispatch('businessSynchronous/syncServiceInstanceByTemplate', {
          params: {
            service_template_id: this.templateId,
            bk_module_ids: this.moduleIds,
            bk_biz_id: this.bizId
          }
        }).then(() => {
          this.$success(this.$t('提交同步成功'))
          this.goBackModule()
        })
          .finally(() => {
            this.confirming = false
          })
      },
      handleModuleCollapseChange(moduleId, collapse) {
        // 打开并且未加载过或者不在加载中状态
        if (!collapse && !this.moduleGroup[moduleId].loaded && !this.moduleGroup[moduleId].loading) {
          this.loadDiffByModule(moduleId)
        }
      },
      goBackModule() {
        this.$routerActions.back()
      },
      handleGoback() {
        this.goBackModule()
      }
    }
  }
</script>

<style lang="scss" scoped>
.service-template-sync-layout {
  .layout-header {
    margin: 20px 24px 0 24px;
  }
  .layout-top {
    display: flex;
    margin: 24px;
    .title {
      font-size: 14px;
    }

    .type-legend {
      display: flex;
      align-items: center;
      font-size: 12px;

      .legend-item {
        margin-right: 30px;
        .dot {
          display: inline-block;
          width: 8px;
          height: 8px;
          border-radius: 50%;
          background-color: #2DCB56;
          margin-right: 2px;
          &.added {
            background-color: #2DCB56;
          }
          &.changed {
            background-color: #FF9C01;
          }
          &.removed {
            background-color: #FF5656;
          }
        }
      }
    }
  }

  .layout-main {
    margin: 24px;
  }

  .module-instance-container {
    background: #fff;
    box-shadow: 0 2px 4px 0 rgba(25, 25, 41, 0.05);
    border-radius: 2px;
    padding: 24px;

    & + .module-instance-container {
      margin-top: 16px;
    }
  }

  .module-instance-group {
    .module-instance-item {
      margin: 24px 16px 0 16px;
    }
  }

  .layout-footer {
    display: flex;
    align-items: center;
    height: 52px;
    padding: 0 24px;
    margin-top: 8px;
    .bk-button {
      min-width: 86px;

      & + .bk-button {
        margin-left: 8px;
      }
    }
    .auth-box + .bk-button {
      margin-left: 8px;
    }
    &.is-sticky {
      background-color: #fff;
      border-top: 1px solid $borderColor;
    }
  }
}
</style>
