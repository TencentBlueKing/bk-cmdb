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
  <div class="topology-tree-wrapper">
    <bk-big-tree class="topology-tree"
      ref="tree"
      v-bind="treeOptions"
      v-bkloading="{
        isLoading: $loading([requestIds.getTopology, requestIds.mainLine, requestIds.moduleApplyStatus])
      }"
      :options="{
        idKey: idGenerator,
        nameKey: 'bk_inst_name',
        childrenKey: 'child'
      }"
      :height="$APP.height - 160 - 44"
      :node-height="36"
      :check-on-click="true"
      :before-select="beforeSelect"
      :before-check="beforeCheck"
      :filter-method="filterMethod"
      :enable-title-tip="true"
      @select-change="handleSelectChange"
      @check-change="handleCheckChange"
      @expand-change="handleExpandChange">
      <div
        class="node-info clearfix"
        :title="getNodeTips(data)"
        :class="{ 'is-selected': node.selected }"
        slot-scope="{ node, data }">
        <i class="internal-node-icon fl"
          v-if="data.default !== 0"
          :class="getInternalNodeClass(node, data)">
        </i>
        <i class="node-icon fl" v-else
          :class="{
            'is-selected': node.selected,
            'is-template': isTemplate(node),
            'is-leaf-icon': node.isLeaf
          }">
          {{modelIconMap[data.bk_obj_id]}}
        </i>
        <span v-show="applyEnabled(node)" class="config-icon fr"><i class="bk-cc-icon icon-cc-selected"></i></span>
        <div class="info-content">
          <span class="node-name" :title="data.bk_inst_name">{{data.bk_inst_name}}</span>
        </div>
      </div>
      <cmdb-table-empty
        slot="empty"
        :stuff="table.stuff"
        @clear="handleClearFilter">
      </cmdb-table-empty>
    </bk-big-tree>
  </div>
</template>

