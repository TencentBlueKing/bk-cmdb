<template>
    <div class="clearfix">
        <dynamic-navigation class="main-navigation"
            @business-change="handleBusinessChange">
        </dynamic-navigation>
        <dynamic-breadcrumbs class="main-breadcrumbs"></dynamic-breadcrumbs>
        <div class="main-layout">
            <div class="main-scroller" v-bkloading="{ isLoading: globalLoading }">
                <router-view class="main-views" :name="view" :key="refreshKey" v-if="view"></router-view>
            </div>
        </div>
    </div>
</template>

<script>
    import { mapGetters } from 'vuex'
    import dynamicNavigation from './dynamic-navigation'
    import dynamicBreadcrumbs from './dynamic-breadcrumbs'
    export default {
        components: {
            dynamicNavigation,
            dynamicBreadcrumbs
        },
        data () {
            return {
                refreshKey: Date.now(),
                meta: this.$route.meta,
                ready: false
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
        methods: {
            handleBusinessChange () {
                if (this.ready) {
                    this.refreshKey = Date.now()
                } else {
                    this.ready = true
                }
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
