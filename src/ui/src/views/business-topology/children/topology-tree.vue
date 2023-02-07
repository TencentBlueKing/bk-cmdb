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
  <section class="tree-layout" v-bkloading="{ isLoading: $loading(Object.values(request)) }">
    <bk-input class="tree-search" v-test-id
      clearable
      right-icon="bk-icon icon-search"
      :placeholder="$t('请输入关键词')"
      v-model.trim="filter">
    </bk-input>
    <bk-big-tree ref="tree" class="topology-tree" v-test-id
      selectable
      display-matched-node-descendants
      :height="$APP.height - 160"
      :node-height="36"
      :options="{
        idKey: getNodeId,
        nameKey: 'bk_inst_name',
        childrenKey: 'child'
      }"
      :lazy-method="lazyGetChildrenNode"
      :lazy-disabled="isLazyDisabledNode"
      @select-change="handleSelectChange"
      @expand-change="handleExpandChange">
      <template #default="{ node, data }">
        <topology-tree-node
          v-bind="{ node, data, isBlueKing, editable, nodeCountType }"
          @create="handleShowCreateDialog">
        </topology-tree-node>
      </template>
    </bk-big-tree>
    <bk-dialog class="bk-dialog-no-padding"
      v-model="createInfo.show"
      :show-footer="false"
      :mask-close="false"
      :width="580"
      @after-leave="handleAfterCancelCreateNode"
      @cancel="handleCancelCreateNode">
      <template v-if="createInfo.nextModelId === 'module'">
        <create-module v-if="createInfo.visible"
          :parent-node="createInfo.parentNode"
          @submit="handleCreateNode"
          @cancel="handleCancelCreateNode">
        </create-module>
      </template>
      <template v-else-if="createInfo.nextModelId === 'set'">
        <create-set v-if="createInfo.visible"
          :parent-node="createInfo.parentNode"
          @submit="handleCreateSetNode"
          @cancel="handleCancelCreateNode">
        </create-set>
      </template>
      <template v-else>
        <create-node v-if="createInfo.visible"
          :next-model-id="createInfo.nextModelId"
          :properties="createInfo.properties"
          :parent-node="createInfo.parentNode"
          @submit="handleCreateNode"
          @cancel="handleCancelCreateNode">
        </create-node>
      </template>
    </bk-dialog>
  </section>
</template>

