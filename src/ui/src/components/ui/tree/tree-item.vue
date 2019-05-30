<template>
    <div class="tree-node"
        v-show="node.parent === null || (node.visible && node.parent.expanded)"
        :class="{
            'is-root': node.parent === null,
            'is-leaf': node.isLeaf,
            'is-first-child': isFirst,
            'is-last-child': isLast,
            'is-expand': node.expanded,
            'is-selected': node.selected,
            'has-link-line': $parent.showLinkLine
        }"
        :style="style"
        @click="$parent.setSelected(node.id, true)">
        <i v-if="!node.isLeaf"
            :class="['node-folder-icon', node.expanded ? node.expandIcon : node.collapseIcon]"
            @click.stop="toggleExpand">
        </i>
        <input type="checkbox" class="node-checkbox"
            v-if="$parent.showCheckbox"
            :checked="node.checked"
            @click.prevent.stop="$parent.setChecked(node.id, !node.checked, true)">
        <i v-if="node.nodeIcon"
            :class="['node-icon', node.nodeIcon]">
        </i>
        <span class="node-content">
            <slot>{{node.name}}</slot>
        </span>
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
            isLast () {
                return false
            },
            isFirst () {
                return this.childIndex === 0 || this.index === 0
            },
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
            toggleExpand () {
                this.node.expanded = !this.node.expanded
                this.$parent.$emit('toggleExpand', this.node.expanded, this.node)
            },
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
    @mixin middleBlock {
        display: inline-block;
        vertical-align: middle;
    }
    .tree-node {
        position: relative;
        height: 32px;
        line-height: 32px;
        font-size: 0;
        cursor: pointer;
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
        }
        &.is-leaf {
            padding-left: 16px;
        }
        &.is-first-child {
            .node-link-line {
                height: 40px;
                top: -16px;
            }
        }
        &.is-last-child {
            .node-link-line {
                height: 24px;
            }
        }
        .node-link-line {
            position: absolute;
            top: -8px;
            width: 1px;
            height: 32px;
            background-color: #f30;
            z-index: 99;
        }
        .node-folder-icon {
            position: relative;
            margin: 0 6px 0 0;
            font-size: 16px;
            z-index: 2;
            @include middleBlock;
        }
        .node-checkbox {
            @include middleBlock;
        }
        .node-icon {
            font-size: 18px;
            @include middleBlock;
        }
        .node-content {
            font-size: 14px;
            @include middleBlock;
        }
    }
</style>
