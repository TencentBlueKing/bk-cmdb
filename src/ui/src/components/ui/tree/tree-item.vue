<template>
    <div class="tree-node clearfix"
        v-show="node.parent === null || (node.visible && node.parent.expanded)"
        :class="{
            'is-root': node.parent === null,
            'is-leaf': node.isLeaf,
            'is-first-child': node.isFirst,
            'is-last-child': node.isLast,
            'is-expand': node.expanded,
            'is-selected': node.selected,
            'has-link-line': $parent.showLinkLine
        }"
        :style="style"
        @click="$parent.setSelected(node.id, true, true)">
        <div class="node-options fl">
            <i v-if="!node.isLeaf"
                :class="['node-folder-icon', node.expanded ? node.expandIcon : node.collapseIcon]"
                @click.stop="$parent.setExpanded(node.id, !node.expanded, true)">
            </i>
            <input type="checkbox" class="node-checkbox"
                v-if="$parent.showCheckbox"
                :checked="node.checked"
                @click.prevent.stop="$parent.setChecked(node.id, !node.checked, true)">
            <i v-if="node.nodeIcon"
                :class="['node-icon', node.nodeIcon]">
            </i>
        </div>
        <div class="node-content">
            <slot>{{node.name}}</slot>
        </div>
    </div>
</template>

<script>
    export default {
        name: 'cmdb-tree-node',
        props: {
            node: {
                type: Object,
                default () {
                    return {}
                }
            }
        },
        data () {
            return {
                line: 0
            }
        },
        computed: {
            style () {
                return {
                    'margin-left': this.node.level * 30 + 'px',
                    '--level': this.node.level,
                    '--line': this.line
                }
            }
        },
        created () {
            this.node.vNode = this
        },
        methods: {
            calulateLine () {
                const {
                    children,
                    isLeaf,
                    expanded
                } = this.node
                if (isLeaf || !expanded) {
                    this.line = 0
                    return
                }
                const firstChild = children[0]
                const lastChild = children[children.length - 1]
                this.line = lastChild.vNode.$el.getBoundingClientRect().bottom - firstChild.vNode.$el.getBoundingClientRect().top
            }
        }
    }
</script>

<style lang="scss" scoped>
    .tree-node {
        position: relative;
        height: 32px;
        line-height: 32px;
        font-size: 0;
        cursor: pointer;
        @include ellipsis;
        &.has-link-line:not(.is-root):before {
            content: "";
            position: absolute;
            width: 22px;
            height: 0;
            border-bottom: 1px dashed #c3cdd7;
            left: -22px;
            top: 15px;
            z-index: 1;
            pointer-events: none;
        }
        &.has-link-line:after {
            position: absolute;
            top: 16px;
            left: 8px;
            width: 0;
            height: calc(var(--line) * 1px);
            border-left: 1px dashed #c3cdd7;
            content: "";
            pointer-events: none;
            z-index: 1;
        }
        &:hover {
            background-color: #f1f7ff;
        }
        &.is-selected {
            background-color: #e1ecff;
            .node-icon {
                color: #fff;
                background-color: #498fe0;
            }
            .node-content {
                color: #498fe0;
            }
        }
        &.is-leaf {
            padding-left: 16px;
        }
        .node-options {
            height: 100%;
            .node-folder-icon {
                position: relative;
                font-size: 16px;
                z-index: 2;
                @include inlineBlock;
            }
            .node-checkbox {
                margin: 0 6px 0 0;
                @include inlineBlock;
            }
            .node-icon {
                margin: 0 6px;
                font-size: 18px;
                @include inlineBlock;
            }
        }
        .node-content {
            font-size: 14px;
            @include ellipsis;
        }
    }
</style>
