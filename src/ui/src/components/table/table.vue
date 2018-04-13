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
    <div class="tabe-wrapper" v-bkloading="{isLoading: isLoading}">
        <div class="table-content">
            <div ref="scrollbar" class="table-scrollbar" :style="maxHeight ? `max-height: ${maxHeight}px` : ''">
                <table class="bk-table has-thead-bordered has-table-hover" :class="{'min-height-control': tableList.length === 0}" border="0" cellpadding="9" cellspacing="0">
                    <thead is="v-thead"
                        :tableHeader="tableHeader"
                        :tableList="tableList"
                        :defaultSort="defaultSort"
                        :chooseId="chooseId"
                        :sortable="sortable"
                        :hasCheckbox="hasCheckbox"
                        :multipleCheck="multipleCheck"
                        @handleTableSortClick="handleTableSortClick"
                        @handleTableAllCheck="handleTableAllCheck"
                    ></thead>
                    <tbody class="pr">
                        <template v-if="tableList.length">
                            <slot name="tableBodyRow"
                                v-for="(item,index) in tableList" 
                                :item="item">
                                <tr class="cp" @click="handleRowClick(item)" :class="{'selected': isRowSelected(item)}">
                                    <slot v-for="header in tableHeader"
                                        :name="header.id"
                                        :item="item">
                                        <td v-if="header.type === 'checkbox'" style="width: 50px;" class="checkbox-wrapper">
                                            <label class="bk-form-checkbox bk-checkbox-small" @click.stop="handleRowChoose(item)">
                                                <input type="checkbox" :value="item[header.id]" :checked="chooseId.indexOf(item[header.id]) !== -1" @change="setChooseId(item[header.id])">
                                            </label>
                                        </td>
                                        <td v-else>{{item[header.id]}}</td>
                                    </slot>
                                </tr>
                            </slot>
                        </template>
                        <template v-else>
                            <slot name="tableEmptyRow">
                                <tr>
                                    <td class="table-empty-col" :colspan="tableHeader.length">{{$t('Common[\'暂时没有数据\']')}}</td>
                                </tr>
                            </slot>
                        </template>
                    </tbody>
                </table>
            </div>  
            <v-pagination ref="pagination" v-if="pagination"
                :pagination="pagination" 
                :tableList="tableList"
                :hasCheckbox="hasCheckbox"
                :chooseId="chooseId"
                @onPageTurning="onPageTurning"
                @onPageSizeChange="onPageSizeChange"
                @handleSizeToggle="handleSizeToggle"
            ></v-pagination>
        </div>
    </div>
</template>

<script>
    import vThead from './thead'
    import vTbody from './tbody'
    import vPagination from './pagination'
    export default {
        components: {
            vThead,
            vTbody,
            vPagination
        },
        props: {
            tableHeader: {
                type: Array,
                required: true
            },
            tableList: {
                type: Array,
                required: true
            },
            defaultSort: {
                type: String,
                default: ''
            },
            sortable: {
                type: Boolean,
                default: true
            },
            pagination: {
                type: Object,
                required: false
            },
            isLoading: {
                type: Boolean,
                default: false
            },
            chooseId: {
                type: Array,
                default () {
                    return []
                }
            },
            maxHeight: {
                type: [Number, String],
                default: 0
            },
            multipleCheck: {
                type: Boolean,
                default: true
            }
        },
        data () {
            return {
                localChosen: []
            }
        },
        computed: {
            hasCheckbox () {
                if (this.tableHeader.length) {
                    return this.tableHeader[0]['type'] === 'checkbox'
                } else {
                    return false
                }
            }
        },
        beforeUpdate () {
            this.calcMaxHeight()
        },
        methods: {
            setChooseId (id) {
                let chooseId = [...this.chooseId]
                if (this.multipleCheck) {
                    let index = chooseId.indexOf(id)
                    if (index === -1) {
                        chooseId.push(id)
                    } else {
                        chooseId.splice(index, 1)
                    }
                    this.$emit('update:chooseId', chooseId)
                } else {
                    if (chooseId.indexOf(id) === -1) {
                        this.$emit('update:chooseId', [id])
                    } else {
                        this.$emit('update:chooseId', [])
                    }
                }
            },
            isRowSelected (item) {
                let isRowSelected = false
                this.tableHeader.map(header => {
                    if (header.type === 'checkbox') {
                        isRowSelected = this.chooseId.indexOf(item[header.id]) !== -1
                    }
                })
                return isRowSelected
            },
            calcMaxHeight () {
                this.$nextTick(() => {
                    let scrollbar = this.$refs.scrollbar
                    if (this.maxHeight) {
                        scrollbar.style.maxHeight = typeof this.maxHeight === 'string' ? this.maxHeight : `${this.maxHeight}px`
                    } else {
                        let scrollbarRect = scrollbar.getBoundingClientRect()
                        let bodyRect = document.body.getBoundingClientRect()
                        let pagingRect = this.$refs.pagination ? this.$refs.pagination.$el.getBoundingClientRect() : {height: 0}
                        let paddingCompensate = 40
                        scrollbar.style.maxHeight = `${bodyRect.height - scrollbarRect.top - pagingRect.height - paddingCompensate}px`
                    }
                })
            },
            handleTableSortClick () {
                this.$emit('handleTableSortClick', ...arguments)
            },
            handleTableAllCheck () {
                this.$emit('handleTableAllCheck', ...arguments)
            },
            handleRowClick () {
                this.$emit('handleRowClick', ...arguments)
            },
            onPageTurning () {
                this.$emit('onPageTurning', ...arguments)
                this.$emit('handlePageTurning', ...arguments)
            },
            onPageSizeChange () {
                this.$emit('onPageSizeChange', ...arguments)
                this.$emit('handlePageSizeChange', ...arguments)
            },
            handleRowChoose () {
                this.$emit('handleRowChoose', ...arguments)
            },
            handleSizeToggle (isOpen) {
                const sizeListHeight = 170
                if (isOpen) {
                    this.$el.style.minHeight = `${this.$el.getBoundingClientRect().height + sizeListHeight}px`
                } else {
                    this.$el.style.minHeight = 'auto'
                }
                this.$emit('handleSizeToggle', isOpen)
            }
        }
    }
