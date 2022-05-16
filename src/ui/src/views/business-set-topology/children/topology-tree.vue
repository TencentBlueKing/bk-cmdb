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
  <section
    class="tree-layout"
    v-bkloading="{
      isLoading: initializing
    }"
  >
    <bk-big-tree
      ref="tree"
      class="topology-tree"
      v-test-id
      selectable
      display-matched-node-descendants
      :height="$APP.height - 160"
      :node-height="36"
      :options="{
        idKey: spliceNodeId,
        nameKey: 'bk_inst_name',
        childrenKey: 'child'
      }"
      :before-select="handleBeforeSelect"
      :lazy-method="loadNodeChildren"
      @select-change="handleSelectChange"
    >
      <div
        :class="['node-info clearfix', { 'is-selected': node.selected }]"
        slot-scope="{ node, data }"
      >
        <!-- 业务集拓扑图标 -->
        <i class="node-icon fl" v-if="data.bk_obj_id === BUILTIN_MODELS.BUSINESS_SET">
          <i class="icon-cc-business-set"></i>
        </i>

        <!-- 内置拓扑图标 -->
        <i
          class="internal-node-icon fl"
          v-else-if="data.default !== 0"
          :class="getBuildInIconClass(node, data)"
        >
        </i>

        <!-- 自定义拓扑图标 -->
        <i
          v-else
          :class="[
            'node-icon fl',
            { 'is-selected': node.selected, 'is-template': isTemplate(node) }
          ]"
        >
          {{ getModelPrefix(data.bk_obj_id) }}
        </i>
        <cmdb-loading
          :class="['node-count fr', { 'is-selected': node.selected }]"
          :loading="['pending', undefined].includes(data.status)"
        >
          {{ getNodeCount(data) }}
        </cmdb-loading>
        <span class="node-name" :title="node.name">{{ node.name }}</span>
      </div>
    </bk-big-tree>
  </section>
</template>

