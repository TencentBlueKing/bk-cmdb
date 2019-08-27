<template>
    <div class="clearfix">
        <dynamic-navigation class="main-navigation"></dynamic-navigation>
        <router-view class="main-views"></router-view>
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
    .main-views {
        overflow: hidden;
    }
</style>
