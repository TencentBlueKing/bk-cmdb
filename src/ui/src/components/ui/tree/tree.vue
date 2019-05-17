<template>
    <ul class="tree-list">
        <cmdb-tree-node v-for="(node, index) in tree"
            :key="layout.getNodeId(node)"
            :node="node"
            :level="1"
            :layout="layout">
        </cmdb-tree-node>
    </ul>
</template>

<script>
    import cmdbTreeNode from './_tree-node.vue'
    import TreeLayout from './_tree-layout.js'
    export default {
        name: 'cmdb-tree',
        components: {
            cmdbTreeNode
        },
        props: {
            tree: {
                type: Array,
                default () {
                    return []
                }
            },
            render: {
                type: Function
            },
            idGenerator: {
                type: Function
            },
            idKey: {
                type: String,
                default: 'id'
            },
            labelKey: {
                type: String,
                default: 'name'
            },
            childrenKey: {
                type: String,
                default: 'children'
            },
            beforeClick: {
                type: Function
            },
            selectable: {
                type: Boolean,
                default: true
            }
        },
        data () {
            const layout = new TreeLayout({
                instance: this
            })
            return {
                layout
            }
        },
        methods: {
            toggleExpanded (id, expanded) {
                this.layout.toggleExpanded(id, expanded)
            },
            selectNode (id) {
                this.layout.selectState(id)
            },
            unselectNode (id) {
                this.layout.unselectNode(id)
            },
            getStateById (id) {
                return this.layout.getStateById(id)
            }
        }
    }
</script>