import Vis from 'vis'
import uuid from 'uuid/v4'
import { svgToImageUrl } from '@/utils/util'
export const color = {
    node: {
        label: '#868b97'
    },
    edge: {
        label: '#868B97'
    }
}
const OPTIONS = {
    nodes: {
        shape: 'image',
        widthConstraint: 62,
        font: {
            color: color.node.label
        }
    },
    edges: {
        color: {
            color: '#c3cdd7',
            highlight: '#3c96ff',
            hover: '#3c96ff',
            opacity: 1
        },
        font: {
            color: color.edge.label,
            background: '#fff'
        },
        smooth: {
            type: 'curvedCW'
        },
        arrows: {
            to: true,
            from: false
        },
        arrowStrikethrough: false
    },
    interaction: {
        hover: true,
        dragNodes: false
    },
    physics: true
}

const TOOL_NODE_OPTION = {
    shape: 'image',
    value: 20,
    widthConstraint: 20,
    heightConstraint: 20,
    scaling: {
        max: 20
    },
    physics: false,
    hidden: false,
    fixed: true,
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
        color: '#ffb23a',
        highlight: '#ffb23a',
        hover: '#ffb23a'
    },
    hoverWidth: 1.5,
    width: 1.5,
    smooth: {
        roundness: 0
    },
    dashes: true
}

