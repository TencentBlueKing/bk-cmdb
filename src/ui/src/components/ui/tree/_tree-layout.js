let ID_SEED = 1
class TreeLayout {
    constructor (options) {
        this.instance = null
        this.flatternStates = []
        this.selectedState = null
        this.expandedStates = []
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

    getStateById (id) {
        return this.flatternStates.find(state => state.id === id)
    }

    selectState (id) {
        this.flatternStates.forEach(state => {
            if (state.id === id) {
                state.selected = true
                this.selectedState = state
            } else {
                state.selected = false
            }
            if (state.node.hasOwnProperty('selected')) {
                state.node.selected = state.selected
            }
        })
    }

    unselectState (id) {
        const state = this.getStateById(id)
        state.selected = false
        this.selectedState = null
        if (state.node.hasOwnProperty('selected')) {
            state.node.selected = false
        }
    }

    toggleExpanded (id, expanded) {
        const state = this.getStateById(id)
        if (state) {
            state.expanded = expanded
            if (expanded) {
                this.expandedStates.push(state)
            } else {
                this.expandedStates = this.expandedStates.filter(state => state.id !== id)
            }
            this.expandedLevel = Math.max.apply(null, this.expandedStates.map(state => state.level))
            if (state.node.hasOwnProperty('expanded')) {
                state.node.expanded = expanded
            }
        } else {
            console.error('state lost, cannot toggle expand with node id:' + id)
        }
    }

    addFlatternState (state) {
        this.flatternStates.push(state)
    }

    removeFlatternState (state) {
        this.flatternStates = this.flatternStates.filter(flatternState => flatternState !== state)
    }

    removeExpandedState (state) {
        this.expandedStates = this.expandedStates.filter(expandedState => expandedState !== state)
    }

    destory (state) {
        this.removeFlatternState(state)
        this.removeExpandedState(state)
    }
}

export default TreeLayout
