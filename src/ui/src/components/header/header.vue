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
        <div class="header-right-contain fr">
            <!-- 导航快速搜索 -->
            <div class="fl clearfix" v-click-outside="resetSearch">
                <transition name="quick-search" @afterEnter="quickSearchTextFocus">
                    <div class="quick-search-contain fl"  v-show="isShowQuickSearch">
                        <div class="dropdown-content-ip fl" 
                            @mouseover="changeSearchTargetListVisible(true)"
                            @mouseout="changeSearchTargetListVisible(false)"
                        >
                            <div class="select-trigger-box pr" @click.stop>
                                <div class="select-trigger pr">
                                    <span class="f14">{{currentSearchTarget.name}}</span>
                                    <i class="bk-icon icon-down-shape f10"></i>
                                    <ul class="select-content" v-show="searchTargetListVisible">
                                        <li v-for="searchTarget in searchTargetList" @click.stop="setSearchTarget(searchTarget)">
                                            <a :class="{'active': currentSearchTarget.id === searchTarget.id}" href="javascript:void(0)">{{searchTarget.name}}</a>
                                        </li>
                                    </ul>
                                </div>
                            </div>
                        </div>
                        <div class="search-content fl" @click.stop.prevent>
                            <input ref="quickSearchText" type="text" name="" value="" :placeholder="`${$t('Common[\'快速查询\']')}...`" class="search-input" v-model.trim="searchText" @keyup.enter="quickSearch">
                        </div>
                    </div>
                </transition>
                <div class="fl quick-search-icon" :class="{'show': !isShowQuickSearch}" @click.stop.prevent="quickSearch" @mouseover="showQuickSearch"><i class="bk-icon icon-search"></i></div>
            </div>
            <div class="language fl" hidden>   
                <i class="icon icon-cc-lang"></i>
                <span class="language-text">{{languageLable}}</span>
                <i class="bk-icon icon-angle-down"></i>
                <ul class="language-box">
                    <li :class="{'active': language==='zh_CN'}" @click="changeLanguage('zh_CN')">简体中文</li>
                    <li :class="{'active': language==='en'}" @click="changeLanguage('en')">English</li>
                </ul>
            </div>
            <div class="user-detail-contain fr pr">
                <div class="dropdown-content-user fl">
                    <div class="select-trigger">
                        <span class="f14">{{userName}}</span>
                        <i class="bk-icon icon-angle-down"></i>
                        <ul class="select-content">
                            <li v-if="isAdmin == 1">
                                <i class="icon-cc-user"></i>{{$t("Common['管理员']")}}
                            </li>
                            <li v-else>
                                <i class="icon-cc-user"></i>{{$t("Common['普通用户']")}}
                            </li>
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
    export default {
        data () {
            return {
                userName: '',
                isAdmin: 0,
                languageLable: '中文',
                searchText: '',
                searchTargetListVisible: false,
                currentSearchTarget: {id: 'ip', name: 'IP'},
                searchTargetList: [{
                    id: 'ip',
                    name: 'IP'
                }],
                isShowQuickSearch: false           // 搜索内容的展示
            }
        },
        watch: {
            searchText (searchText) {
                this.$store.commit('setQuickSearchParams', {
                    text: searchText,
                    type: this.currentSearchTarget.id
                })
            },
            currentSearchTarget (searchTarget) {
                this.$store.commit('setQuickSearchParams', {
                    text: this.searchText,
                    type: searchTarget.id
                })
            }
        },
        methods: {
            /*
                登出
            */
            logOut () {
                window.location.href = `${window.siteUrl}logout`
            },
            changeSearchTargetListVisible (isVisible) {
                this.searchTargetListVisible = isVisible
            },
            setSearchTarget (searchTarget) {
                this.currentSearchTarget = searchTarget
                this.searchTargetListVisible = false
            },
            quickSearch () {
                if (this.searchText.length) {
                    bus.$emit('quickSearch')
                    this.$router.push('/hosts')
                }
            },
            showQuickSearch () {
                this.isShowQuickSearch = true
            },
            quickSearchTextFocus () {
                this.$refs.quickSearchText.focus()
            },
            resetSearch () {
                this.isShowQuickSearch = false
            },
            changeLanguage (language) {
                this.language = language
                this.$i18n.locale = language
                this.$store.commit('setLang', language)
                this.$validator.localize(language)
                if (language === 'zh_CN') {
                    this.languageLable = '中文'
                    this.setLang('zh')
                } else if (language === 'en') {
                    this.languageLable = 'EN'
                    this.setLang('en')
                }
            }
        },
        updated () {
            this.userName = window.userName
            this.isAdmin = window.isAdmin
        },
        created () {
            this.userName = window.userName
            this.isAdmin = window.isAdmin
            bus.$on('setQuickSearchParams', ({type, text}) => {
                if (type === this.currentSearchTarget.id) {
                    this.searchText = text
                }
            })
            const languageTranslate = {
                'zh_cn': 'zh_CN',
                'zh-cn': 'zh_CN',
                'zh': 'zh_CN'
            }
            let language = Cookies.get('blueking_language') || 'zh_CN'
            language = languageTranslate.hasOwnProperty(language) ? languageTranslate[language] : language
            this.language = language
            this.languageLable = language === 'zh_CN' ? '中文' : 'EN'
        }
    }
