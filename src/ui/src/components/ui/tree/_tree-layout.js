let ID_SEED = 1
class TreeLayout {
    constructor (options) {
        this.instance = null
        this.flatternNodes = []
        this.selectedNode = null
        this.expandedNodes = []
        this.expandedLevel = 0
        for (let name in options) {
            if (options.hasOwnProperty(name)) {
                this[name] = options[name]
            }
        }
    }

    getNodeId (node) {
        if (typeof this.instance.idGenerator === 'function') {
            return this.instance.idGenerator(node)
        } else {
            const idKey = this.instance.idKey
            let nodeId
            if (node.hasOwnProperty(idKey)) {
                nodeId = node[idKey]
            } else {
                nodeId = ID_SEED++
            }
            return nodeId
        }
    }

    getNodeById (id) {
        return this.flatternNodes.find(node => node.id === id)
    }

    selectNode (id) {
        this.flatternNodes.forEach(node => {
            if (node.id === id) {
                node.selected = true
                this.selectedNode = node
            } else {
                node.selected = false
            }
        })
    }

    toggleExpanded (id, expanded) {
        const node = this.getNodeById(id)
        if (node) {
            node.expanded = expanded
            if (expanded) {
                this.expandedNodes.push(node)
            } else {
                this.expandedNodes = this.expandedNodes.filter(node => node.id !== id)
            }
            this.expandedLevel = Math.max.apply(null, this.expandedNodes.map(node => node.level))
        }
    }

    addFlatternNode (node) {
        this.flatternNodes.push(node)
    }

    removeFlatternNode (node) {
        this.flatternNodes = this.flatternNodes.filter(flatternNode => flatternNode !== node)
    }

    removeExpandedNode (node) {
        this.expandedNodes = this.expandedNodes.filter(expandedNode => expandedNode !== node)
    }

    destory (node) {
        this.removeFlatternNode(node)
        this.removeExpandedNode(node)
    }
}

export default TreeLayout
