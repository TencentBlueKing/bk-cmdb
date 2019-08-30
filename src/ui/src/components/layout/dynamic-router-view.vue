<template>
    <div class="clearfix">
        <dynamic-navigation class="main-navigation" @business-change="handleBusinessChange"></dynamic-navigation>
        <div class="main-layout">
            <div class="main-scroller">
                <router-view class="main-views" :key="refreshKey"></router-view>
            </div>
        </div>
    </div>
</template>

<script>
    import dynamicNavigation from './dynamic-navigation'
    export default {
        components: {
            dynamicNavigation
        },
        data () {
            return {
                refreshKey: Date.now()
            }
        },
        methods: {
            handleBusinessChange () {
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
