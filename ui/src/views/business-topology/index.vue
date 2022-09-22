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
  <div class="layout" v-bkloading="{ isLoading: $loading(Object.values(request)) }" style="overflow: hidden;">
    <cmdb-resize-layout
      store-id="businessTopoPanel"
      :class="['resize-layout fl', { 'is-collapse': layout.topologyCollapse }]"
      direction="right"
      :handler-offset="3"
      :min="200"
      :max="480"
      :disabled="layout.topologyCollapse">
      <topology-tree ref="topologyTree" :active="activeTab" v-test-id></topology-tree>
      <i class="topology-collapse-icon bk-icon icon-angle-left"
        @click="layout.topologyCollapse = !layout.topologyCollapse">
      </i>
    </cmdb-resize-layout>
    <div class="tab-layout">
      <bk-tab class="topology-tab" type="unborder-card" v-test-id
        :active.sync="activeTab"
        :validate-active="false"
        :before-toggle="handleTabToggle">
        <bk-tab-panel name="hostList" :label="$t('主机列表')">
          <bk-exception class="empty-set" type="empty" scene="part" v-if="emptySet">
            <i18n path="该集群尚未创建模块">
              <template #link>
                <cmdb-auth :auth="{ type: $OPERATION.C_TOPO, relation: [bizId] }">
                  <bk-button text slot-scope="{ disabled }"
                    theme="primary"
                    :disabled="disabled"
                    @click="handleCreateModule">
                    {{$t('立即创建')}}
                  </bk-button>
                </cmdb-auth>
              </template>
            </i18n>
          </bk-exception>
          <host-list v-show="!emptySet" :active="activeTab === 'hostList'" ref="hostList" v-test-id></host-list>
        </bk-tab-panel>

        <bk-tab-panel
          :name="isContainerNode ? 'podList' : 'serviceInstance'"
          :label="$t(isContainerNode ? 'Pod列表' : '服务实例')">
          <pod-list v-if="activeTab === 'podList'" v-test-id></pod-list>
          <template v-else>
            <div class="non-business-module" v-if="!showServiceInstance">
              <div class="tips">
                <i class="bk-cc-icon icon-cc-tips"></i>
                <span>{{$t('非业务模块，无服务实例，请选择业务模块查看')}}</span>
              </div>
            </div>
            <service-instance-view v-else-if="activeTab === 'serviceInstance'" v-test-id></service-instance-view>
          </template>
        </bk-tab-panel>

        <bk-tab-panel name="nodeInfo" :label="$t('节点信息')" render-directive="if">
          <simple-node-info v-if="isSimpleNodeInfo" />
          <container-node-info v-else-if="isContainerNode" />
          <service-node-info v-else :active="activeTab === 'nodeInfo'" />
        </bk-tab-panel>
      </bk-tab>
    </div>
    <router-subview></router-subview>
  </div>
</template>

