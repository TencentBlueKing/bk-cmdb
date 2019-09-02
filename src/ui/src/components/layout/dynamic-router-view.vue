<template>
    <div class="clearfix">
        <dynamic-navigation class="main-navigation" @business-change="handleBusinessChange"></dynamic-navigation>
        <dynamic-breadcumbs class="main-breadcumbs"></dynamic-breadcumbs>
        <div class="main-layout">
            <div class="main-scroller">
                <router-view class="main-views" :key="refreshKey"></router-view>
            </div>
        </div>
    </div>
</template>

<script>
    import dynamicNavigation from './dynamic-navigation'
    import dynamicBreadcumbs from './dynamic-breadcumbs'
    export default {
        components: {
            dynamicNavigation,
            dynamicBreadcumbs
        },
        data () {
            return {
                refreshKey: Date.now()
            }
        },
        methods: {
            handleBusinessChange () {
                this.$nextTick(() => {
                    this.refreshKey = Date.now()
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
        height: calc(100% - 58px);
    }
    .main-breadcumbs {
        overflow: hidden;
    }
    .main-scroller {
        height: 100%;
        overflow: auto;
    }
    .main-views {
        height: calc(100% - 58px);
        padding: 20px;
        min-width: 1106px;
    }
</style>
