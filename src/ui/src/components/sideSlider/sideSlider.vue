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
    <div class="slidebar-wrapper" v-show="isShow">
        <div class="sideslider" :style="'width: ' + limitTheWidth + 'px'">
            <h3 class="title">
                <i :class="title.icon" class="vb"></i><span class="vm">{{title.text}}</span>
                <span class="close" @click="cancel">关闭</span>
            </h3>
            <div class="content clearfix">
                <div class="left-list search-wrapper-hidden">
                    <div class="search-wrapper">
                        <div class="the-host">
                            隐藏属性
                        </div>
                        <div class="search-field">
                            <input type="text" name="" value="" placeholder="搜索属性" v-model.trim="searchText">
                        </div>
                    </div>

                    <ul class="list-wrapper">
                        <li v-for="(item, index) in forSelectionList" @click="addItem(index)">
                            {{item.PropertyName}}
                            <i class="bk-icon icon-angle-right"></i>
                        </li>
                    </ul>
                </div>
                <div class="right-list pr20">
                    <div class="title">
                        <div class="search-wrapper">
                            已显示属性
                        </div>
                    </div>
                    <div :class="['model-content', {'content-left-hidden' : isShow}]" >
                        <bk-tab :active-name="'curveSetting'" class="content-left">
                            <bk-tabpanel name="curveSetting" title="常规">
                                <slot name="contentRight"></slot>
                            </bk-tabpanel>
                            <bk-tabpanel name="exceptionDetecting" title="关联">
                                <slot name="contentRight"></slot>
                            </bk-tabpanel>
                        </bk-tab>
                        <!-- <ul class="content-left">
                            <li class="active">常规</li>
                            <li>关联</li>
                            <li class="add">
                                <span class="fl f10 plus">
                                    <i class="bk-icon icon-plus"></i>
                                </span>
                                <span class="fr f10 edit">
                                    <i class="bk-icon icon-edit"></i>
                                </span>
                            </li>
                        </ul> -->
                        <div slot="contentRight">
                            <!-- <ul class="content-right" id="sort-wrapper"> -->
                                <!-- <li v-for="(item, index) in hasSelectionList">
                                    <i class="icon-triple-dot"></i><span>{{item.PropertyName}}</span><i class="bk-icon icon-close" @click="removeItem(index)"></i>
                                </li> -->
                                <!-- <li
                                    v-for="(item, index) in hasSelectionList"
                                    :key="item.PropertyId"
                                    v-dragging="{item: item, list: hasSelectionList, group: 'item'}">
                                        <i class="icon-triple-dot">
                                    </i><span>{{item.PropertyName}}</span><i class="bk-icon icon-close" @click="removeItem(index)"></i>
                                </li> -->
                                <draggable class="content-right" v-model="hasSelectionList" :options="{animation: 150}">
                                    <div v-for="(item, index) in hasSelectionList" :key="index" class="item">
                                        <i class="icon-triple-dot"></i><span>{{item.PropertyName}}</span><i class="bk-icon icon-close" @click="removeItem(index)"></i>
                                    </div>
                                </draggable>
                                <!-- <li v-for="(item, index) in hasSelectionList">
                                    <i class="icon-triple-dot"></i><span>{{item.PropertyName}}</span><i class="bk-icon icon-close" @click="removeItem(index)"></i>
                                </li> -->
                            <!-- </ul> -->
                        </div>
                    </div>
                </div>
            </div>
            <div class="bk-form-item bk-form-action content-button">
                <bk-button class="btn apply" type="primary" title="应用" @click="apply">
                    应用
                </bk-button>
                <bk-button class="btn reinstate cancel" type="default" title="取消" @click="cancel">
                    取消
                </bk-button>
                <!-- 功能有问题 -->
                <!-- <bk-button class="btn reinstate cancel" type="default" title="恢复默认" @click="reset">
                    恢复默认
                </bk-button> -->
            </div>
        </div>
    </div>
</template>

