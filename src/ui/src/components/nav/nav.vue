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
    <div class=" nav-wrapper bk-sidebar nav-contain" id="bk-sidebar" :class="{'slide-close':isClose,'slide-open':!isClose}">
        <div class="slide-switch" @click="backToIndex">
            <i class="nav-title-img" :class="$t('Common.logo')"></i><i class="icon-cc-triangle-slider" ></i>
        </div>
        <div class="nav-list">
            <ul>
                <template v-for="(navType, navIndex) in navigationOrder">
                    <li v-if="navigation[navType] && navigation[navType]['authorized'] && navigation[navType]['children'].length"
                        :class="{'open': navigation[navType]['isOpen']}"
                        :key="navIndex">
                        <a href="javascript:void(0)" @click="toggleNav(navigation[navType])">
                            <span class="icon-box">
                                <i :class="navigation[navType]['icon']"></i>
                            </span>
                            <span class="nav-name text-hd"
                                :title="navigation[navType]['i18n'] ? $t(navigation[navType]['i18n']) : navigation[navType]['name']">
                                   {{navigation[navType]['i18n'] ? $t(navigation[navType]['i18n']) : navigation[navType]['name']}} 
                            </span>
                            <span class="angle-box">
                                <i class="bk-icon icon-angle-down angle"></i>
                            </span>
                        </a>
                        <div class="flex-subnavs" :style="{
                            'height': calcSubNavHeight(navigation[navType]),
                            'display': isClose ? 'none' : 'block'
                        }">
                            <router-link exact
                                v-for="(subNav, subNavIndex) in navigation[navType]['children']"
                                v-if="subNav.authorized && !subNav.isPaused"
                                :key="subNavIndex"
                                :to="subNav.path"
                                :title="subNav.i18n ? $t(subNav.i18n) : subNav.name">
                                {{subNav.i18n ? $t(subNav.i18n) : subNav.name}} 
                            </router-link>
                        </div>
                    </li>
                </template>
            </ul>
            <div class="copyright-contain">
                <div class="copyright-line"></div>
                <p class="copyright-text">Copyright © 2012-{{year}} Tencent  BlueKing. All Rights Reserved.<br>腾讯蓝鲸&nbsp;版权所有</p>
            </div>
        </div>
    </div>
