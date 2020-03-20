<template>
    <header class="header-layout">
        <div class="logo">
            <router-link class="logo-link" to="/index">
                {{$t('蓝鲸配置平台')}}
            </router-link>
        </div>
        <nav class="header-nav">
            <router-link class="header-link"
                v-for="nav in menu"
                :to="{ name: nav.id }"
                :key="nav.id">
                {{$t(nav.i18n)}}
            </router-link>
        </nav>
        <section class="header-info">
            <bk-popover class="info-item"
                theme="light header-info-popover"
                trigger="click"
                animation="fade"
                placement="bottom-end"
                :arrow="false"
                :tippy-options="{
                    animateFill: false
                }">
                <span class="info-user">
                    <span class="user-name">{{userName}}</span>
                    <i class="user-icon bk-icon icon-angle-down"></i>
                </span>
                <template slot="content">
                    <a class="question-link" href="javascript:void(0)"
                        @click="handleLogout">
                        <i class="icon-cc-logout"></i>
                        {{$t('注销')}}
                    </a>
                </template>
            </bk-popover>
            <bk-popover class="info-item"
                theme="light header-info-popover"
                trigger="click"
                animation="fade"
                placement="bottom-end"
                :arrow="false"
                :tippy-options="{
                    animateFill: false
                }">
                <i class="question-icon icon-cc-default"></i>
                <template slot="content">
                    <a class="question-link" target="_blank" href="http://docs.bk.tencent.com/product_white_paper/cmdb/">{{$t('帮助文档')}}</a>
                    <a class="question-link" target="_blank" href="https://github.com/Tencent/bk-cmdb">{{$t('开源社区')}}</a>
                </template>
            </bk-popover>
        </section>
    </header>
</template>

<script>
    import menu from '@/dictionary/menu'
    import { mapGetters } from 'vuex'
    export default {
        data () {
            return {
                menu: menu
            }
        },
        computed: {
            ...mapGetters(['userName'])
        },
        methods: {
            handleLogout () {
                this.$http.post(`${window.API_HOST}logout`, {
                    'http_scheme': window.location.protocol.replace(':', '')
                }).then(data => {
                    window.location.href = data.url
                })
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
        flex: 3;
        font-size: 0;
        white-space: nowrap;
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
    .header-info {
        flex: 1;
        text-align: right;
        white-space: nowrap;
        @include middleBlockHack;
    }
    .info-item {
        @include inlineBlock;
        margin: 0 25px 0 0;
        text-align: left;
        font-size: 0;
        cursor: pointer;
        .tippy-active {
            .bk-icon {
                color: #fff;
            }
            .user-icon {
                transform: rotate(-180deg);
            }
        }
        .question-icon {
            font-size: 16px;
            color: #DCDEE5;
            &:hover {
                color: #fff;
            }
        }
        .info-user {
            font-size: 14px;
            font-weight: bold;
            color: #fff;
            .user-name {
                max-width: 150px;
                @include inlineBlock;
                @include ellipsis;
            }
            .user-icon {
                margin-left: -4px;
                transition: transform .2s linear;
                font-size: 20px;
                color: #fff;
            }
        }
    }
    .question-link {
        display: block;
        padding: 0 20px;
        line-height: 40px;
        font-size: 14px;
        white-space: nowrap;
        &:hover {
            background-color: #f1f7ff;
            color: #3a84ff;
        }
    }
</style>

<style>
    .tippy-tooltip.header-info-popover-theme {
        padding: 0 !important;
        overflow: hidden;
        border-radius: 2px !important;
    }
</style>