<script>
  import { mapGetters, mapState } from 'vuex'
  import Bus from '@/utils/bus'
  import { addResizeListener, removeResizeListener } from '@/utils/resize-events'
  import topologyInstanceService, { requestIds as topologyrequestIds } from '@/service/topology/instance.js'
  import { CONFIG_MODE } from '@/service/service-template/index.js'

  export default {
    props: {
      treeOptions: {
        type: Object,
        default: () => ({})
      },
      action: {
        type: String,
        default: 'batch-edit'
      }
    },
    data() {
      return {
        treeData: [],
        treeStat: {},
        mainLine: [],
        createInfo: {
          show: false,
          visible: false,
          properties: [],
          parentNode: null,
          nextModelId: null
        },
        requestIds: {
          searchNode: Symbol('searchNode'),
          moduleApplyStatus: Symbol('moduleApplyStatus'),
          mainLine: Symbol('mainLine'),
          ...topologyrequestIds
        },
        nodeIconMap: {
          1: 'icon-cc-host-free-pool',
          2: 'icon-cc-host-breakdown',
          default: 'icon-cc-host-free-pool'
        },
        table: {
          stuff: {
            type: 'search',
            payload: {
              emptyText: this.$t('bk.table.emptyText')
            }
          }
        }
      }
    },
    computed: {
      ...mapGetters(['supplierAccount']),
      ...mapState('hostApply', ['ruleDraft']),
      business() {
        return this.$store.state.objectBiz.bizId
      },
      targetId() {
        return Number(this.$route.query.id)
      },
      propertyMap() {
        return this.$store.state.businessTopology.propertyMap
      },
      mainLineModels() {
        const models = this.$store.getters['objectModelClassify/models']
        return this.mainLine.map(data => models.find(model => model.bk_obj_id === data.bk_obj_id))
      },
      modelIconMap() {
        const map = {}
        this.mainLineModels.forEach((model) => {
          // eslint-disable-next-line prefer-destructuring
          map[model.bk_obj_id] = model.bk_obj_name[0]
        })
        return map
      },
      isDel() {
        return this.action === 'batch-del'
      },
      isEdit() {
        return this.ac
      }
    },
    watch: {
      action() {
        this.setNodeDisabled()
      }
    },
    async created() {
      Bus.$on('host-apply-topology-search', this.handleSearch)
      const [data, mainLine] = await Promise.all([
        topologyInstanceService.geFulltWithStat(this.business),
        this.getMainLine()
      ])

      // 从另外一侧关闭应用规则需要更新对应一侧的状态数据
      Bus.$on('host-apply-closed', (mode, id) => {
        if (mode === CONFIG_MODE.TEMPLATE) {
          this.setModuleApplyStatusByTemplate(this.treeStat?.withTemplateModuleIdMap.get(id))
        }
      })

      this.treeData = data
      this.mainLine = mainLine
      this.treeStat = this.getTreeStat()
      this.$nextTick(() => {
        this.setDefaultState(data)
      })

      this.setModuleApplyStatusByTemplate()
    },
    activated() {
      const treeNode = this.$refs.tree.getNodeById(`module_${this.targetId}`)
      if (treeNode) {
        this.$emit('selected', treeNode)
      }
    },
    mounted() {
      addResizeListener(this.$el, this.handleResize)
    },
    beforeDestroy() {
      Bus.$off('host-apply-topology-search', this.handleSearch)
      removeResizeListener(this.$el, this.handleResize)
    },
    methods: {
      async handleSearch(params) {
        try {
          if (params.query_filter.rules.length) {
            const keywordRuleIndex = params.query_filter.rules.findIndex(item => item.field === 'keyword')

            if (keywordRuleIndex === -1) {
              // 不存在关键字则只采用接口搜索
              const data = await this.searchNode(params)
              this.$refs.tree.filter({ remote: data })
            } else {
              // 先取出关键字的参数值
              const keyword = params.query_filter.rules[keywordRuleIndex]?.value

              // 接口搜索需要去掉keyword参数
              params.query_filter.rules.splice(keywordRuleIndex, 1)

              // 两者都存在，混合搜索
              if (params.query_filter.rules.length) {
                const data = await this.searchNode(params)
                this.$refs.tree.filter({
                  keyword,
                  remote: data,
                })
              } else {
                // 仅存在关键字搜索
                this.$refs.tree.filter({ keyword })
              }
            }
          } else {
            // 清空搜索
            this.$refs.tree.filter()
          }
          this.$refs.tree.removeChecked({ emitEvent: false })
        } catch (e) {
          console.error(e)
        }
      },
      searchNode(params) {
        return this.$store.dispatch('hostApply/searchNode', {
          bizId: this.business,
          params,
          config: {
            requestId: this.requestIds.searchNode
          }
        })
      },
      filterMethod({ remote: remoteData, keyword }, node) {
        // eslint-disable-next-line newline-per-chained-call
        const keywordFilter = (keyword, node) => String(node.name).toLowerCase().indexOf(keyword.toLowerCase()) > -1
        // eslint-disable-next-line max-len
        const remoteFilter = (remoteData, node) => remoteData.some(item => item.bk_inst_id === node.data.bk_inst_id && item.bk_obj_id === node.data.bk_obj_id)

        if (remoteData && keyword) {
          return keywordFilter(keyword, node) && remoteFilter(remoteData, node)
        }

        if (remoteData) {
          return remoteFilter(remoteData, node)
        }

        if (keyword) {
          return keywordFilter(keyword, node)
        }
      },
      async setDefaultState(data) {
        if (!data?.length) {
          return
        }

        this.$refs.tree.setData(data)

        let defaultNodeId
        const { firstModule } = this.treeStat
        if (!isNaN(this.targetId)) {
          defaultNodeId = `module_${this.targetId}`
        } else if (firstModule) {
          defaultNodeId = this.idGenerator(firstModule)
        }
        if (defaultNodeId) {
          const treeNode = this.$refs.tree.getNodeById(defaultNodeId)
          if (treeNode.data.service_template_id) {
            await this.setModuleApplyStatusByTemplate([treeNode.data.bk_inst_id])
          }
          this.$refs.tree.setSelected(defaultNodeId, { emitEvent: true })
          // 展开父级
          this.$refs.tree.setExpanded(treeNode?.parent.id, { emitEvent: true })
        }
      },
      getTreeStat() {
        const stat = {
          firstModule: null,
          levels: {},
          noRuleIds: [],
          // 模板id为键的模块id列表
          withTemplateModuleIdMap: new Map(),
          withTemplateHostApplyIds: []
        }
        const findModule = function (data, parent) {
          // eslint-disable-next-line no-restricted-syntax
          for (const item of data) {
            stat.levels[item.bk_inst_id] = parent ? (stat.levels[parent.bk_inst_id] + 1) : 0
            if (item.bk_obj_id === 'module') {
              if (item.service_template_id) {
                if (stat.withTemplateModuleIdMap.has(item.service_template_id)) {
                  stat.withTemplateModuleIdMap.get(item.service_template_id).push(item.bk_inst_id)
                } else {
                  stat.withTemplateModuleIdMap.set(item.service_template_id, [item.bk_inst_id])
                }
                if (item.service_template_host_apply_enabled) {
                  stat.withTemplateHostApplyIds.push(item.bk_inst_id)
                }
              }
              if (item.host_apply_rule_count === 0) {
                stat.noRuleIds.push(item.bk_inst_id)
              }
              if (!stat.firstModule) {
                stat.firstModule = item
              }
            } else if (item.child.length) {
              const match = findModule(item.child, item)
              if (match && !stat.firstModule) {
                stat.firstModule = item
              }
            }
          }
        }
        findModule(this.treeData)
        return stat
      },
      setNodeDisabled() {
        const { withTemplateHostApplyIds, noRuleIds } = this.treeStat

        // 处于批量操作状态disabled存在模板配置的节点
        if (withTemplateHostApplyIds?.length) {
          const nodeIds = withTemplateHostApplyIds.map(id => `module_${id}`)
          this.$refs.tree.setDisabled(nodeIds, { emitEvent: true, disabled: Boolean(this.action) })
        }

        // 批量删除操作状态disabled不存在模块配置的节点
        if (noRuleIds?.length) {
          // 仅需处理未配置模板的节点
          const nodeIds = noRuleIds.filter(id => !withTemplateHostApplyIds.includes(id)).map(id => `module_${id}`)
          this.$refs.tree.setDisabled(nodeIds, { emitEvent: true, disabled: this.isDel })
        }
      },
      updateNodeStatus(id, { isClose, isClear }) {
        const nodeData = this.$refs.tree.getNodeById(`module_${id}`).data
        if (isClose) {
          nodeData.host_apply_enabled = false
        }
        if (isClear) {
          nodeData.host_apply_rule_count = 0
        }
        this.treeStat = this.getTreeStat()
      },
      getMainLine() {
        return this.$store.dispatch('objectMainLineModule/searchMainlineObject', {
          config: {
            requestId: this.requestIds.mainLine
          }
        })
      },
      idGenerator(data) {
        return `${data.bk_obj_id}_${data.bk_inst_id}`
      },
      applyEnabled(node) {
        return this.isModule(node) && (node.data.host_apply_enabled || node.data.service_template_host_apply_enabled)
      },
      isTemplate(node) {
        return node.data.service_template_id || node.data.set_template_id
      },
      isModule(node) {
        return node.data.bk_obj_id === 'module'
      },
      beforeSelect(node) {
        return this.isModule(node)
      },
      beforeCheck() {
        return Boolean(this.action)
      },
      getNodeTips(nodeData) {
        if (nodeData.service_template_host_apply_enabled) {
          return this.$t('需在模板中配置')
        }
        if (nodeData.host_apply_rule_count === 0 && this.isDel) {
          return this.$t('暂无策略')
        }
        return ''
      },
      getInternalNodeClass(node, data) {
        const clazz = []
        clazz.push(this.nodeIconMap[data.default] || this.nodeIconMap.default)
        if (node.selected) {
          clazz.push('is-selected')
        }
        return clazz
      },
      handleSelectChange(node) {
        this.$emit('selected', node)
      },
      handleCheckChange(id, checked) {
        this.$emit('checked', id, checked)
      },
      handleResize() {
        this.$refs.tree.resize()
      },
      async handleExpandChange(node) {
        if (!node.state.expanded) {
          return
        }

        if (!node.children[0].expanded) {
          const lastNodeLevel = node.children[0].level
          Bus.$emit('topologyTree/expandChange', lastNodeLevel)
        }
      },
      async setModuleApplyStatusByTemplate(ids) {
        const { withTemplateModuleIdMap } = this.treeStat
        let withTemplateModuleIds = []
        for (const moduleIds of withTemplateModuleIdMap.values()) {
          withTemplateModuleIds.push(...moduleIds)
        }
        if (ids) {
          withTemplateModuleIds = ids
        }

        if (!withTemplateModuleIds?.length) {
          return
        }

        try {
          const result = await this.$store.dispatch('hostApply/getModuleApplyStatusByTemplate', {
            params: {
              bk_biz_id: this.business,
              bk_module_ids: withTemplateModuleIds
            },
            config: {
              requestId: this.requestIds.moduleApplyStatus
            }
          })
          withTemplateModuleIds.forEach((id) => {
            const nodeData = this.$refs.tree.getNodeById(`module_${id}`).data
            const statusItem = result.find(item => item.bk_module_id === id)
            this.$set(nodeData, 'service_template_host_apply_enabled', statusItem.host_apply_enabled)
          })
        } catch (error) {
          console.log(error)
        }

        this.treeStat = this.getTreeStat()
      },
      handleClearFilter() {
        this.$emit('clearFilter', [])
      }
    }
  }
