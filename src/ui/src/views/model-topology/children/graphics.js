import Vis from 'vis'
import uuid from 'uuid/v4'
const OPTIONS = {
    nodes: {
        shape: 'image',
        widthConstraint: 62
    },
    edges: {
        color: {
            color: '#c3cdd7',
            highlight: '#3c96ff',
            hover: '#3c96ff'
        },
        font: {
            background: '#fff'
        },
        smooth: {
            type: 'curvedCW',
            roundness: 0.1
        },
        arrows: {
            to: {
                scaleFactor: 0.6
            },
            from: {
                scaleFactor: 0.6
            }
        },
        arrowStrikethrough: false
    },
    interaction: {
        hover: true
    },
    layout: {
        randomSeed: 522492360
    },
    physics: true
}

const TOOL_NODE_OPTION = {
    shape: 'box',
    shapeProperties: {
        borderRadius: 12
    },
    heightConstraint: 14,
    scaling: {
        max: 24
    },
    physics: false,
    hidden: true,
    fixed: true,
    color: {
        background: 'background: rgba(24, 24, 24, .8)'
    },
    font: {
        color: '#fff'
    },
    group: 'tool'
}

const SHADOW_NODE_OPTION = {
    shape: 'dot',
    value: 0,
    size: 0,
    physics: false,
    fixed: true,
    scaling: {
        min: 0,
        max: 0
    },
    group: 'shadow'
}

const SHADOW_EDGE_OPTION = {
    dashed: true,
    color: {
        color: '#ffb23a'
    }
}

const DEFAULT_STATE = {
    editMode: true,
    stabilizing: false,
    timer: null,
    ready: false,
    hidden: true,
    from: null,
    to: null
}
const Utils = {
    getMergedData (data, defaultOption) {
        data = Array.isArray(data) ? data : [data]
        return data.map(option => Object.assign({}, defaultOption, option))
    }
}

export default class Graphics {
    constructor (container, {nodes, edges}) {
        this.state = {...DEFAULT_STATE}

        this.nodes = new Vis.DataSet()
        this.edges = new Vis.DataSet()

        this.normalNodeIds = nodes.map(node => node.id)
        this.overlapNodes = []

        this.addEdgeTriggerId = uuid()
        this.addEdgeTrigger = null
        this.deleteNodeTriggerId = uuid()
        this.deleteNodeTrigger = null
        this.triggerIds = [this.addEdgeTriggerId, this.deleteNodeTriggerId]

        this.shadowNode = null
        this.shadowEdge = null

        this.addNodes(nodes)
        this.addEdges(edges)
        this.addToolNodes([{
            id: this.addEdgeTriggerId,
            label: '创建关联'
        }, {
            id: this.deleteNodeTriggerId,
            label: '删除'
        }])
        this.network = new Vis.Network(container, {
            nodes: this.nodes,
            edges: this.edges
        }, OPTIONS)
        this.network.on('hoverNode', data => this.handleNodeHover(data))
        this.network.on('blurNode', data => this.handleNodeBlur(data))
        this.network.on('selectNode', data => this.handleNodeSelect(data))
        this.network.on('dragStart', data => this.handleDragStart(data))
        this.network.on('dragEnd', data => this.handleDragEnd(data))
        this.network.on('startStabilizing', data => this.handleStartStabilizing(data))
        this.network.on('stabilized', data => this.handleStabilized(data))
        this.network.stabilize()
        this.network.fit()
    }

    handleNodeHover (data) {
        const nodeId = data.node
        if (this.state.ready) {
            return false
        }
        if (this.state.timer) {
            clearTimeout(this.state.timer)
            this.state.timer = null
        }
        if (!this.triggerIds.includes(nodeId)) {
            this.updateToolNodePosition(nodeId)
        }
    }

    handleNodeBlur (data) {
        this.schedulerUpdateToolNodes({
            hidden: true,
            immediate: false
        })
    }

    handleNodeSelect (data) {
        if (data.nodes.includes(this.addEdgeTriggerId)) { // 设置开始连线状态
            this.createShadowNode(data)
            this.state.ready = true
            this.state.to = this.shadowNode.id
            this.createShadowEdge(this.shadowNode.id)
        } else if (this.state.ready) {
            this.destroyShadowEdge()
            this.destroyShadowNode()
            this.createEdges(data)
            this.resetstate()
        } else if (data.nodes.includes(this.deleteNodeTriggerId)) {
            this.removeNode()
        }
        this.schedulerUpdateToolNodes({
            hidden: true,
            immediate: true
        })
    }

    handleDragStart (data) {
        if (!data.nodes.length) { return false }
        if (this.triggerIds.includes(data.nodes[0])) {
            return false
        }
        this.nodes.update({
            id: data.nodes[0],
            fixed: false
        })
        if (this.state.editMode) {
            this.schedulerUpdateToolNodes({
                hidden: true,
                immediate: true
            })
        }
    }

