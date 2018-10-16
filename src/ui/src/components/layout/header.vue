<template>
    <header class="header-layout clearfix" 
        :class="{'nav-sticked': navStick}">
        <div class="breadcrumbs fl">
            <i class="breadcrumbs-back bk-icon icon-arrows-left" href="javascript:void(0)"
                v-if="showBack"
                @click="back"></i>
            <h2 class="breadcrumbs-current">{{$classify.i18n ? $t($classify.i18n) : $classify.name}}</h2>
            <i v-if="$classify.id === 'custom_query'" class="bk-icon icon-info-circle" v-tooltip="{content: $t('CustomQuery[\'保存后的查询可通过接口调用生效\']'), classes: 'custom-query-header-tooltip'}"></i>
        </div>
        <div class="header-options fr">
            <div class="user" v-click-outside="handleCloseUser">
                <p class="user-name" @click="isShowUserDropdown = !isShowUserDropdown">
                    {{userName}}({{userRole}})
                    <i class="user-name-angle bk-icon icon-angle-down"
                        :class="{dropped: isShowUserDropdown}">
                    </i>
                </p>
                <transition name="toggle-slide">
                    <ul class="user-dropdown" v-show="isShowUserDropdown">
                        <li class="user-dropdown-item" @click="logOut">
                            <i class="icon-cc-logout"></i>
                            {{$t("Common['注销']")}}
                        </li>
                    </ul>
                </transition>
            </div>
            <div class="helper" v-click-outside="handleCloseHelper">
                <i class="helper-icon bk-icon icon-question-circle" @click="isShowHelper = !isShowHelper"></i>
                <div class="helper-list" v-show="isShowHelper">
                    <a href="http://docs.bk.tencent.com/product_white_paper/cmdb/" target="_blank" class="helper-link"
                        @click="isShowHelper = false">
                        {{$t('Common["帮助文档"]')}}
                    </a>
                    <a href="https://github.com/Tencent/bk-cmdb" target="_blank" class="helper-link"
                        @click="isShowHelper = false">
                        {{$t('Common["开源社区"]')}}
                    </a>
                </div>
            </div>
        </div>
    </header>
</template>

<script>
    import { mapGetters } from 'vuex'
    export default {
        data () {
            return {
                isShowUserDropdown: false,
                isShowHelper: false
            }
        },
        computed: {
            ...mapGetters(['site', 'userName', 'admin', 'showBack', 'navStick']),
            userRole () {
                return this.admin ? this.$t('Common["管理员"]') : this.$t('Common["普通用户"]')
            }
        },
        methods: {
            // 回退路由
            back () {
                this.$store.commit('setHeaderStatus', {
                    back: false
                })
                this.$router.back()
            },
            // 退出登陆
            logOut () {
                this.$http.post(`${window.API_HOST}logout`, {
                    'http_scheme': window.location.protocol.replace(':', '')
                }).then(data => {
                    window.location.href = data.url
                })
            },
            handleCloseUser () {
                this.isShowUserDropdown = false
            },
            handleCloseHelper () {
                this.isShowHelper = false
            }
        }
    }
</script>

<style lang="scss" scoped>
    .header-layout{
        position: relative;
        height: 61px;
        padding: 0 0 0 60px;
        border-bottom: 1px solid $cmdbBorderColor;
        background-color: #fff;
        transition: padding .1s ease-in;
        z-index: 1000;
        &.nav-sticked{
            padding-left: 240px;
        }
    }
    .breadcrumbs{
        line-height: 60px;
        position: relative;
        margin: 0 0 0 25px;
        font-size: 0;
        .breadcrumbs-back{
            display: inline-block;
            vertical-align: middle;
            width: 24px;
            height: 24px;
            line-height: 24px;
            text-align: center;
            font-size: 16px;
            font-weight: bold;
            cursor: pointer;
            &:hover{
                color: #3c96ff;
            }
        }
        .breadcrumbs-current{
            margin: 0;
            padding: 0;
            display: inline-block;
            vertical-align: middle;
            font-size: 16px;
            font-weight: normal;
        }
        .icon-info-circle {
            margin-left: 5px;
            font-size: 16px;
        }
    }
    .header-options {
        text-align: right;
    }
    .user{
        display: inline-block;
        vertical-align: top;
        font-size: 0;
        line-height: 60px;
        position: relative;
        .user-name{
            padding: 0 20px;
            margin: 0;
            font-size: 14px;
            font-weight: bold;
            color: rgba(115,121,135,1);
            cursor: pointer;
            .user-name-angle{
                display: inline-block;
                font-size: 12px;
                margin: 0 2px;
                color: $cmdbTextColor;
                transition: transform .2s linear;
                &.dropped{
                    transform: rotate(-180deg);
                }
            }
        }
        .user-dropdown{
            position: absolute;
            width: 100px;
            top: 55px;
            right: 20px;
            padding: 10px 0;
            line-height: 45px;
            font-size: 14px;
            background-color: #fff;
            box-shadow: 0 1px 5px 0 rgba(12,34,59, .1);
            z-index: 1;
            .user-dropdown-item{
                padding: 0 0 0 12px;
                text-align: left;
                cursor: pointer;
                &:hover{
                    background-color: #f1f7ff;
                    color: #498fe0;
                }
            }
        }
    }
    .helper {
        position: relative;
        display: inline-block;
        width: 60px;
        text-align: center;
        vertical-align: top;
        line-height: 60px;
        border-left: 1px solid #ebf0f5;
        .helper-icon {
            font-size: 20px;
            cursor: pointer;
            &:hover {
                color: #0082ff;
            }
        }
        .helper-list {
            position: absolute;
            top: 55px;
            right: 1px;
            text-align: left;
            line-height: 40px;
            background-color: #fff;
            border-radius: 2px;
            box-shadow: 0 1px 5px 0 rgba(12,34,59, .1);
            .helper-link {
                display: block;
                padding: 0 20px;
                font-size: 14px;
                white-space: nowrap;
                &:hover {
                    background-color: #f1f7ff;
                    color: #498fe0;
                }
            }
        }
    }
</style>

<style lang="scss">
    .tooltip.custom-query-header-tooltip {
        .tooltip-inner {
            margin-top: 5px;
            background: white;
            color: $cmdbTextColor;
            padding: 30px;
            max-width: 300px;
            box-shadow: 0px 1px 6px 0px rgba(0, 0, 0, 0.3);
        }
        .tooltip-arrow {
            width: 0;
            height: 0;
            margin-top: 5px;
            border-style: solid;
            position: absolute;
            border-color: rgba(0, 0, 0, 0.3);
            z-index: 2;
            &:before {
                content: "";
                border-width: 0 5px 5px 5px;
                border-left-color: transparent !important;
                border-right-color: transparent !important;
                border-top-color: transparent !important;
                top: -5px;
                left: calc(50% - 5px);
                width: 0;
                height: 0;
                margin-top: 5px;
                border-style: solid;
                position: absolute;
                border-color: white;
                z-index: 1;
            }
        }
    }
</style>
