<template>
    <header class="header-layout">
        <div class="logo">
            <router-link class="logo-link" exact to="/index">
                {{$t('蓝鲸配置平台')}}
            </router-link>
        </div>
        <nav class="header-nav">
            <router-link class="header-link" exact
                v-for="nav in headerNav"
                :to="nav.id"
                :key="nav.id">
                {{$t(nav.i18n)}}
            </router-link>
        </nav>
    </header>
</template>

<script>
    import { HEADER_NAV } from '@/dictionary/menu'
    export default {
        data () {
            return {
                headerNav: HEADER_NAV
            }
        }
    }
</script>

<style lang="scss" scoped>
    .header-layout {
        position: relative;
        display: flex;
        height: 58px;
        background-color: #182132;
        z-index: 1002;
    }
    .logo {
        flex: 292px 0 0;
        font-size: 0;
        .logo-link {
            display: inline-block;
            vertical-align: middle;
            height: 58px;
            line-height: 58px;
            margin-left: 23px;
            padding-left: 38px;
            color: #fff;
            font-size: 16px;
            background: url("../../assets/images/logo.svg") no-repeat 0 center;
        }
    }
    .header-nav {
        flex: 1;
        font-size: 0;
        .header-link {
            display: inline-block;
            vertical-align: middle;
            height: 58px;
            line-height: 58px;
            padding: 0 25px;
            color: #979BA5;
            font-size: 14px;
            &:hover {
                background-color: rgba(49, 64, 94, .5);
                color: #fff;
            }
            &.router-link-active {
                background-color: rgba(49, 64, 94, 1);
                color: #fff;
            }
        }
    }
</style>