</script>

<style media="screen" lang="scss" scoped>
    $borderColor: #e7e9ef; //边框色
    $primaryColor: #737987;
    .header-wrapper{
        height: 50px;
        background: #fbfbfb;
        border-bottom: 1px solid $borderColor;
        padding:7px 20px 0 11px;
        font-size: 14px;
        position: fixed;
        width: 100%;
        top: 0;
        z-index: 1000;
        min-width: 1024px;
        .header-right-contain{
            .select-trigger-box{
                top: 0;
                background: transparent;
                padding-bottom: 20px;
            }
            .dropdown-content-ip{
                width: 100px;
                .select-trigger{
                    cursor:pointer;
                    span{
                        color: #bec6de;
                        letter-spacing: 1px;
                        cursor:pointer;
                        padding-left: 9px;
                        padding-right: 30px;
                    }
                    i{
                        color:#bec6de;
                        vertical-align: middle;
                        position: absolute;
                        top: 10px;
                        right: 15px;
                    }
                }
            }
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
            .quick-search-contain{
                height:32px;
                line-height:32px;
                background:#f9f9f9;
                border-radius:2px;
                margin: 2px -25px 0 0;
                background: #f2f2f2;
                color: $primaryColor;
                transform-origin: right center;
                transition: transform .2s ease-in;
                .select-trigger {
                    >span{
                        color: $primaryColor;
                    }
                    >i{
                        color: $primaryColor;
                        
                    }
                }
                .search-content{
                    width:240px;
                    display:inline-block;
                    position:relative;
                    color: $primaryColor;
                    input{
                        width:100%;
                        color: $primaryColor;
                        height:21px;
                        padding: 0 40px 0 10px;
                        border:none;
                        border-left: 1px solid #bec6de;
                        background: #f2f2f2;
                        outline: none;
                        border-left-color:#bec6de;
                        &:focus{
                            border-color: #bec6de!important;
                        }
                    }
                }
            }
            .quick-search-icon{
                width: 24px;
                height: 32px;
                line-height: 32px;
                cursor: pointer;
                z-index: 2;
                position: relative;
                .icon-search{
                    font-size: 16px;
                }
            }
            .language{
                position: relative;
                margin: -7px 0 0 40px;
                height: 50px;
                line-height: 50px;
                cursor: pointer;
                &:hover{
                    color: #498fe0;
                    .icon{
                        color: #498fe0;
                    }
                }
                .icon{
                    font-size: 16px;
                    color: #c3cdd7;
                }
                .language-text{
                    white-space: nowrap;
                }
                .bk-icon.icon-angle-down{
                    font-size: 12px;
                    margin-left: 2px;
                }
                &:hover{
                    .language-box{
                        display: block;
                    }
                }
                img{
                    margin-bottom: 6px;
                }
                .language-box{
                    padding: 10px 0;
                    position: absolute;
                    display: none;
                    width: 100px;
                    left: 100%;
                    top: 45px;
                    margin-left: calc(-50px - 50%);
                    background: #fff;
                    z-index: 1;
                    color: $primaryColor;
                    text-align: center;
                    border-radius: 3px;
                    box-shadow: 0px 1px 5px 0px rgba(12, 34, 59, 0.1);
                    li{
                        cursor: pointer;
                        height: 40px;
                        line-height: 40px;
                        &:hover,
                        &.active{
                            background: #f1f7ff;
                            color: #498fe0;
                        }
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
    .quick-search-enter,
    .quick-search-leave-to{
        transform: scaleX(0);
    }
</style>