<script>
  import { mapGetters, mapState } from 'vuex'
  import Bus from '@/utils/bus'
  import RouterQuery from '@/router/query'
  import { addResizeListener, removeResizeListener } from '@/utils/resize-events'
  import FilterStore from '@/components/filters/store'
  import { sortTopoTree } from '@/utils/tools'
  import CmdbLoading from '@/components/loading/loading'
  import { TopologyService } from '@/service/business-set/topology.js'
  import CombineRequest from '@/api/combine-request.js'
  import to from 'await-to-js'
  import { BUILTIN_MODELS } from '@/dictionary/model-constants.js'

  export default {
    name: 'topology-tree',
    components: {
      CmdbLoading,
    },
    props: {
      /**
       * 统计类型
       */
      countType: {
        type: String,
        required: true,
        validator(value) {
          return ['host_count', 'service_instance_count'].includes(value)
        },
      },
      /**
       * 拓扑节点模型，用来补充拓扑树中的节点图标
       */
      topologyModels: {
        type: Array,
        default: () => [],
        required: true,
      },
    },
    data() {
      this.BUILTIN_MODELS = BUILTIN_MODELS

      return {
        nodeIconMap: {
          1: 'icon-cc-host-free-pool',
          2: 'icon-cc-host-breakdown',
          default: 'icon-cc-host-free-pool',
        },
        initializing: false
      }
    },
    computed: {
      ...mapState('bizSet', ['bizSetName', 'bizSetId', 'bizId']),
      ...mapGetters('businessHost', ['selectedNode'])
    },
    watch: {
      countType() {
        this.countBizSetSum()
      }
    },
    created() {
      Bus.$on('refresh-count', this.refreshCount)
      Bus.$on('refresh-count-by-node', this.refreshCountByNode)
      this.initTopology()
    },
    mounted() {
      addResizeListener(this.$el, this.handleResize)
    },
    beforeDestroy() {
      Bus.$off('refresh-count', this.refreshCount)
      Bus.$off('refresh-count-by-node', this.refreshCountByNode)
      removeResizeListener(this.$el, this.handleResize)
    },
    methods: {
      /**
       * 是否为叶子节点
       * @param {string} modelId 模型 ID
       */
      isLeaf(modelId) {
        return modelId === BUILTIN_MODELS.MODULE
      },
      /**
       * 初始化拓扑，从 Route query 中获取拓扑路径，并加载上次定位的节点
       */
      async initTopology() {
        const topoPathQueryStr = this.$route.query.topoPath
        const topoPath = topoPathQueryStr?.split(',') || null
        let topoTreeData = null
        const expandedIds = []

        this.initializing = true

        const currentBizSetTopo = await this.loadTopologyData(this.bizSetId, BUILTIN_MODELS.BUSINESS_SET, this.bizSetId)

        let currentBizNode = null

        // 检测业务是否存在，如果查询条件中的业务不存在则清除查询条件刷新当前页面
        if (this.bizId) {
          currentBizNode = currentBizSetTopo.find(item => item.bk_inst_id === this.bizId)
          if (!currentBizNode) {
            this.clearQueryReload()
            return
          }
        }

        if (topoPath) {
          let parentData = null

          // 根据 topoPath 遍历出拓扑树数据
          for (let pathIndex = 0; pathIndex < topoPath.length; pathIndex++) {
            const path = topoPath[pathIndex].split('-')
            const modelId = path[0]
            const instanceId = Number(path[1])

            if (this.isLeaf(modelId)) break

            const [err, data] = await to(this.loadTopologyData(this.bizSetId, modelId, instanceId))

            // 获取数据失败则清除条件刷新页面，非法id或者数据被删除
            if (err) {
              this.clearQueryReload()
              return
            }

            if (!topoTreeData) {
              topoTreeData = data
            } else {
              const foundNode = parentData.find(parentItem => parentItem.bk_inst_id === instanceId)
              if (foundNode && parentData && data) {
                foundNode.child = data
              }
            }

            parentData = data

            if (data) {
              expandedIds.push(...data.map(item => this.spliceNodeId(item)))
            }
          }
        }

        // 初始化根节点
        const root = [
          {
            bk_obj_id: BUILTIN_MODELS.BUSINESS_SET,
            bk_obj_name: '',
            bk_inst_id: this.bizSetId,
            bk_inst_name: this.bizSetName,
            host_count: 0,
            service_instance_count: 0,
            child: currentBizSetTopo,
          },
        ]

        if (topoTreeData && currentBizNode) {
          currentBizNode.child = topoTreeData
        }

        this.$refs.tree.setData(root)

        const bizNodes = this.$refs.tree.nodes[0].children
        const [firstBizNode] = bizNodes

        const expanedNodes = currentBizNode ? expandedIds.map(id => this.$refs.tree.getNodeById(id)) : []

        this.loadNodeCount([...bizNodes, ...expanedNodes])
        this.setDefaultExpandedNode(this.getNodeFromQuery() || firstBizNode)

        this.initializing = false
      },
      /**
       * 设置默认展开节点
       * @param {Object} node 默认展开节点
       */
      setDefaultExpandedNode(node) {
        if (node) {
          const { tree } = this.$refs
          const siblingNodes = node.parent.children
          const childrenNodes = node.children

          const childrenHasLeaf = childrenNodes.some(child => this.isLeaf(child.data.bk_obj_id))
          const siblingHasLeaf = siblingNodes.some(child => this.isLeaf(child.data.bk_obj_id))

          tree.setExpanded(node.id)
          tree.setSelected(node.id, { emitEvent: true })

          if (childrenHasLeaf) {
            childrenNodes.forEach((childNode) => {
              if (this.isLeaf(childNode.data.bk_obj_id) && !childNode.state.expanded) {
                tree.setExpanded(childNode.id)
              }
            })
          }

          if (siblingHasLeaf) {
            siblingNodes.forEach((siblingNode) => {
              if (this.isLeaf(siblingNode.data.bk_obj_id) && !siblingNode.state.expanded) {
                tree.setExpanded(siblingNode.id)
              }
            })
          }

          !this.initialized
            && this.$nextTick(() => {
              this.initialized = true
              const index = tree.visibleNodes.indexOf(node)
              tree.$refs.virtualScroll.scrollPageByIndex(index)
            })
        }
      },
      /**
       * 懒加载子节点
       */
      async loadNodeChildren(node) {
        const { bk_obj_id: modelId, bk_inst_id: instanceId } = node.data

        if (this.isLeaf(modelId)) return {}

        const [, data] = await to(this.loadTopologyData(this.bizSetId, modelId, instanceId))

        const leafIds = data
          ?.filter(({ bk_obj_id: objId }) => this.isLeaf(objId))
          .map(item => this.spliceNodeId(item))

        // big-tree 缺少 after-lazy-load 的 hook，所以需要手动把任务置后
        setTimeout(() => {
          this.loadNodeCount([node, ...node.children])
        }, 0)

        return {
          data,
          leaf: leafIds,
        }
      },
      /**
       * 加载拓扑数据，并进行排序、附加模型名称
       */
      loadTopologyData(bizSetId, parentModelId, parentInstanceId) {
        return TopologyService.findChildren({
          bizSetId,
          parentModelId,
          parentInstanceId,
        }).then((data) => {
          sortTopoTree(data, 'bk_inst_name')

          if (parentModelId === BUILTIN_MODELS.BUSINESS) {
            data.sort(a => (a.bk_obj_id === BUILTIN_MODELS.SET && a.default === 1 ? -1 : 0))
          }

          return data
        })
      },
      getModelPrefix(modelId) {
        const modelName = this.topologyModels?.find(model => model.bk_obj_id === modelId)?.bk_obj_name

        return modelName?.[0] || ''
      },
      /**
       * 使业务集根节点不可点击
       */
      handleBeforeSelect({ data }) {
        if (data.bk_obj_id === BUILTIN_MODELS.BUSINESS_SET) {
          return false
        }
        return true
      },
      /**
       * 从 Query 中获取节点
       */
      getNodeFromQuery() {
        // 选中 query 中的节点
        const nodeId = RouterQuery.get('node', '')

        if (nodeId) {
          const node = this.$refs.tree.getNodeById(nodeId)

          if (node) {
            return node
          }
        }

        return null
      },
      /**
       * 拼接节点 ID
       * @param {object} data 节点数据
       */
      spliceNodeId(data) {
        return `${data.bk_obj_id}-${data.bk_inst_id}`
      },
      /**
       * 获取内置的节点 Icon 类名
       */
      getBuildInIconClass(node, data) {
        const nodeClass = []

        nodeClass.push(this.nodeIconMap[data.default] || this.nodeIconMap.default)

        if (node.selected) {
          nodeClass.push('is-selected')
        }

        return nodeClass
      },
      getBizSetRootNode() {
        return this.$refs?.tree.getNodeById(`${BUILTIN_MODELS.BUSINESS_SET}-${this.bizSetId}`)
      },
      /**
       * 处理节点选择事件
       */
      handleSelectChange(node) {
        const oldId = this.$route.query.node
        const newId = node.id
        const expandedBizId = this.findBizIdByNode(node)

        this.$store.commit('bizSet/setBizId', expandedBizId)
        this.$store.commit('businessHost/setSelectedNode', node)
        Bus.$emit('toggle-host-filter', false)

        const query = {
          node: newId,
          bizId: expandedBizId,
          topoPath: this.genTopoPathByNode(node).join(','),
          page: 1,
          _t: Date.now(),
        }

        RouterQuery.set(query)

        // 避免切换目录时收起目录
        this.$refs.tree.setExpanded(newId)

        // 切换节点时，如果存在筛选条件需要清除
        if (FilterStore.hasCondition && oldId !== newId) {
          FilterStore.setActiveCollection(null)
        }
      },
      /**
       * 生成节点的完整拓扑路径
       *  @param {Object} node BigTree 实例节点
       */
      genTopoPathByNode(node) {
        const path = []
        let currentNode = node

        while (currentNode.parent) {
          path.push(this.spliceNodeId(currentNode.data))
          currentNode = currentNode.parent
        }

        return path.reverse()
      },
      /**
       * 获取节点所属的业务 ID
       * @param {Object} node BigTree 实例节点
       */
      findBizIdByNode(node) {
        let parentBizId = null
        let currentNode = node

        while (currentNode.parent) {
          if (currentNode.data.bk_obj_id === BUILTIN_MODELS.BUSINESS) {
            parentBizId = currentNode.data.bk_inst_id
          }
          currentNode = currentNode.parent
        }

        return parentBizId
      },
      /**
       * 加载节点数量统计
       * @param {array} targetNodes 需要统计的节点
       * @param {boolean} force 是否强制加载统计状态
       */
      async loadNodeCount(targetNodes, force = false) {
        let nodes = null

        const doneStatuses = ['pending', 'finished']

        if (force) {
          nodes = targetNodes.filter(({ data }) => !doneStatuses.includes(data.status)
            && data.bk_obj_id !== BUILTIN_MODELS.BUSINESS_SET)
        } else {
          nodes = targetNodes.filter(({ data }) => data.bk_obj_id !== BUILTIN_MODELS.BUSINESS_SET)
        }

        if (!nodes.length) {
          this.setNodesStatus([this.getBizSetRootNode()], 'finished')
          return
        }

        this.setNodesStatus(nodes, 'pending')

        const allCondition = nodes.map(({ data }) => ({
          bk_obj_id: data.bk_obj_id,
          bk_inst_id: data.bk_inst_id,
        }))
        const requests = await CombineRequest.setup(
          Symbol(),
          condition => TopologyService.getInstanceCount(this.bizSetId, condition, { globalError: false }),
          { segment: 10, concurrency: 1 }
        ).add(allCondition)

        for (const req of requests) {
          await req
            .then((statData) => {
              this.setNodeCount(nodes, statData[0].value)
            })
            .catch(() => {
              this.setNodesStatus(nodes, 'error')
            })
        }

        this.countBizSetSum()
      },
      /**
       * 计算业务集下实例总数
       */
      countBizSetSum() {
        if (!this.$refs?.tree) {
          return
        }

        const bizSetSum = this.$refs.tree?.nodes
          .filter(node => node.data.bk_obj_id === BUILTIN_MODELS.BUSINESS)
          .reduce((acc, node) => {
            const count = this.getNodeCount(node.data)
            return acc + count
          }, 0)

        const rootNode = this.getBizSetRootNode()

        rootNode.data[this.countType] = bizSetSum

        this.setNodesStatus([rootNode], 'finished')
      },
      /**
       * 设置所有节点的加载状态
       * @param {array} nodes 节点数组
       * @param {string} status 节点状态
       */
      setNodesStatus(nodes, status) {
        nodes.forEach((node) => {
          if (node?.data) {
            this.$set(node.data, 'status', status)
          }
        })
      },
      /**
       * 设置节点的统计数据
       * @param {array} nodes 节点数组
       * @param {string} data 最新的节点数据
       */
      setNodeCount(nodes, data) {
        nodes.forEach((node) => {
          const foundStat = data?.find(item => item.bk_obj_id === node.data.bk_obj_id
            && item.bk_inst_id === node.data.bk_inst_id)
          if (foundStat) {
            this.$set(node.data, 'host_count', foundStat.host_count)
            this.$set(
              node.data,
              'service_instance_count',
              foundStat.service_instance_count
            )
            this.$set(node.data, 'status', 'finished')
          }
        })
      },
      /**
       * 获取节点数量统计
       */
      getNodeCount(data) {
        return data[this.countType] ?? 0
      },
      /**
       * 判断节点是否为模板节点
       */
      isTemplate(node) {
        return node.data.service_template_id || node.data.set_template_id
      },
      /**
       * 刷新指定主机关联的节点的数量统计
       */
      async refreshCount({ hosts, target }) {
        const nodes = []

        if (target) {
          const node = this.$refs.tree.getNodeById(this.spliceNodeId(target))
          node && nodes.push(node, ...node.parents)
        }

        hosts.forEach(({ module: modules }) => {
          modules.forEach((module) => {
            const node = this.$refs.tree.getNodeById(`module-${module.bk_module_id}`)
            node && nodes.push(node, ...node.parents)
          })
        })

        const nodeSet = new Set()
        const uniqueNodes = nodes.filter((node) => {
          if (nodeSet.has(node)) return false
          nodeSet.add(node)
          return true
        })

        this.loadNodeCount(uniqueNodes, true)
      },
      /**
       * 刷新指定节点的数量统计
       */
      refreshCountByNode(node) {
        const currentNode = node || this.selectedNode
        const nodes = []
        const treeNode = this.$refs.tree.getNodeById(currentNode.id)

        if (treeNode) {
          nodes.push(treeNode, ...treeNode.parents)
        }

        this.loadNodeCount(nodes, true)
      },
      handleResize() {
        this.$refs.tree.resize()
      },
      clearQueryReload() {
        this.$routerActions.redirect({
          ...this.$route,
          query: {},
          reload: true
        })
      }
    },
  }
