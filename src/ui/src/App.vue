<template>
    <div id="app" :bk-language="$i18n.locale">
        <div class="browser-tips" v-if="showBrowserTips">
            <span class="tips-text">{{$t('您的浏览器非Chrome，建议您使用最新版本的Chrome浏览，以保证最好的体验效果')}}</span>
            <i class="tips-icon bk-icon icon-close-circle-shape" @click="showBrowserTips = false"></i>
        </div>
        <the-header></the-header>
        <the-nav class="nav-layout"></the-nav>
        <main class="main-layout" v-bkloading="{ isLoading: globalLoading }">
            <div class="admin-tips" v-if="false">
                <span class="tips-text">{{$t('欢迎来到蓝鲸配置平台全局管理中心！您所做的操作将影响公共部分内容，请谨慎操作')}}</span>
                <i class="bk-icon icon-close"></i>
            </div>
            <div ref="mainScroller" class="main-scroller" @scroll="execMainScrollListener($event)">
                <router-view class="views-layout"></router-view>
            </div>
        </main>
        <the-permission-modal ref="permissionModal"></the-permission-modal>
    </div>
</template>

<script>
    import theHeader from '@/components/layout/header'
    import theNav from '@/components/layout/nav'
    import thePermissionModal from '@/components/modal/permission'
    import { execMainScrollListener, execMainResizeListener } from '@/utils/main-scroller'
    import { addResizeListener, removeResizeListener } from '@/utils/resize-events'
    import { mapGetters } from 'vuex'
    export default {
        name: 'app',
        components: {
            theHeader,
            theNav,
            thePermissionModal
        },
        data () {
            const showBrowserTips = window.navigator.userAgent.toLowerCase().indexOf('chrome') === -1
            return {
                showBrowserTips,
                execMainScrollListener
            }
        },
        computed: {
            ...mapGetters(['globalLoading']),
            ...mapGetters('userCustom', ['usercustom', 'firstEntryKey', 'classifyNavigationKey'])
        },
        mounted () {
            this.$store.commit('setFeatureTipsParams')
            addResizeListener(this.$refs.mainScroller, execMainResizeListener)
            addResizeListener(this.$el, this.calculateAppHeight)
            window.permissionModal = this.$refs.permissionModal
        },
        beforeDestroy () {
            removeResizeListener(this.$refs.mainScroller, execMainResizeListener)
            removeResizeListener(this.$el, this.calculateAppHeight)
        },
        methods: {
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
    .admin-tips {
        position: absolute;
        top: 0;
        left: 0;
        width: 100%;
        height: 40px;
        line-height: 38px;
        text-align: center;
        color: #3a84ff;
        border-top: 1px solid #3a84ff;
        border-bottom: 1px solid #3a84ff;
        background: rgba(58, 132, 255, .13);
        z-index: 9999;
    }
    .nav-layout{
        position: relative;
        float: left;
        height: 100%;
        margin: -61px 0 0 0;
        z-index: 1001;
    }
    .main-layout{
        height: calc(100% - 61px);
        overflow: hidden;
        position: relative;
    }
    .main-scroller{
        height: 100%;
        overflow: auto;
        @include scrollbar;
    }
    .views-layout{
        min-height: 100%;
        min-width: 1100px;
        padding: 20px;
    }
</style>
