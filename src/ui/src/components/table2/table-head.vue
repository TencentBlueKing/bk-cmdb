<template>
    <table class="cc-table-head">
        <colgroup>
            <col v-for="width in layout.colgroup" :width="width">
            <col v-if="layout.scrollY" :width="table.gutterWidth">
        </colgroup>
        <tr :class="{'has-gutter': layout.scrollY}">
            <template v-for="(column, index) in columns">
                <th v-if="column.type === 'checkbox'" class="header-checkbox">
                    <label :for="getCheckboxId(column)" class="bk-form-checkbox bk-checkbox-small">
                        <input type="checkbox" :id="getCheckboxId(column)">
                    </label>
                </th>
                <th v-else
                    v-bind="column.attr"
                    @mousemove="handleMouseMove($event, column)"
                    @mousedown.left="handleMouseDown($event, column)"
                    @mouseout="handleMouseOut($event, column)">
                    <div class="head-label">{{column.name}}</div>
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
                dragging: false,
                draggingColumn: null,
                dragState: {}
            }
        },
        computed: {
            table () {
                return this.layout.table
            },
            columns () {
                return this.layout.columns
            }
        },
        mounted () {
            this.createResizeProxy()
        },
        methods: {
            getCheckboxId (column) {
                return `table-${this.layout.id}-${column.id}-checkbox`
            },
            handleMouseMove (event, column) {
                let $column = event.currentTarget
                let $next = $column.nextElementSibling
                // 鼠标划过单元格，如果不是最后一个则判断是否显示resize的鼠标图案
                if (column.flex && !this.dragging && $next && !$next.classList.contains('header-gutter')) {
                    let rect = $column.getBoundingClientRect()
                    const bodyStyle = document.body.style
                    if (rect.width >= 70 && rect.right - event.pageX < 8) {
                        bodyStyle.cursor = 'col-resize'
                        $column.style.cursor = 'col-resize'
                        $column.style.borderRight = '1px solid #4d597d'
                        this.draggingColumn = $column
                    } else {
                        bodyStyle.cursor = 'default'
                        $column.style.cursor = 'default'
                        $column.style.borderRight = '1px solid #dde4eb'
                        this.draggingColumn = null
                    }
                }
            },
            // 鼠标移出，恢复默认鼠标图案
            handleMouseOut (event, column) {
                if (this.draggingColumn) {
                    this.draggingColumn.style.borderRight = '1px solid #dde4eb'
                }
                document.body.style.cursor = 'default'
            },
            handleMouseDown (event, column) {
                // 鼠标按下，如果是拖拽状态进行resize
                if (this.draggingColumn) {
                    this.dragging = true
                    const $table = this.table.$el
                    const tableRect = $table.getBoundingClientRect()
                    const columnRect = this.draggingColumn.getBoundingClientRect()
                    const minLeft = columnRect.left - tableRect.left + 70
                    this.dragState = {
                        startMouseLeft: event.clientX,
                        startLeft: columnRect.right - tableRect.left,
                        startColumnLeft: columnRect.left - tableRect.left
                    }
                    this.resizeProxy.style.left = this.dragState.startLeft - 2 + 'px'
                    this.resizeProxy.style.visibility = 'visible'
                    document.onselectstart = () => { return false }
                    document.ondragstart = () => { return false }
                    // 动态设置拖拽线的位置
                    const handleMouseMove = (event) => {
                        const deltaLeft = event.clientX - this.dragState.startMouseLeft
                        const proxyLeft = this.dragState.startLeft + deltaLeft
                        const maxLeft = tableRect.width - 4 - (this.layout.scrollY ? this.table.gutterWidth : 0)
                        this.resizeProxy.style.left = Math.min(maxLeft, Math.max(minLeft, proxyLeft)) + 'px'
                    }
                    // 鼠标释放时，更新布局，重置拖拽状态
                    const handleMouseUp = (event) => {
                        if (this.dragging) {
                            const { startColumnLeft, startLeft } = this.dragState
                            const finalLeft = parseInt(this.resizeProxy.style.left, 10)
                            const columnWidth = finalLeft - startColumnLeft
                            this.layout.updateColumnWidth(column, columnWidth)

                            this.dragging = false
                            this.draggingColumn = null
                            this.dragState = {}
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
                padding: 0 20px;
                border: 1px solid $tableBorderColor;
                border-left: none;
                &:last-child{
                    border-right: none;
                }
                &.header-checkbox{
                    padding: 0;
                }
                &.header-gutter{
                    padding: 0;
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