</script>

<style lang="scss" scoped>
  .topology-tree-wrapper {
    ::v-deep .bk-scroll-home .bk-bottom-scroll {
      display: none;
    }
  }
  .topology-tree {
      padding: 10px 0;
      margin-right: 2px;
      .node-info {
        min-width: 100px;
          .node-icon {
              width: 22px;
              height: 22px;
              line-height: 21px;
              text-align: center;
              font-style: normal;
              font-size: 12px;
              margin: 8px 8px 0 6px;
              border-radius: 50%;
              background-color: #c4c6cc;
              color: #fff;
              &.is-template {
                  background-color: #97aed6;
              }
              &.is-selected {
                  background-color: #3a84ff;
              }
          }
          .config-icon {
              position: relative;
              top: 6px;
              right: 4px;
              padding: 0 5px;
              height: 18px;
              line-height: 17px;
              color: #979ba5;
              font-size: 26px;
              text-align: center;
              color: #2dcb56;
          }
          .internal-node-icon{
              width: 20px;
              height: 20px;
              line-height: 20px;
              text-align: center;
              margin: 8px 4px 8px 0;
              &.is-selected {
                  color: #3a84ff;
              }
          }
          .info-content {
              display: flex;
              align-items: center;
              line-height: 36px;
              font-size: 14px;
              .node-name {
                  @include ellipsis;
                  margin-right: 8px;
              }
          }
      }

      .empty {
          position: absolute;
          display: flex;
          height: calc(100% - 30px);
          width: 100%;
          justify-content: center;
          align-items: center;
      }
  }
</style>