<script>
  import TopologyTree from './children/topology-tree.vue'
  import HostList from './host/host-list.vue'
  import SimpleNodeInfo from './children/simple-node-info.vue'
  import ServiceNodeInfo from './children/service-node-info.vue'
  import ContainerNodeInfo from './children/container-node-info.vue'
  import { mapGetters } from 'vuex'
  import Bus from '@/utils/bus.js'
  import RouterQuery from '@/router/query'
  import ServiceInstanceView from './service-instance/view'
  import PodList from './pod/pod-list.vue'
  export default {
    components: {
      TopologyTree,
      HostList,
      SimpleNodeInfo,
      ServiceNodeInfo,
      ContainerNodeInfo,
      ServiceInstanceView,
      PodList
    },
    data() {
      return {
        activeTab: RouterQuery.get('tab', 'hostList'),
        layout: {
          topologyCollapse: false
        },
        request: {
          mainline: Symbol('mainline'),
          properties: Symbol('properties')
        }
      }
    },
    computed: {
      ...mapGetters(['supplierAccount']),
      ...mapGetters('objectBiz', ['bizId']),
      ...mapGetters('businessHost', ['selectedNode']),
      showServiceInstance() {
        return this.selectedNode && this.selectedNode.data.bk_obj_id === 'module' && this.selectedNode.data.default === 0
      },
      isSimpleNodeInfo() {
        return this.selectedNode && (this.selectedNode.data.default !== 0 || this.selectedNode.data.is_folder)
      },
      emptySet() {
        return this.selectedNode && this.selectedNode.data.bk_obj_id === 'set'
          && this.selectedNode.children && !this.selectedNode.children.length
      },
      isContainerNode() {
        return !!this.selectedNode?.data?.is_container
      }
    },
    watch: {
      activeTab(tab) {
        this.$nextTick(() => {
          // 仅保留公用的参数重置路由
          RouterQuery.setAll({
            tab,
            node: RouterQuery.get('node'),
            topo_path: this.isContainerNode ? RouterQuery.get('topo_path') : undefined,
            _f: RouterQuery.get('_f'),
            _t: Date.now()
          })
        })
      },
      emptySet(value) {
        if (!value) {
          this.$nextTick(() => {
            this.$refs.hostList.doLayoutTable()
          })
        }
      }
    },
    async created() {
      this.unwatch = RouterQuery.watch('tab', (value = 'hostList') => {
        this.activeTab = value
      })
      try {
        const topologyModels = await this.getTopologyModels()
        const properties = await this.getProperties(topologyModels)
        this.$store.commit('businessHost/setTopologyModels', topologyModels)
        this.$store.commit('businessHost/setPropertyMap', properties)
        this.$store.commit('businessHost/resolveCommonRequest')
      } catch (e) {
        console.error(e)
      }
    },
    beforeDestroy() {
      this.$store.commit('businessHost/clear')
      this.unwatch()
    },
    methods: {
      handleTabToggle() {
        Bus.$emit('toggle-host-filter', false)
        return true
      },
      handleCreateModule() {
        this.$refs.topologyTree.showCreateDialog(this.selectedNode)
      },
      getTopologyModels() {
        return this.$store.dispatch('objectMainLineModule/searchMainlineObject', {
          config: {
            requestId: this.request.mainline
          }
        })
      },
      getProperties(models) {
        return this.$store.dispatch('objectModelProperty/batchSearchObjectAttribute', {
          injectId: 'host',
          params: {
            bk_biz_id: this.bizId,
            bk_obj_id: {
              $in: models.map(model => model.bk_obj_id)
            },
            bk_supplier_account: this.supplierAccount
          },
          config: {
            requestId: this.request.properties
          }
        })
      }
    }
  }
</script>

<style lang="scss" scoped>
    .resize-layout {
        position: relative;
        width: 286px;
        height: 100%;
        padding-top: 10px;
        border-right: 1px solid $cmdbLayoutBorderColor;
        &.is-collapse {
            width: 0 !important;
            border-right: none;
            .topology-collapse-icon:before {
                display: inline-block;
                transform: rotate(180deg);
            }
        }
        .topology-collapse-icon {
            position: absolute;
            left: 100%;
            top: 50%;
            width: 16px;
            height: 100px;
            line-height: 100px;
            background: $cmdbLayoutBorderColor;
            border-radius: 0px 12px 12px 0px;
            transform: translateY(-50%);
            text-align: center;
            font-size: 20px;
            color: #fff;
            cursor: pointer;
            text-indent: -2px;
            &:hover {
                background: #699DF4;
            }
        }
    }
    .tab-layout {
        height: 100%;
        overflow: hidden;
        .topology-tab {
            /deep/ {
                .bk-tab-header {
                    padding: 0;
                    margin: 0 20px;
                }
                .bk-tab-section {
                  height: calc(100% - 50px);
                }
            }
        }
    }

    .non-business-module {
        display: flex;
        height: 80%;
        justify-content: center;
        align-items: center;
        .tips {
            font-size: 14px;
            .bk-cc-icon {
                font-size: 16px;
                margin-top: -2px;
            }
        }
    }

    .empty-set {
        height: 80%;
        justify-content: center;
    }
</style>
