<template>
    <div class="tree-node"
        v-if="showNode"
        :class="{
            'is-root': !parent,
            'is-leaf': isLeaf,
            'is-last-child': isLast,
            'is-first-child': isFirst,
            'is-expand': node.expanded,
            'has-link-line': tree.showLinkLine,
            'has-sibling': parent && parent.children.length > 1
        }"
        :style="style">
        <i v-if="!isLeaf"
            :class="['node-folder-icon', folderIcon]"
            @click="toggleExpand">
        </i>
        <input type="checkbox" class="node-checkbox" v-if="tree.showCheckbox">
        <i v-if="node.icon.node"
            :class="['node-icon', node.icon.node]">
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
                line: 0,
                timer: null
            }
        },
        computed: {
            tree () {
                return this.$parent
            },
            children () {
                return this.tree.nodes.filter(node => node.parent === this.node)
            },
            index () {
                return this.tree.nodes.indexOf(this.node)
            },
            parent () {
                return this.node.parent
            },
            showNode () {
                return !this.parent || (this.node.visible && this.parent.expanded)
            },
            isLeaf () {
                return !this.children.length
            },
            isLast () {
                return this.childIndex === this.parentChildren.length - 1
            },
            isFirst () {
                return this.childIndex === 0
            },
            parentChildren () {
                const parentChildren = this.parent ? this.parent.vNode.children : []
                return parentChildren
            },
            childIndex () {
                return this.parentChildren.indexOf(this.node)
            },
            style () {
                const paddingLeft = this.node.level * 30
                const folderIconWidth = 16
                return {
                    'padding-left': (this.isLeaf ? paddingLeft + folderIconWidth : paddingLeft) + 'px',
                    '--level': this.node.level,
                    '--line': this.line
                }
            },
            folderIcon () {
                return this.node.expanded ? this.node.icon.expand : this.node.icon.collapse
            },
            parentVisible () {
                return !this.parent || (this.parent.expanded && this.parent.visible)
            }
        },
        watch: {
            parentVisible (visible) {
                this.node.setState('visible', visible)
            },
            children () {
                this.handleDescendantsChange()
            }
        },
        created () {
            this.node.vNode = this
            this.checkState()
        },
        methods: {
            toggleSelect (selected) {},
            toggleExpand () {
                this.node.expanded = !this.node.expanded
                this.handleDescendantsChange()
                this.parent && this.parent.vNode.handleDescendantsChange()
            },
            checkState () {
                if (this.parent && this.node.expanded) {
                    this.parent.expanded = true
                    this.parent.vNode.checkState()
                    this.parent.vNode.handleDescendantsChange()
                    setTimeout(() => {
                        this.calcLine()
                    }, 0)
                }
            },
            calcLine () {
                const children = this.children
                if (this.isLeaf || !this.node.expanded) {
                    this.line = 0
                    return
                }
                const firstChild = children[0]
                const lastChild = children[children.length - 1]
                this.line = lastChild.vNode.$el.getBoundingClientRect().bottom - firstChild.vNode.$el.getBoundingClientRect().top
                this.parent && this.parent.vNode.handleDescendantsChange()
            },
            handleDescendantsChange () {
                this.$nextTick(() => {
                    this.calcLine()
                })
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
        &:not(.is-root):before {
            content: "";
            position: absolute;
            width: 22px;
            height: 0;
            border-bottom: 1px dashed #c3cdd7;
            left: calc((var(--level) - 1) * 30px + 8px);
            top: 15px;
            z-index: 100;
            pointer-events: none;
        }
        &:after {
            position: absolute;
            top: 16px;
            left: calc(var(--level) * 30px + 8px);
            width: 0;
            height: calc(var(--line) * 1px);
            border-left: 1px dashed #c3cdd7;
            content: "";
            pointer-events: none;
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
            font-size: 16px;
            background-color: #fff;
            cursor: pointer;
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
