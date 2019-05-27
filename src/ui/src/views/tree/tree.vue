<template>
    <div class="cmdb-tree">
        <tree-item v-for="node in nodes"
            :key="node.id"
            :node="node">
            <slot :node="node"></slot>
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
            defaultExpandNode: {
                type: Array,
                default () {
                    return []
                }
            }
        },
        data () {
            return {
                nodes: [],
                map: {}
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
                        parent: parent
                    }, this)
                    node.setState('expanded', this.defaultExpandNode.includes(node.id))
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
            addNode (nodeData, parentId, leading = true) {
                const parentNode = this.getNodeById(parentId)
                if (!parentNode) {
                    throw new Error('Unexpected parent id, add node failed')
                }
                const parentVNode = parentNode.vNode
                const data = Array.isArray(nodeData) ? nodeData : [nodeData]
                const insertIndex = (leading ? parentVNode.index : parentVNode.index + parentVNode.children.length) + 1
                const nodes = data.map((datum, index) => {
                    return new TreeNode(datum, {
                        level: parentNode.level + 1,
                        parent: parentNode,
                        children: []
                    }, this)
                })
                nodes.forEach(node => {
                    this.$set(this.map, node.id, node)
                })
                this.nodes.splice(insertIndex, 0, ...nodes)
            },
            removeNode (nodeId) {
                const ids = Array.isArray(nodeId) ? nodeId : [nodeId]
                const deletedIndex = []
                const deletedNodeParent = []
                ids.forEach(id => {
                    const node = this.getNodeById(id)
                    if (node) {
                        deletedIndex.push(id)
                        this.$delete(this.map, id)
                        this.nodes.splice(node.index, 1)
                        if (node.parent) {
                            node.parent.children.splice(node.vNode.childIndex, 1)
                            deletedNodeParent.push(node.parent)
                        }
                    }
                })
                this.nodes.slice(Math.min(...deletedIndex)).forEach(node => {
                    const referenceIndex = deletedIndex.findIndex(index => index > node.index)
                    const changedStep = referenceIndex > -1 ? referenceIndex - 1 : deletedIndex.length
                    node.index = node.index - changedStep
                })
                setTimeout(() => {
                    deletedNodeParent.forEach(parent => {
                        parent.vNode.handleDescendantsChange()
                    })
                })
            }
        }
    }
</script>
