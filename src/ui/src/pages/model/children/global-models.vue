<template>
    <div class="global-model" v-bkloading="{isLoading: loading}"></div>
</template>
<script>
    import Vis from 'vis'
    import { generateObjIcon as GET_OBJ_ICON } from '@/utils/util'
    export default {
        data () {
            return {
                loading: true,
                networkInstance: null,
                networkDataSet: {
                    nodes: null,
                    edges: null
                },
                network: {
                    nodes: null,
                    edges: null,
                    options: {
                        nodes: {
                            shape: 'image',
                            widthConstraint: 55,
                            shadow: {
                                enabled: true,
                                color: 'rgba(0,0,0,0.1)',
                                x: 0,
                                y: 2,
                                size: 4
                            }
                        },
                        edges: {
                            color: {
                                color: '#c3cdd7',
                                highlight: '#3c96ff'
                            },
                            smooth: {
                                type: 'curvedCW',
                                roundness: 0
                            },
                            arrows: {
                                to: {
                                    scaleFactor: 0.6
                                },
                                from: {
                                    scaleFactor: 0.6
                                }
                            }
                        },
                        physics: {
                            enabled: true,
                            barnesHut: {
                                avoidOverlap: 0.5,
                                springLength: 150
                            }
                        }
                    }
                }
            }
        },
        computed: {
            noPositionNodes () {
                return this.network.nodes.filter(node => {
                    const position = node.data.position
                    return position.x === null && position.y === null
                })
            }
        },
        mounted () {
            this.initNetwork()
        },
        methods: {
            async initNetwork () {
                try {
                    const response = await this.$axios.post('objects/topographics/scope_type/global/scope_id/0/action/search')
                    if (response.result) {
                        this.setNodes(response.data)
                        this.setEdges(response.data)
                        this.networkInstance = new Vis.Network(this.$el, {
                            nodes: this.networkDataSet.nodes,
                            edges: this.networkDataSet.edges
                        }, this.network.options)
                        this.addListener()
                    } else {
                        this.$alertMsg(response['bk_error_msg'])
                    }
                } catch (e) {
                    this.$alertMsg(e.message)
                }
            },
            // 设置节点数据
            setNodes (data) {
                this.network.nodes = data.map(nodeData => {
                    const node = {
                        id: nodeData['bk_obj_id'],
                        image: `data:image/svg+xml;charset=utf-8,${encodeURIComponent(GET_OBJ_ICON({
                            name: nodeData['node_name'],
                            backgroundColor: '#fff',
                            fontColor: nodeData['ispre'] ? '#6894c8' : '#868b97'
                        }))}`,
                        data: nodeData
                    }
                    if (nodeData['position']['x'] !== null && nodeData['position']['y'] !== null) {
                        node.physics = false
                        node.x = nodeData['position']['x']
                        node.y = nodeData['position']['y']
                    }
                    return node
                })
                this.networkDataSet.nodes = new Vis.DataSet(this.network.nodes)
            },
            // 设置连线数据
            setEdges (data) {
                let edges = []
                data.forEach(node => {
                    if (Array.isArray(node.assts) && node.assts.length) {
                        node.assts.forEach(asst => {
                            const twoWayAsst = this.getTwoWayAsst(node, asst, edges)
                            if (twoWayAsst) { // 双向关联，将已存在的线改为双向
                                twoWayAsst.arrows = 'to,from'
                                twoWayAsst.label = [twoWayAsst.label, asst['bk_object_att_id']].join(',\n')
                            } else {
                                edges.push({
                                    from: node['bk_obj_id'],
                                    to: asst['bk_obj_id'],
                                    arrows: 'to',
                                    label: asst['bk_object_att_id']
                                })
                            }
                        })
                    }
                })
                this.network.edges = edges
                this.networkDataSet.edges = new Vis.DataSet(this.network.edges)
            },
            // 设置单节点位置并更新其节点位置信息
            setSingleNodePosition () {
                const edges = this.network.edges
                const noPositionSingleNode = this.network.nodes.filter(node => !edges.some(edge => edge.from === node.id || edge.to === node.id))
                const fixedNodeY = this.getSingleNodePositionY()
                const fixedDistance = 70
                const compensateDistance = noPositionSingleNode.length % 2 === 0 ? fixedDistance / 2 : 0
                const middleIndex = Math.floor(noPositionSingleNode.length / 2)
                noPositionSingleNode.forEach((node, index) => {
                    node.x = (index - middleIndex) * fixedDistance - compensateDistance
                    node.y = fixedNodeY
                })
                if (noPositionSingleNode.length) {
                    this.networkDataSet.nodes.update(noPositionSingleNode)
                }
            },
            // 获取单节点摆放的初始化Y轴坐标
            getSingleNodePositionY () {
                const asstNodes = this.network.nodes.filter(node => this.network.edges.some(edge => edge.to === node.id || edge.from === node.id))
                if (asstNodes.length) {
                    const asstNodePositions = this.networkInstance.getPositions(asstNodes.map(node => node.id))
                    return Math.min(...asstNodes.map(node => asstNodePositions[node.id]['y'])) - 70
                }
                return 0
            },
            // 判断是否是双向关联A.to = B.from && A.from = B.to
            getTwoWayAsst (node, asst, edges) {
                return edges.find(edge => edge.from === asst['bk_obj_id'] && edge.to === node['bk_obj_id'])
            },
            // 加载节点icon并更新
            loadNodeImage () {
                this.network.nodes.forEach(node => {
                    let image = new Image()
                    image.onload = () => {
                        this.networkDataSet.nodes.update({
                            id: node.id,
                            image: `data:image/svg+xml;charset=utf-8,${encodeURIComponent(GET_OBJ_ICON(image, {
                                name: node.data['node_name'],
                                fontColor: node.data.ispre ? '#6894c8' : '#868b97',
                                iconColor: node.data.ispre ? '#6894c8' : '#868b97',
                                backgroundColor: '#fff'
                            }))}`
                        })
                    }
                    image.src = `${window.location.origin}/static/svg/${node['data']['bk_obj_icon'].substr(5)}.svg`
                })
            },
            // 批量更新节点位置信息
            updateNodePosition (updateNodes) {
                if (!updateNodes.length) return
                const nodePositions = this.networkInstance.getPositions(updateNodes.map(node => node.id))
                const updateParams = updateNodes.map(node => {
                    const nodeData = node.data
                    return {
                        'bk_obj_id': node.id,
                        'bk_inst_id': nodeData['bk_inst_id'],
                        'node_type': nodeData['node_type'],
                        'position': {
                            x: nodePositions[node.id]['x'],
                            y: nodePositions[node.id]['y']
                        }
                    }
                })
                this.$axios.post('objects/topographics/scope_type/global/scope_id/0/action/update', updateParams).then(res => {
                    if (!res.result) {
                        this.$alert(res['bk_error_msg'])
                    }
                })
            },
            // 拓扑稳定后执行事件
            // 1.取消物理模拟
            // 2.配置拖拽结束监听，更新位置信息
            // 3.设置无位置信息的单节点位置
            // 4.加载节点图标
            listenerCallback () {
                this.networkInstance.setOptions({
                    physics: {
                        enabled: false
                    }
                })
                this.networkInstance.on('dragEnd', (params) => {
                    if (params.nodes.length) {
                        this.updateNodePosition(this.networkDataSet.nodes.get(params.nodes))
                    }
                    this.networkInstance.unselectAll()
                })
                this.setSingleNodePosition()
                this.loadNodeImage()
                this.networkInstance.fit()
                this.loading = false
            },
            addListener () {
                const networkInstance = this.networkInstance
                networkInstance.once('stabilized', this.listenerCallback)
                if (!this.noPositionNodes.length) {
                    this.listenerCallback()
                }
            }
        }

    }
</script>
<style lang="scss" scoped>
    .global-model{
        width: 100%;
        height: 100%;
        background-color: #f4f5f8;
        background-image: linear-gradient(#eef1f5 1px, transparent 0), linear-gradient(90deg, #eef1f5 1px, transparent 0);
        background-size: 10px 10px;
    }
</style>