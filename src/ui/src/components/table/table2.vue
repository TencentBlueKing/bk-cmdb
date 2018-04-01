/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2017-2018 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and limitations under the License.
 */

<template lang="html">
    <div class="" v-bkloading="{isLoading: isLoading}">
        <div class="table-content">
            <div class="table-scrollbar">
                <table class="bk-table has-thead-bordered has-table-hover" :class="{'min-height-control':tableList.length === 0}" table border="0" cellpadding="0" cellspacing="0">
                    <thead is="v-thead"
                        :cantBeSelected="cantBeSelected"
                        :hasCheckbox="hasCheckbox"
                        :tableHeader="tableHeader"
                        :tableList="tableList"
                        :defaultSort="defaultSort"
                        :chooseId="chooseId"
                        :sortable="sortable"
                        @handleTableSortClick="handleTableSortClick"
                        @handleTableAllCheck="handleTableAllCheck"
                    ></thead>
                    <tbody class="pr" v-if="tableList.length !== 0" >
                        <template >
                            <tr class="cp" v-for="(item, index) in tableList" :class="{'selected': chooseId.indexOf(item.HostID) !== -1}">
                                <td v-show="hasCheckbox" style="width: 50px;">
                                    <label class="bk-form-checkbox bk-checkbox-small">
                                        <input type="checkbox" :value="item.HostID" v-model="chooseId">
                                    </label>
                                </td>
                                <td v-for="header in tableHeader" @click="showDetail(item)" >
                                    <template v-if='header.id === customize.id'>
                                        <button type="button" :class="['info-btn main-btn', {'vice-btn': item[header.id] === 0}]" @click="changeBind(item)">
                                            <span v-if="item[header.id] === 0">未绑定</span>
                                            <span v-else>已绑定</span>
                                        </button>
                                        <!-- <th v-if="extractName(field.name) === '__slot'">
                                            {{field.title}}
                                        </th> -->
                                    </template>
                                    <template v-else-if="header.id === '__slot__eventpush__statistics' && item.statistics">
                                        <i class="circle" :class="[{'danger':item.statistics.failure},{'success':!item.statistics.failure}]"></i>失败 {{item.statistics.failure}} / 总量 {{item.statistics.total}}
                                    </template>
                                    <template v-else-if="header.id === '__slot__eventpush__setting'">
                                        <i class="icon-cc-edit mr20" @click="editEvent(item)"></i><i class="icon-cc-del" @click="delEvent(item)"></i>
                                    </template>
                                    <template v-else>{{item[header.id]}}</template>
                                </td>
                            </tr>
                        </template>
                    </tbody>
                    <tbody class="table-empty"  v-else >
                        <tr>
                            <td class="table-empty-col" :colspan="hasCheckbox ? tableHeader.length + 1 : tableHeader.length">
                                <p>暂时没有数据</p>
                            </td>
                        </tr>
                    </tbody>
                </table>
            </div>
            <div class="table-page-contain clearfix">
                <div class="page-info fl" v-if="paging.totalPage">
                    <span class="mr20" v-if="hasCheckbox">已选<span> {{count}} </span>行</span>
                    <!-- <span class="ml20 mr20">显示 {{(paging.current - 1) * paging.showItem + 1}}-{{paging.current * paging.showItem}} 行</span> -->
                    <span> 第 {{paging.current}} 页 / 总 {{paging.totalPage}} 页</span>
                    <div v-show="hasCheckbox" style="display:inline-block"></span></div>
                    <span class="ml20 mr20">
                        每页显示<span class="select_page_setting mr5">
                            <bk-select 
                                :selected.sync="defaultPageSetting"
                                :list="pagelist">
                                <bk-select-option
                                    v-for="(option, index) of pagelist"
                                    :key="index"
                                    :value="option.value"
                                    :label="option.label">
                                </bk-select-option>
                            </bk-select>
                        </span>
                        行
                    </span>
                </div>
                <div class="bk-page bk-page-compact">
                    <ul class="pagination">
                        <li class="page-item" v-show="paging.current != 1" @click="paging.current-- && pageTuring(paging.current)" ><a href="javascript:;" class="page-button"><i class="icon-cc-angle-left"></i></a></li>
                        <li class="page-item" v-for="index in pages" @click="pageTuring(index)" :class="{'cur-page':paging.current == index}" :key="index">
                            <a href="javascript:;" class="page-button" >{{index}}</a>
                        </li>
                        <li class="page-item" v-show="paging.totalPage != paging.current && paging.totalPage != 0 " @click="paging.current++ && pageTuring(paging.current)"><a href="javascript:;" class="page-button" ><i class="icon-cc-angle-right"></i></a></li>
                    </ul>
                </div>
            </div>
        </div>
    </div>
