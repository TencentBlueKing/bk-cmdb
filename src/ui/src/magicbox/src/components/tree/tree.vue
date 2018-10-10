<template>
    <ul class="bk-tree" :class="{'bk-has-border-tree': isBorder}">
        <li v-for="(item, index) in data"
            @drop="drop(item, $event)"
            @dragover="dragover($event)"
            :key="item[nodeKey] ? item[nodeKey] : item.name"
            :class="{leaf: isLeaf(item), 
            'tree-first-node': !parent && index === 0, 
            'tree-only-node': !parent && data.length === 1, 
            'tree-second-node': !parent && index === 1,
            'single': !multiple}"
            v-show="item.hasOwnProperty('visible') ? item.visible : true">
            <div :class="['tree-drag-node', !multiple ? 'tree-singe' : '']"
                 :draggable="draggable" @dragstart="drag(item, $event)">
                <span 
                    @click="expandNode(item)" 
                    v-if="!item.parent ||item.children && item.children.length || item.async" 
                    :class="['bk-icon', 'tree-expanded-icon', item.expanded ? 'icon-down-shape' : 'icon-right-shape']">
                </span>
                <label :class="[item.halfcheck ? 'bk-form-half-checked' : 'bk-form-checkbox','bk-checkbox-small', 'mr5',!item.parent ||item.children && item.children.length ? 'parent-node-left' : '']" v-if="multiple && !item.nocheck">
                    <input type="checkbox"
                        v-if='multiple'
                        :disabled="item.disabled"
                        v-model="item.checked"
                        @change.stop="changeCheckStatus(item, $event)">
                </label>
                <div class="tree-node" @click="triggerExpand(item)">
                    <span class="node-icon bk-icon" v-if="item.icon || item.openedIcon || item.closedIcon" :class="setNodeIcon(item)"></span>
                    <div class="bk-spin-loading bk-spin-loading-mini bk-spin-loading-primary loading" v-if="item.loading && item.expanded">
                        <div class="rotate rotate1"></div>
                        <div class="rotate rotate2"></div>
                        <div class="rotate rotate3"></div>
                        <div class="rotate rotate4"></div>
                        <div class="rotate rotate5"></div>
                        <div class="rotate rotate6"></div>
                        <div class="rotate rotate7"></div>
                        <div class="rotate rotate8"></div>
                    </div>
                    <Render :node="item" :tpl ='tpl'/>
                </div>
            </div>
            <collapse-transition>
                <bk-tree v-if="!isLeaf(item)"
                    @dropTreeChecked='nodeCheckStatusChange'
                    @async-load-nodes='asyncLoadNodes'
                    @on-expanded='onExpanded'
                    @on-click='onClick'
                    @on-check='onCheck'
                    @on-drag-node='onDragNode'
                    :dragAfterExpanded="dragAfterExpanded"
                    :draggable="draggable"
                    v-show="item.expanded"
                    :tpl ="tpl"
                    :data="item.children"
                    :halfcheck='halfcheck'
                    :parent ='item'
                    :isDeleteRoot ='isDeleteRoot'
                    :multiple="multiple">
                </bk-tree>
            </collapse-transition>
        </li>
    </ul>
