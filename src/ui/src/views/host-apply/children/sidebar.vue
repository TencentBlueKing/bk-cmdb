<template>
    <div class="host-apply-sidebar">
        <div class="tree-wrapper">
            <div class="searchbar">
                <div class="search-select">
                    <search-select-mix></search-select-mix>
                </div>
                <div class="action-menu" v-show="!actionMode">
                    <bk-dropdown-menu
                        @show="showBatchDropdown = true"
                        @hide="showBatchDropdown = false"
                        font-size="medium">
                        <div :class="['dropdown-trigger', { selected: actionMode }]" slot="dropdown-trigger">
                            <span>{{$t('批量操作')}}</span>
                            <i :class="['bk-icon icon-angle-down', { 'icon-flip': showBatchDropdown }]"></i>
                        </div>
                        <ul class="bk-dropdown-list" slot="dropdown-content">
                            <li>
                                <cmdb-auth :auth="{ type: $OPERATION.U_HOST_APPLY, relation: [bizId] }">
                                    <a
                                        href="javascript:;"
                                        slot-scope="{ disabled }"
                                        :class="{ disabled }"
                                        @click="handleBatchAction('batch-edit')"
                                    >
                                        {{$t('批量编辑')}}
                                    </a>
                                </cmdb-auth>
                            </li>
                            <li>
                                <cmdb-auth :auth="{ type: $OPERATION.U_HOST_APPLY, relation: [bizId] }">
                                    <a
                                        href="javascript:;"
                                        slot-scope="{ disabled }"
                                        :class="{ disabled }"
                                        @click="handleBatchAction('batch-del')"
                                    >
                                        {{$t('批量删除')}}
                                    </a>
                                </cmdb-auth>
                            </li>
                        </ul>
                    </bk-dropdown-menu>
                </div>
            </div>
            <topology-tree
                ref="topologyTree"
                :tree-options="treeOptions"
                :action="actionMode"
                :checked="checkedNodes"
                @selected="handleTreeSelected"
                @checked="handleTreeChecked"
            ></topology-tree>
        </div>
        <div class="checked-list" v-show="showCheckedPanel">
            <div class="panel-hd">
                <div class="panel-title">
                    <i18n path="已选择N个模块">
                        <em class="checked-num" place="count">{{checkedList.length}}</em>
                    </i18n>
                    <a href="javascript:;" class="clear-all" @click="handleClearChecked">{{$t('清空')}}</a>
                </div>
            </div>
            <div class="panel-bd">
                <dl class="module-list">
                    <div class="module-item" v-for="item in checkedList" :key="item.bk_inst_id">
                        <dt class="module-name">{{item.bk_inst_name}}</dt>
                        <dd class="module-path" :title="item.path.join(' / ')">{{item.path.join(' / ')}}</dd>
                        <dd class="module-icon"><span>{{$i18n.locale === 'en' ? 'M' : '模'}}</span></dd>
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
                    {{$t(actionMode === 'batch-del' ? '去删除' : '去编辑')}}
                </bk-button>
                <bk-button theme="default" @click="handleCancelEdit">{{$t('取消')}}</bk-button>
            </div>
        </div>
    </div>
</template>

<script>
    import searchSelectMix from './search-select-mix'
    import topologyTree from './topology-tree'
    import { MENU_BUSINESS_HOST_APPLY_EDIT } from '@/dictionary/menu-symbol'
    import { mapGetters } from 'vuex'
    export default {
        components: {
            searchSelectMix,
            topologyTree
        },
        data () {
            return {
                treeOptions: {
                    showCheckbox: false,
                    selectable: true,
                    checkOnClick: false,
                    checkOnlyAvailableStrictly: false,
                    displayMatchedNodeDescendants: true
                },
                actionMode: '',
                showCheckedPanel: false,
                checkedList: [],
                showBatchDropdown: false
            }
        },
        computed: {
            ...mapGetters('objectBiz', ['bizId']),
            topologyTree () {
                return this.$refs.topologyTree
            },
            checkedNodes () {
                if (this.actionMode === 'batch-edit') {
                    return this.checkedList.map(data => this.topologyTree.idGenerator(data))
                }
                return []
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
                this.treeOptions.selectable = false
                this.treeOptions.checkOnClick = true
                this.treeOptions.checkOnlyAvailableStrictly = true
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
                this.$routerActions.redirect({
                    name: MENU_BUSINESS_HOST_APPLY_EDIT,
                    query: {
                        mid: checkedIds.join(','),
                        batch: 1,
                        action: this.actionMode
                    },
                    history: true
                })
            },
            handleCancelEdit () {
                this.treeOptions.showCheckbox = false
                this.treeOptions.selectable = true
                this.treeOptions.checkOnClick = false
                this.showCheckedPanel = false
                this.treeOptions.checkOnlyAvailableStrictly = false
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
            margin-left: 8px;

            .dropdown-trigger {
                border: 1px solid #c4c6cc;
                border-radius: 2px;
                padding: 0 8px;
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

                .icon-angle-down {
                    font-size: 22px;
                    margin: -3px -5px 0 -4px
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
                margin: .1em;
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
                    right: 8px;
                    top: 10px;
                    width: 28px;
                    height: 28px;
                    text-align: center;
                    line-height: 28px;

                    a {
                        color: #c4c6cc;
                        &:hover {
                            color: #979ba5;
                        }
                    }
                }
            }
        }
    }

    .bk-dropdown-list {
        .auth-box {
            width: 100%;
        }

        > li a {
            display: block;
            height: 32px;
            line-height: 33px;
            padding: 0 16px;
            color: #63656e;
            font-size: 14px;
            text-decoration: none;
            white-space: nowrap;

            &:hover {
                background-color: #eaf3ff;
                color: #3a84ff;
            }

            &.disabled {
                color: #c4c6cc;
            }
        }
    }
</style>
