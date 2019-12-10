<template>
    <div class="host-apply-sidebar">
        <div class="tree-wrapper">
            <div class="searchbar">
                <div class="search-select">
                    <search-select-mix></search-select-mix>
                </div>
                <div class="action-menu" v-show="!actionMode">
                    <bk-dropdown-menu>
                        <div :class="['dropdown-trigger', { selected: actionMode }]" slot="dropdown-trigger">
                            <i class="bk-cc-icon icon-cc-menu"></i>
                        </div>
                        <ul class="bk-dropdown-list" slot="dropdown-content">
                            <li><a href="javascript:;" @click="handleBatchAction('batch-edit')">批量编辑</a></li>
                            <li><a href="javascript:;" @click="handleBatchAction('batch-del')">批量删除</a></li>
                        </ul>
                    </bk-dropdown-menu>
                </div>
            </div>
            <topology-tree
                ref="topologyTree"
                :tree-options="treeOptions"
                :action="actionMode"
                @selected="handleTreeSelected"
                @checked="handleTreeChecked"
            ></topology-tree>
        </div>
        <div class="checked-list" v-show="showCheckedPanel">
            <div class="panel-hd">
                <div class="panel-title">
                    已选择<em class="checked-num">{{checkedList.length}}</em>模块
                    <a href="javascript:;" class="clear-all" @click="handleClearChecked">清空</a>
                </div>
            </div>
            <div class="panel-bd">
                <dl class="module-list">
                    <div class="module-item" v-for="item in checkedList" :key="item.bk_inst_id">
                        <dt class="module-name">{{item.bk_inst_name}}</dt>
                        <dd class="module-path" :title="item.path.join(' / ')">{{item.path.join(' / ')}}</dd>
                        <dd class="module-icon"><span>模</span></dd>
                        <dd class="action-icon">
                            <a href="javascript:;" @click="handleRemoveChecked(item.bk_inst_id)">
                                <i class="bk-icon icon-close"></i>
                            </a>
                        </dd>
                    </div>
                </dl>
            </div>
            <div class="panel-ft">
                <bk-button theme="primary" :disabled="!checkedList.length" @click="handleGoEdit">
                    {{$t(actionMode === 'batch-del' ? '下一步' : '去编辑')}}
                </bk-button>
                <bk-button theme="default" @click="handleCancelEdit">取消</bk-button>
            </div>
        </div>
    </div>
</template>

<script>
    import searchSelectMix from './search-select-mix'
    import topologyTree from './topology-tree'
    export default {
        components: {
            searchSelectMix,
            topologyTree
        },
        props: {
        },
        data () {
            return {
                treeOptions: {
                    showCheckbox: false
                },
                actionMode: '',
                showCheckedPanel: false,
                checkedList: []
            }
        },
        computed: {
            topologyTree () {
                return this.$refs.topologyTree
            }
        },
        watch: {
            actionMode (value) {
                this.$emit('action-change', value)
            }
        },
        methods: {
            setApplyClosed (moduleId, isClear) {
                this.topologyTree.updateNodeStatus(moduleId, isClear)
            },
            removeChecked () {
                const tree = this.topologyTree.$refs.tree
                tree.removeChecked({ emitEvent: true })
            },
            async handleBatchAction (actionMode) {
                this.actionMode = actionMode
                this.showCheckedPanel = true
                this.treeOptions.showCheckbox = true
            },
            handleTreeSelected (node) {
                this.$emit('module-selected', node.data)
            },
            handleTreeChecked (ids, target) {
                const treeData = this.topologyTree.treeData
                const modules = []
                const findModuleNode = function (data, parent) {
                    data.forEach(item => {
                        item.path = parent ? [...parent.path, item.bk_inst_name] : [item.bk_inst_name]
                        if (item.bk_obj_id === 'module' && ids.includes(`module_${item.bk_inst_id}`)) {
                            modules.push(item)
                        }
                        if (item.child) {
                            findModuleNode(item.child, item)
                        }
                    })
                }
                findModuleNode(treeData)

                this.checkedList = modules
            },
            handleRemoveChecked (id) {
                const tree = this.topologyTree.$refs.tree
                const checkedIds = this.checkedList.filter(item => item.bk_inst_id !== id).map(item => `module_${item.bk_inst_id}`)
                tree.removeChecked({ emitEvent: true })
                tree.setChecked(checkedIds, { emitEvent: true, beforeCheck: true, checked: true })
            },
            handleClearChecked () {
                this.removeChecked()
            },
            handleGoEdit () {
                const checkedIds = this.checkedList.map(item => item.bk_inst_id)
                this.$router.push({
                    name: 'hostApplyEdit',
                    query: {
                        mid: checkedIds.join(','),
                        batch: 1,
                        action: this.actionMode
                    }
                })
            },
            handleCancelEdit () {
                this.treeOptions.showCheckbox = false
                this.showCheckedPanel = false
                this.actionMode = ''
                this.removeChecked()
            }
        }
    }