</template>
<script>
import Render from './render'
import CollapseTransition from './collapse-transition'
export default {
    name: 'bk-tree',
    props: {
        data: {
            type: Array,
            default: () => []
        },
        parent: {
            type: Object,
            default: () => null
        },
        multiple: {
            type: Boolean,
            default: false
        },
        nodeKey: {
            type: String,
            default: 'id'
        },
        draggable: {
            type: Boolean,
            default: false
        },
        hasBorder: {
            type: Boolean,
            default: false
        },
        dragAfterExpanded: {
            type: Boolean,
            default: true
        },
        isDeleteRoot: {
            type: Boolean,
            default: false
        },
        tpl: Function
    },
    components: {Render, CollapseTransition},
    data () {
        return {
            halfcheck: true,
            isBorder: this.hasBorder,
            bkTreeDrag: {}
        }
    },
    watch: {
        data () {
            this.initTreeData()
        },
        hasBorder (value) {
            this.isBorder = !!value
        }
    },
    mounted () {
        /**
        * @event monitor 子节点 selected event
        */
        this.$on('childChecked', (node, checked) => {
            if (node.children && node.children.length) {
                for (let child of node.children) {
                    if (!child.disabled) {
                        this.$set(child, 'checked', checked)
                    }
                    this.$emit('on-check', child, checked)
                }
            }
        })

        /**
        * @event monitor 父节点 selected event
        */
        this.$on('parentChecked', (node, checked) => {
            if (!node.disabled) {
                this.$set(node, 'checked', checked)
            }
            if (!node.parent) return false
            let someBortherNodeChecked = node.parent.children.some(node => node.checked)
            let allBortherNodeChecked = node.parent.children.every(node => node.checked)
            if (this.halfcheck) {
                allBortherNodeChecked ? this.$set(node.parent, 'halfcheck', false) : someBortherNodeChecked ? this.$set(node.parent, 'halfcheck', true) : this.$set(node.parent, 'halfcheck', false)
                if (!checked && someBortherNodeChecked) {
                    this.$set(node.parent, 'halfcheck', true)
                    return false
                }
                this.$emit('parentChecked', node.parent, checked)
            }
            else {
                if (checked && allBortherNodeChecked) this.$emit('parentChecked', node.parent, checked)
                if (!checked) this.$emit('parentChecked', node.parent, checked)
            }
        })

        /**
        * @event monitor 节点 selected event
        */
        this.$on('on-check', (node, checked) => {
            this.$emit('parentChecked', node, checked)
            this.$emit('childChecked', node, checked)
            this.$emit('dropTreeChecked', node, checked)
        })

        /**
        * @event monitor 节点过滤时 可见/不可见 visible event
        */
        this.$on('toggleshow', (node, isShow) => {
            this.$set(node, 'visible', isShow)
            if (isShow && node.parent) {
                this.$emit('toggleshow', node.parent, isShow)
            }
        })
        this.$on('cancelSelected', (root) => {
            for (let child of root.$children) {
                for (let node of child.data) {
                    child.$set(node, 'selected', false)
                }
                if (child.$children) child.$emit('cancelSelected', child)
            }
        })
        this.initTreeData()
    },
    destroyed () {
        this.$delete(window, 'bkTreeDrag')
    },
    methods: {
        /**
         * 拖拽时生成随机guid将节点暂存window['bkTreeDrag'][guid]上
         */
        gid () {
            return 'xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx'.replace(/[xy]/g, c => {
                let r = Math.random() * 16 | 0
                let v = c === 'x' ? r : (r & 0x3 | 0x8)
                return v.toString(16)
            })
        },

        /**
         * 设置拖拽节点
         */
        setDragNode (id, node) {
            window['bkTreeDrag'] = {}
            window['bkTreeDrag'][id] = node
        },

        /**
         * 获取拖拽节点
         */
        getDragNode (id) {
            return window['bkTreeDrag'][id]
        },

        /**
         * 节点是否可拖拽到某个节点下的标识
         */
        hasInGenerations (root, node) {
            if (root.hasOwnProperty('children') && root.children) {
                for (let rn of root.children) {
                    if (rn === node) return true
                    if (rn.children) return this.hasInGenerations(rn, node)
                }
                return false
            }
        },

        /**
         * 设置节点icon(父节点和子节点的传参不同)
         *
         * @param {Object} node 当前节点
         */
        setNodeIcon (node) {
            if (node.children && node.children.length) {
                if (node.expanded) {
                    return node.openedIcon
                } else {
                    return node.closedIcon
                }
            } else {
                return node.icon
            }
        },
        
        /**
        * 节点拖拽
        *
        * @param {Object} node 当前拖拽节点
        * @param {Object} ev   $event
        */
        drop (node, ev) {
            ev.preventDefault()
            ev.stopPropagation()
            let gid = ev.dataTransfer.getData('gid')
            console.log(gid)
            let drag = this.getDragNode(gid)
            // if drag node's parent is enter node or root node
            if (drag.parent === node || drag.parent === null || drag === node) return false
            // drag from parent node to child node
            if (this.hasInGenerations(drag, node)) return false
            let dragHost = drag.parent.children
            if (node.children && node.children.indexOf(drag) === -1) {
                node.children.push(drag)
                dragHost.splice(dragHost.indexOf(drag), 1)
            } else {
                this.$set(node, 'openedIcon', 'icon-folder-open')
                this.$set(node, 'closedIcon', 'icon-folder')
                this.$set(node, 'children', [drag])
                dragHost.splice(dragHost.indexOf(drag), 1)
            }
            this.$set(node, 'expanded', this.dragAfterExpanded)
            this.$emit('on-drag-node', {dragNode: drag, targetNode: node})
        },
        drag (node, ev) {
            let gid = this.gid()
            this.setDragNode(gid, node)
            ev.dataTransfer.setData('gid', gid)
        },
        dragover (ev) {
            ev.preventDefault()
            ev.stopPropagation()
        },

        /*
        * 数据初始化
        */
        initTreeData () {
            for (let node of this.data) {
                this.$set(node, 'parent', this.parent)
                if (node.children && node.children.length) {
                    if (node.hasOwnProperty('disabled')) {
                        this.$delete(node, 'disabled')
                    }
                    if (node.hasOwnProperty('icon')) {
                        this.$delete(node, 'icon')
                    }
                } else {
                    if (node.hasOwnProperty('openedIcon')) {
                        this.$delete(node, 'openedIcon')
                    }
                    if (node.hasOwnProperty('closedIcon')) {
                        this.$delete(node, 'closedIcon')
                    }
                }
                if (this.multiple) {
                    if (node.hasOwnProperty('selected')) {
                        this.$delete(node, 'selected')
                    }
                } else {
                    if (node.hasOwnProperty('checked')) {
                        this.$delete(node, 'checked')
                    }
                }
            }
        },

        /**
        * 节点展开/收起
        *
        * @param {Object} node 当前节点
        */
        expandNode (node) {
            this.$set(node, 'expanded', !node.expanded)
            if (node.async && !node.children) {
                this.$emit('async-load-nodes', node)
            }
            if (node.children && node.children.length) {
                this.$emit('on-expanded', node, node.expanded)
            }
        },
        onExpanded (node) {
            if (node.children && node.children.length) {
                this.$emit('on-expanded', node, node.expanded)
            }
        },
        triggerExpand (item) {
            if (!item.parent || (item.children && item.children.length) || item.async) {
                this.expandNode(item)
            }
        },
 
        /**
        * 异步加载节点
        *
        * @param {Object} node 当前点击节点
        */
        asyncLoadNodes (node) {
            if (node.async && !node.children) {
                this.$emit('async-load-nodes', node)
            }
        },

        /**
        * 判断是否子节点
        *
        * @param {Object} node 当前节点
        */
        isLeaf (node) {
            return !(node.children && node.children.length) && node.parent
        },

        /**
        * 添加单个节点
        *
        * @param {Object} parent 父节点
        * @param {Object} newnode  新节点
        */
        addNode (parent, newNode) {
            let addnode = {}
            this.$set(parent, 'expanded', true)
            if (typeof newNode === 'undefined') {
                throw new ReferenceError('newNode is required but undefined')
            }
            if (typeof newNode === 'object' && !newNode.hasOwnProperty('name')) {
                throw new ReferenceError('the name property is missed')
            }
            if (typeof newNode === 'object' && !newNode.hasOwnProperty(this.nodeKey)) {
                throw new ReferenceError('the nodeKey property is missed')
            }
            if (typeof newNode === 'object' && newNode.hasOwnProperty('name') && newNode.hasOwnProperty(this.nodeKey)) {
                addnode = Object.assign({}, newNode)
            }
            if (this.isLeaf(parent)) {
                this.$set(parent, 'children', [])
                parent.children.push(addnode)
            } else {
                parent.children.push(addnode)
            }
            this.$emit('addNode', { parentNode: parent, newNode: newNode })
        },

        /**
        * 添加多个节点
        *
        * @param {Object} parent 父节点
        * @param {Array} newChildren  子节点数组
        */
        addNodes (parent, newChildren) {
            for (let n of newChildren) {
                this.addNode(parent, n)
            }
        },

        /**
        * 节点点击事件
        *
        * @param {Object} node 当前点击节点
        */
        onClick (node) {
            this.$emit('on-click', node)
        },

        /**
        * 节点复选框 change 事件
        *
        * @param {Object} node 当前节点
        */
        onCheck (node, checked) {
            this.$emit('on-check', node, checked)
        },

        /**
        * 复选框状态改变时告知父组件
        *
        * @param {Object} node 改变状态的节点
        * @param {Boolean} checked 选中/非选中
        */
        nodeCheckStatusChange (node, checked) {
            this.$emit('dropTreeChecked', node, checked)
        },
        
        /**
        * 拖拽结束事件
        *
        * @param {Object} event $event
        */
        onDragNode (event) {
            this.$emit('on-drag-node', event)
        },

        /**
        * 删除节点
        *
        * @param {Object} parent 父节点
        * @param {Object} node 当前节点
        */
        delNode (parent, node) {
            if (parent === null || typeof parent === 'undefined') {
                // isDeleteRoot 为false时不可删除根节点
                if (this.isDeleteRoot) {
                    this.data.splice(0, 1)
                } else {
                    throw new ReferenceError('the root element can\'t deleted!')
                }
            } else {
                parent.children.splice(parent.children.indexOf(node), 1)
            }
            this.$emit('delNode', { parentNode: parent, delNode: node })
        },

        /**
        * 复选框change事件
        *
        * @param {Object} node 当前节点
        * @param {Object} event $event
        */
        changeCheckStatus (node, $event) {
            this.$emit('on-check', node, $event.target.checked)
        },

        /**
        * 节点selected
        *
        * @param {Object} node 当前节点
        */
        nodeSelected (node) {
            const getRoot = (el) => {
                if (el.$parent.$el.nodeName === 'UL') {
                    el = el.$parent
                    return getRoot(el)
                } return el
            }
            let root = getRoot(this)
            if (!this.multiple) {
                for (let rn of root.data || []) {
                    this.$set(rn, 'selected', false)
                    this.$emit('cancelSelected', root)
                }
            }
            // 当为多选时 必须通过选择复选框触发
            // if (this.multiple) this.$set(node, 'checked', !node.selected)
            this.$set(node, 'selected', !node.selected)
            this.$emit('on-click', node)
        },

        /**
        * 节点数据处理
        *
        * @param {Object} opt 参数设置
        * @param {Array}  data 根节点或子节点数组
        * @param {Array/String}  keyParton 自定义键值
        */
        nodeDataHandler (opt, data, keyParton) {
            data = data || this.data
            let res = []
            let keyValue = keyParton
            for (const node of data) {
                for (const [key, value] of Object.entries(opt)) {
                    if (node[key] === value) {
                        if (!keyValue.length || !keyValue) {
                            let n = Object.assign({}, node)
                            delete n['parent']
                            if (!(n.children && n.children.length)) {
                                res.push(n)
                            }
                        } else {
                            let n = {}
                            if (Object.prototype.toString.call(keyValue) === '[object Array]') {
                                for (let i = 0; i < keyValue.length; i++) {
                                    if (node.hasOwnProperty(keyValue[i])) {
                                        n[keyValue[i]] = node[keyValue[i]]
                                    }
                                }
                            }
                            if (Object.prototype.toString.call(keyValue) === '[object String]') {
                                n[keyValue] = node[keyValue]
                            }
                            if (!(node.children && node.children.length)) {
                                res.push(n)
                            }
                        }
                    }
                }
                if (node.children && node.children.length) {
                    res = res.concat(this.nodeDataHandler(opt, node.children, keyValue))
                }
            }
            return res
        },

        /**
        * 获取所需节点
        *
        * @param {Array/String}  keyParton 自定义键值
        */
        getNode (keyParton) {
            if (!this.multiple) {
                return this.nodeDataHandler({selected: true}, this.data, keyParton)
            } else {
                return this.nodeDataHandler({checked: true}, this.data, keyParton)
            }
        },

        /**
        * 节点过滤
        *
        * @param {String/Function} filter 过滤器
        * @param {Object} data 所需过滤的数据
        */
        searchNode (filter, data) {
            data = data || this.data
            for (const node of data) {
                let searched = filter ? (typeof filter === 'function' ? filter(node) : node['name'].indexOf(filter) > -1) : false
                this.$set(node, 'searched', searched)
                this.$set(node, 'visible', false)
                this.$emit('toggleshow', node, filter ? searched : true)
                if (node.children && node.children.length) {
                    if (searched) this.$set(node, 'expanded', true)
                    this.searchNode(filter, node.children)
                }
            }
        }
    }
}
</script>
<style lang="scss">
    @import '../../bk-magic-ui/src/tree.scss';
</style>
