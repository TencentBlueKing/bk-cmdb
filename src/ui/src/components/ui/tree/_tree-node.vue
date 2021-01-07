<template>
    <li class="tree-node-item"
        :class="{
            'tree-node-leaf': leaf,
            'tree-node-expanded': state.expanded,
            'tree-node-hidden': state.hidden
        }">
        <div class="tree-node-info-layout clearfix"
            :class="{
                'tree-node-info-layout-line': level !== 1 && !(state.parent && state.parent.state.hidden),
                'tree-node-info-layout-line-short-v': index === 0,
                'tree-node-info-layout-line-short-h': children.length
            }">
            <i class="tree-node-expanded-icon fl" v-if="!leaf && !state.hidden"
                :class="[state.expanded ? 'icon-cc-rect-sub': 'icon-cc-rect-add']"
                @click="handleExpandClick(!state.expanded)">
            </i>
            <tree-node-info class="tree-node-info"
                v-if="!state.hidden"
                :class="{
                    'tree-node-info-leaf': leaf,
                    'tree-node-info-selected': state.selected
                }"
                :state="state"
                :layout="layout"
                @click.native="handleNodeClick">
            </tree-node-info>
        </div>
        <ul class="tree-node-children"
            v-if="!leaf"
            v-show="state.expanded"
            :class="{'tree-node-children-root': level === 1}">
            <cmdb-tree-node v-for="(child, index) in children"
                :key="layout.getNodeId(child)"
                :index="index"
                :node="child"
                :layout="layout"
                :level="level + 1">
            </cmdb-tree-node>
        </ul>
    </li>
</template>

<script>
    import treeNodeInfo from './_tree-node-info.js'
    export default {
        name: 'cmdb-tree-node',
        components: {
            treeNodeInfo
        },
        props: {
            node: Object,
            layout: Object,
            level: Number,
            index: Number
        },
        data () {
            const basicState = {
                disabled: false,
                expanded: false,
                visible: false,
                selected: false,
                hidden: false
            }
            const state = {
                ...basicState,
                id: this.layout.getNodeId(this.node),
                level: this.level,
                parent: this.level === 1 ? null : this.$parent,
                node: this.node
            }
            return {
                basicState,
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
        watch: {
            node: {
                deep: true,
                handler (val, oldVal) {
                    let stateChanged = Object.keys(this.basicState).some(key => val[key] !== oldVal[key])
                    if (stateChanged) {
                        this.updateBasicState()
                    }
                }
            },
            'state.disabled' (disabled) {
                this.treeInstance.$emit('on-disable-change', this.node, this.state, disabled)
            },
            'state.expanded' (expanded) {
                this.treeInstance.$emit('on-expand', this.node, this.state, expanded)
            },
            'state.visible' (visible) {
                this.treeInstance.$emit('on-visible-change', this.node, this.state, visible)
            },
            'state.selected' (selected) {
                if (selected) {
                    this.treeInstance.$emit('on-selected', this.node, this.state)
                } else {
                    this.treeInstance.$emit('on-cancel-selected', this.node, this.state)
                }
            },
            'state.hidden' (hidden) {
                this.treeInstance.$emit('on-hidden-change', this.node, this.state, hidden)
            }
        },
        created () {
            this.$nextTick(() => {
                this.layout.addFlatternState(this.state)
                this.updateBasicState()
            })
        },
        beforeDestroy () {
            this.layout.destory(this.state)
        },
        methods: {
            updateBasicState () {
                const node = this.node
                const basicState = this.basicState
                const stateHandler = {
                    selected: this.handleNodeClick,
                    expanded: this.handleExpandClick
                }
                const changedState = {
                    id: this.layout.getNodeId(node)
                }
                for (let key in basicState) {
                    if (node.hasOwnProperty(key) && node[key] !== this.state[key]) {
                        changedState[key] = node[key]
                    }
                }
                Object.assign(this.state, changedState)
                for (let key in changedState) {
                    if (stateHandler.hasOwnProperty(key)) {
                        stateHandler[key](changedState[key])
                    }
                }
            },
            handleExpandClick (expanded) {
                this.layout.toggleExpanded(this.state.id, expanded)
            },
            async handleNodeClick () {
                if (typeof this.treeInstance.beforeClick === 'function') {
                    let confirm
                    try {
                        confirm = await Promise.resolve(this.treeInstance.beforeClick(this.node, this.state))
                    } catch (e) {
                        confirm = e
                    }
                    if (!confirm) {
                        this.layout.unselectState(this.state.id)
                        return false
                    }
                }
                if (this.treeInstance.selectable) {
                    const selectedState = this.layout.selectedState
                    if (!selectedState || selectedState.id !== this.state.id) {
                        this.layout.selectState(this.state.id)
                    }
                }
                this.treeInstance.$emit('on-click', this.node, this.state)
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
        .tree-node-expanded-icon{
            display: block;
            margin: 5px 0 0 0;
            font-size: 16px;
            color: #c3cdd7;
            cursor: pointer;
            position: relative;
            z-index: 2;
            &:hover{
                color: #3c96ff;
            }
        }
        &:last-child > .tree-node-children:before {
            display: none;
        }
    }
    .tree-node-info-layout{
        position: relative;
        &.tree-node-info-layout-line:before {
            position: absolute;
            top: -18px;
            left: -23px;
            width: 36px;
            height: 30px;
            content: '';
            border-bottom: 1px dashed #d3d8e7;
            border-left: 1px dashed #d3d8e7;
            z-index: 1;
            pointer-events: none;
        }
        &.tree-node-info-layout-line-short-v:before {
            height: 22px;
            top: -10px;
        }
        &.tree-node-info-layout-line-short-h:before {
            width: 20px;
        }
        .tree-node-info{
            height: 24px;
            padding: 0 0 0 14px;
            line-height: 24px;
            font-size: 14px;
            overflow: hidden;
            cursor: pointer;
            position: relative;
            z-index: 1;
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
        margin: 0 0 0 30px;
        &:not(.tree-node-children-root):before {
            position: absolute;
            top: 15px;
            left: -23px;
            width: 0;
            height: calc(100% - 26px);
            content: '';
            border-left: 1px dashed #d3d8e7;
            z-index: 1;
            pointer-events: none;
        }
    }
</style>