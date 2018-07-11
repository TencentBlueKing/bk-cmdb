let tableIdSeed = 0

class TableLayout {
    constructor (options) {
        this.id = tableIdSeed++
        this.colgroup = []
        this.table = null
        this.columns = []
        this.scrollY = false
        this.bodyWidth = null
        for (let name in options) {
            if (options.hasOwnProperty(name)) {
                this[name] = options[name]
            }
        }
    }

    checkScrollY () {
        let $bodyWrapper = this.table.$refs.bodyLayout
        let $body = this.table.$refs.body.$el
        this.scrollY = $body.offsetHeight > $bodyWrapper.offsetHeight
    }

    update () {
        const wrapperWidth = this.table.width
        let columns = this.columns
        let bodyMinWidth = 0
        let bodyWidth = wrapperWidth ? wrapperWidth - 2 : this.table.$refs.bodyLayout.offsetWidth
        let flexColumns = columns.filter(({flex, dragged}) => flex && !dragged)
        if (flexColumns.length) { // 存在不固定宽度的列
            columns.forEach(column => {
                bodyMinWidth += column.width || column.minWidth
            })
            let gutterWidth = this.scrollY ? this.table.gutterWidth : 0
            if (bodyMinWidth <= bodyWidth - gutterWidth) { // 表格的最小宽度小于容器宽度，需对弹性单元格平均分配剩余宽度
                let totalFlexWidth = bodyWidth - bodyMinWidth - gutterWidth
                if (flexColumns.length === 1) { // 只要一个弹性单元格，则其分配全部剩余宽度
                    flexColumns[0].realWidth = flexColumns[0].minWidth + totalFlexWidth
                } else { // 多个弹性单元格，前面n-1个平均分配向下取整的宽度，第n个获得总弹性宽度减去前面n-1总和的宽度
                    let allflexColumnsWidth = flexColumns.reduce((prev, column) => prev + column.minWidth, 0)
                    let flexRatio = totalFlexWidth / allflexColumnsWidth
                    let notLastWidth = 0
                    flexColumns.forEach((column, index) => {
                        if (index !== (flexColumns.length - 1)) {
                            let flexWidth = Math.floor(column.minWidth * flexRatio)
                            notLastWidth += flexWidth
                            column.realWidth = column.minWidth + flexWidth
                        }
                    })
                    let lastColumn = flexColumns[flexColumns.length - 1]
                    lastColumn.realWidth = lastColumn.minWidth + totalFlexWidth - notLastWidth
                }
            } else {
                columns.forEach(column => {
                    column.realWidth = column.width || column.minWidth
                })
            }
        } else if (columns.length) {
            columns.forEach(column => {
                column.realWidth = column.width || column.minWidth
                bodyMinWidth += column.realWidth
            })
            const gutterWidth = this.scrollY ? this.table.gutterWidth : 0
            if (bodyMinWidth <= bodyWidth - gutterWidth) {
                const lastIndex = columns.length - 1
                let notLastWidth = 0
                columns.forEach((column, index) => {
                    if (index !== lastIndex) {
                        notLastWidth += column.realWidth
                    }
                })
                columns[lastIndex].realWidth = bodyWidth - notLastWidth - gutterWidth
            }
        }
        this.bodyWidth = Math.max(bodyMinWidth, bodyWidth)
        this.updateColgroup()
    }

    updateColgroup () {
        this.colgroup = this.columns.map(column => {
            return column.realWidth
        })
    }

    updateColumnWidth (column, width) {
        const sortable = column.sortable
        column.width = width
        column.realWidth = width
        column.dragged = true
        column.sortable = false
        this.update()
        setTimeout(() => {
            column.sortable = sortable
        }, 0)
    }

    updateColumns () {
        const table = this.table
        const newColumns = []
        table.header.forEach((head, index) => {
            if (this.columns[index] && head === this.columns[index].head) {
                newColumns[index] = this.columns[index]
            } else {
                newColumns[index] = {
                    head: head,
                    id: head[table.valueKey],
                    name: head[table.labelKey],
                    attr: head.attr || {},
                    type: head.type || 'text',
                    width: head.width,
                    minWidth: 100,
                    realWidth: typeof head.width === 'number' ? head.width : 100,
                    flex: typeof head.width !== 'number',
                    sortable: table.sortable && (head.hasOwnProperty('sortable') ? head.sortable : true),
                    sortKey: head.hasOwnProperty('sortKey') ? head.sortKey : head[table.valueKey],
                    dragged: false
                }
            }
        })
        this.columns = newColumns
    }
}

export default TableLayout
