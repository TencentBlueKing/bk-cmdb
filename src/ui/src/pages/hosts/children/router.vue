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
    <div class="tab-content">
        <div class="attribute-list">
            <div class="title clearfix">
                <h3 class="fl">/etc/hosts</h3>
                <div class="content fr clearfix">
                    <div class="info-title">{{$t('HostResourcePool[\'最近更新时间\']')}}：</div>
                    <div class="info-detail">{{getHostSnapshot.upTime}}</div>
                </div>
            </div>
        </div>
        <div class="host-content">
            <div class="host-show">
                <ul>
                    <li v-for="(item, index) in getHostSnapshot.route" :key="index">
                        {{item}}
                    </li>
                </ul>
            </div>
        </div>
    </div>
</template>

<script>
    import {mapGetters} from 'vuex'
    export default {
        computed: {
            ...mapGetters([
                'getHostSnapshot'
            ])
        }
    }
</script>

<style lang="scss" scoped>
    .tab-content{
        padding-top: 62px!important;
        font-size: 14px;
        padding-left:40px;
        height: 100%;
        overflow: auto;
        &::-webkit-scrollbar{
            width: 6px;
            height: 5px;
        }
        &::-webkit-scrollbar-thumb{
            border-radius: 20px;
            background: #a5a5a5;
        }
        .attribute-list{
            margin-bottom:20px;
            .title{
                h3{
                    color:#bec6de;
                    font-size:14px;
                    font-weight: bold;
                    line-height:1;
                    margin: 0;
                }
                .content{
                    font-size: 12px;
                    padding-right: 30px;
                    .info-title{
                        float: left;
                        color: #bec6de;
                    }
                    .info-detail{
                        float: left;
                        color: #6b7baa;
                    }
                }
            }
        }
        .host-content{
            margin-right: 30px;
            .host-show{
                width:100%;
                min-height:13px;
                background:#454f6f;
                padding:6px 20px 20px 20px;
                ul{
                    li{
                        line-height:1;
                        padding-top:14px;
                        color:#bec6de;
                        word-break: break-all;
                        word-wrap: break-word;
                    }
                }
            }
        }
    }
</style>
