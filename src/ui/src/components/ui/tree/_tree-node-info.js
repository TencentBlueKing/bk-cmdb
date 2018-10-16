export default {
    name: 'tree-node-info',
    props: ['state', 'layout'],
    render (h) {
        const treeInstance = this.layout.instance
        const state = this.state
        const node = state.node
        const nodeId = state.id
        if (treeInstance.$scopedSlots[nodeId]) {
            return treeInstance.$scopedSlots[nodeId]({node, state})
        } else if (treeInstance.$scopedSlots.default) {
            return treeInstance.$scopedSlots.default({node, state})
        } else if (typeof treeInstance.render === 'function') {
            return treeInstance.render(h, node, state)
        } else {
            return (<div>{node[treeInstance.labelKey]}</div>)
        }
    }
}
