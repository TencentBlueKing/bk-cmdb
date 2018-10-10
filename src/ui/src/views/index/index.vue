<template>
    <div class="index-wrapper">
        <v-search class="index-search"></v-search>
        <v-recently ref="recently"></v-recently>
        <v-classify></v-classify>
        <p class="copyright" ref="copyright">
            Copyright © 2012-{{year}} Tencent BlueKing. All Rights Reserved. 腾讯蓝鲸 版权所有
        </p>
    </div>
</template>

<script>
    import vSearch from './children/search'
    import vRecently from './children/recently'
    import vClassify from './children/classify'
    import getScrollbarWidth from '@/utils/scrollbar-width.js'
    import { addMainScrollListener, removeMainScrollListener, addMainResizeListener, removeMainResizeListener } from '@/utils/main-scroller'
    export default {
        components: {
            vSearch,
            vRecently,
            vClassify
        },
        data () {
            const year = (new Date()).getFullYear()
            return {
                year,
                scrollHandler: null,
                resizeHandler: null
            }
        },
        beforeRouteLeave (to, from, next) {
            this.$refs.recently.updateRecently(to.path)
            next()
        },
        created () {
            const calcCopyrightPosition = ($scroller) => {
                const scrollerRect = $scroller.getBoundingClientRect()
                const scrollerHeight = scrollerRect.height
                const scrollerWidth = scrollerRect.width
                const scrollerTop = $scroller.scrollTop
                const copyrightHeight = this.$refs.copyright.getBoundingClientRect().height
                const scrollbarWidth = scrollerWidth === ($scroller.scrollWidth + getScrollbarWidth()) ? 0 : getScrollbarWidth()
                this.$refs.copyright.style.top = scrollerTop + scrollerHeight - copyrightHeight - scrollbarWidth + 'px'
            }
            this.scrollHandler = event => {
                calcCopyrightPosition(event.target)
            }
            this.resizeHandler = () => {
                calcCopyrightPosition(document.querySelector('.main-scroller'))
            }
            addMainScrollListener(this.scrollHandler)
            addMainResizeListener(this.resizeHandler)
        },
        mounted () {
            this.resizeHandler()
        },
        beforeDestroy () {
            removeMainScrollListener(this.scrollHandler)
            removeMainResizeListener(this.resizeHandler)
        }
    }
</script>

<style lang="scss" scoped>
    .index-wrapper{
        position: relative;
        background-color: #f5f6fa;
    }
    .index-search{
        width: 50%;
        margin: 0 auto;
        padding: 40px 0 50px;
    }
    .copyright{
        position: absolute;
        width: calc(100% - 50px);
        height: 43px;
        left: 25px;
        bottom: 0;
        line-height: 42px;
        font-size: 12px;
        text-align: center;
        color: rgba(116, 120, 131, 0.5);
        border-top: 1px solid rgba(116, 120, 131, 0.2);
        background-color: #f5f6fa;
    }
</style>