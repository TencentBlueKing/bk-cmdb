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
            bus.$on('nodeClick', (activeNode, activeParentNode, activeRootNode, activeLevel, nodeId, treeId) => {
                if (this.treeId === treeId) {
                    this.activeNodeId = nodeId
                    this.$emit('nodeClick', activeNode, activeParentNode, activeRootNode, activeLevel)
                }
            })
            bus.$on('nodeToggle', (isOpen, node, parentNode, rootNode, level, nodeId, treeId) => {
                if (this.treeId === treeId) {
                    this.$emit('nodeToggle', isOpen, node, parentNode, rootNode, level, nodeId)
                }
            })
            bus.$on('addNode', (node, parentNode, rootNode, level, nodeId, treeId) => {
                if (this.treeId === treeId) {
                    this.$emit('addNode', node, parentNode, rootNode, level)
                }
            })
        },
        components: {
            vTreeItem
        }
    }
</script>