const DEFAULT_STATE = {
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
    constructor (container, { nodes, edges }) {
        this.editMode = false
        this.state = { ...DEFAULT_STATE }
        this.listeners = {}
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
        this.tempEdge = null

        this.addNodes(nodes)
        this.addEdges(this.createAssignedEdges(edges))
        this.addToolNodes([{
            id: this.addEdgeTriggerId,
            icon: 'cc-line',
            label: 'line'
        }, {
            id: this.deleteNodeTriggerId,
            icon: 'cc-del',
            label: 'del'
        }]).then(() => {
            this.network = new Vis.Network(container, {
                nodes: this.nodes,
                edges: this.edges
            }, OPTIONS)
            this.network.on('click', data => this.handleClick(data))
            this.network.on('hoverNode', data => this.handleNodeHover(data))
            this.network.on('hoverEdge', data => this.handleEdgeHover(data))
            this.network.on('blurNode', data => this.handleNodeBlur(data))
            this.network.on('blurEdge', data => this.handleEdgeBlur(data))
            this.network.on('dragStart', data => this.handleDragStart(data))
            this.network.on('dragEnd', data => this.handleDragEnd(data))
            this.network.on('startStabilizing', data => this.handleStartStabilizing(data))
            this.network.on('stabilized', data => this.handleStabilized(data))
            this.network.stabilize()
            this.network.fit()
        })
    }

    on (type, listener) {
        if (typeof listener === 'function') {
            this.listeners[type] = listener
        }
    }

    fire (type) {
        const listener = this.listeners[type]
        if (listener) {
            return listener(...[].slice.call(arguments, 1))
        }
        return true
    }

    changeMode (isEditMode) {
        this.editMode = isEditMode
        this.network.setOptions({
            interaction: {
                dragNodes: isEditMode
            }
        })
    }

    async handleRemoveNode (nodeId) {
        try {
            const edges = this.getNodeEdges(nodeId)
            const result = await this.fire('deleteNode', nodeId, edges)
            if (result) {
                this.nodes.remove(nodeId)
            }
        } catch (e) {
            console.log(e)
        }
    }

    handleClick (data) {
        if (data.nodes.length) {
            if (data.nodes.includes(this.addEdgeTriggerId)) { // 设置开始连线状态
                this.createShadowNode(data)
                this.state.ready = true
                this.state.to = this.shadowNode.id
                this.createShadowEdge(this.shadowNode.id)
            } else if (this.state.ready) {
                this.destroyShadowEdge()
                this.destroyShadowNode()
                this.createTempEdge(data)
                this.resetState()
            } else if (data.nodes.includes(this.deleteNodeTriggerId)) {
                this.handleRemoveNode(this.state.from)
            }
            this.schedulerUpdateToolNodes({
                hidden: true,
                immediate: true
            })
        } else if (this.state.ready) {
            this.destroyShadowEdge()
            this.destroyShadowNode()
            this.resetState()
        } else if (data.edges.length) {
            const edgeId = data.edges[0]
            const edge = this.getEdge(edgeId)
            this.fire('edgeClick', edge)
            this.setCursor('default')
        }
    }

    handleNodeHover (data) {
        const nodeId = data.node
        if (!this.editMode) {
            this.network.unselectAll()
            this.network.selectNodes([data.node])
            return false
        } else {
            this.setCursor('pointer')
        }
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
        this.setCursor('default')
        this.network.unselectAll()
        this.schedulerUpdateToolNodes({
            hidden: true,
            immediate: false
        })
    }

    handleEdgeHover (data) {
        this.setCursor('pointer')
    }

    handleEdgeBlur (data) {
        this.setCursor('default')
    }

    handleDragStart (data) {
        if (!data.nodes.length) {
            return false
        }
        if (this.triggerIds.includes(data.nodes[0])) {
            return false
        }
        this.nodes.update({
            id: data.nodes[0],
            fixed: false
        })
        if (this.editMode) {
            this.schedulerUpdateToolNodes({
                hidden: true,
                immediate: true
            })
        }
    }

    handleDragEnd (data) {
        if (!data.nodes.length) {
            return false
        }
        const [nodeId] = data.nodes
        if (this.triggerIds.includes(nodeId)) {
            return false
        }
        this.bounceOverlapNodes(nodeId)
        this.fire('dragNode', nodeId, this.network.getPositions([nodeId])[nodeId])
        if (this.editMode) {
            this.updateToolNodePosition(nodeId)
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
            return { id, fixed: true }
        }))
        this.fire('stabilized', this.network.getPositions())
    }

    handleStartStabilizing (data) {
        this.state.stabilizing = true
        this.schedulerUpdateToolNodes({
            hidden: true,
            immediate: true
        })
    }

    bounceOverlapNodes (referenceId) {
        const referenceBox = this.network.getBoundingBox(referenceId)
        const overlapNodes = [referenceId]
        this.normalNodeIds.forEach(id => {
            if (id === referenceId) {
                return false
            }
            const targetBox = this.network.getBoundingBox(id)
            const isNotOverlap
                = targetBox.top > referenceBox.bottom
                || targetBox.right < referenceBox.left
                || targetBox.bottom < referenceBox.top
                || targetBox.left > referenceBox.right
            if (!isNotOverlap) {
                overlapNodes.push(id)
            }
        })
        if (overlapNodes.length > 1) {
            this.nodes.update(overlapNodes.map(id => {
                return { id, fixed: false }
            }))
            this.overlapNodes = overlapNodes
            this.network.setOptions({ physics: true })
        } else {
            this.nodes.update({
                id: referenceId,
                fixed: true
            })
        }
    }

    updateToolNodePosition (refNodeId) {
        if (this.state.stabilizing) {
            return false
        }
        const nodeR = OPTIONS.nodes.widthConstraint / 2
        const toolR = TOOL_NODE_OPTION.widthConstraint / 2
        const { x, y } = this.network.getPositions([refNodeId])[refNodeId]
        const deltaXY = Math.sqrt(2) / 2 * (nodeR + toolR)
        const addEdgeTriggerX = x + deltaXY
        const addEdgeTriggerY = y - deltaXY
        this.network.moveNode(this.addEdgeTriggerId, addEdgeTriggerX, addEdgeTriggerY)
        this.network.moveNode(this.deleteNodeTriggerId, addEdgeTriggerX + toolR, y)
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

    schedulerUpdateToolNodes ({ hidden, immediate }) {
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
        if (!this.shadowNode) {
            return false
        }
        this.nodes.update({
            id: this.shadowNode.id,
            ...this.convertNodePosition(event)
        })
    }

    addNodes (data, defaultOption = {}) {
        const mergedData = Utils.getMergedData(data, defaultOption)
        this.nodes.add(mergedData)
    }

    removeNodes (ids) {
        return this.nodes.remove(ids)
    }

    // vis在拥有相同方向箭头的情况下，连线会重叠，通过重新分配方向及弧度使线不重叠
    createAssignedEdges (edges) {
        const groups = {}
        // 1.将两两节点相同的连线聚类
        edges.forEach(edge => {
            const forward = `${edge.from}-${edge.to}`
            const reverse = `${edge.to}-${edge.from}`
            const groupId = groups[forward] ? forward : reverse
            if (groups[groupId]) {
                groups[groupId].push(edge)
            } else {
                groups[groupId] = [edge]
            }
        })
        // 2.平均分配不同方向的连线，并将分配的连线from、to、箭头方向反转
        const assignedEdges = []
        Object.keys(groups).forEach(groupId => {
            const fromId = groups[groupId][0]['from']
            const edges = this.reassignEdges(groups[groupId])
            Array.prototype.push.apply(assignedEdges, edges)
        })
        return assignedEdges
    }

    // 重新分配连续的起点终点，箭头方向
    reassignEdges (edges) {
        if (!edges.length) {
            return edges
        }
        const fromId = edges[0]['from']
        const forwardEdges = edges.filter(edge => edge.from === fromId)
        const reverseEdges = edges.filter(edge => edge.to === fromId)
        const forwardCount = forwardEdges.length
        const reverseCount = reverseEdges.length
        const lessCount = Math.abs(forwardCount - reverseCount)
        if (lessCount > 1) {
            const MoreEdges = forwardCount > reverseCount ? forwardEdges : reverseEdges
            const lessEdges = forwardCount > reverseCount ? reverseEdges : forwardEdges
            const countToAssign = Math.floor(lessCount / 2)
            const edgesToAssign = MoreEdges.splice(MoreEdges.length - countToAssign)
            Array.prototype.push.apply(lessEdges, edgesToAssign.map(edge => {
                const arrows = edge.arrows || { from: false, to: true }
                return {
                    ...edge,
                    from: edge.to,
                    to: edge.from,
                    arrows: {
                        from: arrows.to,
                        to: arrows.from
                    }
                }
            }))
        }
        this.setEdgeRoundness(forwardEdges, forwardCount, reverseCount)
        this.setEdgeRoundness(reverseEdges, forwardCount, reverseCount)
        return [...forwardEdges, ...reverseEdges]
    }

    setEdgeRoundness (edges, forwardCount, reverseCount) {
        const averageCount = Math.floor((forwardCount + reverseCount) / 2)
        const edgesCount = edges.length
        edges.forEach((edge, index) => {
            edge.smooth = {
                roundness: (edgesCount > averageCount) ? (index * 0.1) : ((index + 1) * 0.1)
            }
        })
        return edges
    }

    addEdges (data) {
        this.edges.update(data)
    }

    async addToolNodes (data) {
        const nodes = await Promise.all(data.map(config => {
            return new Promise(resolve => {
                const image = new Image()
                image.onload = () => {
                    resolve({
                        id: config.id,
                        image: svgToImageUrl(image, {
                            iconColor: '#fff',
                            backgroundColor: 'rgba(24, 24, 24, .8)'
                        })
                    })
                }
                image.onerror = () => {
                    resolve({
                        id: config.id,
                        shape: 'box',
                        label: config.label
                    })
                }
                image.src = `${window.location.origin}/static/svg/${config.icon}.svg`
            })
        }))
        this.addNodes(nodes, TOOL_NODE_OPTION)
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

    showNode (nodeId, position = {}) {
        this.nodes.update({
            id: nodeId,
            hidden: false,
            ...position
        })
        this.bounceOverlapNodes(nodeId)
    }

    getNodeEdges (nodeId) {
        const edgeData = this.edges['_data']
        const edges = []
        Object.keys(edgeData).forEach(id => {
            const edge = edgeData[id]
            if (edge.from === nodeId || edge.to === nodeId) {
                edges.push(edge)
            }
        })
        return edges
    }

    async createTempEdge (data) {
        try {
            const toNode = data.nodes[0]
            const tempEdge = {
                id: uuid(),
                from: this.state.from,
                to: toNode,
                color: SHADOW_EDGE_OPTION.color
            }
            const existEdges = this.getExisitEdges(tempEdge)
            const reassignEdges = this.reassignEdges([tempEdge, ...existEdges])
            this.edges.update(reassignEdges)
            this.network.unselectAll()
            this.tempEdge = tempEdge
            const result = await this.fire('addEdge', {
                id: tempEdge.id,
                from: this.state.from,
                to: toNode
            })
            if (result) {
                this.createRealEdge(result, existEdges)
            } else {
                this.destroyTempEdge()
            }
        } catch (e) {
            this.destroyTempEdge()
        }
    }

    createRealEdge (data, existEdges) {
        if (this.tempEdge) {
            this.edges.remove(this.tempEdge.id)
            const reassignEdges = this.reassignEdges([{
                ...this.tempEdge,
                ...data,
                color: OPTIONS.edges.color
            }, ...existEdges])
            this.edges.update(reassignEdges)
            this.tempEdge = null
        }
    }

    destroyTempEdge () {
        const tempEdge = this.tempEdge
        if (tempEdge) {
            const existEdges = this.getExisitEdges(tempEdge).filter(edge => edge.id !== tempEdge.id)
            const reassignEdges = this.reassignEdges(existEdges)
            this.edges.remove(tempEdge.id)
            this.edges.update(reassignEdges)
        }
    }

    getExisitEdges ({ from, to }) {
        const existEdges = []
        const edgeData = this.edges['_data']
        Object.keys(edgeData).find(key => {
            const edge = edgeData[key]
            const isSame = edge.from === from && edge.to === to
            const isReverse = edge.to === from && edge.from === to
            if (isSame || isReverse) {
                existEdges.push({ ...edge })
            }
        })
        return existEdges
    }

    getEdge (edgeId) {
        return this.edges['_data'][edgeId]
    }

    deleteEdge (edgeId) {
        this.edges.remove(edgeId)
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

    convertNodePosition (event) {
        const pointer = this.network.body.functions.getPointer({
            x: event.clientX,
            y: event.clientY
        })
        return this.network.canvas.DOMtoCanvas(pointer)
    }

    resetState () {
        this.state = { ...DEFAULT_STATE }
    }

    updateOptions (options) {
        this.network.setOptions(options)
    }

    updateEdges (config) {
        this.edges.update(config)
    }

    resize () {
        this.network.fit()
    }

    zoom (type = 'in') {
        const ratio = type === 'in' ? 1.2 : 0.8
        this.network.moveTo({
            scale: this.network.getScale() * ratio,
            animation: {
                duration: 100,
                easingFunction: 'easeInOutCubic'
            }
        })
    }

    setCursor (cursor) {
        document.body.style.cursor = cursor
    }
}
