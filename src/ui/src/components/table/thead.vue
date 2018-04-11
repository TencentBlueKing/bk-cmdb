/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2017-2018 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and limitations under the License.
 */

<template>
    <thead :class="['table-header', {'noSelect':cantBeSelected}]">
        <tr>
            <template v-if="localTableHeader.length">
                <template v-for="item in localTableHeader">
                    <th width="50" class="tc checkbox-wrapper" v-if="item.type === 'checkbox'">
                        <label class="bk-form-checkbox bk-checkbox-small" v-show="tableList.length && multipleCheck">
                            <input type="checkbox" v-model="isAllCheck">
                        </label>
                    </th>
                    <th v-else
                        v-bind="item.attr"
                        @mousedown.left="handleMouseDown($event, item)"
                        @mousemove="handleMouseMove($event, item)"
                        @mouseout="handleMouseOut($event, item)">
                        <span class="sort-box">
                            <span>{{$t(item.name)}}</span>
                            <span v-if="sortable && (item.sortable !== false) && item.type !== 'checkbox'">
                                <i class="sort-angle ascing" :class="{'cur-sort' : currentSort === item.id && currentOrder === ASC}" @click="handleTableSortClick($event,item,ASC)"></i>
                                <i class="sort-angle descing" :class="{'cur-sort' : currentSort === item.id && currentOrder === DESC}" @click="handleTableSortClick($event,item,DESC)"></i>
                            </span>
                        </span>
                    </th>
                </template>
            </template>
            <template v-else>
                <th height="40"></th>
            </template>
        </tr>
    </thead>
