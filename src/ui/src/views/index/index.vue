<template>
    <div class="index-layout" :class="{
        'is-sticky': sticky
    }">
        <div class="sticky-layout" ref="stickyLayout">
            <the-search v-if="site.fullTextSearch === 'off'"></the-search>
            <search-input v-else></search-input>
            <the-recently></the-recently>
        </div>
        <div ref="stickyProxy" v-show="sticky"></div>
        <the-classify></the-classify>
        <the-map></the-map>
        <cmdb-main-inject class="copyright" ref="copyright">
            Copyright © 2012-{{year}} Tencent BlueKing. All Rights Reserved. 腾讯蓝鲸 版权所有. {{site.buildVersion}}
        </cmdb-main-inject>
    </div>
</template>

<script>
    import theSearch from './children/search'
    import searchInput from './children/search-input'
    import theRecently from './children/recently'
    import theClassify from './children/classify'
    import theMap from './children/map'
    import cmdbMainInject from '@/components/layout/main-inject'
    import { mapGetters } from 'vuex'
    import {
        addMainResizeListener,
        removeMainResizeListener,
        addMainScrollListener,
        removeMainScrollListener
    } from '@/utils/main-scroller'
    export default {
        name: 'cmdb-index',
        components: {
            theSearch,
            searchInput,
            theRecently,
            theClassify,
            theMap,
            cmdbMainInject
        },
        data () {
            return {
                year: (new Date()).getFullYear(),
                sticky: false,
                resizeHandler: null,
                scrollHandler: null
            }
        },
        computed: {
            ...mapGetters(['site'])
        },
        mounted () {
            this.initResizeListener()
            this.initScrollListener()
        },
        beforeDestroy () {
            removeMainResizeListener(this.resizeHandler)
            removeMainScrollListener(this.scrollHandler)
        },
        methods: {
            initResizeListener () {
                const $copyright = this.$refs.copyright.$el
                this.resizeHandler = event => {
                    const target = event.target
                    if (target.offsetWidth < 1100) {
                        $copyright.style.bottom = '8px'
                    } else {
                        $copyright.style.bottom = 0
                    }
                }
                addMainResizeListener(this.resizeHandler)
                this.resizeHandler({ target: document.querySelector('.main-scroller') })
            },
            initScrollListener () {
                this.scrollHandler = event => {
                    const target = event.target
                    this.sticky = target.scrollTop > 170
                    if (this.sticky) {
                        this.$refs.stickyLayout.style.top = target.scrollTop + 'px'
                    }
                }
                this.$refs.stickyProxy.style.height = this.$refs.stickyLayout.offsetHeight + 'px'
                addMainScrollListener(this.scrollHandler)
            }
        }
    }
</script>

<style lang="scss" scoped>
    .index-layout {
        overflow: auto;
        padding: 0 0 50px;
        background-color: #f5f6fa;
        position: relative;
        z-index: 1;
        &.is-sticky {
            .sticky-layout {
                position: absolute;
                top: 50px;
                left: 0;
                width: 100%;
                padding-top: 50px;
                box-shadow: 0 0 8px 1px rgba(0, 0, 0, 0.03);
                background-color: #f5f6fa;
            }
        }
    }
    .sticky-layout {
        padding: 220px 0 27px;
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
        z-index: 2;
    }
</style>
