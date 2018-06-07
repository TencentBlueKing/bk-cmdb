<template>
    <li :class="['tree-item', {'tree-item-leaf': !isFolder, 'tree-item-open': open, 'tree-item-hide-root': hideRoot && level === 1}]">
        <div :class="['item-node', `item-node-level-${level}`, {'item-node-hide-root': hideRoot && level === 2}]" v-if="!(hideRoot && level === 1)">
            <i class="tree-icon icon fl" :class="[open ? 'icon-cc-rect-sub' : 'icon-cc-rect-add']" v-if="isFolder && !treeData.isLoading" @click="toggleFolder"></i>
            <!-- <i class="tree-icon fl" :class="[open ? 'tree-icon-open' : 'tree-icon-close']" v-if="isFolder && !treeData.isLoading" @click="toggleFolder"></i> -->
            <i class="tree-icon tree-icon-loading fl" v-if="isFolder && treeData.isLoading"></i>
            <div :class="['node-info', {'active': active, 'node-info-not-folder': !isFolder}]" @click="nodeClick">
                <i v-if="treeData['default'] === 1" class="tree-icon fl icon-cc-host-free-pool"></i>
                <i v-else-if="treeData['default'] === 2" class="tree-icon fl icon-cc-host-breakdown"></i>
                <i v-else class="tree-icon tree-icon-name fl">{{treeData['default'] ? treeData['bk_inst_name'][0]: treeData['bk_obj_name'][0]}}</i>
                <template v-if="rootNode['bk_inst_name'] === '蓝鲸'">
                    <i v-tooltip="{content: $t('Common[\'关键业务不能够修改\']'), classes: 'topo-tip'}" class="tree-icon tree-icon-add fr disabled" v-if="isShowAdd && active">{{$t('Common[\'新增\']')}}</i>
                </template>
                <template v-else>
                    <i class="tree-icon tree-icon-add fr" :class="{'disabled': rootNode['bk_inst_name'] === '蓝鲸'}" v-if="isShowAdd && active" @click.stop="addNode" :title="`${$t('Common[\'新增\']')}-${optionModel['bk_next_name']}`">{{$t('Common[\'新增\']')}}</i>
                </template>
                <div class="node-name" :title="treeData['bk_inst_name']">{{treeData['bk_inst_name']}}</div>
            </div>  
        </div>
        <ul v-show="open" v-if="isFolder" class="item-list">
            <tree-item v-for="(childTreeData, index) in treeData.child"
                :key="index"
                :nodeIndex="index"
                :parent="treeData"
                :treeData="childTreeData"
                :initNode="initNode"
                :rootNode="rootNode"
                :treeId="treeId"
                :activeNodeId="activeNodeId"
                :model="model"
                :hideRoot="hideRoot"
                :level="nextLevel">
          </tree-item>
        </ul>
  </li>
