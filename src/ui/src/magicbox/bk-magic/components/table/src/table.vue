<template>
    <table :class="['bk-table has-thead-bordered', className]">
        <thead>
            <tr>
                <template v-for="field in tableFields">
                    <template v-if="isSpecialField(field.name)">
                        <th v-if="extractName(field.name) === '__checkbox'">
                            <label class="bk-form-checkbox">
                                <input type="checkbox" @change="toggleAllCheckboxes" v-model="selectAll" :checked="selectAll">
                            </label>
                        </th>
                        <th v-if="extractName(field.name) === '__slot'">
                            {{field.title}}
                        </th>
                    </template>
                    <template v-else>
                        <th :id="'_' + field.name" :class="[{'bk-table-sortable': field.sortable}]">
                            {{field.title}}
                            <span class="bk-sort-box" v-if="field.sortable">
                      <i class="bk-icon icon-angle-up ascing" :class="[{'cur-sort': sortParams && sortParams[field.name] == 'asc'}]" @click="ascField(field.name)"></i>
                      <i class="bk-icon icon-angle-down descing" :class="[{'cur-sort': sortParams && sortParams[field.name] == 'desc'}]" @click="descField(field.name)"></i>
                    </span>
                        </th>
                    </template>
                </template>
            </tr>
        </thead>
        <tbody>
            <template v-if="curPageData.length">
                <tr v-for="(item, index) in curPageData">
                    <template v-for="field in tableFields">
                        <template v-if="isSpecialField(field.name)">
                            <td v-if="extractName(field.name) === '__checkbox'">
                                <label class="bk-form-checkbox">
                                    <input type="checkbox" @change="toggleCheckbox(item)" :checked="checkRowSelected(item)">
                                </label>
                            </td>
                            <td v-if="extractName(field.name) === '__slot'">
                                <slot :row-data="item" :row-index="index" :name="extractArgs(field.name)"></slot>
                            </td>
                        </template>
                        <template v-else>
                            <td v-html="item[field.name]"></td>
                        </template>
                    </template>
                </tr>
            </template>
            <template v-else>
                <tr>
                    <td :colspan="tableFields.length">
                        <div class="bk-message-box">
                            <p class="message empty-message">暂无数据</p>
                        </div>
                    </td>
                </tr>
            </template>
        </tbody>
    </table>