</template>

<script type="text/javascript">
    import $ from 'jquery'
    import vThead from './thead'
    import qs from 'qs'
    export default {
        props: {
            tableHeader: {
                type: Array,
                default: [
                    {
                        id: '123',
                        name: 'ID'
                    },
                    {
                        id: 'PropertyName',
                        name: '业务名'
                    }
                ]
            },
            tableList: {
                default: function () {
                    return {
                        'SetName': '1',
                        'SetID': '1'
                    }
                }
            },
            paging: {
                default () {
                    return {
                        current: 1,
                        showItem: 10,
                        totalPage: 1,
                        count: 1,
                        sort: 'HostName'
                    }
                }
            },
            isLoading: {
                default: false
            },
            defaultSort: {
                type: String
            },
            hasCheckbox: {
                default: true
            },
            customize: {
                default () {
                    return {
                        name: '',
                        list: []
                    }
                }
            },
            sortable: {
                type: Boolean,
                default: true,
                required: false
            },
            hasDetail: {
                type: Boolean,
                default: true,
                required: false
            }
        },
        data () {
            return {
                defaultPageSetting: 10,
                pagelist: [{
                    value: 10,
                    label: 10
                }, {
                    value: 20,
                    label: 20
                }, {
                    value: 50,
                    label: 50
                }, {
                    value: 100,
                    label: 100
                }],
                cantBeSelected: false,
                chooseId: [],
                checked: false,
                modalAttribute: '',
                isShowRightList: true,          // 筛选框
                isDetailShow: false,            // 表格详情弹窗显示状态
                isSliderShow: false,            // 配置字段弹框
                count: 0                      // 选中主机台数
            }
        },
        computed: {
            pages: function () {
                let pag = []
                if (this.paging.current < this.paging.showItem) { // 如果当前的激活的项 小于要显示的条数
                    // 总页数和要显示的条数那个大就显示多少条
                    let i = Math.min(this.paging.showItem, this.paging.totalPage)
                    while (i) {
                        pag.unshift(i--)
                    }
                } else { // 当前页数大于显示页数了
                    let middle = this.paging.current - Math.floor(this.paging.showItem / 2) // 从哪里开始
                    let i = this.paging.showItem
                    if (middle > (this.paging.totalPage - this.paging.showItem)) {
                        middle = (this.paging.totalPage - this.paging.showItem) + 1
                    }
                    while (i--) {
                        pag.push(middle++)
                    }
                }
                return pag
            }
        },
        watch: {
            /*
                选中某个或几个主机
            */
            chooseId () {
                this.count = this.chooseId.length
                this.$emit('choose', this.chooseId, this.count)
            },
            defaultPageSetting () {
                this.$emit('resetTablePage', this.defaultPageSetting)
            }
        },
        methods: {
            /*
                表头点击排序
            */
            handleTableSortClick (sort) {
                this.$emit('handleTableSortClick', sort)
            },
            /*
                表格跨页全选
            */
            handleTableAllCheck (isChecked) {
                this.$emit('handleTableAllCheck', isChecked)
            },
            create () {
                this.isDetailShow = true
            },
            cancel () {
                this.isDetailShow = false
            },
            /*
                表格翻页
            */
            pageTuring (index) {
                this.paging.current = index
                this.$emit('pageTuring', index)
            },
            /*
                显示表格详情弹窗
            */
            showDetail (hid) {
                if (this.hasDetail) {
                    this.$emit('show-detail', hid)
                }
            },
            /*
                点击按钮
            */
            changeBind (hid) {
                this.$emit('changeBind', hid)
            },
            /*
                关闭表格弹窗详情
            */
            closePopDetail () {
                this.isDetailShow = false
            },
            selected (id, data) {
                console.log(data.name)
            },
            /*
                获取表格某项的具体信息
            */
            getTableDetail (hid) {
                let params = {
                    query_fields: [{
                        object_id: 'host',
                        fields: ['InnerIP', 'OuterIP', 'AssetID', 'SN', 'HostName', 'Description']
                    }]
                }
                this.$axios.post('model/query/' + 'host' + '/' + hid.HostID + '/', qs.stringify(params)).then((res) => {
                    if (res.result) {
                        console.log(res.data)
                    } else {
                        this.$bkInfo({
                            statusOpts: {
                                title: res['bk_error_msg'],
                                subtitle: false
                            },
                            type: 'error'
                        })
                    }
                })
            },
            /*
                获取模型属性
            */
            getModalAttribute () {
                this.$axios.get('model/attribute/Host/').then((res) => {
                    if (res.result) {
                        this.modalAttribute = res.data
                    } else {
                        this.$alertMsg(res['bk_error_msg'])
                    }
                })
            },
            /*
                计算表格最大高度，防止表格高度过高引起页面错乱
            */
            calcMaxHeight (selector) {
                const minusPX = {
                    default: 300,
                    process: 220,
                    organization: 260
                }
                let path = this.$route.fullPath.split('/')[1]
                path = minusPX.hasOwnProperty(path) ? path : 'default'
                this.$nextTick(() => {
                    this.$el.querySelector(selector).style.maxHeight = document.body.getBoundingClientRect().height - minusPX[path] + 'px'
                })
            },
            /*
                退订事件
            */
            delEvent (item) {
                this.$emit('delEvent', item)
            },
            /*
                编辑事件
            */
            editEvent (item) {
                this.$emit('editEvent', item)
            }
        },
        mounted () {
            this.calcMaxHeight('.table-scrollbar')
        },
        components: {
            vThead
        }
    }
