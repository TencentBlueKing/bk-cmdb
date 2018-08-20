<template>
    <table class="cc-table-head">
        <colgroup>
            <col v-for="width in layout.colgroup" :width="width">
            <col v-if="layout.scrollY" :width="table.gutterWidth">
        </colgroup>
        <tr :class="{'has-gutter': layout.scrollY}">
            <template v-for="(column, index) in columns">
                <th v-if="column.type === 'checkbox'" class="header-checkbox">
                    <label :for="getCheckboxId(column)" class="bk-form-checkbox bk-checkbox-small" v-if="table.multipleCheck && table.list.length">
                        <template v-if="layout.table.isCheckboxShow">
                            <input type="checkbox" v-if="table.crossPageCheck" ref="crossPageCheckbox"
                                :id="getCheckboxId(column)"
                                :disabled="!table.list.length"
                                @change="handleCheckAll($event, column)">
                            <input type="checkbox" v-else
                                :id="getCheckboxId(column)"
                                :disabled="!table.list.length"
                                :checked="allChecked"
                                @change="handleCheckAll($event, column)">
                        </template>
                    </label>
                </th>
                <th v-else
                    v-bind="column.attr"
                    :class="{'sortable': column.sortable}"
                    @mousemove="handleMouseMove($event, column)"
                    @mousedown.left="handleMouseDown($event, column)"
                    @mouseout="handleMouseOut($event, column)"
                    @click="handleSort(column)">
                    <div :class="['head-label', 'fl', {'has-sort': column.sortable}]" :title="column.name">{{column.name}}</div>
                    <div class="head-sort fl" v-if="column.sortable">
                        <i :class="['head-sort-angle', 'ascing', {'active': column.sortKey === sortStore.sort && sortStore.order === 'asc'}]"
                            @click.stop="handleSort(column, 'asc')">
                        </i>
                        <i :class="['head-sort-angle', 'descing', {'active': column.sortKey === sortStore.sort && sortStore.order === 'desc'}]"
                            @click.stop="handleSort(column, 'desc')">
                        </i>
                    </div>
                </th>
            </template>
            <th v-if="layout.scrollY" class="header-gutter"></th>
        </tr>
    </table>
</template>

<script>
    export default {
        props: {
            layout: Object
        },
        data () {
            return {
                dragStore: {
                    dragging: false,
                    column: null,
                    state: {}
                },
                sortStore: {
                    order: null,
                    sort: null
                }
            }
        },
        computed: {
            table () {
                return this.layout.table
            },
            allChecked () {
                return this.table.list.length && this.table.checked.length === this.table.list.length
            },
            columns () {
                return this.layout.columns
            },
            crossPageAllChecked () {
                if (this.table.crossPageCheck) {
                    return this.table.checked.length === this.table.pagination.count
                }
                return false
            }
        },
        watch: {
            crossPageAllChecked (checked) {
                if (this.$refs.crossPageCheckbox && this.$refs.crossPageCheckbox.length) {
                    this.$refs.crossPageCheckbox[0].checked = checked
                }
            }
        },
        created () {
            this.sortStore.sort = this.layout.table.defaultSort
        },
        mounted () {
            this.createResizeProxy()
        },
        methods: {
            getCheckboxId (column) {
                return `table-${this.layout.id}-${column.id}-checkbox`
            },
            getSortOrder (column) {
                const sortKey = column.sortKey
                if (sortKey === this.sortStore.sort) {
                    return this.sortStore.order === 'asc' ? 'desc' : 'asc'
                }
                return 'desc'
            },
            handleCheckAll ($event, column) {
                const isChecked = $event.target.checked
                let checked = []
                if (isChecked) {
                    checked = this.table.list.map(item => item[column[this.table.valueKey]])
                }
                this.table.$emit('update:checked', checked)
                this.table.$emit('handleCheckAll', isChecked, column.head)
            },
            handleSort (column, order) {
                if (!column.sortable) return
                if (!order) {
                    order = this.getSortOrder(column)
                }
                const sortKey = column.sortKey
                if (order === this.sortStore.order && sortKey === this.sortStore.sort) {
                    this.sortStore.order = null
                    this.sortStore.sort = this.table.defaultSort
                    this.table.$emit('handleSortChange', this.table.defaultSort)
                } else {
                    this.sortStore.order = order
                    this.sortStore.sort = sortKey
                    this.table.$emit('handleSortChange', order === 'asc' ? sortKey : `-${sortKey}`)
                }
            },
            handleMouseMove (event, column) {
                let $column = event.currentTarget
                let $next = $column.nextElementSibling
                // 鼠标划过单元格，如果不是最后一个则判断是否显示resize的鼠标图案
                if (column.flex && !this.dragStore.dragging && $next && !$next.classList.contains('header-gutter')) {
                    let rect = $column.getBoundingClientRect()
                    const bodyStyle = document.body.style
                    if (rect.width >= 70 && rect.right - event.pageX < 8) {
                        bodyStyle.cursor = 'col-resize'
                        $column.style.cursor = 'col-resize'
                        $column.style.borderRight = '1px solid #4d597d'
                        this.dragStore.column = $column
                    } else {
                        bodyStyle.cursor = 'default'
                        $column.style.cursor = column.sortable ? 'pointer' : 'default'
                        $column.style.borderRight = '1px solid #dde4eb'
                        this.dragStore.column = null
                    }
                }
            },
            // 鼠标移出，恢复默认鼠标图案
            handleMouseOut (event, column) {
                if (this.dragStore.column) {
                    this.dragStore.column.style.borderRight = '1px solid #dde4eb'
                }
                document.body.style.cursor = 'default'
            },
            handleMouseDown (event, column) {
                // 鼠标按下，如果是拖拽状态进行resize
                if (this.dragStore.column) {
                    this.dragStore.dragging = true
                    const $table = this.table.$el
                    const tableRect = $table.getBoundingClientRect()
                    const columnRect = this.dragStore.column.getBoundingClientRect()
                    const minLeft = columnRect.left - tableRect.left + 100
                    this.dragStore.state = {
                        startMouseLeft: event.clientX,
                        startLeft: columnRect.right - tableRect.left,
                        startColumnLeft: columnRect.left - tableRect.left
                    }
                    this.resizeProxy.style.left = this.dragStore.state.startLeft + 'px'
                    this.resizeProxy.style.visibility = 'visible'
                    document.onselectstart = () => { return false }
                    document.ondragstart = () => { return false }
                    // 动态设置拖拽线的位置
                    const handleMouseMove = (event) => {
                        const deltaLeft = event.clientX - this.dragStore.state.startMouseLeft
                        const proxyLeft = this.dragStore.state.startLeft + deltaLeft
                        const maxLeft = tableRect.width - 4 - (this.layout.scrollY ? this.table.gutterWidth : 0)
                        this.resizeProxy.style.left = Math.min(maxLeft, Math.max(minLeft, proxyLeft)) + 'px'
                    }
                    // 鼠标释放时，更新布局，重置拖拽状态
                    const handleMouseUp = (event) => {
                        if (this.dragStore.dragging) {
                            const { startColumnLeft, startLeft } = this.dragStore.state
                            const finalLeft = parseInt(this.resizeProxy.style.left, 10)
                            const columnWidth = finalLeft - startColumnLeft
                            this.layout.updateColumnWidth(column, columnWidth)

                            this.dragStore.dragging = false
                            this.dragStore.column = null
                            this.dragStore.state = {}
                            this.resizeProxy.style.visibility = 'hidden'
                            document.body.style.cursor = 'default'
                        }
                        document.removeEventListener('mousemove', handleMouseMove)
                        document.removeEventListener('mouseup', handleMouseUp)
                        document.onselectstart = null
                        document.ondragstart = null
                    }
                    document.addEventListener('mousemove', handleMouseMove)
                    document.addEventListener('mouseup', handleMouseUp)
                }
            },
            createResizeProxy () {
                this.resizeProxy = document.createElement('i')
                this.resizeProxy.classList.add('col-resize-proxy')
                this.table.$el.appendChild(this.resizeProxy)
            }
        }
    }
