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
                            shape: 'circle',
                            widthConstraint: 55,
                            borderWidth: 0,
                            color: {
                                background: '#6b7baa'
                            },
                            font: {
                                color: '#fff'
                            }
                        },
                        edges: {
                            color: {
                                color: '#6b7baa',
                                highlight: '#6b7baa',
                                hover: '#6b7baa'
                            },
                            smooth: {
                                type: 'curvedCW',
                                roundness: 0
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
                    this.$alertMsg(e)
                }
            },
            setNodes (data) {
                this.network.nodes = data.map(nodeData => {
                    const node = {
                        id: nodeData['bk_obj_id'],
                        label: nodeData['node_name'],
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
            setSingleNodePosition () {
                const edges = this.network.edges
                const noPositionSingleNode = this.network.nodes.filter(node => {
                    const isSingle = !edges.some(edge => edge.from === node.id || edge.to === node.id)
                    const hasPosition = node.hasOwnProperty('x') && node.hasOwnProperty('y')
                    return isSingle && !hasPosition
                })
                const fixedNodeY = -1 * this.$el.offsetHeight / 2
                const fixedDistance = 70
                noPositionSingleNode.forEach((node, index) => {
                    node.x = index % 2 === 0 ? index / 2 * fixedDistance : -1 * (index + 1) / 2 * fixedDistance
                    node.y = fixedNodeY
                })
                this.networkDataSet.nodes.update(noPositionSingleNode)
            },
            // 判断是否是双向关联A.to = B.from && A.from = B.to
            getTwoWayAsst (node, asst, edges) {
                return edges.find(edge => edge.from === asst['bk_obj_id'] && edge.to === node['bk_obj_id'])
            },
            loadNodeImage () {
                this.network.nodes.forEach(node => {
                    let image = new Image()
                    image.onload = () => {
                        this.networkDataSet.nodes.update({
                            id: node.id,
                            shape: 'image',
                            image: `data:image/svg+xml;charset=utf-8,${encodeURIComponent(GET_OBJ_ICON(image, {
                                name: node.label,
                                backgroundColor: '#6b7baa',
                                color: '#fff'
                            }))}`,
                            font: {
                                size: 0
                            }
                        })
                    }
                    image.src = `${window.location.origin}/static/svg/${node['data']['bk_obj_icon'].substr(5)}.svg`
                })
            },
            updateInitPosition () {
                let updateNodes = this.network.nodes.filter(node => {
                    const position = node.data.position
                    return position.x === null && position.y === null
                })
                this.updateNodePosition(updateNodes)
            },
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
            addListener () {
                const networkInstance = this.networkInstance
                networkInstance.once('stabilized', () => {
                    networkInstance.setOptions({
                        physics: {
                            enabled: false
                        }
                    })
                    networkInstance.on('dragEnd', (params) => {
                        if (params.nodes.length) {
                            this.updateNodePosition(this.networkDataSet.nodes.get(params.nodes))
                        }
                    })
                    this.setSingleNodePosition()
                    this.loadNodeImage()
                    this.updateInitPosition()
                    this.loading = false
                })
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