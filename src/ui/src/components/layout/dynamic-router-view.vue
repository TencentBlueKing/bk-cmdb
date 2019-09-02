<template>
    <div class="clearfix">
        <dynamic-navigation class="main-navigation" @business-change="handleBusinessChange"></dynamic-navigation>
        <div class="main-layout">
            <div class="main-scroller">
                <dynamic-breadcumbs v-show="$route.meta.showBreadcumbs"></dynamic-breadcumbs>
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
        height: 100%;
    }
    .main-scroller {
        height: 100%;
        overflow: auto;
    }
    .main-views {
        height: calc(100% - 58px);
        padding: 20px;
        min-width: 1106px;
        background-color: #FAFBFD;
    }
</style>