<script type="text/javascript">
    import $ from 'jquery'
    import draggable from 'vuedraggable'
    export default {
        data () {
            return {
                myArray: [
                    {
                        id: 1,
                        name: 'aaa'
                    },
                    {
                        id: 2,
                        name: 'bbb'
                    },
                    {
                        id: 3,
                        name: 'ccc'
                    },
                    {
                        id: 4,
                        name: 'ddd'
                    }
                ],
                searchText: '',
                forSelectionList: [],
                hasSelectionList: [],
                curSelectionList: [],
                filterClassify: {
                    id: '',
                    label: ''
                },
                classifyList: [],
                curTempSelectionList: []
            }
        },
        watch: {
            isShow: function (newVal) {
                if (newVal) {
                    $('.slidebar-wrapper').addClass('sideSliderShow')
                    setTimeout(() => {
                        $('.sideslider').addClass('slideIn')
                    }, 0)
                    setTimeout(() => {
                        $('.slidebar-wrapper').removeClass('sideSliderShow')
                    }, 500)
                }
            },
            forSelection () {
                this.forSelectionList = this.forSelection
                this.curTempSelectionList = this.forSelectionList.concat()
            },
            searchText: function (val) {
                this.forSelectionList = []
                let searchText = this.searchText
                if (searchText === '') {  // 搜索条件为空时 当前收藏列表为全部收藏列表
                    this.forSelectionList = this.curTempSelectionList
                    return
                }
                for (var i = 0, list = this.curTempSelectionList; i < list.length; i++) {
                    if (list[i].PropertyName.indexOf(searchText) !== -1) {
                        this.forSelectionList.push(list[i])
                    }
                }
            },
            'filterClassify.id' () {
                this.getClassifyList(this.filterClassify.id)
            },
            tableHeader (val) {
                this.tableHeader = this.tableHeader ? this.tableHeader : {}
                this.hasSelectionList = JSON.parse(JSON.stringify(this.tableHeader))
            }
        },
        props: {
            title: {
                default: function () {
                    return {
                        text: '主机筛选项设置',
                        icon: 'icon-cc-list'
                    }
                }
            },
            isShow: {
                default: false
            },
            /*
                滑动框显示方式 默认slide 层中层的话用direct
            */
            showType: {
                default: 'slide'
            },
            forSelection: {
                default: function () {
                    return {
                        'ObjId': 'Host'
                    }
                }
            },
            tableHeader: {
                default: function () {
                    return {
                        'PropertyName': '123'
                    }
                }
            },
            hasReset: {
                default: false
            },
            limitTheWidth: {
                default: 800
            }
        },
        methods: {
            selected () {

            },
            /*
                添加到已选
                index: 当前项的index
            */
            addItem (index) {
                this.hasSelectionList = this.hasSelectionList.concat(this.forSelectionList.splice(index, 1))
            },
            /*
                删除已选
                index: 当前项的index
            */
            removeItem (index) {
                if (this.hasSelectionList.length === 1) {
                    this.$alertMsg('至少选择一个字段')
                } else {
                    this.forSelectionList = this.forSelectionList.concat(this.hasSelectionList.splice(index, 1))
                }
            },
            /*
                取消 收起弹框
            */
            cancel () {
                $('.slidebar-wrapper').addClass('sideSliderShow')
                $('.sideslider').removeClass('slideIn')
                setTimeout(() => {
                    this.$emit('cancel')
                    $('.slidebar-wrapper').removeClass('sideSliderShow')
                }, 500)
            },
            /*
                恢复默认
            */
            reset () {
                this.$emit('reset')
            },
            /*
                应用
            */
            apply () {
                let hasSelectionList = JSON.parse(JSON.stringify(this.hasSelectionList))
                this.$emit('apply', hasSelectionList)
                this.cancel()
            },
            getClassifyList (index) {
                let params = {}
                if (index !== '' || index !== null) {
                    params = {
                        'ClassificationType': 'inner',
                        'ClassificationName': '',
                        'ClassificationId': ''
                    }
                }
                this.$axios.post('object/classification/0/objects', params).then((res) => {
                    if (res.result) {
                        this.classifyList = res.data
                    } else {
                        this.$alertMsg(res['bk_error_msg'])
                    }
                })
            }
        },
        components: {
            draggable
        }
    }
