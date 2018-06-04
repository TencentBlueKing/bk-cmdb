<template>
    <ul class="tree-list">
       <v-tree-item v-if="treeData['bk_inst_id']"
        :treeData="treeData"
        :initNode="initNode"
        :rootNode="treeData"
        :activeNodeId="activeNodeId"
        :treeId="treeId"
        :model="model"
        :hideRoot="hideRoot"
        :level="1"></v-tree-item>
    </ul>
</template>

<script type="text/javascript">
    import vTreeItem from './treeItem.v2'
    import bus from '@/eventbus/bus'
    export default {
        props: {
            treeData: Object,
            initNode: Object,
            hideRoot: Boolean,
            model: {
                type: Array,
                default () {
                    return []
                }
            }
        },
        data () {
            return {
                activeNodeId: null
            }
        },
        computed: {
            treeId () {
                return Math.random()
            }
        },
        created () {
            bus.$on('nodeClick', (activeNode, nodeOptions) => {
                if (this.treeId === nodeOptions.treeId) {
                    this.activeNodeId = nodeOptions.nodeId
                    this.$emit('nodeClick', activeNode, nodeOptions)
                }
            })
            bus.$on('nodeToggle', (isOpen, node, nodeOptions) => {
                if (this.treeId === nodeOptions.treeId) {
                    this.$emit('nodeToggle', isOpen, node, nodeOptions)
                }
            })
            bus.$on('addNode', (node, nodeOptions) => {
                if (this.treeId === nodeOptions.treeId) {
                    this.$emit('addNode', node, nodeOptions)
                }
            })
        },
        components: {
            vTreeItem
        }
    }
</script>