</template>
<script>
    export default {
        props: {
            cantBeSelected: {
                type: Boolean,
                default: true
            },
            hasCheckbox: {
                type: Boolean,
                default: true
            },
            tableHeader: {
                type: Array,
                required: true
            },
            tableList: {
                type: Array,
                default: () => {
                    return []
                }
            },
            defaultSort: {
                type: String,
                default: ''
            },
            chooseId: {
                type: Array,
                default: () => {
                    return []
                }
            },
            sortable: {
                type: Boolean,
                required: false,
                default: true
            },
            multipleCheck: Boolean
        },
        data () {
            return {
                isAllCheck: false,
                currentSort: this.defaultSort,
                currentOrder: '',
                ASC: Symbol('symbol for ascending'),
                DESC: Symbol('symbol for descending'),
                dragging: false,
                draggingColumn: null,
                dragState: {},
                resizeProxy: null,
                resizeProxyVisible: false,
                localTableHeader: [] // 兼容旧版table不在header中显式指定checkbox对象所用
            }
        },
        watch: {
            resizeProxyVisible (isVisible) {
                this.resizeProxy.style.visibility = isVisible ? 'visible' : 'hidden'
            },
            isAllCheck (isChecked) {
                this.$emit('handleTableAllCheck', isChecked)
            },
            chooseId (chooseIdArray) {
                if (!chooseIdArray.length) {
                    this.isAllCheck = false
                }
            },
            '$route.fullPath' () {
                this.isAllCheck = false
            },
            tableHeader (tableHeader) {
                this.initLocalTableHeader()
            }
        },
        methods: {
            initLocalTableHeader () {
                let localTableHeader = this.tableHeader.slice(0)
                if (this.hasCheckbox && this.tableHeader.length && this.tableHeader[0]['type'] !== 'checkbox') {
                    localTableHeader.unshift({
                        id: 'HostID',
                        name: 'ID',
                        type: 'checkbox'
                    })
                }
                this.localTableHeader = localTableHeader
            },
            handleTableSortClick (event, item, sortOrder) {
                let sortKey = item.sortKey ? item.sortKey : item.id
                if (sortOrder === this.currentOrder && sortKey === this.currentSort) {
                    this.currentOrder = ''
                    this.currentSort = this.defaultSort
                    this.$emit('handleTableSortClick', this.defaultSort)
                } else {
                    this.currentOrder = sortOrder
                    this.currentSort = sortKey
                    this.$emit('handleTableSortClick', sortOrder === this.ASC ? sortKey : ('-' + sortKey))
                }
            },
            handleMouseDown (event, item) {
                if (this.draggingColumn) {
                    this.dragging = true
                    const table = this.$parent
                    const tableEl = table.$el
                    const tableLeft = tableEl.getBoundingClientRect().left
                    const columnRect = this.draggingColumn.getBoundingClientRect()
                    const minLeft = columnRect.left - tableLeft + 70

                    this.dragState = {
                        startMouseLeft: event.clientX,
                        startLeft: columnRect.right - tableLeft,
                        startColumnLeft: columnRect.left - tableLeft,
                        tableLeft: tableLeft
                    }
                    this.resizeProxy.style.left = this.dragState.startLeft + 'px'
                    this.resizeProxyVisible = true

                    document.onselectstart = function () { return false }
                    document.ondragstart = function () { return false }
                    // 动态调整参考线的位置
                    const handleMouseMove = (event) => {
                        const deltaLeft = event.clientX - this.dragState.startMouseLeft
                        const proxyLeft = this.dragState.startLeft + deltaLeft
                        this.resizeProxy.style.left = Math.max(minLeft, proxyLeft) + 'px'
                    }
                    const handleMouseUp = (event) => {
                        if (this.dragging) {
                            // 以参考线的位置计算被调整列的宽度，并更新表格各列宽度
                            const { startColumnLeft, startLeft } = this.dragState
                            const finalLeft = parseInt(this.resizeProxy.style.left, 10)
                            const columnWidth = finalLeft - startColumnLeft
                            this.updateHeaderLayout(columnWidth)

                            this.dragging = false
                            this.draggingColumn = null
                            this.dragState = {}
                            this.resizeProxyVisible = false
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
            handleMouseMove (event, item) {
                // 边缘检测，当鼠标出现在表头列的右边界时提示可拖拽
                let column = event.currentTarget
                if (!this.dragging && column.nextElementSibling) {
                    let rect = column.getBoundingClientRect()
                    const bodyStyle = document.body.style
                    if (rect.width > 50 && rect.right - event.pageX < 8) {
                        bodyStyle.cursor = 'col-resize'
                        column.style.cursor = 'col-resize'
                        column.style.borderRight = '1px solid #4d597d'
                        this.draggingColumn = column
                    } else if (!this.dragging) {
                        bodyStyle.cursor = 'default'
                        column.style.cursor = 'default'
                        column.style.borderRight = '1px solid #dde4eb'
                        this.draggingColumn = null
                    }
                }
            },
            handleMouseOut (event, item) {
                if (this.draggingColumn) {
                    this.draggingColumn.style.borderRight = '1px solid #dde4eb'
                }
                document.body.style.cursor = 'default'
            },
            updateHeaderLayout (draggingColumnWidth) {
                // 为了不让列的宽度被表格自动调整出现闪烁
                // 被拖拽之前的列宽度不变，之后的列被自动调整后需不小于55px
                let fixedColumn = this.draggingColumn.previousElementSibling
                let adaptationColumn = this.draggingColumn.nextElementSibling
                while (fixedColumn) {
                    fixedColumn.style.width = fixedColumn.getBoundingClientRect().width + 'px'
                    fixedColumn = fixedColumn.previousElementSibling
                }
                this.draggingColumn.style.width = draggingColumnWidth + 'px'
                this.$nextTick(() => {
                    while (adaptationColumn) {
                        adaptationColumn.style.width = Math.max(adaptationColumn.getBoundingClientRect().width, 70) + 'px'
                        adaptationColumn = adaptationColumn.nextElementSibling
                    }
                })
            },
            createResizeProxy () {
                this.resizeProxy = document.createElement('i')
                this.resizeProxy.className = 'col-resize-proxy'
                this.$parent.$el.appendChild(this.resizeProxy)
            }
        },
        beforeMount () {
            this.initLocalTableHeader()
        },
        mounted () {
            this.createResizeProxy()
        }
    }
</script>
<style lang="scss" scoped>
    .bk-sort-box{
        right: 0;
    }
    .table-header {
        tr{
            th:not(.checkbox-wrapper){
                position: relative;
                overflow: hidden;
                text-overflow: ellipsis;
                white-space: nowrap;
            }
            th:first-child{
                border-left: none !important;
            }
            th:last-child{
                border-right: none !important;
            }
            th.checkbox-wrapper{
                padding: 0;
                margin: 0;
            }
            th label{
                padding: 10px 19px;
                margin-right: 0;
                cursor: pointer;
            }
        }
    }
    .sort-box{
        position: relative;
        padding: 0;
        .sort-angle{
            position: absolute;
            left: 100%;
            width: 0;
            height: 0;
            margin: 0 0 0 5px;
            border: 5px solid transparent;
            cursor: pointer;
            &.ascing{
                border-bottom-color: #c3cdd7;
                top: -1px;
                &:hover{
                    border-bottom-color: #737987;
                }
                &.cur-sort{
                    border-bottom-color: #4b8fe0;
                }
            }
            &.descing{
                border-top-color: #c3cdd7;
                bottom: -2px;
                &:hover{
                    border-top-color: #737987;
                }
                &.cur-sort{
                    border-top-color: #4b8fe0;
                }
            }
        }
    }
</style>
<style lang="scss">
    i.col-resize-proxy{
        visibility: hidden;
        width: 1px;
        height: 100%;
        position: absolute;
        top: 0;
        left: 0;
        background-color: #d1d5e0;
        pointer-events: none;
    }
</style>