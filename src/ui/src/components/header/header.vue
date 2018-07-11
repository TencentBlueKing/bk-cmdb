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
    <div class="header-wrapper clearfix">
        <img v-if="$i18n.locale === 'zh_CN'" class="logo" src="@/common/svg/logo-cn.svg" alt="蓝鲸配置平台" @click="gotoIndex">
        <img v-else class="logo" src="@/common/svg/logo-en.svg" alt="Configuration System" @click="gotoIndex">
        <div class="header-right-contain fr">
            <div class="user-detail-contain fr pr">
                <div class="dropdown-content-user fl">
                    <div class="select-trigger">
                        <span class="f14">{{isAdmin == 1 ? userName + '（' + $t("Common['管理员']") + '）' : userName}}</span>
                        <i class="bk-icon icon-angle-down"></i>
                        <ul class="select-content">
                            <li @click="logOut">
                                <i class="icon-cc-logout"></i>{{$t("Common['注销']")}}
                            </li>
                        </ul>
                    </div>
                </div>
            </div>
        </div>
    </div>
</template>

<script type="text/javascript">
    import bus from '@/eventbus/bus'
    import Cookies from 'js-cookie'
    import {mapGetters} from 'vuex'
    export default {
        data () {
            return {
                userName: '',
                isAdmin: 0
            }
        },
        methods: {
            /*
                登出
            */
            logOut () {
                window.location.href = `${window.siteUrl}logout`
            },
            gotoIndex () {
                this.$router.push('/index')
            }
        },
        updated () {
            this.userName = window.userName
            this.isAdmin = window.isAdmin
        },
        created () {
            this.userName = window.userName
            this.isAdmin = window.isAdmin
        }
    }
</script>
<style lang="scss" scoped>
    $borderColor: #e7e9ef; //边框色
    $primaryColor: #737987;
    .header-wrapper{
        height: 60px;
        background: #fbfbfb;
        border-bottom: 1px solid $borderColor;
        font-size: 14px;
        position: absolute;
        left: 0;
        top: 0;
        width: 100%;
        padding: 0 0 0 60px;
        z-index: 1200;
        .logo{
            height: 27px;
            margin: 8px;
            cursor: pointer;
        }
        .header-right-contain{
            padding: 4px 20px 0 0;
            .select-content{
                position: absolute;
                top: 38px;
                right: 0;
                min-width: 100px;
                padding: 10px 0;
                background-color: #fff;
                box-shadow: 0px 1px 5px 0px rgba(12, 34, 59, .1);
                z-index: 10;
                li{
                    list-style:none;
                    &:first-child{
                        cursor: default;
                    }
                    a{
                        color: $primaryColor;
                    }
                    display: block;
                    height: 45px;
                    padding-left: 12px;
                    line-height: 45px;
                    text-decoration: none;
                    font-size: 14px;
                    color: $primaryColor;
                    background-color: #fff;
                    i{
                        margin-right: 8px;
                    }
                    &:hover,
                    &.active{
                        background-color: #f1f7ff;
                        color: #498fe0;
                    }
                }
            }
            .user-detail-contain{
                height:36px;
                line-height:36px;
                padding: 0;
                margin-left: 25px;
                .dropdown-content-user{
                    cursor: pointer;
                    &:hover{
                        .select-content{
                            display:block;
                        }
                    }
                }
                .select-trigger{
                    // color:$primaryColor;
                    &:hover{
                        color: #498fe0;
                    }
                    span{
                        letter-spacing:1px;
                        padding:20px 0 15px 8px;
                    }
                    .bk-icon.icon-angle-down{
                        font-size: 12px;
                        margin-left: 2px;
                    }
                }
                .avatar-contain{
                    width:34px;
                    height:34px;
                    border-radius:50%;
                    img{
                        width:100%;
                        height:100%;
                    }
                }
                .select-content{
                    display: none;
                }
            }
        }
    }
</style>
