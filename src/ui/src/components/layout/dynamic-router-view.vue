<template>
    <div class="clearfix">
        <dynamic-navigation class="main-navigation"></dynamic-navigation>
        <dynamic-breadcrumbs class="main-breadcrumbs" v-if="$route.meta.layout.breadcrumbs"></dynamic-breadcrumbs>
        <div class="main-layout">
            <div class="main-scroller" v-bkloading="{ isLoading: globalLoading }" ref="scroller">
                <router-view class="main-views" :name="view"></router-view>
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
    export default {
        components: {
            dynamicNavigation,
            dynamicBreadcrumbs
        },
        data () {
            return {
                refreshKey: Date.now(),
                meta: this.$route.meta,
                scrollerObserver: null
            }
        },
        computed: {
            ...mapGetters(['globalLoading']),
            view () {
                return this.meta.view
            }
        },
        watch: {
            $route (val) {
                this.meta = this.$route.meta
            }
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
            },
            scrollerObserverHandler () {
                this.$nextTick(() => {
                    const scroller = this.$refs.scroller
                    if (scroller) {
                        const gutter = scroller.offsetHeight - scroller.clientHeight
                        this.$store.commit('setAppHeight', this.$root.$el.offsetHeight - gutter)
                        this.$store.commit('setScrollerState', {
                            scrollbar: scroller.scrollHeight > scroller.offsetHeight
                        })
                    }
                })
            }
        }
    }
</script>

<style lang="scss" scoped>
    .main-navigation {
        float: left;
    }
    .main-layout {
        position: relative;
        overflow: hidden;
        height: 100%;
        margin-top: -58px;
        z-index: 99;
    }
    .main-breadcrumbs {
        overflow: hidden;
        position: relative;
        background-color: #fafbfd;
        margin-right: 17px;
        z-index: 100;
    }
    .main-scroller {
        height: 100%;
        overflow: auto;
    }
    .main-views {
        height: calc(100% - 58px);
        margin-top: 58px;
        min-width: 1106px;
    }
</style>
