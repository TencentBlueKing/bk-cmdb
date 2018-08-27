<template>
    <li class="tree-node-item"
        :class="{
            'tree-node-leaf': leaf,
            'tree-node-expanded': state.expanded,
            'tree-node-hidden': state.hidden
        }">
        <div class="tree-node-info-layout clearfix" :class="{'tree-node-info-layout-root': level === 1}">
            <i class="tree-node-expanded-icon fl" v-if="!leaf"
                :class="[state.expanded ? 'icon-cc-rect-sub': 'icon-cc-rect-add']"
                @click="handleExpandClick">
            </i>
            <tree-node-info class="tree-node-info" @click.native="handleNodeClick"
                :class="{
                    'tree-node-info-leaf': leaf,
                    'tree-node-info-selected': state.selected
                }"
                :state="state"
                :layout="layout">
            </tree-node-info>
        </div>
        <ul v-if="!leaf" v-show="state.expanded" class="tree-node-children">
            <cmdb-tree-node v-for="(child, index) in children"
                :key="index"
                :node="child"
                :layout="layout"
                :level="level + 1">
            </cmdb-tree-node>
        </ul>
    </li>
</template>

<script>
    import TreeLayout from './_tree-layout.js'
    import treeNodeInfo from './_tree-node-info.js'
    export default {
        name: 'cmdb-tree-node',
        components: {
            treeNodeInfo
        },
        props: {
            node: {
                type: Object,
                required: true
            },
            layout: {
                validator (val) {
                    return val instanceof TreeLayout
                },
                required: true
            },
            level: {
                type: Number,
                required: true
            }
        },
        data () {
            const state = {
                id: this.layout.getNodeId(this.node),
                disabled: false,
                expanded: false,
                hidden: false,
                selected: false,
                level: this.level,
                parent: this.level === 1 ? null : this.$parent,
                node: this.node
            }
            return {
                state
            }
        },
        computed: {
            treeInstance () {
                return this.layout.instance
            },
            children () {
                return this.node[this.treeInstance.childrenKey] || []
            },
            leaf () {
                return !this.children.length
            }
        },
        created () {
            this.layout.addFlatternNode(this.state)
            if (this.node.selected) {
                this.handleNodeClick()
            }
            if (this.node.expanded) {
                this.handleExpandClick()
            }
        },
        methods: {
            handleExpandClick () {
                this.layout.toggleExpanded(this.state.id, !this.state.expanded)
                this.treeInstance.$emit('on-expand', this.node, this.state)
            },
            handleNodeClick () {
                this.layout.selectNode(this.state.id)
                this.treeInstance.$emit('on-selected', this.node, this.state)
            }
        }
    }
</script>

<style lang="scss" scoped>
    .tree-node-item {
        position: relative;
        white-space: nowrap;
        font-size: 0;
        margin: 8px 0;
        &.tree-node-expanded:before{
                position: absolute;
                left: 7px;
                top: 19px;
                width: 0;
                height: calc(100% - 31px);
                content: '';
                border-left: 1px dashed #d3d8e7;
                z-index: 1;
        }
        .tree-node-expanded-icon{
            display: block;
            margin: 5px 0 0 0;
            font-size: 14px;
            color: #c3cdd7;
            cursor: pointer;
            &:hover{
                color: #3c96ff;
            }
        }
    }
    .tree-node-info-layout{
        position: relative;
        &:not(.tree-node-info-layout-root):before{
            position: absolute;
            top: 12px;
            left: -15px;
            width: 20px;
            height: 0;
            content: '';
            border-top: 1px dashed #d3d8e7;
            z-index: 1;
        }
        .tree-node-info{
            height: 24px;
            padding: 0 0 0 14px;
            line-height: 24px;
            font-size: 14px;
            overflow: hidden;
            cursor: pointer;
            &:hover{
                background-color: #f1f7ff;
                color: #498fe0;
            }
            &.tree-node-info-leaf{
                margin: 0 0 0 14px;
            }
            &.tree-node-info-selected{
                background-color: #e2efff;
                color: #498fe0;
            }
        }
    }
    .tree-node-children{
        margin: 0 0 0 24px;
    }
    .test{
        display: inline-block;
    }
</style>