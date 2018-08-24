<template>
    <header class="header-layout clearfix" 
        :class="{'nav-sticked': navStick}">
        <div class="breadcrumbs fl">
            <i class="breadcrumbs-back bk-icon icon-arrows-left" href="javascript:void(0)"
                v-if="showBack"
                @click="back"></i>
            <h2 class="breadcrumbs-current">{{$classify.i18n ? $t($classify.i18n) : $classify.name}}</h2>
        </div>
        <div class="user fr">
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
    </header>
</template>

<script>
    import { mapGetters } from 'vuex'
    export default {
        data () {
            return {
                isShowUserDropdown: false
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
                window.location.href = this.site.login
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
        &-back{
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
        &-current{
            margin: 0;
            padding: 0;
            display: inline-block;
            vertical-align: middle;
            font-size: 16px;
            font-weight: normal;
        }
    }
    .user{
        font-size: 0;
        line-height: 60px;
        position: relative;
        &-name{
            padding: 0 20px;
            margin: 0;
            font-size: 14px;
            font-weight: bold;
            color: rgba(115,121,135,1);
            cursor: pointer;
            &-angle{
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
        &-dropdown{
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
            &-item{
                padding: 0 0 0 12px;
                cursor: pointer;
                &:hover{
                    background-color: #f1f7ff;
                    color: #498fe0;
                }
            }
        }
    }
</style>