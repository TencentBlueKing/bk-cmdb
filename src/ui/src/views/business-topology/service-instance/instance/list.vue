<template>
    <bk-table class="instance-table" ref="instanceTable"
        v-bkloading="{ isLoading: $loading(request.getList) }"
        :row-class-name="getRowClassName"
        :data="list"
        :pagination="pagination"
        :max-height="$APP.height - 250"
        @page-change="handlePageChange"
        @page-limit-change="handlePageLimitChange"
        @expand-change="handleExpandChange"
        @selection-change="handleSelectionChange"
        @row-click="handleRowClick">
        <bk-table-column type="selection" prop="id"></bk-table-column>
        <bk-table-column type="expand" prop="expand" width="15" :before-expand-change="beforeExpandChange">
            <div slot-scope="{ row }" v-bkloading="{ isLoading: row.pending }">
                <expand-list :service-instance="row" @update-list="handleExpandResolved(row, ...arguments)"></expand-list>
            </div>
        </bk-table-column>
        <bk-table-column :label="$t('实例名称')" prop="name" width="340">
            <list-cell-name slot-scope="{ row }" :row="row"
                @edit="handleEditName(row)"
                @success="handleEditNameSuccess(row, ...arguments)"
                @cancel="handleCancelEditName(row)">
            </list-cell-name>
        </bk-table-column>
        <bk-table-column :label="$t('进程数量')" prop="process_count" width="240" :resizable="false">
            <list-cell-count slot-scope="{ row }" :row="row" @refresh-count="handleRefreshCount"></list-cell-count>
        </bk-table-column>
        <bk-table-column :label="$t('标签')" prop="tag" min-width="150">
            <list-cell-tag slot-scope="{ row }" :row="row" @update-labels="handleUpdateLabels(row, ...arguments)"></list-cell-tag>
        </bk-table-column>
        <bk-table-column :label="$t('操作')" :min-width="$i18n.locale === 'en' ? 200 : 150">
            <list-cell-operation slot-scope="{ row }" :row="row" @refresh-count="handleRefreshCount"></list-cell-operation>
        </bk-table-column>
    </bk-table>
</template>

