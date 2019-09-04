<template>
    <div class="clearfix">
        <dynamic-navigation class="main-navigation" @business-change="handleBusinessChange"></dynamic-navigation>
        <dynamic-breadcumbs class="main-breadcumbs"></dynamic-breadcumbs>
        <div class="main-layout">
            <div class="main-scroller">
                <router-view class="main-views" :key="refreshKey" v-if="shouldRenderSubView"></router-view>
            </div>
        </div>
    </div>
</template>

<script>
    import dynamicNavigation from './dynamic-navigation'
    import dynamicBreadcumbs from './dynamic-breadcumbs'
    import { MENU_BUSINESS } from '@/dictionary/menu-symbol'
    export default {
        components: {
            dynamicNavigation,
            dynamicBreadcumbs
        },
        data () {
            return {
                refreshKey: Date.now(),
                businessSelected: false
            }
        },
        computed: {
            shouldRenderSubView () {
                if (this.$route.matched && this.$route.matched[0].name === MENU_BUSINESS) {
                    return this.businessSelected
                }
                return true
            }
        },
        methods: {
            handleBusinessChange () {
                this.businessSelected = true
                this.refreshKey = Date.now()
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
    .main-breadcumbs {
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
        padding: 0 20px 10px;
        min-width: 1106px;
    }
</style>