    handleDragEnd (data) {
        if (!data.nodes.length) { return false }
        if (this.triggerIds.includes(data.nodes[0])) { return false }
        this.bounceOverlapNodes(data)
        if (this.state.editMode) {
            this.updateToolNodePosition(data.nodes[0])
        }
    }

    handleStabilized (data) {
        this.state.stabilizing = false
        this.network.setOptions({
            physics: false,
            nodes: {
                fixed: true
            }
        })
        this.nodes.update(this.overlapNodes.map(id => {
            return {id, fixed: true}
        }))
    }

    handleStartStabilizing (data) {
        this.state.stabilizing = true
        this.schedulerUpdateToolNodes({
            hidden: true,
            immediate: true
        })
    }

    bounceOverlapNodes (data) {
        const referenceId = data.nodes[0]
        const referenceBox = this.network.getBoundingBox(referenceId)
        const overlapNodes = [referenceId]
        this.normalNodeIds.forEach(id => {
            if (id === referenceId) { return false }
            const targetBox = this.network.getBoundingBox(id)
            const isNotOverlap =
                targetBox.top > referenceBox.bottom ||
                targetBox.right < referenceBox.left ||
                targetBox.bottom < referenceBox.top ||
                targetBox.left > referenceBox.right
            if (!isNotOverlap) {
                overlapNodes.push(id)
            }
        })
        if (overlapNodes.length > 1) {
            this.nodes.update(overlapNodes.map(id => {
                return {id, fixed: false}
            }))
            this.overlapNodes = overlapNodes
            this.network.setOptions({physics: true})
        } else {
            this.nodes.update({
                id: referenceId,
                fixed: true
            })
        }
    }

    updateToolNodePosition (refNodeId) {
        if (this.state.stabilizing) { return false }
        const nodeR = OPTIONS.nodes.widthConstraint / 2
        const {x, y} = this.network.getPositions([refNodeId])[refNodeId]
        const addEdgeRect = this.getNodeRect(this.addEdgeTriggerId)
        const addEdgeTriggerX = x + Math.sqrt(2) / 2 * nodeR + addEdgeRect.width / 2
        const addEdgeTriggerY = y - Math.sqrt(2) / 2 * nodeR
        this.network.moveNode(this.addEdgeTriggerId, addEdgeTriggerX, addEdgeTriggerY)
        this.network.moveNode(this.deleteNodeTriggerId, addEdgeTriggerX + 30, addEdgeTriggerY + 30)
        this.state.from = refNodeId
        this.schedulerUpdateToolNodes({
            hidden: false,
            immediate: true
        })
    }

    getNodeRect (nodeId) {
        const box = this.network.getBoundingBox(nodeId)
        return {
            width: box.right - box.left,
            height: box.bottom - box.top
        }
    }

    schedulerUpdateToolNodes ({hidden, immediate}) {
        if (this.state.timer) {
            clearTimeout(this.state.timer)
            this.state.timer = null
        }
        if (immediate) {
            this.updateToolNodes(hidden)
        } else {
            this.state.timer = setTimeout(() => {
                this.updateToolNodes(hidden)
            }, 300)
        }
    }

    updateToolNodes (hidden) {
        this.nodes.update([{
            id: this.addEdgeTriggerId,
            hidden
        }, {
            id: this.deleteNodeTriggerId,
            hidden
        }])
        this.state.timer = null
        this.state.hidden = hidden
    }

    shadowNodeFollowMouse (event) {
        if (!this.shadowNode) { return false }
        const pointer = this.network.body.functions.getPointer({
            x: event.clientX,
            y: event.clientY
        })
        this.nodes.update({
            id: this.shadowNode.id,
            ...this.network.canvas.DOMtoCanvas(pointer)
        })
    }

    addNodes (data, defaultOption = {}) {
        const mergedData = Utils.getMergedData(data, defaultOption)
        this.nodes.add(mergedData)
    }

    removeNodes (ids) {
        return this.nodes.remove(ids)
    }

    addEdges (data) {
        this.edges.add(data)
    }

    addToolNodes (data) {
        this.addNodes(data, TOOL_NODE_OPTION)
    }

    createShadowNode (data) {
        this.shadowNode = {
            id: uuid(),
            ...data.pointer.canvas,
            ...SHADOW_NODE_OPTION
        }
        this.nodes.add(this.shadowNode)
    }

    destroyShadowNode () {
        this.nodes.remove(this.shadowNode.id)
        this.shadowNode = null
    }

    createEdges (data) {
        this.edges.add({
            from: this.state.from,
            to: data.nodes[0]
        })
        this.network.unselectAll()
    }

    createShadowEdge () {
        this.shadowEdge = {
            id: uuid(),
            ...SHADOW_EDGE_OPTION,
            from: this.state.from,
            to: this.state.to
        }
        this.edges.add(this.shadowEdge)
    }

    destroyShadowEdge () {
        this.edges.remove(this.shadowEdge.id)
        this.shadowEdge = null
    }

    resetstate () {
        this.state = {...DEFAULT_STATE}
    }
}
