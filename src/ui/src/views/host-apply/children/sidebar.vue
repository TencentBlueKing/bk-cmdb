<template>
    <div class="host-apply-sidebar">
        <div class="searchbar">
            <div class="search-select">
                <search-select-mix></search-select-mix>
            </div>
            <div class="action-menu">
                <bk-dropdown-menu>
                    <bk-button type="primary" class="is-icon" slot="dropdown-trigger">
                        <i class="bk-cc-icon icon-cc-list"></i>
                    </bk-button>
                    <ul class="bk-dropdown-list" slot="dropdown-content">
                        <li><a href="javascript:;" @click="handleBatchEdit">批量编辑</a></li>
                        <li><a href="javascript:;" @click="handleBatchDel">批量删除</a></li>
                    </ul>
                </bk-dropdown-menu>
            </div>
        </div>
        <topology-tree
            :tree-options="treeOptions"
            @selected="handleTreeSelected"
        ></topology-tree>
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
        data () {
            return {
                treeOptions: {
                    showCheckbox: false
                },
                actionMode: ''
            }
        },
        methods: {
            handleBatchEdit () {
                this.actionMode = 'batch-edit'
                this.treeOptions.showCheckbox = !this.treeOptions.showCheckbox
                this.$emit('update:actionMode', this.actionMode)
            },
            handleBatchDel () {
                this.actionMode = 'batch-del'
                this.treeOptions.showCheckbox = !this.treeOptions.showCheckbox
                this.$emit('update:actionMode', this.actionMode)
            },
            handleTreeSelected (node) {
                this.$emit('module-selected', node.data)
                console.log('handleTreeSelected', node)
            }
        }
    }
</script>

<style lang="scss" scoped>
    .host-apply-sidebar {

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
        }
    }
</style>
