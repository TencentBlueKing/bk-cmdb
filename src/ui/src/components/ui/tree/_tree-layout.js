let ID_SEED = 1
class TreeLayout {
    constructor (options) {
        this.instance = null
        this.flatternStates = {}
        this.selectedState = null
        this.expandedStates = {}
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
        const state = this.flatternStates[id] || null
        return state
    }

    selectState (id) {
        const flatternStates = this.flatternStates
        for (let key in this.flatternStates) {
            const state = flatternStates[key]
            if (state.id === id) {
                state.selected = true
                this.selectedState = state
            } else {
                state.selected = false
            }
            if (state.node.hasOwnProperty('selected')) {
                state.node.selected = state.selected
            }
        }
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
                this.expandedStates[id] = state
            } else {
                delete this.expandedStates[id]
            }
            if (state.node.hasOwnProperty('expanded')) {
                state.node.expanded = expanded
            }
        } else {
            console.error('state lost, cannot toggle expand with node id:' + id)
        }
    }

    addFlatternState (state) {
        this.flatternStates[state.id] = state
    }

    removeFlatternState (state) {
        delete this.flatternStates[state.id]
    }

    removeExpandedState (state) {
        delete this.expandedStates[state.id]
    }

    destory (state) {
        this.removeFlatternState(state)
        this.removeExpandedState(state)
    }
}

export default TreeLayout
