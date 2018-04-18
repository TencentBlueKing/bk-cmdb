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
    <div id="app">
        <div class="error-message-content is-chrome" v-show="isChromeShow">
            <span>{{$t('Common["您的浏览器非Chrome，建议您使用最新版本的Chrome浏览，以保证最好的体验效果"]')}}</span><i class="bk-icon icon-close-circle-shape" @click="closeInfo"></i>
        </div>
        <i class="bk-icon icon-dedent" id="iconCloseNav" :class="{'close': isShow}" @click="closeNav"></i>
        <v-navigation :isClose="isShow"></v-navigation>
        <v-header @quickCheck="quickCheck"></v-header>
        <div id="content-wrapper" :class="{'content-wrapper':!isShow,'content-control':isShow}" v-bkloading="{isLoading: globalLoading}">
            <router-view/>
        </div>
    </div>
</template>
<script>
    import vHeader from '@/components/header/header'
    import vNavigation from '@/components/nav/nav'
    import { mapActions, mapGetters } from 'vuex'
    export default {
        name: 'app',
        components: {
            vHeader,
            vNavigation
        },
        data () {
            return {
                isChromeShow: true,
                isShow: false,
                searchVal: ''
            }
        },
        computed: {
            ...mapGetters(['globalLoading'])
        },
        methods: {
            ...mapActions(['getAllClassify']),
            /*
                导航伸缩时，右侧内容变化
            */
            closeNav () {
                this.isShow = !this.isShow
            },
            quickCheck (searchVal) {
                this.searchVal = searchVal
            },
            closeInfo () {
                this.isChromeShow = false
            }
        },
        created () {
            this.isChromeShow = window.navigator.userAgent.toLowerCase().indexOf('chrome') === -1
        }
    }
</script>
<style lang="scss" scoped>
    $primaryColor: #737987;
    #iconCloseNav.icon-dedent{
        position: absolute;
        left: 240px;
        top: 18px;
        font-size: 16px;
        color:$primaryColor;
        cursor: pointer;
        transition: left .5s;
        z-index: 1001;
        &.close{
            left: 80px;
            transform: rotate(180deg);
        }
    }
</style>
<style lang="scss">
    @import './common/scss/common.scss';
    @import './common/icon/cc-icon/style.css';
    @import './common/icon/bk-icon-2.0/iconfont.css';
    @import './magicbox/bk-magicbox-ui/bk.scss';
    .clearfix:after{
        content: '';
        font-size: 0;
    }
    .error-message-content{
        position: fixed;
        top: 0;
        width: 100%;
        height: 40px;
        line-height: 40px;
        text-align: center;
        background: #f8f6db;
        z-index: 99999;
        span{
            color: #ff5656;
            margin-right: 20px;
        }
        i{
            cursor: pointer;
            color: #ff5656;
        }
    }
</style>