</template>
<script>
    import {mapGetters, mapActions} from 'vuex'
    import STATIC_NAV from './staticNav.json'
    export default {
        props: {
            isClose: {
                type: Boolean,
                default: false
            }
        },
        data () {
            return {
                year: (new Date()).getFullYear(),
                navigation: {},
                navigationOrder: [], // 导航item显示顺序
                hideClassify: ['bk_host_manage', 'bk_biz_topo'], // 需要隐藏的导航
                loadedStatus: {
                    classify: false,
                    authority: false
                }
            }
        },
        computed: {
            ...mapGetters([
                'bkSupplierAccount',
                'isAdmin',
                'authority',
                'adminAuthority',
                'allClassify'
            ]),
            ready () {
                return this.loadedStatus.classify && this.loadedStatus.authority && this.authority && this.adminAuthority
            }
        },
        watch: {
            navigation (navigation) {
                this.checkAuthority()
                this.$store.commit('setNavigation', Object.assign({}, navigation))
            },
            isClose (isClose) {
                Object.keys(this.navigation).map(navType => {
                    this.navigation[navType]['isOpen'] = isClose
                })
                this.checkCurrentNav()
            },
            '$route.fullPath' (fullPath) {
                this.checkAuthority()
                this.checkCurrentNav()
            },
            ready (ready) {
                if (ready) {
                    this.setNavigation()
                    this.checkAuthority()
                    if (window.location.hash === '#/') {
                        this.checkCurrentNav()
                    }
                }
            }
        },
        methods: {
            ...mapActions(['getAuthority', 'getAllClassify']),
            setNavigation () {
                let navigation = STATIC_NAV
                let navigationOrder = ['globalbusi', 'backConfig']
                let {sys_config: sysconfig, model_config: modelconfig} = this.isAdmin ? this.adminAuthority : this.authority
                let {global_busi: globalbusi, back_config: backConfig} = sysconfig || {}
                if (modelconfig) {
                    // 遍历具有权限的模型分类
                    for (let modelClassifyId in modelconfig) {
                        // 从所有分类中提取模型分类的基本信息
                        for (let i = 0; i < this.allClassify.length; i++) {
                            let classify = this.allClassify[i]
                            if (classify['bk_classification_id'] === modelClassifyId) {
                                // 找到模型分类基本信息后作为一级菜单添加到导航
                                if (this.hideClassify.indexOf(modelClassifyId) === -1) {
                                    navigation[modelClassifyId] = {
                                        id: classify['bk_classification_id'],
                                        name: classify['bk_classification_name'],
                                        icon: classify['bk_classification_icon'],
                                        authorized: true,
                                        isOpen: false,
                                        children: []
                                    }
                                    // 遍历具有权限的模型分类下的模型
                                    for (let modelId in modelconfig[modelClassifyId]) {
                                        let modelAuthority = modelconfig[modelClassifyId][modelId]
                                        // 从所有分类的模型中找到具有权限的模型的基本信息
                                        for (let j = 0; j < classify['bk_objects'].length; j++) {
                                            let modelObj = classify['bk_objects'][j]
                                            // 找到模型基本信息后作为二级菜单
                                            if (modelObj['bk_obj_id'] === modelId) {
                                                navigation[modelClassifyId]['children'].push({
                                                    id: modelObj['bk_obj_id'],
                                                    name: modelObj['bk_obj_name'],
                                                    path: `/organization/${modelObj['bk_obj_id']}`,
                                                    authorized: true,
                                                    authority: modelAuthority || [],
                                                    isPaused: modelObj['bk_ispaused']
                                                })
                                                break
                                            }
                                        }
                                    }
                                    // 控制一级菜单显示顺序
                                    if (modelClassifyId === 'bk_organization') {
                                        navigationOrder.splice(1, 0, modelClassifyId)
                                    } else {
                                        navigationOrder.splice(navigationOrder.length - 1, 0, modelClassifyId)
                                    }
                                }
                                break
                            }
                        }
                    }
                }
                // 判断主机管理中的权限，仅做权限有的判断
                if (globalbusi && globalbusi.length) {
                    navigation.globalbusi.children.map(config => {
                        if (globalbusi.indexOf(config.id) !== -1) {
                            config.authorized = true
                        }
                    })
                }
                // 判断后台配置中的权限
                if (backConfig && backConfig.length) {
                    navigation.backConfig.authorized = true
                    navigation.backConfig.children.map(config => {
                        config.authorized = backConfig.indexOf(config.id) !== -1
                    })
                }
                this.navigationOrder = navigationOrder
                this.navigation = navigation
            },
            // 切换菜单的展开收起状态
            toggleNav (targetNav) {
                if (!this.isClose) {
                    targetNav['isOpen'] = !targetNav['isOpen']
                    Object.keys(this.navigation).map(navType => {
                        if (this.navigation[navType]['id'] !== targetNav['id']) {
                            this.navigation[navType]['isOpen'] = false
                        }
                    })
                }
            },
            // 初始化当前URL对应的菜单
            checkCurrentNav () {
                let currentPath = this.$route.path
                Object.keys(this.navigation).map(navType => {
                    this.navigation[navType]['isOpen'] = this.navigation[navType]['children'].some(subNav => subNav.path === currentPath)
                })
            },
            // 计算子菜单高度，用于做展开收起的动画
            calcSubNavHeight (nav) {
                if (nav['isOpen']) {
                    let authorizedCount = 0
                    nav['children'].map(subNav => {
                        if (subNav.authorized && !subNav.isPaused) {
                            authorizedCount++
                        }
                    })
                    return `${authorizedCount * 36}px`
                } else {
                    return 0
                }
            },
            checkAuthority () {
                if (this.authority && Object.keys(this.authority).length && this.allClassify.length) {
                    let currentPath = this.$route.path
                    // let currentPath = this.$route.fullPath
                    if (currentPath !== '/403' && currentPath !== '/404') {
                        let isAuthorized = false
                        for (let navType in this.navigation) {
                            this.navigation[navType]['children'].map(({authorized, path}) => {
                                if (path === currentPath && authorized) {
                                    isAuthorized = true
                                }
                            })
                            if (isAuthorized) break
                        }
                        if (!isAuthorized) this.$router.push('/403')
                    }
                }
            },
            // 导航菜单整体的展开与收起状态切换
            backToIndex () {
                this.$router.push('/')
            }
        },
        async mounted () {
            // 非管理员需要通过接口获取权限，管理员直接通过this.$store.state.common.adminAuthority获取
            await this.getAllClassify().then(() => { this.loadedStatus.classify = true })
            await this.getAuthority().then(() => { this.loadedStatus.authority = true })
        }
    }