</script>

<style media="screen" lang="scss" scoped>
    $borderColor: #e7e9ef; //边框色
    $defaultColor: #ffffff; //默认
    $primaryColor: #fafbfd; //主要
    $fnMainColor: #bec6de; //文案主要颜色
    $primaryHoverColor: #6b7baa; //鼠标移上 主要颜色
    $tableBorderColor: #dde4eb;
    .table-content{
        border: 1px solid $tableBorderColor;
        .table-scrollbar{
            width: 100%;
            overflow: auto;
            @include scrollbar;
        }
        .bk-tab2-content{
            height: 100%;
        }
        .bk-table{
            position: relative;
            border-top: none;
            border: none;
            border-collapse: collapse;
            table-layout:fixed;
            .table-header {
                tr{
                    th{
                        position: relative;
                    }
                    th:first-child{
                        padding: 0 !important;
                        margin: 0 !important;
                    }
                    th label{
                        padding: 10px 19px!important;
                        margin-right: 0;
                        cursor: pointer;
                    }
                }
            }
            tbody {
               tr{
                    td label{
                        display: block;
                        margin: 0;
                        padding: 0;
                        line-height: 40px;
                        cursor: pointer;
                    }
                }
                tr:hover{
                    background-color: #f1f7ff;
                }
            }

            tr th:first-child{
                border-left: none !important;
            }
            .checkbox-wrapper{
                text-align: center;
                padding: 0;
                input{
                    margin: 0;
                }
            }
            td:not(.checkbox-wrapper),th:not(.checkbox-wrapper){
                overflow: hidden;
                text-overflow:ellipsis;
                white-space: nowrap;
            }
            tr th{
                border-left: none !important;
            }
            tr th:last-child {
                border-right: none !important
            }
        }
        .min-height-control{
            min-height: 260px;
            background: #ffffff;
        }
        .table-empty{
            font-size: 14px;
            color: #6b7baa;
            .table-empty-col{
                height: 220px;
                background: #fff;
                vertical-align: middle;
                text-align: center;
            }
        }
        .table-page-contain{
            width:100%;
            padding:5px 20px;
            background:$primaryColor;
            .bk-page{
                height: 32px;
                >ul{
                    height: 32px;
                }
                .page-item{
                    min-width: 32px;
                    height: 32px;
                    line-height: 32px;
                }
            }
            .bk-page-compact{
                float:right;
            }
            .page-info{
                font-size:12px;
                color:#838fb6;
            }
        }
    }
    .select_page_setting {
        display: inline-block;
        width: 60px;
        margin-left: 10px;
        .bk-select-input{
            height: 30px;
            line-height: 30px;
            .icon-angle-down{
                top: 9px;
            }
        }
    }
    .table-empty-col{
        vertical-align: middle;
        text-align: center;
        background-color: #fff;
    }
</style>