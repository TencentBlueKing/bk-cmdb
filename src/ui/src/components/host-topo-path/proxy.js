/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2017-2022 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
 */

import Vue from 'vue'
import http from '@/api'
import containerHostService from '@/service/container/host'
import hostSearchService from '@/service/host/search'

export default new Vue({
  data() {
    return {
      results: new Map(),
      queue: [],
      isContainerSearchMode: false,
      isResourceAssigned: false,
      containerTopoPaths: []
    }
  },
  watch: {
    queue(queue) {
      queue.length && this.dispatchRequest()
    }
  },
  methods: {
    search(data) {
      return new Promise((resolve) => {
        this.queue.push({
          resolve,
          data
        })
      })
    },
    dispatchRequest() {
      this.requestTimer && clearTimeout(this.requestTimer)
      this.requestTimer = setTimeout(async () => {
        // 一次执行请求后将队列清空
        const queue = this.queue.splice(0)
        // 由于不能跨业务批量查询，此处拆分不同业务下的模块
        const requestMap = new Map()

        const hostIds = []

        queue.forEach((meta) => {
          const requestSet = requestMap.get(meta.data.bk_biz_id)
          if (requestSet) {
            requestSet.add(meta)
          } else {
            const newRequestSet = new Set()
            newRequestSet.add(meta)
            requestMap.set(meta.data.bk_biz_id, newRequestSet)
          }

          if (!hostIds.some(id => id === meta.data.hostId)) {
            hostIds.push(meta.data.hostId)
          }
        })

        // 容器拓扑路径
        if (this.isResourceAssigned) {
          const { info: paths } = await containerHostService.getNodePath({ ids: hostIds })
          this.containerTopoPaths = paths || []

          // 已分配并且为容器搜索模式，拓扑路径采用通过主机ids获取的方式
          if (this.isContainerSearchMode) {
            // 通过主机获取获取传统拓扑路径
            const normalTopoPaths = await hostSearchService.getTopoPath({ ids: hostIds }) || []

            queue.forEach((meta) => {
              const containerPaths = this.containerTopoPaths.find(item => item.bk_host_id === meta.data.hostId)?.paths

              // 将数据格式统一为`find/topopath/biz/${bizId}`接口格式，便于后续使用
              const hostNormalTopoList = normalTopoPaths.find(item => item.id === meta.data.hostId)?.topo_path
              const hostNormalTopoPaths = hostNormalTopoList.map(topo => ({
                topo_node: topo.find(item => item.bk_obj_id === 'module'),
                topo_path: topo
              }))

              meta.resolve({
                normal: hostNormalTopoPaths,
                container: containerPaths
              })
            })
          }
        }

        if (!this.isContainerSearchMode) {
          // 不同业务分别发起请求
          requestMap.forEach((requestSet, bizId) => {
            const newModules = this.seperateNode(bizId, requestSet)
            if (!newModules.length) {
              this.resolvePromise(bizId, requestSet)
            } else {
              http.post(`find/topopath/biz/${bizId}`, {
                topo_nodes: newModules
              }).then((data) => {
                this.setResults(bizId, data.nodes)
                this.resolvePromise(bizId, requestSet)
              })
            }
          })
        }
      }, 200)
    },
    // 将业务下的模块去重，并去掉已缓存的模块id
    seperateNode(bizId, requestSet) {
      const modules = []
      requestSet.forEach(({ data }) => {
        modules.push(...(data.modules ?? []))
      })
      const uniqueModules = [...new Set(modules)]
      // eslint-disable-next-line max-len
      const newModules = uniqueModules.filter(moduleId => !(this.results.has(bizId) && this.results.get(bizId).has(moduleId)))
      return newModules.map(moduleId => ({
        bk_obj_id: 'module',
        bk_inst_id: moduleId
      }))
    },
    // 将请求结果用map缓存起来
    setResults(bizId, nodes) {
      nodes.forEach((node) => {
        const resultMap = this.results.get(bizId)
        if (resultMap) {
          resultMap.set(node.topo_node.bk_inst_id, node)
        } else {
          const newResultMap = new Map()
          newResultMap.set(node.topo_node.bk_inst_id, node)
          this.results.set(bizId, newResultMap)
        }
      })
    },
    // resolve返回给每个组件的Promise
    resolvePromise(bizId, requestSet) {
      const resultMap = this.results.get(bizId) || new Map()
      requestSet.forEach((meta) => {
        const result = meta.data.modules?.map(moduleId => resultMap.get(moduleId) || {
          bk_biz_id: bizId,
          topo_node: {
            bk_inst_id: moduleId,
            bk_obj_id: 'module'
          },
          topo_path: []
        })

        const containerPaths = this.containerTopoPaths.find(item => item.bk_host_id === meta.data.hostId)?.paths
        meta.resolve({
          normal: result,
          container: containerPaths
        })
      })
    }
  }
})