</script>
<style media="screen" lang="scss" scoped>
    .bk-sidebar{
        background:#334162;
        line-height:72px;
    }
    .dn{
        display: none!important;
    }
    .nav-contain{
        position:fixed;
        top:0;
        left:0;
        bottom:0;
        z-index:1201;
        .slide-switch{
            transition: all .5s;
            padding:0 10px;
            border-bottom:none;
            line-height:62px;
            text-align:center;
            position:relative;
            height: 62px;
            .triangle-slider{
                position: absolute;
                top: 24px;
                right: 5px;
            }
            >.bk-icon{
                position:absolute;
                right:10px;
                top:28px;
                color:#424c6b;
                font-size:12px;
            }
        }
        .flex-subnavs {
            display: block;
            overflow: hidden;
            transition: height .5s cubic-bezier(.23, 1, .23, 1);
            a.on{
                background:#424c6b;
                color:#fff;
            }
            a{
                height: 36px;
                line-height: 36px;
                &:hover,
                &.router-link-exact-active{
                    color:#fff;
                    background:#283556;
                }
            }
        }
        .nav-list {
            height: calc(100% - 60px);
            position: relative;
            ul{
                height: calc(100% - 120px);
                overflow: auto;
                &::-webkit-scrollbar {
                    width: 6px;
                    height: 5px;
                    &-thumb {
                        border-radius: 20px;
                        background: #a5a5a5;
                        box-shadow: inset 0 0 6px hsla(0,0%,80%,.3);
                    }
                }
            }
            a{
                &:hover{
                    background: #283556;
                }
            }
            a.on{
               background: #424c6b;
            }
            .icon-box {
                line-height: 48px;
                .bk-icon,[class*="icon-cc"]{
                    font-size:16px;
                    display: inline-block;
                    vertical-align: baseline;
                }
            }
            .bk-icon{
                font-size:12px;
            }
            .copyright-contain{
                position: absolute;
                bottom: 33px;
                left: 0;
                margin: 0 20px;
                padding: 18px 0 0 0;
                text-align: center;
                white-space: normal;
                letter-spacing: 0;
                line-height: 16px;
                font-size: 12px;
                color: rgba(255,255,255,.34);
                .copyright-line{
                    border-top: 1px solid #333e5d;
                    border-bottom: 1px solid rgba(255,255,255,.15);
                }
                .copyright-text{
                    width: 188px;
                    padding: 18px 0 0 0;
                    margin: 0 0 0 -4px;
                }
            }
        }
        li{
            a{
                border-bottom: none;
                color:#c9d0e6;
                &.router-link-exact-active{
                    color: #ffffff;
                }
            }
            &:hover >a{
               border-left:none;
            }
            >a{
                height: 48px;
                line-height: 48px;
                border-left:none;
            }
        }
        &.slide-open {
            transition: all 0.5s;
            width:220px;
            text-overflow: ellipsis;
            white-space: nowrap;
            -o-text-overflow: ellipsis;
            overflow: hidden;
            li{
                >a .icon-box{
                    vertical-align: middle;
                    padding:0 4px 0 38px;
                }
                >a .angle-box{
                    position: absolute;
                    right:22px;
                    font-size:12px;
                    transition: all .5s;
                    transform: rotate(90deg);
                }
            }
            li.open{
                background-color: #2f3c5d;
                >a .angle-box{
                    transform: rotate(0deg);
                }
            }
            .flex-subnavs {
                background:#2f3c5d;
                padding-left: 0;
                a{
                    padding-left:80px;
                    padding-right: 10px;
                    overflow: hidden;
                    text-overflow: ellipsis;
                    white-space: nowrap;
                    -o-text-overflow: ellipsis;

                }
            }
            &:hover >a{
               border-left:none;
            }
            >a{
                border-left:none;
            }
            .icon-cc-triangle-slider{
                font-size: 18px;
                position: absolute;
                right: 5px;
                top: 23px;
            }
            .nav-title-img{
                display: inline-block;
                background: url(../../common/images/nav_title.png) no-repeat;
                background-size: 100%;
                width: 173px;
                height: 31px;
                &.zh{
                    background: url(../../common/images/nav-title-zh.png) no-repeat;
                    background-size: 100%;
                }
            }
        }
        &.slide-close {
            transition: all 0.5s;
            width:60px;
            .nav-list,.nav-list ul{
                overflow: visible;
                li > a {
                    background-color: transparent;
                    &:hover{
                        background-color: #2f3c5d;
                    }
                }
            }
            li{
                .nav-name{
                    color: #fff;
                }
                &:hover{
                    .nav-name{
                        position: absolute;
                        background:#2f3c5d;
                        width: 150px;
                    }
                }
            }
            li:hover .flex-subnavs{
                display: block !important;
                height: auto !important;
            }
            .flex-subnavs {
                width: 150px;
                padding-bottom: 10px;
                background: #2f3c5d;
                a {
                    background: #2f3c5d;
                    color: #fff;
                    padding-left:37px;
                    // border-left: 1px solid #5e6d96;
                    overflow: hidden;
                    text-overflow: ellipsis;
                    white-space: nowrap;
                    -o-text-overflow: ellipsis;
                    &:hover{
                        background: #283556;
                    }
                    &.router-link-exact-active{
                        background:#283556;
                        color: #fff;
                    }
                }
            }
            .nav-list {
                .angle-box{
                    width: 0;
                    opacity: 0;
                    visibility: hidden;
                }
                .copyright-contain{
                    display: none;
                }
            }
            .nav-title-img{
                background: url(../../common/images/nav-title-close.png) no-repeat center;
                background-size: 100%;
                width: 27px;
                height: 30px;
                display: inline-block;
                background-size: 100%;
            }
            .icon-cc-triangle-slider{
                font-size: 0;
            }

        }
        .open > a {
            background:#2f3c5d;
            border-left:none;
        }
        .nav-list .nav-name{
            display: inline-block;
            vertical-align: middle;
            height: 48px;
            line-height: 48px;
            border-bottom: none;
            font-size:14px;
            font-weight:bold;
            transition: unset !important;
            width: 120px;
            padding-left: 10px;
            position: initial;
        }
    }
    .text-hd{
        width: 100px;
        overflow: hidden;
        text-overflow: ellipsis;
        white-space: nowrap;
        -o-text-overflow: ellipsis;
        overflow: hidden;
    }
</style>
<style lang="scss" scoped>
    body[lang="en"]{
        .nav-list{
            .bk-icon,[class*="icon-cc"]{
                vertical-align: 1px;
            }
        }
    }
</style>
