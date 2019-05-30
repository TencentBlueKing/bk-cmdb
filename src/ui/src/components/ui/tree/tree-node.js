import {
    getNodeId,
    getNodeIcon
} from './utils.js'

export default class TreeNode {
    constructor (data, options, tree) {
        this.data = data
        this.tree = tree
        this._vNode = null

        const treeOptions = tree.nodeOptions
        this.id = getNodeId(data, tree)
        this.name = data[treeOptions.nameKey]
        this.icon = getNodeIcon(data, tree)

        this.level = options.level
        this.index = options.index
        this.parent = options.parent
        this.isLeaf = true
        this.children = []

        this.state = {
            checked: false,
            expanded: false,
            visible: true
        }
        this.checked = tree.defaultCheckedNodes.includes(this.id)
        this.expanded = tree.defaultExpandedNodes.includes(this.id)

        this.timer = null
    }

    set vNode (vNode) {
        this._vNode = vNode
        if (this.expanded) {
            this.recaculateLinkLine()
        }
    }

    get vNode () {
        return this._vNode
    }

    get parents () {
        console.log(this.parent, this.parent && this.parent.parents)
        if (this.parent) {
            return []
        }
        return [...this.parent.parents, this.parent]
    }

    get collapseIcon () {
        return this.icon.collapse
    }

    get selected () {
        return this.tree.selected === this.id
    }

    get expandIcon () {
        return this.icon.expand
    }

    get nodeIcon () {
        return this.icon.node
    }

    set checked (checked) {
        if (this.state.checked === checked) {
            return false
        }
        this.state.checked = checked
    }

    get checked () {
        return this.state.checked
    }

    set expanded (expanded) {
        if (this.state.expanded === expanded) {
            return false
        }
        this.state.expanded = expanded
        if (expanded && this.parent) {
            this.parent.expanded = true
        }
        this.children.forEach(node => {
            node.visible = expanded
        })
        this.recaculateLinkLine()
    }

    get expanded () {
        return this.state.expanded
    }

    set visible (visible) {
        if (this.state.visible === visible) {
            return false
        }
        this.state.visible = visible
        this.children.forEach(node => {
            node.visible = visible
        })
    }

    get visible () {
        return this.state.visible
    }

    recaculateLinkLine () {
        if (this.tree.showLinkLine) {
            const needsCalculateNodes = this.tree.needsCalculateNodes
            if (needsCalculateNodes.includes(this)) {
                return false
            }
            needsCalculateNodes.push(this)
            this.parent && this.parent.recaculateLinkLine()
        }
    }

    appendChild (node, trailing = true) {
        const nodes = Array.isArray(node) ? node : [node]
        const oldLength = this.children.length
        if (trailing) {
            this.children.push(...nodes)
            nodes.forEach((node, index) => {
                node.childIndex = oldLength + index
            })
        } else {
            this.children.unshift(...nodes)
            this.children.forEach((node, index) => {
                node.childIndex = index
            })
        }

        this.isLeaf = false
        this.expanded = true
        this.recaculateLinkLine()
        return nodes
    }

    removeChild (node) {
        const nodes = Array.isArray(node) ? node : [node]
        const removedChildIndex = []
        const removedIndex = []
        nodes.forEach(node => {
            const childIndex = node.childIndex
            removedChildIndex.push(childIndex)
            removedIndex.push(node.index)
            this.children.splice(childIndex, 1)
        })
        const minIndex = Math.min(...removedChildIndex)
        this.children.slice(minIndex).forEach((node, index) => {
            node.childIndex = minIndex + index
        })
        this.isLeaf = !this.children.length
        this.recaculateLinkLine()
        return nodes
    }
}
