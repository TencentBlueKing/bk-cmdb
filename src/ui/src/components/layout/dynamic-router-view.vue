<template>
    <div class="clearfix">
        <dynamic-navigation class="main-navigation"></dynamic-navigation>
        <div class="main-layout">
            <div class="main-scroller">
                <router-view class="main-views"></router-view>
            </div>
        </div>
    </div>
</template>

<script>
    import dynamicNavigation from './dynamic-navigation'
    export default {
        components: {
            dynamicNavigation
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
        height: 100%;
        padding: 0 20px;
        min-width: 1106px;
    }
</style>