</script>

<style lang="scss" scoped>
    table.cc-table-head{
        background-color: #fafbfd;
        color: $textColor;
        text-align: left;
        border-collapse: separate;
        border-spacing: 0;
        table-layout: fixed;
        tr {
            th{
                height: 40px;
                padding: 0 16px;
                border: 1px solid $tableBorderColor;
                border-left: none;
                line-height: normal;
                &:last-child{
                    border-right: none;
                }
                &.header-checkbox{
                    padding: 0;
                }
                &.header-gutter{
                    padding: 0;
                }
                &.sortable{
                    cursor: pointer;
                }
            }
        }
        tr.has-gutter{
            th:nth-last-child(2) {
                border-right: none;
            }
        }
    }
    .bk-form-checkbox{
        display: block;
        width: 100%;
        height: 100%;
        text-align: center;
        padding: 0;
        cursor: pointer;
        &:before{
            content: '';
            width: 0;
            height: 100%;
            display: inline-block;
            vertical-align: middle;
        }
        input[type='checkbox'] {
            display: inline-block;
            vertical-align: middle;
            &:disabled{
                cursor: not-allowed;
            }
        }
    }
    .head-label{
        font-size: 14px;
        max-width: 100%;
        @include ellipsis;
        &.has-sort{
            max-width: calc(100% - 25px);
        }
    }
    .head-sort{
        width: 20px;
        margin: 0 0 0 5px;
        .head-sort-angle{
            display: block;
            width: 0;
            height: 0;
            border: 5px solid transparent;
            cursor: pointer;
            &.ascing{
                border-bottom-color: #c3cdd7;
                margin-bottom: 2px;
                &:hover{
                    border-bottom-color: #737987;
                }
                &.active{
                    border-bottom-color: #4b8fe0;
                }
            }
            &.descing{
                border-top-color: #c3cdd7;
                margin-top: 2px;
                &:hover{
                    border-top-color: #737987;
                }
                &.active{
                    border-top-color: #4b8fe0;
                }
            }
        }
    }
</style>
<style lang="scss">
    i.col-resize-proxy{
        visibility: hidden;
        width: 0;
        height: 100%;
        position: absolute;
        top: 0;
        left: 0;
        border-left: 1px dashed #d1d5e0;
        pointer-events: none;
    }
</style>