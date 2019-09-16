<template>
    <div class="clearfix">
        <dynamic-navigation class="main-navigation"
            @business-change="handleBusinessChange"
            @business-empty="handleBusinessEmpty">
        </dynamic-navigation>
        <dynamic-breadcrumbs class="main-breadcrumbs"></dynamic-breadcrumbs>
        <div class="main-layout">
            <div class="main-scroller" v-bkloading="{ isLoading: globalLoading }">
                <router-view class="main-views" :key="refreshKey" v-if="shouldRenderSubView"></router-view>
                <router-view class="main-views" name="requireBusiness" v-if="showRequireBusiness"></router-view>
            </div>
        </div>
    </div>
</template>

<script>
    import { mapGetters } from 'vuex'
    import dynamicNavigation from './dynamic-navigation'
    import dynamicBreadcrumbs from './dynamic-breadcrumbs'
    import { MENU_BUSINESS } from '@/dictionary/menu-symbol'
    export default {
        components: {
            dynamicNavigation,
            dynamicBreadcrumbs
        },
        data () {
            return {
                refreshKey: Date.now(),
                businessSelected: false,
                showRequireBusiness: false
            }
        },
        computed: {
            ...mapGetters(['globalLoading']),
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
                this.showRequireBusiness = false
                this.refreshKey = Date.now()
            },
            handleBusinessEmpty () {
                this.businessSelected = false
                this.showRequireBusiness = true
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
