<template>
    <div class="cmdb-tree"></div>
</template>

<script>
    import TreeNode from './tree-node.js'
    export default {
        name: 'cmdb-tree',
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
            }
        },
        data () {
            return {
                flatternTree: [],
                nodeMap: {}
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
                try {
                    const flatternTree = []
                    data.forEach(nodeData => {
                        const node = new TreeNode(nodeData, this.getNodeOptions({
                            level: 0,
                            parent: null
                        }))

                        flatternTree.push(node)
                        this.$set(this.nodeMap, node.id, node)

                        this.recurrenceNodes(flatternTree, node)
                    })
                    this.flatternTree = flatternTree
                    this.updateIndex()
                } catch (e) {
                    console.error(e)
                    this.flatternTree = []
                }
            },
            recurrenceNodes (tree, parent) {
                const children = parent.data[this.nodeOptions.childrenKey]
                if (Array.isArray(children)) {
                    children.forEach(nodeData => {
                        const node = new TreeNode(nodeData, this.getNodeOptions({
                            level: parent.level + 1,
                            parent: parent.id
                        }))

                        tree.push(node)
                        this.$set(this.nodeMap, node.id, node)

                        this.recurrenceNodes(tree, node)
                    })
                }
            },
            getNodeById (id) {
                return this.nodeMap[id]
            },
            removeNode (id) {
                const node = this.getNodeById(id)
                this.flatternTree.splice(node.index, 1)
                this.$delete(this.nodeMap, id)
                this.updateIndex()
            },
            addNode (data, parentId = null, trailing = true) {
                if (parentId !== null) {
                    const parentNode = this.getNodeById(parentId)
                    if (parentNode) {
                        const children = parentNode.data[this.nodeOptions.childrenKey]
                        const indexIncrement = Array.isArray(children) ? children.length : 0
                        const index = trailing ? parentNode.index + indexIncrement + 1 : parentNode.index + 1
                        const node = new TreeNode(data, this.getNodeOptions({
                            level: parentNode.level + 1
                        }))
                        this.flatternTree.splice(index - 1, 0, node)
                        this.$set(this.nodeMap, node.id, node)
                        this.updateIndex()
                    } else {
                        throw new Error('Cant not find parent node with id:', parentId)
                    }
                } else {
                    const node = new TreeNode(data, this.getNodeOptions({
                        level: 0
                    }))
                    const index = trailing ? this.flatternTree.length - 1 : 0
                    this.flatternTree.splice(index, 0, node)
                    this.$set(this.nodeMap, node.id, node)
                    this.updateIndex()
                }
            },
            getNodeOptions (options = {}) {
                return Object.assign({}, this.nodeOptions, options)
            },
            updateIndex (index) {
                this.flatternTree.forEach((node, index) => {
                    node.setValue('index', index)
                })
            }
        }
    }
</script>
