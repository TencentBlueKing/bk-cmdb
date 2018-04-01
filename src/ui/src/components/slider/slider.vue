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
    <div class="slidebar-wrapper business-wrapper" v-show="isSliderShow">
        <div class="sideslider">
            <h3 class="title">
                <i class="mr10" :class="title.icon"></i><span class="vm">{{title.text}}</span>
                <span class="close" @click="closeSlider">关闭</span>
            </h3>
            <slot name="btn"></slot>
            <slot name="content"></slot>
            <slot name="footer"></slot>
        </div>
    </div>
</template>

<script type="text/javascript">
    import $ from 'jquery'
    export default {
        data () {
            return {
                isSliderShow: false
            }
        },
        props: {
            title: {
                default: function () {
                    return {
                        text: '业务(模型管理)',
                        icon: 'icon-cc-list'
                    }
                }
            },
            isShow: {
                default: false
            }
        },
        watch: {
            isShow: function (newVal) {
                if (newVal) {
                    this.showSlider()
                } else {
                    this.isSliderShow = false
                    // this.closeSlider()
                }
            }
        },
        methods: {
            showSlider () {
                this.isSliderShow = true
                $('.slidebar-wrapper').addClass('sideSliderShow')
                setTimeout(() => {
                    $('.sideslider').addClass('slideIn')
                }, 0)
                setTimeout(() => {
                    $('.slidebar-wrapper').removeClass('sideSliderShow')
                }, 300)
            },
            closeSlider () {
                $('.slidebar-wrapper').addClass('sideSliderShow')
                $('.sideslider').removeClass('slideIn')
                // setTimeout(() => {
                this.$emit('closeSlider')
                    // $('.slidebar-wrapper').removeClass('sideSliderShow')
                // }, 300)
            }
        },
        mounted () {

        }
    }
</script>

<style media="screen" lang="scss" scoped>
    $primaryColor: #6b7baa;
    $lineColor: #e7e9ef;
    .slidebar-wrapper{
        font-size: 14px;
        -webkit-transition: all 0.5s;
        -moz-transition: all 0.5s;
        -ms-transition: all 0.5s;
        -o-transition: all 0.5s;
        transition: all 0.5s;
        width: 100%;
        height: 100%;
        z-index: 1300;
        position: fixed;
        top: 0;
        right: 0;
        bottom: 0;
        left: 0;
        background-color: rgba(0,0,0,0.6);
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
            position: relative;
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
            .icon-mainframe{
                position: relative;
                top: 0px;
                margin-right: 5px;
                display: inline-block;
                width: 19px;
                height: 19px;
            }
        }
        >.content{
            height: calc(100% - 122px);
            border-top: 1px solid #e7e9ef;
            .bk-tab2{
                .bk-tab2-content{
                height: calc(100% - 58px) !important;
                overflow-y:auto;
                overflow-x:hidden;
                }
            }
        }
    }
</style>
