<template>
    <div class="sticky-row"
        :class="{ expanded }"
        :style="{
            top: top + 'px'
        }">
        <div class="sticky-row-cell selection-cell" :style="getCellStyle('id')">
            <bk-checkbox class="selection" :checked="checked" @change="handleCheck"></bk-checkbox>
        </div>
        <div class="sticky-row-cell center expand-cell" :style="getCellStyle('expand')">
            <i class="expand-icon bk-icon icon-play-shape" @click="handleExpand">
            </i>
        </div>
        <div class="sticky-row-cell" :style="getCellStyle('name')">
            <list-cell-name :row="row"></list-cell-name>
        </div>
        <div class="sticky-row-cell" :style="getCellStyle('process_count')">
            <list-cell-count :row="row"></list-cell-count>
        </div>
        <div class="sticky-row-cell" :style="getCellStyle('tag')">
            <list-cell-tag :row="row" @update-labels="handleUpdateLabels(row, ...arguments)"></list-cell-tag>
        </div>
    </div>
</template>

<script>
    import ListCellName from './list-cell-name'
    import ListCellCount from './list-cell-count'
    import ListCellTag from './list-cell-tag'
    export default {
        components: {
            ListCellName,
            ListCellCount,
            ListCellTag
        },
        props: {
            table: Object
        },
        data () {
            return {
                row: {},
                updateScrollPosition: false,
                prevScrollTop: 0,
                top: 0
            }
        },
        computed: {
            columns () {
                return this.table.columns
            },
            checked () {
                return this.table.store.states.selection.includes(this.row)
            },
            expanded () {
                return this.table.store.states.expandRows.includes(this.row)
            }
        },
        watch: {
            updateScrollPosition () {
                this.top = 0
            }
        },
        mounted () {
            this.table.$refs.bodyWrapper.addEventListener('scroll', this.updatePostion)
        },
        beforeDestroy () {
            this.table.$refs.bodyWrapper.removeEventListener('scroll', this.updatePostion)
        },
        methods: {
            updatePostion (event) {
                // const newScrollTop = event.target.scrollTop
                // const delta = newScrollTop - this.scrollTop
                // if (delta > 0 && this.updateScrollPosition) {
                //     this.top = this.top - delta
                // }
                // this.scrollTop = newScrollTop
            },
            getCellStyle (property) {
                const column = this.columns.find(column => column.property === property)
                if (column) {
                    return {
                        width: (column.realWidth || column.width) + 'px'
                    }
                }
                return {}
            },
            handleCheck () {
                this.table.toggleRowSelection(this.row, !this.checked)
            },
            handleExpand () {
                this.table.toggleRowExpansion(this.row, !this.expanded)
            },
            handleUpdateLabels (row, labels) {
                row.labels = labels
            }
        }
    }
</script>

<style lang="scss" scoped>
    .sticky-row {
        position: relative;
        display: flex;
        margin-bottom: -42px;
        height: 41px;
        align-items: center;
        border-bottom: 1px solid $borderColor;
        background-color: #f0f1f5;
        z-index: -1;
        &.expanded {
            position: sticky;
            top: 0;
            z-index: 1;
            .expand-icon {
                transform: rotate(90deg);
            }
        }
        &:hover {
            /deep/ {
                .instance-dot-menu {
                    display: inline-block;
                }
                .tag-edit {
                    visibility: visible;
                }
                .tag-empty {
                    display: none;
                }
            }
        }
    }
    .sticky-row-cell {
        height: 100%;
        display: flex;
        align-items: center;
        padding: 0 15px;
        &.center {
            justify-content: center;
        }
        &.expand-cell,
        &.selection-cell {
            padding: 0;
        }
        .selection {
            display: flex;
            width: 100%;
            height: 100%;
            justify-content: center;
            align-items: center;
        }
        .expand-icon {
            color: #c4c6cc;
            cursor: pointer;
            transition: transform .2s ease-in-out;
        }
    }
</style>
