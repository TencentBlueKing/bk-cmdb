<template>
    <div id="app" :bk-language="$i18n.locale"
        :class="{
            'no-breadcrumb': hideBreadcrumbs,
            'main-full-screen': mainFullScreen
        }">
        <div class="browser-tips" v-if="showBrowserTips">
            <span class="tips-text">{{$t('您的浏览器非Chrome，建议您使用最新版本的Chrome浏览，以保证最好的体验效果')}}</span>
            <i class="tips-icon bk-icon icon-close-circle-shape" @click="showBrowserTips = false"></i>
        </div>
        <the-header></the-header>
        <router-view class="views-layout" v-bkloading="{ isLoading: isIndex && globalLoading }"></router-view>
        <the-permission-modal ref="permissionModal"></the-permission-modal>
        <the-login-modal ref="loginModal"
            v-if="loginUrl"
            :login-url="loginUrl"
            :success-url="loginSuccessUrl">
        </the-login-modal>
        <cmdb-business-selector v-if="businessSelectorVisible" hidden
            @on-select="resolveBusinessSelectorPromise"
            @business-empty="resolveBusinessSelectorPromise">
        </cmdb-business-selector>
    </div>
</template>

<script>
    import theHeader from '@/components/layout/header'
    import thePermissionModal from '@/components/modal/permission'
    import theLoginModal from '@blueking/paas-login'
    // import { execMainScrollListener, execMainResizeListener } from '@/utils/main-scroller'
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
        data () {
            const showBrowserTips = window.navigator.userAgent.toLowerCase().indexOf('chrome') === -1
            const siteLoginUrl = window.Site.login
            const loginStrIndex = siteLoginUrl.indexOf('login')
            let loginModalUrl
            if (loginStrIndex > -1) {
                loginModalUrl = siteLoginUrl.substring(0, loginStrIndex) + 'login/plain'
            }
            return {
                showBrowserTips,
                loginUrl: loginModalUrl,
                loginSuccessUrl: window.Site.url + 'static/login_success.html'
                // execMainScrollListener
            }
        },
        computed: {
            ...mapGetters(['globalLoading', 'businessSelectorVisible', 'mainFullScreen']),
            ...mapGetters('userCustom', ['usercustom', 'firstEntryKey', 'classifyNavigationKey']),
            isIndex () {
                return this.$route.name === MENU_INDEX
            },
            hideBreadcrumbs () {
                return !(this.$route.meta.layout || {}).breadcrumbs
            }
        },
        mounted () {
            this.$store.commit('setFeatureTipsParams')
            // addResizeListener(this.$refs.mainScroller, execMainResizeListener)
            addResizeListener(this.$el, this.calculateAppHeight)
            window.permissionModal = this.$refs.permissionModal
            window.loginModal = this.$refs.loginModal
        },
        beforeDestroy () {
            // removeResizeListener(this.$refs.mainScroller, execMainResizeListener)
            removeResizeListener(this.$el, this.calculateAppHeight)
        },
        methods: {
            resolveBusinessSelectorPromise (val) {
                this.$store.commit('resolveBusinessSelectorPromise', !!val)
            },
            calculateAppHeight () {
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
            .nav-layout {
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
                margin-top: 0
            }
            .main-views {
                height: 100%;
                margin-top: 0;
            }
        }
    }
</style>
