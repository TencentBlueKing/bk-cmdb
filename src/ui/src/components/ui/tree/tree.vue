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
            showCheckbox: Boolean,
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
            beforeSelect: Function
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
                console.log(this.nodes.length)
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
                const insertIndex = (trailing ? parent.index + parent.children.length : parent.index) + 1
                const data = Array.isArray(nodeData) ? nodeData : [nodeData]
                const nodes = data.map(datum => {
                    return new TreeNode(datum, {
                        level: parent.level + 1,
                        parent: parent
                    }, this)
                })
                parent.appendChild(nodes, trailing)
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
                const deletedIndex = []
                ids.forEach(id => {
                    const node = this.getNodeById(id)
                    if (node) {
                        deletedIndex.push(id)
                        this.$delete(this.map, id)
                        this.nodes.splice(node.index, 1)
                        if (node.parent) {
                            node.parent.removeChild(node)
                        }
                    }
                })
                const changedIndex = Math.min(...deletedIndex)
                this.nodes.slice(changedIndex).forEach((node, index) => {
                    node.index = changedIndex + index
                })
            },
            async setSelected (nodeId, byEvent = false) {
                if (nodeId === this.selected) {
                    return false
                }
                const node = this.getNodeById(nodeId)
                if (byEvent && typeof this.beforeSelect === 'function') {
                    const confirmSelect = await this.beforeSelect(node)
                    if (confirmSelect) {
                        this.selected = nodeId
                        this.$emit('select-change', node)
                    }
                } else {
                    this.selected = nodeId
                }
            },
            async setChecked (nodeId, checked, byEvent = false) {
                const node = this.getNodeById(nodeId)
                if (!node) {
                    throw new Error('Unexpected node id, set checked failed.')
                }
                const index = this.checked.indexOf(nodeId)
                if ((checked && index > -1) || (!checked && index < 0)) {
                    return false
                }
                if (byEvent && typeof this.beforeCheck === 'function') {
                    const confirmCheck = await this.beforeCheck(node)
                    if (!confirmCheck) {
                        return false
                    }
                }
                node.checked = checked
                if (checked) {
                    this.checked.push(nodeId)
                } else {
                    this.checked.splice(index, 1)
                }
                if (byEvent) {
                    this.$emit('check-change', this.checked, node)
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
