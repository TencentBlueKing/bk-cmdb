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
            'is-disabled': node.disabled,
            'has-link-line': $parent.showLinkLine
        }"
        :style="style"
        @click="handleNodeClick">
        <div class="node-options fl">
            <i v-if="!node.isLeaf"
                :class="['node-folder-icon', node.expanded ? node.expandIcon : node.collapseIcon]"
                @click.stop="handleNodeExpand">
            </i>
            <span class="node-checkbox"
                v-if="node.hasCheckbox"
                :class="{
                    'is-disabled': node.disabled,
                    'is-checked': node.checked,
                    'is-indeterminate': node.indeterminate
                }"
                @click.stop="handleNodeCheck">
            </span>
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
            handleNodeCheck () {
                if (this.node.disabled) {
                    return false
                }
                this.$parent.setChecked(this.node.id, {
                    checked: !this.node.checked,
                    emitEvent: true,
                    beforeCheck: true
                })
            },
            handleNodeExpand () {
                this.$parent.setExpanded(this.node.id, {
                    expanded: !this.node.expanded,
                    emitEvent: true
                })
            },
            handleNodeClick () {
                if (this.node.disabled) {
                    return false
                }
                this.$parent.$emit('node-click', this.node)
                this.$parent.setSelected(this.node.id, {
                    emitEvent: true,
                    beforeSelect: true
                })
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
                background-color: #3a84ff;
            }
            .node-content {
                color: #3a84ff;
            }
            .node-folder-icon {
                color: #3a84ff !important;
            }
        }
        &.is-leaf {
            padding-left: 16px;
        }
        &.is-disabled {
            cursor: not-allowed;
        }
        .node-options {
            height: 100%;
            .node-folder-icon {
                @include inlineBlock;
                position: relative;
                font-size: 14px;
                color: #c4c6cc;
                z-index: 2;
            }
            .node-checkbox {
                @include inlineBlock;
                position: relative;
                width: 16px;
                height: 16px;
                margin: 0 6px 0 0;
                border: 1px solid #979ba5;
                border-radius: 2px;
                &.is-checked {
                    border-color: #3a84ff;
                    background-color: #3a84ff;
                    background-clip: border-box;
                    &:after {
                        content: "";
                        position: absolute;
                        top: 1px;
                        left: 4px;
                        width: 4px;
                        height: 8px;
                        border: 2px solid #fff;
                        border-left: 0;
                        border-top: 0;
                        transform-origin: center;
                        transform: rotate(45deg) scaleY(1);
                    }
                    &.is-disabled {
                        background-color: #dcdee5;
                    }
                }
                &.is-disabled {
                    border-color: #dcdee5;
                    cursor: not-allowed;
                }
            }
            .node-icon {
                @include inlineBlock;
                margin: 0 6px;
                font-size: 18px;
            }
        }
        .node-content {
            @include ellipsis;
            font-size: 14px;
        }
    }
</style>
