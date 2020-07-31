<template>
    <div class="clearfix">
        <dynamic-navigation class="main-navigation" v-show="!isEntry"></dynamic-navigation>
        <dynamic-breadcrumbs class="main-breadcrumbs" ref="breadcrumbs" v-if="showBreadcrumbs"></dynamic-breadcrumbs>
        <div class="main-layout">
            <div class="main-scroller" ref="scroller">
                <router-view class="main-views" :name="view" ref="view"></router-view>
            </div>
        </div>
    </div>
</template>

<script>
    import { mapGetters } from 'vuex'
    import dynamicNavigation from './dynamic-navigation'
    import dynamicBreadcrumbs from './dynamic-breadcrumbs'
    import {
        addResizeListener,
        removeResizeListener
    } from '@/utils/resize-events'
    import { MENU_ENTRY, MENU_ADMIN } from '@/dictionary/menu-symbol'
    import throttle from 'lodash.throttle'
    export default {
        components: {
            dynamicNavigation,
            dynamicBreadcrumbs
        },
        data () {
            return {
                refreshKey: Date.now(),
                meta: this.$route.meta,
                scrollerObserver: null,
                scrollerObserverHandler: null
            }
        },
        computed: {
            ...mapGetters(['globalLoading']),
            view () {
                return this.meta.view
            },
            isEntry () {
                const [topRoute] = this.$route.matched
                return topRoute && [MENU_ENTRY, MENU_ADMIN].includes(topRoute.name)
            },
            showBreadcrumbs () {
                return this.$route.meta.layout && this.$route.meta.layout.breadcrumbs
            }
        },
        watch: {
            $route (val) {
                this.meta = this.$route.meta
            }
        },
        created () {
            this.scrollerObserverHandler = throttle(() => {
                const scroller = this.$refs.scroller
                if (scroller) {
                    const gutter = scroller.offsetHeight - scroller.clientHeight
                    this.$store.commit('setAppHeight', this.$root.$el.offsetHeight - gutter)
                    this.$store.commit('setScrollerState', {
                        scrollbar: scroller.scrollHeight > scroller.offsetHeight
                    })
                }
            }, 300, { leading: false, trailing: true })
        },
        mounted () {
            addResizeListener(this.$refs.scroller, this.scrollerObserverHandler)
            this.addScrollerObserver()
        },
        beforeDestory () {
            removeResizeListener(this.$refs.scroller, this.scrollerObserverHandler)
            this.scrollerObserver && this.scrollerObserver.disconnect()
        },
        methods: {
            addScrollerObserver () {
                this.scrollerObserver = new MutationObserver(this.scrollerObserverHandler)
                this.scrollerObserver.observe(this.$refs.scroller, {
                    attributes: true,
                    childList: true,
                    subtree: true
                })
            }
        }
    }
</script>

<style lang="scss" scoped>
    .main-navigation {
        float: left;
    }
    .main-breadcrumbs {
        overflow: hidden;
        position: relative;
        background-color: #fafbfd;
        margin-right: 17px;
        z-index: 100;
        ~ .main-layout {
            margin-top: -53px;
        }
    }
    .main-layout {
        position: relative;
        overflow: hidden;
        height: 100%;
        z-index: 99;
    }
    .main-scroller {
        height: 100%;
        overflow: auto;
    }
    .main-views {
        position: relative;
        height: calc(100% - 52px);
        border-top: 1px solid $borderColor;
        margin-top: 52px;
        min-width: 1089px;
    }
</style>
