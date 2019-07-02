<template>
    <div class="cmdb-tree">
        <tree-item v-for="node in nodes"
            :key="node.id"
            :node="node">
            <slot
                :node="node"
                :data="node.data">
            </slot>
        </tree-item>
    </div>
</template>

<script>
    import TreeNode from './tree-node.js'
    import treeItem from './tree-item.vue'
    export default {
        name: 'cmdb-tree',
        components: {
            treeItem
        },
        props: {
            data: {
                type: Array,
                default () {
                    return []
                }
            },
            options: {
                type: Object,
                default () {
                    return {}
                }
            },
            lazy: Boolean,
            selectable: {
                type: Boolean,
                default: true
            },
            showCheckbox: [Boolean, Function],
            showLinkLine: {
                type: Boolean,
                default: true
            },
            expandIcon: {
                type: String,
                default: 'bk-icon icon-folder-open'
            },
            collapseIcon: {
                type: String,
                default: 'bk-icon icon-folder'
            },
            nodeIcon: {
                type: [String, Function]
            },
            defaultExpandedNodes: {
                type: Array,
                default () {
                    return []
                }
            },
            defaultCheckedNodes: {
                type: Array,
                default () {
                    return []
                }
            },
            defaultSelectedNode: {
                type: [String, Number],
                default: null
            },
            defaultDisabledNodes: {
                type: Array,
                default () {
                    return []
                }
            },
            beforeSelect: Function,
            beforeCheck: Function
        },
        data () {
            return {
                nodes: [],
                map: {},
                selected: this.defaultSelectedNode,
                checked: [...this.defaultCheckedNodes],
                needsCalculateNodes: [],
                calculateTimer: null
            }
        },
        computed: {
            nodeOptions () {
                const nodeOptions = {
                    idKey: 'id',
                    nameKey: 'name',
                    childrenKey: 'children'
                }
                return Object.assign(nodeOptions, this.options)
            }
        },
        watch: {
            needsCalculateNodes () {
                this.calulateLine()
            }
        },
        created () {
            this.setData(this.data)
        },
        methods: {
            setData (data) {
                const nodes = []
                const map = {}
                this.recurrenceNodes(data, null, nodes, map)
                this.nodes = nodes
                this.map = map
            },
            recurrenceNodes (data, parent, nodes, map) {
                data.forEach((datum, index) => {
                    const node = new TreeNode(datum, {
                        level: parent ? parent.level + 1 : 0,
                        parent: parent,
                        index: nodes.length
                    }, this)
                    if (parent) {
                        node.childIndex = parent.children.length
                        parent.children.push(node)
                        parent.isLeaf = false
                    }
                    nodes.push(node)
                    map[node.id] = node

                    const children = datum[this.nodeOptions.childrenKey]
                    if (Array.isArray(children) && children.length) {
                        this.recurrenceNodes(children, node, nodes, map)
                    }
                })
            },
            getNodeById (id) {
                return this.map[id]
            },
            addNode (nodeData, parentId, trailing = true) {
                const parent = this.getNodeById(parentId)
                if (!parent) {
                    throw new Error('Unexpected parent id, add node failed')
                }
                const children = parent.children
                const offset = typeof trailing === 'number' ? trailing : trailing ? children.length : 0
                let insertIndex
                if (offset > 0) {
                    const referenceChild = children[offset - 1]
                    const referenceChildWithDescendants = [referenceChild, ...referenceChild.descendants]
                    const referenceNode = referenceChildWithDescendants[referenceChildWithDescendants.length - 1]
                    insertIndex = referenceNode.index + 1
                } else {
                    insertIndex = parent.index + 1
                }
                const data = Array.isArray(nodeData) ? nodeData : [nodeData]
                const nodes = data.map(datum => {
                    return new TreeNode(datum, {
                        level: parent.level + 1,
                        parent: parent
                    }, this)
                })
                parent.appendChild(nodes, offset)
                nodes.forEach(node => {
                    this.$set(this.map, node.id, node)
                })
                this.nodes.splice(insertIndex, 0, ...nodes)
                this.nodes.slice(insertIndex).forEach((node, index) => {
                    node.index = insertIndex + index
                })
            },
            removeNode (nodeId) {
                const ids = Array.isArray(nodeId) ? nodeId : [nodeId]
                const nodes = []
                ids.forEach(id => {
                    const node = this.getNodeById(id)
                    if (node) {
                        nodes.push(node)
                    }
                })
                // 从最大的node.index开始倒序splice
                nodes.sort((M, N) => N.index - M.index)
                nodes.forEach(node => {
                    const removeNodes = [node, ...node.descendants]
                    this.nodes.splice(node.index, removeNodes.length)
                    if (node.parent) {
                        node.parent.removeChild(node)
                    }
                })
                const minChangedIndex = Math.min(...nodes.map(node => node.index))
                this.nodes.slice(minChangedIndex).forEach((node, index) => {
                    node.index = minChangedIndex + index
                })
            },
            async setSelected (nodeId, options = {}) {
                if (!this.selectable || nodeId === this.selected) {
                    return false
                }
                const mergeOptions = {
                    emitEvent: false,
                    beforeSelect: true,
                    ...options
                }
                const node = this.getNodeById(nodeId)
                if (mergeOptions.beforeSelect && typeof this.beforeSelect === 'function') {
                    const response = await this.beforeSelect(node)
                    if (!response) {
                        return false
                    }
                }
                this.selected = nodeId
                if (mergeOptions.emitEvent) {
                    this.$emit('select-change', node)
                }
            },
            async removeChecked (options = {}) {
                const mergeOptions = {
                    emitEvent: true,
                    ...options
                }
                this.checked.forEach(id => {
                    const node = this.getNodeById(id)
                    node.checked = false
                })
                this.checked = []
                if (mergeOptions.emitEvent) {
                    this.$emit('check-change', [], [])
                }
            },
            async setChecked (nodeId, options = {}) {
                const ids = Array.isArray(nodeId) ? nodeId : [nodeId]
                if (ids.length) {
                    const mergeOptions = {
                        emitEvent: false,
                        beforeCheck: true,
                        checked: true,
                        ...options
                    }
                    const nodes = ids.map(id => this.getNodeById(id))
                    if (mergeOptions.beforeCheck && typeof this.beforeCheck === 'function') {
                        const response = await this.beforeCheck(nodes.length > 1 ? nodes : nodes[0], mergeOptions.checked)
                        if (!response) {
                            return false
                        }
                    }
                    nodes.forEach(node => {
                        node.checked = mergeOptions.checked
                    })
                    if (mergeOptions.checked) {
                        this.checked = [...new Set([...this.checked, ...ids])]
                    } else {
                        this.checked = this.checked.filter(id => !ids.includes(id))
                    }
                    if (mergeOptions.emitEvent) {
                        this.$emit('check-change', this.checked, this.checked.map(id => this.getNodeById(id)))
                    }
                }
            },
            setExpanded (nodeId, options = {}) {
                const mergeOptions = {
                    expanded: true,
                    emitEvent: false,
                    ...options
                }
                const node = this.getNodeById(nodeId)
                if (!node) {
                    throw new Error('Unexpected node id, set expanded failed.')
                }
                node.expanded = mergeOptions.expanded
                if (mergeOptions.emitEvent) {
                    this.$emit('expand-change', mergeOptions.expanded, node)
                }
            },
            setDisabled (nodeId, options = {}) {
                const mergeOptions = {
                    disabled: true,
                    emitEvent: false,
                    ...options
                }
                const node = this.getNodeById(nodeId)
                if (!node) {
                    throw new Error('Unexpected node id, set disabled failed.')
                }
                node.disabled = mergeOptions.disabled
                if (mergeOptions.emitEvent) {
                    this.$emit('disable-change', mergeOptions.disabled, node)
                }
            },
            calulateLine () {
                this.calculateTimer && clearTimeout(this.calculateTimer)
                if (this.needsCalculateNodes.length) {
                    this.calculateTimer = setTimeout(() => {
                        this.needsCalculateNodes.forEach(node => {
                            node.vNode.calulateLine()
                        })
                        this.needsCalculateNodes.splice(0)
                    }, 0)
                } else {
                    this.calculateTimer = null
                }
            }
        }
    }
</script>