<script>
    import { mapGetters } from 'vuex'
    import ListCellName from './list-cell-name'
    import ListCellCount from './list-cell-count'
    import ListCellTag from './list-cell-tag'
    import ListStickyRow from './list-sticky-row'
    import ListCellOperation from './list-cell-operation'
    import ExpandList from './expand-list'
    import RouterQuery from '@/router/query'
    import Bus from '../common/bus'
    import Vue from 'vue'
    import store from '@/store'
    import i18n from '@/i18n'
    import LabelBatchDialog from './dialog/label-batch-dialog.js'
    export default {
        components: {
            ListCellName,
            ListCellCount,
            ListCellTag,
            ListCellOperation,
            ExpandList
        },
        data () {
            return {
                list: [],
                selection: [],
                pagination: this.$tools.getDefaultPaginationConfig(),
                filters: [],
                request: {
                    getList: Symbol('getList')
                }
            }
        },
        computed: {
            ...mapGetters(['supplierAccount']),
            ...mapGetters('objectBiz', ['bizId']),
            ...mapGetters('businessHost', ['selectedNode']),
            searchKey () {
                const nameFilter = this.filters.find(data => data.id === 'name')
                return nameFilter ? nameFilter.values[0].name : ''
            },
            searchTag () {
                const tagFilters = []
                this.filters.forEach(data => {
                    if (data.id === 'name') return
                    if (data.hasOwnProperty('condition')) {
                        tagFilters.push({
                            key: data.condition.id,
                            operator: 'in',
                            values: data.values.map(value => value.name)
                        })
                    } else {
                        const [{ id }] = data.values
                        tagFilters.push({
                            key: id,
                            operator: 'exists',
                            values: []
                        })
                    }
                })
                return tagFilters
            }
        },
        watch: {
            selectedNode () {
                this.handlePageChange(1)
            }
        },
        created () {
            this.unwatch = RouterQuery.watch(['page', 'limit', '_t', 'view'], ({
                page = this.pagination.current,
                limit = this.pagination.limit,
                view = 'instance'
            }) => {
                if (view !== 'instance') {
                    return false
                }
                this.pagination.page = parseInt(page)
                this.pagination.limit = parseInt(limit)
                this.getList()
            }, { immediate: true, throttle: true })
            Bus.$on('expand-all-change', this.handleExpandAllChange)
            Bus.$on('filter-change', this.handleFilterChange)
            Bus.$on('batch-edit-labels', this.handleBatchEditLabels)
        },
        beforeDestroy () {
            Bus.$off('expand-all-change', this.handleExpandAllChange)
            Bus.$off('filter-change', this.handleFilterChange)
            Bus.$off('batch-edit-labels', this.handleBatchEditLabels)
            this.unwatch()
            // this.removeIntersectionObserver()
        },
        methods: {
            handleFilterChange (filters) {
                this.filters = filters
                RouterQuery.set({
                    page: 1,
                    _t: Date.now()
                })
            },
            async getList () {
                try {
                    // this.removeIntersectionObserver()
                    // this.removeStickyRow()
                    const { count, info } = await this.$store.dispatch('serviceInstance/getModuleServiceInstances', {
                        params: {
                            bk_biz_id: this.bizId,
                            bk_module_id: this.selectedNode.data.bk_inst_id,
                            page: this.$tools.getPageParams(this.pagination),
                            search_key: this.searchKey,
                            selectors: this.searchTag,
                            with_name: true
                        },
                        config: {
                            requestId: this.request.getList,
                            cancelPrevious: true
                        }
                    })
                    this.list = info.map(data => ({ ...data, pending: true, editing: { name: false } }))
                    this.pagination.count = count
                } catch (error) {
                    this.list = []
                    this.pagination.count = 0
                    console.error(error)
                } finally {
                    this.$nextTick(() => {
                        // this.initIntersectionObserver()
                        // this.injectStickyRow()
                    })
                }
            },
            initIntersectionObserver () {
                if (!this.list.length) return false
                const bodyWrapper = this.$refs.instanceTable.$refs.bodyWrapper
                const rows = Array.from(bodyWrapper.querySelectorAll('.instance-table-row'))
                this.referenceObserver = this.createIntersectionObserver(entry => {
                    const index = rows.indexOf(entry.target)
                    this.stickyRow.row = this.list[index] || {}
                    this.stickyRow.updateScrollPosition = false
                }, { root: bodyWrapper, threshold: [0, 1] })
                // this.positionObserver = this.createIntersectionObserver(entry => {
                //     this.stickyRow.updateScrollPosition = true
                //     console.log(entry.target)
                // }, { root: bodyWrapper, rootMargin: '-42px 0px 0px 0px', threshold: [0, 1] })
                rows.forEach(row => {
                    this.referenceObserver.observe(row)
                    // this.positionObserver.observe(row)
                })
            },
            createIntersectionObserver (callback, options = {}) {
                return new IntersectionObserver(entries => {
                    const referenceEntry = entries.find(entry => {
                        const { isIntersecting, rootBounds, boundingClientRect } = entry
                        return isIntersecting && rootBounds.top > boundingClientRect.top
                    })
                    referenceEntry && callback(referenceEntry)
                }, options)
            },
            updateIntersectionObserver (row) {
                const bodyWrapper = this.$refs.instanceTable.$refs.bodyWrapper
                const tableRows = Array.from(bodyWrapper.querySelectorAll('.instance-table-row'))
                const index = this.list.indexOf(row)
                const tableRow = tableRows[index]
                const willExpand = !tableRow.classList.contains('expanded')
                if (willExpand) {
                    setTimeout(() => {
                        this.referenceObserver.observe(tableRow)
                        this.referenceObserver.observe(tableRow.nextElementSibling)
                    }, 0)
                } else {
                    this.referenceObserver.unobserve(tableRow)
                    this.referenceObserver.unobserve(tableRow.nextElementSibling)
                }
            },
            removeIntersectionObserver () {
                this.referenceObserver && this.referenceObserver.disconnect()
                this.positionObserver && this.positionObserver.disconnect()
                this.referenceObserver = null
                this.positionObserver = null
            },
            injectStickyRow () {
                if (!this.list.length) return false
                const bodyWrapper = this.$refs.instanceTable.$refs.bodyWrapper
                if (!this.stickyRow) {
                    const StickyRow = Vue.extend(ListStickyRow)
                    this.stickyRow = new StickyRow({
                        store,
                        i18n,
                        propsData: {
                            table: this.$refs.instanceTable
                        }
                    })
                    this.stickyRow.$mount()
                }
                this.stickyRow.row = this.list[0]
                bodyWrapper.insertBefore(this.stickyRow.$el, bodyWrapper.firstElementChild)
            },
            removeStickyRow () {
                if (!this.stickyRow) return false
                try {
                    const bodyWrapper = this.$refs.instanceTable.$refs.bodyWrapper
                    bodyWrapper.scrollTop = 0
                    bodyWrapper.removeChild(this.stickyRow.$el) // may be not child, catch error
                    this.stickyRow.row = {}
                    this.stickyRow.updateScrollPosition = false
                } catch (error) {}
            },
            getRowClassName ({ row }) {
                const className = ['instance-table-row']
                if (!row.process_count) {
                    className.push('disabled')
                }
                return className.join(' ')
            },
            handlePageChange (page) {
                this.pagination.current = page
                RouterQuery.set({
                    page: page,
                    _t: Date.now()
                })
            },
            handlePageLimitChange (limit) {
                this.pagination.limit = limit
                this.pagination.page = 1
                RouterQuery.set({
                    limit: limit,
                    page: 1,
                    _t: Date.now()
                })
            },
            beforeExpandChange ({ row }) {
                return !!row.process_count
            },
            handleExpandAllChange (expanded) {
                this.list.forEach(row => {
                    row.process_count && this.$refs.instanceTable.toggleRowExpansion(row, expanded)
                })
            },
            async handleExpandChange (row, expandedRows) {
                row.pending = expandedRows.includes(row)
            },
            handleSelectionChange (selection) {
                this.selection = selection
                Bus.$emit('instance-selection-change', selection)
            },
            handleRowClick (row) {
                this.$refs.instanceTable.toggleRowExpansion(row)
            },
            handleExpandResolved (row, list) {
                row.pending = false
                this.handleRefreshCount(row, list.length)
                if (!row.process_count) {
                    this.$refs.instanceTable.toggleRowExpansion(row, false)
                }
            },
            handleRefreshCount (row, newCount) {
                row.process_count = newCount
            },
            handleUpdateLabels (row, labels) {
                row.labels = labels
                Bus.$emit('update-labels')
            },
            handleBatchEditLabels () {
                LabelBatchDialog.show({
                    serviceInstances: this.selection,
                    updateCallback: (removedKeys, newLabels) => {
                        this.selection.forEach(instance => {
                            instance.labels && removedKeys.forEach(key => delete instance.labels[key])
                            instance.labels = Object.assign({}, instance.labels, newLabels)
                        })
                    }
                })
            },
            handleEditName (row) {
                this.list.forEach(row => (row.editing.name = false))
                row.editing.name = true
            },
            handleEditNameSuccess (row, value) {
                row.name = value
                row.editing.name = false
            },
            handleCancelEditName (row) {
                row.editing.name = false
            }
        }
    }
</script>

<style lang="scss" scoped>
    .instance-table {
        /deep/ {
            .instance-table-row {
                &:hover,
                &.expanded {
                    td {
                        background-color: #f0f1f5;
                    }
                    .tag-item {
                        background-color: #dcdee5;
                    }
                }
                &:hover {
                    // 试用操作列替代聚合菜单
                    // .instance-dot-menu {
                    //     display: inline-block;
                    // }
                    .tag-edit {
                        visibility: visible;
                    }
                    .tag-empty {
                        display: none;
                    }
                    .instance-name:not(.disabled) {
                        .name-edit {
                            visibility: visible;
                        }
                    }
                }
                &.disabled {
                    .bk-table-expand-icon {
                        display: none;
                        cursor: not-allowed;
                    }
                }
                td {
                    position: sticky;
                    top: 0;
                    z-index: 1;
                    background-color: #fff;
                }
            }
            .bk-table-expanded-cell {
                padding-left: 80px !important;
            }
        }
    }
</style>
