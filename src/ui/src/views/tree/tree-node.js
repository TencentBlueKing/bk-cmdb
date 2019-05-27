import {
    getNodeId,
    getNodeIcon
} from './utils.js'

const DEFAULT_STATE = {
    level: 0,
    index: 0,
    parent: null,
    isLeaf: false,
    expanded: false,
    checked: false,
    children: [],
    visible: true
}

export default class TreeNode {
    constructor (data, state, tree) {
        this.data = data
        this.vNode = null

        const options = tree.nodeOptions
        this.id = getNodeId(data, tree)
        this.name = data[options.nameKey]
        this.icon = getNodeIcon(data, tree)

        Object.keys(DEFAULT_STATE).forEach(key => {
            this[key] = state.hasOwnProperty(key) ? state[key] : DEFAULT_STATE[key]
        })
    }

    setState (key, value) {
        if (DEFAULT_STATE.hasOwnProperty(key)) {
            this[key] = value
        }
    }
}