</script>

<style lang="scss" scoped>
    .host-apply-sidebar {
        position: relative;
        height: 100%;
        padding: 10px 0;

        .tree-wrapper {
            height: 100%;
        }
    }
    .searchbar {
        display: flex;
        padding: 0 10px;

        .search-select {
            flex: 1;
        }
        .action-menu {
            flex: none;
            margin-left: 12px;

            .dropdown-trigger {
                border: 1px solid #c4c6cc;
                border-radius: 2px;
                width: 32px;
                height: 32px;
                text-align: center;
                line-height: 32px;
                cursor: pointer;
                &:hover {
                    border-color: #979ba5;
                    color: #63656e;
                }
                &:active {
                    border-color: #3a84ff;
                    color: #3a84ff;
                }
                .icon-cc-menu {
                    color: #c4c6cc;
                    font-size: 20px;
                    vertical-align: unset;
                }

                &.selected {
                    border: 1px solid #3a84ff;
                    .icon-cc-menu {
                        color: #3a84ff;
                    }
                }
            }
        }
    }

    .checked-list {
        position: absolute;
        width: 290px;
        height: 100%;
        left: 100%;
        top: 0;
        z-index: 1000;
        border-left: 1px solid #dcdee5;
        border-right: 1px solid #dcdee5;

        .panel-hd,
        .panel-ft {
            background: #fafbfd;
        }
        .panel-hd {
            height: 52px;
            line-height: 52px;
            padding: 0 12px;
            border-bottom: 1px solid #dcdee5;
        }
        .panel-title {
            position: relative;
            font-size: 14px;
            color: #63656e;

            .checked-num {
                font-style: normal;
                font-weight: bold;
                color: #2dcb56;
                padding: .1em;
            }

            .clear-all {
                position: absolute;
                right: 0;
                top: 0;
                color: #3a84ff;
            }
        }
        .panel-bd {
            height: calc(100% - 52px - 60px);
            background: #fff;
            @include scrollbar-y;
        }
        .panel-ft {
            height: 60px;
            line-height: 58px;
            text-align: center;
            border-top: 1px solid #dcdee5;
            .bk-button {
                min-width: 86px;
                margin: 0 3px;
            }
        }

        .module-list {
            .module-item {
                position: relative;
                padding: 8px 42px;

                &:hover {
                    background: #f0f1f5;
                }

                .module-name {
                    font-size: 14px;
                    color: #63656e;
                }
                .module-path {
                    font-size: 12px;
                    color: #c4c6cc;
                    @include ellipsis;
                }
                .module-icon {
                    position: absolute;
                    left: 12px;
                    top: 8px;
                    font-size: 12px;
                    border-radius: 50%;
                    background-color: #c4c6cc;
                    width: 22px;
                    height: 22px;
                    line-height: 21px;
                    text-align: center;
                    font-style: normal;
                    color: #fff;
                }
                .action-icon {
                    position: absolute;
                    right: 12px;
                    top: 12px;
                }
            }
        }
    }
</style>