</script>
<style media="screen" lang="scss" scoped>
    $primaryColor: #6b7baa;
    $lineColor: #e7e9ef;
    .slidebar-wrapper{
        font-size: 14px;
        position: fixed;
        -webkit-transition: all 0.5s;
        -moz-transition: all 0.5s;
        -ms-transition: all 0.5s;
        -o-transition: all 0.5s;
        transition: all 0.5s;
        left: 0;
        top: 0;
        width: 100%;
        height: 100%;
        background: rgba(0, 0, 0, 0.3);
        z-index: 1300;
        &.sideSliderShow{
            overflow: hidden;
        }
        .sideslider{
            transition: all .2s linear;
            position: absolute;
            top: 0;
            bottom: 0;
            right: -830px;
            width: 800px;
            background: #fff;
            box-shadow: -4px 0px 6px 0px rgba(0, 0, 0, 0.06);
            &.slideIn{
                right: 0;
            }
            >.title{
                padding-left: 20px;
                line-height: 60px;
                color: #4d597d;
                font-size: 14px;
                height: 60px;
                font-weight: bold;
                margin: 0;
                background: #f9f9f9;
                .close{
                    position: absolute;
                    top: 0;
                    left: -30px;
                    width: 30px;
                    height: 60px;
                    padding: 10px 7px 0;
                    background-color: #ef4c4c;
                    box-shadow: -2px 0 2px 0 rgba(0, 0, 0, 0.2);
                    cursor: pointer;
                    color: #fff;
                    font-size: 14px;
                    line-height: 20px;
                    font-weight: normal;
                    &:hover{
                        background: #e13d3d;
                    }
                }
                >i{
                    position: relative;
                    top: 2px;
                    margin-right: 5px;
                    display: inline-block;
                }
                .icon-mainframe{
                    width: 12px;
                    height: 15px;
                }
            }
            >.content{
                height: calc(100% - 122px);
                border-top: 1px solid #e7e9ef;
                .left-list{
                    float: left;
                    width: 50%;
                    height: 100%;
                    border-right: 1px solid #e7e9ef;
                    .list-wrapper{
                        height: calc(100% - 78px);
                        overflow: auto;
                        padding: 15px 0 0 0;
                        &::-webkit-scrollbar{
                        width: 6px;
                        height: 5px;
                        }
                        &::-webkit-scrollbar-thumb{
                            border-radius: 20px;
                            background: #a5a5a5;
                        }
                    }
                    &.search-wrapper-hidden{
                        .search-wrapper{
                            .text{
                                display:none;
                            }
                            .search{
                                display:none;
                            }
                            .the-host{
                                width:122px;
                                display:inline-block;
                                line-height: 36px;
                            }
                            .search-field{
                                width:120px;
                                // display:inline-block;
                                float: right;
                                input{
                                    border:1px solid #e7e9ef;
                                    width:100%;
                                    height:36px;
                                    line-height:36px;
                                    outline:none;
                                    padding:0 15px;
                                }

                            }
                        }
                    }
                    .search-wrapper{
                        width: 100%;
                        height: 78px;
                        padding: 20px;
                        .select-box{
                            float: left;
                            width: 163px;
                            height: 37px;
                            margin-right: 10px;
                            &.open{
                                .bk-selector-icon{
                                    top: 17px;
                                }
                            }
                        }
                        .search{
                            float: left;
                            input{
                                width: 131px;
                                height: 37px;
                                border: 1px solid #e7e9ef;
                                border-radius: 2px;
                                padding: 0 12px;
                                font-size: 14px;
                                color: #bec6de;
                            }
                            &.search2{
                                float: right;
                                input{
                                    width: 180px;
                                }
                            }
                        }
                        .text{
                            float: left;
                            line-height: 37px;
                            margin-left: 9px;
                        }
                    }
                    .list-wrapper{
                        border-top: 1px solid #e7e9ef;
                        li{
                            height: 42px;
                            line-height: 42px;
                            color: $primaryColor;
                            font-size: 14px;
                            padding-left: 27px;
                            cursor: pointer;
                            &:hover{
                                background: #f9f9f9;
                            }
                            i{
                                float: right;
                                margin-top: 12px;
                                margin-right: 18px;
                                color: #bec6de;
                            }
                        }
                    }
                }
                .right-list{
                    float: left;
                    height: 100%;
                    width: 50%;
                    .bk-tab2{
                        border: none !important;
                    }
                    .content-left-hidden{
                        .content-left{
                            display:none;
                        }
                        .content-right{
                            width:100% !important;
                            padding: 15px 0 0 0;
                            .content-center{
                                width:232px;
                                 margin-right:36px;
                                .input-number{
                                    width:100%;
                                    height:36px;
                                    line-height:36px;
                                    outline:none;
                                    border:1px solid #e7e9ef;
                                    background:#f9f9f9;
                                    color:#bec6de;
                                    padding:0 15px;
                                    &.disbale{
                                        cursor:not-allowed;
                                    }
                                    &::-webkit-input-placeholder{
                                        font-family: "Microsoft YaHei";
                                        color:#c3cdd7;
                                    }
                                    &:-moz-placeholder{
                                        font-family: "Microsoft YaHei";
                                        color:#c3cdd7;
                                    }
                                    &::-moz-placeholder{
                                        font-family: "Microsoft YaHei";
                                        color:#c3cdd7;
                                    }
                                    &:-ms-input-placeholder{
                                        font-family: "Microsoft YaHei";
                                        color:#c3cdd7;
                                    }
                                }
                            }
                        }
                    }
                    .title{
                        height: 79px;
                        line-height: 78px;
                        padding: 0 20px;
                        width: 430px;
                        border-bottom: 1px solid #e7e9ef;
                        .list-wrapper{
                            border-top: 1px solid #e7e9ef;
                        }
                    }
                    .model-content{
                        color: $primaryColor;
                        height: calc(100% - 79px);
                        overflow: auto;
                        &::-webkit-scrollbar{
                        width: 6px;
                        height: 5px;
                        }
                        &::-webkit-scrollbar-thumb{
                            border-radius: 20px;
                            background: #a5a5a5;
                        }
                        .content-left{
                            float: left;
                            width: 108px;
                            height: 100%;
                            text-align: center;
                            background: #f9f9f9;
                            // padding: 15px 0 0 0;
                            li{
                                height: 43px;
                                line-height: 42px;
                                border-bottom: 1px solid #fff;
                                &.active{
                                    background: #fff;
                                }
                                &.add{
                                    cursor: pointer;
                                    .plus{
                                        width:55px;
                                        border-right: 1px solid #fff;
                                    }
                                    .edit{
                                        width:53px;
                                        text-align:center;
                                    }
                                    i{
                                        color:#bec6de;
                                    }
                                }
                            }
                        }
                        .content-right{
                            // float: left;
                            width: calc(100% - 108px);
                            >.item{
                                height: 43px;
                                line-height: 42px;
                                padding-left: 30px;
                                cursor: move;
                                &:hover{
                                    background: #f9f9f9;
                                }
                                .icon-triple-dot{
                                    position: relative;
                                    top: -1px;
                                    display: inline-block;
                                    width: 4px;
                                    height: 14px;
                                    margin-right: 10px;
                                    background: url(../../common/images/icon/icon-triple-dot.png);
                                }
                                .icon-close{
                                    float: right;
                                    font-size: 12px;
                                    margin-top: 5px;
                                    // margin-right: 20px;
                                    padding: 10px;
                                    cursor: pointer;
                                    &:hover{
                                        color: red
                                    }
                                }
                            }
                        }
                    }
                }
            }
        }
        .content-button{
            background: #f9f9f9;
            height: 62px;
            padding: 14px 20px;
            font-size: 0;
            .btn{
                font-size: 14px;
                width: 110px;
                height: 34px;
                line-height: 34px;
                border-radius: 0;
                margin-right: 10px;
                border: 0;
                border-radius: 2px;
                &.apply{
                    &:hover{
                        background: #4d597d;
                    }
                }
                &.cancel{
                    border: 1px solid #e6e9f2;
                    &:hover{
                        color: #6b7baa;
                        border-color: #6b7baa;

                    }
                }

            }
            .info{
                float: right;
                font-size: 14px;
                height: 34px;
                line-height: 34px;
                cursor: pointer;
                input{
                    position: relative;
                    margin-right: 4px;
                }
            }
        }
    }
    .bk-tab2 .bk-tab2-head .bk-tab2-nav > li{
        width: 108px !important;
    }
</style>
