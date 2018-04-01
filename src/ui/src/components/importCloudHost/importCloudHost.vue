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
    <div class="slidebar-wrapper" v-show="isShowCloud">
        <div class="sideslider">
            <h3 class="title">
                <i :class="title.icon"></i>{{title.text}}
            </h3>
            <div class="content import-cloud-contain">
                <ul class="list clearfix">
                    <li></li>
                    <li></li>
                    <li></li>
                </ul>
            </div>
            <div class="bk-form-item bk-form-action content-button">
                <bk-button class="btn fr vice-btn" type="default" title="取消" @click="cancelPop">
                    取消
                </bk-button>
            </div>
        </div>
    </div>
</template>

<script type="text/javascript">
    import $ from 'jquery'
    export default {
        data () {
            return {

            }
        },
        watch: {
            isShowCloud: function (newVal) {
                if (newVal) {
                    $('.slidebar-wrapper').addClass('sideSliderShow')
                    setTimeout(() => {
                        $('.sideslider').addClass('slideIn')
                    }, 0)
                    setTimeout(() => {
                        $('.slidebar-wrapper').removeClass('sideSliderShow')
                    }, 500)
                }
            }
        },
        props: {
            title: {
                default: function () {
                    return {
                        text: '导入主机',
                        icon: 'icon-cloud'
                    }
                }
            },
            isShowCloud: {
                default: false
            }
        },
        methods: {
            /*
                取消 收起弹框
            */
            cancelPop () {
                $('.slidebar-wrapper').addClass('sideSliderShow')
                $('.sideslider').removeClass('slideIn')
                setTimeout(() => {
                    this.$emit('cancelPop')
                    $('.slidebar-wrapper').removeClass('sideSliderShow')
                }, 500)
            },
            /*
                tap切换回调函数
            */
            tabChanged () {

            }
        }
    }
</script>

<style media="screen" lang="scss" scoped>
    $primaryColor: #6b7baa;
    $lineColor: #e7e9ef;
    $borderColor: #e7e9ef; //边框色
    $defaultColor: #ffffff; //默认
    $primaryColor: #f9f9f9; //主要
    $fnMainColor: #bec6de; //文案主要颜色
    $primaryHoverColor: #6b7baa; // 主要颜色
    .main-btn{  //主要按钮
        background: $primaryHoverColor;
        &:hover{
            background: #4d597d;
        }
    }
    .vice-btn{  //次要按钮 取消按钮
        border: 1px solid #e6e9f2;
        color:  $primaryHoverColor;
        cursor: pointer;
        &:hover{
            border-color: $primaryHoverColor;
        }
    }
    .icon-btn{  //单纯图标的按钮
        background: #ffffff;
        color: $primaryHoverColor;
        cursor: pointer;
        &:hover{
            background: $primaryHoverColor;
            color: $defaultColor;
        }
    }
    .no-border-btn{    //无边框按钮
        background: #fff;
        color: $primaryHoverColor;
        cursor: pointer;
        &:hover{
            background: $primaryHoverColor;
            color: #fff;
        }
    }
    .slidebar-wrapper{
        position: relative;
        font-size: 14px;
        position: absolute;
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
        z-index: 1230;
        &.sideSliderShow{
            overflow: hidden;
        }
        .sideslider{
            transition: all .2s linear;
            position: absolute;
            top: 0;
            bottom: 0;
            right: -800px;
            width: 800px;
            background: #fff;
            box-shadow: -4px 0px 6px 0px rgba(0, 0, 0, 0.06);
            &.slideIn{
                right: 0;
            }
            >.title{
                padding-left: 20px;
                line-height: 60px;
                color: #212232;
                font-size: 16px;
                height: 60px;
                font-weight: normal;
                margin: 0;
                background: #f9f9f9;
                .icon-cloud{
                    position: relative;
                    top: 2px;
                    margin-right: 5px;
                    display: inline-block;
                    width: 16px;
                    height: 15px;
                }
            }
            >.content{
                height: calc(100% - 122px);
                border-top: 1px solid #e7e9ef;
                &.import-cloud-contain{
                    padding:35px;
                    ul{
                        li{
                            float: left;
                            width: 196px;
                            height: 196px;
                            line-height:1;
                            margin-right:40px;
                            border:1px solid $lineColor;
                            cursor:pointer;
                            &:nth-child(1) {
                                background: url("../../common/images/tencent-cloud.png") no-repeat center center;
                                &:hover{
                                    background: url("../../common/images/tencent-cloud-on.png") no-repeat center center;
                                    box-shadow: 0px 0px 5px #f0f1f3;
                                }
                            }
                            &:nth-child(2) {
                                background: url("../../common/images/ali-cloud.png") no-repeat center center;
                                &:hover{
                                    background: url("../../common/images/ali-cloud-on.png") no-repeat center center;
                                    box-shadow: 0px 0px 5px #f0f1f3;
                                }
                            }
                            &:nth-child(3) {
                                background: url("../../common/images/aws.png") no-repeat center center;
                                margin-right:0;
                                &:hover{
                                    background: url("../../common/images/aws-on.png") no-repeat center center;
                                    box-shadow: 0px 0px 5px #f0f1f3;
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
        .bk-tab2{
            border:none;
        }
    }
</style>
