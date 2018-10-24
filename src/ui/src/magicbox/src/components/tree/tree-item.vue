<template>
    <div>
        <li> 
            <span :title="node.title"  @click="toggle(node)" :class="spanClass(node)">
                <i v-if="node.children && node.children.length > 0" class="bk-icon" :class='{"icon-folder": !node.isOpen, "icon-folder-open": node.isOpen}'></i>
                <i v-else class="bk-icon" :class="iconClass(node)" ></i>
                <span>{{node.name}}</span>
            </span>
            <a v-for="(btn, index) in node.buttons" :key="index" href="javascript:" @click="btnClick(btn, node)" v-if="node.buttons && node.buttons.length">
                <i class="bk-icon" :class="buttonIcon(btn)"></i>
            </a>
            <ul v-show="node.isOpen && (node.children && node.children.length)">
                <tree-item v-for="(item, index) in node.children" :node.sync="item" :key="index" :is-multiple="isMultiple"></tree-item>
            </ul>
        </li>
    </div>
</template>
<script>
    import bus from './bus'
    export default {
        name: 'treeItem',
        props: {
            node: {
                type: Object
            },
            isMultiple: {
                type: Boolean,
                default: false
            }
        },
        methods: {
            init () {
                if ((!this.node.hasOwnProperty('isSelected') && !this.node.children) || (this.node.children && !this.node.children.length)) {
                    this.$set(this.node, 'isSelected', false)
                }
            },
            iconClass (node) {
                if (node.isOpen) {
                    return node.openedIcon || node.icon || 'icon-file'
                } else {
                    return node.closedIcon || node.icon || 'icon-file'
                }
            },
            spanClass (node) {
                if (node.isSelected) {
                    return 'selected'
                } else {
                    return 'normal'
                }
            },
            buttonIcon (btn) {
                if (btn.type === 'add') {
                    return btn.icon || 'icon-plus'
                }
                if (btn.type === 'delete') {
                    return btn.icon || 'icon-delete'
                }
            },
            toggle (node) {
                if (node.hasOwnProperty('isOpen') && node.children && node.children.length > 0) {
                    node.isOpen = !node.isOpen
                }
                if (node.hasOwnProperty('isSelected')) {
                    if (this.isMultiple) {
                        node.isSelected = !node.isSelected
                    }
                    bus.$emit('select', node)
                }
            },
            btnClick (btn, node) {
                if (btn.type === 'add') {
                    if (!btn.data || (btn.data && !btn.data.length)) {
                        return
                    }
                    if (node.hasOwnProperty('isSelected')) {
                        this.$delete(node, 'isSelected')
                    }
                    if (!node.hasOwnProperty('children')) {
                        this.$set(node, 'children', [])
                    }
                    if (!node.hasOwnProperty('isOpen')) {
                        this.$set(node, 'isOpen', true)
                    } else {
                        node.isOpen = true
                    }
                    btn.data.forEach(item => {
                        node.children.push(item)
                    })
                }
                if (btn.type === 'delete') {
                    bus.$emit('delete', node)
                }
            }
        },
        created () {
            this.init()
        }
    }
</script>
