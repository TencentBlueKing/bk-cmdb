import Vue from 'vue'
import http from '@/api'

export default new Vue({
    data () {
        return {
            results: new Map(),
            queue: []
        }
    },
    watch: {
        queue (queue) {
            queue.length && this.dispatchRequest()
        }
    },
    methods: {
        search (data) {
            return new Promise(resolve => {
                this.queue.push({
                    resolve,
                    data
                })
            })
        },
        dispatchRequest () {
            this.requestTimer && clearTimeout(this.requestTimer)
            this.requestTimer = setTimeout(() => {
                // 一次执行请求后将队列清空
                const queue = this.queue.splice(0)
                // 由于不能跨业务批量查询，此处拆分不同业务下的模块
                const requestMap = new Map()
                queue.forEach(meta => {
                    const requestSet = requestMap.get(meta.data.bk_biz_id)
                    if (requestSet) {
                        requestSet.add(meta)
                    } else {
                        const newRequestSet = new Set()
                        newRequestSet.add(meta)
                        requestMap.set(meta.data.bk_biz_id, newRequestSet)
                    }
                })
                // 不同业务分别发起请求
                requestMap.forEach((requestSet, bizId) => {
                    const newModules = this.seperateNode(bizId, requestSet)
                    if (!newModules.length) {
                        this.resolvePromise(bizId, requestSet)
                    } else {
                        http.post(`find/topopath/biz/${bizId}`, {
                            topo_nodes: newModules
                        }).then(data => {
                            this.setResults(bizId, data.nodes)
                            this.resolvePromise(bizId, requestSet)
                        })
                    }
                })
            }, 200)
        },
        // 将业务下的模块去重，并去掉已缓存的模块id
        seperateNode (bizId, requestSet) {
            const modules = []
            requestSet.forEach(({ data }) => {
                modules.push(...data.modules)
            })
            const uniqueModules = [...new Set(modules)]
            const newModules = uniqueModules.filter(moduleId => {
                return !(this.results.has(bizId) && this.results.get(bizId).has(moduleId))
            })
            return newModules.map(moduleId => ({
                bk_obj_id: 'module',
                bk_inst_id: moduleId
            }))
        },
        // 将请求结果用map缓存起来
        setResults (bizId, nodes) {
            nodes.forEach(node => {
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
        resolvePromise (bizId, requestSet) {
            const resultMap = this.results.get(bizId) || new Map()
            requestSet.forEach(meta => {
                const result = meta.data.modules.map(moduleId => {
                    return resultMap.get(moduleId) || {
                        bk_biz_id: bizId,
                        topo_node: {
                            bk_inst_id: moduleId,
                            bk_obj_id: 'module'
                        },
                        topo_path: []
                    }
                })
                meta.resolve(result)
            })
        }
    }
})
