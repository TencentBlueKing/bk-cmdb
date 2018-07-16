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
    <div id="app" class="clearfix">
        <div class="error-message-content is-chrome" v-show="isChromeShow">
            <span>{{$t('Common["您的浏览器非Chrome，建议您使用最新版本的Chrome浏览，以保证最好的体验效果"]')}}</span><i class="bk-icon icon-close-circle-shape" @click="closeInfo"></i>
        </div>
        <v-header></v-header>
        <v-nav class="fl"></v-nav>
        <div class="main-container">
            <div class="main-wrapper">
                <div class="content-wrapper" v-bkloading="{isLoading: globalLoading}">
                    <router-view/>
                </div>
            </div>
        </div>
    </div>
</template>
<script>
    import vHeader from '@/components/header/header-v3'
    import vNavigation from '@/components/nav/nav'
    import vNav from '@/components/nav/nav-v3'
    import { mapGetters } from 'vuex'
    export default {
        name: 'app',
        components: {
            vHeader,
            vNavigation,
            vNav
        },
        data () {
            return {
                isChromeShow: true
            }
        },
        computed: {
            ...mapGetters(['globalLoading'])
        },
        methods: {
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
    .main-container{
        overflow: hidden;
        height: 100%;
        position: relative;
        .main-wrapper{
            height: 100%;
            padding-top: 61px;
            overflow: auto;
            position: relative;
            .content-wrapper{
                height: 100%;
                min-width: 1060px;
            }
        }
    }
</style>
<style lang="scss">
    @import './common/scss/common.scss';
    @import './common/icon/cc-icon/style.css';
    @import './common/icon/bk-icon-2.0/iconfont.css';
    @import './magicbox/bk-magicbox-ui/bk.scss';
</style>