</template>
<script type="text/javascript">
    import bus from '@/eventbus/bus'
    export default {
        name: 'treeItem',
        props: {
            treeData: Object,
            level: Number,
            parent: Object,
            initNode: Object,
            rootNode: Object,
            nodeIndex: Number,
            model: Array,
            treeId: {
                required: true
            },
            activeNodeId: {
                required: true
            },
            hideRoot: Boolean
        },
        data () {
            return {
                open: false,
                active: false
            }
        },
        computed: {
            isFolder () {
                if (this.treeData.hasOwnProperty('isFolder')) {
                    return this.treeData.isFolder
                } else if (this.treeData['bk_obj_id'] === 'module' || (Array.isArray(this.treeData.child) && !this.treeData.child.length)) {
                    return false
                }
                return true
            },
            nodeId () {
                return Math.random()
            },
            nextLevel () {
                return this.level + 1
            },
            isShowAdd () {
                return this.model.length && this.treeData['bk_obj_id'] !== 'module' && !this.treeData['default']
            },
            optionModel () {
                return this.model.find(model => {
                    return model['bk_obj_id'] === this.treeData['bk_obj_id']
                })
            },
            nodeOptions () {
                return {
                    parent: this.parent,
                    rootNode: this.rootNode,
                    level: this.level,
                    nodeId: this.nodeId,
                    treeId: this.treeId,
                    isOpen: this.open,
                    index: this.nodeIndex ? this.nodeIndex : 0
                }
            }
        },
        watch: {
            initNode (initNode) {
                if (initNode && Object.keys(initNode).length) {
                    this.init()
                }
            },
            activeNodeId (activeNodeId) {
                this.active = activeNodeId === this.nodeId
            },
            open (open) {
                this.treeData.isOpen = open
                bus.$emit('nodeToggle', this.open, this.treeData, this.nodeOptions)
            }
        },
        methods: {
            toggleFolder () {
                this.open = !this.open
            },
            nodeClick () {
                if (!this.active) {
                    this.active = true
                    bus.$emit('nodeClick', this.treeData, this.nodeOptions)
                }
            },
            addNode () {
                bus.$emit('addNode', this.treeData, this.nodeOptions)
            },
            init () {
                let {
                    level,
                    open,
                    active,
                    bk_inst_id: bkInstId
                } = this.initNode
                if (this.level === level && this.treeData['bk_inst_id'] === bkInstId) {
                    this.open = open
                    this.active = active
                    if (active) {
                        bus.$emit('nodeClick', this.treeData, this.nodeOptions)
                    }
                }
            }
        },
        created () {
            this.$set(this.treeData, 'isOpen', this.open)
            this.$set(this.treeData, 'level', this.level)
            this.$set(this.treeData, 'index', this.nodeIndex ? this.nodeIndex : 0)
            this.$watch('treeData.isOpen', (isOpen) => {
                this.open = isOpen
            })
        },
        mounted () {
            if (this.initNode && Object.keys(this.initNode).length) {
                this.init()
            }
        }
    }
</script>
<style lang="scss" scoped>
    .tree-item{
        position: relative;
        white-space: nowrap;
        &.tree-item-open:not(.tree-item-leaf):not(.tree-item-hide-root):before{
            content: '';
            width: 0;
            height: calc(100% - 15px);
            position: absolute;
            left: 8px;
            top: 4px;
            border-left: 1px dashed #d3d8e7;
            z-index: 1;
        }
    }
    .item-node{
        font-size: 0;
        height: 24px;
        line-height: 24px;
        position: relative;
        margin: 8px 0;
        &:not(.item-node-level-1):not(.item-node-hide-root):before{
            content: '';
            width: 20px;
            height: 0;
            border-top: 1px dashed #d3d8e7;
            position: absolute;
            top: 12px;
            left: -15px;
            z-index: 1;
        }
        .tree-icon{
            font-size: 12px;
            font-style: normal;
            width: 16px;
            height: 16px;
            line-height: 16px;
            text-align: center;
            margin: 4px 4px 4px 0;
            position: relative;
            z-index: 2;
            &.icon-cc-host-free-pool,
            &.icon-cc-host-breakdown{
                font-size: 16px;
            }
            &.icon{
                font-size: 14px;
                background: #fff;
                color: #c3cdd7;
                cursor: pointer;
                &:hover{
                    color: #3c96ff;
                }
            }
            &.tree-icon-loading {
                background: url(../../common/images/icon/loading_2_16x16.gif) no-repeat 0 0 #fff;
            }
            &.tree-icon-name{
                line-height: 14px;
                border-radius: 5px;
                background: #c3cdd7;
                color: #fff;
                margin: 4px 0;
            }
            &.tree-icon-add{
                width: auto;
                height: 18px;
                line-height: 18px;
                margin: 3px 4px;
                padding: 0 6px;
                background: #3c96ff;
                color:#fff;
                border-radius: 4px;
                &.disabled{
                    background: #fafafa;
                    color: #ccc;
                    cursor: default;
                }
            }
        }
        .node-info{
            padding: 0 0 0 10px;
            overflow: hidden;
            cursor: pointer;
            &:hover{
                background: #f1f7ff;
                color: #498fe0;
                .tree-icon-name{
                    background: #50abff;
                }
            }
            &.active{
                background: #e2efff;
                color: #498fe0;
                .tree-icon-name{
                    background: #498fe0;
                }
                .icon-cc-host-free-pool,
                .icon-cc-host-breakdown{
                    color: #ffb400;
                }
            }
            &.node-info-not-folder {
                margin-left: 20px;
            }
            .node-name{
                font-size: 12px;
                padding: 0 0 0 6px;
                @include ellipsis;
            }
        }
    }
    .item-list{
        margin: 0 0 0 24px;
    }
</style>