<script>
  import { mapGetters } from 'vuex'
  import debounce from 'lodash.debounce'
  import CreateNode from './create-node.vue'
  import CreateSet from './create-set.vue'
  import CreateModule from './create-module.vue'
  import Bus from '@/utils/bus'
  import RouterQuery from '@/router/query'
  import { addResizeListener, removeResizeListener } from '@/utils/resize-events'
  import FilterStore from '@/components/filters/store'
  import TopologyTreeNode from './topology-tree-node.vue'
  import {
    MENU_BUSINESS_HOST_AND_SERVICE,
  } from '@/dictionary/menu-symbol'
  import topologyInstanceService from '@/service/topology/instance'
  import { BUILTIN_MODELS } from '@/dictionary/model-constants'
  import { isWorkload } from '@/service/container/common'
  import { CONTAINER_OBJECTS } from '@/dictionary/container'

  export default {
    components: {
      CreateNode,
      CreateSet,
      CreateModule,
      TopologyTreeNode
    },
    props: {
      active: {
        type: String,
        required: true
      }
    },
    data() {
      return {
        isBlueKing: false,
        filter: RouterQuery.get('keyword', ''),
        handleFilter: () => ({}),
        nodeCountType: 'host_count',
        request: {
          instance: Symbol('instance'),
          internal: Symbol('internal'),
          property: Symbol('property')
        },
        createInfo: {
          show: false,
          visible: false,
          properties: [],
          parentNode: null,
          nextModelId: null
        },
        editable: false,
        timer: null
      }
    },
    computed: {
      ...mapGetters('objectBiz', ['bizId']),
      ...mapGetters('businessHost', ['topologyModels', 'propertyMap']),
      ...mapGetters('businessHost', ['selectedNode'])
    },
    watch: {
      filter(value) {
        this.handleFilter()
        RouterQuery.set('keyword', value)
      },
      active: {
        immediate: true,
        handler(value) {
          const map = {
            hostList: 'host_count',
            serviceInstance: 'service_instance_count',
            podList: 'pod_count',
            nodeInfo: 'host_count'
          }
          if (Object.keys(map).includes(value)) {
            this.nodeCountType = map[value]
          }
        }
      },
      isBlueKing(flag) {
        if (flag) {
          this.getBlueKingEditStatus()
          clearInterval(this.timer)
          this.timer = setInterval(this.getBlueKingEditStatus, 1000 * 60)
        }
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
      this.destroyWatcher()
      Bus.$off('refresh-count', this.refreshCount)
      Bus.$off('refresh-count-by-node', this.refreshCountByNode)
      clearInterval(this.timer)
      removeResizeListener(this.$el, this.handleResize)
    },
    methods: {
      async initTopology() {
        try {
          const [topology, internal, container] = await Promise.all([
            this.getInstanceTopology(),
            this.getInternalTopology(),
            this.getContainerTopology()
          ])

          const { topo: containerTopo, leafIds: containerLeafIds } = container

          const root = topology[0] || {}

          const children = root.child || []

          const idlePool = {
            bk_obj_id: 'set',
            bk_inst_id: internal.bk_set_id,
            bk_inst_name: internal.bk_set_name,
            default: internal.default,
            is_idle_set: true,
            child: internal.module.map(module => ({
              bk_obj_id: 'module',
              bk_inst_id: module.bk_module_id,
              bk_inst_name: module.bk_module_name,
              default: module.default
            }))
          }
          children.unshift(idlePool)

          // 容器拓扑追加至底部
          children.push(...containerTopo)

          this.isBlueKing = root.bk_inst_name === '蓝鲸'

          this.$refs.tree.setData(topology)

          containerLeafIds.forEach(id => this.$refs.tree.setExpanded(id))

          this.createWatcher()
        } catch (e) {
          console.error(e)
        }
      },
      createWatcher() {
        this.nodeUnwatch = RouterQuery.watch('node', this.setDefaultState, { immediate: true })
        this.filterUnwatch = RouterQuery.watch('keyword', (value) => {
          this.filter = value
        })
        this.handleFilter = debounce(() => {
          this.$refs.tree.filter(this.filter)
          this.filter && this.setNodeCount(this.$refs.tree.visibleNodes)
        }, 300)
      },
      destroyWatcher() {
        this.nodeUnwatch && this.nodeUnwatch()
        this.filterUnwatch && this.filterUnwatch()
      },
      setDefaultState() {
        // 非业务拓扑主页面不触发设置节点选中等，防止查询条件非预期的被清除
        if (this.$route.name !== MENU_BUSINESS_HOST_AND_SERVICE) {
          return
        }
        const defaultNode = this.getDefaultNode()
        if (defaultNode) {
          const { tree } = this.$refs
          tree.setExpanded(defaultNode.id)
          tree.setSelected(defaultNode.id, { emitEvent: true })
          this.handleDefaultExpand(defaultNode)
          // 仅对第一次设置时调整滚动位置
          !this.initialized && this.$nextTick(() => {
            this.initialized = true
            const index = tree.visibleNodes.indexOf(defaultNode)
            tree.$refs.virtualScroll.scrollPageByIndex(index)
          })
        }
      },
      getDefaultNode() {
        // 选中指定的节点
        const queryNodeId = RouterQuery.get('node', '')
        if (queryNodeId) {
          const node = this.$refs.tree.getNodeById(queryNodeId)
          if (node) {
            return node
          }
        }
        // 从其他页面跳转过来需要筛选节点，例如：删除集群模板中的服务模板
        const keyword = RouterQuery.get('keyword', '')
        if (keyword) {
          const [firstMatchedNode] = this.$refs.tree.filter(keyword.trim())
          if (firstMatchedNode) {
            return firstMatchedNode
          }
        }
        // 选中第一个节点
        const [firstNode] = this.$refs.tree.nodes
        return firstNode || null
      },
      getInstanceTopology() {
        return this.$store.dispatch('objectMainLineModule/getInstTopoInstanceNum', {
          bizId: this.bizId,
          config: {
            requestId: this.request.instance
          }
        })
      },
      getInternalTopology() {
        return this.$store.dispatch('objectMainLineModule/getInternalTopo', {
          bizId: this.bizId,
          config: {
            requestId: this.request.internal
          }
        })
      },
      async getContainerTopology() {
        const topoPath = this.$route.query.topo_path
        const topoPathQueue = topoPath?.split(',') || []

        const leadPath = topoPathQueue?.[0]?.split('-')
        const leadObj = leadPath?.[0]
        const currentClusterId = leadPath?.[1]

        // 容器拓扑的第1层级为cluster，先获取cluster拓扑
        const clusterTopo = await topologyInstanceService.getContainerTopo({
          bizId: this.bizId,
          params: {
            bk_reference_obj_id: BUILTIN_MODELS.BUSINESS,
            bk_reference_id: this.bizId
          }
        })

        let asyncTopoTreeData = null
        let parentData = null
        const leafIds = []

        // 如果查询参数打头的不是cluster则认为不处于容器拓扑，直接返回第1层级的cluster拓扑即可
        if (leadObj !== CONTAINER_OBJECTS.CLUSTER) {
          return { topo: clusterTopo, leafIds }
        }

        // 遍历查询参数的拓扑路径，获取对应的节点数据，以还原拓扑
        for (let i = 0; i < topoPathQueue.length; i++) {
          const path = topoPathQueue[i]?.split('-')
          const objId = path[0]
          const instId = Number(path[1])

          if (this.isContainerLeaf(objId) || objId === CONTAINER_OBJECTS.FOLDER) {
            break
          }

          const data = await topologyInstanceService.getContainerTopo({
            bizId: this.bizId,
            params: {
              bk_reference_obj_id: objId,
              bk_reference_id: instId
            }
          })

          // 加载第1个层级的时候，将数据赋予asyncTopoTreeData，因为此处为引用赋值，此后对data追加child时数据
          if (!asyncTopoTreeData) {
            asyncTopoTreeData = data
          } else {
            // 之后加载的每个层级，找到与之对应在的parent并添加为child
            const foundNode = parentData.find(parent => parent.bk_inst_id === instId)
            if (foundNode && data) {
              foundNode.child = data
            }
          }

          // 数据作为下一次的parent
          parentData = data

          // 记录所有的叶子节点id
          if (data) {
            const nodeIds = data
              .filter(item => this.isContainerLeaf(item.bk_obj_id) || item.is_folder)
              .map(item => this.getNodeId(item))
            leafIds.push(...nodeIds)
          }
        }

        // 找到当前的那个cluster，将异步获取的topo追加为child
        const currentClusterTopo = clusterTopo.find(topo => topo.bk_inst_id === Number(currentClusterId))
        if (currentClusterTopo) {
          currentClusterTopo.child = asyncTopoTreeData
        }

        return { topo: clusterTopo, leafIds }
      },
      getNodeId(data) {
        // folder实际是不存在instid的，默认给了一个999，所以需要再加上上级id以确保节点id全局唯一
        if (data.is_folder) {
          return `${data.bk_obj_id}-${data.bk_inst_id}-${data.ref_id}`
        }

        return `${data.bk_obj_id}-${data.bk_inst_id}`
      },
      isContainerNode(node) {
        return node.data.is_container
      },
      isContainerFolder(node) {
        return node.data.is_folder
      },
      isLazyDisabledNode(node) {
        return !this.isContainerNode(node)
      },
      isContainerLeaf(objId) {
        return isWorkload(objId)
      },
      async lazyGetChildrenNode(node) {
        if (this.isContainerLeaf(node.data.bk_obj_id) || this.isContainerFolder(node)) {
          return {}
        }

        const topoList = await topologyInstanceService.getContainerTopo({
          bizId: this.bizId,
          params: {
            bk_reference_obj_id: node.data.bk_obj_id,
            bk_reference_id: node.data.bk_inst_id
          }
        })

        // 指定哪些是叶子节点，叶子节点不会再显示展开的小箭头
        const leafIds = topoList?.filter(item => item.is_workload || item.is_folder)?.map(item => this.getNodeId(item))

        // 待数据添加为树节点后，设置节点对应的统计数据
        setTimeout(() => {
          this.setNodeCount([node, ...node.children])
        }, 0)

        return {
          data: topoList,
          leaf: leafIds
        }
      },
      handleSelectChange(node) {
        this.$store.commit('businessHost/setSelectedNode', node)
        Bus.$emit('toggle-host-filter', false)

        const oldId = this.$route.query.node
        const oldTab = this.$route.query.tab
        // 服务实例视图参数
        const oldView = this.$route.query.view
        const newId = node.id

        let tab = oldTab
        let view = oldView
        if (node.data?.is_container && oldTab === 'serviceInstance') {
          tab = ''
          view = ''
        }
        if (!node.data?.is_container && oldTab === 'podList') {
          tab = ''
        }

        const query = {
          node: newId,
          tab,
          view,
          page: 1,
          _t: Date.now()
        }
        if (this.isContainerNode(node)) {
          query.topo_path = this.genTopoPathByNode(node).join(',')
        } else {
          query.topo_path = undefined
        }
        RouterQuery.set(query)

        if (FilterStore.hasCondition && oldId !== newId) {
          FilterStore.setActiveCollection(null)
        }
      },
      genTopoPathByNode(node) {
        const path = []
        let currentNode = node

        while (currentNode.parent) {
          path.push(this.getNodeId(currentNode.data))
          currentNode = currentNode.parent
        }

        return path.reverse()
      },
      handleDefaultExpand(node) {
        const nodes = []
        let parentNode = node
        while (parentNode) {
          nodes.push(...parentNode.children)
          if (!parentNode.parent) {
            nodes.push(parentNode)
          }
          parentNode = parentNode.parent
        }
        this.setNodeCount(nodes)
      },
      handleExpandChange(node) {
        if (!node.expanded || this.isContainerNode(node)) return
        this.setNodeCount([node, ...node.children])
      },
      async setNodeCount(targetNodes, force = false) {
        const nodes = force
          ? targetNodes
          : targetNodes.filter(({ data }) => !['pending', 'finished'].includes(data.status))

        if (!nodes.length) return

        nodes.forEach(({ data }) => this.$set(data, 'status', 'pending'))

        const normalNodes = []
        const containerNodes = []
        nodes.forEach((node) => {
          if (this.isContainerNode(node)) {
            containerNodes.push(node)
          } else {
            normalNodes.push(node)
          }
        })

        // targetNodes可能同时存在有传统节点和容器节点，仅当存在对应的节点时才去获取其统计数据
        if (normalNodes.length) {
          this.setNormalNodeCount(normalNodes)
        }

        if (containerNodes.length) {
          this.setContainerNodeCount(containerNodes)
        }
      },
      async setNormalNodeCount(nodes) {
        try {
          const result = await this.$store.dispatch('objectMainLineModule/getTopoStatistics', {
            bizId: this.bizId,
            params: {
              condition: nodes.map(({ data }) => ({ bk_obj_id: data.bk_obj_id, bk_inst_id: data.bk_inst_id }))
            }
          })
          nodes.forEach(({ data }) => {
            // eslint-disable-next-line
            const count = result.find(count => count.bk_obj_id === data.bk_obj_id && count.bk_inst_id === data.bk_inst_id)
            this.$set(data, 'status', 'finished')
            this.$set(data, 'host_count', count.host_count)
            this.$set(data, 'service_instance_count', count.service_instance_count)
          })
        } catch (error) {
          console.error(error)
          nodes.forEach((node) => {
            this.$set(node.data, 'status', 'error')
          })
        }
      },
      async setContainerNodeCount(nodes) {
        try {
          const params = {
            bizId: this.bizId,
            params: {
              resource_info: nodes.map(({ data }) => ({
                kind: data.bk_obj_id,
                // folder传递的是上级clusterid
                id: data.is_folder ? data.ref_id : data.bk_inst_id
              }))
            }
          }
          const { hostStats, podStats } = await topologyInstanceService.getContainerTopoNodeStats(params)
          nodes.forEach(({ data }) => {
            const finder = (item) => {
              if (data.is_folder) {
                return item.kind === data.bk_obj_id && item.id === data.ref_id
              }
              return item.kind === data.bk_obj_id && item.id === data.bk_inst_id
            }
            const hostStat = hostStats.find(finder)
            const podStat = podStats.find(finder)
            this.$set(data, 'status', 'finished')
            this.$set(data, 'host_count', hostStat?.count)
            this.$set(data, 'pod_count', podStat?.count)
          })
        } catch (error) {
          console.error(error)
          nodes.forEach((node) => {
            this.$set(node.data, 'status', 'error')
          })
        }
      },
      async getBlueKingEditStatus() {
        try {
          this.editable = await this.$store.dispatch('getBlueKingEditStatus', {
            config: {
              globalError: false
            }
          })
          this.$store.commit('businessHost/setBlueKingEditable', this.editable)
        } catch (_) {
          this.editable = false
        }
      },
      async handleShowCreateDialog(node) {
        const nodeModel = this.topologyModels.find(data => data.bk_obj_id === node.data.bk_obj_id)
        const nextModelId = nodeModel.bk_next_obj
        this.createInfo.nextModelId = nextModelId
        this.createInfo.parentNode = node
        this.createInfo.show = true
        this.createInfo.visible = true
        let properties = this.propertyMap[nextModelId]
        if (!properties) {
          const action = 'objectModelProperty/searchObjectAttribute'
          properties = await this.$store.dispatch(action, {
            params: {
              bk_biz_id: this.bizId,
              bk_obj_id: nextModelId,
              bk_supplier_account: this.$store.getters.supplierAccount
            },
            config: {
              requestId: this.request.property
            }
          })
          if (!['set', 'module'].includes(nextModelId)) {
            this.$store.commit('businessHost/setProperties', {
              id: nextModelId,
              properties
            })
          }
        }
        const primaryKey = { set: 'bk_set_id', module: 'bk_module_id' }[nextModelId] || 'bk_inst_id'
        this.createInfo.properties = properties.filter(property => property.bk_property_id !== primaryKey)
      },
      handleAfterCancelCreateNode() {
        this.createInfo.visible = false
        this.createInfo.properties = []
        this.createInfo.parentNode = null
        this.createInfo.nextModelId = null
      },
      handleCancelCreateNode() {
        this.createInfo.show = false
      },
      async handleCreateNode(value) {
        try {
          const { parentNode } = this.createInfo
          const formData = {
            ...value,
            bk_biz_id: this.bizId,
            bk_parent_id: parentNode.data.bk_inst_id
          }
          const { nextModelId } = this.createInfo
          const nextModel = this.topologyModels.find(model => model.bk_obj_id === nextModelId)
          const handlerMap = {
            set: this.createSet,
            module: this.createModule
          }
          const data = await (handlerMap[nextModelId] || this.createCommonInstance)(formData)
          const nodeData = {
            default: 0,
            child: [],
            bk_obj_name: nextModel.bk_obj_name,
            bk_obj_id: nextModel.bk_obj_id,
            host_count: 0,
            service_instance_count: 0,
            service_template_id: value.service_template_id,
            status: 'finished',
            ...data
          }
          this.$refs.tree.addNode(nodeData, parentNode.id, parentNode.data.bk_obj_id === 'biz' ? 1 : 0)
          this.$success(this.$t('新建成功'))
          this.createInfo.show = false
        } catch (e) {
          console.error(e)
        }
      },
      async handleCreateSetNode(value) {
        try {
          const { parentNode } = this.createInfo
          const nextModel = this.topologyModels.find(model => model.bk_obj_id === 'set')
          const formData = (value.sets || []).map(set => ({
            ...set,
            bk_biz_id: this.bizId,
            bk_parent_id: parentNode.data.bk_inst_id
          }))
          const data = await this.createSet(formData)
          const insertBasic = parentNode.data.bk_obj_id === 'biz' ? 1 : 0
          data && data.forEach((set, index) => {
            if (set.data) {
              const nodeData = {
                default: 0,
                child: [],
                bk_obj_name: nextModel.bk_obj_name,
                bk_obj_id: nextModel.bk_obj_id,
                host_count: 0,
                service_instance_count: 0,
                service_template_id: value.service_template_id,
                bk_inst_id: set.data.bk_set_id,
                bk_inst_name: set.data.bk_set_name,
                set_template_id: value.set_template_id,
                status: 'finished'
              }
              this.$refs.tree.addNode(nodeData, parentNode.id, insertBasic + index)
              if (value.set_template_id) {
                this.addModulesInSetTemplate(nodeData, set.data.bk_set_id)
              }
            } else {
              this.$error(set.error_message)
            }
          })
          this.$success(this.$t('新建成功'))
          this.createInfo.show = false
        } catch (e) {
          console.error(e)
        }
      },
      async addModulesInSetTemplate(parentNodeData, id) {
        const modules = await this.$store.dispatch('objectModule/searchModule', {
          bizId: this.bizId,
          setId: id,
          params: { bk_biz_id: this.bizId },
          config: {
            requestId: 'searchModule'
          }
        })
        const parentNodeId = this.getNodeId(parentNodeData)
        const nextModel = this.topologyModels.find(model => model.bk_obj_id === 'module')
        modules.info && modules.info.forEach((_module) => {
          const nodeData = {
            default: 0,
            child: [],
            bk_obj_name: nextModel.bk_obj_name,
            bk_obj_id: nextModel.bk_obj_id,
            host_count: 0,
            service_instance_count: 0,
            service_template_id: _module.service_template_id,
            bk_inst_id: _module.bk_module_id,
            bk_inst_name: _module.bk_module_name,
            status: 'finished'
          }
          this.$refs.tree.addNode(nodeData, parentNodeId, 0)
        })
      },
      async createSet(value) {
        const data = await this.$store.dispatch('objectSet/createset', {
          bizId: this.bizId,
          params: {
            sets: value.map(set => ({
              ...set,
              bk_supplier_account: this.supplierAccount
            }))
          }
        })
        return data || []
      },
      async createModule(value) {
        const data = await this.$store.dispatch('objectModule/createModule', {
          bizId: this.bizId,
          setId: this.createInfo.parentNode.data.bk_inst_id,
          params: {
            ...value,
            bk_biz_id: this.bizId,
            bk_supplier_account: this.supplierAccount
          }
        })
        return {
          bk_inst_id: data.bk_module_id,
          bk_inst_name: data.bk_module_name
        }
      },
      async createCommonInstance(value) {
        const data = await this.$store.dispatch('objectCommonInst/createInst', {
          objId: this.createInfo.nextModelId,
          params: value
        })
        return {
          bk_inst_id: data.bk_inst_id,
          bk_inst_name: data.bk_inst_name
        }
      },
      async refreshCount({ hosts, target }) {
        const nodes = []
        if (target) {
          const node = this.$refs.tree.getNodeById(`${target.data.bk_obj_id}-${target.data.bk_inst_id}`)
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
        this.setNodeCount(uniqueNodes, true)
      },
      refreshCountByNode(node) {
        const currentNode = node || this.selectedNode
        const nodes = []
        const treeNode = this.$refs.tree.getNodeById(currentNode.id)
        if (treeNode) {
          nodes.push(treeNode, ...treeNode.parents)
        }
        this.setNodeCount(nodes, true)
      },
      handleResize() {
        this.$refs.tree.resize()
      }
    }
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
    }
</style>
