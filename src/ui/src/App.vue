<!--
 * Tencent is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2017-2022 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
-->

<template>
  <div id="app" v-bkloading="{ isLoading: globalLoading }" :bk-language="$i18n.locale"
    :class="{
      'no-breadcrumb': hideBreadcrumbs,
      'main-full-screen': mainFullScreen
    }">
    <div class="browser-tips" v-if="showBrowserTips">
      <span class="tips-text">{{$t('您的浏览器非Chrome，建议您使用最新版本的Chrome浏览，以保证最好的体验效果')}}</span>
      <i class="tips-icon bk-icon icon-close-circle-shape" @click="showBrowserTips = false"></i>
    </div>
    <the-header></the-header>
    <router-view class="views-layout" :name="topView" ref="topView"></router-view>
    <the-permission-modal ref="permissionModal"></the-permission-modal>
    <the-login-modal ref="loginModal"
      v-if="loginUrl"
      :login-url="loginUrl"
      :success-url="loginSuccessUrl">
    </the-login-modal>
  </div>
</template>

<script>
  import theHeader from '@/components/layout/header'
  import thePermissionModal from '@/components/modal/permission'
  import theLoginModal from '@blueking/paas-login'
  import { addResizeListener, removeResizeListener } from '@/utils/resize-events'
  import { MENU_INDEX } from '@/dictionary/menu-symbol'
  import { mapGetters } from 'vuex'
  export default {
    name: 'app',
    components: {
      theHeader,
      thePermissionModal,
      theLoginModal
    },
    data() {
      const showBrowserTips = window.navigator.userAgent.toLowerCase().indexOf('chrome') === -1
      return {
        showBrowserTips,
        loginSuccessUrl: `${window.location.origin}/static/login_success.html`
      }
    },
    computed: {
      ...mapGetters(['globalLoading', 'mainFullScreen']),
      ...mapGetters('userCustom', ['usercustom', 'firstEntryKey', 'classifyNavigationKey']),
      isIndex() {
        return this.$route.name === MENU_INDEX
      },
      hideBreadcrumbs() {
        return !(this.$route.meta.layout || {}).breadcrumbs
      },
      topView() {
        const [topRoute] = this.$route.matched
        return (topRoute && topRoute.meta.view) || 'default'
      },
      loginUrl() {
        if (process.env.NODE_ENV === 'development') {
          return ''
        }
        const siteLoginUrl = this.$Site.login || ''
        const [loginBaseUrl] = siteLoginUrl.split('?')
        if (loginBaseUrl) {
          return `${loginBaseUrl}plain`
        }
        return ''
      }
    },
    mounted() {
      addResizeListener(this.$el, this.calculateAppHeight)
      window.permissionModal = this.$refs.permissionModal
      window.loginModal = this.$refs.loginModal

      // 在body标签添加语言标识属性，用于插入到body下的内容进行国际化处理
      document.body.setAttribute('lang', this.$i18n.locale)
    },
    beforeDestroy() {
      removeResizeListener(this.$el, this.calculateAppHeight)
    },
    methods: {
      calculateAppHeight() {
        this.$store.commit('setAppHeight', this.$el.offsetHeight)
      }
    }
  }
</script>
<style lang="scss" scoped>
    #app{
        height: 100%;
    }
    .browser-tips{
        position: fixed;
        top: 0;
        left: 0;
        width: 100%;
        height: 40px;
        line-height: 40px;
        text-align: center;
        color: #ff5656;
        background-color: #f8f6db;
        z-index: 99999;
        .tips-text{
            margin: 0 20px 0 0 ;
        }
        .tips-icon{
            cursor: pointer;
        }
    }
    .views-layout{
        height: calc(100% - 58px);
    }
    // 主内容区全屏
    .main-full-screen {
        /deep/ {
            .header-layout,
            .nav-layout,
            .breadcrumbs-layout {
                display: none;
            }
        }
        .views-layout {
            height: 100%;
        }
    }
    .no-breadcrumb {
        /deep/ {
            .main-layout {
               height: 100%;
            }
        }
    }
</style>