</script>

<style media="screen">
    .noSelect {
        -webkit-touch-callout: none;
        -webkit-user-select: none;
        -khtml-user-select: none;
        -moz-user-select: none;
        -ms-user-select: none;
        user-select: none;
    }
</style>

<style media="screen" lang="scss" scoped>
    $borderColor: #e7e9ef; //边框色
    $defaultColor: #ffffff; //默认
    $primaryColor: #fafbfd; //主要
    $fnMainColor: #bec6de; //文案主要颜色
    $primaryHoverColor: #6b7baa; //鼠标移上 主要颜色
    $tableBorderColor: #dde4eb;
    .circle{
        display: inline-block;
        vertical-align: baseline;
        width: 8px;
        height: 8px;
        margin-right: 5px;
        border-radius: 50%;
        &.success{
            background: #30d878;
        }
        &.danger{
            background: #ef4c4c;
        }
    }
    .table-content{
        border: 1px solid $tableBorderColor;
        .table-scrollbar{
            width: 100%;
            overflow: auto;
            &::-webkit-scrollbar{
                width: 6px;
                height: 8px;
            }
            &::-webkit-scrollbar-thumb{
                border-radius: 20px;
                background: #a5a5a5;
                opacity: 0.8;
            }
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
                    td:first-child{
                        padding: 0 !important;
                        margin: 0 !important;
                    }
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
                padding: 4px 0;
                font-size:12px;
                color:#c3cdd7;
            }
        }
    }
    .vice-btn:hover{
        background: #fafafa;
        color: #6b7baa;
    }
</style>
