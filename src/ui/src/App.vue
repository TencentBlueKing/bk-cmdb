<template>
    <div id="app">
        <div class="browser-tips" v-if="showBrowserTips">
            <span class="tips-text">{{$t('Common["您的浏览器非Chrome，建议您使用最新版本的Chrome浏览，以保证最好的体验效果"]')}}</span>
            <i class="tips-icon bk-icon icon-close-circle-shape" @click="showBrowserTips = false"></i>
        </div>
        <the-header></the-header>
        <the-nav class="nav-layout"></the-nav>
        <main class="main-layout" v-bkloading="{isLoading: globalLoading}">
            <div ref="mainScroller" class="main-scroller" @scroll="execMainScrollListener($event)">
                <router-view class="views-layout"></router-view>
            </div>
        </main>
    </div>
</template>

<script>
    import theHeader from '@/components/layout/header'
    import theNav from '@/components/layout/nav'
    import { execMainScrollListener, execMainResizeListener } from '@/utils/main-scroller'
    import { addResizeListener, removeResizeListener } from '@/utils/resize-events'
    import { mapGetters } from 'vuex'
    export default {
        name: 'app',
        components: {
            theHeader,
            theNav
        },
        data () {
            const showBrowserTips = window.navigator.userAgent.toLowerCase().indexOf('chrome') === -1
            return {
                showBrowserTips,
                execMainScrollListener
            }
        },
        computed: {
            ...mapGetters(['globalLoading'])
        },
        mounted () {
            addResizeListener(this.$refs.mainScroller, execMainResizeListener)
        },
        beforeDestroy () {
            removeResizeListener(this.$refs.mainScroller, execMainResizeListener)
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
        z-index: 9999;
        .tips-text{
            margin: 0 20px 0 0 ;
        }
        .tips-icon{
            cursor: pointer;
        }
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
    }
    .views-layout{
        min-height: 100%;
        min-width: 1100px;
        padding: 20px;
    }
</style>