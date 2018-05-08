let tableIdSeed = 0

class TableLayout {
    constructor (options) {
        this.id = tableIdSeed++
        this.colgroup = []
        this.table = null
        this.columns = []
        this.scrollY = false
        for (let name in options) {
            if (options.hasOwnProperty(name)) {
                this[name] = options[name]
            }
        }
    }

    checkScrollY () {
        this.table.$nextTick(() => {
            let $bodyWrapper = this.table.$refs.bodyWrapper
            let $body = this.table.$refs.body.$el
            this.scrollY = $body.offsetHeight > $bodyWrapper.offsetHeight
        })
    }

    update () {
        this.table.$nextTick(() => {
            let columns = this.columns
            let bodyMinWidth = 0
            let bodyWidth = this.table.$refs.bodyWrapper.offsetWidth
            let flexColumns = columns.filter(({flex, dragging}) => flex && !dragging)
            if (flexColumns.length) { // 存在不固定宽度的列
                columns.forEach(column => {
                    bodyMinWidth += column.width
                })
                let gutterWidth = this.scrollY ? this.table.gutterWidth : 0
                if (bodyMinWidth <= bodyWidth - gutterWidth) { // 表格的最小宽度小于容器宽度，需对弹性单元格平均分配剩余宽度
                    let totalFlexWidth = bodyWidth - bodyMinWidth
                    if (flexColumns.length === 1) { // 只要一个弹性单元格，则其分配全部剩余宽度
                        flexColumns[0].realWidth = flexColumns[0].width + totalFlexWidth
                    } else { // 多个弹性单元格，前面n-1个平均分配向下取整的宽度，第n个获得总弹性宽度减去前面n-1总和的宽度
                        let allflexColumnsWidth = flexColumns.reduce((prev, column) => prev + column.width, 0)
                        let flexRatio = totalFlexWidth / allflexColumnsWidth
                        let notLastWidth = 0
                        flexColumns.forEach((column, index) => {
                            if (index !== (flexColumns.length - 1)) {
                                let flexWidth = Math.floor(column.width * flexRatio)
                                notLastWidth += flexWidth
                                column.realWidth = column.width + flexWidth
                            }
                        })
                        let lastColumn = flexColumns[flexColumns.length - 1]
                        lastColumn.realWidth = lastColumn.width + totalFlexWidth - notLastWidth
                    }
                }
            } else {
                columns.forEach(column => {
                    bodyMinWidth += column.realWidth
                })
            }
            const lastIndex = columns.length - 1
            columns.forEach((column, index) => {
                column.dragging = false
            })
            this.updateTableWidth(Math.max(bodyMinWidth, bodyWidth))
            this.updateColgroup()
        })
    }

    updateTableWidth (width) {
        let headWidth = width
        let bodyWidth = this.scrollY ? width - this.table.gutterWidth : width
        this.table.$refs.head.$el.style.width = headWidth + 'px'
        this.table.$refs.body.$el.style.width = bodyWidth + 'px'
    }

    updateColgroup () {
        this.colgroup = this.columns.map(column => {
            return column.realWidth
        })
    }

    updateColumnWidth (column, width) {
        column.width = width
        column.realWidth = width
        column.dragging = true
        this.update()
    }

    doLayout () {
        console.log('doLayout')
        if (this.columns.length) {
            this.checkScrollY()
            this.update()
        }
    }
}

export default TableLayout
