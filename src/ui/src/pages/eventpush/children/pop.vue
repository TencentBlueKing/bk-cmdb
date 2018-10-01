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
    <div class="pop-wrapper" v-show="isShow">
        <div class="pop-master">
            <div class="pop-box">
                <div class="title">{{$t('EventPush["推送测试"]')}}</div>
                <div class="content">
                    <p class="info">{{$t('EventPush["配置平台以POST方法推送以下示例数据到您配置的URL"]')}}</p>
                    <div class="content-box">
                        <ul>
                            <li>{</li>
                            <li class="pl20">"event_id": 100,</li>
                            <li class="pl20">"event_type": "instdata",</li>
                            <li class="pl20">"action": "update",</li>
                            <li class="pl20">"action_time": "2018-01-01 18:00 +08:00",</li>
                            <li class="pl20">"obj_type": "set",</li>
                            <li class="pl20">"cur_data": {},</li>
                            <li class="pl20">"pre_data": {},</li>
                            <li class="pl20">"dstb_id": 1,</li>
                            <li class="pl20">"dstb_time": "2018-01-01 18:00 +08:00",</li>
                            <li class="pl20">"subscription_id": 16</li>
                            <li>}</li>
                        </ul>
                    </div>
                    <div class="btn-group">
                        <bk-button type="primary" class="btn" @click="testTelnet" :loading="$loading('testPush')">{{$t('EventPush["只测试连通性"]')}}</bk-button>
                        <bk-button type="primary" class="btn" @click="testPing" :loading="$loading('testPush')">{{$t('EventPush["推送测试"]')}}</bk-button>
                        <bk-button type="default" class="btn vice-btn" @click="closePop">{{$t('Common["取消"]')}}</bk-button>
                    </div>
                </div>
            </div>
            <div class="result-info" v-show="isResultShow" v-bkloading="{isLoading: $loading('testPush')}">
                <p class="text-success" v-if="resultInfo.result"><i class="bk-icon icon-check-circle-shape"></i>{{$t('EventPush["推送成功"]')}}</p>
                <p class="text-danger" v-else><i class="bk-icon icon-close-circle-shape"></i>{{$t('EventPush["推送失败"]')}}</p>
                <template v-if="resultInfo.result">
                    <ul class="result-data" :class="{'close': !isResultOpen}" v-if="typeof(resultInfo.data)==='string'">
                        <li>{{resultInfo.data}}</li>
                    </ul>
                    <ul class="result-data" :class="{'close': !isResultOpen}" v-else>
                        <li v-for="(value, key, index) in resultInfo.data" :key="index">{{key}}：{{value}}</li>
                    </ul>
                </template>
                <template v-else-if="!resultInfo.result && resultInfo['bk_error_msg']">
                    <ul class="result-data error">
                        <li>{{resultInfo['bk_error_msg']}}</li>
                    </ul>
                </template>
                <a href="javascript:void(0)" class="result-slide" 
                    :class="{'close': !isResultOpen}"
                    v-if="resultInfo.result" 
                    @click="isResultOpen=!isResultOpen">
                    {{isResultOpen ? $t('EventPush["收起"]') : $t('EventPush["展开"]')}}
                </a>
            </div>
        </div>
    </div>
</template>

<script>
    export default {
        props: {
            callbackURL: {
                default: ''
            },
            isShow: {
                default: false,
                type: Boolean
            }
        },
        data () {
            return {
                isResultOpen: false,
                isResultShow: false,
                resultInfo: {
                    result: false,
                    bk_error_msg: '',
                    data: ''
                }
            }
        },
        watch: {
            isShow (isShow) {
                if (!isShow) {
                    this.isResultOpen = false
                }
            }
        },
        methods: {
            /*
                测试推送
            */
            testPing () {
                this.isResultShow = true
                this.$axios.post(`event/subscribe/ping`, {
                    callback_url: this.callbackURL
                }, {id: 'testPush'}).then(res => {
                    this.resultInfo = res
                })
            },
            /*
                测试连通性
            */
            testTelnet () {
                this.isResultShow = true
                this.$axios.post(`event/subscribe/telnet`, {
                    callback_url: this.callbackURL
                }, {id: 'testPush'}).then(res => {
                    this.resultInfo = res
                })
            },
            closePop () {
                this.isResultShow = false
                this.$emit('closePop')
            }
        }
    }
</script>

<style lang="scss" scoped>
    p{
        margin: 0;
    }
    .pop-wrapper{
        position: fixed;
        top: 154px;
        right: 102px;
        box-shadow: 0px 2px 9.6px 0.4px rgba(12, 34, 59, 0.13);
        background: #fff;
        width: 595px;
        .pop-box{
            padding: 20px;
            .title{
                font-weight: bold;
                line-height: 1;
            }
            .content{
                padding: 20px 20px 0;
                .info{
                    line-height: 1;
                }
                .content-box{
                    margin: 10px 0 20px;
                    padding: 20px;
                    background: #455070;
                    color: #bfc6e0;
                    border-radius: 2px;
                    font-size: 12px;
                    line-height: 20px;
                    height: 203px;
                    overflow-y: auto;
                    &::-webkit-scrollbar{
                        width: 6px;
                        height: 5px;
                    }
                    &::-webkit-scrollbar-thumb{
                        border-radius: 20px;
                        background: #a5a5a5;
                    }
                }
                .btn-group{
                    text-align: right;
                    font-size: 0;
                    .btn{
                        margin-left: 11px;
                        width: 112px;
                        height: 34px;
                        padding: 0;
                        line-height: 32px;
                        text-align: center;
                        &.cancel:hover{
                            color: #6b7baa;
                            border: 1px solid #6b7baa;
                        }
                    }
                }
            }
        }
        .result-info{
            border-top: 1px solid #e7e8e8;
            background: #f9f9f9;
            padding: 17px 40px;
            font-size: 14px;
            line-height: 16px;
            .bk-icon{
                font-size: 16px;
                margin-right: 7px;
                vertical-align: -2px;
            }
            .text-success{
                color: #34d97b;
            }
            .text-danger{
                color: #ef4c4c;
            }
            .result-data{
                line-height: 18px;
                font-size: 12px;
                padding: 9px 23px 0;
                max-height: 81px;
                overflow: auto;
                &.close{
                    display: none;
                }
                &.error{
                    color: #ef4c4c;
                }
                &::-webkit-scrollbar {
                    width: 6px;
                    height: 5px;
                }
                &::-webkit-scrollbar-thumb {
                    border-radius: 20px;
                    background: #a5a5a5;
                    -webkit-box-shadow: inset 0 0 6px hsla(0,0%,80%,.3);
                }
            }
            .result-slide{
                display: block;
                text-align: center;
                text-decoration: none;
                color: #6b7baa;
                font-size: 12px;
                line-height: 12px;
                margin-top: 7px;
                &.close{
                    margin-top: -14px;
                    &:before{
                        transform: rotate(180deg);
                    }
                }
                &:before{
                    content: '';
                    display: inline-block;
                    width: 11px;
                    height: 10px;
                    margin-right: 9px;
                    background-image: url('../../../common/images/icon/icon-result-slide.png');
                }
            }
        }
    }
</style>