</template>
<script>
    /**
     * bk-table
     */
    export default {
        name: 'bk-table',
        props: {
            className: {
                type: String,
                default: ''
            },
            trackBy: {
                type: String,
                default: 'id'
            },
            fields: {
                type: Array,
                required: true
            },
            pageSize: {
                type: Number,
                default: 0
            },
            curPage: {
                type: Number,
                default: 1
            },
            apiUrl: {
                type: String,
                default: ''
            },
            data: {
                type: Array,
                default () {
                    return []
                }
            },
            queryParams: {
                type: Object,
                default () {
                    return {
                        sort: '',
                        keyword: ''
                    }
                }
            },
            httpParams: {
                type: Object,
                default () {
                    return {}
                }
            },
            css: {
                type: String,
                default: ''
            }
        },
        data () {
            return {
                tablePageSize: 1,
                tableTotalPage: 1,
                tableCurPage: 1,
                loading: true,
                selectAll: false,
                curPageData: [],
                eventPrefix: '',
                tableData: [],
                rowSelected: [],
                tableFields: [],
                sortParams: null
            }
        },
        created () {
            this.init()
        },
        methods: {
            emitEvent (eventName, args) {
                this.$emit(this.eventPrefix + eventName, args)
            },
            isSpecialField (fieldName) {
                return fieldName.slice(0, 2) === '__'
            },
            extractName (str) {
                return str.split(':')[0].trim()
            },
            renderTitle (field) {
                let _html = field.title
                return _html
            },
            extractArgs (string) {
                return string.split(':')[1]
            },
            checkRowSelected (data) {
                let trackIndex = this.trackBy
                let key = data[trackIndex]
                return this.isSelectedRow(key)
            },
            isSelectedRow (key) {
                return this.rowSelected.indexOf(key) >= 0
            },
            selectRowById (key) {
                if (!this.isSelectedRow(key)) {
                    this.rowSelected.push(key)
                }
            },
            unSelectedRowById (key) {
                this.rowSelected = this.rowSelected.filter((item) => {
                    return item !== key
                })
            },
            toggleCheckbox (item) {
                let trackIndex = this.trackBy
                let key = item[trackIndex]
                if (this.isSelectedRow(key)) {
                    this.unSelectedRowById(key)
                } else {
                    this.selectRowById(key)
                }
                this.selectAll = false
            },
            toggleAllCheckboxes () {
                let isChecked = this.selectAll
                let trackIndex = this.trackBy
                if (isChecked) {
                    this.selectAll = true
                    this.curPageData.forEach((item) => {
                        this.selectRowById(item[trackIndex])
                    })
                } else {
                    this.selectAll = false
                    this.curPageData.forEach((item) => {
                        this.unSelectedRowById(item[trackIndex])
                    })
                }
            },
            removeRowByCurIndex (data, curPageRowIndex) {
                this.removeRow(data)
            },
            removeRowById (id) {
                let trackIndex = this.trackBy
                let removeData = null
                this.tableData = this.tableData.filter((item) => {
                    if (item[trackIndex] === id) {
                        removeData = item
                    }
                    return item[trackIndex] !== id
                })
                return removeData
            },
            removeRow (data, curPageRowIndex) {
                let trackIndex = this.trackBy
                let removeData = null
                if (data.length) {
                    data.forEach((item) => {
                        let id = item[trackIndex]
                        removeData = this.removeRowById(id)
                    })
                } else {
                    let id = data[trackIndex]
                    removeData = this.removeRowById(id)
                }
                this.reloadCurPage()
                this.emitEvent('remove-row', {
                    rowData: removeData,
                    totalPage: this.tableTotalPage,
                    curPage: this.tableCurPage
                })
            },
            removeSelectedRow () {
                let rowSelected = this.rowSelected
                let self = this
                rowSelected.forEach((id) => {
                    self.removeRowById(id)
                })
                this.reloadCurPage()
                this.emitEvent('remove-row', {
                    totalPage: this.totalPage
                })
            },
            orderBy (data, sortKey, sortDirection) {
                if (sortDirection === 'desc') {
                    data.sort((value1, value2) => {
                        if (value1[sortKey] < value2[sortKey]) {
                            return 1
                        } else if (value1[sortKey] > value2[sortKey]) {
                            return -1
                        } else {
                            return 0
                        }
                    })
                } else {
                    data.sort((value1, value2) => {
                        if (value1[sortKey] < value2[sortKey]) {
                            return -1
                        } else if (value1[sortKey] > value2[sortKey]) {
                            return 1
                        } else {
                            return 0
                        }
                    })
                }

                return data
            },
            getDataByPage (page) {
                let startIndex = (page - 1) * this.tablePageSize
                let endIndex = page * this.tablePageSize
                if (startIndex < 0) {
                    startIndex = 0
                }
                if (endIndex > this.tableData.length) {
                    endIndex = this.tableData.length
                }
                let data = this.tableData.slice(startIndex, endIndex)
                if (this.sortParams) {
                    data = this.orderBy(data, this.sortParams.sortKey, this.sortParams.sortDirection)
                }
                this.selectAll = false
                return data
            },
            reloadCurPage () {
                this.initPageConf()
                if (this.tableCurPage > this.tableTotalPage) {
                    this.tableCurPage = this.tableTotalPage
                }
                this.curPageData = this.getDataByPage(this.tableCurPage)
            },
            renderCurPageData (page) {
                this.curPageData = this.getDataByPage(page)
            },
            ascField (name) {
                let sortParams = this.sortParams
                if (sortParams && sortParams[name] === 'asc') {
                    this.sortParams = null
                } else {
                    this.sortParams = {
                        sortKey: name,
                        sortDirection: 'asc'
                    }
                    this.sortParams[name] = 'asc'
                }
                this.reloadCurPage()
            },
            descField (name) {
                let sortParams = this.sortParams
                if (sortParams && sortParams[name] === 'desc') {
                    this.sortParams = null
                } else {
                    this.sortParams = {
                        sortKey: name,
                        sortDirection: 'desc'
                    }
                    this.sortParams[name] = 'desc'
                }
                this.reloadCurPage()
            },
            initPageConf () {
                let total = this.tableData.length
                this.tableCurPage = this.curPage
                if (this.pageSize) {
                    this.tablePageSize = this.pageSize
                    this.tableTotalPage = Math.ceil(total / this.pageSize)
                } else {
                    this.tablePageSize = total
                    this.tableTotalPage = 1
                }
            },
            init () {
                this.tableData = this.data
                this.initPageConf()
                this.curPageData = this.getDataByPage(this.tableCurPage)

                if (typeof this.fields === 'undefined') {
                    this.warn('请配置feields')
                    return false
                }

                this.tableFields = []
                let obj

                this.fields.forEach((field, i) => {
                    if (typeof field === 'string') {
                        obj = {
                            name: field,
                            title: field,
                            titleClass: '',
                            callback: null,
                            visible: true
                        }
                    } else {
                        obj = {
                            name: field.name,
                            title: field.title,
                            sortable: field.sortable,
                            titleClass: field.titleClass,
                            callback: field.callback,
                            visible: field.visible
                        }
                    }
                    this.tableFields.push(obj)
                })
            }
        }
    }
</script>
