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
    <transition name="slide">
        <div class="slidebar-wrapper business-wrapper" v-show="isShow" @click.self="quickClose">
            <div class="sideslider" :style="{width: `${width}px`}">
                <slot name="title">
                    <h3 class="title">
                        <span class="vm">{{title.text}}</span>
                    </h3>
                </slot>
                <slot name="content"></slot>
            </div>
        </div>
    </transition>
</template>

<script type="text/javascript">
    export default {
        props: {
            /*
                标题
            */
            title: {
                default: function () {
                    return {
                        text: '',
                        icon: ''
                    }
                }
            },
            /*
                弹窗显示状态
            */
            isShow: {
                type: Boolean,
                default: false
            },
            /*
                弹窗宽度
            */
            width: {
                default: 800
            },
            /*
                是否支持点击空白处关闭
            */
            hasQuickClose: {
                type: Boolean,
                default: true
            },
            /*
                是否需要显示二次确认弹窗
            */
            hasCloseConfirm: {
                type: Boolean,
                default: false
            },
            /*
                二次确认弹窗显示状态
            */
            isCloseConfirmShow: {
                type: Boolean,
                default: false
            }
        },
        watch: {
            isShow (isShow) {
                if (isShow) {
                    setTimeout(() => {
                        this.$emit('shown')
                    }, 200)
                } else {
                    setTimeout(() => {
                        this.$emit('close')
                    }, 200)
                }
            }
        },
        methods: {
            closeSlider () {
                if (this.hasCloseConfirm) {
                    this.$emit('closeSlider')
                    this.$nextTick(() => {
                        if (this.isCloseConfirmShow) {
                            this.$bkInfo({
                                title: this.$t('Common["退出会导致未保存信息丢失，是否确认？"]'),
                                confirmFn: () => {
                                    this.$emit('update:isShow', false)
                                }
                            })
                        } else {
                            this.$emit('update:isShow', false)
                        }
                    })
                } else {
                    this.$emit('update:isShow', false)
                }
            },
            quickClose () {
                if (this.hasQuickClose) {
                    this.closeSlider()
                }
            }
        }
    }
</script>

<style media="screen" lang="scss" scoped>
    $primaryColor: #6b7baa;
    $lineColor: #e7e9ef;
    .slide-leave,
    .slide-enter-active {
        .sideslider{
            transition: all linear .2s;
            right: 0;
        }
    }
    .slide-enter,
    .slide-leave-active{
        .sideslider{
            right: -100%;
        }
    }
    .slidebar-wrapper{
        font-size: 14px;
        transition: all .2s;
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
        right: 0;
        background: #fff;
        box-shadow: -4px 0px 6px 0px rgba(0, 0, 0, 0.06);
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
        >.title{
            position: relative;
            padding-left: 20px;
            line-height: 60px;
            color: #333948;
            font-size: 14px;
            height: 60px;
            font-weight: bold;
            margin: 0;
            background: #f9f9f9;
            @include ellipsis;
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
        }
    }
</style>