</script>

<style lang="scss" scoped>
.tree-layout {
  overflow: hidden;
}
.tree-search {
  display: block;
  width: auto;
  margin: 0 20px;
}
.topology-tree {
  padding: 10px 0;
  margin-right: 2px;
  @include scrollbar-y(6px);
  .node-icon {
    display: block;
    width: 20px;
    height: 20px;
    margin: 8px 4px 8px 0;
    border-radius: 50%;
    background-color: #c4c6cc;
    line-height: 1.666667;
    text-align: center;
    font-size: 12px;
    font-style: normal;
    color: #fff;
    &.is-template {
      background-color: #97aed6;
    }
    &.is-selected {
      background-color: #3a84ff;
    }
    .icon-cc-business-set{
      vertical-align: 0;
      font-size: 16px;
    }
  }
  .node-name {
    display: block;
    height: 36px;
    line-height: 36px;
    overflow: hidden;
    @include ellipsis;
  }
  .node-count {
    padding: 0 5px;
    margin: 9px 20px 9px 4px;
    height: 18px;
    line-height: 17px;
    border-radius: 2px;
    background-color: #f0f1f5;
    color: #979ba5;
    font-size: 12px;
    text-align: center;
    &.is-selected {
      background-color: #a2c5fd;
      color: #fff;
    }
    &.loading {
      background-color: transparent;
    }
  }
  .internal-node-icon {
    width: 20px;
    height: 20px;
    line-height: 20px;
    text-align: center;
    margin: 8px 4px 8px 0;
    &.is-selected {
      color: #ffb400;
    }
  }
}
</style>
